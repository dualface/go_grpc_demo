# go_grpc_demo 架构说明

使用 protobuf 3 定义消息和服务，由 gRPC 框架处理连接请求。

## 目录结构

```text
ROOT/
|   compile-proto-files.cmd     编译 .proto
|   main.go
|   README.md
|   server.toml                 服务器配置文件
|
+---cmd                         使用 github.com/spf13/cobra 提供命令行支持
|       client.go               测试客户端
|       root.go
|       serve.go                启动服务器
|
+---config
|       global_config.go        全局配置
|
+---protobuf                    所有 .proto 定义文件
|       auth.proto              Auth 服务的定义
|       demo.proto              Demo 服务的定义
|
\---server                      服务端代码
    |   client.go               interface Client
    |   hub.go                  interface Hub
    |   mysql.go                简单的 MySQL 连接封装
    |   redis.go                简单的 Redis 连接封装
    |
    +---impl
    |       client_impl.go
    |       grpc_stats_handler.go           处理 gPRC 连接状态
    |       grpc_transport_credentials.go   定制的 gRPC 连接认证
    |       hub_impl.go
    |       start_server.go
    |
    \---rpc
        |   auth.pb.go
        |   auth_grpc.pb.go
        |   demo.pb.go
        |   demo_grpc.pb.go
        |
        \---impl
                auth_impl.go    Auth 服务实现
                demo_impl.go    Demo 服务实现
```

## 客户端连接到服务器

客户端连接到服务器后：

- 构造一个 Client struct
- 服务端开启一个计时器
- 如果客户端没有在超时之前通过 Auth 消息完成验证，客户端连接将被强制断开
    - Auth 消息中包含用于验证的 Token 字符串，具体验证可以随意搭配

## Auth 机制存在的问题

demo 使用了一种简单的验证机制，既客户端连接后，需要在超时之前调用 Auth 方法，并传入 Token。

这种方式存在的问题：

- 超时之前连接会被保持，遭受攻击时，会占用大量服务器连接
- 因为超时或 Auth 失败，服务端强制 Close 连接后，客户端可以通过断线重连继续重新尝试

如何改进：

- 应该使用 SSL/TLS 等更安全的验证机制，但对环境配置和客户端要求更为麻烦

\- EOF \-
