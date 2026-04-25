// middleware 提供 Gin 中间件，包括 JWT 鉴权等。
package middleware

import (
	"ElainaBlog/internal/common"
	"ElainaBlog/internal/common/model"
	"net/http"
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
			appErr := model.ErrInternal.WithDetail("jwt service not initialized")
			c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
			c.Abort()
			return
		}

		authHeader := strings.TrimSpace(c.GetHeader("Authorization"))
		if authHeader == "" {
			appErr := model.ErrUnauthorized.WithDetail("missing Authorization header")
			c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
			c.Abort()
			return
		}

		tokenString := extractBearerToken(authHeader)
		if tokenString == "" {
			appErr := model.ErrUnauthorized.WithDetail("invalid Authorization header format")
			c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
			c.Abort()
			return
		}

		claims, err := m.JwtAuthService.ParseAndVerifyAccessToken(tokenString)
		if err != nil {
			appErr := model.ErrUnauthorized.WithDetail(err.Error())
			c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
			c.Abort()
			return
		}

		c.Set(common.CtxUserIDKey, claims.UserID)
		c.Set(common.CtxClaimsKey, claims)
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

// OptionalAuth 可选鉴权：有 token 就解析，无 token 则放行。
func (m *JwtAuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := strings.TrimSpace(c.GetHeader("Authorization"))
		if authHeader == "" {
			c.Next()
			return
		}

		tokenString := extractBearerToken(authHeader)
		if tokenString == "" || m == nil || m.JwtAuthService == nil {
			c.Next()
			return
		}

		claims, err := m.JwtAuthService.ParseAndVerifyAccessToken(tokenString)
		if err != nil {
			c.Next()
			return
		}

		c.Set(common.CtxUserIDKey, claims.UserID)
		c.Set(common.CtxClaimsKey, claims)
		c.Next()
	}
}

// WriteUnauthorized 是给非中间件场景的快捷返回。
func WriteUnauthorized(c *gin.Context, detail any) {
	appErr := model.ErrUnauthorized.WithDetail(detail)
	c.JSON(http.StatusUnauthorized, model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
}
