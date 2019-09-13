# gin 框架实战

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

# golang 环境安装

    golang下载地址:
        https://golang.google.cn/dl/

    以go最新版本go1.13版本为例
    https://dl.google.com/go/go1.13.linux-amd64.tar.gz
    1、linux环境，下载go1.12.8.linux-amd64.tar.gz
        cd /usr/local/
        sudo wget https://dl.google.com/go/go1.13.linux-amd64.tar.gz
        sudo tar zxvf go1.13.linux-amd64.tar.gz
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

# 设置 goproxy 代理

    go version >= 1.13
    设置goproxy代理
    vim ~/.bashrc添加如下内容:
    export GOPROXY=https://goproxy.io,direct
    或者
    export GOPROXY=https://goproxy.cn,direct

    让bashrc生效
    source ~/.bashrc

    go version < 1.13
    vim ~/.bashrc添加如下内容：
    export GOPROXY=https://goproxy.io
    或者使用 export GOPROXY=https://athens.azurefd.net
    或者使用 export GOPROXY=https://mirrors.aliyun.com/goproxy/
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

# 关于 redisgo 调优

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

# docker 运行

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

        web图形化查看
            1. $ sudo apt install graphviz
            2. go tool pprof profile /home/heige/pprof/pprof.go-api.samples.cpu.002.pb.gz
            3. (pprof) web

        prometheus性能监控
        http://localhost:2338/metrics

# wrk 工具压力测试

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

        压力测试/api/info接口
        $ wrk -t 8 -c 100 -d 1m --latency http://localhost:1338/api/info
        Running 1m test @ http://localhost:1338/api/info
        8 threads and 100 connections
        Thread Stats Avg Stdev Max +/- Stdev
        Latency 21.69ms 48.75ms 604.71ms 97.39%
        Req/Sec 833.19 149.83 1.76k 78.02%
        Latency Distribution
        50% 15.34ms
        75% 18.86ms
        90% 29.00ms
        99% 317.16ms
        391027 requests in 1.00m, 69.73MB read
        Requests/sec: 6507.18
        Transfer/sec: 1.16MB
        平均每个请求 15-30ms 处理完毕

    3、metrics性能分析
        http://localhost:2338/metrics

# 关于 http 超时的限制

    不恰当的http.Server设置，以及未设置超时处理，可能导致http net.Conn连接泄漏，从而出现太多的文件句柄
    最为直接的原因，就导致服务异常，无法正常响应，出现too many open files的问题，解决方案参考main.go
    压力测试：
    $ wrk -t 8 -c 400 -d 20s http://localhost:1338/indexRunning 20s test @ http://localhost:1338/index
      8 threads and 400 connections
      Thread Stats   Avg      Stdev     Max   +/- Stdev
        Latency    50.61ms   31.75ms 283.06ms   67.54%
        Req/Sec     0.99k   263.19     3.06k    85.16%
      156615 requests in 20.05s, 22.40MB read
    Requests/sec:   7809.62
    Transfer/sec:      1.12MB

# db 压力测试

    $ cd mytest
    $ wrk -t 8 -d 5m -c 400 http://localhost:1338/v1/data
    Running 5m test @ http://localhost:1338/v1/data
      8 threads and 400 connections


    查看文件句柄fd情况
    $ ps -ef | grep "go run"
    heige    14009 13055  0 13:36 pts/8    00:00:00 go run main.go

    $ lsof -p 14009 | wc -l
    12
    $ lsof -p 14009
    COMMAND   PID  USER   FD      TYPE DEVICE  SIZE/OFF     NODE NAME
    go      14009 heige  cwd       DIR    8,1      4096  2490371 /web/go/go-api
    go      14009 heige  rtd       DIR    8,1      4096        2 /
    go      14009 heige  txt       REG    8,1  14613596 22809199 /usr/local/go/bin/go
    go      14009 heige  mem       REG    8,1   2030544 24252912 /lib/x86_64-linux-gnu/libc-2.27.so
    go      14009 heige  mem       REG    8,1    144976 24253540 /lib/x86_64-linux-gnu/libpthread-2.27.so
    go      14009 heige  mem       REG    8,1    170960 24252893 /lib/x86_64-linux-gnu/ld-2.27.so
    go      14009 heige    0u      CHR  136,8       0t0       11 /dev/pts/8
    go      14009 heige    1u      CHR  136,8       0t0       11 /dev/pts/8
    go      14009 heige    2u      CHR  136,8       0t0       11 /dev/pts/8
    go      14009 heige    3w      REG    8,1 139838601  7209351 /home/heige/.cache/go-build/log.txt
    go      14009 heige    4u  a_inode   0,13         0    10638 [eventpoll]

    压力测试过程中，查看mysql
    $ lsof -i TCP | grep mysql | wc -l
    42

    $ lsof -i :3306 | wc -l
    261
    $ lsof -i :3306 | wc -l
    60

    正在建立连接通信的mysql
    $ lsof -i :3306 | grep ESTABLISHED | wc -l
    60

    查看mysql建立tcp个数
    $ lsof  -i -sTCP:ESTABLISHED | grep mysql | wc -l
    107

    查看1338建立连接的个数
    $ lsof -i :1338 | wc -l
    802

    查看进程里的mysql连接情况
    $ lsof -p 14009 -i | grep mysql | wc -l
    71

    比较配置文件中的空闲连接和mysql实际连接的个数,基本上一样
    $ lsof -p 14009 -i | grep mysql | wc -l
    60

    当大规模的请求过来的时候，fd数量里面上涨
    $ lsof -p 14009 -i | wc -l
    886

    当请求下来后，查看fd情况
    $ lsof -p 14009 -i | wc -l
    89

    当请求结束后，查看mysql连接情况
    $ lsof -p 14009 -i | grep mysql | wc -l
      60
    $ netstat -an | grep TIME_WAIT | grep 3306 | wc -l
    0

    $ netstat -ae|grep mysql | wc -l
    122

    $ netstat -an | grep TIME_WAIT | grep 3306
    $ netstat -an|awk '/tcp/ {print $6}'|sort|uniq -c
        134 ESTABLISHED
          1 FIN_WAIT1
         25 LISTEN
          3 SYN_SENT
          1 TIME_WAIT
    $ netstat -n | awk '/^tcp/ {++S[$NF]} END {for(a in S) print a, S[a]}'
        TIME_WAIT 1
        ESTABLISHED 146
        LAST_ACK 1
        SYN_SENT 2

    查看TIME_WAIT数量，$ netstat -ant| grep -i time_wait
    $ netstat -an | grep -c TIME_WAIT
    2

    $ ls -l /proc/14009/fd | wc -l
    6

    查看进程的fd情况
    $ lsof -p 14009 | wc -l
    12
    端口连接情况
    $ lsof -p 14009 -i :1338 | wc -l
    13

    经过压力测试表明gorm mysql连接池方式，当请求过大时候，超过空闲的连接数，就会新建连接句柄放入连接池中
    当请求下来后，mysql tcp都会降下来，golang进程的fd句柄也降下来了。

    压力测试结果：
    $ wrk -t 8 -d 5m -c 400 http://localhost:1338/v1/data
    Running 5m test @ http://localhost:1338/v1/data
     8 threads and 400 connections
     Thread Stats   Avg      Stdev     Max   +/- Stdev
       Latency   231.31ms  129.26ms   1.64s    76.03%
       Req/Sec   227.90    110.72   780.00     64.81%
     535769 requests in 5.00m, 87.88MB read
    Requests/sec:   1785.47
    Transfer/sec:    299.90KB

