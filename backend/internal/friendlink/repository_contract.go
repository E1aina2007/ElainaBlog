// repository_contract.go 定义 friendlink 仓储层接口，用于依赖注入和测试 mock
package friendlink

// Repository 接口封装友情链接数据访问操作。
type Repository interface {
	GetByID(id int64) (*FriendLinkVO, error)
	GetList() ([]*FriendLinkVO, error)
	Create(link *FriendLink) (int64, error)
	Update(link *FriendLink) error
	Delete(id int64) error
}
