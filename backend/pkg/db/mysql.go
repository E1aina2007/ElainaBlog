// mysql.go 初始化 MySQL 数据库连接 (GORM)，提供全局 DB 实例并执行 AutoMigrate 建表
package db

import (
	"ElainaWeb/global"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitMysql() *gorm.DB {
	mysqlCfg := global.Config.Mysql

	// 使用配置中的日志级别初始化 GORM
	db, err := gorm.Open(mysql.Open(mysqlCfg.GetDSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.LogLevel(mysqlCfg.GetLogLevel())),
	})
	if err != nil {
		global.Logger.Error("MySQL 连接失败", zap.Error(err))
		panic("MySQL 连接失败")
	}

	// 设置连接池参数
	sqlDB, err2 := db.DB()
	if err2 != nil {
		global.Logger.Error("获取底层 sql.DB 实例失败", zap.Error(err2))
		panic("获取底层 sql.DB 实例失败")
	}
	sqlDB.SetMaxIdleConns(mysqlCfg.MaxIdleConns) // 最大空闲连接数
	sqlDB.SetMaxOpenConns(mysqlCfg.MaxOpenConns) // 最大打开连接数

	return db
}
