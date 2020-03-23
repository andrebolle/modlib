package utils

import (
	"math"
)

// L L-System
type L struct {
	Seed  string
	Angle float64
	Rules map[rune]string
}

//GenLString GenLString
func GenLString(l L, count uint) (oldSeed string) {
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

// Turtle {x,y,heading}
type Turtle struct {
	x, y, heading, d, theta float64
}

const coordMax int = 50000

//Lsystem Lsystem
func Lsystem(lString string, angle float64) (*[coordMax]float32, int) {

	// Coord array
	var floats [coordMax]float32

	// Create a Turtle
	t := Turtle{-0.5, 0.5, 0, 0.1, angle}

	tally := 0
	floats[0] = float32(t.x)
	floats[1] = float32(t.y)
	tally = 2
	for i := 0; i < len(lString); i++ {
		command := lString[i]
		// fmt.Printf("Do: %c\n", command)
		switch command {
		case 'F':
			x1 := t.x + t.d*math.Cos(t.heading)
			y1 := t.y + t.d*math.Sin(t.heading)
			floats[tally] = float32(x1)
			floats[tally+1] = float32(y1)
			t.x, t.y = x1, y1
			tally += 2
		case '+':
			t.heading += t.theta
		case '-':
			t.heading -= t.theta
		case 'A':
		case 'B':
		default:
			panic("unrecognized character")
		}
	}

	return &floats, tally
}
