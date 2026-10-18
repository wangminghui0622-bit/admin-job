package logic

import (
	"context"

	"github.com/wangminghui0622-bit/admin-common/dal/model"
	"github.com/wangminghui0622-bit/admin-job/internal/svc"
	"github.com/wangminghui0622-bit/admin-job/internal/utils/dberrorhandler"
	"github.com/wangminghui0622-bit/admin-job/job"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateTaskLogic {
	return &CreateTaskLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CreateTask writes the task to the database only. The periodic task manager
// picks it up within one sync interval (no asynq client use here).
func (l *CreateTaskLogic) CreateTask(in *job.TaskInfo) (*job.BaseIDResp, error) {
	task := &model.SysTask{}
	if in.Name != nil {
		task.Name = *in.Name
	}
	if in.TaskGroup != nil {
		task.TaskGroup = *in.TaskGroup
	}
	if in.CronExpression != nil {
		task.CronExpression = *in.CronExpression
	}
	if in.Pattern != nil {
		task.Pattern = *in.Pattern
	}
	if in.Payload != nil {
		task.Payload = *in.Payload
	}
	if in.Status != nil {
		status := int32(*in.Status)
		task.Status = &status
	}

	if err := l.svcCtx.DB.WithContext(l.ctx).Create(task).Error; err != nil {
		return nil, dberrorhandler.DefaultError(l.Logger, err, in)
	}

	return &job.BaseIDResp{Id: uint64(task.ID), Msg: "success"}, nil
}
