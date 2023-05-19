---
layout: post
title: kubernetes Ceph RBD
categories: kubernetes
description: kubernetes Ceph RBD
keywords: kubernetes
---

> 使用 Ceph RBD 做为 Kubernetes 后端存储

# Ceph 安装部署

> 由于这里我们使用 RBD 所以我们使用到的组件为 Ceph.mon, Ceph.osd, 这两个组件就可以了。 Ceph.mds 为 cephfs 所需组件


## 部署环境

```
# Ceph.Mon = 2n+1 个 ，3个的情况下只能掉线一个，如果同时2个掉线
# 集群会出现无法仲裁，集群会一直等待 Ceph.Mon 恢复超过半数。

172.16.1.37  Ceph.admin + Ceph.Mon
172.16.1.38  Ceph.Mon + Ceph.osd
172.16.1.39  Ceph.Mon + Ceph.osd
172.16.1.40  Ceph.osd

```


## 初始化环境

```
# 1. ceph 对时间要求很严格， 一定要同步所有的服务器时间

# 2. 配置 hosts

cat << "EOF" >> /etc/hosts
# --- Ceph Hosts -----
172.16.1.37  ceph-node-1
172.16.1.38  ceph-node-2
172.16.1.39  ceph-node-3
172.16.1.40  ceph-node-4
# --- Ceph Hosts -----
EOF


# 3. 配置 hostname

hostnamectl --static set-hostname ceph-node-1
hostnamectl --static set-hostname ceph-node-2
hostnamectl --static set-hostname ceph-node-3
hostnamectl --static set-hostname ceph-node-4


# 4. Ceph.admin 节点 配置无密码ssh

ssh-keygen
ssh-copy-id ceph-node-1 -p33
ssh-copy-id ceph-node-2 -p33
ssh-copy-id ceph-node-3 -p33
ssh-copy-id ceph-node-4 -p33

Now try logging into the machine, with:   "ssh -p '33' 'ceph-node-2'"



# 5. 增加 ~/.ssh/config 文件, 否则 ceph-deploy 创建集群 报端口错误
[root@ceph-node-1 ~]# vi ~/.ssh/config

Host ceph-node-1
    Hostname ceph-node-1
    Port 33
Host ceph-node-2
    Hostname ceph-node-2
    Port 33
Host ceph-node-3
    Hostname ceph-node-3
    Port 33
Host ceph-node-4
    Hostname ceph-node-4
    Port 33
```


## 安装 Ceph-deploy

```
# 管理节点 安装 ceph-deploy 管理工具
# 配置 官方 的 Ceph 源

rpm --import https://download.ceph.com/keys/release.asc
rpm -Uvh --replacepkgs https://download.ceph.com/rpm-jewel/el7/noarch/ceph-release-1-0.el7.noarch.rpm

# 安装 epel 源
rpm -Uvh http://mirrors.ustc.edu.cn/centos/7/extras/x86_64/Packages/epel-release-7-9.noarch.rpm


[root@ceph-node-1 ~]# yum makecache

[root@ceph-node-1 ~]# yum -y install ceph-deploy ntpdate

```

## 创建 Ceph-Mon

```
# 创建集群目录，用于存放配置文件，证书等信息

mkdir -p /opt/ceph-cluster

cd /opt/ceph-cluster/

# 创建ceph-mon 节点
[root@ceph-node-1 ceph-cluster]# ceph-deploy new ceph-node-1 ceph-node-2 ceph-node-3


Monitor initial members are ['ceph-node-1', 'ceph-node-2', 'ceph-node-3']


# 查看配置文件

[root@ceph-node-1 ceph-cluster]# cat ceph.conf 
[global]
fsid = 211a5d26-0377-4d1c-8555-c4793cef83c1
mon_initial_members = ceph-node-1, ceph-node-2, ceph-node-3
mon_host = 172.16.1.37,172.16.1.38,172.16.1.39
auth_cluster_required = cephx
auth_service_required = cephx
auth_client_required = cephx


# 修改 osd 的副本数，既数据保存N份。

echo 'osd_pool_default_size = 2' >> ./ceph.conf


# 注: 如果文件系统为 ext4 请添加

echo 'osd max object name len = 256' >> ./ceph.conf
echo 'osd max object namespace len = 64' >> ./ceph.conf

```

