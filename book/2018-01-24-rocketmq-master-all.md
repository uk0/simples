---
layout: post
categories: java rocketmq
title: rocketmq 多Master集群搭建
date: 2018-01-24 15:01:53 +0800
description: 搭建一个多Master的rocketmq 集群
keywords: rocketmq 
---

 搭建一套多个Master的集群（RocketMQ）
    为了项目测试使用。

## quick start 

  * 准备RocketMQ 4.2.0 tar
  * 准备JDK 1.8

## config

  * 三个broker的配置 broker-a，broker-b，broker-c
    broker-a.properties
    broker-b.properties
    broker-c.properties
    修改下面配置文件的 `brokerName`即可
 
    ```bash
        #所属集群名字
        brokerClusterName=rocketmq-aom
        #broker名字，注意此处不同的配置文件填写的不一样
        brokerName=broker-a|b|c
        #0 表示 Master，>0 表示 Slave
        brokerId=0
        #nameServer地址，分号分割
        namesrvAddr=tod1:9876;tod2:9876;tod3:9876
        #在发送消息时，自动创建服务器不存在的topic，默认创建的队列数
        defaultTopicQueueNums=4
        #是否允许 Broker 自动创建Topic，建议线下开启，线上关闭
        autoCreateTopicEnable=true
        #是否允许 Broker 自动创建订阅组，建议线下开启，线上关闭
        autoCreateSubscriptionGroup=true
        #Broker 对外服务的监听端口
        listenPort=10911
        #删除文件时间点，默认凌晨 4点
        deleteWhen=04
        #文件保留时间，默认 48 小时
        fileReservedTime=120
        #commitLog每个文件的大小默认1G
        mapedFileSizeCommitLog=1073741824
        #ConsumeQueue每个文件默认存30W条，根据业务情况调整
        mapedFileSizeConsumeQueue=300000
        #destroyMapedFileIntervalForcibly=120000
        #redeleteHangedFileInterval=120000
        #检测物理文件磁盘空间
        diskMaxUsedSpaceRatio=88
        #存储路径
        storePathRootDir=/home/aom/data/rocketmq/store
        #commitLog 存储路径
        storePathCommitLog=/home/aom/data/rocketmq/store/commitlog
        #消费队列存储路径存储路径
        storePathConsumeQueue=/home/aom/data/rocketmq/store/consumequeue
        #消息索引存储路径
        storePathIndex=/home/aom/data/rocketmq/store/index
        #checkpoint 文件存储路径
        storeCheckpoint=/home/aom/data/rocketmq/store/checkpoint
        #abort 文件存储路径
        abortFile=/home/aom/data/rocketmq/store/abort
        #限制的消息大小
        maxMessageSize=65536
        #flushCommitLogLeastPages=4
        #flushConsumeQueueLeastPages=2
        #flushCommitLogThoroughInterval=10000
        #flushConsumeQueueThoroughInterval=60000
        #Broker 的角色
        #- ASYNC_MASTER 异步复制Master
        #- SYNC_MASTER 同步双写Master
        #- SLAVE
        brokerRole=ASYNC_MASTER
        #刷盘方式
        #- ASYNC_FLUSH 异步刷盘
        #- SYNC_FLUSH 同步刷盘
        flushDiskType=ASYNC_FLUSH
        #checkTransactionMessageEnable=false
        #发消息线程池数量
        #sendMessageThreadPoolNums=128
        #拉消息线程池数量
        #pullMessageThreadPoolNums=128
     ```

## start
  * 修改脚本参数

    ```bash
       #runbroker.sh
        JAVA_OPT="${JAVA_OPT} -server -Xms1g  -Xmx1g -Xmn512m -XX:PermSize=128m -XX:MaxPermSize=320m"
       #runserver.sh
       JAVA_OPT="${JAVA_OPT} -server -Xms1g -Xmx1g -Xmn512m -XX:PermSize=128m -XX:MaxPermSize=320m
    ```
  * 启动RocketMQ

     ```bash

       #先启动Nameserv（三台机器）
        nohup sh mqnamesrv &
        nohup sh mqnamesrv &
        nohup sh mqnamesrv &
       # 启动 broker（三台机器）
        nohup sh mqbroker -c /home/rocketmq/conf/2m-noslave/broker-a.properties >/dev/null 2>&1 & 

        nohup sh mqbroker -c /home/rocketmq/conf/2m-noslave/broker-b.properties >/dev/null 2>&1 & 

        nohup sh mqbroker -c /home/rocketmq/conf/2m-noslave/broker-c.properties >/dev/null 2>&1 & 
     ```

## test

  * test RocketMQ

    ```text
        [root@tod1 bin]# ./mqadmin clusterList -n 192.168.1.149:9876
        Java HotSpot(TM) 64-Bit Server VM warning: ignoring option PermSize=128m; support was removed in 8.0
        Java HotSpot(TM) 64-Bit Server VM warning: ignoring option MaxPermSize=128m; support was removed in 8.0
        #Cluster Name     #Broker Name            #BID  #Addr                  #Version                #InTPS(LOAD)       #OutTPS(LOAD) #PCWait(ms) #Hour #SPACE
        rocketmq-aom      broker-a                0     192.168.1.149:10911    V4_2_0_SNAPSHOT          0.00(0,0ms)         0.00(0,0ms)          0 421335.24 0.0128
        rocketmq-aom      broker-b                0     192.168.1.150:10911    V4_2_0_SNAPSHOT          0.00(0,0ms)         0.00(0,0ms)          0 421335.24 0.0116
        rocketmq-aom      broker-c                0     192.168.1.151:10911    V4_2_0_SNAPSHOT          0.00(0,0ms)         0.00(0,0ms)          0 421335.24 0.0116
    ```

转载请注明出处，本文采用 [CC4.0](http://creativecommons.org/licenses/by-nc-nd/4.0/) 协议授权
