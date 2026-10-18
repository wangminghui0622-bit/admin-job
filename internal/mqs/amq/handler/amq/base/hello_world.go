package base

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/wangminghui0622-bit/admin-common/dal/model"
	"github.com/wangminghui0622-bit/admin-common/errorx"
	"github.com/wangminghui0622-bit/admin-job/internal/enum/taskresult"
	"github.com/wangminghui0622-bit/admin-job/internal/mqs/amq/types/pattern"
	"github.com/wangminghui0622-bit/admin-job/internal/mqs/amq/types/payload"
	"github.com/wangminghui0622-bit/admin-job/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type HelloWorldHandler struct {
	svcCtx *svc.ServiceContext
}

func NewHelloWorldHandler(svcCtx *svc.ServiceContext) *HelloWorldHandler {
	return &HelloWorldHandler{svcCtx: svcCtx}
}

// ProcessTask loads the task id lazily on every run instead of at startup.
// This fixes the reference implementation's panic when the task row does not
// exist yet (the service can be started before initDatabase is called).
func (l *HelloWorldHandler) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var task model.SysTask
	if err := l.svcCtx.DB.WithContext(ctx).
		Where("pattern = ?", pattern.RecordHelloWorld).
		First(&task).Error; err != nil {
		logx.Errorw("failed to load task info", logx.Field("detail", err.Error()))
		return errorx.New(errorx.InvalidArgument, "target not exist")
	}

	var p payload.HelloWorldPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("failed to unmarshal the payload %q: %w", string(t.Payload()), err)
	}

	startTime := time.Now()
	fmt.Printf("Hi! %s\n", p.Name)
	finishTime := time.Now()

	taskLog := &model.SysTaskLog{
		StartedAt:    startTime,
		FinishedAt:   finishTime,
		Result:       int32(taskresult.Success),
		TaskTaskLogs: &task.ID,
	}
	if err := l.svcCtx.DB.WithContext(ctx).Create(taskLog).Error; err != nil {
		logx.Errorw("failed to save task log to database", logx.Field("detail", err.Error()))
		return errorx.New(errorx.ServerError, "database error")
	}

	return nil
}
