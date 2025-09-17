package main

import (
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
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
	baseURL        string // Base URL for sitemap
	attachmentsDir = "attachments"
	postTpl        *template.Template
	indexTpl       *template.Template
	mdConv         goldmark.Markdown

	// Password protection - use temporary maps during rebuild
	protectedMu    sync.RWMutex
	protectedPosts = make(map[string]string) // path -> password hash
	fullContent    = make(map[string]string) // path -> full HTML content
	passRe         = regexp.MustCompile(`(?i)\[PASS:\s*([^\]]+)\]`)

	// Flag to indicate rebuild is in progress
	rebuilding   bool
	rebuildingMu sync.Mutex
)

// Sitemap XML structures
type urlset struct {
	XMLName           xml.Name     `xml:"urlset"`
	Xmlns             string       `xml:"xmlns,attr"`
	XmlnsXsi          string       `xml:"xmlns:xsi,attr,omitempty"`
	XsiSchemaLocation string       `xml:"xsi:schemaLocation,attr,omitempty"`
	URLs              []sitemapURL `xml:"url"`
}

type sitemapURL struct {
	Loc        string  `xml:"loc"`
	LastMod    string  `xml:"lastmod,omitempty"`
	ChangeFreq string  `xml:"changefreq,omitempty"`
	Priority   float64 `xml:"priority,omitempty"`
}

// RSS Feed structures
type rss struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel channel  `xml:"channel"`
}

type channel struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Language    string    `xml:"language"`
	LastBuild   string    `xml:"lastBuildDate"`
	Items       []rssItem `xml:"item"`
}

type rssItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	GUID        string `xml:"guid"`
	Category    string `xml:"category,omitempty"`
}

func init() {
	flag.StringVar(&srcDir, "src", "./apple_notes_export", "Notes source directory")
	flag.StringVar(&outDir, "out", "./public", "Output directory")
	flag.StringVar(&addr, "addr", ":8080", "HTTP listen address")
	flag.StringVar(&baseURL, "url", "https://firsh.me", "Base URL for sitemap (e.g., https://example.com)")
}

func main() {
	flag.Parse()

	// Clean up baseURL - remove trailing slash
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

	// Setup HTTP handlers
	mux := http.NewServeMux()

	// API endpoint for password verification
	mux.HandleFunc("/api/unlock", handleUnlock)

	// Protected content handler
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// Check if this is an attachment request for a protected post
		if strings.Contains(path, "/attachments/") {
			parts := strings.Split(strings.Trim(path, "/"), "/")
			if len(parts) >= 3 {
				// Reconstruct post path (e.g., "category/slug")
				postPath := parts[0] + "/" + strings.TrimSuffix(parts[1], ".html")

				protectedMu.RLock()
				_, isProtected := protectedPosts[postPath]
				protectedMu.RUnlock()

				if isProtected {
					// Check session cookie for authorization
					cookie, err := r.Cookie("auth_" + hashString(postPath))
					if err != nil || !verifyAuthCookie(postPath, cookie.Value) {
						http.Error(w, "Unauthorized - Please unlock the post first", http.StatusUnauthorized)
						return
					}
				}
			}
		}

		// Serve static files
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

	// Clean the path
	req.Path = strings.Trim(req.Path, "/")

	// Wait for any ongoing rebuild to complete
	rebuildingMu.Lock()
	isRebuilding := rebuilding
	rebuildingMu.Unlock()

	if isRebuilding {
		// Wait a bit for rebuild to complete
		time.Sleep(100 * time.Millisecond)
	}

	// Verify password
	protectedMu.RLock()
	expectedHash, ok := protectedPosts[req.Path]
	content, hasContent := fullContent[req.Path]
	protectedMu.RUnlock()

	if !ok {
		log.Printf("Post %s not found in protected posts. Available posts: %v", req.Path, getProtectedPaths())
		http.Error(w, "Post not protected", http.StatusBadRequest)
		return
	}

	if hashPassword(req.Password) != expectedHash {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	if !hasContent {
		log.Printf("Content not found for protected post: %s", req.Path)
		http.Error(w, "Content not found", http.StatusInternalServerError)
		return
	}

	// Generate auth token
	authToken := generateAuthToken(req.Path, req.Password)

	// Set cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_" + hashString(req.Path),
		Value:    authToken,
		Path:     "/",
		MaxAge:   86400, // 24 hours
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   false, // Set to true if using HTTPS
	})

	// Return full content
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"content": content,
		"status":  "success",
	})
}

