// Copyright 2014 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Renders a textured spinning cube using GLFW 3 and OpenGL 4.1 core forward-compatible profile.
package main

import (
	_ "image/png"
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/purelazy/modlib/internal/utils"
)

func main() {

	// The Vertices
	floats, indices, stride, posOffset, texOffset, normOffset := utils.OJBLoader("suzanne.obj")

	// Window, Camera
	window, cam := utils.GetWindowAndCamera(800, 600)
	defer window.Destroy()

	// Program
	lighting := utils.NewProgram(utils.ReadShader("Lighting.vs.glsl"), utils.ReadShader("Lighting.fs.glsl"))
	defer gl.DeleteProgram(lighting)
	gl.UseProgram(lighting)

	// Vertex Attribute locations
	aPosLocation := uint32(gl.GetAttribLocation(lighting, gl.Str("aPos\x00")))
	aUVLocation := uint32(gl.GetAttribLocation(lighting, gl.Str("aUV\x00")))
	aNormalLocation := uint32(gl.GetAttribLocation(lighting, gl.Str("aNormal\x00")))

	// Retrieve uniform locations
	uModelLocation := gl.GetUniformLocation(lighting, gl.Str("uModel\x00"))
	uViewLocation := gl.GetUniformLocation(lighting, gl.Str("uView\x00"))
	uProjectionLocation := gl.GetUniformLocation(lighting, gl.Str("uProjection\x00"))
	uTexLocation := gl.GetUniformLocation(lighting, gl.Str("uTex\x00"))
	uViewPosLocation := gl.GetUniformLocation(lighting, gl.Str("uViewPos\x00"))
	uLightColourLocation := gl.GetUniformLocation(lighting, gl.Str("uLightColor\x00"))
	uLightPosLocation := gl.GetUniformLocation(lighting, gl.Str("uLightPos\x00"))

	// Compute static uniform values
	projection := mgl32.Perspective(cam.Fovy, cam.Aspect, cam.Near, cam.Far)
	lightColor := mgl32.Vec3{1, 1, 1}
	lightPos := mgl32.Vec3{3, 3, 3}

	// Set static uniform values
	gl.UniformMatrix4fv(uProjectionLocation, 1, false, &projection[0])
	gl.Uniform1i(uTexLocation, 0)
	gl.Uniform3fv(uLightPosLocation, 1, &lightPos[0])
	gl.Uniform3fv(uLightColourLocation, 1, &lightColor[0])

	// LoadTexture
	texture := utils.NewTexture("square.png")
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB, texture.Width,
		texture.Height, 0, gl.RGBA, gl.UNSIGNED_BYTE, unsafe.Pointer(&texture.RGBA.Pix[0]))

	// Vertex Array Object
	var modelVAO uint32
	gl.GenVertexArrays(1, &modelVAO)
	gl.BindVertexArray(modelVAO)

	// OpenGL objects
	var cubeVAO, cubeVBO, cubeEBO uint32
	// Gens
	gl.GenVertexArrays(1, &cubeVAO)
	gl.GenBuffers(1, &cubeVBO)
	gl.GenBuffers(1, &cubeEBO)

	// Binds
	gl.BindVertexArray(cubeVAO)
	gl.BindBuffer(gl.ARRAY_BUFFER, cubeVBO)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, cubeEBO)

	// Fill buffers
	gl.BufferData(gl.ARRAY_BUFFER, len(*floats)*4, gl.Ptr(&(*floats)[0]), gl.STATIC_DRAW)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(*indices)*4, unsafe.Pointer(&(*indices)[0]), gl.STATIC_DRAW)

	// Describe and Enable vertex attributes
	gl.EnableVertexAttribArray(aPosLocation)
	gl.VertexAttribPointer(aPosLocation, 3, gl.FLOAT, false, int32(stride), gl.PtrOffset(posOffset))

	gl.EnableVertexAttribArray(aUVLocation)
	gl.VertexAttribPointer(aUVLocation, 2, gl.FLOAT, false, int32(stride), gl.PtrOffset(texOffset))

	gl.EnableVertexAttribArray(aNormalLocation)
	gl.VertexAttribPointer(aNormalLocation, 3, gl.FLOAT, false, int32(stride), gl.PtrOffset(normOffset))

	// Pre-Draw Setup
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
		angle += elapsed * 0.1
		model := mgl32.HomogRotate3D(float32(angle), mgl32.Vec3{0, 1, 0})
		view := mgl32.LookAtV(cam.Position, cam.Position.Add(cam.Forward), cam.Up)

		// Set dynamic uniforms
		gl.UniformMatrix4fv(uViewLocation, 1, false, &view[0])
		gl.UniformMatrix4fv(uModelLocation, 1, false, &model[0])
		gl.Uniform3fv(uViewPosLocation, 1, &cam.Position[0])

		// Clear, Draw, Swap, Poll
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.DrawElements(gl.TRIANGLES, int32(len(*indices)), gl.UNSIGNED_INT, gl.PtrOffset(0))
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
