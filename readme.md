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

# golang环境安装
    1、linux环境，下载go1.12.5.linux-amd64.tar.gz
        cd /usr/local/
        sudo wget https://dl.google.com/go/go1.12.5.linux-amd64.tar.gz
        sudo tar zxvf go1.12.5.linux-amd64.tar.gz
        创建golang需要的目录
        sudo mkdir /mygo
        sudo mkdir /mygo/bin
        sudo mkdir /mygo/src
        sudo mkdir /mygo/pkg
    2、设置环境变量vim ~/.bashrc 或者sudo vim /etc/profile
        export GOROOT=/usr/local/go
        export GOOS=linux
        export GOPATH=/mygo
        export GOSRC=$GOPATH/src
        export GOBIN=$GOPATH/bin
        export GOPKG=$GOPATH/pkg
        #开启go mod机制
        export GO111MODULE=auto

        #禁用cgo模块
        export CGO_ENABLED=0

        export PATH=$GOROOT/bin:$GOBIN:$PATH
    3、source ~/.bashrc 生效配置
# 设置goproxy代理
    vim ~/.bashrc添加如下内容：
    export GOPROXY=https://goproxy.io
    或者使用 export GOPROXY=https://athens.azurefd.net
    让bashrc生效
    source ~/.bashrc
  
# 开始运行
    go mod tidy #安装golang module包
    go run main.go
    访问localhost:1338

# 线上部署
    方法1：
        请用supervior启动二进制文件，参考go-api.ini文件
    方法2：
        采用docker运行二进制文件
        
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

    性能监控
        浏览器访问http://localhost:2338/debug/pprof，就可以查看
    在命令终端查看：
        查看profile
            go tool pprof http://localhost:2338/debug/pprof/profile?seconds=60
            (pprof) top 10 --cum --sum

            每一列的含义：
            flat：给定函数上运行耗时
            flat%：同上的 CPU 运行耗时总比例
            sum%：给定函数累积使用 CPU 总比例
            cum：当前函数加上它之上的调用运行总耗时
            cum%：同上的 CPU 运行耗时总比例

        它会收集30s的性能profile,可以用go tool查看
            go tool pprof profile /home/heige/pprof/pprof.go-api.samples.cpu.002.pb.gz
        查看heap和goroutine
            查看活动对象的内存分配情况
            go tool pprof http://localhost:2338/debug/pprof/heap
            go tool pprof http://localhost:2338/debug/pprof/goroutine
        
        prometheus性能监控
        http://localhost:2338/metrics

# wrk工具压力测试
    https://github.com/wg/wrk
    
    ubuntu系统安装如下
    1、安装wrk
        # 安装 make 工具
        sudo apt-get install make git
        
        # 安装 gcc编译环境
        sudo apt-get install build-essential
        sudo mkdir /web/
        sudo chown -R $USER /web/
        cd /web/
        git clone https://github.com/wg/wrk.git
        # 开始编译
        cd /web/wrk
        make
    2、wrk压力测试
        $ wrk -c 100 -t 8 -d 2m http://localhost:1338/index
        Running 2m test @ http://localhost:1338/index
        8 threads and 100 connections
        Thread Stats   Avg      Stdev     Max   +/- Stdev
            Latency    19.50ms   40.88ms 829.98ms   96.82%
            Req/Sec     0.89k   166.70     1.68k    71.41%
        829464 requests in 2.00m, 118.66MB read
        Socket errors: connect 0, read 0, write 0, timeout 96
        Requests/sec:   6911.09
        Transfer/sec:      0.99MB
    3、metrics性能分析
        http://localhost:2338/metrics
# 版权
    MIT
