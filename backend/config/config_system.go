package config

import "fmt"

type SystemConfig struct {
	Host          string `yaml:"host" json:"host"`                     // 主机地址
	Port          int    `yaml:"port" json:"port"`                     // 端口号
	Env           string `yaml:"env" json:"env"`                       // 运行环境
	RouterPrefix  string `yaml:"router_prefix" json:"router_prefix"`   // 路由前缀
	UseMultiPoint bool   `yaml:"use_multipoint" json:"use_multipoint"` // 是否使用多点登录
	SessionsKey   string `yaml:"sessions_key" json:"sessions_key"`     // Sessions密钥
	OssType       string `yaml:"oss_type" json:"oss_type"`             // 对象存储类型
}

func (s SystemConfig) Addr() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}
