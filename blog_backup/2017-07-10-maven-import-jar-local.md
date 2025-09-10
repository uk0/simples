---
layout: post
title:  Maven引入本地Jar包进行编译项目 [来自项目中的经历]
categories: Idea,java
description: 回顾
keywords: Idea, java，Mavne
---


发这个贴的原因，是因为 是因为 是因为啥来着，忘了 那就是因为你，因为你 。


## 摘要: 记录下自己遇到的问题，和解决办法

```text
  先上图，毕竟是无图无真相的时代：
```
![](http://112firshme11224.test.upcdn.net/blog/tmp/maven-import-jar.png)


![](http://112firshme11224.test.upcdn.net/blog/tmp/maven-import-jar-2.png)

* 上图已经标注出两个需要修改的位置。

* `maven clean package` 会加载到这个jar包。用普通的引入方式`Maven`打包会失败。