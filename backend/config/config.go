package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Db           DbConfig           `yaml:"db" json:"db"`                     // 数据库配置
	Auth         AuthConfig         `yaml:"auth" json:"auth"`                 // 认证配置
	Zap          ZapConfig          `yaml:"zap" json:"zap"`                   // 日志配置
	Server       ServerConfig       `yaml:"server" json:"server"`             // 系统配置
	Redis        RedisConfig        `yaml:"redis" json:"redis"`               // Redis 配置
	Upload       UploadConfig       `yaml:"upload" json:"upload"`             // 上传配置
	Smtp         SmtpConfig         `yaml:"smtp" json:"smtp"`                 // SMTP 配置
	Verification VerificationConfig `yaml:"verification" json:"verification"` //验证码配置
	Admin        AdminConfig        `yaml:"admin" json:"admin"`               // 管理员配置

	Dev bool `yaml:"dev"` // 区分开发和生产环境
}

var GlobalConfig *Config // 全局配置实例

func LoadConfigFromYml(path string) error {
	var config Config

	ymlConfig, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(ymlConfig, &config)
	if err != nil {
		return err
	}
	if config.Dev {
		log.Println("开发模式")
	} else {
		log.Println("生产模式")
	}

	if config.Auth.AccessTokenExpiryTime == "" {
		config.Auth.AccessTokenExpiryTime = "2h"
	}
	if config.Auth.RefreshTokenExpiryTime == "" {
		config.Auth.RefreshTokenExpiryTime = "7d"
	}
	if config.Auth.Issuer == "" {
		config.Auth.Issuer = "Elaina"
	}

	// 上传配置默认值
	if config.Upload.AvatarPath == "" {
		config.Upload.AvatarPath = "uploads/avatars"
	}
	if config.Upload.AvatarSize == 0 {
		config.Upload.AvatarSize = 5
	}
	if config.Upload.RateLimit == 0 {
		config.Upload.RateLimit = 20 // 每小时最多上传 20 张
	}
	if config.Upload.RateWindow == 0 {
		config.Upload.RateWindow = 3600 // 默认 1 小时窗口
	}
	if config.Upload.CleanupCron == "" {
		config.Upload.CleanupCron = "0 3 * * *" // 每天凌晨 3 点
	}
	// CleanupEnabled 默认为 true
	config.Upload.CleanupEnabled = true

	GlobalConfig = &config
	return nil
}

// getConfigDir 获取配置目录路径
// 优先使用 CONFIG_DIR 环境变量（Docker 挂载场景），否则回退到项目根目录的 config/backend/
func getConfigDir() string {
	if dir := os.Getenv("CONFIG_DIR"); dir != "" {
		return dir
	}
	return "../config/backend"
}

func CheckMode() string {
	configDir := getConfigDir()

	envPath := configDir + "/.env"
	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatalf("加载配置文件失败，请从 config/.env.example 复制并填写：%v", err)
	}

	mode := os.Getenv("MODE")
	var yamlName string
	switch mode {
	case "dev":
		yamlName = "config.dev.yaml"
	case "prod":
		yamlName = "config.prod.yaml"
	default:
		log.Fatalf("环境变量MODE错误: %v,请检查 .env 文件", mode)
	}

	path := configDir + "/" + yamlName
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Fatalf("配置文件不存在: %s，请从 config/config.example.yaml 复制并填写", path)
	}
	return path
}
