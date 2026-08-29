// Package jwt 提供 JWT 鉴权与管理员权限校验中间件（执行层）。
package jwt

import (
	authsvc "ElainaBlog/internal/auth"
	"ElainaBlog/internal/response"
	"strings"

	"github.com/gin-gonic/gin"
)

// RequireAuth JWT 鉴权中间件：从 Authorization 头或 Cookie 中提取 token，
// 校验通过后将 UserID 写入 gin.Context，校验失败则返回 401 并终止请求。
func RequireAuth(tokenMgr authsvc.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		if tokenMgr == nil {
			c.JSON(response.ErrInternal.HTTPStatus(), response.ApiErrorResponse(response.ErrInternal.Code, response.ErrInternal.Message, nil))
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
			c.JSON(response.ErrUnauthorized.HTTPStatus(), response.ApiErrorResponse(response.ErrUnauthorized.Code, response.ErrUnauthorized.Message, nil))
			c.Abort()
			return
		}

		claims, err := tokenMgr.ParseAndVerifyAccessToken(tokenString)
		if err != nil {
			c.JSON(response.ErrTokenInvalid.HTTPStatus(), response.ApiErrorResponse(response.ErrTokenInvalid.Code, response.ErrTokenInvalid.Message, nil))
			c.Abort()
			return
		}

		c.Set(authsvc.CtxUserIDKey, claims.UserID)
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
