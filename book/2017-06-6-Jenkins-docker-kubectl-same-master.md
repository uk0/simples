---
layout: post
title: jenkins + docker + git 持续集成
categories: [jenkins, docker]
description: jenkins + docker + git 持续集成
keywords: jenkins, docker
---

> git push 以后， jenkins 自动触发 代码打包，生成docker image , docker push 到 仓库，发布到环境里。


# 安装jenkins

这里不建议用 Docker 镜像，因为下面 Jenkins 自身会需要调用 Docker 来启动任务。

Jenkins + Docker + kubectl + git + Blue Ocean。

## 导入 jenkins 源

```

wget -O /etc/yum.repos.d/jenkins.repo http://pkg.jenkins.io/redhat-stable/jenkins.repo
rpm --import http://pkg.jenkins.io/redhat-stable/jenkins.io.key


yum -y install jenkins 
```


## 修改jenkins配置

```
vi /etc/sysconfig/jenkins

# 修改jenkins 目录
JENKINS_HOME="/opt/jenkins"

# 修改jenkins 端口
JENKINS_PORT="9999"
```

## 移动目录

```
# 将目录移动过来，否则程序报错
mv /var/lib/jenkins /opt/
```


## 启动服务

```
systemctl start jenkins
systemctl enable jenkins
```



## 访问WEB UI

```
http://myip:9999/ 
```



生成密钥

```
# 切换用户
su jenkins


# 生成key
ssh-keygen -t rsa -b 4096 -C "jenkins@git"

# 查看key信息
cat /home/jenkins/.ssh/id_rsa.pub
```

## jenkins 后台配置

进入 jenkins --> Credentials --> Add Credentials
![此处输入图片的描述][1]



选择 系统管理 -- > 管理插件

1. 添加 Gradle Plugin 插件
2. 添加 Git plugin 插件



```
常用插件

Build WIth Parameters   # 执行 构建 前手工输入参数

pipeline

Deploy Plugin   # build war 包以后部署

Email Extension Plugin  # 邮件发送

Multiple SCMs Plugin #多项目构建工具

```




下载慢可直接下载 hpi 文件，通过高级 导入插件安装

选择 系统管理 -- > Global Tool Configuration

安装JDK

![描述][2]


安装 Gradle 

![此处输入图片的描述][3]


安装 Git

![此处输入图片的描述][4]



创建项目  选择 自由风格 的项目

源码管理选择 Git

![此处输入图片的描述][5]



构建 选择  Invoke Gradle script

![此处输入图片的描述][6]




## 构建触发器


```
# 勾选 Poll SCM

# 每两分钟检查一次git代码是否有更新
H/2 * * * *
```









## 配置 邮件

首先必须安装 Email Extension Plugin 插件

系统管理 --> 系统设置 -- > Jenkins Location

配置系统管理员邮件地址 --- >  xxx@163.com


配置 Extended E-mail Notification


SMTP Server = 

点击高级

勾选 Use SMTP Authentication

输入 发送 用户 与 密码

填写 SMTP port


Default Content Type 选择 HTML (text/html)


```
Default Subject =  构建通知:$PROJECT_NAME - Build # $BUILD_NUMBER - $BUILD_STATUS!


Default Content = 

<b style="font-size:12px">(本邮件是程序自动下发的，请勿回复，<span style="color:red">请相关人员fix it,重新提交到git 构建</span>)<br></b><hr>

<b style="font-size: 12px;">项目名称：$PROJECT_NAME<br></b><hr>

<b style="font-size: 12px;">构建编号：$BUILD_NUMBER<br></b><hr>

<b style="font-size: 12px;">GIT版本号：${GIT_REVISION}<br></b><hr>

<b style="font-size: 12px;">构建状态：$BUILD_STATUS<br></b><hr>

<b style="font-size: 12px;">触发原因：${CAUSE}<br></b><hr>

<b style="font-size: 12px;">构建日志地址：<a href="${BUILD_URL}console">${BUILD_URL}console</a><br></b><hr>

<b style="font-size: 12px;">构建地址：<a href="$BUILD_URL">$BUILD_URL</a><br></b><hr>

<b style="font-size: 12px;">变更集:${JELLY_SCRIPT,template="html"}<br></b><hr>

```

![此处输入图片的描述][8]






项目 -- > 构建后操作 --- > 添加 Editable Email Notification

拉到最下面 --- > 点击 Advanced Settings...

![此处输入图片的描述][9]

Recipient List 添加 收件邮箱 多个邮件以空格 隔开

![此处输入图片的描述][10]



原文作者 [小炒肉](https://jicki.me/ "Title") .
 

  [1]: http://jicki.me/images/posts/jenkins/4.png
  [2]: http://jicki.me/images/posts/jenkins/1.png
  [3]: http://jicki.me/images/posts/jenkins/2.png
  [4]: http://jicki.me/images/posts/jenkins/3.png
  [5]: http://jicki.me/images/posts/jenkins/5.png
  [6]: http://jicki.me/images/posts/jenkins/6.png
  [7]: http://jicki.me/images/posts/jenkins/7.png
  [8]: http://jicki.me/images/posts/jenkins/8.png
  [9]: http://jicki.me/images/posts/jenkins/9.png
  [10]: http://jicki.me/images/posts/jenkins/10.png


