package category

import (
	"database/sql"
	"time"
)

type Category struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CategoryVO struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetCategoryByID(id int64) (*CategoryVO, error) {
	row := r.db.QueryRow("SELECT id, name FROM category WHERE id = ? AND is_deleted = 0", id)
	var category CategoryVO
	err := row.Scan(&category.ID, &category.Name)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *Repository) GetCategoryByName(name string) (*CategoryVO, error) {
	row := r.db.QueryRow("SELECT id, name FROM category WHERE name = ? AND is_deleted = 0", name)
	var category CategoryVO
	err := row.Scan(&category.ID, &category.Name)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *Repository) GetCategoryList() ([]*CategoryVO, error) {
	rows, err := r.db.Query("SELECT id, name FROM category WHERE is_deleted = 0")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*CategoryVO
	for rows.Next() {
		var category CategoryVO
		err := rows.Scan(&category.ID, &category.Name)
		if err != nil {
			return nil, err
		}
		categories = append(categories, &category)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *Repository) CreateCategory(name string) (*CategoryVO, error) {
	result, err := r.db.Exec("INSERT INTO category (name) VALUES (?)", name)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return r.GetCategoryByID(id)
}

func (r *Repository) UpdateCategory(id int64, name string) (*CategoryVO, error) {
	_, err := r.db.Exec("UPDATE category SET name = ? WHERE id = ? AND is_deleted = 0", name, id)
	if err != nil {
		return nil, err
	}
	return r.GetCategoryByID(id)
}

func (r *Repository) DeleteCategory(id int64) error {
	_, err := r.db.Exec("UPDATE category SET is_deleted = 1 WHERE id = ? AND is_deleted = 0", id)
	return err
}
