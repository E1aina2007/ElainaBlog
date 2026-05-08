package config

import "fmt"

type ServerConfig struct {
	Host        string   `yaml:"host" json:"host"`               // 主机地址
	Port        int      `yaml:"port" json:"port"`               // 端口号
	Env         string   `yaml:"env" json:"env"`                 // 运行环境
	CorsOrigins []string `yaml:"cors_origins" json:"cors_origins"` // 允许的跨域来源
}

func (s ServerConfig) GetAddress() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}
