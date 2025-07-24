---
layout: post
categories: tunnelbroker
title: tunnelbroker ipv6 pool
date: 2025-07-28 23:48:46 +0800
description: 使用 tunnelbroker 构建自己的IPv6 代理池
keywords: tunnelbroker 阿里云 
---



#### 准备


* 阿里云轻量服务器 【Centos 8.2】
* tunnelbroker 账号
* tunnelbroker 通道



#### 脚本如下


* sudo vim /etc/sysctl.conf 编辑系统设置ipv6相关
    
    ```bash
    
    net.ipv6.conf.all.disable_ipv6 = 0
    net.ipv6.conf.default.disable_ipv6 = 0
    net.ipv6.conf.lo.disable_ipv6 = 0
    net.ipv6.conf.eth0.disable_ipv6 = 0
    
    
    ```
  
* sudo sysctl -p 重载系统设置

* 获取Example Config

    > 2001:xxx:x:xxx::2/64 = [Client IPv6 Address]
    > 172.xxx.xxx.xx = [服务器内网IP]

* 服务器上执行
    
  ```bash
    
  sudo modprobe ipv6
  sudo ip tunnel add he-ipv6 mode sit remote 216.66.22.2 local 172.xxx.xxx.xx ttl 255
  sudo ip link set he-ipv6 up
  sudo ip addr add 2001:xxx:x:xxx::2/64 dev he-ipv6
  sudo ip route add ::/0 dev he-ipv6
  sudo ip -f inet6 addr
    
  ```

* 验证 

    ```bash
    
    curl ipv6.ip.sb
    
    ```
  
* 收尾

```bash
sudo sysctl -w net.ipv6.ip_nonlocal_bind=1 #开启不限制绑定

# 这个地址在配置里面可以找到 【Routed IPv6 Prefixes】
sudo ip -6 route add local 2001:xxx:x:xxx::/64 dev lo #添加本地回环接口
# 测试
curl --interface 2001:xxx:x:xxx::3 ipv6.ip.sb
```


* 开启代理


```bash
#amd 

wget https://github.com/deanxv/go-proxy-ipv6-pool/releases/download/v1.0.0/go-proxy-ipv6-pool-linux-amd64

#arm

wget https://github.com/deanxv/go-proxy-ipv6-pool/releases/download/v1.0.0/go-proxy-ipv6-pool-linux-arm64


chmod +x go-proxy-ipv6-pool-linux-amd64

# cidr = 【Routed IPv6 Prefixes】
nohup ./go-proxy-ipv6-pool-linux-amd64 --port 51422 --cidr 2001:xxx:x:xxx::/64 > proxy.log 2>&1 &

curl -x http://xx.xx.xx.xx:51422 http://6.ipw.cn/

# 会收到变化的ipv6地址
```

