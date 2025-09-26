package models

import "time"

// Post represents a blog post
type Post struct {
	Category    string
	Title       string
	File        string
	Date        time.Time
	IsProtected bool
	Description string
}

// CategoryGroup represents a group of posts by category
type CategoryGroup struct {
	Name  string
	Posts []Post
}

// ProtectionInfo contains protection details for a post
type ProtectionInfo struct {
	Path               string
	PasswordHash       string
	FullContent        string
	ProtectedFilePaths []string
}

// NoteMeta reflects the meta.txt content with attachment details
type NoteMeta struct {
	NoteID              int
	Title               string
	Folder              string
	AttachmentsDeclared *int
	Attachments         map[string]AttachmentMeta
}

// AttachmentMeta contains metadata for an attachment
type AttachmentMeta struct {
	ID               string
	SavedAs          string
	OriginalFilename string
	SourceFile       string
	Type             string
	Position         int
	Path             string
}

// Sitemap structures
type Urlset struct {
	XMLName           string       `xml:"urlset"`
	Xmlns             string       `xml:"xmlns,attr"`
	XmlnsXsi          string       `xml:"xmlns:xsi,attr,omitempty"`
	XsiSchemaLocation string       `xml:"xsi:schemaLocation,attr,omitempty"`
	URLs              []SitemapURL `xml:"url"`
}

type SitemapURL struct {
	Loc        string  `xml:"loc"`
	LastMod    string  `xml:"lastmod,omitempty"`
	ChangeFreq string  `xml:"changefreq,omitempty"`
	Priority   float64 `xml:"priority,omitempty"`
}

// RSS structures
type RSS struct {
	XMLName string  `xml:"rss"`
	Version string  `xml:"version,attr"`
	Channel Channel `xml:"channel"`
}

type Channel struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Language    string    `xml:"language"`
	LastBuild   string    `xml:"lastBuildDate"`
	Items       []RSSItem `xml:"item"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	GUID        string `xml:"guid"`
	Category    string `xml:"category,omitempty"`
}
