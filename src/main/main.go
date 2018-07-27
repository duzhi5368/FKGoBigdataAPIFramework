package main

import (
	"api_main"
	"cmd_line"
	"github.com/alexflint/go-arg"
	. "slog"
)

func main() {
	arg.MustParse(&cmd_line.G_AppCommandLine)
	Log.Println(cmd_line.G_AppCommandLine.DumpInfo())

	var errChan = make(chan error)
	Log.Println("Current API Server verison: " + api_main.APIVERSION)
	go func() {
		errChan <- api_main.StartServe(cmd_line.G_AppCommandLine)
	}()
	Log.Println("Servers start successfull...")
	Log.Error(<-errChan)
}
