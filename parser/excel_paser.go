package parser

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"log"
	"path/filepath"
)

// GenerateExcelWarpHTML 生成用于在Web上展示Excel内容的模块
// macOS Sequoia 风格，支持多实例隔离、缩放控制
func GenerateExcelWarpHTML(excelId, excelPath string) string {
	safeExcelId := template.HTMLEscapeString(excelId)
	safeExcelPath := template.HTMLEscapeString(excelPath)
	log.Println("Generating Excel viewer for FileName:", safeExcelPath)
	h := md5.Sum([]byte(safeExcelId))
	uniqueId := fmt.Sprintf("excel_%x", h)[:12]
	fileName := filepath.Base(excelPath)

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
#%[1]s-wrapper .excel-zoom-controls{display:flex;align-items:center;gap:6px;background:rgba(0,0,0,0.2);padding:4px;border-radius:8px}
#%[1]s-wrapper .excel-zoom-btn{width:28px;height:28px;border-radius:6px;background:rgba(255,255,255,0.1);border:1px solid rgba(255,255,255,0.06);color:rgba(255,255,255,0.8);cursor:pointer;display:flex;align-items:center;justify-content:center;transition:all 0.2s;font-size:12px}
#%[1]s-wrapper .excel-zoom-btn:hover{background:rgba(255,255,255,0.15);transform:scale(1.05)}
#%[1]s-wrapper .excel-zoom-value{font-size:11px;color:rgba(255,255,255,0.8);min-width:45px;text-align:center;font-weight:500}
#%[1]s-wrapper .excel-content{padding:0;min-height:600px;background:rgba(0,0,0,0.15);position:relative}
#%[1]s-wrapper .excel-loading{display:flex;align-items:center;justify-content:center;padding:100px 20px;color:rgba(255,255,255,0.95);font-size:13px;font-weight:500;flex-direction:column;position:absolute;top:50%%;left:50%%;transform:translate(-50%%,-50%%);z-index:10}
#%[1]s-wrapper .excel-spinner{border:3px solid rgba(255,255,255,0.2);border-top:3px solid rgba(10,132,255,0.95);border-radius:50%%;width:40px;height:40px;animation:excel-spin-%[1]s 0.7s linear infinite;margin-bottom:16px;box-shadow:0 0 12px rgba(10,132,255,0.3)}
@keyframes excel-spin-%[1]s{0%%{transform:rotate(0deg)}100%%{transform:rotate(360deg)}}
#%[1]s-wrapper .excel-error{display:none;align-items:center;justify-content:center;padding:60px 20px;color:rgba(255,69,58,1);font-size:13px;font-weight:600;text-align:center;flex-direction:column;gap:12px;position:absolute;top:50%%;left:50%%;transform:translate(-50%%,-50%%);width:100%%;z-index:10}
#%[1]s-wrapper .excel-stats{display:flex;align-items:center;gap:20px;padding:12px 20px;background:rgba(0,0,0,0.2);border-top:1px solid rgba(255,255,255,0.06);font-size:11px;color:rgba(255,255,255,0.7);flex-wrap:wrap}
#%[1]s-wrapper .excel-stat{display:flex;align-items:center;gap:6px}
#%[1]s-wrapper .excel-stat-label{font-weight:500;color:rgba(255,255,255,0.5)}
#%[1]s-wrapper .excel-stat-value{font-weight:600;color:rgba(10,132,255,1);font-family:'SF Mono','Monaco','Consolas',monospace}
#%[1]s-wrapper .excel-viewer-container{height:100%%;min-height:600px;overflow:auto;transform-origin:0 0;transition:transform 0.2s ease}
#%[1]s-wrapper .amis-scope{background:transparent!important}
#%[1]s-wrapper .office-viewer{background:transparent!important;border:none!important}
#%[1]s-wrapper .ov-excel{background:transparent!important}
#%[1]s-wrapper .ov-excel-header{background:rgba(10,132,255,0.15)!important;border:none!important;color:rgba(255,255,255,0.98)!important;font-weight:600!important;backdrop-filter:blur(10px);padding:10px!important;font-size:13px!important}
#%[1]s-wrapper .ov-excel-cell{border:1px solid rgba(255,255,255,0.08)!important;color:rgba(255,255,255,0.9)!important;background:rgba(0,0,0,0.12)!important;padding:10px 12px!important;font-size:13px!important;line-height:1.4!important}
#%[1]s-wrapper .ov-excel-row:hover .ov-excel-cell{background:rgba(10,132,255,0.1)!important}
#%[1]s-wrapper .ov-excel-sheet-tab{background:rgba(255,255,255,0.08)!important;border:1px solid rgba(255,255,255,0.06)!important;color:rgba(255,255,255,0.7)!important;border-radius:6px!important;margin-right:6px!important;padding:8px 14px!important;transition:all 0.2s!important;font-size:12px!important;font-weight:500!important}
#%[1]s-wrapper .ov-excel-sheet-tab:hover{background:rgba(255,255,255,0.12)!important;color:rgba(255,255,255,0.9)!important;transform:translateY(-1px)!important}
#%[1]s-wrapper .ov-excel-sheet-tab-active{background:linear-gradient(135deg,rgba(10,132,255,0.25) 0%%,rgba(0,122,255,0.35) 100%%)!important;color:rgba(10,132,255,1)!important;border-color:rgba(10,132,255,0.4)!important;font-weight:600!important;box-shadow:0 2px 8px rgba(10,132,255,0.3)!important}
#%[1]s-wrapper .ov-excel-sheet-tabs{background:rgba(0,0,0,0.2)!important;border:none!important;border-top:1px solid rgba(255,255,255,0.06)!important;padding:12px 16px!important}
#%[1]s-wrapper .ov-excel-canvas-container{scrollbar-width:thin;scrollbar-color:rgba(255,255,255,0.2) transparent}
#%[1]s-wrapper .ov-excel-canvas-container::-webkit-scrollbar{width:10px;height:10px}
#%[1]s-wrapper .ov-excel-canvas-container::-webkit-scrollbar-track{background:rgba(0,0,0,0.1);border-radius:5px}
#%[1]s-wrapper .ov-excel-canvas-container::-webkit-scrollbar-thumb{background:rgba(255,255,255,0.2);border-radius:5px;border:2px solid rgba(0,0,0,0.1)}
#%[1]s-wrapper .ov-excel-canvas-container::-webkit-scrollbar-thumb:hover{background:rgba(255,255,255,0.3)}
@media (max-width:768px){
#%[1]s-wrapper{margin:10px}
#%[1]s-wrapper .excel-header{padding:12px 16px}
#%[1]s-wrapper .excel-title{font-size:13px}
#%[1]s-wrapper .excel-toolbar{width:100%%;justify-content:flex-start}
#%[1]s-wrapper .excel-btn{padding:5px 10px;font-size:11px}
#%[1]s-wrapper .excel-content{min-height:400px}
}
@media (prefers-color-scheme:light){
#%[1]s-wrapper .excel-container{background:rgba(250,250,252,0.95);border-color:rgba(0,0,0,0.08)}
#%[1]s-wrapper .excel-header{background:rgba(0,0,0,0.03);border-bottom-color:rgba(0,0,0,0.06)}
#%[1]s-wrapper .excel-title{color:rgba(0,0,0,0.95);text-shadow:0 1px 2px rgba(255,255,255,0.5)}
#%[1]s-wrapper .excel-btn{background:rgba(0,0,0,0.06);border-color:rgba(0,0,0,0.08);color:rgba(0,0,0,0.85)}
#%[1]s-wrapper .excel-content{background:rgba(0,0,0,0.02)}
#%[1]s-wrapper .excel-stats{background:rgba(0,0,0,0.03);border-top-color:rgba(0,0,0,0.06);color:rgba(0,0,0,0.7)}
#%[1]s-wrapper .excel-zoom-controls{background:rgba(0,0,0,0.05)}
#%[1]s-wrapper .excel-zoom-btn{background:rgba(0,0,0,0.04);border-color:rgba(0,0,0,0.06);color:rgba(0,0,0,0.7)}
#%[1]s-wrapper .ov-excel-header{background:rgba(0,122,255,0.1)!important;color:rgba(0,0,0,0.95)!important}
#%[1]s-wrapper .ov-excel-cell{background:rgba(255,255,255,0.9)!important;color:rgba(0,0,0,0.9)!important;border-color:rgba(0,0,0,0.08)!important}
#%[1]s-wrapper .ov-excel-row:hover .ov-excel-cell{background:rgba(0,122,255,0.05)!important}
#%[1]s-wrapper .ov-excel-sheet-tab{background:rgba(0,0,0,0.04)!important;color:rgba(0,0,0,0.6)!important;border-color:rgba(0,0,0,0.06)!important}
#%[1]s-wrapper .ov-excel-sheet-tabs{background:rgba(0,0,0,0.03)!important;border-top-color:rgba(0,0,0,0.06)!important}
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
<button id="%[1]s-zoom-fit" class="excel-zoom-btn" title="Fit" style="font-size:10px;width:auto;padding:0 8px;">Fit</button>
</div>
<button id="%[1]s-fullscreen-btn" class="excel-btn" title="Fullscreen">🔍 Fullscreen</button>
<a href="%[2]s" download class="excel-btn excel-btn-primary" title="Download">⬇ Download</a>
<a href="%[2]s" target="_blank" rel="noopener noreferrer" class="excel-btn" title="Open">🔗 Open</a>
</div>
</div>
<div class="excel-content" id="%[1]s-content">
<div class="excel-loading" id="%[1]s-loading">
<div class="excel-spinner"></div>
<div>Loading Excel...</div>
</div>
<div id="%[1]s-error" class="excel-error">
<div>❌ Failed to load Excel</div>
<div style="font-size:11px;margin-top:8px;">Please check if the file is accessible</div>
</div>
<div id="%[1]s-viewer" class="excel-viewer-container"></div>
</div>
<div class="excel-stats" id="%[1]s-stats" style="display:none;">
<div class="excel-stat">
<span class="excel-stat-label">Format:</span>
<span class="excel-stat-value">XLSX</span>
</div>
<div class="excel-stat">
<span class="excel-stat-label">Zoom:</span>
<span class="excel-stat-value" id="%[1]s-zoom-stat">100%%</span>
</div>
<div class="excel-stat">
<span class="excel-stat-label">Sheets:</span>
<span class="excel-stat-value" id="%[1]s-sheets-stat">-</span>
</div>
</div>
</div>
</div>
<script>
(function(){
const UNIQUE_ID='%[1]s';
const EXCEL_URL='%[2]s';
const state={zoom:1.0,amisInstance:null,sheetsCount:0};
const elements={
wrapper:document.getElementById(UNIQUE_ID+'-wrapper'),
content:document.getElementById(UNIQUE_ID+'-content'),
viewer:document.getElementById(UNIQUE_ID+'-viewer'),
loading:document.getElementById(UNIQUE_ID+'-loading'),
error:document.getElementById(UNIQUE_ID+'-error'),
stats:document.getElementById(UNIQUE_ID+'-stats'),
zoomValue:document.getElementById(UNIQUE_ID+'-zoom-value'),
zoomStat:document.getElementById(UNIQUE_ID+'-zoom-stat'),
sheetsStat:document.getElementById(UNIQUE_ID+'-sheets-stat'),
zoomInBtn:document.getElementById(UNIQUE_ID+'-zoom-in'),
zoomOutBtn:document.getElementById(UNIQUE_ID+'-zoom-out'),
zoomFitBtn:document.getElementById(UNIQUE_ID+'-zoom-fit'),
fullscreenBtn:document.getElementById(UNIQUE_ID+'-fullscreen-btn')
};
if(!elements.wrapper){
console.error('Excel viewer not found');
return;
}
function hideLoading(){
if(elements.loading)elements.loading.style.display='none';
if(elements.stats)elements.stats.style.display='flex';
}
function showError(msg){
if(elements.loading)elements.loading.style.display='none';
if(elements.error){
elements.error.style.display='flex';
if(msg){
const errorDiv=elements.error.querySelector('div:last-child');
if(errorDiv)errorDiv.textContent=msg;
}
}
}
function updateZoom(){
if(elements.viewer){
elements.viewer.style.transform='scale('+state.zoom+')';
const zoomPercent=Math.round(state.zoom*100);
elements.zoomValue.textContent=zoomPercent+'%%';
elements.zoomStat.textContent=zoomPercent+'%%';
const width=100/state.zoom;
const height=100/state.zoom;
elements.viewer.style.width=width+'%%';
elements.viewer.style.height=height+'%%';
}
}
function onZoomIn(){
state.zoom=Math.min(state.zoom*1.2,3.0);
updateZoom();
}
function onZoomOut(){
state.zoom=Math.max(state.zoom/1.2,0.5);
updateZoom();
}
function onZoomFit(){
state.zoom=1.0;
updateZoom();
}
function loadAmis(callback){
if(typeof amisRequire!=='undefined'){
callback();
return;
}
if(document.getElementById('amis-sdk-script')){
let attempts=0;
const checkAmis=setInterval(()=>{
attempts++;
if(typeof amisRequire!=='undefined'){
clearInterval(checkAmis);
callback();
}else if(attempts>50){
clearInterval(checkAmis);
showError('Timeout loading amis');
}
},100);
return;
}
const css=document.createElement('link');
css.rel='stylesheet';
css.href='https://unpkg.com/amis@6.7.0/sdk/sdk.css';
document.head.appendChild(css);
const script=document.createElement('script');
script.id='amis-sdk-script';
script.src='https://unpkg.com/amis@6.7.0/sdk/sdk.js';
script.onload=()=>{
console.log('amis SDK loaded');
setTimeout(()=>{
if(typeof amisRequire!=='undefined'){
callback();
}else{
showError('amis initialization failed');
}
},200);
};
script.onerror=()=>{
console.error('Failed to load amis SDK');
showError('Failed to load library');
};
document.head.appendChild(script);
}
function renderExcel(){
try{
if(typeof amisRequire==='undefined'){
showError('amis not available');
return;
}
const amis=amisRequire('amis/embed');
if(!amis||typeof amis.embed!=='function'){
showError('amis.embed not found');
return;
}
state.amisInstance=amis.embed(
'#'+UNIQUE_ID+'-viewer',
{
type:'page',
body:{
type:'office-viewer',
src:EXCEL_URL,
display:'normal'
}
}
);
setTimeout(()=>{
hideLoading();
const sheetTabs=elements.wrapper.querySelectorAll('.ov-excel-sheet-tab');
if(sheetTabs.length>0){
state.sheetsCount=sheetTabs.length;
elements.sheetsStat.textContent=state.sheetsCount;
}
console.log('Excel viewer initialized with',state.sheetsCount,'sheets');
},1500);
}catch(err){
console.error('Render error:',err);
showError('Failed to render: '+err.message);
}
}
if(elements.zoomInBtn){
elements.zoomInBtn.addEventListener('click',onZoomIn);
}
if(elements.zoomOutBtn){
elements.zoomOutBtn.addEventListener('click',onZoomOut);
}
if(elements.zoomFitBtn){
elements.zoomFitBtn.addEventListener('click',onZoomFit);
}
if(elements.fullscreenBtn){
elements.fullscreenBtn.addEventListener('click',()=>{
const target=elements.content;
if(target.requestFullscreen)target.requestFullscreen();
else if(target.webkitRequestFullscreen)target.webkitRequestFullscreen();
else if(target.msRequestFullscreen)target.msRequestFullscreen();
});
}
document.addEventListener('keydown',(e)=>{
if(!elements.wrapper)return;
const activeElement=document.activeElement;
if(activeElement&&(activeElement.tagName==='INPUT'||activeElement.tagName==='TEXTAREA'))return;
const rect=elements.wrapper.getBoundingClientRect();
const inViewport=rect.top<window.innerHeight&&rect.bottom>0;
if(!inViewport)return;
switch(e.key){
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
loadAmis(renderExcel);
})();
</script>`, uniqueId, safeExcelPath, fileName)
}
