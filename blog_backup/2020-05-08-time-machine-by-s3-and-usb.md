---
layout: post
categories: macos
title: Time Machine Backup Fix ERROR
date: 2020-05-08 17:49:27 +0800
description: macos catinlina 10.15.4 修复Time Machine 备份错误问题
keywords: time machine
---


### 介绍背景

* 系统配置
![](http://112firshme11224.test.upcdn.net/blog/2020-05-08-17-50-43.png!/watermark/text/Zmlyc2gubWU=/font/simyou/align/southeast/size/18/color/FFFFFF/margin/5x5)


* Script  `https://github.com/uk0/TimeMachine-sparsebundle`




#### 问题介绍

* usb 外接硬盘 无法备份 （或者第二备份失败） 长时间插入硬盘 不会自动备份等问题。
* 将数据备份到 AWS 
* 工具 goofys test-macos ./test


#### 测试 


* Mount s3 

```bash
goofys test-macos ./test
```

* 运行makeImage.sh 脚本创建Image

```bash
 sh  ./makeImage.sh  600 /Volumes/tm
```

![](http://112firshme11224.test.upcdn.net/blog/2020-05-08-18-06-14.png!/watermark/text/Zmlyc2gubWU=/font/simyou/align/southeast/size/18/color/FFFFFF/margin/5x5)


* 在/Volumes/tm生成类似上图所示的文件（设备名字可能不同）
* 双击这个文件会发现多了一个磁盘

![](http://112firshme11224.test.upcdn.net/blog/2020-05-08-18-08-06.png!/watermark/text/Zmlyc2gubWU=/font/simyou/align/southeast/size/18/color/FFFFFF/margin/5x5)


* 最后告诉系统要使用这个硬盘来进行备份

```bash
 sudo tmutil setdestination /Volumes/Time\ Machine\ Backups
```

* 成功备份截图。

![](http://112firshme11224.test.upcdn.net/blog/2020-05-08-18-25-20.png!/watermark/text/Zmlyc2gubWU=/font/simyou/align/southeast/size/18/color/FFFFFF/margin/5x5)


* 以上的方法 适用与移动硬盘直接格式化后TimeMachine 无法备份，可以 用这个办法试试，将移动硬盘当作S3即可。



转载请注明出处，本文采用 [CC4.0](http://creativecommons.org/licenses/by-nc-nd/4.0/) 协议授权
