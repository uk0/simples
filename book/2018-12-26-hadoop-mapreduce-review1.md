---
layout: post
categories: mapreduce
title: MapReduce 基于文件的 join数据清洗
date: 2018-12-26 17:51:32 +0800
description: map reduce join for hdfs
keywords: mapreduce hadoop hdfs
---


这只是一次突然需求导致的问题，解决的比较唐突，如果哪里有不对的地方欢迎指正，不希望误导他人。


## MapReduce 实现基于文件的Join

* 临时整理勿喷.

```java

package com.huaching.task;

import org.apache.hadoop.conf.Configuration;
import org.apache.hadoop.fs.FileSystem;
import org.apache.hadoop.fs.Path;
import org.apache.hadoop.io.Text;
import org.apache.hadoop.mapred.JobConf;
import org.apache.hadoop.mapreduce.Job;
import org.apache.hadoop.mapreduce.Reducer;
import org.apache.hadoop.mapreduce.Mapper;
import org.apache.hadoop.mapreduce.lib.input.MultipleInputs;
import org.apache.hadoop.mapreduce.lib.input.TextInputFormat;
import org.apache.hadoop.mapreduce.lib.output.FileOutputFormat;
import org.apache.hadoop.mapreduce.lib.output.LazyOutputFormat;
import org.apache.hadoop.mapreduce.lib.output.MultipleOutputs;
import org.apache.hadoop.mapreduce.lib.output.TextOutputFormat;

import java.io.IOException;
import java.net.URI;

/**
 * @author firshme
 * @version 2018-12-19.
 * @url github.com/uk0
 * @project Task
 * @since JDK1.8.
 *
 * # 多个cvs文件 以key 相同进行分组  并且合并关键字
 */
public class JoinReducesSimple {

    public static class Mapper1 extends Mapper<Object, Text, Text, Text> {
        int counter = 0;

        // key index
        @Override
        public void map(Object key, Text value, Context context)
                throws IOException, InterruptedException {
            String record = value.toString();
            String[] parts = record.split(",");
            StringBuilder lineData = new StringBuilder();
            for (String sk : parts) {
                lineData.append(sk).append(",");
            }
            counter++;
            if (counter % 10000 == 0) {
                System.out.println("Mapper1---->" + lineData);
            }
            context.write(new Text(parts[0]), new Text("pay\t" + lineData.substring(0, lineData.length() - 1)));
        }
    }



    public static class Mapper2 extends Mapper<Object, Text, Text, Text> {
        int counter = 0;

        @Override
        public void map(Object key, Text value, Context context)
                throws IOException, InterruptedException {
            String record = value.toString();
            String[] parts = record.split(",");
            StringBuilder lineData = new StringBuilder();
            for (String sk : parts) {
                lineData.append(sk).append(",");
            }
            counter++;
            if (counter % 10000 == 0) {
                System.out.println("Mapper2---->" + lineData);
            }
            context.write(new Text(parts[4]), new Text("deposit\t" + lineData.substring(0, lineData.length() - 1)));
        }
    }

    public static class Mapper3 extends Mapper<Object, Text, Text, Text> {
        int counter = 0;

        @Override
        public void map(Object key, Text value, Context context)
                throws IOException, InterruptedException {
            String record = value.toString();
            String[] parts = record.split(",");
            StringBuilder lineData = new StringBuilder();
            for (String sk : parts) {
                lineData.append(sk).append(",");
            }
            counter++;
            if (counter % 10000 == 0) {
                System.out.println("Mapper3---->" + lineData);
            }
            context.write(new Text(parts[0]), new Text("car_user\t" + lineData.substring(0, lineData.length() - 1)));
        }
    }

    public static class Reducer1 extends Reducer<Text, Text, Text, Text> {

        int counter = 0;

        public void reduce(Text key, Iterable<Text> values, Context context)
                throws IOException, InterruptedException {
            if (counter == 0) {
                counter++;
                String title = "用户ID,手机号(如没有手机则空),注册时间,缴纳押金时间,第一次停车时间,第一次停车地点,累计消费金额(不含押金),使用过停车场数量,停车订单累计次数,最后一次停车时,最后一次支付时间,最后一次退押金时间(如现在未退押金显示押金金额)";
                context.write(new Text("Number,"), new Text(title));
            }
            counter++;
            String car_user = "";
            String deposit = "";
            String pay = "";

            // out
            String userID = "";
            String mobile = "";
            String firstParking = "";
            String firstPark = "";
            String sumPrice = "";
            String parks = "";
            String count = "";
            String lastParking = "";
            String lastPay = "";
            String depositTime = "";
            String register_time = "";
            String lastRefundDeposit = "";

            String park = "";
            for (Text t : values) {
                String[] parts = t.toString().split("\t");
                if (parts[0].trim().equals("car_user")) {
                    car_user = parts[1];
                    String [] car_userArrays = car_user.split(",");
                    if (car_userArrays.length>=24) {
                        userID = car_userArrays[0];
                        mobile = car_userArrays[2];
                        register_time = car_userArrays[19];
                    }
                } else if (parts[0].trim().equals("deposit")) {
                    deposit = parts[1];
                    String[] depositArray = deposit.split(",");
                    depositTime = depositArray[2];  //押金创建时间
                    if ("0.00".equals(depositArray[3])) {
                        lastRefundDeposit = depositArray[1]; // 已经退了押金查看退的时间
                    } else {
                        lastRefundDeposit = depositArray[3];  //显示金额
                    }

                } else if (parts[0].trim().equals("pay")) {
                    pay = parts[1];
                    String[] payArrays = pay.split(",");
                    firstParking = payArrays[3];
                    firstPark = payArrays[4];
                    sumPrice = payArrays[5];
                    parks = payArrays[1];
                    count = payArrays[2];
                    lastParking = payArrays[6];
                    lastPay = payArrays[7];
                }
                else if (parts[0].trim().equals("park")) {
                    park = parts[1];
                }
            }
            context.write(new Text(String.valueOf(counter) + ","), new Text(userID + "," + mobile + "," + register_time + "," + depositTime + "," + firstParking + "," + firstPark + "," + sumPrice + "," + parks + "," + count + "," + lastParking + "," + lastPay + "," + lastRefundDeposit+","+park));
        }
    }

    public static void main(String[] args) throws Exception {
        Configuration conf = new Configuration();
        JobConf jobConf = new JobConf(conf, JoinReducesSimple.class);
        Job job = Job.getInstance(jobConf);
        job.setJarByClass(JoinReducesSimple.class);

        job.setReducerClass(Reducer1.class);

        job.setOutputKeyClass(Text.class);
        job.setOutputValueClass(Text.class);

        MultipleInputs.addInputPath(job, new Path("/public/csv/parkingPayment/parkingPayment-r-00000"), TextInputFormat.class, Mapper1.class);
        MultipleInputs.addInputPath(job, new Path("/private/parking_ths_deposit_user.csv"), TextInputFormat.class, Mapper2.class);
        MultipleInputs.addInputPath(job, new Path("/private/parking_car_user.csv"), TextInputFormat.class, Mapper3.class);

        String outFile = "/public/csv/job1/";
        Path path = new Path(outFile);
        FileSystem fs = FileSystem.get(URI.create(outFile), conf);
        if (fs.exists(path)) {
            fs.delete(path, true);
        }
        MultipleOutputs.addNamedOutput(job, "job", TextOutputFormat.class, Text.class, Text.class);

        LazyOutputFormat.setOutputFormatClass(job, TextOutputFormat.class);
        FileOutputFormat.setOutputPath(job, path);
        System.exit(job.waitForCompletion(true) ? 0 : 1);
    }
}

```



## 大致实现方式

```java
/***
    MapReduce 分为 map 阶段 reduce 阶段
    Map 阶段可以理解为将数据根据指定的key 进行整理
    reduce 阶段接到key 以及 value
    多个map 同时进行读取数据 并且以相应的key 进行分组
    reduce 阶段读取到 key value  value 里面取巧的加了一下 相应Map数据的标记，为了获取相应的数据。

    map1data = [id,y,x,z]
    map2data = [y,id,z]
    map3data = [f,y,x,z,id]


     Map1  Map2  Map3
      |     |      |
     key[0] key[1] key[4]
     \      |      /
      \     |     /
          reduce
            |
        key[0],value[map1,map2,map3 value传递的数据]

        reduce的时候多个Map的Key已经被分组，
            这个时候你读取value 可能是他们其中某一个map 的value
            进行相应的操作即可

***/
```



转载请注明出处，本文采用 [CC4.0](http://creativecommons.org/licenses/by-nc-nd/4.0/) 协议授权
