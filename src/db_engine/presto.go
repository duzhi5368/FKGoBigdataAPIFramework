package db_engine

import (
	"config"
	"database/sql"
	_ "github.com/prestodb/presto-go-client/presto"
	"slog"
	"strconv"
	"utils"
)

type PrestoHostParameters struct {
	IsUsePrestoDBTLS bool
	PrestoDBHost     string
	PrestoDBPort     int
	PrestoDBUer      string
	PrestoDBPassword string
	PrestoDBCatalog  string
	PrestoDBSchema   string
}

// ! important!
// PrestoGo only support SELECT and SHOW like sql.
func PrestoEngine(hostParams PrestoHostParameters) (*sql.DB, error) {
	var IsUsePrestoDBTLS = config.Config.PrestoDBTLS
	var PrestoDBUser = utils.If(hostParams.PrestoDBUer == "", config.Config.PrestoDBUser, hostParams.PrestoDBUer).(string)
	var PrestoDBPassword = utils.If(hostParams.PrestoDBPassword == "", config.Config.PrestoDBPassword, hostParams.PrestoDBPassword).(string)
	var PrestoDBHost = utils.If(hostParams.PrestoDBHost == "", config.Config.PrestoDBIP, hostParams.PrestoDBHost).(string)
	var PrestoDBPort = utils.If(hostParams.PrestoDBPort == 0, config.Config.PrestoDBPort, hostParams.PrestoDBPort).(int)
	var PrestoDBCatalog = utils.If(hostParams.PrestoDBCatalog == "", config.Config.PrestoDBCatalog, hostParams.PrestoDBCatalog).(string)
	var PrestoDBSchema = utils.If(hostParams.PrestoDBSchema == "", config.Config.PrestoDBSchema, hostParams.PrestoDBSchema).(string)

	var dsnStr = "http://"
	if IsUsePrestoDBTLS {
		dsnStr = "https://"
	}
	dsnStr += PrestoDBUser
	if PrestoDBPassword != "" {
		dsnStr += ":" + config.Config.PrestoDBPassword
	}
	dsnStr += "@"
	dsnStr += PrestoDBHost
	dsnStr += ":"
	dsnStr += strconv.Itoa(PrestoDBPort)
	dsnStr += "?"
	dsnStr += "catalog="
	dsnStr += PrestoDBCatalog
	dsnStr += "&"
	dsnStr += "schema="
	dsnStr += PrestoDBSchema

	slog.Log.Info("presto dsn: " + dsnStr)
	return sql.Open("presto", dsnStr)
}
