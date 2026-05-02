package article

import (
	"ElainaBlog/config/db"
	"ElainaBlog/internal/common"
	"ElainaBlog/internal/common/model"
	"ElainaBlog/internal/user"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service     *Service
	userService *user.Service // 用于验证管理员权限
}

func NewController(userService *user.Service) *Controller {
	repo := NewRepository(db.DBPool)
	service := NewService(repo)
	return &Controller{service: service, userService: userService}
}

type CreateArticleRequest struct {
	Title      string `json:"title"`
	Summary    string `json:"summary"`
	Content    string `json:"content"`
	CategoryID *int64 `json:"category_id"` // nil 表示未分类
	Cover      string `json:"cover"`
	IsTop      bool   `json:"is_top"`
	IsDraft    bool   `json:"is_draft"`
}

type UpdateArticleRequest struct {
	ID         int64  `json:"id"`
	Title      string `json:"title"`
	Summary    string `json:"summary"`
	Content    string `json:"content"`
	CategoryID *int64 `json:"category_id"`
	Cover      string `json:"cover"`
	IsTop      bool   `json:"is_top"`
	IsDraft    bool   `json:"is_draft"`
}

type DeleteArticleRequest struct {
	ID int64 `json:"id"`
}

func (ctl *Controller) CreateArticle(c *gin.Context) {
	if !common.RequireAdmin(c, ctl.userService.CheckIsAdmin) {
		return
	}

	var req CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := model.ErrInvalidParams.WithDetail(err.Error())
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	userID := c.GetInt64(common.CtxUserIDKey)
	articleID, err := ctl.service.CreateArticle(&CreateArticleParams{
		UserID:     userID,
		CategoryID: req.CategoryID,
		Title:      req.Title,
		Summary:    req.Summary,
		Content:    req.Content,
		Cover:      req.Cover,
		IsTop:      req.IsTop,
		IsDraft:    req.IsDraft,
	})
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

	c.JSON(http.StatusOK, model.ApiSuccessResponse(gin.H{"id": articleID}))
}

func (ctl *Controller) UpdateArticle(c *gin.Context) {
	if !common.RequireAdmin(c, ctl.userService.CheckIsAdmin) {
		return
	}

	var req UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := model.ErrInvalidParams.WithDetail(err.Error())
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	err := ctl.service.UpdateArticle(&UpdateArticleParams{
		ID:         req.ID,
		CategoryID: req.CategoryID,
		Title:      req.Title,
		Summary:    req.Summary,
		Content:    req.Content,
		Cover:      req.Cover,
		IsTop:      req.IsTop,
		IsDraft:    req.IsDraft,
	})
	if err != nil {
		switch err {
		case ErrInvalidParams:
			appErr := model.ErrInvalidParams.WithDetail(err.Error())
			c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		case ErrArticleNotFound:
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

func (ctl *Controller) DeleteArticle(c *gin.Context) {
	if !common.RequireAdmin(c, ctl.userService.CheckIsAdmin) {
		return
	}

	var req DeleteArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := model.ErrInvalidParams.WithDetail(err.Error())
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	err := ctl.service.DeleteArticle(&DeleteArticleParams{ID: req.ID})
	if err != nil {
		switch err {
		case ErrInvalidParams:
			appErr := model.ErrInvalidParams.WithDetail(err.Error())
			c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		case ErrArticleNotFound:
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

// GetList 文章列表（公开），支持分页和分类筛选
func (ctl *Controller) GetList(c *gin.Context) {
	// 解析分页参数
	page := 1
	pageSize := 10
	if p, err := strconv.Atoi(c.Query("page")); err == nil && p > 0 {
		page = p
	}
	if ps, err := strconv.Atoi(c.Query("pageSize")); err == nil && ps > 0 {
		pageSize = ps
	}

	// 解析分类筛选参数
	var categoryID *int64
	if catIDStr := c.Query("categoryId"); catIDStr != "" {
		if catID, err := strconv.ParseInt(catIDStr, 10, 64); err == nil && catID > 0 {
			categoryID = &catID
		}
	}

	result, err := ctl.service.GetArticleList(&ArticleListParams{
		CategoryID: categoryID,
		Page:       page,
		PageSize:   pageSize,
	})
	if err != nil {
		appErr := model.ErrInternal.WithDetail(err.Error())
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	c.JSON(http.StatusOK, model.ApiSuccessResponse(result))
}

// GetByID 文章详情（公开）
func (ctl *Controller) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		appErr := model.ErrInvalidParams.WithDetail("无效的文章 ID")
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	article, err := ctl.service.GetArticleByID(id)
	if err != nil {
		switch err {
		case ErrArticleNotFound:
			appErr := model.ErrNotFound.WithDetail(err.Error())
			c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		default:
			appErr := model.ErrInternal.WithDetail(err.Error())
			c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		}
		return
	}

	// 异步增加浏览量（不阻塞响应）
	go ctl.service.IncrementViewCount(id)

	c.JSON(http.StatusOK, model.ApiSuccessResponse(article))
}
