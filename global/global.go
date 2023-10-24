package global

import (
	"data-transform/config"
	"github.com/go-redis/redis/v8"
	rabbitmq "github.com/kgip/go-rabbit-template"
	"github.com/kgip/redis-lock/lock"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	DB           *gorm.DB
	Config       = &config.Config{}
	LOG          *zap.Logger
	Redis        *redis.Client
	LockOperator *lock.RedisLockOperator //redis分布式锁
	MQ           *rabbitmq.RabbitTemplate
)
