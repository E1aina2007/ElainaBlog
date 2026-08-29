package comment

import (
	"ElainaBlog/internal/auth"
	"ElainaBlog/internal/response"
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AdminChecker 管理员权限检查接口。
type AdminChecker interface {
	CheckIsAdmin(ctx context.Context, userID int64) (bool, error)
}

type Controller struct {
	service      *Service
	adminChecker AdminChecker
}

func NewController(service *Service, adminChecker AdminChecker) *Controller {
	return &Controller{service: service, adminChecker: adminChecker}
}

func (ctl *Controller) GetList(c *gin.Context) {
	articleIDStr := c.Param("article_id")
	articleID, err := strconv.ParseInt(articleIDStr, 10, 64)
	if err != nil || articleID <= 0 {
		appErr := response.ErrInvalidParams.WithDetail("无效的文章 ID")
		c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	list, err := ctl.service.GetCommentList(c.Request.Context(), articleID)
	if err != nil {
		c.JSON(response.ErrInternal.HTTPStatus(), response.ApiErrorResponse(response.ErrInternal.Code, response.ErrInternal.Message, nil))
		return
	}

	c.JSON(http.StatusOK, response.ApiSuccessResponse(list))
}

// GetAllList 管理员获取所有评论列表
func (ctl *Controller) GetAllList(c *gin.Context) {
	list, err := ctl.service.GetAllCommentList(c.Request.Context())
	if err != nil {
		c.JSON(response.ErrInternal.HTTPStatus(), response.ApiErrorResponse(response.ErrInternal.Code, response.ErrInternal.Message, nil))
		return
	}

	c.JSON(http.StatusOK, response.ApiSuccessResponse(list))
}

func (ctl *Controller) CreateComment(c *gin.Context) {
	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := response.ErrInvalidParams.WithDetail("请求参数格式错误")
		c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	userID := c.GetInt64(auth.CtxUserIDKey)
	commentID, err := ctl.service.CreateComment(c.Request.Context(), &CreateCommentParams{
		ArticleID:        req.ArticleID,
		UserID:           userID,
		ReplyToUserID:    req.ReplyToUserID,
		ReplyToCommentID: req.ReplyToCommentID,
		Content:          req.Content,
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

	c.JSON(http.StatusOK, response.ApiSuccessResponse(gin.H{"id": commentID}))
}

func (ctl *Controller) DeleteComment(c *gin.Context) {
	var req DeleteCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := response.ErrInvalidParams.WithDetail("请求参数格式错误")
		c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	userID := c.GetInt64(auth.CtxUserIDKey)
	isAdmin, _ := ctl.adminChecker.CheckIsAdmin(c.Request.Context(), userID)

	err := ctl.service.DeleteComment(c.Request.Context(), &DeleteCommentParams{ID: req.ID}, userID, isAdmin)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidParams):
			appErr := response.ErrInvalidParams.WithDetail("无效的评论 ID")
			c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		case errors.Is(err, ErrCommentNotFound):
			appErr := response.ErrNotFound.WithDetail("评论不存在")
			c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		case errors.Is(err, ErrNoPermission):
			appErr := response.ErrForbidden.WithDetail("仅评论作者或管理员可删除")
			c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		default:
			c.JSON(response.ErrInternal.HTTPStatus(), response.ApiErrorResponse(response.ErrInternal.Code, response.ErrInternal.Message, nil))
		}
		return
	}

	c.JSON(http.StatusOK, response.ApiSuccessResponse(nil))
}