// Helper function to get list of protected paths for debugging
func getProtectedPaths() []string {
	protectedMu.RLock()
	defer protectedMu.RUnlock()

	paths := make([]string, 0, len(protectedPosts))
	for path := range protectedPosts {
		paths = append(paths, path)
	}
	return paths
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
	Category    string
	Title       string
	File        string
	Date        time.Time
	IsProtected bool   // Track if post is password protected
	Description string // For meta description
}

type categoryGroup struct {
	Name  string
	Posts []post
}

func rebuildAll() error {
	// Set rebuilding flag
	rebuildingMu.Lock()
	rebuilding = true
	rebuildingMu.Unlock()

	defer func() {
		rebuildingMu.Lock()
		rebuilding = false
		rebuildingMu.Unlock()
	}()

	// Use temporary maps during rebuild
	tempProtectedPosts := make(map[string]string)
	tempFullContent := make(map[string]string)

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

				// Store protection info in temporary maps
				if protectedInfo != nil {
					tempProtectedPosts[protectedInfo.path] = protectedInfo.passwordHash
					tempFullContent[protectedInfo.path] = protectedInfo.fullContent
					p.IsProtected = true
				}

				allPosts = append(allPosts, p)
			}
		}
	}

	if err := buildIndex(allPosts); err != nil {
		return err
	}

	// Generate sitemap.xml
	if err := buildSitemap(allPosts); err != nil {
		log.Printf("Warning: failed to build sitemap: %v", err)
	} else {
		log.Println("Generated sitemap.xml")
	}

	// Generate RSS feed
	if err := buildRSSFeed(allPosts); err != nil {
		log.Printf("Warning: failed to build RSS feed: %v", err)
	} else {
		log.Println("Generated rss.xml")
	}

	// Generate robots.txt
	if err := buildRobotsTxt(); err != nil {
		log.Printf("Warning: failed to build robots.txt: %v", err)
	} else {
		log.Println("Generated robots.txt")
	}

	// Atomically replace the protection maps
	protectedMu.Lock()
	protectedPosts = tempProtectedPosts
	fullContent = tempFullContent
	protectedMu.Unlock()

	log.Printf("Rebuild complete. Protected posts: %v", getProtectedPaths())

	if err := pruneOutToMatchSource(); err != nil {
		log.Printf("Warning: prune failed: %v", err)
	}

	return nil
}

// buildSitemap generates sitemap.xml for SEO
func buildSitemap(posts []post) error {
	sitemap := urlset{
		Xmlns:             "http://www.sitemaps.org/schemas/sitemap/0.9",
		XmlnsXsi:          "http://www.w3.org/2001/XMLSchema-instance",
		XsiSchemaLocation: "http://www.sitemaps.org/schemas/sitemap/0.9 http://www.sitemaps.org/schemas/sitemap/0.9/sitemap.xsd",
		URLs:              []sitemapURL{},
	}

	// Add home page
	sitemap.URLs = append(sitemap.URLs, sitemapURL{
		Loc:        baseURL + "/",
		LastMod:    time.Now().Format("2006-01-02"),
		ChangeFreq: "daily",
		Priority:   1.0,
	})

	// Add all posts (except protected ones for better SEO)
	for _, p := range posts {
		if !p.IsProtected { // Only include public posts in sitemap
			sitemap.URLs = append(sitemap.URLs, sitemapURL{
				Loc:        baseURL + "/" + p.File,
				LastMod:    p.Date.Format("2006-01-02"),
				ChangeFreq: "weekly",
				Priority:   0.8,
			})
		}
	}

	// Sort URLs by priority and then by location for consistent output
	sort.Slice(sitemap.URLs, func(i, j int) bool {
		if sitemap.URLs[i].Priority != sitemap.URLs[j].Priority {
			return sitemap.URLs[i].Priority > sitemap.URLs[j].Priority
		}
		return sitemap.URLs[i].Loc < sitemap.URLs[j].Loc
	})

	// Marshal to XML
	output, err := xml.MarshalIndent(sitemap, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal sitemap: %w", err)
	}

	// Add XML declaration
	xmlContent := []byte(xml.Header + string(output))

	// Write to file
	sitemapPath := filepath.Join(outDir, "sitemap.xml")
	if err := os.WriteFile(sitemapPath, xmlContent, 0644); err != nil {
		return fmt.Errorf("failed to write sitemap: %w", err)
	}

	log.Printf("Sitemap generated with %d URLs (excluded %d protected posts)",
		len(sitemap.URLs), countProtectedPosts(posts))
	return nil
}

