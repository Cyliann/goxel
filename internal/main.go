package main

import (
	"Cyliann/goxel/internal/app"
	"runtime"
)

func main() {
	runtime.LockOSThread()
	app := app.New()
	app.Run()
	defer app.Close()
}
