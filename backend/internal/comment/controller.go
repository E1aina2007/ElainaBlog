package comment

import (
	"ElainaBlog/internal/common"
	"ElainaBlog/internal/common/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AdminChecker 管理员权限检查接口。
type AdminChecker interface {
	CheckIsAdmin(userID int64) (bool, error)
}

type Controller struct {
	service      *Service
	adminChecker AdminChecker
}

func NewController(service *Service, adminChecker AdminChecker) *Controller {
	return &Controller{service: service, adminChecker: adminChecker}
}

type CreateCommentRequest struct {
	ArticleID        int64  `json:"article_id"`
	ReplyToUserID    *int64 `json:"reply_to_user_id"`
	ReplyToCommentID *int64 `json:"reply_to_comment_id"`
	Content          string `json:"content"`
}

type DeleteCommentRequest struct {
	ID int64 `json:"id"`
}

func (ctl *Controller) GetList(c *gin.Context) {
	articleIDStr := c.Param("article_id")
	articleID, err := strconv.ParseInt(articleIDStr, 10, 64)
	if err != nil || articleID <= 0 {
		appErr := model.ErrInvalidParams.WithDetail("无效的文章 ID")
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	list, err := ctl.service.GetCommentList(articleID)
	if err != nil {
		c.JSON(model.ErrInternal.HTTPStatus(), model.ApiErrorResponse(model.ErrInternal.Code, model.ErrInternal.Message, nil))
		return
	}

	c.JSON(http.StatusOK, model.ApiSuccessResponse(list))
}

// GetAllList 管理员获取所有评论列表
func (ctl *Controller) GetAllList(c *gin.Context) {
	list, err := ctl.service.GetAllCommentList()
	if err != nil {
		c.JSON(model.ErrInternal.HTTPStatus(), model.ApiErrorResponse(model.ErrInternal.Code, model.ErrInternal.Message, nil))
		return
	}

	c.JSON(http.StatusOK, model.ApiSuccessResponse(list))
}

func (ctl *Controller) CreateComment(c *gin.Context) {
	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := model.ErrInvalidParams.WithDetail("请求参数格式错误")
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	userID := c.GetInt64(common.CtxUserIDKey)
	commentID, err := ctl.service.CreateComment(&CreateCommentParams{
		ArticleID:        req.ArticleID,
		UserID:           userID,
		ReplyToUserID:    req.ReplyToUserID,
		ReplyToCommentID: req.ReplyToCommentID,
		Content:          req.Content,
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

	c.JSON(http.StatusOK, model.ApiSuccessResponse(gin.H{"id": commentID}))
}

func (ctl *Controller) DeleteComment(c *gin.Context) {
	var req DeleteCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := model.ErrInvalidParams.WithDetail("请求参数格式错误")
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	userID := c.GetInt64(common.CtxUserIDKey)

	// 检查评论是否存在，并校验本人或管理员
	comment, err := ctl.service.GetCommentByID(req.ID)
	if err != nil {
		switch err {
		case ErrCommentNotFound:
			appErr := model.ErrNotFound.WithDetail("评论不存在")
			c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		default:
			c.JSON(model.ErrInternal.HTTPStatus(), model.ApiErrorResponse(model.ErrInternal.Code, model.ErrInternal.Message, nil))
		}
		return
	}

	if comment.UserID != userID {
		isAdmin, err := ctl.adminChecker.CheckIsAdmin(userID)
		if err != nil {
			c.JSON(model.ErrInternal.HTTPStatus(), model.ApiErrorResponse(model.ErrInternal.Code, model.ErrInternal.Message, nil))
			return
		}
		if !isAdmin {
			appErr := model.ErrForbidden.WithDetail("仅评论作者或管理员可删除")
			c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
			return
		}
	}

	if err := ctl.service.DeleteComment(&DeleteCommentParams{ID: req.ID}); err != nil {
		c.JSON(model.ErrInternal.HTTPStatus(), model.ApiErrorResponse(model.ErrInternal.Code, model.ErrInternal.Message, nil))
		return
	}

	c.JSON(http.StatusOK, model.ApiSuccessResponse(nil))
}
