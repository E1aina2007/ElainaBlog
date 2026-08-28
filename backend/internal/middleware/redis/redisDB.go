package redis

import (
	"ElainaBlog/internal/config"
	"context"

	"github.com/redis/go-redis/v9"
)

var DefaultClient *redis.Client

func InitRedis(cfg *config.RedisConfig) error {
	DefaultClient = redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})
	_, err := DefaultClient.Ping(context.Background()).Result()
	return err
}
