package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/purelazy/modlib/internal/utils"
)

const coordCount int = 1000

func createVertices() (*[coordCount]float32, int) {

	var floats [coordCount]float32

	for i := 0; i < coordCount; i += 2 {
		floats[i] = rand.Float32()*2 - 1
		floats[i+1] = rand.Float32()*2 - 1
	}

	return &floats, coordCount
}

func main() {

	// Its all about vertices
	rand.Seed(time.Now().UTC().UnixNano())
	floats, coordCount := createVertices()
	points := coordCount / 2
	fmt.Println(len(floats))
	fmt.Println("int(unsafe.Sizeof(floats))", int(unsafe.Sizeof(floats)))
	fmt.Println("Point count: ", points)

	// Stick to this thread
	runtime.LockOSThread()

	// Create a window
	win := utils.CreateWindow("Points", 800, 600)

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
	gl.BufferData(gl.ARRAY_BUFFER, coordCount*4, unsafe.Pointer(&floats[0]), gl.STATIC_DRAW)

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
