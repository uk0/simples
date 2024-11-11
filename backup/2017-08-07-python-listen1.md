---
layout: post
title:  使用python把图片存入数据库.
categories: Python,lambda
description: 笔记
keywords: Python,lambda
---


没事看看文章而已～

# 使用python把图片存入数据库

一般情况下我们是把图片存储在文件系统中，而只在数据库中存储文件路径的，但是有时候也会有特殊的需求：把图片二进制存入数据库。

今天我们采用的是python+mysql的方式

MYSQL 是支持把图片存入数据库的，也相应的有一个专门的字段 BLOB (Binary Large Object)，即较大的二进制对象

还有个更大的存二进制的LONGBLOB；

这里需要注意：尽量把字段设置大一些，因为如果设置的字段长度过小，就会出现图片只显示一部分的情况。第二：如果数据量大的话尽量避免使用这种方式进行，因为mysq
l对于大数据的查询速度会很慢。

下面上代码：


```python
#!/usr/bin/python
#-*- coding: UTF-8 -*-

import MySQLdb as mysql
import sys
try:
    #读取图片文件
    fp = open("./test.jpg")
    img = fp.read()
    fp.close()
except IOError,e:
    print "Error %d %s" % (e.args[0],e.args[1])
    sys.exit(1)
try:
    #mysql连接
    conn = mysql.connect(host='localhost',user='root',passwd='123456',db='test')
    cursor = conn.cursor()
    #注意使用Binary()函数来指定存储的是二进制
    cursor.execute("INSERT INTO images SET data='%s'" % mysql.Binary(img))
    #如果数据库没有设置自动提交，这里要提交一下
    conn.commit()
    cursor.close()
    #关闭数据库连接
    conn.close()
except mysql.Error,e:
    print "Error %d %s" % (e.args[0],e.args[1])
    sys.exit(1)

```




* 以上代码经过验证可以直接拿去修改调用
* Owner `breakEval13`
* https://github.com/breakEval13