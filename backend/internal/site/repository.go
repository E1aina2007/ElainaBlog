package site

import (
	"time"

	"gorm.io/gorm"
)

// MySQLRepository 实现 site.Repository 接口，使用 GORM 存储。
type MySQLRepository struct {
	db *gorm.DB
}

// NewRepository 创建站点仓储实例。
func NewRepository(db *gorm.DB) *MySQLRepository {
	return &MySQLRepository{db: db}
}

// DashboardStats 仪表盘统计数据
type DashboardStats struct {
	ArticleCount int64 `json:"article_count"`
	CommentCount int64 `json:"comment_count"`
	UserCount    int64 `json:"user_count"`
	TodayPV      int64 `json:"today_pv"`
	TodayUV      int64 `json:"today_uv"`
	YesterdayPV  int64 `json:"yesterday_pv"`
	YesterdayUV  int64 `json:"yesterday_uv"`
}

// GetDashboardStats 获取仪表盘统计数据
func (r *MySQLRepository) GetDashboardStats() (*DashboardStats, error) {
	var stats DashboardStats

	if err := r.db.Table("article").Where("is_deleted = 0 AND is_draft = 0").Count(&stats.ArticleCount).Error; err != nil {
		return nil, err
	}
	if err := r.db.Table("comment").Where("is_deleted = 0").Count(&stats.CommentCount).Error; err != nil {
		return nil, err
	}
	if err := r.db.Table("`user`").Where("is_deleted = 0").Count(&stats.UserCount).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}

// AuthorStats 关于作者页统计数据
type AuthorStats struct {
	ArticleCount int64 `json:"article_count"`
	CommentCount int64 `json:"comment_count"`
	TotalViews   int64 `json:"total_views"`
	DaysSince    int   `json:"days_since"`
}

// GetAuthorStats 获取作者页统计数据
func (r *MySQLRepository) GetAuthorStats() (*AuthorStats, error) {
	var stats AuthorStats

	if err := r.db.Table("article").Where("is_deleted = 0 AND is_draft = 0").Count(&stats.ArticleCount).Error; err != nil {
		return nil, err
	}
	if err := r.db.Table("comment").Where("is_deleted = 0").Count(&stats.CommentCount).Error; err != nil {
		return nil, err
	}
	if err := r.db.Table("article").Where("is_deleted = 0 AND is_draft = 0").Select("COALESCE(SUM(view_count), 0)").Scan(&stats.TotalViews).Error; err != nil {
		return nil, err
	}

	// 建站天数
	var earliest *time.Time
	if err := r.db.Table("article").Where("is_deleted = 0").Select("MIN(created_at)").Scan(&earliest).Error; err != nil {
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
