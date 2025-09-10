---
layout: post
title:  hadoop HDFS常用文件操作命令
categories: shell,hadoop,hdfs
description: 笔记
keywords: shell,hdfs
---


工作当中接触到了 hadoop HDFS,操作命令有些不熟悉，所以记录下来以备不时之需。


## 命令基本格式:

```bash
hadoop fs -cmd < args >
```
# 1.ls


 * 列出hdfs文件系统根目录下的目录和文件

```bash
hadoop fs -ls  /

```

 * 列出hdfs文件系统所有的目录和文件

```bash
hadoop fs -ls -R /
```



# 2.put

 * hdfs file的父目录一定要存在，否则命令不会执行
```bash
hadoop fs -put < local file > < hdfs file >

```

 * hdfs dir 一定要存在，否则命令不会执行

```bash
hadoop fs -put  < local file or dir >...< hdfs dir >
```

 * 从键盘读取输入到hdfs file中，按Ctrl+D结束输入，hdfs file不能存在，否则命令不会执行

```bash
hadoop fs -put - < hdsf  file>
```

# 2.1.moveFromLocal
 * 与put相类似，命令执行后源文件 local src 被删除，也可以从从键盘读取输入到hdfs file中

```bash
hadoop fs -moveFromLocal  < local src > ... < hdfs dst >
```

# 2.2.copyFromLocal
 * 与put相类似，也可以从从键盘读取输入到hdfs file中

```bash
hadoop fs -copyFromLocal  < local src > ... < hdfs dst >
```

# 3.get

 * local file不能和 hdfs file名字不能相同，否则会提示文件已存在，没有重名的文件会复制到本地

```bash
hadoop fs -get < hdfs file > < local file or dir>
```
 * 拷贝多个文件或目录到本地时，本地要为文件夹路径
   `注意`：如果用户不是root， local 路径要为用户文件夹下的路径，否则会出现权限问题

```bash
hadoop fs -get < hdfs file or dir > ... < local  dir >
```

# 3.1.copyToLocal

```bash
hadoop fs -copyToLocal < local src > ... < hdfs dst >
```

# 4.rm
 * 每次可以删除多个文件或目录

```bash
hadoop fs -rm < hdfs file > ...
hadoop fs -rm -r < hdfs dir>...
```

# 5.mkdir

```bash
# 只能一级一级的建目录，父目录不存在的话使用这个命令会报错
hadoop fs -mkdir < hdfs path>
# 所创建的目录如果父目录不存在就创建该父目录
hadoop fs -mkdir -p < hdfs path>
```

# 6.getmerge

```bash
# 将hdfs指定目录下所有文件排序后合并到local指定的文件中，文件不存在时会自动创建，文件存在时会覆盖里面的内容
hadoop fs -getmerge < hdfs dir >  < local file >

# 加上nl后，合并到local file中的hdfs文件之间会空出一行

hadoop fs -getmerge -nl  < hdfs dir >  < local file >
```

# 7.cp

```bash
 hadoop fs -cp  < hdfs file >  < hdfs file >
# 目标文件不能存在，否则命令不能执行，相当于给文件重命名并保存，源文件还存在
 hadoop fs -cp < hdfs file or dir >... < hdfs dir >
# 目标文件夹要存在，否则命令不能执行
```

# 8.mv

```bash
hadoop fs -mv < hdfs file >  < hdfs file >
# 目标文件不能存在，否则命令不能执行，相当于给文件重命名并保存，源文件不存在
hadoop fs -mv  < hdfs file or dir >...  < hdfs dir >

# 源路径有多个时，目标路径必须为目录，且必须存在。
#  `注意`：跨文件系统的移动（local到hdfs或者反过来）都是不允许的
```

# 9.count

```bash
# 统计hdfs对应路径下的目录个数，文件个数，文件总计大小
# 显示为目录个数，文件个数，文件总计大小，输入路径
hadoop fs -count < hdfs path >
```

# 10.du
* 显示hdfs对应路径下每个文件夹和文件的大小
```bash
hadoop fs -du < hdsf path>
```
* 显示hdfs对应路径下所有文件和的大小
```bash
hadoop fs -du -s < hdsf path>
```
* 显示hdfs对应路径下每个文件夹和文件的大小,文件的大小用方便阅读的形式表示，例如用64M代替67108864
```bash
hadoop fs -du - h < hdsf path>
```

# 11.text
* 将文本文件或某些格式的非文本文件通过文本格式输出
```bash
hadoop fs -text < hdsf file>
```

# 12.setrep
* 改变一个文件在hdfs中的副本个数，上述命令中数字3为所设置的副本个数，-R选项可以对一个人目录下的所有目录+文件递归执行改变副本个数的操作

```bash
hadoop fs -setrep -R 3 < hdfs path >
```

# 13.stat
* 返回对应路径的状态信息
  `[format]`可选参数有：`%b`（文件大小），`%o`（Block大小），`%n`（文件名），`%r`（副本个数），`%y`（最后一次修改日期和时间）
  可以这样书写`hadoop fs -stat` %b%o%n < hdfs path >，不过不建议，这样每个字符输出的结果不是太容易分清楚

```bash
hdoop fs -stat [format] < hdfs path >
```

# 14.tail
* 在标准输出中显示文件末尾的`1KB`数据
```bash
hadoop fs -tail < hdfs file >
```

# 15.archive


```bash
hadoop archive -archiveName name.har -p < hdfs parent dir > < src >* < hdfs dst >
```

* 命令中参数name：压缩文件名，自己任意取；< hdfs parent dir > ：压缩文件所在的父目录；`< src >`：要压缩的文件名；`< hdfs dst >`：压缩文件存放路径
  *示例：`hadoop archive` -archiveName hadoop.har -p /user 1.txt 2.txt /des
  示例中将hdfs中/user目录下的文件`1.txt`，`2.txt`压缩成一个名叫`hadoop.har`的文件存放在`hdfs`中/des目录下，如果1.txt，2.txt不写就是将/user目录下所有的目录和文件压缩成一个名叫hadoop.har的文件存放在hdfs中/des目录下
  显示har的内容可以用如下命令：

```bash
hadoop fs -ls /des/flink-des.jar
```

* 显示har压缩的是那些文件可以用如下命令

```bash
hadoop fs -ls -R har:///des/hadoop.har
```
* 注意：`har`文件不能进行二次压缩。如果想给.har加文件，只能找到原来的文件，重新创建一个。har文件中原来文件的数据并没有变化，har文件真正的作用是减少`NameNode`和`DataNode`过多的空间浪费。