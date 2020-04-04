package main

import (
	"runtime"
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/purelazy/modlib/internal/utils"
)

func main() {

	// Its all about vertices
	var vertices = [...]mgl32.Vec2{
		// Triangle 1
		{-0.90, -0.90},
		{0.85, -0.90},
		{-0.90, 0.85},
		// Triangle 2
		{0.90, -0.85},
		{0.90, 0.90},
		{-0.85, 0.90},
	}

	// Stick to this thread
	runtime.LockOSThread()

	// Create a window
	win := utils.CreateWindow("Triangles", 800, 600)

	// Create and install(Use) a shader
	program, _ := utils.CreateVF(BasicVS, BasicFS)
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
	gl.BufferData(gl.ARRAY_BUFFER, int(unsafe.Sizeof(vertices)), unsafe.Pointer(&vertices), gl.STATIC_DRAW)

	shaderLocation := uint32(gl.GetAttribLocation(program, gl.Str("pos\x00")))
	gl.VertexAttribPointer(shaderLocation, 2, gl.FLOAT, false, 0, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(shaderLocation)
	// ----------------------------------------------

	// Clear screen
	gl.ClearColor(0, 0, 0.2, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	// Draw
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(vertices)))

	// Swap
	win.SwapBuffers()

	// Poll for window close
	for !win.ShouldClose() {
		glfw.PollEvents()
	}
}
