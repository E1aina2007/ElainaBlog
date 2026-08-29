package redis

import (
	"context"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

// 封禁以独立键存储（security:banned:<ip>），从而支持按成员设置 TTL：
// 登录失败触发的自动封禁带过期时间到点自动解封，管理员手动封禁不过期。
// 键值存储封禁时间（Unix 秒），供管理端展示
const bannedIPKeyPrefix = "security:banned:"

// BannedIP 封禁条目
type BannedIP struct {
	IP         string `json:"ip"`
	BannedAt   int64  `json:"banned_at"`   // 封禁时间（Unix 秒）
	TTLSeconds int64  `json:"ttl_seconds"` // 剩余秒数；-1 表示永久封禁
}

func bannedIPKey(ip string) string {
	return bannedIPKeyPrefix + ip
}

// BanIP 将 IP 加入封禁并记录封禁时间，ttl 为 0 表示永久封禁（由管理员手动解封）
func BanIP(client *redis.Client, ip string, ttl time.Duration) error {
	return client.Set(context.Background(), bannedIPKey(ip), strconv.FormatInt(time.Now().Unix(), 10), ttl).Err()
}

// UnbanIP 解封IP
func UnbanIP(client *redis.Client, ip string) error {
	return client.Del(context.Background(), bannedIPKey(ip)).Err()
}

// IsIPBanned 检查 IP 是否被封禁
func IsIPBanned(client *redis.Client, ip string) bool {
	n, err := client.Exists(context.Background(), bannedIPKey(ip)).Result()
	if err != nil {
		// Redis 故障时放行，与限流的 fail-open 策略一致
		return false
	}
	return n > 0
}

// GetBannedIPs 获取所有被封禁的 IP 及其封禁时间与剩余时长
func GetBannedIPs(client *redis.Client) ([]BannedIP, error) {
	var keys []string
	var cursor uint64
	for {
		batch, next, err := client.Scan(context.Background(), cursor, bannedIPKeyPrefix+"*", 100).Result()
		if err != nil {
			return nil, err
		}
		keys = append(keys, batch...)
		cursor = next
		if cursor == 0 {
			break
		}
	}
	if len(keys) == 0 {
		return []BannedIP{}, nil
	}

	// 批量取键值（封禁时间）与 PTTL（剩余时长）
	ctx := context.Background()
	pipe := client.Pipeline()
	getCmds := make([]*redis.StringCmd, len(keys))
	ttlCmds := make([]*redis.DurationCmd, len(keys))
	for i, key := range keys {
		getCmds[i] = pipe.Get(ctx, key)
		ttlCmds[i] = pipe.PTTL(ctx, key)
	}
	if _, err := pipe.Exec(ctx); err != nil && err != redis.Nil {
		return nil, err
	}

	ips := make([]BannedIP, 0, len(keys))
	for i, key := range keys {
		if ttlCmds[i].Val() == -2*time.Second {
			continue // 键在 SCAN 与读取间隙过期，跳过
		}
		bannedAt, _ := strconv.ParseInt(getCmds[i].Val(), 10, 64)
		ttlSeconds := int64(-1) // PTTL 为 -1 即永不过期
		if d := ttlCmds[i].Val(); d > 0 {
			ttlSeconds = int64(math.Ceil(d.Seconds()))
		}
		ips = append(ips, BannedIP{
			IP:         strings.TrimPrefix(key, bannedIPKeyPrefix),
			BannedAt:   bannedAt,
			TTLSeconds: ttlSeconds,
		})
	}
	return ips, nil
}

