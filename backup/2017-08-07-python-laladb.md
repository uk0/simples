---
layout: post
title:  Python 写一个简易的数据库.
categories: Python,lambda
description: 笔记
keywords: Python,lambda
---


没事看看文章而已～

# 教你用100多行写一个数据库（附源码）

本文介绍的是以为中国的IT资深人士写的一个简单的数据库，没有我们使用的数据库那么强大，但是值得大家借鉴。可以用在特定环境中，更加灵活方便。

数据库的名字叫WawaDB，是用python实现的。由此可见python是灰常强大啊！

## 简介

记录日志的需求一般是这样的：

只追加，不修改，写入按时间顺序写入；

大量写，少量读，查询一般查询一个时间段的数据；

MongoDB的固定集合很好的满足了这个需求，但是MongoDB占内存比较大，有点儿火穿蚊子，小题大做的感觉。

WawaDB的思路是每写入1000条日志，在一个索引文件里记录下当前的时间和日志文件的偏移量。

然后按时间询日志时，先把索引加载到内存中，用二分法查出时间点的偏移量，再打开日志文件seek到指定位置，这样就能很快定位用户需要的数据并读取，而不需要遍历整
个日志文件。

## 性能

Core 2 P8400,2.26GHZ,2G内存，32 bit win7

写入测试：

模拟1分钟写入10000条数据，共写入5个小时的数据， 插入300万条数据，每条数据54个字符，用时2分51秒



读取测试:读取指定时间段内包含某个子串的日志

数据范围 遍历数据量 结果数 用时（秒）

5小时 300万 604 6.6

2小时 120万 225 2.7

1小时 60万 96 1.3

30分钟 30万 44 0.6

## 索引

只对日志记录的时间做索引， 简介里大概说了下索引的实现，二分查找肯定没B Tree效率高，但一般情况下也差不了一个数量级，而且实现特别简单。

因为是稀疏索引，并不是每条日志都有索引记录它的偏移量，所以读取数据时要往前多读一些数据，防止漏读，等读到真正所需的数据时再真正给用户返回数据。

如下图，比如用户要读取25到43的日志，用二分法找25，找到的是30所在的点，

索引：0         10        20        30        40        50
日志：|.........|.........|.........|.........|.........|>>>a = [0, 10, 20, 30,
40, 50]>>>bisect.bisect_left(a, 35)>>>3>>>a[3]>>>30>>>bisect.bisect_left(a,
43)>>>5>>>a[5]>>>50

所以我们要往前倒一些，从20（30的前一个刻度）开始读取日志，21，22，23，24读取后因为比25小，所以扔掉, 读到25,26,27,...后返回给用户

读取到40（50的前一个刻度）后就要判断当前数据是否大于43了，如果大于43（返回全开区间的数据），就要停止读了。

整体下来我们只操作了大文件的很少一部分就得到了用户想要的数据。

## 缓冲区

为了减少写入日志时大量的磁盘写，索引在append日志时，把buffer设置成了10k，系统默认应该是4k。

同理，为了提高读取日志的效率，读取的buffer也设置了10k，也需要根据你日志的大小做适当调整。

索引的读写设置成了行buffer，每满一行都要flush到磁盘上，防止读到不完整的索引行（其实实践证明，设置了行buffer，还是能读到半拉的行）。

## 查询

啥？要支持SQL，别闹了，100行代码怎么支持SQL呀。

现在查询是直接传入一个lambada表达式，系统遍历指定时间范围内的数据行时，满足用户的lambada条件才会返回给用户。

当然这样会多读取很多用户不需要的数据，而且每行都要进行lambda表达式的运算，不过没办法，简单就是美呀。

以前我是把一个需要查询的条件和日志时间，日志文件偏移量都记录在索引里，这样从索引里查找出符合条件的偏移量，然后每条数据都如日志文件里seek一次，read一
次。这样好处只有一个，就是读取的数据量少了，但缺点有两个：

索引文件特别大，不方便加载到内存中

每次读取都要先seek，貌似缓冲区用不上，特别慢，比连续读一个段的数据，并用lambda过滤慢四五倍

## 写入

前面说过了，只append，不修改数据，而且每行日志最前面是时间戳。

多线程



查询数据，可以多线程同时查询，每次查询都会打开一个新的日志文件的描述符，所以并行的多个读取不会打架。

写入的话，虽然只是append操作，但不确认多线程对文件进行append操作是否安全，所以建议用一个队列，一个专用线程进行写入。

## 锁

没有任何锁。

## 排序

默认查询出来的数据是按时间正序排列，如需其它排序，可取到内存后用python的sorted函数排序，想怎么排就怎么排。



## 100多行的数据库代码


