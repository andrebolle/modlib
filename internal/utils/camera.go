package utils

import (
	"math"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

//Camera : position, rotation etc
type Camera struct {
	Position mgl32.Vec3
	LookAt   mgl32.Vec3
	Up       mgl32.Vec3
	Aspect   float32 // width/height
	Fovy     float32
	Near     float32
	Far      float32
}

var cam Camera

// Convert WASD keys to new camera position
func moveCamera(forward, horizontal, vertical float32) {
	right := (cam.LookAt).Cross(cam.Up)
	cam.Position = cam.Position.
		Add(cam.LookAt.Mul(forward)).
		Add(right.Mul(horizontal)).
		Add(cam.Up.Mul(vertical))
}

//CreateCam Yo!
func CreateCam() (cam Camera) {

	// // Field of View (along the Y axis)
	// fovy := mgl32.DegToRad(45.0)
	// // The aspect ratio
	// aspectRatio := float32(windowWidth) / float32(windowHeight)
	// // The near and far clipping distances
	// var nearClip float32 = 0.01
	// var farClip float32 = 12
	// // Perspective generates a Perspective Matrix.
	// projection := mgl32.Perspective(fovy, aspectRatio, nearClip, farClip)

	cam = Camera{
		Position: mgl32.Vec3{3, 3, 3},
		LookAt:   mgl32.Vec3{0, 0, 0},
		Up:       mgl32.Vec3{0, 1, 0},
		Aspect:   1,
		Fovy:     math.Pi / 4,
		Near:     0.1,
		Far:      12,
	}
	return
}

// // Used to move the model along the z-axis
// var zoom float32 = 0

func keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	dt := float32(0.01)
	switch key {
	case glfw.KeyW:
		moveCamera(dt, 0.0, 0.0)
	case glfw.KeyS:
		moveCamera(-dt, 0.0, 0.0)
	case glfw.KeyA:
		moveCamera(0, -dt, 0.0)
	case glfw.KeyD:
		moveCamera(0, dt, 0.0)
	}

}
