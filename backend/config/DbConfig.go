package config

import (
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
)

type DbConfig struct {
	SqlName      string `yaml:"sqlname" json:"sqlname"`               // Sql名
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

// GetDSN 通过 mysql.Config 结构体安全地生成 DSN（Data Source Name）连接字符串
// 相比手动拼接字符串，使用 mysql.Config 的优势在于：
//   - FormatDSN() 会自动对密码等字段中的特殊字符（@、:、/ 等）进行转义，避免 DSN 解析出错
//   - 通过结构体字段赋值，语义清晰、不易拼错格式
//
// 生成的 DSN 格式形如：
//
//	username:password@tcp(host:port)/dbname?parseTime=true&loc=Local&collation=utf8mb4_general_ci
func (m DbConfig) GetDSN() string {
	cfg := mysql.Config{
		User:                 m.Username,                                     // 数据库登录用户名
		Passwd:               m.Password,                                     // 数据库登录密码
		Net:                  "tcp",                                          // 网络协议，一般为 tcp；也支持 unix（Unix Socket）
		Addr:                 net.JoinHostPort(m.Host, strconv.Itoa(m.Port)), // 拼接为 "host:port" 格式的地址
		DBName:               m.DBName,                                       // 要连接的目标数据库名
		ParseTime:            true,                                           // 将 MySQL 的 DATE/DATETIME 自动解析为 Go 的 time.Time
		Loc:                  time.Local,                                     // 时区设置为系统本地时区
		Collation:            "utf8mb4_0900_ai_ci",                           // 字符排序规则，utf8mb4 支持完整 Unicode（含 emoji）
		AllowNativePasswords: true,                                           // 允许 MySQL 原生密码认证（mysql_native_password）
	}

	// 如果 YAML 中的 config 字段不为空，将其解析为 key=value 键值对
	// 追加到 cfg.Params 中，FormatDSN() 会把它们拼到 DSN 的 ? 后面
	if m.Config != "" {
		cfg.Params = parseParams(m.Config)
	}

	return cfg.FormatDSN()
}

// parseParams 将 URL query 风格的参数字符串解析为 map[string]string
// 输入示例: "charset=utf8mb4&parseTime=True&loc=Local"
// 输出示例: map[string]string{"charset": "utf8mb4", "parseTime": "True", "loc": "Local"}
//
// 解析流程：
//  1. 按 "&" 分割为多个 "key=value" 片段
//  2. 对每个片段用 strings.Cut 按第一个 "=" 拆分为键和值
//  3. 只有成功拆分（即包含 "="）的片段才会写入 map，跳过格式不合法的片段
func parseParams(raw string) map[string]string {
	params := make(map[string]string)
	for _, pair := range strings.Split(raw, "&") {
		if key, val, ok := strings.Cut(pair, "="); ok {
			params[key] = val
		}
	}
	return params
}
