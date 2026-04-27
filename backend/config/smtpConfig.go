package config

type SmtpConfig struct {
	Host         string `yaml:"host" json:"host"`
	Port         int    `yaml:"port" json:"port"`
	From         string `yaml:"from" json:"from"`
	To           string `yaml:"to" json:"to"`
	Verification string `yaml:"verification" json:"verification"`
	SSL          bool   `yaml:"ssl" json:"ssl"`
	Timeout      int    `yaml:"timeout" json:"timeout"`
}
