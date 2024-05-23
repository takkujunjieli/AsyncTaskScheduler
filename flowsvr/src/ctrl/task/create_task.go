package task

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/niuniumart/asyncflow/flowsvr/src/constant"
	"github.com/niuniumart/asyncflow/flowsvr/src/ctrl/ctrlmodel"
	"github.com/niuniumart/asyncflow/flowsvr/src/db"
	"github.com/niuniumart/asyncflow/taskutils/rpc/model"

	"github.com/gin-gonic/gin"
	"github.com/niuniumart/gosdk/handler"
	"github.com/niuniumart/gosdk/martlog"
)

// CreateTaskHandler
type CreateTaskHandler struct {
	Req    model.CreateTaskReq
	Resp   model.CreateTaskResp
	UserId string
}

// CreateTask
func CreateTask(c *gin.Context) {
	var hd CreateTaskHandler
	defer func() {
		hd.Resp.Msg = constant.GetErrMsg(hd.Resp.Code)
		c.JSON(http.StatusOK, hd.Resp)
	}()
	// get user id
	hd.UserId = c.Request.Header.Get(constant.HEADER_USERID)
	//
	if err := c.ShouldBind(&hd.Req); err != nil {
		martlog.Errorf("CreateTask shouldBind err %s", err.Error())
		hd.Resp.Code = constant.ERR_SHOULD_BIND
		return
	}
	//
	handler.Run(&hd)
}

// HandleInput
func (p *CreateTaskHandler) HandleInput() error {
	if p.Req.TaskData.TaskType == "" {
		martlog.Errorf("input invalid")
		p.Resp.Code = constant.ERR_INPUT_INVALID
		return constant.ERR_HANDLE_INPUT
	}
	if p.Req.TaskData.Priority != nil {
		if *p.Req.TaskData.Priority > db.MAX_PRIORITY || *p.Req.TaskData.Priority < 0 {
			p.Resp.Code = constant.ERR_INPUT_INVALID
			martlog.Errorf("input invalid")
			return constant.ERR_HANDLE_INPUT
		}
	}
	return nil
}

// HandleProcess
func (p *CreateTaskHandler) HandleProcess() error {
	martlog.Infof("into HandleProcess")
	var err error
	var taskTableName string
	//
	//
	var taskPos *db.TaskPos
	taskTableName = db.GetTaskTableName(p.Req.TaskData.TaskType)
	taskPos, err = db.TaskPosNsp.GetTaskPos(db.DB, taskTableName)
	if err != nil {
		p.Resp.Code = constant.ERR_GET_TASK_POS
		martlog.Errorf("db.TaskPosNsp.GetTaskPos err: %s", err.Error())
		return err
	}
	if taskPos == nil {
		martlog.Errorf("db.TaskPosNsp.GetTaskPos failed. TaskTableName : %s", taskTableName)
		return errors.New("Get task pos failed.  TaskTableName : " + taskTableName)
	}
	taskCfg, err := db.TaskTypeCfgNsp.GetTaskTypeCfg(db.DB, p.Req.TaskData.TaskType)
	if err != nil {
		p.Resp.Code = constant.ERR_GET_TASK_SET_POS_FROM_DB
		martlog.Errorf("visit t_task_type_cfg err %s", err.Error())
		return err
	}
	scheduleEndPosStr := fmt.Sprintf("%d", taskPos.ScheduleEndPos)
	if err != nil {
		martlog.Errorf("db.TaskPosNsp.GetTaskPos %s", err.Error())
		return err
	}
	var task = new(db.Task)
	p.Req.TaskData.MaxRetryNum = taskCfg.MaxRetryNum
	p.Req.TaskData.MaxRetryInterval = taskCfg.MaxRetryInterval
	//
	p.Req.TaskData.OrderTime = time.Now().Unix()
	if p.Req.TaskData.Priority != nil {
		p.Req.TaskData.OrderTime -= int64(*p.Req.TaskData.Priority)
	}
	//
	err = ctrlmodel.FillTaskModel(&p.Req.TaskData, task, scheduleEndPosStr)
	if err != nil {
		p.Resp.Code = constant.ERR_CREATE_TASK
		martlog.Errorf("db.TaskPosNsp.GetTaskPos %s", err.Error())
		return err
	}
	//
	err = db.TaskNsp.Create(db.DB, p.Req.TaskData.TaskType, scheduleEndPosStr, task)
	if err != nil {
		martlog.Errorf("db.TaskNsp.Create %s", err.Error())
		p.Resp.Code = constant.ERR_CREATE_TASK
		return err
	}
	//
	p.Resp.TaskId = task.TaskId
	return nil
}
