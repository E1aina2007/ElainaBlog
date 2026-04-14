// main.go 应用入口，负责加载配置、初始化数据库连接并启动 Gin HTTP 服务
package main

import (
	"ElainaWeb/config"
	"ElainaWeb/config/db"
	"ElainaWeb/internal/server"
	"ElainaWeb/pkg/zaplogger"
	"log"

	"go.uber.org/zap"
)

func init() {
	// 1. 加载配置文件
	cfgPath := "../config.yaml"
	err := config.LoadConfigFromYml(cfgPath)
	if err != nil {
		log.Fatalf("配置文件加载失败: %v", err)
	}

	// 2. 初始化 Zap 日志
	zaplogger.Logger = zaplogger.InitLogger()

	// 3. 初始化数据库
	err = db.InitDB(&config.GlobalConfig.Db)
	if err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}
}

func main() {
	err := server.RunServer()
	if err != nil {
		zaplogger.Logger.Fatal("服务器启动失败: ", zap.Error(err))
	}
}
