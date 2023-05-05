package config

type Config struct {
	Server `mapstructure:"server"`
	Mysql  `mapstructure:"mysql"`
	Redis  `mapstructure:"redis"`
	Zap    `mapstructure:"zap"`
}
