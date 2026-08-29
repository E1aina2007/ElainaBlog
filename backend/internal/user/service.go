package user

import (
	"ElainaBlog/internal/config"
	"ElainaBlog/internal/auth"
	"ElainaBlog/internal/mail"
	cache "ElainaBlog/internal/middleware/redis"
	"ElainaBlog/internal/util/verifycode"
	"context"
	"crypto/subtle"
	"errors"
	"fmt"
	"strings"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"golang.org/x/crypto/bcrypt"
)

// UserStore 用户数据的窄接口（消费者侧定义，参照 comment 包模式），
// 便于 service 层单元测试时以假实现替换；*Repository 天然满足此接口
type UserStore interface {
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, id int64) (*User, error)
	GetUserList(ctx context.Context) ([]*User, error)
	CreateUser(ctx context.Context, user *User) (int64, error)
	UpdateProfile(ctx context.Context, id int64, username, email, avatar string) error
	UpdatePassword(ctx context.Context, id int64, newPassword string) error
	DeleteUser(ctx context.Context, id int64) error
}

type Service struct {
	repo     UserStore
	rdb      *goredis.Client    // 可选，用于验证码存储
	tokenMgr auth.TokenManager  // 可选，用于 JWT 签发
}

const (
	cacheKeyAdminPrefix = "cache:user:admin:"
	cacheTTLAdmin       = time.Hour
)

func NewService(repo UserStore, redis *goredis.Client, tokenMgr auth.TokenManager) *Service {
	return &Service{repo: repo, rdb: redis, tokenMgr: tokenMgr}
}

var (
	ErrInvalidParams      = errors.New("无效的参数")
	ErrInvalidLoginParams = errors.New("邮箱或密码不能为空")
	ErrUsernameExists     = errors.New("用户名已存在")
	ErrEmailExists        = errors.New("邮箱已注册")
	ErrUserNotFound       = errors.New("用户不存在")
	ErrPasswordMismatch   = errors.New("密码错误")
	ErrSamePassword       = errors.New("新密码不能与旧密码相同")
	ErrForbidden          = errors.New("无权限执行此操作")
	ErrDBNotInitialized   = errors.New("数据库未初始化")
	ErrResendTooFrequent  = errors.New("发送过于频繁，请稍后再试")
	ErrCodeExpired        = errors.New("验证码已过期或不存在")
	ErrCodeMismatch       = errors.New("验证码错误")
)

func (s *Service) CreateUser(ctx context.Context, params CreateUserParams) (int64, error) {
	if s == nil || s.repo == nil {
		return 0, ErrDBNotInitialized
	}

	username := strings.TrimSpace(params.Username)
	password := strings.TrimSpace(params.Password)
	email := strings.TrimSpace(params.Email)
	if username == "" || password == "" || email == "" {
		return 0, ErrInvalidParams
	}

	// 正则校验邮箱、用户名、密码格式
	if err := ValidateEmail(email); err != nil {
		return 0, err
	}
	if err := ValidateUsername(username); err != nil {
		return 0, err
	}
	if err := ValidatePassword(password); err != nil {
		return 0, err
	}

	// 非管理员创建（即普通注册）需要校验验证码
	if !params.IsAdmin {
		code := strings.TrimSpace(params.Code)
		if code == "" {
			return 0, ErrInvalidParams
		}
		storedCode, err := cache.GetVerificationCode(s.rdb, email)
		if err != nil {
			return 0, ErrCodeExpired
		}
		// 常数时间比对，避免时序侧信道
		if subtle.ConstantTimeCompare([]byte(storedCode), []byte(code)) != 1 {
			return 0, ErrCodeMismatch
		}
		_ = cache.DeleteVerificationCode(s.rdb, email)
	}

	// 检查用户名是否已存在
	existing, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, err
	}
	if existing != nil {
		return 0, ErrUsernameExists
	}

	// 检查邮箱是否已注册
	existing, err = s.repo.GetUserByEmail(ctx, email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, err
	}
	if existing != nil {
		return 0, ErrEmailExists
	}

	// bcrypt 哈希密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	return s.repo.CreateUser(ctx, &User{
		Username: username,
		Password: string(hashedPassword),
		Email:    email,
		Avatar:   strings.TrimSpace(params.Avatar),
		IsAdmin:  params.IsAdmin,
	})
}

