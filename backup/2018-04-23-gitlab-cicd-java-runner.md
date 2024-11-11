---
layout: post
categories: git lab
title: Gitlab JC/CD 第一篇
date: 2018-04-23 10:45:02 +0800
description: spring-cloud java cicd gitlab
keywords: 
---


## just now 

* 不再介绍如何安装`Gitlab`
* 从`runner` 开始介绍


## 配置Runner

* 管理员用户登录 第一步

![](http://112firshme11224.test.upcdn.net/demos/e0a5ba5a-3b77-4871-a06d-b280040fcf8d.png)


* 第二步进入Runner配置项

![](http://112firshme11224.test.upcdn.net/demos/f5f17d52-ff64-44a9-8f10-61a171ddbaf4.png)


* 获得runner的注册信息


![](http://112firshme11224.test.upcdn.net/demos/7a0076af-57c3-4886-9bfb-e7edc4f03750.png)


* 启动一个Runner 并且配置

```bash

sudo docker run -d --name gitlab-java --restart always \
  -v /srv/gitlab-java/config:/etc/gitlab-runner \
  -v /var/run/docker.sock:/var/run/docker.sock \
  gitlab/gitlab-runner:latest

```

* 注册runner 到 Gitlab

```bash
sudo docker exec -it gitlab-java gitlab-runner register \
  --url "http://gitlab.emosa.com/" \
  --registration-token "hyNYzSsoJKjGqy1y4D-Q" \
  --executor docker \
  --description "Java-Runner" \
  --docker-image "debian" \
  --docker-volumes /var/run/docker.sock:/var/run/docker.sock


# gitlab-runner register
Please enter the gitlab-ci coordinator URL:
# 示例：http://gitlab.emosa.com/
Please enter the gitlab-ci token for this runner:
# hyNYzSsoJKjGqy1y4D-Q
Please enter the gitlab-ci description for this runner:
# Java-Runner
Please enter the gitlab-ci tags for this runner (comma separated):
# java
这个地方还有一个配置 默认即可
#false
Whether to run untagged builds [true/false]:
# true
Please enter the executor: docker, parallels, shell, kubernetes, docker-ssh, ssh, virtualbox, docker+machine, docker-ssh+machine:
# docker
Please enter the default Docker image (e.g. ruby:2.1):
# debian

```

* 检查启动结果

![](http://112firshme11224.test.upcdn.net/demos/a8e8f6f4-2567-4434-af95-4de02ee55c8a.png)


* Gitlab Runner 就配置完成了。


转载请注明出处，本文采用 [CC4.0](http://creativecommons.org/licenses/by-nc-nd/4.0/) 协议授权
