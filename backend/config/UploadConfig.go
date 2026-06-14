package config

type UploadConfig struct {
	Size       int    `yaml:"size" json:"size"`               // 最大上传大小（MB）
	Path       string `yaml:"path" json:"path"`               // 本地存储路径
	AvatarPath string `yaml:"avatar_path" json:"avatar_path"` // 头像存储路径
	AvatarSize int    `yaml:"avatar_size" json:"avatar_size"` // 头像最大上传大小（MB）

	// 上传频率限制
	RateLimit  int `yaml:"rate_limit" json:"rate_limit"`   // 每小时最多上传图片数量，默认 20
	RateWindow int `yaml:"rate_window" json:"rate_window"` // 频率限制时间窗口（秒），默认 3600

	// 孤儿图片清理
	CleanupEnabled bool   `yaml:"cleanup_enabled" json:"cleanup_enabled"` // 是否启用孤儿图片清理，默认 true
	CleanupCron    string `yaml:"cleanup_cron" json:"cleanup_cron"`       // 清理任务 cron 表达式，默认 "0 3 * * *"（每天凌晨3点）
}
