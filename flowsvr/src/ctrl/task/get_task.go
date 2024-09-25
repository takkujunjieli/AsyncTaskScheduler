package task

import (
	"net/http"

	"github.com/takkujunjieli/AsyncTaskScheduler/flowsvr/src/ctrl/ctrlmodel"
	"github.com/takkujunjieli/AsyncTaskScheduler/taskutils/rpc/model"

	"github.com/takkujunjieli/AsyncTaskScheduler/flowsvr/src/constant"
	"github.com/takkujunjieli/AsyncTaskScheduler/flowsvr/src/db"

	"github.com/niuniumart/gosdk/martlog"

	"github.com/niuniumart/gosdk/tools"

	"github.com/gin-gonic/gin"
	"github.com/niuniumart/gosdk/handler"
)

// GetTaskHandler
type GetTaskHandler struct {
	Req    model.GetTaskReq
	Resp   model.GetTaskResp
	UserId string
}

// GetTask
func GetTask(c *gin.Context) {
	var hd GetTaskHandler
	defer func() {
		hd.Resp.Msg = constant.GetErrMsg(hd.Resp.Code)
		c.JSON(http.StatusOK, hd.Resp)
	}()
	//
	hd.UserId = c.Request.Header.Get(constant.HEADER_USERID)
	if err := c.ShouldBind(&hd.Req); err != nil {
		martlog.Errorf("GetTaskHandler shouldBind err %s", err.Error())
		hd.Resp.Code = constant.ERR_SHOULD_BIND
		return
	}
	martlog.Infof("GetTaskHandler Req %s", tools.GetFmtStr(hd.Req))
	handler.Run(&hd)
}

// HandleInput
func (p *GetTaskHandler) HandleInput() error {
	if p.Req.TaskId == "" {
		martlog.Errorf("input invalid")
		p.Resp.Code = constant.ERR_INPUT_INVALID
		return constant.ERR_HANDLE_INPUT
	}
	return nil
}

// HandleProcess
func (p *GetTaskHandler) HandleProcess() error {
	dbTaskData, err := db.TaskNsp.Find(db.DB, p.Req.TaskId)
	if err != nil {
		martlog.Errorf("db.TaskNsp.GetTask %s", err.Error())
		p.Resp.Code = constant.ERR_GET_TASK_INFO
		return err
	}
	var task = &model.TaskData{}
	ctrlmodel.FillTaskResp(dbTaskData, task)
	p.Resp.TaskData = task
	return nil
}
