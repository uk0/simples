package parser

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"path/filepath"
	"strings"
)

// GenerateVideoWarpHTML 企业级视频播放器
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
#%[1]s-wrapper .vjs-load-progress{background:rgba(255,255,255,0.3)!important}
#%[1]s-wrapper .vjs-play-progress{background:#3498db!important}
#%[1]s-wrapper .video-poster{position:absolute;top:0;left:0;width:100%%;height:100%%;object-fit:cover;z-index:1;cursor:pointer;transition:opacity 0.3s ease;border-radius:8px}
#%[1]s-wrapper .video-poster.hidden{opacity:0;pointer-events:none;z-index:-1}
#%[1]s-wrapper .play-button-overlay{position:absolute;top:50%%;left:50%%;transform:translate(-50%%,-50%%);width:80px;height:80px;background:rgba(0,0,0,0.7);border:3px solid #fff;border-radius:50%%;cursor:pointer;z-index:2;display:none;align-items:center;justify-content:center;transition:all 0.3s ease}
#%[1]s-wrapper .play-button-overlay.show{display:flex}
#%[1]s-wrapper .play-button-overlay:hover{background:rgba(52,152,219,0.9);transform:translate(-50%%,-50%%) scale(1.1)}
#%[1]s-wrapper .play-button-overlay::after{content:'';width:0;height:0;border-left:25px solid #fff;border-top:15px solid transparent;border-bottom:15px solid transparent;margin-left:5px}
#%[1]s-wrapper .vjs-big-play-button{display:none!important}
#%[1]s-wrapper .vjs-play-progress,#%[1]s-wrapper .vjs-volume-level{background-color:#3498db}
#%[1]s-wrapper .video-loading{position:absolute;top:50%%;left:50%%;transform:translate(-50%%,-50%%);text-align:center;color:#fff;font-size:14px;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,sans-serif;z-index:15;background:rgba(0,0,0,0.8);padding:20px;border-radius:8px;min-width:150px;display:none}
#%[1]s-wrapper .video-spinner{border:3px solid rgba(255,255,255,0.2);border-top:3px solid #3498db;border-radius:50%%;width:40px;height:40px;animation:video-spin-%[1]s 1s linear infinite;margin:0 auto 15px}
@keyframes video-spin-%[1]s{0%%{transform:rotate(0deg)}100%%{transform:rotate(360deg)}}
#%[1]s-wrapper .video-loading-progress{margin-top:10px;height:4px;background:rgba(255,255,255,0.2);border-radius:2px;overflow:hidden;display:none}
#%[1]s-wrapper .video-loading-progress-bar{height:100%%;background:#3498db;width:0%%;transition:width 0.3s ease}
#%[1]s-wrapper .video-error{position:absolute;top:50%%;left:50%%;transform:translate(-50%%,-50%%);text-align:center;padding:20px;color:#e74c3c;font-size:14px;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,sans-serif;background:rgba(255,255,255,0.95);border-radius:8px;box-shadow:0 2px 10px rgba(0,0,0,0.1);z-index:10;max-width:80%%;display:none}
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
<div class="video-loading-progress" id="%[1]s-loading-progress">
<div class="video-loading-progress-bar" id="%[1]s-loading-progress-bar"></div>
</div>
</div>
<canvas id="%[1]s-poster" class="video-poster" style="display:none;"></canvas>
<div id="%[1]s-play-btn" class="play-button-overlay"></div>
<div class="video-toolbar">
<button class="video-btn" id="%[1]s-pip-btn" title="Picture in Picture">PiP</button>
<a href="%[2]s" target="_blank" rel="noopener noreferrer" class="video-btn" title="Open in new tab">Open</a>
</div>
<div class="video-codec-info" id="%[1]s-codec-info">Loading...</div>
<video id="%[1]s-player" class="video-js vjs-default-skin vjs-16-9" controls preload="metadata" crossorigin="anonymous">
<source src="%[2]s" type="%[3]s">
<p class="vjs-no-js">To view this video please enable JavaScript</p>
</video>
<div class="video-error" id="%[1]s-error">
<div class="video-error-title">Failed to load video</div>
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
instances:{},videojsLoaded:false,videojsLoading:false,loadCallbacks:[],keyboardHandlerAttached:false,
loadVideoJS:function(cb){
if(this.videojsLoaded){cb(true);return;}
this.loadCallbacks.push(cb);
if(this.videojsLoading)return;
this.videojsLoading=true;
const css=document.querySelector('link[href*="video-js.min.css"]')||(()=>{const l=document.createElement('link');l.rel='stylesheet';l.href='https://cdn.jsdelivr.net/npm/video.js@8.10.0/dist/video-js.min.css';document.head.appendChild(l);return l;})();
const script=document.querySelector('script[src*="video.min.js"]')||(()=>{const s=document.createElement('script');s.src='https://cdn.jsdelivr.net/npm/video.js@8.10.0/dist/video.min.js';s.onload=()=>{this.videojsLoaded=true;this.videojsLoading=false;this.loadCallbacks.forEach(c=>c(true));this.loadCallbacks=[];};s.onerror=()=>{this.videojsLoading=false;this.loadCallbacks.forEach(c=>c(false));this.loadCallbacks=[];};document.head.appendChild(s);return s;})();
if(typeof videojs!=='undefined'){this.videojsLoaded=true;this.videojsLoading=false;this.loadCallbacks.forEach(c=>c(true));this.loadCallbacks=[];}
},
attachKeyboardHandler:function(){
if(this.keyboardHandlerAttached)return;
this.keyboardHandlerAttached=true;
document.addEventListener('keydown',e=>{
const a=document.activeElement;
if(a&&(a.tagName==='INPUT'||a.tagName==='TEXTAREA'||a.isContentEditable))return;
let t=null;
for(const i in this.instances){
const inst=this.instances[i];
if(!inst.player||!inst.elements.wrapper)continue;
if(inst.player.isFullscreen()){t=inst.player;break;}
}
if(!t){
for(const i in this.instances){
const inst=this.instances[i];
if(!inst.player||!inst.elements.wrapper)continue;
const r=inst.elements.wrapper.getBoundingClientRect();
if(r.top<window.innerHeight&&r.bottom>0){t=inst.player;break;}
}
}
if(!t)return;
switch(e.key){
case ' ':case 'k':t.paused()?t.play():t.pause();e.preventDefault();break;
case 'f':t.isFullscreen()?t.exitFullscreen():t.requestFullscreen();e.preventDefault();break;
case 'Escape':t.isFullscreen()&&t.exitFullscreen();break;
case 'm':t.muted(!t.muted());e.preventDefault();break;
case 'ArrowLeft':t.currentTime(Math.max(0,t.currentTime()-5));e.preventDefault();break;
case 'ArrowRight':t.currentTime(Math.min(t.duration()||0,t.currentTime()+5));e.preventDefault();break;
}
});
},
registerInstance:function(id,inst){this.instances[id]=inst;this.attachKeyboardHandler();},
unregisterInstance:function(id){
if(this.instances[id]){
if(this.instances[id].player)try{this.instances[id].player.dispose();}catch(e){}
delete this.instances[id];
}
}
};
}
class VideoPlayer{
constructor(uid,url,type){
this.uniqueId=uid;this.videoUrl=url;this.videoType=type;
this.player=null;this.posterGenerated=false;this.retryCount=0;this.maxRetries=3;
this.isDestroyed=false;this.eventListeners=[];this.loadingState='init';
this.isPlaying=false;this.bufferingTimeout=null;
this.elements={
wrapper:document.getElementById(uid+'-wrapper'),loading:document.getElementById(uid+'-loading'),
loadingText:document.getElementById(uid+'-loading-text'),loadingProgress:document.getElementById(uid+'-loading-progress'),
loadingProgressBar:document.getElementById(uid+'-loading-progress-bar'),error:document.getElementById(uid+'-error'),
errorMsg:document.getElementById(uid+'-error-msg'),retryBtn:document.getElementById(uid+'-retry'),
videoEl:document.getElementById(uid+'-player'),posterCanvas:document.getElementById(uid+'-poster'),
playBtn:document.getElementById(uid+'-play-btn'),pipBtn:document.getElementById(uid+'-pip-btn'),
codecInfo:document.getElementById(uid+'-codec-info')
};
this.init();
}
init(){
window.__VideoPlayerManager.registerInstance(this.uniqueId,this);
if(this.elements.retryBtn)this.addEventListener(this.elements.retryBtn,'click',()=>this.retry());
this.loadPlayer();
}
addEventListener(el,ev,h){if(el){el.addEventListener(ev,h);this.eventListeners.push({el,ev,h});}}
removeAllEventListeners(){this.eventListeners.forEach(o=>o.el.removeEventListener(o.ev,o.h));this.eventListeners=[];}
showLoading(t){
if(this.loadingState==='loading')return;
this.loadingState='loading';
this.elements.loading.style.display='block';
if(this.elements.loadingText)this.elements.loadingText.textContent=t||'Loading...';
if(this.elements.loadingProgress)this.elements.loadingProgress.style.display='none';
clearTimeout(this.bufferingTimeout);
this.bufferingTimeout=setTimeout(()=>{if(this.loadingState==='loading')this.showError('Buffering timeout.',true);},30000);
}
updateLoadingProgress(p){
if(!this.elements.loadingProgress)return;
this.elements.loadingProgress.style.display='block';
this.elements.loadingProgressBar.style.width=p+'%';
if(this.elements.loadingText)this.elements.loadingText.textContent='Buffering... '+Math.round(p)+'%';
}
hideLoading(){
if(this.loadingState!=='loading')return;
this.loadingState='loaded';
clearTimeout(this.bufferingTimeout);this.bufferingTimeout=null;
this.elements.loading.style.display='none';
if(this.elements.loadingProgress)this.elements.loadingProgress.style.display='none';
}
showPlayButton(){
    if(this.isPlaying) return; // 播放中禁止显示
    if(this.elements.playBtn) this.elements.playBtn.classList.add('show');
}
hidePlayButton(){if(this.elements.playBtn)this.elements.playBtn.classList.remove('show');this.elements.playBtn.classList.add('hide');}
showError(m,r=true){this.hideLoading();this.elements.error.style.display='block';if(this.elements.errorMsg)this.elements.errorMsg.textContent=m;if(this.elements.retryBtn)this.elements.retryBtn.style.display=r?'inline-block':'none';}
hideError(){if(this.elements.error)this.elements.error.style.display='none';}
retry(){
if(this.retryCount>=this.maxRetries){this.showError('Max retries reached.',false);return;}
this.retryCount++;this.hideError();this.showLoading('Retrying... ('+this.retryCount+'/'+this.maxRetries+')');
if(this.player)try{this.player.dispose();}catch(e){}
this.player=null;this.posterGenerated=false;this.isPlaying=false;
setTimeout(()=>this.loadPlayer(),1000);
}
loadPlayer(){
this.showLoading('Initializing player...');
window.__VideoPlayerManager.loadVideoJS(s=>{if(this.isDestroyed)return;s?this.initVideoJS():this.fallbackToNative();});
}
generatePoster(){
if(this.posterGenerated||!this.elements.videoEl)return;
const v=this.elements.videoEl,c=this.elements.posterCanvas;
const onData=()=>{if(this.posterGenerated||v.duration<=0)return;setTimeout(()=>{v.currentTime=Math.min(v.duration*0.1,5);},100);};
const onSeek=()=>{
if(this.posterGenerated)return;
try{c.width=v.videoWidth||640;c.height=v.videoHeight||360;
const ctx=c.getContext('2d');ctx.drawImage(v,0,0,c.width,c.height);c.style.display='block';
this.posterGenerated=true;v.currentTime=0;this.showPlayButton();this.hideLoading();
}catch(e){console.warn('Poster failed:',e);this.showPlayButton();this.hideLoading();}
v.removeEventListener('seeked',onSeek);
};
this.addEventListener(v,'loadeddata',onData);
this.addEventListener(v,'seeked',onSeek);
}
initVideoJS(){
if(this.isDestroyed)return;
try{
this.player=videojs(this.uniqueId+'-player',{controls:true,autoplay:false,preload:'metadata',fluid:true,aspectRatio:'16:9',playbackRates:[0.5,0.75,1,1.25,1.5,2],html5:{vhs:{overrideNative:true}},controlBar:{volumePanel:{inline:false}}});
this.player.ready(()=>{
if(this.isDestroyed)return;
this.showLoading('Probing format...');
this.probeVideoCodec().then(()=>{this.hideLoading();});
this.hideError();this.generatePoster();this.bindVideoJSEvents();this.setupBufferMonitoring();
});
this.player.on('loadeddata',()=>{this.hideLoading();this.showPlayButton();});
this.player.on('canplay',()=>{this.hideLoading();this.showPlayButton();});
this.player.on('waiting',()=>{if(this.isPlaying)this.showLoading('Buffering...');});
this.player.on('canplaythrough',()=>{this.hideLoading();this.showPlayButton();});
this.player.on('play',()=>{this.isPlaying=true;this.hideLoading();this.hidePlayButton();if(this.elements.posterCanvas)this.elements.posterCanvas.classList.add('hidden');});
this.player.on('playing',()=>{this.isPlaying=true;this.hidePlayButton();this.hideLoading();});
this.player.on('pause',()=>{this.isPlaying=false;if(this.player.currentTime()>0&&!this.player.ended())this.showPlayButton();});
this.player.on('ended',()=>{this.isPlaying=false;if(this.elements.posterCanvas)this.elements.posterCanvas.classList.remove('hidden');this.showPlayButton();});
this.player.on('seeking',()=>{const t=this.player.currentTime(),b=this.player.buffered();if(!this.isTimeBuffered(t,b))this.showLoading('Seeking...');this.hidePlayButton();});
// === initVideoJS seeked 事件替换 ===
this.player.on('seeked',()=>{
    const t=this.player.currentTime(),b=this.player.buffered();
    const a=this.getBufferAhead(t,b);
    if(!this.isTimeBuffered(t,b)||a<2){
        this.showLoading('Waiting for buffer...');
    }else{
        this.hideLoading(); // 只隐藏缓冲遮罩
        this.hidePlayButton();
    }
});
this.player.on('error',()=>{const e=this.player.error();this.showError((e&&e.message)||'Failed to load video',this.retryCount<this.maxRetries);});
}catch(e){console.error('VideoJS init failed:',e);this.fallbackToNative();}
}
setupBufferMonitoring(){
if(!this.player)return;
const v=this.player.el().querySelector('video');
if(!v)return;
this.addEventListener(v,'progress',()=>{if(v.buffered.length>0&&v.duration){const p=(v.buffered.end(v.buffered.length-1)/v.duration)*100;if(p<100&&this.loadingState==='loading')this.updateLoadingProgress(p);}});
}
fallbackToNative(){
if(this.isDestroyed)return;
this.showLoading('Initializing native player...');
if(!this.elements.videoEl)return;
const bar=this.elements.videoEl.parentElement.querySelector('.vjs-control-bar');
if(bar)bar.style.display='none';
this.elements.videoEl.style.display='block';this.elements.videoEl.controls=true;
this.addEventListener(this.elements.videoEl,'loadedmetadata',()=>{this.hideError();this.probeVideoCodec();this.generatePoster();});
this.addEventListener(this.elements.videoEl,'loadeddata',()=>{this.hideLoading();this.showPlayButton();});
this.addEventListener(this.elements.videoEl,'canplay',()=>{this.hideLoading();this.showPlayButton();});
this.addEventListener(this.elements.videoEl,'waiting',()=>{if(this.isPlaying)this.showLoading('Buffering...');});
this.addEventListener(this.elements.videoEl,'progress',()=>{const v=this.elements.videoEl;if(v.buffered.length>0&&v.duration){const p=(v.buffered.end(v.buffered.length-1)/v.duration)*100;if(p<100&&this.loadingState==='loading')this.updateLoadingProgress(p);}});
this.addEventListener(this.elements.videoEl,'play',()=>{this.isPlaying=true;this.hideLoading();this.hidePlayButton();if(this.elements.posterCanvas)this.elements.posterCanvas.classList.add('hidden');});
this.addEventListener(this.elements.videoEl,'playing',()=>{this.isPlaying=true;this.hidePlayButton();this.hideLoading();});
this.addEventListener(this.elements.videoEl,'pause',()=>{this.isPlaying=false;if(this.elements.videoEl.currentTime>0&&!this.elements.videoEl.ended)this.showPlayButton();});
this.addEventListener(this.elements.videoEl,'ended',()=>{this.isPlaying=false;if(this.elements.posterCanvas)this.elements.posterCanvas.classList.remove('hidden');this.showPlayButton();});
let last=0;
const loop=()=>{
if(this.isDestroyed||!this.elements.videoEl)return;
const v=this.elements.videoEl,ct=v.currentTime;
if(Math.abs(ct-last)>1&&!v.paused&&!v.seeking){
const b=v.buffered;
if(!this.isTimeBuffered(ct,b))this.showLoading('Seeking...');
this.hidePlayButton();
}
last=ct;
// === fallbackToNative loop 替换 ===
if(this.loadingState==='loading'&&this.isPlaying){
    const b=v.buffered;
    if(this.isTimeBuffered(ct,b)&&this.getBufferAhead(ct,b)>=1){
        this.hideLoading(); // 只隐藏缓冲遮罩
        this.hidePlayButton(); // 强制隐藏
    }
}
requestAnimationFrame(loop);
};
requestAnimationFrame(loop);
this.bindNativeEvents();
}
bindVideoJSEvents(){
if(this.elements.posterCanvas)this.addEventListener(this.elements.posterCanvas,'click',()=>{if(this.player)this.player.play();});
if(this.elements.playBtn)this.addEventListener(this.elements.playBtn,'click',()=>{if(this.player)this.player.paused()?this.player.play():this.player.pause();});
if(this.elements.pipBtn)this.addEventListener(this.elements.pipBtn,'click',e=>{e.stopPropagation();if(!this.player)return;const v=this.player.el().querySelector('video');if(document.pictureInPictureEnabled&&v){document.pictureInPictureElement===v?document.exitPictureInPicture():v.requestPictureInPicture().catch(()=>{});}});
}
bindNativeEvents(){
if(this.elements.posterCanvas)this.addEventListener(this.elements.posterCanvas,'click',()=>{if(this.elements.videoEl)this.elements.videoEl.play();});
if(this.elements.playBtn)this.addEventListener(this.elements.playBtn,'click',()=>{if(this.elements.videoEl)this.elements.videoEl.paused()?this.elements.videoEl.play():this.elements.videoEl.pause();});
if(this.elements.pipBtn)this.addEventListener(this.elements.pipBtn,'click',e=>{e.stopPropagation();if(document.pictureInPictureEnabled&&this.elements.videoEl){document.pictureInPictureElement===this.elements.videoEl?document.exitPictureInPicture():this.elements.videoEl.requestPictureInPicture().catch(()=>{});}});
}
async probeVideoCodec(){
if(!this.elements.videoEl)return null;
return new Promise(r=>{
const v=this.elements.videoEl;
const onMeta=()=>{v.removeEventListener('loadedmetadata',onMeta);fetch(VIDEO_URL,{method:'HEAD'}).then(res=>{r({supportsRange:res.headers.get('Accept-Ranges')==='bytes'});}).catch(()=>r({}));};
v.addEventListener('loadedmetadata',onMeta);
setTimeout(()=>{v.removeEventListener('loadedmetadata',onMeta);r({});},3000);
});
}
isTimeBuffered(t,b){if(!b||b.length===0)return false;for(let i=0;i<b.length;i++)if(t>=b.start(i)&&t<=b.end(i))return true;return false;}
getBufferAhead(t,b){if(!b||b.length===0)return 0;for(let i=0;i<b.length;i++)if(t>=b.start(i)&&t<=b.end(i))return b.end(i)-t;return 0;}
destroy(){
this.isDestroyed=true;
clearTimeout(this.bufferingTimeout);
this.removeAllEventListeners();
window.__VideoPlayerManager.unregisterInstance(this.uniqueId);
}
}
const player=new VideoPlayer(UNIQUE_ID,VIDEO_URL,VIDEO_TYPE);
window.addEventListener('beforeunload',()=>player.destroy());
})();
</script>`, uniqueId, safeVideoPath, videoType)
}
func getVideoMimeType(ext string) string {
	m := map[string]string{
		".mp4": "video/mp4", ".webm": "video/webm", ".ogg": "video/ogg", ".ogv": "video/ogg",
		".mov": "video/quicktime", ".avi": "video/x-msvideo", ".wmv": "video/x-ms-wmv",
		".flv": "video/x-flv", ".mkv": "video/x-matroska", ".m4v": "video/x-m4v",
		".3gp": "video/3gpp", ".3g2": "video/3gpp2", ".mpg": "video/mpeg", ".mpeg": "video/mpeg",
		".ts": "video/mp2t", ".m3u8": "application/x-mpegURL",
	}
	if t, ok := m[ext]; ok {
		return t
	}
	return "video/mp4"
}
