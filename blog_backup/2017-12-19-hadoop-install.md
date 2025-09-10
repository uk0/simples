---
layout: post
title: Hadoop集群配置HA。
categories: Hadoop
description: Hadoop集群配置zookeeper HA
keywords: Hadoop, zookeeper
---

# hadoop集群配置HA


# 准备 

 * JDK1.8
 * NTP
 * zookeeper
 * NTPDATE
 * 配置目录


```bash
scp -r zookeeper/ dscn$:/home/zookeeper

scp -r hadoop/ dscn$:/home/hadoop

mkdir  -p /home/program/hadoop/  data logs pids  # hadoop 文件路径

mkdir  -p  /home/program/zookeeper data logs pids # zookeeper 文件

```

```bash 
hdfs haadmin -transitionToActive nn1 --forcemanual
解决启动节点都为 `standby`

```

# 环境变量 Hadoop

```bash
  echo "
        export JAVA_HOME=/usr/java/jdk1.8.0_73
        export HADOOP_HOME=/home/dscs/hadoop
        export HADOOP_LOG_DIR=/home/program/hadoop/logs
        export HADOOP_PID_DIR=/home/program/hadoop/pids
        export HADOOP_IDENT_STRING=dscs
        export YARN_LOG_DIR=/home/program/hadoop/logs
        export YARN_PID_DIR=/home/program/hadoop/pids
        export YARN_IDENT_STRING=dscs
        export CLASSPATH=.:$JAVA_HOME/lib/dt.jar:$JAVA_HOME/lib/tools.jar
        PATH=$HADOOP_HOME/bin:$HADOOP_HOME/sbin:$JAVA_HOME/bin:$PATH
        " >>~/.bash_profile
    source ~/.bash_profile
```



# 环境变量 Java 

```bash
    export JAVA_HOME=$(dirname $(dirname $(readlink -f $(which javac))))
```


# zookeeper config 

```bash

    dataDir=/usr/local/program/zookeeper/data
    dataLogDir=/usr/local/program/zookeeper/logs
    clientPort=2181
    tickTime=2000
    initLimit=5
    syncLimit=2
    server.1=dscn1:2888:3888
    server.2=dscn2:2888:3888
    server.3=dscn3:2888:3888

```


```bash
    for i in 1 2 3
    do 

        ssh dscn$i mkdir -p /home/program/zookeeper/data/ echo $i > /home/program/zookeeper/data/myid

    done 
```

# hadoop 启动流程 


 * 启动 zookeeper (3个节点)
        
        ```bash
            ./zkServer.sh start ../conf/zoo.cfg
        ```
 * 启动 JournalNode (3个节点) 
        ```bash
            sbin/hadoop-daemon.sh start journalnode
        ```
 * 在其中一个namenode上格式化：
        ```bash
             bin/hdfs namenode -format -bjsxt 
        ```
 * 启动刚刚格式化的namenode:

        ```bash
            sbin/hadoop-daemon.sh start namenode
        ```
 * 在另一个没有格式化的namenode上执行
      
      ```bash
            bin/hdfs namenode -bootstrapStandby 
            
      ```
 * 启动第二个namenode
     
     ```bash
         sbin/hadoop-daemon.sh start namenode
     ```

 * 启动第二个datanode(3个节点)
     
     ```bash
         sbin/hadoop-daemon.sh start datanode
     ```


 * 在其中一个namenode上初始化zkfc：

    ```bash
        bin/hdfs zkfc -formatZK
    ```

 * 启动 JobHistory

  ```bash
    sbin/mr-jobhistory-daemon.sh start historyserver
  ```

 # 异常描述 
  * Hadoop2.0的HA机制判断两者元数据不一致时，为了防止脑裂，所以使两者都处于standby状态,
  * 场景：两台namenode都是standby。
  * zk集群里只是保存，hadoop集群的ha状态，而不保存数据，所以在ha机制出现问题时，可以使用此条命令：
  * 停止集群 hadoop

  ```bash
        hdfs zkfc -formatZK
  ```

 # 插曲 scp 发现存储不够了


 ```bash

 [root@dscn1 hadoop]# vgs
  VG   #PV #LV #SN Attr   VSize    VFree
  VG00   1  10   0 wz--n- <199.51g <1.26g


  lvextend -L +173G /dev/mapper/VG00-lvroot


  xfs_growfs /dev/mapper/VG00-lvroot

 ```


# 读取文件里面的文件路径参数【转译】

```bash

HBASE_TMP_DIR=$(cat $HBASE_PACKAGE_PATH/config | grep HBASE_TMP_DIR | awk -F "=" '{print$2}')

HBASE_TMP_DIR_STR=$(echo $HBASE_TMP_DIR | sed 's#\/#\\\/#g')

sed -i "s/HBASE_TMP_DIR/$HBASE_TMP_DIR_STR/" $HBASE_HOME/conf/hbase-site.xml

```


# 读取节点配置合理配置主从节点

```bash
OLD_IFS="$IFS"
IFS=","
arr=($HBASE_REGION)
IFS="$OLD_IFS"
for s in ${arr[@]}
do
    echo "$s" >>$HBASE_HOME/conf/regionservers
done

```