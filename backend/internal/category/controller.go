package category

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

type CreateCategoryRequest struct {
	Name string `json:"name"`
}

type UpdateCategoryRequest struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type DeleteCategoryRequest struct {
	ID int64 `json:"id"`
}

// GetList 获取所有分类
func (ctl *Controller) GetList(c *gin.Context) {
	list, err := ctl.service.GetCategoryList()
	if err != nil {
		appErr := model.ErrInternal.WithDetail(err.Error())
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	c.JSON(http.StatusOK, model.ApiSuccessResponse(list))
}

// Create 创建分类（管理员）
func (ctl *Controller) Create(c *gin.Context) {
	if !common.RequireAdmin(c, ctl.userService.CheckIsAdmin) {
		return
	}

	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := model.ErrInvalidParams.WithDetail(err.Error())
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	vo, err := ctl.service.CreateCategory(CreateCategoryParams{Name: req.Name})
	if err != nil {
		switch err {
		case ErrInvalidParams:
			appErr := model.ErrInvalidParams.WithDetail(err.Error())
			c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		case ErrCategoryExists:
			appErr := model.ErrConflict.WithDetail(err.Error())
			c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		default:
			appErr := model.ErrInternal.WithDetail(err.Error())
			c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		}
		return
	}

	c.JSON(http.StatusOK, model.ApiSuccessResponse(vo))
}

// Update 更新分类（管理员）
func (ctl *Controller) Update(c *gin.Context) {
	if !common.RequireAdmin(c, ctl.userService.CheckIsAdmin) {
		return
	}

	var req UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := model.ErrInvalidParams.WithDetail(err.Error())
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	vo, err := ctl.service.UpdateCategory(UpdateCategoryParams{ID: req.ID, Name: req.Name})
	if err != nil {
		switch err {
		case ErrInvalidParams:
			appErr := model.ErrInvalidParams.WithDetail(err.Error())
			c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		case ErrCategoryNotFound:
			appErr := model.ErrNotFound.WithDetail(err.Error())
			c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		case ErrCategoryExists:
			appErr := model.ErrConflict.WithDetail(err.Error())
			c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		default:
			appErr := model.ErrInternal.WithDetail(err.Error())
			c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		}
		return
	}

	c.JSON(http.StatusOK, model.ApiSuccessResponse(vo))
}

// Delete 删除分类（管理员）
func (ctl *Controller) Delete(c *gin.Context) {
	if !common.RequireAdmin(c, ctl.userService.CheckIsAdmin) {
		return
	}

	var req DeleteCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := model.ErrInvalidParams.WithDetail(err.Error())
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	err := ctl.service.DeleteCategory(req.ID)
	if err != nil {
		switch err {
		case ErrInvalidParams:
			appErr := model.ErrInvalidParams.WithDetail(err.Error())
			c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		case ErrCategoryNotFound:
			appErr := model.ErrNotFound.WithDetail(err.Error())
			c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		default:
			appErr := model.ErrInternal.WithDetail(err.Error())
			c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		}
		return
	}

	c.JSON(http.StatusOK, model.ApiSuccessResponse(nil))
}
