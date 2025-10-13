package parser

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"log"
	"path/filepath"
	"strings"
)

// GeneratePPTWarpHTML 生成用于在Web上展示PPT内容的模块
// 使用 pptx-preview (esm.sh CDN)，macOS Sequoia 风格
// 修复中文文件名问题，支持多实例和重试机制
func GeneratePPTWarpHTML(pptId, pptPath string) string {
	safePptId := template.HTMLEscapeString(pptId)
	safePptPath := template.HTMLEscapeString(pptPath)
	log.Println("Generating PPT viewer for FileName:", safePptPath)
	h := md5.Sum([]byte(safePptId))
	uniqueId := fmt.Sprintf("ppt_%x", h)[:12]
	fileName := filepath.Base(pptPath)
	ext := strings.ToLower(filepath.Ext(fileName))

	// 对文件名进行额外的 JavaScript 安全转义
	safeFileName := template.JSEscapeString(fileName)

	return fmt.Sprintf(`<div id="%[1]s-wrapper" class="ppt-viewer-wrapper" style="position:relative;max-width:100%%;margin:20px auto;display:block;">
<style>
#%[1]s-wrapper{position:relative;width:100%%;max-width:1400px;margin:20px auto;display:block;clear:both;font-family:-apple-system,BlinkMacSystemFont,'SF Pro Text','Segoe UI',Roboto,sans-serif}
#%[1]s-wrapper *{box-sizing:border-box}
#%[1]s-wrapper .ppt-container{position:relative;width:100%%;background:rgba(30,30,32,0.75);backdrop-filter:blur(40px) saturate(180%%);-webkit-backdrop-filter:blur(40px) saturate(180%%);border:1px solid rgba(255,255,255,0.12);border-radius:16px;box-shadow:0 8px 32px rgba(0,0,0,0.24),0 0 1px rgba(255,255,255,0.12) inset;overflow:hidden;transition:all 0.3s ease}
#%[1]s-wrapper .ppt-container:hover{background:rgba(40,40,42,0.8);box-shadow:0 12px 48px rgba(0,0,0,0.32),0 0 1px rgba(255,255,255,0.18) inset}
#%[1]s-wrapper .ppt-header{background:rgba(0,0,0,0.3);padding:16px 20px;border-bottom:1px solid rgba(255,255,255,0.08);display:flex;align-items:center;justify-content:space-between;gap:12px;flex-wrap:wrap}
#%[1]s-wrapper .ppt-title{font-size:14px;font-weight:600;color:rgba(255,255,255,0.98);display:flex;align-items:center;gap:10px;text-shadow:0 1px 3px rgba(0,0,0,0.3);flex:1;min-width:200px}
#%[1]s-wrapper .ppt-icon{font-size:20px;flex-shrink:0}
#%[1]s-wrapper .ppt-title-text{flex:1;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
#%[1]s-wrapper .ppt-toolbar{display:flex;align-items:center;gap:8px;flex-wrap:wrap}
#%[1]s-wrapper .ppt-btn{background:rgba(255,255,255,0.12);border:1px solid rgba(255,255,255,0.08);color:rgba(255,255,255,0.95);padding:6px 12px;border-radius:8px;cursor:pointer;font-size:12px;font-weight:500;transition:all 0.2s;display:inline-flex;align-items:center;gap:6px;text-decoration:none;box-shadow:0 2px 6px rgba(0,0,0,0.15);white-space:nowrap}
#%[1]s-wrapper .ppt-btn:hover{background:rgba(255,255,255,0.2);border-color:rgba(255,255,255,0.15);transform:translateY(-1px);box-shadow:0 3px 10px rgba(0,0,0,0.2)}
#%[1]s-wrapper .ppt-btn:active{transform:translateY(0)}
#%[1]s-wrapper .ppt-btn:disabled{opacity:0.3;cursor:not-allowed;transform:none}
#%[1]s-wrapper .ppt-btn-primary{background:linear-gradient(135deg,#007AFF 0%%,#0051D5 100%%);border-color:transparent;box-shadow:0 3px 10px rgba(0,122,255,0.4)}
#%[1]s-wrapper .ppt-btn-primary:hover{background:linear-gradient(135deg,#0A84FF 0%%,#0060EA 100%%);box-shadow:0 4px 14px rgba(0,122,255,0.5)}
#%[1]s-wrapper .ppt-btn-retry{background:linear-gradient(135deg,#34C759 0%%,#30B350 100%%);border-color:transparent;color:#fff}
#%[1]s-wrapper .ppt-btn-retry:hover{background:linear-gradient(135deg,#40D762 0%%,#35C755 100%%)}
#%[1]s-wrapper .ppt-controls{display:flex;align-items:center;gap:8px}
#%[1]s-wrapper .ppt-page-info{display:flex;align-items:center;gap:8px;padding:4px 12px;background:rgba(0,0,0,0.2);border-radius:8px;font-size:12px;color:rgba(255,255,255,0.9);font-weight:500;font-family:'SF Mono','Monaco','Consolas',monospace}
#%[1]s-wrapper .ppt-nav-btn{width:32px;height:32px;flex-shrink:0;border-radius:50%%;background:rgba(255,255,255,0.12);border:1px solid rgba(255,255,255,0.08);color:rgba(255,255,255,0.9);cursor:pointer;display:flex;align-items:center;justify-content:center;transition:all 0.2s;font-size:14px}
#%[1]s-wrapper .ppt-nav-btn:hover:not(:disabled){background:rgba(255,255,255,0.2);transform:scale(1.05)}
#%[1]s-wrapper .ppt-nav-btn:disabled{opacity:0.3;cursor:not-allowed}
#%[1]s-wrapper .ppt-content{padding:30px 20px;min-height:650px;background:rgba(0,0,0,0.15);display:flex;flex-direction:column;align-items:center;justify-content:flex-start;position:relative;overflow:auto}
#%[1]s-wrapper #%[1]s-slide-viewer{width:100%%;max-width:1200px;margin:0 auto;background:#fff;border-radius:12px;box-shadow:0 12px 48px rgba(0,0,0,0.3);overflow:visible;padding:20px}
#%[1]s-wrapper #%[1]s-slide-viewer>*{width:100%%!important;height:auto!important}
#%[1]s-wrapper #%[1]s-slide-viewer canvas,#%[1]s-wrapper #%[1]s-slide-viewer svg,#%[1]s-wrapper #%[1]s-slide-viewer img{display:block!important;max-width:100%%!important;width:100%%!important;height:auto!important;margin:0 auto 20px!important;box-shadow:0 4px 16px rgba(0,0,0,0.15)!important;border-radius:8px!important}
#%[1]s-wrapper .ppt-loading{display:flex;align-items:center;justify-content:center;padding:100px 20px;color:rgba(255,255,255,0.95);font-size:13px;font-weight:500;flex-direction:column;position:absolute;top:50%%;left:50%%;transform:translate(-50%%,-50%%);z-index:10}
#%[1]s-wrapper .ppt-spinner{border:3px solid rgba(255,255,255,0.2);border-top:3px solid rgba(10,132,255,0.95);border-radius:50%%;width:40px;height:40px;animation:ppt-spin-%[1]s 0.7s linear infinite;margin-bottom:16px;box-shadow:0 0 12px rgba(10,132,255,0.3)}
@keyframes ppt-spin-%[1]s{0%%{transform:rotate(0deg)}100%%{transform:rotate(360deg)}}
#%[1]s-wrapper .ppt-error{display:none;align-items:center;justify-content:center;padding:40px 20px;color:rgba(255,69,58,1);font-size:13px;font-weight:600;text-align:center;flex-direction:column;gap:12px;position:absolute;top:50%%;left:50%%;transform:translate(-50%%,-50%%);width:90%%;max-width:400px;z-index:10;background:rgba(255,255,255,0.95);border-radius:12px;box-shadow:0 8px 24px rgba(0,0,0,0.2)}
#%[1]s-wrapper .ppt-error-msg{font-size:12px;color:rgba(0,0,0,0.6);margin-top:8px}
#%[1]s-wrapper .ppt-error-actions{display:flex;gap:10px;margin-top:12px}
#%[1]s-wrapper .ppt-stats{display:flex;align-items:center;gap:20px;padding:12px 20px;background:rgba(0,0,0,0.2);border-top:1px solid rgba(255,255,255,0.06);font-size:11px;color:rgba(255,255,255,0.7);flex-wrap:wrap}
#%[1]s-wrapper .ppt-stat{display:flex;align-items:center;gap:6px}
#%[1]s-wrapper .ppt-stat-label{font-weight:500;color:rgba(255,255,255,0.5)}
#%[1]s-wrapper .ppt-stat-value{font-weight:600;color:rgba(10,132,255,1);font-family:'SF Mono','Monaco','Consolas',monospace}
#%[1]s-wrapper .ppt-slide-wrapper{position:relative;width:100%%;margin-bottom:20px}
#%[1]s-wrapper .ppt-slide-number{position:absolute;top:10px;right:10px;background:rgba(0,0,0,0.5);color:rgba(255,255,255,0.9);padding:4px 8px;border-radius:4px;font-size:11px;font-weight:500;z-index:1}
#%[1]s-wrapper .ppt-view-mode{display:flex;gap:8px;align-items:center}
#%[1]s-wrapper .ppt-view-btn{padding:4px 8px;border-radius:6px;background:rgba(255,255,255,0.08);border:1px solid rgba(255,255,255,0.06);color:rgba(255,255,255,0.7);font-size:11px;cursor:pointer;transition:all 0.2s}
#%[1]s-wrapper .ppt-view-btn.active{background:rgba(10,132,255,0.2);border-color:rgba(10,132,255,0.3);color:rgba(10,132,255,1)}
#%[1]s-wrapper .ppt-view-btn:hover{background:rgba(255,255,255,0.12);color:rgba(255,255,255,0.9)}
@media (max-width:768px){
#%[1]s-wrapper{margin:10px}
#%[1]s-wrapper .ppt-header{padding:12px 16px}
#%[1]s-wrapper .ppt-title{font-size:13px}
#%[1]s-wrapper .ppt-toolbar{width:100%%;justify-content:flex-start}
#%[1]s-wrapper .ppt-btn{padding:5px 10px;font-size:11px}
#%[1]s-wrapper .ppt-content{min-height:400px;padding:15px 10px}
#%[1]s-wrapper #%[1]s-slide-viewer{padding:10px}
}
@media (prefers-color-scheme:light){
#%[1]s-wrapper .ppt-container{background:rgba(250,250,252,0.95);border-color:rgba(0,0,0,0.08)}
#%[1]s-wrapper .ppt-header{background:rgba(0,0,0,0.03);border-bottom-color:rgba(0,0,0,0.06)}
#%[1]s-wrapper .ppt-title{color:rgba(0,0,0,0.95);text-shadow:0 1px 2px rgba(255,255,255,0.5)}
#%[1]s-wrapper .ppt-btn{background:rgba(0,0,0,0.06);border-color:rgba(0,0,0,0.08);color:rgba(0,0,0,0.85)}
#%[1]s-wrapper .ppt-content{background:rgba(0,0,0,0.02)}
#%[1]s-wrapper .ppt-page-info{background:rgba(0,0,0,0.05);color:rgba(0,0,0,0.8)}
#%[1]s-wrapper .ppt-nav-btn{background:rgba(0,0,0,0.06);border-color:rgba(0,0,0,0.08);color:rgba(0,0,0,0.8)}
#%[1]s-wrapper .ppt-stats{background:rgba(0,0,0,0.03);border-top-color:rgba(0,0,0,0.06);color:rgba(0,0,0,0.7)}
}
</style>
<div class="ppt-container">
<div class="ppt-header">
<div class="ppt-title">
<span class="ppt-icon">📊</span>
<span class="ppt-title-text">%[3]s</span>
</div>
<div class="ppt-toolbar">
<div class="ppt-view-mode">
<button class="ppt-view-btn active" id="%[1]s-single-view" title="Single Page">📄 Single</button>
<button class="ppt-view-btn" id="%[1]s-all-view" title="All Pages">📚 All</button>
</div>
<div class="ppt-controls">
<button class="ppt-nav-btn" id="%[1]s-prev-btn" title="Previous Slide" disabled>◀</button>
<div class="ppt-page-info">
<span id="%[1]s-current-page">1</span>
<span>/</span>
<span id="%[1]s-total-pages">-</span>
</div>
<button class="ppt-nav-btn" id="%[1]s-next-btn" title="Next Slide" disabled>▶</button>
</div>
<button class="ppt-btn" id="%[1]s-fullscreen-btn" title="Fullscreen">🔍 Fullscreen</button>
<a href="%[2]s" download class="ppt-btn ppt-btn-primary" title="Download">⬇ Download</a>
<a href="%[2]s" target="_blank" rel="noopener noreferrer" class="ppt-btn" title="Open">🔗 Open</a>
</div>
</div>
<div class="ppt-content" id="%[1]s-content">
<div class="ppt-loading" id="%[1]s-loading">
<div class="ppt-spinner"></div>
<div>Loading PowerPoint...</div>
</div>
<div class="ppt-error" id="%[1]s-error">
<div>❌ Failed to load PowerPoint</div>
<div class="ppt-error-msg" id="%[1]s-error-msg">Please check if the file format is correct</div>
<div class="ppt-error-actions">
<button class="ppt-btn ppt-btn-retry" id="%[1]s-retry-btn">🔄 Retry</button>
<a href="%[2]s" download class="ppt-btn">⬇ Download</a>
</div>
</div>
<div id="%[1]s-slide-viewer" style="display:none;"></div>
</div>
<div class="ppt-stats" id="%[1]s-stats" style="display:none;">
<div class="ppt-stat">
<span class="ppt-stat-label">Total Slides:</span>
<span class="ppt-stat-value" id="%[1]s-total-stat">0</span>
</div>
<div class="ppt-stat">
<span class="ppt-stat-label">Current:</span>
<span class="ppt-stat-value" id="%[1]s-current-stat">1</span>
</div>
<div class="ppt-stat">
<span class="ppt-stat-label">Format:</span>
<span class="ppt-stat-value">%[4]s</span>
</div>
<div class="ppt-stat" id="%[1]s-retry-stat" style="display:none;">
<span class="ppt-stat-label">Retry:</span>
<span class="ppt-stat-value" id="%[1]s-retry-count">0</span>
</div>
</div>
</div>
</div>
<script type="module">
import {init} from 'https://esm.sh/pptx-preview@1.0.5';
(function(){
'use strict';
const UNIQUE_ID='%[1]s';
const PPT_URL='%[2]s';
const FILE_NAME='%[5]s';
if(!window.__PPTViewerManager){
window.__PPTViewerManager={
instances:{},
keyHandlerAttached:false,
attachKeyboardHandler:function(){
if(this.keyHandlerAttached)return;
this.keyHandlerAttached=true;
document.addEventListener('keydown',(e)=>{
const activeElement=document.activeElement;
if(activeElement&&(activeElement.tagName==='INPUT'||activeElement.tagName==='TEXTAREA'||activeElement.isContentEditable))return;
for(const id in this.instances){
const inst=this.instances[id];
if(!inst||!inst.elements.wrapper||inst.viewMode!=='single')continue;
const rect=inst.elements.wrapper.getBoundingClientRect();
if(rect.top<window.innerHeight&&rect.bottom>0){
let handled=false;
switch(e.key){
case 'ArrowLeft':
case 'PageUp':
if(inst.currentSlide>1)inst.showSlide(inst.currentSlide-1);
handled=true;
break;
case 'ArrowRight':
case 'PageDown':
case ' ':
if(inst.currentSlide<inst.totalSlides)inst.showSlide(inst.currentSlide+1);
handled=true;
break;
case 'Home':
inst.showSlide(1);
handled=true;
break;
case 'End':
inst.showSlide(inst.totalSlides);
handled=true;
break;
case 'f':
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
delete this.instances[id];
}
};
}
const viewer={
uniqueId:UNIQUE_ID,
pptUrl:PPT_URL,
fileName:FILE_NAME,
pptxPreviewer:null,
currentSlide:1,
totalSlides:0,
viewMode:'single',
slideWrappers:[],
retryCount:0,
maxRetries:3,
elements:{
wrapper:document.getElementById(UNIQUE_ID+'-wrapper'),
content:document.getElementById(UNIQUE_ID+'-content'),
loading:document.getElementById(UNIQUE_ID+'-loading'),
error:document.getElementById(UNIQUE_ID+'-error'),
errorMsg:document.getElementById(UNIQUE_ID+'-error-msg'),
retryBtn:document.getElementById(UNIQUE_ID+'-retry-btn'),
slideViewer:document.getElementById(UNIQUE_ID+'-slide-viewer'),
stats:document.getElementById(UNIQUE_ID+'-stats'),
currentPage:document.getElementById(UNIQUE_ID+'-current-page'),
totalPages:document.getElementById(UNIQUE_ID+'-total-pages'),
currentStat:document.getElementById(UNIQUE_ID+'-current-stat'),
totalStat:document.getElementById(UNIQUE_ID+'-total-stat'),
retryStat:document.getElementById(UNIQUE_ID+'-retry-stat'),
retryCountEl:document.getElementById(UNIQUE_ID+'-retry-count'),
prevBtn:document.getElementById(UNIQUE_ID+'-prev-btn'),
nextBtn:document.getElementById(UNIQUE_ID+'-next-btn'),
fullscreenBtn:document.getElementById(UNIQUE_ID+'-fullscreen-btn'),
singleViewBtn:document.getElementById(UNIQUE_ID+'-single-view'),
allViewBtn:document.getElementById(UNIQUE_ID+'-all-view')
},
hideLoading:function(){
if(this.elements.loading)this.elements.loading.style.display='none';
if(this.elements.slideViewer)this.elements.slideViewer.style.display='block';
if(this.elements.stats)this.elements.stats.style.display='flex';
},
showError:function(msg,canRetry){
if(this.elements.loading)this.elements.loading.style.display='none';
if(this.elements.error){
this.elements.error.style.display='flex';
if(this.elements.errorMsg)this.elements.errorMsg.textContent=msg||'Please check if the file format is correct';
if(this.elements.retryBtn)this.elements.retryBtn.style.display=canRetry?'inline-flex':'none';
}
},
hideError:function(){
if(this.elements.error)this.elements.error.style.display='none';
},
updateRetryInfo:function(){
if(this.retryCount>0){
if(this.elements.retryStat)this.elements.retryStat.style.display='flex';
if(this.elements.retryCountEl)this.elements.retryCountEl.textContent=this.retryCount+'/'+this.maxRetries;
}
},
updatePageInfo:function(){
if(this.elements.currentPage)this.elements.currentPage.textContent=this.currentSlide;
if(this.elements.totalPages)this.elements.totalPages.textContent=this.totalSlides||'-';
if(this.elements.currentStat)this.elements.currentStat.textContent=this.currentSlide;
if(this.elements.totalStat)this.elements.totalStat.textContent=this.totalSlides||0;
if(this.viewMode==='single'){
if(this.elements.prevBtn)this.elements.prevBtn.disabled=this.currentSlide<=1;
if(this.elements.nextBtn)this.elements.nextBtn.disabled=this.currentSlide>=this.totalSlides;
}else{
if(this.elements.prevBtn)this.elements.prevBtn.disabled=true;
if(this.elements.nextBtn)this.elements.nextBtn.disabled=true;
}
},
showSlide:function(slideNum){
if(slideNum<1||slideNum>this.totalSlides)return;
this.currentSlide=slideNum;
if(this.viewMode==='single'){
this.slideWrappers.forEach((wrapper,index)=>{
wrapper.style.display=(index===slideNum-1)?'block':'none';
});
}
this.updatePageInfo();
},
switchViewMode:function(mode){
this.viewMode=mode;
if(mode==='single'){
this.elements.singleViewBtn.classList.add('active');
this.elements.allViewBtn.classList.remove('active');
this.slideWrappers.forEach((wrapper,index)=>{
wrapper.style.display=(index===this.currentSlide-1)?'block':'none';
});
}else{
this.elements.allViewBtn.classList.add('active');
this.elements.singleViewBtn.classList.remove('active');
this.slideWrappers.forEach(wrapper=>{
wrapper.style.display='block';
});
}
this.updatePageInfo();
},
retry:function(){
if(this.retryCount>=this.maxRetries){
this.showError('Maximum retry attempts reached. Please download the file instead.',false);
return;
}
this.retryCount++;
this.updateRetryInfo();
this.hideError();
this.slideWrappers=[];
if(this.elements.slideViewer){
this.elements.slideViewer.innerHTML='';
this.elements.slideViewer.style.display='none';
}
if(this.elements.loading){
this.elements.loading.style.display='flex';
const loadingText=this.elements.loading.querySelector('div:last-child');
if(loadingText)loadingText.textContent='Retrying... ('+this.retryCount+'/'+this.maxRetries+')';
}
setTimeout(()=>this.renderPPT(),1000);
},
renderPPT:async function(){
try{
const containerWidth=this.elements.content.offsetWidth-40;
const maxWidth=Math.min(containerWidth,1200);
this.pptxPreviewer=init(this.elements.slideViewer,{
width:maxWidth,
height:Math.round(maxWidth*9/16)
});
const response=await fetch(this.pptUrl);
if(!response.ok){
throw new Error('HTTP '+response.status);
}
const arrayBuffer=await response.arrayBuffer();
await this.pptxPreviewer.preview(arrayBuffer);
let attempts=0;
const maxAttempts=20;
const checkSlides=setInterval(()=>{
attempts++;
const slides=this.elements.slideViewer.querySelectorAll('canvas');
const svgSlides=this.elements.slideViewer.querySelectorAll('svg');
const imgSlides=this.elements.slideViewer.querySelectorAll('img');
let foundSlides=[];
if(slides.length>0){
foundSlides=Array.from(slides);
}else if(svgSlides.length>0){
foundSlides=Array.from(svgSlides);
}else if(imgSlides.length>0){
foundSlides=Array.from(imgSlides);
}
if(foundSlides.length>0||attempts>=maxAttempts){
clearInterval(checkSlides);
if(foundSlides.length===0){
const allChildren=this.elements.slideViewer.children;
if(allChildren.length>0){
foundSlides=Array.from(allChildren);
}else{
this.showError('No slides found. File may be corrupted.',this.retryCount<this.maxRetries);
return;
}
}
this.totalSlides=foundSlides.length;
foundSlides.forEach((slide,index)=>{
const wrapper=document.createElement('div');
wrapper.className='ppt-slide-wrapper';
wrapper.style.position='relative';
wrapper.style.marginBottom='20px';
wrapper.style.display=index===0?'block':'none';
slide.style.maxWidth='100%%';
slide.style.width='100%%';
slide.style.height='auto';
slide.style.display='block';
if(slide.parentNode===this.elements.slideViewer){
this.elements.slideViewer.insertBefore(wrapper,slide);
wrapper.appendChild(slide);
}else{
const clonedSlide=slide.cloneNode(true);
clonedSlide.style.maxWidth='100%%';
clonedSlide.style.width='100%%';
clonedSlide.style.height='auto';
clonedSlide.style.display='block';
wrapper.appendChild(clonedSlide);
this.elements.slideViewer.appendChild(wrapper);
}
const slideNumber=document.createElement('div');
slideNumber.className='ppt-slide-number';
slideNumber.textContent='Slide '+(index+1);
wrapper.appendChild(slideNumber);
this.slideWrappers.push(wrapper);
});
this.currentSlide=1;
this.updatePageInfo();
this.hideLoading();
}
},200);
}catch(err){
console.error('PPT loading error:',err);
this.showError('Failed to load: '+err.message,this.retryCount<this.maxRetries);
}
},
init:function(){
if(!this.elements.wrapper||!this.elements.slideViewer){
console.error('PPT viewer elements not found');
return;
}
window.__PPTViewerManager.register(this.uniqueId,this);
if(this.elements.prevBtn){
this.elements.prevBtn.addEventListener('click',()=>{
if(this.viewMode==='single'&&this.currentSlide>1){
this.showSlide(this.currentSlide-1);
}
});
}
if(this.elements.nextBtn){
this.elements.nextBtn.addEventListener('click',()=>{
if(this.viewMode==='single'&&this.currentSlide<this.totalSlides){
this.showSlide(this.currentSlide+1);
}
});
}
if(this.elements.singleViewBtn){
this.elements.singleViewBtn.addEventListener('click',()=>{
this.switchViewMode('single');
});
}
if(this.elements.allViewBtn){
this.elements.allViewBtn.addEventListener('click',()=>{
this.switchViewMode('all');
});
}
if(this.elements.fullscreenBtn){
this.elements.fullscreenBtn.addEventListener('click',()=>{
const target=this.elements.slideViewer;
if(target){
if(target.requestFullscreen){
target.requestFullscreen();
}else if(target.webkitRequestFullscreen){
target.webkitRequestFullscreen();
}else if(target.msRequestFullscreen){
target.msRequestFullscreen();
}
}
});
}
if(this.elements.retryBtn){
this.elements.retryBtn.addEventListener('click',()=>this.retry());
}
this.renderPPT();
},
destroy:function(){
window.__PPTViewerManager.unregister(this.uniqueId);
}
};
viewer.init();
window.addEventListener('beforeunload',()=>viewer.destroy());
})();
</script>`, uniqueId, safePptPath, fileName, strings.ToUpper(strings.TrimPrefix(ext, ".")), safeFileName)
}
