// repository_contract.go 定义 user 仓储层接口，用于依赖注入和测试 mock
package user

// Repository 接口封装用户数据访问操作。
type Repository interface {
	GetUserByUsername(username string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id int64) (*User, error)
	GetUserList() ([]*User, error)
	CreateUser(user *User) (int64, error)
	UpdateProfile(id int64, username, email, avatar string) error
	UpdatePassword(id int64, newPassword string) error
	DeleteUser(id int64) error
}
