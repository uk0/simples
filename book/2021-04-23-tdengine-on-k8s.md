---
layout: post
categories: kubernetes
title: Tdengine on k8s
date: 2021-04-23 23:48:46 +0800
description: 将Tdengine 部署在k8s集群上,TDengine 集群模式
keywords: tdengine kuubernetes
---


将 Tdengine 集群模式部署到 Kubernetes ，支持扩节点，后面我会将自动脚本提交到 我的github。


#### 准备

* Tdengine 2.0.20 源代码
* K8S集群(存储 nfs ceph)
* yaml 文件


#### 为集群创建 存储（使用nfs）

```yaml

# 这里面提前创建好了  storage = nfs
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: testtdpvc1
spec:
  accessModes:
    - ReadWriteMany
  resources:
    requests:
      storage: "20"
  volumeName:
  storageClassName: nfs
---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: testtdpvc2
spec:
  accessModes:
    - ReadWriteMany
  resources:
    requests:
      storage: "20"
  volumeName:
  storageClassName: nfs
---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: testtdpvc3
spec:
  accessModes:
    - ReadWriteMany
  resources:
    requests:
      storage: "20"
  volumeName:
  storageClassName: nfs


```


#### 准备集群配置文件


```yaml

apiVersion: v1
kind: ConfigMap
metadata:
  name: tdmnode
  labels:
    app: tdmnode
data:
  ms1.cnf: |
    ########################################################
    #                                                      #
    #                  TDengine Configuration              #
    #   Any questions, please email support@taosdata.com   #
    #                                                      #
    ########################################################

    # first fully qualified domain name (FQDN) for TDengine system
    firstEp                   tdmnode-0.default.svc.cluster.local:6030
    secondEp                  tdmnode-1.default.svc.cluster.local:6030
    # local fully qualified domain name (FQDN)
    fqdn                      tdmnode-0.default.svc.cluster.local

    # first port number for the connection (12 continuous UDP/TCP port number are used)
    serverPort                6030

    # log file's director
    logDir                    /etc/taos/log

    # data file's directory
    dataDir                   /etc/taos/data

    # temporary file's directory
    tempDir                   /tmp/

    # the arbitrator's fully qualified domain name (FQDN) for TDengine system, for cluster only
    # arbitrator                arbitrator_hostname:6042

    # number of threads per CPU core
    # numOfThreadsPerCore       1.0

    # number of threads to commit cache data
    numOfCommitThreads        4

    # the proportion of total CPU cores available for query processing
    # 2.0: the query threads will be set to double of the CPU cores.
    # 1.0: all CPU cores are available for query processing [default].
    # 0.5: only half of the CPU cores are available for query.
    # 0.0: only one core available.
    # ratioOfQueryCores        1.0

    # the last_row/first/last aggregator will not change the original column name in the result fields
    # keepColumnName            0

    # number of management nodes in the system
    # numOfMnodes               3

    # enable/disable backuping vnode directory when removing vnode
    # vnodeBak                  1

    # enable/disable installation / usage report
    # telemetryReporting        1

    # enable/disable load balancing
    # balance                   1

    # role for dnode. 0 - any, 1 - mnode, 2 - dnode
    # role                      0

    # max timer control blocks
    # maxTmrCtrl                512

    # time interval of system monitor, seconds
    # monitorInterval           30

    # number of seconds allowed for a dnode to be offline, for cluster only
    # offlineThreshold          8640000

    # RPC re-try timer, millisecond
    # rpcTimer                  300

    # RPC maximum time for ack, seconds.
    # rpcMaxTime                600

    # time interval of dnode status reporting to mnode, seconds, for cluster only
    # statusInterval            1

    # time interval of heart beat from shell to dnode, seconds
    # shellActivityTimer        3

    # minimum sliding window time, milli-second
    # minSlidingTime            10

    # minimum time window, milli-second
    # minIntervalTime           10

    # maximum delay before launching a stream computation, milli-second
    # maxStreamCompDelay        20000

    # maximum delay before launching a stream computation for the first time, milli-second
    # maxFirstStreamCompDelay   10000

    # retry delay when a stream computation fails, milli-second
    # retryStreamCompDelay      10

    # the delayed time for launching a stream computation, from 0.1(default, 10% of whole computing time window) to 0.9
    # streamCompDelayRatio      0.1

    # max number of vgroups per db, 0 means configured automatically
    # maxVgroupsPerDb           0

    # max number of tables per vnode
    # maxTablesPerVnode         1000000

    # cache block size (Mbyte)
    # cache                     16

    # number of cache blocks per vnode
    # blocks                    6

    # number of days per DB file
    days                  3650

    # number of days to keep DB file
    # keep                  3650

    # minimum rows of records in file block
    # minRows               100

    # maximum rows of records in file block
    # maxRows               4096

    # the number of acknowledgments required for successful data writing
    # quorum                1

    # enable/disable compression
    # comp                  2

    # write ahead log (WAL) level, 0: no wal; 1: write wal, but no fysnc; 2: write wal, and call fsync
    # walLevel              1

    # if walLevel is set to 2, the cycle of fsync being executed, if set to 0, fsync is called right away
    # fsync                 3000

    # number of replications, for cluster only
    # replica               1

    # the compressed rpc message, option:
    #  -1 (no compression)
    #   0 (all message compressed),
    # > 0 (rpc message body which larger than this value will be compressed)
    # compressMsgSize       -1

    # max length of an SQL
    # maxSQLLength          65480

    # the maximum number of records allowed for super table time sorting
    # maxNumOfOrderedRes    100000

    # system time zone
    # timezone              Asia/Shanghai (CST, +0800)

    # system locale
    # locale                en_US.UTF-8

    # default system charset
    # charset               UTF-8

    # max number of connections allowed in dnode
    # maxShellConns         5000

    # max number of connections allowed in client
    # maxConnections        5000

    # stop writing logs when the disk size of the log folder is less than this value
    # minimalLogDirGB       0.1

    # stop writing temporary files when the disk size of the tmp folder is less than this value
    # minimalTmpDirGB       0.1

    # if disk free space is less than this value, taosd service exit directly within startup process
    # minimalDataDirGB      0.1

    # One mnode is equal to the number of vnode consumed
    # mnodeEqualVnodeNum    4

    # enbale/disable http service
    # http                  1

    # enable/disable system monitor
    # monitor               1

    # enable/disable recording the SQL statements via restful interface
    # httpEnableRecordSql   0

    # number of threads used to process http requests
    # httpMaxThreads        2

    # maximum number of rows returned by the restful interface
    # restfulRowLimit       10240

    # The following parameter is used to limit the maximum number of lines in log files.
    # max number of lines per log filters
    # numOfLogLines         10000000

    # enable/disable async log
    # asyncLog              1

    # time of keeping log files, days
    # logKeepDays           0


    # The following parameters are used for debug purpose only.
    # debugFlag 8 bits mask: FILE-SCREEN-UNUSED-HeartBeat-DUMP-TRACE_WARN-ERROR
    # 131: output warning and error
    # 135: output debug, warning and error
    # 143: output trace, debug, warning and error to log
    # 199: output debug, warning and error to both screen and file
    # 207: output trace, debug, warning and error to both screen and file

    # debug flag for all log type, take effect when non-zero value
    debugFlag             135

    # debug flag for meta management messages
    # mDebugFlag            135

    # debug flag for dnode messages
    # dDebugFlag            135

    # debug flag for sync module
    # sDebugFlag            135

    # debug flag for WAL
    # wDebugFlag            135

    # debug flag for SDB
    # sdbDebugFlag          135

    # debug flag for RPC
    # rpcDebugFlag          131

    # debug flag for TAOS TIMER
    # tmrDebugFlag          131

    # debug flag for TDengine client
    # cDebugFlag            131

    # debug flag for JNI
    # jniDebugFlag          131

    # debug flag for storage
    # uDebugFlag            131

    # debug flag for http server
    # httpDebugFlag         131

    # debug flag for monitor
    # monDebugFlag          131

    # debug flag for query
    # qDebugFlag            131

    # debug flag for vnode
    # vDebugFlag            131

    # debug flag for TSDB
    # tsdbDebugFlag         131

    # debug flag for continue query
    # cqDebugFlag           131

    # enable/disable recording the SQL in taos client
    # enableRecordSql    0

    # generate core file when service crash
    # enableCoreFile        1

    # maximum display width of binary and nchar fields in the shell. The parts exceeding this limit will be hidden
    # maxBinaryDisplayWidth 30

    # enable/disable stream (continuous query)
    # stream                1

    # in retrieve blocking model, only in 50% query threads will be used in query processing in dnode
    # retrieveBlockingModel    0

    # the maximum allowed query buffer size in MB during query processing for each data node
    # -1 no limit (default)
    # 0  no query allowed, queries are disabled
    # queryBufferSize         -1

  ms2.cnf: |
    ########################################################
    #                                                      #
    #                  TDengine Configuration              #
    #   Any questions, please email support@taosdata.com   #
    #                                                      #
    ########################################################

    # first fully qualified domain name (FQDN) for TDengine system
    firstEp                   tdmnode-0.default.svc.cluster.local:6030
    secondEp                  tdmnode-1.default.svc.cluster.local:6030
    # local fully qualified domain name (FQDN)
    fqdn                      tdmnode-1.default.svc.cluster.local

    # first port number for the connection (12 continuous UDP/TCP port number are used)
    serverPort                6030

    # log file's director
    logDir                    /etc/taos/log

    # data file's directory
    dataDir                   /etc/taos/data

    # temporary file's directory
    tempDir                   /tmp/

    # the arbitrator's fully qualified domain name (FQDN) for TDengine system, for cluster only
    # arbitrator                arbitrator_hostname:6042

    # number of threads per CPU core
    # numOfThreadsPerCore       1.0

    # number of threads to commit cache data
    numOfCommitThreads        4

    # the proportion of total CPU cores available for query processing
    # 2.0: the query threads will be set to double of the CPU cores.
    # 1.0: all CPU cores are available for query processing [default].
    # 0.5: only half of the CPU cores are available for query.
    # 0.0: only one core available.
    # ratioOfQueryCores        1.0

    # the last_row/first/last aggregator will not change the original column name in the result fields
    # keepColumnName            0

    # number of management nodes in the system
    # numOfMnodes               3

    # enable/disable backuping vnode directory when removing vnode
    # vnodeBak                  1

    # enable/disable installation / usage report
    # telemetryReporting        1

    # enable/disable load balancing
    # balance                   1

    # role for dnode. 0 - any, 1 - mnode, 2 - dnode
    # role                      0

    # max timer control blocks
    # maxTmrCtrl                512

    # time interval of system monitor, seconds
    # monitorInterval           30

    # number of seconds allowed for a dnode to be offline, for cluster only
    # offlineThreshold          8640000

    # RPC re-try timer, millisecond
    # rpcTimer                  300

    # RPC maximum time for ack, seconds.
    # rpcMaxTime                600

    # time interval of dnode status reporting to mnode, seconds, for cluster only
    # statusInterval            1

    # time interval of heart beat from shell to dnode, seconds
    # shellActivityTimer        3

    # minimum sliding window time, milli-second
    # minSlidingTime            10

    # minimum time window, milli-second
    # minIntervalTime           10

    # maximum delay before launching a stream computation, milli-second
    # maxStreamCompDelay        20000

    # maximum delay before launching a stream computation for the first time, milli-second
    # maxFirstStreamCompDelay   10000

    # retry delay when a stream computation fails, milli-second
    # retryStreamCompDelay      10

    # the delayed time for launching a stream computation, from 0.1(default, 10% of whole computing time window) to 0.9
    # streamCompDelayRatio      0.1

    # max number of vgroups per db, 0 means configured automatically
    # maxVgroupsPerDb           0

    # max number of tables per vnode
    # maxTablesPerVnode         1000000

    # cache block size (Mbyte)
    # cache                     16

    # number of cache blocks per vnode
    # blocks                    6

    # number of days per DB file
    days                  3650

    # number of days to keep DB file
    # keep                  3650

    # minimum rows of records in file block
    # minRows               100

    # maximum rows of records in file block
    # maxRows               4096

    # the number of acknowledgments required for successful data writing
    # quorum                1

    # enable/disable compression
    # comp                  2

    # write ahead log (WAL) level, 0: no wal; 1: write wal, but no fysnc; 2: write wal, and call fsync
    # walLevel              1

    # if walLevel is set to 2, the cycle of fsync being executed, if set to 0, fsync is called right away
    # fsync                 3000

    # number of replications, for cluster only
    # replica               1

    # the compressed rpc message, option:
    #  -1 (no compression)
    #   0 (all message compressed),
    # > 0 (rpc message body which larger than this value will be compressed)
    # compressMsgSize       -1

    # max length of an SQL
    # maxSQLLength          65480

    # the maximum number of records allowed for super table time sorting
    # maxNumOfOrderedRes    100000

    # system time zone
    # timezone              Asia/Shanghai (CST, +0800)

    # system locale
    # locale                en_US.UTF-8

    # default system charset
    # charset               UTF-8

    # max number of connections allowed in dnode
    # maxShellConns         5000

    # max number of connections allowed in client
    # maxConnections        5000

    # stop writing logs when the disk size of the log folder is less than this value
    # minimalLogDirGB       0.1

    # stop writing temporary files when the disk size of the tmp folder is less than this value
    # minimalTmpDirGB       0.1

    # if disk free space is less than this value, taosd service exit directly within startup process
    # minimalDataDirGB      0.1

    # One mnode is equal to the number of vnode consumed
    # mnodeEqualVnodeNum    4

    # enbale/disable http service
    # http                  1

    # enable/disable system monitor
    # monitor               1

    # enable/disable recording the SQL statements via restful interface
    # httpEnableRecordSql   0

    # number of threads used to process http requests
    # httpMaxThreads        2

    # maximum number of rows returned by the restful interface
    # restfulRowLimit       10240

    # The following parameter is used to limit the maximum number of lines in log files.
    # max number of lines per log filters
    # numOfLogLines         10000000

    # enable/disable async log
    # asyncLog              1

    # time of keeping log files, days
    # logKeepDays           0


    # The following parameters are used for debug purpose only.
    # debugFlag 8 bits mask: FILE-SCREEN-UNUSED-HeartBeat-DUMP-TRACE_WARN-ERROR
    # 131: output warning and error
    # 135: output debug, warning and error
    # 143: output trace, debug, warning and error to log
    # 199: output debug, warning and error to both screen and file
    # 207: output trace, debug, warning and error to both screen and file

    # debug flag for all log type, take effect when non-zero value
    debugFlag             135

    # debug flag for meta management messages
    # mDebugFlag            135

    # debug flag for dnode messages
    # dDebugFlag            135

    # debug flag for sync module
    # sDebugFlag            135

    # debug flag for WAL
    # wDebugFlag            135

    # debug flag for SDB
    # sdbDebugFlag          135

    # debug flag for RPC
    # rpcDebugFlag          131

    # debug flag for TAOS TIMER
    # tmrDebugFlag          131

    # debug flag for TDengine client
    # cDebugFlag            131

    # debug flag for JNI
    # jniDebugFlag          131

    # debug flag for storage
    # uDebugFlag            131

    # debug flag for http server
    # httpDebugFlag         131

    # debug flag for monitor
    # monDebugFlag          131

    # debug flag for query
    # qDebugFlag            131

    # debug flag for vnode
    # vDebugFlag            131

    # debug flag for TSDB
    # tsdbDebugFlag         131

    # debug flag for continue query
    # cqDebugFlag           131

    # enable/disable recording the SQL in taos client
    # enableRecordSql    0

    # generate core file when service crash
    # enableCoreFile        1

    # maximum display width of binary and nchar fields in the shell. The parts exceeding this limit will be hidden
    # maxBinaryDisplayWidth 30

    # enable/disable stream (continuous query)
    # stream                1

    # in retrieve blocking model, only in 50% query threads will be used in query processing in dnode
    # retrieveBlockingModel    0

    # the maximum allowed query buffer size in MB during query processing for each data node
    # -1 no limit (default)
    # 0  no query allowed, queries are disabled
    # queryBufferSize         -1
  ms3.cnf: |
    ########################################################
    #                                                      #
    #                  TDengine Configuration              #
    #   Any questions, please email support@taosdata.com   #
    #                                                      #
    ########################################################


    # first fully qualified domain name (FQDN) for TDengine system
    firstEp                   tdmnode-0.default.svc.cluster.local:6030
    secondEp                  tdmnode-1.default.svc.cluster.local:6030
    # local fully qualified domain name (FQDN)
    fqdn                      tdmnode-2.default.svc.cluster.local

    # first port number for the connection (12 continuous UDP/TCP port number are used)
    serverPort                6030


    # log file's director
    logDir                    /etc/taos/log

    # data file's directory
    dataDir                   /etc/taos/data

    # temporary file's directory
    tempDir                   /tmp/

    # the arbitrator's fully qualified domain name (FQDN) for TDengine system, for cluster only
    # arbitrator                arbitrator_hostname:6042

    # number of threads per CPU core
    # numOfThreadsPerCore       1.0

    # number of threads to commit cache data
    numOfCommitThreads        4

    # the proportion of total CPU cores available for query processing
    # 2.0: the query threads will be set to double of the CPU cores.
    # 1.0: all CPU cores are available for query processing [default].
    # 0.5: only half of the CPU cores are available for query.
    # 0.0: only one core available.
    # ratioOfQueryCores        1.0

    # the last_row/first/last aggregator will not change the original column name in the result fields
    # keepColumnName            0

    # number of management nodes in the system
    # numOfMnodes               3

    # enable/disable backuping vnode directory when removing vnode
    # vnodeBak                  1

    # enable/disable installation / usage report
    # telemetryReporting        1

    # enable/disable load balancing
    # balance                   1

    # role for dnode. 0 - any, 1 - mnode, 2 - dnode
    # role                      0

    # max timer control blocks
    # maxTmrCtrl                512

    # time interval of system monitor, seconds
    # monitorInterval           30

    # number of seconds allowed for a dnode to be offline, for cluster only
    # offlineThreshold          8640000

    # RPC re-try timer, millisecond
    # rpcTimer                  300

    # RPC maximum time for ack, seconds.
    # rpcMaxTime                600

    # time interval of dnode status reporting to mnode, seconds, for cluster only
    # statusInterval            1

    # time interval of heart beat from shell to dnode, seconds
    # shellActivityTimer        3

    # minimum sliding window time, milli-second
    # minSlidingTime            10

    # minimum time window, milli-second
    # minIntervalTime           10

    # maximum delay before launching a stream computation, milli-second
    # maxStreamCompDelay        20000

    # maximum delay before launching a stream computation for the first time, milli-second
    # maxFirstStreamCompDelay   10000

    # retry delay when a stream computation fails, milli-second
    # retryStreamCompDelay      10

    # the delayed time for launching a stream computation, from 0.1(default, 10% of whole computing time window) to 0.9
    # streamCompDelayRatio      0.1

    # max number of vgroups per db, 0 means configured automatically
    # maxVgroupsPerDb           0

    # max number of tables per vnode
    # maxTablesPerVnode         1000000

    # cache block size (Mbyte)
    # cache                     16

    # number of cache blocks per vnode
    # blocks                    6

    # number of days per DB file
    days                  3650

    # number of days to keep DB file
    # keep                  3650

    # minimum rows of records in file block
    # minRows               100

    # maximum rows of records in file block
    # maxRows               4096

    # the number of acknowledgments required for successful data writing
    # quorum                1

    # enable/disable compression
    # comp                  2

    # write ahead log (WAL) level, 0: no wal; 1: write wal, but no fysnc; 2: write wal, and call fsync
    # walLevel              1

    # if walLevel is set to 2, the cycle of fsync being executed, if set to 0, fsync is called right away
    # fsync                 3000

    # number of replications, for cluster only
    # replica               1

    # the compressed rpc message, option:
    #  -1 (no compression)
    #   0 (all message compressed),
    # > 0 (rpc message body which larger than this value will be compressed)
    # compressMsgSize       -1

    # max length of an SQL
    # maxSQLLength          65480

    # the maximum number of records allowed for super table time sorting
    # maxNumOfOrderedRes    100000

    # system time zone
    # timezone              Asia/Shanghai (CST, +0800)

    # system locale
    # locale                en_US.UTF-8

    # default system charset
    # charset               UTF-8

    # max number of connections allowed in dnode
    # maxShellConns         5000

    # max number of connections allowed in client
    # maxConnections        5000

    # stop writing logs when the disk size of the log folder is less than this value
    # minimalLogDirGB       0.1

    # stop writing temporary files when the disk size of the tmp folder is less than this value
    # minimalTmpDirGB       0.1

    # if disk free space is less than this value, taosd service exit directly within startup process
    # minimalDataDirGB      0.1

    # One mnode is equal to the number of vnode consumed
    # mnodeEqualVnodeNum    4

    # enbale/disable http service
    # http                  1

    # enable/disable system monitor
    # monitor               1

    # enable/disable recording the SQL statements via restful interface
    # httpEnableRecordSql   0

    # number of threads used to process http requests
    # httpMaxThreads        2

    # maximum number of rows returned by the restful interface
    # restfulRowLimit       10240

    # The following parameter is used to limit the maximum number of lines in log files.
    # max number of lines per log filters
    # numOfLogLines         10000000

    # enable/disable async log
    # asyncLog              1

    # time of keeping log files, days
    # logKeepDays           0


    # The following parameters are used for debug purpose only.
    # debugFlag 8 bits mask: FILE-SCREEN-UNUSED-HeartBeat-DUMP-TRACE_WARN-ERROR
    # 131: output warning and error
    # 135: output debug, warning and error
    # 143: output trace, debug, warning and error to log
    # 199: output debug, warning and error to both screen and file
    # 207: output trace, debug, warning and error to both screen and file

    # debug flag for all log type, take effect when non-zero value
    debugFlag             135

    # debug flag for meta management messages
    # mDebugFlag            135

    # debug flag for dnode messages
    # dDebugFlag            135

    # debug flag for sync module
    # sDebugFlag            135

    # debug flag for WAL
    # wDebugFlag            135

    # debug flag for SDB
    # sdbDebugFlag          135

    # debug flag for RPC
    # rpcDebugFlag          131

    # debug flag for TAOS TIMER
    # tmrDebugFlag          131

    # debug flag for TDengine client
    # cDebugFlag            131

    # debug flag for JNI
    # jniDebugFlag          131

    # debug flag for storage
    # uDebugFlag            131

    # debug flag for http server
    # httpDebugFlag         131

    # debug flag for monitor
    # monDebugFlag          131

    # debug flag for query
    # qDebugFlag            131

    # debug flag for vnode
    # vDebugFlag            131

    # debug flag for TSDB
    # tsdbDebugFlag         131

    # debug flag for continue query
    # cqDebugFlag           131

    # enable/disable recording the SQL in taos client
    # enableRecordSql    0

    # generate core file when service crash
    # enableCoreFile        1

    # maximum display width of binary and nchar fields in the shell. The parts exceeding this limit will be hidden
    # maxBinaryDisplayWidth 30

    # enable/disable stream (continuous query)
    # stream                1

    # in retrieve blocking model, only in 50% query threads will be used in query processing in dnode
    # retrieveBlockingModel    0

    # the maximum allowed query buffer size in MB during query processing for each data node
    # -1 no limit (default)
    # 0  no query allowed, queries are disabled
    # queryBufferSize         -1



```



