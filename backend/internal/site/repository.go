package site

import "database/sql"

type Site struct {
	Key   string
	Value string
}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetList() (map[string]string, error) {
	rows, err := r.db.Query("SELECT `key`, `value` FROM `site_config`")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]string)
	for rows.Next() {
		var site Site
		err := rows.Scan(&site.Key, &site.Value)
		if err != nil {
			return nil, err
		}
		result[site.Key] = site.Value
	}
	return result, nil
}

func (r *Repository) Update(configs map[string]string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	for key, value := range configs {
		_, err := tx.Exec("UPDATE `site_config` SET `value` = ? WHERE `key` = ?", value, key)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// DashboardStats 仪表盘统计数据
type DashboardStats struct {
	ArticleCount int64 `json:"article_count"`
	CommentCount int64 `json:"comment_count"`
	UserCount    int64 `json:"user_count"`
}

// GetDashboardStats 获取仪表盘统计数据
func (r *Repository) GetDashboardStats() (*DashboardStats, error) {
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
