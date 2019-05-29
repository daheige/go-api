# gin框架实战
    基于gin框架封装而成的mvc框架，可用于go api开发。
# 目录结构
    .
    ├── app
    │   ├── controller  控制器
    │   ├── logic       业务逻辑层
    │   ├── middleware  中间件层
    │   └── routes      路由层设置
    ├── app.yaml        配置文件
    ├── config          配置文件设置
    │   └── bootstartp.go
    ├── go.mod          go.mod
    ├── go.sum
    ├── LICENSE
    ├── logs            日志目录，可以自定义到别的路径中
    │   └── app.2019-05-06.log
    ├── main.go         程序入口文件
# 关于redisgo调优
    区分两种使用场景：
    1.高频调用的场景，需要尽量压榨redis的性能： 
        调高MaxIdle的大小，该数目小于maxActive，由于作为一个缓冲区一样的存在
        扩大缓冲区自然没有问题，调高MaxActive，考虑到服务端的支持上限，尽量调高
        IdleTimeout由于是高频使用场景，设置短一点也无所谓，需要注意的一点是MaxIdle
        设置的长了队列中的过期连接可能会增多，这个时候IdleTimeout也要相应变化
    2.低频调用的场景，调用量远未达到redis的负载，稳定性为重： 
        MaxIdle可以设置的小一些
        IdleTimeout相应地设置小一些
        MaxActive随意，够用就好，容易检测到异常

# docker运行
    1.构建golang二进制文件
        $ sh bin/app-build
    2.构建docker镜像
        $ docker build -t go-api:v1 .
    3.运行docker容器
        $ docker run -it -d -p 1338:1338 --name=go-api-server -v /web/go/go-api/logs:/go/logs go-api:v1
    4.访问localhost:1338，查看页面
    
    可以通过以下方式运行
    sudo mkdir -p /data/www/go-api/logs
    sudo mkdir -p /data/www/go-api/conf
    sudo cp app.yaml /data/www/go-api/conf/
    sudo chown -R $USER /data/www/go-api/
    docker run -it -d -p 1336:1338 -v /data/www/go-api/logs:/go/logs -v /data/www/go-api/conf:/go/conf go-api:v1
    这样就可以在任意目录中运行docker容器
# 版权
    MIT
