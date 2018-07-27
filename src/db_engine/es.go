package db_engine

import (
	"config"
	"fmt"
	"gopkg.in/olivere/elastic.v5"
	"slog"
	"utils"
)

type ESHostParameters struct {
	EsHost string
	EsPort int
}

func ESEngine(hostParams ESHostParameters) (*elastic.Client, error) {
	var ESHost = utils.If(hostParams.EsHost == "", config.Config.ESHost, hostParams.EsHost).(string)
	var ESPort = utils.If(hostParams.EsPort == 0, config.Config.ESPort, hostParams.EsPort).(int)
	var dsnStr = fmt.Sprintf("http://%s:%d", ESHost, ESPort)
	slog.Log.Info("es dsn: " + dsnStr)
	return elastic.NewClient(elastic.SetURL(dsnStr))
}
