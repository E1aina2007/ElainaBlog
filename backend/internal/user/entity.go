package user

import "time"

// User 用户实体
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

// UserVO 返回给前端的用户视图对象
type UserVO struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Avatar    string    `json:"avatar"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Code     string `json:"code"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UpdateProfileRequest 更新个人资料请求
type UpdateProfileRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
}

// UpdatePasswordRequest 修改密码请求
type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// DeleteUserRequest 删除用户请求
type DeleteUserRequest struct {
	UserID int64 `json:"user_id"`
}

// RefreshTokenRequest 刷新令牌请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// SendCodeRequest 发送验证码请求
type SendCodeRequest struct {
	Email string `json:"email"`
}

// ResetPasswordRequest 重置密码请求
type ResetPasswordRequest struct {
	Email       string `json:"email"`
	Code        string `json:"code"`
	NewPassword string `json:"new_password"`
}

// CreateUserParams 创建用户参数
type CreateUserParams struct {
	Username string
	Password string
	Email    string
	Avatar   string
	IsAdmin  bool
	Code     string
}

// UpdateProfileParams 更新个人资料参数
type UpdateProfileParams struct {
	UserID   int64
	Username string
	Email    string
	Avatar   string
}

// LoginParams 登录参数
type LoginParams struct {
	Email    string
	Password string
}

// LoginResult 登录结果
type LoginResult struct {
	UserID       int64
	Email        string
	AccessToken  string
	RefreshToken string
}
