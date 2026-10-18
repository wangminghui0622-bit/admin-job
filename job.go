package main

import (
	"flag"
	"fmt"

	"github.com/wangminghui0622-bit/admin-job/internal/config"
	"github.com/wangminghui0622-bit/admin-job/internal/mqs/amq/task/dynamicperiodictask"
	"github.com/wangminghui0622-bit/admin-job/internal/mqs/amq/task/mqtask"
	"github.com/wangminghui0622-bit/admin-job/internal/server"
	"github.com/wangminghui0622-bit/admin-job/internal/svc"
	"github.com/wangminghui0622-bit/admin-job/job"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/job.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		job.RegisterJobServer(grpcServer, server.NewJobServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})

	serviceGroup := service.NewServiceGroup()
	defer func() {
		serviceGroup.Stop()
		logx.Close()
	}()

	// Single process runs three services: zrpc server + asynq worker + DPTask manager.
	serviceGroup.Add(s)
	if c.AsynqConf.Enable {
		serviceGroup.Add(mqtask.NewMQTask(ctx))
	}
	if c.TaskConf.EnableDPTask {
		serviceGroup.Add(dynamicperiodictask.NewDPTask(ctx))
	}

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	serviceGroup.Start()
}
