---
layout: post
categories: docker
title: 为 Docker for Mac 配置 HTTP 代理
date: 2018-04-16 13:44:26 +0800
description: 为 Docker for Mac 配置 HTTP 代理 方便Pull 镜像，提高速度。
keywords: dockerproxy bash
---


## 需求

使用 Docker for Mac 时，让 docker pull/push 走本地的代理端口。



## 简介
Docker 是一种虚拟容器，可以提高开发过程中环境的一致性，便于维护，并且大大促进了知识以及经验的分享和积累。Docker for Mac 是 Docker 专门为 macOS 开发的图形客户端，使用起来相当方便，不过如果从 docker 仓库拉取镜像，那速度真是感人，2.0G 的镜像，估计一天一夜都下载不完，中途中断更是让人崩溃。更悲剧的是，docker 的网络流量无法走设置好的代理，比如本地的 http://127.0.0.1:8118，无论是在 Docker for Mac 图形界面里设置，或者在环境变量里设置 http_proxy/https_proxy 。另外，DaoCloud 的 Docker 加速器以及阿里云的 Docker 国内镜像服务，都有一定的效果，但速度不甚理想，而且新的 docker 镜像可能还没来得及缓存。


## 解决方案

1. 把本地代理监听地址调整到 0.0.0.0

在本地设置代理时，把监听地址从 127.0.0.1 调整到 0.0.0.0。不同的软件设置方法不太相同，我用的privoxy。

* Mac 安装 `privoxy`


```bash

brew install privoxy
# brew
/usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"

# Addr: https://brew.sh

```
  * privoxy 配置文件

```bash
    cat /usr/local/etc/privoxy/config
    # 转发地址
    forward-socks5   /               localhost:1083 .
    # 监听地址
    listen-address  0.0.0.0:8118
    # local network do not use proxy
    forward         192.168.*.*/     .
    forward            10.*.*.*/     .
    forward           127.*.*.*/     .
```

2. 将本机的 lo0 网络接口设置成任意 IP

lo0 接口被称作 loopback，默认是 127.0.0.1，这里可以设置任意可用 IP 地址：

```bash
sudo ifconfig lo0 alias 10.200.10.1/24    # 24是掩码
```

3. 在 Docker for Mac 里设置代理 你可以看看自己的本地代理是在监听哪个端 `privoxy`

```bash
 # 监听地址
    listen-address  0.0.0.0:8118
```
不同的软件使用的端口可能不一样，根据情况设置，最后的代理地址就是上面第二部分的 IP 加上这个端口，如下图：

![](http://112firshme11224.test.upcdn.net/demos/f6838f96-60db-4123-a907-450351efda7a.png)


4. 结果测试图

![](http://112firshme11224.test.upcdn.net/demos/d9600bf0-ce6a-41f6-a50b-ce569ea6bfb8.png)

转载请注明出处，本文采用 [CC4.0](http://creativecommons.org/licenses/by-nc-nd/4.0/) 协议授权
