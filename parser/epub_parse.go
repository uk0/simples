package parser

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"log"
	"path/filepath"
	"strings"
)

// GenerateEPUBViewerHTML 生成 EPUB 电子书查看器
// 基于 epub.js 核心渲染引擎，样式参考 macOS 阅读器设计
func GenerateEPUBViewerHTML(bookId, bookPath string) string {
	safeBookId := template.HTMLEscapeString(bookId)
	safeBookPath := template.HTMLEscapeString(bookPath)
	log.Println("Generating EPUB viewer for file:", safeBookPath)
	h := md5.Sum([]byte(safeBookId))
	uniqueId := fmt.Sprintf("epub_%x", h)[:12]
	fileName := filepath.Base(bookPath)
	ext := strings.ToLower(filepath.Ext(fileName))
	// 确认是 EPUB 格式
	if ext != ".epub" {
		log.Printf("Warning: File %s is not .epub format, got %s", fileName, ext)
	}
	displayName := strings.TrimSuffix(fileName, ext)
	if displayName == "" {
		displayName = fileName
	}
	return fmt.Sprintf(`<div id="%[1]s-wrapper" class="epub-viewer-wrapper" style="position:relative;max-width:100%%;margin:20px auto;display:block;">
<style>
#%[1]s-wrapper{position:relative;width:100%%;max-width:1400px;margin:20px auto;display:block;clear:both;font-family:-apple-system,BlinkMacSystemFont,'SF Pro Text','Segoe UI',Roboto,sans-serif}
#%[1]s-wrapper *{box-sizing:border-box}
#%[1]s-wrapper .epub-container{position:relative;width:100%%;background:rgba(30,30,32,0.95);backdrop-filter:blur(40px) saturate(180%%);-webkit-backdrop-filter:blur(40px) saturate(180%%);border:1px solid rgba(255,255,255,0.12);border-radius:16px;box-shadow:0 8px 32px rgba(0,0,0,0.24),0 0 1px rgba(255,255,255,0.12) inset;overflow:hidden;transition:all 0.3s ease}
#%[1]s-wrapper .epub-header{background:rgba(0,0,0,0.3);padding:16px 20px;border-bottom:1px solid rgba(255,255,255,0.08);display:flex;align-items:center;justify-content:space-between;gap:12px;flex-wrap:wrap}
#%[1]s-wrapper .epub-title{font-size:14px;font-weight:600;color:rgba(255,255,255,0.98);display:flex;align-items:center;gap:10px;text-shadow:0 1px 3px rgba(0,0,0,0.3);flex:1;min-width:200px}
#%[1]s-wrapper .epub-icon{font-size:20px;flex-shrink:0}
#%[1]s-wrapper .epub-title-text{flex:1;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
#%[1]s-wrapper .epub-badge{background:rgba(52,199,89,0.2);border:1px solid rgba(52,199,89,0.3);color:rgba(52,199,89,1);padding:4px 10px;border-radius:6px;font-size:10px;font-weight:600;text-transform:uppercase;letter-spacing:0.5px}
#%[1]s-wrapper .epub-toolbar{display:flex;align-items:center;gap:8px;flex-wrap:wrap}
#%[1]s-wrapper .epub-btn{background:rgba(255,255,255,0.12);border:1px solid rgba(255,255,255,0.08);color:rgba(255,255,255,0.95);padding:6px 12px;border-radius:8px;cursor:pointer;font-size:12px;font-weight:500;transition:all 0.2s;display:inline-flex;align-items:center;gap:6px;text-decoration:none;box-shadow:0 2px 6px rgba(0,0,0,0.15);white-space:nowrap}
#%[1]s-wrapper .epub-btn:hover{background:rgba(255,255,255,0.2);border-color:rgba(255,255,255,0.15);transform:translateY(-1px);box-shadow:0 3px 10px rgba(0,0,0,0.2)}
#%[1]s-wrapper .epub-btn:active{transform:translateY(0)}
#%[1]s-wrapper .epub-btn:disabled{opacity:0.3;cursor:not-allowed;transform:none}
#%[1]s-wrapper .epub-btn-primary{background:linear-gradient(135deg,#34C759 0%%,#28A745 100%%);border-color:transparent;box-shadow:0 3px 10px rgba(52,199,89,0.4)}
#%[1]s-wrapper .epub-btn-primary:hover{background:linear-gradient(135deg,#30D158 0%%,#2DB84C 100%%);box-shadow:0 4px 14px rgba(52,199,89,0.5)}
#%[1]s-wrapper .epub-main{position:relative;height:700px;overflow:hidden}
#%[1]s-wrapper #%[1]s-viewer{width:100%%;height:100%%;position:relative;overflow:hidden;background:#fff}
#%[1]s-wrapper #%[1]s-loading{position:absolute;top:50%%;left:50%%;transform:translate(-50%%,-50%%);z-index:100;background:rgba(255,255,255,0.95);padding:20px 30px;border-radius:8px;box-shadow:0 4px 12px rgba(0,0,0,0.1);font-size:14px;color:#333}
#%[1]s-wrapper .epub-controls{background:rgba(0,0,0,0.85);backdrop-filter:blur(10px);padding:12px 20px;border-top:1px solid rgba(255,255,255,0.06);display:flex;align-items:center;justify-content:center;gap:12px}
#%[1]s-wrapper .epub-nav-btn{background:rgba(255,255,255,0.15);border:none;color:rgba(255,255,255,0.95);padding:6px 12px;border-radius:6px;cursor:pointer;font-size:11px;transition:all 0.2s;display:flex;align-items:center;gap:4px}
#%[1]s-wrapper .epub-nav-btn:hover{background:rgba(255,255,255,0.25)}
#%[1]s-wrapper .epub-nav-btn:disabled{opacity:0.3;cursor:not-allowed}
#%[1]s-wrapper .epub-progress-container{display:flex;align-items:center;gap:10px;flex:1;max-width:400px}
#%[1]s-wrapper .epub-page-slider{flex:1;height:4px;-webkit-appearance:none;appearance:none;background:rgba(255,255,255,0.2);border-radius:2px;outline:none}
#%[1]s-wrapper .epub-page-slider::-webkit-slider-thumb{-webkit-appearance:none;appearance:none;width:16px;height:16px;background:rgba(52,199,89,1);border-radius:50%%;cursor:pointer}
#%[1]s-wrapper .epub-page-slider::-moz-range-thumb{width:16px;height:16px;background:rgba(52,199,89,1);border-radius:50%%;cursor:pointer;border:none}
#%[1]s-wrapper .epub-page-info{font-size:11px;color:rgba(255,255,255,0.9);font-family:'SF Mono','Monaco','Consolas',monospace;min-width:80px;text-align:center}
#%[1]s-wrapper .epub-settings{display:flex;gap:6px}
#%[1]s-wrapper .epub-footer{background:rgba(0,0,0,0.2);padding:12px 20px;border-top:1px solid rgba(255,255,255,0.06);display:flex;align-items:center;justify-content:space-between;font-size:11px;color:rgba(255,255,255,0.7)}
#%[1]s-wrapper .epub-error{display:none;position:absolute;top:50%%;left:50%%;transform:translate(-50%%,-50%%);z-index:100;background:#fff;padding:40px;border-radius:12px;box-shadow:0 8px 24px rgba(0,0,0,0.15);text-align:center;max-width:400px}
#%[1]s-wrapper .epub-error h3{color:#ff453a;margin:0 0 10px 0}
#%[1]s-wrapper .epub-error p{color:#666;margin:10px 0}
@media (max-width:768px){
#%[1]s-wrapper{margin:10px}
#%[1]s-wrapper .epub-header{padding:12px 16px}
#%[1]s-wrapper .epub-main{height:500px}
}
</style>
<div class="epub-container">
<div class="epub-header">
<div class="epub-title">
<span class="epub-icon">📚</span>
<span class="epub-title-text" id="%[1]s-title">%[3]s</span>
<span class="epub-badge">EPUB</span>
</div>
<div class="epub-toolbar">
<button id="%[1]s-fontDown" class="epub-btn">A-</button>
<button id="%[1]s-fontUp" class="epub-btn">A+</button>
<button id="%[1]s-reset" class="epub-btn">🔄 重置</button>
<a href="%[2]s" download class="epub-btn epub-btn-primary">⬇ 下载</a>
</div>
</div>
<div class="epub-main">
<div id="%[1]s-viewer">
<div id="%[1]s-loading">正在加载书籍...</div>
</div>
<div id="%[1]s-error" class="epub-error">
<h3>❌ 加载失败</h3>
<p id="%[1]s-error-msg">无法加载 EPUB 文件</p>
<div style="margin-top:20px">
<a href="%[2]s" download class="epub-btn epub-btn-primary">⬇ 下载文件</a>
</div>
</div>
</div>
<div class="epub-controls">
<button id="%[1]s-prev" class="epub-nav-btn">◀</button>
<div class="epub-progress-container">
<span id="%[1]s-pageInfo" class="epub-page-info">---</span>
<input type="range" id="%[1]s-pageSlider" class="epub-page-slider" min="0" max="100" value="0">
</div>
<button id="%[1]s-next" class="epub-nav-btn">▶</button>
<div class="epub-settings">
<input type="number" id="%[1]s-pageInput" min="1" value="1" style="width:50px;padding:4px;border-radius:4px;border:1px solid rgba(255,255,255,0.3);background:rgba(255,255,255,0.1);color:#fff;font-size:11px">
<button id="%[1]s-pageGo" class="epub-btn" style="padding:4px 8px;font-size:11px">跳转</button>
</div>
</div>
<div class="epub-footer">
<div>Powered by <strong>EPUB.js</strong></div>
<div>使用 ← → 或触摸滑动翻页</div>
</div>
</div>
</div>
<!-- 先加载 JSZip -->
<script src="https://cdnjs.cloudflare.com/ajax/libs/jszip/3.7.1/jszip.min.js"></script>
<!-- 然后加载 epub.js -->
<script src="https://cdn.jsdelivr.net/npm/epubjs/dist/epub.min.js"></script>
<script>
(function(){
'use strict';
const UNIQUE_ID = '%[1]s';
const BOOK_PATH = '%[2]s';
// 构建完整URL
function getBookUrl() {
    const path = BOOK_PATH;
    if (path.startsWith('http://') || path.startsWith('https://')) {
        return path;
    }
    const protocol = window.location.protocol;
    const host = window.location.host;
    const basePath = window.location.pathname.substring(0, window.location.pathname.lastIndexOf('/'));
    if (path.startsWith('/')) {
        return protocol + '//' + host + path;
    } else {
        return protocol + '//' + host + basePath + '/' + path;
    }
}
const BOOK_URL = getBookUrl();
console.log('[EPUB] Book URL:', BOOK_URL);
// 存储键
const pageKey = UNIQUE_ID + '_page';
const cfiKey = UNIQUE_ID + '_cfi';
const fontKey = UNIQUE_ID + '_fontSize';
// 状态
let book, rendition, locations;
let currentFontSize = parseInt(localStorage.getItem(fontKey)) || 100;
let lastLocation = localStorage.getItem(cfiKey) || null;
// DOM 元素
const elements = {
    viewer: document.getElementById(UNIQUE_ID + '-viewer'),
    loading: document.getElementById(UNIQUE_ID + '-loading'),
    error: document.getElementById(UNIQUE_ID + '-error'),
    errorMsg: document.getElementById(UNIQUE_ID + '-error-msg'),
    title: document.getElementById(UNIQUE_ID + '-title'),
    prev: document.getElementById(UNIQUE_ID + '-prev'),
    next: document.getElementById(UNIQUE_ID + '-next'),
    fontUp: document.getElementById(UNIQUE_ID + '-fontUp'),
    fontDown: document.getElementById(UNIQUE_ID + '-fontDown'),
    reset: document.getElementById(UNIQUE_ID + '-reset'),
    pageInfo: document.getElementById(UNIQUE_ID + '-pageInfo'),
    pageSlider: document.getElementById(UNIQUE_ID + '-pageSlider'),
    pageInput: document.getElementById(UNIQUE_ID + '-pageInput'),
    pageGo: document.getElementById(UNIQUE_ID + '-pageGo')
};
function showError(msg) {
    console.error('[EPUB] Error:', msg);
    elements.loading.style.display = 'none';
    elements.error.style.display = 'block';
    if (elements.errorMsg) {
        elements.errorMsg.textContent = msg;
    }
}
function loadBook(source) {
    console.log('[EPUB] Loading book...');
    elements.loading.style.display = 'block';
    // 销毁旧的渲染器
    if (rendition) {
        rendition.destroy();
    }
    try {
        // 创建 book 实例
        book = ePub(source);
        console.log('[EPUB] Book created');
        // 创建渲染器
        rendition = book.renderTo(elements.viewer, {
            width: "100%%",
            height: "100%%"
        });
        console.log('[EPUB] Rendition created');
        // 恢复字体大小
        rendition.themes.fontSize(currentFontSize + "%%");
        // 显示书籍
        rendition.display(lastLocation || undefined).then(function() {
            console.log('[EPUB] Book displayed');
            elements.loading.style.display = 'none';
        }).catch(function(error) {
            console.error('[EPUB] Display error:', error);
            showError('显示内容失败: ' + error.message);
        });
        // book.ready 处理
        book.ready.then(function() {
            console.log('[EPUB] Book ready');
            // 生成位置信息
            return book.locations.generate(1000);
        }).then(function(locs) {
            console.log('[EPUB] Locations generated:', locs.length);
            locations = book.locations;
            updatePageDisplay();
        }).catch(function(error) {
            console.warn('[EPUB] Locations error:', error);
        });
        // 加载元数据
        book.loaded.metadata.then(function(meta) {
            console.log('[EPUB] Metadata:', meta);
            if (meta.title && elements.title) {
                elements.title.textContent = meta.title;
            }
        }).catch(function(error) {
            console.warn('[EPUB] Metadata error:', error);
        });
        // 监听位置变化
        rendition.on("relocated", function(location) {
            console.log('[EPUB] Relocated');
            lastLocation = location.start.cfi;
            localStorage.setItem(cfiKey, lastLocation);
            updatePageDisplay();
        });
    } catch(error) {
        console.error('[EPUB] Load error:', error);
        showError('加载失败: ' + error.message);
    }
}
function updatePageDisplay() {
    if (!locations) return;
    const currentPage = locations.locationFromCfi(lastLocation) || 0;
    const totalPages = locations.total;
    elements.pageSlider.max = totalPages;
    elements.pageSlider.value = currentPage;
    elements.pageInfo.textContent = currentPage + ' / ' + totalPages;
    console.log('[EPUB] Page:', currentPage, '/', totalPages);
}
// 导航事件
elements.prev.onclick = function() {
    if (rendition) rendition.prev();
};
elements.next.onclick = function() {
    if (rendition) rendition.next();
};
// 滑块导航
elements.pageSlider.addEventListener("input", function(e) {
    if (!locations) return;
    const loc = locations.cfiFromLocation(parseInt(e.target.value));
    if (loc) rendition.display(loc);
});
// 字体控制
elements.fontUp.onclick = function() {
    currentFontSize += 10;
    if (rendition) {
        rendition.themes.fontSize(currentFontSize + "%%");
        localStorage.setItem(fontKey, currentFontSize);
    }
};
elements.fontDown.onclick = function() {
    currentFontSize = Math.max(50, currentFontSize - 10);
    if (rendition) {
        rendition.themes.fontSize(currentFontSize + "%%");
        localStorage.setItem(fontKey, currentFontSize);
    }
};
// 页面跳转
elements.pageGo.onclick = function() {
    if (!locations) return;
    const pageNum = parseInt(elements.pageInput.value, 10);
    const totalPages = locations.total;
    if (pageNum >= 1 && pageNum <= totalPages) {
        const cfi = locations.cfiFromLocation(pageNum);
        if (cfi) rendition.display(cfi);
    }
};
// 输入框回车
elements.pageInput.addEventListener("keydown", function(e) {
    if (e.key === "Enter") {
        elements.pageGo.click();
    }
});
// 重置
elements.reset.onclick = function() {
    localStorage.removeItem(cfiKey);
    localStorage.removeItem(fontKey);
    currentFontSize = 100;
    lastLocation = null;
    loadBook(BOOK_URL);
};
// 键盘导航
document.addEventListener("keydown", function(e) {
    if (!rendition) return;
    if (e.target.tagName === 'INPUT' || e.target.tagName === 'TEXTAREA') return;
    if (e.key === "ArrowLeft") {
        rendition.prev();
    } else if (e.key === "ArrowRight") {
        rendition.next();
    }
});
// 触摸导航
let touchStartX = 0;
elements.viewer.addEventListener("touchstart", function(e) {
    touchStartX = e.changedTouches[0].screenX;
});
elements.viewer.addEventListener("touchend", function(e) {
    const touchEndX = e.changedTouches[0].screenX;
    if (!rendition) return;
    if (touchEndX - touchStartX > 50) {
        rendition.prev();
    } else if (touchStartX - touchEndX > 50) {
        rendition.next();
    }
});
// 开始加载 - 使用 fetch 获取文件
console.log('[EPUB] Fetching book...');
fetch(BOOK_URL)
    .then(function(resp) {
        if (!resp.ok) {
            throw new Error("HTTP " + resp.status);
        }
        console.log('[EPUB] Fetch successful');
        return resp.blob();
    })
    .then(function(blob) {
        console.log('[EPUB] Got blob:', blob.size, 'bytes');
        loadBook(blob);
    })
    .catch(function(error) {
        console.error('[EPUB] Fetch error:', error);
        showError('无法获取文件: ' + error.message);
    });
})();
</script>`, uniqueId, safeBookPath, displayName)
}
