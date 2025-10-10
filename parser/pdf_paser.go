package parser

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"log"
	"path/filepath"
)

// GeneratePDFWarpHTML 生成用于在Web上展示PDF内容的模块
// macOS Sequoia 风格，与 PPT parser 保持完全一致
func GeneratePDFWarpHTML(pdfId, pdfPath string) string {
	safePdfId := template.HTMLEscapeString(pdfId)
	safePdfPath := template.HTMLEscapeString(pdfPath)
	log.Println("Generating PDF viewer for FileName:", safePdfPath)
	h := md5.Sum([]byte(safePdfId))
	uniqueId := fmt.Sprintf("pdf_%x", h)[:12]
	fileName := filepath.Base(pdfPath)

	return fmt.Sprintf(`<div id="%[1]s-wrapper" class="pdf-viewer-wrapper" style="position:relative;max-width:100%%;margin:20px auto;display:block;">
<style>
#%[1]s-wrapper{position:relative;width:100%%;max-width:1400px;margin:20px auto;display:block;clear:both;font-family:-apple-system,BlinkMacSystemFont,'SF Pro Text','Segoe UI',Roboto,sans-serif}
#%[1]s-wrapper *{box-sizing:border-box}
#%[1]s-wrapper .pdf-container{position:relative;width:100%%;background:rgba(30,30,32,0.75);backdrop-filter:blur(40px) saturate(180%%);-webkit-backdrop-filter:blur(40px) saturate(180%%);border:1px solid rgba(255,255,255,0.12);border-radius:16px;box-shadow:0 8px 32px rgba(0,0,0,0.24),0 0 1px rgba(255,255,255,0.12) inset;overflow:hidden;transition:all 0.3s ease}
#%[1]s-wrapper .pdf-container:hover{background:rgba(40,40,42,0.8);box-shadow:0 12px 48px rgba(0,0,0,0.32),0 0 1px rgba(255,255,255,0.18) inset}
#%[1]s-wrapper .pdf-header{background:rgba(0,0,0,0.3);padding:16px 20px;border-bottom:1px solid rgba(255,255,255,0.08);display:flex;align-items:center;justify-content:space-between;gap:12px;flex-wrap:wrap}
#%[1]s-wrapper .pdf-title{font-size:14px;font-weight:600;color:rgba(255,255,255,0.98);display:flex;align-items:center;gap:10px;text-shadow:0 1px 3px rgba(0,0,0,0.3);flex:1;min-width:200px}
#%[1]s-wrapper .pdf-icon{font-size:20px;flex-shrink:0}
#%[1]s-wrapper .pdf-title-text{flex:1;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
#%[1]s-wrapper .pdf-toolbar{display:flex;align-items:center;gap:8px;flex-wrap:wrap}
#%[1]s-wrapper .pdf-btn{background:rgba(255,255,255,0.12);border:1px solid rgba(255,255,255,0.08);color:rgba(255,255,255,0.95);padding:6px 12px;border-radius:8px;cursor:pointer;font-size:12px;font-weight:500;transition:all 0.2s;display:inline-flex;align-items:center;gap:6px;text-decoration:none;box-shadow:0 2px 6px rgba(0,0,0,0.15);white-space:nowrap}
#%[1]s-wrapper .pdf-btn:hover{background:rgba(255,255,255,0.2);border-color:rgba(255,255,255,0.15);transform:translateY(-1px);box-shadow:0 3px 10px rgba(0,0,0,0.2)}
#%[1]s-wrapper .pdf-btn:active{transform:translateY(0)}
#%[1]s-wrapper .pdf-btn:disabled{opacity:0.3;cursor:not-allowed;transform:none}
#%[1]s-wrapper .pdf-btn-primary{background:linear-gradient(135deg,#007AFF 0%%,#0051D5 100%%);border-color:transparent;box-shadow:0 3px 10px rgba(0,122,255,0.4)}
#%[1]s-wrapper .pdf-btn-primary:hover{background:linear-gradient(135deg,#0A84FF 0%%,#0060EA 100%%);box-shadow:0 4px 14px rgba(0,122,255,0.5)}
#%[1]s-wrapper .pdf-controls{display:flex;align-items:center;gap:8px}
#%[1]s-wrapper .pdf-page-info{display:flex;align-items:center;gap:8px;padding:4px 12px;background:rgba(0,0,0,0.2);border-radius:8px;font-size:12px;color:rgba(255,255,255,0.9);font-weight:500;font-family:'SF Mono','Monaco','Consolas',monospace}
#%[1]s-wrapper .pdf-nav-btn{width:32px;height:32px;flex-shrink:0;border-radius:50%%;background:rgba(255,255,255,0.12);border:1px solid rgba(255,255,255,0.08);color:rgba(255,255,255,0.9);cursor:pointer;display:flex;align-items:center;justify-content:center;transition:all 0.2s;font-size:14px}
#%[1]s-wrapper .pdf-nav-btn:hover:not(:disabled){background:rgba(255,255,255,0.2);transform:scale(1.05)}
#%[1]s-wrapper .pdf-nav-btn:disabled{opacity:0.3;cursor:not-allowed}
#%[1]s-wrapper .pdf-jump-input{width:50px;padding:2px 6px;border-radius:4px;border:1px solid rgba(255,255,255,0.2);background:rgba(0,0,0,0.2);color:rgba(255,255,255,0.9);font-size:11px;font-family:'SF Mono','Monaco','Consolas',monospace;text-align:center}
#%[1]s-wrapper .pdf-jump-input:focus{outline:none;border-color:rgba(10,132,255,0.5);background:rgba(0,0,0,0.3)}
#%[1]s-wrapper .pdf-content{padding:30px 20px;background:rgba(0,0,0,0.15);position:relative;overflow:visible;display:flex;align-items:center;justify-content:center}
#%[1]s-wrapper .pdf-canvas-container{width:100%%;max-width:1200px;background:#fff;border-radius:12px;box-shadow:0 12px 48px rgba(0,0,0,0.3);position:relative;overflow:hidden;display:flex;align-items:center;justify-content:center;min-height:600px}
#%[1]s-wrapper .pdf-canvas{max-width:100%%!important;height:auto!important;display:block;margin:0 auto}
#%[1]s-wrapper .pdf-loading{display:flex;align-items:center;justify-content:center;padding:100px 20px;color:rgba(255,255,255,0.95);font-size:13px;font-weight:500;flex-direction:column;position:absolute;top:50%%;left:50%%;transform:translate(-50%%,-50%%);z-index:10}
#%[1]s-wrapper .pdf-spinner{border:3px solid rgba(255,255,255,0.2);border-top:3px solid rgba(10,132,255,0.95);border-radius:50%%;width:40px;height:40px;animation:pdf-spin-%[1]s 0.7s linear infinite;margin-bottom:16px;box-shadow:0 0 12px rgba(10,132,255,0.3)}
@keyframes pdf-spin-%[1]s{0%%{transform:rotate(0deg)}100%%{transform:rotate(360deg)}}
#%[1]s-wrapper .pdf-error{display:none;align-items:center;justify-content:center;padding:60px 20px;color:rgba(255,69,58,1);font-size:13px;font-weight:600;text-align:center;flex-direction:column;gap:12px;position:absolute;top:50%%;left:50%%;transform:translate(-50%%,-50%%);width:100%%;z-index:10}
#%[1]s-wrapper .pdf-stats{display:flex;align-items:center;gap:20px;padding:12px 20px;background:rgba(0,0,0,0.2);border-top:1px solid rgba(255,255,255,0.06);font-size:11px;color:rgba(255,255,255,0.7);flex-wrap:wrap}
#%[1]s-wrapper .pdf-stat{display:flex;align-items:center;gap:6px}
#%[1]s-wrapper .pdf-stat-label{font-weight:500;color:rgba(255,255,255,0.5)}
#%[1]s-wrapper .pdf-stat-value{font-weight:600;color:rgba(10,132,255,1);font-family:'SF Mono','Monaco','Consolas',monospace}
#%[1]s-wrapper .pdf-zoom-controls{display:flex;align-items:center;gap:6px;background:rgba(0,0,0,0.2);padding:4px;border-radius:8px}
#%[1]s-wrapper .pdf-zoom-btn{width:28px;height:28px;border-radius:6px;background:rgba(255,255,255,0.1);border:1px solid rgba(255,255,255,0.06);color:rgba(255,255,255,0.8);cursor:pointer;display:flex;align-items:center;justify-content:center;transition:all 0.2s;font-size:12px}
#%[1]s-wrapper .pdf-zoom-btn:hover{background:rgba(255,255,255,0.15);transform:scale(1.05)}
#%[1]s-wrapper .pdf-zoom-value{font-size:11px;color:rgba(255,255,255,0.8);min-width:45px;text-align:center;font-weight:500}
@media (max-width:768px){
#%[1]s-wrapper{margin:10px}
#%[1]s-wrapper .pdf-header{padding:12px 16px}
#%[1]s-wrapper .pdf-title{font-size:13px}
#%[1]s-wrapper .pdf-toolbar{width:100%%;justify-content:flex-start}
#%[1]s-wrapper .pdf-btn{padding:5px 10px;font-size:11px}
#%[1]s-wrapper .pdf-content{padding:15px 10px}
}
@media (prefers-color-scheme:light){
#%[1]s-wrapper .pdf-container{background:rgba(250,250,252,0.95);border-color:rgba(0,0,0,0.08)}
#%[1]s-wrapper .pdf-header{background:rgba(0,0,0,0.03);border-bottom-color:rgba(0,0,0,0.06)}
#%[1]s-wrapper .pdf-title{color:rgba(0,0,0,0.95);text-shadow:0 1px 2px rgba(255,255,255,0.5)}
#%[1]s-wrapper .pdf-btn{background:rgba(0,0,0,0.06);border-color:rgba(0,0,0,0.08);color:rgba(0,0,0,0.85)}
#%[1]s-wrapper .pdf-content{background:rgba(0,0,0,0.02)}
#%[1]s-wrapper .pdf-page-info{background:rgba(0,0,0,0.05);color:rgba(0,0,0,0.8)}
#%[1]s-wrapper .pdf-nav-btn{background:rgba(0,0,0,0.06);border-color:rgba(0,0,0,0.08);color:rgba(0,0,0,0.8)}
#%[1]s-wrapper .pdf-stats{background:rgba(0,0,0,0.03);border-top-color:rgba(0,0,0,0.06);color:rgba(0,0,0,0.7)}
#%[1]s-wrapper .pdf-jump-input{background:rgba(0,0,0,0.05);border-color:rgba(0,0,0,0.1);color:rgba(0,0,0,0.8)}
#%[1]s-wrapper .pdf-zoom-controls{background:rgba(0,0,0,0.05)}
#%[1]s-wrapper .pdf-zoom-btn{background:rgba(0,0,0,0.04);border-color:rgba(0,0,0,0.06);color:rgba(0,0,0,0.7)}
}
</style>
<div class="pdf-container">
<div class="pdf-header">
<div class="pdf-title">
<span class="pdf-icon">📄</span>
<span class="pdf-title-text">%[3]s</span>
</div>
<div class="pdf-toolbar">
<div class="pdf-controls">
<button id="%[1]s-prev" class="pdf-nav-btn" disabled title="Previous Page">◀</button>
<div class="pdf-page-info">
<span id="%[1]s-current-page">1</span>
<span>/</span>
<span id="%[1]s-total-pages">-</span>
</div>
<button id="%[1]s-next" class="pdf-nav-btn" title="Next Page">▶</button>
<input id="%[1]s-jump" class="pdf-jump-input" type="number" min="1" placeholder="#"/>
<button id="%[1]s-jump-btn" class="pdf-btn" title="Go to Page">Go</button>
</div>
<div class="pdf-zoom-controls">
<button id="%[1]s-zoom-out" class="pdf-zoom-btn" title="Zoom Out">−</button>
<span class="pdf-zoom-value" id="%[1]s-zoom-value">100%%</span>
<button id="%[1]s-zoom-in" class="pdf-zoom-btn" title="Zoom In">+</button>
<button id="%[1]s-zoom-fit" class="pdf-zoom-btn" title="Fit Page" style="font-size:10px;width:auto;padding:0 8px;">Fit</button>
</div>
<button id="%[1]s-fullscreen-btn" class="pdf-btn" title="Fullscreen">🔍 Fullscreen</button>
<a href="%[2]s" download class="pdf-btn pdf-btn-primary" title="Download">⬇ Download</a>
<a href="%[2]s" target="_blank" rel="noopener noreferrer" class="pdf-btn" title="Open">🔗 Open</a>
</div>
</div>
<div class="pdf-content" id="%[1]s-content">
<div id="%[1]s-loading" class="pdf-loading">
<div class="pdf-spinner"></div>
<div>Loading PDF...</div>
</div>
<div id="%[1]s-error" class="pdf-error">
<div>❌ Failed to load PDF</div>
<div style="font-size:11px;margin-top:8px;">Please check if the file is accessible</div>
</div>
<div class="pdf-canvas-container" id="%[1]s-canvas-container">
<canvas id="%[1]s-canvas" class="pdf-canvas" style="display:none;"></canvas>
</div>
</div>
<div class="pdf-stats" id="%[1]s-stats" style="display:none;">
<div class="pdf-stat">
<span class="pdf-stat-label">Total Pages:</span>
<span class="pdf-stat-value" id="%[1]s-total-stat">0</span>
</div>
<div class="pdf-stat">
<span class="pdf-stat-label">Current:</span>
<span class="pdf-stat-value" id="%[1]s-current-stat">1</span>
</div>
<div class="pdf-stat">
<span class="pdf-stat-label">Format:</span>
<span class="pdf-stat-value">PDF</span>
</div>
<div class="pdf-stat">
<span class="pdf-stat-label">Zoom:</span>
<span class="pdf-stat-value" id="%[1]s-zoom-stat">100%%</span>
</div>
</div>
</div>
</div>
<script>
(function(){
const UNIQUE_ID='%[1]s';
const PDF_URL='%[2]s';
const state={pdfDoc:null,pageNum:1,pageRendering:false,pageNumPending:null,scale:1.0,baseScale:1.0,loadedPages:new Set()};
const elements={
canvas:document.getElementById(UNIQUE_ID+'-canvas'),
canvasContainer:document.getElementById(UNIQUE_ID+'-canvas-container'),
loading:document.getElementById(UNIQUE_ID+'-loading'),
error:document.getElementById(UNIQUE_ID+'-error'),
stats:document.getElementById(UNIQUE_ID+'-stats'),
prevBtn:document.getElementById(UNIQUE_ID+'-prev'),
nextBtn:document.getElementById(UNIQUE_ID+'-next'),
totalPagesEl:document.getElementById(UNIQUE_ID+'-total-pages'),
currentPageEl:document.getElementById(UNIQUE_ID+'-current-page'),
totalStat:document.getElementById(UNIQUE_ID+'-total-stat'),
currentStat:document.getElementById(UNIQUE_ID+'-current-stat'),
jumpInput:document.getElementById(UNIQUE_ID+'-jump'),
jumpBtn:document.getElementById(UNIQUE_ID+'-jump-btn'),
downloadBtn:document.querySelector('a[download].pdf-btn-primary'),
fullscreenBtn:document.getElementById(UNIQUE_ID+'-fullscreen-btn'),
zoomInBtn:document.getElementById(UNIQUE_ID+'-zoom-in'),
zoomOutBtn:document.getElementById(UNIQUE_ID+'-zoom-out'),
zoomFitBtn:document.getElementById(UNIQUE_ID+'-zoom-fit'),
zoomValue:document.getElementById(UNIQUE_ID+'-zoom-value'),
zoomStat:document.getElementById(UNIQUE_ID+'-zoom-stat')
};
const ctx=elements.canvas.getContext('2d');
function hideLoading(){
elements.loading.style.display='none';
elements.canvas.style.display='block';
elements.stats.style.display='flex';
}
function showError(msg){
elements.loading.style.display='none';
elements.error.style.display='flex';
if(msg){
const errorDiv=elements.error.querySelector('div:last-child');
if(errorDiv)errorDiv.textContent=msg;
}
}
function loadPDFJS(callback){
if(typeof pdfjsLib!=='undefined'){
callback();
return;
}
if(document.getElementById('pdfjs-script')){
waitForPDFJS(callback);
return;
}
const script=document.createElement('script');
script.id='pdfjs-script';
script.src='https://cdn.jsdelivr.net/npm/pdfjs-dist@3.11.174/build/pdf.min.js';
script.onload=function(){
if(typeof pdfjsLib!=='undefined'){
pdfjsLib.GlobalWorkerOptions.workerSrc='https://cdn.jsdelivr.net/npm/pdfjs-dist@3.11.174/build/pdf.worker.min.js';
}
console.log('PDF.js loaded');
callback();
};
script.onerror=function(){
console.error('Failed to load PDF.js');
showError('Failed to load PDF.js library');
};
document.head.appendChild(script);
}
function waitForPDFJS(callback){
let attempts=0;
const checkPDFJS=setInterval(function(){
attempts++;
if(typeof pdfjsLib!=='undefined'){
clearInterval(checkPDFJS);
callback();
}else if(attempts>50){
clearInterval(checkPDFJS);
showError('PDF.js loading timeout');
}
},100);
}
function updateButtons(){
elements.prevBtn.disabled=state.pageNum<=1;
elements.nextBtn.disabled=state.pageNum>=state.pdfDoc.numPages;
}
function updateZoomDisplay(){
const zoomPercent=Math.round(state.scale*100);
elements.zoomValue.textContent=zoomPercent+'%%';
elements.zoomStat.textContent=zoomPercent+'%%';
}
function renderPage(num){
if(state.pageRendering){
state.pageNumPending=num;
return;
}
state.pageRendering=true;
elements.loading.style.display='flex';
elements.canvas.style.display='none';
state.pdfDoc.getPage(num).then(function(page){
const viewport=page.getViewport({scale:state.scale});
elements.canvas.height=viewport.height;
elements.canvas.width=viewport.width;
const renderContext={canvasContext:ctx,viewport:viewport};
const renderTask=page.render(renderContext);
renderTask.promise.then(function(){
state.pageRendering=false;
state.loadedPages.add(num);
hideLoading();
elements.currentPageEl.textContent=num;
elements.currentStat.textContent=num;
updateButtons();
updateZoomDisplay();
console.log('Page',num,'rendered at scale',state.scale);
if(state.pageNumPending!==null){
const pending=state.pageNumPending;
state.pageNumPending=null;
renderPage(pending);
}
}).catch(function(error){
console.error('Render error:',error);
showError('Failed to render page');
state.pageRendering=false;
});
}).catch(function(error){
console.error('Get page error:',error);
showError('Failed to load page');
state.pageRendering=false;
});
}
function onPrevPage(){
if(state.pageNum<=1)return;
state.pageNum--;
renderPage(state.pageNum);
}
function onNextPage(){
if(state.pageNum>=state.pdfDoc.numPages)return;
state.pageNum++;
renderPage(state.pageNum);
}
function onJumpPage(){
const jumpTo=parseInt(elements.jumpInput.value);
if(!isNaN(jumpTo)&&jumpTo>=1&&jumpTo<=state.pdfDoc.numPages){
state.pageNum=jumpTo;
renderPage(state.pageNum);
elements.jumpInput.value='';
}
}
function onZoomIn(){
state.scale=Math.min(state.scale*1.2,3.0);
renderPage(state.pageNum);
}
function onZoomOut(){
state.scale=Math.max(state.scale/1.2,0.5);
renderPage(state.pageNum);
}
function onZoomFit(){
state.scale=state.baseScale;
renderPage(state.pageNum);
}
function initPDF(){
loadPDFJS(function(){
const loadingTask=pdfjsLib.getDocument({
url:PDF_URL,
rangeChunkSize:65536,
disableAutoFetch:true,
disableStream:false,
cMapUrl:'https://cdn.jsdelivr.net/npm/pdfjs-dist@3.11.174/cmaps/',
cMapPacked:true
});
loadingTask.onProgress=function(progress){
if(progress.total>0){
const percent=Math.round((progress.loaded/progress.total)*100);
elements.loading.innerHTML='<div class="pdf-spinner"></div><div>Loading PDF... '+percent+'%%</div>';
}
};
loadingTask.promise.then(function(pdfDoc){
state.pdfDoc=pdfDoc;
elements.totalPagesEl.textContent=pdfDoc.numPages;
elements.totalStat.textContent=pdfDoc.numPages;
elements.jumpInput.max=pdfDoc.numPages;
pdfDoc.getPage(1).then(function(page){
const viewport=page.getViewport({scale:1.0});
const containerWidth=elements.canvasContainer.offsetWidth-40;
state.baseScale=containerWidth/viewport.width;
state.scale=state.baseScale;
renderPage(state.pageNum);
});
console.log('PDF loaded. Total pages:',pdfDoc.numPages);
}).catch(function(error){
console.error('PDF loading error:',error);
showError('Failed to load PDF: '+error.message);
});
});
}
elements.prevBtn.addEventListener('click',onPrevPage);
elements.nextBtn.addEventListener('click',onNextPage);
elements.jumpBtn.addEventListener('click',onJumpPage);
elements.jumpInput.addEventListener('keydown',function(e){
if(e.key==='Enter')onJumpPage();
});
elements.fullscreenBtn.addEventListener('click',function(){
const target=elements.canvasContainer;
if(target.requestFullscreen)target.requestFullscreen();
else if(target.webkitRequestFullscreen)target.webkitRequestFullscreen();
else if(target.msRequestFullscreen)target.msRequestFullscreen();
});
elements.zoomInBtn.addEventListener('click',onZoomIn);
elements.zoomOutBtn.addEventListener('click',onZoomOut);
elements.zoomFitBtn.addEventListener('click',onZoomFit);
document.addEventListener('keydown',function(e){
const wrapper=document.getElementById(UNIQUE_ID+'-wrapper');
if(!wrapper)return;
const activeElement=document.activeElement;
if(activeElement&&(activeElement.tagName==='INPUT'||activeElement.tagName==='TEXTAREA'))return;
const rect=wrapper.getBoundingClientRect();
const inViewport=rect.top<window.innerHeight&&rect.bottom>0;
if(!inViewport)return;
switch(e.key){
case 'ArrowLeft':
case 'PageUp':
onPrevPage();
e.preventDefault();
break;
case 'ArrowRight':
case 'PageDown':
case ' ':
onNextPage();
e.preventDefault();
break;
case '+':
case '=':
onZoomIn();
e.preventDefault();
break;
case '-':
case '_':
onZoomOut();
e.preventDefault();
break;
case '0':
onZoomFit();
e.preventDefault();
break;
case 'f':
elements.fullscreenBtn.click();
e.preventDefault();
break;
}
});
initPDF();
})();
</script>`, uniqueId, safePdfPath, fileName)
}
