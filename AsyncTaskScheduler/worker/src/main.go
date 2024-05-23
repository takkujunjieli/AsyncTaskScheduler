package main

import (
	"encoding/json"
	"fmt"

	"github.com/niuniumart/asyncflow/taskutils/constant"
	"github.com/niuniumart/asyncflow/worker/src/initialise"
	"github.com/niuniumart/asyncflow/worker/src/tasksdk"
	"github.com/niuniumart/gosdk/martlog"
	"github.com/niuniumart/gosdk/response"
	"github.com/niuniumart/gosdk/tools"
)

func main() {
	larkTask := tasksdk.TaskHandler{
		TaskType: "lark",
		NewProc:  func() tasksdk.TaskIntf { return new(LarkTask) },
	}
	tasksdk.RegisterHandler(&larkTask)
	initialise.InitResource()
	tasksdk.InitSvr("http://127.0.0.1:41555", "")
	var taskMgr = tasksdk.TaskMgr{
		TaskType: "lark",
	}
	taskMgr.Schedule()
}

type LarkReq struct {
	Msg      string
	FromAddr string
	ToAddr   string
}

// LarkTask
type LarkTask struct {
	tasksdk.TaskBase
	ContextData *LarkTaskContext
}

// LartTaskContext
type LarkTaskContext struct {
	ReqBody *LarkReq
	UserId  string
}

// ContextLoad
func (p *LarkTask) ContextLoad() error {
	martlog.Infof("run lark task %s", p.TaskId)
	err := json.Unmarshal([]byte(p.TaskContext), &p.ContextData)
	if err != nil {
		martlog.Errorf("json unmarshal for context err %s", err.Error())
		return response.RESP_JSON_UNMARSHAL_ERROR
	}
	if p.ContextData.ReqBody == nil {
		p.ContextData.ReqBody = new(LarkReq)
	}
	return nil
}

// HandleProcess
func (p *LarkTask) HandleProcess() error {
	fmt.Println("task ", tools.GetFmtStr(*p))
	switch p.TaskStage {
	case "sendmsg":
		p.TaskStage = "record"
		p.SetContextLocal(p.ContextData)
		fallthrough
	case "record":
		fmt.Println("come here")
		p.TaskStage = "record"
		p.Base().Status = int(constant.TASK_STATUS_SUCC)

	default:
		p.Base().Status = int(constant.TASK_STATUS_FAILED)
	}
	return nil
}

// HandleFinish
func (p *LarkTask) HandleFinish() {
}
