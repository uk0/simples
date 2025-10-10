package parser

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"log"
	"path/filepath"
	"strings"
)

// GenerateSCSViewerHTML 生成基于 HOOPS Communicator 的 SCS 查看器
// 专门用于渲染 SOLIDWORKS .scs (HOOPS Stream Cache) 格式文件
// 支持多实例，每个实例完全独立
func GenerateSCSViewerHTML(cadId, cadPath string) string {
	safeCadId := template.HTMLEscapeString(cadId)
	safeCadPath := template.HTMLEscapeString(cadPath)
	log.Println("Generating SCS HOOPS viewer for file:", safeCadPath)

	h := md5.Sum([]byte(safeCadId))
	uniqueId := fmt.Sprintf("scs_%x", h)[:12]
	fileName := filepath.Base(cadPath)
	ext := strings.ToLower(filepath.Ext(fileName))

	// 确定文件类型
	fileType := "SCS"
	displayName := fileName
	if ext == ".scs" {
		fileType = "HOOPS SCS"
	} else if ext == ".scz" {
		fileType = "HOOPS SCZ"
	} else if ext == ".hsf" {
		fileType = "HOOPS HSF"
	}

	return fmt.Sprintf(`<div id="%[1]s-wrapper" class="hoops-viewer-wrapper" style="position:relative;max-width:100%%;margin:20px auto;display:block;">
<style>
#%[1]s-wrapper{position:relative;width:100%%;max-width:1400px;margin:20px auto;display:block;clear:both;font-family:-apple-system,BlinkMacSystemFont,'SF Pro Text','Segoe UI',Roboto,sans-serif}
#%[1]s-wrapper *{box-sizing:border-box}
#%[1]s-wrapper .hoops-container{position:relative;width:100%%;background:rgba(30,30,32,0.75);backdrop-filter:blur(40px) saturate(180%%);-webkit-backdrop-filter:blur(40px) saturate(180%%);border:1px solid rgba(255,255,255,0.12);border-radius:16px;box-shadow:0 8px 32px rgba(0,0,0,0.24),0 0 1px rgba(255,255,255,0.12) inset;overflow:hidden;transition:all 0.3s ease}
#%[1]s-wrapper .hoops-container:hover{background:rgba(40,40,42,0.8);box-shadow:0 12px 48px rgba(0,0,0,0.32),0 0 1px rgba(255,255,255,0.18) inset}
#%[1]s-wrapper .hoops-header{background:rgba(0,0,0,0.3);padding:16px 20px;border-bottom:1px solid rgba(255,255,255,0.08);display:flex;align-items:center;justify-content:space-between;gap:12px;flex-wrap:wrap}
#%[1]s-wrapper .hoops-title{font-size:14px;font-weight:600;color:rgba(255,255,255,0.98);display:flex;align-items:center;gap:10px;text-shadow:0 1px 3px rgba(0,0,0,0.3);flex:1;min-width:200px}
#%[1]s-wrapper .hoops-icon{font-size:20px;flex-shrink:0}
#%[1]s-wrapper .hoops-title-text{flex:1;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
#%[1]s-wrapper .hoops-badge{background:rgba(10,132,255,0.2);border:1px solid rgba(10,132,255,0.3);color:rgba(10,132,255,1);padding:4px 10px;border-radius:6px;font-size:10px;font-weight:600;text-transform:uppercase;letter-spacing:0.5px}
#%[1]s-wrapper .hoops-toolbar{display:flex;align-items:center;gap:8px;flex-wrap:wrap}
#%[1]s-wrapper .hoops-btn{background:rgba(255,255,255,0.12);border:1px solid rgba(255,255,255,0.08);color:rgba(255,255,255,0.95);padding:8px 14px;border-radius:8px;cursor:pointer;font-size:12px;font-weight:500;transition:all 0.2s;display:inline-flex;align-items:center;gap:6px;text-decoration:none;box-shadow:0 2px 6px rgba(0,0,0,0.15);white-space:nowrap}
#%[1]s-wrapper .hoops-btn:hover{background:rgba(255,255,255,0.2);border-color:rgba(255,255,255,0.15);transform:translateY(-1px);box-shadow:0 3px 10px rgba(0,0,0,0.2)}
#%[1]s-wrapper .hoops-btn:active{transform:translateY(0)}
#%[1]s-wrapper .hoops-btn:disabled{opacity:0.5;cursor:not-allowed;transform:none}
#%[1]s-wrapper .hoops-btn-primary{background:linear-gradient(135deg,#007AFF 0%%,#0051D5 100%%);border-color:transparent;box-shadow:0 3px 10px rgba(0,122,255,0.4)}
#%[1]s-wrapper .hoops-btn-primary:hover{background:linear-gradient(135deg,#0A84FF 0%%,#0060EA 100%%);box-shadow:0 4px 14px rgba(0,122,255,0.5)}
#%[1]s-wrapper .hoops-btn-group{display:flex;gap:4px;background:rgba(0,0,0,0.2);padding:4px;border-radius:8px}
#%[1]s-wrapper .hoops-btn-small{padding:6px 10px;font-size:11px;min-width:auto}
#%[1]s-wrapper .hoops-btn-active{background:rgba(10,132,255,0.3);border-color:rgba(10,132,255,0.5);color:rgba(10,132,255,1)}
#%[1]s-wrapper .hoops-viewer-content{position:relative;height:700px;background:#16213e;overflow:hidden}
#%[1]s-wrapper .hoops-canvas-container{width:100%%;height:100%%;position:relative}
#%[1]s-wrapper .hoops-sidebar{position:absolute;top:0;left:0;width:280px;height:100%%;background:rgba(0,0,0,0.9);backdrop-filter:blur(20px);border-right:1px solid rgba(255,255,255,0.1);overflow:hidden;z-index:20;transform:translateX(-100%%);transition:transform 0.3s ease;display:flex;flex-direction:column}
#%[1]s-wrapper .hoops-sidebar.active{transform:translateX(0)}
#%[1]s-wrapper .hoops-sidebar-header{padding:16px;border-bottom:1px solid rgba(255,255,255,0.1);display:flex;justify-content:space-between;align-items:center;flex-shrink:0;background:rgba(0,0,0,0.5)}
#%[1]s-wrapper .hoops-sidebar-title{font-size:13px;font-weight:600;color:rgba(255,255,255,0.95);display:flex;align-items:center;gap:8px}
#%[1]s-wrapper .hoops-sidebar-close{background:rgba(255,255,255,0.1);border:none;color:rgba(255,255,255,0.9);width:24px;height:24px;border-radius:6px;cursor:pointer;display:flex;align-items:center;justify-content:center;transition:all 0.2s;font-size:16px;line-height:1}
#%[1]s-wrapper .hoops-sidebar-close:hover{background:rgba(255,69,58,0.3);color:rgba(255,69,58,1)}
#%[1]s-wrapper .hoops-sidebar-content{padding:12px;flex:1;overflow-y:auto;overflow-x:hidden}
#%[1]s-wrapper .hoops-sidebar-content::-webkit-scrollbar{width:6px}
#%[1]s-wrapper .hoops-sidebar-content::-webkit-scrollbar-track{background:rgba(255,255,255,0.05);border-radius:3px}
#%[1]s-wrapper .hoops-sidebar-content::-webkit-scrollbar-thumb{background:rgba(255,255,255,0.2);border-radius:3px}
#%[1]s-wrapper .hoops-sidebar-content::-webkit-scrollbar-thumb:hover{background:rgba(255,255,255,0.3)}
#%[1]s-wrapper .hoops-tree{list-style:none;padding:0;margin:0}
#%[1]s-wrapper .hoops-tree-item{padding:8px 10px;margin:2px 0;border-radius:6px;cursor:pointer;font-size:12px;color:rgba(255,255,255,0.85);display:flex;align-items:center;gap:8px;transition:all 0.2s;user-select:none;position:relative}
#%[1]s-wrapper .hoops-tree-item:hover{background:rgba(255,255,255,0.12)}
#%[1]s-wrapper .hoops-tree-item.selected{background:rgba(10,132,255,0.25);color:rgba(10,132,255,1);font-weight:600}
#%[1]s-wrapper .hoops-tree-item.hidden{opacity:0.4}
#%[1]s-wrapper .hoops-tree-toggle{width:16px;height:16px;display:inline-flex;align-items:center;justify-content:center;font-size:10px;flex-shrink:0;transition:transform 0.2s;cursor:pointer}
#%[1]s-wrapper .hoops-tree-toggle.expanded{transform:rotate(90deg)}
#%[1]s-wrapper .hoops-tree-icon{font-size:14px;flex-shrink:0}
#%[1]s-wrapper .hoops-tree-label{flex:1;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
#%[1]s-wrapper .hoops-tree-children{margin-left:20px;list-style:none;padding:0;display:none}
#%[1]s-wrapper .hoops-tree-children.expanded{display:block}
#%[1]s-wrapper .hoops-loading{display:flex;align-items:center;justify-content:center;padding:100px 20px;color:rgba(255,255,255,0.95);font-size:13px;font-weight:500;flex-direction:column;position:absolute;top:50%%;left:50%%;transform:translate(-50%%,-50%%);z-index:100;width:80%%;max-width:400px;background:rgba(0,0,0,0.85);border-radius:12px;pointer-events:none}
#%[1]s-wrapper .hoops-spinner{border:3px solid rgba(255,255,255,0.2);border-top:3px solid rgba(10,132,255,0.95);border-radius:50%%;width:40px;height:40px;animation:hoops-spin-%[1]s 0.7s linear infinite;margin-bottom:16px;box-shadow:0 0 12px rgba(10,132,255,0.3)}
@keyframes hoops-spin-%[1]s{0%%{transform:rotate(0deg)}100%%{transform:rotate(360deg)}}
#%[1]s-wrapper .hoops-loading-text{margin-top:8px;font-size:14px}
#%[1]s-wrapper .hoops-loading-detail{font-size:11px;color:rgba(255,255,255,0.6);margin-top:6px}
#%[1]s-wrapper .hoops-error{display:none;align-items:center;justify-content:center;padding:60px 20px;color:rgba(255,69,58,1);font-size:13px;font-weight:600;text-align:center;flex-direction:column;gap:12px;position:absolute;top:50%%;left:50%%;transform:translate(-50%%,-50%%);width:90%%;max-width:500px;z-index:100;background:rgba(0,0,0,0.9);border-radius:12px;border:2px solid rgba(255,69,58,0.5)}
#%[1]s-wrapper .hoops-error-icon{font-size:48px;margin-bottom:12px}
#%[1]s-wrapper .hoops-error-detail{font-size:11px;color:rgba(255,255,255,0.7);font-weight:400;margin-top:8px;line-height:1.5;max-width:400px}
#%[1]s-wrapper .hoops-info-panel{position:absolute;bottom:20px;right:20px;background:rgba(0,0,0,0.85);backdrop-filter:blur(10px);padding:12px 16px;border-radius:10px;font-size:11px;color:rgba(255,255,255,0.9);border:1px solid rgba(255,255,255,0.1);box-shadow:0 4px 16px rgba(0,0,0,0.3);z-index:5;display:none;transition:opacity 0.3s ease}
#%[1]s-wrapper .hoops-info-item{display:flex;justify-content:space-between;margin-bottom:6px;gap:20px}
#%[1]s-wrapper .hoops-info-item:last-child{margin-bottom:0}
#%[1]s-wrapper .hoops-info-label{color:rgba(255,255,255,0.6)}
#%[1]s-wrapper .hoops-info-value{color:rgba(10,132,255,1);font-weight:600;font-family:'SF Mono','Monaco','Consolas',monospace}
#%[1]s-wrapper .hoops-footer{background:rgba(0,0,0,0.2);padding:12px 20px;border-top:1px solid rgba(255,255,255,0.06);display:flex;align-items:center;justify-content:space-between;font-size:11px;color:rgba(255,255,255,0.7);flex-wrap:wrap;gap:10px}
#%[1]s-wrapper .hoops-powered{display:flex;align-items:center;gap:6px;font-weight:500}
#%[1]s-wrapper .hoops-powered-logo{background:rgba(124,77,255,0.2);border:1px solid rgba(124,77,255,0.3);color:rgba(124,77,255,1);padding:2px 8px;border-radius:4px;font-size:9px;font-weight:700;letter-spacing:0.5px}
#%[1]s-wrapper .hoops-controls-hint{color:rgba(255,255,255,0.5);font-size:10px}
@media (max-width:768px){
#%[1]s-wrapper{margin:10px}
#%[1]s-wrapper .hoops-header{padding:12px 16px}
#%[1]s-wrapper .hoops-title{font-size:13px}
#%[1]s-wrapper .hoops-toolbar{width:100%%;justify-content:flex-start}
#%[1]s-wrapper .hoops-btn{padding:6px 10px;font-size:11px}
#%[1]s-wrapper .hoops-viewer-content{height:500px}
#%[1]s-wrapper .hoops-sidebar{width:240px}
#%[1]s-wrapper .hoops-info-panel{bottom:10px;right:10px;font-size:10px;padding:8px 12px}
#%[1]s-wrapper .hoops-footer{font-size:10px;padding:10px 16px}
}
</style>
<div class="hoops-container">
<div class="hoops-header">
<div class="hoops-title">
<span class="hoops-icon">🏭</span>
<span class="hoops-title-text" title="%[3]s">%[3]s</span>
<span class="hoops-badge">%[4]s</span>
</div>
<div class="hoops-toolbar">
<button id="%[1]s-btn-tree" class="hoops-btn hoops-btn-small" title="模型树" disabled>🌳 Tree</button>
<div class="hoops-btn-group">
<button id="%[1]s-btn-fit" class="hoops-btn hoops-btn-small" title="适配视图" disabled>🔄 Fit</button>
<button id="%[1]s-btn-wireframe" class="hoops-btn hoops-btn-small" title="线框模式" disabled>📐 Wire</button>
<button id="%[1]s-btn-shaded" class="hoops-btn hoops-btn-small hoops-btn-active" title="着色模式" disabled>🎨 Shade</button>
</div>
<button id="%[1]s-btn-measure" class="hoops-btn hoops-btn-small" title="测量工具" disabled>📏 Measure</button>
<button id="%[1]s-btn-section" class="hoops-btn hoops-btn-small" title="剖切" disabled>✂️ Section</button>
<button id="%[1]s-btn-fullscreen" class="hoops-btn hoops-btn-small" title="全屏" disabled>🔍 Full</button>
<a href="%[2]s" download class="hoops-btn hoops-btn-primary hoops-btn-small" title="下载">⬇️ Down</a>
</div>
</div>
<div class="hoops-viewer-content">
<div id="%[1]s-sidebar" class="hoops-sidebar">
<div class="hoops-sidebar-header">
<div class="hoops-sidebar-title">
<span>📦</span>
<span>Model Tree</span>
</div>
<button class="hoops-sidebar-close" id="%[1]s-sidebar-close">✕</button>
</div>
<div class="hoops-sidebar-content">
<ul id="%[1]s-tree" class="hoops-tree"></ul>
</div>
</div>
<div id="%[1]s-viewer" class="hoops-canvas-container"></div>
<div id="%[1]s-loading" class="hoops-loading">
<div class="hoops-spinner"></div>
<div class="hoops-loading-text">Initializing HOOPS Viewer...</div>
<div class="hoops-loading-detail">Loading model from server</div>
</div>
<div id="%[1]s-error" class="hoops-error">
<div class="hoops-error-icon">❌</div>
<div>Failed to load model</div>
<div class="hoops-error-detail" id="%[1]s-error-detail">
Please check if the file format is correct (.scs, .scz, .hsf)
</div>
</div>
<div id="%[1]s-info" class="hoops-info-panel">
<div class="hoops-info-item">
<span class="hoops-info-label">Format:</span>
<span class="hoops-info-value">%[4]s</span>
</div>
<div class="hoops-info-item">
<span class="hoops-info-label">Nodes:</span>
<span class="hoops-info-value" id="%[1]s-nodes">-</span>
</div>
<div class="hoops-info-item">
<span class="hoops-info-label">Render:</span>
<span class="hoops-info-value">WebGL</span>
</div>
</div>
</div>
<div class="hoops-footer">
<div class="hoops-powered">
<span>Powered by</span>
<span class="hoops-powered-logo">HOOPS</span>
</div>
<div class="hoops-controls-hint">
🖱️ 左键旋转 · 右键平移 · 滚轮缩放
</div>
</div>
</div>
</div>
<script type="text/javascript" src="https://cdn.jsdelivr.net/gh/techsoft3d/hoops-web-viewer@latest/hoops_web_viewer.js"></script>
<script>
(function(){
'use strict';
const UNIQUE_ID='%[1]s';
const MODEL_URL='%[2]s';
const FILE_TYPE='%[4]s';
const ViewerState={
hwv:null,
wireframeMode:false,
shadedMode:true,
measureActive:false,
sectionActive:false,
treeVisible:false,
cuttingPlaneIds:[],
selectedNodeId:null,
isLoading:true,
isLoaded:false,
loadTimeout:null
};
const elements={
wrapper:document.getElementById(UNIQUE_ID+'-wrapper'),
viewerContainer:document.getElementById(UNIQUE_ID+'-viewer'),
loading:document.getElementById(UNIQUE_ID+'-loading'),
loadingText:null,
loadingDetail:null,
error:document.getElementById(UNIQUE_ID+'-error'),
errorDetail:document.getElementById(UNIQUE_ID+'-error-detail'),
info:document.getElementById(UNIQUE_ID+'-info'),
nodesEl:document.getElementById(UNIQUE_ID+'-nodes'),
btnTree:document.getElementById(UNIQUE_ID+'-btn-tree'),
btnFit:document.getElementById(UNIQUE_ID+'-btn-fit'),
btnWireframe:document.getElementById(UNIQUE_ID+'-btn-wireframe'),
btnShaded:document.getElementById(UNIQUE_ID+'-btn-shaded'),
btnMeasure:document.getElementById(UNIQUE_ID+'-btn-measure'),
btnSection:document.getElementById(UNIQUE_ID+'-btn-section'),
btnFullscreen:document.getElementById(UNIQUE_ID+'-btn-fullscreen'),
sidebar:document.getElementById(UNIQUE_ID+'-sidebar'),
sidebarClose:document.getElementById(UNIQUE_ID+'-sidebar-close'),
tree:document.getElementById(UNIQUE_ID+'-tree')
};
if(!elements.wrapper){
console.error('['+UNIQUE_ID+'] Wrapper element not found');
return;
}
if(elements.loading){
elements.loadingText=elements.loading.querySelector('.hoops-loading-text');
elements.loadingDetail=elements.loading.querySelector('.hoops-loading-detail');
}
function enableAllButtons(){
const buttons=[
elements.btnTree,
elements.btnFit,
elements.btnWireframe,
elements.btnShaded,
elements.btnMeasure,
elements.btnSection,
elements.btnFullscreen
];
buttons.forEach(function(btn){
if(btn)btn.disabled=false;
});
console.log('['+UNIQUE_ID+'] All buttons enabled');
}
function hideLoading(){
ViewerState.isLoading=false;
if(ViewerState.loadTimeout){
clearTimeout(ViewerState.loadTimeout);
ViewerState.loadTimeout=null;
}
if(elements.loading){
elements.loading.style.display='none';
console.log('['+UNIQUE_ID+'] Loading overlay hidden');
}
if(elements.info){
elements.info.style.display='block';
console.log('['+UNIQUE_ID+'] Info panel shown');
}
enableAllButtons();
}
function showError(message,detail){
ViewerState.isLoading=false;
if(ViewerState.loadTimeout){
clearTimeout(ViewerState.loadTimeout);
ViewerState.loadTimeout=null;
}
if(elements.loading){
elements.loading.style.display='none';
}
if(elements.error){
elements.error.style.display='flex';
if(detail&&elements.errorDetail){
elements.errorDetail.textContent=detail;
}
}
console.error('['+UNIQUE_ID+'] Error:',message,detail);
}
function updateLoadingText(text,detail){
if(!ViewerState.isLoading){
return;
}
if(text&&elements.loadingText){
elements.loadingText.textContent=text;
}
if(detail&&elements.loadingDetail){
elements.loadingDetail.textContent=detail;
}
console.log('['+UNIQUE_ID+'] Loading:',text,detail||'');
}
async function initHOOPS(){
try{
if(typeof Communicator==='undefined'){
showError('HOOPS Library not loaded','Failed to load hoops_web_viewer.js from CDN');
return;
}
updateLoadingText('Initializing HOOPS Communicator...','Setting up WebGL renderer');
ViewerState.loadTimeout=setTimeout(function(){
console.warn('['+UNIQUE_ID+'] Load timeout reached, forcing UI update');
if(ViewerState.isLoading&&!ViewerState.isLoaded){
showError('Loading timeout','Model took too long to load. Please refresh and try again.');
}else if(ViewerState.isLoading&&ViewerState.isLoaded){
hideLoading();
}
},30000);
ViewerState.hwv=new Communicator.WebViewer({
containerId:UNIQUE_ID+'-viewer',
model:'',
rendererType:Communicator.RendererType.Client
});
let modelRendered=false;
ViewerState.hwv.setCallbacks({
sceneReady:function(){
console.log('['+UNIQUE_ID+'] HOOPS scene ready');
updateLoadingText('Scene ready, loading model...','Fetching SCS file');
setupEnvironment();
loadModel();
},
modelLoadBegin:function(){
console.log('['+UNIQUE_ID+'] Model load begin');
updateLoadingText('Loading model...','Parsing geometry data');
},
modelLoadEnd:function(){
console.log('['+UNIQUE_ID+'] Model load complete');
ViewerState.isLoaded=true;
hideLoading();
setTimeout(function(){
fitView();
updateModelInfo();
buildModelTree();
},100);
},
modelLoadFailure:function(error){
console.error('['+UNIQUE_ID+'] Model load failed:',error);
showError('Failed to load model',error?error.toString():'Unknown error');
},
modelStructureReady:function(){
console.log('['+UNIQUE_ID+'] Model structure ready');
updateModelInfo();
},
frameDrawn:function(){
if(!modelRendered&&ViewerState.isLoaded){
modelRendered=true;
console.log('['+UNIQUE_ID+'] First frame drawn');
hideLoading();
}
},
selectionArray:function(selectionEvents){
if(selectionEvents.length>0){
const selection=selectionEvents[0];
const nodeId=selection.getSelection().getNodeId();
ViewerState.selectedNodeId=nodeId;
updateTreeSelection(nodeId);
console.log('['+UNIQUE_ID+'] Selected node:',nodeId);
}
}
});
ViewerState.hwv.start();
console.log('['+UNIQUE_ID+'] HOOPS WebViewer started');
}catch(error){
console.error('['+UNIQUE_ID+'] HOOPS init error:',error);
showError('Initialization failed',error.message||'Unknown initialization error');
}
}
function setupEnvironment(){
if(!ViewerState.hwv||!ViewerState.hwv.view)return;
try{
const view=ViewerState.hwv.view;
view.setBackgroundColor(
new Communicator.Color(22,33,62),
new Communicator.Color(26,26,46)
);
view.setAntiAliasingMode(Communicator.AntiAliasingMode.SMAA);
view.setAmbientOcclusionEnabled(true);
view.setAmbientOcclusionRadius(5);
console.log('['+UNIQUE_ID+'] Environment configured');
}catch(error){
console.error('['+UNIQUE_ID+'] Setup environment error:',error);
}
}
async function loadModel(){
if(!ViewerState.hwv||!ViewerState.hwv.model)return;
try{
updateLoadingText('Loading SCS file...','Downloading from server');
await ViewerState.hwv.model.clear();
const result=await ViewerState.hwv.model.loadSubtreeFromScsFile(
ViewerState.hwv.model.getRootNode(),
MODEL_URL
);
console.log('['+UNIQUE_ID+'] Model loaded successfully');
if(result!==null){
ViewerState.isLoaded=true;
setTimeout(function(){
hideLoading();
fitView();
updateModelInfo();
buildModelTree();
},300);
}
}catch(error){
console.error('['+UNIQUE_ID+'] Load model error:',error);
showError('Model loading failed',error.message||'Invalid SCS file or network error');
}
}
function updateModelInfo(){
if(!ViewerState.hwv||!ViewerState.hwv.model)return;
try{
const rootNode=ViewerState.hwv.model.getRootNode();
const children=ViewerState.hwv.model.getNodeChildren(rootNode);
if(elements.nodesEl){
elements.nodesEl.textContent=children.length.toString();
}
console.log('['+UNIQUE_ID+'] Model info: nodes='+children.length);
}catch(error){
console.error('['+UNIQUE_ID+'] Update model info error:',error);
}
}
function buildModelTree(){
if(!ViewerState.hwv||!ViewerState.hwv.model||!elements.tree)return;
try{
const model=ViewerState.hwv.model;
const rootNode=model.getRootNode();
elements.tree.innerHTML='';
function createTreeNode(nodeId,parentUl,depth){
const nodeName=model.getNodeName(nodeId)||('Part_'+nodeId);
const children=model.getNodeChildren(nodeId);
const hasChildren=children.length>0;
const li=document.createElement('li');
const itemDiv=document.createElement('div');
itemDiv.className='hoops-tree-item';
itemDiv.dataset.nodeId=nodeId;
itemDiv.style.paddingLeft=(depth*8)+'px';
if(hasChildren){
const toggle=document.createElement('span');
toggle.className='hoops-tree-toggle';
toggle.textContent='▶';
toggle.onclick=function(e){
e.stopPropagation();
const childrenUl=li.querySelector('.hoops-tree-children');
if(childrenUl){
const isExpanded=toggle.classList.contains('expanded');
if(isExpanded){
toggle.classList.remove('expanded');
childrenUl.classList.remove('expanded');
}else{
toggle.classList.add('expanded');
childrenUl.classList.add('expanded');
}
}
};
itemDiv.appendChild(toggle);
}else{
const spacer=document.createElement('span');
spacer.className='hoops-tree-toggle';
spacer.textContent='';
spacer.style.width='16px';
itemDiv.appendChild(spacer);
}
const icon=document.createElement('span');
icon.className='hoops-tree-icon';
icon.textContent=hasChildren?'📦':'🔩';
itemDiv.appendChild(icon);
const label=document.createElement('span');
label.className='hoops-tree-label';
label.textContent=nodeName;
label.title=nodeName;
itemDiv.appendChild(label);
itemDiv.onclick=function(e){
e.stopPropagation();
ViewerState.selectedNodeId=nodeId;
if(ViewerState.hwv&&ViewerState.hwv.view){
ViewerState.hwv.view.isolateNodes([nodeId]);
setTimeout(function(){
ViewerState.hwv.view.fitNodes([nodeId]);
},100);
}
updateTreeSelection(nodeId);
};
li.appendChild(itemDiv);
parentUl.appendChild(li);
if(hasChildren){
const childrenUl=document.createElement('ul');
childrenUl.className='hoops-tree-children';
li.appendChild(childrenUl);
children.forEach(function(childId){
createTreeNode(childId,childrenUl,depth+1);
});
}
}
const children=model.getNodeChildren(rootNode);
children.forEach(function(nodeId){
createTreeNode(nodeId,elements.tree,0);
});
console.log('['+UNIQUE_ID+'] Model tree built with',children.length,'root nodes');
}catch(error){
console.error('['+UNIQUE_ID+'] Build tree error:',error);
}
}
function updateTreeSelection(nodeId){
if(!elements.tree)return;
const allItems=elements.tree.querySelectorAll('.hoops-tree-item');
allItems.forEach(function(item){
item.classList.remove('selected');
});
const selectedItem=elements.tree.querySelector('[data-node-id="'+nodeId+'"]');
if(selectedItem){
selectedItem.classList.add('selected');
selectedItem.scrollIntoView({behavior:'smooth',block:'nearest'});
}
}
function toggleTree(){
ViewerState.treeVisible=!ViewerState.treeVisible;
if(elements.sidebar){
if(ViewerState.treeVisible){
elements.sidebar.classList.add('active');
if(elements.btnTree){
elements.btnTree.classList.add('hoops-btn-active');
}
}else{
elements.sidebar.classList.remove('active');
if(elements.btnTree){
elements.btnTree.classList.remove('hoops-btn-active');
}
}
}
console.log('['+UNIQUE_ID+'] Tree visible:',ViewerState.treeVisible);
}
function fitView(){
if(!ViewerState.hwv||!ViewerState.hwv.view||!ViewerState.isLoaded)return;
try{
ViewerState.hwv.view.fitWorld();
console.log('['+UNIQUE_ID+'] View fitted');
}catch(error){
console.error('['+UNIQUE_ID+'] Fit view error:',error);
}
}
function toggleWireframe(){
if(!ViewerState.hwv||!ViewerState.hwv.view||!ViewerState.isLoaded)return;
try{
ViewerState.wireframeMode=!ViewerState.wireframeMode;
const view=ViewerState.hwv.view;
if(ViewerState.wireframeMode){
view.setDrawMode(Communicator.DrawMode.Wireframe);
if(elements.btnWireframe){
elements.btnWireframe.classList.add('hoops-btn-active');
}
if(elements.btnShaded){
elements.btnShaded.classList.remove('hoops-btn-active');
}
}else{
view.setDrawMode(Communicator.DrawMode.Shaded);
if(elements.btnWireframe){
elements.btnWireframe.classList.remove('hoops-btn-active');
}
if(elements.btnShaded){
elements.btnShaded.classList.add('hoops-btn-active');
}
}
console.log('['+UNIQUE_ID+'] Wireframe:',ViewerState.wireframeMode);
}catch(error){
console.error('['+UNIQUE_ID+'] Toggle wireframe error:',error);
}
}
function toggleShaded(){
if(!ViewerState.hwv||!ViewerState.hwv.view||!ViewerState.isLoaded)return;
try{
ViewerState.shadedMode=!ViewerState.shadedMode;
const view=ViewerState.hwv.view;
if(ViewerState.shadedMode){
view.setDrawMode(Communicator.DrawMode.Shaded);
if(elements.btnShaded){
elements.btnShaded.classList.add('hoops-btn-active');
}
}else{
view.setDrawMode(Communicator.DrawMode.HiddenLine);
if(elements.btnShaded){
elements.btnShaded.classList.remove('hoops-btn-active');
}
}
console.log('['+UNIQUE_ID+'] Shaded:',ViewerState.shadedMode);
}catch(error){
console.error('['+UNIQUE_ID+'] Toggle shaded error:',error);
}
}
function toggleMeasure(){
if(!ViewerState.hwv||!ViewerState.hwv.operatorManager||!ViewerState.isLoaded)return;
try{
ViewerState.measureActive=!ViewerState.measureActive;
if(ViewerState.measureActive){
ViewerState.hwv.operatorManager.set(Communicator.OperatorId.MeasurePointPointDistance,0);
if(elements.btnMeasure){
elements.btnMeasure.classList.add('hoops-btn-active');
}
}else{
ViewerState.hwv.operatorManager.clear();
if(elements.btnMeasure){
elements.btnMeasure.classList.remove('hoops-btn-active');
}
}
console.log('['+UNIQUE_ID+'] Measure:',ViewerState.measureActive);
}catch(error){
console.error('['+UNIQUE_ID+'] Toggle measure error:',error);
}
}
function toggleSection(){
if(!ViewerState.hwv||!ViewerState.hwv.cuttingManager||!ViewerState.isLoaded)return;
try{
ViewerState.sectionActive=!ViewerState.sectionActive;
if(ViewerState.sectionActive){
const plane=new Communicator.Plane();
const camera=ViewerState.hwv.view.getCamera();
const target=camera.getTarget();
plane.normal=new Communicator.Point3(0,0,1);
plane.d=-target.z;
const planeId=ViewerState.hwv.cuttingManager.addCuttingPlane(plane);
ViewerState.cuttingPlaneIds.push(planeId);
if(elements.btnSection){
elements.btnSection.classList.add('hoops-btn-active');
}
}else{
ViewerState.cuttingPlaneIds.forEach(function(id){
ViewerState.hwv.cuttingManager.removeCuttingPlane(id);
});
ViewerState.cuttingPlaneIds=[];
if(elements.btnSection){
elements.btnSection.classList.remove('hoops-btn-active');
}
}
console.log('['+UNIQUE_ID+'] Section:',ViewerState.sectionActive);
}catch(error){
console.error('['+UNIQUE_ID+'] Toggle section error:',error);
}
}
function toggleFullscreen(){
const container=elements.viewerContainer?elements.viewerContainer.parentElement:null;
if(!container)return;
try{
if(document.fullscreenElement){
document.exitFullscreen();
}else{
if(container.requestFullscreen){
container.requestFullscreen();
}else if(container.webkitRequestFullscreen){
container.webkitRequestFullscreen();
}else if(container.mozRequestFullScreen){
container.mozRequestFullScreen();
}else if(container.msRequestFullscreen){
container.msRequestFullscreen();
}
}
console.log('['+UNIQUE_ID+'] Fullscreen toggled');
}catch(error){
console.error('['+UNIQUE_ID+'] Fullscreen error:',error);
}
}
if(elements.btnTree){
elements.btnTree.addEventListener('click',toggleTree);
}
if(elements.btnFit){
elements.btnFit.addEventListener('click',fitView);
}
if(elements.btnWireframe){
elements.btnWireframe.addEventListener('click',toggleWireframe);
}
if(elements.btnShaded){
elements.btnShaded.addEventListener('click',toggleShaded);
}
if(elements.btnMeasure){
elements.btnMeasure.addEventListener('click',toggleMeasure);
}
if(elements.btnSection){
elements.btnSection.addEventListener('click',toggleSection);
}
if(elements.btnFullscreen){
elements.btnFullscreen.addEventListener('click',toggleFullscreen);
}
if(elements.sidebarClose){
elements.sidebarClose.addEventListener('click',toggleTree);
}
if(typeof Communicator!=='undefined'){
console.log('['+UNIQUE_ID+'] Communicator already loaded, initializing...');
initHOOPS();
}else{
console.log('['+UNIQUE_ID+'] Waiting for Communicator to load...');
let retries=0;
const maxRetries=30;
const checkInterval=setInterval(function(){
if(typeof Communicator!=='undefined'){
clearInterval(checkInterval);
console.log('['+UNIQUE_ID+'] Communicator loaded after '+retries+' retries');
initHOOPS();
}else if(++retries>=maxRetries){
clearInterval(checkInterval);
showError('HOOPS library failed to load','Please check your internet connection and refresh the page');
}
},500);
}
console.log('['+UNIQUE_ID+'] Viewer instance initialized');
})();
</script>`, uniqueId, safeCadPath, displayName, fileType)
}
