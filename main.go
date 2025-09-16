package main

import (
	"bytes"
	"crypto/md5"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"retroblog/parser"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/yuin/goldmark"
	ghtml "github.com/yuin/goldmark/renderer/html"
)

var (
	srcDir         string
	outDir         string
	addr           string
	attachmentsDir = "attachments" // attachments folder inside each note dir
	postTpl        *template.Template
	indexTpl       *template.Template
	mdConv         goldmark.Markdown
)

func init() {
	flag.StringVar(&srcDir, "src", "./apple_notes_export", "Notes source directory")
	flag.StringVar(&outDir, "out", "./public", "Output directory")
	flag.StringVar(&addr, "addr", ":8080", "HTTP listen address")
}

func main() {
	flag.Parse()
	loadTemplates()
	// Allow raw HTML in markdown so that we can inject HTML tags for attachments.
	mdConv = goldmark.New(
		goldmark.WithRendererOptions(
			ghtml.WithUnsafe(),
		),
	)

	if err := rebuildAll(); err != nil {
		log.Fatalf("initial build failed: %v", err)
	}
	log.Println("initial build done.")

	go watch()

	log.Printf("serving %s at http://%s …", outDir, addr)
	http.Handle("/", http.FileServer(http.Dir(outDir)))
	log.Fatal(http.ListenAndServe(addr, nil))
}

func loadTemplates() {
	var err error
	postTpl, err = template.ParseFiles("tpl/post.html")
	if err != nil {
		log.Fatalf("post template: %v", err)
	}
	indexTpl, err = template.ParseFiles("tpl/index.html")
	if err != nil {
		log.Fatalf("index template: %v", err)
	}
}

type post struct {
	Category string
	Title    string
	File     string
	Date     time.Time
}

type categoryGroup struct {
	Name  string
	Posts []post
}

func rebuildAll() error {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return err
	}

	// Copy base assets (tpl/base -> public)
	if err := copyBaseAssets(); err != nil {
		return err
	}

	// Copy any top-level assets from src root (css, images, favicon, etc.)
	if err := copyTopLevelAssets(); err != nil {
		return err
	}

	var allPosts []post

	// Only scan first-level directories as categories
	cats, err := os.ReadDir(srcDir)
	if err != nil {
		return err
	}

	for _, ce := range cats {
		if !ce.IsDir() {
			continue
		}
		catName := ce.Name()
		catPath := filepath.Join(srcDir, catName)

		// Create category directory in output
		if err := os.MkdirAll(filepath.Join(outDir, catName), 0755); err != nil {
			return err
		}

		// Scan each subdir under category for a note.md
		entries, err := os.ReadDir(catPath)
		if err != nil {
			log.Printf("Failed to read category dir %s: %v", catPath, err)
			continue
		}

		for _, pe := range entries {
			if !pe.IsDir() {
				continue
			}
			postDir := filepath.Join(catPath, pe.Name())
			noteFile := filepath.Join(postDir, "note.md")
			if fi, err := os.Stat(noteFile); err == nil && !fi.IsDir() {
				p, err := buildPage(noteFile, catName)
				if err != nil {
					log.Printf("Failed to build page from %s: %v", noteFile, err)
					continue
				}
				allPosts = append(allPosts, p)
			}
		}
	}

	// Build index after all pages
	if err := buildIndex(allPosts); err != nil {
		return err
	}

	// After rebuild, prune outDir to keep it in sync with srcDir + tpl/base
	if err := pruneOutToMatchSource(); err != nil {
		log.Printf("Warning: prune failed: %v", err)
	}

	return nil
}

