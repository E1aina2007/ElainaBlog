package notification

import "time"

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

// MarkReadRequest 标记已读请求
type MarkReadRequest struct {
	ID int64 `json:"id"`
}

// DeleteRequest 删除通知请求
type DeleteRequest struct {
	ID int64 `json:"id"`
}
