package message

import (
	"ElainaBlog/internal/auth"
	"ElainaBlog/internal/response"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdminChecker 管理员权限检查接口。
type AdminChecker interface {
	CheckIsAdmin(ctx context.Context, userID int64) (bool, error)
}

type Controller struct {
	service      *Service
	adminChecker AdminChecker
}

func NewController(service *Service, adminChecker AdminChecker) *Controller {
	return &Controller{service: service, adminChecker: adminChecker}
}

func (ctl *Controller) GetList(c *gin.Context) {
	list, err := ctl.service.GetList(c.Request.Context(), 50)
	if err != nil {
		c.JSON(response.ErrInternal.HTTPStatus(), response.ApiErrorResponse(response.ErrInternal.Code, response.ErrInternal.Message, nil))
		return
	}
	c.JSON(http.StatusOK, response.ApiSuccessResponse(list))
}

func (ctl *Controller) Create(c *gin.Context) {
	var req CreateMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(response.ErrInvalidParams.HTTPStatus(), response.ApiErrorResponse(response.ErrInvalidParams.Code, response.ErrInvalidParams.Message, nil))
		return
	}

	userID := c.GetInt64(auth.CtxUserIDKey)
	msgID, err := ctl.service.Create(c.Request.Context(), userID, req.Content)
	if err != nil {
		switch err {
		case ErrInvalidParams:
			appErr := response.ErrInvalidParams.WithDetail(err.Error())
			c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		default:
			c.JSON(response.ErrInternal.HTTPStatus(), response.ApiErrorResponse(response.ErrInternal.Code, response.ErrInternal.Message, nil))
		}
		return
	}

	c.JSON(http.StatusOK, response.ApiSuccessResponse(gin.H{"id": msgID}))
}

func (ctl *Controller) Delete(c *gin.Context) {
	var req DeleteMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(response.ErrInvalidParams.HTTPStatus(), response.ApiErrorResponse(response.ErrInvalidParams.Code, response.ErrInvalidParams.Message, nil))
		return
	}

	userID := c.GetInt64(auth.CtxUserIDKey)

	msg, err := ctl.service.GetByID(c.Request.Context(), req.ID)
	if err != nil {
		switch err {
		case ErrMessageNotFound:
			appErr := response.ErrNotFound.WithDetail(err.Error())
			c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		default:
			c.JSON(response.ErrInternal.HTTPStatus(), response.ApiErrorResponse(response.ErrInternal.Code, response.ErrInternal.Message, nil))
		}
		return
	}

	if msg.UserID != userID {
		isAdmin, err := ctl.adminChecker.CheckIsAdmin(c.Request.Context(), userID)
		if err != nil {
			c.JSON(response.ErrInternal.HTTPStatus(), response.ApiErrorResponse(response.ErrInternal.Code, response.ErrInternal.Message, nil))
			return
		}
		if !isAdmin {
			appErr := response.ErrForbidden.WithDetail("仅留言作者或管理员可删除")
			c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
			return
		}
	}

	if err := ctl.service.Delete(c.Request.Context(), req.ID); err != nil {
		c.JSON(response.ErrInternal.HTTPStatus(), response.ApiErrorResponse(response.ErrInternal.Code, response.ErrInternal.Message, nil))
		return
	}

	c.JSON(http.StatusOK, response.ApiSuccessResponse(nil))
}
