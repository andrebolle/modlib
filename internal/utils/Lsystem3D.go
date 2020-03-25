package utils

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl32"
)

// L3D L-System
type L3D struct {
	Seed  string
	Angle float64
	Rules map[rune]string
}

//GenLString3D GenLString
func GenLString3D(l L3D, count uint) (oldSeed string) {
	oldSeed = l.Seed
	for i := uint(0); i < count; i++ {
		newSeed := ""
		for _, n := range oldSeed {
			newSeed += l.Rules[n]
		}
		oldSeed = newSeed
	}
	return
}

// Turtle3D {x,y,heading}
type Turtle3D struct {
	position          mgl32.Vec3
	direction         mgl32.Vec3
	right             mgl32.Vec3
	heading, d, theta float64
}

const coordMax3D int = 50000

//Lsystem3D Lsystem
func Lsystem3D(lString string, angle float64) (*[coordMax3D]float32, int) {
	const coordinates = 3

	// Coord array
	var floats [coordMax3D]float32

	// Create a Turtle
	t := Turtle3D{
		position:  mgl32.Vec3{0, 0, 0},
		direction: mgl32.Vec3{0, 1, 0},
		right:     mgl32.Vec3{1, 0, 0},
		heading:   0,
		d:         0.1,
		theta:     angle}

	tally := 0
	floats[0] = float32(t.position.X())
	floats[1] = float32(t.position.Y())
	floats[2] = float32(t.position.Z())
	tally = coordinates
	for i := 0; i < len(lString); i++ {
		command := lString[i]
		// fmt.Printf("Do: %c\n", command)
		switch command {
		// case 'F':
		// 	x1 := t.position.X() + float32(t.d*math.Cos(t.heading))
		// 	y1 := t.position.Y() + float32(t.d*math.Sin(t.heading))
		// 	z1 := t.position.Z()
		// 	floats[tally] = float32(x1)
		// 	floats[tally+1] = float32(y1)
		// 	floats[tally+2] = float32(z1)
		// 	t.position[0], t.position[1], t.position[2] = x1, y1, z1
		// 	tally += coordinates
		case 'F':
			newPos := t.position.Add(t.direction.Mul(float32(t.d)))
			floats[tally] = float32(newPos.X())
			floats[tally+1] = float32(newPos.Y())
			floats[tally+2] = float32(newPos.Z())
			t.position = newPos
			tally += coordinates

		// case '+':
		// 	t.heading += t.theta
		case '+':
			axis := t.direction
			axis = axis.Cross(t.right)
			quatRotate := mgl32.QuatRotate(float32(angle), axis)
			t.direction = quatRotate.Rotate(t.direction)
			t.right = quatRotate.Rotate(t.right)
			t.direction = t.direction.Normalize()
			t.right = t.right.Normalize()

		// case '-':
		// 	t.heading -= t.theta
		case '-':
			axis := t.direction
			axis = axis.Cross(t.right)
			quatRotate := mgl32.QuatRotate(float32(-angle), axis)
			t.direction = quatRotate.Rotate(t.direction)
			t.right = quatRotate.Rotate(t.right)
			t.direction = t.direction.Normalize()
			t.right = t.right.Normalize()

		case 'A':
		case 'B':
		default:
			panic("unrecognized character")
		}
		fmt.Println("Pos: ", t.position)
		fmt.Println("Dir: ", t.direction)
	}

	return &floats, tally
}
