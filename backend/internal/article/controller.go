package article

import (
	"ElainaBlog/internal/auth"
	"ElainaBlog/internal/response"
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AdminChecker 管理员权限检查接口，避免直接依赖 user.Service。
type AdminChecker interface {
	CheckIsAdmin(ctx context.Context, userID int64) (bool, error)
}

type Controller struct {
	service       *Service
	adminChecker  AdminChecker
}

func NewController(service *Service, adminChecker AdminChecker) *Controller {
	return &Controller{service: service, adminChecker: adminChecker}
}

func (ctl *Controller) CreateArticle(c *gin.Context) {
	var req CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := response.ErrInvalidParams.WithDetail("请求参数格式错误")
		c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	userID := c.GetInt64(auth.CtxUserIDKey)
	articleID, err := ctl.service.CreateArticle(c.Request.Context(), &CreateArticleParams{
		UserID:     userID,
		CategoryID: req.CategoryID,
		Title:      req.Title,
		Summary:    req.Summary,
		Content:    req.Content,
		Tags:       req.Tags,
		IsTop:      req.IsTop,
		IsDraft:    req.IsDraft,
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

	c.JSON(http.StatusOK, response.ApiSuccessResponse(gin.H{"id": articleID}))
}

func (ctl *Controller) UpdateArticle(c *gin.Context) {
	var req UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := response.ErrInvalidParams.WithDetail("请求参数格式错误")
		c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	userID := c.GetInt64(auth.CtxUserIDKey)
	isAdmin, _ := ctl.adminChecker.CheckIsAdmin(c.Request.Context(), userID)

	err := ctl.service.UpdateArticle(c.Request.Context(), &UpdateArticleParams{
		ID:         req.ID,
		CategoryID: req.CategoryID,
		Title:      req.Title,
		Summary:    req.Summary,
		Content:    req.Content,
		Tags:       req.Tags,
		IsTop:      req.IsTop,
		IsDraft:    req.IsDraft,
	}, userID, isAdmin)
	if err != nil {
		switch err {
		case ErrInvalidParams:
			appErr := response.ErrInvalidParams.WithDetail(err.Error())
			c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		case ErrArticleNotFound:
			appErr := response.ErrNotFound.WithDetail("资源不存在")
			c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		case ErrNoPermission:
			appErr := response.ErrForbidden.WithDetail("无权限")
			c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		default:
			c.JSON(response.ErrInternal.HTTPStatus(), response.ApiErrorResponse(response.ErrInternal.Code, response.ErrInternal.Message, nil))
		}
		return
	}

	c.JSON(http.StatusOK, response.ApiSuccessResponse(nil))
}

func (ctl *Controller) DeleteArticle(c *gin.Context) {
	var req DeleteArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := response.ErrInvalidParams.WithDetail("请求参数格式错误")
		c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	userID := c.GetInt64(auth.CtxUserIDKey)
	isAdmin, _ := ctl.adminChecker.CheckIsAdmin(c.Request.Context(), userID)

	err := ctl.service.DeleteArticle(c.Request.Context(), &DeleteArticleParams{ID: req.ID}, userID, isAdmin)
	if err != nil {
		switch err {
		case ErrInvalidParams:
			appErr := response.ErrInvalidParams.WithDetail(err.Error())
			c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		case ErrArticleNotFound:
			appErr := response.ErrNotFound.WithDetail("资源不存在")
			c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		case ErrNoPermission:
			appErr := response.ErrForbidden.WithDetail("无权限")
			c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		default:
			c.JSON(response.ErrInternal.HTTPStatus(), response.ApiErrorResponse(response.ErrInternal.Code, response.ErrInternal.Message, nil))
		}
		return
	}

	c.JSON(http.StatusOK, response.ApiSuccessResponse(nil))
}

// ToggleTop 切换文章置顶状态（仅管理员），只修改 is_top 字段不影响其他字段
func (ctl *Controller) ToggleTop(c *gin.Context) {
	var req ToggleTopRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := response.ErrInvalidParams.WithDetail("请求参数格式错误")
		c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	userID := c.GetInt64(auth.CtxUserIDKey)
	isAdmin, _ := ctl.adminChecker.CheckIsAdmin(c.Request.Context(), userID)

	err := ctl.service.ToggleTop(c.Request.Context(), req.ID, req.IsTop, isAdmin)
	if err != nil {
		switch err {
		case ErrInvalidParams:
			appErr := response.ErrInvalidParams.WithDetail(err.Error())
			c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		case ErrArticleNotFound:
			appErr := response.ErrNotFound.WithDetail("资源不存在")
			c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		case ErrNoPermission:
			appErr := response.ErrForbidden.WithDetail("无权限")
			c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		default:
			c.JSON(response.ErrInternal.HTTPStatus(), response.ApiErrorResponse(response.ErrInternal.Code, response.ErrInternal.Message, nil))
		}
		return
	}

	c.JSON(http.StatusOK, response.ApiSuccessResponse(nil))
}

// parseListParams 解析列表分页、分类筛选和排序参数
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
	sortBy := c.DefaultQuery("sortBy", "latest")
	return &ArticleListParams{CategoryID: categoryID, SortBy: sortBy, Page: page, PageSize: pageSize}, nil
}

// GetList 文章列表（公开），支持分页和分类筛选
func (ctl *Controller) GetList(c *gin.Context) {
	params, err := parseListParams(c)
	if err != nil {
		c.JSON(response.ErrInternal.HTTPStatus(), response.ApiErrorResponse(response.ErrInternal.Code, response.ErrInternal.Message, nil))
		return
	}
	result, err := ctl.service.GetArticleList(c.Request.Context(), params)
	if err != nil {
		c.JSON(response.ErrInternal.HTTPStatus(), response.ApiErrorResponse(response.ErrInternal.Code, response.ErrInternal.Message, nil))
		return
	}
	c.JSON(http.StatusOK, response.ApiSuccessResponse(result))
}

// GetMyList 当前用户的文章列表（含草稿）
func (ctl *Controller) GetMyList(c *gin.Context) {
	params, err := parseListParams(c)
	if err != nil {
		c.JSON(response.ErrInternal.HTTPStatus(), response.ApiErrorResponse(response.ErrInternal.Code, response.ErrInternal.Message, nil))
		return
	}
	userID := c.GetInt64(auth.CtxUserIDKey)
	result, err := ctl.service.GetUserArticleList(c.Request.Context(), userID, params)
	if err != nil {
		c.JSON(response.ErrInternal.HTTPStatus(), response.ApiErrorResponse(response.ErrInternal.Code, response.ErrInternal.Message, nil))
		return
	}
	c.JSON(http.StatusOK, response.ApiSuccessResponse(result))
}

// GetAdminList 文章列表（管理员），包含草稿
func (ctl *Controller) GetAdminList(c *gin.Context) {
	params, err := parseListParams(c)
	if err != nil {
		c.JSON(response.ErrInternal.HTTPStatus(), response.ApiErrorResponse(response.ErrInternal.Code, response.ErrInternal.Message, nil))
		return
	}
	result, err := ctl.service.GetAdminArticleList(c.Request.Context(), params)
	if err != nil {
		c.JSON(response.ErrInternal.HTTPStatus(), response.ApiErrorResponse(response.ErrInternal.Code, response.ErrInternal.Message, nil))
		return
	}
	c.JSON(http.StatusOK, response.ApiSuccessResponse(result))
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
		appErr := response.ErrInvalidParams.WithDetail(err.Error())
		c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	article, err := ctl.service.GetArticleByID(c.Request.Context(), id)
	if err != nil {
		switch err {
		case ErrArticleNotFound:
			appErr := response.ErrNotFound.WithDetail("资源不存在")
			c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		default:
			c.JSON(response.ErrInternal.HTTPStatus(), response.ApiErrorResponse(response.ErrInternal.Code, response.ErrInternal.Message, nil))
		}
		return
	}

	// 异步增加浏览量（IP 去重，不阻塞响应）
	go ctl.service.IncrementViewCount(context.Background(), id, c.ClientIP())

	c.JSON(http.StatusOK, response.ApiSuccessResponse(article))
}

// GetAdminByID 文章详情（管理员，包含草稿）
func (ctl *Controller) GetAdminByID(c *gin.Context) {
	id, err := parseArticleID(c)
	if err != nil {
		appErr := response.ErrInvalidParams.WithDetail(err.Error())
		c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	article, err := ctl.service.GetArticleByIDIncludeDraft(c.Request.Context(), id)
	if err != nil {
		switch err {
		case ErrArticleNotFound:
			appErr := response.ErrNotFound.WithDetail("资源不存在")
			c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		default:
			c.JSON(response.ErrInternal.HTTPStatus(), response.ApiErrorResponse(response.ErrInternal.Code, response.ErrInternal.Message, nil))
		}
		return
	}

	c.JSON(http.StatusOK, response.ApiSuccessResponse(article))
}

// GetMyByID 文章详情（当前用户，含草稿，仅限自己的文章）
func (ctl *Controller) GetMyByID(c *gin.Context) {
	id, err := parseArticleID(c)
	if err != nil {
		appErr := response.ErrInvalidParams.WithDetail(err.Error())
		c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	article, err := ctl.service.GetArticleByIDIncludeDraft(c.Request.Context(), id)
	if err != nil {
		switch err {
		case ErrArticleNotFound:
			appErr := response.ErrNotFound.WithDetail("资源不存在")
			c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		default:
			c.JSON(response.ErrInternal.HTTPStatus(), response.ApiErrorResponse(response.ErrInternal.Code, response.ErrInternal.Message, nil))
		}
		return
	}

	// 校验文章归属：非管理员只能查看自己的文章
	userID := c.GetInt64(auth.CtxUserIDKey)
	isAdmin, _ := ctl.adminChecker.CheckIsAdmin(c.Request.Context(), userID)
	if !isAdmin && article.UserID != userID {
		appErr := response.ErrForbidden.WithDetail("无权限访问此文章")
		c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	c.JSON(http.StatusOK, response.ApiSuccessResponse(article))
}

// Search 文章全文搜索（公开）
func (ctl *Controller) Search(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		appErr := response.ErrInvalidParams.WithDetail("搜索关键词不能为空")
		c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	page := 1
	pageSize := 10
	if p, err := strconv.Atoi(c.Query("page")); err == nil && p > 0 {
		page = p
	}
	if ps, err := strconv.Atoi(c.Query("pageSize")); err == nil && ps > 0 {
		pageSize = ps
	}

	result, err := ctl.service.SearchArticles(c.Request.Context(), keyword, page, pageSize)
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

	c.JSON(http.StatusOK, response.ApiSuccessResponse(result))
}

// GetArticleUV 获取文章的独立访客数（管理员）
func (ctl *Controller) GetArticleUV(c *gin.Context) {
	id, err := parseArticleID(c)
	if err != nil {
		appErr := response.ErrInvalidParams.WithDetail(err.Error())
		c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	uv, err := ctl.service.GetArticleUV(c.Request.Context(), id)
	if err != nil {
		c.JSON(response.ErrInternal.HTTPStatus(), response.ApiErrorResponse(response.ErrInternal.Code, response.ErrInternal.Message, nil))
		return
	}

	c.JSON(http.StatusOK, response.ApiSuccessResponse(gin.H{"uv": uv}))
}
