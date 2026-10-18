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

type GetTaskLogListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTaskLogListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTaskLogListLogic {
	return &GetTaskLogListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetTaskLogListLogic) GetTaskLogList(in *job.TaskLogListReq) (*job.TaskLogListResp, error) {
	db := l.svcCtx.DB.WithContext(l.ctx).Model(&model.SysTaskLog{})

	if in.TaskId != nil {
		db = db.Where("task_task_logs = ?", *in.TaskId)
	}
	if in.Result != nil && *in.Result != 0 {
		db = db.Where("result = ?", *in.Result)
	}

	var total int64
	if err := db.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, dberrorhandler.DefaultError(l.Logger, err, in)
	}

	page, pageSize := in.Page, in.PageSize
	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}

	var logs []*model.SysTaskLog
	if err := db.Session(&gorm.Session{}).
		Order("id DESC").
		Offset(int((page - 1) * pageSize)).
		Limit(int(pageSize)).
		Find(&logs).Error; err != nil {
		return nil, dberrorhandler.DefaultError(l.Logger, err, in)
	}

	data := make([]*job.TaskLogInfo, 0, len(logs))
	for _, v := range logs {
		data = append(data, taskLogInfoFromModel(v))
	}

	return &job.TaskLogListResp{Total: uint64(total), Data: data}, nil
}
