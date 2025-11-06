# x-ui 分流


> 前提HK作为最终入口

* HK部署X-UI （日常Google查找文献）
* US部署X-UI （用于ChatGPT）
* 备份HK节点的配置文件

* HK 节点面板JSON配置如下(全量JSON)：

```json

{
  "api": {
    "services": [
      "HandlerService",
      "LoggerService",
      "StatsService"
    ],
    "tag": "api"
  },
  "inbounds": [
    {
      "listen": "127.0.0.1",
      "port": 62789,
      "protocol": "dokodemo-door",
      "settings": {
        "address": "127.0.0.1"
      },
      "tag": "api"
    }
  ],
  "outbounds": [
    {
      "protocol": "vless",
      "settings": {
        "vnext": [
          {
            "address": "1.1.1.1",
            "port": 55555,
            "users": [
              {
                "id": "xxxxxxxxxxxxxxxxxxx",
                "encryption": "none"
              }
            ]
          }
        ]
      },
      "streamSettings": {
        "network": "kcp",
        "security": "none",
        "kcpSettings": {
          "mtu": 1350,
          "tti": 50,
          "uplinkCapacity": 12,
          "downlinkCapacity": 100,
          "congestion": false,
          "header": {
            "type": "srtp"
          },
          "seed": "asdasd"
        }
      },
      "mux": {
        "enabled": false
      },
      "tag": "chatgpt-outbound"
    },
    {
      "protocol": "freedom",
      "settings": {}
    },
    {
      "protocol": "blackhole",
      "settings": {},
      "tag": "blocked"
    }
  ],
  "policy": {
    "system": {
      "statsInboundDownlink": true,
      "statsInboundUplink": true
    }
  },
  "routing": {
    "rules": [
      {
        "type": "field",
        "domain": [
          "geosite:netflix",
          "chatgpt.com",
          "openai.com",
          "oaistatic.com",
          "oaiusercontent.com",
          "openai.com.cdn.cloudflare.net",
          "openaiapi-site.azureedge.net",
          "openaicom-api-bdcpf8c6d2e9atf6.z01.azurefd.net",
          "openaicomproductionae4b.blob.core.windows.net",
          "production-openaicom-storage.azureedge.net",
          "claude.ai",
          "anthropic.com"
        ],
        "outboundTag": "chatgpt-outbound"
      },
      {
        "inboundTag": [
          "api"
        ],
        "outboundTag": "api",
        "type": "field"
      },
      {
        "ip": [
          "geoip:private"
        ],
        "outboundTag": "blocked",
        "type": "field"
      },
      {
        "outboundTag": "blocked",
        "protocol": [
          "bittorrent"
        ],
        "type": "field"
      }
    ]
  },
  "stats": {}
}

```


* HK节点分段JSON解释

> 在outbounds下面添加以下内容即可

```json
{
      "protocol": "vless",
      "settings": {
        "vnext": [
          {
            "address": "1.1.1.1",
            "port": 55555,
            "users": [
              {
                "id": "xxxxxxxxxxxxxxxxxxx",
                "encryption": "none"
              }
            ]
          }
        ]
      },
      "streamSettings": {
        "network": "kcp",
        "security": "none",
        "kcpSettings": {
          "mtu": 1350,
          "tti": 50,
          "uplinkCapacity": 12,
          "downlinkCapacity": 100,
          "congestion": false,
          "header": {
            "type": "srtp"
          },
          "seed": "asdasd"
        }
      },
      "mux": {
        "enabled": false
      },
      "tag": "chatgpt-outbound"
    }
```

>在routing的rules下面添加以下内容即可

```json

{
        "type": "field",
        "domain": [
          "geosite:netflix",
          "chatgpt.com",
          "openai.com",
          "oaistatic.com",
          "oaiusercontent.com",
          "openai.com.cdn.cloudflare.net",
          "openaiapi-site.azureedge.net",
          "openaicom-api-bdcpf8c6d2e9atf6.z01.azurefd.net",
          "openaicomproductionae4b.blob.core.windows.net",
          "production-openaicom-storage.azureedge.net",
          "claude.ai",
          "anthropic.com"
        ],
        "outboundTag": "chatgpt-outbound"
      }
```



* US节点只需要创建好`VLESS`节点即可。

* 现在使用HK也可以访问ChatGPT了。