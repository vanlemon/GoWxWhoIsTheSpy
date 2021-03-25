#! /bin/bash

cd $(dirname $0)

RUN_NAME="lmf.mortal.spy"

time=$(date "+%Y-%m-%d %H:%M:%S")

mv output "output_${time}" # 伪删除 output
mkdir -p output/bin output/conf # 重新创建 output
cp script/bootstrap.sh output # 赋值执行脚本
chmod +x output/bootstrap.sh # 执行脚本权限为可执行
find conf/ -type f ! -name "*_local.*" | xargs -I{} cp {} output/conf/ # 赋值所有配置文件

go build -a -o output/bin/${RUN_NAME} # 构建可执行文件
