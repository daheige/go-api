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
# 关于gin版本
    如果go build编译的时候提示gin/json找不到，可以将gin版本改成1.3.0或等待官方升级到更高的版本。
# 版权
    MIT