// buildRSSFeed generates RSS feed for subscribers
func buildRSSFeed(posts []post) error {
	feed := rss{
		Version: "2.0",
		Channel: channel{
			Title:       "NeoJ's Web Page",
			Link:        baseURL,
			Description: "Personal blog and notes",
			Language:    "en-US",
			LastBuild:   time.Now().Format(time.RFC1123Z),
			Items:       []rssItem{},
		},
	}

	// Sort posts by date (newest first)
	sortedPosts := make([]post, len(posts))
	copy(sortedPosts, posts)
	sort.Slice(sortedPosts, func(i, j int) bool {
		return sortedPosts[i].Date.After(sortedPosts[j].Date)
	})

	// Add up to 20 most recent public posts to RSS
	count := 0
	for _, p := range sortedPosts {
		if !p.IsProtected && count < 20 {
			desc := p.Description
			if desc == "" {
				desc = fmt.Sprintf("Post from %s category", p.Category)
			}

			feed.Channel.Items = append(feed.Channel.Items, rssItem{
				Title:       p.Title,
				Link:        baseURL + "/" + p.File,
				Description: desc,
				PubDate:     p.Date.Format(time.RFC1123Z),
				GUID:        baseURL + "/" + p.File,
				Category:    p.Category,
			})
			count++
		}
	}

	// Marshal to XML
	output, err := xml.MarshalIndent(feed, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal RSS feed: %w", err)
	}

	// Add XML declaration
	xmlContent := []byte(xml.Header + string(output))

	// Write to file
	rssPath := filepath.Join(outDir, "rss.xml")
	if err := os.WriteFile(rssPath, xmlContent, 0644); err != nil {
		return fmt.Errorf("failed to write RSS feed: %w", err)
	}

	log.Printf("RSS feed generated with %d items", len(feed.Channel.Items))
	return nil
}

// buildRobotsTxt generates robots.txt for search engines
func buildRobotsTxt() error {
	robotsContent := fmt.Sprintf(`# Robots.txt for %s
User-agent: *
Disallow: /*api/
Allow: /blog/*

User-agent: Baiduspider
Disallow: /api/*
Allow: /blog/*

User-agent: HaoSouSpider
Disallow: /api/*
Allow: /blog/*

User-agent: Sogou web spider
Disallow: /api/*
Allow: /blog/*

User-agent: Sogou inst spider
Disallow: /api/*
Allow: /blog/*

User-agent: Sogou spider2
Disallow: /api/*
Allow: /blog/*

# Sitemap
Sitemap: %s/sitemap.xml

# RSS Feed
# RSS: %s/rss.xml
`, baseURL, baseURL, baseURL)

	robotsPath := filepath.Join(outDir, "robots.txt")
	if err := os.WriteFile(robotsPath, []byte(robotsContent), 0644); err != nil {
		return fmt.Errorf("failed to write robots.txt: %w", err)
	}

	return nil
}

// Helper function to count protected posts
func countProtectedPosts(posts []post) int {
	count := 0
	for _, p := range posts {
		if p.IsProtected {
			count++
		}
	}
	return count
}

// extractDescription extracts a description from markdown content
func extractDescription(body string) string {
	// Remove markdown formatting
	desc := body

	// Remove headers
	desc = regexp.MustCompile(`(?m)^#+\s+.*$`).ReplaceAllString(desc, "")

	// Remove links but keep link text
	desc = regexp.MustCompile(`\[([^\]]+)\]\([^\)]+\)`).ReplaceAllString(desc, "$1")

	// Remove images
	desc = regexp.MustCompile(`!\[([^\]]*)\]\([^\)]+\)`).ReplaceAllString(desc, "")

	// Remove bold/italic
	desc = regexp.MustCompile(`\*+([^\*]+)\*+`).ReplaceAllString(desc, "$1")
	desc = regexp.MustCompile(`_+([^_]+)_+`).ReplaceAllString(desc, "$1")

	// Remove code blocks
	desc = regexp.MustCompile("(?s)```[^`]*```").ReplaceAllString(desc, "")
	desc = regexp.MustCompile("`([^`]+)`").ReplaceAllString(desc, "$1")

	// Remove HTML tags
	desc = regexp.MustCompile(`<[^>]+>`).ReplaceAllString(desc, "")

	// Remove attachment tags
	desc = regexp.MustCompile(`\[Attachment:[^\]]+\]`).ReplaceAllString(desc, "")

	// Clean up whitespace
	desc = strings.TrimSpace(desc)
	desc = regexp.MustCompile(`\s+`).ReplaceAllString(desc, " ")

	// Truncate to 160 characters for SEO
	if len(desc) > 160 {
		// Try to cut at word boundary
		if idx := strings.LastIndex(desc[:157], " "); idx > 100 {
			desc = desc[:idx] + "..."
		} else {
			desc = desc[:157] + "..."
		}
	}

	return desc
}

