---
layout: post
categories: vgs lve
title: LVM 管理之一：扩容VG/LV  
date: 2017-12-27 19:40:23 +0800
description: LVM 管理之一：扩容VG/LV  
keywords: vgs 扩容vgs
---

#动态扩容 VG

* 1.1 查看硬盘信息

```bash
[root@pgb lvm]# fdisk -l

Disk /dev/hda: 19.3 GB, 19327352832 bytes
255 heads, 63 sectors/track, 2349 cylinders
Units = cylinders of 16065 * 512 = 8225280 bytes

   Device Boot      Start         End      Blocks   Id  System
/dev/hda1   *           1          13      104391   83  Linux
/dev/hda2              14         144     1052257+  82  Linux swap / Solaris
/dev/hda3             145        2349    17711662+  83  Linux

Disk /dev/hdb: 2147 MB, 2147483648 bytes
16 heads, 63 sectors/track, 4161 cylinders
Units = cylinders of 1008 * 512 = 516096 bytes

   Device Boot      Start         End      Blocks   Id  System
/dev/hdb1               1        1985     1000408+  83  Linux
/dev/hdb2            1986        4161     1096704   83  Linux

Disk /dev/hdd: 1073 MB, 1073741824 bytes
16 heads, 63 sectors/track, 2080 cylinders
Units = cylinders of 1008 * 512 = 516096 bytes

   Device Boot      Start         End      Blocks   Id  System
/dev/hdd1               1        2080     1048288+  83  Linux


[root@pgb lvm]# pvscan
  PV /dev/hdb1   VG vg01_pgdata   lvm2 [976.00 MB / 972.00 MB free]
  PV /dev/hdd1   VG vg01_pgdata   lvm2 [1020.00 MB / 0    free]
  Total: 2 [1.95 GB] / in use: 2 [1.95 GB] / in no VG: 0 [0   ]  
  
  #备注：根据 fdisk 和 pvscan 命令输出，知道 /dev/hdb2 还没有加入 VG， 可以使用，接下来将 /dev/hdb2 加入  VG vg01_pgdata。
  
```        
        
* 1.2 查看 VG 信息

```bash

[root@pgb lvm]# vgdisplay
  --- Volume group ---
  VG Name               vg01_pgdata
  System ID             
  Format                lvm2
  Metadata Areas        2
  Metadata Sequence No  2
  VG Access             read/write
  VG Status             resizable
  MAX LV                0
  Cur LV                1
  Open LV               1
  Max PV                0
  Cur PV                2
  Act PV                2
  VG Size               1.95 GB
  PE Size               4.00 MB
  Total PE              499
  Alloc PE / Size       256 / 1.00 GB
  Free  PE / Size       243 / 972.00 MB
  VG UUID               B5pg8R-2AGm-6DEp-K7HK-TI1I-HC3h-gWx32m

```
  
  
* 1.3 格式化文件系统

```bash
[root@pgb lvm]# mkfs -t ext3 -c /dev/hdb2
mke2fs 1.39 (29-May-2006)
Filesystem label=
OS type: Linux
Block size=4096 (log=2)
Fragment size=4096 (log=2)
137088 inodes, 274176 blocks
13708 blocks (5.00%) reserved for the super user
First data block=0
Maximum filesystem blocks=281018368
9 block groups
32768 blocks per group, 32768 fragments per group
15232 inodes per group
Superblock backups stored on blocks: 
        32768, 98304, 163840, 229376

Checking for bad blocks (read-only test): done                                
Writing inode tables: done                            
Creating journal (8192 blocks): done
Writing superblocks and filesystem accounting information: done

This filesystem will be automatically checked every 28 mounts or
180 days, whichever comes first.  Use tune2fs -c or -i to override.
```
         

* 1.4 创建PV    

```bash
[root@pgb lvm]# pvcreate /dev/hdb2
  Physical volume "/dev/hdb2" successfully created  
```

  
  
* 1.5 在线扩容 VG

```bash
[root@pgb lvm]# vgs
  VG          #PV #LV #SN Attr   VSize VFree  
  vg01_pgdata   2   1   0 wz--n- 1.95G 972.00M

[root@pgb lvm]# vgextend vg01_pgdata /dev/hdb2
  Volume group "vg01_pgdata" successfully extended
  
[root@pgb lvm]# vgs
  VG          #PV #LV #SN Attr   VSize VFree
  vg01_pgdata   3   1   0 wz--n- 2.99G 1.99G

```

* 1.6 再次查看 VG，查看是否扩容

