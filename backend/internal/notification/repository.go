package notification

import (
	"ElainaBlog/config/db"
	"time"
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

// MySQLRepository 实现 notification.Repository 接口，使用 MySQL 存储。
type MySQLRepository struct {
	db db.DBTX
}

// NewRepository 创建通知仓储实例。
func NewRepository(db db.DBTX) *MySQLRepository {
	return &MySQLRepository{db: db}
}

func (r *MySQLRepository) Create(n *Notification) (int64, error) {
	result, err := r.db.Exec(`INSERT INTO notification (user_id, type, title, content, target_id)
		VALUES (?, ?, ?, ?, ?)`, n.UserID, n.Type, n.Title, n.Content, n.TargetID)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *MySQLRepository) GetByUserID(userID int64, onlyUnread bool) ([]*NotificationVO, error) {
	query := `SELECT id, type, title, content, target_id, is_read, created_at
		FROM notification WHERE user_id = ? AND is_deleted = 0`
	args := []interface{}{userID}

	if onlyUnread {
		query += " AND is_read = 0"
	}

	query += " ORDER BY created_at DESC LIMIT 50"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notifications := make([]*NotificationVO, 0)
	for rows.Next() {
		var vo NotificationVO
		err := rows.Scan(&vo.ID, &vo.Type, &vo.Title, &vo.Content, &vo.TargetID, &vo.IsRead, &vo.CreatedAt)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, &vo)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return notifications, nil
}

func (r *MySQLRepository) GetUnreadCount(userID int64) (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM notification
		WHERE user_id = ? AND is_read = 0 AND is_deleted = 0`, userID).Scan(&count)
	return count, err
}

func (r *MySQLRepository) MarkAsRead(id int64, userID int64) error {
	_, err := r.db.Exec(`UPDATE notification SET is_read = 1
		WHERE id = ? AND user_id = ? AND is_deleted = 0`, id, userID)
	return err
}

func (r *MySQLRepository) MarkAllAsRead(userID int64) error {
	_, err := r.db.Exec(`UPDATE notification SET is_read = 1
		WHERE user_id = ? AND is_read = 0 AND is_deleted = 0`, userID)
	return err
}

func (r *MySQLRepository) Delete(id int64, userID int64) error {
	_, err := r.db.Exec(`UPDATE notification SET is_deleted = 1
		WHERE id = ? AND user_id = ? AND is_deleted = 0`, id, userID)
	return err
}
