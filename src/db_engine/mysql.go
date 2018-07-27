package db_engine

import (
	"config"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"slog"
	"utils"
)

type MysqlHostParameters struct {
	MysqlHost     string
	MysqlPort     int
	MysqlDBName   string
	MysqlUser     string
	MysqlPassword string
}

func MysqlEngine(hostParams MysqlHostParameters) (*xorm.Engine, error) {
	var MysqlHost = utils.If(hostParams.MysqlHost == "", config.Config.MysqlHost, hostParams.MysqlHost).(string)
	var MysqlPort = utils.If(hostParams.MysqlPort == 0, config.Config.MysqlPort, hostParams.MysqlPort).(int)
	var MysqlDBName = utils.If(hostParams.MysqlDBName == "", config.Config.MysqlDBName, hostParams.MysqlDBName).(string)
	var MysqlUser = utils.If(hostParams.MysqlUser == "", config.Config.MysqlUser, hostParams.MysqlUser).(string)
	var MysqlPassword = utils.If(hostParams.MysqlPassword == "", config.Config.MysqlHost, hostParams.MysqlPassword).(string)

	var dsnStr = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", MysqlUser, MysqlPassword, MysqlHost, MysqlPort, MysqlDBName)
	slog.Log.Info("mysql dsn: " + dsnStr)
	engine, err := xorm.NewEngine("mysql", dsnStr)

	if err != nil {
		return nil, err
	}
	engine.ShowSQL(config.Config.DebugMode)
	return engine, nil
}
