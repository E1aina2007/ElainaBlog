package siteconfig

import (
	"ElainaBlog/config/db"
)

type SiteConfig struct {
	ID      int64  `json:"id"`
	KeyName string `json:"key_name"`
	Value   string `json:"value"`
}

// MySQLRepository 实现 siteconfig.Repository 接口，使用 MySQL 存储。
type MySQLRepository struct {
	db db.DBTX
}

// NewRepository 创建站点配置仓储实例。
func NewRepository(db db.DBTX) *MySQLRepository {
	return &MySQLRepository{db: db}
}

func (r *MySQLRepository) GetAll() ([]*SiteConfig, error) {
	rows, err := r.db.Query(`SELECT id, key_name, value FROM site_config WHERE is_deleted = 0`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	configs := make([]*SiteConfig, 0)
	for rows.Next() {
		var c SiteConfig
		if err := rows.Scan(&c.ID, &c.KeyName, &c.Value); err != nil {
			return nil, err
		}
		configs = append(configs, &c)
	}
	return configs, rows.Err()
}

func (r *MySQLRepository) GetByKey(key string) (*SiteConfig, error) {
	var c SiteConfig
	err := r.db.QueryRow(`SELECT id, key_name, value FROM site_config WHERE key_name = ? AND is_deleted = 0`, key).Scan(
		&c.ID, &c.KeyName, &c.Value)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *MySQLRepository) Upsert(key, value string) error {
	_, err := r.db.Exec(`
		INSERT INTO site_config (key_name, value) VALUES (?, ?)
		ON DUPLICATE KEY UPDATE value = VALUES(value), updated_at = CURRENT_TIMESTAMP
	`, key, value)
	return err
}

func (r *MySQLRepository) DeleteByKey(key string) error {
	_, err := r.db.Exec(`UPDATE site_config SET is_deleted = 1 WHERE key_name = ? AND is_deleted = 0`, key)
	return err
}
