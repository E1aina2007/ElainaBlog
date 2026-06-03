package article

import (
	"ElainaBlog/internal/common"
	"ElainaBlog/internal/common/model"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AdminChecker 管理员权限检查接口，避免直接依赖 user.Service。
type AdminChecker interface {
	CheckIsAdmin(userID int64) (bool, error)
}

type Controller struct {
	service       *Service
	adminChecker  AdminChecker
}

func NewController(service *Service, adminChecker AdminChecker) *Controller {
	return &Controller{service: service, adminChecker: adminChecker}
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
	var req CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := model.ErrInvalidParams.WithDetail("请求参数格式错误")
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
			c.JSON(model.ErrInternal.HTTPStatus(), model.ApiErrorResponse(model.ErrInternal.Code, model.ErrInternal.Message, nil))
		}
		return
	}

	c.JSON(http.StatusOK, model.ApiSuccessResponse(gin.H{"id": articleID}))
}

func (ctl *Controller) UpdateArticle(c *gin.Context) {
	var req UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := model.ErrInvalidParams.WithDetail("请求参数格式错误")
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	userID := c.GetInt64(common.CtxUserIDKey)
	isAdmin, _ := ctl.adminChecker.CheckIsAdmin(userID)

	err := ctl.service.UpdateArticle(&UpdateArticleParams{
		ID:         req.ID,
		CategoryID: req.CategoryID,
		Title:      req.Title,
		Summary:    req.Summary,
		Content:    req.Content,
		Cover:      req.Cover,
		IsTop:      req.IsTop,
		IsDraft:    req.IsDraft,
	}, userID, isAdmin)
	if err != nil {
		switch err {
		case ErrInvalidParams:
			appErr := model.ErrInvalidParams.WithDetail(err.Error())
			c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		case ErrArticleNotFound:
			appErr := model.ErrNotFound.WithDetail("资源不存在")
			c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		case ErrNoPermission:
			appErr := model.ErrForbidden.WithDetail("无权限")
			c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		default:
			c.JSON(model.ErrInternal.HTTPStatus(), model.ApiErrorResponse(model.ErrInternal.Code, model.ErrInternal.Message, nil))
		}
		return
	}

	c.JSON(http.StatusOK, model.ApiSuccessResponse(nil))
}

func (ctl *Controller) DeleteArticle(c *gin.Context) {
	var req DeleteArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := model.ErrInvalidParams.WithDetail("请求参数格式错误")
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	userID := c.GetInt64(common.CtxUserIDKey)
	isAdmin, _ := ctl.adminChecker.CheckIsAdmin(userID)

	err := ctl.service.DeleteArticle(&DeleteArticleParams{ID: req.ID}, userID, isAdmin)
	if err != nil {
		switch err {
		case ErrInvalidParams:
			appErr := model.ErrInvalidParams.WithDetail(err.Error())
			c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		case ErrArticleNotFound:
			appErr := model.ErrNotFound.WithDetail("资源不存在")
			c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		case ErrNoPermission:
			appErr := model.ErrForbidden.WithDetail("无权限")
			c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		default:
			c.JSON(model.ErrInternal.HTTPStatus(), model.ApiErrorResponse(model.ErrInternal.Code, model.ErrInternal.Message, nil))
		}
		return
	}

	c.JSON(http.StatusOK, model.ApiSuccessResponse(nil))
}

// parseListParams 解析列表分页和分类筛选参数
func parseListParams(c *gin.Context) (*ArticleListParams, error) {
	page := 1
	pageSize := 10
	if p, err := strconv.Atoi(c.Query("page")); err == nil && p > 0 {
		page = p
	}
	if ps, err := strconv.Atoi(c.Query("pageSize")); err == nil && ps > 0 {
		pageSize = ps
	}
	var categoryID *int64
	if catIDStr := c.Query("categoryId"); catIDStr != "" {
		if catID, err := strconv.ParseInt(catIDStr, 10, 64); err == nil && catID > 0 {
			categoryID = &catID
		}
	}
	return &ArticleListParams{CategoryID: categoryID, Page: page, PageSize: pageSize}, nil
}

