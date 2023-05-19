---
layout: post
categories: docker
title: install docker for centos6 
date: 2017-06-30 21:05:11 +0800
description: 
keywords: 
---


零、前言
本来，对于 CentOS 系列，Docker 官方要求要 CentOS7.0 及以上系统版本，但是有时候迫不得已，还是要在已有的 CentOS6.x 的系统上安装。

比如我遇到的：要在一台已有的 CentOS6.5 的服务器上部署一个 Java 应用，该 Java 应用基于 Java8 和 Mysql5.6 开发，都用到了相应的特性。但是，已有的 CentOS6.5 上已经在跑着几个 PHP 和 Java 应用了，装的 JDK 是 JDK7，Mysql 装的是 Mysql5.1，第一感觉就是要 GG 了！

好在，Docker 强大的软硬件隔离特性要发挥威力了！好，就走 Docker 这条路！剩下的问题也就变成了 Docker & Docker Compose 的安装问题了！

本文基于以下两篇文章而成，并结合自己遇到的坑写的。
- How To Install Docker on CentOS 6
- Docker and docker-compose in CentOS 6

## 一、安装 Docker
        这一节主要来自于：How To Install Docker on CentOS 6

        以下 命令都是在 root 权限下运行的。

        Docker是Enterprise Linux（EPEL）的额外包的一部分，EPEL是用于RHEL发行版的非标准包的社区库。首先，我们将安装EPEL仓库：

 ```bash
 rpm -iUvh http://dl.fedoraproject.org/pub/epel/6/x86_64/epel-release-6-8.noarch.rpm
 ```

### 更新 yum 源：
  ```bash
  yum update -y
  ```

   * 现在让我们通过安装 docker-io 软件包来安装 Docker：

   ```bash
    yum -y install docker-io
   ```

   * 安装完成后，我们需要启动 Docker 守护进程：

    ```bash
    service docker start
    ```

   * 最后，可选地，我们让 Docker 在服务器启动时启动：

   ```bash
    chkconfig docker on
   ```
 
   * 经过以上步骤，Docker 就已经安装成功了


## 二、安装 Docker Compose
* 这一节主要来自于：Docker and docker-compose in CentOS 6

Docker Compose 是 Python 写的一个可以同时管理多个 Docker容器 的工具。因为是 CentOS6.x，所以没办法直接安装该工具，要通过 Python 的 pip 管理器工具来安装。

一般 CentOS6.5 都自带了 Pyhton2.6，所以要先安装 Python2.7。

### 2.1、安装开发工具和Python2.7
```bash
# yum groupinstall "Development"
# yum install zlib-devel bzip2-devel openssl-devel ncurses-devel sqlite-devel
# wget https://www.python.org/ftp/python/2.7.8/Python-2.7.8.tar.xz
# tar -xvf Python-2.7.8.tar.xz
# cd Python-2.7.8
# ./configure --prefix=/usr/local
# make && make altinstall
# mv /usr/bin/python /usr/bin/python2.6.6
# ln -s /usr/local/bin/python2.7 /usr/bin/python
vi /usr/bin/yum CHANGE to python2.6.6
```

>以上是通过源码方式安装了 Python2.7，倒数第三行的命令，把系统原有的 python 程序重命名成 python2.6.6；倒数第二行命令，把 python2.7 程序链接/替换到原有的 python 程序。倒数第一行并不是一个命令，只是说明要把 /usr/bin/yum yum 程序的第一行声明，改成原有的 python2.6.6。即：

```bash

#!/usr/bin/python
import sys
......

#modify

#!/usr/bin/python2.6.6
import sys
......
```
* 一句话来说就是：升级了 Python 版本，但是 yum 程序还是用回原来的 Python 版本。

### 2.2、安装 pip

```bash
#原文里的 get-pip.py 这个下文件的下载地址过期了，新地址是https://bootstrap.pypa.io/get-pip.py所以命令改为如下：

wget https://bootstrap.pypa.io/get-pip.py
python get-pip.py 
```
2.3、安装 docker-compose
```bash
 pip install docker-compose
```

* 顺利的话，以上步骤走完，就完成了 docker-compose 的安装了，验证一下是否安装成功了，查看版本信息：

```bash
docker-compose version
```

## 三、我遇到的坑

### 3.1、pip 版本问题

* 执行 2.3节 的 pip install docker-compose 是遇到错误提示：

> pkg_resources.DistributionNotFound: The 'pip==8.1.0' distribution was not found and is required by the application
经 关于pip安装时提示pkg_resources.DistributionNotFound 错误问题 和CentOS升级Python2.7启发，我运行 whereis pip 命令发现调用的是 /usr/bin/pip，查看其内容：


```bash
cat /usr/bin/pip
```

打印出如下：
```python

#!/usr/bin/python
# EASY-INSTALL-ENTRY-SCRIPT: 'pip==8.1.0','console_scripts','pip'
__requires__ = 'pip==8.1.0'
import sys
from pkg_resources import load_entry_point

if __name__ == '__main__':
    sys.exit(
        load_entry_point('pip==8.1.0', 'console_scripts', 'pip')()
    )
```
但是考虑到：

```bash
 ll /usr/local/lib/python2.7/site-packages | grep pip
打印如下：

# drwxr-xr-x 10 root root  4096 Jun 16 13:08 pip
# drwxr-xr-x  2 root root  4096 Jun 16 13:08 pip-9.0.1.dist-info
```

就把 /usr/bin/pip 里的所有 8.1.0 改为 9.0.1，然后 pip 就正常了。具体为啥，我也不太清楚了，不懂 Python 的路过。。

我猜是因为，python get-pip.py 安装了 pip 之后，但是 /usr/bin/pip 里面调用的还是旧版本好的 pip 导致的。
3.2、docker-compose 版本问题
上面 2.3节 的 pip install docker-compose 执行之后就安装了 docker-compose。
但是当运行 docker-compose build 的时候，就会提示我们的 Dcoker 版本太低要求升级 Docker！默认安装的 docker-compose 版本太高了！显然我们目前的处境是无法再升级 Docker 的了，只能降 docker-compose 的版本。

* 先看一下我们已经安装的 Docker 版本：
```bash
docker -v
```

```text
Docker version 1.7.1, build 786b29d/1.7.1
Google 了一下，在 StackOverflow 上找到个降 docker-compose 版本的方法：ERROR: The Docker Engine version is less than the minimum required by Compose
经查 Docker Compose Github Docs，发现 docker-compose 1.5.2 版本是兼容 Docker 1.7.1 的：Note that Compose 1.5.2 requires Docker 1.7.1 or later.。
好了，开始降级 docker-compose，先卸载：
```

```bash
pip uninstall docker-compose
```
再安装指定版本：
```bash
 pip install docker-compose==1.5.2
```
 * 至此，docker-compose 降版本成功！

转载请注明出处，本文采用 [CC4.0](http://creativecommons.org/licenses/by-nc-nd/4.0/) 协议授权
