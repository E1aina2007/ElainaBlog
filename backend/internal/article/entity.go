package article

import "time"

// ArticleVO 返回给前端的文章视图对象
type ArticleVO struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	Username     string    `json:"author_name" gorm:"column:username"`
	Avatar       string    `json:"author_avatar" gorm:"column:author_avatar"`
	IsAdmin      bool      `json:"author_is_admin" gorm:"column:author_is_admin"`
	CategoryID   *int64    `json:"category_id"`
	CategoryName string    `json:"category_name" gorm:"column:category_name"`
	Title        string    `json:"title"`
	Summary      string    `json:"summary"`
	Content      string    `json:"content"`
	Tags         string    `json:"tags"`
	IsTop        bool      `json:"is_top"`
	IsDraft      bool      `json:"is_draft"`
	ViewCount    int       `json:"view_count"`
	UVCount      int64     `json:"uv_count"`
	CommentCount int       `json:"comment_count" gorm:"column:comment_count"`
	CreatedAt    time.Time `json:"created_at"`
}

// CreateArticleRequest 创建文章请求
type CreateArticleRequest struct {
	Title      string `json:"title"`
	Summary    string `json:"summary"`
	Content    string `json:"content"`
	Tags       string `json:"tags"`
	CategoryID *int64 `json:"category_id"` // nil 表示未分类
	IsTop      bool   `json:"is_top"`
	IsDraft    bool   `json:"is_draft"`
}

// UpdateArticleRequest 更新文章请求
type UpdateArticleRequest struct {
	ID         int64  `json:"id"`
	Title      string `json:"title"`
	Summary    string `json:"summary"`
	Content    string `json:"content"`
	Tags       string `json:"tags"`
	CategoryID *int64 `json:"category_id"`
	IsTop      bool   `json:"is_top"`
	IsDraft    bool   `json:"is_draft"`
}

// DeleteArticleRequest 删除文章请求
type DeleteArticleRequest struct {
	ID int64 `json:"id"`
}

// ToggleTopRequest 切换文章置顶请求
type ToggleTopRequest struct {
	ID    int64 `json:"id"`
	IsTop bool  `json:"is_top"`
}

// CreateArticleParams 创建文章参数
type CreateArticleParams struct {
	UserID     int64
	CategoryID *int64
	Title      string
	Summary    string
	Content    string
	Tags       string
	IsTop      bool
	IsDraft    bool
}

// UpdateArticleParams 更新文章参数
type UpdateArticleParams struct {
	ID         int64
	CategoryID *int64
	Title      string
	Summary    string
	Content    string
	Tags       string
	IsTop      bool
	IsDraft    bool
}

// DeleteArticleParams 删除文章参数
type DeleteArticleParams struct {
	ID int64
}

// ArticleListParams 文章列表查询参数
type ArticleListParams struct {
	CategoryID *int64
	SortBy     string // "latest"（默认） | "popular"
	Page       int
	PageSize   int
}

// ArticleListResult 文章列表返回结果
type ArticleListResult struct {
	List  []*ArticleVO `json:"list"`
	Total int          `json:"total"`
}
