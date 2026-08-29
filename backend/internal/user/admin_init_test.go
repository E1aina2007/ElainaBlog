package user

import (
	"context"
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"ElainaBlog/internal/config"
)

func adminConfig(password string) config.AdminConfig {
	return config.AdminConfig{Username: "admin", Email: "admin@example.com", Password: password}
}

func TestEnsureAdmin(t *testing.T) {
	ctx := context.Background()

	t.Run("未配置密码", func(t *testing.T) {
		store := newFakeUserStore()
		svc := NewService(store, nil, nil)
		err := EnsureAdmin(ctx, svc, adminConfig(""))
		if err == nil || !strings.Contains(err.Error(), "管理员密码未配置") {
			t.Errorf("应提示密码未配置，得到 %v", err)
		}
	})

	t.Run("首次创建管理员", func(t *testing.T) {
		store := newFakeUserStore()
		svc := NewService(store, nil, nil)
		if err := EnsureAdmin(ctx, svc, adminConfig("passw0rd123")); err != nil {
			t.Fatalf("创建管理员失败: %v", err)
		}
		if len(store.created) != 1 {
			t.Fatalf("应创建一条用户记录，得到 %d", len(store.created))
		}
		created := store.created[0]
		if !created.IsAdmin {
			t.Error("创建的账号应为管理员")
		}
		if created.Username != "admin" || created.Email != "admin@example.com" {
			t.Errorf("账号信息不符: username=%s email=%s", created.Username, created.Email)
		}
		if bcrypt.CompareHashAndPassword([]byte(created.Password), []byte("passw0rd123")) != nil {
			t.Error("密码应为 bcrypt 哈希")
		}
	})

	t.Run("密码一致时跳过", func(t *testing.T) {
		store := newFakeUserStore()
		store.byUsername["admin"] = &User{ID: 1, Username: "admin", Email: "admin@example.com", IsAdmin: true, Password: hashPassword(t, "passw0rd123")}
		svc := NewService(store, nil, nil)
		if err := EnsureAdmin(ctx, svc, adminConfig("passw0rd123")); err != nil {
			t.Fatalf("不应报错: %v", err)
		}
		if len(store.created) != 0 || len(store.passwordUpdates) != 0 {
			t.Error("密码一致时不应创建或重置")
		}
	})

	t.Run("密码变更时重置", func(t *testing.T) {
		store := newFakeUserStore()
		store.byUsername["admin"] = &User{ID: 3, Username: "admin", Email: "admin@example.com", IsAdmin: true, Password: hashPassword(t, "old-passw0rd")}
		svc := NewService(store, nil, nil)
		if err := EnsureAdmin(ctx, svc, adminConfig("new-passw0rd9")); err != nil {
			t.Fatalf("重置失败: %v", err)
		}
		hashed, ok := store.passwordUpdates[3]
		if !ok {
			t.Fatal("应调用仓库层重置密码")
		}
		if bcrypt.CompareHashAndPassword([]byte(hashed), []byte("new-passw0rd9")) != nil {
			t.Error("重置后的密码应为新密码的 bcrypt 哈希")
		}
	})

	t.Run("用户名被普通用户占用时跳过", func(t *testing.T) {
		store := newFakeUserStore()
		store.byUsername["admin"] = &User{ID: 5, Username: "admin", Email: "someone@example.com", IsAdmin: false, Password: hashPassword(t, "whatever123")}
		svc := NewService(store, nil, nil)
		if err := EnsureAdmin(ctx, svc, adminConfig("passw0rd123")); err != nil {
			t.Fatalf("应静默跳过: %v", err)
		}
		if len(store.passwordUpdates) != 0 {
			t.Error("不应重置普通用户的密码")
		}
	})

	t.Run("邮箱被占用时跳过创建", func(t *testing.T) {
		store := newFakeUserStore()
		store.byEmail["admin@example.com"] = &User{ID: 6, Username: "someone", Email: "admin@example.com", IsAdmin: false, Password: hashPassword(t, "whatever123")}
		svc := NewService(store, nil, nil)
		if err := EnsureAdmin(ctx, svc, adminConfig("passw0rd123")); err != nil {
			t.Fatalf("应静默跳过: %v", err)
		}
		if len(store.created) != 0 {
			t.Error("不应创建新用户")
		}
	})
}
