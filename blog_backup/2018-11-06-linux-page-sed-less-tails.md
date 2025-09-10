---
layout: post
categories: linux
title: Linux awk less tail grep 使用
date: 2018-11-06 17:00:35 +0800
description: 学习使用linux小工具
keywords: awk sed grep 日志
---


## Linux awk less grep sed 等命令使用

* shell demo

```bash

set time = 12:34:56
set hr = `echo $time | awk '{split($0,a,":" ); print a[1]}'` # = 12

set sec = `echo $time | awk '{split($0,a,":" ); print a[3]}'` # = 56

set hms = `echo $time | awk '{split($0,a,":" ); print a[1], a[2], a[3]}'`# = 12 34 56

# 获得5 - 10 line 并且用 `;` 分隔每一行  获得第个元素
sed -n '5,10p' xvideos.com-db.csv | awk '{split($0,a,";" ); print a[1]}'

sed -n '5,10p' xvideos.com-db.csv | awk '{split($0,a,";" ); print a[1] a[2]}'

#从第3000行开始，显示1000行。即显示3000~3999行
cat filename | tail -n +3000 | head -n 1000

#显示1000行到3000行

cat filename| head -n 3000 | tail -n +1000 

tail -n 1000 #：显示最后1000行

tail -n +1000 #：从1000行开始显示，显示1000行以后的

head -n 1000 #：显示前面1000行

tail -400f demo.log #监控最后400行日志文件的变化 等价与 tail -n 400 -f （-f参数是实时）

less demo.log #查看日志文件，支持上下滚屏，查找功能

uniq -c demo.log  #标记该行重复的数量，不重复值为1

grep 'INFO' demo.log     #在文件demo.log中查找所有包行INFO的行

grep -o 'order-fix.curr_id:\([0-9]\+\)' demo.log    #-o选项只提取order-fix.curr_id:xxx的内容（而不是一整行），并输出到屏幕上
grep -c 'ERROR' demo.log   #输出文件demo.log中查找所有包行ERROR的行的数量

# 输出demo.log中的某个日期中的ERROR的行
sed -n '/^2011-08-23.*ERROR/p' demolog.log

#指定执行的sed文件
sed -f demo.sed2 demo.log
```
* demo.sed2

```bash
#n   #这一行用法和命令中的-n一样意思，就是默认不输出
#demo.sed2
#下面的一行是替换指令，就是把19位长的日期和INFO/ERROR,id,和后面的一截提取出来，然后用@分割符把这4个字段重新按顺序组合
s/^\([-\: 0-9]\{19\}\).*\(INFO\|ERROR\) .*order-fix.curr_id:\([0-9]\+\),\(.*$\)/\1@\3@\2@\4/p


#排序功能 -t表示用@作为分割符，-k表示用分割出来的第几个域排序(不要漏掉后面的,2/,3/,1，详细意思看下面的参考链接，这里不做详述)
sed -f test.sed demolog.log | sort -t@ -k2,2n -k3,3r -k1,1  #n为按数字排序，r为倒序


awk 'BEGIN{FS="@"} {print $2,$3}' demo.log_after_sort   #BEGIN中预处理的是，把@号作为行的列分割符,把分割后的行的第2，3列输出

```

* 对指定时间范围内的日志进行统计，包括输出INFO，ERROR总数，记录总数，每个订单记录分类统计

```bash
sed -f demo.sed demolog.log | sort -t@ -k2,2n -k3,3r -k1,1 | awk -f demo.awk

```

* demo.awk

```bash
#下面的例子是作为命令行输入的，利用单引号作为换行标记，这样就不用另外把脚本写进文件调用了
awk '
BEGIN {
FS="@"
}

{
if ($3 == "INFO") {info_count++}
if ($3 == "ERROR") {error_count++}

}

END {
print "order total count:"NR           #NR是awk内置变量，是遍历的当前行号，到了END区域自然行号就等于总数了
printf("INFO count:%d ERROR count:%d\n",info_count,error_count)
} ' demo.log_after_sort

```


```bash
ll -lrth #:按照更改时间倒序排列，最新文件在下边

ll -lrSh #:按照文件大小倒序排列，最大文件在下边
grep --color # :高亮查询关键字
```
* 在大多数情况下` awk` 的 `print` 语句可以完成任务，但有时我们还需要更多。在那些情况下，awk 提供了两个我们熟知的函数 `printf()` 和 `sprintf()`。是的，如同其它许多 awk 部件一样，这些函数等同于相应的` C 语言函数`。`printf()` 会将格式化字符串打印到 stdout，而 sprintf()函数返回根据`printf`格式说明指定的格式化的字符串，它格式化数据但不输出数据。a w k提供函数`printf`，拥有几种不同的格式化输出功能。例如按列输出、左对齐或右对齐方式。
* `printf()`函数基本语法是`printf()`（`格式控制符`，`参数`） ，格式控制字符通常在引号里。类似C语言，awk printf格式有如下：
    ```c
    %c //ASCII字符
    %d //整数
    %e //浮点数，科学记数法
    %f //浮点数，例如（1 2 3 . 4 4）
    %g //awk决定使用哪种浮点数转换 e或者f
    %o //八进制数
    %s //字符串
    %x //十六进制数
    ```
* 下面来试试这些个格式：
```bash
echo 97 | awk '{printf("%c\n", $0)}'
A
```
* 类似`C语言`的格式化输出：
```bash
awk 'BEGIN{FS=":"}{printf("%-15s%s\n", $1, $3)}' group_file2
wireshark    987
usbmon       986
jackuser     985
vboxusers    984
aln         1001
```
* `sprintf`()函数返回根据`printf`格式说明指定的格式化的字符串，它格式化数据但不输出数据。所以需要将`sprintf`返回的数据保存在变量里面再输出

```bash
awk 'BEGIN{FS=":";ORS=""}{var=sprintf("%s\n", $1);print var}' group_file2
wireshark
usbmon
jackuser
vboxusers
aln
```


*  清除 Javascript 脚本里面的所有console(包含即删除整行)

```bash
 sed -i "" "/console/d" app.txt # Mac （Mac默认会要求你操作源文件的时候备份）

# 可以配合 find 使用
 find . -name "*.js"
 sed -i "/console/d" app.txt # Linux

#Mac 

find . -name "*.js" | awk '{print$1}' | xargs -L1  -I NAME sed -i ""  "/console/d" NAME

# 如果失败或者遇到问题
git checkout . && git clean -xdf


```



一个Linux 大佬 https://blog.csdn.net/imxiangzi/article/details/50387073

转载请注明出处，本文采用 [CC4.0](http://creativecommons.org/licenses/by-nc-nd/4.0/) 协议授权
