---
layout: post
title: 使用travis-ci 来持续集成jekyll静态博客。
categories: Kubernetes,docker,travis
description: CICD。
keywords: Java, GitHub Pages
---

CICD (travis) for my blog 

# 使用travis-ci 来持续集成jekyll静态博客.

##简介
Travis CI提供一个在线的持续集成服务，用来构建托管在github上的代码。
许多知名的开源项目使用它来自动构建测试代码,它支持hexo的运行环境node.js。
原理很简单，Travis会在你每一次在github上提交代码后，生成一个Linux虚拟机来运行你配置好的任务,
生成和部署hexo只需要一个命令 hexo gd ，但是Travis CI需要有我们的github项目上传代码的权限，
所以我们需要使用SSH key来使Travis CI拥有权限。还有一些其他的问题也都参开hexo作者的博文解决了，
用Travis CI自动部署网站到Github。

```bash
#安装travis-ci的命令行工具

gem install travis

```
```bash
#登录Travis CI
travis login --auto #这里可以用token登录。

```

```bash
#使用travis命令行工具加密私钥
travis encrypt-file id_rsa --add
```

```bash
#travis CI解密配置,这部分内容需要配置在 .travis.yml 里,注意修改秘钥

- openssl aes-256-cbc -K $encrypted_26b4962af0e7_key -iv $encrypted_26b4962af0e7_iv
  -in id_rsa.enc -out ~/.ssh/id_rsa -d
```
.travis.yml

```yaml

language: node_js

node_js:
- '4.1'
before_install:
- openssl aes-256-cbc -K $encrypted_26b4962af0e7_key -iv $encrypted_26b4962af0e7_iv
  -in id_rsa.enc -out ~/.ssh/id_rsa -d
- chmod 600 ~/.ssh/id_rsa
- eval $(ssh-agent)
- ssh-add ~/.ssh/id_rsa
- cp ssh_config ~/.ssh/config
- git config --global user.name "pangjian"
- git config --global user.email "pangjian1010@gmail.com"

install:
- npm install hexo-cli -g
- npm install

script:
- hexo -version
- hexo clean
- hexo g
- hexo deploy
```

* 转载请注明出处：https://firsh.me/2017/06/13/travis-ci-myblog/

