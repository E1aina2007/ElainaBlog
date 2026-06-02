package site

import (
	"ElainaBlog/config/db"
	"database/sql"
	"time"
)

// MySQLRepository 实现 site.Repository 接口，使用 MySQL 存储。
type MySQLRepository struct {
	db db.DBTX
}

// NewRepository 创建站点仓储实例。
func NewRepository(db db.DBTX) *MySQLRepository {
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

	// 文章数（不包含草稿和已删除）
	err := r.db.QueryRow("SELECT COUNT(*) FROM article WHERE is_deleted = 0 AND is_draft = 0").Scan(&stats.ArticleCount)
	if err != nil {
		return nil, err
	}

	// 评论数（不包含已删除）
	err = r.db.QueryRow("SELECT COUNT(*) FROM comment WHERE is_deleted = 0").Scan(&stats.CommentCount)
	if err != nil {
		return nil, err
	}

	// 用户数（不包含已删除）
	err = r.db.QueryRow("SELECT COUNT(*) FROM `user` WHERE is_deleted = 0").Scan(&stats.UserCount)
	if err != nil {
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

	err := r.db.QueryRow("SELECT COUNT(*) FROM article WHERE is_deleted = 0 AND is_draft = 0").Scan(&stats.ArticleCount)
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow("SELECT COUNT(*) FROM comment WHERE is_deleted = 0").Scan(&stats.CommentCount)
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow("SELECT COALESCE(SUM(view_count), 0) FROM article WHERE is_deleted = 0 AND is_draft = 0").Scan(&stats.TotalViews)
	if err != nil {
		return nil, err
	}

	// 建站天数：取最早一篇文章的创建时间距今天数，无文章则为 0
	var earliest sql.NullTime
	err = r.db.QueryRow("SELECT MIN(created_at) FROM article WHERE is_deleted = 0").Scan(&earliest)
	if err != nil {
		return nil, err
	}
	if earliest.Valid {
		stats.DaysSince = int(time.Since(earliest.Time).Hours() / 24)
		if stats.DaysSince < 1 {
			stats.DaysSince = 1
		}
	}

	return &stats, nil
}
