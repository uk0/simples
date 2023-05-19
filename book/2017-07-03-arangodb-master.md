---
layout: post
title: arangodb/linux下搭建，可以搭建集群模式。
categories: arangodb,linux
description: linux
keywords: arangodb, linux,kubernetes
---


发这个贴的原因，是因为 圈叉三角方块 ，圈叉三角方块 。


## arangodb安装以及集群部署(集群采用kubernetes，Docker)
Arangodb

>https://github.com/arangodb

现在版本：3.1.23


:smile:


## 1. 安装部署


![arangodb](http://112firshme11224.test.upcdn.net/posts/arangodb/WX20170703-180430@2x.png)


## 2.使用介绍 （Building your own ArangoDB image）


[![dockeri.co](http://dockeri.co/image/_/arangodb)](https://registry.hub.docker.com/_/arangodb/)

We are auto generating docker images via our build system so the Dockerfile is a template. To build your own ArangoDB image:

```bash
cp Dockerfile.templ Dockerfile

```

Adjust @VERSION@ in the Dockerfile to the version of arangodb you want to have and issue:

```bash
docker build -t arangodb .
```

This will create an image named `arangodb`.

## Reference documentation

For user documentation please refer to our official docker hub documentation:

https://hub.docker.com/_/arangodb/

# 3.docker-compose mini 集群启动方式
   ` 各位稍等，稍后整理一并奉上。`
   
   - 大体思路,Docker-Compose 通过link进行docker容器的内部链接，来实现集群模式。
   
>Over