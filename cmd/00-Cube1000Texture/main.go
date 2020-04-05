// Copyright 2014 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Renders a textured spinning cube using GLFW 3 and OpenGL 4.1 core forward-compatible profile.
package main

import (
	_ "image/png"
	"log"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/purelazy/modlib/internal/utils"
)

func main() {

	// Window, Camera
	window, cam := utils.GetWindowAndCamera()
	defer window.Destroy()

	// Program
	program, _ := utils.CreateVF(utils.ReadShader("textureMVP.vs.glsl"), utils.ReadShader("textureMVP.fs.glsl"))
	defer gl.DeleteProgram(program)
	gl.UseProgram(program)

	// Get and Bind Locations
	projection := mgl32.Perspective(cam.Fovy, cam.Aspect, cam.Near, cam.Far)

	projectionUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	camera := mgl32.LookAtV(mgl32.Vec3{3, 3, 3}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	cameraUniform := gl.GetUniformLocation(program, gl.Str("camera\x00"))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

	model := mgl32.Ident4()
	modelUniform := gl.GetUniformLocation(program, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	textureUniform := gl.GetUniformLocation(program, gl.Str("tex\x00"))
	gl.Uniform1i(textureUniform, 0)

	// LoadTexture, ActiveTexture, BindTexture
	texture, err := utils.LoadTexture("square.png")
	if err != nil {
		log.Fatalln(err)
	}
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	// VAO - Gen and Bind
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	// Buffers: Gen, Bind, VertexAttribPointer, EnableVertexAttribArray, BufferData
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(utils.CubeUV)*4, gl.Ptr(utils.CubeUV), gl.STATIC_DRAW)

	vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))

	texCoordAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))

	// Pre Draw Setup
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)

	angle := 0.0
	previousTime := glfw.GetTime()

	for !window.ShouldClose() {

		// Update the View Transform, because the Camera/Model may have moved
		time := glfw.GetTime()
		elapsed := time - previousTime
		previousTime = time
		angle += elapsed
		model = mgl32.HomogRotate3D(float32(angle), mgl32.Vec3{0, 1, 0})
		view := mgl32.LookAtV(cam.Position, cam.Position.Add(cam.Forward), cam.Up)

		gl.UniformMatrix4fv(cameraUniform, 1, false, &view[0])
		gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

		// Clear, Draw, Swap, Poll
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.DrawArrays(gl.TRIANGLES, 0, 6*2*3)
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
