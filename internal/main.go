package main

import (
	"Cyliann/goxel/internal/app"
	"runtime"

	"github.com/charmbracelet/log"
)

func init() {
	log.SetLevel(log.DebugLevel)
	runtime.LockOSThread()
}

func main() {
	app := app.New()
	app.Run()
	defer app.Close()
}
