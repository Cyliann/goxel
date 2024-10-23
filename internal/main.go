package main

import (
	"Cyliann/goxel/internal/app"
	"runtime"

	"github.com/charmbracelet/log"
)

func main() {
	log.SetLevel(log.DebugLevel)
	runtime.LockOSThread()
	app := app.New()
	app.Run()
	defer app.Close()
}
