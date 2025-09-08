package main

import (
	"bytes"
	"flag"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/yuin/goldmark"
)

var (
	srcDir   string
	outDir   string
	addr     string
	postTpl  *template.Template
	indexTpl *template.Template
	mdConv   goldmark.Markdown
)

func init() {
	flag.StringVar(&srcDir, "src", "./pre", "Markdown & asset source directory")
	flag.StringVar(&outDir, "out", "./public", "Output directory")
	flag.StringVar(&addr, "addr", ":8080", "HTTP listen address")
}

func main() {
	flag.Parse()
	loadTemplates()
	mdConv = goldmark.New() // 默认即可，足够精简快速

	if err := rebuildAll(); err != nil {
		log.Fatalf("initial build failed: %v", err)
	}
	log.Println("initial build done.")

	go watch()

	log.Printf("serving %s at http://%s …", outDir, addr)
	http.Handle("/", http.FileServer(http.Dir(outDir)))
	log.Fatal(http.ListenAndServe(addr, nil))
}

/* ---------------------- templates ---------------------- */

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

/* ---------------------- rebuild ------------------------ */

type post struct {
	Title string
	File  string
	Date  time.Time
}

func rebuildAll() error {
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return err
	}
	if err := copyAssets(); err != nil {
		return err
	}

	var posts []post
	err := filepath.WalkDir(srcDir, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil || d.IsDir() || !strings.HasSuffix(d.Name(), ".md") {
			return walkErr
		}
		p, err := buildPage(path)
		if err != nil {
			return err
		}
		posts = append(posts, p)
		return nil
	})
	if err != nil {
		return err
	}
	return buildIndex(posts)
}

func buildPage(mdPath string) (post, error) {
	data, err := os.ReadFile(mdPath)
	if err != nil {
		return post{}, err
	}

	title := firstHeading(string(data), filepath.Base(mdPath))

	var buf bytes.Buffer
	if err := mdConv.Convert(data, &buf); err != nil {
		return post{}, err
	}

	// 复制文中引用图片
	ensureImgs(string(data))

	htmlName := strings.TrimSuffix(filepath.Base(mdPath), ".md") + ".html"
	out := filepath.Join(outDir, htmlName)
	f, err := os.Create(out)
	if err != nil {
		return post{}, err
	}
	defer f.Close()

	if err := postTpl.Execute(f, map[string]any{
		"Title": title,
		"Body":  template.HTML(buf.String()),
	}); err != nil {
		return post{}, err
	}

	st, _ := os.Stat(mdPath)
	return post{Title: title, File: htmlName, Date: st.ModTime()}, nil
}

func buildIndex(posts []post) error {
	sort.Slice(posts, func(i, j int) bool { return posts[i].Date.After(posts[j].Date) })
	f, err := os.Create(filepath.Join(outDir, "index.html"))
	if err != nil {
		return err
	}
	defer f.Close()
	return indexTpl.Execute(f, map[string]any{
		"Posts": posts,
		"Now":   time.Now().Format("2006-01-02 15:04"),
	})
}

/* ---------------------- assets ------------------------- */

func copyAssets() error {
	return filepath.WalkDir(srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || strings.HasSuffix(d.Name(), ".md") {
			return err
		}
		rel, _ := filepath.Rel(srcDir, path)
		dest := filepath.Join(outDir, rel)
		if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
			return err
		}
		return copyFile(path, dest)
	})
}

func copyFile(src, dest string) error {
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

var imgRe = regexp.MustCompile(`!\[[^\]]*\]\(([^)]+)\)`)

func ensureImgs(md string) {
	for _, m := range imgRe.FindAllStringSubmatch(md, -1) {
		p := strings.TrimSpace(m[1])
		if p == "" || strings.HasPrefix(p, "http") {
			continue
		}
		src := filepath.Join(srcDir, p)
		dest := filepath.Join(outDir, p)
		if _, err := os.Stat(dest); err == nil {
			continue
		}
		if err := os.MkdirAll(filepath.Dir(dest), 0755); err == nil {
			_ = copyFile(src, dest)
		}
	}
}

/* ---------------------- watcher ------------------------ */

func watch() {
	w, _ := fsnotify.NewWatcher()
	defer w.Close()
	_ = w.Add(srcDir)
	for {
		select {
		case ev := <-w.Events:
			if ev.Has(fsnotify.Write | fsnotify.Create | fsnotify.Remove | fsnotify.Rename) {
				log.Printf("change: %s", ev)
				if err := rebuildAll(); err != nil {
					log.Printf("rebuild error: %v", err)
				}
			}
		case err := <-w.Errors:
			log.Printf("watch error: %v", err)
		}
	}
}

/* ---------------------- helpers ------------------------ */

func firstHeading(md, fallback string) string {
	for _, l := range strings.Split(md, "\n") {
		l = strings.TrimSpace(l)
		if strings.HasPrefix(l, "#") {
			return strings.TrimSpace(strings.TrimLeft(l, "#"))
		}
	}
	return strings.TrimSuffix(fallback, ".md")
}
