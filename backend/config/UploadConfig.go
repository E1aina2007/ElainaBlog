package config

type UploadConfig struct {
	Size int    `yaml:"size" json:"size"` // 最大上传大小（MB）
	Path string `yaml:"path" json:"path"` // 本地存储路径
}
