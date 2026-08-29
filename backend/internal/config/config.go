package config

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Db           DbConfig           `yaml:"db" json:"db"`                     // 数据库配置
	Auth         AuthConfig         `yaml:"auth" json:"auth"`                 // 认证配置
	Server       ServerConfig       `yaml:"server" json:"server"`             // 系统配置
	Redis        RedisConfig        `yaml:"redis" json:"redis"`               // Redis 配置
	Upload       UploadConfig       `yaml:"upload" json:"upload"`             // 上传配置
	Smtp         SmtpConfig         `yaml:"smtp" json:"smtp"`                 // SMTP 配置
	Verification VerificationConfig `yaml:"verification" json:"verification"` //验证码配置
	Admin        AdminConfig        `yaml:"admin" json:"admin"`               // 管理员配置

	Dev bool // 是否开发模式，仅由环境变量 MODE 决定，不读取 yaml
}

// AdminConfig 管理员初始化配置
type AdminConfig struct {
	Username string `yaml:"username" json:"username"` // 管理员用户名
	Email    string `yaml:"email" json:"email"`       // 管理员邮箱
	Password string `yaml:"password" json:"password"` // 管理员密码（必填，启动时自动创建/重置）
}

// AuthConfig 认证配置
type AuthConfig struct {
	AccessTokenSecret      string `yaml:"access_token_secret" json:"access_token_secret"`             // 访问令牌密钥
	RefreshTokenSecret     string `yaml:"refresh_token_secret" json:"refresh_token_secret"`           // 刷新令牌密钥
	AccessTokenExpiryTime  string `yaml:"access_token_expiry_time" json:"access_token_expiry_time"`   // 访问令牌过期时间
	RefreshTokenExpiryTime string `yaml:"refresh_token_expiry_time" json:"refresh_token_expiry_time"` // 刷新令牌过期时间
	Issuer                 string `yaml:"issuer" json:"issuer"`                                       // 签发者
}

// DbConfig 数据库连接配置
type DbConfig struct {
	Host     string `yaml:"host" json:"host"`         // 数据库主机
	Port     int    `yaml:"port" json:"port"`         // 数据库端口
	Username string `yaml:"username" json:"username"` // 数据库用户
	Password string `yaml:"password" json:"password"` // 数据库密码
	DBName   string `yaml:"db_name" json:"db_name"`   // 数据库名
}

// RedisConfig Redis 配置
type RedisConfig struct {
	Addr     string `yaml:"address" json:"address"`
	Password string `yaml:"password" json:"password"`
	DB       int    `yaml:"db" json:"db"`
	PoolSize int    `yaml:"pool_size" json:"pool_size"`
}

// ServerConfig 服务配置
type ServerConfig struct {
	Port int `yaml:"port" json:"port"` // 端口号
}

// GetAddress 返回监听地址 :port
func (s ServerConfig) GetAddress() string {
	return fmt.Sprintf(":%d", s.Port)
}

// SmtpConfig SMTP 配置
type SmtpConfig struct {
	Host         string `yaml:"host" json:"host"`
	Port         int    `yaml:"port" json:"port"`
	From         string `yaml:"from" json:"from"`
	Verification string `yaml:"verification" json:"verification"`
	SSL          bool   `yaml:"ssl" json:"ssl"`
	Timeout      int    `yaml:"timeout" json:"timeout"`
}

// UploadConfig 上传配置
type UploadConfig struct {
	Size       int    `yaml:"size" json:"size"`               // 最大上传大小（MB）
	Path       string `yaml:"path" json:"path"`               // 本地存储路径
	AvatarPath string `yaml:"avatar_path" json:"avatar_path"` // 头像存储路径
	AvatarSize int    `yaml:"avatar_size" json:"avatar_size"` // 头像最大上传大小（MB）

	// 上传频率限制
	RateLimit  int `yaml:"rate_limit" json:"rate_limit"`   // 每小时最多上传图片数量，默认 20
	RateWindow int `yaml:"rate_window" json:"rate_window"` // 频率限制时间窗口（秒），默认 3600

	// 孤儿图片清理
	CleanupEnabled bool   `yaml:"cleanup_enabled" json:"cleanup_enabled"` // 是否启用孤儿图片清理，默认 true
	CleanupCron    string `yaml:"cleanup_cron" json:"cleanup_cron"`       // 清理任务 cron 表达式，默认 "0 3 * * *"（每天凌晨3点）
}

