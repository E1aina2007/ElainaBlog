package main

import (
	"ElainaBlog/config"
	"ElainaBlog/config/db"
	"ElainaBlog/internal/article"
	"ElainaBlog/internal/common"
	"ElainaBlog/internal/router"
	"ElainaBlog/internal/user"
	"ElainaBlog/pkg/rdb"
	"ElainaBlog/pkg/zaplogger"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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