func (s *Service) Login(ctx context.Context, params LoginParams) (*LoginResult, error) {
	if s == nil || s.repo == nil {
		return nil, ErrDBNotInitialized
	}

	email := strings.TrimSpace(params.Email)
	password := strings.TrimSpace(params.Password)

	if email == "" || password == "" {
		return nil, ErrInvalidLoginParams
	}

	// 校验邮箱格式
	if err := ValidateEmail(email); err != nil {
		return nil, err
	}

	// 通过邮箱查询用户
	u, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// bcrypt 比对密码
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return nil, ErrPasswordMismatch
	}

	// 签发 access token
	accessToken, err := s.tokenMgr.GenerateAccessToken(u.ID)
	if err != nil {
		return nil, err
	}

	// 签发 refresh token
	refreshToken, err := s.tokenMgr.GenerateRefreshToken(u.ID)
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		UserID:       u.ID,
		Email:        u.Email,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (*User, error) {
	if s == nil || s.repo == nil {
		return nil, ErrDBNotInitialized
	}
	if id <= 0 {
		return nil, ErrInvalidParams
	}

	u, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return u, nil
}

func (s *Service) GetList(ctx context.Context) ([]*User, error) {
	if s == nil || s.repo == nil {
		return nil, ErrDBNotInitialized
	}
	return s.repo.GetUserList(ctx)
}

func (s *Service) UpdateProfile(ctx context.Context, params UpdateProfileParams) error {
	if s == nil || s.repo == nil {
		return ErrDBNotInitialized
	}

	username := strings.TrimSpace(params.Username)
	email := strings.TrimSpace(params.Email)
	avatar := strings.TrimSpace(params.Avatar)

	if params.UserID <= 0 || username == "" || email == "" {
		return ErrInvalidParams
	}

	// 正则校验用户名、邮箱格式
	if err := ValidateUsername(username); err != nil {
		return err
	}
	if err := ValidateEmail(email); err != nil {
		return err
	}

	// 检查用户是否存在
	_, err := s.repo.GetUserByID(ctx, params.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	// 检查用户名是否被其他人占用
	existing, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if existing != nil && existing.ID != params.UserID {
		return ErrUsernameExists
	}

	// 检查邮箱是否被其他人占用
	existing, err = s.repo.GetUserByEmail(ctx, email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if existing != nil && existing.ID != params.UserID {
		return ErrEmailExists
	}

	return s.repo.UpdateProfile(ctx, params.UserID, username, email, avatar)
}

func (s *Service) UpdatePassword(ctx context.Context, userID int64, oldPassword, newPassword string) error {
	if s == nil || s.repo == nil {
		return ErrDBNotInitialized
	}

	oldPassword = strings.TrimSpace(oldPassword)
	newPassword = strings.TrimSpace(newPassword)
	if userID <= 0 || oldPassword == "" || newPassword == "" {
		return ErrInvalidParams
	}

	// 校验新密码格式
	if err := ValidatePassword(newPassword); err != nil {
		return err
	}

	// 查询用户
	u, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(oldPassword)); err != nil {
		return ErrPasswordMismatch
	}

	// 新旧密码不能相同
	if oldPassword == newPassword {
		return ErrSamePassword
	}

	// 哈希新密码并更新
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.repo.UpdatePassword(ctx, userID, string(hashedPassword))
}

func (s *Service) DeleteUser(ctx context.Context, operatorID, targetID int64) error {
	if s == nil || s.repo == nil {
		return ErrDBNotInitialized
	}
	if operatorID <= 0 || targetID <= 0 {
		return ErrInvalidParams
	}

	// 校验操作者是否为管理员
	operator, err := s.repo.GetUserByID(ctx, operatorID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}
	if !operator.IsAdmin {
		return ErrForbidden
	}

	// 校验目标用户是否存在
	_, err = s.repo.GetUserByID(ctx, targetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	return s.repo.DeleteUser(ctx, targetID)
}

// CheckIsAdmin 检查用户是否为管理员（优先查 Redis 缓存）
func (s *Service) CheckIsAdmin(ctx context.Context, userID int64) (bool, error) {
	if s == nil || s.repo == nil {
		return false, ErrDBNotInitialized
	}

	// 尝试从 Redis 读取
	if s.rdb != nil {
		val, err := s.rdb.Get(ctx, cacheKeyAdminPrefix+fmt.Sprintf("%d", userID)).Result()
		if err == nil {
			return val == "1", nil
		}
	}

	// 缓存未命中，查库
	u, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, ErrUserNotFound
		}
		return false, err
	}

	// 写入缓存
	if s.rdb != nil {
		v := "0"
		if u.IsAdmin {
			v = "1"
		}
		s.rdb.Set(ctx, cacheKeyAdminPrefix+fmt.Sprintf("%d", userID), v, cacheTTLAdmin)
	}

	return u.IsAdmin, nil
}

