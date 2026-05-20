package user

import (
	"ElainaBlog/internal/common"
	"ElainaBlog/internal/common/model"
	"net/http"
	"time"

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
		appErr := model.ErrInternal.WithDetail("用户控制器未初始化")
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := model.ErrInvalidParams.WithDetail("请求参数格式错误")
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
		case ErrCodeExpired:
			c.JSON(model.ErrCodeExpired.HTTPStatus(), model.ApiErrorResponse(model.ErrCodeExpired.Code, model.ErrCodeExpired.Message, model.ErrCodeExpired))
		case ErrCodeMismatch:
			c.JSON(model.ErrCodeMismatch.HTTPStatus(), model.ApiErrorResponse(model.ErrCodeMismatch.Code, model.ErrCodeMismatch.Message, model.ErrCodeMismatch))
		case ErrEmailFormat, ErrEmailTooLong:
			c.JSON(model.ErrEmailFormat.HTTPStatus(), model.ApiErrorResponse(model.ErrEmailFormat.Code, model.ErrEmailFormat.Message, model.ErrEmailFormat))
		case ErrUsernameFormat:
			c.JSON(model.ErrUsernameFormat.HTTPStatus(), model.ApiErrorResponse(model.ErrUsernameFormat.Code, model.ErrUsernameFormat.Message, model.ErrUsernameFormat))
		case ErrPasswordLength, ErrPasswordChars, ErrPasswordNeedLetter, ErrPasswordNeedDigit:
			c.JSON(model.ErrPasswordLength.HTTPStatus(), model.ApiErrorResponse(model.ErrPasswordLength.Code, model.ErrPasswordLength.Message, model.ErrPasswordLength))
		case ErrUsernameExists:
			c.JSON(model.ErrUsernameExists.HTTPStatus(), model.ApiErrorResponse(model.ErrUsernameExists.Code, model.ErrUsernameExists.Message, model.ErrUsernameExists))
		case ErrEmailExists:
			c.JSON(model.ErrEmailExists.HTTPStatus(), model.ApiErrorResponse(model.ErrEmailExists.Code, model.ErrEmailExists.Message, model.ErrEmailExists))
		default:
			c.JSON(model.ErrInternal.HTTPStatus(), model.ApiErrorResponse(model.ErrInternal.Code, model.ErrInternal.Message, nil))
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
		appErr := model.ErrInternal.WithDetail("用户控制器未初始化")
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := model.ErrInvalidParams.WithDetail("请求参数格式错误")
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
			c.JSON(model.ErrEmailFormat.HTTPStatus(), model.ApiErrorResponse(model.ErrEmailFormat.Code, model.ErrEmailFormat.Message, model.ErrEmailFormat))
			return
		case ErrUserNotFound, ErrPasswordMismatch:
			// 登录失败统一返回邮箱或密码错误，不透露具体原因
			c.JSON(model.ErrPasswordMismatch.HTTPStatus(), model.ApiErrorResponse(model.ErrPasswordMismatch.Code, model.ErrPasswordMismatch.Message, model.ErrPasswordMismatch))
			return
		default:
			c.JSON(model.ErrInternal.HTTPStatus(), model.ApiErrorResponse(model.ErrInternal.Code, model.ErrInternal.Message, nil))
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
		c.JSON(model.ErrInternal.HTTPStatus(), model.ApiErrorResponse(model.ErrInternal.Code, model.ErrInternal.Message, nil))
		return
	}

	c.JSON(http.StatusOK, model.ApiSuccessResponse(u.ToVO()))
}

func (ctl *Controller) GetList(c *gin.Context) {
	// 验证管理员权限：从 JWT 中取出当前用户 ID，校验是否为管理员
	userID := c.GetInt64(common.CtxUserIDKey)
	isAdmin, err := ctl.service.CheckIsAdmin(userID)
	if err != nil {
		c.JSON(model.ErrInternal.HTTPStatus(), model.ApiErrorResponse(model.ErrInternal.Code, model.ErrInternal.Message, nil))
		return
	}
	if !isAdmin {
		appErr := model.ErrForbidden.WithDetail("仅管理员可访问")
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	users, err := ctl.service.GetList()
	if err != nil {
		c.JSON(model.ErrInternal.HTTPStatus(), model.ApiErrorResponse(model.ErrInternal.Code, model.ErrInternal.Message, nil))
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
		appErr := model.ErrInvalidParams.WithDetail("请求参数格式错误")
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
		case ErrUsernameExists:
			c.JSON(model.ErrUsernameExists.HTTPStatus(), model.ApiErrorResponse(model.ErrUsernameExists.Code, model.ErrUsernameExists.Message, model.ErrUsernameExists))
		case ErrEmailExists:
			c.JSON(model.ErrEmailExists.HTTPStatus(), model.ApiErrorResponse(model.ErrEmailExists.Code, model.ErrEmailExists.Message, model.ErrEmailExists))
		case ErrEmailFormat, ErrEmailTooLong:
			c.JSON(model.ErrEmailFormat.HTTPStatus(), model.ApiErrorResponse(model.ErrEmailFormat.Code, model.ErrEmailFormat.Message, model.ErrEmailFormat))
		case ErrUsernameFormat:
			c.JSON(model.ErrUsernameFormat.HTTPStatus(), model.ApiErrorResponse(model.ErrUsernameFormat.Code, model.ErrUsernameFormat.Message, model.ErrUsernameFormat))
		case ErrUserNotFound:
			c.JSON(model.ErrUserNotFound.HTTPStatus(), model.ApiErrorResponse(model.ErrUserNotFound.Code, model.ErrUserNotFound.Message, model.ErrUserNotFound))
		default:
			c.JSON(model.ErrInternal.HTTPStatus(), model.ApiErrorResponse(model.ErrInternal.Code, model.ErrInternal.Message, nil))
		}
		return
	}

	c.JSON(http.StatusOK, model.ApiSuccessResponse(nil))
}

