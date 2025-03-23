#!/bin/bash

OS_NAME=Linux
# 删除旧的可执行文件
if [ -f server ]; then
    rm server
fi

echo "Operating System: $OS_NAME"
# 根据操作系统运行应用
case "$OS_NAME" in
    Darwin)
        echo "Building on macOS..."
        go build -ldflags "-s -w" -o oneapi
        ;;
    Linux)
        echo "Running on Linux..."
        #docker exec -itd gva-server /bin/sh /go/src/start_server.sh

        ;;
    *)
        echo "Unsupported OS: $OS_NAME"
        exit 1
        ;;
esac