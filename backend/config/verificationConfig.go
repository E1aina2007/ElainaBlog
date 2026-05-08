package config

type VerificationConfig struct {
	CodeLength     int `yaml:"code_length" json:"code_length"`         // 验证码长度
	ExpireTime     int `yaml:"expire_time" json:"expire_time"`         // 过期时间
	ResendInterval int `yaml:"resend_interval" json:"resend_interval"` //重发间隔
}


