package parser

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"log"
	"path/filepath"
	"strings"
)

// GeneratePICWarpHTML 生成用于在Web上展示图片内容的模块
// 企业级商业可用版本：支持多实例隔离、智能加载、重试机制、完善的错误处理
// 使用 Viewer.js 进行优化展示，支持所有图片格式（包括GIF、WebP、AVIF等）
func GeneratePICWarpHTML(picId, picPath string) string {
	safePicId := template.HTMLEscapeString(picId)
	safePicPath := template.HTMLEscapeString(picPath)
	log.Println("Generating PIC viewer for FileName:", safePicPath)

	h := md5.Sum([]byte(safePicId))
	uniqueId := fmt.Sprintf("pic_%x", h)[:12]

	fileName := filepath.Base(picPath)
	ext := strings.ToLower(filepath.Ext(fileName))
	safeFileName := template.JSEscapeString(fileName)

	return fmt.Sprintf(`<div id="%[1]s-wrapper" class="pic-viewer-wrapper" style="max-width:100%%;margin:20px auto;">
<style>
#%[1]s-wrapper{display:flex;justify-content:center;align-items:center;position:relative;min-height:200px}
#%[1]s-wrapper .pic-container{position:relative;display:inline-block;max-width:100%%;cursor:zoom-in;transition:transform 0.2s ease}
#%[1]s-wrapper .pic-container:hover{transform:scale(1.02)}
#%[1]s-wrapper .pic-image{max-width:100%%;height:auto;display:block;border-radius:12px;box-shadow:0 4px 24px rgba(0,0,0,0.12);transition:all 0.3s ease;background:#f8f9fa}
#%[1]s-wrapper .pic-image:hover{box-shadow:0 8px 40px rgba(0,0,0,0.18);transform:translateY(-2px)}
#%[1]s-wrapper .pic-loading{text-align:center;padding:60px 20px;color:#6c757d;font-size:14px;font-family:-apple-system,BlinkMacSystemFont,'SF Pro Text','Segoe UI',Roboto,sans-serif;display:flex;flex-direction:column;align-items:center;gap:16px}
#%[1]s-wrapper .pic-loading-text{font-weight:500}
#%[1]s-wrapper .pic-loading-hint{font-size:12px;color:#95a5a6;margin-top:8px}
#%[1]s-wrapper .pic-spinner{border:3px solid rgba(52,152,219,0.2);border-top:3px solid #3498db;border-radius:50%%;width:44px;height:44px;animation:pic-spin-%[1]s 0.8s cubic-bezier(0.4,0,0.2,1) infinite;box-shadow:0 0 15px rgba(52,152,219,0.3)}
@keyframes pic-spin-%[1]s{0%%{transform:rotate(0deg)}100%%{transform:rotate(360deg)}}
#%[1]s-wrapper .pic-error{text-align:center;padding:50px 20px;font-size:14px;font-family:-apple-system,BlinkMacSystemFont,'SF Pro Text','Segoe UI',Roboto,sans-serif;background:rgba(231,76,60,0.05);border-radius:12px;border:2px dashed rgba(231,76,60,0.3)}
#%[1]s-wrapper .pic-error-icon{font-size:48px;margin-bottom:16px}
#%[1]s-wrapper .pic-error-title{color:#e74c3c;font-weight:600;font-size:16px;margin-bottom:8px}
#%[1]s-wrapper .pic-error-message{color:#7f8c8d;font-size:13px;margin-bottom:20px;line-height:1.6}
#%[1]s-wrapper .pic-error-actions{display:flex;gap:12px;justify-content:center;flex-wrap:wrap}
#%[1]s-wrapper .pic-toolbar{position:absolute;top:12px;right:12px;display:flex;gap:8px;opacity:0;transition:opacity 0.3s ease;z-index:10}
#%[1]s-wrapper:hover .pic-toolbar{opacity:1}
#%[1]s-wrapper .pic-btn{background:rgba(0,0,0,0.65);backdrop-filter:blur(12px);-webkit-backdrop-filter:blur(12px);border:1px solid rgba(255,255,255,0.1);color:#fff;padding:8px 14px;border-radius:8px;cursor:pointer;font-size:12px;font-weight:500;transition:all 0.2s cubic-bezier(0.4,0,0.2,1);text-decoration:none;display:inline-flex;align-items:center;gap:6px;font-family:-apple-system,BlinkMacSystemFont,'SF Pro Text','Segoe UI',Roboto,sans-serif;box-shadow:0 2px 8px rgba(0,0,0,0.2);white-space:nowrap}
#%[1]s-wrapper .pic-btn:hover{background:rgba(0,0,0,0.85);transform:translateY(-2px);box-shadow:0 4px 12px rgba(0,0,0,0.3)}
#%[1]s-wrapper .pic-btn:active{transform:translateY(0)}
#%[1]s-wrapper .pic-btn-retry{background:linear-gradient(135deg,#27ae60 0%%,#229954 100%%);border:none}
#%[1]s-wrapper .pic-btn-retry:hover{background:linear-gradient(135deg,#2ecc71 0%%,#27ae60 100%%)}
#%[1]s-wrapper .pic-btn-primary{background:linear-gradient(135deg,#3498db 0%%,#2980b9 100%%);border:none}
#%[1]s-wrapper .pic-btn-primary:hover{background:linear-gradient(135deg,#5dade2 0%%,#3498db 100%%)}
#%[1]s-wrapper .pic-info{position:absolute;bottom:12px;left:12px;background:rgba(0,0,0,0.65);backdrop-filter:blur(12px);color:#fff;padding:8px 14px;border-radius:8px;font-size:11px;opacity:0;transition:opacity 0.3s ease;font-family:'SF Mono','Monaco','Consolas',monospace;border:1px solid rgba(255,255,255,0.1)}
#%[1]s-wrapper:hover .pic-info{opacity:1}
#%[1]s-wrapper .pic-retry-badge{position:absolute;top:12px;left:12px;background:rgba(241,196,15,0.9);color:#fff;padding:6px 12px;border-radius:6px;font-size:11px;font-weight:600;box-shadow:0 2px 8px rgba(241,196,15,0.3)}
.viewer-backdrop{background-color:rgba(0,0,0,0.92)!important;backdrop-filter:blur(8px)!important}
.viewer-toolbar>ul>li{background-color:rgba(255,255,255,0.12)!important;border-radius:8px;margin:0 4px!important;transition:all 0.2s!important}
.viewer-toolbar>ul>li:hover{background-color:rgba(255,255,255,0.25)!important;transform:translateY(-2px)!important}
.viewer-title{background:rgba(0,0,0,0.6)!important;backdrop-filter:blur(10px)!important;padding:12px 20px!important;border-radius:8px!important}
@media (max-width:768px){
#%[1]s-wrapper{margin:10px}
#%[1]s-wrapper .pic-toolbar{opacity:1;position:static;justify-content:center;margin-top:12px;gap:6px}
#%[1]s-wrapper .pic-btn{padding:6px 10px;font-size:11px}
#%[1]s-wrapper .pic-info{position:static;opacity:1;margin-top:10px;text-align:center}
#%[1]s-wrapper .pic-container:hover{transform:none}
}
</style>
<div class="pic-loading" id="%[1]s-loading">
<div class="pic-spinner"></div>
<div class="pic-loading-text">Loading image...</div>
<div class="pic-loading-hint">Please wait a moment</div>
</div>
<div class="pic-container" id="%[1]s-container" style="display:none;">
<div class="pic-retry-badge" id="%[1]s-retry-badge" style="display:none;">Retry <span id="%[1]s-retry-count">0</span></div>
<img id="%[1]s-image" class="pic-image" data-src="%[2]s" alt="Image Preview">
<div class="pic-toolbar">
<button class="pic-btn pic-btn-primary" id="%[1]s-fullscreen-btn" title="Fullscreen View">🔍 View</button>
<a href="%[2]s" download class="pic-btn" title="Download Image">⬇ Download</a>
<a href="%[2]s" target="_blank" rel="noopener noreferrer" class="pic-btn" title="Open in New Tab">🔗 Open</a>
</div>
<div class="pic-info" id="%[1]s-info">%[4]s</div>
</div>
<div class="pic-error" id="%[1]s-error" style="display:none;">
<div class="pic-error-icon">🖼️</div>
<div class="pic-error-title">Failed to Load Image</div>
<div class="pic-error-message" id="%[1]s-error-msg">Unable to load the image. Please check your connection.</div>
<div class="pic-error-actions">
<button class="pic-btn pic-btn-retry" id="%[1]s-retry-btn">🔄 Retry</button>
<a href="%[2]s" download class="pic-btn">⬇ Download Instead</a>
<a href="%[2]s" target="_blank" class="pic-btn">🔗 Open Direct</a>
</div>
</div>
</div>
<script>
(function(){
'use strict';
const UNIQUE_ID='%[1]s';
const IMAGE_URL='%[2]s';
const FILE_NAME='%[3]s';
if(!window.__PICViewerManager){
window.__PICViewerManager={
instances:{},
viewerJSLoaded:false,
viewerJSLoading:false,
loadCallbacks:[],
loadViewerJS:function(callback){
if(this.viewerJSLoaded){
callback(true);
return;
}
this.loadCallbacks.push(callback);
if(this.viewerJSLoading)return;
this.viewerJSLoading=true;
const existingCSS=document.getElementById('viewerjs-global-css');
const existingJS=document.querySelector('script[src*="viewerjs"]');
if(existingCSS&&existingJS&&typeof Viewer!=='undefined'){
this.viewerJSLoaded=true;
this.viewerJSLoading=false;
this.loadCallbacks.forEach(cb=>cb(true));
this.loadCallbacks=[];
return;
}
if(!existingCSS){
const cssLink=document.createElement('link');
cssLink.id='viewerjs-global-css';
cssLink.rel='stylesheet';
cssLink.href='https://cdn.jsdelivr.net/npm/viewerjs@1.11.6/dist/viewer.min.css';
document.head.appendChild(cssLink);
}
if(!existingJS){
const script=document.createElement('script');
script.src='https://cdn.jsdelivr.net/npm/viewerjs@1.11.6/dist/viewer.min.js';
script.onload=()=>{
this.viewerJSLoaded=true;
this.viewerJSLoading=false;
this.loadCallbacks.forEach(cb=>cb(true));
this.loadCallbacks=[];
};
script.onerror=()=>{
this.viewerJSLoading=false;
this.loadCallbacks.forEach(cb=>cb(false));
this.loadCallbacks=[];
};
document.head.appendChild(script);
}else{
this.waitForViewer();
}
},
waitForViewer:function(){
let attempts=0;
const check=setInterval(()=>{
attempts++;
if(typeof Viewer!=='undefined'){
clearInterval(check);
this.viewerJSLoaded=true;
this.viewerJSLoading=false;
this.loadCallbacks.forEach(cb=>cb(true));
this.loadCallbacks=[];
}else if(attempts>60){
clearInterval(check);
this.viewerJSLoading=false;
this.loadCallbacks.forEach(cb=>cb(false));
this.loadCallbacks=[];
}
},100);
},
registerInstance:function(id,instance){
this.instances[id]=instance;
},
unregisterInstance:function(id){
if(this.instances[id]){
if(this.instances[id].viewer){
try{this.instances[id].viewer.destroy();}catch(e){}
}
delete this.instances[id];
}
}
};
}
class PICViewer{
constructor(uniqueId,imageUrl,fileName){
this.uniqueId=uniqueId;
this.imageUrl=imageUrl;
this.fileName=fileName;
this.viewer=null;
this.retryCount=0;
this.maxRetries=3;
this.isDestroyed=false;
this.loadStartTime=Date.now();
this.elements={
wrapper:document.getElementById(uniqueId+'-wrapper'),
loading:document.getElementById(uniqueId+'-loading'),
container:document.getElementById(uniqueId+'-container'),
image:document.getElementById(uniqueId+'-image'),
error:document.getElementById(uniqueId+'-error'),
errorMsg:document.getElementById(uniqueId+'-error-msg'),
retryBtn:document.getElementById(uniqueId+'-retry-btn'),
retryBadge:document.getElementById(uniqueId+'-retry-badge'),
retryCountEl:document.getElementById(uniqueId+'-retry-count'),
fullscreenBtn:document.getElementById(uniqueId+'-fullscreen-btn'),
info:document.getElementById(uniqueId+'-info')
};
this.eventListeners=[];
this.init();
}
init(){
window.__PICViewerManager.registerInstance(this.uniqueId,this);
this.bindEvents();
this.loadImage();
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
const loadingText=this.elements.loading?.querySelector('.pic-loading-text');
const loadingHint=this.elements.loading?.querySelector('.pic-loading-hint');
if(loadingText)loadingText.textContent=text||'Loading image...';
if(hint!==undefined&&loadingHint)loadingHint.textContent=hint||'';
}
hideLoading(){
if(this.elements.loading){
this.elements.loading.style.display='none';
}
}
showContainer(){
if(this.elements.container){
this.elements.container.style.display='inline-block';
}
this.hideError();
}
showError(message,canRetry=true){
this.hideLoading();
if(this.elements.container){
this.elements.container.style.display='none';
}
if(this.elements.error){
this.elements.error.style.display='block';
if(this.elements.errorMsg){
this.elements.errorMsg.textContent=message||'Unable to load the image. Please check your connection.';
}
if(this.elements.retryBtn){
this.elements.retryBtn.style.display=canRetry?'inline-flex':'none';
}
}
}
hideError(){
if(this.elements.error){
this.elements.error.style.display='none';
}
}
updateRetryInfo(){
if(this.retryCount>0){
if(this.elements.retryBadge){
this.elements.retryBadge.style.display='block';
}
if(this.elements.retryCountEl){
this.elements.retryCountEl.textContent=this.retryCount+'/'+this.maxRetries;
}
}
}
loadImage(){
if(this.isDestroyed)return;
if(this.retryCount>0){
this.updateLoadingText('Retrying... ('+this.retryCount+'/'+this.maxRetries+')','Attempting to reload the image');
}
this.loadStartTime=Date.now();
const img=this.elements.image;
if(!img)return;
const imgUrl=img.getAttribute('data-src');
if(!imgUrl)return;
const tempImg=new Image();
this.addEventListener(tempImg,'load',()=>{
if(this.isDestroyed)return;
const loadTime=((Date.now()-this.loadStartTime)/1000).toFixed(2);
img.src=imgUrl;
img.removeAttribute('data-src');
this.hideLoading();
this.showContainer();
this.hideError();
this.updateRetryInfo();
if(this.elements.info){
const sizeInfo=tempImg.naturalWidth+'×'+tempImg.naturalHeight;
this.elements.info.textContent=this.fileName+' • '+sizeInfo+' • '+loadTime+'s';
}
this.initViewerJS();
});
this.addEventListener(tempImg,'error',()=>{
if(this.isDestroyed)return;
if(this.retryCount<this.maxRetries){
this.showError('Failed to load image. Click retry to try again ('+this.retryCount+'/'+this.maxRetries+').',true);
}else{
this.showError('Failed to load image after '+this.maxRetries+' attempts. The image may be unavailable or corrupted.',false);
}
});
tempImg.src=imgUrl;
}
retry(){
if(this.retryCount>=this.maxRetries){
this.showError('Maximum retry attempts reached. Please download the image or open it directly.',false);
return;
}
this.retryCount++;
this.hideError();
if(this.elements.loading){
this.elements.loading.style.display='flex';
}
setTimeout(()=>{
if(this.elements.image){
this.elements.image.removeAttribute('src');
this.elements.image.setAttribute('data-src',this.imageUrl);
}
this.loadImage();
},1000);
}
initViewerJS(){
if(this.isDestroyed)return;
window.__PICViewerManager.loadViewerJS((success)=>{
if(this.isDestroyed)return;
if(!success){
console.warn('Viewer.js failed to load, basic functionality only');
return;
}
try{
if(typeof Viewer==='undefined'){
console.warn('Viewer not available');
return;
}
this.viewer=new Viewer(this.elements.image,{
inline:false,
button:true,
navbar:false,
title:true,
toolbar:{
zoomIn:1,
zoomOut:1,
oneToOne:1,
reset:1,
prev:0,
play:0,
next:0,
rotateLeft:1,
rotateRight:1,
flipHorizontal:1,
flipVertical:1
},
tooltip:true,
movable:true,
zoomable:true,
rotatable:true,
scalable:true,
transition:true,
fullscreen:true,
keyboard:true,
url:'src',
ready:()=>{
console.log('Viewer ready for:',this.fileName);
},
show:()=>{
console.log('Viewer shown');
},
hide:()=>{
console.log('Viewer hidden');
}
});
}catch(err){
console.error('Viewer init error:',err);
}
});
}
showViewer(){
if(this.viewer){
this.viewer.show();
}else{
this.initViewerJS();
setTimeout(()=>{
if(this.viewer)this.viewer.show();
},300);
}
}
bindEvents(){
if(this.elements.image){
this.addEventListener(this.elements.image,'click',()=>this.showViewer());
}
if(this.elements.fullscreenBtn){
this.addEventListener(this.elements.fullscreenBtn,'click',(e)=>{
e.stopPropagation();
this.showViewer();
});
}
if(this.elements.retryBtn){
this.addEventListener(this.elements.retryBtn,'click',()=>this.retry());
}
}
destroy(){
this.isDestroyed=true;
this.removeAllEventListeners();
if(this.viewer){
try{this.viewer.destroy();}catch(e){}
this.viewer=null;
}
window.__PICViewerManager.unregisterInstance(this.uniqueId);
}
}
try{
const viewer=new PICViewer(UNIQUE_ID,IMAGE_URL,FILE_NAME);
window.addEventListener('beforeunload',()=>viewer.destroy());
}catch(error){
console.error('Failed to initialize PIC viewer:',error);
}
})();
</script>`, uniqueId, safePicPath, safeFileName, strings.ToUpper(strings.TrimPrefix(ext, ".")))
}
