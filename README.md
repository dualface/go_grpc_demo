# go_grpc_demo 架构说明

使用 protobuf 3 定义消息和服务，由 gRPC 框架处理连接请求。


## 目录结构

```text
ROOT/
  +--- cmd/                         使用 github.com/spf13/cobra 提供命令行支持
  |  +--- root.go
  |  \--- serve.go                  实现 serve 命令
  |
  +--- config/                      保存所有配置文件
  |  \--- global_config.go
  |
  +--- protobuf/                    保存所有 .proto 定义文件
  |  +--- hello.proto               Hello 消息和服务的定义
  |  +--- auth.proto                Auth 消息和服务的定义
  |
  +--- server
  |  +--- rpc/                      protoc 编译生成的文件放入此目录
  |  |  +--- hello.pb.go
  |  |  +--- hello_grpc.pb.go
  |  |  +--- auth.pb.go
  |  |  \--- auth_grpc.pb.go
  |  |
  |  +--- impl/                     服务端架构的主要实现
  |  |  +--- start_server.go        启动服务器的 StartServer() 实现
  |  |  +--- hub_impl.go
  |  |  +--- client_impl.go
  |  |  \--- room_impl.go
  |  |
  |  +--- hub.go                    Hub 维护所有客户端连接，以及客户端连接的验证
  |  +--- client.go                 Client 封装客户端连接
  |  \--- room.go                   验证后的 Client 会加入 Room，方便对 Client 进行分组
  |
  +--- main.go                      启动应用入口
  +--- server.toml                  配置文件，由 global_config.go 访问
  \--- README.md
```

## 客户端连接到服务器

客户端连接到服务器后：
- 构造一个 Client struct
- 服务端开启一个计时器
- 如果客户端没有在超时之前通过 Auth 消息完成验证，客户端连接将被强制断开
  - Auth 消息中包含用于验证的 Token 字符串，具体验证可以随意搭配
  - demo 中仅仅检查 Token 是否在 Redis 中存在
- 验证通过后 Client 被注册到 Hub 中
- Hub 会根据配置，将 Client 分配到一个 Room 中
  - 默认分配规则是查找第一个符合负载要求的 Room
  - 如果所有 Room 都满员，新建一个 Room
  - 如果达到 Room 数量限制，客户端连接将被强制断开
- 客户端断开连接时，会自动从 Hub 和 Room 中注销


## Redis 在架构中的作用

如果是多进程部署，需要依赖 Redis 支持：
- 跨进程数据共享
- 分布式锁
- 分布式队列
