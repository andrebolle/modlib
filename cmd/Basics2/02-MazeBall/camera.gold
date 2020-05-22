package main

import (
	"github.com/go-gl/mathgl/mgl32"
)

// Frustum Frustum
type Frustum struct {
	Aspect float32 // width/height
	Fovy   float32
	Near   float32
	Far    float32
}

// Camera Camera
type Camera struct {
	Position  mgl32.Vec3
	Direction mgl32.Vec3
	Right     mgl32.Vec3
	Up        mgl32.Vec3
	Frustum
}

var (
	zero = mgl32.Vec3{0, 0, 0}
)

// NewCamera NewCamera
func NewCamera(pos, lookat, up mgl32.Vec3, aspect, fovy, near, far float32) Camera {

	camera := Camera{
		// Position of camera in world space
		Position: mgl32.Vec3{0, 0, 4},
		// Direction of camera in world space
		Direction: mgl32.Vec3{0, 0, -1},
		// Represents the positive x-axis of the camera space.
		Right:   mgl32.Vec3{1, 0, 0},
		Up:      mgl32.Vec3{0, 1, 0},
		Frustum: Frustum{aspect, fovy, near, far},
	}

	return camera

}
