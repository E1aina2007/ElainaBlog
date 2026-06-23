package user

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Email     string    `json:"email"`
	Avatar    string    `json:"avatar"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// 返回给前端
type UserVO struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Avatar    string    `json:"avatar"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
}

func (u *User) ToVO() *UserVO {
	return &UserVO{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Avatar:    u.Avatar,
		IsAdmin:   u.IsAdmin,
		CreatedAt: u.CreatedAt,
	}
}

// MySQLRepository 实现 user.Repository 接口，使用 GORM 存储。
type MySQLRepository struct {
	db *gorm.DB
}

// NewRepository 创建用户仓储实例。
func NewRepository(db *gorm.DB) *MySQLRepository {
	return &MySQLRepository{db: db}
}

func (r *MySQLRepository) GetUserByUsername(username string) (*User, error) {
	var user User
	err := r.db.Table("`user`").
		Where("username = ? AND is_deleted = 0", username).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *MySQLRepository) GetUserByEmail(email string) (*User, error) {
	var user User
	err := r.db.Table("`user`").
		Where("email = ? AND is_deleted = 0", email).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *MySQLRepository) GetUserByID(id int64) (*User, error) {
	var user User
	err := r.db.Table("`user`").
		Where("id = ? AND is_deleted = 0", id).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *MySQLRepository) GetUserList() ([]*User, error) {
	var users []*User
	err := r.db.Table("`user`").
		Where("is_deleted = 0").
		Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *MySQLRepository) CreateUser(user *User) (int64, error) {
	if err := r.db.Table("`user`").Create(user).Error; err != nil {
		return 0, err
	}
	return user.ID, nil
}

func (r *MySQLRepository) UpdateProfile(id int64, username, email, avatar string) error {
	return r.db.Table("`user`").
		Where("id = ? AND is_deleted = 0", id).
		Updates(map[string]any{
			"username": username,
			"email":    email,
			"avatar":   avatar,
		}).Error
}

func (r *MySQLRepository) UpdatePassword(id int64, newPassword string) error {
	return r.db.Table("`user`").
		Where("id = ? AND is_deleted = 0", id).
		Update("password", newPassword).Error
}

func (r *MySQLRepository) DeleteUser(id int64) error {
	return r.db.Table("`user`").
		Where("id = ? AND is_deleted = 0", id).
		Update("is_deleted", 1).Error
}

// GetAdminUserIDs 获取所有管理员用户的 ID 列表
func (r *MySQLRepository) GetAdminUserIDs() ([]int64, error) {
	var ids []int64
	err := r.db.Table("`user`").
		Where("is_admin = 1 AND is_deleted = 0").
		Pluck("id", &ids).Error
	return ids, err
}
