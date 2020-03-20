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

	cam := utils.CreateCam()
	window := utils.CreateWindow("WASD Cube", 600, 600)
	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	// Compile shader
	program, _ := utils.CreateVF(VertexShader, FragmentShader)
	gl.UseProgram(program)

	// Projection
	projection := mgl32.Perspective(cam.Fovy, cam.Aspect, cam.Near, cam.Far)
	projectionUniform := gl.GetUniformLocation(program, gl.Str("P\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	// View
	view := mgl32.LookAtV(cam.Position, cam.LookAt, cam.Up)
	cameraUniform := gl.GetUniformLocation(program, gl.Str("V\x00"))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &view[0])

	// Model
	//model := mgl32.Ident4()
	var model mgl32.Mat4
	modelUniform := gl.GetUniformLocation(program, gl.Str("M\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	// Texture
	textureUniform := gl.GetUniformLocation(program, gl.Str("tex\x00"))
	gl.Uniform1i(textureUniform, 0)

	// Load Texture
	_, err := utils.LoadTexture("square.png")
	if err != nil {
		log.Fatalln(err)
	}

	// Configure the vertex data
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(utils.CubeVertices)*4, gl.Ptr(utils.CubeVertices), gl.STATIC_DRAW)

	vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))

	texCoordAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))

	// Use depth test (or not if you want it to look weird)
	//gl.Enable(gl.DEPTH_TEST)
	//gl.DepthFunc(gl.LESS)

	// Set clear colour
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)

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
