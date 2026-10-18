package config

import (
	"time"

	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	DatabaseDSN string
	RedisConf   struct {
		Host string
		Type string
		Pass string
	}
	AsynqConf AsynqConf
	TaskConf  TaskConf
}

type AsynqConf struct {
	Enable       bool          `json:",default=true"`
	Addr         string        `json:",default=127.0.0.1:6379"`
	Pass         string        `json:",optional"`
	DB           int           `json:",default=0"`
	Concurrency  int           `json:",default=20"`
	SyncInterval time.Duration `json:",default=10s"`
}

type TaskConf struct {
	EnableScheduledTask bool `json:",default=false"`
	EnableDPTask        bool `json:",default=true"`
}
