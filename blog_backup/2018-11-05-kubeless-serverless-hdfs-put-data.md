---
layout: post
categories: hdfs
title: seerverless操作hdfs上传数据。
date: 2018-11-05 20:00:33 +0800
description: serverless 通过 python 操作hdfs
keywords: hdfs
---

这只是一次尝试。。。。

## kubernetes serverless 操作 Hadoop Hdfs 

* serverless Python 操作hdfs

```bash
# /usr/bin/env python
# coding=utf-8

import paramiko
import os
from hdfs import Client
from bottle import route, run
from bottle import request

client = Client("http://hdfs-web-svc.cloudera:50070", root="/", timeout=100, session=False)
# 文件上传的HTML模板，这里没有额外去写html模板了，直接写在这里，方便点吧
@route('/upload')
def upload():
    return '''
        <html>
            <head>
            </head>
            <body>
                <form action"/upload" method="post" enctype="multipart/form-data">
                    <input type="file" name="data" />
                    <input type="submit" value="Upload" />
                </form>
            </body>
        </html>
    '''


# 文件上传，overwrite=True为覆盖原有的文件，
# 如果不加这参数，当服务器已存在同名文件时，将返回“IOError: File exists.”错误
@route('/upload', method='POST')
def doupload():
    upload = request.files.getall('data')
    for meta in upload:
        buf = meta.file.read()
        name, ext = os.path.splitext(meta.filename)
        if ext not in ('.jar', '.tar'):
            return 'File extension not allowed. type [JAR]'
        path = '/task/mr/' + name + ext
        # ssh root@agent-2.agent-svc.cloudera.svc.cluster.local
        print("--------------save Task To" + path)
        with client.write(path, overwrite=True) as writer:
            writer.write(buf);

    return 'ok'
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
service: put-taskjar

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
  upload:
    handler: task.upload

```

* 具体信息看第一条关于 serverless 的文章。


转载请注明出处，本文采用 [CC4.0](http://creativecommons.org/licenses/by-nc-nd/4.0/) 协议授权