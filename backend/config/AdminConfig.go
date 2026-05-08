package config

type AdminConfig struct {
	Username string `yaml:"username" json:"username"` // 管理员用户名
	Email    string `yaml:"email" json:"email"`       // 管理员邮箱
	Password string `yaml:"password" json:"password"` // 管理员密码（可选，命令行参数优先）
}
