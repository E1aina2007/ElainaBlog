// main.go 应用入口，负责加载配置、初始化数据库连接并启动 Gin HTTP 服务
package main

import (
	"ElainaWeb/config"
	"ElainaWeb/global"
	"ElainaWeb/internal/model"
	"ElainaWeb/internal/server"
	"ElainaWeb/pkg/db"
	"ElainaWeb/pkg/zaplogger"
)

func main() {
	global.Config = config.ConfigInit()
	global.Logger = zaplogger.InitLogger()
	global.DB = db.InitMysql()
	global.Redis = model.ConnectRedis()

	defer global.Redis.Close()

	server.RunServer()
}
