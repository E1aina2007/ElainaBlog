package notification

import (
	"ElainaBlog/internal/auth"
	"ElainaBlog/internal/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service *Service
}

func NewController(service *Service) *Controller {
	return &Controller{service: service}
}

// GetList 获取通知列表
func (ctl *Controller) GetList(c *gin.Context) {
	userID := c.GetInt64(auth.CtxUserIDKey)
	onlyUnread := c.Query("unread") == "1"

	list, err := ctl.service.GetList(c.Request.Context(), userID, onlyUnread)
	if err != nil {
		c.JSON(response.ErrInternal.HTTPStatus(), response.ApiErrorResponse(response.ErrInternal.Code, response.ErrInternal.Message, nil))
		return
	}
	c.JSON(http.StatusOK, response.ApiSuccessResponse(list))
}

// GetUnreadCount 获取未读通知数量
func (ctl *Controller) GetUnreadCount(c *gin.Context) {
	userID := c.GetInt64(auth.CtxUserIDKey)

	count, err := ctl.service.GetUnreadCount(c.Request.Context(), userID)
	if err != nil {
		c.JSON(response.ErrInternal.HTTPStatus(), response.ApiErrorResponse(response.ErrInternal.Code, response.ErrInternal.Message, nil))
		return
	}
	c.JSON(http.StatusOK, response.ApiSuccessResponse(gin.H{"count": count}))
}

// MarkAsRead 标记单条通知为已读
func (ctl *Controller) MarkAsRead(c *gin.Context) {
	var req MarkReadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := response.ErrInvalidParams.WithDetail("请求参数格式错误")
		c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	userID := c.GetInt64(auth.CtxUserIDKey)
	if err := ctl.service.MarkAsRead(c.Request.Context(), req.ID, userID); err != nil {
		switch err {
		case ErrInvalidParams:
			appErr := response.ErrInvalidParams.WithDetail(err.Error())
			c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		default:
			c.JSON(response.ErrInternal.HTTPStatus(), response.ApiErrorResponse(response.ErrInternal.Code, response.ErrInternal.Message, nil))
		}
		return
	}
	c.JSON(http.StatusOK, response.ApiSuccessResponse(nil))
}

// MarkAllAsRead 标记所有通知为已读
func (ctl *Controller) MarkAllAsRead(c *gin.Context) {
	userID := c.GetInt64(auth.CtxUserIDKey)
	if err := ctl.service.MarkAllAsRead(c.Request.Context(), userID); err != nil {
		c.JSON(response.ErrInternal.HTTPStatus(), response.ApiErrorResponse(response.ErrInternal.Code, response.ErrInternal.Message, nil))
		return
	}
	c.JSON(http.StatusOK, response.ApiSuccessResponse(nil))
}

// Delete 删除通知
func (ctl *Controller) Delete(c *gin.Context) {
	var req DeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := response.ErrInvalidParams.WithDetail("请求参数格式错误")
		c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	userID := c.GetInt64(auth.CtxUserIDKey)
	if err := ctl.service.Delete(c.Request.Context(), req.ID, userID); err != nil {
		switch err {
		case ErrInvalidParams:
			appErr := response.ErrInvalidParams.WithDetail(err.Error())
			c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		default:
			c.JSON(response.ErrInternal.HTTPStatus(), response.ApiErrorResponse(response.ErrInternal.Code, response.ErrInternal.Message, nil))
		}
		return
	}
	c.JSON(http.StatusOK, response.ApiSuccessResponse(nil))
}
