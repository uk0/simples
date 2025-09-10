---
layout: post
title: ehcache集群配置terracotta。
categories: mongo
description: ehcache集群配置terracotta
keywords: ehcache, terracotta
---

# ehcache集群配置terracotta

 * ehcache for maven conf

```xml
        <dependency>
            <groupId>org.ehcache</groupId>
            <artifactId>ehcache</artifactId>
            <version>${ehcache.version}</version>
        </dependency>
        <dependency>
            <groupId>org.ehcache</groupId>
            <artifactId>ehcache-clustered</artifactId>
            <version>${ehcache.version}</version>
        </dependency>

```

* tc-config.xml (ehcache3.4.0)版本已经被合并了里面包含`terracotta`

```xml
<?xml version="1.0" encoding="UTF-8"?>
<tc-config xmlns="http://www.terracotta.org/config"
           xmlns:ohr="http://www.terracotta.org/config/offheap-resource">
  <!--
    This is the default Terracotta server configuration file for the Ehcache kit.

    It defines a single offheap resource of 512MB to be used for caching data.

    It also defines a single server, but you can add another one to benefit from high availability.
  -->

  <plugins>
    <config>
      <ohr:offheap-resources>
        <ohr:resource name="main" unit="MB">512</ohr:resource>
      </ohr:offheap-resources>
    </config>
  </plugins>

  <servers>
    <server host="10.11.0.224" name="s1">
      <!--
        Indicates the location for logs files - %H will resolve to user home directory.
        Note that relative path will be resolved from the location of this configuration file.
       -->
      <logs>/home/hadmin/data/logs/ehcache/ehcache-server/terracotta-logs-s1</logs>

      <!--
        This port is used by clients to communicate to the server.
        Its value is actually the default one and is thus omitted.
      -->
      <tsa-port>9410</tsa-port>

      <!--
        This port is used for server to server communication.
        Its value is actually the default one and is thus omitted.
      -->
      <tsa-group-port>9430</tsa-group-port>
    </server>
    <server host="10.11.0.225" name="s2">
      <!--
                   Indicates the location for logs files - %H will resolve to user home directory.
        Note that relative path will be resolved from the location of this configuration file.
       -->
      <logs>/home/hadmin/data/logs/ehcache/ehcache-server/terracotta-logs-s2</logs>

      <!--
                   This port is used by clients to communicate to the server.
        Its value is actually the default one and is thus omitted.
      -->
      <tsa-port>9410</tsa-port>

      <!--
                   This port is used for server to server communication.
        Its value is actually the default one and is thus omitted.
      -->
      <tsa-group-port>9430</tsa-group-port>
    </server>

    <!--
      Below a sample server definition that will give HA to the cluster if run on a different host.

      Servers know how to communicate between each others because each version of the config file on each host
      will list all servers in it.
    -->
    <!--<server host="otherhost" name="other-server">-->
      <!--<logs>logs</logs>-->
      <!--<tsa-port>9410</tsa-port>-->
    <!--</server>-->

    <!--
      Indicates how much time a server taking over after a failure in an active will wait
      to allow existing clients to reconnect. The time unit is seconds.
    -->
    <client-reconnect-window>120</client-reconnect-window>
  </servers>
</tc-config>

```


* 以上的配置文件为两个节点


* 分别在两个节点上运行相应命令。

```bash
# 1
nohup ./bin/start-tc-server.sh  -f conf/tc-config.xml  -n s1 &
# 2
nohup ./bin/start-tc-server.sh  -f conf/tc-config.xml  -n s2 &

```

* 这样两台基本的Ehcache 服务集群就搭建成功了。

* 文章写的比较着急，不忙了在继续整理，望理解。

* 测试代码


```java


package todcloud.utils.ehcache;

/**
 * Created by zhangjianxin on 2017/12/4.
 * Github Breakeval13
 * blog firsh.me
 */
import java.net.URI;

import org.ehcache.Cache;
import org.ehcache.CacheManager;
import org.slf4j.Logger;

import static java.net.URI.create;
import static org.ehcache.clustered.client.config.builders.ClusteredResourcePoolBuilder.clusteredDedicated;
import static org.ehcache.clustered.client.config.builders.ClusteringServiceConfigurationBuilder.cluster;
import static org.ehcache.config.builders.CacheConfigurationBuilder.newCacheConfigurationBuilder;
import static org.ehcache.config.builders.CacheManagerBuilder.newCacheManagerBuilder;
import static org.ehcache.config.builders.ResourcePoolsBuilder.heap;
import static org.ehcache.config.units.MemoryUnit.MB;
import static org.slf4j.LoggerFactory.getLogger;

public class ClusteredProgrammatic {
    private static final Logger LOGGER = getLogger(ClusteredProgrammatic.class);

    public static void main(String[] args) {
        LOGGER.info("Creating clustered cache manager");
        final URI uri = create("terracotta://10.11.0.224:9410/s1");
        try (CacheManager cacheManager = newCacheManagerBuilder()
                .with(cluster(uri).autoCreate().defaultServerResource("main"))
                .withCache("basicCache",
                        newCacheConfigurationBuilder(Long.class, String.class,
                                heap(100).offheap(1, MB).with(clusteredDedicated(512, MB))))
                .build(true)) {
            Cache<Long, String> basicCache = cacheManager.getCache("basicCache", Long.class, String.class);

            LOGGER.info("Putting to cache");
//            basicCache.put(100L, "da one!");
//            basicCache.put(1L, "da one1L!");
            System.out.println(basicCache.get(1L));
            LOGGER.info("Closing cache manager");
        }

        LOGGER.info("Exiting");
    }
}

```