func buildPage(noteMDPath, category string) (post, error) {
	postDir := filepath.Dir(noteMDPath)
	slug := filepath.Base(postDir)

	// 1) Load meta.txt (optional but recommended)
	metaPath := filepath.Join(postDir, "meta.txt")
	var meta NoteMeta
	var metaErr error
	if fi, err := os.Stat(metaPath); err == nil && !fi.IsDir() {
		meta, metaErr = ParseMetaFile(metaPath)
		if metaErr != nil {
			log.Printf("Warning: failed to parse meta.txt at %s: %v", metaPath, metaErr)
		}
	} else {
		log.Printf("Warning: meta.txt not found in %s, attachments won't be auto-transformed", postDir)
	}

	// 2) Read note.md as markdown source
	bodyBytes, err := os.ReadFile(noteMDPath)
	if err != nil {
		return post{}, err
	}
	body := string(bodyBytes)

	// 3) Transform [Attachment: ID] placeholders based on meta
	if meta.Attachments != nil && len(meta.Attachments) > 0 {
		body = transformAttachmentTagsByMeta(body, slug, meta)
	}

	// 4) Convert markdown to HTML
	var buf bytes.Buffer
	if err := mdConv.Convert([]byte(body), &buf); err != nil {
		return post{}, err
	}

	// 5) Prepare output path: outDir/<category>/<slug>.html
	dirOut := filepath.Join(outDir, category)
	if err := os.MkdirAll(dirOut, 0755); err != nil {
		return post{}, err
	}
	htmlName := slug + ".html"
	outFile := filepath.Join(dirOut, htmlName)

	f, err := os.Create(outFile)
	if err != nil {
		return post{}, err
	}
	defer f.Close()

	title := meta.Title
	if strings.TrimSpace(title) == "" {
		// fallback to first heading line from markdown, if any
		if t := extractFirstHeading(body); t != "" {
			title = t
		} else {
			title = slug
		}
	}

	if err := postTpl.Execute(f, map[string]any{
		"Title":    title,
		"Body":     template.HTML(buf.String()),
		"Category": category,
	}); err != nil {
		return post{}, err
	}

	// 6) Copy attachments folder if it exists
	srcAttach := filepath.Join(postDir, attachmentsDir)
	if fi, err := os.Stat(srcAttach); err == nil && fi.IsDir() {
		dstAttach := filepath.Join(outDir, category, slug, attachmentsDir)
		if err := os.MkdirAll(dstAttach, 0755); err != nil {
			log.Printf("Failed to create attachments dir %s: %v", dstAttach, err)
		} else {
			err := filepath.WalkDir(srcAttach, func(p string, d fs.DirEntry, err error) error {
				if err != nil || d.IsDir() {
					return err
				}
				relName := filepath.Base(p)
				dest := filepath.Join(dstAttach, relName)
				return copyFile(p, dest)
			})
			if err != nil {
				log.Printf("Error copying attachments: %v", err)
			}
		}
	}

	// 7) Use note.md modification time as date
	st, _ := os.Stat(noteMDPath)
	return post{
		Category: category,
		Title:    title,
		File:     filepath.ToSlash(filepath.Join(category, htmlName)),
		Date:     st.ModTime(),
	}, nil
}

func buildIndex(posts []post) error {
	// Group by category
	groups := map[string][]post{}
	for _, p := range posts {
		groups[p.Category] = append(groups[p.Category], p)
	}

	// Prepare sorted categories
	var cats []categoryGroup
	for cat, ps := range groups {
		// Sort posts in each category by date desc
		sort.Slice(ps, func(i, j int) bool { return ps[i].Date.After(ps[j].Date) })
		cats = append(cats, categoryGroup{Name: cat, Posts: ps})
	}

	// Sort categories by name
	sort.Slice(cats, func(i, j int) bool { return cats[i].Name < cats[j].Name })

	// Render index.html
	f, err := os.Create(filepath.Join(outDir, "index.html"))
	if err != nil {
		return err
	}
	defer f.Close()

	return indexTpl.Execute(f, map[string]any{
		"Categories": cats,
		"Now":        time.Now().Format("2006-01-02 15:04"),
	})
}

