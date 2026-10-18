// Package periodicconfig provides the asynq periodic task configurations from the database.
package periodicconfig

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/wangminghui0622-bit/admin-common/dal/model"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

// statusNormal is the value of sys_tasks.status that means the task is enabled.
const statusNormal = 1

// GormConfigProvider reads task configurations from sys_tasks with gorm.
// The PeriodicTaskManager calls GetConfigs every SyncInterval, so changes made
// through the CRUD RPCs take effect within one sync interval.
type GormConfigProvider struct {
	DB *gorm.DB
}

func NewGormConfigProvider(db *gorm.DB) *GormConfigProvider {
	return &GormConfigProvider{DB: db}
}

func (p *GormConfigProvider) GetConfigs() ([]*asynq.PeriodicTaskConfig, error) {
	var tasks []*model.SysTask
	if err := p.DB.WithContext(context.Background()).
		Where("status = ?", statusNormal).
		Find(&tasks).Error; err != nil {
		// Keep the periodic task manager alive before the database is initialized:
		// the next sync will pick the tasks up once the tables exist.
		logx.Errorw("failed to query tasks for periodic task manager",
			logx.Field("detail", err.Error()),
			logx.Field("recommend", "you may need to initialize the database"))
		return nil, nil
	}

	configs := make([]*asynq.PeriodicTaskConfig, 0, len(tasks))
	for _, v := range tasks {
		configs = append(configs, &asynq.PeriodicTaskConfig{
			Cronspec: v.CronExpression,
			Task:     asynq.NewTask(v.Pattern, []byte(v.Payload)),
		})
	}

	return configs, nil
}
