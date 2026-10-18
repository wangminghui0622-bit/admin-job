package logic

import (
	"context"

	"github.com/wangminghui0622-bit/admin-common/dal/model"
	"github.com/wangminghui0622-bit/admin-job/internal/svc"
	"github.com/wangminghui0622-bit/admin-job/internal/utils/dberrorhandler"
	"github.com/wangminghui0622-bit/admin-job/job"
	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteTaskLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteTaskLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteTaskLogLogic {
	return &DeleteTaskLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteTaskLogLogic) DeleteTaskLog(in *job.IDsReq) (*job.BaseResp, error) {
	if len(in.Ids) == 0 {
		return &job.BaseResp{Msg: "success"}, nil
	}

	if err := l.svcCtx.DB.WithContext(l.ctx).
		Where("id IN ?", in.Ids).
		Delete(&model.SysTaskLog{}).Error; err != nil {
		return nil, dberrorhandler.DefaultError(l.Logger, err, in)
	}

	return &job.BaseResp{Msg: "success"}, nil
}
