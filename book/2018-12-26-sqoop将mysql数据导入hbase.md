---
layout: post
categories: hbase sqoop
title: 使用Sqoop将mysql数据导入到Hbase[整表,部分数据]
date: 2018-12-26 18:08:22 +0800
description: sqoop 迁移数据
keywords: sqoop
---


使用sqoop 将数据从mysql 导入到hbase 分为 整张表 和 部分数据，其中部分数据基于 where条件等。


### 相关的命令

```bash
# 整个表导入
./bin/sqoop import --driver com.mysql.jdbc.Driver --connect "jdbc:mysql://rm-adfagga.mysql.rds.aliyuncs.com/parking?zeroDateTimeBehavior=convertToNull" --username user1 --password 123123 --table parking_ths_car_record --hbase-table parking_ths_car_record_all --column-family id --hbase-row-key id --split-by park_id --hbase-create-table

./bin/sqoop import --driver com.mysql.jdbc.Driver --connect "jdbc:mysql://rm-asfagag.mysql.rds.aliyuncs.com:3306/pp_parking?zeroDateTimeBehavior=convertToNull" --username user1 --password 123123 --table parking_car_user --hbase-table car_user --column-family id --hbase-row-key id --hbase-create-table

# 执行sql （任意）

./bin/sqoop eval --connect "jdbc:mysql://rm-asdasdfffaf.mysql.rds.aliyuncs.com/parking" --username user1 --password 123123  --query "SELECT * FROM parking_ths_car_record LIMIT 10"


# 驱动更换位置
/opt/cloudera/parcels/CDH-5.14.4-1.cdh5.14.4.p0.3/lib/sqoop/lib/
/opt/cloudera/parcels/CDH-5.14.4-1.cdh5.14.4.p0.3/lib/hadoop/lib/


./bin/sqoop import --driver com.mysql.jdbc.Driver --connect "jdbc:mysql://rm-aasdasdasd.mysql.rds.aliyuncs.com/pp_parking?zeroDateTimeBehavior=convertToNull" --username user1 --password 123123 --table parking_ths_car_record --hbase-table parking_ths_car_record --column-family id --hbase-row-key id --split-by park_id --hbase-create-table

# 导出部分数据需要Mysql驱动版本较高

# 导出到文件
./bin/sqoop import --connect "jdbc:mysql://rm-aaaaaaaaa.mysql.rds.aliyuncs.com/pp_parking?zeroDateTimeBehavior=convertToNull" --username user1 --password 123123 --query "select * from parking_ths_car_record  where createtime  BETWEEN '2018-09-01 00:00:00' AND '2018-12-26 00:00:00' AND \$CONDITIONS " --split-by id  --target-dir /public/csv/mysql_export


#导出到hbase
./bin/sqoop import --connect "jdbc:mysql://rm-aaasdadadad.mysql.rds.aliyuncs.com/pp_parking?zeroDateTimeBehavior=convertToNull" --username user1 --password 123123 --query "select * from parking_ths_car_record  where createtime  BETWEEN '2018-09-01 00:00:00' AND '2018-12-26 00:00:00' AND \$CONDITIONS " --hbase-table 9_12_parking_ths_car_record --column-family id --hbase-row-key id --split-by park_id --hbase-create-table

```


转载请注明出处，本文采用 [CC4.0](http://creativecommons.org/licenses/by-nc-nd/4.0/) 协议授权
