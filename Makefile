# Makefile
APP_NAME="blog"
BUILD_TIME=$(date +%Y%m%d-%H%M%S)
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")


.PHONY: build

# 1) 构建目标
build:
	@echo "[INFO]  build GoServer ..."
	echo "🚀 Building ${APP_NAME} with Go 1.25..."
	GOEXPERIMENT=greenteagc CGO_ENABLED=0 GOOS=linux  GOARCH=amd64 go build  -ldflags="-s -w -funcalign=32 -X main.Version=${BUILD_TIME} \
                                                                                                                -X main.GitCommit=${GIT_COMMIT}" \
                                                                                                                -o ${APP_NAME}
