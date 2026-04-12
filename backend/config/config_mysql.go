package config

import (
	"strconv"
	"strings"

	"gorm.io/gorm/logger"
)

type MySQLConfig struct {
	Host         string `yaml:"host" json:"host"`                     // 数据库主机
	Port         int    `yaml:"port" json:"port"`                     // 数据库端口
	Config       string `yaml:"config" json:"config"`                 // 数据库配置
	Username     string `yaml:"username" json:"username"`             // 数据库用户
	Password     string `yaml:"password" json:"password"`             // 数据库密码
	DBName       string `yaml:"db_name" json:"db_name"`               // 数据库名
	MaxIdleConns int    `yaml:"max_idle_conns" json:"max_idle_conns"` // 最大空闲连接数
	MaxOpenConns int    `yaml:"max_open_conns" json:"max_open_conns"` // 最大打开连接数
	LogMode      string `yaml:"log_mode" json:"log_mode"`             // 日志模式
}

// GetDSN 返回数据库连接字符串
// 格式: username:password@tcp(host:port)/dbname?config
// Config 字段通常包含以下参数（以 & 分隔）:
//   - charset=utf8mb4       : 使用 utf8mb4 字符集，支持存储 emoji 等完整 Unicode 字符
//   - parseTime=True        : 将 MySQL 的 DATE/DATETIME 类型自动解析为 Go 的 time.Time，
//                              而非默认的 []byte/string；GORM 依赖此参数正确处理时间字段
//   - loc=Local             : 设置时区为系统本地时区，确保时间值与服务器时区一致
func (m MySQLConfig) GetDSN() string {
	return m.Username + ":" + m.Password + "@tcp(" + m.Host + ":" + strconv.Itoa(m.Port) + ")/" + m.DBName + "?" + m.Config
}

// GetLogLevel 返回 GORM 日志级别
func (m MySQLConfig) GetLogLevel() logger.LogLevel {
	switch strings.ToLower(m.LogMode) {
	case "silent":
		return logger.Silent
	case "error":
		return logger.Error
	case "warn":
		return logger.Warn
	case "info":
		return logger.Info
	default:
		return logger.Info
	}
}
