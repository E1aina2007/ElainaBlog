// main.go 应用入口，负责加载配置、初始化数据库连接并启动 Gin HTTP 服务
package main

import (
	"ElainaBlog/config"
	"ElainaBlog/config/db"
	"ElainaBlog/internal/common"
	"ElainaBlog/pkg/rdb"
	"ElainaBlog/pkg/zaplogger"
	"fmt"
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

	// 3. 连接数据库（不自动迁移，由各命令自行决定）
	err = db.ConnectDB(&config.GlobalConfig.Db)
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
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
	case "migrate":
		if err := db.Migrate(); err != nil {
			log.Fatalf("数据库迁移失败: %v", err)
		}
		fmt.Println("数据库迁移完成")
	case "rollback":
		if err := db.DownMigrate(); err != nil {
			log.Fatalf("数据库回滚失败: %v", err)
		}
		fmt.Println("数据库回滚完成")
	case "initSystem":
		initSystem()
	case "runServer":
		if err := db.Migrate(); err != nil {
			log.Fatalf("数据库迁移失败: %v", err)
		}
		if err := runServer(); err != nil {
			log.Fatalf("服务器启动失败：%v", err)
		}
	default:
		log.Fatalf("未知的命令")
	}
}
