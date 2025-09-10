---
layout: post
title: 将“yyyyMMdd”格式的时间字符串转换为“yyyy-MM-dd HH:mm:ss”格式（Java）[原创]
categories: java
description: crate.io
keywords: crate.io,databases
---


将“yyyyMMdd”格式的时间字符串转换为“yyyy-MM-dd HH:mm:ss”格式（Java）


date：2017-10-31
author：zhangjianxin

#
[*] 解决方案


  * start
  ```java
    date = DateUtils.parseDate(value, new String[]{"yyyyMMdd"});
    value = DateFormatUtils.format(date, "yyyy-MM-dd HH:mm:ss");
  ```


* 以上操作经过验证可以直接拿去。
* auther `breakEval13`
* https://github.com/breakEval13