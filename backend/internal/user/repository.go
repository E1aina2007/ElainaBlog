package user

import (
	"context"

	"gorm.io/gorm"
)

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

// Repository 使用 GORM 存储用户数据。
type Repository struct {
	db *gorm.DB
}

// NewRepository 创建用户仓储实例。
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).Table("user").
		Where("username = ? AND is_deleted = 0", username).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).Table("user").
		Where("email = ? AND is_deleted = 0", email).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetUserByID(ctx context.Context, id int64) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).Table("user").
		Where("id = ? AND is_deleted = 0", id).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetUserList(ctx context.Context) ([]*User, error) {
	var users []*User
	err := r.db.WithContext(ctx).Table("user").
		Where("is_deleted = 0").
		Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *Repository) CreateUser(ctx context.Context, user *User) (int64, error) {
	if err := r.db.WithContext(ctx).Table("user").Create(user).Error; err != nil {
		return 0, err
	}
	return user.ID, nil
}

func (r *Repository) UpdateProfile(ctx context.Context, id int64, username, email, avatar string) error {
	return r.db.WithContext(ctx).Table("user").
		Where("id = ? AND is_deleted = 0", id).
		Updates(map[string]any{
			"username": username,
			"email":    email,
			"avatar":   avatar,
		}).Error
}

func (r *Repository) UpdatePassword(ctx context.Context, id int64, newPassword string) error {
	return r.db.WithContext(ctx).Table("user").
		Where("id = ? AND is_deleted = 0", id).
		Update("password", newPassword).Error
}

func (r *Repository) DeleteUser(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Table("user").
		Where("id = ? AND is_deleted = 0", id).
		Update("is_deleted", 1).Error
}

// GetAdminUserIDs 获取所有管理员用户的 ID 列表
func (r *Repository) GetAdminUserIDs(ctx context.Context) ([]int64, error) {
	var ids []int64
	err := r.db.WithContext(ctx).Table("user").
		Where("is_admin = 1 AND is_deleted = 0").
		Pluck("id", &ids).Error
	return ids, err
}
