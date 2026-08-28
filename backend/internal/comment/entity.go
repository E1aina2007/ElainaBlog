package comment

import "time"

// Comment 评论实体
type Comment struct {
	ID               int64     `json:"id"`
	ArticleID        int64     `json:"article_id"`
	UserID           int64     `json:"user_id"`
	ReplyToUserID    *int64    `json:"reply_to_user_id"`
	ReplyToUsername  *string   `json:"reply_to_username"`
	ReplyToCommentID *int64    `json:"reply_to_comment_id"`
	ReplyToContent   *string   `json:"reply_to_content"`
	Content          string    `json:"content"`
	CreatedAt        time.Time `json:"created_at"`
}

// CommentVO 返回给前端的评论视图对象
type CommentVO struct {
	ID               int64     `json:"id"`
	ArticleID        int64     `json:"article_id"`
	UserID           int64     `json:"user_id"`
	ReplyToUserID    *int64    `json:"reply_to_user_id"`
	ReplyToUsername  *string   `json:"reply_to_username"`
	ReplyToCommentID *int64    `json:"reply_to_comment_id"`
	ReplyToContent   *string   `json:"reply_to_content"`
	Username         string    `json:"username"`
	Avatar           string    `json:"avatar"`
	IsAdmin          bool      `json:"is_admin"`
	Content          string    `json:"content"`
	CreatedAt        time.Time `json:"created_at"`
}

// CreateCommentRequest 创建评论请求
type CreateCommentRequest struct {
	ArticleID        int64  `json:"article_id"`
	ReplyToUserID    *int64 `json:"reply_to_user_id"`
	ReplyToCommentID *int64 `json:"reply_to_comment_id"`
	Content          string `json:"content"`
}

// DeleteCommentRequest 删除评论请求
type DeleteCommentRequest struct {
	ID int64 `json:"id"`
}

// CreateCommentParams 创建评论参数
type CreateCommentParams struct {
	ArticleID        int64
	UserID           int64
	ReplyToUserID    *int64
	ReplyToCommentID *int64
	Content          string
}

// DeleteCommentParams 删除评论参数
type DeleteCommentParams struct {
	ID int64
}
