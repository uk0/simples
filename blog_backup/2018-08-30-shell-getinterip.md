---
layout: post
categories: shell
title: Shell 获得外网IP
date: 2018-08-30 14:33:33 +0800
description: shell 获得服务器的外网地址
keywords: shell
---

阿里云外网服务器如何获得外网IP。



### 看这里 


```bash
#Function: gtIp
#Author: zhangjianxin
#Time: 2018-08-30
#! /bin/bash


bashpath=$(cd `dirname $0`; pwd)
brokerName=dev_rocketmq_4.3.0
#########################setting#######################
URL="http://ip.cn/"
Result=`curl -o /dev/null -s -m 10 --connect-timeout 10 -w %{http_code} $URL`
Test=`echo $Result`
if [ "$Test" == "200" ]
then
	curl -s -o $bashpath/ip.html ${URL}
	echo "search inter ip "
        ip=$(cat ${bashpath}/ip.html  | sed  's/：/ /g' | awk '{print$3}')
	echo $ip
    rm ${bashpath}/ip.html
fi
```



```bash
    awk '!seen[$0]++' <filename> # 去重行
```

转载请注明出处，本文采用 [CC4.0](http://creativecommons.org/licenses/by-nc-nd/4.0/) 协议授权
