package category

import (
	"ElainaBlog/config/db"
)

type CategoryVO struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	ArticleCount int    `json:"article_count"`
}

// MySQLRepository 实现 category.Repository 接口，使用 MySQL 存储。
type MySQLRepository struct {
	db db.DBTX
}

// NewRepository 创建分类仓储实例。
func NewRepository(db db.DBTX) *MySQLRepository {
	return &MySQLRepository{db: db}
}

func (r *MySQLRepository) GetCategoryByID(id int64) (*CategoryVO, error) {
	row := r.db.QueryRow("SELECT id, name FROM category WHERE id = ? AND is_deleted = 0", id)
	var category CategoryVO
	err := row.Scan(&category.ID, &category.Name)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *MySQLRepository) GetCategoryByName(name string) (*CategoryVO, error) {
	row := r.db.QueryRow("SELECT id, name FROM category WHERE name = ? AND is_deleted = 0", name)
	var category CategoryVO
	err := row.Scan(&category.ID, &category.Name)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *MySQLRepository) GetCategoryList() ([]*CategoryVO, error) {
	rows, err := r.db.Query(`
		SELECT c.id, c.name, COUNT(a.id) AS article_count
		FROM category c
		LEFT JOIN article a ON a.category_id = c.id AND a.is_deleted = 0 AND a.is_draft = 0
		WHERE c.is_deleted = 0
		GROUP BY c.id, c.name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]*CategoryVO, 0)
	for rows.Next() {
		var category CategoryVO
		err := rows.Scan(&category.ID, &category.Name, &category.ArticleCount)
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

func (r *MySQLRepository) CreateCategory(name string) (*CategoryVO, error) {
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

func (r *MySQLRepository) UpdateCategory(id int64, name string) (*CategoryVO, error) {
	_, err := r.db.Exec("UPDATE category SET name = ? WHERE id = ? AND is_deleted = 0", name, id)
	if err != nil {
		return nil, err
	}
	return r.GetCategoryByID(id)
}

func (r *MySQLRepository) DeleteCategory(id int64) error {
	_, err := r.db.Exec("UPDATE category SET is_deleted = 1 WHERE id = ? AND is_deleted = 0", id)
	return err
}
