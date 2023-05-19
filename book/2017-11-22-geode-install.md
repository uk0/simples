---
layout: post
title: 为流处理加上一层高性能缓存。
categories: Apache
description: 搭建分布式缓存集群。
keywords: Apache geode, geode
---

# Apache geode install 

 * 下载二进制文件

```bash
wget http://apache.osuosl.org/geode/1.3.0/apache-geode-1.3.0.zip

```

* 复制到其他多个节点进行解压


```bash

unzip apache-geode-1.3.0.zip

```

* 修改为集群模式（注意：支持的模式有分组模式，集群模式，集群+分组模式）

* 二进制文件所在 `/home/hadmin/geode/`

* 创建相应的配置文件路径

```bash
mkdir -p cluster_config/cluster
```

* 将`/home/hadmin/geode/config`内的文件复制到刚才创建的目录内，并且进行相应的修改文件名。

```bash

gemfire.properties	mv -> cluster.properties

cache.xml	mv -> cluster.xml

```


* 结果

```bash

[root@dscn022 geode]# ls cluster_config/cluster/
cluster.properties  cluster.xml

```

* 启动主节点 locator 

* 其他节点 connect --locator=10.11.0.224[10334]

```text

create region --name=example-region --type=REPLICATE_PERSISTENT

start server --name=server1 --server-port=40411
start server --name=server2 --server-port=40412
start server --name=server3 --server-port=40413
start server --name=server4 --server-port=40414
start server --name=server5 --server-port=40415

query --query=" SELECT  * FROM /example-region h WHERE h.hoursPerWeek < 40 "



locator 定位器，相当于master-slave中的master，或者zookeeper，主要用于管理集群，和链接不同的server。

　　gfsh> start locator --name=locator1

server 服务器，可以部署在同一台机器，也可以部署在不同机器。在不同的机器上启动时，需要先用connect连接已启动的locator

　　connect --locator=ip[locator的port]

　　start server --name=server1

region 数据区域，或者叫表，是数据存储的基本单位，以下创建一个在集群内自动复制的，自动持久化的region，并持久化数据

　　create region --name=regionA --type=REPLICATE_PERSISTENT

　　put --region=regionA --key="1" --value="one"

OQL 类SQL的脚本，用来查数

　　query --query="select * from /regionA"

以上命令的执行默认是以集群为范围的，如果要单机执行，需要修改apache-geode\config\gemfire.properties文件中的属性：enable-cluster-configuration=true，改为false。



gfsh show missing-disk-stores 查询 disk 


start server --name=server13 --locators=10.11.0.224[10334]

start server --name=server20 --server-port=40431  --locators=10.11.0.224[10334]

```
