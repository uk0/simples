---
layout: post
categories: httpd apache
title: Apache Httpd 2.4.43 Install 
date: 2020-04-03 16:10:24 +0800
description: 安装httpd
keywords: 安装httpd
---


安装httpd 过程记录。



### 准备

* 系统 Linux x86_64 gcc 版本4.8
* apr + apr-utils + pcre + httpd-2.4.43 `源码包`
* `yum install -y libxml2-devel`
* `yum install -y gcc gcc-c++`

### 安装开始

#### apr

```bash
wget -c http://archive.apache.org/dist/apr/apr-1.6.2.tar.gz
tar -xf apr-1.6.2.tar.gz
cd apr-1.6.2
./configure --prefix=/usr/local/apr
make && make install

```


#### apr-utils

```bash

wget -c http://archive.apache.org/dist/apr/apr-util-1.6.1.tar.gz
tar -xf apr-util-1.6.1.tar.gz
cd apr-util-1.6.1
./configure --prefix=/usr/local/apr-util --with-apr=/usr/local/apr/bin/apr-1-config
make && make install

```

#### pcre

```bash
wget -c https://ftp.pcre.org/pub/pcre/pcre-8.41.tar.gz
tar -xf pcre-8.41.tar.gz
cd pcre-8.41
./configure --prefix=/usr/local/pcre
make && make install
```

#### httpd

```bash
wget https://ftp.riken.jp/net/apache//httpd/httpd-2.4.43.tar.gz
tar -zxvf httpd-2.4.43.tar.gz
cd httpd-2.4.43
cd ..
cp -r apr-1.6.2 httpd-2.4.43/srclib/apr
cp -r apr-util-1.6.1  httpd-2.4.43/srclib/apr-util

./configure --prefix=/usr/local/apache2 --with-apr=/usr/local/apr --with-apr-util=/usr/local/apr-util --with-pcre=/usr/local/pcre --enable-mods-shared=most --enable-so --with-included-apr
make && make install

```




转载请注明出处，本文采用 [CC4.0](http://creativecommons.org/licenses/by-nc-nd/4.0/) 协议授权
