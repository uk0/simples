package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// OpenAI API 配置
const (
	OpenAIAPIURL = "https://open.bigmodel.cn/api/paas/v4/chat/completions"
	DefaultModel = "glm-4.5-flash"
)

// OpenAI 请求结构
type OpenAIRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenAI 响应结构
type OpenAIResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error,omitempty"`
}

// OG Meta 信息结构
type OGMetaInfo struct {
	Title       string `json:"og:title"`
	Description string `json:"og:description"`
	Type        string `json:"og:type"`
	URL         string `json:"og:url"`
	Image       string `json:"og:image"`
	SiteName    string `json:"og:site_name"`
	Locale      string `json:"og:locale"`
}

// HTMLExtractResult HTML 提取结果
type HTMLExtractResult struct {
	TextContent string
	Images      []string
	BaseURL     string
}

// OGGenerator 主结构
type OGGenerator struct {
	APIKey     string
	HTTPClient *http.Client
	Model      string
}

// NewOGGenerator 创建新的生成器实例
func NewOGGenerator(apiKey string) *OGGenerator {
	return &OGGenerator{
		APIKey: apiKey,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		Model: DefaultModel,
	}
}

// FetchHTMLFromURL 从 URL 获取 HTML 内容
func (og *OGGenerator) FetchHTMLFromURL(urlStr string) (string, error) {
	resp, err := og.HTTPClient.Get(urlStr)
	if err != nil {
		return "", fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP error: status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	return string(body), nil
}

// normalizeImageURL 标准化图片 URL
func (og *OGGenerator) normalizeImageURL(imgSrc, baseURL string) string {
	imgSrc = strings.TrimSpace(imgSrc)
	if imgSrc == "" {
		return ""
	}

	// 如果已经是完整的 URL
	if strings.HasPrefix(imgSrc, "http://") || strings.HasPrefix(imgSrc, "https://") {
		return imgSrc
	}

	// 解析 base URL
	base, err := url.Parse(baseURL)
	if err != nil {
		return imgSrc
	}

	// 解析图片 URL
	imgURL, err := url.Parse(imgSrc)
	if err != nil {
		return imgSrc
	}

	// 合并 URL
	fullURL := base.ResolveReference(imgURL)
	return fullURL.String()
}

// isValidImageURL 判断是否是有效的图片 URL
func (og *OGGenerator) isValidImageURL(imgSrc string) bool {
	imgSrc = strings.ToLower(imgSrc)

	// 排除一些常见的无效图片
	if strings.Contains(imgSrc, "data:image") {
		return false
	}
	if strings.Contains(imgSrc, "blank.") {
		return false
	}
	if strings.Contains(imgSrc, "spacer.") {
		return false
	}
	if strings.Contains(imgSrc, "pixel.") {
		return false
	}

	// 检查是否是图片扩展名
	validExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp", ".svg", ".bmp"}
	for _, ext := range validExts {
		if strings.HasSuffix(imgSrc, ext) {
			return true
		}
	}

	// 如果 URL 中包含常见的图片路径关键词
	imageKeywords := []string{".png", ".jpg", ".webp", ".jpeg", ".gif"}
	for _, keyword := range imageKeywords {
		if strings.Contains(imgSrc, keyword) {
			return true
		}
	}

	return false
}

// ExtractTextFromHTML 从 HTML 中提取文本内容和图片
func (og *OGGenerator) ExtractTextFromHTML(htmlContent string, baseURL string) (*HTMLExtractResult, error) {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var textBuilder strings.Builder
	var title string
	var metaDesc string
	var images []string
	imageSet := make(map[string]bool) // 用于去重

	var extractContent func(*html.Node)
	extractContent = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "script", "style", "noscript", "iframe":
				return
			case "title":
				if n.FirstChild != nil {
					title = n.FirstChild.Data
				}
			case "meta":
				var name, content, property string
				for _, attr := range n.Attr {
					if attr.Key == "name" {
						name = attr.Val
					}
					if attr.Key == "property" {
						property = attr.Val
					}
					if attr.Key == "content" {
						content = attr.Val
					}
				}
				// 提取 meta 描述
				if (name == "description" || name == "og:description") && content != "" {
					metaDesc = content
				}
				// 提取 og:image
				if property == "og:image" && content != "" {
					normalizedURL := og.normalizeImageURL(content, baseURL)
					if normalizedURL != "" && !imageSet[normalizedURL] {
						images = append(images, normalizedURL)
						imageSet[normalizedURL] = true
					}
				}
			case "img":
				// 提取 img 标签的 src
				var src, alt string
				for _, attr := range n.Attr {
					if attr.Key == "src" {
						src = attr.Val
					}
					if attr.Key == "alt" {
						alt = attr.Val
					}
				}
				if src != "" && og.isValidImageURL(src) {
					normalizedURL := og.normalizeImageURL(src, baseURL)
					if normalizedURL != "" && !imageSet[normalizedURL] {
						images = append(images, normalizedURL)
						imageSet[normalizedURL] = true
					}
				}
				// 将 alt 文本也加入内容
				if alt != "" {
					textBuilder.WriteString(alt)
					textBuilder.WriteString(" ")
				}
			case "link":
				// 提取 link rel="image_src"
				var rel, href string
				for _, attr := range n.Attr {
					if attr.Key == "rel" {
						rel = attr.Val
					}
					if attr.Key == "href" {
						href = attr.Val
					}
				}
				if rel == "image_src" && href != "" {
					normalizedURL := og.normalizeImageURL(href, baseURL)
					if normalizedURL != "" && !imageSet[normalizedURL] {
						images = append(images, normalizedURL)
						imageSet[normalizedURL] = true
					}
				}
			}
		}
		if n.Type == html.TextNode {
			text := strings.TrimSpace(n.Data)
			if text != "" {
				textBuilder.WriteString(text)
				textBuilder.WriteString(" ")
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractContent(c)
		}
	}

	extractContent(doc)

	// 构建最终文本
	result := ""
	if title != "" {
		result += "Title: " + title + "\n\n"
	}
	if metaDesc != "" {
		result += "Meta Description: " + metaDesc + "\n\n"
	}

	bodyText := textBuilder.String()
	// 限制文本长度，避免 token 超限
	if len(bodyText) > 3000 {
		bodyText = bodyText[:3000] + "..."
	}
	result += "Content: " + bodyText

	return &HTMLExtractResult{
		TextContent: result,
		Images:      images,
		BaseURL:     baseURL,
	}, nil
}

