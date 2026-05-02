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

// GetAuthorInfo 获取作者公开信息（公开接口）
func (ctl *Controller) GetAuthorInfo(c *gin.Context) {
	info, err := ctl.service.GetAuthorInfo()
	if err != nil {
		appErr := model.ErrInternal.WithDetail(err.Error())
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	c.JSON(http.StatusOK, model.ApiSuccessResponse(info))
}

// GetAuthorStats 获取作者统计数据（公开接口，默认查询管理员）
func (ctl *Controller) GetAuthorStats(c *gin.Context) {
	// 获取管理员列表（简单实现：获取第一个管理员）
	users, err := ctl.userService.GetList()
	if err != nil {
		appErr := model.ErrInternal.WithDetail(err.Error())
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	var adminID int64
	for _, u := range users {
		if u.IsAdmin {
			adminID = u.ID
			break
		}
	}

	if adminID == 0 {
		appErr := model.ErrNotFound.WithDetail("未找到管理员用户")
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	stats, err := ctl.userService.GetAuthorStats(adminID)
	if err != nil {
		switch err {
		case user.ErrUserNotFound:
			appErr := model.ErrNotFound.WithDetail(err.Error())
			c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		case user.ErrForbidden:
			appErr := model.ErrForbidden.WithDetail(err.Error())
			c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		default:
			appErr := model.ErrInternal.WithDetail(err.Error())
			c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		}
		return
	}

	c.JSON(http.StatusOK, model.ApiSuccessResponse(stats))
}

// GetDashboardStats 获取仪表盘统计数据（管理员）
func (ctl *Controller) GetDashboardStats(c *gin.Context) {
	if !common.RequireAdmin(c, ctl.userService.CheckIsAdmin) {
		return
	}
	stats, err := ctl.service.GetDashboardStats()
	if err != nil {
		appErr := model.ErrInternal.WithDetail(err.Error())
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}
	c.JSON(http.StatusOK, model.ApiSuccessResponse(stats))
}