```bash

[root@pgb lvm]# vgdisplay
  --- Volume group ---
  VG Name               vg01_pgdata
  System ID             
  Format                lvm2
  Metadata Areas        3
  Metadata Sequence No  3
  VG Access             read/write
  VG Status             resizable
  MAX LV                0
  Cur LV                1
  Open LV               1
  Max PV                0
  Cur PV                3
  Act PV                3
  VG Size               2.99 GB
  PE Size               4.00 MB
  Total PE              766
  Alloc PE / Size       256 / 1.00 GB
  Free  PE / Size       510 / 1.99 GB
  VG UUID               B5pg8R-2AGm-6DEp-K7HK-TI1I-HC3h-gWx32m    
  
  #  备注：现在 vg01_pgdata 大小为 3 GB 左右，已成功扩容 1 GB。
```
   

# 动态扩容 LV

      * 目标给已在线上使用的LV 扩容，在以下例子中，给目录 /database/pgdata1 扩容 512 M。

* 2.1 查看目录使用情况

```bash
[root@pgb lvm]# df -hv
Filesystem            Size  Used Avail Use% Mounted on
/dev/hda3              17G  9.8G  5.8G  64% /
/dev/hda1              99M   18M   76M  20% /boot
tmpfs                 217M     0  217M   0% /dev/shm
none                  217M  104K  217M   1% /var/lib/xenstored
/dev/mapper/vg01_pgdata-lv_pgdata1
                     1008M   34M  924M   4% /database/pgdata1  
```
 
* 2.2 查看所属 VG 信息

```bash
[root@pgb lvm]# vgdisplay
  --- Volume group ---
  VG Name               vg01_pgdata
  System ID             
  Format                lvm2
  Metadata Areas        3
  Metadata Sequence No  3
  VG Access             read/write
  VG Status             resizable
  MAX LV                0
  Cur LV                1
  Open LV               1
  Max PV                0
  Cur PV                3
  Act PV                3
  VG Size               2.99 GB
  PE Size               4.00 MB
  Total PE              766
  Alloc PE / Size       256 / 1.00 GB
  Free  PE / Size       510 / 1.99 GB
  VG UUID               B5pg8R-2AGm-6DEp-K7HK-TI1I-HC3h-gWx32m
  
#备注：从上面看出，VG vg01_pgdata 最大可用空间为 2.99 GB, 目前已分配 1 GB，还剩余 1.99 GB 可以分配。
```
  
  
* 2.3 增加 LV 大小 

```bash
[root@pgb lvm]#  lvextend  -L +512M /dev/mapper/vg01_pgdata-lv_pgdata1 
  Extending logical volume lv_pgdata1 to 1.50 GB
  Logical volume lv_pgdata1 successfully resized
 

[root@pgb lvm]# df -hv
Filesystem            Size  Used Avail Use% Mounted on
/dev/hda3              17G  9.9G  5.8G  64% /
/dev/hda1              99M   18M   76M  20% /boot
tmpfs                 217M     0  217M   0% /dev/shm
none                  217M  104K  217M   1% /var/lib/xenstored
/dev/mapper/vg01_pgdata-lv_pgdata1
                     1008M   34M  924M   4% /database/pgdata1
                     
  #备注： LV 扩容成功，但目录 /database/pgdata1 大小仍然为 1G，没有变化。还需要 resize2fs 命令处理下。


```
         
         
 * 2.4 resize2fs

```bash
[root@pgb lvm]# resize2fs -f /dev/mapper/vg01_pgdata-lv_pgdata1
resize2fs 1.39 (29-May-2006)
Filesystem at /dev/mapper/vg01_pgdata-lv_pgdata1 is mounted on /database/pgdata1; on-line resizing required
Performing an on-line resize of /dev/mapper/vg01_pgdata-lv_pgdata1 to 524288 (4k) blocks.
The filesystem on /dev/mapper/vg01_pgdata-lv_pgdata1 is now 524288 blocks long.
```


* 2.5 再次查看

```bash
[root@pgb lvm]# df -hv
Filesystem            Size  Used Avail Use% Mounted on
/dev/hda3              17G  9.8G  5.8G  64% /
/dev/hda1              99M   18M   76M  20% /boot
tmpfs                 217M     0  217M   0% /dev/shm
none                  217M  104K  217M   1% /var/lib/xenstored
/dev/mapper/vg01_pgdata-lv_pgdata1
                      1.5G   34M  1.4G   3% /database/pgdata1      
``` 
                      
  备注：目录  /database/pgdata1  空间果然变大了。    



转载请注明出处，本文采用 [CC4.0](http://creativecommons.org/licenses/by-nc-nd/4.0/) 协议授权
