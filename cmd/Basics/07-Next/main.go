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

	// We need an OpenGL context, a window and a camera
	window, cam := utils.GetWindowAndCamera(800, 600)
	defer window.Destroy()

	// We need something to draw
	floats, indices, stride, posOffset, texOffset, normOffset := utils.OJBLoader("cube.obj")

	// Load the texture for the model
	texture := utils.NewTexture("square.png")
	// Load cubemap texture
	//cubemapTexture := loadCubemap(utils.Faces)
	cubemapTexture := utils.LoadCubemap(utils.Faces)

	// Program
	lighting := utils.NewProgram(utils.ReadShader("Lighting.vs.glsl"), utils.ReadShader("Lighting.fs.glsl"))
	cubemap := utils.NewProgram(utils.ReadShader("cubemap.vs.glsl"), utils.ReadShader("cubemap.fs.glsl"))
	defer gl.DeleteProgram(lighting)
	defer gl.DeleteProgram(cubemap)

	// -------------------------------------- Model ---------------------------------
	gl.UseProgram(lighting)

	// Retrieve vertex attribute locations
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

	// Compute and set shader constants for "lighting" program
	projection := mgl32.Perspective(cam.Fovy, cam.Aspect, cam.Near, cam.Far)
	lightColor := mgl32.Vec3{1, 1, 1}
	lightPos := mgl32.Vec3{3, 3, 3}

	gl.UniformMatrix4fv(uProjectionLocation, 1, false, &projection[0])
	gl.Uniform1i(uTexLocation, 0)
	gl.Uniform3fv(uLightPosLocation, 1, &lightPos[0])
	gl.Uniform3fv(uLightColourLocation, 1, &lightColor[0])

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

	//  --------------------------------------------------- Skybox
	var skyboxVAO, skyboxVBO uint32
	gl.GenVertexArrays(1, &skyboxVAO)
	gl.GenBuffers(1, &skyboxVBO)
	gl.BindVertexArray(skyboxVAO)
	gl.BindBuffer(gl.ARRAY_BUFFER, skyboxVBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(utils.SkyboxVertices)*4, gl.Ptr(&(utils.SkyboxVertices)[0]), gl.STATIC_DRAW)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))

	gl.UseProgram(cubemap)

	// Retrieve uniform locations
	uViewCubemapLocation := gl.GetUniformLocation(cubemap, gl.Str("uView\x00"))
	uProjectionCubemapLocation := gl.GetUniformLocation(cubemap, gl.Str("uProjection\x00"))
	uTexCubemapLocation := gl.GetUniformLocation(cubemap, gl.Str("uTex\x00"))

	// Set uniform values which do not change
	gl.UniformMatrix4fv(uProjectionCubemapLocation, 1, false, &projection[0])
	gl.Uniform1i(uTexCubemapLocation, 0)

	// Pre-Draw Setup
	gl.Enable(gl.DEPTH_TEST)
	//gl.Enable(gl.CULL_FACE)
	gl.FrontFace(gl.CCW) // CCW is default
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(1, 0, 0, 1.0)

	angle := 0.0
	previousTime := glfw.GetTime()

	for !window.ShouldClose() {

		// Update the View Transform, because the Camera/Model may have moved
		time := glfw.GetTime()
		elapsed := time - previousTime
		previousTime = time
		angle += elapsed * 0.1

		// Clear
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Calculate dynamic uniforms
		model := mgl32.HomogRotate3D(float32(angle), mgl32.Vec3{0, 1, 0})
		view := mgl32.LookAtV(cam.Position, cam.Position.Add(cam.Forward), cam.Up)

		// ---------------------------------------- Draw the skybox
		// Draw the Skybox first with no depth writing do it will always be
		// drawn at the background of all the other objects.
		gl.DepthMask(false)

		// Load the cubemap program
		gl.UseProgram(cubemap)

		// Remove translation from the view matrix. i.e. the skybox never moves.
		viewWithoutTranslation := view.Mat3().Mat4()
		gl.UniformMatrix4fv(uViewCubemapLocation, 1, false, &viewWithoutTranslation[0])
		//gl.UniformMatrix4fv(uProjectionCubemapLocation, 1, false, &projection[0])
		// skybox cube
		gl.BindVertexArray(skyboxVAO)
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_CUBE_MAP, cubemapTexture)
		gl.DrawArrays(gl.TRIANGLES, 0, 36)
		gl.BindVertexArray(0)

		// ----------------------------------- Draw the model
		// Do write to the depth buffer
		gl.DepthMask(true)
		gl.UseProgram(lighting)

		gl.UniformMatrix4fv(uViewLocation, 1, false, &view[0])
		gl.UniformMatrix4fv(uModelLocation, 1, false, &model[0])
		gl.Uniform3fv(uViewPosLocation, 1, &cam.Position[0])

		gl.BindVertexArray(cubeVAO)
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture.ID)
		gl.DrawElements(gl.TRIANGLES, int32(len(*indices)), gl.UNSIGNED_INT, gl.PtrOffset(0))
		//gl.DrawElements(gl.TRIANGLES, 0, gl.UNSIGNED_INT, gl.PtrOffset(0))
		gl.BindVertexArray(0)

		// Swap and Poll
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