func (ctl *Controller) UpdatePassword(c *gin.Context) {
	userID := c.GetInt64(common.CtxUserIDKey)

	var req UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := model.ErrInvalidParams.WithDetail("请求参数格式错误")
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	err := ctl.service.UpdatePassword(userID, req.OldPassword, req.NewPassword)
	if err != nil {
		switch err {
		case ErrPasswordMismatch:
			c.JSON(model.ErrPasswordMismatch.HTTPStatus(), model.ApiErrorResponse(model.ErrPasswordMismatch.Code, model.ErrPasswordMismatch.Message, model.ErrPasswordMismatch))
		case ErrSamePassword:
			c.JSON(http.StatusBadRequest, model.ApiErrorResponse(400010, "same password", nil))
		case ErrPasswordLength, ErrPasswordChars, ErrPasswordNeedLetter, ErrPasswordNeedDigit:
			c.JSON(model.ErrPasswordLength.HTTPStatus(), model.ApiErrorResponse(model.ErrPasswordLength.Code, model.ErrPasswordLength.Message, model.ErrPasswordLength))
		default:
			c.JSON(model.ErrInternal.HTTPStatus(), model.ApiErrorResponse(model.ErrInternal.Code, model.ErrInternal.Message, nil))
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
		appErr := model.ErrInvalidParams.WithDetail("请求参数格式错误")
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	err := ctl.service.DeleteUser(operatorID, req.UserID)
	if err != nil {
		switch err {
		case ErrForbidden:
			c.JSON(model.ErrForbidden.HTTPStatus(), model.ApiErrorResponse(model.ErrForbidden.Code, model.ErrForbidden.Message, model.ErrForbidden))
		case ErrUserNotFound:
			c.JSON(model.ErrUserNotFound.HTTPStatus(), model.ApiErrorResponse(model.ErrUserNotFound.Code, model.ErrUserNotFound.Message, model.ErrUserNotFound))
		default:
			c.JSON(model.ErrInternal.HTTPStatus(), model.ApiErrorResponse(model.ErrInternal.Code, model.ErrInternal.Message, nil))
		}
		return
	}

	c.JSON(http.StatusOK, model.ApiSuccessResponse(nil))
}

func (ctl *Controller) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := model.ErrInvalidParams.WithDetail("请求参数格式错误")
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	claims, err := common.JwtAuth.ParseAndVerifyRefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(model.ErrRefreshTokenInvalid.HTTPStatus(), model.ApiErrorResponse(model.ErrRefreshTokenInvalid.Code, model.ErrRefreshTokenInvalid.Message, model.ErrRefreshTokenInvalid))
		return
	}

	// 吊销旧的 refresh token（一次性使用）
	if claims.JTI != "" {
		ttl := time.Until(claims.ExpiresAt.Time)
		if ttl > 0 {
			common.BlacklistToken(claims.JTI, ttl)
		}
	}

	accessToken, err := common.JwtAuth.GenerateAccessToken(claims.UserID)
	if err != nil {
		c.JSON(model.ErrInternal.HTTPStatus(), model.ApiErrorResponse(model.ErrInternal.Code, model.ErrInternal.Message, nil))
		return
	}

	refreshToken, err := common.JwtAuth.GenerateRefreshToken(claims.UserID)
	if err != nil {
		c.JSON(model.ErrInternal.HTTPStatus(), model.ApiErrorResponse(model.ErrInternal.Code, model.ErrInternal.Message, nil))
		return
	}

	c.JSON(http.StatusOK, model.ApiSuccessResponse(gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}))
}

func (ctl *Controller) SendCode(c *gin.Context) {
	var req SendCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := model.ErrInvalidParams.WithDetail("请求参数格式错误")
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	err := ctl.service.SendVerificationCode(req.Email)
	if err != nil {
		switch err {
		case ErrEmailFormat, ErrEmailTooLong:
			c.JSON(model.ErrEmailFormat.HTTPStatus(), model.ApiErrorResponse(model.ErrEmailFormat.Code, model.ErrEmailFormat.Message, model.ErrEmailFormat))
		case ErrResendTooFrequent:
			c.JSON(model.ErrResendTooFrequent.HTTPStatus(), model.ApiErrorResponse(model.ErrResendTooFrequent.Code, model.ErrResendTooFrequent.Message, model.ErrResendTooFrequent))
		default:
			c.JSON(model.ErrInternal.HTTPStatus(), model.ApiErrorResponse(model.ErrInternal.Code, model.ErrInternal.Message, nil))
		}
		return
	}

	c.JSON(http.StatusOK, model.ApiSuccessResponse(nil))
}
