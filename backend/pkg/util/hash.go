package util

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// EmailToAvatarHash 根据邮箱生成头像文件名哈希
// 同一邮箱始终生成相同的哈希值，用于头像文件命名
func EmailToAvatarHash(email string) string {
	email = strings.ToLower(strings.TrimSpace(email))
	hash := sha256.Sum256([]byte(email))
	return hex.EncodeToString(hash[:8]) // 取前8字节，生成16位十六进制字符串
}
