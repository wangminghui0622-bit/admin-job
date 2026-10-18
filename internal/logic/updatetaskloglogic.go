package logic

import (
	"context"
	"time"

	"github.com/wangminghui0622-bit/admin-common/dal/model"
	"github.com/wangminghui0622-bit/admin-common/errorx"
	"github.com/wangminghui0622-bit/admin-job/internal/svc"
	"github.com/wangminghui0622-bit/admin-job/internal/utils/dberrorhandler"
	"github.com/wangminghui0622-bit/admin-job/job"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateTaskLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateTaskLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTaskLogLogic {
	return &UpdateTaskLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UpdateTaskLog updates only finished_at and result (SetNotNil style).
func (l *UpdateTaskLogLogic) UpdateTaskLog(in *job.TaskLogInfo) (*job.BaseResp, error) {
	if in.Id == nil {
		return nil, errorx.New(errorx.InvalidArgument, "id is required")
	}

	db := l.svcCtx.DB.WithContext(l.ctx)

	var count int64
	if err := db.Model(&model.SysTaskLog{}).Where("id = ?", *in.Id).Count(&count).Error; err != nil {
		return nil, dberrorhandler.DefaultError(l.Logger, err, in)
	}
	if count == 0 {
		return nil, errorx.New(errorx.InvalidArgument, "target not exist")
	}

	updates := map[string]any{}
	if in.FinishedAt != nil {
		updates["finished_at"] = time.UnixMilli(*in.FinishedAt)
	}
	if in.Result != nil {
		updates["result"] = int32(*in.Result)
	}

	if len(updates) > 0 {
		if err := db.Model(&model.SysTaskLog{}).Where("id = ?", *in.Id).
			Updates(updates).Error; err != nil {
			return nil, dberrorhandler.DefaultError(l.Logger, err, in)
		}
	}

	return &job.BaseResp{Msg: "success"}, nil
}