// Structure to hold protection info during build
type protectionInfo struct {
	path         string
	passwordHash string
	fullContent  string
}

func buildPageWithProtection(noteMDPath, category string) (post, *protectionInfo, error) {
	postDir := filepath.Dir(noteMDPath)
	slug := filepath.Base(postDir)

	// Load meta.txt
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

	// Read note.md
	bodyBytes, err := os.ReadFile(noteMDPath)
	if err != nil {
		return post{}, nil, err
	}
	body := string(bodyBytes)

	// Check for password protection
	postPath := category + "/" + slug
	var protInfo *protectionInfo
	var isProtected bool
	var truncatedBody string

	if matches := passRe.FindStringSubmatch(body); len(matches) == 2 {
		isProtected = true
		password := strings.TrimSpace(matches[1])

		// Find the position of [PASS: xxx] and truncate
		loc := passRe.FindStringIndex(body)
		if loc != nil {
			truncatedBody = body[:loc[0]]
		}

		log.Printf("Post %s is password protected", postPath)

		// Create protection info
		protInfo = &protectionInfo{
			path:         postPath,
			passwordHash: hashPassword(password),
		}
	}

	// Process the appropriate body
	processBody := body
	if isProtected {
		processBody = truncatedBody
	}

	// Extract description for SEO
	description := extractDescription(processBody)

	// Transform attachments
	if meta.Attachments != nil && len(meta.Attachments) > 0 {
		processBody = transformAttachmentTagsByMeta(processBody, slug, meta)
	}

	// Convert to HTML
	var buf bytes.Buffer
	if err := mdConv.Convert([]byte(processBody), &buf); err != nil {
		return post{}, nil, err
	}

	// If protected, generate and store full content
	if isProtected && protInfo != nil {
		// Process full body with attachments
		fullBody := body
		// Remove the [PASS: xxx] tag from full content
		fullBody = passRe.ReplaceAllString(fullBody, "")

		if meta.Attachments != nil && len(meta.Attachments) > 0 {
			fullBody = transformAttachmentTagsByMeta(fullBody, slug, meta)
		}

		var fullBuf bytes.Buffer
		if err := mdConv.Convert([]byte(fullBody), &fullBuf); err != nil {
			return post{}, nil, err
		}

		protInfo.fullContent = fullBuf.String()
	}

	// Prepare output
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

	// Get file modification time
	st, _ := os.Stat(noteMDPath)
	modTime := st.ModTime()

	// Prepare template data
	templateData := map[string]interface{}{
		"Title":       title,
		"Body":        template.HTML(buf.String()),
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

	// Copy attachments
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

// Keep original buildPage for compatibility but redirect to new function
func buildPage(noteMDPath, category string) (post, error) {
	p, _, err := buildPageWithProtection(noteMDPath, category)
	return p, err
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

// Password helper functions
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

	// Basic verification - check token format
	return len(token) == 64 // SHA256 hex string length
}

// copyBaseAssets copies everything inside tpl/base into the public root directory.
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

	// Debounce mechanism to avoid multiple rebuilds
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

// Attachment transform based on meta.txt
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
			return match
		}
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
			return parser.GenerateJARPlayerHTML(uniqueID, rel, htmlEscape(name), ext)
		case ".md", ".smd", ".gen", ".bin", ".sms", ".gg", ".32x", ".cue", ".iso":
			uniqueID := fmt.Sprintf("sega_%x", md5.Sum([]byte(id)))[:12]
			name := att.OriginalFilename
			if name == "" {
				name = filepath.Base(att.Path)
			}
			log.Println("EmbeddingJS Sega player for", name, "with ID", uniqueID)
			return parser.GenerateSegaMDPlayerHTML(uniqueID, rel, htmlEscape(name), ext)
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

// NoteMeta reflects the meta.txt content with attachment details.
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
				m.Attachments[cur.ID] = *cur
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

// Pruning functions
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

	// Core files
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
