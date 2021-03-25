#!/usr/bin/env bash

# latest:5.8
# docker pull mysql

cd `dirname $0`

export DOCKER_ROOT=/Users/limengfan/Desktop/mysql
#export DOCKER_ROOT=/home/lighthouse/mysql
rm -rf $DOCKER_ROOT
mkdir -p $DOCKER_ROOT/conf $DOCKER_ROOT/logs $DOCKER_ROOT/data

export SQL_PATH=$PWD/../sql

docker rm -f mmm
docker run -p something:3306 --name mmm -v $DOCKER_ROOT/conf:/etc/mysql/conf.d -v $DOCKER_ROOT/logs:/logs -v $DOCKER_ROOT/data:/var/lib/mysql -v $SQL_PATH:/sql -e MYSQL_ROOT_PASSWORD=something -d mysql
# docker exec -it mmm /bin/bash
docker exec -i mmm /bin/bash << EOF
# 等待mysql启动，等待文件挂载
echo "wait..."
sleep 30s
echo "init"
cat /sql/init.sql /sql/user.sql > /sql/all.sql
mysql -uroot -something -h127.0.0.1 -P3306 < /sql/all.sql
echo "finish"
EOF
exit

#mysql 日志
#
#```
#show global variables like "%genera%"
#set global general_log = on;
#/var/lib/mysql/73529d1b3007.log
#```
