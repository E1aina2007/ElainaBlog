package notification

import (
	"time"

	"gorm.io/gorm"
)

// Notification 通知实体
type Notification struct {
	ID       int64  `json:"id"`
	UserID   int64  `json:"user_id"`
	Type     string `json:"type"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	TargetID int64  `json:"target_id"`
	IsRead   bool   `json:"is_read"`
}

// NotificationVO 返回给前端的视图对象
type NotificationVO struct {
	ID        int64     `json:"id"`
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	TargetID  int64     `json:"target_id"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}

// MySQLRepository 实现 notification.Repository 接口，使用 GORM 存储。
type MySQLRepository struct {
	db *gorm.DB
}

// NewRepository 创建通知仓储实例。
func NewRepository(db *gorm.DB) *MySQLRepository {
	return &MySQLRepository{db: db}
}

func (r *MySQLRepository) Create(n *Notification) (int64, error) {
	if err := r.db.Table("notification").Create(n).Error; err != nil {
		return 0, err
	}
	return n.ID, nil
}

func (r *MySQLRepository) GetByUserID(userID int64, onlyUnread bool) ([]*NotificationVO, error) {
	query := r.db.Table("notification").
		Select("id", "type", "title", "content", "target_id", "is_read", "created_at").
		Where("user_id = ? AND is_deleted = 0", userID)

	if onlyUnread {
		query = query.Where("is_read = 0")
	}

	var notifications []*NotificationVO
	err := query.Order("created_at DESC").Limit(50).Scan(&notifications).Error
	return notifications, err
}

func (r *MySQLRepository) GetUnreadCount(userID int64) (int, error) {
	var count int64
	err := r.db.Table("notification").
		Where("user_id = ? AND is_read = 0 AND is_deleted = 0", userID).
		Count(&count).Error
	return int(count), err
}

func (r *MySQLRepository) MarkAsRead(id int64, userID int64) error {
	return r.db.Table("notification").
		Where("id = ? AND user_id = ? AND is_deleted = 0", id, userID).
		Update("is_read", 1).Error
}

func (r *MySQLRepository) MarkAllAsRead(userID int64) error {
	return r.db.Table("notification").
		Where("user_id = ? AND is_read = 0 AND is_deleted = 0", userID).
		Update("is_read", 1).Error
}

func (r *MySQLRepository) Delete(id int64, userID int64) error {
	return r.db.Table("notification").
		Where("id = ? AND user_id = ? AND is_deleted = 0", id, userID).
		Update("is_deleted", 1).Error
}