#### 节点配置文件 1

```yaml

apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: tdmnode-0
spec:
  selector:
    matchLabels:
      app: tdmnode-0
  serviceName: "tdmnode-0"
  replicas: 1
  template:
    metadata:
      labels:
        app: tdmnode-0
    spec:
      containers:
        - name: tdmnode-0
          image: tdengine/tdengine:latest
          ports:
            - containerPort: 6030
            - containerPort: 6035
            - containerPort: 6041
            - containerPort: 6031
            - containerPort: 6032
            - containerPort: 6033
            - containerPort: 6034
            - containerPort: 6036
            - containerPort: 6037
            - containerPort: 6038
            - containerPort: 6039
            - containerPort: 6040
          command:
            - bash
            - "-c"
            - |
              set -ex
              apt install inetutils-ping dnsutils -y
              # nodes
              export nodes=3
              export end_index=$(($nodes - 1))
              while :; do
                sleep 1
                export index=0
                for i in `seq 0 $end_index`; do
                 tmpip=`dig tdmnode-$i.default.svc.cluster.local | grep "ANSWER SECTION" -A 1 | awk 'NR>1 {print$5}'`
                 if [ -z "$tmpip" ]; then
                    echo -e "\033[31m wait Node tdmnode-$i ip is null  \033[0m";
                 else
                    index=$(($index + 1));
                 fi
                done

                if [ $index -ge $nodes ] ; then
                      echo -e "\033[32m ips is ready .Starting now Tdengine Taosd.... \033[0m"
                      for i in `seq 0 $end_index`; do
                          tmpip=`dig tdmnode-$i.default.svc.cluster.local | grep "ANSWER SECTION" -A 1 | awk 'NR>1 {print$5}'`
                          echo  "${tmpip}       tdmnode-$i.default.svc.cluster.local" >> /etc/hosts ;
                          done
                      break
                 else
                     index=0
                     echo -e "\033[31m all node not get ready pls wait....  \033[0m";
                 fi
              done
              taosd
          env:
            - name: MY_POD_IP
              valueFrom:
                fieldRef:
                  fieldPath: status.podIP
          volumeMounts:
            - name: config-map
              mountPath: /etc/taos/taos.cfg
              subPath: ms1.cnf
            - name: datadir
              mountPath: /etc/taos/data
      volumes:
        - name: config-map
          configMap:
            name: tdmnode
            items:
              - key: ms1.cnf
                path: ms1.cnf
        - name: datadir
          persistentVolumeClaim:
            claimName: testtdpvc3

---
apiVersion: v1
kind: Service
metadata:
  name: tdmnode-0
  labels:
    app: tdmnode-0
spec:
  clusterIP: None
  ports:
    - name: port6030
      port: 6030
      protocol: TCP
      targetPort: 6030
    - name: port6031
      port: 6031
      protocol: TCP
      targetPort: 6031
    - name: port6032
      port: 6032
      protocol: TCP
      targetPort: 6032
    - name: port6033
      port: 6033
      protocol: TCP
      targetPort: 6033
    - name: port6034
      port: 6034
      protocol: TCP
      targetPort: 6034
    - name: port6035
      port: 6035
      protocol: TCP
      targetPort: 6035
    - name: port6036
      port: 6036
      protocol: TCP
      targetPort: 6036
    - name: port6037
      port: 6037
      protocol: TCP
      targetPort: 6037
    - name: port6038
      port: 6038
      protocol: TCP
      targetPort: 6038
    - name: port6039
      port: 6039
      protocol: TCP
      targetPort: 6039
    - name: port6040
      port: 6040
      protocol: TCP
      targetPort: 6040
    - name: port6041
      port: 6041
      protocol: TCP
      targetPort: 6041
  selector:
    app: tdmnode-0
  sessionAffinity: None
  type: ClusterIP

```

