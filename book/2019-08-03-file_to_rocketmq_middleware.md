---
layout: post
categories: rocketmq golang
title: 读取数据发送到MQ
date: 2019-08-03 14:01:12 +0800
description: 监控文件夹变动产生新文件读取该文件内容每一行都发送到RocketMQ
keywords: RocketMQ Golang
---


数据通过FTP上传到服务器要求将数据发送到RocketMQ原版本用Java写的（不是本人）TPS上不去资源消耗大等问题。



## 缘由

 * 数据通过FTP上传到服务器要求将数据发送到RocketMQ原版本用Java写的（不是本人）TPS上不去资源消耗大等问题。
 * 解决方案：采用Golang 重新实现，用到了Golang的 Go关键字 WaitGroup MMAP等。
 * 目前TPS 在虚拟机上跑 TPS稳定 5600 - 4700 
 * 链接如下：https://github.com/uk0/file_to_rocketmq_middleware


## 核心代码介绍

```go
package main

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/zjykzk/rocketmq-client-go/log"
	"github.com/zjykzk/rocketmq-client-go/message"
	"github.com/zjykzk/rocketmq-client-go/producer"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"sync/atomic"
	"time"
	"unsafe"
)

type Client struct {
	P            *producer.Producer
	Group        string
	NamesrvAddrs string
	DataChan     chan Message
}

func newLogger(filename string) (log.Logger, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		println("create file error", err.Error())
		return nil, err
	}

	logger := log.New(file, "", log.Ldefault)
	logger.Level = log.Ldebug

	return logger, nil
}

func (cli *Client) SendMsg(stati *statiBenchmarkProducerSnapshot) {
	now := time.Now()
	var count = 1;
	for messageBody := range cli.DataChan {
		m := &message.Message{Topic: messageBody.topic, Body: []byte(messageBody.line)}

		r, err := cli.P.SendSync(m)
		if err != nil {
			logs.Debug("Send with sync error:%s\n", err)
			continue
		}

		// 取模
		//var BatchLines = []message.Data{};
		//BatchLines = append(BatchLines, message.Data{Body: []byte(messageBody.line)})
		//if count%defaultSendLimitbatch == 0 {
		//
		//	r, err3 := client.SendBatchSync(&message.Batch{Topic: messageBody.topic, Datas: BatchLines})
		//	if err3 != nil {
		//		logs.Debug("SendWithBatchSync error:%s\n", err3)
		//		continue
		//	}
		//	logs.Debug(fmt.Sprintf("BatchSendMessageResult: %s FileName = %s  Process = %f%%  ProcessCount = %d , TotalCount = %d", r.OffsetID, messageBody.FileName, float32(messageBody.ProcessCount)/float32(messageBody.TotalCount)*100, messageBody.ProcessCount, messageBody.TotalCount))
		//}

		//selector
		//r, err2 := client.SendSyncWithSelector(m, messageQueueSelector{}, count)
		//if err2 != nil {
		//	logs.Debug("SendWithSelector error:%s\n", err2)
		//	continue
		//}

		if r.Status == producer.OK {
			atomic.AddInt64(&stati.receiveResponseSuccessCount, 1)
			atomic.AddInt64(&stati.sendRequestSuccessCount, 1)
			currentRT := int64(time.Since(now) / time.Millisecond)
			atomic.AddInt64(&stati.sendMessageSuccessTimeTotal, currentRT)
			prevRT := atomic.LoadInt64(&stati.sendMessageMaxRT)
			for currentRT > prevRT {
				if atomic.CompareAndSwapInt64(&stati.sendMessageMaxRT, prevRT, currentRT) {
					break
				}
				prevRT = atomic.LoadInt64(&stati.sendMessageMaxRT)
			}
			if count%1000 == 0 {
				logs.Debug(fmt.Sprintf("SendMessageResultQueueOffset : %d FileName = %s  Process = %f%%  ProcessCount = %d , TotalCount = %d", r.QueueOffset, messageBody.FileName, float32(messageBody.ProcessCount)/float32(messageBody.TotalCount)*100, messageBody.ProcessCount, messageBody.TotalCount))
			}
		}
		count++;
	}
	logs.Debug("DONE")
}

func NewMQSender(groupId string,nameSvr string) (cli *Client, err error) {
	cli = &Client{
		Group:groupId,
		NamesrvAddrs:nameSvr,
		DataChan:data,
	}
	stati := statiBenchmarkProducerSnapshot{}
	snapshots := produceSnapshots{cur: &stati}


	logger, err := newLogger("producer.log")
	if err != nil {
		logs.Debug("new logger of producer.loge error:%s\n", err)
		return
	}

	client := producer.New(groupId, strings.Split(nameSvr, ";"), logger)

	err = client.Start()
	if err != nil {
		logs.Debug("start producer error: ", err)
		return
	}
	//defer client.Shutdown()

	cli.P = client
	// 根据CPU来协商启动协程数
	for i := 0; i < runtime.GOMAXPROCS(runtime.NumCPU()); i++ {
		go func() {
			waitgroup.Add(1)
			cli.SendMsg(&stati)
			waitgroup.Done()
		}()
	}
	// snapshot
	go func() {
		waitgroup.Add(1)
		defer waitgroup.Done()
		ticker := time.NewTicker(time.Second)
		for {
			select {
			case <-ticker.C:
				snapshots.TakeSnapshot()
			}
		}
	}()
	// print statistic
	go func() {
		waitgroup.Add(1)
		defer waitgroup.Done()
		ticker := time.NewTicker(time.Second * 10)
		for {
			select {
			case <-ticker.C:
				snapshots.TrintStati()
			}
		}
	}()
	return
}

type messageQueueSelector struct{}

func (s messageQueueSelector) Select(mqs []*message.Queue, m *message.Message, arg interface{}) *message.Queue {
	orderID := arg.(int)
	return mqs[orderID%len(mqs)]
}

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxMask = 1<<6 - 1 // All 1-bits, as many as 6
)

// 随机
var src = rand.NewSource(time.Now().UnixNano())

func RandStringBytesMaskImprSrc(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for 10 characters!
	for i, cache, remain := n-1, src.Int63(), 10; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), 10
		}
		b[i] = letterBytes[int(cache&letterIdxMask)%len(letterBytes)]
		i--
		cache >>= 6
		remain--
	}
	return *(*string)(unsafe.Pointer(&b))
}


```
* 调用方式 `	_, _ = NewMQSender(TConfig.SendGroup, TConfig.RocketMQNameserver)
`

转载请注明出处，本文采用 [CC4.0](http://creativecommons.org/licenses/by-nc-nd/4.0/) 协议授权
