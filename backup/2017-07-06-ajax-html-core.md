---
layout: post
title:  Ajax封装类,所有的Ajax都需要通过此方法进行发送 以便进行相应的Ajax方式管理
categories: ajax,httprequest
description: 回顾
keywords: javascript, ajax
---


发这个贴的原因，是因为 是因为 是因为啥来着，忘了 那就是因为你，因为你 。


## 代码

```javascript
/**
 * @author 第一作者：张建新
 * Ajax封装类,所有的Ajax都需要通过此方法进行发送 以便进行相应的Ajax方式管理
 * 注意 $.support.cors = true; 用来解决Ajax跨域问题，此方法需要和后台数据接口相配合
 * successFun,errorFun,beforeSendFun,completeFun
 * 示例写法
 * var successFun = function functionName(data)
 * {
 * 		对Ajax解析过来的数据进行渲染
 * }
 * 将successFun作为参数传入方法中
 * @param {Object} actionUrl 数据请求地址
 * @param {Object} parems Ajax参数列表
 * @param {Object} sendMethod 发送方式
 * @param {Object} successFun 请求成功后回调函数
 * @param {Object} errorFun 请求错误时回调函数
 * @param {Object} beforeSendFun 在发送请求之前回调函数
 * @param {Object} completeFun 请求完成后回调函数 (请求成功或失败之后均调用)
 */
function ajaxLoad(actionUrl,parems,sendMethod,successFun,errorFun,beforeSendFun,completeFun){
	/**解决ajax跨域问题*/
	$.support.cors = true;
	$.ajax( {
		url :actionUrl, /**请求的url地址*/
		dataType : "json", /**返回格式为json*/
		async : true,/**请求是否异步，默认为异步，这也是ajax重要特性 */
		data : parems, /**参数值*/
		type :sendMethod, /**请求方式*/
		/** 状态*/
		beforeSend :beforeSendFun,
		success : successFun,
		complete : completeFun,
		error :errorFun
	});
}

```


## Toast显示框

```javascript

/*
 * 是否正在Toast通知
 */
 
var isToasting = false;
/*
 * Toast通知
 */
function Toast(toastText){
	/*如果正在显示通知*/
	if(isToasting){
		/*延时3秒回调方法*/
		setTimeout(function(){
			Toast(toastText)
		},3000);
	}else{
		/*设置正在显示通知*/
		isToasting = true;
		/*创建通知控件*/
		var toast = document.createElement("div");
		/*设置通知文本*/
		toast.innerHTML = toastText;
		/*设置通知样式*/
		toast.className = "Toast";
		/*延时3秒销毁控件*/
		setTimeout(function(){	
			document.body.removeChild(toast);
			isToasting = false;
		},3000);
		/*显示控件*/
		document.body.appendChild(toast);
	}
}

```


## Echart echarts-utils

