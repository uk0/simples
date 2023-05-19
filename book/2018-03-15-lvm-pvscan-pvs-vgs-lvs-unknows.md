---
layout: post
categories: linux
title: LVM unknown device how to recover?
date: 2018-03-15 14:36:05 +0800
description: I have a server with two hard drives that I thought I had correctly installed with LVM, until I discovered that the second hard-drive did not actually seem to be used. I investigated the problem and followed some instructions found online, but the problem got worse. Apparently, my initial mistake was to remove the physical volume with pvremove when I should have used mvreduce.
keywords: lvm vgs pv lv unknowsDevice
---


## LV扩容


   ### 过程如下：

   *  挂载一块硬盘4T

        ```bash
            fdisk /dev/sdb

            send n

            send p

            send w;

            mkfs -t ext4 -c /dev/sdb1

            wrating....


            pvcreate /dev/sdb1


            vgextend rhel /dev/sdb1


            lvextend -L +2000G  /dev/mapper/lvs-home

            #根据文件格式进行重置home大小

            xfs_growfs /dev/mapper/lvs-home

                or

            resize2fs -p /dev/mapper/lvs-home

            其中遇到的问题如下：

            创建完成PV-->mkfs(格盘)-->vgs(warning Device with uuid XXXXXX-xxxxxx-xxx-xxxx-xxxxx)

            解决办法：
            pvcreate --uuid "UUID" --restorefile /etc/lvm/archie/rhel_000001-1511980835.vg  /dev/sdb1 -ff

            如果解决不了用第二种方式：
            vgreduce --removemissing --verbose myVG_NAME

        ```


转载请注明出处，本文采用 [CC4.0](http://creativecommons.org/licenses/by-nc-nd/4.0/) 协议授权
