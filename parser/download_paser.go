package parser

import (
	"fmt"
	"html/template"
	"path/filepath"
	"strings"
)

// GenerateDownloadLinkHTML 渲染单个下载卡片（MacOS 26 风格，轻量美化）
func GenerateDownloadLinkHTML(uniqueID, rel, fileName string) string {
	safeRel := template.HTMLEscapeString(rel)
	safeFileName := template.HTMLEscapeString(fileName)

	ext := strings.ToLower(filepath.Ext(fileName))
	fileIcon := getFileIcon(ext)
	fileType := getFileType(ext)

	return fmt.Sprintf(`
<div class="download-link-wrapper" style="position:relative;max-width:100%%;margin:10px auto;">
  <style>
    /* ========== 卡片容器 ========== */
    .dl-card {
      position: relative;
      display: flex;
      align-items: center;
      gap: 12px;
      width: 88%%;
      max-width: 960px;
      margin: 0 auto;
      padding: 12px 16px;
      background: rgba(30, 30, 32, 0.70);
      backdrop-filter: blur(26px) saturate(180%%);
      -webkit-backdrop-filter: blur(26px) saturate(180%%);
      border: 1px solid rgba(255, 255, 255, 0.12);
      border-radius: 12px;
      text-decoration: none;
      box-shadow: 0 4px 14px rgba(0, 0, 0, 0.18), inset 0 0 1px rgba(255, 255, 255, 0.08);
      transition: transform 0.15s ease, box-shadow 0.2s ease, border-color 0.2s ease, background 0.2s ease;
    }
    
    .dl-card:hover {
      transform: translateY(-1px);
      background: rgba(40, 40, 42, 0.78);
      border-color: rgba(255, 255, 255, 0.18);
      box-shadow: 0 8px 22px rgba(0, 0, 0, 0.22), inset 0 0 1px rgba(255, 255, 255, 0.12);
    }
    
    .dl-card:active {
      transform: translateY(0);
    }

    /* ========== 文件图标 ========== */
    .dl-icon {
      width: 40px;
      height: 40px;
      min-width: 40px;
      border-radius: 10px;
      display: flex;
      align-items: center;
      justify-content: center;
      background: linear-gradient(135deg, rgba(10, 132, 255, 0.14), rgba(0, 122, 255, 0.22));
      border: 1px solid rgba(10, 132, 255, 0.18);
      box-shadow: 0 2px 6px rgba(0, 122, 255, 0.18);
      font-size: 20px;
    }

    /* ========== 文件信息区域 ========== */
    .dl-info {
      min-width: 0;
      flex: 1;
      display: flex;
      flex-direction: column;
      gap: 2px;
      overflow: hidden;
    }
    
    .dl-name {
      font: 600 13px/1.25 -apple-system, BlinkMacSystemFont, "SF Pro Text", "Segoe UI", Roboto, sans-serif;
      color: rgba(255, 255, 255, 0.98);
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
      letter-spacing: -0.01em;
      text-shadow: 0 1px 2px rgba(0, 0, 0, 0.2);
    }
    
    .dl-type {
      font: 500 10px/1 -apple-system, BlinkMacSystemFont, "SF Pro Text", sans-serif;
      color: rgba(255, 255, 255, 0.58);
      letter-spacing: 0.7px;
      text-transform: uppercase;
    }

    /* ========== 下载按钮 (macOS 26 毛玻璃风格) ========== */
    .dl-btn {
      /* 布局 */
      flex-shrink: 0;
      display: inline-flex;
      align-items: center;
      justify-content: center;
      gap: 8px;
      margin-left: auto;
      padding: 8px 16px;
      border-radius: 11px;
      
      /* 毛玻璃背景 */
      background: linear-gradient(180deg, rgba(10, 132, 255, 0.95) 0%%, rgba(0, 106, 220, 0.92) 100%%);
      backdrop-filter: blur(20px) saturate(180%%);
      -webkit-backdrop-filter: blur(20px) saturate(180%%);
      
      /* 多层边框效果 */
      border: 1.5px solid rgba(255, 255, 255, 0.28);
      
      /* 多层阴影系统 */
      box-shadow: 
        0 8px 20px rgba(10, 132, 255, 0.35),
        0 2px 8px rgba(0, 0, 0, 0.24),
        inset 0 1px 1px rgba(255, 255, 255, 0.35),
        inset 0 -1px 1px rgba(0, 0, 0, 0.15);
      
      /* 文字样式 */
      color: #fff;
      font: 600 12px/1 -apple-system, BlinkMacSystemFont, "SF Pro Text", sans-serif;
      text-shadow: 0 1px 2px rgba(0, 0, 0, 0.3);
      
      /* 交互 */
      cursor: pointer;
      white-space: nowrap;
      transition: all 0.18s cubic-bezier(0.4, 0, 0.2, 1);
    }
    
    .dl-btn:hover {
      transform: translateY(-1.5px);
      background: linear-gradient(180deg, rgba(10, 132, 255, 1) 0%%, rgba(0, 106, 220, 0.98) 100%%);
      border-color: rgba(255, 255, 255, 0.38);
      box-shadow: 
        0 12px 28px rgba(10, 132, 255, 0.45),
        0 4px 12px rgba(0, 0, 0, 0.28),
        inset 0 1px 1px rgba(255, 255, 255, 0.4),
        inset 0 -1px 1px rgba(0, 0, 0, 0.15);
    }
    
    .dl-btn:active {
      transform: translateY(0);
      background: linear-gradient(180deg, rgba(0, 106, 220, 0.95) 0%%, rgba(10, 132, 255, 0.92) 100%%);
      box-shadow: 
        0 4px 12px rgba(10, 132, 255, 0.3),
        0 1px 4px rgba(0, 0, 0, 0.2),
        inset 0 1px 2px rgba(0, 0, 0, 0.2),
        inset 0 -1px 1px rgba(255, 255, 255, 0.15);
    }
    
    .dl-btn:focus-visible {
      outline: 2px solid rgba(96, 165, 250, 0.8);
      outline-offset: 3px;
    }

    /* ========== 响应式设计 (小屏) ========== */
    @media (max-width: 768px) {
      .dl-card {
        padding: 10px 12px;
        gap: 10px;
      }
      
      .dl-icon {
        width: 36px;
        height: 36px;
        min-width: 36px;
        font-size: 18px;
      }
      
      .dl-name {
        font-size: 12px;
      }
      
      .dl-btn {
        padding: 7px 12px;
        font-size: 11px;
        gap: 6px;
      }
    }

    /* ========== 明亮主题适配 ========== */
    @media (prefers-color-scheme: light) {
      /* 卡片 */
      .dl-card {
        background: rgba(250, 250, 252, 0.86);
        border-color: rgba(0, 0, 0, 0.08);
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.10), inset 0 0 1px rgba(0, 0, 0, 0.04);
      }
      
      .dl-card:hover {
        background: rgba(245, 245, 247, 0.94);
        border-color: rgba(0, 0, 0, 0.12);
        box-shadow: 0 8px 20px rgba(0, 0, 0, 0.14), inset 0 0 1px rgba(0, 0, 0, 0.08);
      }
      
      /* 文字 */
      .dl-name {
        color: rgba(0, 0, 0, 0.92);
        text-shadow: none;
      }
      
      .dl-type {
        color: rgba(0, 0, 0, 0.55);
      }
      
      /* 图标 */
      .dl-icon {
        background: linear-gradient(135deg, rgba(0, 122, 255, 0.08), rgba(0, 122, 255, 0.15));
        border-color: rgba(0, 122, 255, 0.12);
        box-shadow: none;
      }
      
      /* 按钮 */
      .dl-btn {
        background: linear-gradient(180deg, rgba(0, 122, 255, 0.96) 0%%, rgba(0, 106, 220, 0.94) 100%%);
        border: 1.5px solid rgba(255, 255, 255, 0.5);
        box-shadow: 
          0 6px 18px rgba(0, 122, 255, 0.32),
          0 2px 6px rgba(0, 0, 0, 0.12),
          inset 0 1px 1px rgba(255, 255, 255, 0.5),
          inset 0 -1px 1px rgba(0, 0, 0, 0.1);
      }
      
      .dl-btn:hover {
        background: linear-gradient(180deg, rgba(0, 122, 255, 1) 0%%, rgba(0, 106, 220, 0.98) 100%%);
        border-color: rgba(255, 255, 255, 0.6);
        box-shadow: 
          0 10px 24px rgba(0, 122, 255, 0.4),
          0 3px 10px rgba(0, 0, 0, 0.15),
          inset 0 1px 1px rgba(255, 255, 255, 0.55),
          inset 0 -1px 1px rgba(0, 0, 0, 0.1);
      }
    }
  </style>

  <a class="dl-card" href="%s" download="%s" title="Download %s" aria-label="Download %s">
    <div class="dl-icon" aria-hidden="true">%s</div>
    <div class="dl-info">
      <div class="dl-name">%s</div>
      <div class="dl-type">%s file</div>
    </div>
    <div class="dl-btn">
      <svg width="14" height="14" viewBox="0 0 24 24" fill="none" aria-hidden="true">
        <path d="M12 3v12m0 0l-4-4m4 4l4-4M5 21h14" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
      </svg>
      Download
    </div>
  </a>
</div>`,
		safeRel, safeFileName, safeFileName, safeFileName,
		fileIcon, safeFileName, fileType,
	)
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
