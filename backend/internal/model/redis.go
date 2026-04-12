package model

import (
	"ElainaWeb/global"
	"context"
	"os"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func ConnectRedis() *redis.Client {
	redisCfg := global.Config.Redis

	client := redis.NewClient(&redis.Options{
		Addr:     redisCfg.Addr,
		Password: redisCfg.Password,
		DB:       redisCfg.DB,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		global.Logger.Error("Redis连接失败", zap.Error(err))
		os.Exit(1)
	}
	return client
}
