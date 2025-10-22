package parser

import (
	"crypto/md5"
	"fmt"
	"html/template"

	"path/filepath"
	"strings"
)

// GenerateExcelWarpHTML 生成用于在Web上展示Excel内容的模块
// 企业级商业可用版本：支持多实例隔离、智能加载、重试机制、完善的错误处理
// macOS Sequoia 风格，支持缩放控制、多工作表切换
func GenerateExcelWarpHTML(excelId, excelPath string) string {
	safeExcelId := template.HTMLEscapeString(excelId)
	safeExcelPath := template.HTMLEscapeString(excelPath)

	h := md5.Sum([]byte(safeExcelId))
	uniqueId := fmt.Sprintf("excel_%x", h)[:12]
	fileName := filepath.Base(excelPath)
	ext := strings.ToLower(filepath.Ext(fileName))
	safeFileName := template.JSEscapeString(fileName)

	return fmt.Sprintf(`<div id="%[1]s-wrapper" class="excel-viewer-wrapper" style="position:relative;max-width:100%%;margin:20px auto;display:block;">
<style>
#%[1]s-wrapper{position:relative;width:100%%;max-width:1400px;margin:20px auto;display:block;clear:both;font-family:-apple-system,BlinkMacSystemFont,'SF Pro Text','Segoe UI',Roboto,sans-serif}
#%[1]s-wrapper *{box-sizing:border-box}
#%[1]s-wrapper .excel-container{position:relative;width:100%%;background:rgba(30,30,32,0.75);backdrop-filter:blur(40px) saturate(180%%);-webkit-backdrop-filter:blur(40px) saturate(180%%);border:1px solid rgba(255,255,255,0.12);border-radius:16px;box-shadow:0 8px 32px rgba(0,0,0,0.24),0 0 1px rgba(255,255,255,0.12) inset;overflow:hidden;transition:all 0.3s ease}
#%[1]s-wrapper .excel-container:hover{background:rgba(40,40,42,0.8);box-shadow:0 12px 48px rgba(0,0,0,0.32),0 0 1px rgba(255,255,255,0.18) inset}
#%[1]s-wrapper .excel-header{background:rgba(0,0,0,0.3);padding:16px 20px;border-bottom:1px solid rgba(255,255,255,0.08);display:flex;align-items:center;justify-content:space-between;gap:12px;flex-wrap:wrap}
#%[1]s-wrapper .excel-title{font-size:14px;font-weight:600;color:rgba(255,255,255,0.98);display:flex;align-items:center;gap:10px;text-shadow:0 1px 3px rgba(0,0,0,0.3);flex:1;min-width:200px}
#%[1]s-wrapper .excel-icon{font-size:20px;flex-shrink:0}
#%[1]s-wrapper .excel-title-text{flex:1;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
#%[1]s-wrapper .excel-toolbar{display:flex;align-items:center;gap:8px;flex-wrap:wrap}
#%[1]s-wrapper .excel-btn{background:rgba(255,255,255,0.12);border:1px solid rgba(255,255,255,0.08);color:rgba(255,255,255,0.95);padding:6px 12px;border-radius:8px;cursor:pointer;font-size:12px;font-weight:500;transition:all 0.2s;display:inline-flex;align-items:center;gap:6px;text-decoration:none;box-shadow:0 2px 6px rgba(0,0,0,0.15);white-space:nowrap}
#%[1]s-wrapper .excel-btn:hover{background:rgba(255,255,255,0.2);border-color:rgba(255,255,255,0.15);transform:translateY(-1px);box-shadow:0 3px 10px rgba(0,0,0,0.2)}
#%[1]s-wrapper .excel-btn:active{transform:translateY(0)}
#%[1]s-wrapper .excel-btn:disabled{opacity:0.3;cursor:not-allowed;transform:none}
#%[1]s-wrapper .excel-btn-primary{background:linear-gradient(135deg,#007AFF 0%%,#0051D5 100%%);border-color:transparent;box-shadow:0 3px 10px rgba(0,122,255,0.4)}
#%[1]s-wrapper .excel-btn-primary:hover{background:linear-gradient(135deg,#0A84FF 0%%,#0060EA 100%%);box-shadow:0 4px 14px rgba(0,122,255,0.5)}
#%[1]s-wrapper .excel-btn-retry{background:linear-gradient(135deg,#34C759 0%%,#30B350 100%%);border-color:transparent;color:#fff}
#%[1]s-wrapper .excel-btn-retry:hover{background:linear-gradient(135deg,#40D762 0%%,#35C755 100%%)}
#%[1]s-wrapper .excel-zoom-controls{display:flex;align-items:center;gap:6px;background:rgba(0,0,0,0.2);padding:4px;border-radius:8px}
#%[1]s-wrapper .excel-zoom-btn{width:28px;height:28px;border-radius:6px;background:rgba(255,255,255,0.1);border:1px solid rgba(255,255,255,0.06);color:rgba(255,255,255,0.8);cursor:pointer;display:flex;align-items:center;justify-content:center;transition:all 0.2s;font-size:14px;font-weight:600}
#%[1]s-wrapper .excel-zoom-btn:hover{background:rgba(255,255,255,0.2);transform:scale(1.1)}
#%[1]s-wrapper .excel-zoom-btn:disabled{opacity:0.3;cursor:not-allowed;transform:none}
#%[1]s-wrapper .excel-zoom-value{font-size:11px;color:rgba(255,255,255,0.9);min-width:50px;text-align:center;font-weight:600;font-family:'SF Mono','Monaco','Consolas',monospace}
#%[1]s-wrapper .excel-content{padding:0;min-height:600px;background:rgba(0,0,0,0.15);position:relative;overflow:hidden}
#%[1]s-wrapper .excel-loading{display:flex;align-items:center;justify-content:center;padding:100px 20px;color:rgba(255,255,255,0.95);font-size:13px;font-weight:500;flex-direction:column;position:absolute;top:50%%;left:50%%;transform:translate(-50%%,-50%%);z-index:10;text-align:center;gap:16px}
#%[1]s-wrapper .excel-loading-text{font-weight:600}
#%[1]s-wrapper .excel-loading-hint{font-size:11px;color:rgba(255,255,255,0.6);margin-top:4px}
#%[1]s-wrapper .excel-spinner{border:3px solid rgba(255,255,255,0.2);border-top:3px solid rgba(10,132,255,0.95);border-radius:50%%;width:44px;height:44px;animation:excel-spin-%[1]s 0.7s linear infinite;box-shadow:0 0 12px rgba(10,132,255,0.3)}
@keyframes excel-spin-%[1]s{0%%{transform:rotate(0deg)}100%%{transform:rotate(360deg)}}
#%[1]s-wrapper .excel-error{display:none;align-items:center;justify-content:center;padding:50px 20px;font-size:14px;font-weight:600;text-align:center;flex-direction:column;gap:16px;position:absolute;top:50%%;left:50%%;transform:translate(-50%%,-50%%);width:90%%;max-width:450px;z-index:10;background:rgba(255,255,255,0.98);border-radius:16px;box-shadow:0 12px 40px rgba(0,0,0,0.25);border:1px solid rgba(0,0,0,0.08)}
#%[1]s-wrapper .excel-error-icon{font-size:48px;margin-bottom:8px}
#%[1]s-wrapper .excel-error-title{color:rgba(231,76,60,1);font-size:16px;font-weight:700}
#%[1]s-wrapper .excel-error-message{font-size:13px;color:rgba(0,0,0,0.65);margin-top:8px;line-height:1.6;max-width:350px}
#%[1]s-wrapper .excel-error-actions{display:flex;gap:12px;margin-top:12px;flex-wrap:wrap;justify-content:center}
#%[1]s-wrapper .excel-stats{display:flex;align-items:center;gap:20px;padding:12px 20px;background:rgba(0,0,0,0.2);border-top:1px solid rgba(255,255,255,0.06);font-size:11px;color:rgba(255,255,255,0.7);flex-wrap:wrap}
#%[1]s-wrapper .excel-stat{display:flex;align-items:center;gap:6px}
#%[1]s-wrapper .excel-stat-label{font-weight:500;color:rgba(255,255,255,0.5)}
#%[1]s-wrapper .excel-stat-value{font-weight:600;color:rgba(10,132,255,1);font-family:'SF Mono','Monaco','Consolas',monospace}
#%[1]s-wrapper .excel-viewer-container{height:100%%;min-height:600px;overflow:auto;transform-origin:0 0;transition:transform 0.3s cubic-bezier(0.4,0,0.2,1)}
#%[1]s-wrapper .excel-retry-badge{position:absolute;top:16px;left:16px;background:rgba(241,196,15,0.95);color:#fff;padding:8px 14px;border-radius:8px;font-size:11px;font-weight:700;box-shadow:0 3px 12px rgba(241,196,15,0.4);z-index:5}
#%[1]s-wrapper .amis-scope{background:transparent!important}
#%[1]s-wrapper .office-viewer{background:transparent!important;border:none!important}
#%[1]s-wrapper .ov-excel{background:transparent!important}
#%[1]s-wrapper .ov-excel-header{background:rgba(10,132,255,0.18)!important;border:none!important;color:rgba(255,255,255,0.98)!important;font-weight:700!important;backdrop-filter:blur(10px);padding:12px!important;font-size:13px!important;text-shadow:0 1px 2px rgba(0,0,0,0.2)}
#%[1]s-wrapper .ov-excel-cell{border:1px solid rgba(255,255,255,0.08)!important;color:rgba(255,255,255,0.92)!important;background:rgba(0,0,0,0.15)!important;padding:10px 14px!important;font-size:13px!important;line-height:1.5!important;transition:background 0.2s ease!important}
#%[1]s-wrapper .ov-excel-row:hover .ov-excel-cell{background:rgba(10,132,255,0.12)!important}
#%[1]s-wrapper .ov-excel-sheet-tab{background:rgba(255,255,255,0.1)!important;border:1px solid rgba(255,255,255,0.08)!important;color:rgba(255,255,255,0.75)!important;border-radius:8px!important;margin-right:8px!important;padding:10px 16px!important;transition:all 0.2s cubic-bezier(0.4,0,0.2,1)!important;font-size:12px!important;font-weight:600!important}
#%[1]s-wrapper .ov-excel-sheet-tab:hover{background:rgba(255,255,255,0.15)!important;color:rgba(255,255,255,0.95)!important;transform:translateY(-2px)!important;box-shadow:0 4px 12px rgba(0,0,0,0.2)!important}
#%[1]s-wrapper .ov-excel-sheet-tab-active{background:linear-gradient(135deg,rgba(10,132,255,0.3) 0%%,rgba(0,122,255,0.4) 100%%)!important;color:rgba(10,132,255,1)!important;border-color:rgba(10,132,255,0.5)!important;font-weight:700!important;box-shadow:0 4px 16px rgba(10,132,255,0.35)!important;text-shadow:0 1px 2px rgba(0,0,0,0.1)!important}
#%[1]s-wrapper .ov-excel-sheet-tabs{background:rgba(0,0,0,0.25)!important;border:none!important;border-top:1px solid rgba(255,255,255,0.08)!important;padding:14px 20px!important;backdrop-filter:blur(10px)}
#%[1]s-wrapper .ov-excel-canvas-container{scrollbar-width:thin;scrollbar-color:rgba(255,255,255,0.25) transparent}
#%[1]s-wrapper .ov-excel-canvas-container::-webkit-scrollbar{width:12px;height:12px}
#%[1]s-wrapper .ov-excel-canvas-container::-webkit-scrollbar-track{background:rgba(0,0,0,0.15);border-radius:6px}
#%[1]s-wrapper .ov-excel-canvas-container::-webkit-scrollbar-thumb{background:rgba(255,255,255,0.25);border-radius:6px;border:2px solid rgba(0,0,0,0.15)}
#%[1]s-wrapper .ov-excel-canvas-container::-webkit-scrollbar-thumb:hover{background:rgba(255,255,255,0.35)}
@media (max-width:768px){
#%[1]s-wrapper{margin:10px}
#%[1]s-wrapper .excel-header{padding:12px 16px}
#%[1]s-wrapper .excel-title{font-size:13px}
#%[1]s-wrapper .excel-toolbar{width:100%%;justify-content:flex-start}
#%[1]s-wrapper .excel-btn{padding:5px 10px;font-size:11px}
#%[1]s-wrapper .excel-content{min-height:400px}
#%[1]s-wrapper .excel-zoom-btn{width:32px;height:32px}
#%[1]s-wrapper .excel-error{width:95%%;padding:30px 20px}
}
@media (prefers-color-scheme:light){
#%[1]s-wrapper .excel-container{background:rgba(250,250,252,0.95);border-color:rgba(0,0,0,0.08)}
#%[1]s-wrapper .excel-header{background:rgba(0,0,0,0.03);border-bottom-color:rgba(0,0,0,0.06)}
#%[1]s-wrapper .excel-title{color:rgba(0,0,0,0.95);text-shadow:0 1px 2px rgba(255,255,255,0.5)}
#%[1]s-wrapper .excel-btn{background:rgba(0,0,0,0.06);border-color:rgba(0,0,0,0.08);color:rgba(0,0,0,0.85)}
#%[1]s-wrapper .excel-content{background:rgba(0,0,0,0.02)}
#%[1]s-wrapper .excel-stats{background:rgba(0,0,0,0.03);border-top-color:rgba(0,0,0,0.06);color:rgba(0,0,0,0.7)}
#%[1]s-wrapper .excel-zoom-controls{background:rgba(0,0,0,0.05)}
#%[1]s-wrapper .excel-zoom-btn{background:rgba(0,0,0,0.05);border-color:rgba(0,0,0,0.08);color:rgba(0,0,0,0.7)}
#%[1]s-wrapper .excel-loading{color:rgba(0,0,0,0.85)}
#%[1]s-wrapper .excel-loading-hint{color:rgba(0,0,0,0.5)}
#%[1]s-wrapper .excel-spinner{border-color:rgba(0,0,0,0.1);border-top-color:rgba(0,122,255,0.9)}
#%[1]s-wrapper .ov-excel-header{background:rgba(0,122,255,0.12)!important;color:rgba(0,0,0,0.95)!important;text-shadow:none}
#%[1]s-wrapper .ov-excel-cell{background:rgba(255,255,255,0.95)!important;color:rgba(0,0,0,0.9)!important;border-color:rgba(0,0,0,0.08)!important}
#%[1]s-wrapper .ov-excel-row:hover .ov-excel-cell{background:rgba(0,122,255,0.06)!important}
#%[1]s-wrapper .ov-excel-sheet-tab{background:rgba(0,0,0,0.05)!important;color:rgba(0,0,0,0.65)!important;border-color:rgba(0,0,0,0.08)!important}
#%[1]s-wrapper .ov-excel-sheet-tabs{background:rgba(0,0,0,0.03)!important;border-top-color:rgba(0,0,0,0.08)!important}
}
</style>
<div class="excel-container">
<div class="excel-header">
<div class="excel-title">
<span class="excel-icon">📊</span>
<span class="excel-title-text">%[3]s</span>
</div>
<div class="excel-toolbar">
<div class="excel-zoom-controls">
<button id="%[1]s-zoom-out" class="excel-zoom-btn" title="Zoom Out">−</button>
<span class="excel-zoom-value" id="%[1]s-zoom-value">100%%</span>
<button id="%[1]s-zoom-in" class="excel-zoom-btn" title="Zoom In">+</button>
<button id="%[1]s-zoom-fit" class="excel-zoom-btn" title="Reset Zoom" style="font-size:10px;width:auto;padding:0 10px;font-weight:600;">Reset</button>
</div>
<button id="%[1]s-fullscreen-btn" class="excel-btn" title="Fullscreen">🔍 Fullscreen</button>
<a href="%[2]s" download class="excel-btn excel-btn-primary" title="Download">⬇ Download</a>
<a href="%[2]s" target="_blank" rel="noopener noreferrer" class="excel-btn" title="Open">🔗 Open</a>
</div>
</div>
<div class="excel-content" id="%[1]s-content">
<div class="excel-retry-badge" id="%[1]s-retry-badge" style="display:none;">Retry <span id="%[1]s-retry-count">0</span>/%[5]s</div>
<div class="excel-loading" id="%[1]s-loading">
<div class="excel-spinner"></div>
<div class="excel-loading-text">Loading Excel...</div>
<div class="excel-loading-hint">Initializing viewer</div>
</div>
<div id="%[1]s-error" class="excel-error">
<div class="excel-error-icon">📊</div>
<div class="excel-error-title">Failed to Load Excel</div>
<div class="excel-error-message" id="%[1]s-error-msg">Unable to load the Excel file. Please try again.</div>
<div class="excel-error-actions">
<button class="excel-btn excel-btn-retry" id="%[1]s-retry-btn">🔄 Retry</button>
<a href="%[2]s" download class="excel-btn">⬇ Download File</a>
<a href="%[2]s" target="_blank" class="excel-btn">🔗 Open Direct</a>
</div>
</div>
<div id="%[1]s-viewer" class="excel-viewer-container"></div>
</div>
<div class="excel-stats" id="%[1]s-stats" style="display:none;">
<div class="excel-stat">
<span class="excel-stat-label">Format:</span>
<span class="excel-stat-value">%[4]s</span>
</div>
<div class="excel-stat">
<span class="excel-stat-label">Zoom:</span>
<span class="excel-stat-value" id="%[1]s-zoom-stat">100%%</span>
</div>
<div class="excel-stat">
<span class="excel-stat-label">Sheets:</span>
<span class="excel-stat-value" id="%[1]s-sheets-stat">-</span>
</div>
<div class="excel-stat" id="%[1]s-retry-stat" style="display:none;">
<span class="excel-stat-label">Retries:</span>
<span class="excel-stat-value" id="%[1]s-retry-stat-value">0</span>
</div>
</div>
</div>
</div>
<script>
(function(){
'use strict';
const UNIQUE_ID='%[1]s';
const EXCEL_URL='%[2]s';
const FILE_NAME='%[6]s';
const MAX_RETRIES=3;
if(!window.__ExcelViewerManager){
window.__ExcelViewerManager={
instances:{},
amisLoaded:false,
amisLoading:false,
loadCallbacks:[],
keyHandlerAttached:false,
loadAmis:function(callback){
if(this.amisLoaded){
callback(true);
return;
}
this.loadCallbacks.push(callback);
if(this.amisLoading)return;
this.amisLoading=true;
const existingCSS=document.querySelector('link[href*="amis"][href*="sdk.css"]');
const existingJS=document.getElementById('amis-global-sdk');
if(existingCSS&&existingJS&&typeof amisRequire!=='undefined'){
this.amisLoaded=true;
this.amisLoading=false;
this.loadCallbacks.forEach(cb=>cb(true));
this.loadCallbacks=[];
return;
}
if(!existingCSS){
const css=document.createElement('link');
css.rel='stylesheet';
css.href='https://unpkg.com/amis@6.7.0/sdk/sdk.css';
document.head.appendChild(css);
}
if(!existingJS){
const script=document.createElement('script');
script.id='amis-global-sdk';
script.src='https://unpkg.com/amis@6.7.0/sdk/sdk.js';
script.onload=()=>{
setTimeout(()=>{
if(typeof amisRequire!=='undefined'){
this.amisLoaded=true;
this.amisLoading=false;
this.loadCallbacks.forEach(cb=>cb(true));
this.loadCallbacks=[];
}else{
this.amisLoading=false;
this.loadCallbacks.forEach(cb=>cb(false));
this.loadCallbacks=[];
}
},300);
};
script.onerror=()=>{
this.amisLoading=false;
this.loadCallbacks.forEach(cb=>cb(false));
this.loadCallbacks=[];
};
document.head.appendChild(script);
}else{
this.waitForAmis();
}
},
waitForAmis:function(){
let attempts=0;
const check=setInterval(()=>{
attempts++;
if(typeof amisRequire!=='undefined'){
clearInterval(check);
this.amisLoaded=true;
this.amisLoading=false;
this.loadCallbacks.forEach(cb=>cb(true));
this.loadCallbacks=[];
}else if(attempts>60){
clearInterval(check);
this.amisLoading=false;
this.loadCallbacks.forEach(cb=>cb(false));
this.loadCallbacks=[];
}
},100);
},
attachKeyboardHandler:function(){
if(this.keyHandlerAttached)return;
this.keyHandlerAttached=true;
document.addEventListener('keydown',(e)=>{
const activeElement=document.activeElement;
if(activeElement&&(activeElement.tagName==='INPUT'||activeElement.tagName==='TEXTAREA'||activeElement.isContentEditable))return;
for(const id in this.instances){
const inst=this.instances[id];
if(!inst||!inst.elements.wrapper)continue;
const rect=inst.elements.wrapper.getBoundingClientRect();
if(rect.top<window.innerHeight&&rect.bottom>0){
let handled=false;
switch(e.key){
case '+':
case '=':
inst.onZoomIn();
handled=true;
break;
case '-':
case '_':
inst.onZoomOut();
handled=true;
break;
case '0':
inst.onZoomFit();
handled=true;
break;
case 'f':
case 'F':
if(inst.elements.fullscreenBtn)inst.elements.fullscreenBtn.click();
handled=true;
break;
}
if(handled){
e.preventDefault();
break;
}
}
}
});
},
register:function(id,inst){
this.instances[id]=inst;
this.attachKeyboardHandler();
},
unregister:function(id){
if(this.instances[id]){
if(this.instances[id].amisInstance){
try{
const viewer=document.getElementById(id+'-viewer');
if(viewer)viewer.innerHTML='';
}catch(e){}
}
delete this.instances[id];
}
}
};
}
class ExcelViewer{
constructor(uniqueId,excelUrl,fileName){
this.uniqueId=uniqueId;
this.excelUrl=excelUrl;
this.fileName=fileName;
this.zoom=1.0;
this.minZoom=0.3;
this.maxZoom=3.0;
this.zoomStep=0.15;
this.amisInstance=null;
this.sheetsCount=0;
this.retryCount=0;
this.maxRetries=MAX_RETRIES;
this.isDestroyed=false;
this.isLoading=false;
this.eventListeners=[];
this.elements={
wrapper:document.getElementById(uniqueId+'-wrapper'),
content:document.getElementById(uniqueId+'-content'),
viewer:document.getElementById(uniqueId+'-viewer'),
loading:document.getElementById(uniqueId+'-loading'),
loadingText:null,
loadingHint:null,
error:document.getElementById(uniqueId+'-error'),
errorMsg:document.getElementById(uniqueId+'-error-msg'),
retryBtn:document.getElementById(uniqueId+'-retry-btn'),
retryBadge:document.getElementById(uniqueId+'-retry-badge'),
retryCountEl:document.getElementById(uniqueId+'-retry-count'),
stats:document.getElementById(uniqueId+'-stats'),
retryStat:document.getElementById(uniqueId+'-retry-stat'),
retryStatValue:document.getElementById(uniqueId+'-retry-stat-value'),
zoomValue:document.getElementById(uniqueId+'-zoom-value'),
zoomStat:document.getElementById(uniqueId+'-zoom-stat'),
sheetsStat:document.getElementById(uniqueId+'-sheets-stat'),
zoomInBtn:document.getElementById(uniqueId+'-zoom-in'),
zoomOutBtn:document.getElementById(uniqueId+'-zoom-out'),
zoomFitBtn:document.getElementById(uniqueId+'-zoom-fit'),
fullscreenBtn:document.getElementById(uniqueId+'-fullscreen-btn')
};
this.elements.loadingText=this.elements.loading?.querySelector('.excel-loading-text');
this.elements.loadingHint=this.elements.loading?.querySelector('.excel-loading-hint');
this.init();
}
init(){
if(!this.elements.wrapper){
console.error('Excel viewer: Required elements not found');
return;
}
window.__ExcelViewerManager.register(this.uniqueId,this);
this.bindEvents();
this.loadExcel();
}
addEventListener(element,event,handler){
if(!element)return;
element.addEventListener(event,handler);
this.eventListeners.push({element,event,handler});
}
removeAllEventListeners(){
this.eventListeners.forEach(({element,event,handler})=>{
try{element.removeEventListener(event,handler);}catch(e){}
});
this.eventListeners=[];
}
updateLoadingText(text,hint){
if(this.elements.loadingText)this.elements.loadingText.textContent=text||'Loading Excel...';
if(hint!==undefined&&this.elements.loadingHint)this.elements.loadingHint.textContent=hint||'';
}
hideLoading(){
if(this.elements.loading)this.elements.loading.style.display='none';
if(this.elements.stats)this.elements.stats.style.display='flex';
}
showError(message,canRetry=true){
this.hideLoading();
if(this.elements.error){
this.elements.error.style.display='flex';
if(this.elements.errorMsg){
this.elements.errorMsg.textContent=message||'Unable to load the Excel file. Please try again.';
}
if(this.elements.retryBtn){
this.elements.retryBtn.style.display=canRetry?'inline-flex':'none';
}
}
}
hideError(){
if(this.elements.error)this.elements.error.style.display='none';
}
updateRetryInfo(){
if(this.retryCount>0){
if(this.elements.retryBadge){
this.elements.retryBadge.style.display='block';
}
if(this.elements.retryCountEl){
this.elements.retryCountEl.textContent=this.retryCount;
}
if(this.elements.retryStat){
this.elements.retryStat.style.display='flex';
}
if(this.elements.retryStatValue){
this.elements.retryStatValue.textContent=this.retryCount+'/'+this.maxRetries;
}
}
}
updateZoom(){
if(this.elements.viewer){
this.elements.viewer.style.transform='scale('+this.zoom+')';
const zoomPercent=Math.round(this.zoom*100);
if(this.elements.zoomValue)this.elements.zoomValue.textContent=zoomPercent+'%%';
if(this.elements.zoomStat)this.elements.zoomStat.textContent=zoomPercent+'%%';
const width=100/this.zoom;
const height=100/this.zoom;
this.elements.viewer.style.width=width+'%%';
this.elements.viewer.style.height=height+'%%';
}
if(this.elements.zoomInBtn){
this.elements.zoomInBtn.disabled=this.zoom>=this.maxZoom;
}
if(this.elements.zoomOutBtn){
this.elements.zoomOutBtn.disabled=this.zoom<=this.minZoom;
}
}
onZoomIn(){
if(this.zoom<this.maxZoom){
this.zoom=Math.min(this.zoom+this.zoomStep,this.maxZoom);
this.updateZoom();
}
}
onZoomOut(){
if(this.zoom>this.minZoom){
this.zoom=Math.max(this.zoom-this.zoomStep,this.minZoom);
this.updateZoom();
}
}
onZoomFit(){
this.zoom=1.0;
this.updateZoom();
}
retry(){
if(this.isLoading)return;
if(this.retryCount>=this.maxRetries){
this.showError('Maximum retry attempts reached ('+this.maxRetries+'). The file may be corrupted or unavailable.',false);
return;
}
this.retryCount++;
this.updateRetryInfo();
this.hideError();
if(this.elements.viewer){
this.elements.viewer.innerHTML='';
}
if(this.elements.loading){
this.elements.loading.style.display='flex';
}
this.amisInstance=null;
this.zoom=1.0;
setTimeout(()=>this.loadExcel(),1000);
}
loadExcel(){
if(this.isDestroyed||this.isLoading)return;
this.isLoading=true;
const retryText=this.retryCount>0?' (Retry '+this.retryCount+'/'+this.maxRetries+')':'';
this.updateLoadingText('Loading amis library...'+retryText,'Preparing Excel viewer');
window.__ExcelViewerManager.loadAmis((success)=>{
if(this.isDestroyed)return;
if(!success){
this.isLoading=false;
if(this.retryCount<this.maxRetries){
this.showError('Failed to load Excel viewer library. Click retry to try again.',true);
}else{
this.showError('Failed to load viewer library after '+this.maxRetries+' attempts.',false);
}
return;
}
this.updateLoadingText('Initializing viewer...'+retryText,'Loading Excel file');
this.renderExcel();
});
}
renderExcel(){
if(this.isDestroyed)return;
try{
if(typeof amisRequire==='undefined'){
throw new Error('amis library not available');
}
const amis=amisRequire('amis/embed');
if(!amis||typeof amis.embed!=='function'){
throw new Error('amis.embed function not found');
}
this.amisInstance=amis.embed(
'#'+this.uniqueId+'-viewer',
{
type:'page',
body:{
type:'office-viewer',
src:this.excelUrl,
display:'normal'
}
}
);
setTimeout(()=>{
if(this.isDestroyed)return;
const sheetTabs=this.elements.wrapper.querySelectorAll('.ov-excel-sheet-tab');
if(sheetTabs.length>0){
this.sheetsCount=sheetTabs.length;
if(this.elements.sheetsStat){
this.elements.sheetsStat.textContent=this.sheetsCount;
}
}
this.hideLoading();
this.hideError();
this.updateZoom();
this.isLoading=false;
},1800);
}catch(err){
console.error('Excel render error:',err);
this.isLoading=false;
const errorMsg='Failed to render Excel: '+(err.message||'Unknown error');
if(this.retryCount<this.maxRetries){
this.showError(errorMsg+'. Click retry to try again.',true);
}else{
this.showError(errorMsg,false);
}
}
}
bindEvents(){
if(this.elements.zoomInBtn){
this.addEventListener(this.elements.zoomInBtn,'click',()=>this.onZoomIn());
}
if(this.elements.zoomOutBtn){
this.addEventListener(this.elements.zoomOutBtn,'click',()=>this.onZoomOut());
}
if(this.elements.zoomFitBtn){
this.addEventListener(this.elements.zoomFitBtn,'click',()=>this.onZoomFit());
}
if(this.elements.fullscreenBtn){
this.addEventListener(this.elements.fullscreenBtn,'click',()=>{
const target=this.elements.content;
if(!target)return;
try{
if(document.fullscreenElement){
document.exitFullscreen();
}else if(target.requestFullscreen){
target.requestFullscreen();
}else if(target.webkitRequestFullscreen){
target.webkitRequestFullscreen();
}else if(target.msRequestFullscreen){
target.msRequestFullscreen();
}
}catch(e){
console.error('Fullscreen error:',e);
}
});
}
if(this.elements.retryBtn){
this.addEventListener(this.elements.retryBtn,'click',()=>this.retry());
}
}
destroy(){
this.isDestroyed=true;
this.isLoading=false;
this.removeAllEventListeners();
if(this.elements.viewer){
this.elements.viewer.innerHTML='';
}
this.amisInstance=null;
window.__ExcelViewerManager.unregister(this.uniqueId);
}
}
try{
const viewer=new ExcelViewer(UNIQUE_ID,EXCEL_URL,FILE_NAME);
window.addEventListener('beforeunload',()=>viewer.destroy());
}catch(error){
console.error('Failed to initialize Excel viewer:',error);
}
})();
</script>`, uniqueId, safeExcelPath, fileName, strings.ToUpper(strings.TrimPrefix(ext, ".")), "3", safeFileName)
}
