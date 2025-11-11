package config

import (
	"html/template"
	"strings"
	"sync"
	"time"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	ghtml "github.com/yuin/goldmark/renderer/html"

	mermaid "go.abhg.dev/goldmark/mermaid"
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
		goldmark.WithExtensions(
			extension.GFM,
			extension.Table,
			extension.TaskList,
			extension.Strikethrough,
			extension.CJK,
			highlighting.NewHighlighting(
				highlighting.WithStyle("monokai"),
				highlighting.WithFormatOptions(
					chromahtml.WithLineNumbers(true),
				),
			),
			&mermaid.Extender{},
			extension.NewTypographer(
				extension.WithTypographicSubstitutions(extension.TypographicSubstitutions{
					extension.LeftSingleQuote:  []byte("&sbquo;"),
					extension.RightSingleQuote: nil, // nil disables a substitution
				}),
			),
		),
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

func isUpdated(created, modified time.Time) bool {
	if modified.IsZero() || created.IsZero() {
		return false
	}
	// 3 天时间窗口
	sevenDays := 3 * 24 * time.Hour
	// 如果修改时间比创建时间晚，且在 3 天之内，返回 true
	return modified.After(created) && modified.Sub(created) <= sevenDays
}

func LoadTemplates() {
	PostTpl = template.Must(template.New("post.html").Funcs(template.FuncMap{
		"isUpdated": isUpdated,
	}).ParseFiles("tpl/post.html"))

	IndexTpl = template.Must(template.New("index.html").Funcs(template.FuncMap{
		"isUpdated": isUpdated,
	}).ParseFiles("tpl/index.html"))
}
