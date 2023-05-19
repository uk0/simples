---
layout: post
categories: awk
title: linux硬盘检查大小awk脚本全自动执行,以及日志抽取
date: 2018-04-04 17:08:52 +0800
description: 硬盘监测
keywords: awk linux 硬盘检查
---


## DISK 监测剩余量


```bash
#!/bin/bash
temp=0.0
DISK_SIZE=0.0
for i in {0..26}
do
    echo -e "dscn0$i Disk"
    echo "剩余可用大小GB"
    temp=`ssh dscn0$i df -m|grep hadoop | awk '$4~/^[0-9]/ {split($4,array,"[A-Z]");b+=array[1]} END {print b/1024}'`
    DISK_SIZE=$(echo "$DISK_SIZE + $temp" | bc)

    echo -e "总数 剩余 占比 目录"
    ssh dscn0$i df -h | grep hadoop  | awk '{print$2, $4, $5, $6}'
done
echo -e "\033[32m 剩余总大小: $DISK_SIZE GB \033[0m"

```
## bash 样式

```bash

    PS1='\[\e[1;32m\]\u@\h\[\e[m\]:\[\e[1;34m\]\W\[\e[1;33m\]\$\[\e[m\] '
    
```

##  抽取一天的日志到另一个日志文件

```bash
#!/bin/bash

TIME_LINE=$(date +%Y%m%d%H%M%S)


function getOne(){
    # 设定变量
    if [ -z $1 ];then
        echo -e "\033[31m Error: log Name is Null  \033[0m"
        echo -e "\033[32m Simple : ./readlogs.sh getOne xxx.log  \033[0m"
        exit 1
    fi
    log=$1
    s=`date -d yesterday +%F` #前一天
    #s=$(date +%F) #今天
    e=$(date +%F) #今天
    # 根据条件获得开始和结束的行号
    sl=`cat -n $log | grep $s | head -1 | sed 's/^[ \t]*//g' | cut -f 1`
    el=`cat -n $log | grep $e | tail -1 | sed 's/^[ \t]*//g' | cut -f 1`
    echo -e "\033[32m start Line $sl \033[0m"
    echo -e "\033[32m end Line $el \033[0m"
    # 获取结果并输出到 `date`.log 文件
    sed -n "$sl,$el"'p' $log >> check_auto_$TIME_LINE.log
}

case $1 in

"getOne")
	getOne $2
	;;
*)
	echo -e "\033[1m usage: \n \t  [getOne] \033[0m"
	echo -e "\033[1m DESC: \n \t getOne : 获得一天内的文件日志，并且进行导出。 \033[0m"
	exit 2 # Command to come out of the program with status 1
	;;
esac
exit 0

```

转载请注明出处，本文采用 [CC4.0](http://creativecommons.org/licenses/by-nc-nd/4.0/) 协议授权
