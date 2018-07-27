package slog

import (
	"github.com/qiniu/log"
	"testing"
)

func TestLog(t *testing.T) {
	Log.SetOutputLevel(log.Lerror)
	Log.Println("hello, world")
	Log.Error("error")
	Log.SetOutputLevel(log.Ldebug)
	Log.Println("hello, world")
	Log.Error("error")
	t.Log("success")
}
