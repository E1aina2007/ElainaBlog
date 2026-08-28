package friendlink

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

// GetList 获取所有友情链接（公开）
func (ctl *Controller) GetList(c *gin.Context) {
	list, err := ctl.service.GetList(c.Request.Context())
	if err != nil {
		c.JSON(response.ErrInternal.HTTPStatus(), response.ApiErrorResponse(response.ErrInternal.Code, response.ErrInternal.Message, nil))
		return
	}
	c.JSON(http.StatusOK, response.ApiSuccessResponse(list))
}

// Create 创建友情链接（管理员）
func (ctl *Controller) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := response.ErrInvalidParams.WithDetail("请求参数格式错误")
		c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	vo, err := ctl.service.Create(c.Request.Context(), CreateParams{
		Name:        req.Name,
		URL:         req.URL,
		Avatar:      req.Avatar,
		Description: req.Description,
		SortOrder:   req.SortOrder,
	})
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
	c.JSON(http.StatusOK, response.ApiSuccessResponse(vo))
}

// Update 更新友情链接（管理员）
func (ctl *Controller) Update(c *gin.Context) {
	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := response.ErrInvalidParams.WithDetail("请求参数格式错误")
		c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	vo, err := ctl.service.Update(c.Request.Context(), UpdateParams{
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
			appErr := response.ErrInvalidParams.WithDetail(err.Error())
			c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		case ErrLinkNotFound:
			appErr := response.ErrNotFound.WithDetail("友链不存在")
			c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		default:
			c.JSON(response.ErrInternal.HTTPStatus(), response.ApiErrorResponse(response.ErrInternal.Code, response.ErrInternal.Message, nil))
		}
		return
	}
	c.JSON(http.StatusOK, response.ApiSuccessResponse(vo))
}

// Delete 删除友情链接（管理员）
func (ctl *Controller) Delete(c *gin.Context) {
	var req DeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := response.ErrInvalidParams.WithDetail("请求参数格式错误")
		c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	err := ctl.service.Delete(c.Request.Context(), req.ID)
	if err != nil {
		switch err {
		case ErrInvalidParams:
			appErr := response.ErrInvalidParams.WithDetail(err.Error())
			c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		case ErrLinkNotFound:
			appErr := response.ErrNotFound.WithDetail("友链不存在")
			c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		default:
			c.JSON(response.ErrInternal.HTTPStatus(), response.ApiErrorResponse(response.ErrInternal.Code, response.ErrInternal.Message, nil))
		}
		return
	}
	c.JSON(http.StatusOK, response.ApiSuccessResponse(nil))
}
