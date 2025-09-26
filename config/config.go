package config

import (
	"html/template"
	"log"
	"strings"
	"sync"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	ghtml "github.com/yuin/goldmark/renderer/html"
)

var (
	// Configuration
	SrcDir         string
	OutDir         string
	Addr           string
	BaseURL        string
	AttachmentsDir = "attachments"

	// Templates
	PostTpl  *template.Template
	IndexTpl *template.Template

	// Markdown converter
	MdConv goldmark.Markdown

	// Protected content management
	ProtectedMu          sync.RWMutex
	ProtectedPosts       = make(map[string]string)
	FullContent          = make(map[string]string)
	ProtectedAttachments = make(map[string]bool)

	// Rebuild state
	Rebuilding   bool
	RebuildingMu sync.Mutex
)

func Initialize() {
	BaseURL = strings.TrimRight(BaseURL, "/")
	LoadTemplates()
	MdConv = goldmark.New(
		goldmark.WithExtensions(extension.GFM, extension.Table, extension.TaskList, extension.Strikethrough, extension.CJK),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			ghtml.WithHardWraps(),
			ghtml.WithXHTML(),
			ghtml.WithUnsafe(),
		),
	)
}

func LoadTemplates() {
	var err error
	PostTpl, err = template.ParseFiles("tpl/post.html")
	if err != nil {
		log.Fatalf("post template: %v", err)
	}

	IndexTpl, err = template.ParseFiles("tpl/index.html")
	if err != nil {
		log.Fatalf("index template: %v", err)
	}
}
