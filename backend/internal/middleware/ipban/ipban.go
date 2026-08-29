// Package ipban 提供基于 Redis 的 IP 封禁检查中间件
package ipban

import (
	banredis "ElainaBlog/internal/middleware/redis"
	"ElainaBlog/internal/response"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// Middleware 拦截被封禁的 IP：命中返回 403 并终止请求。
// 作用于全部 /api 路由，client 为 nil 或 Redis 故障时放行（fail-open）
func Middleware(client *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		if client == nil {
			c.Next()
			return
		}
		if banredis.IsIPBanned(client, c.ClientIP()) {
			appErr := response.ErrForbidden.WithDetail("IP 已被封禁，如有疑问请联系管理员")
			c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
			c.Abort()
			return
		}
		c.Next()
	}
}