// copyBaseAssets copies everything inside tpl/base into the public root directory.
func copyBaseAssets() error {
	baseDir := filepath.Join("tpl", "base")
	fi, err := os.Stat(baseDir)
	if err != nil {
		// If base dir doesn't exist, nothing to copy.
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if !fi.IsDir() {
		return nil
	}

	return filepath.WalkDir(baseDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(baseDir, path)
		if err != nil {
			return err
		}
		dest := filepath.Join(outDir, rel)
		return copyFile(path, dest)
	})
}

func copyTopLevelAssets() error {
	return filepath.WalkDir(srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}

		// Only copy top-level files, not files in subdirectories
		if filepath.Dir(path) != srcDir {
			return nil
		}

		// Skip note files
		if strings.HasSuffix(strings.ToLower(d.Name()), ".txt") || strings.HasSuffix(strings.ToLower(d.Name()), ".md") {
			return nil
		}

		dest := filepath.Join(outDir, d.Name())
		return copyFile(path, dest)
	})
}

func copyFile(src, dest string) error {
	// Ensure the destination directory exists
	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return err
	}

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}

func watch() {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		log.Printf("Watch init error: %v", err)
		return
	}
	defer w.Close()

	addTree := func(root string) {
		filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
			if err == nil && d.IsDir() {
				if err := w.Add(p); err != nil {
					log.Printf("Failed to watch directory %s: %v", p, err)
				}
			}
			return nil
		})
	}

	// Watch source notes and templates so that base assets changes also trigger rebuild
	addTree(srcDir)
	if _, err := os.Stat("tpl"); err == nil {
		addTree("tpl")
	}

	for {
		select {
		case ev, ok := <-w.Events:
			if !ok {
				return
			}

			if ev.Has(fsnotify.Create) {
				if fi, err := os.Stat(ev.Name); err == nil && fi.IsDir() {
					if err := w.Add(ev.Name); err != nil {
						log.Printf("Failed to watch new directory %s: %v", ev.Name, err)
					}
				}
			}

			if ev.Has(fsnotify.Write) || ev.Has(fsnotify.Create) ||
				ev.Has(fsnotify.Remove) || ev.Has(fsnotify.Rename) {
				log.Printf("Change detected: %s", ev)
				if err := rebuildAll(); err != nil {
					log.Printf("Rebuild error: %v", err)
				} else {
					log.Printf("Rebuild successful")
				}
			}

		case err, ok := <-w.Errors:
			if !ok {
				return
			}
			log.Printf("Watch error: %v", err)
		}
	}
}

// ================= Attachment transform based on meta.txt =================

var attachmentRefRe = regexp.MustCompile(`\[Attachment:\s*([A-Za-z0-9_-]+)\]`)

