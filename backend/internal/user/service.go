package user

import (
	"ElainaBlog/config"
	"ElainaBlog/internal/common"
	"ElainaBlog/pkg/mail"
	"ElainaBlog/pkg/rdb"
	"ElainaBlog/pkg/util"
	"database/sql"
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

type CreateUserParams struct {
	Username string
	Password string
	Email    string
	Avatar   string
	IsAdmin  bool
	Code     string
}

type UpdateProfileParams struct {
	UserID   int64
	Username string
	Email    string
	Avatar   string
}

type LoginParams struct {
	Email    string
	Password string
}

type LoginResult struct {
	UserID       int64
	Email        string
	AccessToken  string
	RefreshToken string
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

func (s *Service) CreateUser(params CreateUserParams) (int64, error) {
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
		storedCode, err := rdb.GetVerificationCode(email)
		if err != nil {
			return 0, ErrCodeExpired
		}
		if storedCode != code {
			return 0, ErrCodeMismatch
		}
		_ = rdb.DeleteVerificationCode(email)
	}

	// 检查用户名是否已存在
	existing, err := s.repo.GetUserByUsername(username)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}
	if existing != nil {
		return 0, ErrUsernameExists
	}

	// 检查邮箱是否已注册
	existing, err = s.repo.GetUserByEmail(email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
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

	return s.repo.CreateUser(&User{
		Username: username,
		Password: string(hashedPassword),
		Email:    email,
		Avatar:   strings.TrimSpace(params.Avatar),
		IsAdmin:  params.IsAdmin,
	})
}

func (s *Service) Login(params LoginParams) (*LoginResult, error) {
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
	u, err := s.repo.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// bcrypt 比对密码
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return nil, ErrPasswordMismatch
	}

	// 签发 access token
	accessToken, err := common.JwtAuth.GenerateAccessToken(u.ID)
	if err != nil {
		return nil, err
	}

	// 签发 refresh token
	refreshToken, err := common.JwtAuth.GenerateRefreshToken(u.ID)
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

func (s *Service) GetByID(id int64) (*User, error) {
	if s == nil || s.repo == nil {
		return nil, ErrDBNotInitialized
	}
	if id <= 0 {
		return nil, ErrInvalidParams
	}

	u, err := s.repo.GetUserByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return u, nil
}

func (s *Service) GetList() ([]*User, error) {
	if s == nil || s.repo == nil {
		return nil, ErrDBNotInitialized
	}
	return s.repo.GetUserList()
}

func (s *Service) UpdateProfile(params UpdateProfileParams) error {
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
	_, err := s.repo.GetUserByID(params.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}
		return err
	}

	// 检查用户名是否被其他人占用
	existing, err := s.repo.GetUserByUsername(username)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	if existing != nil && existing.ID != params.UserID {
		return ErrUsernameExists
	}

	// 检查邮箱是否被其他人占用
	existing, err = s.repo.GetUserByEmail(email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	if existing != nil && existing.ID != params.UserID {
		return ErrEmailExists
	}

	return s.repo.UpdateProfile(params.UserID, username, email, avatar)
}

func (s *Service) UpdatePassword(userID int64, oldPassword, newPassword string) error {
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
	u, err := s.repo.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
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

	return s.repo.UpdatePassword(userID, string(hashedPassword))
}

func (s *Service) DeleteUser(operatorID, targetID int64) error {
	if s == nil || s.repo == nil {
		return ErrDBNotInitialized
	}
	if operatorID <= 0 || targetID <= 0 {
		return ErrInvalidParams
	}

	// 校验操作者是否为管理员
	operator, err := s.repo.GetUserByID(operatorID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}
		return err
	}
	if !operator.IsAdmin {
		return ErrForbidden
	}

	// 校验目标用户是否存在
	_, err = s.repo.GetUserByID(targetID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}
		return err
	}

	return s.repo.DeleteUser(targetID)
}

func (s *Service) CheckIsAdmin(userID int64) (bool, error) {
	if s == nil || s.repo == nil {
		return false, ErrDBNotInitialized
	}
	u, err := s.repo.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, ErrUserNotFound
		}
		return false, err
	}
	return u.IsAdmin, nil
}

func (s *Service) SendVerificationCode(email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return ErrInvalidParams
	}

	// 校验邮箱格式
	if err := ValidateEmail(email); err != nil {
		return err
	}

	limited, err := rdb.IsDuringInterval(email)
	if err != nil {
		return err
	}
	if limited {
		return ErrResendTooFrequent
	}

	cfg := config.GlobalConfig.Verification
	code := util.GenerateCode(cfg.CodeLength)

	expiry := time.Duration(cfg.ExpireTime) * time.Second
	interval := time.Duration(cfg.ResendInterval) * time.Second
	if err := rdb.SetVerificationCode(email, code, expiry, interval); err != nil {
		return err
	}

	return mail.SendVerificationCode(email, code)
}

// GetAuthorInfo 获取管理员作者信息（从 site_config 读取）
func (s *Service) GetAuthorInfo() (*AuthorInfoVO, error) {
	// 从 site_config 表读取作者信息配置
	// 这里通过 site 模块获取，但 service 层不直接依赖其他模块的 repository
	// 暂时返回一个占位实现，实际应与 site 模块协作
	return nil, errors.New("请通过 site 模块获取作者信息")
}

// GetAuthorStats 获取作者统计数据
func (s *Service) GetAuthorStats(userID int64) (*AuthorStatsVO, error) {
	if s == nil || s.repo == nil {
		return nil, ErrDBNotInitialized
	}
	if userID <= 0 {
		return nil, ErrInvalidParams
	}

	// 校验用户是否存在且为管理员（只有管理员才是作者）
	u, err := s.repo.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	if !u.IsAdmin {
		return nil, ErrForbidden
	}

	return s.repo.GetAuthorStats(userID)
}
