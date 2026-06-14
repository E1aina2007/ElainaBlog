// upload_limit.go 上传频率限制中间件，基于用户 ID 进行限制
package middleware

import (
	"ElainaBlog/internal/common"
	"ElainaBlog/internal/common/model"
	"ElainaBlog/pkg/rdb"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// UploadLimitMiddleware 基于用户 ID 的上传频率限制中间件。
type UploadLimitMiddleware struct {
	rdb rdb.RedisClient
}

// NewUploadLimitMiddleware 创建上传频率限制中间件实例。
func NewUploadLimitMiddleware(redis rdb.RedisClient) *UploadLimitMiddleware {
	return &UploadLimitMiddleware{rdb: redis}
}

// Limit 返回一个 Gin 中间件，基于登录用户 ID 进行上传频率限制。
// maxUploads: 时间窗口内最大上传次数
// window: 时间窗口大小
func (m *UploadLimitMiddleware) Limit(maxUploads int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 JWT 上下文获取用户 ID
		userID, exists := c.Get(common.CtxUserIDKey)
		if !exists {
			// 未登录用户不进行限制（由 auth 中间件处理）
			c.Next()
			return
		}

		key := fmt.Sprintf("upload:limit:%d:%d", userID.(int64), time.Now().Unix()/int64(window.Seconds()))

		ctx := context.Background()
		count, err := m.rdb.Incr(ctx, key).Result()
		if err != nil {
			// Redis 故障时放行，不影响正常上传
			c.Next()
			return
		}

		// 首次请求时设置过期时间
		if count == 1 {
			m.rdb.Expire(ctx, key, window)
		}

		if count > int64(maxUploads) {
			c.JSON(http.StatusTooManyRequests, model.ApiErrorResponse(
				model.ErrTooManyRequests.Code,
				fmt.Sprintf("上传过于频繁，每 %v 最多上传 %d 张图片", window, maxUploads),
				model.ErrTooManyRequests,
			))
			c.Abort()
			return
		}

		c.Next()
	}
}