func transformAttachmentTagsByMeta(body, slug string, meta NoteMeta) string {
	return attachmentRefRe.ReplaceAllStringFunc(body, func(match string) string {
		m := attachmentRefRe.FindStringSubmatch(match)
		if len(m) != 2 {
			return match
		}
		id := m[1]
		att, ok := meta.Attachments[id]
		if !ok {
			// No mapping found, keep original text
			return match
		}
		// att.Path is relative to note dir (e.g., "attachments/111680A9.jpg")
		// HTML page lives at: out/<category>/<slug>.html
		// Copied file lives at: out/<category>/<slug>/attachments/...
		rel := filepath.ToSlash(filepath.Join(slug, att.Path))
		ext := strings.ToLower(filepath.Ext(att.Path))
		ext_v2 := strings.ToLower(filepath.Ext(att.OriginalFilename))

		log.Printf("Transforming attachment ID %s (%s) ext (%s) to HTML tag", id, att.Path, ext_v2)

		switch ext_v2 {
		case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp", ".svg":
			alt := att.OriginalFilename
			if alt == "" {
				alt = att.SavedAs
			}
			return `<img src="` + rel + `" alt="` + htmlEscapeAttr(alt) + `" loading="lazy" />`
		case ".pdf":
			// Prefer an embed with a fallback link
			return `<object data="` + rel + `" type="application/pdf" width="100%" height="800">
  <a href="` + rel + `" target="_blank" rel="noopener">Open PDF</a>
</object>`
		case ".mp4", ".webm", ".ogg", ".mov":
			return `<video controls style="max-width:100%;">
  <source src="` + rel + `" type="video/` + strings.TrimPrefix(ext_v2, ".") + `">
  <a href="` + rel + `">Download video</a>
</video>`
		case ".mp3", ".wav", ".flac", ".m4a", ".aac", ".oga":
			return `<audio controls>
  <source src="` + rel + `" type="audio/` + strings.TrimPrefix(ext_v2, ".") + `">
  <a href="` + rel + `">Download audio</a>
</audio>`

		case ".nes":
			uniqueID := fmt.Sprintf("nes_%x", md5.Sum([]byte(id)))[:12]
			name := att.OriginalFilename
			if name == "" {
				name = filepath.Base(att.Path)
			}
			log.Println("Embedding NES player for", name, "with ID", uniqueID)
			return parser.GenerateNESPlayerHTML(uniqueID, rel, htmlEscape(name), ext)

		case ".pbp", ".chd", ".7z":
			uniqueID := fmt.Sprintf("pbp_%x", md5.Sum([]byte(id)))[:12]
			name := att.OriginalFilename
			if name == "" {
				name = filepath.Base(att.Path)
			}
			log.Println("Embedding PSX player for", name, "with ID", uniqueID)
			return parser.GeneratePlayStationPlayerHTML(uniqueID, rel, htmlEscape(name), ext)
		case ".gba":
			uniqueID := fmt.Sprintf("gba_%x", md5.Sum([]byte(id)))[:12]
			name := att.OriginalFilename
			if name == "" {
				name = filepath.Base(att.Path)
			}
			log.Println("Embedding GBA player for", name, "with ID", uniqueID)
			return parser.GenerateGBAPlayerHTML(uniqueID, rel, htmlEscape(name), ext)
		case ".jar":
			uniqueID := fmt.Sprintf("jar_%x", md5.Sum([]byte(id)))[:12]
			name := att.OriginalFilename
			if name == "" {
				name = filepath.Base(att.Path)
			}
			log.Println("Embedding Jar player for", name, "with ID", uniqueID)
			html := parser.GenerateJARPlayerHTML(uniqueID, rel, htmlEscape(name), ext)
			return html
		case ".md", ".smd", ".gen", ".bin", ".sms", ".gg", ".32x", ".cue", ".iso":
			uniqueID := fmt.Sprintf("sega_%x", md5.Sum([]byte(id)))[:12]
			name := att.OriginalFilename
			if name == "" {
				name = filepath.Base(att.Path)
			}
			log.Println("EmbeddingJS Sega player for", name, "with ID", uniqueID)
			html := parser.GenerateSegaMDPlayerHTML(uniqueID, rel, htmlEscape(name), ext)
			return html
		default:
			name := att.OriginalFilename
			if name == "" {
				name = filepath.Base(att.Path)
			}
			return `<a href="` + rel + `" download>` + htmlEscape(name) + `</a>`
		}
	})
}

func extractFirstHeading(md string) string {
	for _, line := range strings.Split(md, "\n") {
		l := strings.TrimSpace(line)
		if strings.HasPrefix(l, "#") {
			// strip leading # and spaces
			return strings.TrimSpace(strings.TrimLeft(l, "#"))
		}
	}
	return ""
}

