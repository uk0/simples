---
layout: post
categories: dejavu python
title: 基于Python DeJavu 进音乐指纹识别，识别歌曲
date: 2018-07-01 01:24:05 +0800
description: 识别音乐
keywords: 
---

这只是一个Demo

## 开始

* 将音频数据分析的结果存储到数据库

```python
#!/usr/bin/python
# -*- coding: utf-8 -*-
from dejavu import Dejavu

config = {
    "database": {
        "host": "127.0.0.1",
        "user": "root",
        "passwd": "passw0rd",
        "db": "dejavu",
    },
    "database_type": "mysql",
}

# 将数据提取到Mysql （音乐的指纹）

djv = Dejavu(config)

djv.fingerprint_directory("/Users/zhangjianxin/home/mp3/", [".mp3"])
print(djv.db.get_num_fingerprints())

```

* 监听麦克风进行识别歌曲

```python
#!/usr/bin/python
# -*- coding: utf-8 -*-


from dejavu import Dejavu
from dejavu.recognize import FileRecognizer, MicrophoneRecognizer

config = {
    "database": {
        "host": "127.0.0.1",
        "user": "root",
        "passwd": "passw0rd",
        "db": "dejavu",
    },
    "database_type": "mysql"
}
djv = Dejavu(config)

# 识别从mic输入的音频 只获得5秒的音频
secs = 5
song = djv.recognize(MicrophoneRecognizer, seconds=secs)
if song is None:
    print "No Match"
else:
    print(song)
```


## 需要用到的工具

```bash
python 2.7

brew install portaudio
brew install ffmpeg
brew install mysql-connector-c # my_config.h
brew install mysql
#!!!  mysql-connector-c (my_config.h) ->  mysql
ln -s  ../../../../mysql-connector-c/6.1.11/include/my_config.h ./


brew link mysql # create links

pip install mysql-python
pip install numpy==1.10.2 # 版本兼容性问题
pip install PyDejavu

```

![](http://112firshme11224.test.upcdn.net/demos/788db197-f796-42d2-8ea5-ef655160e7c5.png)

转载请注明出处，本文采用 [CC4.0](http://creativecommons.org/licenses/by-nc-nd/4.0/) 协议授权

