---
layout: post
categories: linux
title: CentOS6 下升级默认的OpenSSH操作记录（升级到OpenSSH_8.1p1)
date: 2019-12-30 17:46:12 +0800
description: CentOS6.9下升级默认的OpenSSH
keywords: 升级 openssh
---


 升级当前机器的openssh 到最新版本来规避漏洞导致服务器收到攻击。



###  准备和检查

* 检查系统版本
* 检查openssh 版本
* 下载openssh 源码包


```bash
查看操作系统版本
[root@Centos6 ~]# cat /etc/redhat-release       #或者执行"cat /etc/issue"
CentOS release 6.9 (Final)
  
查看默认的OpenSSH版本
[root@Centos6 ~]# ssh -V
OpenSSH_5.3p1, OpenSSL 1.0.1e-fips 11 Feb 2013
 
openssl默认版本如下，这个是关系到下面编译安装openssh的成功与否。
[root@Centos6 ~]# openssl version -a
OpenSSL 1.0.1e-fips 11 Feb 2013
 
备份ssh目录（此步非常重要，一定要记得提前做备份）
[root@Centos6 ~]# cp -rf /etc/ssh /etc/ssh.bak
  
安装telnet，记得提前部署telnet远程登录环境(root用户登录)，避免ssh升级出现问题，导致无法远程管理！
[root@Centos6 ~]# yum install telnet-server
[root@Centos6 ~]# vim /etc/xinetd.d/telnet       
# default: on
# description: The telnet server serves telnet sessions; it uses \
#       unencrypted username/password pairs for authentication.
service telnet
{
        flags           = REUSE
        socket_type     = stream      
        wait            = no
        user            = root
        server          = /usr/sbin/in.telnetd
        log_on_failure  += USERID
        disable         = no                     #将默认的yes修改为no。开启telnet服务功能，否则telnet启动后，23端口就会起不来！
}
  
重启telnet服务
[root@Centos6 ~]# /etc/init.d/xinetd start
Starting xinetd:                                           [  OK  ]
[root@Centos6 ~]# /etc/init.d/xinetd restart
Stopping xinetd:                                           [  OK  ]
Starting xinetd:                                           [  OK  ]
  
telnet默认用于远程登录的端口是23（21是默认的ftp端口、22是默认的ssh端口、23是默认的telnet端口）
[root@Centos6 ~]# lsof -i:23
COMMAND    PID USER   FD   TYPE   DEVICE SIZE/OFF NODE NAME
xinetd    2489 root    5u  IPv6 22131982      0t0  TCP *:telnet (LISTEN)

```

##### 默认情况下，linux不允许root用户以telnet方式登录linux主机，若要允许root用户登录，可采取以下两种方法中的任何一种方法：

```bash
1）第一种方法：修改securetty文件
增加pts配置。如果登录用户较多，需要更多的pts/*。
[root@Centos6 ~]# vim /etc/securetty
......
pts/0
pts/1
pts/2
  
2）第二种方法：移除securetty文件
验证规则设置在/etc/security文件中，该文件定义root用户只能在tty1-tty6的终端上记录，删除该文件或者将其改名即可避开验证规则实现root用户远程登录。
[root@Centos6 ~]# rm -rf /etc/securetty
  
以上两种方法中的任意一种设置后，在客户端使用telnet远程登录目标服务器(使用root用户)都是可以的！

```


#### 升级组件

```bash
[root@Centos6 ~]# yum install -y gcc openssl-devel pam-devel rpm-build tcp_wrappers-devel

# downlaod 
wget http://openbsd.hk/pub/OpenBSD/OpenSSH/portable/openssh-8.1p1.tar.gz

tar -zvxf openssh-8.1p1.tar.gz -C /opt/
cd /opt/openssh-8.1p1
./configure --prefix=/usr --sysconfdir=/etc/ssh --with-pam --with-zlib --with-md5-passwords --with-tcp-wrappers
make && make install


## 安装后提示：

/etc/ssh/ssh_config already exists, install will not overwrite
/etc/ssh/sshd_config already exists, install will not overwrite
/etc/ssh/moduli already exists, install will not overwrite
ssh-keygen: generating new host keys: ECDSA ED25519
/usr/sbin/sshd -t -f /etc/ssh/sshd_config
/etc/ssh/sshd_config line 81: Unsupported option GSSAPIAuthentication
/etc/ssh/sshd_config line 83: Unsupported option GSSAPICleanupCredentials


## 修改配置文件,允许root登录
[root@Centos6 openssh-7.6p1]# sed -i '/^#PermitRootLogin/s/#PermitRootLogin yes/PermitRootLogin yes/' /etc/ssh/sshd_config
[root@Centos6 openssh-7.6p1]# cat /etc/ssh/sshd_config|grep RootLogin
PermitRootLogin yes


#重启OpenSSH
[root@Centos6 openssh-7.6p1]# service sshd restart                 
Stopping sshd:                                             [  OK  ]
Starting sshd: /etc/ssh/sshd_config line 81: Unsupported option GSSAPIAuthentication
/etc/ssh/sshd_config line 83: Unsupported option GSSAPICleanupCredentials
                                                           [  OK  ]
#如上重启OpenSSH出现的告警错误，解决办法如下：
#将/etc/ssh/sshd_config文件中以上行数内容注释下即可
[root@Centos6 openssh-7.6p1]# sed -i '/^GSSAPICleanupCredentials/s/GSSAPICleanupCredentials yes/#GSSAPICleanupCredentials yes/' /etc/ssh/sshd_config
[root@Centos6 openssh-7.6p1]# sed -i '/^GSSAPIAuthentication/s/GSSAPIAuthentication yes/#GSSAPIAuthentication yes/' /etc/ssh/sshd_config
[root@Centos6 openssh-7.6p1]# sed -i '/^GSSAPIAuthentication/s/GSSAPIAuthentication no/#GSSAPIAuthentication no/' /etc/ssh/sshd_config

#次重启OpenSSH服务，则不会出现错误了
[root@Centos6 openssh-7.6p1]# service sshd restart
Stopping sshd:                                             [  OK  ]
Starting sshd:                                             [  OK  ]
# 升级后版本
[root@Centos6 openssh-7.6p1]# ssh -V
OpenSSH_7.6p1, OpenSSL 1.0.1e-fips 11 Feb 2013
```




转载请注明出处，本文采用 [CC4.0](http://creativecommons.org/licenses/by-nc-nd/4.0/) 协议授权
