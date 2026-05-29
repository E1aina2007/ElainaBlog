package user

import (
	"ElainaBlog/config"
	"ElainaBlog/internal/common"
	"ElainaBlog/internal/common/model"
	"fmt"
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

const (
	accessTokenCookie  = "access_token"
	refreshTokenCookie = "refresh_token"
)

func parseDuration(s string, defaultDur time.Duration) time.Duration {
	// Handle "Nd" (days) format which time.ParseDuration doesn't support
	if len(s) > 1 && s[len(s)-1] == 'd' {
		var days int
		if _, err := fmt.Sscanf(s, "%dd", &days); err == nil {
			return time.Duration(days) * 24 * time.Hour
		}
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return defaultDur
	}
	return d
}

func (ctl *Controller) setTokenCookies(c *gin.Context, accessToken, refreshToken string) {
	secure := !config.GlobalConfig.Dev
	accessTTL := parseDuration(config.GlobalConfig.Auth.AccessTokenExpiryTime, 2*time.Hour)
	refreshTTL := parseDuration(config.GlobalConfig.Auth.RefreshTokenExpiryTime, 7*24*time.Hour)
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     accessTokenCookie,
		Value:    accessToken,
		MaxAge:   int(accessTTL.Seconds()),
		Path:     "/",
		Secure:   secure,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     refreshTokenCookie,
		Value:    refreshToken,
		MaxAge:   int(refreshTTL.Seconds()),
		Path:     "/",
		Secure:   secure,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
}

func (ctl *Controller) clearTokenCookies(c *gin.Context) {
	secure := !config.GlobalConfig.Dev
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     accessTokenCookie,
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		Secure:   secure,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     refreshTokenCookie,
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		Secure:   secure,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
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

	ctl.setTokenCookies(c, result.AccessToken, result.RefreshToken)
	c.JSON(http.StatusOK, model.ApiSuccessResponse(gin.H{
		"user_id": result.UserID,
		"email":   result.Email,
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
	refreshTokenStr, err := c.Cookie(refreshTokenCookie)
	if err != nil || refreshTokenStr == "" {
		var req RefreshTokenRequest
		if bindErr := c.ShouldBindJSON(&req); bindErr == nil {
			refreshTokenStr = req.RefreshToken
		}
	}

	if refreshTokenStr == "" {
		appErr := model.ErrInvalidParams.WithDetail("缺少 refresh token")
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	claims, err := common.JwtAuth.ParseAndVerifyRefreshToken(refreshTokenStr)
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

	ctl.setTokenCookies(c, accessToken, refreshToken)
	c.JSON(http.StatusOK, model.ApiSuccessResponse(nil))
}

func (ctl *Controller) Logout(c *gin.Context) {
	// 吊销 access token
	if tokenStr, err := c.Cookie(accessTokenCookie); err == nil && tokenStr != "" {
		if claims, err := common.JwtAuth.ParseAndVerifyAccessToken(tokenStr); err == nil && claims.JTI != "" {
			ttl := time.Until(claims.ExpiresAt.Time)
			if ttl > 0 {
				common.BlacklistToken(claims.JTI, ttl)
			}
		}
	}

	// 吊销 refresh token
	if tokenStr, err := c.Cookie(refreshTokenCookie); err == nil && tokenStr != "" {
		if claims, err := common.JwtAuth.ParseAndVerifyRefreshToken(tokenStr); err == nil && claims.JTI != "" {
			ttl := time.Until(claims.ExpiresAt.Time)
			if ttl > 0 {
				common.BlacklistToken(claims.JTI, ttl)
			}
		}
	}

	ctl.clearTokenCookies(c)
	c.JSON(http.StatusOK, model.ApiSuccessResponse(nil))
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
