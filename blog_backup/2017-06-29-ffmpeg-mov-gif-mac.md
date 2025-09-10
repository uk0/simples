---
layout: post
title: ffmpeg转换GIF图片参数,MOV转GIF[常用]
categories: Java, Scala,Hadoop,Flink
description: Java, Scala,Hadoop,Flink
keywords: Java, Scala,Hadoop,Flink
---

ffmpeg转换GIF图片参数,这里介绍的是MOV转为GIF来方便放到Blog上。




## 1. 介绍

ffmpeg除了在音视频方面牛叉之外，在转换gif动态图片方面也是不弱的，修改分辨率，改变帧率，颜色模式，添加水印字幕等等均不在话下；

下面几个操作命令是我这几天一一测试过，有需要的拿去，ffmpeg真的超级强大，轻而易举就能实现想要的效果；


## 2. 从视频第三十秒处开始截取视频，截取十秒钟片段

```bash
ffmpeg -ss 0 -t 150 -i Flinkdemo.mov  -s 800x600 -f gif -r 1 b.gif
```

![描述](http://112firshme11224.test.upcdn.net/posts/ffmpeg/QQ20170629-151034@2x.png)


## 3.频率限制
限制GIF体积大小

直接输出的gif体积会相对比较大，压缩不是很厉害，像在QQ中发送gif体积不能大于6M，所以要使用参数来降低gif图片体积大小；

-s 输出分辨率
-r 帧率

```bash

ffmpeg -i A.wmv -s 320x240  -r 10 B.gif

```

## 4.RGB颜色限制,输出RGB24位颜色gif图片

```bash
ffmpeg -i A.wmv -pix_fmt rgb24   B.gif
```

## 5.给gif动态图片添加字幕

```bash

ffmpeg -i A.gif   -vf "ass=test.ass"   B.gif

```