func htmlEscapeAttr(s string) string {
	// Minimal escaping for attributes
	s = strings.ReplaceAll(s, `"`, "&quot;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}

func htmlEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}

// ===================== meta.txt parser =====================

// NoteMeta reflects the meta.txt content with attachment details.
type NoteMeta struct {
	NoteID              int
	Title               string
	Folder              string
	AttachmentsDeclared *int
	Attachments         map[string]AttachmentMeta // keyed by ID
}

type AttachmentMeta struct {
	ID               string
	SavedAs          string
	OriginalFilename string
	SourceFile       string
	Type             string
	Position         int
	Path             string
}

var (
	reNoteID      = regexp.MustCompile(`(?m)^Note ID:\s*(\d+)\s*$`)
	reTitle       = regexp.MustCompile(`(?m)^Title:\s*(.+)\s*$`)
	reFolder      = regexp.MustCompile(`(?m)^Folder:\s*(.+)\s*$`)
	reAttachCount = regexp.MustCompile(`(?m)^Attachments:\s*(\d+)\s*file\(s\)\s*$`)

	reAttachmentStart = regexp.MustCompile(`(?i)^\s*\d+\.\s*Saved\s*as:\s*(.+?)\s*$`)
	reKeyVal          = regexp.MustCompile(`^\s*([A-Za-z ]+):\s*(.*?)\s*$`)
	reFirstInt        = regexp.MustCompile(`(\d+)`)
)

func ParseMetaFile(path string) (NoteMeta, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return NoteMeta{}, err
	}
	return ParseMeta(string(data))
}

func ParseMeta(text string) (NoteMeta, error) {
	m := NoteMeta{Attachments: map[string]AttachmentMeta{}}

	t := normalizeNewlines(text)

	// Required header fields
	if id, err := findInt(reNoteID, t); err == nil {
		m.NoteID = id
	} else {
		return m, errors.New("meta missing Note ID")
	}
	if s, err := findString(reTitle, t); err == nil {
		m.Title = s
	} else {
		return m, errors.New("meta missing Title")
	}
	if s, err := findString(reFolder, t); err == nil {
		m.Folder = s
	} else {
		return m, errors.New("meta missing Folder")
	}

	// Optional attachment count
	if ms := reAttachCount.FindStringSubmatch(t); len(ms) == 2 {
		if v, err := strconv.Atoi(ms[1]); err == nil {
			m.AttachmentsDeclared = &v
		}
	}

	// Extract attachment section
	attachSec := extractAttachmentSection(t)
	if strings.TrimSpace(attachSec) == "" {
		return m, nil
	}

	lines := strings.Split(attachSec, "\n")
	var cur *AttachmentMeta

	flush := func() {
		if cur != nil {
			if cur.ID != "" {
				m.Attachments[cur.ID] = *cur
			} else if cur.SavedAs != "" {
				// fallback key if ID missing
				key := strings.TrimSuffix(cur.SavedAs, filepath.Ext(cur.SavedAs))
				m.Attachments[key] = *cur
			}
			cur = nil
		}
	}

	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}
		if ms := reAttachmentStart.FindStringSubmatch(line); len(ms) == 2 {
			flush()
			cur = &AttachmentMeta{SavedAs: strings.TrimSpace(ms[1])}
			continue
		}
		if cur != nil {
			if kv := reKeyVal.FindStringSubmatch(line); len(kv) == 3 {
				key := strings.ToLower(strings.TrimSpace(kv[1]))
				val := strings.TrimSpace(kv[2])
				switch key {
				case "original filename":
					cur.OriginalFilename = val
				case "source file":
					cur.SourceFile = val
				case "type":
					cur.Type = val
				case "position":
					if mm := reFirstInt.FindStringSubmatch(val); len(mm) == 2 {
						if iv, err := strconv.Atoi(mm[1]); err == nil {
							cur.Position = iv
						}
					}
				case "path":
					cur.Path = val
				case "id":
					cur.ID = val
				}
			}
		}
	}
	flush()

	return m, nil
}

func normalizeNewlines(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	return s
}

func findString(re *regexp.Regexp, s string) (string, error) {
	m := re.FindStringSubmatch(s)
	if len(m) != 2 {
		return "", errors.New("not found")
	}
	return strings.TrimSpace(m[1]), nil
}

func findInt(re *regexp.Regexp, s string) (int, error) {
	st, err := findString(re, s)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(st)
}

