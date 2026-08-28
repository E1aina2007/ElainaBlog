// main.go 应用入口，负责加载配置、初始化数据库连接并启动 Gin HTTP 服务
package main

import (
	"ElainaBlog/internal/auth"
	"ElainaBlog/internal/config"
	"ElainaBlog/internal/db"
	cache "ElainaBlog/internal/middleware/redis"
	"ElainaBlog/internal/user"
	"context"
	"log"
)

func init() {
	// 1. 加载配置文件（含 .env 与 CONFIG_PATH 选择，见 config.LoadConfig）
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("配置文件加载失败: %v", err)
	}
	cfg := config.GlobalConfig

	// 2. 初始化数据库连接
	if err := db.ConnectDB(&cfg.Db); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	if err := cache.InitRedis(&cfg.Redis); err != nil {
		log.Fatalf("Redis初始化失败：%v", err)
	}

	// 3. 初始化 JWT 服务
	auth.InitJwtAuth()

	// 4. 初始化管理员（不存在则创建，配置密码变更时重置）
	userService := user.NewService(user.NewRepository(db.DB), cache.DefaultClient, auth.JwtAuth)
	if err := user.EnsureAdmin(context.Background(), userService, cfg.Admin); err != nil {
		log.Fatalf("初始化管理员失败: %v", err)
	}

}

func main() {
	if err := runServer(); err != nil {
		log.Fatalf("服务器启动失败：%v", err)
	}
}
