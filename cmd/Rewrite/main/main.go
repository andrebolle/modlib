package main

import (
	"runtime"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/purelazy/modlib/cmd/Rewrite/main/app"
	"github.com/purelazy/modlib/cmd/Rewrite/main/input"
	"github.com/purelazy/modlib/cmd/Rewrite/main/window"
)

const (
	windowWidth  = 500
	windowHeight = 500
)

func init() {
	runtime.LockOSThread()
}

func main() {

	app := new(app.App)
	app.Win = window.NewWindow(windowWidth, windowHeight)
	defer glfw.Terminate()

	for !app.Win.ShouldClose() {

		input.PollKeyboard(app)

		// TODO
		glfw.PollEvents()
	}
}
