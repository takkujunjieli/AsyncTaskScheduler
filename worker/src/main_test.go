package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/niuniumart/gosdk/tools"
	"github.com/smartystreets/goconvey/convey"
	"github.com/takkujunjieli/AsyncTaskScheduler/flowsvr/src/config"
	"github.com/takkujunjieli/AsyncTaskScheduler/flowsvr/src/db"
	"github.com/takkujunjieli/AsyncTaskScheduler/taskutils/rpc"
	"github.com/takkujunjieli/AsyncTaskScheduler/taskutils/rpc/model"
)

func TestCreateTask(t *testing.T) {
	config.TestFilePath = "../config/config-test.toml"
	config.Init()
	db.InitDB()
	fmt.Println("aaa   ************* ")
	convey.Convey("TestCreateTask", t, func() {
		// case 1: input err
		var rpc rpc.TaskRpc
		rpc.Host = "http://127.0.0.1:41555"
		var reqBody = new(model.CreateTaskReq)
		reqBody.TaskData.TaskType = "lark"
		reqBody.TaskData.TaskStage = "sendmsg"
		reqBody.TaskData.UserId = "niuniu"
		reqBody.TaskData.Status = 1
		var ltctx = LarkTaskContext{
			ReqBody: &LarkReq{Msg: "nice to meet u", FromAddr: "fish", ToAddr: "cat"},
		}
		ltctxStr, _ := json.Marshal(ltctx)
		reqBody.TaskData.TaskContext = string(ltctxStr)
		resp, err := rpc.CreateTask(reqBody)
		fmt.Println(tools.GetFmtStr(resp))
		fmt.Println(err)
		convey.So(err, convey.ShouldBeNil)

	})
}
