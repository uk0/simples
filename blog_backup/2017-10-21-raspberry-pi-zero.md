---
layout: post
title: raspberry pi zero 安装系统，以及添加USB识别网卡模式[badUSB]
categories: python
description: raspberry
keywords: raspberry,python
---


一个用Python如何将自己的树莓派打造成“渗透测试神器”


date：2017-10-21
author：zhangjianxin

[TOC]

# 1 、 初始化树莓派zero。

[*] 第一步安装系统，

   * 有以下设备
      * raspberry pi zero
      * SSK 转接卡
      * 闪迪SD卡8GClass10

   ![](http://112firshme11224.test.upcdn.net/images/1508598613256.jpg)


* 第一次没有插入SD卡的时候

```bash
➜  ~ df -h
Filesystem                          Size   Used  Avail Capacity iused      ifree %iused  Mounted on
/dev/disk1                         112Gi   93Gi   18Gi    84% 1444634 4293522645    0%   /
devfs                              186Ki  186Ki    0Bi   100%     642          0  100%   /dev
map -hosts                           0Bi    0Bi    0Bi   100%       0          0  100%   /net
map auto_home                        0Bi    0Bi    0Bi   100%       0          0  100%   /home
localhost:/fMn-YFzQfnHjvi04JpRLvc  112Gi  112Gi    0Bi   100%       0          0  100%   /Volumes/MobileBackups

```



* 插入sd卡的时候（多了一个/dev/disk2s1 ）

![](http://112firshme11224.test.upcdn.net/images/df-h-sd-card.png)


```bash
➜  ~ df -h
Filesystem                          Size   Used  Avail Capacity iused      ifree %iused  Mounted on
/dev/disk1                         112Gi   93Gi   18Gi    84% 1444634 4293522645    0%   /
devfs                              186Ki  186Ki    0Bi   100%     642          0  100%   /dev
map -hosts                           0Bi    0Bi    0Bi   100%       0          0  100%   /net
map auto_home                        0Bi    0Bi    0Bi   100%       0          0  100%   /home
localhost:/fMn-YFzQfnHjvi04JpRLvc  112Gi  112Gi    0Bi   100%       0          0  100%   /Volumes/MobileBackups
/dev/disk2s1                       7.4Gi  2.4Mi  7.4Gi     1%      76     243827    0%   /Volumes/usb

```

* 通过diskutils工具进行umont这个SD卡(为了防止其他程序写入)

 ![](http://112firshme11224.test.upcdn.net/images/umount-mac.png)

    ```bash
    ➜  ~ diskutil unmountDisk /dev/disk2s1
    Unmount of all volumes on disk2 was successful
    ```


* 通过dd 写入镜像到sd卡。(10min)

 ![](http://112firshme11224.test.upcdn.net/images/dd-sd-car-success.png)

    ```bash
    ➜  ~ sudo dd bs=512  if=2017-02-16-raspbian-jessie.img  of=/dev/disk2s1
    ```

 > Ethernet Gadget模式即USB网卡模式，
   比较常见的就是通过android手机的usb接口对pc进行网络共享的一种手段，
   一般电脑都会自动识别这种设备，并开启网卡进行共享。

* 1、修改配置文件
  为了进行usb连接，我们需要修改BOOT文件夹下的`config.txt` 与`cmdline.txt`配置文件。
  首先在`config.txt`最末行处换行添加如下代码，打开usb网卡模式：

    ```bash
    sudo vim config.txt
    ```

 * 增加内容如下：

    ```text
    dtoverlay=dwc2
    ```

    ![](http://112firshme11224.test.upcdn.net/images/vim-config.png)


 * `cmdline.txt`文件中找到`rootwait`字段，并在其后面空格添加如下信息，在打开系统时开启usb网卡模式。

    ```bash
    sudo vim cmdline.txt
    ```

 * 增加内容如下：

    ```text
      modules-load=dwc2,g_ether
    ```

  ![](http://112firshme11224.test.upcdn.net/images/vim-cmdline.png)

 * 通过usb连接设备树莓派zero有两个`MicroUSB`口，一个是`电源`插口，职司供电的功能，
   另外一个就是usb接口，它除了供电以外还提供`OTG`的功能，我们也是通过这个接口来连接pc
   ssh连接到树莓派  默认账号为`pi` 密码为`raspberry`

 * 将树莓派的第二个USB口链接，并且链接到MacBook的USB口(Win10都可以)。
 * 等待启动，树莓派zero的灯成黄色急闪状态(在load内核。)
 * 在你的网络内看到以下显示，就可以ssh链接树莓派zero了。

    ![](http://112firshme11224.test.upcdn.net/images/usb-otg.png)

 * 看到如上图：`ssh pi@raspberrypi.local` 输入密码链接即可。

 * 配置网络地址，wifi网卡
    ```bash
    pi@raspberrypi:~ $ sudo nano /etc/wpa_supplicant/wpa_supplicant.conf
    ```
 * 添加以下内容
    ```text
    network={
    ssid="WiFi-B"
    psk="12345678"
    key_mgmt=WPA-PSK
    priority=1
    }
    -----------------
    #ssid:网络的ssid

    #psk:密码

    #priority:连接优先级，数字越大优先级越高（不可以是负数）

    #scan_ssid:连接隐藏WiFi时需要指定该值为1
    ```

 * 保存配置 `control + o 后 Enter`

 * 退出 `control + x `

 * 重启ZeroW `sudo reboot`

 * ping Test

   ![](http://112firshme11224.test.upcdn.net/images/ssh-new-ping.png)


# 2 、配置和下载安装程序

 * 安装软件
   ```bash
    sudo apt-get update
   ```

    ![](http://112firshme11224.test.upcdn.net/images/zero-update.png)

    ```bash
    sudo apt-get install git john
    ```

    ![](http://112firshme11224.test.upcdn.net/images/zero-install-git.png)

 * 克隆代码

   ```bash
     git clone --recursive  http://github.com/mame82/P4wnP1
   ```

    ![](http://112firshme11224.test.upcdn.net/images/git-clone-bad-usb.png)

 * 安装

    ```bash
    cd P4wnP1/
    .／install.sh
    ```

 * 等待 下载大约300m数据，编译需要20分钟。

   ![](http://112firshme11224.test.upcdn.net/images/install-y-success.jpg)

 * 提示

    ```text
       当然最重要的一点是安装后完毕wifi会修改为P4wnP1
       密码是MaMe82-P4wnP1
       ssh连接的地址是172.24.0.1账号和密码还是树莓派的初始密码
       建议修改初始密码！小心被黑(bao)吃(ju)黑(hua)！！！
    ```

 * 将USB设备插到你要渗透的机器上，你会搜索到一个Wifi热点`P4wnP1`,密码`MaMe82-P4wnP1`

 * 链接后进行SSH链接
    `ssh pi@172.24.0.1`即可。注：蓝牙ssh连接的地址是172.26.0.1


 * 配置文件地址如下：
     ```bash
     pi@MAME82-P4WNP1:~ $ nano P4wnP1/setup.cfg
     ```

     ![](http://112firshme11224.test.upcdn.net/images/hack-setting.png)

     * 目前有的功能列表

     ![](http://112firshme11224.test.upcdn.net/images/windows-python-script-hack.png)

     * 测试shell,可以直接获取权限，图片：

     ![](http://112firshme11224.test.upcdn.net/images/success-bad-usb-hack.png)


# 3 成品图：

 ![](http://112firshme11224.test.upcdn.net/images/success-png.jpg)

* 以上操作经过验证可以直接拿去。
* auther `breakEval13`
* https://github.com/breakEval13