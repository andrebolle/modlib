package main

import (
	"fmt"
	"math"
	"os"
	"runtime"
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/purelazy/modlib/internal/utils"
)

func main() {

	// Wire this calling goroutine to its current operating system thread.
	runtime.LockOSThread()

	fmt.Println("Use W,A,S,D,R,F to move forward, left, backward, right, up and down.")

	// // -------------------------Create the Vertices ------------------------

	// //Hilbert https://en.wikipedia.org/wiki/Hilbert_curve#Representation_as_Lindenmayer_system
	// var hilbert = utils.L{
	// 	Seed:  "A",
	// 	Angle: math.Pi / 2,
	// 	Rules: map[rune]string{
	// 		'A': "-BF+AFA+FB-",
	// 		'B': "+AF-BFB-FA+",
	// 		'F': "F",
	// 		'-': "-",
	// 		'+': "+"}}

	// // Generate the points
	// floatArray, coordCount := utils.Lsystem(hilbert, 6)
	// fmt.Println("coordCount", coordCount)
	// points := coordCount / 2

	// ------------------------ Camera and Window  ----------------------------------

	// Set the width of the window in pixels
	const width float32 = 1600

	// Get a camera
	cam := utils.Cam()

	// Find the window height to match the camera aspect ratio
	height := int(width / cam.Aspect)
	fmt.Println("Width x Height", width, height)
	fmt.Println("cam.Aspect", cam.Aspect)

	// Create a window
	win := utils.CreateWindow(os.Args[0], int(width), int(height))

	// ------------------------ Keyboard Input ----------------------------------

	// Define the keyboard input callback function
	keyCallback := func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		dt := float32(0.1)
		switch key {
		case glfw.KeyW:
			cam.Position = cam.Position.Add(mgl32.Vec3{0, 0, -dt})
		case glfw.KeyS:
			cam.Position = cam.Position.Add(mgl32.Vec3{0, 0, dt})
		case glfw.KeyA:
			cam.Position = cam.Position.Add(mgl32.Vec3{-dt, 0, 0})
		case glfw.KeyD:
			cam.Position = cam.Position.Add(mgl32.Vec3{dt, 0, 0})
		case glfw.KeyR:
			cam.Position = cam.Position.Add(mgl32.Vec3{0, dt, 0})
		case glfw.KeyF:
			cam.Position = cam.Position.Add(mgl32.Vec3{0, -dt, 0})
		}
	}

	win.SetKeyCallback(keyCallback)

	// ------------------------ Load Shader ----------------------------------

	// Create and install(Use) a shader
	program, _ := utils.CreateVF(utils.MVPVertShader, utils.MVPFragShader)
	defer gl.DeleteProgram(program)
	gl.UseProgram(program)

	// --------------------------- Uniforms ----------------------------------

	// Get Uniform locations
	modelUniform := gl.GetUniformLocation(program, gl.Str("model\x00"))
	viewLocation := gl.GetUniformLocation(program, gl.Str("view\x00"))
	projectionUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))

	// Set Model
	model := mgl32.Ident4()
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	// Update Projection
	projection := mgl32.Perspective(cam.Fovy, cam.Aspect, cam.Near, cam.Far)
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	positionLocation := uint32(gl.GetAttribLocation(program, gl.Str("position\x00")))

	// ------------------------ Create Vertex Array --------------------------

	var array uint32
	// array -> vertex array
	gl.GenVertexArrays(1, &array)
	gl.BindVertexArray(array)

	// ------------------------ Create Buffer --------------------------

	var buffer uint32
	// buffer -> buffer
	gl.GenBuffers(1, &buffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, buffer)

	// ------------------------ Describe data Data --------------------------
	// size: components per vertex attribute
	gl.VertexAttribPointer(positionLocation, 2, gl.FLOAT, false, 0, gl.PtrOffset(0))

	// ------------------------ Enable this attributes
	// Causes position to be passed to the shader
	gl.EnableVertexAttribArray(positionLocation)

	// ------------------------- Pre Draw Setup --------------------------------
	gl.ClearColor(0, 0, 0, 1.0)
	gl.PointSize(2)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	angle := float64(-1)
	// Poll for window close
	for !win.ShouldClose() {

		// -------------------------Create the Vertices ------------------------

		//Hilbert https://en.wikipedia.org/wiki/Hilbert_curve#Representation_as_Lindenmayer_system
		var hilbert = utils.L{
			Seed:  "A",
			Angle: math.Pi/2 + angle,
			Rules: map[rune]string{
				'A': "-BF+AFA+FB-",
				'B': "+AF-BFB-FA+",
				'F': "F",
				'-': "-",
				'+': "+"}}

		// Generate the points
		floatArray, coordCount := utils.Lsystem(hilbert, 6)
		//fmt.Println("coordCount", coordCount)
		points := coordCount / 2

		// ------------------------ Copy Data --------------------------
		gl.BufferData(gl.ARRAY_BUFFER, coordCount*4, unsafe.Pointer(&floatArray[0]), gl.STATIC_DRAW)

		// Update View
		view := mgl32.LookAtV(cam.Position, cam.Position.Add(cam.Direction), cam.Up)

		gl.UniformMatrix4fv(viewLocation, 1, false, &view[0])

		// Clear screen
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Draw
		gl.DrawArrays(gl.LINE_STRIP, 0, int32(points))
		angle += 0.01

		// Swap
		win.SwapBuffers()

		glfw.PollEvents()
	}
}
