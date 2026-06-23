package siteconfig

import (
	"gorm.io/gorm"
)

type SiteConfig struct {
	ID      int64  `json:"id"`
	KeyName string `json:"key_name"`
	Value   string `json:"value"`
}

// MySQLRepository 实现 siteconfig.Repository 接口，使用 GORM 存储。
type MySQLRepository struct {
	db *gorm.DB
}

// NewRepository 创建站点配置仓储实例。
func NewRepository(db *gorm.DB) *MySQLRepository {
	return &MySQLRepository{db: db}
}

func (r *MySQLRepository) GetAll() ([]*SiteConfig, error) {
	var configs []*SiteConfig
	err := r.db.Table("site_config").
		Select("id", "key_name", "value").
		Where("is_deleted = 0").
		Scan(&configs).Error
	return configs, err
}

func (r *MySQLRepository) GetByKey(key string) (*SiteConfig, error) {
	var c SiteConfig
	err := r.db.Table("site_config").
		Select("id", "key_name", "value").
		Where("key_name = ? AND is_deleted = 0", key).
		Scan(&c).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *MySQLRepository) Upsert(key, value string) error {
	return r.db.Exec(`
		INSERT INTO site_config (key_name, value) VALUES (?, ?)
		ON DUPLICATE KEY UPDATE value = VALUES(value), updated_at = CURRENT_TIMESTAMP
	`, key, value).Error
}

func (r *MySQLRepository) DeleteByKey(key string) error {
	return r.db.Table("site_config").
		Where("key_name = ? AND is_deleted = 0", key).
		Update("is_deleted", 1).Error
}
