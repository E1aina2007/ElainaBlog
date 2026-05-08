package util

import (
	"math/rand"
	"strings"
)

// GenerateCode 生成指定长度的数字验证码
func GenerateCode(length int) string {
	var sb strings.Builder

	for i := 0; i < length; i++ {
		sb.WriteByte(byte('0' + rand.Intn(10)))
	}
	return sb.String()
}

