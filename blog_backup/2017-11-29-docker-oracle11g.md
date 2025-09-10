---
layout: post
title: docker run Oracle11g。
categories: oracle
description: oracle for container
keywords: container, docker,oracle
---

# oracle running for docker

 * find aggregate

```bash
docker run -d \
  -p 49160:22 \
  -p 49161:1521 \
  -p 8080:8080 \
  -e ORACLE_ENABLE_XDB=true \
  -e ORACLE_ALLOW_REMOTE=true \
  -e ORACLE_PASSWORD_VERIFY=true  \
  wnameless/oracle-xe-11g
```

![](http://112firshme11224.test.upcdn.net/images/oracle.png)

### Installation(with Ubuntu 16.04)
```
docker pull wnameless/oracle-xe-11g
```

Run with 22 and 1521 ports opened:
```
docker run -d -p 49160:22 -p 49161:1521 wnameless/oracle-xe-11g
```

Run this, if you want the database to be connected remotely:
```
docker run -d -p 49160:22 -p 49161:1521 -e ORACLE_ALLOW_REMOTE=true wnameless/oracle-xe-11g
```

By default, the password verification is disable(password never expired). If you want it back, run this:
```
docker run -d -p 49160:22 -p 49161:1521 -e ORACLE_PASSWORD_VERIFY=true wnameless/oracle-xe-11g
```

![](http://112firshme11224.test.upcdn.net/images/oracle2.png)


For performance concern, you may want to disable the disk asynch IO:
```
docker run -d -p 49160:22 -p 49161:1521 -e ORACLE_DISABLE_ASYNCH_IO=true wnameless/oracle-xe-11g
```

For XDB user, run this:
```
docker run -d -p 49160:22 -p 49161:1521 -p 8080:8080 -e ORACLE_ENABLE_XDB=true wnameless/oracle-xe-11g
```

Check if localhost:8080 work
```
curl -XGET "http://localhost:8080"
```
You will show
```
<!DOCTYPE HTML PUBLIC "-//IETF//DTD HTML 2.0//EN">
<HTML><HEAD>
<TITLE>401 Unauthorized</TITLE>
</HEAD><BODY><H1>Unauthorized</H1>
</BODY></HTML>
```

```
# Login http://localhost:8080 with following credential:
username: XDB
password: xdb
```

Connect database with following setting:
```
hostname: localhost
port: 49161
sid: xe
username: system
password: oracle
```

Password for SYS & SYSTEM
```
oracle
```

Login by SSH
```
ssh root@localhost -p 49160
password: admin
```

Support custom DB Initialization
```
# Dockerfile
FROM wnameless/oracle-xe-11g

ADD init.sql /docker-entrypoint-initdb.d/
```


![](http://112firshme11224.test.upcdn.net/images/oracle3.png)