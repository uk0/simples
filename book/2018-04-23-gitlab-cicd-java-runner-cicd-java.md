---
layout: post
categories: git lab
title: Gitlab JC/CD 第二篇
date: 2018-04-23 10:45:02 +0800
description: spring-cloud java cicd gitlab
keywords: spring git lab java
---


## just now
* 上个文章介绍了 如何配置Gitlab 的Runner，这次介绍如何将Runner 利用起来。
* 创建一个`Java`项目，配置它的`gitlab-ci`文件.
* 下面的这个`yaml`文件是用来配置从源码到 Jar/Tar的一个过程，在提交到另一个仓库的yaml
* 还有另一个 `yaml`是从 源码-->Jar／tar-->镜像

```yaml
# cache 这个参数用于定义全局那些文件将被 cache 到下一个  stages
# 调试开启
before_script:
 - pwd
 - env
 ##
 ## Assuming you created the SSH_KNOWN_HOSTS variable, uncomment the
 ## follo wing two lines.
 ##
 - mkdir -p /root/.ssh/
 - chmod  700 /root/.ssh/
 - echo -e "StrictHostKeyChecking no \nUserKnownHostsFile /dev/null" > ~/.ssh/config
cache:
  key: $CI_PROJECT_NAME-$CI_COMMIT_REF_NAME-$CI_COMMIT_SHA
  paths:
    - build/
    - /data/repo
stages:
  - build-jar
  - build-release
  - build-dev
variables:
  DOCKER_DRIVER: overlay2
  MAVEN_OPTS: "-Dmaven.repo.local=/data/repo"
  CI_DEBUG_TRACE: "true"
build-Java:
  image: "registry.cn-hangzhou.aliyuncs.com/emos_prod/centos7-jdk8-maven3-git-1.8:latest"
  stage: build-jar
  script:
    - chmod u+x ./maven-build-all.sh
    - ./maven-build-all.sh
    - ls -a build/
  tags:
    - build_dev

release-jar-release:
  image: "registry.cn-hangzhou.aliyuncs.com/emos_prod/centos7-jdk8-maven3-git-1.8:latest"
  stage: build-release
  script:
    - git clone http://$GITLAB_USER:$GITLAB_PASS@gitlab-cicd.com/release/build-space.git
    - cd build-space && rm -rf * && cp -r ../build/* .
    - git config --global user.name "root"
    - git config --global user.email "root@emos.com"
    - git add --all .
    - git commit -m "Gitlab CI Auto Builder master"
    - git push --force --quiet http://$GITLAB_USER:$GITLAB_PASSgitlab-cicd.com/release/build-space.git master:master
  tags:
    - build_dev
  only:
    - master

release-jar-dev:
  image: "registry.cn-hangzhou.aliyuncs.com/emos_prod/centos7-jdk8-maven3-git-1.8:latest"
  stage: build-dev
  script:
    - git clone http://$GITLAB_USER:$GITLAB_PASS@gitlab-cicd.com/dev/build-space.git
    - cd build-space && rm -rf * && cp -r ../build/* .
    - git config user.name "root"
    - git config user.email "root@emos.com"
    - git add --all .
    - git commit -m "Gitlab CI Auto Builder dev"
    - git push --force --quiet http://$GITLAB_USER:$GITLAB_PASS@gitlab-cicd.com/dev/build-space.git dev:dev
  tags:
    - build_dev
  only:
    - dev
```

* 文件主要是为了给Gitlab中的某个项目绑定一个`Job`运行这个Job的就是我们在上次讲的在`gitlab` `runner` 。

* 看一下目录结构

![](http://112firshme11224.test.upcdn.net/demos/7382a109-79c5-41c3-91ba-fa5d946ac61f.png)

* 这个总体配置很简单都是`yaml`文件规范，主要还是项介绍一下里面的`cache` cache 是做CICD避免不掉的东西，可以用来 将编译好的文件传送到下一个 `stage` 实现方式大概是 将你要`cache`的包 打成一个`zip`包,启动下一个`stage`在进行`unzip`到当初的目录地址。

* 在以上的配置中 `GITLAB_USER`,`GITLAB_PASS`

![](http://112firshme11224.test.upcdn.net/demos/66b3a3ce-0df4-4c69-8c43-f3738001aaa1.png)

* 设置运行的环境变量

![](http://112firshme11224.test.upcdn.net/demos/6c890392-8690-4384-bc57-ffc38f9cae86.png)

* 推荐创建全局Cache 。

转载请注明出处，本文采用 [CC4.0](http://creativecommons.org/licenses/by-nc-nd/4.0/) 协议授权
