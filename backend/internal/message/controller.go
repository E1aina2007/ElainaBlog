package message

import (
	"ElainaBlog/config/db"
	"ElainaBlog/internal/common"
	"ElainaBlog/internal/common/model"
	"ElainaBlog/internal/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service     *Service
	userService *user.Service
}

func NewController(userService *user.Service) *Controller {
	repo := NewRepository(db.DBPool)
	service := NewService(repo)
	return &Controller{service: service, userService: userService}
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
		appErr := model.ErrInternal.WithDetail(err.Error())
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}
	c.JSON(http.StatusOK, model.ApiSuccessResponse(list))
}

func (ctl *Controller) Create(c *gin.Context) {
	var req CreateMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := model.ErrInvalidParams.WithDetail(err.Error())
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
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
			appErr := model.ErrInternal.WithDetail(err.Error())
			c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		}
		return
	}

	c.JSON(http.StatusOK, model.ApiSuccessResponse(gin.H{"id": msgID}))
}

func (ctl *Controller) Delete(c *gin.Context) {
	var req DeleteMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := model.ErrInvalidParams.WithDetail(err.Error())
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
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
			appErr := model.ErrInternal.WithDetail(err.Error())
			c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		}
		return
	}

	if msg.UserID != userID {
		isAdmin, err := ctl.userService.CheckIsAdmin(userID)
		if err != nil || !isAdmin {
			appErr := model.ErrForbidden.WithDetail("仅留言作者或管理员可删除")
			c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
			return
		}
	}

	if err := ctl.service.Delete(req.ID); err != nil {
		appErr := model.ErrInternal.WithDetail(err.Error())
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	c.JSON(http.StatusOK, model.ApiSuccessResponse(nil))
}
