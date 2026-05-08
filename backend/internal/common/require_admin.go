package common

import (
	"ElainaBlog/internal/common/model"

	"github.com/gin-gonic/gin"
)

// RequireAdmin 公共管理员权限校验 helper。
// checkAdmin 接受用户 ID，返回是否为管理员；校验失败直接写响应并返回 false。
func RequireAdmin(c *gin.Context, checkAdmin func(int64) (bool, error)) bool {
	userID := c.GetInt64(CtxUserIDKey)
	isAdmin, err := checkAdmin(userID)
	if err != nil {
		appErr := model.ErrInternal.WithDetail(err.Error())
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return false
	}
	if !isAdmin {
		appErr := model.ErrForbidden.WithDetail("仅管理员可操作")
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return false
	}
	return true
}