// VerificationConfig 验证码配置
type VerificationConfig struct {
	CodeLength     int `yaml:"code_length" json:"code_length"`         // 验证码长度
	ExpireTime     int `yaml:"expire_time" json:"expire_time"`         // 过期时间
	ResendInterval int `yaml:"resend_interval" json:"resend_interval"` //重发间隔
}

var GlobalConfig *Config // 全局配置实例，由入口加载配置后赋值（兼容现有调用）

// LoadConfig 加载配置：加载 .env、按 CONFIG_PATH 选择 yaml，并赋值 GlobalConfig
// CONFIG_PATH 未设置时默认读取 configs/config.dev.yaml
func LoadConfig() error {
	configDir := "configs"
	if dir := os.Getenv("CONFIG_DIR"); dir != "" {
		configDir = dir
	}

	// 加载 .env（可选，不存在时使用默认配置）
	if err := godotenv.Load(); err != nil {
		log.Println("未找到 .env 文件，使用默认配置:", err)
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = configDir + "/config.dev.yaml"
	}

	cfg, err := Load(configPath)
	if err != nil {
		return err
	}
	ensureJwtSecrets(&cfg)
	GlobalConfig = &cfg
	if cfg.Dev {
		log.Println("开发模式")
	} else {
		log.Println("生产模式")
	}
	return nil
}

// jwtPlaceholderSecrets 示例配置中的公开占位密钥，绝不能参与真实签发
var jwtPlaceholderSecrets = map[string]struct{}{
	"your-access-secret":  {},
	"your-refresh-secret": {},
}

// ensureJwtSecrets 校验 JWT 密钥：空值或示例占位值时生成本次运行专用的随机密钥。
// 临时密钥仅当前进程有效，重启后所有已签发 token 失效（单实例下即全员重新登录）；
// 生产部署必须配置固定密钥。
func ensureJwtSecrets(cfg *Config) {
	fields := []struct {
		name   string
		secret *string
	}{
		{"access", &cfg.Auth.AccessTokenSecret},
		{"refresh", &cfg.Auth.RefreshTokenSecret},
	}
	for _, f := range fields {
		s := strings.TrimSpace(*f.secret)
		if s != "" {
			if _, isPlaceholder := jwtPlaceholderSecrets[s]; !isPlaceholder {
				continue
			}
		}
		generated, err := randomSecret()
		if err != nil {
			log.Fatalf("生成 JWT %s 临时密钥失败: %v", f.name, err)
		}
		*f.secret = generated
		log.Printf("JWT %s 密钥未配置或为示例占位值，已生成本次运行专用的临时密钥；重启后所有登录态将失效，生产部署请配置固定密钥", f.name)
	}
}

// randomSecret 生成 32 字节随机密钥的十六进制表示
func randomSecret() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// Load 读取 yaml 配置文件，并用环境变量覆盖对应字段（环境变量优先）
func Load(filename string) (Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return Config{}, fmt.Errorf("failed to read config file: %w", err)
	}

	cfg := Config{}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("failed to parse config %s: %w", filename, err)
	}

	OverrideWithEnv(&cfg)
	applyDefaults(&cfg)
	return cfg, nil
}

