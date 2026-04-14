package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Db     DbConfig     `yaml:"db" json:"db"`         // 数据库配置
	Auth   AuthConfig   `yaml:"auth" json:"auth"`     // 认证配置
	Zap    ZapConfig    `yaml:"zap" json:"zap"`       // 日志配置
	Server ServerConfig `yaml:"server" json:"server"` // 系统配置
	Redis  RedisConfig  `yaml:"redis" json:"redis"`   // Redis 配置
	Email  EmailConfig 	`yaml:"email" json:"email"`
}

var GlobalConfig *Config // 全局配置实例

func LoadConfigFromYml(ymlPath string) error {
	var config Config

	ymlConfig, err := os.ReadFile(ymlPath)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(ymlConfig, &config)
	if err != nil {
		return err
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
