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
