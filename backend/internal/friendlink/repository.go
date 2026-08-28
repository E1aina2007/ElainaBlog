package friendlink

import (
	"context"

	"gorm.io/gorm"
)

// Repository 使用 GORM 存储友情链接数据。
type Repository struct {
	db *gorm.DB
}

// NewRepository 创建友情链接仓储实例。
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*FriendLinkVO, error) {
	var vo FriendLinkVO
	err := r.db.WithContext(ctx).Table("friend_link").
		Select("id", "name", "url", "avatar", "description", "sort_order").
		Where("id = ? AND is_deleted = 0", id).
		First(&vo).Error
	if err != nil {
		return nil, err
	}
	return &vo, nil
}

func (r *Repository) GetList(ctx context.Context) ([]*FriendLinkVO, error) {
	var links []*FriendLinkVO
	err := r.db.WithContext(ctx).Table("friend_link").
		Select("id", "name", "url", "avatar", "description", "sort_order").
		Where("is_deleted = 0").
		Order("sort_order DESC, id ASC").
		Scan(&links).Error
	return links, err
}

func (r *Repository) Create(ctx context.Context, link *FriendLink) (int64, error) {
	if err := r.db.WithContext(ctx).Table("friend_link").Create(link).Error; err != nil {
		return 0, err
	}
	return link.ID, nil
}

func (r *Repository) Update(ctx context.Context, link *FriendLink) error {
	return r.db.WithContext(ctx).Table("friend_link").
		Where("id = ? AND is_deleted = 0", link.ID).
		Updates(map[string]any{
			"name":        link.Name,
			"url":         link.URL,
			"avatar":      link.Avatar,
			"description": link.Description,
			"sort_order":  link.SortOrder,
		}).Error
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Table("friend_link").
		Where("id = ? AND is_deleted = 0", id).
		Update("is_deleted", 1).Error
}
