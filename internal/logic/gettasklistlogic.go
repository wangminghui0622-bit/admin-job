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

type GetTaskListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTaskListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTaskListLogic {
	return &GetTaskListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetTaskListLogic) GetTaskList(in *job.TaskListReq) (*job.TaskListResp, error) {
	db := l.svcCtx.DB.WithContext(l.ctx).Model(&model.SysTask{})
	if in.Name != nil {
		db = db.Where("name LIKE ?", "%"+*in.Name+"%")
	}
	if in.TaskGroup != nil {
		db = db.Where("task_group LIKE ?", "%"+*in.TaskGroup+"%")
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

	var tasks []*model.SysTask
	if err := db.Session(&gorm.Session{}).
		Order("id DESC").
		Offset(int((page - 1) * pageSize)).
		Limit(int(pageSize)).
		Find(&tasks).Error; err != nil {
		return nil, dberrorhandler.DefaultError(l.Logger, err, in)
	}

	data := make([]*job.TaskInfo, 0, len(tasks))
	for _, v := range tasks {
		data = append(data, taskInfoFromModel(v))
	}

	return &job.TaskListResp{Total: uint64(total), Data: data}, nil
}
