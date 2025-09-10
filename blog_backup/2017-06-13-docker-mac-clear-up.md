---
layout: post
title: docker clear up for mac 
categories: docker
description: bash shell docker clear up for mac 
keywords: docker
---


# 1、Main

> 清理mac上Docker占用的磁盘空间

## 1.1、内容：

```text
Docker依赖Linux系统的cgroup实现，在mac系统中运行的时候，Docker会启动一个虚拟机中的Linux内核，并在硬盘上放一个 qcow2 格式的磁盘镜像文件。这个文件会随着Docker的使用不断膨胀，即使删除不用的Docker Image和Container也不会缩小。我在测试完一个自动化工具的Dockerfile改写之后，Docker.qcow2文件就达到42GB了。

Google搜了一圏发现一个叫Théo Chamley的法国人写了一个脚本释放Docker.qcow2文件占用的空间。基本原理是用 docker save 命令保存要保留的image，然后关闭Docker，删除Docker.qcow2，再启动Docker，它会自动重建，最后用 docker load 命令恢复保留的image。
```


## 命令行

```bash
 sudo sh docker-clean.sh $(docker images | awk 'NR>1 {print $1":"$2}')
```


## 脚本

```bash
#!/bin/bash
# sudo sh docker-clean.sh $(docker images | awk 'NR>1 {print $1":"$2}')

IMAGES=$@

echo "This will remove all your current containers and images except for:"
echo ${IMAGES}
read -p "Are you sure? [yes/NO] " -n 1 -r
echo    # (optional) move to a new line
if [[ ! $REPLY =~ ^[Yy]$ ]]
then
    exit 1
fi


TMP_DIR=$(mktemp -d)

pushd $TMP_DIR >/dev/null

open -a Docker
echo "=> Saving the specified images"
for image in ${IMAGES}; do
	echo "==> Saving ${image}"
	tar=$(echo -n ${image} | base64)
	docker save -o ${tar}.tar ${image}
	echo "==> Done."
done

echo "=> Cleaning up"
echo -n "==> Quiting Docker"
osascript -e 'quit app "Docker"'
while docker info >/dev/null 2>&1; do
	echo -n "."
	sleep 1
done;
echo ""

echo "==> Removing Docker.qcow2 file"
rm ~/Library/Containers/com.docker.docker/Data/com.docker.driver.amd64-linux/Docker.qcow2

echo "==> Launching Docker"
open -a Docker
echo -n "==> Waiting for Docker to start"
until docker info >/dev/null 2>&1; do
	echo -n "."
	sleep 1
done;
echo ""

echo "=> Done."

echo "=> Loading saved images"
for image in ${IMAGES}; do
	echo "==> Loading ${image}"
	tar=$(echo -n ${image} | base64)
	docker load -q -i ${tar}.tar || exit 1
	echo "==> Done."
done

popd >/dev/null
rm -r ${TMP_DIR}
```



## 貌似更好

 * This patch fixes a bug which would prevent the online compaction from being enabled. Hopefully now it's possible to turn it on again. I installed this version and then enabled experimental compaction by running the following script:

```bash
#!/bin/bash

cd ~/Library/Containers/com.docker.docker/Data/database/
git checkout master
git reset --hard
mkdir -p com.docker.driver.amd64-linux/disk
echo 262144 > com.docker.driver.amd64-linux/disk/compact-after
echo 262144 > com.docker.driver.amd64-linux/disk/keep-erased
echo -n true > com.docker.driver.amd64-linux/disk/trim
#echo -n 'tcp:9090' > com.docker.driver.amd64-linux/disk/stats
git add com.docker.driver.amd64-linux/disk/compact-after
git add com.docker.driver.amd64-linux/disk/keep-erased
git add com.docker.driver.amd64-linux/disk/trim
#git add com.docker.driver.amd64-linux/disk/stats
git commit -s -m 'Enable on-line compaction'
```
* I then rebooted the app. I created a large temporary file in a container, deleted it, and around 15 minutes later an fstrim /var was triggered from cron and the file became smaller again. (It's also possible to trigger fstrim immediately with a command like docker run --rm --net=host --pid=host --privileged -it justincormack/nsenter1 /sbin/fstrim /var)

* In 17.07 (the upcoming edge release) I'm hoping to switch online compaction on by default. In 17.08 I'm hoping to automatically trigger fstrim when docker images are removed, instead of waiting for a cron job.
