// github.com/purelazy/modlib
// github.com/purelazy/modlib/internal/utils

// Renders a textured spinning cube using GLFW 3 and OpenGL 4.1 core forward-compatible profile.
package main

import (
	"fmt"
	_ "image/png"
	"log"
	"runtime"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/purelazy/modlib/internal/utils"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func main() {

	// Get the camera
	cam := utils.CreateCam()

	// Set the width of the window in pixels
	var width float32 = 600

	// Find the height required so there is no distortion
	height := int(width / cam.Aspect)

	// Create the window
	window := utils.CreateWindow("WASD Cube", int(width), int(height))
	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	// Compile shader
	program, _ := utils.CreateVF(VertexShader, FragmentShader)
	gl.UseProgram(program)

	// Bind Projection
	projection := mgl32.Perspective(cam.Fovy, cam.Aspect, cam.Near, cam.Far)
	projectionUniform := gl.GetUniformLocation(program, gl.Str("P\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	// Bind View
	view := mgl32.LookAtV(cam.Position, cam.LookAt, cam.Up)
	cameraUniform := gl.GetUniformLocation(program, gl.Str("V\x00"))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &view[0])

	// Bind Model - its set in the render loop anyway
	var model mgl32.Mat4
	modelUniform := gl.GetUniformLocation(program, gl.Str("M\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	// Bind Texture
	textureUniform := gl.GetUniformLocation(program, gl.Str("tex\x00"))
	gl.Uniform1i(textureUniform, 0)

	// Load Texture
	_, err := utils.LoadTexture("square.png")
	if err != nil {
		log.Fatalln(err)
	}

	// Get a Vertex Array
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	// ARRAY_BUFFER[vbo] = size in total, pointer to data, usage (DRAW)
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(utils.CubeVertices)*4, gl.Ptr(utils.CubeVertices), gl.STATIC_DRAW)

	// Located where in the shader?
	vertexLocation := uint32(gl.GetAttribLocation(program, gl.Str("vert\x00")))

	// This vertex is: 3 floats, not normalised, stride, offset by 0
	// X, Y, Z, U, V: stride = 20 (5 values * 4 bytes per value)
	stride := int32(20)
	gl.VertexAttribPointer(vertexLocation, 3, gl.FLOAT, false, stride, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(vertexLocation)

	// This texture is: 2 floats, not normalised, stride, offset by 12
	// X, Y, Z, U, V: stride = 20 (5 values * 4 bytes per value)
	textureLocation := uint32(gl.GetAttribLocation(program, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(textureLocation)
	gl.VertexAttribPointer(textureLocation, 2, gl.FLOAT, false, stride, gl.PtrOffset(12))

	// Use depth test (or not if you want it to look weird)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	// Set clear colour
	grey := float32(0.3)
	gl.ClearColor(grey, grey, grey, 1.0)

	angle := 0.0
	previousTime := glfw.GetTime()

	for !window.ShouldClose() {
		// Clear the screen and depth buffer
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Calculate dt
		time := glfw.GetTime()
		dt := time - previousTime
		previousTime = time

		// Update rotation angle
		angle += dt

		// Create new rotation matrix
		model = mgl32.HomogRotate3D(float32(angle), mgl32.Vec3{0, 1, 0})

		// Bind rotation matrix to shader program
		gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

		// Draw cube
		gl.DrawArrays(gl.TRIANGLES, 0, 6*2*3)

		// Show the frambebuffer
		window.SwapBuffers()

		// Deal with new events
		glfw.PollEvents()
	}
}
