package utils

import (
	"fmt"

	"github.com/go-gl/glfw/v3.3/glfw"
)

//var mouseDx, mouseDy float64
var xOld, yOld float64

// CreateMouse CreateMouse
func CreateMouse(win *glfw.Window, cam *Camera) {

	// If you wish to implement mouse motion based camera controls or
	// other input schemes that require unlimited mouse movement, set
	// the cursor mode to GLFW_CURSOR_DISABLED.
	win.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)

	// Get raw input (See https://www.glfw.org/docs/latest/input_guide.html#raw_mouse_motion)
	if glfw.RawMouseMotionSupported() {
		fmt.Println("Using raw mouse motion")
		win.SetInputMode(glfw.RawMouseMotion, glfw.True)
	}

	// 400 and 300 ar magic numbers
	win.SetCursorPos(400, 300)
	xOld, yOld = win.GetCursorPos()
	//fmt.Println(xOld, yOld)

	// The callback function
	mouseCallback := func(w *glfw.Window, xPos float64, yPos float64) {
		//fmt.Println(xPos, yPos)

		mouseDx, mouseDy := xPos-xOld, yPos-yOld
		yaw, pitch := mouseDx/600, mouseDy/600
		//fmt.Println(yaw, pitch)
		xOld, yOld = xPos, yPos
		YawPitchCamera(cam, yaw, pitch)
	}

	win.SetCursorPosCallback(mouseCallback)

}
