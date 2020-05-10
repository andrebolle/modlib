package main

import (
	"runtime"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/purelazy/modlib/internal/window"
)

const (
	windowWidth  = 500
	windowHeight = 500
)

// App App
type App struct {
	WinW, WinH int
	Win        *glfw.Window
}

func init() {
	runtime.LockOSThread()
}

func main() {

	app := App{500, 500, window.NewWindow(windowWidth, windowHeight)}
	defer glfw.Terminate()

	for !app.Win.ShouldClose() {
		// TODO
		glfw.PollEvents()
	}
}
