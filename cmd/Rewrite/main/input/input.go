package input

import (
	"fmt"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/purelazy/modlib/cmd/Rewrite/main/app"
)

// PollKeyboard PollKeyboard
func PollKeyboard(app *app.App) {
	state := app.Win.GetKey(glfw.KeyX)
	if state == glfw.Press {
		fmt.Println("X Key down")
	}
}
