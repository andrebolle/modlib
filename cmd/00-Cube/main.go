package main

import (
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/purelazy/modlib/internal/utils"
)

func main() {

	win, cam := utils.AppInit()
	defer win.Destroy()

	// Create and Use the Shader
	program, _ := utils.CreateVF(utils.ReadShader("mvp.vs.glsl"), utils.ReadShader("mvp.fs.glsl"))
	defer gl.DeleteProgram(program)
	gl.UseProgram(program)

	// Gen and Bind Vertex Array & Buffer
	var array, buffer uint32
	gl.GenVertexArrays(1, &array)
	gl.BindVertexArray(array)
	gl.GenBuffers(1, &buffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, buffer)

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

	// Copy "size" bytes (all) of static vertex data to (ARRAY) buffer
	gl.BufferData(gl.ARRAY_BUFFER, int(unsafe.Sizeof(vertices)), unsafe.Pointer(&vertices), gl.STATIC_DRAW)

	// Get Uniform and Atrribute Locations
	modelLocation := gl.GetUniformLocation(program, gl.Str("model\x00"))
	viewLocation := gl.GetUniformLocation(program, gl.Str("view\x00"))
	projectionLocation := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	positionLocation := uint32(gl.GetAttribLocation(program, gl.Str("pos\x00")))

	// Do VertexAttribPointer & EnableVertexAttribArray
	gl.VertexAttribPointer(positionLocation, 2, gl.FLOAT, false, 0, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(positionLocation)

	// Set MVP as required
	model := mgl32.Ident4()
	gl.UniformMatrix4fv(modelLocation, 1, false, &model[0])
	projection := mgl32.Perspective(cam.Fovy, cam.Aspect, cam.Near, cam.Far)
	gl.UniformMatrix4fv(projectionLocation, 1, false, &projection[0])

	// Set clear colour
	gl.ClearColor(0, 0, 0.2, 1.0)

	//  By default, face culling is disabled.
	// gl.Enable(gl.CULL_FACE)

	// Main Draw Loop
	for !win.ShouldClose() {

		// Update the View Transform, because the Camera may have moved
		view := mgl32.LookAtV(cam.Position, cam.Position.Add(cam.Forward), cam.Up)
		gl.UniformMatrix4fv(viewLocation, 1, false, &view[0])

		// Clear
		gl.Clear(gl.COLOR_BUFFER_BIT)

		// Draw
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(vertices)))

		// Swap
		win.SwapBuffers()

		// Poll
		glfw.PollEvents()
	}
}
