package app

import (
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

// App App
type App struct {
	Win *glfw.Window
}

// NewApp NewApp
func NewApp() *App {

	return new(App)
}

// StartOpenGL StartOpenGL
func StartOpenGL() {
	if err := gl.Init(); err != nil {
		panic(err)
	}
}
