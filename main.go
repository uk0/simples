package main

import (
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
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
	"slices"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/yuin/goldmark"
	ghtml "github.com/yuin/goldmark/renderer/html"
)

var (
	srcDir         string
	outDir         string
	addr           string
	baseURL        string
	attachmentsDir = "attachments"
	postTpl        *template.Template
	indexTpl       *template.Template
	mdConv         goldmark.Markdown

	protectedMu    sync.RWMutex
	protectedPosts = make(map[string]string)
	fullContent    = make(map[string]string)
	// 新增：记录哪些附件路径是受保护的
	protectedAttachments = make(map[string]bool)
	passRe               = regexp.MustCompile(`(?i)\[PASS:\s*([^\]]+)\]`)

	rebuilding   bool
	rebuildingMu sync.Mutex
)

func init() {
	flag.StringVar(&srcDir, "src", "./apple_notes_export", "Notes source directory")
	flag.StringVar(&outDir, "out", "./public", "Output directory")
	flag.StringVar(&addr, "addr", ":8080", "HTTP listen address")
	flag.StringVar(&baseURL, "url", "https://firsh.me", "Base URL for sitemap")
}

func main() {
	flag.Parse()
	baseURL = strings.TrimRight(baseURL, "/")
	loadTemplates()
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

	mux := http.NewServeMux()
	mux.HandleFunc("/api/unlock", handleUnlock)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// 检查是否是附件请求
		if strings.Contains(path, "/attachments/") {
			// 构建完整的附件路径
			fullPath := strings.TrimPrefix(path, "/")

			protectedMu.RLock()
			isProtected := protectedAttachments[fullPath]
			protectedMu.RUnlock()

			if isProtected {
				// 从路径中提取文章路径
				parts := strings.Split(fullPath, "/")
				if len(parts) >= 2 {
					postPath := parts[0] + "/" + parts[1]
					postPath = strings.TrimSuffix(postPath, ".html")

					cookie, err := r.Cookie("auth_" + hashString(postPath))
					if err != nil || !verifyAuthCookie(postPath, cookie.Value) {
						http.Error(w, "Unauthorized - Please unlock the post first", http.StatusUnauthorized)
						return
					}
				}
			}
		}

		http.FileServer(http.Dir(outDir)).ServeHTTP(w, r)
	})

	log.Printf("serving %s at http://%s …", outDir, addr)
	log.Printf("Site base URL: %s", baseURL)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func handleUnlock(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Path     string `json:"path"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	req.Path = strings.Trim(req.Path, "/")

	rebuildingMu.Lock()
	isRebuilding := rebuilding
	rebuildingMu.Unlock()

	if isRebuilding {
		time.Sleep(100 * time.Millisecond)
	}

	protectedMu.RLock()
	expectedHash, ok := protectedPosts[req.Path]
	content, hasContent := fullContent[req.Path]
	protectedMu.RUnlock()

	if !ok {
		http.Error(w, "Post not protected", http.StatusBadRequest)
		return
	}

	if hashPassword(req.Password) != expectedHash {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	if !hasContent {
		http.Error(w, "Content not found", http.StatusInternalServerError)
		return
	}

	authToken := generateAuthToken(req.Path, req.Password)
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_" + hashString(req.Path),
		Value:    authToken,
		Path:     "/",
		MaxAge:   86400,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   false,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"content": content,
		"status":  "success",
	})
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

func rebuildAll() error {
	rebuildingMu.Lock()
	rebuilding = true
	rebuildingMu.Unlock()

	defer func() {
		rebuildingMu.Lock()
		rebuilding = false
		rebuildingMu.Unlock()
	}()

	tempProtectedPosts := make(map[string]string)
	tempFullContent := make(map[string]string)
	tempProtectedAttachments := make(map[string]bool)

	if err := os.MkdirAll(outDir, 0755); err != nil {
		return err
	}

	if err := copyBaseAssets(); err != nil {
		return err
	}

	if err := copyTopLevelAssets(); err != nil {
		return err
	}

	var allPosts []post
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

		if err := os.MkdirAll(filepath.Join(outDir, catName), 0755); err != nil {
			return err
		}

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
				p, protectedInfo, err := buildPageWithProtection(noteFile, catName)
				if err != nil {
					log.Printf("Failed to build page from %s: %v", noteFile, err)
					continue
				}

				if protectedInfo != nil {
					tempProtectedPosts[protectedInfo.path] = protectedInfo.passwordHash
					tempFullContent[protectedInfo.path] = protectedInfo.fullContent

					// 记录受保护的附件路径
					for _, attachPath := range protectedInfo.protectedFilePaths {
						tempProtectedAttachments[attachPath] = true
					}

					p.IsProtected = true
				}

				allPosts = append(allPosts, p)
			}
		}
	}

	if err := buildIndex(allPosts); err != nil {
		return err
	}

	if err := buildSitemap(allPosts); err != nil {
		log.Printf("Warning: failed to build sitemap: %v", err)
	}

	if err := buildRSSFeed(allPosts); err != nil {
		log.Printf("Warning: failed to build RSS feed: %v", err)
	}

	if err := buildRobotsTxt(); err != nil {
		log.Printf("Warning: failed to build robots.txt: %v", err)
	}

	protectedMu.Lock()
	protectedPosts = tempProtectedPosts
	fullContent = tempFullContent
	protectedAttachments = tempProtectedAttachments
	protectedMu.Unlock()

	if err := pruneOutToMatchSource(); err != nil {
		log.Printf("Warning: prune failed: %v", err)
	}

	return nil
}

func buildPageWithProtection(noteMDPath, category string) (post, *protectionInfo, error) {
	postDir := filepath.Dir(noteMDPath)
	slug := filepath.Base(postDir)

	// Load meta.txt
	metaPath := filepath.Join(postDir, "meta.txt")
	var meta NoteMeta
	if fi, err := os.Stat(metaPath); err == nil && !fi.IsDir() {
		meta, _ = ParseMetaFile(metaPath)
	}

	bodyBytes, err := os.ReadFile(noteMDPath)
	if err != nil {
		return post{}, nil, err
	}
	body := string(bodyBytes)
	postPath := category + "/" + slug

	var protInfo *protectionInfo
	var isProtected bool
	var publicBody string
	var protectedBody string
	var protectedAttachmentIDs []string

	// 查找 [PASS: xxx] 标记
	if matches := passRe.FindStringSubmatch(body); len(matches) == 2 {
		isProtected = true
		password := strings.TrimSpace(matches[1])

		loc := passRe.FindStringIndex(body)
		if loc != nil {
			// 分割内容：公开部分和受保护部分
			publicBody = body[:loc[0]]
			protectedBody = body[loc[1]:] // 从 [PASS: xxx] 之后开始

			// 找出受保护部分中的附件ID
			protectedAttachmentIDs = findAttachmentIDs(protectedBody)
			fmt.Println("Protected attachment IDs:", protectedAttachmentIDs)
			// 构建受保护的附件路径列表
			var protectedPaths []string
			for _, id := range protectedAttachmentIDs {
				if att, ok := meta.Attachments[id]; ok {
					// 构建完整的输出路径
					attachPath := filepath.Join(category, slug, attachmentsDir, filepath.Base(att.Path))
					protectedPaths = append(protectedPaths, attachPath)
				}
			}

			protInfo = &protectionInfo{
				path:               postPath,
				passwordHash:       hashPassword(password),
				protectedFilePaths: protectedPaths,
			}

			log.Printf("Post %s has %d protected attachments", postPath, len(protectedPaths))
		}
	} else {
		// 没有密码保护，全部公开
		publicBody = body
	}

	// 处理公开部分
	description := ExtractDescription(publicBody)

	// 转换公开部分的附件
	if meta.Attachments != nil && len(meta.Attachments) > 0 {
		publicBody = transformAttachmentTagsByMeta(publicBody, slug, meta, false)
	}

	// 转换公开部分为HTML
	var publicBuf bytes.Buffer
	if err := mdConv.Convert([]byte(publicBody), &publicBuf); err != nil {
		return post{}, nil, err
	}

	// 如果有受保护内容，处理完整内容
	if isProtected && protInfo != nil {
		fullBody := publicBody + "\n" + protectedBody

		// 转换完整内容的附件
		if meta.Attachments != nil && len(meta.Attachments) > 0 {
			fullBody = transformAttachmentTagsByMeta(fullBody, slug, meta, false)
		}

		var fullBuf bytes.Buffer
		if err := mdConv.Convert([]byte(fullBody), &fullBuf); err != nil {
			return post{}, nil, err
		}

		protInfo.fullContent = fullBuf.String()
	}

	// 生成HTML文件
	dirOut := filepath.Join(outDir, category)
	if err := os.MkdirAll(dirOut, 0755); err != nil {
		return post{}, nil, err
	}

	htmlName := slug + ".html"
	outFile := filepath.Join(dirOut, htmlName)

	f, err := os.Create(outFile)
	if err != nil {
		return post{}, nil, err
	}
	defer f.Close()

	title := meta.Title
	if strings.TrimSpace(title) == "" {
		if t := extractFirstHeading(body); t != "" {
			title = t
		} else {
			title = slug
		}
	}

	st, _ := os.Stat(noteMDPath)
	modTime := st.ModTime()

	templateData := map[string]interface{}{
		"Title":       title,
		"Body":        template.HTML(publicBuf.String()),
		"Category":    category,
		"Date":        modTime,
		"IsProtected": isProtected,
		"PostPath":    postPath,
		"Description": description,
		"BaseURL":     baseURL,
	}

	if err := postTpl.Execute(f, templateData); err != nil {
		return post{}, nil, err
	}

	// 复制附件
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

	return post{
		Category:    category,
		Title:       title,
		File:        filepath.ToSlash(filepath.Join(category, htmlName)),
		Date:        modTime,
		IsProtected: isProtected,
		Description: description,
	}, protInfo, nil
}

// 查找内容中的附件ID
func findAttachmentIDs(content string) []string {
	var ids []string
	matches := attachmentRefRe.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		if len(match) > 1 {
			ids = append(ids, match[1])
		}
	}
	return ids
}

func buildIndex(posts []post) error {
	groups := map[string][]post{}
	for _, p := range posts {
		groups[p.Category] = append(groups[p.Category], p)
	}

	var cats []categoryGroup
	for cat, ps := range groups {
		sort.Slice(ps, func(i, j int) bool { return ps[i].Date.After(ps[j].Date) })
		cats = append(cats, categoryGroup{Name: cat, Posts: ps})
	}

	sort.Slice(cats, func(i, j int) bool { return cats[i].Name < cats[j].Name })

	f, err := os.Create(filepath.Join(outDir, "index.html"))
	if err != nil {
		return err
	}
	defer f.Close()

	return indexTpl.Execute(f, map[string]interface{}{
		"Categories": cats,
		"Now":        time.Now().Format("2006-01-02 15:04"),
	})
}

// Helper functions
func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

func hashString(s string) string {
	hash := md5.Sum([]byte(s))
	return hex.EncodeToString(hash[:])[:16]
}

func generateAuthToken(path, password string) string {
	combined := path + ":" + password + ":" + time.Now().Format("2006-01-02")
	hash := sha256.Sum256([]byte(combined))
	return hex.EncodeToString(hash[:])
}

func verifyAuthCookie(path, token string) bool {
	protectedMu.RLock()
	_, ok := protectedPosts[path]
	protectedMu.RUnlock()

	if !ok {
		return false
	}

	return len(token) == 64
}

func copyBaseAssets() error {
	baseDir := filepath.Join("tpl", "base")
	fi, err := os.Stat(baseDir)
	if err != nil {
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

		if filepath.Dir(path) != srcDir {
			return nil
		}

		if strings.HasSuffix(strings.ToLower(d.Name()), ".txt") || strings.HasSuffix(strings.ToLower(d.Name()), ".md") {
			return nil
		}

		dest := filepath.Join(outDir, d.Name())
		return copyFile(path, dest)
	})
}

func copyFile(src, dest string) error {
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

	addTree(srcDir)
	if _, err := os.Stat("tpl"); err == nil {
		addTree("tpl")
	}

	var debounceTimer *time.Timer
	var debounceMu sync.Mutex

	triggerRebuild := func() {
		debounceMu.Lock()
		defer debounceMu.Unlock()

		if debounceTimer != nil {
			debounceTimer.Stop()
		}

		debounceTimer = time.AfterFunc(500*time.Millisecond, func() {
			log.Printf("Triggering rebuild...")
			if err := rebuildAll(); err != nil {
				log.Printf("Rebuild error: %v", err)
			} else {
				log.Printf("Rebuild successful")
			}
		})
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
				triggerRebuild()
			}

		case err, ok := <-w.Errors:
			if !ok {
				return
			}
			log.Printf("Watch error: %v", err)
		}
	}
}

var attachmentRefRe = regexp.MustCompile(`\[Attachment:\s*([\p{Han}A-Za-z0-9_.-]+)\]`)

func transformAttachmentTagsByMeta(body, slug string, meta NoteMeta, isProtected bool) string {
	return attachmentRefRe.ReplaceAllStringFunc(body, func(match string) string {
		m := attachmentRefRe.FindStringSubmatch(match)
		fmt.Println("Found attachment tag:", match, "->", m)
		if len(m) != 2 {
			return match
		}
		originalFilename := m[1]
		att, ok := meta.Attachments[originalFilename]
		if !ok {
			return match
		}

		rel := filepath.ToSlash(filepath.Join(slug, att.Path))
		ext_v2 := strings.ToLower(filepath.Ext(att.OriginalFilename))

		fmt.Println("Processing attachment:", att.OriginalFilename, "as", rel, "ext:", ext_v2)

		switch ext_v2 {
		case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp", ".svg":
			alt := att.OriginalFilename
			if alt == "" {
				alt = att.SavedAs
			}
			return `<img src="` + rel + `" alt="` + htmlEscapeAttr(alt) + `" loading="lazy" />`
		case ".pdf":
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
			uniqueID := fmt.Sprintf("nes_%x", md5.Sum([]byte(originalFilename)))[:12]
			name := att.OriginalFilename
			if name == "" {
				name = filepath.Base(att.Path)
			}
			return parser.GenerateNESPlayerHTML(uniqueID, rel, htmlEscape(name), strings.ToLower(filepath.Ext(att.Path)))
		case ".pbp", ".chd", ".7z":
			uniqueID := fmt.Sprintf("pbp_%x", md5.Sum([]byte(originalFilename)))[:12]
			name := att.OriginalFilename
			if name == "" {
				name = filepath.Base(att.Path)
			}
			return parser.GeneratePlayStationPlayerHTML(uniqueID, rel, htmlEscape(name), strings.ToLower(filepath.Ext(att.Path)))
		case ".zip":
			uniqueID := fmt.Sprintf("arc_%x", md5.Sum([]byte(originalFilename)))[:12]
			name := att.OriginalFilename
			if name == "" {
				name = filepath.Base(att.Path)
			}
			// name

			if slices.Contains([]string{"kof97.zip", "pbobbl2n.zip", "kof98.zip"}, name) {
				return parser.GenerateArcadePlayerHTML(uniqueID, rel, htmlEscape(name), strings.ToLower(filepath.Ext(att.Path)))
			}
			return `<a href="` + rel + `" download>` + htmlEscape(name) + `</a>`
		case ".gba":
			uniqueID := fmt.Sprintf("gba_%x", md5.Sum([]byte(originalFilename)))[:12]
			name := att.OriginalFilename
			if name == "" {
				name = filepath.Base(att.Path)
			}
			return parser.GenerateGBAPlayerHTML(uniqueID, rel, htmlEscape(name), strings.ToLower(filepath.Ext(att.Path)))
		case ".jar":
			uniqueID := fmt.Sprintf("jar_%x", md5.Sum([]byte(originalFilename)))[:12]
			name := att.OriginalFilename
			if name == "" {
				name = filepath.Base(att.Path)
			}
			return parser.GenerateJARPlayerHTML(uniqueID, rel, htmlEscape(name), strings.ToLower(filepath.Ext(att.Path)))
		case ".md", ".smd", ".gen", ".bin", ".sms", ".gg", ".32x", ".cue", ".iso":
			uniqueID := fmt.Sprintf("sega_%x", md5.Sum([]byte(originalFilename)))[:12]
			name := att.OriginalFilename
			if name == "" {
				name = filepath.Base(att.Path)
			}
			return parser.GenerateSegaMDPlayerHTML(uniqueID, rel, htmlEscape(name), strings.ToLower(filepath.Ext(att.Path)))
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
			return strings.TrimSpace(strings.TrimLeft(l, "#"))
		}
	}
	return ""
}

func htmlEscapeAttr(s string) string {
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

// NoteMeta and related structures remain the same
type NoteMeta struct {
	NoteID              int
	Title               string
	Folder              string
	AttachmentsDeclared *int
	Attachments         map[string]AttachmentMeta
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
	reNoteID          = regexp.MustCompile(`(?m)^Note ID:\s*(\d+)\s*$`)
	reTitle           = regexp.MustCompile(`(?m)^Title:\s*(.+)\s*$`)
	reFolder          = regexp.MustCompile(`(?m)^Folder:\s*(.+)\s*$`)
	reAttachCount     = regexp.MustCompile(`(?m)^Attachments:\s*(\d+)\s*file\(s\)\s*$`)
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

	if ms := reAttachCount.FindStringSubmatch(t); len(ms) == 2 {
		if v, err := strconv.Atoi(ms[1]); err == nil {
			m.AttachmentsDeclared = &v
		}
	}

	attachSec := extractAttachmentSection(t)
	if strings.TrimSpace(attachSec) == "" {
		return m, nil
	}

	lines := strings.Split(attachSec, "\n")
	var cur *AttachmentMeta

	flush := func() {
		if cur != nil {
			if cur.ID != "" {
				//TODO 因为部分加载器无法使用ID获取的文件进行加载，所以这里改为使用OriginalFilename作为key
				m.Attachments[cur.OriginalFilename] = *cur
			} else if cur.SavedAs != "" {
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
	i := strings.Index(s, "Attachment Details:")
	if i == -1 {
		return ""
	}
	return s[i:]
}

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

	add(filepath.Join(outAbs, "index.html"))
	add(filepath.Join(outAbs, "sitemap.xml"))
	add(filepath.Join(outAbs, "rss.xml"))
	add(filepath.Join(outAbs, "robots.txt"))

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
				add(filepath.Join(outAbs, cat, pe.Name()+".html"))
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

func pruneDirRecursive(root, dir string, expected map[string]struct{}) error {
	ents, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, e := range ents {
		p := filepath.Join(dir, e.Name())
		if e.IsDir() {
			if err := pruneDirRecursive(root, p, expected); err != nil {
				return err
			}
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
