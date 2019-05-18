#!/usr/bin/env bash
#初始化目录权限
root_dir=$(cd "$(dirname "$0")"; cd ..; pwd)

mkdir -p $root_dir/logs
chmod 777 -R $root_dir/logs
echo "初始化成功！"

exit 0
