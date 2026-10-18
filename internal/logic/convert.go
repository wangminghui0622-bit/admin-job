package logic

import (
	"github.com/wangminghui0622-bit/admin-common/dal/model"
	"github.com/wangminghui0622-bit/admin-job/job"
)

func pointer[T any](v T) *T {
	return &v
}

func taskInfoFromModel(v *model.SysTask) *job.TaskInfo {
	var status uint32
	if v.Status != nil {
		status = uint32(*v.Status)
	}

	return &job.TaskInfo{
		Id:             pointer(uint64(v.ID)),
		CreatedAt:      pointer(v.CreatedAt.UnixMilli()),
		UpdatedAt:      pointer(v.UpdatedAt.UnixMilli()),
		Status:         pointer(status),
		Name:           pointer(v.Name),
		TaskGroup:      pointer(v.TaskGroup),
		CronExpression: pointer(v.CronExpression),
		Pattern:        pointer(v.Pattern),
		Payload:        pointer(v.Payload),
	}
}

func taskLogInfoFromModel(v *model.SysTaskLog) *job.TaskLogInfo {
	return &job.TaskLogInfo{
		Id:         pointer(uint64(v.ID)),
		StartedAt:  pointer(v.StartedAt.UnixMilli()),
		FinishedAt: pointer(v.FinishedAt.UnixMilli()),
		Result:     pointer(uint32(v.Result)),
	}
}
