// repository_contract.go 定义 notification 仓储层接口，用于依赖注入和测试 mock
package notification

// Repository 接口封装通知数据访问操作。
type Repository interface {
	Create(n *Notification) (int64, error)
	GetByUserID(userID int64, onlyUnread bool) ([]*NotificationVO, error)
	GetUnreadCount(userID int64) (int, error)
	MarkAsRead(id int64, userID int64) error
	MarkAllAsRead(userID int64) error
	Delete(id int64, userID int64) error
}
