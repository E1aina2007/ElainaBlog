// admin_init.go 启动时初始化管理员账号：不存在则按配置创建，配置密码变更时自动重置
package user

import (
	"ElainaBlog/internal/config"
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// EnsureAdmin 确保管理员账号存在，密码以配置 admin.password 为准：
//   - 账号不存在时创建管理员
//   - 账号存在且配置密码与当前密码不一致时重置密码
func EnsureAdmin(ctx context.Context, svc *Service, cfg config.AdminConfig) error {
	username := strings.TrimSpace(cfg.Username)
	if username == "" {
		username = "admin"
	}
	email := strings.TrimSpace(cfg.Email)
	if email == "" {
		email = "admin@admin.com"
	}

	password := cfg.Password
	if password == "" {
		return errors.New("管理员密码未配置，请在 config 的 admin.password 中设置")
	}
	if err := ValidatePassword(password); err != nil {
		return fmt.Errorf("管理员密码格式无效: %w", err)
	}

	// 账号已存在：密码不一致时按配置重置
	if u, err := svc.repo.GetUserByUsername(ctx, username); err == nil && u != nil {
		if !u.IsAdmin {
			log.Printf("用户名 %s 已被普通用户占用，跳过管理员初始化", username)
			return nil
		}
		if bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil {
			log.Println("管理员账号已存在，跳过创建")
			return nil
		}

		hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		if err := svc.repo.UpdatePassword(ctx, u.ID, string(hashed)); err != nil {
			return err
		}
		log.Printf("管理员密码已按配置重置: username=%s", username)
		return nil
	}

	// 邮箱已被占用：跳过创建
	if u, err := svc.repo.GetUserByEmail(ctx, email); err == nil && u != nil {
		log.Printf("管理员邮箱已被占用，跳过创建: email=%s", email)
		return nil
	}

	adminUserID, err := svc.CreateUser(ctx, CreateUserParams{
		Username: username,
		Password: password,
		Email:    email,
		IsAdmin:  true,
	})
	if err != nil {
		if errors.Is(err, ErrUsernameExists) || errors.Is(err, ErrEmailExists) {
			log.Println("管理员账号已存在，跳过创建")
			return nil
		}
		return err
	}

	log.Printf("管理员创建成功，请妥善保管密码: userID=%d, email=%s, username=%s", adminUserID, email, username)
	return nil
}
