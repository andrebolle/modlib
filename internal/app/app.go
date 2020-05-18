package app

import (
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/purelazy/modlib/internal/window"
)

// App App
type App struct {
	Win *glfw.Window
}

// NewApp NewApp
func NewApp(width, height int) *App {
	app := App{window.NewWindow(width, height)}
	if err := gl.Init(); err != nil {
		panic(err)
	}
	return &app
}