```javascript

/*
 * 参数列表
 * title 雷达图标题 字符串
 * subTitle 雷达图子标题 字符串
 * analyzeItems 分析项目 字符串数组 如分析北京 上海 广州 的农业林业牧业副业渔业 ['北京','上海','广州']
 * kindItems 雷达图对比分类 如事故分类 农业 林业 牧业 副业 渔业    ['农业', '林业', '牧业', '副业', '渔业']
 * datas对应行业的参数的二维数组，注意父数组长度对应对比分类 子数组长度应与分析项目长度相同，对应相应的对比项目下的分类的参数 如[[10,20,50],[30,80,100],[500,600,200],[700,800,900],[100,200,300]]
 */
function initRaderMap(title,
	analyzeItems,
	datas, drawPictureDiv,pictrueType,yTitle) {
	var myChart = echarts.init(document.getElementById(drawPictureDiv));

	option = {
		title: {
			text: title,
			x: 'center'
		},
		tooltip: {
			trigger: 'axis'
		},
		legend: {
			orient: 'horizontal',/**默认为 水平布局，垂直布局为 'vertical*/
			padding: 3,
			data: analyzeItems,
			x: 'center',
			y: 'bottom'
		},
		toolbox: {
			show: true,
			feature: {
				mark: {
					show: true
				},
				dataView: {
					show: false,
					readOnly: false
				},
				magicType: {
					show: true,
					type: ['line', 'bar', 'stack', 'tiled']
				},
				restore: {
					show: false
				},
				saveAsImage: {
					show: true
				}
			}
		},
		calculable: true, /**启用拖拽重计算特性*/
		xAxis: [{
			type: 'category',
           /**	boundaryGap: false,*/
			data: getXAxisData(datas)
		}],
		yAxis: [{
			type: 'value',
			name:yTitle,
			axisLabel : {
                formatter: '{value}'      
            }
		}],
		series : getSeriesData(analyzeItems, datas, pictrueType)
	};
	myChart.setOption(option);
}
/**
 * 获取分析指标参数
 * @param {Object} kindItems 雷达图对比分类 如事故分类 农业 林业 牧业 副业 渔业    ['农业', '林业', '牧业', '副业', '渔业']
 * @param {Object} datas datas对应行业的参数的二维数组，注意父数组长度对应对比分类 子数组长度应与分析项目长度相同，对应相应的对比项目下的分类的参数 如[[10,20,50],[30,80,100],[500,600,200],[700,800,900],[100,200,300]]
 */
function getXAxisData(datas) {
	var indicatorArray = new Array(datas.length);
	for (var i = 0; i < indicatorArray.length; i++) {
		indicatorArray[i] = datas[i][0];
		datas[i].splice(0, 1);
	}
	return indicatorArray;
}
/**
 * 获取数据封装参数
 * @param {Object} analyzeItems 分析项目 字符串数组 如分析北京 上海 广州 的农业林业牧业副业渔业 ['北京','上海','广州']
 * @param {Object} datas对应行业的参数的二维数组，注意父数组长度对应对比分类 子数组长度应与分析项目长度相同，对应相应的对比项目下的分类的参数 如[[10,20,50],[30,80,100],[500,600,200],[700,800,900],[100,200,300]]
 */
function getSeriesData(analyzeItems, datas, pictrueType) {
	var seriesData = new Array(analyzeItems.length);
	for (var i = 0; i < analyzeItems.length; i++) {
		var dataItem = new Object();
		dataItem.data = getAnalyzeItem(i, datas);
		dataItem.name = analyzeItems[i];
		dataItem.type = pictrueType;
		dataItem.smooth = true;
		/**	dataItem.itemStyle = {normal: {areaStyle: {type: 'default'}}};*/
		seriesData[i] = dataItem;
	}
	/*此函数不支持IE8 暂时注释*/
	/**console.log(JSON.stringify(seriesData));*/
	return seriesData;
}
/**
 * 获取指定索引的所有分类数据
 * @param {Object} index 指定的索引
 * @param {Object} datas对应行业的参数的二维数组，注意父数组长度对应对比分类 子数组长度应与分析项目长度相同，对应相应的对比项目下的分类的参数 如[[10,20,50],[30,80,100],[500,600,200],[700,800,900],[100,200,300]]
 */
function getAnalyzeItem(index, datas) {
	var value = new Array(datas.length);
	for (var i = 0; i < value.length; i++) {
		value[i] = datas[i][index];
	}
	return value;
}
function drawPiePicture(drawPictureDiv, title, analyzeItems, datas) {
	/** 基于准备好的dom，初始化echarts图表*/
	myChart = echarts.init(document.getElementById(drawPictureDiv));
	var option = {
		title: {
			text: title,
			x: 'center',
			textStyle:{
				color: '#333',
				fontStyle: 'normal',
				fontWeight: 'normal',
				fontFamily: 'sans-serief',
				fontSize: 25
			}
		},
		tooltip: {
			trigger: 'item',
			formatter: "{a} <br/>{b} : {c} ({d}%)"
		},
		legend: {
			orient: 'vertical',
			x: 'left',
			data: analyzeItems,
			textStyle:{
				fontStyle: 'normal',
				fontWeight: 'normal',
				fontFamily: 'sans-serief',
				fontSize: 15
			}
		},
		toolbox: {
			show: true,
			feature: {
				mark: {
					show: true
				},
				dataView: {
					show: true,
					readOnly: false
				},
				magicType: {
					show: true,
					type: ['pie', 'funnel'],
					option: {
						funnel: {
							x: '25%',
							width: '50%',
							funnelAlign: 'left',
							max: 1548
						}
					}
				},
				restore: {
					show: true
				},
				saveAsImage: {
					show: true
				}
			}
		},
		calculable: true,
		series: [{
			type: 'pie',
			radius: '55%',
			center: ['50%', '60%'],
			data: datas,
			itemStyle: {
				normal: {
					label: {
						show: true,
						formatter: '{b} : {c}' + '吨标煤' + ' ({d}%)',
						textStyle:{
							fontStyle: 'normal',
							fontWeight: 'normal',
							fontFamily: 'sans-serief',
							fontSize: 15
						}
					},
					labelLine: {
						show: true
					}
				}
			}
		}]
	};
	/** 为echarts对象加载数据 */
	myChart.setOption(option);
}
function drawColummnPicture(echartsTitle, drawPictureDiv, analyzeItems) {
	/** 基于准备好的dom，初始化echarts图表*/
	myChart = echarts.init(document.getElementById(drawPictureDiv));;
	/** 指定图表的配置项和数据*/
	var option = {
		title: {
			text: echartsTitle,
			x: 'center'
		},
		tooltip: {
			trigger: 'axis'
		},
		legend: {
			orient: 'horizontal',
			padding: 3,
			data: analyzeItems,
			x: 'center',
			y: 'bottom'
		},
		xAxis: [{
			type: 'category',
			data: ['系统效率']
		}],
		yAxis: [{
			type: 'value'
		}],
		series: [{
			name: '',
			type: 'bar',
			data: []
		}]
	};
	/** 为echarts对象加载数据 */
	myChart.setOption(option);
}
function drawColummnMap(title,
						analyzeItems,
						datas, 
						drawPictureDiv,
						yTitle) {
	var myChart = echarts.init(document.getElementById(drawPictureDiv));

	option = {
		title: {
			text: title,
			x: 'center'
		},
		tooltip: {
			trigger: 'axis'
		},
		legend: {
			orient: 'horizontal',
			padding: 1,
			data: analyzeItems,
			x: 'center',
			y: 'bottom',
			itemWidth: 10,             /** 图例图形宽度*/
        	itemHeight: 5,  		   /** 图例图形高度*/
        	textStyle: {
            	fontSize:7                /** 图例文字大小*/
        	}
		},
		toolbox: {
			show: true,
			feature: {
				mark: {
					show: true
				},
				dataView: {
					show: false,
					readOnly: false
				},
				magicType: {
					show: true,
					type: ['line', 'bar', 'stack', 'tiled']
				},
				restore: {
					show: false
				},
				saveAsImage: {
					show: true
				}
			}
		},
		calculable: true, /**启用拖拽重计算特性*/
		xAxis: [{
			type: 'category',
			data: getXAxisDataOne(datas)
		}],
		yAxis: [{
			type: 'value',
			name:yTitle
		}],
		series: getSeriesDataOne(analyzeItems, datas)
	};
	myChart.setOption(option);
}
			function getXAxisDataOne(pictrueDatas) {
				var indicatorArray = new Array(pictrueDatas.length);
				for (var i = 0; i < indicatorArray.length; i++) {
					indicatorArray[i] = pictrueDatas[i][0];
					pictrueDatas[i].splice(0, 1);
				}
				return indicatorArray;
			}
			function getSeriesDataOne(testAnalyzeItems, pictrueDatas) {
				var seriesData = new Array(testAnalyzeItems.length);
				for (var i = 0; i < testAnalyzeItems.length; i++) {
					var dataItem = new Object();
					dataItem.data = getAnalyzeItem(i, pictrueDatas);
					dataItem.name = testAnalyzeItems[i];
					dataItem.type = 'bar';
					dataItem.smooth = true;
					seriesData[i] = dataItem;
				}
				return seriesData;
			}
			function testAnalyzeItems(index, pictrueDatas) {
				var value = new Array(pictrueDatas.length);
				for (var i = 0; i < value.length; i++) {
					value[i] = datas[i][index];
				}
				return value;
			}
```