// CallOpenAI 调用 OpenAI API 生成 OG 信息
func (og *OGGenerator) CallOpenAI(content string, sourceURL string, extractedImage string) (*OGMetaInfo, error) {
	prompt := og.buildPrompt(content, sourceURL, extractedImage)

	reqBody := OpenAIRequest{
		Model: og.Model,
		Messages: []Message{
			{
				Role:    "system",
				Content: "你是一个专业的 SEO 和 Open Graph meta 标签专家。你的任务是分析网页内容并生成完整、准确、吸引人的 Open Graph meta 标签信息。",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		MaxTokens:   800,
		Temperature: 0.1,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", OpenAIAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+og.APIKey)

	resp, err := og.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call OpenAI API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	fmt.Println("------------------------------------")
	fmt.Println(string(body))
	fmt.Println("------------------------------------")
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var openAIResp OpenAIResponse
	if err := json.Unmarshal(body, &openAIResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if openAIResp.Error != nil {
		return nil, fmt.Errorf("OpenAI API error: %s", openAIResp.Error.Message)
	}

	if len(openAIResp.Choices) == 0 {
		return nil, errors.New("no response from OpenAI")
	}

	// 解析 OpenAI 返回的 JSON
	ogInfo, err := og.parseOGInfo(openAIResp.Choices[0].Message.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse OG info: %w", err)
	}

	return ogInfo, nil
}

// buildPrompt 构建 OpenAI 提示词
func (og *OGGenerator) buildPrompt(content string, sourceURL string, extractedImage string) string {
	urlInfo := ""
	if sourceURL != "" {
		urlInfo = fmt.Sprintf("\n\n原始URL: %s", sourceURL)
	}

	imageInfo := ""
	if extractedImage != "" {
		imageInfo = fmt.Sprintf("\n\n已提取的图片URL（请直接使用此URL作为 og:image）: %s", extractedImage)
	} else {
		imageInfo = "\n\n注意：页面中未找到合适的图片，请使用 \"logo.png\" 作为 og:image"
	}

	prma := "请根据以下网页内容，生成完整的 Open Graph (OG) meta 标签信息。%s%s\n\t\t\t网页内容:\n\t\t\t%s\n\t\t请以 JSON 格式返回以下字段（仅返回 JSON，不要其他说明文字，直接返回内容即可,不要```json包裹内容）Response Format:\n\t\t\t\t\t\t\t\t\n\t\t\t\t\t\t\t\t{\n\t\t\t\t\t\t\t\t  \"og:title\": \"网页标题（简洁有力，50-60字符）\",\n\t\t\t\t\t\t\t\t  \"og:description\": \"网页描述（吸引人的摘要，150-160字符）\",\n\t\t\t\t\t\t\t\t  \"og:type\": \"内容类型（如: website, article, product 等）\",\n\t\t\t\t\t\t\t\t  \"og:url\": \"规范化的URL地址\",\n\t\t\t\t\t\t\t\t  \"og:image\": \"代表性图片URL\",\n\t\t\t\t\t\t\t\t  \"og:site_name\": \"网站名称\",\n\t\t\t\t\t\t\t\t  \"og:locale\": \"语言区域（如: zh_CN, en_US）\"\n\t\t\t\t\t\t\t\t}\n\t\t\t\t\t\t\t\t\n\t\t\t\t\t\t\t\t要求：\n\t\t\t\t\t\t\t\t1. 标题要简洁明了，突出核心内容\n\t\t\t\t\t\t\t\t2. 描述要有吸引力，能激发用户点击欲望\n\t\t\t\t\t\t\t\t3. 根据内容智能判断类型（article/website/product等）\n\t\t\t\t\t\t\t\t4. 自动检测内容语言并设置正确的 locale\n\t\t\t\t\t\t\t\t5. 如果是文章类型，描述应该是内容摘要\n\t\t\t\t\t\t\t\t6. 所有字段都必须填写，不能为空\n\t\t\t\t\t\t\t\t7. 如果提供了已提取的图片URL，必须直接使用该URL作为 og:image 的值\n\t\t\t\t\t\t\t\t8. 如果未提供图片URL，使用 \"logo.png\" 作为 og:image 的值"

	return fmt.Sprintf(prma, urlInfo, imageInfo, content)
}

// parseOGInfo 解析 OpenAI 返回的 OG 信息
func (og *OGGenerator) parseOGInfo(content string) (*OGMetaInfo, error) {
	// 提取 JSON 部分（可能被包裹在 markdown 代码块中）
	jsonPattern := regexp.MustCompile("(?s)```(?:json)?\\s*({.*?})\\s*```")
	matches := jsonPattern.FindStringSubmatch(content)

	var jsonStr string
	if len(matches) > 1 {
		jsonStr = matches[1]
	} else {
		// 尝试直接解析
		jsonStr = strings.TrimSpace(content)
	}

	var ogInfo OGMetaInfo
	if err := json.Unmarshal([]byte(jsonStr), &ogInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal OG info: %w, content: %s", err, jsonStr)
	}

	return &ogInfo, nil
}

// GenerateOGMetaHTML 生成 OG Meta HTML 标签
func (og *OGGenerator) GenerateOGMetaHTML(ogInfo *OGMetaInfo) string {
	var builder strings.Builder

	builder.WriteString("<!-- Open Graph Meta Tags -->\n")

	if ogInfo.Title != "" {
		builder.WriteString(fmt.Sprintf("<meta property=\"og:title\" content=\"%s\" />\n",
			html.EscapeString(ogInfo.Title)))
	}

	if ogInfo.Description != "" {
		builder.WriteString(fmt.Sprintf("<meta property=\"og:description\" content=\"%s\" />\n",
			html.EscapeString(ogInfo.Description)))
	}

	if ogInfo.Type != "" {
		builder.WriteString(fmt.Sprintf("<meta property=\"og:type\" content=\"%s\" />\n",
			html.EscapeString(ogInfo.Type)))
	}

	if ogInfo.URL != "" {
		builder.WriteString(fmt.Sprintf("<meta property=\"og:url\" content=\"%s\" />\n",
			html.EscapeString(ogInfo.URL)))
	}

	if ogInfo.Image != "" {
		builder.WriteString(fmt.Sprintf("<meta property=\"og:image\" content=\"%s\" />\n",
			html.EscapeString(ogInfo.Image)))
	}

	if ogInfo.SiteName != "" {
		builder.WriteString(fmt.Sprintf("<meta property=\"og:site_name\" content=\"%s\" />\n",
			html.EscapeString(ogInfo.SiteName)))
	}

	if ogInfo.Locale != "" {
		builder.WriteString(fmt.Sprintf("<meta property=\"og:locale\" content=\"%s\" />\n",
			html.EscapeString(ogInfo.Locale)))
	}

	return builder.String()
}

// GenerateFromURL 从 URL 生成 OG Meta 信息
func (og *OGGenerator) GenerateFromURL(urlStr string) (string, *OGMetaInfo, error) {
	// 1. 获取 HTML 内容
	htmlContent, err := og.FetchHTMLFromURL(urlStr)
	if err != nil {
		return "", nil, err
	}

	// 2. 提取文本和图片
	extractResult, err := og.ExtractTextFromHTML(htmlContent, urlStr)
	if err != nil {
		return "", nil, err
	}

	// 3. 选择图片：优先使用第一张提取到的图片，否则使用 logo.png
	selectedImage := "logo.png"
	if len(extractResult.Images) > 0 {
		selectedImage = extractResult.Images[0]
		fmt.Printf("✓ 从页面中提取到 %d 张图片，使用第一张: %s\n", len(extractResult.Images), selectedImage)
	} else {
		fmt.Println("✗ 页面中未找到合适的图片，使用默认: logo.png")
	}

	// 4. 调用 OpenAI 生成 OG 信息
	ogInfo, err := og.CallOpenAI(extractResult.TextContent, urlStr, selectedImage)
	if err != nil {
		return "", nil, err
	}

	// 5. 确保使用选定的图片
	ogInfo.Image = selectedImage

	// 6. 如果没有设置 URL，使用源 URL
	if ogInfo.URL == "" {
		ogInfo.URL = urlStr
	}

	// 7. 生成 HTML 标签
	metaTags := og.GenerateOGMetaHTML(ogInfo)

	return metaTags, ogInfo, nil
}

// GenerateFromHTML 从 HTML 内容生成 OG Meta 信息
func (og *OGGenerator) GenerateFromHTML(htmlContent string, sourceURL string) (string, *OGMetaInfo, error) {
	// 1. 提取文本和图片
	extractResult, err := og.ExtractTextFromHTML(htmlContent, sourceURL)
	if err != nil {
		return "", nil, err
	}

	// 2. 选择图片：优先使用第一张提取到的图片，否则使用 logo.png
	selectedImage := "logo.png"
	if len(extractResult.Images) > 0 {
		selectedImage = extractResult.Images[0]
		fmt.Printf("✓ 从 HTML 中提取到 %d 张图片，使用第一张: %s\n", len(extractResult.Images), selectedImage)
	} else {
		fmt.Println("✗ HTML 中未找到合适的图片，使用默认: logo.png")
	}

	// 3. 调用 OpenAI 生成 OG 信息
	ogInfo, err := og.CallOpenAI(extractResult.TextContent, sourceURL, selectedImage)
	if err != nil {
		return "", nil, err
	}

	// 4. 确保使用选定的图片
	ogInfo.Image = selectedImage

	// 5. 生成 HTML 标签
	metaTags := og.GenerateOGMetaHTML(ogInfo)

	return metaTags, ogInfo, nil
}

// PrintOGInfo 打印 OG 信息（格式化输出）
func PrintOGInfo(ogInfo *OGMetaInfo) {
	fmt.Println("\n=== Open Graph Meta 信息 ===")
	fmt.Printf("标题 (og:title): %s\n", ogInfo.Title)
	fmt.Printf("描述 (og:description): %s\n", ogInfo.Description)
	fmt.Printf("类型 (og:type): %s\n", ogInfo.Type)
	fmt.Printf("URL (og:url): %s\n", ogInfo.URL)
	fmt.Printf("图片 (og:image): %s\n", ogInfo.Image)
	fmt.Printf("站点名称 (og:site_name): %s\n", ogInfo.SiteName)
	fmt.Printf("语言区域 (og:locale): %s\n", ogInfo.Locale)
	fmt.Println("============================\n")
}
