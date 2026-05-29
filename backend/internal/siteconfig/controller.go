package siteconfig

import (
	"ElainaBlog/config/db"
	"ElainaBlog/internal/common/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service *Service
}

func NewController() *Controller {
	repo := NewRepository(db.DBPool)
	service := NewService(repo)
	return &Controller{service: service}
}

type UpsertConfigsRequest struct {
	Configs map[string]string `json:"configs"`
}

type DeleteConfigRequest struct {
	Key string `json:"key"`
}

// GetPublicConfigs 公开接口：返回前端需要的公开配置
func (ctl *Controller) GetPublicConfigs(c *gin.Context) {
	configs, err := ctl.service.GetAllMap()
	if err != nil {
		c.JSON(model.ErrInternal.HTTPStatus(), model.ApiErrorResponse(model.ErrInternal.Code, model.ErrInternal.Message, nil))
		return
	}
	c.JSON(http.StatusOK, model.ApiSuccessResponse(configs))
}

// GetQuotes 公开接口：返回随机语句列表
func (ctl *Controller) GetQuotes(c *gin.Context) {
	quotes, err := ctl.service.GetQuotes()
	if err != nil {
		c.JSON(model.ErrInternal.HTTPStatus(), model.ApiErrorResponse(model.ErrInternal.Code, model.ErrInternal.Message, nil))
		return
	}
	c.JSON(http.StatusOK, model.ApiSuccessResponse(quotes))
}

// GetAll 管理员接口：获取所有配置
func (ctl *Controller) GetAll(c *gin.Context) {
	configs, err := ctl.service.GetAllMap()
	if err != nil {
		c.JSON(model.ErrInternal.HTTPStatus(), model.ApiErrorResponse(model.ErrInternal.Code, model.ErrInternal.Message, nil))
		return
	}
	c.JSON(http.StatusOK, model.ApiSuccessResponse(configs))
}

// Upsert 管理员接口：批量更新配置
func (ctl *Controller) Upsert(c *gin.Context) {
	var req UpsertConfigsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(model.ErrInvalidParams.HTTPStatus(), model.ApiErrorResponse(model.ErrInvalidParams.Code, model.ErrInvalidParams.Message, nil))
		return
	}

	if err := ctl.service.Upsert(req.Configs); err != nil {
		c.JSON(model.ErrInternal.HTTPStatus(), model.ApiErrorResponse(model.ErrInternal.Code, model.ErrInternal.Message, nil))
		return
	}

	c.JSON(http.StatusOK, model.ApiSuccessResponse(nil))
}

// Delete 管理员接口：删除配置
func (ctl *Controller) Delete(c *gin.Context) {
	var req DeleteConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(model.ErrInvalidParams.HTTPStatus(), model.ApiErrorResponse(model.ErrInvalidParams.Code, model.ErrInvalidParams.Message, nil))
		return
	}

	if err := ctl.service.Delete(req.Key); err != nil {
		c.JSON(model.ErrInternal.HTTPStatus(), model.ApiErrorResponse(model.ErrInternal.Code, model.ErrInternal.Message, nil))
		return
	}

	c.JSON(http.StatusOK, model.ApiSuccessResponse(nil))
}
