package rdb

import (
	"context"
)

const bannedIPSetKey = "security:banned_ips"

// BanIP 将 IP 加入封禁集合
func BanIP(client RedisClient, ip string) error {
	return client.SAdd(context.Background(), bannedIPSetKey, ip).Err()
}

// UnbanIP 从封禁集合中移除 IP
func UnbanIP(client RedisClient, ip string) error {
	return client.SRem(context.Background(), bannedIPSetKey, ip).Err()
}

// IsIPBanned 检查 IP 是否被封禁
func IsIPBanned(client RedisClient, ip string) bool {
	exists, err := client.SIsMember(context.Background(), bannedIPSetKey, ip).Result()
	if err != nil {
		return false
	}
	return exists
}

// GetBannedIPs 获取所有被封禁的 IP
func GetBannedIPs(client RedisClient) ([]string, error) {
	return client.SMembers(context.Background(), bannedIPSetKey).Result()
}
