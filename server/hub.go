package server

import (
	"context"
	"go_grpc_demo/config"
)

type (
	// Hub 提供全局配置、房间查找等功能
	Hub interface {
		// Run 启动
		Run()

		// GlobalConfig 返回全局设置
		GlobalConfig() *config.GlobalConfig

		// NumberOfClients 返回在线客户端总数
		NumberOfClients() int

		// GetClientByContext 从 Context 中提取 ClientId，然后再查找 Client
		GetClientByContext(ctx context.Context) (Client, bool)

		// GetAuthedClientByContext 从 Context 中提取 ClientId，然后再查找 Client，并且确保客户端已经通过验证
		GetAuthedClientByContext(ctx context.Context) (Client, bool)
	}
)
