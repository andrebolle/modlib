package utils

// https://github.com/patrickhadlaw/cpp-opengl/blob/master/src/Camera.cpp

import (
	"math"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

//Camera : position, rotation etc
type Camera struct {
	Position   mgl32.Vec3
	Forward    mgl32.Vec3
	Right      mgl32.Vec3
	Projection mgl32.Mat4
	View       mgl32.Mat4
	Up         mgl32.Vec3
	Aspect     float32 // width/height
	Fovy       float32
	Near       float32
	Far        float32
	Paused     bool
}

//Cam Cam
func Cam() *Camera {

	return &Camera{
		Position:   mgl32.Vec3{0, 0, 4},
		Forward:    mgl32.Vec3{0, 0, -1},
		Right:      mgl32.Vec3{1, 0, 0},
		Projection: mgl32.Ident4(),
		View:       mgl32.Ident4(),
		Up:         mgl32.Vec3{0, 1, 0},
		Aspect:     1, // width/height
		Fovy:       math.Pi / 4,
		Near:       0.01,
		Far:        1000,
		Paused:     false,
	}
}

//MoveCamera - Basic WASD with EC (Up/Down) and Space to Pause
func MoveCamera(cam *Camera, action glfw.Action, key glfw.Key) {
	dt := float32(0.05)

	// Check for Key Presses and repeats
	if action == glfw.Press || action == glfw.Repeat {
		switch key {
		case glfw.KeyW:
			cam.Position = cam.Position.Add(cam.Forward.Mul(dt))
		case glfw.KeyS:
			cam.Position = cam.Position.Add(cam.Forward.Mul(-dt))
		case glfw.KeyA:
			cam.Position = cam.Position.Add(cam.Right.Mul(-dt))
		case glfw.KeyD:
			cam.Position = cam.Position.Add(cam.Right.Mul(dt))
		case glfw.KeyE:
			cam.Position = cam.Position.Add(mgl32.Vec3{0, dt, 0})
		case glfw.KeyC:
			cam.Position = cam.Position.Add(mgl32.Vec3{0, -dt, 0})
		case glfw.KeySpace:
			cam.Paused = !cam.Paused
		}
	}
	//cam.Position[1] = 0.5
	// if cam.Position.Y() < 0.2 {
	// 	cam.Position[1] = 0.2
	// }

}

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