// OverrideWithEnv 用环境变量覆盖 yaml 中的配置（密码、地址等）
// 环境变量名见 configs/.env.example，未设置时保留 yaml 原值
func OverrideWithEnv(cfg *Config) {
	if cfg == nil {
		return
	}

	// MODE 仅从环境变量读取（不读 yaml），默认开发模式
	switch mode := os.Getenv("MODE"); mode {
	case "dev", "development", "":
		cfg.Dev = true
	case "prod", "production":
		cfg.Dev = false
	default:
		log.Printf("未知 MODE: %q，默认使用开发模式", mode)
		cfg.Dev = true
	}

	// 数据库
	cfg.Db.Host = envString("MYSQL_HOST", cfg.Db.Host)
	cfg.Db.Port = envInt("MYSQL_PORT", cfg.Db.Port)
	cfg.Db.Username = envString("MYSQL_USER", cfg.Db.Username)
	// MYSQL_PASSWORD 仅在没有设置 MYSQL_ROOT_PASSWORD 时生效，避免优先级歧义
	if v := os.Getenv("MYSQL_PASSWORD"); v != "" && os.Getenv("MYSQL_ROOT_PASSWORD") == "" {
		cfg.Db.Password = v
	}
	if v := os.Getenv("MYSQL_ROOT_PASSWORD"); v != "" {
		cfg.Db.Password = v
	}
	cfg.Db.DBName = envString("MYSQL_DATABASE", cfg.Db.DBName)

	// Redis
	cfg.Redis.Addr = envString("REDIS_ADDRESS", cfg.Redis.Addr)
	cfg.Redis.Password = envString("REDIS_PASSWORD", cfg.Redis.Password)
	cfg.Redis.DB = envInt("REDIS_DB", cfg.Redis.DB)

	// JWT 密钥
	cfg.Auth.AccessTokenSecret = envString("JWT_ACCESS_SECRET", cfg.Auth.AccessTokenSecret)
	cfg.Auth.RefreshTokenSecret = envString("JWT_REFRESH_SECRET", cfg.Auth.RefreshTokenSecret)

	// SMTP
	cfg.Smtp.Host = envString("SMTP_HOST", cfg.Smtp.Host)
	cfg.Smtp.Port = envInt("SMTP_PORT", cfg.Smtp.Port)
	cfg.Smtp.From = envString("SMTP_FROM", cfg.Smtp.From)
	cfg.Smtp.Verification = envString("SMTP_PASSWORD", cfg.Smtp.Verification)
	cfg.Smtp.SSL = envBool("SMTP_SSL", cfg.Smtp.SSL)

	// 服务
	cfg.Server.Port = envInt("SERVER_PORT", cfg.Server.Port)

	// 上传目录
	cfg.Upload.Path = envString("UPLOAD_PATH", cfg.Upload.Path)
	cfg.Upload.AvatarPath = envString("UPLOAD_AVATAR_PATH", cfg.Upload.AvatarPath)

	// 管理员
	cfg.Admin.Username = envString("ADMIN_USERNAME", cfg.Admin.Username)
	cfg.Admin.Email = envString("ADMIN_EMAIL", cfg.Admin.Email)
	cfg.Admin.Password = envString("ADMIN_PASSWORD", cfg.Admin.Password)
}

// applyDefaults 填充未配置的默认值
func applyDefaults(cfg *Config) {
	if cfg.Auth.AccessTokenExpiryTime == "" {
		cfg.Auth.AccessTokenExpiryTime = "2h"
	}
	if cfg.Auth.RefreshTokenExpiryTime == "" {
		cfg.Auth.RefreshTokenExpiryTime = "7d"
	}
	if cfg.Auth.Issuer == "" {
		cfg.Auth.Issuer = "Elaina"
	}

	// 上传配置默认值
	if cfg.Upload.AvatarPath == "" {
		cfg.Upload.AvatarPath = "uploads/avatars"
	}
	if cfg.Upload.AvatarSize == 0 {
		cfg.Upload.AvatarSize = 5
	}
	if cfg.Upload.RateLimit == 0 {
		cfg.Upload.RateLimit = 20 // 每小时最多上传 20 张
	}
	if cfg.Upload.RateWindow == 0 {
		cfg.Upload.RateWindow = 3600 // 默认 1 小时窗口
	}
	if cfg.Upload.CleanupCron == "" {
		cfg.Upload.CleanupCron = "0 3 * * *" // 每天凌晨 3 点
	}
	// CleanupEnabled 默认为 true
	cfg.Upload.CleanupEnabled = true
}

// envString 返回环境变量值；未设置或为空时返回当前值
func envString(name, current string) string {
	if v := os.Getenv(name); v != "" {
		return v
	}
	return current
}

// envInt 解析整数环境变量；为空或非法时返回当前值
func envInt(name string, current int) int {
	v := os.Getenv(name)
	if v == "" {
		return current
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		log.Printf("环境变量 %s 不是有效整数：%q，已忽略", name, v)
		return current
	}
	return n
}

// envBool 解析布尔环境变量；为空或非法时返回当前值
func envBool(name string, current bool) bool {
	v := os.Getenv(name)
	if v == "" {
		return current
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		log.Printf("环境变量 %s 不是有效布尔值：%q，已忽略", name, v)
		return current
	}
	return b
}
