---
layout: post
categories: Linux
title: 一次磁盘挂载
date: 2019-02-18 19:19:53 +0800
description: 磁盘挂载Lvm
keywords: Linux
---

服务器10台每台16块盘 其中两块系统盘,剩下的皆为数据盘，因为涉及到扩容，所以有一个大目录采用lvm挂载。


### 一次磁盘挂载（Centos7）

```bash
#!/usr/bin/env bash
################## base ###############
bashPath=$(cd `dirname $0`; pwd)

${bashPath}/lvms.sh

sleep 2

${bashPath}/parted.sh
```

* lvms.sh

```bash
#!/bin/bash
yum install -y lvm2

## disk--->pv--->vg--->lv

disk_index=14                  # 第14块盘

partition=/data            # 定义最终挂载的名称

j=`echo $disk_index|awk '{printf "%c",97+$disk_index}'`

## 直接disk
fdisk /dev/sd$j << FORMAT
d
w
FORMAT

parted /dev/sd$j <<ESXU
mklabel gpt
yes
mkpart primary 0 100%
ignore
quit
ESXU

echo -e "\033[32m disk successfully lvm $j \033[0m"

mkdir -p $partition

pvcreate /dev/sd${j}1

vgcreate vg_data_${j}  /dev/sd${j}1

lvcreate  -l 100%FREE  -n  lv_data_${j}  vg_data_${j}

mkfs.ext4  /dev/vg_data_${j}/lv_data_${j}

mount="/dev/mapper/vg_data_${j}-lv_data_${j}      ${partition}  ext4    defaults        0       0"

echo $mount >>/etc/fstab                #写入分区表

sleep 1s

mount -a


```

* parted.sh

```bash
#!/bin/bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin
export PATH
i=1
while [ $i -lt 14 ]                  #硬盘数量,除系统盘之外是13块
do
j=`echo $i|awk '{printf "%c",97+$i}'` #系统盘是sda1,如果是其它的需要修改脚本97=a

unmount -v /dev/sd$j

parted /dev/sd$j <<FORMAT
mklabel gpt
mkpart primary 0 100%
ignore
quit
FORMAT
mkfs.ext4 -T largefile  /dev/sd${j}1    #格式化磁盘
mkdir -p /data/disk${i}
mount="/dev/sd${j}1       /data/disk${i}  ext4    defaults        0       0"
rm -rf /data/disk${i}/*
echo $mount >>/etc/fstab                #写入分区表
i=$(($i+1))
done
echo -e "\033[32m *****Formating  and Mounting have finished wait 5s **** \033[0m"

chmod -R 777 /data/*

sleep 5

mount -a
```



转载请注明出处，本文采用 [CC4.0](http://creativecommons.org/licenses/by-nc-nd/4.0/) 协议授权
