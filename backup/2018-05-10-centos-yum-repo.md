---
layout: post
categories: linux yum
title: 修改CentOS默认yum源为国内yum镜像源
date: 2018-05-10 14:05:42 +0800
description: 有时候CentOS默认的yum源不一定是国内镜像，导致yum在线安装及更新速度不是很理想。这时候需要将yum源设置为国内镜像站点。国内主要开源的开源镜像站点应该是网易和阿里云了.
keywords: centos7 yum源
---


有时候CentOS默认的yum源不一定是国内镜像，导致yum在线安装及更新速度不是很理想。这时候需要将yum源设置为国内镜像站点。国内主要开源的开源镜像站点应该是网易和阿里云了.


## 修改CentOS默认yum源为mirrors.163.com

* 首先备份系统自带yum源配置文件/etc/yum.repos.d/CentOS-Base.repo


```bash
[root@localhost ~]# mv /etc/yum.repos.d/CentOS-Base.repo /etc/yum.repos.d/CentOS-Base.repo.backup
```

* 进入yum源配置文件所在的文件夹

```bash
[root@localhost ~]# cd /etc/yum.repos.d/
```
* 下载163的yum源配置文件

```bash
# 自己根据版本进行下载
# 163
wget -O /etc/yum.repos.d/CentOS-Base.repo http://mirrors.163.com/.help/CentOS7-Base-163.repo
wget -O /etc/yum.repos.d/CentOS-Base.repo  http://mirrors.163.com/.help/CentOS6-Base-163.repo
wget -O /etc/yum.repos.d/CentOS-Base.repo  http://mirrors.163.com/.help/CentOS5-Base-163.repo
# aliyun
wget -O /etc/yum.repos.d/CentOS-Base.repo  http://mirrors.aliyun.com/repo/Centos-7.repo
wget -O /etc/yum.repos.d/CentOS-Base.repo  http://mirrors.aliyun.com/repo/Centos-5.repo
wget -O /etc/yum.repos.d/CentOS-Base.repo  http://mirrors.aliyun.com/repo/Centos-6.repo

# 运行yum makecache生成缓存

yum makecache
```

# 更新动作

```bash
#更新系统就会看到以下mirrors.163.com信息
已加载插件：fastestmirror, refresh-packagekit, security
设置更新进程Loading mirror speeds from cached hostfile
* base: mirrors.163.com
* extras: mirrors.163.com
* updates: mirrors.163.com
```






转不注明出处，老子就不写了，本文采用 [CC4.0](http://creativecommons.org/licenses/by-nc-nd/4.0/) 协议授权
