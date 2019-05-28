#/bin/bash
#编译golang可执行文件
root_dir=$(cd "$(dirname "$0")"; cd ..; pwd)

export GO111MODULE=auto

#禁用cgo模块
export CGO_ENABLED=0
cd $root_dir
go build -a -installsuffix cgo -o go-api