```python

    # -*- coding:utf-8 -*-
    import os
    import time
    import bisect
    import itertools
    from datetime import datetime
    import logging

    default_data_dir = './data/'
    default_write_buffer_size = 1024*10
    default_read_buffer_size = 1024*10
    default_index_interval = 1000

    def ensure_data_dir():
        if not os.path.exists(default_data_dir):
            os.makedirs(default_data_dir)

    def init():
        ensure_data_dir()

    class WawaIndex:
        def __init__(self, index_name):
            self.fp_index = open(os.path.join(default_data_dir, index_name + '.index'), 'a+', 1)
            self.indexes, self.offsets, self.index_count = [], [], 0
            self.__load_index()

        def __update_index(self, key, offset):
            self.indexes.append(key)
            self.offsets.append(offset)

        def __load_index(self):
            self.fp_index.seek(0)
            for line in self.fp_index:
                try:
                    key, offset  = line.split()
                    self.__update_index(key, offset)
                except ValueError: # 索引如果没有flush的话，可能读到有半行的数据
                    pass

        def append_index(self, key, offset):
            self.index_count += 1
            if self.index_count % default_index_interval == 0:
                self.__update_index(key, offset)
                self.fp_index.write('%s %s %s' % (key, offset, os.linesep))

        def get_offsets(self, begin_key, end_key):
            left = bisect.bisect_left(self.indexes, str(begin_key))
            right = bisect.bisect_left(self.indexes, str(end_key))
            left, right = left - 1, right - 1
            if left < 0: left = 0
            if right < 0: right = 0
            if right > len(self.indexes) - 1: right = len(self.indexes) - 1
            logging.debug('get_index_range:%s %s %s %s %s %s', self.indexes[0], self.indexes[-1], begin_key, end_key, left, right)
            return self.offsets[left], self.offsets[right]


    class WawaDB:
        def __init__(self, db_name):
            self.db_name = db_name
            self.fp_data_for_append = open(os.path.join(default_data_dir, db_name + '.db'), 'a', default_write_buffer_size)
            self.index = WawaIndex(db_name)

        def __get_data_by_offsets(self, begin_key, end_key, begin_offset, end_offset):
            fp_data = open(os.path.join(default_data_dir, self.db_name + '.db'), 'r', default_read_buffer_size)
            fp_data.seek(int(begin_offset))

            line = fp_data.readline()
            find_real_begin_offset = False
            will_read_len, read_len = int(end_offset) - int(begin_offset), 0
            while line:
                read_len += len(line)
                if (not find_real_begin_offset) and  (line < str(begin_key)):
                    line = fp_data.readline()
                    continue
                find_real_begin_offset = True
                if (read_len >= will_read_len) and (line > str(end_key)): break
                yield line.rstrip('\r\n')
                line = fp_data.readline()

        def append_data(self, data, record_time=datetime.now()):
            def check_args():
                if not data:
                    raise ValueError('data is null')
                if not isinstance(data, basestring):
                    raise ValueError('data is not string')
                if data.find('\r') != -1 or data.find('\n') != -1:
                    raise ValueError('data contains linesep')

            check_args()

            record_time = time.mktime(record_time.timetuple())
            data = '%s %s %s' % (record_time, data, os.linesep)
            offset = self.fp_data_for_append.tell()
            self.fp_data_for_append.write(data)
            self.index.append_index(record_time, offset)

        def get_data(self, begin_time, end_time, data_filter=None):
            def check_args():
                if not (isinstance(begin_time, datetime) and isinstance(end_time, datetime)):
                    raise ValueError('begin_time or end_time is not datetime')

            check_args()

            begin_time, end_time = time.mktime(begin_time.timetuple()), time.mktime(end_time.timetuple())
            begin_offset, end_offset = self.index.get_offsets(begin_time, end_time)

            for data in self.__get_data_by_offsets(begin_time, end_time, begin_offset, end_offset):
                if data_filter:
                    if data_filter(data):
                        yield data
                else:
                    yield data

    def test():
        from datetime import datetime, timedelta
        import uuid, random
        logging.getLogger().setLevel(logging.NOTSET)

        def time_test(test_name):
            def inner(f):
                def inner2(*args, **kargs):
                    start_time = datetime.now()
                    result = f(*args, **kargs)
                    print '%s take time:%s' % (test_name, (datetime.now() - start_time))
                    return result
                return inner2
            return inner

        @time_test('gen_test_data')
        def gen_test_data(db):
            now = datetime.now()
            begin_time = now - timedelta(hours=5)
            while begin_time < now:
                print begin_time
                for i in range(10000):
                    db.append_data(str(random.randint(1,10000))+ ' ' +str(uuid.uuid1()), begin_time)
                begin_time += timedelta(minutes=1)

        @time_test('test_get_data')
        def test_get_data(db):
            begin_time = datetime.now() - timedelta(hours=3)
            end_time = begin_time + timedelta(minutes=120)
            results = list(db.get_data(begin_time, end_time, lambda x: x.find('1024') != -1))
            print 'test_get_data get %s results' % len(results)

        @time_test('get_db')
        def get_db():
            return WawaDB('test')

        if not os.path.exists('./data/test.db'):
            db = get_db()
            gen_test_data(db)
            #db.index.fp_index.flush()

        db = get_db()
        test_get_data(db)

    init()

    if __name__ == '__main__':
        test()
```




* 以上代码经过验证可以直接拿去修改调用
* Owner `breakEval13`
* https://github.com/breakEval13