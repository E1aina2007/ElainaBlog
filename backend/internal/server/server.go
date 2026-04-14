package server

import (
	"ElainaWeb/config"
	"ElainaWeb/internal/router"
	"ElainaWeb/pkg/zaplogger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type server interface {
	ListenAndServe() error
}

func RunServer() error {
	// 根据运行环境设置 Gin 模式
	switch config.GlobalConfig.Server.Env {
	case "production":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	// 创建 Gin 引擎并注册中间件
	r := gin.New()
	r.Use(gin.Recovery())

	// 注册路由
	router.RouterInit(r)

	// 初始化服务器并启动
	address := config.GlobalConfig.Server.GetAddress()
	s := initServer(address, r)
	zaplogger.Logger.Info("服务器启动中", zap.String("address", address), zap.String("env", config.GlobalConfig.Server.Env))

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
