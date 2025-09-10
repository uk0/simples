---
layout: post
title: Centos kafka 消息队列 
categories: kafka
description: Centos kafka 消息队列 
keywords: kafka
---


# kafka 消息队列

> kafka new Version

## 安装配置

查询下载最新版本 kafka 
http://kafka.apache.org/downloads.html

```
wget http://mirror.bit.edu.cn/apache/kafka/0.8.2.0/kafka-0.8.2.0-src.tgz
tar zxvf kafka-0.8.2.0-src.tgz
mv kafka-0.8.2.0-src /opt/local/kafka
cd /opt/local/kafka
./gradlew jar
```


提示: 
错误: 找不到或无法加载主类 org.gradle.wrapper.GradleWrapperMain
需要先安装 gradle

```
http://www.scala-lang.org/files/archive/scala-2.10.4.tgz
tar zxvf scala-2.10.4.tgz
mv scala-2.10.4 /usr/lib64/scala
```

设置环境变量 

```
vi /etc/profile
```

```
export SACLA_HOME=/usr/lib64/scala/
export PATH=$SACLA_HOME/bin:$PATH
source /etc/profile
```

然后再执行 gradlew jar

```
./gradlew jarAll
```


如果 jarAll 会报错，java 版本不能为1.8 不然会报不兼容的错误，请使用1.7版本


./gradlew jar --stacktrace  --info --debug

 

创建日志目录

```
mkdir -p /opt/local/kafka/logs
```
 

编辑配置文件

 
```
vim config/server.properties
```


将 log.dirs=/tmp/kafka-logs
改为

```
log.dirs=/opt/local/kafka/logs
```

将 zookeeper.connect=localhost:2181
改为

```
zookeeper.connect=172.24.0.100:2181
```

启动程序

```
nohup /opt/local/kafka/bin/zookeeper-server-start.sh /opt/local/kafka/config/zookeeper.properties &
nohup /opt/local/kafka/bin/kafka-server-start.sh /opt/local/kafka/config/server.properties &
```


创建主题

```
/opt/local/kafka/bin/kafka-topics.sh --create --zookeeper 192.168.20.200:2181 --replication-factor 1 --partitions 1 --topic LJ200
```

查看现有主题

```
/opt/local/kafka/bin/kafka-topics.sh --list --zookeeper 192.168.20.200:2181
```