package slog

import (
	"config"
	"github.com/qiniu/log"
	"os"
)

var (
	Log = log.New(os.Stdout, "", log.Ldefault)
)

func init() {
	// 调试阶段，打印所有日志
	if config.Config.DebugMode {
		Log.Level = log.Ldebug
	} else {
		Log.Level = log.Linfo
	}
}
