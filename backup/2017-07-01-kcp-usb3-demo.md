---
layout: post
title: K2P 加 USB 2.0 3.0 过程 -- 改造成功，【3.0未实现】[转载]
categories: K2P, 路由器
description: armv7, DDR3,DIY
keywords: armv7, DDR3,DIY
---


发这个贴的原因，只是想跟大家一起探讨一下，也可以为想加USB的朋友们做个参考，同时也欢迎大家来陶侃吐槽。
由于时间的原因，所以，加USB的过程，可能要几天来完成，请大家理解。。。
首先，吐槽一下某东，618凌晨撸的K2P，623晚10点才收到，那个心情捉急啊，某东，我对你太失望了。

特别鸣谢以下的帖子：samson9913
斐讯K2P MT7621AT芯片硬改加USB3.0 2.0---全网首发

http://www.right.com.cn/forum/thread-217002-1-1.html

不好意思，为了不让该帖子沉下去，造福需要改机的机友们，所以加上了回复才能查看的内容。


## 1.改造电源模块，改为固定电压输出，同时已经与：
    把可调电阻换成两个22K（223）贴片电阻串联44K左右，这样输出为5.15伏左右，如果使用47K（473），则电压会比较高5.3~5.5左右。
    更换的原因，是怕用久了可调电阻阻值不知什么时候会发生变化。
   
   ![](http://www.right.com.cn/forum/data/attachment/forum/201706/29/223357fqj9q2tl5hchmhh2.jpg)
   
## 2、电源模块摆放位置：
      刚好在角落两个天线口之间，一点都不松动的哦
      
   ![](http://www.right.com.cn/forum/data/attachment/forum/201706/29/224211fb90mnx9mnnb2zda.jpg)
   
##3、连接USB2.0与3.0的数据线
    （3.0查了很多资料，发现干扰2.4G严重，屏蔽要求很高，手头上没有屏蔽线或屏蔽贴纸，所以放弃3.0的接线，直接2.0的线，改后，两个USB都是2.0来用。）
     同时加了个固定座（方便接线）。
     USB母座的接线图忘记拍了。
   ![](http://www.right.com.cn/forum/data/attachment/forum/201706/29/224746kdjdgkynkg7yefof.jpg)
   
   ![](http://www.right.com.cn/forum/data/attachment/forum/201706/29/224817eqrrgz2q7tf2fk2m.jpg)
   
##4、USB座的固定：
     由于口开好后，USB公座插入拔出时，母座肯定会被带出来的，所以本人想了一下的这种方法，请看图：
     用铜片与USB的贴片焊死，就不会动了
   ![](http://www.right.com.cn/forum/data/attachment/forum/201706/29/225323z90920f4vffr0fb4.jpg)
   
   ![](http://www.right.com.cn/forum/data/attachment/forum/201706/29/225323jc4kb9lt7alc4zlb.jpg)
   
   ![](http://www.right.com.cn/forum/data/attachment/forum/201706/29/225520gzuns9hhhuhocheh.jpg)
   
## 5. 以上独立电源、USB接线完毕。加电接移动硬盘能正常开机，至于U盘与移动硬盘的测试，又需要下个回合了。
   
       下一回合：刷固件测试USB。
       敬请期待。。。
      
## 6、刷固件：
     【2017-06-26】斐讯K2P官方固件定制版，WEB直刷+不死BOOT，加adb、S-S R等【V1.21】
       http://www.right.com.cn/forum/thread-217181-1-1.html
       Web直刷，上面帖子有过程。
   ![](http://www.right.com.cn/forum/data/attachment/forum/201706/30/215353xbob7gz4pj3bpkjd.jpg)
   
   ![](http://www.right.com.cn/forum/data/attachment/forum/201706/30/215538to2umb71ovruemek.jpg)
   
   ![](http://www.right.com.cn/forum/data/attachment/forum/201706/30/215837moczren7er7lmpn7.jpg)
## 7.sda1 与 usbsda1 是32G的U盘，usbsdb1 是IDE的移动硬盘。
    该固件由于还有共享的问题没有解决，所以，本人直接测了一下读取的速度。
    U盘的如下：
   ![](http://www.right.com.cn/forum/data/attachment/forum/201706/30/220148i64ea6weahua8a84.jpg)
   
   
>至此，K2P 加装 USB 完美结束（可惜由于USB3.0的干扰问，没有继续）。