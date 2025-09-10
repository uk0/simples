---
layout: post
title:  Handtable封装类,所有的Table都需要通过此工具操作表格
categories: Handtable,HandtableExcel
description: 回顾
keywords: javascript, Handtable
---


发这个贴的原因，是因为 是因为 是因为啥来着，忘了 那就是因为你，因为你 。


## 代码

```javascript
/**
 * @author 张建新
 * @数据表格加载
 * @callback 无
 * @DataType tableTools.init(@Coulnms) <-调用方式
 * @Coulnms List<Map[CID,TYPE]> result
 * (@所在div,@列数据(以及类型),@表格装载数据,@表头数据(中文),@是否激活右键菜单(右键可以开启只读模式),@是否表格只读)
 **/
var HandtableTools = {
	/**
	 * HandtableTools 工具初始化方法
	 * @param {Object} tableContainerDiv 表格载体Div
	 * @param {Object} tableColumnInfo 表格列信息
	 * @param {Object} contextMenuFlag 是否显示右键列表
	 * @param {Object} readOnlyFlag 表格是否只读
	 */
	initColumnsType: function(tableContainerDiv, tableColumnInfo, contextMenuFlag, readOnlyFlag) {
		/*初始化表格基本信息*/
		HandtableTools.tableContainerDiv = tableContainerDiv;
		HandtableTools.tableColumnInfo = tableColumnInfo;
		HandtableTools.contextMenuFlag = contextMenuFlag;
		HandtableTools.readOnlyFlag = readOnlyFlag;
	},
	/**
	 * 当前表格容器Div
	 */
	tableContainerDiv: null,
	/**
	 * 当前表格列基本信息
	 */
	tableColumnInfo: null,
	/**
	 * 是否显示右键列表
	 */
	contextMenuFlag: null,
	/**
	 * 表是否只读
	 */
	readOnlyFlag: null,
	/*
	 * 最终插入的数据 
	 */
	insertData: null,
	/**
	 * 表格控制器
	 */
	handsontableControl: null,
	/**
	 *初始化表格数据 
	 * @param {Object} gridData
	 */
	initTableData: function(gridData) {
		/**
		 * 表格列信息
		 */
		var columnsInfo = [];
		/**
		 * 表头中文信息
		 */
		var columnsHeadersInfo = [];
		/**
		 * 初始化插入数据
		 */
		HandtableTools.insertData = new Array();
		/**
		 * 遍历表格数据 加载到插入数据
		 */
		for (var i = 0; i < gridData.length; i++) {
			HandtableTools.insertData[i] = new Object();
			for (var j = 0; j < HandtableTools.tableColumnInfo.length; j++) {
				/**eval('HandtableTools.insertData[' + i +'].' + HandtableTools.tableColumnInfo[j].fieid + "=null");*/
				HandtableTools.insertData[i][HandtableTools.tableColumnInfo[j].fieid.toUpperCase()] = 0;
			}
		}
		/**
		 * 设置表格列信息
		 */
		for (var i = 0; i < HandtableTools.tableColumnInfo.length; i++) {
			columnsInfo.push({
				data: '',/**真实列名*/
				type: '',/**数据类型*/
				readOnly: '',/**列只读*/
				comment:''/**注释*/
			
			});
		}
		/**清空数据*/
		columnsHeadersInfo.length = 0;
		/**解析表格列信息 标注中文列名*/
		for (var c = 0; c < HandtableTools.tableColumnInfo.length; c++) {
			/**标注中文列名*/
			columnsHeadersInfo.push(HandtableTools.tableColumnInfo[c].alias);
		}
		
		for (var c = 0; c < HandtableTools.tableColumnInfo.length; c++) {
			/**设置列真实Data*/
			columnsInfo[c]['data'] = HandtableTools.tableColumnInfo[c].fieid.toUpperCase();
			/**设置列类型*/
			columnsInfo[c]['type'] = "" + HandtableTools.tableColumnInfo[c].genre + "";
			/**设置comments*/
			columnsInfo[c]['comment'] = "" + HandtableTools.tableColumnInfo[c].comat + "";
			/**判断是否是主键 如果是主键不允许修改*/
			if (HandtableTools.tableColumnInfo[c].ispk == "t" && HandtableTools.tableColumnInfo[c].ispk != "") {
				columnsInfo[c]['readOnly'] = true; /**可以修改列*/
			} else {
				columnsInfo[c]['readOnly'] = false; /**列不能修改*/
			}
		}
		console.log(HandtableTools.tableContainerDiv);
		HandtableTools.obs = columnsInfo;
		HandtableTools.handsontableControl = new Handsontable(HandtableTools.tableContainerDiv, {
			colHeaders: columnsHeadersInfo, /**显示列数据中文列名*/
			rowHeaders: true, /**显示行号*/
			/**minSpareRows: 1,/**预留行*/
			stretchH: 'all', /**自适应*/
			columnSorting: true, /**序列化*/
			/**contextMenu: contextMenu_flag,/**是否开启右键*/
			contextMenu:HandtableTools.contextMenuFlag,
			comments: false, /**开启单个表格注解*/
			readOnly: HandtableTools.readOnlyFlag, /**表格是否只读*/
			className: "htCenter htMiddle", /**基本样式*/
			removeRowPlugin: false, /**删除行插件*/
			fixedColumnsLeft: 3,/**设置固定列*/
			autoColumnSize:true, /**自适应列*/
			manualColumnMove: true, /**左右移动*/
			manualRowMove: false,	/**数据跟随移动*/
			columns: columnsInfo/**表格列信息*/
		});
		HandtableTools.handsontableControl.clear;/**清除数据*/
		HandtableTools.handsontableControl.destroy;/**清除DOM对象*/
		HandtableTools.handsontableControl.loadData(gridData); /**表格加载数据*/
	},
	obs: [],
	objData: function() {
		var tempData = HandtableTools.handsontableControl.getData();
		for (var i = 0; i < tempData.length; i++) {
			for (var j = 0; j < HandtableTools.obs.length; j++) {
				if(tempData[i][j]==null||tempData[i][j]==''){
					HandtableTools.insertData[i][HandtableTools.obs[j]["data"].toUpperCase()] = 0.0;
				}else{
					HandtableTools.insertData[i][HandtableTools.obs[j]["data"].toUpperCase()] = tempData[i][j];
				}
			}
		}
	}
}
```