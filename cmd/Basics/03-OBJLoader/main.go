// Copyright 2014 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Renders a textured spinning cube using GLFW 3 and OpenGL 4.1 core forward-compatible profile.
package main

import (
	"fmt"
	_ "image/png"
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/purelazy/modlib/internal/utils"
)

func main() {

	// Window, Camera
	window, cam := utils.GetWindowAndCamera(800, 600)
	defer window.Destroy()

	// Let's see the number of textures that can be accessed by the fragment shader.
	var textureUnits int32

	gl.GetIntegerv(gl.MAX_TEXTURE_IMAGE_UNITS, &textureUnits)
	fmt.Println("Texture units for Fragment Shader", textureUnits)
	gl.GetIntegerv(gl.MAX_COMBINED_TEXTURE_IMAGE_UNITS, &textureUnits)
	fmt.Println("MAX_COMBINED Texture units", textureUnits)

	// Program
	vs := utils.ReadShader("textureMVP.vs.glsl")
	fs := utils.ReadShader("textureMVP.fs.glsl")
	program, _ := utils.CreateVF(vs, fs)
	defer gl.DeleteProgram(program)
	gl.UseProgram(program)

	// Get uniform locations using GetUniformLocation
	// Bind using UniformMatrix4fv, Uniform1i
	// Model
	model := mgl32.Ident4()
	modelUniform := gl.GetUniformLocation(program, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	// View gets updated in the main loop
	cameraUniform := gl.GetUniformLocation(program, gl.Str("view\x00"))

	// Perspective
	projection := mgl32.Perspective(cam.Fovy, cam.Aspect, cam.Near, cam.Far)
	projectionUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	// Texture
	textureUniform := gl.GetUniformLocation(program, gl.Str("tex\x00"))
	// The value of a sampler variable is a reference to a texture unit.
	const textureUnitObject int32 = 0
	gl.Uniform1i(textureUniform, textureUnitObject)

	// LoadTexture
	_, err := utils.LoadTexture("square.png")
	if err != nil {
		panic(err)
	}

	// The Vertices
	floats, indices, stride, posOffset, texOffset, _ := OJBLoader("suzanne.obj")

	// Vertex Array Object
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	// Buffers: Gen, Bind, VertexAttribPointer, EnableVertexAttribArray, BufferData
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(*floats)*4, gl.Ptr(&(*floats)[0]), gl.STATIC_DRAW)

	// Generate a buffer for the indices as well
	var ebo uint32
	gl.GenBuffers(1, &ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(*indices)*4, unsafe.Pointer(&(*indices)[0]), gl.STATIC_DRAW)

	vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, int32(stride), gl.PtrOffset(posOffset))

	texCoordAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, int32(stride), gl.PtrOffset(texOffset))

	// Pre Draw Setup
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0, 0, 0, 1.0)

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
		//gl.DrawArrays(gl.TRIANGLES, 0, 6*2*3)
		gl.DrawElements(gl.TRIANGLES, 2904, gl.UNSIGNED_INT, gl.PtrOffset(0))
		//gl.DrawElements(gl.TRIANGLES, 3, gl.UNSIGNED_INT, unsafe.Pointer())
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
