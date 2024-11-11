---
layout: post
title: Python监控Linux服务器登陆
categories: python
description: Python监控Linux服务器登陆
keywords: linux,python
---


一个用Python做监控Linux服务器登陆

#  先看代码
date：2017-10-17
author：zhangjianxin

[TOC]


```python

#!/usr/bin/env python
#encoding=utf-8

from smtplib import SMTP
import subprocess

smtp = "smtp.qq.com"
user = '1234567'
password = 'xxxx'

run_comd = subprocess.Popen('w¦grep pts',shell=True,stdout=subprocess.PIPE)
data = run_comd.stdout.read()
mailb = ["服务器有新用户登录",data]
mailh = ["From: 1234567@qq.com", "To: xxxx@gmail.com", "Subject: 用户登录监控"]
mailmsg = "\r\n\r\n".join(["\r\n".join(mailh), "\r\n".join(mailb)])

#www.iplaypy.com
def send_mail():
    send = SMTP(smtp)
    send.login(user,password)
    result = send.sendmail("1234567@qq.com", ("xxxx@gmail.com",), mailmsg)
    send.quit()
if data == '':
    pass
else:
    send_mail()

```


* 以上操作经过验证可以直接拿去。
* Owner `breakEval13`
* https://github.com/breakEval13