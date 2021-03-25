#! /bin/bash

cd $(dirname $0)

# 一键启动脚本

cd ../GoLogs
git pull origin master

cd ../GoLimiter
git pull origin master

cd ../GoWxWhoIsTheSpy
git pull origin master

./rerun.sh