# 查看机器的 cpu，核数

    CPU总核数 = 物理CPU个数 * 每颗物理CPU的核数
    总逻辑CPU数 = 物理CPU个数 * 每颗物理CPU的核数 * 超线程数

    复制代码
    查看CPU信息（型号）
    # cat /proc/cpuinfo | grep name | cut -f2 -d: | uniq -c
    4  Intel(R) Core(TM) i5-2450M CPU @ 2.50GHz

    # 查看物理CPU个数
    # cat /proc/cpuinfo| grep "physical id"| sort| uniq| wc -l
    1

    # 查看每个物理CPU中core的个数(即核数)
    # cat /proc/cpuinfo| grep "cpu cores"| uniq
    cpu cores	: 2

    # 查看逻辑CPU的个数
    # cat /proc/cpuinfo| grep "processor"| wc -l
    4

    $ top -H -p 14009

    top - 14:18:48 up 1 day, 16:36,  1 user,  load average: 12.52, 8.71, 7.68
    Threads:  10 total,   0 running,  10 sleeping,   0 stopped,   0 zombie
    %Cpu(s): 71.2 us, 20.4 sy,  0.0 ni,  4.0 id,  0.0 wa,  0.0 hi,  4.4 si,  0.0 st
    KiB Mem :  8110128 total,   149228 free,  5636640 used,  2324260 buff/cache
    KiB Swap:   998396 total,   848400 free,   149996 used.  1618124 avail Mem

# 采用 profile 库查看 pprof 性能指标

    import "github.com/pkg/profile"

    在函数里面
    defer profile.Start().Stop()

    参考mytest/app.go，其他性能指标可以看profile源码
    $ go tool pprof -http=:8080 /tmp/profile235146184/cpu.pprof
    [11667:11684:0824/203331.299458:ERROR:browser_process_sub_thread.cc(221)] Waited 3 ms for network service
    open /tmp/go-build321889850/b001/exe/app: no such file or directory

    自动打开浏览器访问
    http://localhost:8080/ui/top

    火焰图： http://localhost:8080/ui/flamegraph

    测试db性能
    $ wrk -t 8 -d 5m -c 400 http://localhost:1338/v1/data
    Running 5m test @ http://localhost:1338/v1/data
      8 threads and 400 connections
      Thread Stats   Avg      Stdev     Max   +/- Stdev
        Latency   164.29ms   96.87ms   1.75s    81.10%
        Req/Sec   322.53    129.82   800.00     66.26%
      762208 requests in 5.00m, 125.03MB read
    Requests/sec:   2540.31
    Transfer/sec:    426.69KB

    请求结束后，退出app.go 会生成cpu.pprof
    2019/08/24 20:57:51 user:  &{2 hello}
    ^C2019/08/24 20:57:58 profile: caught interrupt, stopping profiles
    2019/08/24 20:57:58 exit signal:  interrupt
    2019/08/24 20:57:58 http: Server closed
    2019/08/24 20:57:58 profile: cpu profiling disabled, /tmp/profile682666456/cpu.pprof

    用pprof工具查看
    $ go tool pprof -http=:8080 /tmp/profile682666456/cpu.pprof
    [12616:12634:0824/210007.864297:ERROR:browser_process_sub_thread.cc(221)] Waited 1043 ms for network service
    http://localhost:8080/ui/

# 版权

    MIT
