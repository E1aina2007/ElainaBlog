package config

type EmailConfig struct {
	Host     string `yaml:"host" json:"host"`         // 邮件服务器地址
	Port     int    `yaml:"port" json:"port"`         // 邮件服务器端口
	Username string `yaml:"username" json:"username"` // 邮箱用户名
	Password string `yaml:"password" json:"password"` // 邮箱密码
	From     string `yaml:"from" json:"from"`         // 发件人地址
}
