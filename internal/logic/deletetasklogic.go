package logic

import (
	"context"

	"github.com/wangminghui0622-bit/admin-common/dal/model"
	"github.com/wangminghui0622-bit/admin-job/internal/svc"
	"github.com/wangminghui0622-bit/admin-job/internal/utils/dberrorhandler"
	"github.com/wangminghui0622-bit/admin-job/job"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type DeleteTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteTaskLogic {
	return &DeleteTaskLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// DeleteTask deletes the task logs first, then the tasks, in one transaction.
func (l *DeleteTaskLogic) DeleteTask(in *job.IDsReq) (*job.BaseResp, error) {
	if len(in.Ids) == 0 {
		return &job.BaseResp{Msg: "success"}, nil
	}

	err := l.svcCtx.DB.WithContext(l.ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("task_task_logs IN ?", in.Ids).
			Delete(&model.SysTaskLog{}).Error; err != nil {
			return err
		}
		return tx.Where("id IN ?", in.Ids).Delete(&model.SysTask{}).Error
	})
	if err != nil {
		return nil, dberrorhandler.DefaultError(l.Logger, err, in)
	}

	return &job.BaseResp{Msg: "success"}, nil
}
