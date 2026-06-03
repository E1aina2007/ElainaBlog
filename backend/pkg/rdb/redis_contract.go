// redis_contract.go 定义 Redis 客户端接口，用于依赖注入和测试 mock
package rdb

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisClient 封装当前项目用到的所有 Redis 操作。
// *redis.Client 天然满足此接口。
type RedisClient interface {
	Incr(ctx context.Context, key string) *redis.IntCmd
	Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	Exists(ctx context.Context, keys ...string) *redis.IntCmd
	GetDel(ctx context.Context, key string) *redis.StringCmd
	Scan(ctx context.Context, cursor uint64, match string, count int64) *redis.ScanCmd
	Pipeline() redis.Pipeliner
	SAdd(ctx context.Context, key string, members ...interface{}) *redis.IntCmd
	SRem(ctx context.Context, key string, members ...interface{}) *redis.IntCmd
	SIsMember(ctx context.Context, key string, member interface{}) *redis.BoolCmd
	SCard(ctx context.Context, key string) *redis.IntCmd
	SMembers(ctx context.Context, key string) *redis.StringSliceCmd
	Ping(ctx context.Context) *redis.StatusCmd
	Info(ctx context.Context, sections ...string) *redis.StringCmd
}
