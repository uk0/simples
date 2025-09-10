---
layout: post
title: IDEA Java注释模板设置详解
categories: IDEA,java,mac
description: matching,mac
keywords: java, IDEA,mac
---

IDEA Java注释模板设置详解
设置注释模板的入口： Window->Preference->Java->Code Style->Code Template 然后展开Comments节点就是所有需设置注释的元素啦。现就每一个元素逐一介绍：



## 1、文件(Files)注释标签：

```java

/**   
* @Title: ${file_name}
* @Package ${package_name}
* @Description: ${todo}(用一句话描述该文件做什么)
* @author A18ccms A18ccms_gmail_com   
* @date ${date} ${time}
* @version V1.0   
*/
```

## 2.类型(Types)注释标签（类的注释）：

```java
/**
* @ClassName: ${type_name}
* @Description: ${todo}(这里用一句话描述这个类的作用)
* @author A18ccms a18ccms_gmail_com
* @date ${date} ${time}
* 
* ${tags}
*/
```

## 3.字段(Fields)注释标签：

```java
/**
* @Fields ${field} : ${todo}(用一句话描述这个变量表示什么)
*/
构造函数标签：

/**
* <p>Title: </p>
* <p>Description: </p>
* ${tags}
*/

```

## 4.方法(Constructor & Methods)标签：

```java
/**
* @Title: ${enclosing_method}
* @Description: ${todo}(这里用一句话描述这个方法的作用)
* @param ${tags}    设定文件
* @return ${return_type}    返回类型
* @throws
*/

```
## 5.覆盖方法(Overriding Methods)标签：
```java

/* (非 Javadoc)
* <p>Title: ${enclosing_method}</p>
* <p>Description: </p>
* ${tags}
* ${see_to_overridden}
*/

```

## 6. 代表方法(Delegate Methods)标签：
```java
/**
* ${tags}
* ${see_to_target}
*/
getter方法标签：

/**
* @return ${bare_field_name}
*/

setter方法标签：

/**
* @param ${param} 要设置的 ${bare_field_name}
*/

```

 