#### 节点配置文件 2

```yaml

apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: tdmnode-1
spec:
  selector:
    matchLabels:
      app: tdmnode-1
  serviceName: "tdmnode-1"
  replicas: 1
  template:
    metadata:
      labels:
        app: tdmnode-1
    spec:
      containers:
        - name: tdmnode-1
          image: tdengine/tdengine:latest
          ports:
            - containerPort: 6030
            - containerPort: 6035
            - containerPort: 6041
            - containerPort: 6031
            - containerPort: 6032
            - containerPort: 6033
            - containerPort: 6034
            - containerPort: 6036
            - containerPort: 6037
            - containerPort: 6038
            - containerPort: 6039
            - containerPort: 6040
          command:
            - bash
            - "-c"
            - |
              set -ex
              apt install inetutils-ping dnsutils -y
              # nodes
              export nodes=3
              export end_index=$(($nodes - 1))
              while :; do
                sleep 1
                export index=0
                for i in `seq 0 $end_index`; do
                 tmpip=`dig tdmnode-$i.default.svc.cluster.local | grep "ANSWER SECTION" -A 1 | awk 'NR>1 {print$5}'`
                 if [ -z "$tmpip" ]; then
                    echo -e "\033[31m wait Node tdmnode-$i ip is null  \033[0m";
                 else
                    index=$(($index + 1))
                 fi
                done

                if [ $index -ge $nodes ] ; then
                      echo -e "\033[32m ips is ready .Starting now Tdengine Taosd.... \033[0m"
                      for i in `seq 0 $end_index`; do
                          tmpip=`dig tdmnode-$i.default.svc.cluster.local | grep "ANSWER SECTION" -A 1 | awk 'NR>1 {print$5}'`
                          echo  "${tmpip}       tdmnode-$i.default.svc.cluster.local" >> /etc/hosts ;
                          done
                      break
                 else
                     index=0
                     echo -e "\033[31m all node not get ready pls wait....  \033[0m";
                 fi
              done
              taosd
          env:
            - name: MY_POD_IP
              valueFrom:
                fieldRef:
                  fieldPath: status.podIP
          volumeMounts:
            - name: config-map
              mountPath: /etc/taos/taos.cfg
              subPath: ms2.cnf
            - name: datadir
              mountPath: /etc/taos/data
      volumes:
        - name: config-map
          configMap:
            name: tdmnode
            items:
              - key: ms2.cnf
                path: ms2.cnf
        - name: datadir
          persistentVolumeClaim:
            claimName: testtdpvc2
---
apiVersion: v1
kind: Service
metadata:
  name: tdmnode-1
  labels:
    app: tdmnode-1
spec:
  clusterIP: None
  ports:
    - name: port6030
      port: 6030
      protocol: TCP
      targetPort: 6030
    - name: port6031
      port: 6031
      protocol: TCP
      targetPort: 6031
    - name: port6032
      port: 6032
      protocol: TCP
      targetPort: 6032
    - name: port6033
      port: 6033
      protocol: TCP
      targetPort: 6033
    - name: port6034
      port: 6034
      protocol: TCP
      targetPort: 6034
    - name: port6035
      port: 6035
      protocol: TCP
      targetPort: 6035
    - name: port6036
      port: 6036
      protocol: TCP
      targetPort: 6036
    - name: port6037
      port: 6037
      protocol: TCP
      targetPort: 6037
    - name: port6038
      port: 6038
      protocol: TCP
      targetPort: 6038
    - name: port6039
      port: 6039
      protocol: TCP
      targetPort: 6039
    - name: port6040
      port: 6040
      protocol: TCP
      targetPort: 6040
    - name: port6041
      port: 6041
      protocol: TCP
      targetPort: 6041
  selector:
    app: tdmnode-1
  sessionAffinity: None
  type: ClusterIP

```

