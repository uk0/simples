package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// Build sitemap.xml
func buildSitemap(posts []post) error {
	sitemap := urlset{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  []sitemapURL{},
	}

	sitemap.URLs = append(sitemap.URLs, sitemapURL{
		Loc:        baseURL + "/",
		LastMod:    time.Now().Format("2006-01-02"),
		ChangeFreq: "daily",
		Priority:   1.0,
	})

	for _, p := range posts {
		if !p.IsProtected {
			sitemap.URLs = append(sitemap.URLs, sitemapURL{
				Loc:        baseURL + "/" + p.File,
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
	sitemapPath := filepath.Join(outDir, "sitemap.xml")
	return os.WriteFile(sitemapPath, xmlContent, 0644)
}

// Build RSS feed
func buildRSSFeed(posts []post) error {
	feed := rss{
		Version: "2.0",
		Channel: channel{
			Title:       "Blog Feed",
			Link:        baseURL,
			Description: "Latest posts",
			Language:    "en-US",
			LastBuild:   time.Now().Format(time.RFC1123Z),
			Items:       []rssItem{},
		},
	}

	sortedPosts := make([]post, len(posts))
	copy(sortedPosts, posts)
	sort.Slice(sortedPosts, func(i, j int) bool {
		return sortedPosts[i].Date.After(sortedPosts[j].Date)
	})

	count := 0
	for _, p := range sortedPosts {
		if !p.IsProtected && count < 20 {
			feed.Channel.Items = append(feed.Channel.Items, rssItem{
				Title:       p.Title,
				Link:        baseURL + "/" + p.File,
				Description: p.Description,
				PubDate:     p.Date.Format(time.RFC1123Z),
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
	rssPath := filepath.Join(outDir, "rss.xml")
	return os.WriteFile(rssPath, xmlContent, 0644)
}

// Build robots.txt
func buildRobotsTxt() error {
	robotsContent := fmt.Sprintf(`User-agent: *
Allow: /
Disallow: /api/

Sitemap: %s/sitemap.xml
`, baseURL)

	robotsPath := filepath.Join(outDir, "robots.txt")
	return os.WriteFile(robotsPath, []byte(robotsContent), 0644)
}
