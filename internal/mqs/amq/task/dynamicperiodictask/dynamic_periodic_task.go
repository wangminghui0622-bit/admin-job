package dynamicperiodictask

import (
	"sync"

	"github.com/wangminghui0622-bit/admin-job/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

// DPTask runs the asynq PeriodicTaskManager. The manager polls sys_tasks
// every SyncInterval and keeps the asynq scheduler in sync with the database,
// so task CRUD never has to touch an asynq client directly.
type DPTask struct {
	svcCtx   *svc.ServiceContext
	done     chan struct{}
	stopOnce sync.Once
}

func NewDPTask(svcCtx *svc.ServiceContext) *DPTask {
	return &DPTask{
		svcCtx: svcCtx,
		done:   make(chan struct{}),
	}
}

// Start starts the periodic task manager and blocks until Stop is called.
func (m *DPTask) Start() {
	if err := m.svcCtx.AsynqPTM.Start(); err != nil {
		logx.Errorw("failed to start dynamic periodic task manager", logx.Field("detail", err.Error()))
		return
	}
	logx.Infow("dynamic periodic task manager started")

	<-m.done
}

// Stop gracefully shuts down the periodic task manager.
func (m *DPTask) Stop() {
	m.stopOnce.Do(func() {
		close(m.done)
		m.svcCtx.AsynqPTM.Shutdown()
	})
}
