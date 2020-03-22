package main

import (
	"fmt"
	"math"
)

// L L-System
type L struct {
	seed  string
	rules map[rune]string
}

//KochIsland KochIsland
var KochIsland = L{seed: "F-F-F-F", rules: map[rune]string{'F': "F-F+F+FF-F-F+F", '-': "-", '+': "+"}}

//KochCurve KochCurve
var KochCurve = L{seed: "F-F-F-F", rules: map[rune]string{'F': "F+F--F+F", '-': "-", '+': "+"}}

//Snowflake Koch Snowflake
var Snowflake = L{seed: "F--F--F", rules: map[rune]string{'F': "F+F--F+F", '-': "-", '+': "+"}}

func (l L) gen(count int) (oldSeed string) {
	oldSeed = l.seed
	for i := 0; i < count; i++ {
		newSeed := ""
		for _, n := range oldSeed {
			newSeed += l.rules[n]
		}
		oldSeed = newSeed
	}
	return
}

// Turtle {x,y,heading}
type Turtle struct {
	x, y, heading, d, theta float64
}

const coordMax int = 5000

func createVertices() (*[coordMax]float32, int) {

	// Coord array
	var floats [coordMax]float32

	// Generate snowflake string
	snowflake := Snowflake.gen(4)
	fmt.Println(snowflake)
	fmt.Println(len(snowflake))

	// Create a Turtle
	t := Turtle{-0.5, 0.5, 0, 0.01, math.Pi / 3}

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
		default:
			panic("unrecognized character")
		}
	}

	return &floats, tally
}
