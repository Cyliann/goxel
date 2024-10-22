package main

import "runtime"

func main() {
	runtime.LockOSThread()
	app := App{}
	app.Create()
	app.Run()
	defer app.Close()
}
