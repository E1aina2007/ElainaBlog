package redis

import (
	"context"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

// 封禁以独立键存储（security:banned:<ip>），从而支持按成员设置 TTL：
// 登录失败触发的自动封禁带过期时间到点自动解封，管理员手动封禁不过期
const bannedIPKeyPrefix = "security:banned:"

func bannedIPKey(ip string) string {
	return bannedIPKeyPrefix + ip
}

// BanIP 将 IP 加入封禁，ttl 为 0 表示永久封禁（由管理员手动解封）
func BanIP(client *redis.Client, ip string, ttl time.Duration) error {
	value := "manual"
	if ttl > 0 {
		value = "auto"
	}
	return client.Set(context.Background(), bannedIPKey(ip), value, ttl).Err()
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

// GetBannedIPs 获取所有被封禁的 IP
func GetBannedIPs(client *redis.Client) ([]string, error) {
	var ips []string
	var cursor uint64
	for {
		keys, next, err := client.Scan(context.Background(), cursor, bannedIPKeyPrefix+"*", 100).Result()
		if err != nil {
			return ips, err
		}
		for _, key := range keys {
			ips = append(ips, strings.TrimPrefix(key, bannedIPKeyPrefix))
		}
		cursor = next
		if cursor == 0 {
			break
		}
	}
	return ips, nil
}
