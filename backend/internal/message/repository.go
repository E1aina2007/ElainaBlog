package message

import (
	"ElainaBlog/config/db"
	"time"
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

// MySQLRepository 实现 message.Repository 接口，使用 MySQL 存储。
type MySQLRepository struct {
	db db.DBTX
}

// NewRepository 创建留言仓储实例。
func NewRepository(db db.DBTX) *MySQLRepository {
	return &MySQLRepository{db: db}
}

func (r *MySQLRepository) GetList(limit int) ([]*MessageVO, error) {
	rows, err := r.db.Query(`
		SELECT m.id, m.user_id, u.username, COALESCE(u.avatar,''), u.is_admin, m.content, m.created_at
		FROM message m
		LEFT JOIN `+"`user`"+` u ON m.user_id = u.id
		WHERE m.is_deleted = 0
		ORDER BY m.created_at DESC
		LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := make([]*MessageVO, 0)
	for rows.Next() {
		var vo MessageVO
		if err := rows.Scan(&vo.ID, &vo.UserID, &vo.Username, &vo.Avatar, &vo.IsAdmin, &vo.Content, &vo.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, &vo)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *MySQLRepository) Create(msg *Message) (int64, error) {
	result, err := r.db.Exec(`INSERT INTO message (user_id, content) VALUES (?, ?)`, msg.UserID, msg.Content)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *MySQLRepository) GetByID(id int64) (*Message, error) {
	var msg Message
	err := r.db.QueryRow(`SELECT id, user_id, content, created_at FROM message WHERE id = ? AND is_deleted = 0`, id).Scan(
		&msg.ID, &msg.UserID, &msg.Content, &msg.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

func (r *MySQLRepository) Delete(id int64) error {
	_, err := r.db.Exec(`UPDATE message SET is_deleted = 1 WHERE id = ? AND is_deleted = 0`, id)
	return err
}
