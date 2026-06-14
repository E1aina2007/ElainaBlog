package main

import (
	"ElainaBlog/config"
	"ElainaBlog/config/db"
	"ElainaBlog/internal/article"
	"ElainaBlog/internal/common"
	"ElainaBlog/internal/router"
	"ElainaBlog/internal/upload"
	"ElainaBlog/internal/user"
	"ElainaBlog/pkg/rdb"
	"ElainaBlog/pkg/util"
	"ElainaBlog/pkg/zaplogger"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

func initSystem() {
	adminCfg := config.GlobalConfig.Admin

	// 优先使用命令行参数的密码，其次使用 yaml 配置中的密码
	var adminPassword string
	if len(os.Args) >= 3 {
		adminPassword = os.Args[2]
	} else if adminCfg.Password != "" {
		adminPassword = adminCfg.Password
	} else {
		log.Fatalf("管理员密码未配置，请通过命令行参数传入（go run ./cmd initSystem <password>）或在 yaml 的 admin.password 中设置")
	}

	adminUsername := adminCfg.Username
	adminEmail := adminCfg.Email
	if adminUsername == "" {
		adminUsername = "admin"
	}
	if adminEmail == "" {
		adminEmail = "admin@admin.com"
	}

	userService := user.NewService(user.NewRepository(db.DBPool), rdb.DefaultClient, common.JwtAuth)
	adminUserID, err := userService.CreateUser(user.CreateUserParams{
		Username: adminUsername,
		Password: adminPassword,
		Email:    adminEmail,
		Avatar:   "",
		IsAdmin:  true,
	})
	if err != nil {
		if err == user.ErrUsernameExists || err == user.ErrEmailExists {
			zaplogger.Logger.Info("管理员账号已存在，跳过创建")
			return
		}
		zaplogger.Logger.Fatal("创建管理员失败", zap.Error(err))
	}
	zaplogger.Logger.Info("管理员创建成功，请妥善保管密码",
		zap.Int64("userID", adminUserID),
		zap.String("email", adminEmail),
		zap.String("username", adminUsername),
	)
}

// migrateAvatars 迁移头像文件：将邮箱命名改为哈希命名
func migrateAvatars() {
	uploadDir := config.GlobalConfig.Upload.AvatarPath
	if uploadDir == "" {
		uploadDir = "uploads/avatars"
	}

	userRepo := user.NewRepository(db.DBPool)
	users, err := userRepo.GetUserList()
	if err != nil {
		log.Fatalf("获取用户列表失败: %v", err)
	}

	migrated := 0
	for _, u := range users {
		if u.Avatar == "" {
			continue
		}

		// 从 avatar URL 中提取文件名
		// URL 格式: /uploads/avatars/xxx.jpg
		parts := strings.Split(u.Avatar, "/")
		if len(parts) == 0 {
			continue
		}
		oldFileName := parts[len(parts)-1]
		ext := filepath.Ext(oldFileName)

		// 计算新文件名
		newFileName := util.EmailToAvatarHash(u.Email) + ext
		if oldFileName == newFileName {
			continue // 已经是哈希命名，跳过
		}

		// 重命名文件
		oldPath := filepath.Join(uploadDir, oldFileName)
		newPath := filepath.Join(uploadDir, newFileName)

		if _, err := os.Stat(oldPath); os.IsNotExist(err) {
			fmt.Printf("跳过用户 %s: 头像文件不存在 %s\n", u.Email, oldPath)
			continue
		}

		if err := os.Rename(oldPath, newPath); err != nil {
			fmt.Printf("迁移用户 %s 头像失败: %v\n", u.Email, err)
			continue
		}

		// 更新数据库中的 avatar URL
		newAvatar := "/uploads/avatars/" + newFileName
		if err := userRepo.UpdateProfile(u.ID, u.Username, u.Email, newAvatar); err != nil {
			fmt.Printf("更新用户 %s 头像URL失败: %v\n", u.Email, err)
			continue
		}

		migrated++
		fmt.Printf("迁移用户 %s 头像成功: %s -> %s\n", u.Email, oldFileName, newFileName)
	}

	fmt.Printf("\n迁移完成，共迁移 %d 个头像\n", migrated)
}

func runServer() error {
	// 根据运行环境设置 Gin 模式
	switch config.GlobalConfig.Server.Env {
	case "production":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	case "debug":
		gin.SetMode(gin.DebugMode)
	default:
		zaplogger.Logger.Fatal("未知的运行环境", zap.String("env", config.GlobalConfig.Server.Env))
	}

	// 创建 Gin 引擎并注册中间件
	r := gin.New()
	r.Use(gin.Recovery())

	corsOrigins := config.GlobalConfig.Server.CorsOrigins
	if len(corsOrigins) == 0 {
		corsOrigins = []string{"http://localhost:5173"}
	}
	r.Use(cors.New(cors.Config{
		AllowOrigins:     corsOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 注册路由
	router.RouterInit(r)

	// 启动浏览量定时同步（每 5 分钟将 Redis 缓冲写入 MySQL）
	articleRepo := article.NewRepository(db.DBPool, rdb.DefaultClient)
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			flushed, err := articleRepo.FlushViewCounts()
			if err != nil {
				zaplogger.Logger.Error("浏览量同步失败", zap.Error(err))
			} else if flushed > 0 {
				zaplogger.Logger.Info("浏览量同步完成", zap.Int("articles", flushed))
			}
		}
	}()

	// 启动孤儿图片清理定时任务
	if config.GlobalConfig.Upload.CleanupEnabled {
		imageCleanup := upload.NewImageCleanup(
			config.GlobalConfig.Upload.Path,
			config.GlobalConfig.Upload.AvatarPath,
			articleRepo,
			zaplogger.Logger,
		)

		c := cron.New()
		_, err := c.AddFunc(config.GlobalConfig.Upload.CleanupCron, func() {
			zaplogger.Logger.Info("开始执行孤儿图片清理任务")
			result, err := imageCleanup.CleanupOrphanImages()
			if err != nil {
				zaplogger.Logger.Error("孤儿图片清理失败", zap.Error(err))
			} else {
				zaplogger.Logger.Info("孤儿图片清理完成",
					zap.Int("scanned", result.ScannedFiles),
					zap.Int("referenced", result.ReferencedFiles),
					zap.Int("deleted", result.DeletedFiles),
				)
				if len(result.Errors) > 0 {
					zaplogger.Logger.Warn("清理过程中有错误", zap.Int("error_count", len(result.Errors)))
				}
			}
		})
		if err != nil {
			zaplogger.Logger.Error("注册孤儿图片清理定时任务失败", zap.Error(err))
		} else {
			c.Start()
			zaplogger.Logger.Info("孤儿图片清理定时任务已启动", zap.String("cron", config.GlobalConfig.Upload.CleanupCron))
		}
	}

	// 初始化服务器
	address := config.GlobalConfig.Server.GetAddress()
	s := &http.Server{
		Addr:    address,
		Handler: r,
	}

	// 在独立 goroutine 中启动 HTTP 服务
	go func() {
		zaplogger.Logger.Info("服务器启动中", zap.String("address", address), zap.String("env", config.GlobalConfig.Server.Env))
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zaplogger.Logger.Fatal("服务器异常退出", zap.Error(err))
		}
	}()

	// 监听系统信号，优雅关停
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	zaplogger.Logger.Info("收到退出信号，正在关闭服务器...", zap.String("signal", sig.String()))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		zaplogger.Logger.Error("服务器关闭失败", zap.Error(err))
		return err
	}

	zaplogger.Logger.Info("服务器已关闭")
	return nil
}
