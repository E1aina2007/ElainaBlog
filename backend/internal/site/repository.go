package site

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// Repository 使用 GORM 存储站点统计数据。
type Repository struct {
	db *gorm.DB
}

// NewRepository 创建站点仓储实例。
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// GetDashboardStats 获取仪表盘统计数据
func (r *Repository) GetDashboardStats(ctx context.Context) (*DashboardStats, error) {
	var stats DashboardStats

	if err := r.db.WithContext(ctx).Table("article").Where("is_deleted = 0 AND is_draft = 0").Count(&stats.ArticleCount).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).Table("comment").Where("is_deleted = 0").Count(&stats.CommentCount).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).Table("user").Where("is_deleted = 0").Count(&stats.UserCount).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}

// GetAuthorStats 获取作者页统计数据
func (r *Repository) GetAuthorStats(ctx context.Context) (*AuthorStats, error) {
	var stats AuthorStats

	if err := r.db.WithContext(ctx).Table("article").Where("is_deleted = 0 AND is_draft = 0").Count(&stats.ArticleCount).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).Table("comment").Where("is_deleted = 0").Count(&stats.CommentCount).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).Table("article").Where("is_deleted = 0 AND is_draft = 0").Select("COALESCE(SUM(view_count), 0)").Scan(&stats.TotalViews).Error; err != nil {
		return nil, err
	}

	// 建站天数
	var earliest *time.Time
	if err := r.db.WithContext(ctx).Table("article").Where("is_deleted = 0").Select("MIN(created_at)").Scan(&earliest).Error; err != nil {
		return nil, err
	}
	if earliest != nil {
		stats.DaysSince = int(time.Since(*earliest).Hours() / 24)
		if stats.DaysSince < 1 {
			stats.DaysSince = 1
		}
	}

	return &stats, nil
}
