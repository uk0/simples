package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"retroblog/config"
	"retroblog/models"
	"retroblog/parser"
	"slices"
	"strings"
	"time"
)

var attachmentRefRe = regexp.MustCompile(`\[Attachment:\s*([^\[\]]+?)\s*\]`)

// HashPassword creates SHA256 hash of password
func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

// HashString creates MD5 hash for string (used for cookie names)
func HashString(s string) string {
	hash := md5.Sum([]byte(s))
	return hex.EncodeToString(hash[:])[:16]
}

// GenerateAuthToken generates authentication token
func GenerateAuthToken(path, password string) string {
	combined := path + ":" + password + ":" + time.Now().Format("2006-01-02")
	hash := sha256.Sum256([]byte(combined))
	return hex.EncodeToString(hash[:])
}

// VerifyAuthCookie verifies authentication cookie
func VerifyAuthCookie(path, token string) bool {
	config.ProtectedMu.RLock()
	_, ok := config.ProtectedPosts[path]
	config.ProtectedMu.RUnlock()

	if !ok {
		return false
	}

	return len(token) == 64
}

// CopyFile copies a file from src to dest
func CopyFile(src, dest string) error {
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

// FindAttachmentNames finds attachment names in content
func FindAttachmentNames(content string) []string {
	var names []string
	matches := attachmentRefRe.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		if len(match) > 1 {
			names = append(names, match[1])
		}
	}
	return names
}

// ExtractFirstHeading extracts first heading from markdown
func ExtractFirstHeading(md string) string {
	for _, line := range strings.Split(md, "\n") {
		l := strings.TrimSpace(line)
		if strings.HasPrefix(l, "#") {
			return strings.TrimSpace(strings.TrimLeft(l, "#"))
		}
	}
	return ""
}

// ExtractDescription extracts description from content
func ExtractDescription(body string) string {
	desc := body
	desc = regexp.MustCompile(`(?m)^#+\s+.*$`).ReplaceAllString(desc, "")
	desc = regexp.MustCompile(`\[([^\]]+)\]\([^\)]+\)`).ReplaceAllString(desc, "$1")
	desc = regexp.MustCompile(`!\[([^\]]*)\]\([^\)]+\)`).ReplaceAllString(desc, "")
	desc = regexp.MustCompile(`\*+([^\*]+)\*+`).ReplaceAllString(desc, "$1")
	desc = regexp.MustCompile(`_+([^_]+)_+`).ReplaceAllString(desc, "$1")
	desc = regexp.MustCompile("(?s)```[^`]*```").ReplaceAllString(desc, "")
	desc = regexp.MustCompile("`([^`]+)`").ReplaceAllString(desc, "$1")
	desc = regexp.MustCompile(`<[^>]+>`).ReplaceAllString(desc, "")
	desc = regexp.MustCompile(`\[Attachment:[^\]]+\]`).ReplaceAllString(desc, "")
	desc = strings.TrimSpace(desc)
	desc = regexp.MustCompile(`\s+`).ReplaceAllString(desc, " ")

	if len(desc) > 160 {
		if idx := strings.LastIndex(desc[:157], " "); idx > 100 {
			desc = desc[:idx] + "..."
		} else {
			desc = desc[:157] + "..."
		}
	}

	return desc
}

