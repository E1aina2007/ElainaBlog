package message

import (
	"ElainaBlog/internal/common"
	"ElainaBlog/internal/common/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdminChecker 管理员权限检查接口。
type AdminChecker interface {
	CheckIsAdmin(userID int64) (bool, error)
}

type Controller struct {
	service      *Service
	adminChecker AdminChecker
}

func NewController(service *Service, adminChecker AdminChecker) *Controller {
	return &Controller{service: service, adminChecker: adminChecker}
}

type CreateMessageRequest struct {
	Content string `json:"content"`
}

type DeleteMessageRequest struct {
	ID int64 `json:"id"`
}

func (ctl *Controller) GetList(c *gin.Context) {
	list, err := ctl.service.GetList(50)
	if err != nil {
		c.JSON(model.ErrInternal.HTTPStatus(), model.ApiErrorResponse(model.ErrInternal.Code, model.ErrInternal.Message, nil))
		return
	}
	c.JSON(http.StatusOK, model.ApiSuccessResponse(list))
}

func (ctl *Controller) Create(c *gin.Context) {
	var req CreateMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(model.ErrInvalidParams.HTTPStatus(), model.ApiErrorResponse(model.ErrInvalidParams.Code, model.ErrInvalidParams.Message, nil))
		return
	}

	userID := c.GetInt64(common.CtxUserIDKey)
	msgID, err := ctl.service.Create(userID, req.Content)
	if err != nil {
		switch err {
		case ErrInvalidParams:
			appErr := model.ErrInvalidParams.WithDetail(err.Error())
			c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		default:
			c.JSON(model.ErrInternal.HTTPStatus(), model.ApiErrorResponse(model.ErrInternal.Code, model.ErrInternal.Message, nil))
		}
		return
	}

	c.JSON(http.StatusOK, model.ApiSuccessResponse(gin.H{"id": msgID}))
}

func (ctl *Controller) Delete(c *gin.Context) {
	var req DeleteMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(model.ErrInvalidParams.HTTPStatus(), model.ApiErrorResponse(model.ErrInvalidParams.Code, model.ErrInvalidParams.Message, nil))
		return
	}

	userID := c.GetInt64(common.CtxUserIDKey)

	msg, err := ctl.service.GetByID(req.ID)
	if err != nil {
		switch err {
		case ErrMessageNotFound:
			appErr := model.ErrNotFound.WithDetail(err.Error())
			c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		default:
			c.JSON(model.ErrInternal.HTTPStatus(), model.ApiErrorResponse(model.ErrInternal.Code, model.ErrInternal.Message, nil))
		}
		return
	}

	if msg.UserID != userID {
		isAdmin, err := ctl.adminChecker.CheckIsAdmin(userID)
		if err != nil {
			c.JSON(model.ErrInternal.HTTPStatus(), model.ApiErrorResponse(model.ErrInternal.Code, model.ErrInternal.Message, nil))
			return
		}
		if !isAdmin {
			appErr := model.ErrForbidden.WithDetail("仅留言作者或管理员可删除")
			c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
			return
		}
	}

	if err := ctl.service.Delete(req.ID); err != nil {
		c.JSON(model.ErrInternal.HTTPStatus(), model.ApiErrorResponse(model.ErrInternal.Code, model.ErrInternal.Message, nil))
		return
	}

	c.JSON(http.StatusOK, model.ApiSuccessResponse(nil))
}
