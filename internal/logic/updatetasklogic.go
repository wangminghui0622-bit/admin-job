package logic

import (
	"context"

	"github.com/wangminghui0622-bit/admin-common/dal/model"
	"github.com/wangminghui0622-bit/admin-common/errorx"
	"github.com/wangminghui0622-bit/admin-job/internal/svc"
	"github.com/wangminghui0622-bit/admin-job/internal/utils/dberrorhandler"
	"github.com/wangminghui0622-bit/admin-job/job"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTaskLogic {
	return &UpdateTaskLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UpdateTask updates only the provided (non-nil) fields, SetNotNil style.
func (l *UpdateTaskLogic) UpdateTask(in *job.TaskInfo) (*job.BaseResp, error) {
	if in.Id == nil {
		return nil, errorx.New(errorx.InvalidArgument, "id is required")
	}

	db := l.svcCtx.DB.WithContext(l.ctx)

	var count int64
	if err := db.Model(&model.SysTask{}).Where("id = ?", *in.Id).Count(&count).Error; err != nil {
		return nil, dberrorhandler.DefaultError(l.Logger, err, in)
	}
	if count == 0 {
		return nil, errorx.New(errorx.InvalidArgument, "target not exist")
	}

	updates := map[string]any{}
	if in.Status != nil {
		updates["status"] = int32(*in.Status)
	}
	if in.Name != nil {
		updates["name"] = *in.Name
	}
	if in.TaskGroup != nil {
		updates["task_group"] = *in.TaskGroup
	}
	if in.CronExpression != nil {
		updates["cron_expression"] = *in.CronExpression
	}
	if in.Pattern != nil {
		updates["pattern"] = *in.Pattern
	}
	if in.Payload != nil {
		updates["payload"] = *in.Payload
	}

	if len(updates) > 0 {
		if err := db.Model(&model.SysTask{}).Where("id = ?", *in.Id).
			Updates(updates).Error; err != nil {
			return nil, dberrorhandler.DefaultError(l.Logger, err, in)
		}
	}

	return &job.BaseResp{Msg: "success"}, nil
}
