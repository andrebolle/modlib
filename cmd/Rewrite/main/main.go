package main

import (
	"runtime"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/purelazy/modlib/internal/app"
	"github.com/purelazy/modlib/internal/vao"

	"github.com/purelazy/modlib/internal/geo"
	"github.com/purelazy/modlib/internal/shader"
)

const (
	windowWidth  = 500
	windowHeight = 500
)

func init() {
	runtime.LockOSThread()
}

func main() {

	game := app.NewApp(windowWidth, windowHeight)
	defer glfw.Terminate()

	vao := vao.NewVAO(shader.NewProgram("triangle"), &geo.Triangle, 2)

	vao.Draw()
	game.Win.SwapBuffers()

	for !game.Win.ShouldClose() {

		// input.PollKeyboard(game)
		// game.Win.SwapBuffers()
		glfw.PollEvents()

	}
}
