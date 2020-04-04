package main

import (
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/purelazy/modlib/internal/utils"
)

func main() {

	// ---------------------- Setup

	// A basic Window and Camera
	win, cam := utils.GetWindowAndCamera()
	defer win.Destroy()

	// A basic Shader
	program, _ := utils.CreateVF(utils.ReadShader("mvp.vs.glsl"), utils.ReadShader("mvp.fs.glsl"))
	defer gl.DeleteProgram(program)
	gl.UseProgram(program)

	// A basic Vertex Array and Buffer
	var array, buffer uint32
	gl.GenVertexArrays(1, &array)
	gl.BindVertexArray(array)
	gl.GenBuffers(1, &buffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, buffer)

	// The array of Vertices
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

	// Copy the Vertices to the Buffer
	gl.BufferData(gl.ARRAY_BUFFER, int(unsafe.Sizeof(vertices)), unsafe.Pointer(&vertices), gl.STATIC_DRAW)

	// Get Uniform and Atrribute Locations
	modelLocation := gl.GetUniformLocation(program, gl.Str("model\x00"))
	viewLocation := gl.GetUniformLocation(program, gl.Str("view\x00"))
	projectionLocation := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	positionLocation := uint32(gl.GetAttribLocation(program, gl.Str("pos\x00")))

	// Set VertexAttribPointer & EnableVertexAttribArray
	gl.VertexAttribPointer(positionLocation, 2, gl.FLOAT, false, 0, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(positionLocation)

	// Use UniformMatrix4fv to set MVP as required
	model := mgl32.Ident4()
	gl.UniformMatrix4fv(modelLocation, 1, false, &model[0])
	projection := mgl32.Perspective(cam.Fovy, cam.Aspect, cam.Near, cam.Far)
	gl.UniformMatrix4fv(projectionLocation, 1, false, &projection[0])

	// Pre Draw Setup
	gl.ClearColor(0, 0, 0.2, 1.0)

	// Main Draw Loop
	for !win.ShouldClose() {

		// Update the View Transform, because the Camera may have moved
		view := mgl32.LookAtV(cam.Position, cam.Position.Add(cam.Forward), cam.Up)
		gl.UniformMatrix4fv(viewLocation, 1, false, &view[0])

		// Clear, Draw, Swap, Poll
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(vertices)))
		win.SwapBuffers()
		glfw.PollEvents()
	}
}
