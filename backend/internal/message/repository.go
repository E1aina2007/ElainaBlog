package message

import (
	"context"

	"gorm.io/gorm"
)

// Repository 使用 GORM 存储留言数据。
type Repository struct {
	db *gorm.DB
}

// NewRepository 创建留言仓储实例。
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetList(ctx context.Context, limit int) ([]*MessageVO, error) {
	var messages []*MessageVO
	err := r.db.WithContext(ctx).Table("message m").
		Select(`m.id, m.user_id, u.username, COALESCE(u.avatar,'') AS avatar, u.is_admin, m.content, m.created_at`).
		Joins("LEFT JOIN `user` u ON m.user_id = u.id").
		Where("m.is_deleted = 0").
		Order("m.created_at DESC").
		Limit(limit).
		Scan(&messages).Error
	return messages, err
}

func (r *Repository) Create(ctx context.Context, msg *Message) (int64, error) {
	if err := r.db.WithContext(ctx).Table("message").Create(msg).Error; err != nil {
		return 0, err
	}
	return msg.ID, nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*Message, error) {
	var msg Message
	err := r.db.WithContext(ctx).Table("message").
		Where("id = ? AND is_deleted = 0", id).
		First(&msg).Error
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Table("message").
		Where("id = ? AND is_deleted = 0", id).
		Update("is_deleted", 1).Error
}
