package parser

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"log"
)

// GeneratePICWarpHTML 生成用于在Web上展示图片内容的模块
// 使用 Viewer.js 进行优化展示，修复多实例加载问题
// 支持多实例隔离、点击放大、缩放、旋转、支持所有图片格式（包括GIF）
func GeneratePICWarpHTML(picId, picPath string) string {
	// 确保输入安全
	safePicId := template.HTMLEscapeString(picId)
	safePicPath := template.HTMLEscapeString(picPath)
	log.Println("Generating PIC viewer for FileName:", safePicPath)
	// 生成唯一ID（防止多实例冲突）
	h := md5.Sum([]byte(safePicId))
	uniqueId := fmt.Sprintf("pic_%x", h)[:12]

	return fmt.Sprintf(`<div id="%[1]s-wrapper" class="pic-viewer-wrapper" style="max-width:100%%;margin:20px auto;">
<style>
#%[1]s-wrapper{display:flex;justify-content:center;align-items:center;position:relative}
#%[1]s-wrapper .pic-container{position:relative;display:inline-block;max-width:100%%;cursor:zoom-in;transition:transform 0.2s ease}
#%[1]s-wrapper .pic-container:hover{transform:scale(1.02)}
#%[1]s-wrapper .pic-image{max-width:100%%;height:auto;display:block;border-radius:8px;box-shadow:0 4px 20px rgba(0,0,0,0.1);transition:box-shadow 0.3s ease}
#%[1]s-wrapper .pic-image:hover{box-shadow:0 8px 30px rgba(0,0,0,0.15)}
#%[1]s-wrapper .pic-loading{text-align:center;padding:60px 20px;color:#7f8c8d;font-size:14px;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,sans-serif}
#%[1]s-wrapper .pic-spinner{border:3px solid #ecf0f1;border-top:3px solid #3498db;border-radius:50%%;width:40px;height:40px;animation:pic-spin-%[1]s 1s linear infinite;margin:0 auto 15px}
@keyframes pic-spin-%[1]s{0%%{transform:rotate(0deg)}100%%{transform:rotate(360deg)}}
#%[1]s-wrapper .pic-error{text-align:center;padding:60px 20px;color:#e74c3c;font-size:14px;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,sans-serif}
#%[1]s-wrapper .pic-toolbar{position:absolute;top:10px;right:10px;display:flex;gap:8px;opacity:0;transition:opacity 0.3s ease;z-index:10}
#%[1]s-wrapper:hover .pic-toolbar{opacity:1}
#%[1]s-wrapper .pic-btn{background:rgba(0,0,0,0.6);backdrop-filter:blur(10px);border:none;color:#fff;padding:8px 12px;border-radius:6px;cursor:pointer;font-size:12px;transition:all 0.2s;text-decoration:none;display:inline-flex;align-items:center;gap:5px;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,sans-serif}
#%[1]s-wrapper .pic-btn:hover{background:rgba(0,0,0,0.8);transform:translateY(-1px)}
.viewer-backdrop{background-color:rgba(0,0,0,0.9)!important}
.viewer-toolbar>ul>li{background-color:rgba(255,255,255,0.1)!important;border-radius:6px}
.viewer-toolbar>ul>li:hover{background-color:rgba(255,255,255,0.2)!important}
@media (max-width:768px){
#%[1]s-wrapper{margin:10px}
#%[1]s-wrapper .pic-toolbar{opacity:1;position:static;justify-content:center;margin-top:10px}
}
</style>
<div class="pic-loading" id="%[1]s-loading">
<div class="pic-spinner"></div>
<div>Loading image...</div>
</div>
<div class="pic-container" id="%[1]s-container" style="display:none;">
<img id="%[1]s-image" class="pic-image" src="%[2]s" alt="Image Preview">
<div class="pic-toolbar">
<button class="pic-btn" id="%[1]s-fullscreen-btn" title="Fullscreen View">🔍 View</button>
<a href="%[2]s" download class="pic-btn" title="Download">⬇ Download</a>
<a href="%[2]s" target="_blank" rel="noopener noreferrer" class="pic-btn" title="Open in new tab">🔗 Open</a>
</div>
</div>
<div class="pic-error" id="%[1]s-error" style="display:none;">❌ Failed to load image</div>
</div>
<script>
(function(){
const UNIQUE_ID='%[1]s';
const IMAGE_URL='%[2]s';
const elements={wrapper:document.getElementById(UNIQUE_ID+'-wrapper'),loading:document.getElementById(UNIQUE_ID+'-loading'),container:document.getElementById(UNIQUE_ID+'-container'),image:document.getElementById(UNIQUE_ID+'-image'),error:document.getElementById(UNIQUE_ID+'-error'),fullscreenBtn:document.getElementById(UNIQUE_ID+'-fullscreen-btn')};
let viewer=null;
let loadingHidden=false;
let viewerInitialized=false;
function hideLoading(){
if(!loadingHidden&&elements.loading){
elements.loading.style.display='none';
loadingHidden=true;
}
}
function showContainer(){
if(elements.container){
elements.container.style.display='inline-block';
}
if(elements.error){
elements.error.style.display='none';
}
}
function showError(){
hideLoading();
if(elements.container){
elements.container.style.display='none';
}
if(elements.error){
elements.error.style.display='block';
}
}
function loadViewerJS(callback){
if(typeof Viewer!=='undefined'){
callback();
return;
}
if(document.getElementById('viewerjs-css')){
waitForViewer(callback);
return;
}
const cssLink=document.createElement('link');
cssLink.id='viewerjs-css';
cssLink.rel='stylesheet';
cssLink.href='https://cdn.jsdelivr.net/npm/viewerjs@1.11.6/dist/viewer.min.css';
document.head.appendChild(cssLink);
const script=document.createElement('script');
script.src='https://cdn.jsdelivr.net/npm/viewerjs@1.11.6/dist/viewer.min.js';
script.onload=function(){
console.log('Viewer.js loaded');
callback();
};
script.onerror=function(){
console.error('Failed to load Viewer.js');
};
document.head.appendChild(script);
}
function waitForViewer(callback){
let attempts=0;
const checkViewer=setInterval(function(){
attempts++;
if(typeof Viewer!=='undefined'){
clearInterval(checkViewer);
callback();
}else if(attempts>50){
clearInterval(checkViewer);
console.error('Viewer.js timeout');
}
},100);
}
function initViewer(){
if(viewerInitialized||!elements.image)return;
try{
viewer=new Viewer(elements.image,{inline:false,button:true,navbar:false,title:true,toolbar:{zoomIn:1,zoomOut:1,oneToOne:1,reset:1,prev:0,play:0,next:0,rotateLeft:1,rotateRight:1,flipHorizontal:1,flipVertical:1},tooltip:true,movable:true,zoomable:true,rotatable:true,scalable:true,transition:true,fullscreen:true,keyboard:true,url:'src'});
viewerInitialized=true;
console.log('Viewer.js initialized for:',IMAGE_URL);
}catch(err){
console.error('Viewer init error:',err);
}
}
elements.image.addEventListener('load',function(){
hideLoading();
showContainer();
loadViewerJS(function(){
initViewer();
});
console.log('Image loaded:',IMAGE_URL);
});
elements.image.addEventListener('error',function(){
showError();
console.error('Failed to load image:',IMAGE_URL);
});
if(elements.image.complete&&elements.image.naturalWidth>0){
hideLoading();
showContainer();
loadViewerJS(function(){
initViewer();
});
}
elements.image.addEventListener('click',function(){
if(viewer){
viewer.show();
}
});
if(elements.fullscreenBtn){
elements.fullscreenBtn.addEventListener('click',function(e){
e.stopPropagation();
if(viewer){
viewer.show();
}else{
loadViewerJS(function(){
initViewer();
if(viewer){
viewer.show();
}
});
}
});
}
setTimeout(function(){
if(!loadingHidden&&elements.image.complete){
hideLoading();
showContainer();
}
},3000);
window.addEventListener('beforeunload',function(){
if(viewer){
try{
viewer.destroy();
}catch(e){
console.error('Viewer destroy error:',e);
}
}
});
})();
</script>`, uniqueId, safePicPath)
}
