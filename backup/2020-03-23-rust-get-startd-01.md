---
layout: post
categories: Rust
title: Rust 第一个入门程序
date: 2020-03-23 00:35:55 +0800
description: rust入门介绍，以及第一个程序里面包含了rust多少设计思想
keywords:  rust rust入门程序 rustweb actix
---



记录第一个rust 程序，以及对用到的工具rust语言特性的介绍。


### Rust 官方文档

* https://doc.rust-lang.org/book/


### 安装rust cargo

* 安装 
    ```bash
    ## install rust
    curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
    ## 
    ```
* 官方安装文档

  > https://forge.rust-lang.org/infra/other-installation-methods.html





### 创建一个rust项目

```bash
cargo new web_demo  --bin


> cd web_demo
> tree .
.
├── Cargo.toml
└── src
    └── main.rs

1 directory, 2 files

```

Let’s take a closer look at `Cargo.toml`:

```toml
[package]
name = "web-demo"
version = "0.1.0"
authors = ["张建新 <zhangjianxin@bigsfintech.com>"]
edition = "2018"

# See more keys and their definitions at https://doc.rust-lang.org/cargo/reference/manifest.html

[dependencies]

# 下面是我做测试添加的依赖
actix-web = "2.0"
actix-rt = "1.0"
actix-files = "0.2.1"
actix-session = "0.3.0"
serde_json = "1.0"
serde="1.0.104"
strfmt="0.1.6"
kv = "0.20.2"
base64 = "0.12.0"
```

This is called a manifest, and it contains all of the metadata that Cargo needs to compile your package.

Here’s what’s in `src/main.rs`:

```rust
fn main() {
    println!("Hello, world!");
}
```

### 开始添砖加瓦

