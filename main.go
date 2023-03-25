package main

import (
	"fmt"
	"syscall"
	"embed"

	_ "github.com/ihucos/counter.dev/endpoints"
	"github.com/ihucos/counter.dev/lib"
)


//go:embed static
var staticFS embed.FS

func main() {

	// HOTFIX
	var rLimit syscall.Rlimit
	rLimit.Max = 100307
	rLimit.Cur = 100307
	err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		fmt.Println("Error Setting Rlimit ", err)
	}

	app := lib.NewApp()
	app.ConnectEndpoints(staticFS)
	app.Logger.Println("Start")
	app.Serve()
}
