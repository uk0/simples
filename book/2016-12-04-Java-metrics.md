---
layout: post
title: Java metrics
categories: GitHub
description: 关于Java web 监控系统。
keywords: Java, GitHub Pages
---

项目中有对程序的一个健康检查，以及TPS,QPS等。

## 文档

*Maven

```code
 <!--  metrics * -->
    <dependency>
        <groupId>io.dropwizard.metrics</groupId>
        <artifactId>metrics-core</artifactId>
        <version>3.1.0</version>
    </dependency>
    <dependency>
        <groupId>io.dropwizard.metrics</groupId>
        <artifactId>metrics-servlets</artifactId>
        <version>3.1.0</version>
    </dependency>
    <dependency>
        <groupId>com.codahale.metrics</groupId>
        <artifactId>metrics-healthchecks</artifactId>
        <version>3.0.1</version>
    </dependency>
    <dependency>
        <groupId>com.github.davidb</groupId>
        <artifactId>metrics-influxdb</artifactId>
        <version>0.8.2</version>
    </dependency>
```
* 核心工具类

```html
import com.codahale.metrics.*;
import com.ktyh.think.filter.prometheus.DropwizardExportSetUp;
import io.prometheus.client.CollectorRegistry;
import metrics_influxdb.InfluxdbReporter;
import metrics_influxdb.api.protocols.InfluxdbProtocols;
import org.slf4j.LoggerFactory;
import javax.servlet.ServletRequest;
import javax.servlet.ServletResponse;
import java.util.concurrent.TimeUnit;
/**
 * Created by zhangjianxin on 2016/11/3.
 * metrics core config and static void formart
 */
public class MetricsCore {
    /**
     * 因数据存在争议，改为单例模式。
     */
    private static final MetricRegistry METRIC_REGISTRY = new MetricRegistry();
    private static String HOST_IP_ADDRESS;
    private static String USERNAME;
    private static String PASSWORD;
    private static String DATABASE;
    private static int PORT;
    private static ServletRequest ARG_REQ;
    private static ServletResponse ARG_RESP;
    /**
     * @author 内存使用程度MB
     */
    Gauge<Long> getFreeMemory = () -> toMb(Runtime.getRuntime().freeMemory());
    Gauge<Long> getTotalMemory = () -> toMb(Runtime.getRuntime().totalMemory());

    public static synchronized MetricRegistry getInstance() {
        return METRIC_REGISTRY;
    }

    public static ScheduledReporter influxdbReporter(MetricRegistry metrics) {
        InitInfluxdbConfig();//初始化
        return InfluxdbReporter.forRegistry(metrics)
                .protocol(InfluxdbProtocols.http(HOST_IP_ADDRESS, PORT, USERNAME, PASSWORD, DATABASE))
                .convertRatesTo(TimeUnit.SECONDS)
                .convertDurationsTo(TimeUnit.MILLISECONDS)
                .filter(MetricFilter.ALL)
                .skipIdleMetrics(false)
                .build();
    }


    public static void freeMemoryForJvm(MetricRegistry metrics) {
        metrics.register(MetricRegistry.name("memory_free_mb_size"), new Gauge<Long>() {
            @Override
            public Long getValue() {
                return toMb(Runtime.getRuntime().freeMemory());
            }
        });
    }

    public static void totalMemoryForJvm(MetricRegistry metrics) {
        metrics.register(MetricRegistry.name("memory_total_mb_size"), new Gauge<Long>() {
            @Override
            public Long getValue() {
                return toMb(Runtime.getRuntime().totalMemory());
            }
        });
    }

    //转换为MB
    private static long toMb(long bytes) {
        return bytes / 1024 / 1024;
    }

    /**
     * 初始化时序数据库配置，可以通过反射进行更改
     **/

    public static void InitInfluxdbConfig() {
        PORT = 8085;
        HOST_IP_ADDRESS = "192.168.1.102";
        USERNAME = "root";
        PASSWORD = "root";
        // DATABASE = "Metrics";
        DATABASE = "_internal";
    }

    /**
     * 传递初始化参数，用单例去回调
     *
     * @param "Servlet request and respone"
     */

    public static void requestFullArgs(ServletRequest req, ServletResponse resp) {
        ARG_REQ = req;
        ARG_RESP = resp;
    }

    /**
     * Reporter 数据的展现位置
     *
     * @param metrics
     * @return
     */

    public ConsoleReporter consoleReporter(MetricRegistry metrics) {
        return ConsoleReporter.forRegistry(metrics)
                .convertRatesTo(TimeUnit.SECONDS)
                .convertDurationsTo(TimeUnit.MILLISECONDS)
                .build();
    }

    public Slf4jReporter slf4jReporter(MetricRegistry metrics) {
        return Slf4jReporter.forRegistry(metrics)
                .outputTo(LoggerFactory.getLogger("demo.metrics"))
                .convertRatesTo(TimeUnit.SECONDS)
                .convertDurationsTo(TimeUnit.MILLISECONDS)
                .build();
    }

    public JmxReporter jmxReporter(MetricRegistry metrics, String doMain) {
        return JmxReporter.forRegistry(metrics)
                .inDomain(doMain)
                .convertRatesTo(TimeUnit.SECONDS)
                .convertDurationsTo(TimeUnit.MILLISECONDS)
                .filter(MetricFilter.ALL).build();
    }

    /**
     * Meters会将最近1分钟，5分钟，15分钟的TPS（每秒处理的request数）给打印出来，还有所有时间的TPS。
     *
     * @param metrics
     * @return
     */
    public Meter requestMeter(MetricRegistry metrics, String arg) {
        return metrics.meter(MetricRegistry.name("meter_event_" + arg));
    }

    /**
     * 直方图
     *
     * @param metrics
     * @return
     */

    public Histogram responseSizes(MetricRegistry metrics, String arg) {
        return metrics.histogram(MetricRegistry.name("size_histogram_" + arg));
    }

    /**
     * 计数器
     *
     * @param metrics
     * @return
     */

    public Counter pendingJobs(MetricRegistry metrics, String arg) {
        return metrics.counter(MetricRegistry.name("counter_Job_" + arg));
    }

    /**
     * 计时器
     *
     * @param metrics
     * @return
     */

    public Timer responsesExecureTime(MetricRegistry metrics, String arg) {
        return metrics.timer(MetricRegistry.name("response_timer_" + arg));
    }

    public Meter requestURLCount(MetricRegistry metrics, String arg) {
        return metrics.meter(MetricRegistry.name("request_URL_" + arg));
    }
}

```
##目前这种方式只是将数据推送到influxdb ，以及prometheus 。

1 数据采集

2 采集点(Filter,Controller,Service)

3 数据推送(influxdb,primetheus,jmx,console,cvs等。)