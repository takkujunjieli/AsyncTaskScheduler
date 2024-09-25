package ctrlmodel

import (
	"github.com/takkujunjieli/AsyncTaskScheduler/flowsvr/src/db"
	"github.com/takkujunjieli/AsyncTaskScheduler/taskutils/rpc/model"
)

/**
 * @Description:
 * @receiver sTask
 * @return *db.Task
 * @return error
 */

// Fill db Task
func FillTaskModel(sTask *model.TaskData, task *db.Task, scheduleEndPosStr string) error {
	task.UserId = sTask.UserId
	if sTask.TaskId == "" {
		task.TaskId = db.TaskNsp.GenTaskId(sTask.TaskType, scheduleEndPosStr)
		task.Status = int(db.TASK_STATUS_PENDING)
	} else {
		task.Status = sTask.Status
	}
	task.TaskType = sTask.TaskType
	task.UserId = sTask.UserId
	task.ScheduleLog = sTask.ScheduleLog
	task.TaskStage = sTask.TaskStage
	task.CrtRetryNum = sTask.CrtRetryNum
	task.MaxRetryNum = sTask.MaxRetryNum
	task.MaxRetryInterval = sTask.MaxRetryInterval
	task.TaskContext = sTask.TaskContext
	return nil
}

/**
 * @Description:
 * @receiver sTask
 * @return *db.Task
 * @return error
 */

// FillTaskResp resp
func FillTaskResp(task *db.Task, sTask *model.TaskData) {
	sTask.UserId = task.UserId
	sTask.TaskId = task.TaskId
	sTask.TaskType = task.TaskType
	sTask.Status = task.Status
	sTask.ScheduleLog = task.ScheduleLog
	var priority = task.Priority
	sTask.Priority = &priority
	sTask.TaskStage = task.TaskStage
	sTask.CrtRetryNum = task.CrtRetryNum
	sTask.MaxRetryNum = task.MaxRetryNum
	sTask.MaxRetryInterval = task.MaxRetryInterval
	sTask.TaskContext = task.TaskContext
	sTask.CreateTime = task.CreateTime
	sTask.ModifyTime = task.ModifyTime
}
