// admin_auth.go 管理员权限中间件
package middleware

import (
	"ElainaBlog/internal/common"
	"ElainaBlog/internal/common/model"
	"ElainaBlog/internal/user"

	"github.com/gin-gonic/gin"
)

// AdminAuthMiddleware 管理员权限中间件，持有 userService 实例用于校验管理员身份。
type AdminAuthMiddleware struct {
	userService *user.Service
}

// NewAdminAuthMiddleware 创建管理员权限中间件实例。
func NewAdminAuthMiddleware(userService *user.Service) *AdminAuthMiddleware {
	return &AdminAuthMiddleware{userService: userService}
}

// RequireAdmin 强制管理员权限：校验当前用户是否为管理员，
// 非管理员返回 403 并终止请求。必须在 RequireAuth 之后使用。
func (m *AdminAuthMiddleware) RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetInt64(common.CtxUserIDKey)
		isAdmin, err := m.userService.CheckIsAdmin(userID)
		if err != nil {
			c.JSON(model.ErrInternal.HTTPStatus(), model.ApiErrorResponse(model.ErrInternal.Code, model.ErrInternal.Message, nil))
			c.Abort()
			return
		}
		if !isAdmin {
			appErr := model.ErrForbidden.WithDetail("仅管理员可操作")
			c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
			c.Abort()
			return
		}
		c.Next()
	}
}
