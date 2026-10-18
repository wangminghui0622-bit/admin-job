package mqtask

import (
	"github.com/hibiken/asynq"
	"github.com/wangminghui0622-bit/admin-job/internal/mqs/amq/handler/amq/base"
	"github.com/wangminghui0622-bit/admin-job/internal/mqs/amq/types/pattern"
)

// Register registers task handlers. The task type equals the pattern column
// stored in sys_tasks.
func (m *MQTask) Register() {
	mux := asynq.NewServeMux()

	mux.Handle(pattern.RecordHelloWorld, base.NewHelloWorldHandler(m.svcCtx))

	m.mux = mux
}
