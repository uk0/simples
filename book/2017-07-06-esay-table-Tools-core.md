---
layout: post
title:  简易Table封装类,所有的Table都需要通过此工具操作合并col，row
categories: Handtable,HandtableExcel
description: 回顾
keywords: javascript, Handtable
---


发这个贴的原因，是因为 是因为 是因为啥来着，忘了 那就是因为你，因为你 。


## 代码

```javascript
/**
 * divId:插入表格的divId
 * listTitle:插入表格标题
 * listData:插入表格数据
 */
function addTable(divId, listTitle, listData,n) {
	var num = parseInt(Math.random()*100);
	var tableDivId = "table"+num.toString();
	var div = document.getElementById(divId);
	div.innerHTML = "";
	var strTitle = "";
	var scendTitle = "";
	var strDate = "";
	for (var i = 0; i < listTitle.length; i++) {
		if (listTitle[i].indexOf(":") != -1) {
			var strs = listTitle[i].split(":");
			strTitle += '<th>' + strs[0] + '</th>';
		} else {
			strTitle += '<th>' + listTitle[i] + '</th>';
		}
	}
	for (var i = 0; i < listTitle.length; i++) {
		if (listTitle[i].indexOf(":") != -1) {
			var strs = listTitle[i].split(":");
			scendTitle += '<th>' + strs[1] + '</th>';
		} else {
			scendTitle += '<th>' + listTitle[i] + '</th>';
		}
	}
	/**alert(listData.length);*/
	for (i = 0; i < listData.length; i++) {
		strDate += '<tr>';
		for (j = 0; j < listData[0].length; j++) {
			strDate += '<td>' + listData[i][j] + '</td>';
		}
		strDate += '</tr>';
	}
	var table = '<table class="t1" id='+tableDivId+' width="100%"  cellspacing="0" cellpadding="0" ><thead>' + '<tr>' + strTitle + '</tr><tr>' + scendTitle + '</tr></thead><tbody>' + strDate + '</tbody></table>';
	div.innerHTML = table;
	uniteTdCells(tableDivId,n);
	uniteTdRells(tableDivId);
}
/*
 * 合并行
 * tableId:需要操作表格的Id
 */
function uniteTdCells(tableDivId,n) {
	var table = document.getElementById(tableDivId);
	var nLeftCol = n;
	for (var r = 0; r < table.rows.length; r++) {
		for (var c = 0; c < table.rows[r].cells.length; c++) {
			/**选择要合并的列序数，去掉默认全部合并*/
			if (r >= 2 && c >= nLeftCol) {
				continue;
			}
			for (j = r + 1; j < table.rows.length; j++) {
				var cell1 = table.rows[r].cells[c].innerHTML;
				var cell2 = table.rows[j].cells[c].innerHTML;
				if (cell1 == cell2) {
					table.rows[j].cells[c].style.display = 'none';
					table.rows[r].cells[c].rowSpan++;
				} else {
					break;
				}
			}
		}
	}
}

/*
 *合并列
 *tableId：操作表格ID
 */
function uniteTdRells(tableDivId) {
	var table = document.getElementById(tableDivId);
	for (var c = 0; c < table.rows[0].cells.length; c++) {
		for (j = c + 1; j < table.rows[0].cells.length; j++) {
			var strOne = table.rows[0].cells[c].innerHTML;
			var strTwo = table.rows[0].cells[j].innerHTML;
			if (strOne == strTwo) {
				table.rows[0].cells[j].style.display = 'none';
				table.rows[0].cells[c].colSpan++;
			}
		}
	}
}



function uniteTdCellsOne( n) {
	var table = document.getElementById(tableId);
	var nLeftCol = n;
	for (var r = 0; r < table.rows.length; r++) {
		for (var c = 7; c < table.rows[r].cells.length; c++) {
			/**             	/**选择要合并的列序数，去掉默认全部合并*/
			/** 	        	if (r >= 2 && c >= nLeftCol)*/
			/** 	        	{*/
			/** 	        		continue; */
			/** 	        	}*/
			for (j = r + 1; j < table.rows.length; j++) {
				var cell1 = table.rows[r].cells[c].innerHTML;
				var cell2 = table.rows[j].cells[c].innerHTML;
				if (cell1 == cell2) {
					table.rows[j].cells[c].style.display = 'none';
					table.rows[r].cells[c].rowSpan++;
				} else {
					break;
				}
			}
		}
	}
}

```