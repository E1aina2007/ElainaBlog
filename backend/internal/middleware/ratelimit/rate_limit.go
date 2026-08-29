// rate_limit.go 基于 Redis 的固定窗口速率限制中间件
package ratelimit

import (
	"ElainaBlog/internal/response"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// incrWithTTL 原子执行 INCR 与首次 EXPIRE：
// 两者分开调用时，若 INCR 成功而 EXPIRE 失败（Redis 抖动），
// key 将永不过期，对应 IP 会被永久限流
var incrWithTTL = redis.NewScript(`
local v = redis.call('INCR', KEYS[1])
if v == 1 then
	redis.call('EXPIRE', KEYS[1], ARGV[1])
end
return v
`)

// Limit 基于客户端 IP 的固定窗口速率限制中间件。
// keyPrefix: Redis 键前缀（用于区分不同接口）
// maxRequests: 时间窗口内最大请求数
// window: 时间窗口大小
func Limit(rdb *redis.Client, keyPrefix string, maxRequests int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := "rate_limit:" + keyPrefix + ":" + c.ClientIP()

		count, err := incrWithTTL.Run(c.Request.Context(), rdb, []string{key}, int64(window.Seconds())).Int64()
		if err != nil {
			// Redis 故障时放行，不影响正常请求
			c.Next()
			return
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
