---
layout: post
title: Mongodb Aggregate Sum。
categories: mongo
description: Mongodb sum 查询聚合
keywords: mongo, sum
---

# Mongodb search aggregate SUM

 * find aggregate and where

```bash

db.cabinet.aggregate(
   [
     {
       $group:
         {
           _id: { },
           totalAmount: { $sum: "$cabinetBoxTotal" },
           count: { $sum: 1 }
         }
     }
   ]
)

db.cabinet.aggregate(
   [
       { $match : { cabinetName : "中国农业大学西校区2号柜"} },
  
     {
       $group:
         {
           _id: {},
           totalAmount: { $sum: "$cabinetBoxTotal" },
           count: { $sum: 1 }
         }
     }
   ]
)

```
