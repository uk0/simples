---
layout: post
categories: linux serverspeeder
title: 优秀的VPS TCP加速软件 —— 一键锐速安装脚本（开心版）
date: 2018-05-21 13:49:52 +0800
description: 锐速是一款非常不错的TCP底层加速软件，可以非常方便快速地完成服务器网络的优化，配合 ShadowSocks 效果奇佳。目前锐速官方已经关闭注册了，所以我在这里分享几个开心版（破解）的一键锐速脚本。
keywords: 
---

##  锐速 +   RMTNSS


 * 系统要求Centos7 64 位 内核 3.10
 * RMTNSS https://github.com/uk0/RMTNSS
 * 锐速 https://www.91yun.co/serverspeeder91yun




## quick start


```bash
#检查系统虚拟化支持
yum install virt-what -y
# 例如返回 kvm
virt-what
# 安装新的内核
wget --no-check-certificate -O rskernel.sh https://raw.githubusercontent.com/uxh/shadowsocks_bash/master/rskernel.sh && bash rskernel.sh

# 安装锐速

wget -N --no-check-certificate https://github.com/91yun/serverspeeder/raw/master/serverspeeder-v.sh && bash serverspeeder-v.sh CentOS 7.1 3.10.0-229.1.2.el7.x86_64 x64 3.11.20.4 serverspeeder_3283

#重启锐速
/serverspeeder/bin/serverSpeeder.sh restart
#启动锐速
/serverspeeder/bin/serverSpeeder.sh start
#停止锐速
/serverspeeder/bin/serverSpeeder.sh stop
#查看锐速运行情况
/serverspeeder/bin/serverSpeeder.sh status
# 卸载
chattr -i /serverspeeder/etc/apx* && /serverspeeder/bin/serverSpeeder.sh uninstall -f

```




转载请注明出处，本文采用 [CC4.0](http://creativecommons.org/licenses/by-nc-nd/4.0/) 协议授权
