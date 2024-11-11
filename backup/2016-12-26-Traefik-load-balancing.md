---
layout: post
title: traefik 负载均衡
categories: docker
description: traefik 负载均衡
keywords: docker
---

> Traefik 是一款开源的反向代理与负载均衡工具。它最大的优点是能够与常见的微服务系统直接整合，可以实现自动化动态配置。目前支持Docker, Swarm, Mesos/Marathon, Mesos, Kubernetes, Consul, Etcd, Zookeeper, BoltDB, Rest API等等后端模型

# Traefik 简介

traefik 官网 https://docs.traefik.io/  
github   https://github.com/containous/traefik
   

# traefik 与 kubernetes


## 1.1 kubernetes 环境

```
[root@k8s-node-1 ~]#kubectl get nodes
NAME         STATUS         AGE
k8s-node-1   Ready,master   5d
k8s-node-2   Ready          5d
k8s-node-3   Ready          5d

```


## 1.2 下载 traefik-dm

```
# 定义 node-2 node-3 标签。

kubectl label nodes k8s-node-2 role=traefik-lb
kubectl label nodes k8s-node-3 role=traefik-lb

# 查看此标签的 nodes

[root@k8s-node-1 ~]#kubectl get nodes -l role=traefik-lb
NAME         STATUS    AGE
k8s-node-2   Ready     5d
k8s-node-3   Ready     5d




# 下载yaml 文件

curl -O https://raw.githubusercontent.com/containous/traefik/master/examples/k8s/traefik.yaml
		
		
# yaml 可修改 replicas = node 个数
#      增加 restartPolicy: Always
#	   修改 web ui 端口 containerPort: 8080 修改为 containerPort: 8888
#	                    args: 下面增加 --web.address=:8888
#	   增加 nodeSelector:
#              role: traefik-lb
			   
# 修改后的文件如下：

apiVersion: v1
kind: Deployment
apiVersion: extensions/v1beta1
metadata:
  name: traefik-ingress-controller
  namespace: kube-system
  labels:
    k8s-app: traefik-ingress-lb
spec:
  replicas: 2
  selector:
    matchLabels:
      k8s-app: traefik-ingress-lb
  template:
    metadata:
      labels:
        k8s-app: traefik-ingress-lb
        name: traefik-ingress-lb
    spec:
      terminationGracePeriodSeconds: 60
      hostNetwork: true
	  restartPolicy: Always
      nodeSelector:
        role: traefik-lb
      containers:
      - image: traefik
        name: traefik-ingress-lb
        resources:
          limits:
            cpu: 200m
            memory: 30Mi
          requests:
            cpu: 100m
            memory: 20Mi
        ports:
        - name: http
          containerPort: 80
          hostPort: 80
        - name: admin
          containerPort: 8888
        args:
        - --web
        - --web.address=:8888
        - --kubernetes

		
# 创建 traefik deployment

[root@k8s-node-1 ~]#kubectl create -f traefik.yaml 
deployment "traefik-ingress-controller" created


# 查看 deployment

[root@k8s-node-1 ~]#kubectl get deploy --all-namespaces
NAMESPACE     NAME                         DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
kube-system   traefik-ingress-controller   1         1         1            1           2m

```


## 1.3 部署一个 nginx 应用

```
# 创建一个 nginx yaml 文件

apiVersion: extensions/v1beta1 
kind: Deployment 
metadata: 
  name: nginx-dm
spec: 
  replicas: 2
  template: 
    metadata: 
      labels: 
        name: nginx 
    spec: 
      containers: 
        - name: nginx 
          image: nginx 
          imagePullPolicy: IfNotPresent
          ports: 
            - containerPort: 80

---
apiVersion: v1 
kind: Service
metadata: 
  name: nginx-dm 
spec: 
  ports: 
    - port: 80
      targetPort: 80
      protocol: TCP 
  selector: 
    name: nginx
	

# 创建 deploy 与 svc

[root@k8s-node-1 ~]#kubectl create -f nginx.yaml 



# 查看 nginx-dm

[root@k8s-node-1 ~]#kubectl get svc
NAME         CLUSTER-IP      EXTERNAL-IP   PORT(S)   AGE
kubernetes   10.96.0.1       <none>        443/TCP   5d
nginx-dm     10.105.76.250   <none>        80/TCP    34m
[root@k8s-node-1 ~]#kubectl get deployment
NAME       DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
nginx-dm   2         2         2            2           35m
	
```


## 1.4 创建 traefik ingress

```
# 创建 traefik-ingress yaml 文件


apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: traefik-ingress
spec:
  rules:
  - host: traefik.test.jicki.me
    http:
      paths:
      - path: /
        backend:
          serviceName: nginx-dm
          servicePort: 80


# 创建 ingress	  
[root@k8s-node-1 ~]#kubectl create -f traefik-ingress.yaml


# 查看 ing

[root@k8s-node-1 ~]#kubectl get ing
NAME              HOSTS                   ADDRESS   PORTS     AGE
traefik-ingress   traefik.test.jicki.me             80        32m	
	
```


## 1.5 创建 traefik web ui 

```
# 下载 web ui yaml 文件

curl -O https://raw.githubusercontent.com/containous/traefik/master/examples/k8s/ui.yaml

[root@k8s-node-1 ~]#kubectl create -f ui.yaml

```


## 1.6 访问 nginx 应用


```
[root@k8s-node-1 ~]#curl -I traefik.test.jicki.me
HTTP/1.1 200 OK
Accept-Ranges: bytes
Content-Length: 612
Content-Type: text/html
Date: Mon, 26 Dec 2016 09:18:33 GMT
Etag: "585011f4-264"
Last-Modified: Tue, 13 Dec 2016 15:21:24 GMT
Server: nginx/1.11.7

```



## 1.7 访问 traefik web ui

```
http://10.6.0.187:8888
http://10.6.0.188:8888
```



## 1.8 心跳检测

Null