#!/bin/bash


APP_NAME="blog"
BUILD_TIME=$(date +%Y%m%d-%H%M%S)
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")




GOEXPERIMENT=greenteagc CGO_ENABLED=0 GOOS=linux  GOARCH=amd64 go build  -ldflags="-s -w -funcalign=32 -X main.Version=${BUILD_TIME} \
                                                                                                                -X main.GitCommit=${GIT_COMMIT}" \
                                                                                                                -o ${APP_NAME}


 rsync -av \
--exclude='.git'  \
--exclude='.DS_Store'  \
 /Users/firshme/Desktop/work/simples/tpl root@vm2:/root/blog/


  rsync -av \
--exclude='.git'  \
--exclude='.DS_Store'  \
 /Users/firshme/Desktop/work/simples/blog root@vm2:/root/blog/





