#! /bin/bash

cd $(dirname $0)

# 杀死服务
sudo netstat -lnp|grep 9205

# sudo kill -9 加进程ID
sudo kill -9 $1

