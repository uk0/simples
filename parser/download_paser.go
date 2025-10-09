package parser

import (
	"fmt"
	"html/template"
	"path/filepath"
	"strings"
)

// GenerateDownloadLinkHTML 生成 macOS Sequoia 风格的下载链接
// 支持文件图标、大小显示、悬停效果
func GenerateDownloadLinkHTML(uniqueID, rel, fileName string) string {
	safeRel := template.HTMLEscapeString(rel)
	safeFileName := template.HTMLEscapeString(fileName)

	// 获取文件扩展名
	ext := strings.ToLower(filepath.Ext(fileName))
	fileIcon := getFileIcon(ext)
	fileType := getFileType(ext)

	return fmt.Sprintf(`<div class="download-link-wrapper" style="position:relative;max-width:100%%;margin:12px auto;display:block;">
<style>
.download-link-wrapper{position:relative;width:100%%;max-width:600px;margin:12px auto;display:block}
.download-link-wrapper .download-container{position:relative;width:100%%;background:rgba(30,30,32,0.7);backdrop-filter:blur(30px) saturate(180%%);-webkit-backdrop-filter:blur(30px) saturate(180%%);border:1px solid rgba(255,255,255,0.12);border-radius:12px;padding:14px 16px;transition:all 0.25s ease;display:flex;align-items:center;gap:12px;text-decoration:none;cursor:pointer;box-shadow:0 4px 16px rgba(0,0,0,0.15),0 0 1px rgba(255,255,255,0.08) inset}
.download-link-wrapper .download-container:hover{background:rgba(40,40,42,0.75);border-color:rgba(255,255,255,0.18);box-shadow:0 6px 24px rgba(0,0,0,0.22),0 0 1px rgba(255,255,255,0.12) inset;transform:translateY(-1px)}
.download-link-wrapper .download-container:active{transform:translateY(0);box-shadow:0 2px 8px rgba(0,0,0,0.2)}
.download-link-wrapper .download-icon{width:48px;height:48px;min-width:48px;border-radius:10px;background:linear-gradient(135deg,rgba(10,132,255,0.15) 0%%,rgba(0,122,255,0.25) 100%%);border:1px solid rgba(10,132,255,0.2);display:flex;align-items:center;justify-content:center;font-size:24px;box-shadow:0 2px 8px rgba(0,122,255,0.15)}
.download-link-wrapper .download-info{flex:1;min-width:0}
.download-link-wrapper .download-filename{font-size:13px;font-weight:600;color:rgba(255,255,255,0.98);margin-bottom:3px;white-space:nowrap;overflow:hidden;text-overflow:ellipsis;font-family:-apple-system,BlinkMacSystemFont,'SF Pro Text','Segoe UI',Roboto,sans-serif;letter-spacing:-0.01em;text-shadow:0 1px 2px rgba(0,0,0,0.2)}
.download-link-wrapper .download-type{font-size:10px;color:rgba(255,255,255,0.55);text-transform:uppercase;letter-spacing:0.8px;font-weight:500}
.download-link-wrapper .download-action{display:flex;align-items:center;gap:6px;padding:7px 12px;background:linear-gradient(135deg,#007AFF 0%%,#0051D5 100%%);border-radius:8px;font-size:12px;font-weight:600;color:#fff;white-space:nowrap;transition:all 0.2s;box-shadow:0 2px 8px rgba(0,122,255,0.3);font-family:-apple-system,BlinkMacSystemFont,'SF Pro Text',sans-serif}
.download-link-wrapper .download-container:hover .download-action{background:linear-gradient(135deg,#0A84FF 0%%,#0060EA 100%%);box-shadow:0 3px 12px rgba(0,122,255,0.4);transform:scale(1.02)}
.download-link-wrapper .download-badge{position:absolute;top:8px;right:8px;background:rgba(0,0,0,0.4);backdrop-filter:blur(10px);border:1px solid rgba(255,255,255,0.08);padding:3px 8px;border-radius:6px;font-size:9px;color:rgba(255,255,255,0.75);font-weight:600;letter-spacing:0.5px;text-transform:uppercase;opacity:0;transition:all 0.2s;font-family:-apple-system,BlinkMacSystemFont,'SF Pro Text',sans-serif}
.download-link-wrapper .download-container:hover .download-badge{opacity:1}
@media (max-width:768px){
.download-link-wrapper{margin:8px}
.download-link-wrapper .download-container{padding:12px}
.download-link-wrapper .download-icon{width:40px;height:40px;font-size:20px}
.download-link-wrapper .download-filename{font-size:12px}
.download-link-wrapper .download-action{padding:6px 10px;font-size:11px}
}
@media (prefers-color-scheme:light){
.download-link-wrapper .download-container{background:rgba(250,250,252,0.85);border-color:rgba(0,0,0,0.08)}
.download-link-wrapper .download-container:hover{background:rgba(245,245,247,0.95);border-color:rgba(0,0,0,0.12)}
.download-link-wrapper .download-filename{color:rgba(0,0,0,0.95);text-shadow:0 1px 1px rgba(255,255,255,0.5)}
.download-link-wrapper .download-type{color:rgba(0,0,0,0.5)}
.download-link-wrapper .download-icon{background:linear-gradient(135deg,rgba(0,122,255,0.08) 0%%,rgba(0,122,255,0.15) 100%%);border-color:rgba(0,122,255,0.12)}
.download-link-wrapper .download-badge{background:rgba(255,255,255,0.8);border-color:rgba(0,0,0,0.06);color:rgba(0,0,0,0.65)}
}
</style>
<a href="%s" download="%s" class="download-container" title="Download %s">
<div class="download-badge">Download</div>
<div class="download-icon">%s</div>
<div class="download-info">
<div class="download-filename">%s</div>
<div class="download-type">%s File</div>
</div>
<div class="download-action">
<span>⬇</span>
<span>Download</span>
</div>
</a>
</div>`, safeRel, safeFileName, safeFileName, fileIcon, safeFileName, fileType)
}

