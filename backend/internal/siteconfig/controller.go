package siteconfig

import (
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

// GetPublicConfigs 公开接口：返回前端需要的公开配置
func (ctl *Controller) GetPublicConfigs(c *gin.Context) {
	configs, err := ctl.service.GetAllMap(c.Request.Context())
	if err != nil {
		c.JSON(response.ErrInternal.HTTPStatus(), response.ApiErrorResponse(response.ErrInternal.Code, response.ErrInternal.Message, nil))
		return
	}
	c.JSON(http.StatusOK, response.ApiSuccessResponse(configs))
}

// GetQuotes 公开接口：返回随机语句列表
func (ctl *Controller) GetQuotes(c *gin.Context) {
	quotes, err := ctl.service.GetQuotes(c.Request.Context())
	if err != nil {
		c.JSON(response.ErrInternal.HTTPStatus(), response.ApiErrorResponse(response.ErrInternal.Code, response.ErrInternal.Message, nil))
		return
	}
	c.JSON(http.StatusOK, response.ApiSuccessResponse(quotes))
}

// GetAll 管理员接口：获取所有配置
func (ctl *Controller) GetAll(c *gin.Context) {
	configs, err := ctl.service.GetAllMap(c.Request.Context())
	if err != nil {
		c.JSON(response.ErrInternal.HTTPStatus(), response.ApiErrorResponse(response.ErrInternal.Code, response.ErrInternal.Message, nil))
		return
	}
	c.JSON(http.StatusOK, response.ApiSuccessResponse(configs))
}

// Upsert 管理员接口：批量更新配置
func (ctl *Controller) Upsert(c *gin.Context) {
	var req UpsertConfigsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(response.ErrInvalidParams.HTTPStatus(), response.ApiErrorResponse(response.ErrInvalidParams.Code, response.ErrInvalidParams.Message, nil))
		return
	}

	if err := ctl.service.Upsert(c.Request.Context(), req.Configs); err != nil {
		c.JSON(response.ErrInternal.HTTPStatus(), response.ApiErrorResponse(response.ErrInternal.Code, response.ErrInternal.Message, nil))
		return
	}

	c.JSON(http.StatusOK, response.ApiSuccessResponse(nil))
}

// Delete 管理员接口：删除配置
func (ctl *Controller) Delete(c *gin.Context) {
	var req DeleteConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(response.ErrInvalidParams.HTTPStatus(), response.ApiErrorResponse(response.ErrInvalidParams.Code, response.ErrInvalidParams.Message, nil))
		return
	}

	if err := ctl.service.Delete(c.Request.Context(), req.Key); err != nil {
		c.JSON(response.ErrInternal.HTTPStatus(), response.ApiErrorResponse(response.ErrInternal.Code, response.ErrInternal.Message, nil))
		return
	}

	c.JSON(http.StatusOK, response.ApiSuccessResponse(nil))
}
