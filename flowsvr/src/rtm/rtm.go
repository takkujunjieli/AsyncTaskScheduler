package rtm

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/google/uuid"
	"github.com/niuniumart/gosdk/martlog"
	"github.com/niuniumart/gosdk/requestid"
	"github.com/takkujunjieli/AsyncTaskScheduler/flowsvr/src/config"
	"github.com/takkujunjieli/AsyncTaskScheduler/flowsvr/src/db"
)

// TaskRuntime 短任务运行时
type TaskRuntime struct {
}

// Run 开始运行
func (p *TaskRuntime) Run() {
	go p.run()
}

func (p *TaskRuntime) run() {
	/******  dealLongTimeProcess *******/
	go func() {
		defer func() {
			if err := recover(); err != nil {
				martlog.Errorf("WatTaskRuntime PanicRecover,Error:%s", err)
				//打印调用栈信息
				debug.PrintStack()
				buf := make([]byte, 2048)
				n := runtime.Stack(buf, false)
				stackInfo := fmt.Sprintf("%s", buf[:n])
				martlog.Errorf("panic stack info %s\n", stackInfo)

			}
			p.dealLongTimeProcess()
		}()
		requestIDStr := fmt.Sprintf("%+v", uuid.New())
		requestid.Set(requestIDStr)
		p.dealLongTimeProcess()
	}()
}

func (p *TaskRuntime) dealLongTimeProcess() {
	for {
		martlog.Infof("short task deal long time process")
		t := time.NewTimer(time.Duration(config.Conf.Task.LongProcessInterval))
		<-t.C
		/***** step 1: do get lock   *****/
		//lockKey := SHORT_TASK_LONGTIME_DEAL_LOCK_KEY
		/***** step 2: deal long process do  *****/
		martlog.Infof("schedule do dealTimeoutProcessing")
		p.dealTimeoutProcessing()
		/***** step 3: do unlock *****/
	}
}

func (p *TaskRuntime) dealTimeoutProcessing() {
	taskTypeCfgList, err := db.TaskTypeCfgNsp.GetTaskTypeCfgList(db.DB)
	if err != nil {
		martlog.Errorf("visit t_task_type_cfg err %s", err.Error())
		return
	}
	for _, taskTypeCfg := range taskTypeCfgList {
		p.dealTimeoutProcessingWithType(taskTypeCfg)
	}
}

func (p *TaskRuntime) dealTimeoutProcessingWithType(taskCfg *db.TaskScheduleCfg) {
	taskTableName := db.GetTaskTableName(taskCfg.TaskType)
	taskPos, err := db.TaskPosNsp.GetTaskPos(db.DB, taskTableName)
	if err != nil {
		martlog.Errorf("db.TaskPosNsp.GetTaskPos err %s", err.Error())
		return
	}
	for i := taskPos.ScheduleBeginPos; i <= taskPos.ScheduleEndPos; i++ {
		maxProcessTime := config.Conf.Task.MaxProcessTime
		if int64(taskCfg.MaxProcessingTime) != 0 {
			maxProcessTime = taskCfg.MaxProcessingTime
		}
		taskList, err := db.TaskNsp.GetLongTimeProcessing(db.DB, taskCfg.TaskType,
			fmt.Sprintf("%d", i), maxProcessTime, 1000)
		if err != nil {
			martlog.Errorf("get long time processing err %s", err.Error())
			continue
		}
		for _, ts := range taskList {
			err = db.TaskNsp.SetStatus(db.DB, ts.TaskId, db.TASK_STATUS_PENDING)
			if err != nil {
				martlog.Errorf("deal long time task, save task err %s", err.Error())
				continue
			}
		}
	}
}
