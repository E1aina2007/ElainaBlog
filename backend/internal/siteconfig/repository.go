package siteconfig

import (
	"database/sql"
)

type SiteConfig struct {
	ID      int64  `json:"id"`
	KeyName string `json:"key_name"`
	Value   string `json:"value"`
}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetAll() ([]*SiteConfig, error) {
	rows, err := r.db.Query(`SELECT id, key_name, value FROM site_config WHERE is_deleted = 0`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var configs []*SiteConfig
	for rows.Next() {
		var c SiteConfig
		if err := rows.Scan(&c.ID, &c.KeyName, &c.Value); err != nil {
			return nil, err
		}
		configs = append(configs, &c)
	}
	return configs, rows.Err()
}

func (r *Repository) GetByKey(key string) (*SiteConfig, error) {
	var c SiteConfig
	err := r.db.QueryRow(`SELECT id, key_name, value FROM site_config WHERE key_name = ? AND is_deleted = 0`, key).Scan(
		&c.ID, &c.KeyName, &c.Value)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *Repository) Upsert(key, value string) error {
	_, err := r.db.Exec(`
		INSERT INTO site_config (key_name, value) VALUES (?, ?)
		ON DUPLICATE KEY UPDATE value = VALUES(value), updated_at = CURRENT_TIMESTAMP
	`, key, value)
	return err
}

func (r *Repository) DeleteByKey(key string) error {
	_, err := r.db.Exec(`UPDATE site_config SET is_deleted = 1 WHERE key_name = ? AND is_deleted = 0`, key)
	return err
}
