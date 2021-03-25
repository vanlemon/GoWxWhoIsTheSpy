#! /bin/bash

cd $(dirname $0)

export GOPROXY=https://goproxy.io
#export GOPROXY=https://goproxy.cn
#export GOPROXY=https://mirrors.aliyun.com/goproxy/
export GO111MODULE=on

go mod tidy

go mod download