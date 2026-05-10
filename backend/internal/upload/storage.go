package upload

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

type Storage interface {
	Save(file multipart.File, filename string) (url string, err error)
	SaveAs(file multipart.File, filename string, customName string) (url string, err error)
}

type LocalStorage struct {
	BasePath string
}

func NewLocalStorage(basePath string) *LocalStorage {
	return &LocalStorage{BasePath: basePath}
}

func (l *LocalStorage) Save(file multipart.File, filename string) (string, error) {
	ext := filepath.Ext(filename)

	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("生成文件名失败：%v", err)
	}
	newName := hex.EncodeToString(bytes) + ext

	dateDir := time.Now().Format("2006-01-02")
	dir := filepath.Join(l.BasePath, dateDir)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("创建目录失败：%v", err)
	}

	dstPath := filepath.Join(dir, newName)
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", fmt.Errorf("创建文件失败：%v", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("保存文件失败：%v", err)
	}

	url := "/" + l.BasePath + "/" + dateDir + "/" + newName
	return url, nil
}

func (l *LocalStorage) SaveAs(file multipart.File, filename string, customName string) (string, error) {
	ext := filepath.Ext(filename)
	newName := customName + ext

	if err := os.MkdirAll(l.BasePath, 0755); err != nil {
		return "", fmt.Errorf("创建目录失败：%v", err)
	}

	dstPath := filepath.Join(l.BasePath, newName)
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", fmt.Errorf("创建文件失败：%v", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("保存文件失败：%v", err)
	}

	url := "/" + l.BasePath + "/" + newName
	return url, nil
}

