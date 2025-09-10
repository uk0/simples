---
layout: post
title:  Mongodb作为Flink 的 Sink
categories: Flink,Mongodb
description: 笔记
keywords: Mongodb,flink
---


因工作需求所整合Flink + Mongodb 作为一个Demo 帮助大家跳坑。
Mongodb 采用集群主从模式。

# MongoSink.Java
   * 继承RichSinkFunction，重写`open` 和 `invoke`方法，还有`close`。


```java
package com.e.firsh.spb.sink;

import com.alibaba.fastjson.JSON;
import com.e.firs.spb.utils.Mongo.MongoService;
import org.apache.flink.configuration.Configuration;

import org.apache.flink.streaming.api.functions.sink.RichSinkFunction;

import java.net.UnknownHostException;

/**
 * Created by zhangjianxin on 2017/7/31.
 */
public  class MongoSink extends RichSinkFunction<String> {
    public static  String CollectionName = "collection-a";
    private MongoService mongoService;
    @Override
    public void invoke(String t) {
        try {
            this.mongoService.saveJson(JSON.parseObject(t),CollectionName);
        } catch (UnknownHostException e) {
            e.printStackTrace();
        }
    }
    @Override
    public void open(Configuration config) {
        mongoService = new MongoService();
        try {
            super.open(config);
        } catch (Exception e) {
            e.printStackTrace();
        }
    }
};
```
# MongoService.Java
   * MongoDB Service 处理工具

```java
package com.e.firsh.spb.utils.Mongo;

import com.alibaba.fastjson.JSONObject;

import java.net.UnknownHostException;
import java.util.List;

/**
 * Created by zhangjianxin on 2017/7/31.
 */
public class MongoService {



    private EntryDao EntryDao = new EntryDao();

    /**
     * 保存
     * @param saveJson
     * @throws UnknownHostException
     */
    public void saveJson(JSONObject saveJson, String CollectionName) throws UnknownHostException {
        this.EntryDao.save(CollectionName,saveJson,true);
    }




    /**
     * 更新
     * @param Entry
     * @throws UnknownHostException
     */
    public void update(Entry Entry,String CollectionName) throws UnknownHostException {
        this.EntryDao.update(Entry,CollectionName);
    }

    /**
     * 查询所有
     * @return
     * @throws UnknownHostException
     */
    public List<Entry> findAll(String CollectionName) throws UnknownHostException{
        return this.EntryDao.findAll(CollectionName);
    }


    /**
     * 删除操作
     * @param id
     * @throws UnknownHostException
     */
    public void remove(int id,String CollectionName) throws UnknownHostException{
        this.EntryDao.remove(id,CollectionName);
    }
}

```


# MongoManager.Java
   * 获得mongodb链接池的工具类，获得某个DB的链接。

```java
package com.e.firsh.spb.utils.Mongo;

/**
 * Created by zhangjianxin on 2017/7/31.
 */
import com.mongodb.*;
import com.mongodb.client.MongoDatabase;

import java.util.Arrays;

public class MongoManager {

    private static Mongo mongo = null;
    private static MongoClient mongoClient= null;
    private MongoManager() { }

    static {
        initDBPrompties();
    }

    public static DB getDB(String dbName) {
        return mongo.getDB(dbName);
    }
    public static MongoDatabase getDatabase(String dbName) {
        return mongoClient.getDatabase(dbName);
    }

    /**
     * 初始化连接池
     */
    private static void initDBPrompties() {
        try {

            MongoClientOptions.Builder mcob = MongoClientOptions.builder();
            mcob.connectionsPerHost(1000);
            mcob.socketKeepAlive(true);
            mcob.readPreference(ReadPreference.secondaryPreferred());
            MongoClientOptions mco = mcob.build();

            mongoClient = new MongoClient(Arrays.asList(
                    new ServerAddress("127.0.0.1", 27017),
                    new ServerAddress("127.0.0.1", 27017),
                    new ServerAddress("127.0.0.1", 27017)), mco);


            mongo = new Mongo(MongoDBConstant.HOST, MongoDBConstant.PORT);
            MongoOptions opt = mongo.getMongoOptions();
            opt.connectionsPerHost = MongoDBConstant.POOLSIZE;
            opt.threadsAllowedToBlockForConnectionMultiplier = MongoDBConstant.BLOCKSIZE;
        } catch (MongoException e) {
            e.printStackTrace();
        }
    }
}

```


# EntryDao.java
   * 用来操作实体的工具类

