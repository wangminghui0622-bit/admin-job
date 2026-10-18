package svc

import (
	"github.com/hibiken/asynq"
	"github.com/wangminghui0622-bit/admin-common/dal/query"
	"github.com/wangminghui0622-bit/admin-job/internal/config"
	"github.com/wangminghui0622-bit/admin-job/internal/mqs/amq/types/periodicconfig"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config      config.Config
	DB          *gorm.DB
	Query       *query.Query
	Redis       *redis.Redis
	AsynqServer *asynq.Server
	AsynqPTM    *asynq.PeriodicTaskManager
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := gorm.Open(mysql.Open(c.DatabaseDSN), &gorm.Config{})
	logx.Must(err)

	redisClient := redis.MustNewRedis(redis.RedisConf{
		Host: c.RedisConf.Host,
		Type: c.RedisConf.Type,
		Pass: c.RedisConf.Pass,
	})

	redisConnOpt := asynq.RedisClientOpt{
		Addr:     c.AsynqConf.Addr,
		Password: c.AsynqConf.Pass,
		DB:       c.AsynqConf.DB,
	}

	asynqServer := asynq.NewServer(redisConnOpt, asynq.Config{
		Concurrency: c.AsynqConf.Concurrency,
		// A handler error always counts as a failure (same as the reference implementation).
		IsFailure: func(err error) bool { return true },
	})

	ptm, err := asynq.NewPeriodicTaskManager(asynq.PeriodicTaskManagerOpts{
		RedisConnOpt:               redisConnOpt,
		SyncInterval:               c.AsynqConf.SyncInterval,
		PeriodicTaskConfigProvider: periodicconfig.NewGormConfigProvider(db),
	})
	logx.Must(err)

	return &ServiceContext{
		Config:      c,
		DB:          db,
		Query:       query.Use(db),
		Redis:       redisClient,
		AsynqServer: asynqServer,
		AsynqPTM:    ptm,
	}
}
