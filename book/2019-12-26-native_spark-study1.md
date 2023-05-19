---
layout: post
categories: spark
title: Rust Native Spark Sample Install 
date: 2019-12-26 10:34:51 +0800
description: Native Spark for rust fast spark 
keywords: fast spark native spark
---

本文简单介绍Rust版本的Spark 安装 以及dev版本的测试。


### 介绍 Native Spark For Rust


* github 地址 ： https://github.com/rajasekarv/native_spark


### Sample 

```bash
## clone project 
git clone https://github.com/rajasekarv/native_spark
## Use Ubuntu18.04

    curl -ssf https://sh.rustup.rs | sh  # 如果失败 curl https://sh.rustup.rs > a.sh && ./a.sh -y 强制执行
     apt install openssl
     apt install openssl-dev
     apt install openssl-sys
     apt install libssl-dev
## clone capnproto 
## 序列化工具包
git clone https://github.com/sandstorm-io/capnproto.git
cd capnproto
sudo  apt install autoreconf
sudo  apt install autoconf
sudo  apt install automake
sudo  apt install libtool
autoreconf -i
./configure 
make -j4 check
sudo make install 

## 进入项目测试

cd native_spark
## 安装nightly 版本
rustup install nightly
## 强行覆盖一波
rustup override set nightly
## run example 
cargo run --example make_rdd

## cd native_spark/docker
# build docker 镜像
./build_image.sh
# 启动测试节点 并且将当前的 target 目录映射到 容器内的 /home/dev
./test_cluster.sh

```


## 图示


* 看一下代码

![](http://112firshme11224.test.upcdn.net/blog/2019-12-26-10-46-14.png!100)

* 启动测试节点 

![](http://112firshme11224.test.upcdn.net/blog/2019-12-26-10-47-39.png!100)


* 进入Master


![](http://112firshme11224.test.upcdn.net/blog/2019-12-26-10-48-47.png!100)


* 测试native_spark

![](http://112firshme11224.test.upcdn.net/blog/2019-12-26-10-50-53.png!100)


* 在容器内测试native_spark 集群 `./file_read --deployment_mode distributed --local_ip=0.0.0.0`

```bash
RUST_BACKTRACE=full ./file_read --deployment_mode distributed --local_ip=0.0.0.0
```

![](http://112firshme11224.test.upcdn.net/blog/2019-12-26-11-07-11.png!100)


#### 未完结

转载请注明出处，本文采用 [CC4.0](http://creativecommons.org/licenses/by-nc-nd/4.0/) 协议授权
