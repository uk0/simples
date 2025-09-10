# Makefile

.PHONY: build

# 1) 构建目标
build:
	@echo "[INFO]  build GoServer ..."
	CGO_ENABLED=0 GOOS=linux  GOARCH=amd64 go build -o blog
