---
layout: post
categories: flink java hadoop zookeeper
title: Flink 高可用集群搭建 HA
date: 2018-01-24 16:16:32 +0800
description: 搭建 一个高可用的 Flink
keywords: flink
---

    搭建一个高可用的Flink 集群HA
        用于项目测试

## quick start

 * flink-conf.yaml

    ```yaml

        jobmanager.rpc.port: 6123


        jobmanager.heap.mb: 512


        taskmanager.heap.mb: 512

        taskmanager.numberOfTaskSlots: 2

        taskmanager.memory.preallocate: false


        parallelism.default: 1


        jobmanager.web.address: 0.0.0.0

        jobmanager.web.port: 8081


        jobmanager.archive.fs.dir: hdfs://ns/flink/completed_jobs/

        historyserver.web.address: 0.0.0.0


        historyserver.web.port: 8082


        historyserver.archive.fs.dir: hdfs://ns/flink/completed_jobs/


        historyserver.archive.fs.refresh-interval: 10000



        state.backend: filesystem


        state.backend.fs.checkpointdir: hdfs://ns/flink/checkpoints_backend
        state.backend.fs.savepointdir: hdfs://ns/flink/savepoints_backend
        state.checkpoints.dir: hdfs://ns/flink/checkpoints_data
        state.savepoints.dir: hdfs://ns/flink/savepoints_data


        taskmanager.tmp.dirs: /home/aom/data/flink/tmp
        env.log.dir: /home/aom/data/logs/flink




        fs.hdfs.hadoopconf: /home/aom/hadoop/etc/hadoop


        high-availability: zookeeper
        high-availability.zookeeper.quorum: dscn1:2181,dscn2:2181,dscn3:2181
        high-availability.zookeeper.storageDir: hdfs://ns/flink/recovery
        high-availability.zookeeper.path.root: /flink
    ```
  * masters
  ```bahs
    tod1:8081
    tod2:8081
  ```
  * slaves

  ```bash
    tod1
    tod2
    tod3
  ```

  * scp 到其他节点
  * 启动集群
    ```bash
        start-cluster.sh
    ```


转载请注明出处，本文采用 [CC4.0](http://creativecommons.org/licenses/by-nc-nd/4.0/) 协议授权
