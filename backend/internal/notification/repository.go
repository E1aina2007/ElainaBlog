package notification

import (
	"context"

	"gorm.io/gorm"
)

// Repository 使用 GORM 存储通知数据。
type Repository struct {
	db *gorm.DB
}

// NewRepository 创建通知仓储实例。
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, n *Notification) (int64, error) {
	if err := r.db.WithContext(ctx).Table("notification").Create(n).Error; err != nil {
		return 0, err
	}
	return n.ID, nil
}

func (r *Repository) GetByUserID(ctx context.Context, userID int64, onlyUnread bool) ([]*NotificationVO, error) {
	query := r.db.WithContext(ctx).Table("notification").
		Select("id", "type", "title", "content", "target_id", "is_read", "created_at").
		Where("user_id = ? AND is_deleted = 0", userID)

	if onlyUnread {
		query = query.Where("is_read = 0")
	}

	var notifications []*NotificationVO
	err := query.Order("created_at DESC").Limit(50).Scan(&notifications).Error
	return notifications, err
}

func (r *Repository) GetUnreadCount(ctx context.Context, userID int64) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).Table("notification").
		Where("user_id = ? AND is_read = 0 AND is_deleted = 0", userID).
		Count(&count).Error
	return int(count), err
}

func (r *Repository) MarkAsRead(ctx context.Context, id int64, userID int64) error {
	return r.db.WithContext(ctx).Table("notification").
		Where("id = ? AND user_id = ? AND is_deleted = 0", id, userID).
		Update("is_read", 1).Error
}

func (r *Repository) MarkAllAsRead(ctx context.Context, userID int64) error {
	return r.db.WithContext(ctx).Table("notification").
		Where("user_id = ? AND is_read = 0 AND is_deleted = 0", userID).
		Update("is_read", 1).Error
}

func (r *Repository) Delete(ctx context.Context, id int64, userID int64) error {
	return r.db.WithContext(ctx).Table("notification").
		Where("id = ? AND user_id = ? AND is_deleted = 0", id, userID).
		Update("is_deleted", 1).Error
}