// GetList 文章列表（公开），支持分页和分类筛选
func (ctl *Controller) GetList(c *gin.Context) {
	params, err := parseListParams(c)
	if err != nil {
		c.JSON(model.ErrInternal.HTTPStatus(), model.ApiErrorResponse(model.ErrInternal.Code, model.ErrInternal.Message, nil))
		return
	}
	result, err := ctl.service.GetArticleList(params)
	if err != nil {
		c.JSON(model.ErrInternal.HTTPStatus(), model.ApiErrorResponse(model.ErrInternal.Code, model.ErrInternal.Message, nil))
		return
	}
	c.JSON(http.StatusOK, model.ApiSuccessResponse(result))
}

// GetMyList 当前用户的文章列表（含草稿）
func (ctl *Controller) GetMyList(c *gin.Context) {
	params, err := parseListParams(c)
	if err != nil {
		c.JSON(model.ErrInternal.HTTPStatus(), model.ApiErrorResponse(model.ErrInternal.Code, model.ErrInternal.Message, nil))
		return
	}
	userID := c.GetInt64(common.CtxUserIDKey)
	result, err := ctl.service.GetUserArticleList(userID, params)
	if err != nil {
		c.JSON(model.ErrInternal.HTTPStatus(), model.ApiErrorResponse(model.ErrInternal.Code, model.ErrInternal.Message, nil))
		return
	}
	c.JSON(http.StatusOK, model.ApiSuccessResponse(result))
}

// GetAdminList 文章列表（管理员），包含草稿
func (ctl *Controller) GetAdminList(c *gin.Context) {
	params, err := parseListParams(c)
	if err != nil {
		c.JSON(model.ErrInternal.HTTPStatus(), model.ApiErrorResponse(model.ErrInternal.Code, model.ErrInternal.Message, nil))
		return
	}
	result, err := ctl.service.GetAdminArticleList(params)
	if err != nil {
		c.JSON(model.ErrInternal.HTTPStatus(), model.ApiErrorResponse(model.ErrInternal.Code, model.ErrInternal.Message, nil))
		return
	}
	c.JSON(http.StatusOK, model.ApiSuccessResponse(result))
}

// parseArticleID 解析文章 ID 参数
func parseArticleID(c *gin.Context) (int64, error) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("无效的文章 ID")
	}
	return id, nil
}

// GetByID 文章详情（公开，过滤草稿）
func (ctl *Controller) GetByID(c *gin.Context) {
	id, err := parseArticleID(c)
	if err != nil {
		appErr := model.ErrInvalidParams.WithDetail(err.Error())
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	article, err := ctl.service.GetArticleByID(id)
	if err != nil {
		switch err {
		case ErrArticleNotFound:
			appErr := model.ErrNotFound.WithDetail("资源不存在")
			c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		default:
			c.JSON(model.ErrInternal.HTTPStatus(), model.ApiErrorResponse(model.ErrInternal.Code, model.ErrInternal.Message, nil))
		}
		return
	}

	// 异步增加浏览量（IP 去重，不阻塞响应）
	go ctl.service.IncrementViewCount(id, c.ClientIP())

	c.JSON(http.StatusOK, model.ApiSuccessResponse(article))
}

// GetAdminByID 文章详情（管理员，包含草稿）
func (ctl *Controller) GetAdminByID(c *gin.Context) {
	id, err := parseArticleID(c)
	if err != nil {
		appErr := model.ErrInvalidParams.WithDetail(err.Error())
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	article, err := ctl.service.GetArticleByIDIncludeDraft(id)
	if err != nil {
		switch err {
		case ErrArticleNotFound:
			appErr := model.ErrNotFound.WithDetail("资源不存在")
			c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		default:
			c.JSON(model.ErrInternal.HTTPStatus(), model.ApiErrorResponse(model.ErrInternal.Code, model.ErrInternal.Message, nil))
		}
		return
	}

	c.JSON(http.StatusOK, model.ApiSuccessResponse(article))
}

// GetArticleUV 获取文章的独立访客数（管理员）
func (ctl *Controller) GetArticleUV(c *gin.Context) {
	id, err := parseArticleID(c)
	if err != nil {
		appErr := model.ErrInvalidParams.WithDetail(err.Error())
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	uv, err := ctl.service.GetArticleUV(id)
	if err != nil {
		c.JSON(model.ErrInternal.HTTPStatus(), model.ApiErrorResponse(model.ErrInternal.Code, model.ErrInternal.Message, nil))
		return
	}

	c.JSON(http.StatusOK, model.ApiSuccessResponse(gin.H{"uv": uv}))
}
