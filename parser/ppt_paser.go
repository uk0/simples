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
// 修复无法找到幻灯片的问题
func GeneratePPTWarpHTML(pptId, pptPath string) string {
	safePptId := template.HTMLEscapeString(pptId)
	safePptPath := template.HTMLEscapeString(pptPath)
	log.Println("Generating PPT viewer for FileName:", safePptPath)
	h := md5.Sum([]byte(safePptId))
	uniqueId := fmt.Sprintf("ppt_%x", h)[:12]
	fileName := filepath.Base(pptPath)
	ext := strings.ToLower(filepath.Ext(fileName))

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
#%[1]s-wrapper .ppt-error{display:none;align-items:center;justify-content:center;padding:60px 20px;color:rgba(255,69,58,1);font-size:13px;font-weight:600;text-align:center;flex-direction:column;gap:12px;position:absolute;top:50%%;left:50%%;transform:translate(-50%%,-50%%);width:100%%;z-index:10}
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
<div style="font-size:11px;margin-top:8px;">Please check if the file format is correct</div>
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
</div>
</div>
</div>
<script type="module">
import {init} from 'https://esm.sh/pptx-preview@1.0.5';
(function(){
const UNIQUE_ID='%[1]s';
const PPT_URL='%[2]s';
const elements={wrapper:document.getElementById(UNIQUE_ID+'-wrapper'),content:document.getElementById(UNIQUE_ID+'-content'),loading:document.getElementById(UNIQUE_ID+'-loading'),error:document.getElementById(UNIQUE_ID+'-error'),slideViewer:document.getElementById(UNIQUE_ID+'-slide-viewer'),stats:document.getElementById(UNIQUE_ID+'-stats'),currentPage:document.getElementById(UNIQUE_ID+'-current-page'),totalPages:document.getElementById(UNIQUE_ID+'-total-pages'),currentStat:document.getElementById(UNIQUE_ID+'-current-stat'),totalStat:document.getElementById(UNIQUE_ID+'-total-stat'),prevBtn:document.getElementById(UNIQUE_ID+'-prev-btn'),nextBtn:document.getElementById(UNIQUE_ID+'-next-btn'),fullscreenBtn:document.getElementById(UNIQUE_ID+'-fullscreen-btn'),singleViewBtn:document.getElementById(UNIQUE_ID+'-single-view'),allViewBtn:document.getElementById(UNIQUE_ID+'-all-view')};
if(!elements.wrapper||!elements.slideViewer){
console.error('PPT viewer elements not found');
return;
}
let pptxPreviewer=null;
let currentSlide=1;
let totalSlides=0;
let viewMode='single';
let slideWrappers=[];
function hideLoading(){
if(elements.loading)elements.loading.style.display='none';
if(elements.slideViewer)elements.slideViewer.style.display='block';
if(elements.stats)elements.stats.style.display='flex';
}
function showError(msg){
if(elements.loading)elements.loading.style.display='none';
if(elements.error){
elements.error.style.display='flex';
if(msg){
const errorDiv=elements.error.querySelector('div');
if(errorDiv)errorDiv.textContent='❌ '+msg;
}
}
}
function updatePageInfo(){
if(elements.currentPage)elements.currentPage.textContent=currentSlide;
if(elements.totalPages)elements.totalPages.textContent=totalSlides||'-';
if(elements.currentStat)elements.currentStat.textContent=currentSlide;
if(elements.totalStat)elements.totalStat.textContent=totalSlides||0;
if(viewMode==='single'){
if(elements.prevBtn)elements.prevBtn.disabled=currentSlide<=1;
if(elements.nextBtn)elements.nextBtn.disabled=currentSlide>=totalSlides;
}else{
if(elements.prevBtn)elements.prevBtn.disabled=true;
if(elements.nextBtn)elements.nextBtn.disabled=true;
}
}
function showSlide(slideNum){
if(slideNum<1||slideNum>totalSlides)return;
currentSlide=slideNum;
if(viewMode==='single'){
slideWrappers.forEach((wrapper,index)=>{
if(index===slideNum-1){
wrapper.style.display='block';
}else{
wrapper.style.display='none';
}
});
}
updatePageInfo();
console.log('Showing slide:',currentSlide+'/'+totalSlides);
}
function switchViewMode(mode){
viewMode=mode;
if(mode==='single'){
elements.singleViewBtn.classList.add('active');
elements.allViewBtn.classList.remove('active');
slideWrappers.forEach((wrapper,index)=>{
if(index===currentSlide-1){
wrapper.style.display='block';
}else{
wrapper.style.display='none';
}
});
}else{
elements.allViewBtn.classList.add('active');
elements.singleViewBtn.classList.remove('active');
slideWrappers.forEach(wrapper=>{
wrapper.style.display='block';
});
}
updatePageInfo();
}
async function renderPPT(){
try{
const containerWidth=elements.content.offsetWidth-40;
const maxWidth=Math.min(containerWidth,1200);
console.log('Container width:',maxWidth);
pptxPreviewer=init(elements.slideViewer,{
width:maxWidth,
height:Math.round(maxWidth*9/16)
});
console.log('PPT previewer initialized');
const response=await fetch(PPT_URL);
if(!response.ok){
throw new Error('HTTP '+response.status);
}
const arrayBuffer=await response.arrayBuffer();
console.log('File loaded, size:',arrayBuffer.byteLength);
await pptxPreviewer.preview(arrayBuffer);
console.log('Preview rendered, waiting for DOM...');
let attempts=0;
const maxAttempts=20;
const checkSlides=setInterval(()=>{
attempts++;
console.log('Checking for slides, attempt:',attempts);
const allElements=elements.slideViewer.querySelectorAll('canvas, svg, img, .slide, [class*="slide"], [id*="slide"]');
console.log('Found elements:',allElements.length);
if(allElements.length>0){
console.log('Found all elements:',allElements);
}
const slides=elements.slideViewer.querySelectorAll('canvas');
const svgSlides=elements.slideViewer.querySelectorAll('svg');
const imgSlides=elements.slideViewer.querySelectorAll('img');
let foundSlides=[];
if(slides.length>0){
foundSlides=Array.from(slides);
console.log('Found canvas slides:',slides.length);
}else if(svgSlides.length>0){
foundSlides=Array.from(svgSlides);
console.log('Found SVG slides:',svgSlides.length);
}else if(imgSlides.length>0){
foundSlides=Array.from(imgSlides);
console.log('Found image slides:',imgSlides.length);
}
if(foundSlides.length>0||attempts>=maxAttempts){
clearInterval(checkSlides);
if(foundSlides.length===0){
console.log('Checking all child elements...');
const allChildren=elements.slideViewer.children;
console.log('Direct children:',allChildren.length);
if(allChildren.length>0){
foundSlides=Array.from(allChildren);
totalSlides=allChildren.length;
console.log('Using direct children as slides:',totalSlides);
}else{
showError('No slides found in presentation. The file may be corrupted or in an unsupported format.');
return;
}
}else{
totalSlides=foundSlides.length;
}
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
if(slide.parentNode===elements.slideViewer){
elements.slideViewer.insertBefore(wrapper,slide);
wrapper.appendChild(slide);
}else{
const clonedSlide=slide.cloneNode(true);
clonedSlide.style.maxWidth='100%%';
clonedSlide.style.width='100%%';
clonedSlide.style.height='auto';
clonedSlide.style.display='block';
wrapper.appendChild(clonedSlide);
elements.slideViewer.appendChild(wrapper);
}
const slideNumber=document.createElement('div');
slideNumber.className='ppt-slide-number';
slideNumber.textContent='Slide '+(index+1);
wrapper.appendChild(slideNumber);
slideWrappers.push(wrapper);
});
currentSlide=1;
updatePageInfo();
hideLoading();
console.log('PPTX rendered successfully. Total slides:',totalSlides);
}
},200);
}catch(err){
console.error('Error:',err);
showError('Failed to load: '+err.message);
}
}
if(elements.prevBtn){
elements.prevBtn.addEventListener('click',()=>{
if(viewMode==='single'&&currentSlide>1){
showSlide(currentSlide-1);
}
});
}
if(elements.nextBtn){
elements.nextBtn.addEventListener('click',()=>{
if(viewMode==='single'&&currentSlide<totalSlides){
showSlide(currentSlide+1);
}
});
}
if(elements.singleViewBtn){
elements.singleViewBtn.addEventListener('click',()=>{
switchViewMode('single');
});
}
if(elements.allViewBtn){
elements.allViewBtn.addEventListener('click',()=>{
switchViewMode('all');
});
}
if(elements.fullscreenBtn){
elements.fullscreenBtn.addEventListener('click',()=>{
const target=elements.slideViewer;
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
document.addEventListener('keydown',(e)=>{
if(!elements.wrapper||viewMode!=='single')return;
const activeElement=document.activeElement;
if(activeElement&&(activeElement.tagName==='INPUT'||activeElement.tagName==='TEXTAREA'||activeElement.isContentEditable))return;
const rect=elements.wrapper.getBoundingClientRect();
const inViewport=rect.top<window.innerHeight&&rect.bottom>0;
if(!inViewport)return;
switch(e.key){
case 'ArrowLeft':
case 'PageUp':
if(currentSlide>1)showSlide(currentSlide-1);
e.preventDefault();
break;
case 'ArrowRight':
case 'PageDown':
case ' ':
if(currentSlide<totalSlides)showSlide(currentSlide+1);
e.preventDefault();
break;
case 'Home':
showSlide(1);
e.preventDefault();
break;
case 'End':
showSlide(totalSlides);
e.preventDefault();
break;
case 'f':
if(elements.fullscreenBtn)elements.fullscreenBtn.click();
e.preventDefault();
break;
}
});
renderPPT();
})();
</script>`, uniqueId, safePptPath, fileName, strings.ToUpper(strings.TrimPrefix(ext, ".")))
}
