package logic

import (
	"context"
	_ "embed"

	"github.com/wangminghui0622-bit/admin-common/dal/model"
	"github.com/wangminghui0622-bit/admin-common/errorx"
	"github.com/wangminghui0622-bit/admin-job/internal/mqs/amq/types/pattern"
	"github.com/wangminghui0622-bit/admin-job/internal/svc"
	"github.com/wangminghui0622-bit/admin-job/job"
	"github.com/zeromicro/go-zero/core/logx"
)

// initSQL contains the DDL and seed data generated from admin_zero.
// It only uses CREATE TABLE IF NOT EXISTS / INSERT IGNORE and is safe to run repeatedly.
//
//go:embed init.sql
var initSQL string

// initKey marks that this service has been initialized.
const initKey = "initDatabase:job"

type InitDatabaseLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewInitDatabaseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InitDatabaseLogic {
	return &InitDatabaseLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// InitDatabase creates the job tables, inserts the hello_world seed task, and
// registers the job APIs / casbin policies for role 001.
func (l *InitDatabaseLogic) InitDatabase(in *job.Empty) (*job.BaseResp, error) {
	if ok, err := l.svcCtx.Redis.Exists(initKey); err == nil && ok {
		return &job.BaseResp{Msg: "already initialized"}, nil
	}

	if err := l.svcCtx.DB.WithContext(l.ctx).Exec(initSQL).Error; err != nil {
		logx.Errorw("init sql failed", logx.Field("detail", err.Error()))
		return nil, errorx.New(errorx.ServerError, "init database failed: "+err.Error())
	}
	logx.Infow("init sql done")

	// The INSERT IGNORE above is idempotent, but the table may already contain
	// data. Make sure there is exactly one hello_world task (pattern is unique).
	var count int64
	if err := l.svcCtx.DB.WithContext(l.ctx).Model(&model.SysTask{}).
		Where("pattern = ?", pattern.RecordHelloWorld).Count(&count).Error; err != nil {
		return nil, errorx.New(errorx.ServerError, "database error")
	}
	if count == 0 {
		status := int32(1)
		seed := &model.SysTask{
			Status:         &status,
			Name:           "hello_world",
			TaskGroup:      "base",
			CronExpression: "@every 60s",
			Pattern:        pattern.RecordHelloWorld,
			Payload:        `{"name": "Mike (DPTask 60s)"}`,
		}
		if err := l.svcCtx.DB.WithContext(l.ctx).Create(seed).Error; err != nil {
			logx.Errorw("failed to insert hello_world seed task", logx.Field("detail", err.Error()))
			return nil, errorx.New(errorx.ServerError, "database error")
		}
	}

	if err := l.svcCtx.Redis.Set(initKey, "1"); err != nil {
		return nil, errorx.New(errorx.ServerError, "failed to set init flag: "+err.Error())
	}

	return &job.BaseResp{Msg: "success"}, nil
}
