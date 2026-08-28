package authorprofile

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

// Get 公开接口：获取作者信息
func (ctl *Controller) Get(c *gin.Context) {
	profile, err := ctl.service.Get(c.Request.Context())
	if err != nil {
		c.JSON(response.ErrInternal.HTTPStatus(), response.ApiErrorResponse(response.ErrInternal.Code, response.ErrInternal.Message, nil))
		return
	}
	c.JSON(http.StatusOK, response.ApiSuccessResponse(profile))
}

// Update 管理员接口：更新作者信息
func (ctl *Controller) Update(c *gin.Context) {
	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(response.ErrInvalidParams.HTTPStatus(), response.ApiErrorResponse(response.ErrInvalidParams.Code, response.ErrInvalidParams.Message, nil))
		return
	}

	profile := &AuthorProfile{
		Nickname:             req.Nickname,
		Avatar:               req.Avatar,
		Background:           req.Background,
		Signature:            req.Signature,
		Location:             req.Location,
		Occupation:           req.Occupation,
		School:               req.School,
		Major:                req.Major,
		Email:                req.Email,
		Wechat:               req.Wechat,
		Bio:                  req.Bio,
		TechStackFrontend:    req.TechStackFrontend,
		TechStackBackend:     req.TechStackBackend,
		TechStackEngineering: req.TechStackEngineering,
		SocialGithub:         req.SocialGithub,
		SocialBilibili:       req.SocialBilibili,
	}

	if err := ctl.service.Update(c.Request.Context(), profile); err != nil {
		c.JSON(response.ErrInternal.HTTPStatus(), response.ApiErrorResponse(response.ErrInternal.Code, response.ErrInternal.Message, nil))
		return
	}

	c.JSON(http.StatusOK, response.ApiSuccessResponse(nil))
}
