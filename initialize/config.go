package initialize

import (
	"data-transform/global"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
)

// Config 初始化配置文件
func Config(path string) {
	log.Println("start initialize config")
	viper := viper.New()
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Println(err.Error())
		log.Panic("initialize config failed")
	}
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("config file changed")
		if err := viper.Unmarshal(global.Config); err != nil {
			log.Println(err)
		}
	})
	if err := viper.Unmarshal(global.Config); err != nil {
		log.Panic(err)
	}
	log.Println("finished initializing config")
}
