---
layout: post
title: Javascript 一行能装逼的JavaScript代码 [fork]
categories: javascript
description: 一行能装逼的JavaScript代码
keywords: javascript,java
---


一行能装逼的JavaScript代码。

#  代码学习
date：2017-05-01
author：zhangjianxin

[TOC]

#  1. 一行神奇的JS代码，当时我就震惊了，这不就是传说中的ZB神奇么… … 哈哈。

   * 写本篇文章的缘由是之前看到了一段js代码，如下：

   *  (!(~+[])+{})[--[~+""][+[]]*[~+[]] + ~~!+[]]+({}+[])[[~!+[]]*~+[]]

   * 然后让大家运行，出来的结果让人有点出乎意料，请看：

   ![](http://112firshme11224.test.upcdn.net/blog/66042438.jpg)


   * > 太风骚了有木有！如果有人诋毁前端瞧不起JS的话，那就可以把这段代码发给他了~

    * > 不过话说回来了，这到底是什么原理呢？为什么一堆符号运算结果竟然能是两个字符，而且恰巧还是个sb！

    * > 其实靠的是js的类型转化的一些基本原理，本篇就来揭密”sb”是如何炼成的。相信你如果能把这个理清楚了，以后遇到类型转化之类的题目，就可以瞬间秒杀了。

# 一、JS运算符的优先级
  * 首先要运用到的第一个知识就是js运算符的优先级，因为这么长一段运算看的人眼花，我们必须得先根据优先级分成n小段，然后再各个击破。优先级的排列如下表：

  *  优先级从高到低：
   ![](http://112firshme11224.test.upcdn.net/blog/-365094845.png)


# 二、根据此规则，我们把这一串运算分为以下16个子表达式：

   ![](http://112firshme11224.test.upcdn.net/blog/-238288323.jpg)

   * 运算符用红色标出，有一点可能大家会意识不到，其实中括号[]也是一个运算符，用来通过索引访问数组项，另外也可以访问字符串的子字符，有点类似charAt方法，如：’abcd'[1] // 返回’b’。而且中括号的优先级还是最高的哦。

# 三、JS的类型转化

预处理结束，接下来需要运用的就是javascript的类型转化知识了。我们先说说什么情况下需要进行类型转化。当操作符两边的操作数类型不一致或者不是基本类型（也叫原始类型）时，需要进行类型转化。先按运算符来分一下类：

减号-，乘号*，肯定是进行数学运算，所以操作数需转化为number类型。

加号+，可能是字符串拼接，也可能是数学运算，所以可能会转化为number或string

一元运算，如+[]，只有一个操作数的，转化为number类型

下面来看一下转化规则。

1. 对于非原始类型的，通过ToPrimitive() 将值转换成原始类型：
ToPrimitive(input, PreferredType?)

可选参数PreferredType是Number或者是String。返回值为任何原始值。如果PreferredType是Number，执行顺序如下：

如果input为primitive，返回

否则，input为Object。调用 obj.valueOf()。如果结果是primitive，返回。

否则，调用obj.toString(). 如果结果是primitive，返回

否则，抛出TypeError

如果 PreferredType是String，步骤2跟3互换，如果PreferredType没有，Date实例被设置成String，其他都是Number

2. 通过ToNumber()把值转换成Number，直接看ECMA 9.3的表格
规则如下：
   ![](http://112firshme11224.test.upcdn.net/blog/823282545.png)

3. 通过ToString()把值转化成字符串， 直接看ECMA 9.8的表格

规则如下：
   ![](http://112firshme11224.test.upcdn.net/blog/-376644909.png)



* 以上操作经过验证可以直接拿去。
* Owner `breakEval13`
* https://github.com/breakEval13