---
layout: post
title: ActorDB/linux下搭建，可以搭建集群模式。
categories: ActorDB,linux
description: linux
keywords: ActorDB, linux,kubernetes
---


发这个贴的原因，是因为@#¥%……UIO@！#*……！（@ 


## ActorDB安装以及集群部署(集群采用kubernetes)

ActorDB

现在版本：0.10.21

## 1. 安装部署
* 1.1. 官网下载linux版本安装包（rpm包），上传至服务器并使用`rpm -ivh`进行安装
* 1.2. 安装好后，会自动创建 个文件夹：
    `/etc/actordb`		----	配置文件
    `/var/log/actordb`	----	默认日志目录
    `/var/lib/actordb`	----	默认数据存储目录
* 1.3. 修改配置文件`app.config`以变更数据目录、日志目录以及thrift和Mysql监听的端口等配置信息
     修改配置文件`vm.args`中`-name`的内容指定当前节点的node名称，格式为：<名称>@<ip地址>
     在`vm.args`中增加 `+S 8` 来限制`Erlang`的`Scheduler`只能使用8核，否则可能无法正常启动。
*  1.4. 执行 `actordb start` 启动ActorDB
*  1.5. 修改初始化脚本'init.sql'
    * 1.5.1. `use config` # 使用config库
    * 1.5.2. `insert into groups values ('dscnCluster','cluster')` # 在groups表里增加名为dscnCluster的cluster
    * 1.5.3. `insert into nodes values ('<名称>@<ip地址>', 'dscnCluster')` # 向dscnCluster组内增加成员
    * 1.5.4. 重复3.5.2和3.5.3，直到所有的节点都添加至集群中
    * 1.5.5. `CREATE USER 'root' IDENTIFIED BY '<pwd>'` # 创建root用户，并设置其密码为<pwd>
    * 1.5.6. `commit` # 提交
* 1.6. 执行初始化脚本 `actordb_console -f /etc/actordb/init.sql` 完成初始化
* 1.7. 建立库表
    * 1.7.1. `use schema`
    * 1.7.2. `actor <dbName> [kv]`
    * 1.7.3. `create table <tableName> (field type [constraint], field2 type2 [constraint2].....)` # 与传统DDL相同

## 2. 服务实例
* 2.1. `Cluster`：ActorDB集群中可定义多个`Cluster``，各Cluster间数据是独立的，相当于数据分片
* 2.2. `Node`：每个Cluster中包含多个Node（即ActorDB实例），一个Cluster中的所有Node互为副本，数据同步。

## 3. 存储
* 3.1. 所有的表操作都必须经由某一个（insert可以是多个）Actor，而所有的Actor都有一个所属的`ActorType`。
* 3.2. 用户在定义`ActorType`时，可以定义该``ActorType``所使用的表结构（与传统数据库一样），并可以创建该Type的Actor（数量不限）
* 3.3. 一个`ActorType`中的所有actor共享相同的表结构和表关系，`ActorType`中的所有actor都只能访问自己actor中的数据
    就像微信的收藏功能一样，不用用户共享同样的收藏夹结构，但无法看到其他人收藏的内容。
* 3.4. Actor拥有过多的数据会影响性能，官网建议每个Actor拥有的数据最好控制在1G一下，每个Actor的操作TPS最好低于1000/s
* 3.5. 可以使用user、filename或其他可识别的内容对某一`ActorType`中的数据在逻辑上做进一步的分割
* 3.6. 操作示例：
    `actor <dbName>(<userName>) create;`
    `insert into <tableName> values(<val1>,<val2>......);`
    `commit`

## 4. Thrift客户端连接
* 4.1. 在actordb的github中下载thrift接口生成种子：`adbt.thrift`，使用thrift官网的生成工具生成实现thrift接口的类及相关类
* 4.2. 使用Thrift客户端进行连接，客户端使用详见`thrift`官网。
* 4.3. 使用生成的类作为`Client`，调用sql方法直接执行（sql举例如：`actor metadata(lcy);select count(*) from conf`）。API包括了同步和异步的接口。

:smile:

>感谢公司多位同事的帮助,此文档来自公司一位大神(妹子) 