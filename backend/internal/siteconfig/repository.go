package siteconfig

import (
	"context"

	"gorm.io/gorm"
)

// Repository 使用 GORM 存储站点配置数据。
type Repository struct {
	db *gorm.DB
}

// NewRepository 创建站点配置仓储实例。
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetAll(ctx context.Context) ([]*SiteConfig, error) {
	var configs []*SiteConfig
	err := r.db.WithContext(ctx).Table("site_config").
		Select("id", "key_name", "value").
		Where("is_deleted = 0").
		Scan(&configs).Error
	return configs, err
}

func (r *Repository) GetByKey(ctx context.Context, key string) (*SiteConfig, error) {
	var c SiteConfig
	err := r.db.WithContext(ctx).Table("site_config").
		Select("id", "key_name", "value").
		Where("key_name = ? AND is_deleted = 0", key).
		First(&c).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *Repository) Upsert(ctx context.Context, key, value string) error {
	return r.db.WithContext(ctx).Exec(`
		INSERT INTO site_config (key_name, value) VALUES (?, ?)
		ON DUPLICATE KEY UPDATE value = VALUES(value), updated_at = CURRENT_TIMESTAMP
	`, key, value).Error
}

func (r *Repository) DeleteByKey(ctx context.Context, key string) error {
	return r.db.WithContext(ctx).Table("site_config").
		Where("key_name = ? AND is_deleted = 0", key).
		Update("is_deleted", 1).Error
}
