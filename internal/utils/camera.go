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
	Direction  mgl32.Vec3
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
		Direction:  mgl32.Vec3{0, 0, -1},
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
			cam.Position = cam.Position.Add(mgl32.Vec3{0, 0, -dt})
		case glfw.KeyS:
			cam.Position = cam.Position.Add(mgl32.Vec3{0, 0, dt})
		case glfw.KeyA:
			cam.Position = cam.Position.Add(mgl32.Vec3{-dt, 0, 0})
		case glfw.KeyD:
			cam.Position = cam.Position.Add(mgl32.Vec3{dt, 0, 0})
		case glfw.KeyE:
			cam.Position = cam.Position.Add(mgl32.Vec3{0, dt, 0})
		case glfw.KeyC:
			cam.Position = cam.Position.Add(mgl32.Vec3{0, -dt, 0})
		case glfw.KeySpace:
			cam.Paused = !cam.Paused
		}
	}

}
