---
layout: post
title: ssh链接报错 Unable to negotiate with IP
categories: linux,sshd,ssh
description: matching,ssh
keywords: ssh, matching
---

ssh Unable to negotiate with 182.158.68.5 port 22: no matching host key type found. Their offer: ssh-dss



## 介绍

* 今天帮别人部署东西，远程连接服务器的时候发现发生了错误：。

```bash
 Unable to negotiate with 182.158.68.5 port 22: no matching host key type found. Their offer: ssh-dss
```

## 解决方案


>The version of OpenSSH included in 16.04 disables ssh-dss. There's a neat page with legacy information that includes this issue: http://www.openssh.com/legacy.html
>In a nutshell, you should add the option -oHostKeyAlgorithms=+ssh-dss to the SSH command:

```bash
ssh -oHostKeyAlgorithms=+ssh-dss root@182.158.68.5
```
> 成功登录

```bash
➜  ~ ssh -oHostKeyAlgorithms=+ssh-dss root@182.158.68.5
The authenticity of host '182.158.68.55 (182.158.68.5)' can't be established.
DSA key fingerprint is SHA256:LVWhbgb8q5TdW2QEsFjOVNp8ekW8UqPVE/KCbMOI4CM.
Are you sure you want to continue connecting (yes/no)? 
```

 