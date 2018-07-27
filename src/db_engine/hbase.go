package db_engine

import (
	"config"
	"fmt"
	"github.com/sdming/goh"
	"slog"
	"utils"
)

type HBaseHostParameters struct {
	HBaseHost string
	HBasePort int
}

func HBaseEngine(hostParams HBaseHostParameters) (*goh.HClient, error) {
	var HBaseHost = utils.If(hostParams.HBaseHost == "", config.Config.HbaseHost, hostParams.HBaseHost).(string)
	var HBasePort = utils.If(hostParams.HBasePort == 0, config.Config.HbasePort, hostParams.HBasePort).(int)
	var dsnStr = fmt.Sprintf("%s:%d", HBaseHost, HBasePort)
	slog.Log.Info("hbase dsn: " + dsnStr)
	return goh.NewTcpClient(dsnStr, goh.TBinaryProtocol, false)
}
