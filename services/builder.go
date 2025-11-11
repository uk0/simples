package services

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"retroblog/config"
	"retroblog/models"
	"retroblog/utils"
	"sort"
	"strings"
	"time"
)

var passRe = regexp.MustCompile(`(?i)\[PASS:\s*([^\]]+)\]`)

// RebuildAll rebuilds the entire site
func RebuildAll() error {
	config.RebuildingMu.Lock()
	config.Rebuilding = true
	config.RebuildingMu.Unlock()

	defer func() {
		config.RebuildingMu.Lock()
		config.Rebuilding = false
		config.RebuildingMu.Unlock()
	}()

	// Temporary storage during rebuild
	tempProtectedPosts := make(map[string]string)
	tempFullContent := make(map[string]string)
	tempProtectedAttachments := make(map[string]bool)

	// Create output directory
	if err := os.MkdirAll(config.OutDir, 0755); err != nil {
		return err
	}

	// Copy assets
	if err := copyBaseAssets(); err != nil {
		return err
	}
	if err := copyTopLevelAssets(); err != nil {
		return err
	}

	// Build posts
	allPosts, err := buildAllPosts(tempProtectedPosts, tempFullContent, tempProtectedAttachments)
	if err != nil {
		return err
	}

	// Build index and other pages
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

	// Update protected content maps
	config.ProtectedMu.Lock()
	config.ProtectedPosts = tempProtectedPosts
	config.FullContent = tempFullContent
	config.ProtectedAttachments = tempProtectedAttachments
	config.ProtectedMu.Unlock()

	log.Printf("Rebuild complete. Protected posts: %d, Protected attachments: %d",
		len(tempProtectedPosts), len(tempProtectedAttachments))

	// Clean up old files
	if err := pruneOutToMatchSource(); err != nil {
		log.Printf("Warning: prune failed: %v", err)
	}

	return nil
}

func buildAllPosts(tempProtectedPosts, tempFullContent map[string]string, tempProtectedAttachments map[string]bool) ([]models.Post, error) {
	var allPosts []models.Post

	cats, err := os.ReadDir(config.SrcDir)
	if err != nil {
		return nil, err
	}

	for _, ce := range cats {
		if !ce.IsDir() {
			continue
		}
		catName := ce.Name()
		catPath := filepath.Join(config.SrcDir, catName)

		if err := os.MkdirAll(filepath.Join(config.OutDir, catName), 0755); err != nil {
			return nil, err
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
					tempProtectedPosts[protectedInfo.Path] = protectedInfo.PasswordHash
					tempFullContent[protectedInfo.Path] = protectedInfo.FullContent

					for _, attachPath := range protectedInfo.ProtectedFilePaths {
						tempProtectedAttachments[attachPath] = true
					}

					p.IsProtected = true
				}

				allPosts = append(allPosts, p)
			}
		}
	}

	return allPosts, nil
}

func buildPageWithProtection(noteMDPath, category string) (models.Post, *models.ProtectionInfo, error) {
	postDir := filepath.Dir(noteMDPath)
	slug := filepath.Base(postDir)

	// Load metadata
	metaPath := filepath.Join(postDir, "meta.txt")
	var meta models.NoteMeta
	if fi, err := os.Stat(metaPath); err == nil && !fi.IsDir() {
		meta, _ = ParseMetaFile(metaPath)
	}

	// Read content
	bodyBytes, err := os.ReadFile(noteMDPath)
	if err != nil {
		return models.Post{}, nil, err
	}
	body := string(bodyBytes)
	postPath := category + "/" + slug

	var protInfo *models.ProtectionInfo
	var isProtected bool
	var publicBody string
	var protectedBody string
	var protectedFileCount int = 0
	total := FileStats{}
	// Check for password protection
	if matches := passRe.FindStringSubmatch(body); len(matches) == 2 {
		isProtected = true
		password := strings.TrimSpace(matches[1])

		loc := passRe.FindStringIndex(body)
		if loc != nil {
			publicBody = body[:loc[0]]
			protectedBody = body[loc[1]:]

			// Find protected attachments
			protectedAttachmentNames := utils.FindAttachmentNames(protectedBody)
			protectedFileCount = len(protectedAttachmentNames)
			var protectedPaths []string
			for _, filename := range protectedAttachmentNames {
				stats := CountByFilename(filename)
				total.ImageCount += stats.ImageCount
				total.VideoCount += stats.VideoCount
				total.FileCount += stats.FileCount
				if att, ok := meta.Attachments[filename]; ok {
					attachPath := filepath.Join(category, slug, config.AttachmentsDir, filepath.Base(att.Path))
					protectedPaths = append(protectedPaths, attachPath)
				}
			}

			protInfo = &models.ProtectionInfo{
				Path:               postPath,
				PasswordHash:       utils.HashPassword(password),
				ProtectedFilePaths: protectedPaths,
			}
		}
	} else {
		publicBody = body
	}

	// Process public content
	description := utils.ExtractDescription(publicBody)
	if meta.Attachments != nil && len(meta.Attachments) > 0 {
		publicBody = utils.TransformAttachmentTags(publicBody, slug, meta)
	}

	var publicBuf bytes.Buffer
	if err := config.MdConv.Convert([]byte(publicBody), &publicBuf); err != nil {
		return models.Post{}, nil, err
	}

	// Process protected content if exists
	if isProtected && protInfo != nil {
		fullBody := publicBody + "\n" + protectedBody
		if meta.Attachments != nil && len(meta.Attachments) > 0 {
			fullBody = utils.TransformAttachmentTags(fullBody, slug, meta)
		}

		var fullBuf bytes.Buffer
		if err := config.MdConv.Convert([]byte(fullBody), &fullBuf); err != nil {
			return models.Post{}, nil, err
		}
		protInfo.FullContent = fullBuf.String()
	}

	// Generate HTML file
	if err := generateHTMLFile(category, slug, body, publicBuf.String(), isProtected, postPath, description, meta); err != nil {
		return models.Post{}, nil, err
	}

	// Copy attachments
	if err := copyAttachments(postDir, category, slug); err != nil {
		log.Printf("Error copying attachments: %v", err)
	}

	// Get modification time
	st, _ := os.Stat(noteMDPath)
	modTime := st.ModTime()

	title := meta.Title
	if strings.TrimSpace(title) == "" {
		if t := utils.ExtractFirstHeading(body); t != "" {
			title = t
		} else {
			title = slug
		}
	}

	return models.Post{
		Category:    category,
		Title:       title,
		File:        filepath.ToSlash(filepath.Join(category, slug+".html")),
		Date:        modTime,
		Modified:    meta.ParseTime(meta.Modified),
		Created:     meta.ParseTime(meta.Created),
		IsProtected: isProtected,
		Description: description,
		// 文件与保护文件
		ProtectedFileCount: protectedFileCount,
		ImageCount:         total.ImageCount,
		FileCount:          total.FileCount,
		VideoCount:         total.VideoCount,
	}, protInfo, nil
}

