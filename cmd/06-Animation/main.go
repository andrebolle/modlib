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

	// Lock this calling goroutine to its current operating system thread.
	runtime.LockOSThread()

	// Print instructions
	fmt.Println("Use W,A,S,D,R,F to move forward, left, backward, right, up and down.")

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

	var pause bool = false
	// Define the keyboard input callback function
	keyCallback := func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		// dt is the x, y, z move distance per press (or repeated press by holding the key down)
		dt := float32(0.05)

		// Check for Key Presses and repeats
		if action == glfw.Press || action == glfw.Repeat {
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
			case glfw.KeySpace:
				pause = !pause
			}
		}

	}

	// Set Keyboard Callback function
	win.SetKeyCallback(keyCallback)

	// Create and Use the Shader
	program, _ := utils.CreateVF(utils.MVPVertShader, utils.MVPFragShader)
	defer gl.DeleteProgram(program)
	gl.UseProgram(program)

	// Get Uniform and Atrribute Locations
	modelUniform := gl.GetUniformLocation(program, gl.Str("model\x00"))
	viewLocation := gl.GetUniformLocation(program, gl.Str("view\x00"))
	projectionUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	positionLocation := uint32(gl.GetAttribLocation(program, gl.Str("position\x00")))

	// Set Model
	model := mgl32.Ident4()
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	// Set Projection
	projection := mgl32.Perspective(cam.Fovy, cam.Aspect, cam.Near, cam.Far)
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	// Create a Vertex Array
	var array uint32
	gl.GenVertexArrays(1, &array)
	gl.BindVertexArray(array)

	// Create a Buffer
	var buffer uint32
	gl.GenBuffers(1, &buffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, buffer)

	// Describe the Position
	// size: 2 floats per Position
	gl.VertexAttribPointer(positionLocation, 2, gl.FLOAT, false, 0, gl.PtrOffset(0))

	// Causes position to be passed to the shader
	gl.EnableVertexAttribArray(positionLocation)

	// Define your L-System
	// Hilbert https://en.wikipedia.org/wiki/Hilbert_curve#Representation_as_Lindenmayer_system
	var hilbert = utils.L{
		Seed:  "A",
		Angle: math.Pi / 2,
		Rules: map[rune]string{
			'A': "-BF+AFA+FB-",
			'B': "+AF-BFB-FA+",
			'F': "F",
			'-': "-",
			'+': "+"}}

	// Generate the L-system string
	snowflake := utils.GenLString(hilbert, 6)
	// Create a varying angle
	angle := float64(0)

	// Set Clear Colour
	gl.ClearColor(0, 0, 0, 1.0)

	// Set PointSize and/or LineWidth as required
	//gl.PointSize(4)
	gl.LineWidth(2)

	// Depth Test (if required)
	// gl.Enable(gl.DEPTH_TEST)
	// gl.DepthFunc(gl.LESS)

	// Main loop
	for !win.ShouldClose() {

		// Create the Geometry
		floatArray, coordCount := utils.Lsystem(snowflake, angle)
		points := coordCount / 2

		// Copy the Geometry
		gl.BufferData(gl.ARRAY_BUFFER, coordCount*4, unsafe.Pointer(&floatArray[0]), gl.STATIC_DRAW)

		// Update the View Transform, because the Camera may have moved
		view := mgl32.LookAtV(cam.Position, cam.Position.Add(cam.Direction), cam.Up)
		gl.UniformMatrix4fv(viewLocation, 1, false, &view[0])

		// Set Clear gl.COLOR_BUFFER_BIT and gl.DEPTH_BUFFER_BIT as required
		gl.Clear(gl.COLOR_BUFFER_BIT)

		// Draw the Geometry
		gl.DrawArrays(gl.LINE_STRIP, 0, int32(points))

		// Change the Hilbert Angle, ready for the next Frame
		if pause == false {
			angle += 0.0005
		}

		// Swap
		win.SwapBuffers()

		// Poll
		glfw.PollEvents()
	}
}