## 安装 Ceph

```
# 可通过 ceph-deploy 安装，也可以登陆node 本地安装
# 命令: ceph-deploy install {ceph-node}[{ceph-node} ...]

[root@ceph-node-1 ceph-cluster]# ceph-deploy install ceph-node-1 ceph-node-2 ceph-node-3 ceph-node-4

# 使用 cpeh-deploy 安装，会一直使用官方源下载，超级慢
# 如果出现超时，请自行在每个服务器中使用 yum 安装

# 替换 ceph 源 为 163 源
sed -i 's/download\.ceph\.com/mirrors\.163\.com\/ceph/g' /etc/yum.repos.d/ceph.repo


yum -y install ceph ceph-radosgw


# 检测 安装

[root@ceph-node-1 ceph-cluster]# ceph --version
ceph version 10.2.7 (50e863e0f4bc8f4b9e31156de690d765af245185)
```


## 初始化 ceph-mon 节点

```
# 请务必在 ceph-cluster 目录下

[root@ceph-node-1 ceph-cluster]# ceph-deploy mon create-initial



[ceph_deploy.gatherkeys][INFO  ] Storing ceph.client.admin.keyring
[ceph_deploy.gatherkeys][INFO  ] Storing ceph.bootstrap-mds.keyring
[ceph_deploy.gatherkeys][INFO  ] keyring 'ceph.mon.keyring' already exists
[ceph_deploy.gatherkeys][INFO  ] Storing ceph.bootstrap-osd.keyring
[ceph_deploy.gatherkeys][INFO  ] Storing ceph.bootstrap-rgw.keyring
[ceph_deploy.gatherkeys][INFO  ] Destroy temp directory /tmp/tmpNfXVDA

```

## 初始化 ceph.osd 节点

```
# 首先创建 存储空间, 如果使用分区，可略过

[root@ceph-node-2 ~]# mkdir -p /opt/ceph-osd-1 && chown ceph:ceph /opt/ceph-osd-1
[root@ceph-node-3 ~]# mkdir -p /opt/ceph-osd-2 && chown ceph:ceph /opt/ceph-osd-2
[root@ceph-node-4 ~]# mkdir -p /opt/ceph-osd-3 && chown ceph:ceph /opt/ceph-osd-3




# 启动 osd

[root@ceph-node-1 ceph-cluster]# ceph-deploy osd prepare ceph-node-2:/opt/ceph-osd-1 ceph-node-3:/opt/ceph-osd-2 ceph-node-4:/opt/ceph-osd-3

```

```
# 激活 所有 osd 节点

[root@ceph-node-1 ceph-cluster]# ceph-deploy osd activate ceph-node-2:/opt/ceph-osd-1 ceph-node-3:/opt/ceph-osd-2 ceph-node-4:/opt/ceph-osd-3
```


```
# 把管理节点的配置文件与keyring同步至其它节点
[root@ceph-node-1 ceph-cluster]# ceph-deploy admin ceph-node-1 ceph-node-2 ceph-node-3 ceph-node-4
```


```
# 没有出现 红色的 ERROR 
# 查看状态

[root@ceph-node-1 ceph-cluster]# ceph -s
    cluster 7ab7e076-24e8-4236-bd4e-88675aab7365
     health HEALTH_OK
     monmap e1: 3 mons at {ceph-node-1=172.16.1.37:6789/0,ceph-node-2=172.16.1.38:6789/0,ceph-node-3=172.16.1.39:6789/0}
            election epoch 6, quorum 0,1,2 ceph-node-1,ceph-node-2,ceph-node-3
     osdmap e15: 3 osds: 3 up, 3 in
            flags sortbitwise,require_jewel_osds
      pgmap v24: 64 pgs, 1 pools, 0 bytes data, 0 objects
            15628 MB used, 8957 GB / 9452 GB avail
                  64 active+clean
```

