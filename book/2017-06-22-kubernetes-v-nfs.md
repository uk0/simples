---
layout: post
title: Kubernetes部分Volume类型介绍及yaml示例
categories: linux,Volume,yaml
description: matching,ssh
keywords: kubernetes, Volume,yaml
---

Kubernetes部分Volume类型介绍及yaml示例



## 1、EmptyDir（本地数据卷

　EmptyDir类型的volume创建于pod被调度到某个宿主机上的时候，而同一个pod内的容器都能读写EmptyDir中的同一个文件。一旦这个pod离开了这个宿主机，EmptyDirr中的数据就会被永久删除。所以目前EmptyDir类型的volume主要用作临时空间，比如Web服务器写日志或者tmp文件需要的临时目录。yaml示例如下：
```bash
[root@k8s-master demon2]# cat test-emptypath.yaml 

```

```yaml
apiVersion: v1
kind: Pod
metadata:
  labels:
    name: test-emptypath
    role: master
  name: test-emptypath
spec:
  containers:
    - name: test-emptypath
      image: registry:5000/back_demon:1.0
      volumeMounts:
       - name: log-storage
         mountPath: /home/laizy/test/
      command:
      - /run.sh
  volumes:
  - name: log-storage
    emptyDir: {}
```

## 2、HostDir（本地数据卷
　HostDir属性的volume使得对应的容器能够访问当前宿主机上的指定目录。例如，需要运行一个访问Docker系统目录的容器，那么就使用/var/lib/docker目录作为一个HostDir类型的volume；或者要在一个容器内部运行CAdvisor，那么就使用/dev/cgroups目录作为一个HostDir类型的volume。一旦这个pod离开了这个宿主机，HostDir中的数据虽然不会被永久删除，但数据也不会随pod迁移到其他宿主机上。因此，需要注意的是，由于各个宿主机上的文件系统结构和内容并不一定完全相同，所以相同pod的HostDir可能会在不同的宿主机上表现出不同的行为。yaml示例如下：

```bash
[root@k8s-master demon2]# cat test-hostpath.yaml 
```


```yaml
apiVersion: v1
kind: Pod
metadata:
  labels:
    name: test-hostpath
    role: master
  name: test-hostpath
spec:
  containers:
    - name: test-hostpath
      image: registry:5000/back_demon:1.0
      volumeMounts:
       - name: ssl-certs
         mountPath: /home/laizy/test/cert
         readOnly: true
      command:
      - /run.sh
  volumes:
  - name: ssl-certs
    hostPath:
     path: /etc/ssl/certs
```
## 3、NFS（网络数据卷）


NFS类型的volume。允许一块现有的网络硬盘在同一个pod内的容器间共享。yaml示例如下：

```bash
[root@k8s-master demon2]# cat test-nfspath.yaml 
```


```yaml
apiVersion: v1
kind: Pod
metadata:
  labels:
    name: test-nfspath
    role: master
  name: test-nfspath
spec:
  containers:
    - name: test-nfspath
      image: registry:5000/back_demon:1.0
      volumeMounts:
       - name: nfs-storage
         mountPath: /home/laizy/test/
      command:
      - /run.sh
  volumes:
  - name: nfs-storage
    nfs:
     server: 192.168.20.47
     path: "/data/disk1"
```
## 4、Secret（信息数据卷）

　Kubemetes提供了Secret来处理敏感数据，比如密码、Token和密钥，相比于直接将敏感数据配置在Pod的定义或者镜像中，Secret提供了更加安全的机制（Base64加密），防止数据泄露。Secret的创建是独立于Pod的，以数据卷的形式挂载到Pod中，Secret的数据将以文件的形式保存，容器通过读取文件可以获取需要的数据。yaml示例如下：

```bash
[root@k8s-master demon2]# cat secret.yaml 
```
```yaml
apiVersion: v1
kind: Secret
metadata:
 name: mysecret
type: Opaque
data:
 username: emhlbnl1
 password: eWFvZGlkaWFv
 ```
 ```bash
[root@k8s-master demon2]# cat test-secret.yaml 
```
```yaml
apiVersion: v1
kind: Pod
metadata:
  labels:
    name: test-secret
    role: master
  name: test-secret
spec:
  containers:
    - name: test-secret
      image: registry:5000/back_demon:1.0
      volumeMounts:
       - name: secret
         mountPath: /home/laizy/secret
         readOnly: true
      command:
      - /run.sh
  volumes:
  - name: secret
    secret:
     secretName: mysecret
```

 