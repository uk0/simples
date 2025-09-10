---
layout: post
title: javascript 实现的贪吃蛇 29line Code[瞎搞]
categories: javascript,贪吃蛇
description: 回顾笔记
keywords: javascript, Script
---


发这个贴的原因，是因为 是因为 是因为啥来着，忘了 那就是因为你，因为你 。


## 摘要: Code(我想装个b)

```html
<!doctype html><html><head></head><body>
<div id="box" style="width:256px;font:25px/25px 宋体;background:#000;color:#9f9;"></div><script>
var arr=[15];
var next=10;/**最好判断一下arr[1]!=arr[0]+x*/
document.onkeydown=function(e){next={'37':-1,'39':1,'38':-10,'40':10}[(e?e:event).keyCode]};
var run=setInterval(function(){
    var list=['-1&&arr[0]%10==0','1&&arr[0]%10==9','-10&&arr[0]<10','10&&arr[0]>189'];
    var temp=arr[0]+next;
    if(list.some(function(x){return eval('next=='+x)})||arr.indexOf(temp)>=0)
        return clearInterval(run)||alert('OVER');
    arr.unshift(temp);
    if(this.food==null||food==temp){
        list=[];/**(list=Array(200)).forEach(function(val,i){list[i]=i});*/
        for(var i=0;i<200;i++)
            list.push(i);
        list=list.filter(function(x){return arr.indexOf(x)<0});
        food=list[Math.round(Math.random()*(list.length-1))];/**food是全局变量*/
    }
    else
        arr.pop();
    var s='';
    for(var j=0;j<200;j++){/**这个循环还可以缩短*/
        if(j%10==0&&j!=0)
            s+='<br/>';
        if(arr.indexOf(j)>=0)
            s+='\u25a0';
        else if(j==food)
            s+='\u2605';
        else
            s+='&emsp;';
    }
    document.getElementById("box").innerHTML=s;
},200);
</script></body></html>
```



