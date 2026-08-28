package category

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

// GetList 获取所有分类
func (ctl *Controller) GetList(c *gin.Context) {
	list, err := ctl.service.GetCategoryList(c.Request.Context())
	if err != nil {
		c.JSON(response.ErrInternal.HTTPStatus(), response.ApiErrorResponse(response.ErrInternal.Code, response.ErrInternal.Message, nil))
		return
	}

	c.JSON(http.StatusOK, response.ApiSuccessResponse(list))
}

// Create 创建分类（管理员）
func (ctl *Controller) Create(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := response.ErrInvalidParams.WithDetail("请求参数格式错误")
		c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	vo, err := ctl.service.CreateCategory(c.Request.Context(), CreateCategoryParams{Name: req.Name})
	if err != nil {
		switch err {
		case ErrInvalidParams:
			appErr := response.ErrInvalidParams.WithDetail(err.Error())
			c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		case ErrCategoryExists:
			appErr := response.ErrConflict.WithDetail("资源已存在")
			c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		default:
			c.JSON(response.ErrInternal.HTTPStatus(), response.ApiErrorResponse(response.ErrInternal.Code, response.ErrInternal.Message, nil))
		}
		return
	}

	c.JSON(http.StatusOK, response.ApiSuccessResponse(vo))
}

// Update 更新分类（管理员）
func (ctl *Controller) Update(c *gin.Context) {
	var req UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := response.ErrInvalidParams.WithDetail("请求参数格式错误")
		c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	vo, err := ctl.service.UpdateCategory(c.Request.Context(), UpdateCategoryParams{ID: req.ID, Name: req.Name})
	if err != nil {
		switch err {
		case ErrInvalidParams:
			appErr := response.ErrInvalidParams.WithDetail(err.Error())
			c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		case ErrCategoryNotFound:
			appErr := response.ErrNotFound.WithDetail("资源不存在")
			c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		case ErrCategoryExists:
			appErr := response.ErrConflict.WithDetail("资源已存在")
			c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		default:
			c.JSON(response.ErrInternal.HTTPStatus(), response.ApiErrorResponse(response.ErrInternal.Code, response.ErrInternal.Message, nil))
		}
		return
	}

	c.JSON(http.StatusOK, response.ApiSuccessResponse(vo))
}

// Delete 删除分类（管理员）
func (ctl *Controller) Delete(c *gin.Context) {
	var req DeleteCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := response.ErrInvalidParams.WithDetail("请求参数格式错误")
		c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	err := ctl.service.DeleteCategory(c.Request.Context(), req.ID)
	if err != nil {
		switch err {
		case ErrInvalidParams:
			appErr := response.ErrInvalidParams.WithDetail(err.Error())
			c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		case ErrCategoryNotFound:
			appErr := response.ErrNotFound.WithDetail("资源不存在")
			c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		default:
			c.JSON(response.ErrInternal.HTTPStatus(), response.ApiErrorResponse(response.ErrInternal.Code, response.ErrInternal.Message, nil))
		}
		return
	}

	c.JSON(http.StatusOK, response.ApiSuccessResponse(nil))
}

// ToggleTop 切换分类置顶状态（管理员）
func (ctl *Controller) ToggleTop(c *gin.Context) {
	var req ToggleTopRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := response.ErrInvalidParams.WithDetail("请求参数格式错误")
		c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	err := ctl.service.ToggleTop(c.Request.Context(), req.ID, req.IsTop)
	if err != nil {
		switch err {
		case ErrInvalidParams:
			appErr := response.ErrInvalidParams.WithDetail(err.Error())
			c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		case ErrCategoryNotFound:
			appErr := response.ErrNotFound.WithDetail("资源不存在")
			c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		default:
			c.JSON(response.ErrInternal.HTTPStatus(), response.ApiErrorResponse(response.ErrInternal.Code, response.ErrInternal.Message, nil))
		}
		return
	}

	c.JSON(http.StatusOK, response.ApiSuccessResponse(nil))
}
