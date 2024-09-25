package db

import (
	"fmt"
	"strings"

	"github.com/takkujunjieli/AsyncTaskScheduler/flowsvr/src/config"

	"github.com/niuniumart/gosdk/martlog"

	"github.com/jinzhu/gorm"
	"github.com/niuniumart/gosdk/gormcli"
)

var DB *gorm.DB

// InitDB
func InitDB() error {
	var err error
	gormcli.Factory.IdleTimeout = 100
	gormcli.Factory.MaxIdleConn = 1
	DB, err = gormcli.Factory.CreateGorm(config.Conf.MySQL.User,
		config.Conf.MySQL.Pwd, config.Conf.MySQL.Url, config.Conf.MySQL.Dbname)
	if err != nil {
		martlog.Errorf("gormcli.Factory.CreateTBassGorm err %s", err.Error())
		return err
	}
	return nil
}

const (
	GORM_DUPLICATE_ERR_KEY = "Duplicate entry"
)

// IsDupErr
func IsDupErr(err error) bool {
	return strings.Contains(err.Error(), GORM_DUPLICATE_ERR_KEY)
}

// GetTaskTableName
func GetTaskTableName(taskType string) string {
	taskTableName := fmt.Sprintf("t_%s_task", taskType)
	return taskTableName
}
