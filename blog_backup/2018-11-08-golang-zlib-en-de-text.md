---
layout: post
categories: zlib golang php
title: Golang PHP zlib (php)加密<=>解密(go).
date: 2018-11-08 20:14:05 +0800
description: text zlib and deflate. php 加密 Golang 解密
keywords: php golang zlib加密解密
---

## 这个问题是来自一个群友的，一开始马虎大意 后来解决问题。

s

### Golang 代码

```golang
package main

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"fmt"
	"io"
)
//进行zlib压缩
func DoZlibCompress(src []byte) []byte {
	var in bytes.Buffer
	w, _ := zlib.NewWriterLevelDict(&in,-1,src)
	w.Write(src)
	w.Close()
	return in.Bytes()
}

func DoZlibCompress2(src []byte) []byte {
	var bufs bytes.Buffer
	w, _ := zlib.NewWriterLevel(&bufs, -1)
	w.Write(src)
	defer w.Flush()
		w.Close()
	return bufs.Bytes()
}

//进行zlib解压缩
func DoZlibUnCompress(compressSrc []byte) []byte {
	b := bytes.NewReader(compressSrc)
	var out bytes.Buffer
	r, _ := zlib.NewReader(b)
	io.Copy(&out, r)
	return out.Bytes()
}

func main() {
	//buff, _ := base64.StdEncoding.DecodeString("eJxLTEoGAAJNASc=")
	buff, _ := base64.StdEncoding.DecodeString("eJxKTEoGBAAA//8CTQEn")
	fmt.Println(string(DoZlibUnCompress(buff)))



	//zip := DoZlibCompress([]byte("abc"))
	zip := DoZlibCompress2([]byte("abc"))
	fmt.Println(zip)

	fmt.Println(base64.StdEncoding.EncodeToString(zip))

	//
	//fmt.Println(string(DoZlibUnCompress(zip)))
}

...
abc
[120 156 74 76 74 6 4 0 0 255 255 2 77 1 39]
eJxKTEoGBAAA//8CTQEn
```



### PHP 代码


```php

<?php

$code =zlib_encode("abc",ZLIB_ENCODING_DEFLATE);
$code = base64_encode($code);
echo  "$code";
// eJxLTEoGAAJNASc=

$code = bin2hex($code);
echo  "$code";
// 654a784c54456f4741414a4e4153633d

?>
```


### 疑问

*  为什么种方式加密的hex 不一致解析的结果却是一致的。
*  Golang 是对Zlib 进行了哪些优化？


测试地址。> http://www.txtwizard.net/compression

https://stackoverflow.com/questions/40893411/golang-python-zlib-difference


https://stackoverflow.com/questions/52767214/compressed-output-differs-from-go-to-ruby-implementation


转载请注明出处，本文采用 [CC4.0](http://creativecommons.org/licenses/by-nc-nd/4.0/) 协议授权
