package rdb

import (
	"context"
	"fmt"
	"time"
)

// codeKey 验证码存储 key，格式: verify:code:<email>
func codeKey(email string) string {
	return fmt.Sprintf("verify:code:%s", email)
}

// limitKey 防刷限制 key，格式: verify:limit:<email>
func limitKey(email string) string {
	return fmt.Sprintf("verify:limit:%s", email)
}

// SetVerificationCode 存储验证码并设置防刷 key
// expiry: 验证码过期时间, interval: 重发间隔
func SetVerificationCode(client RedisClient, email, code string, expiry, interval time.Duration) error {
	ctx := context.Background()
	pipe := client.Pipeline()
	pipe.Set(ctx, codeKey(email), code, expiry)
	pipe.Set(ctx, limitKey(email), "1", interval)

	_, err := pipe.Exec(ctx)
	return err
}

// GetVerificationCode 根据邮箱获取验证码
func GetVerificationCode(client RedisClient, email string) (string, error) {
	return client.Get(context.Background(), codeKey(email)).Result()
}

// DeleteVerificationCode 验证成功后删除验证码
func DeleteVerificationCode(client RedisClient, email string) error {
	return client.Del(context.Background(), codeKey(email)).Err()
}

// IsDuringInterval 检查是否在重发间隔内（防刷）
func IsDuringInterval(client RedisClient, email string) (bool, error) {
	n, err := client.Exists(context.Background(), limitKey(email)).Result()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}