func (s *Service) GetUserEmailByID(ctx context.Context, userID int64) (string, error) {
	if s == nil || s.repo == nil {
		return "", ErrDBNotInitialized
	}
	u, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrUserNotFound
		}
		return "", err
	}
	return u.Email, nil
}

func (s *Service) GetUsernameByID(ctx context.Context, userID int64) (string, error) {
	if s == nil || s.repo == nil {
		return "", ErrDBNotInitialized
	}
	u, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrUserNotFound
		}
		return "", err
	}
	return u.Username, nil
}

func (s *Service) SendVerificationCode(ctx context.Context, email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return ErrInvalidParams
	}

	// 校验邮箱格式
	if err := ValidateEmail(email); err != nil {
		return err
	}

	limited, err := cache.IsDuringInterval(s.rdb, email)
	if err != nil {
		return err
	}
	if limited {
		return ErrResendTooFrequent
	}

	cfg := config.GlobalConfig.Verification
	code := verifycode.GenerateCode(cfg.CodeLength)

	expiry := time.Duration(cfg.ExpireTime) * time.Second
	interval := time.Duration(cfg.ResendInterval) * time.Second
	if err := cache.SetVerificationCode(s.rdb, email, code, expiry, interval); err != nil {
		return err
	}

	return mail.SendVerificationCode(config.GlobalConfig.Smtp, cfg.ExpireTime, email, code)
}

func (s *Service) ResetPassword(ctx context.Context, email, code, newPassword string) error {
	if s == nil || s.repo == nil {
		return ErrDBNotInitialized
	}

	email = strings.TrimSpace(email)
	code = strings.TrimSpace(code)
	newPassword = strings.TrimSpace(newPassword)

	if email == "" || code == "" || newPassword == "" {
		return ErrInvalidParams
	}

	// 校验邮箱格式
	if err := ValidateEmail(email); err != nil {
		return err
	}

	// 校验新密码格式
	if err := ValidatePassword(newPassword); err != nil {
		return err
	}

	// 验证验证码
	storedCode, err := cache.GetVerificationCode(s.rdb, email)
	if err != nil {
		return ErrCodeExpired
	}
	// 常数时间比对，避免时序侧信道
	if subtle.ConstantTimeCompare([]byte(storedCode), []byte(code)) != 1 {
		return ErrCodeMismatch
	}

	// 查询用户
	u, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	// 哈希新密码并更新
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 删除已使用的验证码
	_ = cache.DeleteVerificationCode(s.rdb, email)

	return s.repo.UpdatePassword(ctx, u.ID, string(hashedPassword))
}

