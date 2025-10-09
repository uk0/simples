package parser

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"log"
	"path/filepath"
	"strings"
)

// GenerateAudioWarpHTML 生成用于在Web上展示音频内容的模块
// macOS Sequoia 风格透明毛玻璃设计，高对比度优化
// 原生HTML5音频，完美支持多实例
func GenerateAudioWarpHTML(audioId, audioPath string) string {
	safeAudioId := template.HTMLEscapeString(audioId)
	safeAudioPath := template.HTMLEscapeString(audioPath)
	log.Println("Generating Audio player for FileName:", safeAudioPath)
	h := md5.Sum([]byte(safeAudioId))
	uniqueId := fmt.Sprintf("audio_%x", h)[:12]
	ext := strings.ToLower(filepath.Ext(audioPath))
	audioType := getAudioMimeType(ext)
	fileName := filepath.Base(audioPath)
	fileNameWithoutExt := strings.TrimSuffix(fileName, ext)

	return fmt.Sprintf(`<div id="%[1]s-wrapper" class="audio-viewer-wrapper" style="position:relative;max-width:100%%;margin:20px auto;display:block;">
<style>
#%[1]s-wrapper{position:relative;width:100%%;max-width:600px;margin:20px auto;display:block;clear:both}
#%[1]s-wrapper .audio-container{position:relative;width:100%%;background:rgba(30,30,32,0.75);backdrop-filter:blur(40px) saturate(180%%);-webkit-backdrop-filter:blur(40px) saturate(180%%);border:1px solid rgba(255,255,255,0.12);border-radius:16px;box-shadow:0 8px 32px rgba(0,0,0,0.24),0 0 1px rgba(255,255,255,0.12) inset;padding:16px;transition:all 0.3s ease}
#%[1]s-wrapper .audio-container:hover{background:rgba(40,40,42,0.8);box-shadow:0 12px 48px rgba(0,0,0,0.32),0 0 1px rgba(255,255,255,0.18) inset}
#%[1]s-wrapper .audio-main{display:flex;align-items:center;gap:12px;margin-bottom:12px}
#%[1]s-wrapper .audio-play-btn{width:44px;height:44px;min-width:44px;border-radius:50%%;background:linear-gradient(135deg,#007AFF 0%%,#0051D5 100%%);border:none;color:#fff;font-size:16px;cursor:pointer;display:flex;align-items:center;justify-content:center;transition:all 0.2s;box-shadow:0 4px 12px rgba(0,122,255,0.4),0 1px 2px rgba(0,0,0,0.2)}
#%[1]s-wrapper .audio-play-btn:hover{background:linear-gradient(135deg,#0A84FF 0%%,#0060EA 100%%);transform:scale(1.05);box-shadow:0 6px 16px rgba(0,122,255,0.5),0 1px 2px rgba(0,0,0,0.2)}
#%[1]s-wrapper .audio-play-btn:active{transform:scale(0.98)}
#%[1]s-wrapper .audio-info{flex:1;min-width:0}
#%[1]s-wrapper .audio-title{font-size:13px;font-weight:600;color:rgba(255,255,255,0.98);margin-bottom:3px;white-space:nowrap;overflow:hidden;text-overflow:ellipsis;font-family:-apple-system,BlinkMacSystemFont,'SF Pro Text','Segoe UI',Roboto,sans-serif;text-shadow:0 1px 3px rgba(0,0,0,0.3);letter-spacing:-0.01em}
#%[1]s-wrapper .audio-format{font-size:10px;color:rgba(255,255,255,0.55);text-transform:uppercase;letter-spacing:0.8px;font-weight:500}
#%[1]s-wrapper .audio-visualizer{height:32px;display:flex;align-items:flex-end;justify-content:center;gap:2px;margin-bottom:10px;background:rgba(0,0,0,0.2);border-radius:8px;padding:4px 8px}
#%[1]s-wrapper .audio-bar{width:2px;background:rgba(10,132,255,0.7);border-radius:1px;transition:height 0.1s ease;height:6px;box-shadow:0 0 4px rgba(10,132,255,0.3)}
#%[1]s-wrapper .audio-progress{position:relative;width:100%%;height:4px;background:rgba(255,255,255,0.12);border-radius:2px;margin-bottom:6px;cursor:pointer;box-shadow:0 1px 2px rgba(0,0,0,0.2) inset}
#%[1]s-wrapper .audio-progress-bar{position:absolute;left:0;top:0;height:100%%;background:linear-gradient(90deg,#0A84FF 0%%,#007AFF 100%%);border-radius:2px;transition:width 0.1s linear;box-shadow:0 0 8px rgba(10,132,255,0.4)}
#%[1]s-wrapper .audio-time{display:flex;justify-content:space-between;font-size:10px;color:rgba(255,255,255,0.75);margin-bottom:8px;font-family:-apple-system,BlinkMacSystemFont,'SF Mono','Monaco','Menlo','Consolas',monospace;font-weight:500;letter-spacing:0.02em}
#%[1]s-wrapper .audio-controls{display:flex;align-items:center;gap:6px}
#%[1]s-wrapper .audio-control-btn{width:28px;height:28px;border-radius:50%%;background:rgba(255,255,255,0.12);border:1px solid rgba(255,255,255,0.08);color:rgba(255,255,255,0.9);font-size:11px;cursor:pointer;display:flex;align-items:center;justify-content:center;transition:all 0.2s;box-shadow:0 1px 3px rgba(0,0,0,0.15)}
#%[1]s-wrapper .audio-control-btn:hover{background:rgba(255,255,255,0.2);border-color:rgba(255,255,255,0.15);transform:scale(1.08);box-shadow:0 2px 6px rgba(0,0,0,0.2)}
#%[1]s-wrapper .audio-control-btn:active{transform:scale(0.95)}
#%[1]s-wrapper .audio-volume{display:flex;align-items:center;gap:6px;flex:1}
#%[1]s-wrapper .audio-volume-icon{color:rgba(255,255,255,0.85);font-size:14px;cursor:pointer;width:18px;text-align:center;filter:drop-shadow(0 1px 2px rgba(0,0,0,0.3));transition:all 0.2s}
#%[1]s-wrapper .audio-volume-icon:hover{color:rgba(255,255,255,1);transform:scale(1.1)}
#%[1]s-wrapper .audio-volume-slider{flex:1;height:4px;background:rgba(255,255,255,0.12);border-radius:2px;position:relative;cursor:pointer;box-shadow:0 1px 2px rgba(0,0,0,0.2) inset}
#%[1]s-wrapper .audio-volume-bar{position:absolute;left:0;top:0;height:100%%;background:linear-gradient(90deg,rgba(10,132,255,0.8) 0%%,rgba(0,122,255,0.9) 100%%);border-radius:2px;box-shadow:0 0 6px rgba(10,132,255,0.4)}
#%[1]s-wrapper .audio-toolbar{position:absolute;top:8px;right:8px;display:flex;gap:4px;opacity:0;transition:opacity 0.3s ease;z-index:10}
#%[1]s-wrapper:hover .audio-toolbar{opacity:1}
#%[1]s-wrapper .audio-toolbar-btn{background:rgba(0,0,0,0.35);backdrop-filter:blur(10px);border:1px solid rgba(255,255,255,0.08);color:rgba(255,255,255,0.95);padding:5px 9px;border-radius:7px;cursor:pointer;font-size:10px;transition:all 0.2s;text-decoration:none;display:inline-flex;align-items:center;font-family:-apple-system,BlinkMacSystemFont,'SF Pro Text',Roboto,sans-serif;font-weight:500;box-shadow:0 2px 8px rgba(0,0,0,0.2)}
#%[1]s-wrapper .audio-toolbar-btn:hover{background:rgba(0,0,0,0.5);border-color:rgba(255,255,255,0.15);transform:translateY(-1px);box-shadow:0 3px 12px rgba(0,0,0,0.3)}
#%[1]s-wrapper .audio-loading{position:absolute;top:0;left:0;right:0;bottom:0;background:rgba(0,0,0,0.6);backdrop-filter:blur(20px);border-radius:16px;display:flex;align-items:center;justify-content:center;text-align:center;color:rgba(255,255,255,0.95);font-size:11px;z-index:20;font-weight:500}
#%[1]s-wrapper .audio-spinner{border:2px solid rgba(255,255,255,0.2);border-top:2px solid rgba(10,132,255,0.95);border-radius:50%%;width:24px;height:24px;animation:audio-spin-%[1]s 0.7s linear infinite;margin:0 auto 8px;box-shadow:0 0 8px rgba(10,132,255,0.3)}
@keyframes audio-spin-%[1]s{0%%{transform:rotate(0deg)}100%%{transform:rotate(360deg)}}
#%[1]s-wrapper .audio-error{position:absolute;top:0;left:0;right:0;bottom:0;background:rgba(0,0,0,0.7);backdrop-filter:blur(20px);border-radius:16px;display:none;align-items:center;justify-content:center;text-align:center;color:rgba(255,69,58,1);font-size:11px;z-index:20;font-weight:600}
@media (max-width:768px){
#%[1]s-wrapper{margin:10px}
#%[1]s-wrapper .audio-container{padding:12px}
#%[1]s-wrapper .audio-title{font-size:12px}
#%[1]s-wrapper .audio-toolbar{opacity:1}
#%[1]s-wrapper .audio-play-btn{width:40px;height:40px}
}
@media (prefers-color-scheme:light){
#%[1]s-wrapper .audio-container{background:rgba(250,250,252,0.85);border-color:rgba(0,0,0,0.08)}
#%[1]s-wrapper .audio-title{color:rgba(0,0,0,0.95);text-shadow:0 1px 2px rgba(255,255,255,0.5)}
#%[1]s-wrapper .audio-format{color:rgba(0,0,0,0.5)}
#%[1]s-wrapper .audio-control-btn{background:rgba(0,0,0,0.06);border-color:rgba(0,0,0,0.08);color:rgba(0,0,0,0.85)}
#%[1]s-wrapper .audio-volume-icon{color:rgba(0,0,0,0.8)}
#%[1]s-wrapper .audio-time{color:rgba(0,0,0,0.65)}
}
</style>
<div class="audio-container">
<div class="audio-toolbar">
<a href="%[2]s" target="_blank" rel="noopener noreferrer" class="audio-toolbar-btn" title="Open">🔗</a>
</div>
<div class="audio-loading" id="%[1]s-loading">
<div><div class="audio-spinner"></div><div>Loading...</div></div>
</div>
<div class="audio-main">
<button class="audio-play-btn" id="%[1]s-play" title="Play">▶</button>
<div class="audio-info">
<div class="audio-title" title="%[4]s">%[4]s</div>
<div class="audio-format">%[5]s</div>
</div>
</div>
<div class="audio-visualizer" id="%[1]s-visualizer"></div>
<div class="audio-progress" id="%[1]s-progress">
<div class="audio-progress-bar" id="%[1]s-progress-bar"></div>
</div>
<div class="audio-time">
<span id="%[1]s-current-time">0:00</span>
<span id="%[1]s-duration">0:00</span>
</div>
<div class="audio-controls">
<button class="audio-control-btn" id="%[1]s-prev" title="Rewind 10s">⏮</button>
<button class="audio-control-btn" id="%[1]s-next" title="Forward 10s">⏭</button>
<div class="audio-volume">
<span class="audio-volume-icon" id="%[1]s-volume-icon">🔊</span>
<div class="audio-volume-slider" id="%[1]s-volume-slider">
<div class="audio-volume-bar" id="%[1]s-volume-bar"></div>
</div>
</div>
</div>
<div class="audio-error" id="%[1]s-error">❌ Failed to load</div>
</div>
<audio id="%[1]s-player" preload="metadata" style="display:none;">
<source src="%[2]s" type="%[3]s">
</audio>
</div>
<script>
(function(){
const UNIQUE_ID='%[1]s';
const AUDIO_URL='%[2]s';
const elements={wrapper:document.getElementById(UNIQUE_ID+'-wrapper'),loading:document.getElementById(UNIQUE_ID+'-loading'),error:document.getElementById(UNIQUE_ID+'-error'),audioEl:document.getElementById(UNIQUE_ID+'-player'),playBtn:document.getElementById(UNIQUE_ID+'-play'),prevBtn:document.getElementById(UNIQUE_ID+'-prev'),nextBtn:document.getElementById(UNIQUE_ID+'-next'),progress:document.getElementById(UNIQUE_ID+'-progress'),progressBar:document.getElementById(UNIQUE_ID+'-progress-bar'),currentTime:document.getElementById(UNIQUE_ID+'-current-time'),duration:document.getElementById(UNIQUE_ID+'-duration'),volumeIcon:document.getElementById(UNIQUE_ID+'-volume-icon'),volumeSlider:document.getElementById(UNIQUE_ID+'-volume-slider'),volumeBar:document.getElementById(UNIQUE_ID+'-volume-bar'),visualizer:document.getElementById(UNIQUE_ID+'-visualizer')};
if(!elements.wrapper||!elements.audioEl){
console.error('Audio player elements not found for:',UNIQUE_ID);
return;
}
let isPlaying=false;
let isDragging=false;
let visualizerBars=[];
let loadingHidden=false;
let animationId=null;
function formatTime(seconds){
if(isNaN(seconds)||!isFinite(seconds))return'0:00';
const mins=Math.floor(seconds/60);
const secs=Math.floor(seconds%%60);
return mins+':'+(secs<10?'0':'')+secs;
}
function hideLoading(){
if(!loadingHidden&&elements.loading){
elements.loading.style.display='none';
loadingHidden=true;
}
}
function initVisualizer(){
if(!elements.visualizer)return;
for(let i=0;i<50;i++){
const bar=document.createElement('div');
bar.className='audio-bar';
bar.style.height=Math.random()*12+6+'px';
elements.visualizer.appendChild(bar);
visualizerBars.push(bar);
}
}
function animateVisualizer(){
if(isPlaying){
visualizerBars.forEach(bar=>{
const height=Math.random()*24+6;
bar.style.height=height+'px';
});
animationId=requestAnimationFrame(animateVisualizer);
}else{
if(animationId){
cancelAnimationFrame(animationId);
animationId=null;
}
visualizerBars.forEach(bar=>{
bar.style.height='6px';
});
}
}
function togglePlay(){
if(!elements.audioEl)return;
if(isPlaying){
elements.audioEl.pause();
if(elements.playBtn)elements.playBtn.innerHTML='▶';
isPlaying=false;
animateVisualizer();
}else{
const playPromise=elements.audioEl.play();
if(playPromise!==undefined){
playPromise.then(()=>{
if(elements.playBtn)elements.playBtn.innerHTML='⏸';
isPlaying=true;
animateVisualizer();
}).catch(err=>{
console.error('Play error:',err);
if(elements.error){
elements.error.style.display='flex';
elements.error.textContent='❌ Playback failed';
}
});
}
}
}
function updateProgress(){
if(!elements.audioEl||!elements.progressBar||!elements.currentTime)return;
if(!isDragging&&elements.audioEl.duration){
const percent=(elements.audioEl.currentTime/elements.audioEl.duration)*100;
elements.progressBar.style.width=percent+'%%';
elements.currentTime.textContent=formatTime(elements.audioEl.currentTime);
}
}
function seekAudio(e){
if(!elements.audioEl||!elements.progress)return;
const rect=elements.progress.getBoundingClientRect();
const percent=Math.max(0,Math.min(1,(e.clientX-rect.left)/rect.width));
const time=percent*elements.audioEl.duration;
if(!isNaN(time)&&isFinite(time)){
elements.audioEl.currentTime=time;
updateProgress();
}
}
function updateVolume(e){
if(!elements.audioEl||!elements.volumeSlider||!elements.volumeBar)return;
const rect=elements.volumeSlider.getBoundingClientRect();
const percent=Math.max(0,Math.min(1,(e.clientX-rect.left)/rect.width));
elements.audioEl.volume=percent;
elements.volumeBar.style.width=(percent*100)+'%%';
updateVolumeIcon(percent);
}
function updateVolumeIcon(volume){
if(!elements.volumeIcon)return;
if(volume===0)elements.volumeIcon.textContent='🔇';
else if(volume<0.5)elements.volumeIcon.textContent='🔉';
else elements.volumeIcon.textContent='🔊';
}
function initPlayer(){
if(!elements.audioEl)return;
elements.audioEl.addEventListener('loadedmetadata',function(){
hideLoading();
if(elements.duration){
elements.duration.textContent=formatTime(elements.audioEl.duration);
}
console.log('Audio loaded:',elements.audioEl.duration+'s');
});
elements.audioEl.addEventListener('loadeddata',hideLoading);
elements.audioEl.addEventListener('canplay',hideLoading);
elements.audioEl.addEventListener('canplaythrough',hideLoading);
elements.audioEl.addEventListener('timeupdate',updateProgress);
elements.audioEl.addEventListener('ended',function(){
if(elements.playBtn)elements.playBtn.innerHTML='▶';
isPlaying=false;
elements.audioEl.currentTime=0;
updateProgress();
animateVisualizer();
});
elements.audioEl.addEventListener('error',function(e){
hideLoading();
if(elements.error){
elements.error.style.display='flex';
}
console.error('Audio error:',e);
});
if(elements.playBtn){
elements.playBtn.addEventListener('click',togglePlay);
}
if(elements.progress){
elements.progress.addEventListener('click',seekAudio);
elements.progress.addEventListener('mousedown',function(){isDragging=true;});
}
document.addEventListener('mouseup',function(){isDragging=false;});
document.addEventListener('mousemove',function(e){
if(isDragging)seekAudio(e);
});
if(elements.volumeSlider){
elements.volumeSlider.addEventListener('click',updateVolume);
}
if(elements.volumeIcon){
elements.volumeIcon.addEventListener('click',function(){
if(!elements.audioEl)return;
if(elements.audioEl.volume>0){
elements.audioEl.volume=0;
if(elements.volumeBar)elements.volumeBar.style.width='0%%';
updateVolumeIcon(0);
}else{
elements.audioEl.volume=0.7;
if(elements.volumeBar)elements.volumeBar.style.width='70%%';
updateVolumeIcon(0.7);
}
});
}
if(elements.prevBtn){
elements.prevBtn.addEventListener('click',function(){
if(elements.audioEl){
elements.audioEl.currentTime=Math.max(0,elements.audioEl.currentTime-10);
}
});
}
if(elements.nextBtn){
elements.nextBtn.addEventListener('click',function(){
if(elements.audioEl&&elements.audioEl.duration){
elements.audioEl.currentTime=Math.min(elements.audioEl.duration,elements.audioEl.currentTime+10);
}
});
}
document.addEventListener('keydown',function(e){
if(!elements.wrapper)return;
const activeElement=document.activeElement;
if(activeElement&&(activeElement.tagName==='INPUT'||activeElement.tagName==='TEXTAREA'||activeElement.isContentEditable))return;
const rect=elements.wrapper.getBoundingClientRect();
const inViewport=rect.top<window.innerHeight&&rect.bottom>0;
if(!inViewport)return;
switch(e.key){
case ' ':
togglePlay();
e.preventDefault();
break;
case 'ArrowLeft':
if(elements.audioEl){
elements.audioEl.currentTime=Math.max(0,elements.audioEl.currentTime-5);
}
e.preventDefault();
break;
case 'ArrowRight':
if(elements.audioEl&&elements.audioEl.duration){
elements.audioEl.currentTime=Math.min(elements.audioEl.duration,elements.audioEl.currentTime+5);
}
e.preventDefault();
break;
case 'ArrowUp':
if(elements.audioEl){
elements.audioEl.volume=Math.min(1,elements.audioEl.volume+0.1);
if(elements.volumeBar)elements.volumeBar.style.width=(elements.audioEl.volume*100)+'%%';
updateVolumeIcon(elements.audioEl.volume);
}
e.preventDefault();
break;
case 'ArrowDown':
if(elements.audioEl){
elements.audioEl.volume=Math.max(0,elements.audioEl.volume-0.1);
if(elements.volumeBar)elements.volumeBar.style.width=(elements.audioEl.volume*100)+'%%';
updateVolumeIcon(elements.audioEl.volume);
}
e.preventDefault();
break;
case 'm':
if(elements.volumeIcon)elements.volumeIcon.click();
e.preventDefault();
break;
}
});
initVisualizer();
if(elements.audioEl){
elements.audioEl.volume=0.7;
}
if(elements.volumeBar){
elements.volumeBar.style.width='70%%';
}
setTimeout(hideLoading,3000);
}
if(document.readyState==='loading'){
document.addEventListener('DOMContentLoaded',initPlayer);
}else{
initPlayer();
}
window.addEventListener('beforeunload',function(){
if(animationId){
cancelAnimationFrame(animationId);
}
if(elements.audioEl){
elements.audioEl.pause();
elements.audioEl.src='';
elements.audioEl.load();
}
});
})();
</script>`, uniqueId, safeAudioPath, audioType, fileNameWithoutExt, strings.ToUpper(strings.TrimPrefix(ext, ".")))
}

// getAudioMimeType 根据文件扩展名返回对应的 MIME 类型
func getAudioMimeType(ext string) string {
	mimeTypes := map[string]string{
		".mp3":  "audio/mpeg",
		".wav":  "audio/wav",
		".flac": "audio/flac",
		".m4a":  "audio/mp4",
		".aac":  "audio/aac",
		".oga":  "audio/ogg",
		".ogg":  "audio/ogg",
		".opus": "audio/opus",
		".weba": "audio/webm",
	}
	if mimeType, ok := mimeTypes[ext]; ok {
		return mimeType
	}
	return "audio/mpeg"
}
