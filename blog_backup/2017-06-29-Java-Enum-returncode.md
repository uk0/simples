---
layout: post
title: Java RESTfull状态码以及返回码Enum实现[学习]
categories: Java, RESTfull,Enum
description: Java, RESTfull,Enum
keywords: Java,Enum
---

Java RESTfull 开发中利用枚举类来记录错误码和错误码对应的错误信息，效果很棒。




## 1. Java中枚举类可以用于记录错误码和错误码对应的错误信息，其实现的技巧如下：
```java
package todcloud.utils.Enum;

/**
 * Created by zhangjianxin on 2017/6/29.
 */
public enum ResponseInfo {
    SERVER_ERROR("返回码异常", "6666"),
    SUCCESS("success", "200"), ERROR("Not Fount", "404"), GATEWAY("Fuck Server", "502"), CREATE("Created", "201");


    private String code;
    private String msg;
    private ResponseInfo(String code, String msg) {
        this.code = code;
        this.msg = msg;
    }
    public String getReturnCode() {
        return code;
    }
    public String getReturnMsg() {
        return msg;
    }

    public static String getResponseMsg(String code){
        for(ResponseInfo responseInfo:ResponseInfo.values()){
            if(code.equals(responseInfo.getReturnCode())){
                return responseInfo.getReturnMsg();
            }
        }
        return SERVER_ERROR.getReturnMsg();
    }
}

```

## 2. 调用测试

```java
package todcloud.utils.Enum;

/**
 * Created by zhangjianxin on 2017/6/29.
 */
public class ResponseInfoTest {
    public static void main(String[] args) {
        System.out.println(ResponseInfo.SUCCESS.getReturnCode());
        String code = ResponseInfo.SUCCESS.getReturnCode();
        System.out.println(ResponseInfo.getResponseMsg(code));
    }
}


```

> 通过复习Java中的Enum来完成这个功能。

## 3.失败案例

```java
package todcloud.utils.Enum;

/**
 * Created by zhangjianxin on 2017/6/29.
 */
public enum ReturnCode {
    SUCCESS("success", 200), ERROR("Not Fount", 404), GATEWAY("Fuck Server", 502), CREATE("Created", 201);
     /** 成员变量 */
    private String CodeData;
    private int CodeIndex;
    /**     构造方法*/
    private ReturnCode(String CodeData, int CodeIndex) {
        this.CodeData = CodeData;
        this.CodeIndex = CodeIndex;
    }
     /** 普通方法*/
    public static String getCodeData(int CodeIndex) {
        for (ReturnCode c : ReturnCode.values()) {
            if (c.getCodeIndex() == CodeIndex) {
                return c.CodeData;
            }
        }
        return null;
    }
     /** get set 方法*/
    public String getCodeData() {
        return CodeData;
    }
    public void setCodeData(String CodeData) {
        this.CodeData = CodeData;
    }
    public int getCodeIndex() {
        return CodeIndex;
    }
    public void setCodeIndex(int CodeIndex) {
        this.CodeIndex = CodeIndex;
    }
}
```


```java
package todcloud.utils.Enum;

import com.alibaba.fastjson.JSONObject;

import javax.print.attribute.standard.MediaSize;

/**
 * Created by zhangjianxin on 2017/6/29.
 */
public enum ReturnCode2 implements ReturnCodeIns{
    SUCCESS("success", 200), ERROR("Not Fount", 404), GATEWAY("Fuck Server", 502), CREATE("Created", 201);
    /** 成员变量*/
    private String name;
    private int index;
     /** 构造方法*/
    private ReturnCode2(String name, int index) {
        this.name = name;
        this.index = index;
    }
     /**接口方法*/
    @Override
    public JSONObject getInfo() {
        return this.QueryData(this.name);
    }
     /**接口方法*/
    @Override
    public void print() {
        System.out.println(this.index+":"+this.name);
    }

    public static JSONObject QueryData(String Name){
        System.out.println(Name);
        return null;
    }
}
```


```java
package todcloud.utils.Enum;

import com.alibaba.fastjson.JSON;
import com.alibaba.fastjson.JSONObject;

/**
 * Created by zhangjianxin on 2017/6/29.
 */
public interface ReturnCodeIns {
    void print();
    JSONObject getInfo();
}

```


```java
package todcloud.utils.Enum;

/**
 * Created by zhangjianxin on 2017/6/29.
 */
public class ReturnCodeMain {

    public static void main(String[] args) {
        System.out.println(ReturnCode.getCodeData(200));

    }
}

```
