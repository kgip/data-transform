package config

type Server struct {
	Host     string `mapstructure:"host" json:"host" yaml:"host"`                //服务主机名
	Port     int    `mapstructure:"port" json:"port" yaml:"port"`                //服务端口号
	BasePath string `mapstructure:"base-path" json:"base-path" yaml:"base-path"` //服务访问根地址
}
