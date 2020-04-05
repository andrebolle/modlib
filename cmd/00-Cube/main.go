package main

import (
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/purelazy/modlib/internal/utils"
)

func main() {

	// Window, Camera
	win, cam := utils.GetWindowAndCamera()
	defer win.Destroy()

	// Program
	program, _ := utils.CreateVF(utils.ReadShader("mvp.vs.glsl"), utils.ReadShader("mvp.fs.glsl"))
	defer gl.DeleteProgram(program)
	gl.UseProgram(program)

	// Get Locations
	modelLocation := gl.GetUniformLocation(program, gl.Str("model\x00"))
	viewLocation := gl.GetUniformLocation(program, gl.Str("view\x00"))
	projectionLocation := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	positionLocation := uint32(gl.GetAttribLocation(program, gl.Str("pos\x00")))

	// Bind Locations
	model := mgl32.Ident4()
	gl.UniformMatrix4fv(modelLocation, 1, false, &model[0])
	projection := mgl32.Perspective(cam.Fovy, cam.Aspect, cam.Near, cam.Far)
	gl.UniformMatrix4fv(projectionLocation, 1, false, &projection[0])

	// The Vertices
	// The array of Vertices
	var cube = []float32{
		// Bottom
		-1.0, -1.0, -1.0,
		1.0, -1.0, -1.0,
		-1.0, -1.0, 1.0,
		1.0, -1.0, -1.0,
		1.0, -1.0, 1.0,
		-1.0, -1.0, 1.0,

		// Top
		-1.0, 1.0, -1.0,
		-1.0, 1.0, 1.0,
		1.0, 1.0, -1.0,
		1.0, 1.0, -1.0,
		-1.0, 1.0, 1.0,
		1.0, 1.0, 1.0,

		// Front
		-1.0, -1.0, 1.0,
		1.0, -1.0, 1.0,
		-1.0, 1.0, 1.0,
		1.0, -1.0, 1.0,
		1.0, 1.0, 1.0,
		-1.0, 1.0, 1.0,

		// Back
		-1.0, -1.0, -1.0,
		-1.0, 1.0, -1.0,
		1.0, -1.0, -1.0,
		1.0, -1.0, -1.0,
		-1.0, 1.0, -1.0,
		1.0, 1.0, -1.0,

		// Left
		-1.0, -1.0, 1.0,
		-1.0, 1.0, -1.0,
		-1.0, -1.0, -1.0,
		-1.0, -1.0, 1.0,
		-1.0, 1.0, 1.0,
		-1.0, 1.0, -1.0,

		// Right
		1.0, -1.0, 1.0,
		1.0, -1.0, -1.0,
		1.0, 1.0, -1.0,
		1.0, -1.0, 1.0,
		1.0, 1.0, -1.0,
		1.0, 1.0, 1.0,
	}

	// VAO - Gen and Bind
	var array uint32
	gl.GenVertexArrays(1, &array)
	gl.BindVertexArray(array)

	// Buffers: Gen, Bind, VertexAttribPointer, EnableVertexAttribArray, BufferData
	var buffer uint32
	gl.GenBuffers(1, &buffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, buffer)
	gl.VertexAttribPointer(positionLocation, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(positionLocation)
	gl.BufferData(gl.ARRAY_BUFFER, len(cube)*4, unsafe.Pointer(&cube[0]), gl.STATIC_DRAW)

	// Pre Draw Setup
	gl.ClearColor(0, 0, 0.2, 1.0)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	// Main Draw Loop
	for !win.ShouldClose() {

		// Update the View Transform, because the Camera may have moved
		view := mgl32.LookAtV(cam.Position, cam.Position.Add(cam.Forward), cam.Up)
		gl.UniformMatrix4fv(viewLocation, 1, false, &view[0])

		// Clear, Draw, Swap, Poll
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(cube)))
		win.SwapBuffers()
		glfw.PollEvents()
	}
}
