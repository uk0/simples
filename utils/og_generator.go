package utils

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode"

	"golang.org/x/net/html"
)

/* =========================
   数据结构
========================= */

type OGMetaInfo struct {
	Title       string `json:"og:title"`
	Description string `json:"og:description"`
	Type        string `json:"og:type"`
	URL         string `json:"og:url"`
	Image       string `json:"og:image"`
	SiteName    string `json:"og:site_name"`
	Locale      string `json:"og:locale"`
}

type HTMLExtractResult struct {
	BaseURL     string
	Title       string
	MetaDesc    string
	TextContent string // 纯文本（截断前）
	Lang        string // <html lang=""> 或 og:locale 推断
	Canonical   string
	SiteName    string // og:site_name
	OGImage     string // og:image
	Images      []string
}

/* =========================
   入口与 HTTP
========================= */

type OGGenerator struct {
	HTTPClient *http.Client
}

func NewOGGenerator() *OGGenerator {
	return &OGGenerator{
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (og *OGGenerator) FetchHTMLFromURL(urlStr string) (string, error) {
	resp, err := og.HTTPClient.Get(urlStr)
	if err != nil {
		return "", fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP error: status %d", resp.StatusCode)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read body: %w", err)
	}
	return string(b), nil
}

/* =========================
   解析 HTML
========================= */

func (og *OGGenerator) normalizeImageURL(imgSrc, baseURL string) string {
	imgSrc = strings.TrimSpace(imgSrc)
	if imgSrc == "" {
		return ""
	}
	if strings.HasPrefix(imgSrc, "http://") || strings.HasPrefix(imgSrc, "https://") {
		return imgSrc
	}
	base, err := url.Parse(baseURL)
	if err != nil {
		return imgSrc
	}
	u, err := url.Parse(imgSrc)
	if err != nil {
		return imgSrc
	}
	return base.ResolveReference(u).String()
}

func isLikelyImageURL(u string) bool {
	low := strings.ToLower(u)
	if strings.HasPrefix(low, "data:image") {
		return false
	}
	// 常见无效
	if strings.Contains(low, "blank.") || strings.Contains(low, "spacer.") || strings.Contains(low, "pixel.") {
		return false
	}
	exts := []string{".jpg", ".jpeg", ".png", ".webp", ".gif", ".svg", ".bmp"}
	for _, e := range exts {
		if strings.HasSuffix(low, e) {
			return true
		}
	}
	// 参数里包含图片关键词也放行
	for _, kw := range []string{".jpg", ".jpeg", ".png", ".webp", ".gif"} {
		if strings.Contains(low, kw) {
			return true
		}
	}
	return false
}

func (og *OGGenerator) ExtractTextFromHTML(htmlContent, baseURL string) (*HTMLExtractResult, error) {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return nil, fmt.Errorf("parse HTML: %w", err)
	}

	var (
		title, metaDesc, lang, canonical, siteName, ogImage string
		images                                              []string
		imageSet                                            = map[string]bool{}
		textBuilder                                         strings.Builder
	)

	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "html":
				for _, a := range n.Attr {
					if a.Key == "lang" && lang == "" {
						lang = strings.TrimSpace(a.Val)
					}
				}
			case "head":
				// 继续
			case "title":
				if n.FirstChild != nil && title == "" {
					title = strings.TrimSpace(n.FirstChild.Data)
				}
			case "meta":
				var name, property, content charsetKV
				for _, a := range n.Attr {
					switch a.Key {
					case "name":
						name.Val = a.Val
					case "property":
						property.Val = a.Val
					case "content":
						content.Val = a.Val
					}
				}
				if content.Val == "" {
					break
				}
				switch strings.ToLower(name.Val) {
				case "description":
					if metaDesc == "" {
						metaDesc = strings.TrimSpace(content.Val)
					}
				case "og:description":
					if metaDesc == "" {
						metaDesc = strings.TrimSpace(content.Val)
					}
				}
				switch strings.ToLower(property.Val) {
				case "og:locale":
					if lang == "" {
						lang = strings.TrimSpace(content.Val)
					}
				case "og:sitename", "og:site_name":
					if siteName == "" {
						siteName = strings.TrimSpace(content.Val)
					}
				case "og:image":
					if ogImage == "" {
						ogImage = og.normalizeImageURL(strings.TrimSpace(content.Val), baseURL)
					}
				case "og:title":
					// 如需覆盖 title 可在此处理
				}
			case "link":
				var rel, href string
				for _, a := range n.Attr {
					if a.Key == "rel" {
						rel = a.Val
					}
					if a.Key == "href" {
						href = a.Val
					}
				}
				if strings.EqualFold(rel, "canonical") && canonical == "" {
					canonical = og.makeAbsolute(href, baseURL)
				}
				if strings.EqualFold(rel, "image_src") && href != "" {
					abs := og.normalizeImageURL(href, baseURL)
					if abs != "" && !imageSet[abs] && isLikelyImageURL(abs) {
						images = append(images, abs)
						imageSet[abs] = true
					}
				}
			case "img":
				var src, alt string
				for _, a := range n.Attr {
					if a.Key == "src" {
						src = a.Val
					}
					if a.Key == "alt" {
						alt = a.Val
					}
				}
				if alt != "" {
					textBuilder.WriteString(strings.TrimSpace(alt))
					textBuilder.WriteByte(' ')
				}
				if src != "" {
					abs := og.normalizeImageURL(src, baseURL)
					if abs != "" && !imageSet[abs] && isLikelyImageURL(abs) {
						images = append(images, abs)
						imageSet[abs] = true
					}
				}
			case "script", "style", "noscript", "iframe":
				// 跳过
				return
			}
		}
		if n.Type == html.TextNode {
			t := strings.TrimSpace(n.Data)
			if t != "" {
				textBuilder.WriteString(t)
				textBuilder.WriteByte(' ')
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}

	walk(doc)

	return &HTMLExtractResult{
		BaseURL:     baseURL,
		Title:       title,
		MetaDesc:    metaDesc,
		TextContent: strings.TrimSpace(textBuilder.String()),
		Lang:        lang,
		Canonical:   canonical,
		SiteName:    siteName,
		OGImage:     ogImage,
		Images:      images,
	}, nil
}

