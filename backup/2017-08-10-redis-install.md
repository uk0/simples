---
layout: post
title: Redis安装,For Centos6
categories: Redis
description: 笔记
keywords: Redis,linxu,centos
---


没事看看文章而已～

#  Redis安装
date：2016-11-15
author：YangRaner 

[TOC]

#  1.检查安装依赖程序
```bash
yum install -y tcl gcc-c++ wget
```
#  2.下载Redis
```bash
mkdir /opt/redis/
cd /opt/redis/
wget http://download.redis.io/releases/redis-3.0.4.tar.gz
```
#  3.解压Redis
```bash
tar -xzvf redis-3.0.4.tar.gz
```
#  4.编译安装Redis
```bash
cd redis-3.0.4
make
make install	
-------------------------------------------------------------
make install安装完成后，会在/usr/local/bin目录下生成下面几个可执行文件，它们的作用分别是：
redis-server：Redis服务器端启动程序
redis-cli：Redis客户端操作工具。也可以用telnet根据其纯文本协议来操作
redis-benchmark：Redis性能测试工具
redis-check-aof：数据修复工具
redis-check-dump：检查导出工具
```
#  5.配置Redis
```bash
cp redis.conf /etc/   复制配置文件到/etc/目录
vim /etc/redis.conf    修改redis文件
----------------------------------------------
daemonize yes     修改daemonize配置项为yes，使Redis进程在后台运行
requirepass Ejtone   # 设置密码，以提供远程登陆
```
#  6.启动Redis
```bash
cd /usr/local/bin		进入bin目录下
./redis-server /etc/redis.conf		启动Redis
ps -ef | grep redis			查看redis进程
----------------------------------------------
root  18443  1  0 13:05 ?    00:00:00 ./redis-server *:6379 

```
# 7.添加开机启动项
```bash
# 让Redis在服务器重启后自动启动，需要将启动命令写入开机启动项
echo "/usr/local/bin/redis-server /etc/redis.conf" >>/etc/rc.local
```
**以上完成安装**

# 8.Redis配置参数
> * 在前面的操作中，我们用到了使Redis进程在后台运行的参数，下面介绍其它一些常用的Redis启动参数：**
`daemonize`：是否以后台daemon方式运行
`pidfile`：pid文件位置
`port`：监听的端口号
`timeout`：请求超时时间
`loglevel`：log信息级别
`logfile`：log文件位置
`databases`：开启数据库的数量
`save` * *：保存快照的频率，第一个*表示多长时间，第二个*表示执行多少次写操作。在一定时间内执行一定数量的写操作时，自动保存快照。可设置多个条件。
`rdbcompression`：是否使用压缩
`dbfilename`：数据快照文件名（只是文件名）
`dir`：数据快照的保存目录（仅目录）
`appendonly`：是否开启appendonlylog，开启的话每次写操作会记一条log，这会提高数据抗风险能力，但影响效率。
`appendfsync`：appendonlylog如何同步到磁盘。三个选项，分别是每次写都强制调用fsync、每秒启用一次fsync、不调用fsync等待系统自己同步


# 9.重启Redis
  * 9.1如果Redis已经配置为service服务，可以通过以下方式重启：

```bash
service redis restart
```

  * 9.2如果Redis没有配置为service服务，可以通过以下方式重启：
```bash
/usr/local/bin/redis-cli shutdown
/usr/local/bin/redis-server /etc/redis.conf
```


* 以上代码经过验证可以直接拿去修改调用
* Owner `breakEval13`
* https://github.com/breakEval13