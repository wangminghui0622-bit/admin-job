package logic

import (
	"context"

	"github.com/wangminghui0622-bit/admin-common/dal/model"
	"github.com/wangminghui0622-bit/admin-job/internal/svc"
	"github.com/wangminghui0622-bit/admin-job/internal/utils/dberrorhandler"
	"github.com/wangminghui0622-bit/admin-job/job"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetTaskLogByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTaskLogByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTaskLogByIdLogic {
	return &GetTaskLogByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetTaskLogByIdLogic) GetTaskLogById(in *job.IDReq) (*job.TaskLogInfo, error) {
	var taskLog model.SysTaskLog
	if err := l.svcCtx.DB.WithContext(l.ctx).Where("id = ?", in.Id).First(&taskLog).Error; err != nil {
		return nil, dberrorhandler.DefaultError(l.Logger, err, in)
	}

	return taskLogInfoFromModel(&taskLog), nil
}
