package main

import (
	"encoding/xml"
	"time"
)

// 修改 protectionInfo 结构
type protectionInfo struct {
	path               string   // 文章路径
	passwordHash       string   // 密码哈希
	fullContent        string   // 完整内容（包含受保护部分）
	protectedFilePaths []string // 受保护的附件文件路径列表
}

type post struct {
	Category    string
	Title       string
	File        string
	Date        time.Time
	IsProtected bool
	Description string
}

type categoryGroup struct {
	Name  string
	Posts []post
}

// Sitemap structures
type urlset struct {
	XMLName xml.Name     `xml:"urlset"`
	Xmlns   string       `xml:"xmlns,attr"`
	URLs    []sitemapURL `xml:"url"`
}

type sitemapURL struct {
	Loc        string  `xml:"loc"`
	LastMod    string  `xml:"lastmod,omitempty"`
	ChangeFreq string  `xml:"changefreq,omitempty"`
	Priority   float64 `xml:"priority,omitempty"`
}

// RSS structures
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
	Category    string `xml:"category,omitempty"`
}