```java
package com.e.firsh.spb.utils.Mongo;
import java.net.UnknownHostException;
import java.util.ArrayList;
import java.util.List;
import java.util.UUID;

import com.alibaba.fastjson.JSONObject;
import com.mongodb.*;
import com.mongodb.client.FindIterable;
import com.mongodb.client.MongoCollection;
import com.mongodb.client.MongoCursor;
import com.mongodb.client.MongoDatabase;
import com.mongodb.util.JSON;
import org.bson.Document;

/**
 * Created by zhangjianxin on 2017/7/31.
 */
public class EntryDao {

    /**
     * 删除操作
     * @param id
     * @throws UnknownHostException
     */
    public void remove(int id,String CollectionName) throws UnknownHostException{
        MongoDatabase mongo = MongoManager.getDatabase(MongoDBConstant.DB);
        MongoCollection coll = mongo.getCollection(CollectionName);


        Document baseDBO =new Document();
        baseDBO.put("id", id);

        //删除某一条记录
        coll.deleteOne(baseDBO);
    }


    /**
     *
     * @param CollectionName 集合名
     * @param saveJson  待存入JSON
     * @param useDefaultId  未传入_id时，若为true则使用MongoDB自动生成的_id。若为false则生成UUID作为主键
     * @return
     */
    public JSONObject save(String CollectionName, JSONObject saveJson, boolean useDefaultId) {
        JSONObject resp = new JSONObject();
        try {
            MongoDatabase mongo = MongoManager.getDatabase(MongoDBConstant.DB);
            MongoCollection coll = mongo.getCollection(CollectionName);
            Object _id = null;
           if (!saveJson.containsKey("_id")){
                if (!useDefaultId) {
                    _id = UUID.randomUUID().toString();
                    saveJson.put("Data", _id);
                }
            } else if(saveJson.containsKey("_id"))  {
                _id = saveJson.get("_id");
            }
            Document doc = Document.parse(saveJson.toString());
            coll.insertOne(doc);
            resp.put("Data", _id);
        } catch (MongoTimeoutException e1) {
            e1.printStackTrace();
            resp.put("ReasonMessage",e1.getClass() + ":" + e1.getMessage());
            return resp;
        } catch (Exception e) {
            e.printStackTrace();
            resp.put("ReasonMessage",e.getClass() + ":" + e.getMessage());
        }
        return resp;
    }

}
```

# JsonUtils.java

* JSON的小工具类

```java
package com.tod.spb.utils.Mongo;

import java.io.OutputStream;

import org.codehaus.jackson.map.ObjectMapper;
import org.codehaus.jackson.type.TypeReference;


/**
 * Json转化工具，可以实现java对象和json字符串之间的互相转化<br />
 * Created by zhangjianxin on 2017/7/31.
 */
public class JsonUtils {
    static ObjectMapper objectMapper = new ObjectMapper();

    /**
     * java 对象转换为json 存入流中
     *
     * @param obj
     * @return  s
     */
    public static String toJson(Object obj) {
        String s = "";
        try {
            s = objectMapper.writeValueAsString(obj);
        } catch (Exception e) {
            e.printStackTrace();
        }
        return s;
    }

    /**
     * java 对象转换为json 存入流中
     *
     * @param obj
     * @param out
     */
    public static void toJson(Object obj, OutputStream out) {
        try {
            objectMapper.writeValue(out, obj);
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    /**
     * json 转为java对象
     *
     * @param json
     * @param obj
     */
    @SuppressWarnings({ "rawtypes", "unchecked" })
    public static void fromJson(String json, Object obj, Class valueType) {
        try {
            obj = objectMapper.readValue(json, valueType);
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    /**
     * json 转为java对象
     * @param json
     * @param obj
     * @param valueTypeRef
     */
    @SuppressWarnings("rawtypes")
    public static void fromJson(String json, Object obj, TypeReference valueTypeRef) {
        try {
            obj = objectMapper.readValue(json, valueTypeRef);
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    /**
     * json 转为java对象
     *
     * @param json
     * @return obj
     */
    @SuppressWarnings({ "rawtypes", "unchecked" })
    public static Object fromJson(String json, Class valueType) {
        Object obj = null;
        try {
            obj = objectMapper.readValue(json, valueType);
        } catch (Exception e) {
            e.printStackTrace();
        }
        return obj;
    }

    /**
     * json 转为java对象
     *
     * @param json
     * @return obj
     * @param valueTypeRef
     */
    @SuppressWarnings({ "rawtypes", "hiding" })
    public static <Object> Object fromJson(String json, TypeReference valueTypeRef) {
        Object obj = null;
        try {
            obj = objectMapper.readValue(json, valueTypeRef);
        } catch (Exception e) {
            e.printStackTrace();
        }
        return obj;
    }
}

```

# MongoDBHelper.java

 * 用来测试MongoDB链接成功

```bash
package com.tod.spb.utils.Mongo;
import com.alibaba.fastjson.JSON;
import com.alibaba.fastjson.JSONObject;
import com.mongodb.DB;
import com.mongodb.DBCollection;
import com.mongodb.DBObject;
import java.net.UnknownHostException;


/**
 * Created by zhangjianxin on 2017/7/31.
 */
public class MongoDBHelper {

    public static  String CollectionName = "collection-a";
    /**
     * 保存
     * @param null
     * @throws UnknownHostException
     */
    private MongoService mongoService = new MongoService();


    public void save() throws UnknownHostException{
        JSONObject data = new JSONObject();
        data.put("tableName","1");
        data.put("body","jsonobject");
        data.put("number","12345");
        this.mongoService.saveJson(data,CollectionName);
    }



    public static void main(String[] args) {

        try {
            MongoDBHelper mongoDBHelper = new MongoDBHelper();
            mongoDBHelper.save();
        } catch (UnknownHostException e) {
            e.printStackTrace();
        }
    }
}

```
* Owner `breakEval13`
* https://github.com/breakEval13