func extractAttachmentSection(s string) string {
	// Find the "Attachment Details:" sentinel in a tolerant way
	i := strings.Index(s, "Attachment Details:")
	if i == -1 {
		return ""
	}
	return s[i:]
}

// ===================== pruning (sync public to src) =====================

// pruneOutToMatchSource scans srcDir and tpl/base to compute all expected outputs,
// then removes any files/empty directories under outDir that are not expected.
func pruneOutToMatchSource() error {
	outAbs, err := filepath.Abs(outDir)
	if err != nil {
		return err
	}
	exp, err := computeExpectedOutputs(outAbs)
	if err != nil {
		return err
	}
	return pruneDirRecursive(outAbs, outAbs, exp)
}

func computeExpectedOutputs(outAbs string) (map[string]struct{}, error) {
	expected := make(map[string]struct{})

	add := func(p string) { expected[filepath.Clean(p)] = struct{}{} }

	// index.html
	add(filepath.Join(outAbs, "index.html"))

	// tpl/base files copied to public root
	baseDir := filepath.Join("tpl", "base")
	if fi, err := os.Stat(baseDir); err == nil && fi.IsDir() {
		filepath.WalkDir(baseDir, func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return err
			}
			rel, err := filepath.Rel(baseDir, path)
			if err != nil {
				return err
			}
			add(filepath.Join(outAbs, rel))
			return nil
		})
	}

	// top-level assets from srcDir root (non .md/.txt)
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		low := strings.ToLower(name)
		if strings.HasSuffix(low, ".md") || strings.HasSuffix(low, ".txt") {
			continue
		}
		add(filepath.Join(outAbs, name))
	}

	// posts: out/<cat>/<slug>.html and attachments
	for _, ce := range entries {
		if !ce.IsDir() {
			continue
		}
		cat := ce.Name()
		catPath := filepath.Join(srcDir, cat)
		sub, _ := os.ReadDir(catPath)
		for _, pe := range sub {
			if !pe.IsDir() {
				continue
			}
			noteDir := filepath.Join(catPath, pe.Name())
			noteMD := filepath.Join(noteDir, "note.md")
			if stat, err := os.Stat(noteMD); err == nil && !stat.IsDir() {
				// expected page
				add(filepath.Join(outAbs, cat, pe.Name()+".html"))
				// expected attachments
				srcAttach := filepath.Join(noteDir, attachmentsDir)
				if fi, err := os.Stat(srcAttach); err == nil && fi.IsDir() {
					filepath.WalkDir(srcAttach, func(p string, d fs.DirEntry, err error) error {
						if err != nil || d.IsDir() {
							return err
						}
						add(filepath.Join(outAbs, cat, pe.Name(), attachmentsDir, filepath.Base(p)))
						return nil
					})
				}
			}
		}
	}

	return expected, nil
}

// pruneDirRecursive removes files not in expected and deletes empty directories.
// root must be the absolute outDir. dir is the current directory being pruned.
func pruneDirRecursive(root, dir string, expected map[string]struct{}) error {
	ents, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, e := range ents {
		p := filepath.Join(dir, e.Name())
		if e.IsDir() {
			// Recurse first
			if err := pruneDirRecursive(root, p, expected); err != nil {
				return err
			}
			// After recursion, remove directory if empty
			rem, _ := os.ReadDir(p)
			if len(rem) == 0 {
				if err := os.Remove(p); err == nil {
					log.Printf("Pruned empty directory: %s", relToRoot(root, p))
				}
			}
		} else {
			if _, ok := expected[filepath.Clean(p)]; !ok {
				if err := os.Remove(p); err == nil {
					log.Printf("Pruned file: %s", relToRoot(root, p))
				}
			}
		}
	}
	return nil
}

func relToRoot(root, p string) string {
	rel, err := filepath.Rel(root, p)
	if err != nil {
		return p
	}
	return rel
}
