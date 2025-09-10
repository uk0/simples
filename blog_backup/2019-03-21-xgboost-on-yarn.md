---
layout: post
categories: xgboost
title: xgboost on yarn (cdh,cloudera)
date: 2019-03-21 13:20:52 +0800
description: xgboost install on yarn 
keywords: xgboost
---



Xgboost install by cluster yarn



## yum 安装基础的依赖

```bash
yum install gcc-c++ fuse-devel git  hadoop-libhdfs-devel
```



## 配置环境变量

```bash
export MVN_HOME=/opt/soft/apache-maven-3.6.0
export PATH=$MVN_HOME/bin:$PATH
export XGB_HOME=/home/download
export PATH=$XGB_HOME:$PATH
export HDFS_LIB_PATH=${XGB_HOME}/xgboost-packages/libhdfs
export LD_LIBRARY_PATH=${XGB_HOME}/xgboost-packages/lib64:$JAVA_HOME/jre/lib/amd64/server:/${XGB_HOME}/xgboost-packages/libhdfs:$LD_LIBRARY_PATH
export HADOOP_HOME=/opt/cloudera/parcels/CDH/lib/hadoop
export HADOOP_COMMON_HOME=$HADOOP_HOME
export HADOOP_HDFS_HOME=/opt/cloudera/parcels/CDH/lib/hadoop-hdfs
export HADOOP_MAPRED_HOME=/opt/cloudera/parcels/CDH/lib/hadoop-yarn
export HADOOP_YARN_HOME=$HADOOP_MAPRED_HOME
export HADOOP_CONF_DIR=$HADOOP_HOME/etc/hadoop
```


## xgboost git clone  and install 

```bash
mkdir -p xgboost-packages/lib64

cd xgboost-packages


git clone --recursive https://github.com/dmlc/xgboost
cd xgboost
cp make/config.mk ./
vim config.mk
add line # HADOOP_HOME = /usr/lib/hadoop

mkdir build
cd build
cmake .. USE_HDFS=ON # 开启HDFS
make -j4

```


## install python package on env

```bash

cd python-package/ && python setup.py install 
```
## install jvm-package on env

```bash
cd jvm-packages/ && mvn install:install-file -Dfile=xgboost4j-spark-0.83-jar-with-dependencies.jar -DgroupId=ml.dmlc -DartifactId=xgboost4j-spark -Dversion=0.7 -Dpackaging=jar
```











转载请注明出处，本文采用 [CC4.0](http://creativecommons.org/licenses/by-nc-nd/4.0/) 协议授权
