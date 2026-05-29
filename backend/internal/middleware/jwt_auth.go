// middleware 提供 Gin 中间件，包括 JWT 鉴权等。
package middleware

import (
	"ElainaBlog/internal/common"
	"ElainaBlog/internal/common/model"
	"strings"

	"github.com/gin-gonic/gin"
)

// JwtAuthMiddleware JWT 鉴权中间件，持有 JwtAuthService 实例用于 token 的签发与校验。
type JwtAuthMiddleware struct {
	JwtAuthService *common.JwtAuthService
}

// NewJwtAuthMiddleware 创建 JWT 鉴权中间件实例。
func NewJwtAuthMiddleware(jwtAuthService *common.JwtAuthService) *JwtAuthMiddleware {
	return &JwtAuthMiddleware{JwtAuthService: jwtAuthService}
}

// RequireAuth 强制鉴权：从 Authorization 头中提取 Bearer token，
// 校验通过后将 UserID 和 Claims 写入 gin.Context，校验失败则返回 401 并终止请求。
func (m *JwtAuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if m == nil || m.JwtAuthService == nil {
			c.JSON(model.ErrInternal.HTTPStatus(), model.ApiErrorResponse(model.ErrInternal.Code, model.ErrInternal.Message, nil))
			c.Abort()
			return
		}

		var tokenString string

		authHeader := strings.TrimSpace(c.GetHeader("Authorization"))
		if authHeader != "" {
			tokenString = extractBearerToken(authHeader)
		}

		if tokenString == "" {
			tokenString, _ = c.Cookie("access_token")
		}

		if tokenString == "" {
			c.JSON(model.ErrUnauthorized.HTTPStatus(), model.ApiErrorResponse(model.ErrUnauthorized.Code, model.ErrUnauthorized.Message, nil))
			c.Abort()
			return
		}

		claims, err := m.JwtAuthService.ParseAndVerifyAccessToken(tokenString)
		if err != nil {
			c.JSON(model.ErrTokenInvalid.HTTPStatus(), model.ApiErrorResponse(model.ErrTokenInvalid.Code, model.ErrTokenInvalid.Message, model.ErrTokenInvalid))
			c.Abort()
			return
		}

		c.Set(common.CtxUserIDKey, claims.UserID)
		c.Next()
	}
}

// extractBearerToken 从 "Bearer <token>" 格式的 Authorization 头中提取 token 字符串。
func extractBearerToken(authHeader string) string {
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 {
		return ""
	}
	if !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return strings.TrimSpace(parts[1])
}

