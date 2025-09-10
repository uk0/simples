---
layout: post
categories: cdh
title: kubernetes 部署cdh 通过 NASNFS）进行文件存储。
date: 2018-09-20 15:59:04 +0800
description: 在kubernetes集群部署cdh 并且部署hadoop 服务。
keywords: hadoop cdh kubernetes
---

技术验证篇，通过Kubernetes搭建一个可扩展的 CDH，并且启动hadoop。


## Quick

 * 主要用到`Kubernetes` 阿里云
 * 镜像来源 `https://github.com/windawings/spot-docker` .
 * 启动文件.`yaml`文件
 * NAS 阿里云弹性存储 NAS



##  注意
 * NAS 需要开启`VPC`网络安全组
 * 镜像需要上传到自己的私有仓库
 * Kubernetes 下载镜像需要配置`yaml << imagePullSecrets`,文章结尾会将`yaml`奉上
 
 * Kubernetes默认会创建一个`SLB`（用来给Ingress 使用)，我们用来暴露CDH的`UI`就是需要采用 `Ingress` 来进行暴露。
 * `Ingress` 所在 `Namespace = kube-system`

![](http://112firshme11224.test.upcdn.net/demos/17612368-36b7-49c3-bcbd-6bd7848c849a.png)

*  域名解析 注意是`泛域名`

## 开始部署

```yaml
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: cloudera-sa
  namespace: cloudera
  labels:
    name: sa
    app: cloudera
automountServiceAccountToken: true
---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  namespace: cloudera
  name: cloudera-cr
  labels:
    name: cr
    app: cloudera
rules:
- apiGroups: ["*"]
  resources: ["*"]
  verbs: ["*"]
- nonResourceURLs: ["*"]
  verbs: ["*"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: cloudera-crb
  namespace: cloudera
  labels:
    name: crb
    app: cloudera
subjects:
- kind: ServiceAccount
  name: cloudera-sa
  namespace: cloudera
roleRef:
  kind: ClusterRole
  name: cloudera-cr
  apiGroup: rbac.authorization.k8s.io
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: cloudera-config
  namespace: cloudera
  labels:
    name: cloudera-config
    app: cloudera
data:
  master_host: master-svc
  cron.tab: |
    */30 * * * * ntpdate master-svc && hwclock --systohc

  config.ini: |-
    [General]
    # Hostname of the CM server.
    server_host=master-svc

    # Port that the CM server is listening on.
    server_port=7182

    ## It should not normally be necessary to modify these.
    # Port that the CM agent should listen on.
    # listening_port=9000

    # IP Address that the CM agent should listen on.
    # listening_ip=

    # Hostname that the CM agent reports as its hostname. If unset, will be
    # obtained in code through something like this:
    #
    #   python -c 'import socket; \
    #              print socket.getfqdn(), \
    #                    socket.gethostbyname(socket.getfqdn())'
    #
    # listening_hostname=

    # An alternate hostname to report as the hostname for this host in CM.
    # Useful when this agent is behind a load balancer or proxy and all
    # inbound communication must connect through that proxy.
    # reported_hostname=

    # Port that supervisord should listen on.
    # NB: This only takes effect if supervisord is restarted.
    # supervisord_port=19001

    # Log file.  The supervisord log file will be placed into
    # the same directory.  Note that if the agent is being started via the
    # init.d script, /var/log/cloudera-scm-agent/cloudera-scm-agent.out will
    # also have a small amount of output (from before logging is initialized).
    # log_file=/var/log/cloudera-scm-agent/cloudera-scm-agent.log

    # Persistent state directory.  Directory to store CM agent state that
    # persists across instances of the agent process and system reboots.
    # Particularly, the agent's UUID is stored here.
    # lib_dir=/var/lib/cloudera-scm-agent

    # Parcel directory.  Unpacked parcels will be stored in this directory.
    # Downloaded parcels will be stored in <parcel_dir>/../parcel-cache
    # parcel_dir=/opt/cloudera/parcels

    # Enable supervisord event monitoring.  Used in eager heartbeating, amongst
    # other things.
    # enable_supervisord_events=true

    # Maximum time to wait (in seconds) for all metric collectors to finish
    # collecting data.
    max_collection_wait_seconds=10.0

    # Maximum time to wait (in seconds) when connecting to a local role's
    # webserver to fetch metrics.
    metrics_url_timeout_seconds=30.0

    # Maximum time to wait (in seconds) when connecting to a local TaskTracker
    # to fetch task attempt data.
    task_metrics_timeout_seconds=5.0

    # The list of non-device (nodev) filesystem types which will be monitored.
    monitored_nodev_filesystem_types=nfs,nfs4,tmpfs

    # The list of filesystem types which are considered local for monitoring purposes.
    # These filesystems are combined with the other local filesystem types found in
    # /proc/filesystems
    local_filesystem_whitelist=ext2,ext3,ext4,xfs

    # The largest size impala profile log bundle that this agent will serve to the
    # CM server. If the CM server requests more than this amount, the bundle will
    # be limited to this size. All instances of this limit being hit are logged to
    # the agent log.
    impala_profile_bundle_max_bytes=1073741824

    # The largest size stacks log bundle that this agent will serve to the CM
    # server. If the CM server requests more than this amount, the bundle will be
    # limited to this size. All instances of this limit being hit are logged to the
    # agent log.
    stacks_log_bundle_max_bytes=1073741824

    # The size to which the uncompressed portion of a stacks log can grow before it
    # is rotated. The log will then be compressed during rotation.
    stacks_log_max_uncompressed_file_size_bytes=5242880

    # The orphan process directory staleness threshold. If a diretory is more stale
    # than this amount of seconds, CM agent will remove it.
    orphan_process_dir_staleness_threshold=5184000

    # The orphan process directory refresh interval. The CM agent will check the
    # staleness of the orphan processes config directory every this amount of
    # seconds.
    orphan_process_dir_refresh_interval=3600

    # A knob to control the agent logging level. The options are listed as follows:
    # 1) DEBUG (set the agent logging level to 'logging.DEBUG')
    # 2) INFO (set the agent logging level to 'logging.INFO')
    scm_debug=INFO

    # The DNS resolution collecion interval in seconds. A java base test program
    # will be executed with at most this frequency to collect java DNS resolution
    # metrics. The test program is only executed if the associated health test,
    # Host DNS Resolution, is enabled.
    dns_resolution_collection_interval_seconds=60

    # The maximum time to wait (in seconds) for the java test program to collect
    # java DNS resolution metrics.
    dns_resolution_collection_timeout_seconds=30

    # The directory location in which the agent-wide kerberos credential cache
    # will be created.
    # agent_wide_credential_cache_location=/var/run/cloudera-scm-agent

    [Security]
    # Use TLS and certificate validation when connecting to the CM server.
    use_tls=0

    # The maximum allowed depth of the certificate chain returned by the peer.
    # The default value of 9 matches the default specified in openssl's
    # SSL_CTX_set_verify.
    max_cert_depth=9

    # A file of CA certificates in PEM format. The file can contain several CA
    # certificates identified by
    #
    # -----BEGIN CERTIFICATE-----
    # ... (CA certificate in base64 encoding) ...
    # -----END CERTIFICATE-----
    #
    # sequences. Before, between, and after the certificates text is allowed which
    # can be used e.g. for descriptions of the certificates.
    #
    # The file is loaded once, the first time an HTTPS connection is attempted. A
    # restart of the agent is required to pick up changes to the file.
    #
    # Note that if neither verify_cert_file or verify_cert_dir is set, certificate
    # verification will not be performed.
    # verify_cert_file=

    # Directory containing CA certificates in PEM format. The files each contain one
    # CA certificate. The files are looked up by the CA subject name hash value,
    # which must hence be available. If more than one CA certificate with the same
    # name hash value exist, the extension must be different (e.g. 9d66eef0.0,
    # 9d66eef0.1 etc). The search is performed in the ordering of the extension
    # number, regardless of other properties of the certificates. Use the c_rehash
    # utility to create the necessary links.
    #
    # The certificates in the directory are only looked up when required, e.g. when
    # building the certificate chain or when actually performing the verification
    # of a peer certificate. The contents of the directory can thus be changed
    # without an agent restart.
    #
    # When looking up CA certificates, the verify_cert_file is first searched, then
    # those in the directory. Certificate matching is done based on the subject name,
    # the key identifier (if present), and the serial number as taken from the
    # certificate to be verified. If these data do not match, the next certificate
    # will be tried. If a first certificate matching the parameters is found, the
    # verification process will be performed; no other certificates for the same
    # parameters will be searched in case of failure.
    #
    # Note that if neither verify_cert_file or verify_cert_dir is set, certificate
    # verification will not be performed.
    # verify_cert_dir=

    # PEM file containing client private key.
    # client_key_file=

    # A command to run which returns the client private key password on stdout
    # client_keypw_cmd=

    # If client_keypw_cmd isn't specified, instead a text file containing
    # the client private key password can be used.
    # client_keypw_file=

    # PEM file containing client certificate.
    # client_cert_file=

    ## Location of Hadoop files.  These are the CDH locations when installed by
    ## packages.  Unused when CDH is installed by parcels.
    [Hadoop]
    #cdh_crunch_home=/usr/lib/crunch
    #cdh_flume_home=/usr/lib/flume-ng
    #cdh_hadoop_bin=/usr/bin/hadoop
    #cdh_hadoop_home=/usr/lib/hadoop
    #cdh_hbase_home=/usr/lib/hbase
    #cdh_hbase_indexer_home=/usr/lib/hbase-solr
    #cdh_hcat_home=/usr/lib/hive-hcatalog
    #cdh_hdfs_home=/usr/lib/hadoop-hdfs
    #cdh_hive_home=/usr/lib/hive
    #cdh_httpfs_home=/usr/lib/hadoop-httpfs
    #cdh_hue_home=/usr/share/hue
    #cdh_hue_plugins_home=/usr/lib/hadoop
    #cdh_impala_home=/usr/lib/impala
    #cdh_llama_home=/usr/lib/llama
    #cdh_mr1_home=/usr/lib/hadoop-0.20-mapreduce
    #cdh_mr2_home=/usr/lib/hadoop-mapreduce
    #cdh_oozie_home=/usr/lib/oozie
    #cdh_parquet_home=/usr/lib/parquet
    #cdh_pig_home=/usr/lib/pig
    #cdh_solr_home=/usr/lib/solr
    #cdh_spark_home=/usr/lib/spark
    #cdh_sqoop_home=/usr/lib/sqoop
    #cdh_sqoop2_home=/usr/lib/sqoop2
    #cdh_yarn_home=/usr/lib/hadoop-yarn
    #cdh_zookeeper_home=/usr/lib/zookeeper
    #hive_default_xml=/etc/hive/conf.dist/hive-default.xml
    #webhcat_default_xml=/etc/hive-webhcat/conf.dist/webhcat-default.xml
    #jsvc_home=/usr/libexec/bigtop-utils
    #tomcat_home=/usr/lib/bigtop-tomcat
    #oracle_home=/usr/share/oracle/instantclient

    ## Location of Cloudera Management Services files.
    [Cloudera]
    #mgmt_home=/usr/share/cmf

    ## Location of JDBC Drivers.
    [JDBC]
    cloudera_mysql_connector_jar=/usr/java/latest/mysql-connector-java.jar
    #cloudera_oracle_connector_jar=/usr/share/java/oracle-connector-java.jar
    #By default, postgres jar is found dynamically in $MGMT_HOME/lib
    #cloudera_postgresql_jdbc_jar=
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nfs-deploy
  namespace: cloudera
  labels:
    name: nfs
    app: cloudera
spec:
  replicas: 1
  revisionHistoryLimit: 2
  selector:
    matchLabels:
      name: nfs
      app: cloudera
  template:
    metadata:
      name: nfs-client
      namespace: cloudera
      labels:
        name: nfs
        app: cloudera
    spec:
      tolerations:
      - key: node-role.kubernetes.io/master
        effect: NoSchedule
      serviceAccountName: cloudera-sa
      containers:
      - image: registry.cn-hangzhou.aliyuncs.com/acs/alicloud-nas-controller:v1.8.4
#        imagePullPolicy: Never
        name: nfs-client
#        securityContext:
#          privileged: true
        volumeMounts:
        - name: nfs-client-root
          mountPath: /persistentvolumes
        env:
        - name: PROVISIONER_NAME
          value: alicloud/nas
        - name: NFS_SERVER
          value: 074ca49d80-ak47.cn-hangzhou.nas.aliyuncs.com
        - name: NFS_PATH
          value: /pv/cloudera/agent/
      volumes:
      - name: nfs-client-root
        nfs:
          server: 074ca49d80-ak47.cn-hangzhou.nas.aliyuncs.com
          path: /pv/cloudera/agent/
---
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: cloudera-storage
  namespace: cloudera
  labels:
    name: storage
    app: cloudera
reclaimPolicy: Delete
provisioner: alicloud/nas
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: master-deploy
  namespace: cloudera
  labels:
    name: master-deploy
    app: cloudera
spec:
  replicas: 1
  revisionHistoryLimit: 2
  selector:
    matchLabels:
      name: master
      app: cloudera
  template:
    metadata:
      name: master
      namespace: cloudera
      labels:
        name: master
        app: cloudera
    spec:
      tolerations:
      - effect: NoSchedule
        key: node-role.kubernetes.io/master
      restartPolicy: Always
      serviceAccountName: cloudera-sa
      automountServiceAccountToken: true
      hostname: master-svc
      containers:
      - name: cloudera-master
        image: registry-vpc.cn-hangzhou.aliyuncs.com/huaching_prod/cloudera-master:5.14.2
#        imagePullPolicy: IfNotPresent
        resources:
          limits:
            memory: "4096Mi"
            cpu: "1000m"
        lifecycle:
          preStop:
            exec:
              command: ["/bin/bash","-c","/cloudera-init/run/stop.sh >> /cloudera/log/stop.log 2>&1"]
        volumeMounts:
        #- name: cloudera-init
          #mountPath: /cloudera-init/run/
        - name: mysql
          mountPath: /var/lib/mysql/
        - name: repo
          mountPath: /opt/cloudera/parcel-repo/
        - name: config
          mountPath: /opt/cm/etc/cloudera-scm-agent/
        env:
        - name: HOSTNAME
          valueFrom:
            configMapKeyRef:
              name: cloudera-config
              key: master_host
        - name: POD_IP
          valueFrom:
            fieldRef:
              apiVersion: v1
              fieldPath: status.podIP
        - name: POD_NAME
          valueFrom:
            fieldRef:
              apiVersion: v1
              fieldPath: metadata.name
        ports:
        - containerPort: 7180
          name: ui
        - containerPort: 9000
          name: agent
        - containerPort: 7182
          name: server
        - containerPort: 7191
          name: parcel
        - containerPort: 4434
          name: parcel-ssl
        securityContext:
          privileged: true
      imagePullSecrets:
      - name: registrykey-vpc
      volumes:
      - name: repo
        nfs:
          server: 074ca49d80-ak47.cn-hangzhou.nas.aliyuncs.com
          path: /pv/cloudera/repo/
      - name: mysql
        nfs:
          server: 074ca49d80-ak47.cn-hangzhou.nas.aliyuncs.com
          path: /pv/cloudera/master/
      #- name: cloudera-init
        #configMap:
          #name: cloudera-config
          #defaultMode: 0751
          #items:
          #- key: master.start
            #path: start.sh
          #- key: master.stop
            #path: stop.sh
          #- key: mysql.sql
            #path: mysql.sql
      - name: config
        configMap:
          name: cloudera-config
          items:
          - key: config.ini
            path: config.ini
---
apiVersion: v1
kind: Service
metadata:
  name: master-svc
  namespace: cloudera
  labels:
    name: master-svc
    app: cloudera
spec:
    type: ClusterIP
    clusterIP: None
    ports:
    - name: agent
      port: 9000
      protocol: TCP
      targetPort: 9000
    - name: server
      port: 7182
      protocol: TCP
      targetPort: 7182
    - name: parcel
      port: 7191
      protocol: TCP
      targetPort: 7191
    - name: parcel-ssl
      port: 4434
      protocol: TCP
      targetPort: 4434
    - name: ui
      port: 7180
      protocol: TCP
      targetPort: 7180
    selector:
        name: master
        app: cloudera
---
apiVersion: v1
kind: Service
metadata:
  name: master-external
  namespace: cloudera
  labels:
    name: master-external
    app: cloudera
spec:
    type: ClusterIP
    ports:
    - name: ui
      port: 7180
      protocol: TCP
      targetPort: 7180
    selector:
        name: master
        app: cloudera
---
apiVersion: policy/v1beta1
kind: PodDisruptionBudget
metadata:
  name: agent-pdb
  namespace: cloudera
  labels:
    name: agent-pdb
    app: cloudera
spec:
  selector:
    matchLabels:
      name: agent
      app: cloudera
  maxUnavailable: 1
---
apiVersion: v1
kind: Service
metadata:
  name: agent-svc
  namespace: cloudera
  labels:
    name: agent-svc
    app: cloudera
spec:
    type: ClusterIP
    clusterIP: None
    ports:
    - name: agent
      port: 9000
      protocol: TCP
      targetPort: 9000
    - name: server
      port: 7182
      protocol: TCP
      targetPort: 7182
    - name: parcel
      port: 7191
      protocol: TCP
      targetPort: 7191
    - name: parcel-ssl
      port: 4434
      protocol: TCP
      targetPort: 4434
    - name: parcel-ssl2
      port: 4433
      protocol: TCP
      targetPort: 4433
    selector:
        name: agent
        app: cloudera
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: agent
  namespace: cloudera
  labels:
    name: agent-sts
    app: cloudera
spec:
  selector:
    matchLabels:
      name: agent
      app: cloudera
  serviceName: agent-svc
  updateStrategy:
    type: RollingUpdate
  podManagementPolicy: OrderedReady
  replicas: 3
  revisionHistoryLimit: 2
  template:
    metadata:
      name: cagent
      namespace: cloudera
      labels:
        name: agent
        app: cloudera
    spec:
      tolerations:
      - effect: NoSchedule
        key: node-role.kubernetes.io/master
      restartPolicy: Always
      serviceAccountName: cloudera-sa
      automountServiceAccountToken: true
      containers:
      - name: agent
        image: registry-vpc.cn-hangzhou.aliyuncs.com/huaching_prod/cloudera-node:5.14.2
#        imagePullPolicy: IfNotPresent
        resources:
          limits:
            memory: "4096Mi"
            cpu: "1000m"
        lifecycle:
          preStop:
            exec:
              command: ["/bin/bash","-c","/cloudera-init/run/stop.sh"]
        volumeMounts:
        - name: repo
          mountPath: /opt/cloudera/parcel-repo/
        - name: lib
          mountPath: /opt/cm/lib/cloudera-scm-agent/
        - name: config
          mountPath: /opt/cm/etc/cloudera-scm-agent/
        ports:
        - containerPort: 9000
          name: agent
        - containerPort: 7182
          name: server
        - containerPort: 7191
          name: parcel
        - containerPort: 4434
          name: parcel-ssl
        - containerPort: 4433
          name: parcel-ssl2
        securityContext:
          privileged: true
      imagePullSecrets:
      - name: registrykey-vpc
      volumes:
      - name: config
        configMap:
          name: cloudera-config
          items:
          - key: config.ini
            path: config.ini
          - key: cron.tab
            path: cron.tab
      - name: repo
        nfs:
          server: 074ca49d80-ak47.cn-hangzhou.nas.aliyuncs.com
          path: /pv/cloudera/repo/
  volumeClaimTemplates:
  - metadata:
      name: lib
      namespace: cloudera
    spec:
      #volumeName: lib
      volumeMode: Filesystem
      accessModes: [ "ReadWriteOnce" ]
      storageClassName: cloudera-storage
      resources:
        requests:
          storage: 20Gi
---
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: cdh-master-ing
  namespace: cloudera
spec:
  rules:
  - host: cdh.docs.aliyun.com
    http:
      paths:
      - backend:
          serviceName: master-external
          servicePort: 7180

```

## 注意事项
* 创建`namespace cloudera`
* 镜像来源已经在上方说明了,可以自己去`rebuild` 
* 默认全部不会启动`command脚本可以看出`
* 启动流程进入`Master-deploy ` 去执行`start.sh`启动，根据机器情况，启动时间不等，期间可以用`netstat -tlunp` 查看端口是否占用来验证启动成功与否.
* `CDH` 需要用到Kubernetes的`特权模式`,相关细节请查关键字 `Kubernetes securityContext`
* 部分界面如下：
![](http://112firshme11224.test.upcdn.net/demos/2f2be0a5-ef33-4458-9cd4-433b9a3ebed4.png)


* 服务界面

    ![](http://112firshme11224.test.upcdn.net/demos/3596a98d-8de0-4641-b2b6-e7904061479d.png)

* PVC 界面

    ![](http://112firshme11224.test.upcdn.net/demos/e68c1177-32c0-419b-8acd-eeecbfd702dc.png)

*  CDH 节点界面

    ![](http://112firshme11224.test.upcdn.net/demos/49f7fc4a-2f25-4895-818a-0e6d296983b0.png)

* 节点服务运行情况。
    ![](http://112firshme11224.test.upcdn.net/demos/a054ea06-ffcb-4340-bfc9-7b307d7cbc99.png)






### 感谢 https://github.com/windawings

转载请注明出处，本文采用 [CC4.0](http://creativecommons.org/licenses/by-nc-nd/4.0/) 协议授权
