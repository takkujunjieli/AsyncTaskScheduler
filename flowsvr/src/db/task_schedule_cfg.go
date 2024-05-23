package db

import (
	"time"

	"github.com/jinzhu/gorm"
)

var TaskTypeCfgNsp TaskScheduleCfg

/*
 * @desc: task schedule config
 * @TaskType：
 * @ScheduleLimit：
 * @MaxProcessingTime：
 * @MaxRetryNum:
 * @MaxRetryInterval:
 * @
 */

// TaskScheduleCfg cfg
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

// TableName
func (p *TaskScheduleCfg) TableName() string {
	return "t_schedule_cfg"
}

// Create
func (p *TaskScheduleCfg) Create(db *gorm.DB, task *TaskScheduleCfg) error {
	err := db.Table(p.TableName()).Create(task).Error
	return err
}

// Save
func (p *TaskScheduleCfg) Save(db *gorm.DB, task *TaskScheduleCfg) error {
	err := db.Table(p.TableName()).Save(task).Error
	return err
}

// GetTaskTypeCfg
func (p *TaskScheduleCfg) GetTaskTypeCfg(db *gorm.DB, taskType string) (*TaskScheduleCfg, error) {
	var cfg = new(TaskScheduleCfg)
	err := db.Table(p.TableName()).Where("task_type = ?", taskType).First(&cfg).Error
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

// GetTaskTypeCfgList
func (p *TaskScheduleCfg) GetTaskTypeCfgList(db *gorm.DB) ([]*TaskScheduleCfg, error) {
	var taskTypeCfgList = make([]*TaskScheduleCfg, 0)
	db = db.Table(p.TableName())
	err := db.Find(&taskTypeCfgList).Error
	if err != nil {
		return nil, err
	}
	return taskTypeCfgList, nil
}
