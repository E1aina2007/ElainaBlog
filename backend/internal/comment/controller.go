package comment

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
	userService *user.Service // 校验为本人或管理员
}

func NewController(userService *user.Service) *Controller {
	repo := NewRepository(db.DBPool)
	service := NewService(repo)
	return &Controller{service: service, userService: userService}
}

type CreateCommentRequest struct {
	ArticleID int64  `json:"article_id"`
	Content   string `json:"content"`
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
		appErr := model.ErrInternal.WithDetail(err.Error())
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	c.JSON(http.StatusOK, model.ApiSuccessResponse(list))
}

func (ctl *Controller) CreateComment(c *gin.Context) {
	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := model.ErrInvalidParams.WithDetail(err.Error())
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	userID := c.GetInt64(common.CtxUserIDKey)
	commentID, err := ctl.service.CreateComment(&CreateCommentParams{
		ArticleID: req.ArticleID,
		UserID:    userID,
		Content:   req.Content,
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

	c.JSON(http.StatusOK, model.ApiSuccessResponse(gin.H{"id": commentID}))
}

func (ctl *Controller) DeleteComment(c *gin.Context) {
	var req DeleteCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := model.ErrInvalidParams.WithDetail(err.Error())
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	userID := c.GetInt64(common.CtxUserIDKey)

	// 检查评论是否存在，并校验本人或管理员
	comment, err := ctl.service.GetCommentByID(req.ID)
	if err != nil {
		switch err {
		case ErrCommentNotFound:
			appErr := model.ErrNotFound.WithDetail(err.Error())
			c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		default:
			appErr := model.ErrInternal.WithDetail(err.Error())
			c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		}
		return
	}

	if comment.UserID != userID {
		isAdmin, err := ctl.userService.CheckIsAdmin(userID)
		if err != nil || !isAdmin {
			appErr := model.ErrForbidden.WithDetail("仅评论作者或管理员可删除")
			c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
			return
		}
	}

	if err := ctl.service.DeleteComment(&DeleteCommentParams{ID: req.ID}); err != nil {
		appErr := model.ErrInternal.WithDetail(err.Error())
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	c.JSON(http.StatusOK, model.ApiSuccessResponse(nil))
}
