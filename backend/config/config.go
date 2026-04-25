package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Db     DbConfig     `yaml:"db" json:"db"`         // 数据库配置
	Auth   AuthConfig   `yaml:"auth" json:"auth"`     // 认证配置
	Zap    ZapConfig    `yaml:"zap" json:"zap"`       // 日志配置
	Server ServerConfig `yaml:"server" json:"server"` // 系统配置
	Redis  RedisConfig  `yaml:"redis" json:"redis"`   // Redis 配置
	Email  EmailConfig  `yaml:"email" json:"email"`   // 邮箱配置
	Upload UploadConfig `yaml:"upload" json:"upload"` // 上传配置

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

	GlobalConfig = &config
	return nil
}

func CheckMode() string {
	var path string
	err := godotenv.Load("config/.env")
	if err != nil {
		log.Fatalf("加载配置文件失败：%v", err)
	}

	mode := os.Getenv("MODE")
	switch mode {
	case "dev":
		path = "config/config.dev.yaml"
	case "prod":
		path = "config/config.prod.yaml"
	default:
		log.Fatalf("环境变量MODE错误: %v,请检查.env文件", mode)
	}
	return path
}
