package parser

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"log"
)

// GeneratePDFWarpHTML 生成用于在Web上展示PDF内容的模块
// 使用PDF.js 进行优化展示内容，增加便利性，提升用户体验，优化加载速度
// 支持多实例隔离，按需加载，减少流量浪费
func GeneratePDFWarpHTML(pdfId, pdfPath string) string {
	// 确保输入安全
	safePdfId := template.HTMLEscapeString(pdfId)
	safePdfPath := template.HTMLEscapeString(pdfPath)
	log.Println("Generating PDF player for FileName:", safePdfPath)
	// 生成唯一ID（防止多实例冲突）
	h := md5.Sum([]byte(safePdfId))
	uniqueId := fmt.Sprintf("pdf_%x", h)[:12]

	return fmt.Sprintf(`
<div id="%[1]s-wrapper" class="pdfjs-viewer-wrapper" style="max-width:900px;margin:20px auto;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,sans-serif;">
  <style>
    #%[1]s-wrapper .pdfjs-toolbar {
      background: linear-gradient(135deg, #2c3e50 0%%, #34495e 100%%);
      color: #fff;
      padding: 10px 15px;
      display: flex;
      align-items: center;
      gap: 12px;
      border-radius: 8px 8px 0 0;
      box-shadow: 0 2px 8px rgba(0,0,0,0.15);
      flex-wrap: wrap;
    }
    #%[1]s-wrapper .pdfjs-btn {
      background: rgba(255,255,255,0.15);
      border: none;
      color: #fff;
      padding: 6px 14px;
      border-radius: 5px;
      cursor: pointer;
      font-size: 13px;
      transition: all 0.2s;
      font-weight: 500;
    }
    #%[1]s-wrapper .pdfjs-btn:hover {
      background: rgba(255,255,255,0.25);
      transform: translateY(-1px);
    }
    #%[1]s-wrapper .pdfjs-btn:disabled {
      opacity: 0.5;
      cursor: not-allowed;
    }
    #%[1]s-wrapper .pdfjs-btn-primary {
      background: #3498db;
    }
    #%[1]s-wrapper .pdfjs-btn-primary:hover {
      background: #2980b9;
    }
    #%[1]s-wrapper .pdfjs-btn-download {
      background: #f39c12;
      color: #fff;
    }
    #%[1]s-wrapper .pdfjs-btn-download:hover {
      background: #e67e22;
    }
    #%[1]s-wrapper .pdfjs-page-info {
      margin: 0 10px;
      font-size: 14px;
      color: #ecf0f1;
    }
    #%[1]s-wrapper .pdfjs-jump-input {
      width: 60px;
      padding: 4px 8px;
      border-radius: 4px;
      border: 1px solid rgba(255,255,255,0.3);
      background: rgba(255,255,255,0.1);
      color: #fff;
      font-size: 13px;
    }
    #%[1]s-wrapper .pdfjs-jump-input:focus {
      outline: none;
      border-color: #3498db;
      background: rgba(255,255,255,0.15);
    }
    #%[1]s-wrapper .pdfjs-link {
      color: #f39c12;
      padding: 0 10px;
      text-decoration: none;
      font-size: 13px;
      transition: color 0.2s;
    }
    #%[1]s-wrapper .pdfjs-link:hover {
      color: #e67e22;
    }
    #%[1]s-wrapper .pdfjs-spacer {
      flex: 1;
    }
    #%[1]s-wrapper .pdfjs-canvas-container {
      background: #f5f5f5;
      min-height: 600px;
      border-radius: 0 0 8px 8px;
      box-shadow: 0 4px 16px rgba(0,0,0,0.1);
      overflow-x: auto;
      position: relative;
    }
    #%[1]s-wrapper .pdfjs-canvas {
      width: 100%%;
      max-width: 900px;
      display: block;
      margin: 0 auto;
      box-shadow: 0 2px 12px rgba(0,0,0,0.08);
    }
    #%[1]s-wrapper .pdfjs-loading {
      text-align: center;
      padding: 50px 30px;
      color: #7f8c8d;
      font-size: 14px;
    }
    #%[1]s-wrapper .pdfjs-spinner {
      border: 3px solid #ecf0f1;
      border-top: 3px solid #3498db;
      border-radius: 50%%;
      width: 40px;
      height: 40px;
      animation: spin 1s linear infinite;
      margin: 0 auto 15px;
    }
    @keyframes spin {
      0%% { transform: rotate(0deg); }
      100%% { transform: rotate(360deg); }
    }
    @media (max-width: 768px) {
      #%[1]s-wrapper {
        margin: 10px;
      }
      #%[1]s-wrapper .pdfjs-toolbar {
        padding: 8px 10px;
        gap: 8px;
      }
      #%[1]s-wrapper .pdfjs-btn {
        padding: 5px 10px;
        font-size: 12px;
      }
      #%[1]s-wrapper .pdfjs-page-info {
        font-size: 12px;
      }
    }
  </style>
  <div class="pdfjs-toolbar">
    <button id="%[1]s-prev" class="pdfjs-btn" disabled>⬅ Prev</button>
    <button id="%[1]s-next" class="pdfjs-btn">Next ➡</button>
    <span class="pdfjs-page-info">
      Page <span id="%[1]s-current-page">1</span> / <span id="%[1]s-total-pages">?</span>
    </span>
    <input id="%[1]s-jump" class="pdfjs-jump-input" type="number" min="1" placeholder="Page" />
    <button id="%[1]s-jump-btn" class="pdfjs-btn pdfjs-btn-primary">Go</button>
    <span class="pdfjs-spacer"></span>
    <a href="%[2]s" target="_blank" rel="noopener noreferrer" class="pdfjs-link">🔗 Open</a>
    <button id="%[1]s-download-btn" class="pdfjs-btn pdfjs-btn-download">⬇ Download</button>
  </div>
  <div class="pdfjs-canvas-container">
    <div id="%[1]s-loading" class="pdfjs-loading">
      <div class="pdfjs-spinner"></div>
      <div>Loading PDF...</div>
    </div>
    <canvas id="%[1]s-canvas" class="pdfjs-canvas" style="display:none;"></canvas>
  </div>
</div>
<script src="https://cdn.jsdelivr.net/npm/pdfjs-dist@3.11.174/build/pdf.min.js"></script>
<script>
(function() {
  const UNIQUE_ID = '%[1]s';
  const PDF_URL = '%[2]s';
  // 状态管理
  const state = {
    pdfDoc: null,
    pageNum: 1,
    pageRendering: false,
    pageNumPending: null,
    scale: 1.15
  };
  // DOM 元素
  const elements = {
    canvas: document.getElementById(UNIQUE_ID + '-canvas'),
    loading: document.getElementById(UNIQUE_ID + '-loading'),
    prevBtn: document.getElementById(UNIQUE_ID + '-prev'),
    nextBtn: document.getElementById(UNIQUE_ID + '-next'),
    totalPagesEl: document.getElementById(UNIQUE_ID + '-total-pages'),
    currentPageEl: document.getElementById(UNIQUE_ID + '-current-page'),
    jumpInput: document.getElementById(UNIQUE_ID + '-jump'),
    jumpBtn: document.getElementById(UNIQUE_ID + '-jump-btn'),
    downloadBtn: document.getElementById(UNIQUE_ID + '-download-btn')
  };
  const ctx = elements.canvas.getContext('2d');
  // 更新按钮状态
  function updateButtons() {
    elements.prevBtn.disabled = state.pageNum <= 1;
    elements.nextBtn.disabled = state.pageNum >= state.pdfDoc.numPages;
  }
  // 渲染页面
  function renderPage(num) {
    state.pageRendering = true;
    elements.loading.style.display = 'block';
    elements.canvas.style.display = 'none';
    state.pdfDoc.getPage(num).then(function(page) {
      const viewport = page.getViewport({ scale: state.scale });
      elements.canvas.height = viewport.height;
      elements.canvas.width = viewport.width;
      const renderContext = {
        canvasContext: ctx,
        viewport: viewport
      };
      const renderTask = page.render(renderContext);
      renderTask.promise.then(function() {
        state.pageRendering = false;
        elements.loading.style.display = 'none';
        elements.canvas.style.display = 'block';
        elements.currentPageEl.textContent = num;
        // 更新按钮状态
        updateButtons();
        // 如果有待渲染的页面，继续渲染
        if (state.pageNumPending !== null) {
          renderPage(state.pageNumPending);
          state.pageNumPending = null;
        }
      }).catch(function(error) {
        console.error('Render error:', error);
        elements.loading.innerHTML = '<div style="color:#e74c3c;">Failed to render page.</div>';
        state.pageRendering = false;
      });
    }).catch(function(error) {
      console.error('Get page error:', error);
      elements.loading.innerHTML = '<div style="color:#e74c3c;">Failed to load page.</div>';
      state.pageRendering = false;
    });
  }
  // 队列渲染页面
  function queueRenderPage(num) {
    if (state.pageRendering) {
      state.pageNumPending = num;
    } else {
      renderPage(num);
    }
  }
  // 上一页
  function onPrevPage() {
    if (state.pageNum <= 1) return;
    state.pageNum--;
    queueRenderPage(state.pageNum);
  }
  // 下一页
  function onNextPage() {
    if (state.pageNum >= state.pdfDoc.numPages) return;
    state.pageNum++;
    queueRenderPage(state.pageNum);
  }
  // 跳转页面
  function onJumpPage() {
    const jumpTo = parseInt(elements.jumpInput.value);
    if (!isNaN(jumpTo) && jumpTo >= 1 && jumpTo <= state.pdfDoc.numPages) {
      state.pageNum = jumpTo;
      queueRenderPage(state.pageNum);
      elements.jumpInput.value = '';
    }
  }
  // 绑定事件
  elements.prevBtn.addEventListener('click', onPrevPage);
  elements.nextBtn.addEventListener('click', onNextPage);
  elements.jumpBtn.addEventListener('click', onJumpPage);
  elements.jumpInput.addEventListener('keydown', function(e) {
    if (e.key === 'Enter') onJumpPage();
  });
  elements.downloadBtn.addEventListener('click', function() {
    window.open(PDF_URL, '_blank');
  });
  // 键盘快捷键
  document.addEventListener('keydown', function(e) {
    // 检查焦点是否在当前PDF查看器区域
    const wrapper = document.getElementById(UNIQUE_ID + '-wrapper');
    if (!wrapper || !wrapper.contains(document.activeElement)) {
      return;
    }
    if (e.key === 'ArrowLeft' || e.key === 'PageUp') {
      e.preventDefault();
      onPrevPage();
    } else if (e.key === 'ArrowRight' || e.key === 'PageDown') {
      e.preventDefault();
      onNextPage();
    }
  });
  // 配置 PDF.js Worker
  if (typeof pdfjsLib !== 'undefined') {
    pdfjsLib.GlobalWorkerOptions.workerSrc = 
      'https://cdn.jsdelivr.net/npm/pdfjs-dist@3.11.174/build/pdf.worker.min.js';
  }
  // 加载PDF（优化配置：按需加载，减少流量）
  const loadingTask = pdfjsLib.getDocument({
    url: PDF_URL,
    rangeChunkSize: 65536 * 2,      // 128KB 分块加载
    disableAutoFetch: true,      // 禁用自动获取后续页面
    disableStream: false,        // 启用流式加载
    cMapUrl: 'https://cdn.jsdelivr.net/npm/pdfjs-dist@3.11.174/cmaps/',
    cMapPacked: true,
    standardFontDataUrl: 'https://cdn.jsdelivr.net/npm/pdfjs-dist@3.11.174/standard_fonts/'
  });
  loadingTask.promise.then(function(pdfDoc) {
    state.pdfDoc = pdfDoc;
    elements.totalPagesEl.textContent = pdfDoc.numPages;
    elements.jumpInput.max = pdfDoc.numPages;
    // 仅渲染第一页
    renderPage(state.pageNum);
    console.log('PDF loaded successfully. Total pages:', pdfDoc.numPages);
  }).catch(function(error) {
    console.error('PDF loading error:', error);
    elements.loading.innerHTML = 
      '<div style="color:#e74c3c;">❌ Failed to load PDF<br><small>' + 
      error.message + '</small></div>';
  });
  // 加载进度监听（可选）
  loadingTask.onProgress = function(progress) {
    if (progress.total > 0) {
      const percent = Math.round((progress.loaded / progress.total) * 100);
      elements.loading.innerHTML = 
        '<div class="pdfjs-spinner"></div>' +
        '<div>Loading PDF... ' + percent + '%%</div>';
    }
  };
})();
</script>
`, uniqueId, safePdfPath)
}
