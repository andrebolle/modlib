package utils

// https://github.com/patrickhadlaw/cpp-opengl/blob/master/src/Camera.cpp

import (
	"math"

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
}

//MoveCamera Convert WASD keys to new camera position
func (cam *Camera) MoveCamera(z, x, y float32) {
	//direction := mgl32.Vec3{0, 0, 1}
	//right := (Cam.LookAt).Cross(Cam.Up)
	// Cam.Position = Cam.Position.
	// 	//Add(Cam.LookAt.Mul(forward)).
	// 	Add(direction.Mul(forward)).
	// 	Add(right.Mul(horizontal)).
	// 	Add(Cam.Up.Mul(vertical))
	// fmt.Println(Cam.Position)
	cam.Position = cam.Position.Add(mgl32.Vec3{x, y, z})
	//fmt.Println(Cam.Position)
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
	}
}
