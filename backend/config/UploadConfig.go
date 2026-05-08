package config

type UploadConfig struct {
	Size       int    `yaml:"size" json:"size"`             // 最大上传大小（MB）
	Path       string `yaml:"path" json:"path"`             // 本地存储路径
	AvatarPath string `yaml:"avatar_path" json:"avatar_path"` // 头像存储路径
	AvatarSize int    `yaml:"avatar_size" json:"avatar_size"` // 头像最大上传大小（MB）
}
