package main

import (
	"math/rand"
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
	program, _ := utils.CreateVF(utils.ReadShader("colourMVP.vs.glsl"), utils.ReadShader("colourMVP.fs.glsl"))
	defer gl.DeleteProgram(program)
	gl.UseProgram(program)

	// Get Locations
	modelLocation := gl.GetUniformLocation(program, gl.Str("model\x00"))
	viewLocation := gl.GetUniformLocation(program, gl.Str("view\x00"))
	projectionLocation := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	positionLocation := uint32(gl.GetAttribLocation(program, gl.Str("position\x00")))
	colourLocation := uint32(gl.GetAttribLocation(program, gl.Str("colour\x00")))

	// Bind Locations
	model := mgl32.Ident4()
	gl.UniformMatrix4fv(modelLocation, 1, false, &model[0])
	projection := mgl32.Perspective(cam.Fovy, cam.Aspect, cam.Near, cam.Far)
	gl.UniformMatrix4fv(projectionLocation, 1, false, &projection[0])

	// The Vertices
	cubes := make([]float32, 0)
	colours := make([]float32, 0)

	for z := 0; z < 10; z++ {
		for y := 0; y < 10; y++ {
			for x := 0; x < 10; x++ {
				// Copy the cube (transforming the vertices)
				for i := 0; i < len(utils.Cube); i += 3 {
					// Shift to the right
					var x1, y1, z1 float32
					x1 = utils.Cube[i] + 14*float32(x)
					y1 = utils.Cube[i+1] + 14*float32(y)
					z1 = utils.Cube[i+2] - 14*float32(z)
					cubes = append(cubes, x1, y1, z1)

				}

				r := rand.Float32()
				g := rand.Float32()
				b := rand.Float32()

				// Different Coloured Cube
				for i := 0; i < len(utils.Cube); i += 3 { // Copy cube (2)

					colours = append(colours, r, g, b)
				}
			}
		}
	}

	// VAO - Gen and Bind
	var array uint32
	gl.GenVertexArrays(1, &array)
	gl.BindVertexArray(array)

	// Buffers: Gen, Bind, VertexAttribPointer, EnableVertexAttribArray, BufferData
	var vertexBuffer uint32
	gl.GenBuffers(1, &vertexBuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)
	gl.VertexAttribPointer(positionLocation, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(positionLocation)
	gl.BufferData(gl.ARRAY_BUFFER, len(cubes)*4, unsafe.Pointer(&cubes[0]), gl.STATIC_DRAW)

	// Create a Buffer for the colours
	var colourBuffer uint32
	gl.GenBuffers(1, &colourBuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, colourBuffer)
	gl.VertexAttribPointer(colourLocation, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(colourLocation)
	//gl.BufferData(gl.ARRAY_BUFFER, len(utils.CubeColour)*4, unsafe.Pointer(&utils.CubeColour[0]), gl.STATIC_DRAW)
	gl.BufferData(gl.ARRAY_BUFFER, len(colours)*4, unsafe.Pointer(&colours[0]), gl.STATIC_DRAW)

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
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(cubes)))
		win.SwapBuffers()
		glfw.PollEvents()
	}
}
