---
layout: post
title: Crate DB重启数据恢复Shards较慢解决方案[原创]
categories: java
description: crate.io
keywords: crate.io,databases
---


Crate DB重启数据恢复Shards较慢解决方案“采用ES restFull API”


date：2017-10-31
author：zhangjianxin

# Crate DB重启数据恢复Shards较慢解决方案

[*] 成功图，


   ![](http://112firshme11224.test.upcdn.net/images/crate-success.png)


[*] 解决方案


# 解决Crate 重启恢复Shards速度慢的问题
  * 修改crate配置文件
   ```bash
   es.api.enabled: true
   ```

  * 开始处理
  ```bash
  curl -XPUT http://127.0.0.1:9200/_cluster/settings -d'
  {
  	"transient" : {
  		"cluster.routing.allocation.enable" : "none"
  	}
  }'
  ```

# 感谢`国庆`大佬
# 感谢`禹神`大佬

* 以上操作经过验证可以直接拿去。
* auther `breakEval13`
* https://github.com/breakEval13