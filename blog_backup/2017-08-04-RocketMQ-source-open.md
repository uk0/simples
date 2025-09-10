---
layout: post
title:  Flink的自定义Source 如何给重写的open方法传递参数.
categories: Flink,rocektmq
description: 笔记
keywords: flink,rocektmq
---


因工作需求所整合Flink + RocketMQ,中间遇到许多问题本文章将问题解决。

# 部分代码：
   * Flink程序入口
   * 继承RichSinkFunction，重写`open` 和 `invoke`方法，还有`close`。


```java
 final ParameterTool parameterTool = ParameterTool.fromArgs(args);
        if (parameterTool.getNumberOfParameters() < 3) {
                System.out.println("Missing parameters!");
                System.out.println("\nUsage: RocketMQ \n" +
                        "--topic <topic> \n" +
                        "--cgroup <consumerGroupName> \n" +
                        "--conf <consumerGroupName> \n" +
                        "--hdfs <hdfs path>\n"
                );
                return;
        }
            env.getConfig().enableSysoutLogging();
            env.getConfig().setRestartStrategy(RestartStrategies
                    .fixedDelayRestart(4, 100));
            env.enableCheckpointing(5000);
            env.setMaxParallelism(72);
               /**创建 Configuration 对象**/
            Configuration conf = new Configuration();
            conf.setString("topic",parameterTool.get("topic"));
            conf.setString("cgroup",parameterTool.get("cgroup"));
            conf.setString("conf",parameterTool.get("conf"));
            conf.setString("hdfs",parameterTool.get("hdfs"));
            /**将Configuration存入数据，放入全局执行环境内*/
            env.getConfig().setGlobalJobParameters(conf);
            RocketMQSource rocketMQSource = new RocketMQSource();
            DataStream<String> sourceStream = env.addSource(rocketMQSource);
```
# Source [重点]
   * 自定义Source的一部分代码

```java
 @Override
    public void open(Configuration parameters) throws Exception {
        super.open(parameters);
        /**获取运行环境内的参数**/
        ExecutionConfig.GlobalJobParameters globalParams = getRuntimeContext().getExecutionConfig().getGlobalJobParameters();
        /**执行参数提取**/
        Configuration globConf = (Configuration) globalParams;
        queue = new LinkedBlockingQueue<>(1000);
        MQConfig.getInstance(globConf.getString("conf",null));
        /**拿到存入参数，并且设置默认值**/
        consumer = new DefaultMQPushConsumer( globConf.getString("cgroup", null));
        consumer.setNamesrvAddr( MQConfig.getMQ());
        consumer.setInstanceName(UUID.randomUUID().toString());
        consumer.setMessageModel(MessageModel.CLUSTERING);// 消费模式
        consumer.setConsumeFromWhere(ConsumeFromWhere.CONSUME_FROM_FIRST_OFFSET);
        consumer.registerMessageListener(this);
        consumer.subscribe(globConf.getString("topic", null), "*");
        consumer.setConsumeThreadMax(64);
        consumer.setConsumeThreadMin(8);
        consumer.setConsumeMessageBatchMaxSize(8);//消息数量每次读取的消息数量
        System.out.println("启动线程");
        consumer.start();
        System.out.println("RocketMQ Started.");
        LOG.info("consumeBatchSize:{} pullBatchSize:{} consumeThread:{}",consumer.getConsumeMessageBatchMaxSize(),consumer.getPullBatchSize(),consumer.getConsumeThreadMax());
    }
```

* 以上代码经过验证可以直接拿去修改调用
* Owner `breakEval13`
* https://github.com/breakEval13