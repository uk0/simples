---
layout: post
categories: serverless kubeless sls hadoop cdh hdfs
title: 通过Serverless 获取 hadoop 文件系统的[原创]
date: 2018-10-23 18:53:13 +0800
description: 开放文件下载
keywords: serverless 操作 hdfs
---





## 通过Serverless (Kubeless) 进行操作Hdfs , 实现下载文件.





### code

* Python 

```python
from hdfs import Client
from bottle import route,HTTPResponse
client = Client("http://hdfs-web-svc.cloudera:50070", root="/", timeout=100, session=False)

headers = {}
headers[str("content-type")] = 'application/octet-stream'
headers['Content-Disposition'] = 'attachment;filename="data.csv"'
@route('/download')
def getfile():
    with client.read("/public/csv/parking_ths_car_record/part-m-00000") as reader:
        content = reader.read()
        content = content.decode('utf-8', 'ignore').encode('gbk')
        return HTTPResponse(body=content, status=200, headers=headers)

```


*  serverless.yml

```yaml
# Welcome to Serverless!
#
# For full config options, check the kubeless plugin docs:
#    https://github.com/serverless/serverless-kubeless
#
# For documentation on kubeless itself:
#    http://kubeless.io

# Update the service name below with your own service name
service: getfile

# Please ensure the serverless-kubeless provider plugin is installed globally.
# $ npm install -g serverless-kubeless
#
# ...before installing project dependencies to register this provider.
# $ npm install

provider:
  name: kubeless
  namespace: ${env:K8S_NAMESPACE, 'cloudera'}
  runtime: python2.7

plugins:
  - serverless-kubeless

functions:
  getfile:
    handler: handler.getfile

```


* requirements.txt


```txt
bottle==0.12.13
hdfs==2.1.0
```



* kubernetes Ingress 


```yaml
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: getfile-ing
  namespace: cloudera
spec:
  rules:
  - host: getfile.kube.ibm-testing.com
    http:
      paths:
      - backend:
          serviceName: getfile
          servicePort: 8080
```


### 相关地址

> https://github.com/serverless/

> 更多信息可以关注

> https://github.com/kubeless/kubeless/tree/f49533233cd79b116113ea10454cb2ca234dadfd/docker/runtime



转载请注明出处，本文采用 [CC4.0](http://creativecommons.org/licenses/by-nc-nd/4.0/) 协议授权
