package category

import (
	"context"

	"gorm.io/gorm"
)

// Repository 使用 GORM 存储分类数据。
type Repository struct {
	db *gorm.DB
}

// NewRepository 创建分类仓储实例。
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetCategoryByID(ctx context.Context, id int64) (*CategoryVO, error) {
	var category CategoryVO
	err := r.db.WithContext(ctx).Table("category").
		Select("id", "name", "is_top").
		Where("id = ? AND is_deleted = 0", id).
		First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *Repository) GetCategoryByName(ctx context.Context, name string) (*CategoryVO, error) {
	var category CategoryVO
	err := r.db.WithContext(ctx).Table("category").
		Select("id", "name", "is_top").
		Where("name = ? AND is_deleted = 0", name).
		First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *Repository) GetCategoryList(ctx context.Context) ([]*CategoryVO, error) {
	var categories []*CategoryVO
	err := r.db.WithContext(ctx).Table("category c").
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

func (r *Repository) CreateCategory(ctx context.Context, name string) (*CategoryVO, error) {
	var category CategoryVO
	// 事务绑定同一连接：LAST_INSERT_ID 为连接级变量，池化下跨连接查询会拿到错误 ID
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("INSERT INTO category (name) VALUES (?)", name).Error; err != nil {
			return err
		}
		return tx.Raw("SELECT LAST_INSERT_ID() AS id").Scan(&category).Error
	})
	if err != nil {
		return nil, err
	}
	return r.GetCategoryByID(ctx, category.ID)
}

func (r *Repository) UpdateCategory(ctx context.Context, id int64, name string) (*CategoryVO, error) {
	if err := r.db.WithContext(ctx).Exec("UPDATE category SET name = ? WHERE id = ? AND is_deleted = 0", name, id).Error; err != nil {
		return nil, err
	}
	return r.GetCategoryByID(ctx, id)
}

// ToggleCategoryTop 单独切换分类置顶状态
func (r *Repository) ToggleCategoryTop(ctx context.Context, id int64, isTop bool) error {
	return r.db.WithContext(ctx).Exec("UPDATE category SET is_top = ? WHERE id = ? AND is_deleted = 0", isTop, id).Error
}

func (r *Repository) DeleteCategory(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Exec("UPDATE category SET is_deleted = 1 WHERE id = ? AND is_deleted = 0", id).Error
}
