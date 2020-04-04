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

//KeyCallback KeyCallback
func keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
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
	fmt.Println(cam.Position)
}

var cam *utils.Camera

func main() {

	// Wire this calling goroutine to its current operating system thread.
	runtime.LockOSThread()

	// Get coordinates
	floatArray, coordCount := createVertices()
	points := coordCount / 2

	// Set the width of the window in pixels
	const width float32 = 600

	cam = utils.Cam()

	// Find the window height to match the camera aspect ratio
	height := int(width / cam.Aspect)
	fmt.Println("Width x Height", width, height)
	fmt.Println("cam.Aspect", cam.Aspect)

	// Create a window
	win := utils.CreateWindow(os.Args[0], int(width), int(height))
	win.SetKeyCallback(keyCallback)

	// Create and install(Use) a shader
	program, _ := utils.CreateVF(utils.MVPVertShader, utils.MVPFragShader)
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
	//gl.BufferData(gl.ARRAY_BUFFER, int(unsafe.Sizeof(vertices)), unsafe.Pointer(&vertices[0]), gl.STATIC_DRAW)
	gl.BufferData(gl.ARRAY_BUFFER, coordCount*4, unsafe.Pointer(&floatArray[0]), gl.STATIC_DRAW)

	shaderLocation := uint32(gl.GetAttribLocation(program, gl.Str("position\x00")))
	// size: number of components per generic vertex attribute
	gl.VertexAttribPointer(shaderLocation, 2, gl.FLOAT, false, 0, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(shaderLocation)
	// ----------------------------------------------

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

	// Just before we loop ...
	gl.ClearColor(0, 0, 0, 1.0)
	gl.PointSize(2)

	// Use depth test (or not if you want it to look weird)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	// Poll for window close
	for !win.ShouldClose() {
		// Update View
		view := mgl32.LookAtV(cam.Position, cam.Position.Add(cam.Forward), cam.Up)

		gl.UniformMatrix4fv(viewLocation, 1, false, &view[0])

		// Clear screen
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Draw
		gl.DrawArrays(gl.LINE_STRIP, 0, int32(points))

		// Swap
		win.SwapBuffers()

		glfw.PollEvents()
	}
}
