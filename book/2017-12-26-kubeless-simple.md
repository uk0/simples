---
layout: post
categories: kubeless
title: kubeless原创实验
date: 2017-12-26 01:13:11 +0800
description: kubeless 实验性，不涉及到性能测试，只是为了验证功能
keywords: kubeless serverless kubernetes docker
---



# quick start

 * 实验环境kubernetes 1.8
 * kubeless 3
 * docker 17-ce
 * kubeless-ui latest
 

## 部署function

```bash

 kubeless function deploy get-python  --runtime python2.7 \
                                --from-file hellowithdata.py\
                                --handler hellowithdata.handler \
                                --trigger-http

```

## 部署成功


```bash
➜  python git:(master) ✗ kubeless function ls
NAME      	NAMESPACE	HANDLER              	RUNTIME  	TYPE	TOPIC	DEPENDENCIES	STATUS
function  	default  	hello.handler        	python2.7	    	     	            	1/1 READY
get-python	default  	hellowithdata.handler	python2.7	HTTP	     	            	1/1 READY

```

## UI查看


![](http://112firshme11224.test.upcdn.net/demos/33632748-61ce-41ac-9ede-f9b9a0b2fc7f.png)

## 用命令测试


![](http://112firshme11224.test.upcdn.net/demos/ff6ab7bd-2891-428d-b7f3-4975b07ae896.png)


## curl请求

```bash
# 将kubectl 代理到本地端口
kubectl proxy -p 8080 &
```
* 测试

```bash

➜  python git:(master) ✗ curl -L --data '{"Another": "Echo"}' localhost:8080/api/v1/proxy/namespaces/default/services/get-python:function-port/ --header "Content-Type:application/json"
{"Another": "Echo"}%

```

## 解答


```bash
    --from-file # 可执行文件 .py
    --runtime #执行环境
    --handler # function 执行某个方法
    --trigger-http  # function 模式
    --runtime-image # 默认不需要（如果在离线状态可以指定本地已经存在的镜像）
    --trigger-topic # 个人理解是参数存放的topic 比如流处理 （存放到指定队列）
```



转载请注明出处，本文采用 [CC4.0](http://creativecommons.org/licenses/by-nc-nd/4.0/) 协议授权
