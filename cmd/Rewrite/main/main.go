package main

import (
	"runtime"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/purelazy/modlib/cmd/Rewrite/main/app"
	"github.com/purelazy/modlib/cmd/Rewrite/main/geo"
	"github.com/purelazy/modlib/cmd/Rewrite/main/input"
	"github.com/purelazy/modlib/cmd/Rewrite/main/shader"
	"github.com/purelazy/modlib/cmd/Rewrite/main/vao"
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

	game := app.NewApp()
	game.Win = window.NewWindow(windowWidth, windowHeight)
	defer glfw.Terminate()

	app.StartOpenGL()

	vao := vao.NewVAO(shader.NewProgram("triangle"), &geo.Triangle, 2)

	vao.Draw()
	game.Win.SwapBuffers()

	for !game.Win.ShouldClose() {

		input.PollKeyboard(game)
		game.Win.SwapBuffers()
		glfw.PollEvents()

	}
}