## alertWin 早期作品

```javascript
    /**
     * @author 作者：张建新
     *弹窗封装类 alertWin(title, msg, w, h)
     * 示例写法：alertWin('提示信息','数据库修改失败',200,200);
     * @param {Object} title 标题
     * @param {Object} msg 错误信息以及提示信息
     * @param {Object} w 宽度
     * @param {Object} h 高度
     **/
    function alertWin(title, msg, w, h) {
        var titleheight = "22px"; /** 提示窗口标题高度 */
        var bordercolor = "#364751"; /** 提示窗口的边框颜色 */
        var titlecolor = "#FFFFFF"; /** 提示窗口的标题颜色 */
        var titlebgcolor = "#364751"; /** 提示窗口的标题背景色 */
        var bgcolor = "#FFFFFF"; /** 提示内容的背景色 */
    
        var iWidth = document.documentElement.clientWidth;
        var iHeight = document.documentElement.clientHeight;
        var bgObj = document.createElement("div");
        bgObj.style.cssText = "position:absolute;left:0px;top:0px;width:" + iWidth + "px;height:" + Math.max(document.body.clientHeight, iHeight) + "px;filter:Alpha(Opacity=30);opacity:0.3;background-color:#000000;z-index:101;";
        document.body.appendChild(bgObj);
    
        var msgObj = document.createElement("div");
        msgObj.style.cssText = "position:absolute;font:11px '宋体';top:" + (iHeight - h) / 2 + "px;left:" + (iWidth - w) / 2 + "px;width:" + w + "px;height:" + h + "px;text-align:center;border:1px solid " + bordercolor + ";background-color:" + bgcolor + ";line-height:22px;z-index:102;";
        document.body.appendChild(msgObj);
    
        var table = document.createElement("table"); /** */
        msgObj.appendChild(table);
        table.style.cssText = "margin:0px;border:0px;padding:0px;";
        table.cellSpacing = 0;
        var tr = table.insertRow(-1);
        var titleBar = tr.insertCell(-1);
        titleBar.style.cssText = "width:100%;height:" + titleheight + "px;text-align:left;padding:3px;margin:0px;font:bold 13px '宋体';color:" + titlecolor + ";border:1px solid " + bordercolor + ";cursor:move;background-color:" + titlebgcolor;
        titleBar.style.paddingLeft = "10px";
        titleBar.innerHTML = title;
        var moveX = 0;
        var moveY = 0;
        var moveTop = 0;
        var moveLeft = 0;
        var moveable = false;
        var docMouseMoveEvent = document.onmousemove; 
        var docMouseUpEvent = document.onmouseup;
        titleBar.onmousedown = function() {
            var evt = getEvent();
            moveable = true;
            moveX = evt.clientX;
            moveY = evt.clientY;
            moveTop = parseInt(msgObj.style.top);
            moveLeft = parseInt(msgObj.style.left);
    
            document.onmousemove = function() {
                if (moveable) {
                    var evt = getEvent();
                    var x = moveLeft + evt.clientX - moveX; /***/
                    var y = moveTop + evt.clientY - moveY;
                    if (x > 0 && (x + w < iWidth) && y > 0 && (y + h < iHeight)) {
                        msgObj.style.left = x + "px";
                        msgObj.style.top = y + "px";
                    }
                }
            };
            document.onmouseup = function() {
                if (moveable) {
                    document.onmousemove = docMouseMoveEvent; /** */
                    document.onmouseup = docMouseUpEvent;
                    moveable = false;
                    moveX = 0;
                    moveY = 0;
                    moveTop = 0;
                    moveLeft = 0;
                }
            };
        }
    
        var closeBtn = tr.insertCell(-1);
        closeBtn.style.cssText = "cursor:pointer; padding:2px;background-color:" + titlebgcolor;
        closeBtn.innerHTML = "<span style='font-size:15pt; color:" + titlecolor + ";'>×</span>";
        closeBtn.onclick = function() {
            document.body.removeChild(bgObj);
            document.body.removeChild(msgObj);
        }
        var msgBox = table.insertRow(-1).insertCell(-1);
        msgBox.style.cssText = "font:10pt '宋体';";
        msgBox.colSpan = 2;
        msgBox.innerHTML = msg;
        /**
        获得事件Event对象，用于兼容IE和FireFox 
        * 为了兼容 Ie 进行获取ie 浏览器窗口的操作控制
        */
        function getEvent() {
            return window.event || arguments.callee.caller.arguments[0];
        }
    }
```