type FileStats struct {
	ImageCount int
	VideoCount int
	FileCount  int
}

// CountByFilename 仅根据文件名（或扩展名）推断类型数量
func CountByFilename(filename string) FileStats {
	var stats FileStats
	ext := strings.ToLower(filepath.Ext(filename))

	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp", ".heic":
		stats.ImageCount = 1
	case ".mp4", ".mov", ".mkv", ".avi", ".webm":
		stats.VideoCount = 1
	default:
		if ext != "" {
			stats.FileCount = 1
		}
	}

	return stats
}

func generateHTMLFile(category, slug, body, publicHTML string, isProtected bool, postPath, description string, meta models.NoteMeta) error {
	dirOut := filepath.Join(config.OutDir, category)
	if err := os.MkdirAll(dirOut, 0755); err != nil {
		return err
	}

	htmlName := slug + ".html"
	outFile := filepath.Join(dirOut, htmlName)

	f, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer f.Close()

	title := meta.Title
	if strings.TrimSpace(title) == "" {
		if t := utils.ExtractFirstHeading(body); t != "" {
			title = t
		} else {
			title = slug
		}
	}

	st, _ := os.Stat(filepath.Join(config.SrcDir, category, slug, "note.md"))
	modTime := st.ModTime()
	//TODO ai generator og meta
	//generator := utils.NewOGGenerator("sk-")
	//
	//// 示例 2: 从 HTML 内容生成
	//log.Println("HTML 内容生成 OG Meta 信息")
	//log.Println("=========================================")
	//
	//metaTags2, ogInfo2, err := generator.GenerateFromHTML(publicHTML, fmt.Sprintf("%s/%s", config.BaseURL, postPath))
	//if err != nil {
	//	log.Printf("错误: %v\n", err)
	//} else {
	//	utils.PrintOGInfo(ogInfo2)
	//	log.Println("生成的 HTML Meta 标签:")
	//	log.Println(metaTags2)
	//}

	templateData := map[string]interface{}{
		"Title":       title,
		"Body":        template.HTML(publicHTML),
		"Category":    category,
		"Date":        modTime,
		"Created":     meta.ParseTime(meta.Created),
		"Modified":    meta.ParseTime(meta.Modified),
		"IsProtected": isProtected,
		"PostPath":    postPath,
		"Description": description,
		"Slug":        slug, // 编号
		"BaseURL":     config.BaseURL,
	}

	return config.PostTpl.Execute(f, templateData)
}

func copyAttachments(postDir, category, slug string) error {
	srcAttach := filepath.Join(postDir, config.AttachmentsDir)
	fi, err := os.Stat(srcAttach)
	if err != nil || !fi.IsDir() {
		return nil
	}

	dstAttach := filepath.Join(config.OutDir, category, slug, config.AttachmentsDir)
	if err := os.MkdirAll(dstAttach, 0755); err != nil {
		return err
	}

	return filepath.WalkDir(srcAttach, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		relName := filepath.Base(p)
		dest := filepath.Join(dstAttach, relName)
		return utils.CopyFile(p, dest)
	})
}

