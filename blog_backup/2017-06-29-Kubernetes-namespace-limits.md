---
layout: post
title: Kubernetes命名空间资源限制[学习]
categories: Java, Scala,Hadoop,Flink
description: Java, Scala,Hadoop,Flink
keywords: Java, Scala,Hadoop,Flink
---

Kubernetes命名空间资源限制主要是为了多个团队一起开发来区分资源使用情况。




## 1. ResourceQuota
```yaml

apiVersion: v1
kind: ResourceQuota
metadata:
  name: compute-resources
  namespace: test
spec:
  hard:
    requests.cpu: "1"
    requests.memory: 2Gi
    limits.cpu: "4"
    limits.memory: 12Gi

```

## 2. LimitRange

```yaml

apiVersion: v1
kind: LimitRange
metadata:
  name: limits
spec:
  limits:
  - max:
      cpu: "4"
      memory: 2Gi
    min:
      cpu: 200m
      memory: 6Mi
    maxLimitRequestRatio:
      cpu: 3
      memory: 2
    type: Pod
  - default:
      cpu: 300m
      memory: 200Mi
    defaultRequest:
      cpu: 200m
      memory: 100Mi
    max:
      cpu: "2"
      memory: 1Gi
    min:
      cpu: 100m
      memory: 3Mi
    maxLimitRequestRatio:
      cpu: 5
      memory: 4
    type: Container


```

## 3.映射方案Ceph RBD

```yaml
apiVersion: v1
kind: ReplicationController
metadata:
  name: demo-app-server
  namespace: [namespace]
spec:
  replicas: 1
  selector:
    app: demo-app-server
  template:
    metadata:
      labels:
        app: demo-app-server
    spec:
      containers:
        - name: demo-app-server-container
          image: image.firsh.io/demo-base/appserver:latest
          imagePullPolicy: Always
          ports:
          - containerPort: 8080
          env:
            - name: APP_HOME_CONF_DIR
              value: /home/demo-app/conf
          command: ["java","-jar","/home/demo-app/demoapp.jar"]
          volumeMounts:
                - name: config-volume
                  mountPath: /home/demo-app/conf
                - name: app-log
                  mountPath: /var/log/development
      volumes:
        - name: config-volume
          configMap:
            name: demo-app-server
        - name: app-log
          hostPath:
            path: /kube/demo-apiserver/logs

---
apiVersion: v1
kind: Service
metadata:
 name: demo-app-service
 namespace: demo
spec:
   #type: ClusterIP
   ports:
   - name: demo-app-out-in
     port: 1000
     targetPort: 8080
   externalIPs:
   - [192.168.1.100]
   selector:
     app: demo-app-server
```

> 感谢企业服务部牛总提供的解决方案
