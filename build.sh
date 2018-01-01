#!/bin/sh

set -e
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
	&& docker build -t registry.cn-hangzhou.aliyuncs.com/deepzz/backup . \
        && docker push registry.cn-hangzhou.aliyuncs.com/deepzz/backup