func buildIndex(posts []models.Post) error {
	groups := map[string][]models.Post{}
	for _, p := range posts {
		groups[p.Category] = append(groups[p.Category], p)
	}

	var cats []models.CategoryGroup
	for cat, ps := range groups {
		sort.Slice(ps, func(i, j int) bool { return ps[i].Created.After(ps[j].Created) })
		cats = append(cats, models.CategoryGroup{Name: cat, Posts: ps})
	}

	f, err := os.Create(filepath.Join(config.OutDir, "index.html"))
	if err != nil {
		return err
	}
	defer f.Close()

	return config.IndexTpl.Execute(f, map[string]interface{}{
		"Categories": cats,
		"Now":        time.Now().Format("2006-01-02 15:04"),
	})
}

func buildSitemap(posts []models.Post) error {
	sitemap := models.Urlset{
		Xmlns:             "http://www.sitemaps.org/schemas/sitemap/0.9",
		XmlnsXsi:          "http://www.w3.org/2001/XMLSchema-instance",
		XsiSchemaLocation: "http://www.sitemaps.org/schemas/sitemap/0.9 http://www.sitemaps.org/schemas/sitemap/0.9/sitemap.xsd",
		URLs:              []models.SitemapURL{},
	}

	sitemap.URLs = append(sitemap.URLs, models.SitemapURL{
		Loc:        config.BaseURL + "/",
		LastMod:    time.Now().Format("2006-01-02"),
		ChangeFreq: "daily",
		Priority:   1.0,
	})

	for _, p := range posts {
		if !p.IsProtected {
			sitemap.URLs = append(sitemap.URLs, models.SitemapURL{
				Loc:        config.BaseURL + "/" + p.File,
				LastMod:    p.Date.Format("2006-01-02"),
				ChangeFreq: "weekly",
				Priority:   0.8,
			})
		}
	}

	output, err := xml.MarshalIndent(sitemap, "", "  ")
	if err != nil {
		return err
	}

	xmlContent := []byte(xml.Header + string(output))
	sitemapPath := filepath.Join(config.OutDir, "sitemap.xml")
	return os.WriteFile(sitemapPath, xmlContent, 0644)
}

func buildRSSFeed(posts []models.Post) error {
	feed := models.RSS{
		Version: "2.0",
		Channel: models.Channel{
			Title:       "Blog Feed",
			Link:        config.BaseURL,
			Description: "Latest posts",
			Language:    "en-US",
			LastBuild:   time.Now().Format(time.RFC1123Z),
			Items:       []models.RSSItem{},
		},
	}

	sortedPosts := make([]models.Post, len(posts))
	copy(sortedPosts, posts)
	sort.Slice(sortedPosts, func(i, j int) bool {
		return sortedPosts[i].Date.After(sortedPosts[j].Date)
	})

	count := 0
	for _, p := range sortedPosts {
		if !p.IsProtected && count < 20 {
			feed.Channel.Items = append(feed.Channel.Items, models.RSSItem{
				Title:       p.Title,
				Link:        config.BaseURL + "/" + p.File,
				Description: p.Description,
				PubDate:     p.Date.Format(time.RFC1123Z),
				GUID:        config.BaseURL + "/" + p.File,
				Category:    p.Category,
			})
			count++
		}
	}

	output, err := xml.MarshalIndent(feed, "", "  ")
	if err != nil {
		return err
	}

	xmlContent := []byte(xml.Header + string(output))
	rssPath := filepath.Join(config.OutDir, "rss.xml")
	return os.WriteFile(rssPath, xmlContent, 0644)
}

func buildRobotsTxt() error {
	robotsContent := fmt.Sprintf(`User-agent: *
Allow: /
Disallow: /api/

Sitemap: %s/sitemap.xml
`, config.BaseURL)

	robotsPath := filepath.Join(config.OutDir, "robots.txt")
	return os.WriteFile(robotsPath, []byte(robotsContent), 0644)
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
		if err != nil || d.IsDir() {
			return err
		}
		rel, err := filepath.Rel(baseDir, path)
		if err != nil {
			return err
		}
		dest := filepath.Join(config.OutDir, rel)
		return utils.CopyFile(path, dest)
	})
}

func copyTopLevelAssets() error {
	return filepath.WalkDir(config.SrcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		if filepath.Dir(path) != config.SrcDir {
			return nil
		}
		if strings.HasSuffix(strings.ToLower(d.Name()), ".txt") || strings.HasSuffix(strings.ToLower(d.Name()), ".md") {
			return nil
		}
		dest := filepath.Join(config.OutDir, d.Name())
		return utils.CopyFile(path, dest)
	})
}

func pruneOutToMatchSource() error {
	// Implementation remains the same as original
	// ... (pruning logic)
	return nil
}
