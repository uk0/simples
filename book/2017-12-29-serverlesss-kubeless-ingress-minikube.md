---
layout: post
categories: serverless kubernetes ingress kubeless
title: serverless demo for kubernetes + ingress [原创]
date: 2017-12-29 00:11:45 +0800
description: 通过 Serverless 创建 functionkubeless）以及简化开。
keywords: serverless
---

####  文章介绍，Serverless + kubernetes + ingress + kubeless 快速开发function 简化开发，无服务器接口，让开发更快 更爽


## 工具

  * kubernetes
  * kubeless
  * ingress
  * serverless


## 启动Minikube

  * 配置Docker 代理下载镜像速度能快一点。

```bash
minikube start --docker-env HTTP_PROXY=http://192.168.155.2:8118 \
                 --docker-env HTTPS_PROXY=https://192.168.155.2:8118

```
  * 检查kubernetes 是否已经启动完成。

  ![](http://112firshme11224.test.upcdn.net/demos/d8922ed9-aeb0-4056-9abb-3319cf2b7544.png)

  * 检查ingress
     
  ![](http://112firshme11224.test.upcdn.net/demos/51342058-f68a-4c0f-ba9b-9ff846c8a571.png)
   

  * 检查UI是否已经启动
    
  ![](http://112firshme11224.test.upcdn.net/demos/baee88c2-ccfc-4551-9596-6752c4640826.png)

  * 检查kubelessui[有没有都可以]

  ![](http://112firshme11224.test.upcdn.net/demos/4a80959d-8db3-4a44-bc41-403f1c77f324.png)




## 正题

  * 创建一个空的目录一会要用到

    * 安装 serverless  `npm install serverless` 

    * 安装 serverless-kubeless  `npm install serverless-kubeless`

    * 创建一个serverless function `serverless create --template kubeless-python`
    
    * 看看serverless 支持多少模版

    ![](http://112firshme11224.test.upcdn.net/demos/6b52ee20-569a-4450-b3b9-b495c9e3fd0b.png)

* 执行创建命令

```bash
    serverless create --template kubeless-python
```

* 结果 

```bash
➜  demo2 serverless create --template kubeless-python
    Serverless: Generating boilerplate...
    _______                             __
    |   _   .-----.----.--.--.-----.----|  .-----.-----.-----.
    |   |___|  -__|   _|  |  |  -__|   _|  |  -__|__ --|__ --|
    |____   |_____|__|  \___/|_____|__| |__|_____|_____|_____|
    |   |   |             The Serverless Application Framework
    |       |                           serverless.com, v1.25.0
    -------'

    Serverless: Successfully generated boilerplate for template: "kubeless-python"
    Serverless: NOTE: Please update the "service" property in serverless.yml with your service name
```
* 目录介绍

![](http://112firshme11224.test.upcdn.net/demos/0d3ccbe2-a619-44ca-bf70-267a2c6a10db.png)

> 里面会出现四个文件 ，第一个git的忽略文件(不用git可能没有)
> 第二个 是function 的主体
> 第三个是serverless 部署依赖的yml

* handler.py

```python
import json
def hello(request):
    body = {
        "message": "Go Serverless v1.0! Your function executed successfully!",
        "input": request.json
    }
    response = {
        "statusCode": 200,
        "body": json.dumps(body)
    }
    return response
```
* serverless.yml

```yaml
    # Welcome to Serverless!
    #
    # For full config options, check the kubeless plugin docs:
    #    https://github.com/serverless/serverless-kubeless
    #
    # For documentation on kubeless itself:
    #    http://kubeless.io

    # Update the service name below with your own service name
    service: hello-world

    # Please ensure the serverless-kubeless provider plugin is installed globally.
    # $ npm install -g serverless-kubeless
    #
    # ...before installing project dependencies to register this provider.
    # $ npm install

    provider:
    name: kubeless
    runtime: python2.7

    plugins:
    - serverless-kubeless

    functions:
    demo2:  # 有重名 hello 所以改成demo2
        handler: handler.hello

```

* package.json [可以不用管]


* 以上的Serverless 可以运行了

* 开始部署

* 提示 `serverless` 可以简写  `sls`

```bash
    ➜  demo2 serverless deploy
    Serverless: Packaging service...
    Serverless: Excluding development dependencies...
    Serverless: Deploying function demo2...
    Serverless: Function demo2 successfully deployed

```

* 查看是否部署成功

![](http://112firshme11224.test.upcdn.net/demos/c1313fae-c4dc-4b9f-bcb8-1508b8dec8fe.png)

* 创建`Ingress`

```bash
kubeless ingress create ingress-demo2  -n default --function demo2
# 参数介绍
# ingress-demo 是ingress 的名字
# -n 是kubernetes 的命名空间
# --function 是你要绑定到那个 function
```

* 执行命令 没有错误即可，检查 ingress 是否创建

```bash
    kubeless ingress list
```

![](http://112firshme11224.test.upcdn.net/demos/e2839acc-e728-4f81-99e1-97474ec53a40.png)

* 找到我们的链接

```bash
    demo2.192.168.64.2.nip.io
```

* 用postman进行测试。

![](http://112firshme11224.test.upcdn.net/demos/6e04c70c-ad58-46a5-afc7-498b741323f6.png)

  测试通过



## 总结

  * ingress 安装一定要检查好
  * serverless-kubeless 每次创建一个 fucntion 都需要在当前目录执行安装命令`npm install serverless-kubeless`
  * 目前kubeless 支持的语言比较少(runtime) python nodejs
  * 有兴趣可以看看 `aws lambda` 还有 `fnproject`
  * 教程整理比较匆忙如果有疑问请留言。


转载请注明出处，本文采用 [CC4.0](http://creativecommons.org/licenses/by-nc-nd/4.0/) 协议授权
