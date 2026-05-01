package user

import (
	"ElainaBlog/internal/common"
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

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateProfileRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
}

type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type DeleteUserRequest struct {
	UserID int64 `json:"user_id"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Code     string `json:"code"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type SendCodeRequest struct {
	Email string `json:"email"`
}

// Register 注册接口：创建新用户。
func (ctl *Controller) Register(c *gin.Context) {
	if ctl == nil || ctl.service == nil {
		appErr := model.ErrInternal.WithDetail("user controller not initialized")
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := model.ErrInvalidParams.WithDetail(err.Error())
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	userID, err := ctl.service.CreateUser(CreateUserParams{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
		Code:     req.Code,
	})
	if err != nil {
		switch err {
		case ErrInvalidParams, ErrCodeExpired, ErrCodeMismatch,
			ErrEmailFormat, ErrEmailTooLong, ErrUsernameFormat,
			ErrPasswordLength, ErrPasswordChars, ErrPasswordNeedLetter, ErrPasswordNeedDigit:
			appErr := model.ErrInvalidParams.WithDetail(err.Error())
			c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		case ErrUsernameExists, ErrEmailExists:
			appErr := model.ErrConflict.WithDetail(err.Error())
			c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		default:
			appErr := model.ErrInternal.WithDetail(err.Error())
			c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		}
		return
	}

	c.JSON(http.StatusOK, model.ApiSuccessResponse(gin.H{
		"user_id": userID,
	}))
}

// Login 登录接口：校验账号密码并签发 JWT（access/refresh）。
func (ctl *Controller) Login(c *gin.Context) {
	if ctl == nil || ctl.service == nil {
		appErr := model.ErrInternal.WithDetail("user controller not initialized")
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := model.ErrInvalidParams.WithDetail(err.Error())
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	result, err := ctl.service.Login(LoginParams{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		switch err {
		case ErrInvalidLoginParams, ErrEmailFormat, ErrEmailTooLong:
			appErr := model.ErrInvalidParams.WithDetail(err.Error())
			c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
			return
		case ErrUserNotFound, ErrPasswordMismatch:
			appErr := model.ErrUnauthorized.WithDetail("邮箱或密码错误")
			c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
			return
		default:
			appErr := model.ErrInternal.WithDetail(err.Error())
			c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
			return
		}
	}

	c.JSON(http.StatusOK, model.ApiSuccessResponse(gin.H{
		"user_id":       result.UserID,
		"email":         result.Email,
		"access_token":  result.AccessToken,
		"refresh_token": result.RefreshToken,
	}))
}

func (ctl *Controller) GetProfile(c *gin.Context) {
	userID := c.GetInt64(common.CtxUserIDKey)
	u, err := ctl.service.GetByID(userID)
	if err != nil {
		appErr := model.ErrInternal.WithDetail(err.Error())
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	c.JSON(http.StatusOK, model.ApiSuccessResponse(u.ToVO()))
}

func (ctl *Controller) GetList(c *gin.Context) {
	// 验证管理员权限：从 JWT 中取出当前用户 ID，校验是否为管理员
	userID := c.GetInt64(common.CtxUserIDKey)
	isAdmin, err := ctl.service.CheckIsAdmin(userID)
	if err != nil {
		appErr := model.ErrInternal.WithDetail(err.Error())
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}
	if !isAdmin {
		appErr := model.ErrForbidden.WithDetail("仅管理员可访问")
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	users, err := ctl.service.GetList()
	if err != nil {
		appErr := model.ErrInternal.WithDetail(err.Error())
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	var voList []*UserVO
	for _, u := range users {
		voList = append(voList, u.ToVO())
	}

	c.JSON(http.StatusOK, model.ApiSuccessResponse(voList))
}

func (ctl *Controller) UpdateProfile(c *gin.Context) {
	userID := c.GetInt64(common.CtxUserIDKey)

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := model.ErrInvalidParams.WithDetail(err.Error())
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	err := ctl.service.UpdateProfile(UpdateProfileParams{
		UserID:   userID,
		Username: req.Username,
		Email:    req.Email,
		Avatar:   req.Avatar,
	})
	if err != nil {
		switch err {
		case ErrUsernameExists, ErrEmailExists, ErrInvalidParams,
			ErrUsernameFormat, ErrEmailFormat, ErrEmailTooLong:
			appErr := model.ErrInvalidParams.WithDetail(err.Error())
			c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		case ErrUserNotFound:
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

func (ctl *Controller) UpdatePassword(c *gin.Context) {
	userID := c.GetInt64(common.CtxUserIDKey)

	var req UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := model.ErrInvalidParams.WithDetail(err.Error())
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	err := ctl.service.UpdatePassword(userID, req.OldPassword, req.NewPassword)
	if err != nil {
		switch err {
		case ErrPasswordMismatch:
			appErr := model.ErrUnauthorized.WithDetail(err.Error())
			c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		case ErrSamePassword, ErrInvalidParams,
			ErrPasswordLength, ErrPasswordChars, ErrPasswordNeedLetter, ErrPasswordNeedDigit:
			appErr := model.ErrInvalidParams.WithDetail(err.Error())
			c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		default:
			appErr := model.ErrInternal.WithDetail(err.Error())
			c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		}
		return
	}

	c.JSON(http.StatusOK, model.ApiSuccessResponse(nil))
}

func (ctl *Controller) DeleteUser(c *gin.Context) {
	// 验证管理员权限：操作者必须为管理员，由 service.DeleteUser 内部校验
	operatorID := c.GetInt64(common.CtxUserIDKey)

	var req DeleteUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := model.ErrInvalidParams.WithDetail(err.Error())
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	err := ctl.service.DeleteUser(operatorID, req.UserID)
	if err != nil {
		switch err {
		case ErrForbidden:
			appErr := model.ErrForbidden.WithDetail(err.Error())
			c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		case ErrUserNotFound:
			appErr := model.ErrNotFound.WithDetail(err.Error())
			c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		case ErrInvalidParams:
			appErr := model.ErrInvalidParams.WithDetail(err.Error())
			c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		default:
			appErr := model.ErrInternal.WithDetail(err.Error())
			c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		}
		return
	}

	c.JSON(http.StatusOK, model.ApiSuccessResponse(nil))
}

func (ctl *Controller) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := model.ErrInvalidParams.WithDetail(err.Error())
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	claims, err := common.JwtAuth.ParseAndVerifyRefreshToken(req.RefreshToken)
	if err != nil {
		appErr := model.ErrUnauthorized.WithDetail("refresh token 无效或已过期")
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	accessToken, err := common.JwtAuth.GenerateAccessToken(claims.UserID)
	if err != nil {
		appErr := model.ErrInternal.WithDetail(err.Error())
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	c.JSON(http.StatusOK, model.ApiSuccessResponse(gin.H{
		"access_token": accessToken,
	}))
}

func (ctl *Controller) SendCode(c *gin.Context) {
	var req SendCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := model.ErrInvalidParams.WithDetail(err.Error())
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	err := ctl.service.SendVerificationCode(req.Email)
	if err != nil {
		switch err {
		case ErrInvalidParams, ErrEmailFormat, ErrEmailTooLong:
			appErr := model.ErrInvalidParams.WithDetail(err.Error())
			c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		case ErrResendTooFrequent:
			appErr := model.ErrTooManyRequests.WithDetail(err.Error())
			c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		default:
			appErr := model.ErrInternal.WithDetail(err.Error())
			c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		}
		return
	}

	c.JSON(http.StatusOK, model.ApiSuccessResponse(nil))
}
