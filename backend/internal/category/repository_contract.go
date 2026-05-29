// repository_contract.go 定义 category 仓储层接口，用于依赖注入和测试 mock
package category

// Repository 接口封装分类数据访问操作。
type Repository interface {
	GetCategoryByID(id int64) (*CategoryVO, error)
	GetCategoryByName(name string) (*CategoryVO, error)
	GetCategoryList() ([]*CategoryVO, error)
	CreateCategory(name string) (*CategoryVO, error)
	UpdateCategory(id int64, name string) (*CategoryVO, error)
	DeleteCategory(id int64) error
}
