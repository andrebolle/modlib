package utils

import (
	"fmt"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

//var mouseDx, mouseDy float64
var xOld, yOld float64

// YawPitchCamera YawPitchCamera
func YawPitchCamera(t *Camera, yaw, pitch float64) {

	// Yaw/Turn
	up := t.Forward.Cross(t.Right)
	quatRotate := mgl32.QuatRotate(float32(yaw), up)
	t.Forward = quatRotate.Rotate(t.Forward)
	t.Right = quatRotate.Rotate(t.Right)
	t.Forward = t.Forward.Normalize()
	t.Right = t.Right.Normalize()

	// Pitch
	quatRotate = mgl32.QuatRotate(float32(-pitch), t.Right)
	t.Forward = quatRotate.Rotate(t.Forward)
	t.Forward = t.Forward.Normalize()

}

//MoveCamera - Basic WASD with EC (Up/Down) and Space to Pause
func MoveCamera(win *glfw.Window, cam *Camera, action glfw.Action, key glfw.Key, scancode int) {
	dt := float32(0.1)

	if key == glfw.KeyEscape && action == glfw.Press {
		win.SetShouldClose(true)
	}

	state := win.GetKey(glfw.KeyW)
	if state == glfw.Press {
		cam.Position = cam.Position.Add(cam.Forward.Mul(dt))
	}

	state = win.GetKey(glfw.KeyA)
	if state == glfw.Press {
		cam.Position = cam.Position.Add(cam.Right.Mul(-dt))
	}

	state = win.GetKey(glfw.KeyS)
	if state == glfw.Press {
		cam.Position = cam.Position.Add(cam.Forward.Mul(-dt))
	}

	state = win.GetKey(glfw.KeyD)
	if state == glfw.Press {
		cam.Position = cam.Position.Add(cam.Right.Mul(dt))
	}

	state = win.GetKey(glfw.KeyE)
	if state == glfw.Press {
		cam.Position = cam.Position.Add(mgl32.Vec3{0, dt, 0})
	}

	state = win.GetKey(glfw.KeyC)
	if state == glfw.Press {
		cam.Position = cam.Position.Add(mgl32.Vec3{0, -dt, 0})
	}

	state = win.GetKey(glfw.KeySpace)
	if state == glfw.Press {
		fmt.Println("Toggle Pause")
		cam.Paused = !cam.Paused
	}

	// // Check for Key Presses and repeats
	// if action == glfw.Press || action == glfw.Repeat {
	// 	switch key {
	// 	// case glfw.KeyW:
	// 	// 	cam.Position = cam.Position.Add(cam.Forward.Mul(dt))
	// 	// case glfw.KeyS:
	// 	// 	cam.Position = cam.Position.Add(cam.Forward.Mul(-dt))
	// 	// case glfw.KeyA:
	// 	// 	cam.Position = cam.Position.Add(cam.Right.Mul(-dt))
	// 	// case glfw.KeyD:
	// 	// 	cam.Position = cam.Position.Add(cam.Right.Mul(dt))
	// 	// case glfw.KeyE:
	// 	// 	cam.Position = cam.Position.Add(mgl32.Vec3{0, dt, 0})
	// 	// case glfw.KeyC:
	// 	// 	cam.Position = cam.Position.Add(mgl32.Vec3{0, -dt, 0})
	// 	case glfw.KeySpace:
	// 		cam.Paused = !cam.Paused
	// 	}
	// }
}

// SetPitchYawCallback SetPitchYawCallback
func SetPitchYawCallback(win *glfw.Window, cam *Camera) {

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

// SetWASDCallback SetWASDCallback
func SetWASDCallback(win *glfw.Window, cam *Camera) {
	// Keyboard Setup
	keyCallback := func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		MoveCamera(win, cam, action, key, scancode)
	}
	win.SetKeyCallback(keyCallback)
}
