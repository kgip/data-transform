package main

import (
	"context"
	"data-transform/global"
	"data-transform/initialize"
	"fmt"
	"github.com/DeanThompson/ginpprof"
	"github.com/kgip/redis-lock/adapters"
	"github.com/kgip/redis-lock/lock"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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
	global.DB = initialize.DB()
	//5.初始化MQ
	initialize.MQ()
	//6.初始化任务
	initialize.Task()
	//7.初始化server
	router := initialize.Router()
	ginpprof.Wrap(router)
	s := &http.Server{
		Addr:    global.Config.Host + ":" + fmt.Sprintf("%d", global.Config.Port),
		Handler: router,
	}
	go func() {
		if err := s.ListenAndServe(); err != nil {
			panic(err)
		}
	}()
	//接收服务异常退出信息
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-signalChan
	{
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		s.Shutdown(ctx)
		db, _ := global.DB.DB()
		db.Close()
		global.Redis.Close()
	}
}
