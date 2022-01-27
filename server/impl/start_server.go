package impl

import (
	"context"
	"go_grpc_demo/config"
	"go_grpc_demo/server"
	"log"
)

func StartServer(cfgFile string, printCfg bool) {
	// 载入配置文件
	cfg, err := config.LoadGlobalConfig(cfgFile, printCfg)
	if err != nil {
		log.Fatalf("load config failed, %v", err)
	}

	// 初始化 Redis 连接
	ctx := context.Background()
	if cfg.Redis.Enabled {
		_, err = server.InitRedis(cfg, ctx)
		if err != nil {
			log.Fatalf("connect to redis failed, %v", err)
		}
	}

	// 初始化 MySQL 连接
	if cfg.MySQL.Enabled {
		_, err = server.InitMySQL(cfg, ctx)
		if err != nil {
			log.Fatalf("connect to mysql failed, %v", err)
		}
	}

	// Hub 启动 gRPC，并维护所有 Room 和 Client
	hub := NewHub(cfg)
	hub.Run()
}
