# [X-PROJECT](https://github.com/xgd16/x-project/releases)

<img src="https://goframe.org/download/attachments/1114119/logo2.png?version=1&modificationDate=1684158720965&api=v2" width="300" alt="">

> 基于 GoFrame
～[💈官方文档地址](https://goframe.org/display/gf)

### 💿 支持

``mac`` - arm64 amd64

``windows`` - amd64

``linux`` - amd64

### 💼 编译

##### 1.直接编译
```shell
go build -o ./bin/xProject main.go
```
##### 2.使用 makefile
```shell
make help
```
```shell
make linux-amd64
```

### ⚙️使用扩展地址

1. [gf-x-tool](github.com/xgd16/gf-x-tool) 工具扩展
2. [gf-x-rabbitMQ](github.com/xgd16/gf-x-rabbitMQ) 队列支持扩展

### 🌲目录结构

```
➜  x-project git:(master) ✗ tree                  
.
├── PROJECT_README.md
├── bin // 编译后文件生成位置
├── config.yaml // 系统配置文件
├── go.mod // 扩展支持文件
├── go.sum
├── main.go // 入口文件
├── makefile // make 支持
└── src
    ├── global // 全局可访问变量资源
    │   └── system.go
    ├── lib // 编写需要的扩展位置
    │   └── helper.go
    ├── models // 模型
    ├── service // 服务
    │   ├── cmd
    │   ├── init.go // 在此 注册 service
    │   ├── queue // 基于 rabbitMQ 的队列服务
    │   │   ├── handler
    │   │   └── service.go
    │   └── web // HTTP 服务
    │       ├── controller
    │       ├── route // 编写 HTTP 路由
    │       │   └── api.go
    │       └── service.go
    └── types // 定义类型
```