```
# 查看 osd tree

[root@ceph-node-1 ceph-cluster]# ceph osd tree
ID WEIGHT  TYPE NAME            UP/DOWN REWEIGHT PRIMARY-AFFINITY 
-1 9.23126 root default                                           
-2 3.07709     host ceph-node-2                                   
 0 3.07709         osd.0             up  1.00000          1.00000 
-3 3.07709     host ceph-node-3                                   
 1 3.07709         osd.1             up  1.00000          1.00000 
-4 3.07709     host ceph-node-4                                   
 2 3.07709         osd.2             up  1.00000          1.00000 

```

```
# 设置所有开机启动

systemctl enable ceph-osd.target
systemctl enable ceph-mon.target


```



```
# 重启系统以后重启 osd

#查看处于down 的 osd
ceph osd tree

# 登陆所在node 启动ceph-osd进程 [id 为 tree 查看的 id]

systemctl start ceph-osd@id
```



# Kubernetes Volume Ceph RBD

> 官方 RDB 的文件 https://github.com/kubernetes/kubernetes/tree/master/examples/volumes/rbd

## k8s 安装 ceph-common

```
# 1. 在所有的 k8s node 里面安装 ceph-common
yum -y install ceph-common

# 2. 拷贝 ceph.conf 与 ceph.client.admin.keyring
拷贝到 /etc/ceph/ 目录下

# 3. 配置 kubelet 增加ceph参数
# 这里最好是在创建 k8s 集群的时候就添加好

vim /usr/local/bin/kubelet

# 增加 如下

  -v /etc/ceph:/etc/ceph:ro \
  -v /sbin/modprobe:/sbin/modprobe:ro \
  -v /usr/sbin/modprobe:/usr/sbin/modprobe:ro \
  -v /lib/modules:/lib/modules:ro \

# 重启 kubelet

systemctl restart kubelet.service

# 否则创建 pod 的时候 报错
with: rbd: failed to modprobe rbd error:exit status 1

```




## 修改 ceph 配置

```
# 修改配置
# rbd image有4个 features，layering, exclusive-lock, object-map, fast-diff, deep-flatten
因为目前内核仅支持layering，修改默认配置。

[root@ceph-node-1 ceph-cluster]# echo 'rbd_default_features = 1' >> ./ceph.conf

# 验证一下
[root@ceph-node-1 ceph-cluster]# ceph --show-config|grep rbd|grep features
rbd_default_features = 1

# 把管理节点的配置文件与keyring同步至其它节点
[root@ceph-node-1 ceph-cluster]# ceph-deploy --overwrite-conf admin ceph-node-1 ceph-node-2 ceph-node-3 ceph-node-4
```

## 创建 ceph-secret

```
# 获取 client.admin 的值
[root@ceph-node-1 ceph-cluster]# ceph auth get-key client.admin
AQBpUAxZGmDBBxAAVbgeRss9jv39dE0biTE7qQ==

# 转换成 base64 编码
[root@ceph-node-1 ceph-cluster]# echo "AQBpUAxZGmDBBxAAVbgeRss9jv39dE0biTE7qQ=="|base64
QVFCcFVBeFpHbURCQnhBQVZiZ2VSc3M5anYzOWRFMGJpVEU3cVE9PQo=

# 创建ceph-secret.yaml文件

apiVersion: v1
kind: Secret
metadata:
  name: ceph-secret
data:
  key: QVFCcFVBeFpHbURCQnhBQVZiZ2VSc3M5anYzOWRFMGJpVEU3cVE9PQo=


# 导入 yaml 文件
[root@k8s-node-28 cephrbd]# kubectl apply -f ceph-secret.yaml
secret "ceph-secret" created

# 查看 状态
[root@k8s-node-28 cephrbd]# kubectl get secret |grep ceph
ceph-secret           Opaque                                1         53s
```

