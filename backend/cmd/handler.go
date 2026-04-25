package main

import (
	"ElainaBlog/config"
	"ElainaBlog/config/db"
	"ElainaBlog/internal/router"
	"ElainaBlog/internal/user"
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
	if len(os.Args) < 3 {
		log.Fatalf("请输入管理员密码，用法: go run ./cmd initSystem <password>")
	}
	adminPassword := os.Args[2]

	userService := user.NewService(user.NewRepository(db.DBPool))
	adminUserID, err := userService.CreateUser(user.CreateUserParams{
		Username: "admin",
		Password: adminPassword,
		Email:    "admin@qq.com", // TODO: 生产环境下修改为实际的管理员邮箱
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
	zaplogger.Logger.Info("管理员创建成功", zap.Int64("userID", adminUserID))

	log.Printf("管理员创建成功，用户名: admin，密码: %s", adminPassword)
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
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // TODO：前端地址, 生产环境把 AllowOrigins 改为实际域名。
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 注册路由
	router.RouterInit(r)

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
