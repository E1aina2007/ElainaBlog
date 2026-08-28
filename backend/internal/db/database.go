// database.go 初始化 MySQL 数据库连接，提供全局 *gorm.DB 实例
package db

import (
	"ElainaBlog/internal/config"
	"fmt"
	"time"

	sqldriver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB // DB 全局数据库实例

// ConnectDB 连接数据库并初始化 *gorm.DB
func ConnectDB(dbConfig *config.DbConfig) error {
	// 参照 GoFeed：使用 go-sql-driver/mysql 构造 DSN
	mc := sqldriver.NewConfig()
	mc.User = dbConfig.Username
	mc.Passwd = dbConfig.Password
	mc.Net = "tcp"
	mc.Addr = fmt.Sprintf("%s:%d", dbConfig.Host, dbConfig.Port)
	mc.DBName = dbConfig.DBName
	mc.ParseTime = true
	mc.Loc = time.Local
	mc.Params = map[string]string{"charset": "utf8mb4", "clientFoundRows": "true"}

	gormDB, err := gorm.Open(mysql.Open(mc.FormatDSN()), &gorm.Config{})
	if err != nil {
		return err
	}

	// 获取底层 *sql.DB 配置连接池
	sqlDB, err := gormDB.DB()
	if err != nil {
		return err
	}
	if err = sqlDB.Ping(); err != nil {
		return err
	}

	DB = gormDB
	return nil
}

// Close 关闭数据库连接
func Close(gormDB *gorm.DB) error {
	sqlDB, err := gormDB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
