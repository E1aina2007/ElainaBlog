// admin_auth.go 管理员权限中间件
package middleware

import (
	"ElainaBlog/internal/common"
	"ElainaBlog/internal/common/model"

	"github.com/gin-gonic/gin"
)

// AdminChecker 管理员权限检查接口。
type AdminChecker interface {
	CheckIsAdmin(userID int64) (bool, error)
}

// AdminAuthMiddleware 管理员权限中间件，持有 AdminChecker 实例用于校验管理员身份。
type AdminAuthMiddleware struct {
	adminChecker AdminChecker
}

// NewAdminAuthMiddleware 创建管理员权限中间件实例。
func NewAdminAuthMiddleware(adminChecker AdminChecker) *AdminAuthMiddleware {
	return &AdminAuthMiddleware{adminChecker: adminChecker}
}

// RequireAdmin 强制管理员权限：校验当前用户是否为管理员，
// 非管理员返回 403 并终止请求。必须在 RequireAuth 之后使用。
func (m *AdminAuthMiddleware) RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetInt64(common.CtxUserIDKey)
		isAdmin, err := m.adminChecker.CheckIsAdmin(userID)
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
