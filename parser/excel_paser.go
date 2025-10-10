package parser

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"log"
	"path/filepath"
)

// GenerateExcelWarpHTML 生成用于在Web上展示Excel内容的模块
// 使用 amis office-viewer 组件，macOS Sequoia 无边框风格
// 修复加载问题，使用正确的 CDN 和初始化方式
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
#%[1]s-wrapper .excel-container{position:relative;width:100%%;background:rgba(30,30,32,0.75);backdrop-filter:blur(40px) saturate(180%%);-webkit-backdrop-filter:blur(40px) saturate(180%%);border:1px solid rgba(255,255,255,0.12);border-radius:16px;box-shadow:0 8px 32px rgba(0,0,0,0.24),0 0 1px rgba(255,255,255,0.12) inset;overflow:hidden;transition:all 0.3s ease}
#%[1]s-wrapper .excel-container:hover{background:rgba(40,40,42,0.8);box-shadow:0 12px 48px rgba(0,0,0,0.32),0 0 1px rgba(255,255,255,0.18) inset}
#%[1]s-wrapper .excel-header{background:rgba(0,0,0,0.3);padding:16px 20px;border-bottom:1px solid rgba(255,255,255,0.08);display:flex;align-items:center;justify-content:space-between;gap:12px;flex-wrap:wrap}
#%[1]s-wrapper .excel-title{font-size:14px;font-weight:600;color:rgba(255,255,255,0.98);display:flex;align-items:center;gap:10px;text-shadow:0 1px 3px rgba(0,0,0,0.3)}
#%[1]s-wrapper .excel-icon{font-size:20px}
#%[1]s-wrapper .excel-toolbar{display:flex;align-items:center;gap:8px}
#%[1]s-wrapper .excel-btn{background:rgba(255,255,255,0.12);border:1px solid rgba(255,255,255,0.08);color:rgba(255,255,255,0.95);padding:6px 12px;border-radius:8px;cursor:pointer;font-size:12px;font-weight:500;transition:all 0.2s;display:inline-flex;align-items:center;gap:6px;text-decoration:none;box-shadow:0 2px 6px rgba(0,0,0,0.15)}
#%[1]s-wrapper .excel-btn:hover{background:rgba(255,255,255,0.2);border-color:rgba(255,255,255,0.15);transform:translateY(-1px);box-shadow:0 3px 10px rgba(0,0,0,0.2)}
#%[1]s-wrapper .excel-btn-primary{background:linear-gradient(135deg,#007AFF 0%%,#0051D5 100%%);border-color:transparent;box-shadow:0 3px 10px rgba(0,122,255,0.4)}
#%[1]s-wrapper .excel-btn-primary:hover{background:linear-gradient(135deg,#0A84FF 0%%,#0060EA 100%%);box-shadow:0 4px 14px rgba(0,122,255,0.5)}
#%[1]s-wrapper .excel-content{padding:0;min-height:500px;background:rgba(0,0,0,0.15)}
#%[1]s-wrapper .excel-loading{display:flex;align-items:center;justify-content:center;padding:80px 20px;color:rgba(255,255,255,0.95);font-size:13px;font-weight:500;flex-direction:column}
#%[1]s-wrapper .excel-spinner{border:3px solid rgba(255,255,255,0.2);border-top:3px solid rgba(10,132,255,0.95);border-radius:50%%;width:40px;height:40px;animation:excel-spin-%[1]s 0.7s linear infinite;margin-bottom:16px;box-shadow:0 0 12px rgba(10,132,255,0.3)}
@keyframes excel-spin-%[1]s{0%%{transform:rotate(0deg)}100%%{transform:rotate(360deg)}}
#%[1]s-wrapper .amis-scope{background:transparent!important}
#%[1]s-wrapper .office-viewer{background:transparent!important;border:none!important}
#%[1]s-wrapper .ov-excel{background:transparent!important}
#%[1]s-wrapper .ov-excel-header{background:rgba(10,132,255,0.12)!important;border:none!important;color:rgba(255,255,255,0.98)!important;font-weight:600!important;backdrop-filter:blur(10px)}
#%[1]s-wrapper .ov-excel-cell{border:1px solid rgba(255,255,255,0.06)!important;color:rgba(255,255,255,0.9)!important;background:rgba(0,0,0,0.1)!important;padding:8px 10px!important}
#%[1]s-wrapper .ov-excel-row:hover .ov-excel-cell{background:rgba(10,132,255,0.08)!important}
#%[1]s-wrapper .ov-excel-sheet-tab{background:rgba(255,255,255,0.08)!important;border:1px solid rgba(255,255,255,0.06)!important;color:rgba(255,255,255,0.7)!important;border-radius:6px!important;margin-right:6px!important;padding:8px 14px!important;transition:all 0.2s!important}
#%[1]s-wrapper .ov-excel-sheet-tab:hover{background:rgba(255,255,255,0.12)!important;color:rgba(255,255,255,0.9)!important}
#%[1]s-wrapper .ov-excel-sheet-tab-active{background:linear-gradient(135deg,rgba(10,132,255,0.2) 0%%,rgba(0,122,255,0.3) 100%%)!important;color:rgba(10,132,255,1)!important;border-color:rgba(10,132,255,0.3)!important;font-weight:600!important}
#%[1]s-wrapper .ov-excel-sheet-tabs{background:rgba(0,0,0,0.2)!important;border:none!important;border-bottom:1px solid rgba(255,255,255,0.06)!important;padding:8px 16px!important}
#%[1]s-wrapper .ov-excel-canvas-container{scrollbar-width:thin;scrollbar-color:rgba(255,255,255,0.2) transparent}
#%[1]s-wrapper .ov-excel-canvas-container::-webkit-scrollbar{width:8px;height:8px}
#%[1]s-wrapper .ov-excel-canvas-container::-webkit-scrollbar-track{background:rgba(0,0,0,0.1);border-radius:4px}
#%[1]s-wrapper .ov-excel-canvas-container::-webkit-scrollbar-thumb{background:rgba(255,255,255,0.2);border-radius:4px}
#%[1]s-wrapper .ov-excel-canvas-container::-webkit-scrollbar-thumb:hover{background:rgba(255,255,255,0.3)}
@media (max-width:768px){
#%[1]s-wrapper{margin:10px}
#%[1]s-wrapper .excel-header{padding:12px 16px}
#%[1]s-wrapper .excel-title{font-size:13px}
#%[1]s-wrapper .excel-btn{padding:5px 10px;font-size:11px}
#%[1]s-wrapper .excel-content{min-height:400px}
}
@media (prefers-color-scheme:light){
#%[1]s-wrapper .excel-container{background:rgba(250,250,252,0.95);border-color:rgba(0,0,0,0.08)}
#%[1]s-wrapper .excel-header{background:rgba(0,0,0,0.03);border-bottom-color:rgba(0,0,0,0.06)}
#%[1]s-wrapper .excel-title{color:rgba(0,0,0,0.95);text-shadow:0 1px 2px rgba(255,255,255,0.5)}
#%[1]s-wrapper .excel-btn{background:rgba(0,0,0,0.06);border-color:rgba(0,0,0,0.08);color:rgba(0,0,0,0.85)}
#%[1]s-wrapper .excel-content{background:rgba(0,0,0,0.02)}
#%[1]s-wrapper .ov-excel-header{background:rgba(0,122,255,0.08)!important;color:rgba(0,0,0,0.95)!important}
#%[1]s-wrapper .ov-excel-cell{background:rgba(255,255,255,0.8)!important;color:rgba(0,0,0,0.9)!important;border-color:rgba(0,0,0,0.06)!important}
#%[1]s-wrapper .ov-excel-sheet-tab{background:rgba(0,0,0,0.04)!important;color:rgba(0,0,0,0.6)!important;border-color:rgba(0,0,0,0.06)!important}
#%[1]s-wrapper .ov-excel-sheet-tabs{background:rgba(0,0,0,0.03)!important;border-bottom-color:rgba(0,0,0,0.06)!important}
}
</style>
<div class="excel-container">
<div class="excel-header">
<div class="excel-title">
<span class="excel-icon">📊</span>
<span>%[3]s</span>
</div>
<div class="excel-toolbar">
<a href="%[2]s" download class="excel-btn excel-btn-primary" title="Download">⬇ Download</a>
<a href="%[2]s" target="_blank" rel="noopener noreferrer" class="excel-btn" title="Open">🔗 Open</a>
</div>
</div>
<div class="excel-content" id="%[1]s-content">
<div class="excel-loading" id="%[1]s-loading">
<div class="excel-spinner"></div>
<div>Loading Excel...</div>
</div>
</div>
</div>
</div>
<link rel="stylesheet" href="https://unpkg.com/amis@6.7.0/sdk/sdk.css">
<script src="https://unpkg.com/amis@6.7.0/sdk/sdk.js"></script>
<script>
(function(){
const UNIQUE_ID='%[1]s';
const EXCEL_URL='%[2]s';
const elements={wrapper:document.getElementById(UNIQUE_ID+'-wrapper'),content:document.getElementById(UNIQUE_ID+'-content'),loading:document.getElementById(UNIQUE_ID+'-loading')};
if(!elements.wrapper){
console.error('Excel viewer not found');
return;
}
function hideLoading(){
if(elements.loading){
elements.loading.style.display='none';
}
}
function showError(msg){
if(elements.loading){
elements.loading.innerHTML='<div style="color:rgba(255,69,58,1);">❌ '+msg+'</div>';
}
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
const amisInstance=amis.embed(
'#'+UNIQUE_ID+'-content',
{
type:'page',
body:{
type:'office-viewer',
src:EXCEL_URL,
display:'normal'
}
}
);
hideLoading();
console.log('Excel viewer initialized');
if(amisInstance&&amisInstance.getComponentByName){
console.log('amis instance created successfully');
}
}catch(err){
console.error('Render error:',err);
showError('Failed to render: '+err.message);
}
}
loadAmis(renderExcel);
})();
</script>`, uniqueId, safeExcelPath, fileName)
}
