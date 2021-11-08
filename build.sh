#!/bin/bash

export GO111MODULE=on
export GOPROXY=https://goproxy.cn
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o wilsonchai-webhook

docker build -t registry.cn-beijing.aliyuncs.com/wilsonchai/mutating-webhook:$1 .
docker push registry.cn-beijing.aliyuncs.com/wilsonchai/mutating-webhook:$1
