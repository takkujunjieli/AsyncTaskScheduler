package model

import (
	"time"
)

// CreateTaskReq
type CreateTaskReq struct {
	TaskData TaskData `json:"taskData"`
}

// RespComm
type RespComm struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// CreateTaskResp
type CreateTaskResp struct {
	RespComm
	TaskId string `json:"taskId"`
}

// GetTaskListReq
type GetTaskListReq struct {
	TaskType string `json:"taskType" form:"taskType"`
	Status   int    `json:"status" form:"status"`
	Limit    int    `json:"limit" form:"limit"`
}

// GetTaskListResp
type GetTaskListResp struct {
	RespComm
	TaskList []*TaskData `json:"taskList"`
}

// GetTaskListReq
type HoldTasksReq struct {
	TaskType string `json:"taskType" form:"taskType"`
	Limit    int    `json:"limit" form:"limit"`
}

// GetTaskListResp
type HoldTasksResp struct {
	RespComm
	TaskList []*TaskData `json:"taskList"`
}

// GetTaskReq
type GetTaskReq struct {
	TaskId string `json:"taskId" form:"taskId"`
}

// GetTaskResp
type GetTaskResp struct {
	RespComm
	TaskData *TaskData `json:"taskData"`
}

// GetTaskCountByStatusReq
type GetTaskCountByStatusReq struct {
	TaskType string `json:"taskType" form:"taskType"`
	Status   int    `json:"status" form:"status"`
}

// GetTaskCountByStatusResp
type GetTaskCountByStatusResp struct {
	RespComm
	Count int `json:"count"`
}

// GetTaskScheduleCfgListReq
type GetTaskScheduleCfgListReq struct {
}

// GetTaskScheduleCfgListResp
type GetTaskScheduleCfgListResp struct {
	RespComm
	ScheduleCfgList []*TaskScheduleCfg `json:"scheduleCfgList"`
}

// TaskScheduleCfg
type TaskScheduleCfg struct {
	TaskType          string
	ScheduleLimit     int
	ScheduleInterval  int
	MaxProcessingTime int64
	MaxRetryNum       int
	MaxRetryInterval  int
	CreateTime        *time.Time
	ModifyTime        *time.Time
}

// SetTaskStatusReq
type SetTaskStatusReq struct {
	TaskId       string `json:"taskId"`
	Status       int    `json:"status"`
	NoModifyTime bool   `json:"noModifyTime"`
}

// SetTaskStatusResp
type SetTaskStatusResp struct {
	RespComm
}

// SetTaskReq
type SetTaskReq struct {
	TaskId   string `json:"taskId"`
	TaskData `json:"TaskData"`
	Context  string `json:"context"`
}

// SetTaskResp
type SetTaskResp struct {
	RespComm
}

// TaskData
type TaskData struct {
	UserId           string     `json:"userId"`
	TaskId           string     `json:"taskId"`
	TaskType         string     `json:"taskType"`
	TaskStage        string     `json:"taskStage"`
	Status           int        `json:"status"`
	Priority         *int       `json:"priority"`
	CrtRetryNum      int        `json:"crtRetryNum"`
	MaxRetryNum      int        `json:"maxRetryNum"`
	MaxRetryInterval int        `json:"maxRetryInterval"`
	ScheduleLog      string     `json:"scheduleLog"`
	TaskContext      string     `json:"context"`
	OrderTime        int64      `json:"orderTime"`
	CreateTime       *time.Time `json:"createTime"`
	ModifyTime       *time.Time `json:"modifyTime"`
}
