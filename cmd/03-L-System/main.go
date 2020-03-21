package main

import (
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"time"
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/purelazy/modlib/internal/utils"
)

// L L-System
type L struct {
	seed  string
	rules map[rune]string
}

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

func main() {

	// Wire this calling goroutine to its current operating system thread.
	runtime.LockOSThread()

	// Random Seed
	rand.Seed(time.Now().UTC().UnixNano())

	// Get coordinates
	floatArray, coordCount := createVertices()
	points := coordCount / 2

	// Create a window
	win := utils.CreateWindow("Points", 800, 800)

	// Create and install(Use) a shader
	program, _ := utils.CreateVF(utils.VertexShader, utils.FragmentShader)
	defer gl.DeleteProgram(program)
	gl.UseProgram(program)

	var array uint32
	// array -> vertex array
	gl.GenVertexArrays(1, &array)
	gl.BindVertexArray(array)

	var buffer uint32
	// buffer -> buffer
	gl.GenBuffers(1, &buffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, buffer)

	// copy "size" bytes (all) of vertices to (ARRAY) buffer
	//gl.BufferData(gl.ARRAY_BUFFER, int(unsafe.Sizeof(vertices)), unsafe.Pointer(&vertices[0]), gl.STATIC_DRAW)
	gl.BufferData(gl.ARRAY_BUFFER, coordCount*4, unsafe.Pointer(&floatArray[0]), gl.STATIC_DRAW)

	shaderLocation := uint32(gl.GetAttribLocation(program, gl.Str("position\x00")))
	// size: number of components per generic vertex attribute
	gl.VertexAttribPointer(shaderLocation, 2, gl.FLOAT, false, 0, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(shaderLocation)
	// ----------------------------------------------

	// Clear screen
	gl.ClearColor(0, 0, 0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	// Draw
	gl.PointSize(4)
	gl.DrawArrays(gl.LINE_STRIP, 0, int32(points))

	// Swap
	win.SwapBuffers()

	// Poll for window close
	for !win.ShouldClose() {
		glfw.PollEvents()
	}
}
