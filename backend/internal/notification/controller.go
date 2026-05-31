package notification

import (
	"ElainaBlog/internal/common"
	"ElainaBlog/internal/common/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service *Service
}

func NewController(service *Service) *Controller {
	return &Controller{service: service}
}

type MarkReadRequest struct {
	ID int64 `json:"id"`
}

type DeleteRequest struct {
	ID int64 `json:"id"`
}

// GetList 获取通知列表
func (ctl *Controller) GetList(c *gin.Context) {
	userID := c.GetInt64(common.CtxUserIDKey)
	onlyUnread := c.Query("unread") == "1"

	list, err := ctl.service.GetList(userID, onlyUnread)
	if err != nil {
		c.JSON(model.ErrInternal.HTTPStatus(), model.ApiErrorResponse(model.ErrInternal.Code, model.ErrInternal.Message, nil))
		return
	}
	c.JSON(http.StatusOK, model.ApiSuccessResponse(list))
}

// GetUnreadCount 获取未读通知数量
func (ctl *Controller) GetUnreadCount(c *gin.Context) {
	userID := c.GetInt64(common.CtxUserIDKey)

	count, err := ctl.service.GetUnreadCount(userID)
	if err != nil {
		c.JSON(model.ErrInternal.HTTPStatus(), model.ApiErrorResponse(model.ErrInternal.Code, model.ErrInternal.Message, nil))
		return
	}
	c.JSON(http.StatusOK, model.ApiSuccessResponse(gin.H{"count": count}))
}

// MarkAsRead 标记单条通知为已读
func (ctl *Controller) MarkAsRead(c *gin.Context) {
	var req MarkReadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := model.ErrInvalidParams.WithDetail("请求参数格式错误")
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	userID := c.GetInt64(common.CtxUserIDKey)
	if err := ctl.service.MarkAsRead(req.ID, userID); err != nil {
		switch err {
		case ErrInvalidParams:
			appErr := model.ErrInvalidParams.WithDetail(err.Error())
			c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		default:
			c.JSON(model.ErrInternal.HTTPStatus(), model.ApiErrorResponse(model.ErrInternal.Code, model.ErrInternal.Message, nil))
		}
		return
	}
	c.JSON(http.StatusOK, model.ApiSuccessResponse(nil))
}

// MarkAllAsRead 标记所有通知为已读
func (ctl *Controller) MarkAllAsRead(c *gin.Context) {
	userID := c.GetInt64(common.CtxUserIDKey)
	if err := ctl.service.MarkAllAsRead(userID); err != nil {
		c.JSON(model.ErrInternal.HTTPStatus(), model.ApiErrorResponse(model.ErrInternal.Code, model.ErrInternal.Message, nil))
		return
	}
	c.JSON(http.StatusOK, model.ApiSuccessResponse(nil))
}

// Delete 删除通知
func (ctl *Controller) Delete(c *gin.Context) {
	var req DeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := model.ErrInvalidParams.WithDetail("请求参数格式错误")
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	userID := c.GetInt64(common.CtxUserIDKey)
	if err := ctl.service.Delete(req.ID, userID); err != nil {
		switch err {
		case ErrInvalidParams:
			appErr := model.ErrInvalidParams.WithDetail(err.Error())
			c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		default:
			c.JSON(model.ErrInternal.HTTPStatus(), model.ApiErrorResponse(model.ErrInternal.Code, model.ErrInternal.Message, nil))
		}
		return
	}
	c.JSON(http.StatusOK, model.ApiSuccessResponse(nil))
}
