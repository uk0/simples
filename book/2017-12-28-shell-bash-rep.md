---
layout: post
categories: shell bash 
title: shell bash 模版填充ENV
date: 2017-12-28 21:46:47 +0800
description: 通过环境变量来初始化模版数据
keywords: storm bash shell 替换  模版 
---





# quick  start


```bash
#!/bin/bash
ZK_IP="192.168.1.1"
cat ${template_path} |
awk '$0 !~ /^\s*#.*$/' |
sed 's/[ "]/\\&/g' |
while read -r line;do
    eval echo ${line}
done >${dst_path}
```

转载请注明出处，本文采用 [CC4.0](http://creativecommons.org/licenses/by-nc-nd/4.0/) 协议授权
