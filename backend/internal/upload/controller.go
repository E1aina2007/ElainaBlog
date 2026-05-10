package upload

import (
	"ElainaBlog/internal/common"
	"ElainaBlog/internal/common/model"
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
	allowExts     map[string]bool // 白名单
	userService   UserService
}

type UserService interface {
	GetUserEmailByID(id int64) (string, error)
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
			".gif":  true,
			".webp": true,
		},
		userService: userService,
	}
}

func (ctl *Controller) Upload(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		appErr := model.ErrInvalidParams.WithDetail("缺少上传文件")
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	// 校验文件大小
	if fileHeader.Size > ctl.maxSize {
		appErr := model.ErrInvalidParams.WithDetail("文件大小超出限制")
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	// 校验扩展名
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !ctl.allowExts[ext] {
		appErr := model.ErrInvalidParams.WithDetail("不支持的文件类型")
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	// 打开文件
	file, err := fileHeader.Open()
	if err != nil {
		appErr := model.ErrInternal.WithDetail(err.Error())
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}
	defer file.Close()

	// 保存
	url, err := ctl.storage.Save(file, fileHeader.Filename)
	if err != nil {
		appErr := model.ErrInternal.WithDetail(err.Error())
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	c.JSON(http.StatusOK, model.ApiSuccessResponse(gin.H{"url": url}))
}

func (ctl *Controller) UploadAvatar(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		appErr := model.ErrInvalidParams.WithDetail("缺少上传文件")
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	// 校验文件大小
	if fileHeader.Size > ctl.avatarMaxSize {
		appErr := model.ErrInvalidParams.WithDetail("头像文件大小超出限制")
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	// 校验扩展名
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !ctl.allowExts[ext] {
		appErr := model.ErrInvalidParams.WithDetail("不支持的文件类型")
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	// 从 JWT 中获取用户 ID
	userID, exists := c.Get(common.CtxUserIDKey)
	if !exists {
		appErr := model.ErrUnauthorized.WithDetail("未登录")
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	// 获取用户邮箱
	email, err := ctl.userService.GetUserEmailByID(userID.(int64))
	if err != nil {
		appErr := model.ErrInternal.WithDetail("获取用户信息失败")
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	// 打开文件
	file, err := fileHeader.Open()
	if err != nil {
		appErr := model.ErrInternal.WithDetail(err.Error())
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}
	defer file.Close()

	// 构造自定义文件名：邮箱（新头像会覆盖旧头像）
	customName := email

	// 保存到头像专用目录
	url, err := ctl.avatarStorage.SaveAs(file, fileHeader.Filename, customName)
	if err != nil {
		appErr := model.ErrInternal.WithDetail(err.Error())
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	c.JSON(http.StatusOK, model.ApiSuccessResponse(gin.H{"url": url}))
}
