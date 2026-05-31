// repository_contract.go 定义 authorprofile 仓储层接口，用于依赖注入和测试 mock
package authorprofile

// Repository 接口封装作者信息数据访问操作。
type Repository interface {
	Get() (*AuthorProfile, error)
	Create(p *AuthorProfile) (int64, error)
	Update(p *AuthorProfile) error
}
