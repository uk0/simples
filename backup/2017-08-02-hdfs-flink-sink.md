---
layout: post
title:  hadoop HDFS作为Flink 的 Sink
categories: shell,hadoop,hdfs
description: 笔记
keywords: hdfs,flink
---


因工作需求，学习了Hdfs分布式文件存储系统，所整合Flink + HDFS 作为一个Demo 帮助大家跳坑。
HDFS 采用NS高可用模式。

# HDFSManager.Java
    * 初始化HDFS链接。

```java
package com.e.firsh.spb.utils.hdfs;

import com.alibaba.fastjson.JSONObject;
import com.e.firsh.base.esb.BO;
import com.e.firsh.base.utils.Registry;
import org.apache.hadoop.conf.Configuration;
import org.apache.hadoop.fs.FSDataOutputStream;
import org.apache.hadoop.fs.FileSystem;
import org.apache.hadoop.fs.Path;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import java.io.IOException;
import java.util.Set;

/**
 * The type Hdfs manager.
 */
public class HDFSManager {
    private static Logger logger = LoggerFactory.getLogger(HDFSManager.class);
    private static Configuration configuration;
    private static FileSystem fs;

    public HDFSManager() {
        JSONObject jsonObject = new JSONObject();
        jsonObject.put("fs.defaultFS", "hdfs://ns");
        jsonObject.put("dfs.nameservices", "ns");
        jsonObject.put("dfs.ha.namenodes.ns","nn1,nn2");
        jsonObject.put("dfs.namenode.rpc-address.ns.nn1", "10.11.0.12:9000");
        jsonObject.put("dfs.namenode.rpc-address.ns.nn2", "10.11.0.13:9000");
        jsonObject.put("dfs.client.failover.proxy.provider.ns", "org.apache.hadoop.hdfs.server.namenode.ha.ConfiguredFailoverProxyProvider");
        init(jsonObject);
    }

    public boolean init(JSONObject args) {
        Configuration conf = new Configuration();
        Set<String> itr = args.keySet();
        if (itr != null) {
            for (String pname : itr) {
                String pvalue = args.getString(pname);
                conf.set(pname, pvalue);
            }
        }
        try {
            fs = FileSystem.get(conf);
            return true;
        } catch (Exception e) {
            logger.error(e.getMessage() + ":{}", e);
        }
        return false;
    }


    public BO appendToFile(String destPath, String line) throws IOException {
        BO respBO = BO.createResponseBO();
        Path path = new Path(destPath);
        FSDataOutputStream dos = null;
        try {
            if (!fs.exists(path)) {
                dos = fs.create(path);
            } else {
                dos = fs.append(path);
            }
            byte[] readBuf = line.getBytes("UTF-8");
            dos.write(readBuf, 0, readBuf.length);
            dos.close();
            respBO.setDataNameValue("length", readBuf.length);
        } catch (Exception e) {
            logger.error(e.getMessage() + ":{}", e);
            respBO.setCompleteCode(BO.BO_V_CC_FAILED);
        } finally {
            if (dos != null) {
                dos.close();
            }
        }
        return respBO;
    }

    public static void main(String[] args) {
        HDFSManager hdfsManager = new HDFSManager();
        try {
            hdfsManager.appendToFile("/testmq","hello w");
        } catch (IOException e) {
            e.printStackTrace();
        }
    }

    public static void init2(){
         configuration = new Configuration();
        /**TODO 添加Hadoop配置内容*/
        configuration.set("dfs.namenode.name.dir", "file:///home/admin/data/hadoop/hdfs/name");
        configuration.set("dfs.nameservices", "ns");
        configuration.set("dfs.ha.namenodes.ns", "nn1,nn2");
        configuration.set("dfs.namenode.rpc-address.ns.nn1", "10.11.0.12:9000");
        configuration.set("dfs.namenode.rpc-address.ns.nn2", "10.11.0.13:9000");
        configuration.set("dfs.namenode.shared.edits.dir", "qjournal://10.11.0.12:8485;10.11.0.13:8485;10.11.0.14:8485/ns");
        configuration.set("hadoop.tmp.dir", "/home/admin/data/hadoop/tmp");
        configuration.set("fs.defaultFS", "hdfs://ns");
        configuration.set("dfs.journalnode.edits.dir", "/home/admin/data/hadoop/journal");
        configuration.set("dfs.client.failover.proxy.provider.ns", "org.apache.hadoop.hdfs.server.namenode.ha.ConfiguredFailoverProxyProvider");
        configuration.set("ha.zookeeper.quorum", "10.11.0.12:2181,10.11.0.13:2181,10.11.0.14:2181");
        configuration.set("mapreduce.input.fileinputformat.split.minsize", "10");


        try {
            fs = FileSystem.get(configuration);
        } catch (Exception e) {
            logger.error(e.getMessage() + ":{}", e);
        }
    }
}

```
# HDFSSink.Java

 * 继承`RichSinkFunction`，重写`open` 和 `invoke`方法，还有`close`。

```java
package com.e.firsh.spb.sink;

import com.e.firsh.spb.utils.hdfs.HDFSManager;
import org.apache.flink.configuration.Configuration;
import org.apache.flink.streaming.api.functions.sink.RichSinkFunction;

/**
 * Created by zhangjianxin on 2017/8/1.
 */
public class HDFSSink extends RichSinkFunction<String> {
    private HDFSManager hdfsManager;
    private final static  String  HDFS_PATH ="/testmq";
    @Override
    public void invoke(String t) throws Exception {
        hdfsManager.appendToFile(HDFS_PATH,t);
    }
    @Override
    public void open(Configuration config) {
        hdfsManager = new HDFSManager();
    }
}

```

# 以上代码为简易版经过修改，思路清晰，（Flink自带一种HDFS的链接方式详情见链接：）

* https://ci.apache.org/projects/flink/flink-docs-release-1.4/dev/connectors/filesystem_sink.html