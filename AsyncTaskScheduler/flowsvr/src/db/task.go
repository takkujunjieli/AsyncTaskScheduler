package db

import (
	"fmt"
	"strings"
	"time"

	"github.com/niuniumart/gosdk/martlog"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

var TaskNsp Task

const (
	BEGIN_TABLE_POS = 1
)

// Task
type Task struct {
	Id               uint64
	UserId           string
	TaskId           string
	TaskType         string
	TaskStage        string
	Status           int
	Priority         int
	CrtRetryNum      int
	MaxRetryNum      int
	MaxRetryInterval int
	ScheduleLog      string
	TaskContext      string // task struct json string
	OrderTime        int64
	CreateTime       *time.Time
	ModifyTime       *time.Time
}

func (p *Task) getTableName(taskType, pos string) string {
	return fmt.Sprintf("t_%s_%s_%s", taskType, p.TableName(), pos)
}

// TableName
func (p *Task) TableName() string {
	return "task"
}

// GenTaskId
func (p *Task) GenTaskId(taskType, pos string) string {
	taskType = strings.Replace(taskType, "_", "-", -1)
	return fmt.Sprintf("%+v_%s_%s", uuid.New(), taskType, pos)
}
func (p *Task) getTablePosFromTaskId(taskId string) (string, string) {
	s := strings.Split(taskId, "_")
	switch len(s) {
	case 3:
		s[1] = strings.Replace(s[1], "-", "_", -1)
		return s[1], s[2]
	default:
		martlog.Errorf("big error taskId %s have not match _", taskId)
		return "", ""
	}
}

// BatchSetStatus batch set
func (p *Task) BatchSetStatus(db *gorm.DB,
	taskIdList []string, status TaskEnum) error {
	var dic = map[string]interface{}{
		"status": status,
	}
	tmpTaskId := taskIdList[0]
	taskType, pos := p.getTablePosFromTaskId(tmpTaskId)
	db = db.Table(p.getTableName(taskType, pos)).Where("task_id in (?)", taskIdList).
		UpdateColumns(dic)
	err := db.Error
	if err != nil {
		return err
	}
	return nil
}

// CreateNextTable
func (p *Task) CreateNextTable(db *gorm.DB, taskType, pos string) error {
	newTable := p.getTableName(taskType, pos)
	beginTable := p.getTableName(taskType, fmt.Sprintf("%d", BEGIN_TABLE_POS))
	return db.Exec(fmt.Sprintf("create table %s like %s", newTable, beginTable)).Error
}

// Find
func (p *Task) Find(db *gorm.DB, taskId string) (*Task, error) {
	var data = &Task{}
	taskType, pos := p.getTablePosFromTaskId(taskId)
	err := db.Table(p.getTableName(taskType, pos)).Where("task_id = ?", taskId).First(data).Error
	return data, err
}

// Create
func (p *Task) Create(db *gorm.DB, taskType, pos string, task *Task) error {
	err := db.Table(p.getTableName(taskType, pos)).Create(task).Error
	return err
}

// Save
func (p *Task) Save(db *gorm.DB, task *Task) error {
	taskType, pos := p.getTablePosFromTaskId(task.TaskId)
	err := db.Table(p.getTableName(taskType, pos)).Save(task).Error
	return err
}

// GetTaskList
func (p *Task) GetTaskList(db *gorm.DB,
	pos string, taskType string, status TaskEnum, limit int) ([]*Task, error) {
	var taskList = make([]*Task, 0)
	err := db.
		Table(p.getTableName(taskType, pos)).
		Where("status = ?", status).
		Order("order_time").
		Limit(limit).
		Find(&taskList).Error
	if err != nil {
		return nil, err
	}
	return taskList, nil
}

// GetAliveTaskList
func (p *Task) GetAliveTaskList(db *gorm.DB, taskType, pos string, limit int) ([]*Task, error) {
	var taskList = make([]*Task, 0)
	var statusSet = []TaskEnum{TASK_STATUS_PENDING, TASK_STATUS_PROCESSING}
	err := db.
		Table(p.getTableName(taskType, pos)).
		Order("modify_time").
		Limit(limit).
		Where("status in (?)", statusSet).
		Find(&taskList).Error
	if err != nil {
		return nil, err
	}
	return taskList, nil
}

// GetAliveTaskCount
func (p *Task) GetAliveTaskCount(db *gorm.DB, taskType, pos string) (int, error) {
	return p.getTaskCount(db, taskType, pos,
		[]TaskEnum{TASK_STATUS_PENDING, TASK_STATUS_PROCESSING})
}

// GetAllTaskCount
func (p *Task) GetAllTaskCount(db *gorm.DB, taskType, pos string) (int, error) {
	return p.getTaskCount(db, taskType, pos,
		[]TaskEnum{TASK_STATUS_PENDING, TASK_STATUS_PROCESSING,
			TASK_STATUS_FAILED, TASK_STATUS_SUCC})
}

// GetTaskCountByStatus
func (p *Task) GetTaskCountByStatus(db *gorm.DB, taskType, pos string, status int) (int, error) {
	var count int
	err := db.Table(p.getTableName(taskType, pos)).Where("status = ?", status).Count(&count).Error
	if err != nil {
		return count, err
	}
	return count, nil
}
func (p *Task) getTaskCount(db *gorm.DB, taskType, pos string, statusSet []TaskEnum) (int, error) {
	var count int
	err := db.Table(p.getTableName(taskType, pos)).Where("status in (?)", statusSet).Count(&count).Error
	if err != nil {
		return count, err
	}
	return count, nil
}

// SetStatusPending
func (p *Task) SetStatusPending(db *gorm.DB, taskId string) error {
	return p.SetStatus(db, taskId, TASK_STATUS_PENDING)
}

// SetStatusSucc
func (p *Task) SetStatusSucc(db *gorm.DB, taskId string) error {
	return p.SetStatus(db, taskId, TASK_STATUS_SUCC)
}

// SetStatusFailed
func (p *Task) SetStatusFailed(db *gorm.DB, taskId string) error {
	return p.SetStatus(db, taskId, TASK_STATUS_FAILED)
}

// SetStatus
func (p *Task) SetStatus(db *gorm.DB, taskId string, status TaskEnum) error {
	var dic = map[string]interface{}{
		"status": status,
	}
	taskType, pos := p.getTablePosFromTaskId(taskId)
	err := db.Table(p.getTableName(taskType, pos)).Where("task_id = ?", taskId).Updates(dic).Error
	if err != nil {
		return err
	}
	return nil
}

// SetStatusWithOutModifyTime
func (p *Task) SetStatusWithOutModifyTime(db *gorm.DB, taskId string, status TaskEnum) error {
	taskType, pos := p.getTablePosFromTaskId(taskId)
	err := db.Table(p.getTableName(taskType, pos)).Where("task_id = ?", taskId).UpdateColumn("status", status).Error
	if err != nil {
		return err
	}
	return nil
}

// SetContext
func (p *Task) SetContext(db *gorm.DB, taskId, context string) error {
	var dic = map[string]interface{}{
		"task_context": context,
	}
	taskType, pos := p.getTablePosFromTaskId(taskId)
	err := db.Table(p.getTableName(taskType, pos)).Where("task_id = ?", taskId).Updates(dic).Error
	if err != nil {
		return err
	}
	return nil
}

// GetLongTimeProcessing
func (p *Task) GetLongTimeProcessing(db *gorm.DB,
	taskType, pos string, maxProcessTime int64, limit int) ([]*Task, error) {
	var Tasks = make([]*Task, 0)
	err := db.Table(p.getTableName(taskType, pos)).
		Where("status = ?", TASK_STATUS_PROCESSING).
		Where("unix_timestamp(modify_time) + ? < ?", maxProcessTime, time.Now().Unix()).
		Limit(limit).
		Find(&Tasks).
		Error
	if err != nil {
		return nil, err
	}
	return Tasks, nil
}

// ModifyTimeoutPending
func (p *Task) ModifyTimeoutPending(db *gorm.DB, taskType, pos string, maxProcessTime int64) error {
	var dic = map[string]interface{}{
		"status": TASK_STATUS_PENDING,
	}
	err := db.Table(p.getTableName(taskType, pos)).
		Where("status = ?", TASK_STATUS_PROCESSING).
		Where("unix_timestamp(modify_time) + ? < ?", maxProcessTime, time.Now().Unix()).
		Updates(dic).
		Error
	if err != nil {
		return err
	}
	return nil
}

// IncreaseCrtRetryNum
func (p *Task) IncreaseCrtRetryNum(db *gorm.DB, taskId string) error {
	taskType, pos := p.getTablePosFromTaskId(taskId)
	return db.Table(p.getTableName(taskType, pos)).
		Where("task_id = ?", taskId).
		Update("crt_retry_num", gorm.Expr("crt_retry_num + ?", 1)).Error
}

// BeforeCreate
func (p *Task) BeforeCreate(scope *gorm.Scope) error {
	now := time.Now()
	scope.SetColumn("create_time", now)
	scope.SetColumn("modify_time", now)
	return nil
}

// UpdateTask
func (p *Task) UpdateTask(db *gorm.DB) error {
	taskType, pos := p.getTablePosFromTaskId(p.TaskId)
	err := db.Table(p.getTableName(taskType, pos)).Where("task_id = ?", p.TaskId).
		Where("status <> ? and status <> ?", TASK_STATUS_SUCC, TASK_STATUS_FAILED).Updates(p).Error
	if err != nil {
		return err
	}
	return nil
}

// SetScheduleLog
func (p *Task) SetScheduleLog(db *gorm.DB, ScheduleLog string) error {
	p.ScheduleLog = ScheduleLog
	return p.UpdateTask(db)
}

func (p *Task) BatchSetOwnerStatusWithPendingOutModify(db *gorm.DB,
	taskIdList []string, owner string, status TaskEnum) (int64, error) {
	var dic = map[string]interface{}{
		"status": status,
	}
	if owner != "" {
		dic["owner"] = owner
	}
	tmpTaskId := taskIdList[0]
	taskType, pos := p.getTablePosFromTaskId(tmpTaskId)
	db = db.Table(p.getTableName(taskType, pos)).Where("task_id in (?)", taskIdList).
		Where("status = ?", TASK_STATUS_PENDING).UpdateColumns(dic)
	err := db.Error
	if err != nil {
		return 0, err
	}
	return db.RowsAffected, nil
}

// GetAssignTasksByOwnerStatus
func (p *Task) GetAssignTasksByOwnerStatus(db *gorm.DB,
	taskIdList []string, owner string, status TaskEnum, limit int64) ([]*Task, error) {
	if len(taskIdList) == 0 {
		martlog.Infof("taskId list is empty")
		return nil, nil
	}
	var Tasks = make([]*Task, 0)
	tmpTaskId := taskIdList[0]
	taskType, pos := p.getTablePosFromTaskId(tmpTaskId)
	err := db.Table(p.getTableName(taskType, pos)).
		Where("task_id in (?)", taskIdList).
		Where("owner = ? and status = ?", owner, status).
		Limit(limit).
		Find(&Tasks).
		Error
	if err != nil {
		return nil, err
	}
	return Tasks, nil
}

// ConventTaskIdList
func ConventTaskIdList(tasks []*Task) []string {
	taskIds := make([]string, 0, len(tasks))
	for _, task := range tasks {
		if task != nil {
			taskIds = append(taskIds, task.TaskId)
		}
	}
	return taskIds
}
