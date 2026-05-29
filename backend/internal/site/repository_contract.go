// repository_contract.go 定义 site 仓储层接口，用于依赖注入和测试 mock
package site

// Repository 接口封装站点统计数据访问操作。
type Repository interface {
	GetDashboardStats() (*DashboardStats, error)
	GetAuthorStats() (*AuthorStats, error)
}
