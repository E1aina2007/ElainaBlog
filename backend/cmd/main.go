// main.go 应用入口，负责加载配置、初始化数据库连接并启动 Gin HTTP 服务
package main

import (
	"ElainaBlog/config"
	"ElainaBlog/config/db"
	"ElainaBlog/internal/common"
	"ElainaBlog/pkg/rdb"
	"ElainaBlog/pkg/zaplogger"
	"log"
	"os"
)

func init() {
	// 1. 加载配置文件
	path := config.CheckMode()
	err := config.LoadConfigFromYml(path)
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

	err = rdb.InitRedis(&config.GlobalConfig.Redis)
	if err != nil {
		log.Fatalf("Redis初始化失败：%v", err)
	}

	// 4. 初始化 JWT 服务
	common.InitJwtAuth()
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("请输入命令")
	}
	switch os.Args[1] {
	case "initSystem":
		initSystem()
		return
	case "runServer":
		err := runServer()
		if err != nil {
			log.Fatalf("服务器启动失败：%v", err)
		}
		return
	default:
		log.Fatalf("未知的命令")
	}
}
