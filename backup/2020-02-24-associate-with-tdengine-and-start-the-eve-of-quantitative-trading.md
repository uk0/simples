---
layout: post
categories: TDengine
title: TDengine在量化交易中提供数据分析[更新中]。
date: 2020-02-24 21:18:30 +0800
description: TDengine 在量化交易里面的数据分析实践
keywords: TDengine db python xgboost
---


###  记录一下量化交易的开始


### 准备工具 

* DB TDengine 
* 编程语言 Golang + Python 
* 环境 Unix Linux
* 交易所 Bitmex Huobi Binance Okex
* 每天增量300w数据
* 机器配置 I7-8700 + 32GB DDR4 + GTX2060 + GTX1080 (Xgboost) 磁盘 6TB（SATA6GB/s）


### 编译 TDengine （修改表字段限制）

* 下载或者clone 官方发布的Release

* 编译

```bash
yum install gcc gcc-c++ make maven  -y # centos 7

git clone https://github.com/taosdata/TDengine

cd TDengine

mkdir build && cd build

cmake .. && cmake --build .

make install

#Install the project...
#/usr/bin/cmake -P cmake_install.cmake
#-- Install configuration: "Debug"
#make install script: /root/TDengine-ver-1.6.5.5-fnk/packaging/tools/make_install.sh
#this is centos system
#source directory: /root/TDengine-ver-1.6.5.5-fnk
#binary directory: /root/TDengine-ver-1.6.5.5-fnk/build
#Start to install TDEngine...
#Created symlink from /etc/systemd/system/multi-user.target.wants/taosd.service to /etc/systemd/system/taosd.service.

#TDengine is installed successfully!

#To configure TDengine : edit /etc/taos/taos.cfg
#To start TDengine     : sudo systemctl start taosd
#To access TDengine    : use taos in shell

#TDengine is installed successfully!

# 出现这个说明已经安装好了

```

### 修改Taos DB 配置文件

```bash
# 将数据存储地址修改成 数据盘 日志改为 日志盘
# data file's directory
dataDir               /data

# log file's directory
logDir                /data/log

```

### 创建数据库

```sql
create database bitmex;
use bitmex;
--- 订单订阅数据表
CREATE TABLE bitmex.trade
 (
   ts                   TIMESTAMP,
   systemc_time         BIGINT,
   bitmex_timestamp     TIMESTAMP,
   side                 BINARY(4),
   size                 FLOAT,
   price                FLOAT,
   symbol               BINARY(16)
);


--- 交易深度数据
CREATE TABLE bitmex.order_l1_10 
  ( 
     ts               TIMESTAMP, 
     systemc_time     BIGINT, 
     symbol           BINARY(16), 
     bitmex_timestamp TIMESTAMP, 
     asks_1_p         FLOAT, 
     asks_1_s         FLOAT, 
     asks_2_p         FLOAT, 
     asks_2_s         FLOAT, 
     asks_3_p         FLOAT, 
     asks_3_s         FLOAT, 
     asks_4_p         FLOAT, 
     asks_4_s         FLOAT, 
     asks_5_p         FLOAT, 
     asks_5_s         FLOAT, 
     asks_6_p         FLOAT, 
     asks_6_s         FLOAT, 
     asks_7_p         FLOAT, 
     asks_7_s         FLOAT, 
     asks_8_p         FLOAT, 
     asks_8_s         FLOAT, 
     asks_9_p         FLOAT, 
     asks_9_s         FLOAT, 
     asks_10_p        FLOAT, 
     asks_10_s        FLOAT, 
     bids_1_p         FLOAT, 
     bids_1_s         FLOAT, 
     bids_2_p         FLOAT, 
     bids_2_s         FLOAT, 
     bids_3_p         FLOAT, 
     bids_3_s         FLOAT, 
     bids_4_p         FLOAT, 
     bids_4_s         FLOAT, 
     bids_5_p         FLOAT, 
     bids_5_s         FLOAT, 
     bids_6_p         FLOAT, 
     bids_6_s         FLOAT, 
     bids_7_p         FLOAT, 
     bids_7_s         FLOAT, 
     bids_8_p         FLOAT, 
     bids_8_s         FLOAT, 
     bids_9_p         FLOAT, 
     bids_9_s         FLOAT, 
     bids_10_p        FLOAT, 
     bids_10_s        FLOAT 
  ); 
```

### 启动采集工具 （开始实时的通过WebSocket接收数据）
```bash
# 采集部分两部分

# 采集端 Agent（通过SSL GRPC 将数据传输到 Server）多个 Agent 同时传输

# Server端 将通过GRPC 接收到的数据 通过http接口落库（这里用http接口是因为开发和部署的机器都是Unix的机器，数据库在内网的一台mac上运行。）

```
* 采集服务目录结构

![](http://112firshme11224.test.upcdn.net/blog/2020-02-25-22-06-16.png!100)

* 采集服务日志恢复
  > 下图所看到的是 GRPC服务断掉了，数据将会写入log文件并且有明确的标记，当grpc再次可用，将会启动日志恢复。

![](http://112firshme11224.test.upcdn.net/blog/2020-02-25-22-12-48.png!100)


### 通过Middware将数据实时的放入TDEngine（这里存在一层日志备份未来可能存储到Hdfs走FLink训练）



### 确认数据

```json
{
    "status": "succ",
    "head": [
        "ts",
        "systemc_time",
        "symbol",
        "bitmex_timestamp",
        "asks_1_p",
        "asks_1_s",
        "asks_2_p",
        "asks_2_s",
        "asks_3_p",
        "asks_3_s",
        "asks_4_p",
        "asks_4_s",
        "asks_5_p",
        "asks_5_s",
        "asks_6_p",
        "asks_6_s",
        "asks_7_p",
        "asks_7_s",
        "asks_8_p",
        "asks_8_s",
        "asks_9_p",
        "asks_9_s",
        "asks_10_p",
        "asks_10_s",
        "bids_1_p",
        "bids_1_s",
        "bids_2_p",
        "bids_2_s",
        "bids_3_p",
        "bids_3_s",
        "bids_4_p",
        "bids_4_s",
        "bids_5_p",
        "bids_5_s",
        "bids_6_p",
        "bids_6_s",
        "bids_7_p",
        "bids_7_s",
        "bids_8_p",
        "bids_8_s",
        "bids_9_p",
        "bids_9_s",
        "bids_10_p",
        "bids_10_s"
    ],
    "data": [
        [
            "2020-02-25 21:02:56.000",
            1582635776179,
            "XBTUSD",
            "2020-02-25 13:02:55.729",
            9625,
            49906,
            9625.5,
            82752,
            9626,
            10282,
            9626.5,
            1112,
            9627,
            109168,
            9627.5,
            47425,
            9628,
            84425,
            9628.5,
            114773,
            9629,
            69830,
            9629.5,
            34412,
            9624.5,
            2225271,
            9624,
            3024,
            9623.5,
            11051,
            9623,
            468523,
            9622.5,
            6502,
            9622,
            9985,
            9621.5,
            550000,
            9621,
            9295,
            9620.5,
            726837,
            9620,
            647086
        ]
    ],
    "rows": 1
}
```




### 提供Web界面 空和多 多家交易所对比，以及瞬时趋势


### Python 调用数据开启分析 基础分析 EMA 交叉 （四家交易所同时10秒计算一次）

### 每天计算趋势图

### 钉钉通知。

### 纠正模型

转载请注明出处，本文采用 [CC4.0](http://creativecommons.org/licenses/by-nc-nd/4.0/) 协议授权
