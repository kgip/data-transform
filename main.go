package main

import (
	"data-transform/global"
	"data-transform/initialize"
	"github.com/kgip/redis-lock/adapters"
	"github.com/kgip/redis-lock/lock"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	//1.初始化配置文件
	initialize.Config(global.ConfigPath)
	//2.初始化zap日志
	global.LOG = initialize.Zap()
	//3.初始化redis客户端和redis分布式锁
	global.Redis = initialize.Redis()
	global.LockOperator, _ = lock.NewRedisLockOperator(adapters.NewGoRedisV8Adapter(global.Redis), lock.Config{})
	//4.初始化gorm
	global.DB = initialize.Gorm()
	//7.初始化server
	go func() {
		//接收服务异常退出信息
		signalChan := make(chan os.Signal)
		signal.Notify(signalChan, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		<-signalChan
	}()
	server := initialize.Router()
	server.Run(":80")
}
