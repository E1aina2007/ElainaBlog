// repository_contract.go 定义 message 仓储层接口，用于依赖注入和测试 mock
package message

// Repository 接口封装留言数据访问操作。
type Repository interface {
	GetList(limit int) ([]*MessageVO, error)
	Create(msg *Message) (int64, error)
	GetByID(id int64) (*Message, error)
	Delete(id int64) error
}
