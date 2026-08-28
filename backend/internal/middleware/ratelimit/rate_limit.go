// rate_limit.go 基于 Redis 的固定窗口速率限制中间件
package ratelimit

import (
	"ElainaBlog/internal/response"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// Limit 基于客户端 IP 的固定窗口速率限制中间件。
// keyPrefix: Redis 键前缀（用于区分不同接口）
// maxRequests: 时间窗口内最大请求数
// window: 时间窗口大小
func Limit(redis *redis.Client, keyPrefix string, maxRequests int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		key := fmt.Sprintf("rate_limit:%s:%s", keyPrefix, clientIP)

		ctx := context.Background()
		count, err := redis.Incr(ctx, key).Result()
		if err != nil {
			// Redis 故障时放行，不影响正常请求
			c.Next()
			return
		}

		// 首次请求时设置过期时间
		if count == 1 {
			redis.Expire(ctx, key, window)
		}

		if count > int64(maxRequests) {
			c.JSON(http.StatusTooManyRequests, response.ApiErrorResponse(
				response.ErrTooManyRequests.Code,
				"请求过于频繁，请稍后再试",
				response.ErrTooManyRequests,
			))
			c.Abort()
			return
		}

		c.Next()
	}
}
