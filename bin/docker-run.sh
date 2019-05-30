#!/usr/bin/env bash
#初始化目录权限
root_dir=$(cd "$(dirname "$0")"; cd ..; pwd)

#docker数据卷映射到主机的目录
workdir=$HOME

#docker image版本
version=$1
if [ -z $version ];then
    version=v1
fi

#构建镜像
sh $root_dir/bin/docker-build.sh $version
echo "构建image完成"

#docker容器name名称
containerName=$2
if [ -z $containerName ];then
    containerName=go-api
fi

#创建docker映射到当前主机上的目录
mkdir -p $workdir/go-api/conf
mkdir -p $workdir/go-api/logs
chmod 777 -R $workdir/go-api/logs
cp $root_dir/app.yaml $workdir/go-api/conf/

#停止之前的容器
cnt=`docker container ls -a | grep $containerName | grep -v grep | wc -l`
if [ $cnt -gt 0 ];then
    docker stop $containerName
    docker rm $containerName
fi

#运行容器
docker run -it -d -p 1338:1338 --name $containerName -v $workdir/go-api/logs:/go/logs -v $workdir/go-api/conf:/go/conf go-api-server:$version

echo "docker运行go-api成功!"
echo "访问localhost:1338 开始你的应用之旅吧!"

exit 0
