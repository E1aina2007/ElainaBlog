package upload

import (
	"ElainaBlog/internal/auth"
	"ElainaBlog/internal/response"
	"ElainaBlog/internal/mail"
	"context"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	storage       Storage
	maxSize       int64 // 字节
	avatarStorage Storage
	avatarMaxSize int64           // 字节
	allowExts     map[string]bool // 扩展名白名单
	allowMIMEs    map[string]bool // MIME 类型白名单
	userService   UserService
}

type UserService interface {
	GetUserEmailByID(ctx context.Context, id int64) (string, error)
}

func NewController(storage Storage, maxSizeMB int, avatarStorage Storage, avatarMaxSizeMB int, userService UserService) *Controller {
	return &Controller{
		storage:       storage,
		maxSize:       int64(maxSizeMB) << 20,
		avatarStorage: avatarStorage,
		avatarMaxSize: int64(avatarMaxSizeMB) << 20,
		allowExts: map[string]bool{
			".jpg":  true,
			".jpeg": true,
			".png":  true,
			".webp": true,
		},
		allowMIMEs: map[string]bool{
			"image/jpeg": true,
			"image/png":  true,
			"image/webp": true,
		},
		userService: userService,
	}
}

// validateMIME 读取文件头部字节检测真实 MIME 类型，校验后重置读取位置。
func (ctl *Controller) validateMIME(file io.ReadSeeker) (string, error) {
	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return "", err
	}
	// 重置读取位置，确保后续存储能读到完整文件
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", err
	}
	mimeType := http.DetectContentType(buf[:n])
	// 去除参数部分，如 "image/jpeg; charset=utf-8" -> "image/jpeg"
	if idx := strings.Index(mimeType, ";"); idx != -1 {
		mimeType = strings.TrimSpace(mimeType[:idx])
	}
	return mimeType, nil
}

func (ctl *Controller) Upload(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		appErr := response.ErrInvalidParams.WithDetail("缺少上传文件")
		c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	// 校验文件大小
	if fileHeader.Size > ctl.maxSize {
		appErr := response.ErrInvalidParams.WithDetail("文件大小超出限制")
		c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	// 校验扩展名
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !ctl.allowExts[ext] {
		appErr := response.ErrInvalidParams.WithDetail("不支持的文件类型")
		c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	// 打开文件
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(response.ErrInternal.HTTPStatus(), response.ApiErrorResponse(response.ErrInternal.Code, response.ErrInternal.Message, nil))
		return
	}
	defer file.Close()

	// 校验 MIME 类型（魔数检测）
	mimeType, err := ctl.validateMIME(file)
	if err != nil || !ctl.allowMIMEs[mimeType] {
		appErr := response.ErrInvalidParams.WithDetail("不支持的文件内容类型")
		c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	// 保存
	url, err := ctl.storage.Save(file, fileHeader.Filename)
	if err != nil {
		c.JSON(response.ErrInternal.HTTPStatus(), response.ApiErrorResponse(response.ErrInternal.Code, response.ErrInternal.Message, nil))
		return
	}

	c.JSON(http.StatusOK, response.ApiSuccessResponse(gin.H{"url": url}))
}

func (ctl *Controller) UploadAvatar(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		appErr := response.ErrInvalidParams.WithDetail("缺少上传文件")
		c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	// 校验文件大小
	if fileHeader.Size > ctl.avatarMaxSize {
		appErr := response.ErrInvalidParams.WithDetail("头像文件大小超出限制")
		c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	// 校验扩展名
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !ctl.allowExts[ext] {
		appErr := response.ErrInvalidParams.WithDetail("不支持的文件类型")
		c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	// 从 JWT 中获取用户 ID
	userID, exists := c.Get(auth.CtxUserIDKey)
	if !exists {
		appErr := response.ErrUnauthorized.WithDetail("未登录")
		c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	// 获取用户邮箱
	email, err := ctl.userService.GetUserEmailByID(c.Request.Context(), userID.(int64))
	if err != nil {
		appErr := response.ErrInternal.WithDetail("获取用户信息失败")
		c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	// 打开文件
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(response.ErrInternal.HTTPStatus(), response.ApiErrorResponse(response.ErrInternal.Code, response.ErrInternal.Message, nil))
		return
	}
	defer file.Close()

	// 校验 MIME 类型（魔数检测）
	mimeType, err := ctl.validateMIME(file)
	if err != nil || !ctl.allowMIMEs[mimeType] {
		appErr := response.ErrInvalidParams.WithDetail("不支持的文件内容类型")
		c.JSON(appErr.HTTPStatus(), response.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	// 构造自定义文件名：邮箱哈希（同一邮箱始终生成相同哈希，新头像会覆盖旧头像）
	customName := mail.EmailToAvatarHash(email)

	// 保存到头像专用目录
	url, err := ctl.avatarStorage.SaveAs(file, fileHeader.Filename, customName)
	if err != nil {
		c.JSON(response.ErrInternal.HTTPStatus(), response.ApiErrorResponse(response.ErrInternal.Code, response.ErrInternal.Message, nil))
		return
	}

	c.JSON(http.StatusOK, response.ApiSuccessResponse(gin.H{"url": url}))
}
