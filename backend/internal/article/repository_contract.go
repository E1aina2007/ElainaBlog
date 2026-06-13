// repository_contract.go 定义 article 仓储层接口，用于依赖注入和测试 mock
package article

// Repository 接口封装文章数据访问操作。
type Repository interface {
	GetArticleByID(id int64) (*ArticleVO, error)
	GetArticleByIDIncludeDraft(id int64) (*ArticleVO, error)
	GetArticleList(categoryID *int64, sortBy string, page, pageSize int) ([]*ArticleVO, int, error)
	GetAdminArticleList(categoryID *int64, sortBy string, page, pageSize int) ([]*ArticleVO, int, error)
	GetUserArticleList(userID int64, categoryID *int64, sortBy string, page, pageSize int) ([]*ArticleVO, int, error)
	IncrementViewCount(id int64) error
	IncrementViewCountUnique(id int64, clientIP string) error
	GetViewCountDelta(id int64) int
	GetArticleUV(id int64) (int64, error)
	FlushViewCounts() (int, error)
	CreateArticle(userID int64, categoryID *int64, title, summary, content string, isTop, isDraft bool) (int64, error)
	UpdateArticle(id int64, categoryID *int64, title, summary, content string, isTop, isDraft bool) error
	DeleteArticle(id int64) error
}
