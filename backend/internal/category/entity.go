package category

// CategoryVO 分类视图对象（数据模型）
type CategoryVO struct {
	ID           int64  `json:"id" gorm:"column:id"`
	Name         string `json:"name" gorm:"column:name"`
	IsTop        bool   `json:"is_top" gorm:"column:is_top"`
	ArticleCount int    `json:"article_count" gorm:"column:article_count"`
}

// CreateCategoryRequest 创建分类请求
type CreateCategoryRequest struct {
	Name string `json:"name"`
}

// UpdateCategoryRequest 更新分类请求
type UpdateCategoryRequest struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// DeleteCategoryRequest 删除分类请求
type DeleteCategoryRequest struct {
	ID int64 `json:"id"`
}

// ToggleTopRequest 切换置顶请求
type ToggleTopRequest struct {
	ID    int64 `json:"id"`
	IsTop bool  `json:"is_top"`
}

// CreateCategoryParams 创建分类参数
type CreateCategoryParams struct {
	Name string
}

// UpdateCategoryParams 更新分类参数
type UpdateCategoryParams struct {
	ID   int64
	Name string
}
