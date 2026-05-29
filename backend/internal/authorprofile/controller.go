package authorprofile

import (
	"ElainaBlog/internal/common/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service *Service
}

func NewController(service *Service) *Controller {
	return &Controller{service: service}
}

type UpdateProfileRequest struct {
	Nickname             string `json:"nickname"`
	Avatar               string `json:"avatar"`
	Background           string `json:"background"`
	Signature            string `json:"signature"`
	Location             string `json:"location"`
	Occupation           string `json:"occupation"`
	School               string `json:"school"`
	Major                string `json:"major"`
	Email                string `json:"email"`
	Wechat               string `json:"wechat"`
	Bio                  string `json:"bio"`
	TechStackFrontend    string `json:"tech_stack_frontend"`
	TechStackBackend     string `json:"tech_stack_backend"`
	TechStackEngineering string `json:"tech_stack_engineering"`
	SocialGithub         string `json:"social_github"`
	SocialBilibili       string `json:"social_bilibili"`
}

// Get 公开接口：获取作者信息
func (ctl *Controller) Get(c *gin.Context) {
	profile, err := ctl.service.Get()
	if err != nil {
		c.JSON(model.ErrInternal.HTTPStatus(), model.ApiErrorResponse(model.ErrInternal.Code, model.ErrInternal.Message, nil))
		return
	}
	c.JSON(http.StatusOK, model.ApiSuccessResponse(profile))
}

// Update 管理员接口：更新作者信息
func (ctl *Controller) Update(c *gin.Context) {
	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(model.ErrInvalidParams.HTTPStatus(), model.ApiErrorResponse(model.ErrInvalidParams.Code, model.ErrInvalidParams.Message, nil))
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

	if err := ctl.service.Update(profile); err != nil {
		c.JSON(model.ErrInternal.HTTPStatus(), model.ApiErrorResponse(model.ErrInternal.Code, model.ErrInternal.Message, nil))
		return
	}

	c.JSON(http.StatusOK, model.ApiSuccessResponse(nil))
}
