package config

type ZapConfig struct {
	Level          string `yaml:"level" json:"level"`                       // 日志级别
	Format         string `yaml:"format" json:"format"`                     // 日志格式
	MaxSize        int    `yaml:"max_size" json:"max_size"`                 // 单个日志文件最大大小（MB）
	MaxBackups     int    `yaml:"max_backups" json:"max_backups"`           // 保留的最大日志文件数
	MaxAge         int    `yaml:"max_age" json:"max_age"`                   // 日志文件保留的最大天数
	Compress       bool   `yaml:"compress" json:"compress"`                 // 是否压缩旧日志文件
	FileName       string `yaml:"fileName" json:"fileName"`                 // 日志文件名
	IsConsolePrint bool   `yaml:"is_console_print" json:"is_console_print"` // 是否在控制台打印
}
