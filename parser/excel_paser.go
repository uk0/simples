package parser

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"log"
	"path/filepath"
	"strings"
)

// GenerateExcelWarpHTML 生成用于在Web上展示Excel内容的模块
// 使用 amis + SheetJS 解析和渲染 Excel，macOS Sequoia 风格
// 支持 .xlsx, .xls, .csv 格式
func GenerateExcelWarpHTML(excelId, excelPath string) string {
	safeExcelId := template.HTMLEscapeString(excelId)
	safeExcelPath := template.HTMLEscapeString(excelPath)
	log.Println("Generating Excel viewer for FileName:", safeExcelPath)
	h := md5.Sum([]byte(safeExcelId))
	uniqueId := fmt.Sprintf("excel_%x", h)[:12]
	ext := strings.ToLower(filepath.Ext(excelPath))
	fileName := filepath.Base(excelPath)
	fileNameWithoutExt := strings.TrimSuffix(fileName, ext)

	return fmt.Sprintf(`<div id="%[1]s-wrapper" class="excel-viewer-wrapper" style="position:relative;max-width:100%%;margin:20px auto;display:block;">
<style>
#%[1]s-wrapper{position:relative;width:100%%;max-width:1200px;margin:20px auto;display:block;clear:both;font-family:-apple-system,BlinkMacSystemFont,'SF Pro Text','Segoe UI',Roboto,sans-serif}
#%[1]s-wrapper .excel-container{position:relative;width:100%%;background:rgba(30,30,32,0.75);backdrop-filter:blur(40px) saturate(180%%);-webkit-backdrop-filter:blur(40px) saturate(180%%);border:1px solid rgba(255,255,255,0.12);border-radius:16px;box-shadow:0 8px 32px rgba(0,0,0,0.24),0 0 1px rgba(255,255,255,0.12) inset;overflow:hidden;transition:all 0.3s ease}
#%[1]s-wrapper .excel-container:hover{background:rgba(40,40,42,0.8);box-shadow:0 12px 48px rgba(0,0,0,0.32),0 0 1px rgba(255,255,255,0.18) inset}
#%[1]s-wrapper .excel-header{background:rgba(0,0,0,0.3);padding:16px 20px;border-bottom:1px solid rgba(255,255,255,0.08);display:flex;align-items:center;justify-content:space-between;gap:12px}
#%[1]s-wrapper .excel-title{font-size:14px;font-weight:600;color:rgba(255,255,255,0.98);display:flex;align-items:center;gap:10px;text-shadow:0 1px 3px rgba(0,0,0,0.3)}
#%[1]s-wrapper .excel-icon{font-size:20px}
#%[1]s-wrapper .excel-toolbar{display:flex;align-items:center;gap:8px;flex-wrap:wrap}
#%[1]s-wrapper .excel-btn{background:rgba(255,255,255,0.12);border:1px solid rgba(255,255,255,0.08);color:rgba(255,255,255,0.95);padding:6px 12px;border-radius:8px;cursor:pointer;font-size:12px;font-weight:500;transition:all 0.2s;display:inline-flex;align-items:center;gap:6px;text-decoration:none;box-shadow:0 2px 6px rgba(0,0,0,0.15)}
#%[1]s-wrapper .excel-btn:hover{background:rgba(255,255,255,0.2);border-color:rgba(255,255,255,0.15);transform:translateY(-1px);box-shadow:0 3px 10px rgba(0,0,0,0.2)}
#%[1]s-wrapper .excel-btn:active{transform:translateY(0)}
#%[1]s-wrapper .excel-btn-primary{background:linear-gradient(135deg,#007AFF 0%%,#0051D5 100%%);border-color:transparent;box-shadow:0 3px 10px rgba(0,122,255,0.4)}
#%[1]s-wrapper .excel-btn-primary:hover{background:linear-gradient(135deg,#0A84FF 0%%,#0060EA 100%%);box-shadow:0 4px 14px rgba(0,122,255,0.5)}
#%[1]s-wrapper .excel-tabs{display:flex;align-items:center;gap:4px;padding:0 16px;background:rgba(0,0,0,0.2);border-bottom:1px solid rgba(255,255,255,0.06);overflow-x:auto;scrollbar-width:thin;scrollbar-color:rgba(255,255,255,0.2) transparent}
#%[1]s-wrapper .excel-tabs::-webkit-scrollbar{height:6px}
#%[1]s-wrapper .excel-tabs::-webkit-scrollbar-track{background:transparent}
#%[1]s-wrapper .excel-tabs::-webkit-scrollbar-thumb{background:rgba(255,255,255,0.2);border-radius:3px}
#%[1]s-wrapper .excel-tabs::-webkit-scrollbar-thumb:hover{background:rgba(255,255,255,0.3)}
#%[1]s-wrapper .excel-tab{padding:10px 16px;font-size:12px;font-weight:500;color:rgba(255,255,255,0.7);cursor:pointer;white-space:nowrap;border-bottom:2px solid transparent;transition:all 0.2s}
#%[1]s-wrapper .excel-tab:hover{color:rgba(255,255,255,0.9);background:rgba(255,255,255,0.05)}
#%[1]s-wrapper .excel-tab.active{color:rgba(10,132,255,1);border-bottom-color:#0A84FF;background:rgba(10,132,255,0.1)}
#%[1]s-wrapper .excel-content{padding:0;min-height:400px;max-height:600px;overflow:auto;background:rgba(0,0,0,0.15)}
#%[1]s-wrapper .excel-table-wrapper{width:100%%;overflow:auto}
#%[1]s-wrapper .excel-table{width:100%%;border-collapse:collapse;font-size:12px}
#%[1]s-wrapper .excel-table th{background:rgba(10,132,255,0.15);color:rgba(255,255,255,0.98);font-weight:600;text-align:left;padding:10px 12px;border:1px solid rgba(255,255,255,0.08);white-space:nowrap;position:sticky;top:0;z-index:10;backdrop-filter:blur(10px)}
#%[1]s-wrapper .excel-table td{padding:8px 12px;border:1px solid rgba(255,255,255,0.06);color:rgba(255,255,255,0.9);background:rgba(0,0,0,0.1);white-space:nowrap}
#%[1]s-wrapper .excel-table tr:hover td{background:rgba(10,132,255,0.08)}
#%[1]s-wrapper .excel-table td.number{text-align:right;font-family:'SF Mono','Monaco','Consolas',monospace}
#%[1]s-wrapper .excel-loading{display:flex;align-items:center;justify-content:center;padding:60px 20px;color:rgba(255,255,255,0.95);font-size:13px;font-weight:500}
#%[1]s-wrapper .excel-spinner{border:3px solid rgba(255,255,255,0.2);border-top:3px solid rgba(10,132,255,0.95);border-radius:50%%;width:32px;height:32px;animation:excel-spin-%[1]s 0.7s linear infinite;margin:0 auto 12px;box-shadow:0 0 10px rgba(10,132,255,0.3)}
@keyframes excel-spin-%[1]s{0%%{transform:rotate(0deg)}100%%{transform:rotate(360deg)}}
#%[1]s-wrapper .excel-error{display:none;align-items:center;justify-content:center;padding:60px 20px;color:rgba(255,69,58,1);font-size:13px;font-weight:600;text-align:center}
#%[1]s-wrapper .excel-empty{display:none;align-items:center;justify-content:center;padding:60px 20px;color:rgba(255,255,255,0.6);font-size:13px;text-align:center}
#%[1]s-wrapper .excel-stats{display:flex;align-items:center;gap:20px;padding:12px 20px;background:rgba(0,0,0,0.2);border-top:1px solid rgba(255,255,255,0.06);font-size:11px;color:rgba(255,255,255,0.7);flex-wrap:wrap}
#%[1]s-wrapper .excel-stat{display:flex;align-items:center;gap:6px}
#%[1]s-wrapper .excel-stat-label{font-weight:500;color:rgba(255,255,255,0.5)}
#%[1]s-wrapper .excel-stat-value{font-weight:600;color:rgba(10,132,255,1);font-family:'SF Mono','Monaco','Consolas',monospace}
@media (max-width:768px){
#%[1]s-wrapper{margin:10px}
#%[1]s-wrapper .excel-header{padding:12px 16px;flex-wrap:wrap}
#%[1]s-wrapper .excel-title{font-size:13px}
#%[1]s-wrapper .excel-toolbar{width:100%%;justify-content:flex-start}
#%[1]s-wrapper .excel-btn{padding:5px 10px;font-size:11px}
#%[1]s-wrapper .excel-table{font-size:11px}
#%[1]s-wrapper .excel-table th,#%[1]s-wrapper .excel-table td{padding:6px 8px}
}
@media (prefers-color-scheme:light){
#%[1]s-wrapper .excel-container{background:rgba(250,250,252,0.95);border-color:rgba(0,0,0,0.08)}
#%[1]s-wrapper .excel-header{background:rgba(0,0,0,0.03);border-bottom-color:rgba(0,0,0,0.06)}
#%[1]s-wrapper .excel-title{color:rgba(0,0,0,0.95);text-shadow:0 1px 2px rgba(255,255,255,0.5)}
#%[1]s-wrapper .excel-btn{background:rgba(0,0,0,0.06);border-color:rgba(0,0,0,0.08);color:rgba(0,0,0,0.85)}
#%[1]s-wrapper .excel-tabs{background:rgba(0,0,0,0.03);border-bottom-color:rgba(0,0,0,0.06)}
#%[1]s-wrapper .excel-tab{color:rgba(0,0,0,0.6)}
#%[1]s-wrapper .excel-tab:hover{color:rgba(0,0,0,0.9);background:rgba(0,0,0,0.03)}
#%[1]s-wrapper .excel-content{background:rgba(0,0,0,0.02)}
#%[1]s-wrapper .excel-table th{background:rgba(0,122,255,0.08);color:rgba(0,0,0,0.95)}
#%[1]s-wrapper .excel-table td{background:rgba(255,255,255,0.8);color:rgba(0,0,0,0.9);border-color:rgba(0,0,0,0.06)}
#%[1]s-wrapper .excel-stats{background:rgba(0,0,0,0.03);border-top-color:rgba(0,0,0,0.06);color:rgba(0,0,0,0.7)}
}
</style>
<div class="excel-container">
<div class="excel-header">
<div class="excel-title">
<span class="excel-icon">📊</span>
<span>%[3]s</span>
</div>
<div class="excel-toolbar">
<button class="excel-btn" id="%[1]s-refresh-btn" title="Refresh">🔄 Refresh</button>
<button class="excel-btn excel-btn-primary" id="%[1]s-download-btn" title="Download">⬇ Download</button>
<a href="%[2]s" target="_blank" rel="noopener noreferrer" class="excel-btn" title="Open">🔗 Open</a>
</div>
</div>
<div class="excel-tabs" id="%[1]s-tabs"></div>
<div class="excel-content" id="%[1]s-content">
<div class="excel-loading" id="%[1]s-loading">
<div><div class="excel-spinner"></div><div>Loading Excel file...</div></div>
</div>
<div class="excel-error" id="%[1]s-error">❌ Failed to load Excel file</div>
<div class="excel-empty" id="%[1]s-empty">📭 No data available</div>
<div class="excel-table-wrapper" id="%[1]s-table-wrapper" style="display:none;"></div>
</div>
<div class="excel-stats" id="%[1]s-stats" style="display:none;">
<div class="excel-stat">
<span class="excel-stat-label">Rows:</span>
<span class="excel-stat-value" id="%[1]s-rows">0</span>
</div>
<div class="excel-stat">
<span class="excel-stat-label">Columns:</span>
<span class="excel-stat-value" id="%[1]s-cols">0</span>
</div>
<div class="excel-stat">
<span class="excel-stat-label">Sheet:</span>
<span class="excel-stat-value" id="%[1]s-current-sheet">-</span>
</div>
</div>
</div>
</div>
<script src="https://cdn.jsdelivr.net/npm/xlsx@0.18.5/dist/xlsx.full.min.js"></script>
<script>
(function(){
const UNIQUE_ID='%[1]s';
const EXCEL_URL='%[2]s';
const elements={wrapper:document.getElementById(UNIQUE_ID+'-wrapper'),loading:document.getElementById(UNIQUE_ID+'-loading'),error:document.getElementById(UNIQUE_ID+'-error'),empty:document.getElementById(UNIQUE_ID+'-empty'),tabs:document.getElementById(UNIQUE_ID+'-tabs'),content:document.getElementById(UNIQUE_ID+'-content'),tableWrapper:document.getElementById(UNIQUE_ID+'-table-wrapper'),stats:document.getElementById(UNIQUE_ID+'-stats'),rowsStat:document.getElementById(UNIQUE_ID+'-rows'),colsStat:document.getElementById(UNIQUE_ID+'-cols'),sheetStat:document.getElementById(UNIQUE_ID+'-current-sheet'),refreshBtn:document.getElementById(UNIQUE_ID+'-refresh-btn'),downloadBtn:document.getElementById(UNIQUE_ID+'-download-btn')};
if(!elements.wrapper){
console.error('Excel viewer elements not found');
return;
}
let workbook=null;
let currentSheet=null;
function showLoading(){
if(elements.loading)elements.loading.style.display='flex';
if(elements.error)elements.error.style.display='none';
if(elements.empty)elements.empty.style.display='none';
if(elements.tableWrapper)elements.tableWrapper.style.display='none';
if(elements.stats)elements.stats.style.display='none';
}
function showError(msg){
if(elements.loading)elements.loading.style.display='none';
if(elements.error){
elements.error.style.display='flex';
if(msg)elements.error.innerHTML='❌ '+msg;
}
if(elements.tableWrapper)elements.tableWrapper.style.display='none';
if(elements.stats)elements.stats.style.display='none';
}
function showEmpty(){
if(elements.loading)elements.loading.style.display='none';
if(elements.error)elements.error.style.display='none';
if(elements.empty)elements.empty.style.display='flex';
if(elements.tableWrapper)elements.tableWrapper.style.display='none';
if(elements.stats)elements.stats.style.display='none';
}
function showTable(){
if(elements.loading)elements.loading.style.display='none';
if(elements.error)elements.error.style.display='none';
if(elements.empty)elements.empty.style.display='none';
if(elements.tableWrapper)elements.tableWrapper.style.display='block';
if(elements.stats)elements.stats.style.display='flex';
}
function isNumeric(val){
if(typeof val==='number')return true;
if(typeof val==='string'){
return!isNaN(parseFloat(val))&&isFinite(val);
}
return false;
}
function renderSheet(sheetName){
if(!workbook||!workbook.Sheets[sheetName]){
showEmpty();
return;
}
const worksheet=workbook.Sheets[sheetName];
const jsonData=XLSX.utils.sheet_to_json(worksheet,{header:1,defval:''});
if(!jsonData||jsonData.length===0){
showEmpty();
return;
}
currentSheet=sheetName;
const headers=jsonData[0]||[];
const rows=jsonData.slice(1);
let html='<table class="excel-table"><thead><tr>';
headers.forEach((header,idx)=>{
html+='<th>'+(header||'Column '+(idx+1))+'</th>';
});
html+='</tr></thead><tbody>';
rows.forEach(row=>{
html+='<tr>';
headers.forEach((header,idx)=>{
const cell=row[idx]||'';
const className=isNumeric(cell)?'number':'';
html+='<td class="'+className+'">'+cell+'</td>';
});
html+='</tr>';
});
html+='</tbody></table>';
if(elements.tableWrapper){
elements.tableWrapper.innerHTML=html;
}
if(elements.rowsStat)elements.rowsStat.textContent=rows.length;
if(elements.colsStat)elements.colsStat.textContent=headers.length;
if(elements.sheetStat)elements.sheetStat.textContent=sheetName;
showTable();
console.log('Rendered sheet:',sheetName,'Rows:',rows.length,'Cols:',headers.length);
}
function renderTabs(){
if(!workbook||!elements.tabs)return;
const sheetNames=workbook.SheetNames;
if(!sheetNames||sheetNames.length===0)return;
let html='';
sheetNames.forEach((name,idx)=>{
const activeClass=idx===0?'active':'';
html+='<div class="excel-tab '+activeClass+'" data-sheet="'+name+'">'+name+'</div>';
});
elements.tabs.innerHTML=html;
const tabs=elements.tabs.querySelectorAll('.excel-tab');
tabs.forEach(tab=>{
tab.addEventListener('click',function(){
const sheetName=this.getAttribute('data-sheet');
tabs.forEach(t=>t.classList.remove('active'));
this.classList.add('active');
renderSheet(sheetName);
});
});
}
function loadExcel(){
showLoading();
fetch(EXCEL_URL)
.then(response=>{
if(!response.ok){
throw new Error('Failed to fetch file ('+response.status+')');
}
return response.arrayBuffer();
})
.then(data=>{
try{
workbook=XLSX.read(data,{type:'array'});
if(!workbook||!workbook.SheetNames||workbook.SheetNames.length===0){
showEmpty();
return;
}
renderTabs();
renderSheet(workbook.SheetNames[0]);
console.log('Excel loaded successfully. Sheets:',workbook.SheetNames.length);
}catch(err){
console.error('Parse error:',err);
showError('Failed to parse Excel file');
}
})
.catch(error=>{
console.error('Load error:',error);
showError('Failed to load Excel file');
});
}
function loadXLSX(callback){
if(typeof XLSX!=='undefined'){
callback();
return;
}
if(document.getElementById('xlsx-script')){
const checkXLSX=setInterval(()=>{
if(typeof XLSX!=='undefined'){
clearInterval(checkXLSX);
callback();
}
},100);
return;
}
const script=document.createElement('script');
script.id='xlsx-script';
script.src='https://cdn.jsdelivr.net/npm/xlsx@0.18.5/dist/xlsx.full.min.js';
script.onload=()=>{
console.log('XLSX.js loaded');
callback();
};
script.onerror=()=>{
console.error('Failed to load XLSX.js');
showError('Failed to load library');
};
document.head.appendChild(script);
}
if(elements.refreshBtn){
elements.refreshBtn.addEventListener('click',()=>{
loadXLSX(loadExcel);
});
}
if(elements.downloadBtn){
elements.downloadBtn.addEventListener('click',()=>{
window.open(EXCEL_URL,'_blank');
});
}
loadXLSX(loadExcel);
})();
</script>`, uniqueId, safeExcelPath, fileNameWithoutExt)
}
