---
layout: post
categories: Java Linux Directio
title: Java实现DirectIo文件方式操作文件系统
date: 2019-08-07 09:41:11 +0800
description: Java 采用低消耗写入文件
keywords: directio java
---









## 场景
  * APIServer接口要求较高的并发，而且还要将数据文件存储到本地备份，在低消耗内存CPU情况下提高程序的运行速度以及稳定性。
  * 相关文章 `http://man7.org/linux/man-pages/man2/open.2.html`
  * 实现思路 调用Linux本身的接口 Java采用JNA实现


### 具体方法调用

```java
	private native int open(String pathname, int flags, int mode);
	private native int fcntl(int id, int cmd, int flags);
	private native int write(int fd, Pointer buf, int count);
	private native int posix_memalign(PointerByReference memptr, int alignment, int size);
	private native int close(int fd);
```

### 核心代码

```java
/*
 *    Copyright 2013 Matteo Catena (catena.matteo@gmail.com)
      Release v1 Owner  zhangjianxinnet@gmail.com 
   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
 */
package com.huadi.utils.directio.io;

import java.io.IOException;
import java.io.OutputStream;

import com.sun.jna.Native;
import com.sun.jna.Pointer;
import com.sun.jna.ptr.PointerByReference;

public class JNAOutputStream extends OutputStream {

	//just a couple of them are used, here for eventual future uses
	protected static final int O_RDONLY = 00;
	protected static final int O_WRONLY = 01;
	protected static final int O_RDWR = 02;
	protected static final int O_CREAT = 0100;
	protected static final int O_EXCL = 0200;
	protected static final int O_NOCTTY = 0400;
	protected static final int O_TRUNC = 01000;
	protected static final int O_APPEND = 02000;
	protected static final int O_NONBLOCK = 04000;
	protected static final int O_NDELAY = O_NONBLOCK;
	protected static final int O_SYNC = 010000;
	protected static final int O_ASYNC = 020000;
	protected static final int O_DIRECT = 040000;
	protected static final int O_DIRECTORY = 0200000;
	protected static final int O_NOFOLLOW = 0400000;
	protected static final int O_NOATIME = 01000000;
	protected static final int O_CLOEXEC = 02000000;
	
	//fcntl has also other flags, not presented here
	protected static final int F_SETFL = 4;
	   	
	public static final int BLOCK_SIZE = 512;    
  
	static {
		Native.register("c");
	}    
  
	private native int open(String pathname, int flags, int mode);
	private native int fcntl(int id, int cmd, int flags);
	private native int write(int fd, Pointer buf, int count);
	private native int posix_memalign(PointerByReference memptr, int alignment, int size);
	private native int close(int fd);
  
	private int fd;
	private Pointer bufPnt;
	private byte[] buffer;
	private int pos;
  
	public JNAOutputStream(String pathname, boolean direct) throws IOException {
  			
		if (!direct) {
			fd = open(pathname, O_CREAT|O_TRUNC|O_WRONLY, 00755);
		} else {
//			fd = open(pathname, O_CREAT|O_TRUNC|O_WRONLY|O_DIRECT, 00755);
			fd = open(pathname, O_CREAT|O_APPEND|O_WRONLY|O_DIRECT, 00755);
		}

		PointerByReference pntByRef = new PointerByReference();
		posix_memalign(pntByRef, BLOCK_SIZE, BLOCK_SIZE);
		bufPnt = pntByRef.getValue();    	
		buffer = new byte[BLOCK_SIZE];
	}
			
	@Override
	public void close() throws IOException {
		
		if (pos % BLOCK_SIZE != 0) {
			if (fcntl(fd, F_SETFL, ~O_DIRECT) < 0) {
				throw new IOException("Impossible to modify fd on close()");
			}
			bufPnt.write(0, buffer, 0, pos);
			int rtn = write(fd,bufPnt,pos);
			if (rtn != pos) {
				throw new IOException(rtn+"/"+pos+" bytes written to disk on close()");
			}
		}
		if (close(fd)< 0) {
			throw new IOException("Problems occured while doing close()");
		}
	}	
	
	@Override
	public void write(int b) throws IOException {
		if (pos < BLOCK_SIZE) {
			buffer[pos++] = (byte)b;
		}
		if (pos == BLOCK_SIZE) {
			bufPnt.write(0, buffer, 0, BLOCK_SIZE);
			int rtn = write(fd, bufPnt, BLOCK_SIZE);
			if (rtn != BLOCK_SIZE) {
				throw new IOException(rtn+"/"+BLOCK_SIZE+" bytes written to disk");
			}
			pos = 0;
		}
	}

}

```


