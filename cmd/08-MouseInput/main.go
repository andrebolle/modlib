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

	// Mouse Setup
	// Hide the cursor
	win.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)

	// Raw mouse motion
	// ----------------
	//
	// When the cursor is disabled, raw (unscaled and unaccelerated) mouse motion
	// can be enabled if available.
	//
	// Raw mouse motion is closer to the actual motion of the mouse across a surface.
	// It is not affected by the scaling and acceleration applied to the motion of the
	// desktop cursor. That processing is suitable for a cursor while raw motion is
	// better for controlling for example a 3D camera. Because of this, raw mouse motion
	// is only provided when the cursor is disabled.
	if glfw.RawMouseMotionSupported() {
		fmt.Println("Using raw mouse motion")
		win.SetInputMode(glfw.RawMouseMotion, glfw.True)
	}

	// Set initial cursor position
	win.SetCursorPos(float64(vidMode.Width)/4, float64(vidMode.Height)/4)
	xOld, yOld := win.GetCursorPos()
	var mouseDx, mouseDy float64 = 0, 0
	var yaw, pitch float64 = 0, 0
	mouseCallback := func(w *glfw.Window, xPos float64, yPos float64) {

		mouseDx, mouseDy = xPos-xOld, yPos-yOld
		yaw, pitch = mouseDx/600, mouseDy/600
		//fmt.Println(yaw, pitch)
		xOld, yOld = xPos, yPos
		utils.YawPitchCamera(cam, yaw, pitch)
	}

	win.SetCursorPosCallback(mouseCallback)

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
	// pointer: Specifies an offset in bytes of the first component of the first generic vertex attribute
	// in the array in the data store of the buffer currently bound to the GL_ARRAY_BUFFER target.
	gl.VertexAttribPointer(positionLocation, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))

	// Causes position to be passed to the shader
	gl.EnableVertexAttribArray(positionLocation)

	// // Define your L-System
	// // Hilbert https://en.wikipedia.org/wiki/Hilbert_curve#Representation_as_Lindenmayer_system
	// var hilbert = utils.L3D{
	// 	Seed:  "A",
	// 	Angle: math.Pi / 2,
	// 	Rules: map[rune]string{
	// 		'A': "-BF+AFA+FB-",
	// 		'B': "+AF-BFB-FA+",
	// 		'F': "F",
	// 		'-': "-",
	// 		'+': "+"}}

	var hilbert = utils.L3D{
		Seed:  "X",
		Angle: math.Pi / 2,
		Rules: map[rune]string{
			'X': "^<XF^<XFX-F^>>XFX&F+>>XFX-F>X->",
			'F': "F",
			'-': "-",
			'+': "+",
			'^': "^",
			'&': "&",
			'<': "<",
			'>': ">",
		}}

	// Generate the L-system string
	snowflake := utils.GenLString3D(hilbert, 3)
	//fmt.Println(snowflake)

	// Set Clear Colour
	gl.ClearColor(0, 0, 0, 1.0)

	// Set PointSize and/or LineWidth as required
	//gl.PointSize(4)
	gl.LineWidth(1)

	// Depth Test (if required)
	// gl.Enable(gl.DEPTH_TEST)
	// gl.DepthFunc(gl.LESS)

	// Create first frame
	floatArray, coordCount := utils.Lsystem3D(snowflake, hilbert.Angle)
	points := coordCount / 3

	// Copy the Geometry to the Array Buffer
	gl.BufferData(gl.ARRAY_BUFFER, coordCount*4, unsafe.Pointer(&floatArray[0]), gl.STATIC_DRAW)

	// The Render Loop
	for !win.ShouldClose() {

		// Update the View Transform, because the Camera may have moved
		view := mgl32.LookAtV(cam.Position, cam.Position.Add(cam.Forward), cam.Up)
		gl.UniformMatrix4fv(viewLocation, 1, false, &view[0])

		// Set Clear gl.COLOR_BUFFER_BIT and gl.DEPTH_BUFFER_BIT as required
		gl.Clear(gl.COLOR_BUFFER_BIT)

		// Draw the Geometry
		gl.DrawArrays(gl.LINE_STRIP, 0, int32(points))

		// Swap
		win.SwapBuffers()

		// Poll
		glfw.PollEvents()
	}
}
