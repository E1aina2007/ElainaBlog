// cleanup.go 孤儿图片清理服务
// 扫描 uploads 目录，对比所有有效文章中的图片引用，删除未被引用的文件
package upload

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"go.uber.org/zap"
)

// ArticleContentProvider 提供文章内容的接口
type ArticleContentProvider interface {
	GetAllActiveContents() ([]string, error)
}

// ImageCleanup 孤儿图片清理服务
type ImageCleanup struct {
	uploadPath    string
	avatarPath    string
	contentProvider ArticleContentProvider
	logger        *zap.Logger
}

// NewImageCleanup 创建图片清理服务实例
func NewImageCleanup(uploadPath, avatarPath string, contentProvider ArticleContentProvider, logger *zap.Logger) *ImageCleanup {
	return &ImageCleanup{
		uploadPath:      uploadPath,
		avatarPath:      avatarPath,
		contentProvider: contentProvider,
		logger:          logger,
	}
}

// CleanupResult 清理结果
type CleanupResult struct {
	ScannedFiles  int      // 扫描的文件数
	ReferencedFiles int    // 被引用的文件数
	DeletedFiles  int      // 删除的文件数
	DeletedPaths  []string // 删除的文件路径列表
	Errors        []error  // 删除失败的错误列表
}

// CleanupOrphanImages 清理孤儿图片
// 1. 获取所有有效文章的 content
// 2. 用正则提取所有图片 URL
// 3. 扫描 uploads 目录获取所有文件
// 4. 删除不在引用集合中的文件
func (s *ImageCleanup) CleanupOrphanImages() (*CleanupResult, error) {
	result := &CleanupResult{
		DeletedPaths: make([]string, 0),
		Errors:       make([]error, 0),
	}

	// 1. 获取所有有效文章的 content
	contents, err := s.contentProvider.GetAllActiveContents()
	if err != nil {
		return nil, fmt.Errorf("获取文章内容失败: %w", err)
	}

	// 2. 提取所有被引用的图片 URL
	// 匹配 Markdown 图片语法: ![alt](/uploads/YYYY-MM-DD/filename.ext)
	// 同时匹配可能的变体: ![alt](uploads/...) 或无 alt 文本
	imgRegex := regexp.MustCompile(`!\[.*?\]\((/uploads/\d{4}-\d{2}-\d{2}/[^)]+)\)`)
	referencedImages := make(map[string]bool)

	for _, content := range contents {
		matches := imgRegex.FindAllStringSubmatch(content, -1)
		for _, match := range matches {
			if len(match) > 1 {
				// 转换为文件系统相对路径: /uploads/2024-01-15/xxx.jpg -> 2024-01-15/xxx.jpg
				relPath := strings.TrimPrefix(match[1], "/uploads/")
				referencedImages[relPath] = true
			}
		}
	}

	s.logger.Info("图片引用统计", zap.Int("referenced_count", len(referencedImages)))
	result.ReferencedFiles = len(referencedImages)

	// 3. 扫描 uploads 目录，跳过 avatars 目录
	err = filepath.Walk(s.uploadPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录本身
		if info.IsDir() {
			return nil
		}

		// 跳过 avatars 目录下的文件
		relToBase, _ := filepath.Rel(s.uploadPath, path)
		if strings.HasPrefix(filepath.ToSlash(relToBase), "avatars/") {
			return nil
		}

		result.ScannedFiles++

		// 获取相对于 uploadPath 的路径
		relPath := filepath.ToSlash(relToBase)

		// 4. 不在引用集合中则删除
		if !referencedImages[relPath] {
			if err := os.Remove(path); err != nil {
				result.Errors = append(result.Errors, fmt.Errorf("删除 %s 失败: %w", relPath, err))
				s.logger.Error("删除孤儿图片失败", zap.String("path", relPath), zap.Error(err))
			} else {
				result.DeletedFiles++
				result.DeletedPaths = append(result.DeletedPaths, relPath)
				s.logger.Info("已删除孤儿图片", zap.String("path", relPath))
			}
		}

		return nil
	})

	if err != nil {
		return result, fmt.Errorf("扫描目录失败: %w", err)
	}

	// 清理空目录（从深到浅）
	s.cleanupEmptyDirs(s.uploadPath)

	return result, nil
}

// cleanupEmptyDirs 清理空的日期目录
func (s *ImageCleanup) cleanupEmptyDirs(basePath string) {
	entries, err := os.ReadDir(basePath)
	if err != nil {
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// 跳过 avatars 目录
		if entry.Name() == "avatars" {
			continue
		}

		dirPath := filepath.Join(basePath, entry.Name())

		// 检查目录是否为空
		dirEntries, err := os.ReadDir(dirPath)
		if err != nil {
			continue
		}

		if len(dirEntries) == 0 {
			if err := os.Remove(dirPath); err == nil {
				s.logger.Info("已删除空目录", zap.String("path", dirPath))
			}
		}
	}
}