## 创建 ceph  image

```
# 创建一个 1G 的 image
[root@ceph-node-1 ceph-cluster]# rbd create test-image -s 1G

# 查看 image
[root@ceph-node-1 ceph-cluster]# rbd list
test-image

[root@ceph-node-1 ceph-cluster]# rbd info test-image
rbd image 'test-image':
        size 1024 MB in 256 objects
        order 22 (4096 kB objects)
        block_name_prefix: rbd_data.373b2ae8944a
        format: 2
        features: layering
        flags:

# 这里 format 为 2 , 如果旧系统 不支持 format 2 ，可将 format 设置为 1 。
```

## 创建一个 pv 

```
# 创建 test-pv.yaml

apiVersion: v1
kind: PersistentVolume
metadata:
  name: test-pv
spec:
  capacity:
    storage: 1Gi
  accessModes:
    - ReadWriteOnce
  rbd:
    monitors:
      - 172.16.1.37:6789
      - 172.16.1.38:6789
      - 172.16.1.38:6789
    pool: rbd
    image: test-image
    user: admin
    secretRef:
      name: ceph-secret
    fsType: ext4
    readOnly: false
  persistentVolumeReclaimPolicy: Recycle



# 导入 yaml 文件

[root@k8s-node-28 cephrbd]# kubectl apply -f test-pv.yaml 
persistentvolume "test-pv" created

# 查看 pv
[root@k8s-node-28 cephrbd]# kubectl get pv |grep test-pv
test-pv   1Gi        RWO           Recycle         Available                                      33s
```

## 创建一个 pvc

```
# vi test-pvc.yaml

kind: PersistentVolumeClaim
apiVersion: v1
metadata:
  name: test-claim
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 1Gi


# 导入 yaml 文件

[root@k8s-node-28 cephrbd]# kubectl apply -f test-pvc.yaml 
persistentvolumeclaim "test-claim" created


# 查看 pvc

[root@k8s-node-28 cephrbd]# kubectl get pvc |grep test
test-claim   Bound     test-pv   1Gi        RWO                          31s
```

## 创建一个 deployment

```
vi nginx-deplyment.yaml


apiVersion: extensions/v1beta1 
kind: Deployment 
metadata: 
  name: nginx-dm
spec: 
  replicas: 1
  template: 
    metadata: 
      labels: 
        name: nginx 
    spec: 
      containers: 
        - name: nginx 
          image: nginx:alpine 
          imagePullPolicy: IfNotPresent
          ports: 
            - containerPort: 80
          volumeMounts:
            - name: ceph-rbd-volume
              mountPath: "/usr/share/nginx/html"
      volumes:
      - name: ceph-rbd-volume
        persistentVolumeClaim:
          claimName: test-claim


# 导入 yaml 文件

[root@k8s-node-28 cephrbd]# kubectl apply -f nginx-deplyment.yaml 
deployment "nginx-dm" configured


# 查看 pods

[root@k8s-node-28 cephrbd]# kubectl get pods
NAME                        READY     STATUS    RESTARTS   AGE
nginx-dm-1697116629-j4m1g   1/1       Running   0          8m

```

## 测试

