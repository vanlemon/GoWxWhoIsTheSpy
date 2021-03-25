#! /bin/bash

cd $(dirname $0)

# 重新构建并运行服务
sudo ./build.sh
sudo ./run.sh $1
