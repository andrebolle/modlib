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

func gen(l L, count uint) (oldSeed string) {
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

const coordMax int = 25000

//Lsystem Lsystem
func Lsystem(l L, gens uint) (*[coordMax]float32, int) {

	// Coord array
	var floats [coordMax]float32

	// Generate snowflake string
	snowflake := gen(l, gens)
	// fmt.Println(snowflake)
	// fmt.Println(len(snowflake))

	// Create a Turtle
	t := Turtle{-0.5, 0.5, 0, 0.1, l.Angle}

	tally := 0
	floats[0] = float32(t.x)
	floats[1] = float32(t.y)
	tally = 2
	for i := 0; i < len(snowflake); i++ {
		command := snowflake[i]
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
