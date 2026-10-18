package logic

import (
	"context"
	"time"

	"github.com/wangminghui0622-bit/admin-common/dal/model"
	"github.com/wangminghui0622-bit/admin-job/internal/svc"
	"github.com/wangminghui0622-bit/admin-job/internal/utils/dberrorhandler"
	"github.com/wangminghui0622-bit/admin-job/job"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateTaskLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateTaskLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateTaskLogLogic {
	return &CreateTaskLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CreateTaskLog only takes finished_at and result from the request. started_at
// is filled server side (the reference ent schema also defaults it with
// time.Now), because the column has no database-level default.
func (l *CreateTaskLogLogic) CreateTaskLog(in *job.TaskLogInfo) (*job.BaseIDResp, error) {
	taskLog := &model.SysTaskLog{
		StartedAt:  time.Now(),
		FinishedAt: time.Now(),
	}
	if in.FinishedAt != nil {
		taskLog.FinishedAt = time.UnixMilli(*in.FinishedAt)
	}
	if in.Result != nil {
		taskLog.Result = int32(*in.Result)
	}

	if err := l.svcCtx.DB.WithContext(l.ctx).Create(taskLog).Error; err != nil {
		return nil, dberrorhandler.DefaultError(l.Logger, err, in)
	}

	return &job.BaseIDResp{Id: uint64(taskLog.ID), Msg: "success"}, nil
}
