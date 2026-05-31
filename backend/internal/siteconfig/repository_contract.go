// repository_contract.go 定义 siteconfig 仓储层接口，用于依赖注入和测试 mock
package siteconfig

// Repository 接口封装站点配置数据访问操作。
type Repository interface {
	GetAll() ([]*SiteConfig, error)
	GetByKey(key string) (*SiteConfig, error)
	Upsert(key, value string) error
	DeleteByKey(key string) error
}
