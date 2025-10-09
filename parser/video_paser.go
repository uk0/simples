package parser

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"log"
	"path/filepath"
	"strings"
)

// GenerateVideoWarpHTML 生成用于在Web上展示视频内容的模块
// 使用原生HTML5视频播放器+Video.js增强，支持多种视频格式
// 自动提取视频封面，居中显示，简洁设计，完美全屏支持
func GenerateVideoWarpHTML(videoId, videoPath string) string {
	safeVideoId := template.HTMLEscapeString(videoId)
	safeVideoPath := template.HTMLEscapeString(videoPath)
	log.Println("Generating Video player for FileName:", safeVideoPath)
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
#%[1]s-wrapper .video-poster.hidden{opacity:0;pointer-events:none}
#%[1]s-wrapper .play-button-overlay{position:absolute;top:50%%;left:50%%;transform:translate(-50%%,-50%%);width:80px;height:80px;background:rgba(0,0,0,0.7);border:3px solid #fff;border-radius:50%%;cursor:pointer;z-index:2;display:flex;align-items:center;justify-content:center;transition:all 0.3s ease}
#%[1]s-wrapper .play-button-overlay:hover{background:rgba(52,152,219,0.9);transform:translate(-50%%,-50%%) scale(1.1)}
#%[1]s-wrapper .play-button-overlay.hidden{display:none}
#%[1]s-wrapper .play-button-overlay::after{content:'';width:0;height:0;border-left:25px solid #fff;border-top:15px solid transparent;border-bottom:15px solid transparent;margin-left:5px}
#%[1]s-wrapper .vjs-big-play-button{display:none!important}
#%[1]s-wrapper .vjs-play-progress,#%[1]s-wrapper .vjs-volume-level{background-color:#3498db}
#%[1]s-wrapper .video-loading{position:absolute;top:50%%;left:50%%;transform:translate(-50%%,-50%%);text-align:center;color:#fff;font-size:14px;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,sans-serif;z-index:3}
#%[1]s-wrapper .video-spinner{border:3px solid rgba(255,255,255,0.2);border-top:3px solid #3498db;border-radius:50%%;width:40px;height:40px;animation:video-spin-%[1]s 1s linear infinite;margin:0 auto 15px}
@keyframes video-spin-%[1]s{0%%{transform:rotate(0deg)}100%%{transform:rotate(360deg)}}
#%[1]s-wrapper .video-error{position:absolute;top:50%%;left:50%%;transform:translate(-50%%,-50%%);text-align:center;padding:20px;color:#e74c3c;font-size:14px;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,sans-serif;background:rgba(255,255,255,0.9);border-radius:8px;z-index:10}
#%[1]s-wrapper .video-toolbar{position:absolute;top:10px;right:10px;display:flex;gap:8px;opacity:0;transition:opacity 0.3s ease;z-index:10}
#%[1]s-wrapper:hover .video-toolbar{opacity:1}
#%[1]s-wrapper .video-btn{background:rgba(0,0,0,0.6);backdrop-filter:blur(10px);border:none;color:#fff;padding:8px 12px;border-radius:6px;cursor:pointer;font-size:12px;transition:all 0.2s;text-decoration:none;display:inline-flex;align-items:center;gap:5px;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,sans-serif}
#%[1]s-wrapper .video-btn:hover{background:rgba(0,0,0,0.8);transform:translateY(-1px)}
@media (max-width:768px){
#%[1]s-wrapper{margin:10px}
#%[1]s-wrapper .video-toolbar{opacity:1}
#%[1]s-wrapper .play-button-overlay{width:60px;height:60px}
#%[1]s-wrapper .play-button-overlay::after{border-left-width:20px;border-top-width:12px;border-bottom-width:12px}
}
</style>
<div class="video-container">
<div class="video-inner">
<div class="video-loading" id="%[1]s-loading"><div class="video-spinner"></div><div>Loading video...</div></div>
<canvas id="%[1]s-poster" class="video-poster" style="display:none;"></canvas>
<div id="%[1]s-play-btn" class="play-button-overlay"></div>
<div class="video-toolbar">
<button class="video-btn" id="%[1]s-pip-btn" title="Picture in Picture">📺 PiP</button>
<a href="%[2]s" target="_blank" rel="noopener noreferrer" class="video-btn" title="Open in new tab">🔗 Open</a>
</div>
<video id="%[1]s-player" class="video-js vjs-default-skin vjs-16-9" controls preload="metadata">
<source src="%[2]s" type="%[3]s">
<p class="vjs-no-js">To view this video please enable JavaScript</p>
</video>
<div class="video-error" id="%[1]s-error" style="display:none;">❌ Failed to load video</div>
</div>
</div>
</div>
<link href="https://cdn.jsdelivr.net/npm/video.js@8.10.0/dist/video-js.min.css" rel="stylesheet">
<script src="https://cdn.jsdelivr.net/npm/video.js@8.10.0/dist/video.min.js"></script>
<script>
(function(){
const UNIQUE_ID='%[1]s';
const VIDEO_URL='%[2]s';
const elements={wrapper:document.getElementById(UNIQUE_ID+'-wrapper'),loading:document.getElementById(UNIQUE_ID+'-loading'),error:document.getElementById(UNIQUE_ID+'-error'),videoEl:document.getElementById(UNIQUE_ID+'-player'),posterCanvas:document.getElementById(UNIQUE_ID+'-poster'),playBtn:document.getElementById(UNIQUE_ID+'-play-btn'),pipBtn:document.getElementById(UNIQUE_ID+'-pip-btn')};
let player=null;
let posterGenerated=false;
function generatePoster(){
if(posterGenerated||!elements.videoEl)return;
const videoEl=elements.videoEl;
const canvas=elements.posterCanvas;
const handleLoadedData=function(){
if(posterGenerated)return;
setTimeout(function(){
videoEl.currentTime=Math.min(videoEl.duration*0.1,5);
},100);
};
const handleSeeked=function(){
if(posterGenerated)return;
try{
canvas.width=videoEl.videoWidth;
canvas.height=videoEl.videoHeight;
const ctx=canvas.getContext('2d');
ctx.drawImage(videoEl,0,0,canvas.width,canvas.height);
canvas.style.display='block';
posterGenerated=true;
videoEl.currentTime=0;
console.log('Video poster generated');
}catch(e){
console.error('Poster generation error:',e);
}
videoEl.removeEventListener('seeked',handleSeeked);
};
videoEl.addEventListener('loadeddata',handleLoadedData,{once:true});
videoEl.addEventListener('seeked',handleSeeked);
}
function initPlayer(){
try{
player=videojs(UNIQUE_ID+'-player',{controls:true,autoplay:false,preload:'metadata',fluid:true,aspectRatio:'16:9',playbackRates:[0.5,0.75,1,1.25,1.5,2],controlBar:{volumePanel:{inline:false},children:['playToggle','volumePanel','currentTimeDisplay','timeDivider','durationDisplay','progressControl','remainingTimeDisplay','playbackRateMenuButton','pictureInPictureToggle','fullscreenToggle']}});
player.ready(function(){
console.log('Video.js player ready');
elements.loading.style.display='none';
generatePoster();
bindEvents();
});
player.on('play',function(){
if(elements.posterCanvas)elements.posterCanvas.classList.add('hidden');
if(elements.playBtn)elements.playBtn.classList.add('hidden');
});
player.on('pause',function(){
if(!player.seeking()&&player.currentTime()>0&&!player.ended()&&!player.isFullscreen()){
if(elements.playBtn)elements.playBtn.classList.remove('hidden');
}
});
player.on('ended',function(){
if(elements.posterCanvas)elements.posterCanvas.classList.remove('hidden');
if(elements.playBtn)elements.playBtn.classList.remove('hidden');
});
player.on('error',function(){
const error=player.error();
console.error('Video player error:',error);
elements.loading.style.display='none';
elements.error.style.display='block';
if(error)elements.error.innerHTML='❌ Failed to load video<br><small>'+error.message+'</small>';
});
player.on('fullscreenchange',function(){
const isFS=player.isFullscreen();
console.log(isFS?'Entered fullscreen':'Exited fullscreen');
});
}catch(error){
console.error('Failed to initialize player:',error);
elements.loading.style.display='none';
elements.error.style.display='block';
elements.error.innerHTML='❌ Failed to initialize player<br><small>'+error.message+'</small>';
}
}
function bindEvents(){
if(elements.posterCanvas){
elements.posterCanvas.addEventListener('click',function(){if(player)player.play();});
}
if(elements.playBtn){
elements.playBtn.addEventListener('click',function(){if(player){if(player.paused())player.play();else player.pause();}});
}
if(elements.pipBtn){
elements.pipBtn.addEventListener('click',function(e){
e.stopPropagation();
if(!player)return;
const videoElement=player.el().querySelector('video');
if(document.pictureInPictureEnabled&&videoElement){
if(document.pictureInPictureElement===videoElement)document.exitPictureInPicture();
else videoElement.requestPictureInPicture().catch(err=>console.error('PiP error:',err));
}else alert('Picture-in-Picture is not supported');
});
}
document.addEventListener('keydown',function(e){
if(!player||!elements.wrapper)return;
const activeElement=document.activeElement;
if(activeElement&&(activeElement.tagName==='INPUT'||activeElement.tagName==='TEXTAREA'||activeElement.isContentEditable))return;
const rect=elements.wrapper.getBoundingClientRect();
const inViewport=rect.top<window.innerHeight&&rect.bottom>0;
if(!inViewport&&!player.isFullscreen())return;
switch(e.key){
case ' ':
case 'k':
if(player.paused())player.play();else player.pause();
e.preventDefault();
break;
case 'f':
if(player.isFullscreen())player.exitFullscreen();else player.requestFullscreen();
e.preventDefault();
break;
case 'Escape':
if(player.isFullscreen())player.exitFullscreen();
break;
case 'm':
player.muted(!player.muted());
e.preventDefault();
break;
case 'ArrowLeft':
player.currentTime(Math.max(0,player.currentTime()-5));
e.preventDefault();
break;
case 'ArrowRight':
player.currentTime(Math.min(player.duration()||0,player.currentTime()+5));
e.preventDefault();
break;
}
});
}
function waitForVideoJS(){
if(typeof videojs==='function')initPlayer();
else{console.log('Waiting for Video.js...');setTimeout(waitForVideoJS,100);}
}
if(document.readyState==='loading')document.addEventListener('DOMContentLoaded',waitForVideoJS);
else waitForVideoJS();
window.addEventListener('beforeunload',function(){if(player){try{player.dispose();}catch(e){console.error('Error disposing player:',e);}}});
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
