---
layout: post
title:  Mavne jar包去重以及升级,[来自项目中的经历]
categories: Mavne,java
description: 回顾
keywords: Mavne, java
---


发这个贴的原因，是因为 是因为 是因为啥来着，忘了 那就是因为你，因为你 。


## 错误描述

```text
  因为项目中用到ActorDB所以有了以下的故事，当maven中引入多个jar包，那么避免不了会发生重复的包存在。
 
```
* 如图所示，成功的倒入了·`org.apache.thrift`
  但是工具里面没有找到 `clearBit`，`setBit`等方法。

![](http://112firshme11224.test.upcdn.net/blog/QQ20170708-002142@2x.png)


* 经过查证，的确没有，那么为什么会发生这个事情。

 ![](http://112firshme11224.test.upcdn.net/blog/error.png)
 
 
* 揭开迷雾
    
 ![](http://112firshme11224.test.upcdn.net/blog/success-geterror.png)
 
 * 发现了存在的包竟然是MSF4J的包。
   于是找到了解决方案代码如下：
   
  ```xml
   <dependency>
      <groupId>org.wso2.msf4j</groupId>
      <artifactId>msf4j-all</artifactId>
      <version>2.1.0</version>
      <exclusions>
          <exclusion>
              <artifactId>libthrift</artifactId>
              <groupId>org.apache.thrift</groupId>
          </exclusion>
      </exclusions>
  </dependency>
  ```
  
 * 通过这种方式进行包的摘除，(结果还不是很理想，还是有问题。)
   * 于是仔细查看pom.xml文件的内容，发现了问题，怀疑Maven加载包的顺序是有优先级的(从上到下)
   * 于是将 `org.apache.thrift` 放到了MSF4J包的上面。
   * 成功解决问题：
   ![](http://112firshme11224.test.upcdn.net/blog/th3.png)
   
## 总结
   学习了Mavne的加载顺序，以及排错方式：
```bash
#查找依赖所在的地方
 mvn dependency:tree -Dverbose -Dincludes=org.apache.thrift:libthrift
```

```xml
<!--去除依赖，在外部加载-->
<exclusions>
      <exclusion>
          <artifactId>libthrift</artifactId>
          <groupId>org.apache.thrift</groupId>
      </exclusion>
  </exclusions>
```


