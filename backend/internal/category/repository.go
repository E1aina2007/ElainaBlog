package category

import (
	"gorm.io/gorm"
)

type CategoryVO struct {
	ID           int64  `json:"id" gorm:"column:id"`
	Name         string `json:"name" gorm:"column:name"`
	IsTop        bool   `json:"is_top" gorm:"column:is_top"`
	ArticleCount int    `json:"article_count" gorm:"column:article_count"`
}

// MySQLRepository 实现 category.Repository 接口，使用 GORM 存储。
type MySQLRepository struct {
	db *gorm.DB
}

// NewRepository 创建分类仓储实例。
func NewRepository(db *gorm.DB) *MySQLRepository {
	return &MySQLRepository{db: db}
}

func (r *MySQLRepository) GetCategoryByID(id int64) (*CategoryVO, error) {
	var category CategoryVO
	err := r.db.Table("category").
		Select("id", "name", "is_top").
		Where("id = ? AND is_deleted = 0", id).
		First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *MySQLRepository) GetCategoryByName(name string) (*CategoryVO, error) {
	var category CategoryVO
	err := r.db.Table("category").
		Select("id", "name", "is_top").
		Where("name = ? AND is_deleted = 0", name).
		First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *MySQLRepository) GetCategoryList() ([]*CategoryVO, error) {
	var categories []*CategoryVO
	err := r.db.Table("category c").
		Select("c.id", "c.name", "c.is_top", "COUNT(a.id) AS article_count").
		Joins("LEFT JOIN article a ON a.category_id = c.id AND a.is_deleted = 0 AND a.is_draft = 0").
		Where("c.is_deleted = 0").
		Group("c.id, c.name, c.is_top").
		Order("c.is_top DESC, c.name ASC").
		Scan(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *MySQLRepository) CreateCategory(name string) (*CategoryVO, error) {
	if err := r.db.Exec("INSERT INTO category (name) VALUES (?)", name).Error; err != nil {
		return nil, err
	}
	var category CategoryVO
	if err := r.db.Raw("SELECT LAST_INSERT_ID() AS id").Scan(&category).Error; err != nil {
		return nil, err
	}
	return r.GetCategoryByID(category.ID)
}

func (r *MySQLRepository) UpdateCategory(id int64, name string) (*CategoryVO, error) {
	if err := r.db.Exec("UPDATE category SET name = ? WHERE id = ? AND is_deleted = 0", name, id).Error; err != nil {
		return nil, err
	}
	return r.GetCategoryByID(id)
}

// ToggleCategoryTop 单独切换分类置顶状态
func (r *MySQLRepository) ToggleCategoryTop(id int64, isTop bool) error {
	return r.db.Exec("UPDATE category SET is_top = ? WHERE id = ? AND is_deleted = 0", isTop, id).Error
}

func (r *MySQLRepository) DeleteCategory(id int64) error {
	return r.db.Exec("UPDATE category SET is_deleted = 1 WHERE id = ? AND is_deleted = 0", id).Error
}
