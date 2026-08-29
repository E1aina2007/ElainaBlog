package user

import (
	"ElainaBlog/internal/config"
	"ElainaBlog/internal/auth"
	cache "ElainaBlog/internal/middleware/redis"
	"ElainaBlog/internal/response"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Controller struct {
	service *Service
	rdb     *redis.Client // 可选，用于 token 黑名单
}

func NewController(service *Service, redis *redis.Client) *Controller {
	return &Controller{service: service, rdb: redis}
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
		appErr := response.ErrInternal.WithDetail("用户控制器未初始化")
		c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := response.ErrInvalidParams.WithDetail("请求参数格式错误")
		c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	userID, err := ctl.service.CreateUser(c.Request.Context(), CreateUserParams{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
		Code:     req.Code,
	})
	if err != nil {
		switch err {
		case ErrCodeExpired:
			c.JSON(response.ErrCodeExpired.HTTPStatus(), response.ApiErrorResponse(response.ErrCodeExpired.Code, response.ErrCodeExpired.Message, response.ErrCodeExpired))
		case ErrCodeMismatch:
			c.JSON(response.ErrCodeMismatch.HTTPStatus(), response.ApiErrorResponse(response.ErrCodeMismatch.Code, response.ErrCodeMismatch.Message, response.ErrCodeMismatch))
		case ErrEmailFormat, ErrEmailTooLong:
			c.JSON(response.ErrEmailFormat.HTTPStatus(), response.ApiErrorResponse(response.ErrEmailFormat.Code, response.ErrEmailFormat.Message, response.ErrEmailFormat))
		case ErrUsernameFormat:
			c.JSON(response.ErrUsernameFormat.HTTPStatus(), response.ApiErrorResponse(response.ErrUsernameFormat.Code, response.ErrUsernameFormat.Message, response.ErrUsernameFormat))
		case ErrPasswordLength, ErrPasswordChars, ErrPasswordNeedLetter, ErrPasswordNeedDigit:
			c.JSON(response.ErrPasswordLength.HTTPStatus(), response.ApiErrorResponse(response.ErrPasswordLength.Code, response.ErrPasswordLength.Message, response.ErrPasswordLength))
		case ErrUsernameExists:
			c.JSON(response.ErrUsernameExists.HTTPStatus(), response.ApiErrorResponse(response.ErrUsernameExists.Code, response.ErrUsernameExists.Message, response.ErrUsernameExists))
		case ErrEmailExists:
			c.JSON(response.ErrEmailExists.HTTPStatus(), response.ApiErrorResponse(response.ErrEmailExists.Code, response.ErrEmailExists.Message, response.ErrEmailExists))
		default:
			c.JSON(response.ErrInternal.HTTPStatus(), response.ApiErrorResponse(response.ErrInternal.Code, response.ErrInternal.Message, nil))
		}
		return
	}

	c.JSON(http.StatusOK, response.ApiSuccessResponse(gin.H{
		"user_id": userID,
	}))
}

// Login 登录接口：校验账号密码并签发 JWT（access/refresh）。
func (ctl *Controller) Login(c *gin.Context) {
	if ctl == nil || ctl.service == nil {
		appErr := response.ErrInternal.WithDetail("用户控制器未初始化")
		c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := response.ErrInvalidParams.WithDetail("请求参数格式错误")
		c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	result, err := ctl.service.Login(c.Request.Context(), LoginParams{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		switch err {
		case ErrInvalidLoginParams, ErrEmailFormat, ErrEmailTooLong:
			c.JSON(response.ErrEmailFormat.HTTPStatus(), response.ApiErrorResponse(response.ErrEmailFormat.Code, response.ErrEmailFormat.Message, response.ErrEmailFormat))
			return
		case ErrUserNotFound, ErrPasswordMismatch:
			// 登录失败统一返回邮箱或密码错误，不透露具体原因
			ctl.recordLoginFailure(c)
			c.JSON(response.ErrPasswordMismatch.HTTPStatus(), response.ApiErrorResponse(response.ErrPasswordMismatch.Code, response.ErrPasswordMismatch.Message, response.ErrPasswordMismatch))
			return
		default:
			c.JSON(response.ErrInternal.HTTPStatus(), response.ApiErrorResponse(response.ErrInternal.Code, response.ErrInternal.Message, nil))
			return
		}
	}

	if ctl.rdb != nil {
		cache.ResetLoginFailures(ctl.rdb, c.ClientIP())
	}
	ctl.setTokenCookies(c, result.AccessToken, result.RefreshToken)
	c.JSON(http.StatusOK, response.ApiSuccessResponse(gin.H{
		"user_id": result.UserID,
		"email":   result.Email,
	}))
}

// recordLoginFailure 记录登录失败，窗口内达到阈值时自动封禁该 IP（带 TTL 自动解封）
func (ctl *Controller) recordLoginFailure(c *gin.Context) {
	if ctl.rdb == nil {
		return
	}
	ip := c.ClientIP()
	fails, err := cache.RecordLoginFailure(ctl.rdb, ip)
	if err != nil {
		return // Redis 故障时静默跳过，不阻塞登录响应
	}
	if cache.AutoBanIfAbusive(ctl.rdb, ip, fails) {
		log.Printf("IP %s 登录失败达 %d 次，已自动封禁 %s", ip, fails, cache.LoginBanTTL)
	}
}

func (ctl *Controller) GetProfile(c *gin.Context) {
	userID := c.GetInt64(auth.CtxUserIDKey)
	u, err := ctl.service.GetByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(response.ErrInternal.HTTPStatus(), response.ApiErrorResponse(response.ErrInternal.Code, response.ErrInternal.Message, nil))
		return
	}

	c.JSON(http.StatusOK, response.ApiSuccessResponse(u.ToVO()))
}

func (ctl *Controller) GetList(c *gin.Context) {
	users, err := ctl.service.GetList(c.Request.Context())
	if err != nil {
		c.JSON(response.ErrInternal.HTTPStatus(), response.ApiErrorResponse(response.ErrInternal.Code, response.ErrInternal.Message, nil))
		return
	}

	var voList []*UserVO
	for _, u := range users {
		voList = append(voList, u.ToVO())
	}

	c.JSON(http.StatusOK, response.ApiSuccessResponse(voList))
}

func (ctl *Controller) UpdateProfile(c *gin.Context) {
	userID := c.GetInt64(auth.CtxUserIDKey)

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := response.ErrInvalidParams.WithDetail("请求参数格式错误")
		c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	err := ctl.service.UpdateProfile(c.Request.Context(), UpdateProfileParams{
		UserID:   userID,
		Username: req.Username,
		Email:    req.Email,
		Avatar:   req.Avatar,
	})
	if err != nil {
		switch err {
		case ErrUsernameExists:
			c.JSON(response.ErrUsernameExists.HTTPStatus(), response.ApiErrorResponse(response.ErrUsernameExists.Code, response.ErrUsernameExists.Message, response.ErrUsernameExists))
		case ErrEmailExists:
			c.JSON(response.ErrEmailExists.HTTPStatus(), response.ApiErrorResponse(response.ErrEmailExists.Code, response.ErrEmailExists.Message, response.ErrEmailExists))
		case ErrEmailFormat, ErrEmailTooLong:
			c.JSON(response.ErrEmailFormat.HTTPStatus(), response.ApiErrorResponse(response.ErrEmailFormat.Code, response.ErrEmailFormat.Message, response.ErrEmailFormat))
		case ErrUsernameFormat:
			c.JSON(response.ErrUsernameFormat.HTTPStatus(), response.ApiErrorResponse(response.ErrUsernameFormat.Code, response.ErrUsernameFormat.Message, response.ErrUsernameFormat))
		case ErrUserNotFound:
			c.JSON(response.ErrUserNotFound.HTTPStatus(), response.ApiErrorResponse(response.ErrUserNotFound.Code, response.ErrUserNotFound.Message, response.ErrUserNotFound))
		default:
			c.JSON(response.ErrInternal.HTTPStatus(), response.ApiErrorResponse(response.ErrInternal.Code, response.ErrInternal.Message, nil))
		}
		return
	}

	c.JSON(http.StatusOK, response.ApiSuccessResponse(nil))
}

func (ctl *Controller) UpdatePassword(c *gin.Context) {
	userID := c.GetInt64(auth.CtxUserIDKey)

	var req UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := response.ErrInvalidParams.WithDetail("请求参数格式错误")
		c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	err := ctl.service.UpdatePassword(c.Request.Context(), userID, req.OldPassword, req.NewPassword)
	if err != nil {
		switch err {
		case ErrPasswordMismatch:
			c.JSON(response.ErrPasswordMismatch.HTTPStatus(), response.ApiErrorResponse(response.ErrPasswordMismatch.Code, response.ErrPasswordMismatch.Message, response.ErrPasswordMismatch))
		case ErrSamePassword:
			c.JSON(http.StatusBadRequest, response.ApiErrorResponse(400010, "same password", nil))
		case ErrPasswordLength, ErrPasswordChars, ErrPasswordNeedLetter, ErrPasswordNeedDigit:
			c.JSON(response.ErrPasswordLength.HTTPStatus(), response.ApiErrorResponse(response.ErrPasswordLength.Code, response.ErrPasswordLength.Message, response.ErrPasswordLength))
		default:
			c.JSON(response.ErrInternal.HTTPStatus(), response.ApiErrorResponse(response.ErrInternal.Code, response.ErrInternal.Message, nil))
		}
		return
	}

	c.JSON(http.StatusOK, response.ApiSuccessResponse(nil))
}

func (ctl *Controller) DeleteUser(c *gin.Context) {
	// 验证管理员权限：操作者必须为管理员，由 service.DeleteUser 内部校验
	operatorID := c.GetInt64(auth.CtxUserIDKey)

	var req DeleteUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := response.ErrInvalidParams.WithDetail("请求参数格式错误")
		c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	err := ctl.service.DeleteUser(c.Request.Context(), operatorID, req.UserID)
	if err != nil {
		switch err {
		case ErrForbidden:
			c.JSON(response.ErrForbidden.HTTPStatus(), response.ApiErrorResponse(response.ErrForbidden.Code, response.ErrForbidden.Message, response.ErrForbidden))
		case ErrUserNotFound:
			c.JSON(response.ErrUserNotFound.HTTPStatus(), response.ApiErrorResponse(response.ErrUserNotFound.Code, response.ErrUserNotFound.Message, response.ErrUserNotFound))
		default:
			c.JSON(response.ErrInternal.HTTPStatus(), response.ApiErrorResponse(response.ErrInternal.Code, response.ErrInternal.Message, nil))
		}
		return
	}

	c.JSON(http.StatusOK, response.ApiSuccessResponse(nil))
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
		appErr := response.ErrInvalidParams.WithDetail("缺少 refresh token")
		c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	claims, err := auth.JwtAuth.ParseAndVerifyRefreshToken(refreshTokenStr)
	if err != nil {
		c.JSON(response.ErrRefreshTokenInvalid.HTTPStatus(), response.ApiErrorResponse(response.ErrRefreshTokenInvalid.Code, response.ErrRefreshTokenInvalid.Message, response.ErrRefreshTokenInvalid))
		return
	}

	// 吊销旧的 refresh token（一次性使用）
	if claims.JTI != "" {
		ttl := time.Until(claims.ExpiresAt.Time)
		if ttl > 0 {
			auth.BlacklistToken(ctl.rdb, claims.JTI, ttl)
		}
	}

	accessToken, err := auth.JwtAuth.GenerateAccessToken(claims.UserID)
	if err != nil {
		c.JSON(response.ErrInternal.HTTPStatus(), response.ApiErrorResponse(response.ErrInternal.Code, response.ErrInternal.Message, nil))
		return
	}

	refreshToken, err := auth.JwtAuth.GenerateRefreshToken(claims.UserID)
	if err != nil {
		c.JSON(response.ErrInternal.HTTPStatus(), response.ApiErrorResponse(response.ErrInternal.Code, response.ErrInternal.Message, nil))
		return
	}

	ctl.setTokenCookies(c, accessToken, refreshToken)
	c.JSON(http.StatusOK, response.ApiSuccessResponse(nil))
}

func (ctl *Controller) Logout(c *gin.Context) {
	// 吊销 access token
	if tokenStr, err := c.Cookie(accessTokenCookie); err == nil && tokenStr != "" {
		if claims, err := auth.JwtAuth.ParseAndVerifyAccessToken(tokenStr); err == nil && claims.JTI != "" {
			ttl := time.Until(claims.ExpiresAt.Time)
			if ttl > 0 {
				auth.BlacklistToken(ctl.rdb, claims.JTI, ttl)
			}
		}
	}

	// 吊销 refresh token
	if tokenStr, err := c.Cookie(refreshTokenCookie); err == nil && tokenStr != "" {
		if claims, err := auth.JwtAuth.ParseAndVerifyRefreshToken(tokenStr); err == nil && claims.JTI != "" {
			ttl := time.Until(claims.ExpiresAt.Time)
			if ttl > 0 {
				auth.BlacklistToken(ctl.rdb, claims.JTI, ttl)
			}
		}
	}

	ctl.clearTokenCookies(c)
	c.JSON(http.StatusOK, response.ApiSuccessResponse(nil))
}

func (ctl *Controller) SendCode(c *gin.Context) {
	var req SendCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := response.ErrInvalidParams.WithDetail("请求参数格式错误")
		c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	err := ctl.service.SendVerificationCode(c.Request.Context(), req.Email)
	if err != nil {
		switch err {
		case ErrEmailFormat, ErrEmailTooLong:
			c.JSON(response.ErrEmailFormat.HTTPStatus(), response.ApiErrorResponse(response.ErrEmailFormat.Code, response.ErrEmailFormat.Message, response.ErrEmailFormat))
		case ErrResendTooFrequent:
			c.JSON(response.ErrResendTooFrequent.HTTPStatus(), response.ApiErrorResponse(response.ErrResendTooFrequent.Code, response.ErrResendTooFrequent.Message, response.ErrResendTooFrequent))
		default:
			c.JSON(response.ErrInternal.HTTPStatus(), response.ApiErrorResponse(response.ErrInternal.Code, response.ErrInternal.Message, nil))
		}
		return
	}

	c.JSON(http.StatusOK, response.ApiSuccessResponse(nil))
}

func (ctl *Controller) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := response.ErrInvalidParams.WithDetail("请求参数格式错误")
		c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	err := ctl.service.ResetPassword(c.Request.Context(), req.Email, req.Code, req.NewPassword)
	if err != nil {
		switch err {
		case ErrEmailFormat, ErrEmailTooLong:
			c.JSON(response.ErrEmailFormat.HTTPStatus(), response.ApiErrorResponse(response.ErrEmailFormat.Code, response.ErrEmailFormat.Message, response.ErrEmailFormat))
		case ErrPasswordLength, ErrPasswordChars, ErrPasswordNeedLetter, ErrPasswordNeedDigit:
			c.JSON(response.ErrPasswordLength.HTTPStatus(), response.ApiErrorResponse(response.ErrPasswordLength.Code, response.ErrPasswordLength.Message, response.ErrPasswordLength))
		case ErrCodeExpired:
			c.JSON(response.ErrCodeExpired.HTTPStatus(), response.ApiErrorResponse(response.ErrCodeExpired.Code, response.ErrCodeExpired.Message, response.ErrCodeExpired))
		case ErrCodeMismatch:
			c.JSON(response.ErrCodeMismatch.HTTPStatus(), response.ApiErrorResponse(response.ErrCodeMismatch.Code, response.ErrCodeMismatch.Message, response.ErrCodeMismatch))
		case ErrUserNotFound:
			c.JSON(response.ErrUserNotFound.HTTPStatus(), response.ApiErrorResponse(response.ErrUserNotFound.Code, response.ErrUserNotFound.Message, response.ErrUserNotFound))
		default:
			c.JSON(response.ErrInternal.HTTPStatus(), response.ApiErrorResponse(response.ErrInternal.Code, response.ErrInternal.Message, nil))
		}
		return
	}

	c.JSON(http.StatusOK, response.ApiSuccessResponse(nil))
}
