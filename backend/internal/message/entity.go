package message

import "time"

// Message 留言实体
type Message struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// MessageVO 返回给前端的视图对象
type MessageVO struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Username  string    `json:"username"`
	Avatar    string    `json:"avatar"`
	IsAdmin   bool      `json:"is_admin"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateMessageRequest 创建留言请求
type CreateMessageRequest struct {
	Content string `json:"content"`
}

// DeleteMessageRequest 删除留言请求
type DeleteMessageRequest struct {
	ID int64 `json:"id"`
}