// TransformAttachmentTags transforms attachment tags to HTML
func TransformAttachmentTags(body, slug string, meta models.NoteMeta) string {
	return attachmentRefRe.ReplaceAllStringFunc(body, func(match string) string {
		m := attachmentRefRe.FindStringSubmatch(match)
		if len(m) != 2 {
			return match
		}

		originalFilename := m[1]
		att, ok := meta.Attachments[originalFilename]
		if !ok {
			return match
		}

		rel := filepath.ToSlash(filepath.Join(slug, att.Path))
		ext := strings.ToLower(filepath.Ext(att.OriginalFilename))

		switch ext {
		case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp", ".svg":
			alt := att.OriginalFilename
			if alt == "" {
				alt = att.SavedAs
			}
			uniqueID := fmt.Sprintf("pic_%x", md5.Sum([]byte(originalFilename)))[:12]
			return parser.GeneratePICWarpHTML(uniqueID, rel)

		case ".mp4", ".webm", ".ogg", ".mov":
			uniqueID := fmt.Sprintf("vid_%x", md5.Sum([]byte(originalFilename)))[:12]
			return parser.GenerateVideoWarpHTML(uniqueID, rel)

		case ".mp3", ".wav", ".flac", ".m4a", ".aac", ".oga":
			uniqueID := fmt.Sprintf("aid_%x", md5.Sum([]byte(originalFilename)))[:12]
			return parser.GenerateAudioWarpHTML(uniqueID, rel)

		case ".nes":
			uniqueID := fmt.Sprintf("nes_%x", md5.Sum([]byte(originalFilename)))[:12]
			name := att.OriginalFilename
			if name == "" {
				name = filepath.Base(att.Path)
			}
			return parser.GenerateNESPlayerHTML(uniqueID, rel, htmlEscape(name), strings.ToLower(filepath.Ext(att.Path)))

		case ".pdf":
			log.Println("Generating PDF parser originalFilename: ", originalFilename)
			log.Println("Generating PDF parser SourceFile: ", att.SourceFile)
			log.Println("Generating PDF parser SourceFile: ", att.OriginalFilename)
			log.Println("Generating PDF parser SavedAs: ", att.SavedAs)
			uniqueID := fmt.Sprintf("pdf_%x", md5.Sum([]byte(originalFilename)))[:12]
			return parser.GeneratePDFWarpHTML(uniqueID, htmlEscape(rel))

		case ".csv", ".xlsx", ".xls", ".docx":
			uniqueID := fmt.Sprintf("amis_%x", md5.Sum([]byte(originalFilename)))[:12]
			return parser.GenerateExcelWarpHTML(uniqueID, rel)

		case ".pptx":
			uniqueID := fmt.Sprintf("ppt_%x", md5.Sum([]byte(originalFilename)))[:12]
			return parser.GeneratePPTWarpHTML(uniqueID, rel)
		case ".jar":
			uniqueID := fmt.Sprintf("java_%x", md5.Sum([]byte(originalFilename)))[:12]
			name := att.OriginalFilename
			if name == "" {
				name = filepath.Base(att.Path)
			}
			// 这些是已知的Java游戏，可以直接运行
			EnableJavaGames := []string{
				"bomberman.jar", "bomberman2.jar", "bomberman3.jar",
				"tankwars.jar", "tankwars2.jar",
				"snake.jar", "snake2.jar",
				"minesweeper.jar",
				"tetris.jar", "tetris2.jar", "ForgottenWarrior.jar"}
			if slices.Contains(EnableJavaGames, name) {
				return parser.GenerateJARPlayerHTML(uniqueID, rel, htmlEscape(name), strings.ToLower(filepath.Ext(att.Path)))
			}
			//默认下载
			return parser.GenerateDownloadLinkHTML(uniqueID, rel, htmlEscape(name))

		case ".zip":
			uniqueID := fmt.Sprintf("arc_%x", md5.Sum([]byte(originalFilename)))[:12]
			name := att.OriginalFilename
			if name == "" {
				name = filepath.Base(att.Path)
			}

			if slices.Contains([]string{"kof97.zip", "pbobbl2n.zip", "kof98.zip"}, name) {
				return parser.GenerateArcadePlayerHTML(uniqueID, rel, htmlEscape(name), strings.ToLower(filepath.Ext(att.Path)))
			}
			return fmt.Sprintf(`<a href="%s" download>%s</a>`, rel, htmlEscape(name))
		case ".7z", ".chd":
			uniqueID := fmt.Sprintf("ps_%x", md5.Sum([]byte(originalFilename)))[:12]
			name := att.OriginalFilename
			if name == "" {
				name = filepath.Base(att.Path)
			}
			return parser.GeneratePlayStationPlayerHTML(uniqueID, rel, htmlEscape(name), strings.ToLower(filepath.Ext(att.Path)))
		case ".gba", ".agb":
			uniqueID := fmt.Sprintf("gba_%x", md5.Sum([]byte(originalFilename)))[:12]
			name := att.OriginalFilename
			if name == "" {
				name = filepath.Base(att.Path)
			}
			return parser.GenerateGBAPlayerHTML(uniqueID, rel, htmlEscape(name), strings.ToLower(filepath.Ext(att.Path)))

		case ".md":
			uniqueID := fmt.Sprintf("md_%x", md5.Sum([]byte(originalFilename)))[:12]
			name := att.OriginalFilename
			if name == "" {
				name = filepath.Base(att.Path)
			}
			return parser.GenerateSegaMDPlayerHTML(uniqueID, rel, htmlEscape(name), strings.ToLower(filepath.Ext(att.Path)))

		// Add other game formats...
		default:
			name := att.OriginalFilename
			if name == "" {
				name = filepath.Base(att.Path)
			}
			uniqueID := fmt.Sprintf("any_%x", md5.Sum([]byte(name)))[:12]
			return parser.GenerateDownloadLinkHTML(uniqueID, rel, htmlEscape(name))
		}
	})
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
