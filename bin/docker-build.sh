#/bin/bash
#编译golang可执行文件
root_dir=$(cd "$(dirname "$0")"; cd ..; pwd)

goExecPath=`which go`
if [ -z $goExecPath ];then
    echo "请先安装golang1.12.5+版本，再运行"
fi

export GO111MODULE=auto

#禁用cgo模块
export CGO_ENABLED=0
cd $root_dir
$goExecPath build -a -installsuffix cgo -o go-api

#docker image版本
version=$1
if [ -z $version ];then
    version=v1
fi

cnt=`docker image ls | grep "go-api-server:$version" | wc -l`
if [ $cnt -gt 0 ];then
    #删除之前的image
    docker rmi -f go-api-server:$version
fi

#重新生成镜像
cd $root_dir
docker build -t go-api-server:$version .

echo "构建image go-api-server:$version 成功"

exit 0
