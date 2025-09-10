---
layout: post
categories: oracle java docker
title: oracle for docker dev ENV
date: 2018-01-18 00:43:56 +0800
description: oracle 链解决方案。
keywords: oracle docker 
---

oracle 服务采用 Docker 搭建，用的镜像为`wnameless/oracle-xe-11g` 这里负责整理用来测试时候的问题。


## 下载镜像(pull)


```bash
docker pull wnameless/oracle-xe-11g
```

```bash
docker run -d -p 49160:22 -p 49161:1521 -e ORACLE_ALLOW_REMOTE=true wnameless/oracle-xe-11g
```
By default, the password verification is disable(password never expired). If you want it back, run this:

```bash
docker run -d -p 49160:22 -p 49161:1521 -e ORACLE_PASSWORD_VERIFY=true wnameless/oracle-xe-11g
```

For performance concern, you may want to disable the disk asynch IO:

```bash
docker run -d -p 49160:22 -p 49161:1521 -e ORACLE_DISABLE_ASYNCH_IO=true wnameless/oracle-xe-11g
```
For XDB user, run this:

```bash
docker run -d -p 49160:22 -p 49161:1521 -p 8080:8080 -e ORACLE_ENABLE_XDB=true wnameless/oracle-xe-11g
Check if localhost:8080 work
```


```bash
curl -XGET "http://localhost:8080"
```
You will show

<!DOCTYPE HTML PUBLIC "-//IETF//DTD HTML 2.0//EN">
<HTML><HEAD>
<TITLE>401 Unauthorized</TITLE>
</HEAD><BODY><H1>Unauthorized</H1>
</BODY></HTML>

## Login http://localhost:8080 with following credential:
username: XDB
password: xdb
Connect database with following setting:

hostname: localhost
port: 49161
sid: xe
username: system
password: oracle
Password for SYS & SYSTEM

oracle
Login by SSH

ssh root@localhost -p 49160
password: admin
Support custom DB Initialization

## Dockerfile
FROM wnameless/oracle-xe-11g

ADD init.sql /docker-entrypoint-initdb.d/


## SQLDeveloper

```sql
--- 遇到的问题 ora-01045

---ORA-01045: user ... lacks CREATE SESSION privilege; logon denied
--- This error is thrown if a user wants to log on a database but lacks the create session system privilege.
---In order to give someone this privilege, use:

grant create session to aom;
--- See also other Oracle errors.

--- 创建表空间
CREATE   TABLESPACE aom
DATAFILE '/data/disk1/aom_space.dbf' size 500m
AUTOEXTEND ON
NEXT 200M MAXSIZE 20480M
EXTENT MANAGEMENT LOCAL;
--- oracle 创建DB需要用用户去创建，也就是先要创建用户 并且指定表空间
CREATE USER aom IDENTIFIED BY 123456 DEFAULT TABLESPACE aom_space;

---创建用户并指定表空间和临时表空间
CREATE USER aom IDENTIFIED BY 123456
DEFAULT TABLESPACE aom_space
TEMPORARY TABLESPACE aom_temp;

---授权用户
GRANT CREATE USER,DROP USER,ALTER USER ,CREATE ANY VIEW ,
     DROP ANY VIEW,EXP_FULL_DATABASE,IMP_FULL_DATABASE,
     DBA,CONNECT,RESOURCE,CREATE SESSION TO  aom;


---删除表空间
DROP TABLESPACE aom_space INCLUDING CONTENTS AND DATAFILES;


---查看表空间
SELECT tv.TABLESPACE_NAME "TABLESPACE_NAME",TOTALSPACE "TOTALSPACE/M",FREESPACE "FREESPACE/M",ROUND((1-FREESPACE/TOTALSPACE)*100,2) "USED%"
FROM
(SELECT TABLESPACE_NAME,ROUND(SUM(bytes)/1024/1024) TOTALSPACE FROM    DBA_DATA_FILES GROUP BY TABLESPACE_NAME) tv,
         (SELECT TABLESPACE_NAME,ROUND(SUM(bytes)/1024/1024) FREESPACE FROM DBA_FREE_SPACE GROUP BY TABLESPACE_NAME) fs
WHERE tv.TABLESPACE_NAME=fs.TABLESPACE_NAME;


---查看临时表空间
SELECT TABLESPACE_NAME FROM DBA_TABLESPACES;

---增加表空间大小
ALTER TABLESPACE aom_space ADD DATAFILE '/data/disk1/aom_space.dbf' size 4096M;

---增加临时表空间大小
ALTER DATABASE TEMPFILE '/data/disk1/aom_temp.dbf' RESIZE 8192M;

---删除用户
DROP USER aom CASCADE

---把数据导入不同于原系统的表空间，在导入之后却往往发现，数据被导入了原表空间(下面解决此方法)
grant connect, resource,dba to aom;
---回收用户unlimited tablespace权限，这样就可以导入到用户缺省表空间：
revoke unlimited tablespace from aom;
alter user aom quota 0 on aom_space;
alter user aom quota unlimited on aom_space;
```







转载请注明出处，本文采用 [CC4.0](http://creativecommons.org/licenses/by-nc-nd/4.0/) 协议授权
