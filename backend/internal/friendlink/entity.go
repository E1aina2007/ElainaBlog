package friendlink

// FriendLink 友情链接实体
type FriendLink struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	Avatar      string `json:"avatar"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
}

// FriendLinkVO 返回给前端的视图对象
type FriendLinkVO struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	Avatar      string `json:"avatar"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
}

// CreateRequest 创建友链请求
type CreateRequest struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	Avatar      string `json:"avatar"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
}

// UpdateRequest 更新友链请求
type UpdateRequest struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	Avatar      string `json:"avatar"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
}

// DeleteRequest 删除友链请求
type DeleteRequest struct {
	ID int64 `json:"id"`
}

// CreateParams 创建友链参数
type CreateParams struct {
	Name        string
	URL         string
	Avatar      string
	Description string
	SortOrder   int
}

// UpdateParams 更新友链参数
type UpdateParams struct {
	ID          int64
	Name        string
	URL         string
	Avatar      string
	Description string
	SortOrder   int
}
