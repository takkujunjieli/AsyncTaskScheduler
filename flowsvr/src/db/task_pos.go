package db

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/niuniumart/gosdk/martlog"
	"github.com/niuniumart/gosdk/tools"

	"github.com/jinzhu/gorm"
)

var TaskPosNsp TaskPos

// TaskPos taskpos
type TaskPos struct {
	Id               uint64
	TaskType         string
	ScheduleBeginPos int
	ScheduleEndPos   int
	CreateTime       *time.Time
	ModifyTime       *time.Time
}

// TableName
func (p *TaskPos) TableName() string {
	return "t_schedule_pos"
}

// Create
func (p *TaskPos) Create(db *gorm.DB, task *TaskPos) error {
	err := db.Table(p.TableName()).Create(task).Error
	return err
}

// Save
func (p *TaskPos) Save(db *gorm.DB, task *TaskPos) error {
	err := db.Table(p.TableName()).Save(task).Error
	return err
}

// GetTaskPos
func (p *TaskPos) GetTaskPos(db *gorm.DB, taskSetName string) (*TaskPos, error) {
	var taskPos = new(TaskPos)
	err := db.Table(p.TableName()).Where("task_type = ?", taskSetName).First(&taskPos).Error
	if err != nil {
		return nil, err
	}
	return taskPos, nil
}

// GetRandomSchedulePos
func (p *TaskPos) GetRandomSchedulePos(db *gorm.DB, taskSetName string) (int, error) {
	taskPos, err := p.GetTaskPos(db, taskSetName)
	if err != nil {
		return 0, err
	}
	martlog.Infof("taskPos %s", tools.GetFmtStr(taskPos))
	base := taskPos.ScheduleEndPos - taskPos.ScheduleBeginPos + 1
	pos := rand.Intn(base) + taskPos.ScheduleBeginPos
	martlog.Infof("random schedule pos %d", pos)
	return int(pos), nil
}

// GetBeginSchedulePos
func (p *TaskPos) GetBeginSchedulePos(db *gorm.DB, taskSetName string) (int, error) {
	taskPos, err := p.GetTaskPos(db, taskSetName)
	if err != nil {
		return 0, err
	}
	martlog.Infof("taskPos %s", tools.GetFmtStr(taskPos))
	return int(taskPos.ScheduleBeginPos), nil
}

// GetNextPos
func (p *TaskPos) GetNextPos(pos string) string {
	posInt, err := strconv.Atoi(pos)
	if err != nil {
		martlog.Errorf("pos %s maybe not int", pos)
		return ""
	}
	return fmt.Sprintf("%d", posInt+1)
}

// GetTaskPosList
func (p *TaskPos) GetTaskPosList(db *gorm.DB) ([]*TaskPos, error) {
	var taskList = make([]*TaskPos, 0)
	err := db.Table(p.TableName()).Find(&taskList).Error
	if err != nil {
		return nil, err
	}
	return taskList, nil
}

// BeforeCreate
func (this *TaskPos) BeforeCreate(scope *gorm.Scope) error {
	now := time.Now()
	scope.SetColumn("create_time", now)
	scope.SetColumn("modify_time", now)
	return nil
}