#### 节点配置文件 3


```yaml

apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: tdmnode-2
spec:
  selector:
    matchLabels:
      app: tdmnode-2
  serviceName: "tdmnode-2"
  replicas: 1
  template:
    metadata:
      labels:
        app: tdmnode-2
    spec:
      containers:
        - name: tdmnode-2
          image: tdengine/tdengine:latest
          ports:
            - containerPort: 6030
            - containerPort: 6035
            - containerPort: 6041
            - containerPort: 6031
            - containerPort: 6032
            - containerPort: 6033
            - containerPort: 6034
            - containerPort: 6036
            - containerPort: 6037
            - containerPort: 6038
            - containerPort: 6039
            - containerPort: 6040
          command:
            - bash
            - "-c"
            - |
              set -ex
              apt install inetutils-ping dnsutils -y
              # nodes
              export nodes=3
              export end_index=$(($nodes - 1))
              while :; do
                sleep 1
                export index=0
                for i in `seq 0 $end_index`; do
                 tmpip=`dig tdmnode-$i.default.svc.cluster.local | grep "ANSWER SECTION" -A 1 | awk 'NR>1 {print$5}'`
                 if [ -z "$tmpip" ]; then
                    echo -e "\033[31m wait Node tdmnode-$i ip is null  \033[0m";
                 else
                    index=$(($index + 1))
                 fi
                done

                if [ $index -ge $nodes ] ; then
                      echo -e "\033[32m ips is ready .Starting now Tdengine Taosd.... \033[0m"
                      for i in `seq 0 $end_index`; do
                          tmpip=`dig tdmnode-$i.default.svc.cluster.local | grep "ANSWER SECTION" -A 1 | awk 'NR>1 {print$5}'`
                          echo  "${tmpip}       tdmnode-$i.default.svc.cluster.local" >> /etc/hosts ;
                          done
                      break
                 else
                     index=0
                     echo -e "\033[31m all node not get ready pls wait....  \033[0m";
                 fi
              done
              taosd
          env:
            - name: MY_POD_IP
              valueFrom:
                fieldRef:
                  fieldPath: status.podIP
          volumeMounts:
            - name: config-map
              mountPath: /etc/taos/taos.cfg
              subPath: ms3.cnf
            - name: datadir
              mountPath: /etc/taos/data
      volumes:
        - name: config-map
          configMap:
            name: tdmnode
            items:
              - key: ms3.cnf
                path: ms3.cnf
        - name: datadir
          persistentVolumeClaim:
            claimName: testtdpvc1
---
apiVersion: v1
kind: Service
metadata:
  name: tdmnode-2
  labels:
    app: tdmnode-2
spec:
  clusterIP: None
  ports:
    - name: port6030
      port: 6030
      protocol: TCP
      targetPort: 6030
    - name: port6031
      port: 6031
      protocol: TCP
      targetPort: 6031
    - name: port6032
      port: 6032
      protocol: TCP
      targetPort: 6032
    - name: port6033
      port: 6033
      protocol: TCP
      targetPort: 6033
    - name: port6034
      port: 6034
      protocol: TCP
      targetPort: 6034
    - name: port6035
      port: 6035
      protocol: TCP
      targetPort: 6035
    - name: port6036
      port: 6036
      protocol: TCP
      targetPort: 6036
    - name: port6037
      port: 6037
      protocol: TCP
      targetPort: 6037
    - name: port6038
      port: 6038
      protocol: TCP
      targetPort: 6038
    - name: port6039
      port: 6039
      protocol: TCP
      targetPort: 6039
    - name: port6040
      port: 6040
      protocol: TCP
      targetPort: 6040
    - name: port6041
      port: 6041
      protocol: TCP
      targetPort: 6041
  selector:
    app: tdmnode-2
  sessionAffinity: None
  type: ClusterIP


```



