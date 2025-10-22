package parser

import (
	"crypto/md5"
	"fmt"
	"html/template"

	"path/filepath"
	"strings"
)

// GenerateVideoWarpHTML 生成用于在Web上展示视频内容的模块
// 使用原生HTML5视频播放器+Video.js增强，支持多种视频格式
// 支持WebCodecs API处理特殊编码，Range请求按需加载
// 自动提取视频封面,居中显示,简洁设计,完美全屏支持
func GenerateVideoWarpHTML(videoId, videoPath string) string {
	safeVideoId := template.HTMLEscapeString(videoId)
	safeVideoPath := template.HTMLEscapeString(videoPath)
	h := md5.Sum([]byte(safeVideoId))
	uniqueId := fmt.Sprintf("video_%x", h)[:12]
	ext := strings.ToLower(filepath.Ext(videoPath))
	videoType := getVideoMimeType(ext)
	return fmt.Sprintf(`<div id="%[1]s-wrapper" class="video-viewer-wrapper" style="position:relative;max-width:100%%;margin:20px auto;display:block;">
<style>
#%[1]s-wrapper{position:relative;width:100%%;max-width:900px;margin:20px auto;display:block;clear:both}
#%[1]s-wrapper .video-container{position:relative;width:100%%;background:#000;border-radius:8px;overflow:visible;box-shadow:0 4px 20px rgba(0,0,0,0.1);transition:box-shadow 0.3s ease}
#%[1]s-wrapper .video-container:hover{box-shadow:0 8px 30px rgba(0,0,0,0.15)}
#%[1]s-wrapper .video-inner{position:relative;width:100%%;height:0;padding-bottom:56.25%%;overflow:hidden;border-radius:8px}
#%[1]s-wrapper .video-js{position:absolute!important;top:0!important;left:0!important;width:100%%!important;height:100%%!important;border-radius:8px}
.video-js.vjs-fullscreen,.video-js.vjs-fullscreen.vjs-user-inactive{position:fixed!important;top:0!important;left:0!important;width:100vw!important;height:100vh!important;max-width:none!important;max-height:none!important;z-index:999999!important;border-radius:0!important;margin:0!important;padding:0!important}
.video-js.vjs-fullscreen video,.video-js.vjs-fullscreen .vjs-tech{position:absolute!important;top:0!important;left:0!important;width:100%%!important;height:100%%!important;max-width:none!important;max-height:none!important;object-fit:contain!important}
#%[1]s-wrapper video{position:absolute!important;top:0!important;left:0!important;width:100%%!important;height:100%%!important;object-fit:contain;background:#000}
#%[1]s-wrapper .vjs-control-bar{position:absolute!important;bottom:0!important;width:100%%!important;background:linear-gradient(0deg,rgba(0,0,0,0.8) 0%%,rgba(0,0,0,0.4) 100%%);backdrop-filter:blur(10px);border-radius:0 0 8px 8px;z-index:20}
.video-js.vjs-fullscreen .vjs-control-bar{border-radius:0!important}
#%[1]s-wrapper .video-poster{position:absolute;top:0;left:0;width:100%%;height:100%%;object-fit:cover;z-index:1;cursor:pointer;transition:opacity 0.3s ease;border-radius:8px}
#%[1]s-wrapper .video-poster.hidden{opacity:0;pointer-events:none;z-index:-1}
#%[1]s-wrapper .play-button-overlay{position:absolute;top:50%%;left:50%%;transform:translate(-50%%,-50%%);width:80px;height:80px;background:rgba(0,0,0,0.7);border:3px solid #fff;border-radius:50%%;cursor:pointer;z-index:2;display:none;align-items:center;justify-content:center;transition:all 0.3s ease}
#%[1]s-wrapper .play-button-overlay.show{display:flex}
#%[1]s-wrapper .play-button-overlay:hover{background:rgba(52,152,219,0.9);transform:translate(-50%%,-50%%) scale(1.1)}
#%[1]s-wrapper .play-button-overlay::after{content:'';width:0;height:0;border-left:25px solid #fff;border-top:15px solid transparent;border-bottom:15px solid transparent;margin-left:5px}
#%[1]s-wrapper .vjs-big-play-button{display:none!important}
#%[1]s-wrapper .vjs-play-progress,#%[1]s-wrapper .vjs-volume-level{background-color:#3498db}
#%[1]s-wrapper .video-loading{position:absolute;top:50%%;left:50%%;transform:translate(-50%%,-50%%);text-align:center;color:#fff;font-size:14px;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,sans-serif;z-index:15;background:rgba(0,0,0,0.8);padding:20px;border-radius:8px;min-width:150px}
#%[1]s-wrapper .video-spinner{border:3px solid rgba(255,255,255,0.2);border-top:3px solid #3498db;border-radius:50%%;width:40px;height:40px;animation:video-spin-%[1]s 1s linear infinite;margin:0 auto 15px}
@keyframes video-spin-%[1]s{0%%{transform:rotate(0deg)}100%%{transform:rotate(360deg)}}
#%[1]s-wrapper .video-loading-progress{margin-top:10px;height:4px;background:rgba(255,255,255,0.2);border-radius:2px;overflow:hidden}
#%[1]s-wrapper .video-loading-progress-bar{height:100%%;background:#3498db;width:0%%;transition:width 0.3s ease}
#%[1]s-wrapper .video-error{position:absolute;top:50%%;left:50%%;transform:translate(-50%%,-50%%);text-align:center;padding:20px;color:#e74c3c;font-size:14px;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,sans-serif;background:rgba(255,255,255,0.95);border-radius:8px;box-shadow:0 2px 10px rgba(0,0,0,0.1);z-index:10;max-width:80%%}
#%[1]s-wrapper .video-error-title{font-weight:600;margin-bottom:10px}
#%[1]s-wrapper .video-error-message{font-size:12px;color:#666;margin-bottom:15px}
#%[1]s-wrapper .video-retry-btn{background:#3498db;color:#fff;border:none;padding:8px 16px;border-radius:4px;cursor:pointer;font-size:13px;transition:background 0.2s}
#%[1]s-wrapper .video-retry-btn:hover{background:#2980b9}
#%[1]s-wrapper .video-toolbar{position:absolute;top:10px;right:10px;display:flex;gap:8px;opacity:0;transition:opacity 0.3s ease;z-index:10}
#%[1]s-wrapper:hover .video-toolbar{opacity:1}
#%[1]s-wrapper .video-btn{background:rgba(0,0,0,0.6);backdrop-filter:blur(10px);border:none;color:#fff;padding:8px 12px;border-radius:6px;cursor:pointer;font-size:12px;transition:all 0.2s;text-decoration:none;display:inline-flex;align-items:center;gap:5px;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,sans-serif}
#%[1]s-wrapper .video-btn:hover{background:rgba(0,0,0,0.8);transform:translateY(-1px)}
#%[1]s-wrapper .video-codec-info{position:absolute;bottom:60px;left:10px;background:rgba(0,0,0,0.7);color:#fff;padding:5px 10px;border-radius:4px;font-size:11px;opacity:0;transition:opacity 0.3s;z-index:5}
#%[1]s-wrapper:hover .video-codec-info{opacity:1}
@media (max-width:768px){
#%[1]s-wrapper{margin:10px}
#%[1]s-wrapper .video-toolbar{opacity:1}
#%[1]s-wrapper .play-button-overlay{width:60px;height:60px}
#%[1]s-wrapper .play-button-overlay::after{border-left-width:20px;border-top-width:12px;border-bottom-width:12px}
#%[1]s-wrapper .video-codec-info{display:none}
}
</style>
<div class="video-container">
<div class="video-inner">
<div class="video-loading" id="%[1]s-loading">
<div class="video-spinner"></div>
<div id="%[1]s-loading-text">Initializing...</div>
<div class="video-loading-progress" id="%[1]s-loading-progress" style="display:none;">
<div class="video-loading-progress-bar" id="%[1]s-loading-progress-bar"></div>
</div>
</div>
<canvas id="%[1]s-poster" class="video-poster" style="display:none;"></canvas>
<div id="%[1]s-play-btn" class="play-button-overlay"></div>
<div class="video-toolbar">
<button class="video-btn" id="%[1]s-pip-btn" title="Picture in Picture">📺 PiP</button>
<a href="%[2]s" target="_blank" rel="noopener noreferrer" class="video-btn" title="Open in new tab">🔗 Open</a>
</div>
<div class="video-codec-info" id="%[1]s-codec-info">Loading...</div>
<video id="%[1]s-player" class="video-js vjs-default-skin vjs-16-9" controls preload="auto" crossorigin="anonymous">
<source src="%[2]s" type="%[3]s">
<p class="vjs-no-js">To view this video please enable JavaScript</p>
</video>
<div class="video-error" id="%[1]s-error" style="display:none;">
<div class="video-error-title">❌ Failed to load video</div>
<div class="video-error-message" id="%[1]s-error-msg">Unable to load video content</div>
<button class="video-retry-btn" id="%[1]s-retry">Retry</button>
</div>
</div>
</div>
</div>
<script>
(function(){
'use strict';
const UNIQUE_ID='%[1]s';
const VIDEO_URL='%[2]s';
const VIDEO_TYPE='%[3]s';
if(!window.__VideoPlayerManager){
window.__VideoPlayerManager={
instances:{},
videojsLoaded:false,
videojsLoading:false,
loadCallbacks:[],
keyboardHandlerAttached:false,
loadVideoJS:function(callback){
if(this.videojsLoaded){
callback(true);
return;
}
this.loadCallbacks.push(callback);
if(this.videojsLoading)return;
this.videojsLoading=true;
if(!document.querySelector('link[href*="video-js.min.css"]')){
const link=document.createElement('link');
link.rel='stylesheet';
link.href='https://cdn.jsdelivr.net/npm/video.js@8.10.0/dist/video-js.min.css';
document.head.appendChild(link);
}
if(!document.querySelector('script[src*="video.min.js"]')){
const script=document.createElement('script');
script.src='https://cdn.jsdelivr.net/npm/video.js@8.10.0/dist/video.min.js';
script.onload=()=>{
this.videojsLoaded=true;
this.videojsLoading=false;
this.loadCallbacks.forEach(cb=>cb(true));
this.loadCallbacks=[];
};
script.onerror=()=>{
this.videojsLoading=false;
this.loadCallbacks.forEach(cb=>cb(false));
this.loadCallbacks=[];
};
document.head.appendChild(script);
}else if(typeof videojs!=='undefined'){
this.videojsLoaded=true;
this.videojsLoading=false;
this.loadCallbacks.forEach(cb=>cb(true));
this.loadCallbacks=[];
}
},
attachKeyboardHandler:function(){
if(this.keyboardHandlerAttached)return;
this.keyboardHandlerAttached=true;
document.addEventListener('keydown',function(e){
const activeElement=document.activeElement;
if(activeElement&&(activeElement.tagName==='INPUT'||activeElement.tagName==='TEXTAREA'||activeElement.isContentEditable))return;
let targetPlayer=null;
let targetWrapper=null;
for(const id in window.__VideoPlayerManager.instances){
const instance=window.__VideoPlayerManager.instances[id];
if(!instance.player||!instance.elements.wrapper)continue;
if(instance.player.isFullscreen()){
targetPlayer=instance.player;
targetWrapper=instance.elements.wrapper;
break;
}
}
if(!targetPlayer){
for(const id in window.__VideoPlayerManager.instances){
const instance=window.__VideoPlayerManager.instances[id];
if(!instance.player||!instance.elements.wrapper)continue;
const rect=instance.elements.wrapper.getBoundingClientRect();
const inViewport=rect.top<window.innerHeight&&rect.bottom>0;
if(inViewport){
targetPlayer=instance.player;
targetWrapper=instance.elements.wrapper;
break;
}
}
}
if(!targetPlayer)return;
switch(e.key){
case ' ':
case 'k':
if(targetPlayer.paused())targetPlayer.play();
else targetPlayer.pause();
e.preventDefault();
break;
case 'f':
if(targetPlayer.isFullscreen())targetPlayer.exitFullscreen();
else targetPlayer.requestFullscreen();
e.preventDefault();
break;
case 'Escape':
if(targetPlayer.isFullscreen())targetPlayer.exitFullscreen();
break;
case 'm':
targetPlayer.muted(!targetPlayer.muted());
e.preventDefault();
break;
case 'ArrowLeft':
targetPlayer.currentTime(Math.max(0,targetPlayer.currentTime()-5));
e.preventDefault();
break;
case 'ArrowRight':
targetPlayer.currentTime(Math.min(targetPlayer.duration()||0,targetPlayer.currentTime()+5));
e.preventDefault();
break;
}
});
},
registerInstance:function(id,instance){
this.instances[id]=instance;
this.attachKeyboardHandler();
},
unregisterInstance:function(id){
if(this.instances[id]){
if(this.instances[id].player){
try{this.instances[id].player.dispose();}catch(e){}
}
delete this.instances[id];
}
}
};
}
class VideoPlayer{
constructor(uniqueId,videoUrl,videoType){
this.uniqueId=uniqueId;
this.videoUrl=videoUrl;
this.videoType=videoType;
this.player=null;
this.posterGenerated=false;
this.retryCount=0;
this.maxRetries=3;
this.isDestroyed=false;
this.eventListeners=[];
this.loadingState='init';
this.supportsWebCodecs=false;
this.codecInfo=null;
this.bufferProgress=0;
this.isPlaying=false;
this.elements={
wrapper:document.getElementById(uniqueId+'-wrapper'),
loading:document.getElementById(uniqueId+'-loading'),
loadingText:document.getElementById(uniqueId+'-loading-text'),
loadingProgress:document.getElementById(uniqueId+'-loading-progress'),
loadingProgressBar:document.getElementById(uniqueId+'-loading-progress-bar'),
error:document.getElementById(uniqueId+'-error'),
errorMsg:document.getElementById(uniqueId+'-error-msg'),
retryBtn:document.getElementById(uniqueId+'-retry'),
videoEl:document.getElementById(uniqueId+'-player'),
posterCanvas:document.getElementById(uniqueId+'-poster'),
playBtn:document.getElementById(uniqueId+'-play-btn'),
pipBtn:document.getElementById(uniqueId+'-pip-btn'),
codecInfo:document.getElementById(uniqueId+'-codec-info')
};
this.checkWebCodecsSupport();
this.init();
}
checkWebCodecsSupport(){
this.supportsWebCodecs='VideoDecoder' in window && 'VideoEncoder' in window;
console.log('WebCodecs API support:',this.supportsWebCodecs);
}
init(){
window.__VideoPlayerManager.registerInstance(this.uniqueId,this);
if(this.elements.retryBtn){
this.addEventListener(this.elements.retryBtn,'click',()=>this.retry());
}
this.showLoading('Checking video format...');
this.loadPlayer();
}
addEventListener(element,event,handler){
if(!element)return;
element.addEventListener(event,handler);
this.eventListeners.push({element,event,handler});
}
removeAllEventListeners(){
this.eventListeners.forEach(({element,event,handler})=>{
element.removeEventListener(event,handler);
});
this.eventListeners=[];
}
showLoading(text){
this.loadingState='loading';
if(this.elements.loading){
this.elements.loading.style.display='block';
if(this.elements.loadingText){
this.elements.loadingText.textContent=text||'Loading video...';
}
}
}
updateLoadingProgress(percent){
if(this.elements.loadingProgress&&percent>=0){
this.elements.loadingProgress.style.display='block';
if(this.elements.loadingProgressBar){
this.elements.loadingProgressBar.style.width=percent+'%%';
}
if(this.elements.loadingText){
this.elements.loadingText.textContent='Buffering... '+Math.round(percent)+'%%';
}
}
}
hideLoading(){
this.loadingState='loaded';
if(this.elements.loading){
this.elements.loading.style.display='none';
}
if(this.elements.loadingProgress){
this.elements.loadingProgress.style.display='none';
}
}
showPlayButton(){
if(this.elements.playBtn){
this.elements.playBtn.classList.add('show');
}
}
hidePlayButton(){
if(this.elements.playBtn){
this.elements.playBtn.classList.remove('show');
}
}
showError(message,canRetry=true){
this.hideLoading();
if(this.elements.error){
this.elements.error.style.display='block';
if(this.elements.errorMsg){
this.elements.errorMsg.textContent=message||'Unable to load video content';
}
if(this.elements.retryBtn){
this.elements.retryBtn.style.display=canRetry?'inline-block':'none';
}
}
}
hideError(){
if(this.elements.error)this.elements.error.style.display='none';
}
updateCodecInfo(info){
if(this.elements.codecInfo){
this.elements.codecInfo.textContent=info;
this.codecInfo=info;
}
}
retry(){
if(this.retryCount>=this.maxRetries){
this.showError('Maximum retry attempts reached. Please check your connection.',false);
return;
}
this.retryCount++;
this.hideError();
this.showLoading('Retrying... (Attempt '+this.retryCount+'/'+this.maxRetries+')');
if(this.player){
try{this.player.dispose();}catch(e){}
this.player=null;
}
this.posterGenerated=false;
this.isPlaying=false;
setTimeout(()=>this.loadPlayer(),1000);
}
async probeVideoCodec(){
if(!this.elements.videoEl)return null;
return new Promise((resolve)=>{
const videoEl=this.elements.videoEl;
const handleLoadedMetadata=()=>{
videoEl.removeEventListener('loadedmetadata',handleLoadedMetadata);
const track=videoEl.videoTracks&&videoEl.videoTracks[0];
if(track){
this.updateCodecInfo('Codec: '+(track.label||'Unknown'));
}
const info={
duration:videoEl.duration,
videoWidth:videoEl.videoWidth,
videoHeight:videoEl.videoHeight,
supportsRange:false
};
fetch(this.videoUrl,{method:'HEAD'})
.then(response=>{
info.supportsRange=response.headers.get('Accept-Ranges')==='bytes';
console.log('Server supports Range requests:',info.supportsRange);
resolve(info);
})
.catch(()=>resolve(info));
};
videoEl.addEventListener('loadedmetadata',handleLoadedMetadata);
setTimeout(()=>{
videoEl.removeEventListener('loadedmetadata',handleLoadedMetadata);
resolve(null);
},5000);
});
}
loadPlayer(){
this.showLoading('Initializing player...');
window.__VideoPlayerManager.loadVideoJS((success)=>{
if(this.isDestroyed)return;
if(success){
this.initVideoJS();
}else{
this.fallbackToNative();
}
});
}
generatePoster(){
if(this.posterGenerated||!this.elements.videoEl)return;
const videoEl=this.elements.videoEl;
const canvas=this.elements.posterCanvas;
const handleLoadedData=()=>{
if(this.posterGenerated)return;
setTimeout(()=>{
if(videoEl.duration&&videoEl.duration>0){
videoEl.currentTime=Math.min(videoEl.duration*0.1,5);
}
},100);
};
const handleSeeked=()=>{
if(this.posterGenerated)return;
try{
canvas.width=videoEl.videoWidth||640;
canvas.height=videoEl.videoHeight||360;
const ctx=canvas.getContext('2d');
ctx.drawImage(videoEl,0,0,canvas.width,canvas.height);
canvas.style.display='block';
this.posterGenerated=true;
videoEl.currentTime=0;
this.showPlayButton();
}catch(e){
console.warn('Poster generation failed:',e);
this.showPlayButton();
}
videoEl.removeEventListener('seeked',handleSeeked);
};
this.addEventListener(videoEl,'loadeddata',handleLoadedData);
this.addEventListener(videoEl,'seeked',handleSeeked);
}
initVideoJS(){
if(this.isDestroyed)return;
try{
this.player=videojs(this.uniqueId+'-player',{
controls:true,
autoplay:false,
preload:'auto',
fluid:true,
aspectRatio:'16:9',
playbackRates:[0.5,0.75,1,1.25,1.5,2],
html5:{
vhs:{
withCredentials:false,
handleManifestRedirects:true,
overrideNative:true
},
nativeVideoTracks:false,
nativeAudioTracks:false,
nativeTextTracks:false
},
controlBar:{
volumePanel:{inline:false},
children:['playToggle','volumePanel','currentTimeDisplay','timeDivider','durationDisplay','progressControl','remainingTimeDisplay','playbackRateMenuButton','pictureInPictureToggle','fullscreenToggle']
}
});
this.player.ready(()=>{
if(this.isDestroyed)return;
this.showLoading('Probing video format...');
this.probeVideoCodec().then(info=>{
if(info){
console.log('Video info:',info);
}
this.hideLoading();
});
this.hideError();
this.generatePoster();
this.bindVideoJSEvents();
this.setupBufferMonitoring();
});
this.player.on('loadstart',()=>{
this.showLoading('Loading video...');
});
this.player.on('loadedmetadata',()=>{
this.showLoading('Video ready, loading data...');
});
this.player.on('loadeddata',()=>{
this.hideLoading();
});
this.player.on('canplay',()=>{
this.hideLoading();
});
this.player.on('waiting',()=>{
if(this.loadingState!=='loading'){
this.showLoading('Buffering...');
}
});
this.player.on('canplaythrough',()=>{
this.hideLoading();
});
this.player.on('play',()=>{
this.isPlaying=true;
this.hideLoading();
this.hidePlayButton();
if(this.elements.posterCanvas)this.elements.posterCanvas.classList.add('hidden');
});
this.player.on('playing',()=>{
this.isPlaying=true;
this.hidePlayButton();
this.hideLoading();
});
this.player.on('pause',()=>{
this.isPlaying=false;
if(!this.player.seeking()&&this.player.currentTime()>0&&!this.player.ended()){
this.showPlayButton();
}
});
this.player.on('ended',()=>{
this.isPlaying=false;
if(this.elements.posterCanvas)this.elements.posterCanvas.classList.remove('hidden');
this.showPlayButton();
});
this.player.on('seeking',()=>{
this.hidePlayButton();
});
this.player.on('seeked',()=>{
if(!this.isPlaying&&this.player.paused()){
this.showPlayButton();
}
});
this.player.on('error',()=>{
const error=this.player.error();
console.error('Video error:',error);
if(error&&error.code===4&&this.supportsWebCodecs){
this.showError('Video format not supported by browser. Trying alternative method...',true);
setTimeout(()=>this.tryWebCodecsPlayback(),1000);
}else if(this.retryCount<this.maxRetries){
this.showError((error?error.message:'Failed to load video')+'. Click retry to try again.',true);
}else{
this.showError(error?error.message:'Failed to load video',false);
}
});
}catch(error){
console.error('VideoJS init failed:',error);
this.fallbackToNative();
}
}
setupBufferMonitoring(){
if(!this.player)return;
const videoElement=this.player.el().querySelector('video');
if(!videoElement)return;
this.addEventListener(videoElement,'progress',()=>{
if(videoElement.buffered.length>0&&videoElement.duration){
const bufferedEnd=videoElement.buffered.end(videoElement.buffered.length-1);
const percent=(bufferedEnd/videoElement.duration)*100;
if(percent<100&&this.loadingState==='loading'){
this.updateLoadingProgress(percent);
}
}
});
}
async tryWebCodecsPlayback(){
console.log('Attempting WebCodecs playback...');
this.showLoading('Trying advanced codec support...');
try{
const response=await fetch(this.videoUrl);
if(!response.ok)throw new Error('Failed to fetch video');
this.showError('WebCodecs playback not fully implemented. Please use a compatible video format.',false);
}catch(error){
console.error('WebCodecs playback failed:',error);
this.showError('All playback methods failed. Video format may not be supported.',false);
}
}
fallbackToNative(){
if(this.isDestroyed)return;
console.log('Using native HTML5 player');
this.showLoading('Initializing native player...');
if(!this.elements.videoEl)return;
this.elements.videoEl.style.display='block';
this.elements.videoEl.controls=true;
this.addEventListener(this.elements.videoEl,'loadstart',()=>{
this.showLoading('Loading video...');
});
this.addEventListener(this.elements.videoEl,'loadedmetadata',()=>{
this.showLoading('Video ready...');
this.hideError();
this.probeVideoCodec();
this.generatePoster();
});
this.addEventListener(this.elements.videoEl,'loadeddata',()=>{
this.hideLoading();
});
this.addEventListener(this.elements.videoEl,'waiting',()=>{
this.showLoading('Buffering...');
});
this.addEventListener(this.elements.videoEl,'canplay',()=>{
this.hideLoading();
this.hideError();
if(!this.isPlaying){
this.showPlayButton();
}
});
this.addEventListener(this.elements.videoEl,'canplaythrough',()=>{
this.hideLoading();
});
this.addEventListener(this.elements.videoEl,'progress',()=>{
const videoEl=this.elements.videoEl;
if(videoEl.buffered.length>0&&videoEl.duration){
const bufferedEnd=videoEl.buffered.end(videoEl.buffered.length-1);
const percent=(bufferedEnd/videoEl.duration)*100;
if(percent<100&&this.loadingState==='loading'){
this.updateLoadingProgress(percent);
}
}
});
this.addEventListener(this.elements.videoEl,'error',(e)=>{
console.error('Native video error:',e);
const error=this.elements.videoEl.error;
let errorMessage='Failed to load video';
if(error){
switch(error.code){
case 1:errorMessage='Video loading aborted';break;
case 2:errorMessage='Network error while loading video';break;
case 3:errorMessage='Video decoding failed';break;
case 4:errorMessage='Video format not supported';break;
}
}
if(this.retryCount<this.maxRetries){
this.showError(errorMessage+'. Click retry to try again.',true);
}else{
this.showError(errorMessage,false);
}
});
this.addEventListener(this.elements.videoEl,'play',()=>{
this.isPlaying=true;
this.hideLoading();
this.hidePlayButton();
if(this.elements.posterCanvas)this.elements.posterCanvas.classList.add('hidden');
});
this.addEventListener(this.elements.videoEl,'playing',()=>{
this.isPlaying=true;
this.hidePlayButton();
});
this.addEventListener(this.elements.videoEl,'pause',()=>{
this.isPlaying=false;
if(this.elements.videoEl.currentTime>0&&!this.elements.videoEl.ended){
this.showPlayButton();
}
});
this.addEventListener(this.elements.videoEl,'ended',()=>{
this.isPlaying=false;
if(this.elements.posterCanvas)this.elements.posterCanvas.classList.remove('hidden');
this.showPlayButton();
});
this.addEventListener(this.elements.videoEl,'seeking',()=>{
this.hidePlayButton();
});
this.addEventListener(this.elements.videoEl,'seeked',()=>{
if(!this.isPlaying&&this.elements.videoEl.paused){
this.showPlayButton();
}
});
this.bindNativeEvents();
}
bindVideoJSEvents(){
if(this.elements.posterCanvas){
this.addEventListener(this.elements.posterCanvas,'click',()=>{
if(this.player){
this.player.play();
}
});
}
if(this.elements.playBtn){
this.addEventListener(this.elements.playBtn,'click',()=>{
if(this.player){
if(this.player.paused()){
this.player.play();
}else{
this.player.pause();
}
}
});
}
if(this.elements.pipBtn){
this.addEventListener(this.elements.pipBtn,'click',(e)=>{
e.stopPropagation();
if(!this.player)return;
const videoElement=this.player.el().querySelector('video');
if(document.pictureInPictureEnabled&&videoElement){
if(document.pictureInPictureElement===videoElement){
document.exitPictureInPicture();
}else{
videoElement.requestPictureInPicture().catch(err=>console.error('PiP error:',err));
}
}
});
}
}
bindNativeEvents(){
if(this.elements.posterCanvas){
this.addEventListener(this.elements.posterCanvas,'click',()=>{
if(this.elements.videoEl){
this.elements.videoEl.play();
}
});
}
if(this.elements.playBtn){
this.addEventListener(this.elements.playBtn,'click',()=>{
if(this.elements.videoEl){
if(this.elements.videoEl.paused){
this.elements.videoEl.play();
}else{
this.elements.videoEl.pause();
}
}
});
}
if(this.elements.pipBtn){
this.addEventListener(this.elements.pipBtn,'click',(e)=>{
e.stopPropagation();
if(document.pictureInPictureEnabled&&this.elements.videoEl){
if(document.pictureInPictureElement===this.elements.videoEl){
document.exitPictureInPicture();
}else{
this.elements.videoEl.requestPictureInPicture().catch(err=>console.error('PiP error:',err));
}
}
});
}
}
destroy(){
this.isDestroyed=true;
this.removeAllEventListeners();
window.__VideoPlayerManager.unregisterInstance(this.uniqueId);
}
}
const player=new VideoPlayer(UNIQUE_ID,VIDEO_URL,VIDEO_TYPE);
window.addEventListener('beforeunload',()=>player.destroy());
})();
</script>`, uniqueId, safeVideoPath, videoType)
}

// getVideoMimeType 根据文件扩展名返回对应的 MIME 类型
func getVideoMimeType(ext string) string {
	mimeTypes := map[string]string{
		".mp4":  "video/mp4",
		".webm": "video/webm",
		".ogg":  "video/ogg",
		".ogv":  "video/ogg",
		".mov":  "video/quicktime",
		".avi":  "video/x-msvideo",
		".wmv":  "video/x-ms-wmv",
		".flv":  "video/x-flv",
		".mkv":  "video/x-matroska",
		".m4v":  "video/x-m4v",
		".3gp":  "video/3gpp",
		".3g2":  "video/3gpp2",
		".mpg":  "video/mpeg",
		".mpeg": "video/mpeg",
		".ts":   "video/mp2t",
		".m3u8": "application/x-mpegURL",
	}
	if mimeType, ok := mimeTypes[ext]; ok {
		return mimeType
	}
	return "video/mp4"
}
