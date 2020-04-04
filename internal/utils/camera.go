package utils

// https://github.com/patrickhadlaw/cpp-opengl/blob/master/src/Camera.cpp

import (
	"math"

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