```rust


use actix_files;
use actix_web::http::StatusCode;
use actix_web::{
    guard, middleware, web, App, HttpRequest, HttpResponse, HttpServer, Responder, Result,
};
use kv::*;
use std::collections::HashMap;
use std::{fs, path};
use std::str;
use strfmt::strfmt;
use actix_files::NamedFile;
use std::path::PathBuf;
use base64::{decode, encode};
use serde::{Deserialize, Serialize};


#[derive(Serialize, Deserialize)]
struct ResultResp {
    msg: String,
    code: String,
    data: String,
    data_array :Vec<String>,
}


#[derive(Debug, serde::Serialize, serde::Deserialize, PartialEq)]
struct KvS {
    url: String,
    data: String,
}


#[derive(Debug, serde::Serialize, serde::Deserialize, PartialEq)]
struct History {
    url: String,
    s: String,
    r: String,
    time: String,
    ip: String,
}


#[derive(Deserialize)]
struct RequestData {
    help: bool,
    data: String,
    exec: String,
    args: Vec<String>,
}
#[derive(Deserialize)]
struct RequestDataAddurl {
    help: bool,
    data: String,
    exec: String,
    url: String,
    args: Vec<String>,
}



fn reset(name: &str) -> String {
    let s = format!("./stroe/{}", name);
    //let _ = fs::remove_dir_all(&s);
    s
}



async fn hello() -> impl Responder {
    "hello"
}

//* 修改模板 以及保存模板
async fn modify_template(info: web::Json<RequestData>) -> impl Responder {
    let mut resp = ResultResp {
        msg: "data == null is error ".to_owned(),
        code: format!("Use {} ", info.exec.to_owned()),
        data: "".to_owned(),
        data_array: Vec::new(),
    };

    //if tos::template::CMD_LIST.contains(&&*info.exec) {
        println!("{}", "success".to_owned());

        put_k(info.exec.to_owned(), info.data.to_owned());

        resp.data.push_str("template add success");
    //}

    let j = serde_json::ser::to_string(&resp);
    HttpResponse::Ok()
        .content_type("application/json")
        .header("X-Hdr", "sample");

    j
}

//* 获取单个模板信息
async fn get_template(info: web::Json<RequestData>) -> impl Responder {
    let mut resp = ResultResp {
        msg: "data == null is error ".to_owned(),
        code: format!("Use {} ", info.exec.to_owned()),
        data: "".to_owned(),
        data_array: Vec::new(),
    };

    //if tos::template::CMD_LIST.contains(&&*info.exec) {
        println!("{}", "success".to_owned());
        let mut x =  String::from("");
        match get_k(info.exec.to_owned()) {
            Some(v) => x = v,
            None => panic!(),
        }
        let b = base64::decode(x).unwrap();

        let s = str::from_utf8(&b).unwrap();

        resp.data.push_str(s);
    //}

    let j = serde_json::ser::to_string(&resp);
    HttpResponse::Ok()
        .content_type("application/json")
        .header("X-Hdr", "sample");

    j
}

//* iter keys

fn iter_k() ->Vec<String>{
    let path = reset("kvs");
    let cfg = Config::new(path.clone());
    let store = Store::new(cfg).unwrap();
    let test_iter = store.bucket::<String,String>(Some("stroe")).unwrap();
    let mut v: Vec<String> = vec![];
    for item in test_iter.iter() {
        let item = item.unwrap();
        let key: String = item.key().unwrap();
        let value = item.value::<String>().unwrap();
        println!("key: {}, value: {}", key, value);
        // let b;
        // match  base64::decode(value) {
        //     Ok(v) => b =v,
        //     Err(e) => panic!(),
        // }
        // let s = str::from_utf8(&b).unwrap();
        v.push(key.to_owned());
    }
    v
}

//* 获取所有模板 信息 列表JSON
async fn get_templates() -> impl Responder {
    let  mut resp = ResultResp {
        msg: "data == null is error ".to_owned(),
        code: "0".to_owned(),
        data: "".to_owned(),
        data_array: Vec::new(),
    };
    println!("{}", "success".to_owned());

    resp.data_array = iter_k();

    let j = serde_json::ser::to_string(&resp);
    HttpResponse::Ok()
        .content_type("application/json")
        .header("X-Hdr", "sample");

    j
}

//* 发起请求
/// extract `Info` using serde
async fn start_http_client_post(info: web::Json<RequestData>) -> Result<String> {
    if info.help {
        let message = format!("{}", tos::template::TRACKING_HELP);
        return Result::Ok(message);
    }
    if tos::template::CMD_LIST.contains(&&*info.exec) {
        println!("{}", info.args[0]);
        let mut vars = HashMap::new();
        vars.insert("timestamp".to_string(), info.args[0].to_string());
        vars.insert("tcuId".to_string(), info.args[1].to_string());
        vars.insert("vin".to_string(), info.args[2].to_string());

        let mut x =String::from("");

        match get_k(info.exec.to_owned()) {
            Some(v) => x = v,
            None => panic!(),
        }
        let b = base64::decode(x).unwrap();

        let s = str::from_utf8(&b).unwrap();
        println!("{}", s);

        return Result::Ok(format!("{}", strfmt(&s.to_owned(), &vars).unwrap()));
    }

    return Result::Ok(format!("not found cmd  = {}", info.exec));
}

#[warn(unused_mut)]
fn put_k(K: String, V: String) -> Result<(), Error> {
    // Configure the database
    let path = reset("kvs");
    let cfg = Config::new(path.clone());
    let store = Store::new(cfg).unwrap();
    // raw Integer String
    let test = store.bucket::<String, String>(Some("stroe"))?;
    // Set testing = 123

    test.set(K, V)?;

    let temp_get = test.get("TrackingCmd").unwrap();
    println!("{:?}", temp_get);
    Ok(())
}



//* get kv stroe by 'K

fn get_k(K: String) -> Option<String> {
    // Configure the database
    let path = reset("kvs");
    let cfg = Config::new(path.clone());
    let store = Store::new(cfg).unwrap();

    let test_get = store.bucket::<String, String>(Some("stroe")).unwrap();

    let c = test_get.get(K.to_owned()).unwrap();

    c
}

async fn static_index(req: HttpRequest) -> Result<NamedFile> {
    Ok(actix_files::NamedFile::open(
        "../static/index.html",
    )?)
}

#[actix_rt::main]
async fn main() -> std::io::Result<()> {
    use actix_web::{web, App, HttpServer};

    HttpServer::new(|| {
        App::new()
            .wrap(middleware::Compress::default())
            .data(web::JsonConfig::default().limit(4096)) // <- limit size of the payload (global configuration)
            .service(web::resource("/addTask").route(web::post().to(start_http_client_post)))
            .service(web::resource("/hello").route(web::get().to(hello)))
            .service(web::resource("/addTemplate").route(web::post().to(modify_template)))
            .service(web::resource("/getTemplate").route(web::post().to(get_template)))
            .service(web::resource("/getTemplateList").route(web::get().to(get_templates)))
            .service(actix_files::Files::new(
                "/static",
                "../static/",
            ))
            .default_service(
                web::resource("").route(web::get().to(static_index)).route(
                    web::route()
                        .guard(guard::Not(guard::Get()))
                        .to(|| HttpResponse::MethodNotAllowed()),
                ),
            )
    })
    .bind("127.0.0.1:8088")?
    .run()
    .await
}

```


### 写一个简单的web遇到哪些感觉有意思的东西

* match 大体功能同等于`switch ` 但是加强了枚举
* Some  `一个神奇的枚举种类` （比较有意思的点）
* `|| ` `闭包`
* ? `Error propagation	一般会和 Match Some Restlt<Option> ` 打一套组合拳 主要还是用于传播 （比较有意思的点）
* unwrap  和  `?` 相反的工作 `?`是传播  `unwrap` 摊开数据 解包数据返回。
*  async await Rust的异步并发编程

```rust
async fn print_async() {
     println!("Hello from print_async")
}

fn main() {
     let future = print_async();
     println!("Hello from main");
     futures::block_on(future);
}
// 本质上是 延迟运算
// 在手动 poll future 之前，异步函数内部一行代码也不会执行
```
* await! `是一个编译器内建的 macro ，用来暂停（pause） Future 的执行流程，并且把执行流程交回给调用方 (yield) 和React访问rest接口 一致。`



#### 今天是学习Rust的第一天 开阔了视野

* 文章不定时更新，遇到自己觉得有意思的问题 第一时间更新上来。



转载请注明出处，本文采用 [CC4.0](http://creativecommons.org/licenses/by-nc-nd/4.0/) 协议授权
