package config

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var yamlPath = GetYamlPath()

// GetYamlPath 基于可执行文件所在目录定位 config.yaml
// 本地开发时从 cmd/ 运行，config.yaml 在上级目录 backend/
// 部署时将 config.yaml 放在可执行文件同级目录即可
func GetYamlPath() string {
	// 优先查找工作目录下的 config.yaml
	if _, err := os.Stat("config.yaml"); err == nil {
		return "config.yaml"
	}
	// 兼容从 cmd/ 子目录运行的情况
	if _, err := os.Stat(filepath.Join("..", "config.yaml")); err == nil {
		return filepath.Join("..", "config.yaml")
	}
	// 默认返回工作目录下
	return "config.yaml"
}

type Config struct {
	Mysql  MySQLConfig  `yaml:"mysql" json:"mysql"`   // MySQL 数据库配置
	Jwt    JwtConfig    `yaml:"jwt" json:"jwt"`       // JWT 认证配置
	Zap    ZapConfig    `yaml:"zap" json:"zap"`       // 日志配置
	System SystemConfig `yaml:"system" json:"system"` // 系统配置
	Redis  RedisConfig  `yaml:"redis" json:"redis"`   // Redis 配置
}

func ConfigInit() *Config {
	c := &Config{}
	yamlConfig, err := LoadYAML()
	if err != nil {
		log.Fatalf("加载 YAML 配置文件失败: %v", err)
	}

	err = yaml.Unmarshal(yamlConfig, c)
	if err != nil {
		log.Fatalf("解析 YAML 配置失败: %v", err)
	}
	return c
}

func LoadYAML() ([]byte, error) {
	return os.ReadFile(yamlPath)
}

func SaveYAML(cfg *Config) error {
	byteData, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(yamlPath, byteData, fs.ModePerm)
}
