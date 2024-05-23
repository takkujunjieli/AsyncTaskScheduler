package initialize

import (
	"github.com/gin-gonic/gin"
	"github.com/niuniumart/asyncflow/flowsvr/src/ctrl/task"
	"github.com/niuniumart/asyncflow/flowsvr/src/db"
	"github.com/niuniumart/gosdk/martlog"
)

// InitResource 初始化服务资源
func InitResource() error {
	err := InitInfra()
	if err != nil {
		martlog.Errorf("InitInfra err %s", err.Error())
		return err
	}
	return nil
}

// RegisterRouter 注册路由
func RegisterRouter(router *gin.Engine) {
	{
		// 创建任务接口，前面是路径，后面是执行的函数，跳进去
		router.POST("/create_task", task.CreateTask)
		router.POST("/hold_tasks", task.HoldTasks)
		router.GET("/get_task_list", task.GetTaskList)
		router.GET("/get_task_schedule_cfg_list", task.GetTaskScheduleCfgList)
		router.GET("/get_task", task.GetTask)
		router.POST("/set_task", task.SetTask)
		//logprint.RegisterIgnoreRespLogUrl("/get_task_list")
	}
}

// InitInfra 初始化基础设施
func InitInfra() error {
	err := db.InitDB()
	if err != nil {
		return err
	}
	return nil
}
