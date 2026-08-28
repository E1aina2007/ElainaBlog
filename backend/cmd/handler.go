// handler.go 应用启动入口：初始化 Gin 服务并监听端口
package main

import (
	"ElainaBlog/internal/article"
	"ElainaBlog/internal/auth"
	"ElainaBlog/internal/config"
	"ElainaBlog/internal/db"
	cache "ElainaBlog/internal/middleware/redis"
	"ElainaBlog/internal/router"
	"ElainaBlog/internal/upload"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/robfig/cron/v3"
)

func runServer() error {
	// 创建路由（含 gin 模式、中间件、CORS 与全部 API 路由）
	r := router.New(router.Options{
		Config: *config.GlobalConfig,
		DB:     db.DB,
		Redis:  cache.DefaultClient,
		Jwt:    auth.JwtAuth,
	})

	// 启动浏览量定时同步（每 5 分钟将 Redis 缓冲写入 MySQL）
	articleRepo := article.NewRepository(db.DB, cache.DefaultClient)
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			if _, err := articleRepo.FlushViewCounts(context.Background()); err != nil {
				log.Printf("浏览量同步失败: %v", err)
			}
		}
	}()

	// 启动孤儿图片清理定时任务
	if config.GlobalConfig.Upload.CleanupEnabled {
		imageCleanup := upload.NewImageCleanup(
			config.GlobalConfig.Upload.Path,
			config.GlobalConfig.Upload.AvatarPath,
			articleRepo,
		)

		c := cron.New()
		_, err := c.AddFunc(config.GlobalConfig.Upload.CleanupCron, func() {
			result, err := imageCleanup.CleanupOrphanImages(context.Background())
			if err != nil {
				log.Printf("孤儿图片清理失败: %v", err)
			} else if len(result.Errors) > 0 {
				log.Printf("孤儿图片清理存在 %d 个错误", len(result.Errors))
			}
		})
		if err != nil {
			log.Printf("注册孤儿图片清理定时任务失败: %v", err)
		} else {
			c.Start()
			log.Printf("孤儿图片清理定时任务已启动: %s", config.GlobalConfig.Upload.CleanupCron)
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
		log.Printf("服务器启动中: address=%s", address)
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务器异常退出: %v", err)
		}
	}()

	// 监听系统信号，优雅关停
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	log.Printf("收到退出信号，正在关闭服务器... signal=%s", sig.String())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Printf("服务器关闭失败: %v", err)
		return err
	}

	if err := db.Close(db.DB); err != nil {
		log.Printf("数据库关闭失败: %v", err)
	}
	log.Println("服务器已关闭")
	return nil
}
