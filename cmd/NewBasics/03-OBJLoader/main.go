// Copyright 2014 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Renders a textured spinning cube using GLFW 3 and OpenGL 4.1 core forward-compatible profile.
package main

import (
	_ "image/png"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/purelazy/modlib/internal/utils"
)

func main() {

	// Window, Camera
	window, cam := utils.GetWindowAndCamera(800, 600)
	defer window.Destroy()

	// Object Data
	floats, indices, stride, posOffset, texOffset, _ := OJBLoader("cube.obj")

	// Program, Vertex, Buffer, Index and Texture objects
	program, _ := utils.NewProgram(utils.ReadShader("textureMVP.vs.glsl"), utils.ReadShader("textureMVP.fs.glsl"))
	defer gl.DeleteProgram(program)
	gl.UseProgram(program)

	vao := new(utils.VAO)
	utils.NewArray(vao)

	vbo := new(utils.VBO)
	utils.NewBuffer(vbo, floats)

	ibo := new(utils.IBO)
	utils.NewIndices(ibo, indices)

	tex := new(utils.Texture)
	utils.NewTexture(tex, "square.png")
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, tex.Width, tex.Height, 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(tex.Texture))

	// Attributes
	vao.Attribute(program, "vert", 3, gl.FLOAT, false, int32(stride), gl.PtrOffset(posOffset))
	vao.Attribute(program, "vertTexCoord", 2, gl.FLOAT, false, int32(stride), gl.PtrOffset(texOffset))

	// Uniforms
	// Model
	model := mgl32.Ident4()
	modelLocation := utils.SetUniformMat4(program, "model", &model[0])
	// View
	viewLocation := gl.GetUniformLocation(program, gl.Str("view\x00"))
	// Projection
	projection := mgl32.Perspective(cam.Fovy, cam.Aspect, cam.Near, cam.Far)
	utils.SetUniformMat4(program, "projection", &projection[0])
	// Sampler/Texture Unit
	textureLocation := gl.GetUniformLocation(program, gl.Str("tex\x00"))
	gl.Uniform1i(textureLocation, 0)

	// Pre Draw Setup
	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.CULL_FACE)
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

		gl.UniformMatrix4fv(viewLocation, 1, false, &view[0])
		gl.UniformMatrix4fv(modelLocation, 1, false, &model[0])

		// Clear, Draw, Swap, Poll
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.DrawElements(gl.TRIANGLES, int32(len(*indices)), gl.UNSIGNED_INT, gl.PtrOffset(0))
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
