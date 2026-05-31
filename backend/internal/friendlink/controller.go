package friendlink

import (
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

type CreateRequest struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	Avatar      string `json:"avatar"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
}

type UpdateRequest struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	Avatar      string `json:"avatar"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
}

type DeleteRequest struct {
	ID int64 `json:"id"`
}

// GetList 获取所有友情链接（公开）
func (ctl *Controller) GetList(c *gin.Context) {
	list, err := ctl.service.GetList()
	if err != nil {
		c.JSON(model.ErrInternal.HTTPStatus(), model.ApiErrorResponse(model.ErrInternal.Code, model.ErrInternal.Message, nil))
		return
	}
	c.JSON(http.StatusOK, model.ApiSuccessResponse(list))
}

// Create 创建友情链接（管理员）
func (ctl *Controller) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := model.ErrInvalidParams.WithDetail("请求参数格式错误")
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	vo, err := ctl.service.Create(CreateParams{
		Name:        req.Name,
		URL:         req.URL,
		Avatar:      req.Avatar,
		Description: req.Description,
		SortOrder:   req.SortOrder,
	})
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
	c.JSON(http.StatusOK, model.ApiSuccessResponse(vo))
}

// Update 更新友情链接（管理员）
func (ctl *Controller) Update(c *gin.Context) {
	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := model.ErrInvalidParams.WithDetail("请求参数格式错误")
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	vo, err := ctl.service.Update(UpdateParams{
		ID:          req.ID,
		Name:        req.Name,
		URL:         req.URL,
		Avatar:      req.Avatar,
		Description: req.Description,
		SortOrder:   req.SortOrder,
	})
	if err != nil {
		switch err {
		case ErrInvalidParams:
			appErr := model.ErrInvalidParams.WithDetail(err.Error())
			c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		case ErrLinkNotFound:
			appErr := model.ErrNotFound.WithDetail("友链不存在")
			c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		default:
			c.JSON(model.ErrInternal.HTTPStatus(), model.ApiErrorResponse(model.ErrInternal.Code, model.ErrInternal.Message, nil))
		}
		return
	}
	c.JSON(http.StatusOK, model.ApiSuccessResponse(vo))
}

// Delete 删除友情链接（管理员）
func (ctl *Controller) Delete(c *gin.Context) {
	var req DeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := model.ErrInvalidParams.WithDetail("请求参数格式错误")
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	err := ctl.service.Delete(req.ID)
	if err != nil {
		switch err {
		case ErrInvalidParams:
			appErr := model.ErrInvalidParams.WithDetail(err.Error())
			c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		case ErrLinkNotFound:
			appErr := model.ErrNotFound.WithDetail("友链不存在")
			c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		default:
			c.JSON(model.ErrInternal.HTTPStatus(), model.ApiErrorResponse(model.ErrInternal.Code, model.ErrInternal.Message, nil))
		}
		return
	}
	c.JSON(http.StatusOK, model.ApiSuccessResponse(nil))
}
