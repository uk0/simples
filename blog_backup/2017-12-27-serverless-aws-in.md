---
layout: post
categories: aws
title: Aws Serverless for lambda Test
date: 2017-12-27 14:20:17 +0800
description: Serverless test to Aws lambda 
keywords: lambda serverless aws kubeless
---



# Aws Lambda to Serverless 

![](http://112firshme11224.test.upcdn.net/demos/a43efa7f-f810-4164-970f-8bd1f309650c.png)



# aws lambda function 编辑器

![](http://112firshme11224.test.upcdn.net/demos/8d472d40-ff9c-4108-91be-1364307b5599.png)


# 编辑器的console 以及 测试按钮，语言版本

![](http://112firshme11224.test.upcdn.net/demos/111c3da3-cb1e-48de-b3b7-c286d9fc55e2.png)


# 常用的配置项

![](http://112firshme11224.test.upcdn.net/demos/89ba65d2-c46a-434e-bb96-afa15d02c6f5.png)



![](http://112firshme11224.test.upcdn.net/demos/bfa64818-bf58-45fa-8e2a-51405c187f7e.png)

##  配置AWS IAM执行策略

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "VisualEditor0",
            "Effect": "Allow",
            "Action": [
                "s3:*",
                "iam:*",
                "apigateway:*",
                "cloudformation:ListStacks",
                "logs:*",
                "cloudformation:ValidateTemplate",
                "cloudformation:DescribeStackEvents",
                "cloudformation:CreateStack",
                "lambda:*",
                "cloudformation:UpdateStack",
                "cloudformation:DescribeStackResource",
                "cloudformation:CreateChangeSet",
                "cloudformation:DescribeChangeSet",
                "cloudformation:ExecuteChangeSet",
                "cloudformation:DescribeStacks"
            ],
            "Resource": "*"
        }
    ]
}
```

# demo


![](http://112firshme11224.test.upcdn.net/demos/070f4afd-25c7-4783-a9e3-e77f133daad8.gif)



转载请注明出处，本文采用 [CC4.0](http://creativecommons.org/licenses/by-nc-nd/4.0/) 协议授权
