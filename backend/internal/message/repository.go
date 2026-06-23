package message

import (
	"time"

	"gorm.io/gorm"
)

type Message struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type MessageVO struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Username  string    `json:"username"`
	Avatar    string    `json:"avatar"`
	IsAdmin   bool      `json:"is_admin"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// MySQLRepository 实现 message.Repository 接口，使用 GORM 存储。
type MySQLRepository struct {
	db *gorm.DB
}

// NewRepository 创建留言仓储实例。
func NewRepository(db *gorm.DB) *MySQLRepository {
	return &MySQLRepository{db: db}
}

func (r *MySQLRepository) GetList(limit int) ([]*MessageVO, error) {
	var messages []*MessageVO
	err := r.db.Table("message m").
		Select(`m.id, m.user_id, u.username, COALESCE(u.avatar,'') AS avatar, u.is_admin, m.content, m.created_at`).
		Joins("LEFT JOIN `user` u ON m.user_id = u.id").
		Where("m.is_deleted = 0").
		Order("m.created_at DESC").
		Limit(limit).
		Scan(&messages).Error
	return messages, err
}

func (r *MySQLRepository) Create(msg *Message) (int64, error) {
	if err := r.db.Table("message").Create(msg).Error; err != nil {
		return 0, err
	}
	return msg.ID, nil
}

func (r *MySQLRepository) GetByID(id int64) (*Message, error) {
	var msg Message
	err := r.db.Table("message").
		Where("id = ? AND is_deleted = 0", id).
		First(&msg).Error
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

func (r *MySQLRepository) Delete(id int64) error {
	return r.db.Table("message").
		Where("id = ? AND is_deleted = 0", id).
		Update("is_deleted", 1).Error
}