#### 注意事项

* 1 pvc 策略 非立即回收（重启集群 只删除 `StatefulSet ` 即可）
* 2 监控IP 解析 哪一块目前还不完善 需要一定的时间才能解析到IP 大约 `5s`


#### 内部使用方式


* 1 将Client 文件通过nfs映射到POD 或者其他 类型的 APP内部，通过 脚本初始化，大概内容如下。

```bash

       command：
        - sh 或者 bash
        - "-c"
        - |
        cd /plugin_free/client_td_2.0.20
        ./install_client.sh

```



目前在k8s 集群上已经连续运行了一周，期间多次强行破坏集群验证恢复能力，以及数据存储能力，目前集群数据大约15亿条（单表测试，还在继续写入。）



#### 运行图

* StatefulSet

![](/static/blog/2021-04-24-00-44-25.png)

* service

![](/static/blog/2021-04-24-00-45-23.png)


* taos

![](/static/blog/2021-04-24-00-47-07.png)



* 开始重启测试

![](/static/blog/2021-04-24-00-48-11.png)


* 重启后状态信息

![](/static/blog/2021-04-24-00-48-31.png)


* 端口等待 

![](/static/blog/2021-04-24-00-49-11.png)


* 重启后恢复


![](/static/blog/2021-04-24-00-49-32.png)


* 恢复后 还在继续写入（其中一个库的一个表）

![](/static/blog/2021-04-24-00-52-32.png)


* 补一张 PVC

![](/static/blog/2021-04-24-01-00-37.png)


#### tdengine 采用默认配置文件没有修改时区，镜像也没有修改时区，所以日志时间有误差，忽略即可，（强迫症可以自己加）


###  以上操作均在正式环境在运行，出现问题我会及时更新，感谢您的阅读。


转载请注明出处，本文采用 [CC4.0](http://creativecommons.org/licenses/by-nc-nd/4.0/) 协议授权