### 使用方式 


```java
/*
 *    Copyright 2013 Matteo Catena (catena.matteo@gmail.com)

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
 */
package com.huadi.utils.directio.io;

import java.io.*;
import java.util.Random;


public class Run {

	/**
	 * @ower 张建新
	 * @param args
	 * @throws IOException 
	 */
	private static final int BLOCK_SIZE = 4 * 1024;
	public static void main(String[] args) throws IOException {
		Random rand = new Random(System.nanoTime());
		int runs = Integer.parseInt(args[0]);
		int integers = Integer.parseInt(args[1]);
		long[] sio = new long[2];
		long[] jnaio = new long[2];
		long[] djnaio = new long[2];
		int[] data = new int[integers];
		for (int i = 0; i < integers; i++) {
			data[i] = rand.nextInt();
		}
		for (int i = 0; i < runs; i++) {
			runStd(data, sio);
			runJNA(data, jnaio, false);
			runJNA(data, djnaio, true);
		}
		System.out.println("JNA I/O (ns): "+jnaio[1]/(double)runs+"/"+jnaio[0]/(double)runs);
		System.out.println("O_DIRECT JNA I/O (ns): "+djnaio[1]/(double)runs+"/"+djnaio[0]/(double)runs);
		System.out.println("Standard I/O (ns): "+sio[1]/(double)runs+"/"+sio[0]/(double)runs);
	}

	public static void writeLineToFile2(byte[] data,String filename) throws IOException {
		DataOutputStream out = new DataOutputStream(new JNAOutputStream(filename, true));
		write(out, data);
		out.close();
	}

	
	public static void runStd(int data[], long time[]) throws IOException {
		String filename = "/tmp/"+System.nanoTime()+".bin";
		DataOutputStream out = new DataOutputStream(new BufferedOutputStream(new FileOutputStream(filename), 512));
		time[0] += write(out, data);	
		out.close();
		DataInputStream in = new DataInputStream(new BufferedInputStream(new FileInputStream(filename), 512));
		time[1] += read(in, data);
		in.close();
	}
	
	public static void runJNA(int data[], long time[], boolean o_direct) throws IOException {
		String filename = "/tmp/"+System.nanoTime()+".bin";
		DataOutputStream out = new DataOutputStream(new JNAOutputStream(filename, o_direct));
		time[0] += write(out, data);	
		out.close();
		DataInputStream in = new DataInputStream(new JNAInputStream(filename, o_direct));
		time[1] += read(in, data);
		in.close();
	}	
	
	public static long write(DataOutputStream out, int[] data) throws IOException {
		long t = System.nanoTime();
		for (int i = 0; i < data.length; i++ ) {
			out.writeInt(data[i]);
		}
		out.flush();
		return System.nanoTime() - t;
	}

	public static long write(DataOutputStream out, byte[] data) throws IOException {
		long t = System.nanoTime();
		for (int i = 0; i < data.length; i++ ) {
			out.write(data[i]);
		}
		out.flush();
		return System.nanoTime() - t;
	}


	public static long read(DataInputStream in, int[] data) throws IOException {
		long t = System.nanoTime();
		for (int i = 0; i < data.length; i++) {
			if (data[i] != in.readInt()) {
				throw new IOException("wrong!");
			}
		}
		return System.nanoTime() - t;
	}

}

```

转载请注明出处，本文采用 [CC4.0](http://creativecommons.org/licenses/by-nc-nd/4.0/) 协议授权
