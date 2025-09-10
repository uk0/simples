---
layout: post
categories: docker swarm spring
title: docker swarm 集成 Spring Cloud 验证
date: 2018-04-23 10:18:14 +0800
description: docekr swarm spring cloud 验证
keywords: swarm
---

## 开始

 * 机器选择 Centos 7 4Gb内存
 * systemctl stop firewalld
 * systemctl disable firewalld
 * yum install -y docker-io
 * yum install -y vim
 * systemctl start docker
 * systemctl enable docker


## docker swarm init

 * 任意一台机器执行 `docker swarm init `
 * 讲得到的 结果 `copy`
 * 在其他两台机器上执行
 * 配置 docker swarm 的启动文件
 

```yaml
version: '3'
services:
  eureka1:
    image: registry.ap-northeast-1.aliyuncs.com/emos_prod/emos.register.test:latest
    networks:
      - springcloud-overlay
    ports:
      - "8098:8098"
    environment:
      - ADDITIONAL_EUREKA_SERVER_LIST=http://emos:emos@eureka2:8098/eureka/,http://emos:emos@eureka3:8098/eureka/
    tty: true
    command: ["bash","-i","bin/Entrypoint.sh","start"]
  eureka2:
    image: registry.ap-northeast-1.aliyuncs.com/emos_prod/emos.register.test:latest
    networks:
      - springcloud-overlay
    ports:
      - "8097:8098"
    environment:
      - ADDITIONAL_EUREKA_SERVER_LIST=http://emos:emos@eureka1:8098/eureka/,http://emos:emos@eureka3:8098/eureka/
    tty: true
    command: ["bash","-i","bin/Entrypoint.sh","start"]
  eureka3:
    image: registry.ap-northeast-1.aliyuncs.com/emos_prod/emos.register.test:latest
    networks:
      - springcloud-overlay
    ports:
      - "8096:8098"
    environment:
      - ADDITIONAL_EUREKA_SERVER_LIST=http://emos:emos@eureka1:8098/eureka/,http://emos:emos@eureka2:8098/eureka/
    tty: true
    command: ["bash","-i","bin/Entrypoint.sh","start"]


  emos-config:
    image: registry.ap-northeast-1.aliyuncs.com/emos_prod/emos.config.test:latest
    ports:
      - "8888"
    networks:
      - springcloud-overlay
    environment:
      - EUREKA_SERVER_ADDRESS=http://emos:emos@eureka1:8098/eureka/,http://emos:emos@eureka2:8098/eureka/,http://emos:emos@eureka3:8098/eureka/
    tty: true
    depends_on:
      - eureka1
      - eureka2
      - eureka3
    command: ["bash","-i","bin/Entrypoint.sh","start"]
  ths-account:
    image: registry.ap-northeast-1.aliyuncs.com/emos_prod/emos.account.test:latest
    ports:
      - "8006"
    networks:
      - springcloud-overlay
    depends_on:
      - eureka1
      - eureka2
      - eureka3
      - emos-config
    environment:
      - EUREKA_SERVER_ADDRESS=http://emos:emos@eureka1:8098/eureka/,http://emos:emos@eureka2:8098/eureka/,http://emos:emos@eureka3:8098/eureka/
    tty: true
    command: ["bash","-i","bin/Entrypoint.sh","start"]
  emos-gateway:
    image: registry.ap-northeast-1.aliyuncs.com/emos_prod/emos.gateway.test:latest
    ports:
      - "8002"
    networks:
      - springcloud-overlay
    depends_on:
      - eureka1
      - eureka2
      - eureka3
      - ths-account
      - emos-config
    environment:
      - EUREKA_SERVER_ADDRESS=http://emos:emos@eureka1:8098/eureka/,http://emos:emos@eureka2:8098/eureka/,http://emos:emos@eureka3:8098/eureka/
    tty: true
    command: ["bash","-i","bin/Entrypoint.sh","start"]
networks:
  springcloud-overlay:
    driver: overlay
```

## 注意启动顺序  

```bash
#可以拆开分开执行
docker stack deploy -c  emeos.yml demo
```


### 问题解析

* 其中spring-clou

```yaml

spring:
  cloud:
    config:
      profile: ${config.profile:dev}
      discovery:
        enabled: true
        service-id: ${APPLICATION_CONFIG_NAME}
  application:
      name: ${APPLICATION_NAME}
security:
  basic:
    enabled: false
  user:
    name: emos #eureka 用户名
    password: emos #eureka 密码
eureka:
  client:
    serviceUrl:
      defaultZone: ${EUREKA_SERVER_ADDRESS}
  instance:
    status-page-url-path: /swagger-ui.html
    hostname: ${HOST_NAME}
management:
  security:
    enabled: false

```
* `HOST_NAME` 是要用 `swarm` 启动的服务名来定义的例如：`ths-account`
* `EUREKA_SERVER_ADDRESS` 是高可用服务
* `APPLICATION_NAME` 是当前的 `APPLICATION`的名字需要进行配置
* `APPLICATION_CONFIG_NAME` 是`config`服务的`application`名称 有些服务需要配置此项
* `以上可能不善于表达，仅仅是技术方向的实现验证时间：2018年4月23日`

转载请注明出处，本文采用 [CC4.0](http://creativecommons.org/licenses/by-nc-nd/4.0/) 协议授权
