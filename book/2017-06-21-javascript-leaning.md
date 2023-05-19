---
layout: post
title: Javascript or nodejs 学习笔记。
categories: nodejs,nodejs,学习笔记
description: nodejs。
keywords: Javascript, nodejs
---

这一篇学习笔记，记录了今日所学的一些关于Javascript的高级用法。



## Main
 
* 今天主要学习的是Javascript的原型，闭包，以及this指针。

* 废话不多直接上代码；

## Javascript闭包

```javascript
/**
 * Created by zhangjianxin on 2017/6/21.
 */

function a(){
    var tmp=1;
    return function a_sub(x){
        console.log(x+(tmp++));
    }
}

var b = a;      /**此处a和b是一模一样的,b是从a拷贝了一份过去*/
var c1 = a();   /**此处a加上了括号，于是c1的函数其实就是 a中return回来的函数（此函数可以调用a外面的变量tmp）*/
var c2 = b();   /**为了证明a和b是同样的功能，且是复制了一份*/
c1(3);    /**输出4*/
c1(3);    /**输出5,——说明a并内存并未释放*/

c2(3);    /**输出4——说明a和b是复制关系(两个独立的对象)，内存中互相独立*/
/**另外一个例子***************************************************** * **/
var test = (function(){
    var tmp=1;
    return function(x){
        console.log(tmp+x);
    }
})
var d = test();     /**var objname = function(){...}和 function objname(){...}定义是一样的效果*/
d(3)

```

## Javascript this指针问题

```javascript
/**
 * Created by zhangjianxin on 2017/6/21.
 */



var fuckthis = {
    var1:"asdasd",

    func1:function () {
        var f1 = "1";
        console.log("func1"+ this.var1)
    },
    func2:function () {
        var f2 = "2";
        console.log("func2"+f2)
    },

    func3:function () {
        this.func1();
    }
}

fuckthis.func3();

```

## Javascript原型

```javascript

/**
 * Created by zhangjianxin on 2017/6/21.
 */


function preson() {
    console.log("preson");
}

var p = new preson();

preson.prototype.fuck = function () {
    console.log("fuck prototype");
}

preson.prototype.fuck2 = function () {
    console.log("fuck prototype 2222");
}


p.fuck2();

```