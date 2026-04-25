package upload

import (
	"ElainaBlog/internal/common/model"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	storage   Storage
	maxSize   int64           // 字节
	allowExts map[string]bool // 白名单
}

func NewController(storage Storage, maxSizeMB int) *Controller {
	return &Controller{
		storage: storage,
		maxSize: int64(maxSizeMB) << 20,
		allowExts: map[string]bool{
			".jpg":  true,
			".jpeg": true,
			".png":  true,
			".gif":  true,
			".webp": true,
		},
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