type charsetKV struct{ Val string }

func (og *OGGenerator) makeAbsolute(href, base string) string {
	u, err := url.Parse(href)
	if err != nil {
		return href
	}
	if u.IsAbs() {
		return u.String()
	}
	b, err := url.Parse(base)
	if err != nil {
		return href
	}
	return b.ResolveReference(u).String()
}

/* =========================
   规则/启发式
========================= */

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		v = strings.TrimSpace(v)
		if v != "" {
			return v
		}
	}
	return ""
}

// 描述：优先 meta description；否则从正文切 150–160 字，尽量整句
func buildDescription(metaDesc, plain string) string {
	if metaDesc = strings.TrimSpace(metaDesc); metaDesc != "" {
		return clipTo(metaDesc, 160)
	}
	// 从正文取前几句
	sents := splitSentences(plain)
	var buf strings.Builder
	for _, s := range sents {
		if s == "" {
			continue
		}
		if buf.Len()+len([]rune(s)) > 160 {
			break
		}
		if buf.Len() > 0 {
			buf.WriteByte(' ')
		}
		buf.WriteString(s)
		if buf.Len() >= 120 {
			break
		}
	}
	desc := buf.String()
	if desc == "" {
		desc = plain
	}
	return clipTo(desc, 160)
}

func clipTo(s string, max int) string {
	rs := []rune(strings.TrimSpace(s))
	if len(rs) <= max {
		return string(rs)
	}
	return string(rs[:max-1]) + "…"
}