```
# 查看 分区文件

[root@k8s-node-28 cephrbd]# kubectl exec -it nginx-dm-1697116629-j4m1g -- df -h
Filesystem                Size      Used Available Use% Mounted on
overlay                 115.2G      4.9G    104.4G   4% /
tmpfs                    62.8G         0     62.8G   0% /dev
tmpfs                    62.8G         0     62.8G   0% /sys/fs/cgroup
/dev/mapper/centos-root
                        115.2G      4.9G    104.4G   4% /dev/termination-log
/dev/mapper/centos-root
                        115.2G      4.9G    104.4G   4% /etc/resolv.conf
/dev/mapper/centos-root
                        115.2G      4.9G    104.4G   4% /etc/hostname
/dev/mapper/centos-root
                        115.2G      4.9G    104.4G   4% /etc/hosts
shm                      64.0M         0     64.0M   0% /dev/shm
/dev/rbd0               975.9M      1.3M    907.4M   0% /usr/share/nginx/html
tmpfs                    62.8G     12.0K     62.8G   0% /var/run/secrets/kubernetes.io/serviceaccount
tmpfs                    62.8G         0     62.8G   0% /proc/kcore
tmpfs                    62.8G         0     62.8G   0% /proc/timer_list
tmpfs                    62.8G         0     62.8G   0% /proc/timer_stats
tmpfs                    62.8G         0     62.8G   0% /proc/sched_debug
tmpfs                    62.8G         0     62.8G   0% /sys/firmware


## 写入测试
[root@k8s-node-28 cephrbd]# kubectl exec -it nginx-dm-1697116629-j4m1g -- touch /usr/share/nginx/html/jicki.html

# 删除 pod
[root@k8s-node-28 cephrbd]# kubectl delete pod/nginx-dm-1697116629-j4m1g
pod "nginx-dm-1697116629-j4m1g" deleted

# 新的 pod
[root@k8s-node-28 cephrbd]# kubectl get pods
NAME                        READY     STATUS    RESTARTS   AGE
nginx-dm-1697116629-v2n52   1/1       Running   0          11s

# 查看文件
[root@k8s-node-28 cephrbd]# kubectl exec -it nginx-dm-1697116629-v2n52 -- ls -lt /usr/share/nginx/html
total 16
-rw-r--r--    1 root     root             0 May  9 07:53 jicki.html

```





# FAQ and Error



## FAQ 1


```
# ceph 挂载关于 namespace != default 的情况下


1. 需要配置 secret = namespace

kubectl get secret  可以看到只存在 default 的 认证


2. 需要配置 pv = namespace ，虽然pv 是全局的，但是最好还是配置一下

kubectl get pv --namespace=


3. 需要配置 pvc = namespace 。

kubectl get pvc --namespace=


# 确认以上3点，否则在挂载namespace 的时候会报错

ceph timeout expired waiting for volumes to attach/mount for pod




```



## Error 1

```
# 删除 osd 如果创建出错，请删除 osd

# 使用 ceph osd tree 查看出错的 osd.id

ceph osd out id   提出集群

# 登陆所在 服务器 停止服务

systemctl stop ceph-osd@id

# 然后再执行  crush remove 删除 tree

ceph osd crush remove osd.id

# 最后 执行 auth del 和 rm 删除

ceph auth del osd.id 

ceph osd rm id

# 在运行 tree 已经不存在了
# 如果使用分区挂载 还需要 登陆 osd 节点中 umount 分区
# 使用 lsblk 查看挂载
# 如: /dev/sdb 被挂载

umount /dev/sdb
ceph-disk zap /dev/sdb

```


## Error 2

```
# 如果暴力删除了 osd 然后重建了 osd

ceph -s

报错 stale+active+undersized+degraded+remapped

# 执行 ceph health detail 查看所有报错的信息

# 执行 ceph pg <pgid> query  查看具体信息
# 如果执行 query 报错
# 使用 ceph pg force_create_pg <pgid>  覆盖错误

# 如果 datail 有很多 pgid 出错 使用 for 循环去跑

for pg in `ceph health detail | grep "stale+active+undersized+degraded" | awk '{print $2}' | sort | uniq`; 
do 
  ceph pg force_create_pg $pg
done

# 如果仍然不能解决问题，那么请使用暴力方法
# 重建ceph 集群 删除集群所有文件,重新部署

systemctl stop ceph-osd.target
systemctl stop ceph-mon.target


rm -rf /etc/ceph/*
rm -rf /var/log/ceph/*
rm -rf /var/lib/ceph/*

# osd 清除挂载 目录的里的所有内容

# 卸载 ceph
yum -y remove ceph ceph-radosgw

# 重新安装 
yum -y install ceph ceph-radosgw

# 在 osd 节点 创建目录
mkdir -p /var/lib/ceph/osd/
```

