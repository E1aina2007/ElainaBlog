package site

import (
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

func NewController(service *Service, userService *user.Service) *Controller {
	return &Controller{service: service, userService: userService}
}

func (ctl *Controller) GetList(c *gin.Context) {
	list, err := ctl.service.GetList()
	if err != nil {
		appErr := model.ErrInternal.WithDetail(err.Error())
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	c.JSON(http.StatusOK, model.ApiSuccessResponse(list))
}

// Update 更新站点配置（管理员）
func (ctl *Controller) Update(c *gin.Context) {
	if !common.RequireAdmin(c, ctl.userService.CheckIsAdmin) {
		return
	}
	var configs map[string]string
	if err := c.ShouldBindJSON(&configs); err != nil {
		appErr := model.ErrInvalidParams.WithDetail(err.Error())
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	if err := ctl.service.Update(configs); err != nil {
		appErr := model.ErrInternal.WithDetail(err.Error())
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	c.JSON(http.StatusOK, model.ApiSuccessResponse(nil))
}
