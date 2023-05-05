package initialize

import (
	"data-transform/global"
	"data-transform/model/po"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"regexp"
	"strconv"
	"time"
)

//initialize gorm
func Gorm() (db *gorm.DB) {
	global.LOG.Info("start initialize gorm")
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: global.Config.Dsn(),
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "rp_",
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}
	if mysqlDB, err := db.DB(); err != nil {
		panic(err)
	} else {
		mysqlDB.SetMaxOpenConns(global.Config.MaxOpenConns)
		mysqlDB.SetMaxIdleConns(global.Config.MaxIdleConns)
		//设置空闲连接的存活时间
		lifetime := global.Config.ConnMaxLifetime
		if ok, err := regexp.MatchString("^[0-9]+[s|m|h]{0,1}$", lifetime); err != nil || !ok {
			global.LOG.Warn("use default ConnMaxIdleTime", zap.Any("err", err))
			mysqlDB.SetConnMaxIdleTime(global.ConnMaxIdleTime)
		} else {
			var timeUnit time.Duration
			n := -1
			switch lifetime[len(lifetime)-1:] {
			case "s":
				timeUnit = time.Second
			case "m":
				timeUnit = time.Minute
			case "h":
				timeUnit = time.Hour
			default:
				timeUnit = time.Second
				if n, err = strconv.Atoi(lifetime); err != nil {
					panic(err)
				}
			}
			if n < 0 {
				if n, err = strconv.Atoi(lifetime[:len(lifetime)-1]); err != nil {
					panic(err)
				}
			}
			mysqlDB.SetConnMaxIdleTime(time.Duration(n) * timeUnit)
		}
	}
	InitSchemas(db)
	return db
}

//初始化数据库表
func InitSchemas(db *gorm.DB) {
	zap.S().Info("开始初始化数据库表")
	db.AutoMigrate(
		&po.User{},

	)
	zap.S().Info("数据库表初始化完成")
}