func splitSentences(s string) []string {
	// 简易句子分割：中文句末 + 英文句末
	re := regexp.MustCompile(`(?m)(.*?[\.\!\?。！？]|.+$)`)
	out := re.FindAllString(strings.TrimSpace(s), -1)
	for i := range out {
		out[i] = strings.TrimSpace(out[i])
	}
	return out
}

// 类型：路径含 /blog/ 或正文含“发表于/作者/评论”等 → article，否则 website
func guessType(pageURL, text string) string {
	u, _ := url.Parse(pageURL)
	path := strings.ToLower(u.Path)
	if strings.Contains(path, "/blog/") || strings.Contains(path, "/post/") || strings.Contains(path, "/article") {
		return "article"
	}
	low := strings.ToLower(text)
	if strings.Contains(low, "发表于") || strings.Contains(low, "作者") || strings.Contains(low, "评论") || strings.Contains(low, "阅读") {
		return "article"
	}
	return "website"
}

// locale：优先 <html lang> / og:locale；否则按 CJK 占比判定 zh_CN / en_US
func guessLocale(metaLang string, text string) string {
	lang := strings.ToLower(strings.TrimSpace(metaLang))
	if lang != "" {
		// 规范化常见值
		if lang == "zh" || strings.HasPrefix(lang, "zh-") {
			return "zh_CN"
		}
		if lang == "en" || strings.HasPrefix(lang, "en-") {
			return "en_US"
		}
		return lang
	}
	var total, cjk int
	for _, r := range text {
		if unicode.IsSpace(r) {
			continue
		}
		total++
		if isCJK(r) {
			cjk++
		}
	}
	if total > 0 && float64(cjk)/float64(total) > 0.2 {
		return "zh_CN"
	}
	return "en_US"
}

func isCJK(r rune) bool {
	return unicode.Is(unicode.Han, r) ||
		(r >= 0x3040 && r <= 0x309F) || // Hiragana
		(r >= 0x30A0 && r <= 0x30FF) || // Katakana
		(r >= 0xAC00 && r <= 0xD7AF) // Hangul
}

// 站点名：优先 og:site_name；否则使用 host 的主域（去 www），TitleCase
func guessSiteName(siteName, pageURL string) string {
	if v := strings.TrimSpace(siteName); v != "" {
		return v
	}
	u, err := url.Parse(pageURL)
	if err != nil || u.Host == "" {
		return ""
	}
	host := strings.ToLower(u.Host)
	host = strings.TrimPrefix(host, "www.")
	parts := strings.Split(host, ".")
	if len(parts) >= 2 {
		base := parts[len(parts)-2]
		return strings.Title(base)
	}
	return strings.Title(host)
}

// 图片选择：优先 og:image；否则根据关键词打分（cover/hero/share/banner/og）优先，排除 logo/icon/avatar/sprite
func pickBestImage(ogImage string, imgs []string) string {
	if ogImage != "" {
		return ogImage
	}
	score := func(u string) int {
		l := strings.ToLower(u)
		if !isLikelyImageURL(l) {
			return -1000
		}
		s := 0
		// 正向关键词
		for _, kw := range []string{"og", "cover", "hero", "share", "banner", "card"} {
			if strings.Contains(l, kw) {
				s += 5
			}
		}
		// 负向关键词
		for _, kw := range []string{"logo", "icon", "avatar", "sprite"} {
			if strings.Contains(l, kw) {
				s -= 5
			}
		}
		// 大图猜测
		for _, kw := range []string{"1200", "1024", "800", "640"} {
			if strings.Contains(l, kw) {
				s += 2
			}
		}
		return s
	}
	best, bestScore := "", -10000
	for _, u := range imgs {
		if v := score(u); v > bestScore {
			bestScore, best = v, u
		}
	}
	return best
}

/* =========================
   生成 OG
========================= */

