package mqtask

import (
	"sync"

	"github.com/hibiken/asynq"
	"github.com/wangminghui0622-bit/admin-job/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

// MQTask runs the asynq worker (consumer) side of the job service.
type MQTask struct {
	svcCtx   *svc.ServiceContext
	mux      *asynq.ServeMux
	done     chan struct{}
	stopOnce sync.Once
}

func NewMQTask(svcCtx *svc.ServiceContext) *MQTask {
	return &MQTask{
		svcCtx: svcCtx,
		done:   make(chan struct{}),
	}
}

// Start starts the worker and blocks until Stop is called.
func (m *MQTask) Start() {
	m.Register()

	if err := m.svcCtx.AsynqServer.Start(m.mux); err != nil {
		logx.Errorw("failed to start asynq worker", logx.Field("detail", err.Error()))
		return
	}
	logx.Infow("asynq worker started")

	<-m.done
}

// Stop gracefully shuts down the worker.
func (m *MQTask) Stop() {
	m.stopOnce.Do(func() {
		close(m.done)
		m.svcCtx.AsynqServer.Shutdown()
	})
}
