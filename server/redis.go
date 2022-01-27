package server

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go_grpc_demo/config"
	"log"
	"time"
)

type (
	RedisConnect struct {
		DB  *redis.Client
		CTX context.Context
	}
)

var redisConn *RedisConnect

func InitRedis(cfg *config.GlobalConfig, ctx context.Context) (*RedisConnect, error) {
	if redisConn != nil {
		return nil, fmt.Errorf("[Redis]: already initialized")
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:         cfg.Redis.Addr,
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		DialTimeout:  cfg.Redis.DialTimeout * time.Millisecond,
		ReadTimeout:  cfg.Redis.ReadTimeout * time.Millisecond,
		WriteTimeout: cfg.Redis.WriteTimeout * time.Millisecond,
	})
	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		return nil, fmt.Errorf("[Redis]: redisConn to redis failed, %s", err)
	}

	redisConn = &RedisConnect{DB: rdb, CTX: ctx}
	return redisConn, nil
}

func GetRedis() *RedisConnect {
	if redisConn == nil {
		log.Panicln("[Redis]: not initialized")
	}
	return redisConn
}