func (og *OGGenerator) BuildOG(ex *HTMLExtractResult) *OGMetaInfo {
	// Title
	title := firstNonEmpty(ex.Title, fallbackTitleFromURL(ex.BaseURL), "Untitled")

	// URL / Canonical
	finalURL := firstNonEmpty(ex.Canonical, ex.BaseURL)

	// Description
	desc := buildDescription(ex.MetaDesc, ex.TextContent)

	// Type
	typ := guessType(finalURL, ex.TextContent)

	// Image
	image := pickBestImage(ex.OGImage, ex.Images)

	// SiteName
	site := guessSiteName(ex.SiteName, finalURL)

	// Locale
	loc := guessLocale(ex.Lang, ex.TextContent)

	return &OGMetaInfo{
		Title:       title,
		Description: desc,
		Type:        typ,
		URL:         finalURL,
		Image:       image,
		SiteName:    site,
		Locale:      loc,
	}
}

func fallbackTitleFromURL(pageURL string) string {
	u, err := url.Parse(pageURL)
	if err != nil || u.Path == "" {
		return pageURL
	}
	p := strings.Trim(u.Path, "/")
	if p == "" {
		return u.Host
	}
	seg := p
	if i := strings.LastIndex(p, "/"); i >= 0 && i+1 < len(p) {
		seg = p[i+1:]
	}
	seg = strings.TrimSuffix(seg, ".html")
	seg = strings.ReplaceAll(seg, "-", " ")
	seg = strings.ReplaceAll(seg, "_", " ")
	return strings.Title(seg)
}

/* =========================
   对外方法
========================= */

func (og *OGGenerator) GenerateFromURL(urlStr string) (string, *OGMetaInfo, error) {
	htmlContent, err := og.FetchHTMLFromURL(urlStr)
	if err != nil {
		return "", nil, err
	}
	return og.GenerateFromHTML(htmlContent, urlStr)
}

func (og *OGGenerator) GenerateFromHTML(htmlContent, sourceURL string) (string, *OGMetaInfo, error) {
	ex, err := og.ExtractTextFromHTML(htmlContent, sourceURL)
	if err != nil {
		return "", nil, err
	}
	// 文本过长剪裁，避免后续处理超大
	if len([]rune(ex.TextContent)) > 6000 {
		ex.TextContent = string([]rune(ex.TextContent)[:6000])
	}
	ogInfo := og.BuildOG(ex)
	meta := og.GenerateOGMetaHTML(ogInfo)
	return meta, ogInfo, nil
}

/* =========================
   HTML 片段输出
========================= */

func (og *OGGenerator) GenerateOGMetaHTML(ogInfo *OGMetaInfo) string {
	var b strings.Builder
	b.WriteString("<!-- Open Graph Meta Tags (heuristic) -->\n")
	if ogInfo.Title != "" {
		b.WriteString(fmt.Sprintf("<meta property=\"og:title\" content=\"%s\" />\n", htmlEscape(ogInfo.Title)))
	}
	if ogInfo.Description != "" {
		b.WriteString(fmt.Sprintf("<meta property=\"og:description\" content=\"%s\" />\n", htmlEscape(ogInfo.Description)))
	}
	if ogInfo.Type != "" {
		b.WriteString(fmt.Sprintf("<meta property=\"og:type\" content=\"%s\" />\n", htmlEscape(ogInfo.Type)))
	}
	if ogInfo.URL != "" {
		b.WriteString(fmt.Sprintf("<meta property=\"og:url\" content=\"%s\" />\n", htmlEscape(ogInfo.URL)))
	}
	if ogInfo.Image != "" {
		b.WriteString(fmt.Sprintf("<meta property=\"og:image\" content=\"%s\" />\n", htmlEscape(ogInfo.Image)))
	}
	if ogInfo.SiteName != "" {
		b.WriteString(fmt.Sprintf("<meta property=\"og:site_name\" content=\"%s\" />\n", htmlEscape(ogInfo.SiteName)))
	}
	if ogInfo.Locale != "" {
		b.WriteString(fmt.Sprintf("<meta property=\"og:locale\" content=\"%s\" />\n", htmlEscape(ogInfo.Locale)))
	}
	return b.String()
}
