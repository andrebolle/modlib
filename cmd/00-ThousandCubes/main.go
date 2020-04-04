package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/purelazy/modlib/internal/utils"
)

func main() {

	// Lock this calling goroutine to its current operating system thread.
	runtime.LockOSThread()

	// Initializes the GLFW library
	if err := glfw.Init(); err != nil {
		panic(fmt.Errorf("I could not initialize glfw: %v", err))
	}

	win, cam := utils.FullScreen()
	defer win.Destroy()

	// Write the graphics card name
	fmt.Println(utils.GraphicsCardName())
	// Print instructions
	fmt.Println("Use W,A,S,D,E,C to move forward, left, backward, right, up and down.")

	utils.SetWASDCallback(win, cam)
	utils.SetPitchYawCallback(win, cam)

	// Create and Use the Shader
	program, _ := utils.CreateVF(utils.MVPColourVertShader, utils.MVPColourFragShader)
	defer gl.DeleteProgram(program)
	gl.UseProgram(program)

	// Get Uniform and Atrribute Locations
	modelLocation := gl.GetUniformLocation(program, gl.Str("model\x00"))
	viewLocation := gl.GetUniformLocation(program, gl.Str("view\x00"))
	projectionLocation := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	positionLocation := uint32(gl.GetAttribLocation(program, gl.Str("position\x00")))
	colourLocation := uint32(gl.GetAttribLocation(program, gl.Str("colour\x00")))

	// Set Model
	model := mgl32.Ident4()
	gl.UniformMatrix4fv(modelLocation, 1, false, &model[0])

	// Set Projection
	projection := mgl32.Perspective(cam.Fovy, cam.Aspect, cam.Near, cam.Far)
	gl.UniformMatrix4fv(projectionLocation, 1, false, &projection[0])

	cubeSlice := make([]float32, 0)
	colourSlice := make([]float32, 0)
	rand.Seed(time.Now().UTC().UnixNano())

	for cubeZ := 0; cubeZ < 10; cubeZ++ {
		for cubeY := 0; cubeY < 10; cubeY++ {
			for cubeX := 0; cubeX < 10; cubeX++ {
				// Copy the cube (transforming the vertices)
				for i := 0; i < len(utils.Cube); i += 3 { // Copy cube (2)
					// Shift to the right
					var x1, y1, z1 float32
					x1 = utils.Cube[i] + 14*float32(cubeX)
					y1 = utils.Cube[i+1] + 14*float32(cubeY)
					z1 = utils.Cube[i+2] - 14*float32(cubeZ)
					cubeSlice = append(cubeSlice, x1, y1, z1)
				}

				r := rand.Float32()
				g := rand.Float32()
				b := rand.Float32()

				// Different Coloured Cube
				for i := 0; i < len(utils.CubeColour); i += 3 { // Copy cube (2)

					colourSlice = append(colourSlice, r, g, b)
				}

			}
		}
	}
	// Create a Vertex Array - I will use 2 buffers, on for vertices, one for colours
	var array uint32
	gl.GenVertexArrays(1, &array)
	gl.BindVertexArray(array)

	// Create a Buffer for the vertices
	var vertexBuffer uint32
	gl.GenBuffers(1, &vertexBuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)
	gl.VertexAttribPointer(positionLocation, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(positionLocation)
	gl.BufferData(gl.ARRAY_BUFFER, len(cubeSlice)*4, unsafe.Pointer(&cubeSlice[0]), gl.STATIC_DRAW)

	// Create a Buffer for the colours
	var colourBuffer uint32
	gl.GenBuffers(1, &colourBuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, colourBuffer)
	gl.VertexAttribPointer(colourLocation, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(colourLocation)
	//gl.BufferData(gl.ARRAY_BUFFER, len(utils.CubeColour)*4, unsafe.Pointer(&utils.CubeColour[0]), gl.STATIC_DRAW)
	gl.BufferData(gl.ARRAY_BUFFER, len(colourSlice)*4, unsafe.Pointer(&colourSlice[0]), gl.STATIC_DRAW)

	// Set Clear Colour
	gl.ClearColor(0, 0, 0, 1.0)

	// Set PointSize and/or LineWidth as required
	//gl.PointSize(4)
	gl.LineWidth(1)

	// Depth Test (if required)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	// // Copy the Geometry to the Array Buffer
	// fmt.Println("len(cube)", len(utils.Cube))
	// gl.BufferData(gl.ARRAY_BUFFER, len(utils.Cube)*4, unsafe.Pointer(&utils.Cube[0]), gl.STATIC_DRAW)

	// The Render Loop
	for !win.ShouldClose() {

		// Update the View Transform, because the Camera may have moved
		view := mgl32.LookAtV(cam.Position, cam.Position.Add(cam.Forward), cam.Up)
		gl.UniformMatrix4fv(viewLocation, 1, false, &view[0])

		// Set Clear gl.COLOR_BUFFER_BIT and gl.DEPTH_BUFFER_BIT as required
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Draw the Geometry
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(cubeSlice)/3))

		// Swap
		win.SwapBuffers()

		// Poll
		glfw.PollEvents()
	}
}
