package rdb

import (
	"context"
)

const bannedIPSetKey = "security:banned_ips"

// BanIP 将 IP 加入封禁集合
func BanIP(ip string) error {
	return RedisClient.SAdd(context.Background(), bannedIPSetKey, ip).Err()
}

// UnbanIP 从封禁集合中移除 IP
func UnbanIP(ip string) error {
	return RedisClient.SRem(context.Background(), bannedIPSetKey, ip).Err()
}

// IsIPBanned 检查 IP 是否被封禁
func IsIPBanned(ip string) bool {
	exists, err := RedisClient.SIsMember(context.Background(), bannedIPSetKey, ip).Result()
	if err != nil {
		return false
	}
	return exists
}

// GetBannedIPs 获取所有被封禁的 IP
func GetBannedIPs() ([]string, error) {
	return RedisClient.SMembers(context.Background(), bannedIPSetKey).Result()
}
