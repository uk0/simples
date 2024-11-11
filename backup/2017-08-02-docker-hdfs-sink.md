---
layout: post
title:  Docker Development HDFS for Flink  Sink
categories: HDFS,Mongodb
description: 笔记
keywords: HDFS,flink
---


因工作需求所整合Flink + HDFS 作为一个Demo 帮助大家跳坑。
HDFS + Docker 采用集群主从模式。

# HDFS with Docker
   * Docker image for single `HDFS` node.。
   * `Only for development purpose.`
   * License: `MIT`


# start
    * docker pull mdouchement/hdfs

# Local build

```bash
$ docker build -t mdouchement/hdfs .
```

# Running HDFS container

```bash
# Running and get a Bash interpreter
$ docker run -p 22022:22 -p 8020:8020 -p 50010:50010 -p 50020:50020 -p 50070:50070 -p 50075:50075 -it mdouchement/hdfs

# With NFS
$ docker run -p 22022:22 -p 8020:8020 -p 50010:50010 -p 50020:50020 -p 50070:50070 -p 50075:50075 -p 111:111 -p 2049:2049 -it mdouchement/hdfs

# Running as daemon
$ docker run -p 22022:22 -p 8020:8020 -p 50010:50010 -p 50020:50020 -p 50070:50070 -p 50075:50075 -d mdouchement/hdfs

```


* Ports
     * Portmap -> 111
     *  NFS -> 2049
     *  HDFS namenode -> 8020 (hdfs://localhost:8020)
     *  HDFS datanode -> 50010
     *  HDFS datanode (ipc) -> 50020
     *  HDFS Web browser -> 50070
     *  HDFS datanode (http) -> 50075
     *  HDFS secondary namenode -> 50090
     *  SSH -> 22

# Contributing
  * All PRs are welcome.
      * 1.Fork it ( https://github.com/[my-github-username]/gemsupport/fork )
      * 2.Create your feature branch (git checkout -b my-new-feature)
      * 3.Commit your changes (git commit -am 'Add some feature')
      * 4.Push to the branch (git push origin my-new-feature)
      * 5.Create a new Pull Request

# dev
```java
 /**HDFS Config*/
  configuration.set("fs.default.name", "hdfs://localhost:8020");
```

# `Docker`这个是为了快速测试代码，以及快速开发。
   * 真的很方便，用后直接销毁。
   * 也可以持久化储存数据。
   * 详情请查看`Docker` 官网。


* Owner `breakEval13`
* https://github.com/breakEval13
* `Docker`  https://docker.com