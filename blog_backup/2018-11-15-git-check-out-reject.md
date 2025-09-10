---
layout: post
categories: git
title: git丢弃本地修改的所有文件（新增、删除、修改）
date: 2018-11-15 17:03:19 +0800
description: git丢弃本地修改的所有文件（新增、删除、修改）
keywords: git 丢弃本地修改
---


## git丢弃本地修改的所有文件（新增、删除、修改）


```bash
git checkout . #本地所有修改的。没有的提交的，都返回到原来的状态
git stash #把所有没有提交的修改暂存到stash里面。可用git stash pop回复。
git reset --hard HASH #返回到某个节点，不保留修改。
git reset --soft HASH #返回到某个节点。保留修改

git clean -df #返回到某个节点
git clean #参数
    -n #显示 将要 删除的 文件 和  目录
    -f #删除 文件
    -df #删除 文件 和 目录

git checkout . && git clean -xdf # 清除本地
```


原文：https://blog.csdn.net/leedaning/article/details/51304690 

转载请注明出处，本文采用 [CC4.0](http://creativecommons.org/licenses/by-nc-nd/4.0/) 协议授权
