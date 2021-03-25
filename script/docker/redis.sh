#!/usr/bin/env bash

# latest:5.0.5
# docker pull redis

cd `dirname $0`

export DOCKER_ROOT=/Users/limengfan/Desktop/redis
#export DOCKER_ROOT=/home/lighthouse/redis

rm -rf $DOCKER_ROOT
mkdir -p $DOCKER_ROOT/data

docker rm -f rrr
docker run -p something:6379 --name rrr -v $DOCKER_ROOT/data:/data -d redis redis-server --appendonly yes --requirepass "something"
# docker exec -it rrr /bin/bash
