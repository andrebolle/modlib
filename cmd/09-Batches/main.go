package main

import (
	"fmt"
	"os"
	"runtime"
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

	// Print instructions
	fmt.Println("Use W,A,S,D,E,C to move forward, left, backward, right, up and down.")

	// Camera and Sceen choices (width, height, aspect ratio)
	monitor := glfw.GetPrimaryMonitor()
	vidMode := monitor.GetVideoMode()

	// Set the camera aspect to the screen aspect
	cam := utils.Cam()
	cam.Aspect = float32(vidMode.Width) / float32(vidMode.Height)

	// Print some info
	fmt.Println(vidMode.Width, "x ", vidMode.Height)
	fmt.Println("cam.Aspect", cam.Aspect)

	// Create a window
	win := utils.CreateWindow(os.Args[0], vidMode.Width/2, vidMode.Height/2)
	defer win.Destroy()

	// Write the graphics card name
	fmt.Println(gl.GoStr(gl.GetString(gl.RENDERER)))

	// Keyboard Setup
	keyCallback := func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		utils.MoveCamera(cam, action, key)
	}
	win.SetKeyCallback(keyCallback)

	utils.CreateMouse(win, cam)

	// Create and Use the Shader
	program, _ := utils.CreateVF(utils.MVPColourVertShader, utils.MVPColourFragShader)
	defer gl.DeleteProgram(program)
	gl.UseProgram(program)

	// Get Uniform and Atrribute Locations
	modelUniform := gl.GetUniformLocation(program, gl.Str("model\x00"))
	viewLocation := gl.GetUniformLocation(program, gl.Str("view\x00"))
	projectionUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	positionLocation := uint32(gl.GetAttribLocation(program, gl.Str("position\x00")))
	colourLocation := uint32(gl.GetAttribLocation(program, gl.Str("colour\x00")))

	// Set Model
	model := mgl32.Ident4()
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	// Set Projection
	projection := mgl32.Perspective(cam.Fovy, cam.Aspect, cam.Near, cam.Far)
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	// Create a Vertex Array - I will use 2 buffers, on for vertices, one for colours
	var array uint32
	gl.GenVertexArrays(1, &array)
	gl.BindVertexArray(array)

	// Create a Buffer for the vertices
	var vertexBuffer uint32
	gl.GenBuffers(1, &vertexBuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)
	// Describe the Position
	// size: 2 floats per Position
	// pointer: Specifies an offset in bytes of the first component of the first generic vertex attribute
	// in the array in the data store of the buffer currently bound to the GL_ARRAY_BUFFER target.
	gl.VertexAttribPointer(positionLocation, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))
	// Causes position to be passed to the shader
	gl.EnableVertexAttribArray(positionLocation)
	// Copy the Geometry to the Array Buffer
	fmt.Println("len(cube)", len(utils.Cube))
	gl.BufferData(gl.ARRAY_BUFFER, len(utils.Cube)*4, unsafe.Pointer(&utils.Cube[0]), gl.STATIC_DRAW)

	// Create a Buffer for the colours
	var colourBuffer uint32
	gl.GenBuffers(1, &colourBuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, colourBuffer)
	// Describe the Position
	// size: 2 floats per Position
	// pointer: Specifies an offset in bytes of the first component of the first generic vertex attribute
	// in the array in the data store of the buffer currently bound to the GL_ARRAY_BUFFER target.
	gl.VertexAttribPointer(colourLocation, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))
	// Causes position to be passed to the shader
	gl.EnableVertexAttribArray(colourLocation)
	// Copy the Geometry to the Array Buffer
	fmt.Println("len(cubeColour)", len(utils.CubeColour))
	gl.BufferData(gl.ARRAY_BUFFER, len(utils.CubeColour)*4, unsafe.Pointer(&utils.CubeColour[0]), gl.STATIC_DRAW)

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
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(utils.Cube)/3))

		// Swap
		win.SwapBuffers()

		// Poll
		glfw.PollEvents()
	}
}
