// admin_auth.go 管理员权限中间件
package jwt

import (
	authsvc "ElainaBlog/internal/auth"
	"ElainaBlog/internal/response"
	"ElainaBlog/internal/user"

	"github.com/gin-gonic/gin"
)

// RequireAdmin 管理员权限中间件：校验当前用户是否为管理员，
// 非管理员返回 403 并终止请求。必须在 RequireAuth 之后使用。
func RequireAdmin(userService *user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetInt64(authsvc.CtxUserIDKey)
		isAdmin, err := userService.CheckIsAdmin(c.Request.Context(), userID)
		if err != nil {
			c.JSON(response.ErrInternal.HTTPStatus(), response.ApiErrorResponse(response.ErrInternal.Code, response.ErrInternal.Message, nil))
			c.Abort()
			return
		}
		if !isAdmin {
			appErr := response.ErrForbidden.WithDetail("仅管理员可操作")
			c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
			c.Abort()
			return
		}
		c.Next()
	}
}
