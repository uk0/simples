---
layout: post
title: nodejs + npm 版本切换(Mac/Windows/Linux)
categories: nodejs,npm,mac
description: npm版本切换。
keywords: npm, nodejs, npm版本切换
---

因项目需求需要多次切换npm以及nodejs，所写出这个文章以供参考。

# 使用npm n进行版本切换，对nodejs的版本控制。

## 安装N控制器。
```bash
sudo npm install -g n
```
## 安装nodejs的最新版本

```text

$ n 5.4.1  -->（安装node.js 5.4.1版本）

$ n latest  -->（安装node.js最新版本）

$ n stable   -->（安装node.js稳定版本）

$ n rm 4.0.0   -->（删除某个版本）

```
## 如图


![此处输入图片的描述](http://112firshme11224.test.upcdn.net/blog/QQ20170615-215055@2x.png)



```text

$ sudo n （用上下选择版本会车切换版本）

#以笔者装的为例，输入以上代码后会出现
    node/4.4.7
    node/6.5.0

```
![此处输入图片的描述](http://112firshme11224.test.upcdn.net/QQ20170615-214859@2x.png)