// getFileIcon 根据文件扩展名返回对应的图标
func getFileIcon(ext string) string {
	icons := map[string]string{
		// 文档
		".pdf":  "📄",
		".doc":  "📝",
		".docx": "📝",
		".txt":  "📃",
		".rtf":  "📃",
		".md":   "📋",

		// 表格
		".xls":  "📊",
		".xlsx": "📊",
		".csv":  "📊",

		// 演示
		".ppt":  "📊",
		".pptx": "📊",
		".key":  "📊",

		// 图片
		".jpg":  "🖼",
		".jpeg": "🖼",
		".png":  "🖼",
		".gif":  "🖼",
		".svg":  "🎨",
		".webp": "🖼",
		".ico":  "🖼",
		".bmp":  "🖼",

		// 视频
		".mp4":  "🎬",
		".mov":  "🎬",
		".avi":  "🎬",
		".mkv":  "🎬",
		".webm": "🎬",
		".flv":  "🎬",
		".wmv":  "🎬",

		// 音频
		".mp3":  "🎵",
		".wav":  "🎵",
		".flac": "🎵",
		".m4a":  "🎵",
		".aac":  "🎵",
		".ogg":  "🎵",

		// 压缩包
		".zip": "📦",
		".rar": "📦",
		".7z":  "📦",
		".tar": "📦",
		".gz":  "📦",
		".bz2": "📦",

		// 代码
		".js":    "💻",
		".jsx":   "⚛️",
		".ts":    "💻",
		".tsx":   "⚛️",
		".py":    "🐍",
		".go":    "🔷",
		".java":  "☕",
		".cpp":   "💻",
		".c":     "💻",
		".php":   "🐘",
		".rb":    "💎",
		".swift": "🦅",
		".rs":    "🦀",

		// 配置
		".json": "⚙️",
		".yaml": "⚙️",
		".yml":  "⚙️",
		".xml":  "⚙️",
		".toml": "⚙️",
		".ini":  "⚙️",

		// 其他
		".iso": "💿",
		".dmg": "💿",
		".exe": "⚙️",
		".app": "📱",
		".apk": "🤖",
	}

	if icon, ok := icons[ext]; ok {
		return icon
	}
	return "📁" // 默认文件图标
}

// getFileType 根据文件扩展名返回文件类型描述
func getFileType(ext string) string {
	types := map[string]string{
		// 文档
		".pdf":  "PDF",
		".doc":  "Word",
		".docx": "Word",
		".txt":  "Text",
		".rtf":  "RTF",
		".md":   "Markdown",

		// 表格
		".xls":  "Excel",
		".xlsx": "Excel",
		".csv":  "CSV",

		// 演示
		".ppt":  "PowerPoint",
		".pptx": "PowerPoint",
		".key":  "Keynote",

		// 图片
		".jpg":  "JPEG Image",
		".jpeg": "JPEG Image",
		".png":  "PNG Image",
		".gif":  "GIF Image",
		".svg":  "SVG Image",
		".webp": "WebP Image",
		".ico":  "Icon",
		".bmp":  "Bitmap",

		// 视频
		".mp4":  "MP4 Video",
		".mov":  "QuickTime",
		".avi":  "AVI Video",
		".mkv":  "MKV Video",
		".webm": "WebM Video",
		".flv":  "Flash Video",
		".wmv":  "WMV Video",

		// 音频
		".mp3":  "MP3 Audio",
		".wav":  "WAV Audio",
		".flac": "FLAC Audio",
		".m4a":  "M4A Audio",
		".aac":  "AAC Audio",
		".ogg":  "OGG Audio",

		// 压缩包
		".zip": "ZIP Archive",
		".rar": "RAR Archive",
		".7z":  "7-Zip Archive",
		".tar": "TAR Archive",
		".gz":  "GZip Archive",
		".bz2": "BZ2 Archive",

		// 代码
		".js":    "JavaScript",
		".jsx":   "React JSX",
		".ts":    "TypeScript",
		".tsx":   "React TSX",
		".py":    "Python",
		".go":    "Go",
		".java":  "Java",
		".cpp":   "C++",
		".c":     "C",
		".php":   "PHP",
		".rb":    "Ruby",
		".swift": "Swift",
		".rs":    "Rust",

		// 配置
		".json": "JSON",
		".yaml": "YAML",
		".yml":  "YAML",
		".xml":  "XML",
		".toml": "TOML",
		".ini":  "Config",

		// 其他
		".iso": "Disk Image",
		".dmg": "macOS Image",
		".exe": "Executable",
		".app": "Application",
		".apk": "Android App",
	}

	if fileType, ok := types[ext]; ok {
		return fileType
	}

	// 返回大写的扩展名（去掉点）
	if len(ext) > 1 {
		return strings.ToUpper(ext[1:])
	}
	return "File"
}
