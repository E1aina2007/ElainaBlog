package config

type AuthConfig struct {
	AccessTokenSecret      string `yaml:"access_token_secret" json:"access_token_secret"`             // 访问令牌密钥
	RefreshTokenSecret     string `yaml:"refresh_token_secret" json:"refresh_token_secret"`           // 刷新令牌密钥
	AccessTokenExpiryTime  string `yaml:"access_token_expiry_time" json:"access_token_expiry_time"`   // 访问令牌过期时间
	RefreshTokenExpiryTime string `yaml:"refresh_token_expiry_time" json:"refresh_token_expiry_time"` // 刷新令牌过期时间
	Issuer                 string `yaml:"issuer" json:"issuer"`                                       // 签发者
}
