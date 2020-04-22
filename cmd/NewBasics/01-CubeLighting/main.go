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
	//floats, indices, stride, posOffset, texOffset, _ := OJBLoader("cube.obj")
	floats, indices, stride, posOffset, texOffset, normOffset := OJBLoader("suzanne.obj")

	// Program, Vertex, Buffer, Index and Texture objects
	program, _ := utils.NewProgram(utils.ReadShader("NormSpec.vs.glsl"), utils.ReadShader("NormSpec.fs.glsl"))
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
	vao.Attribute(program, "aPos", 3, gl.FLOAT, false, int32(stride), gl.PtrOffset(posOffset))
	vao.Attribute(program, "aUV", 2, gl.FLOAT, false, int32(stride), gl.PtrOffset(texOffset))
	vao.Attribute(program, "aNormal", 3, gl.FLOAT, false, int32(stride), gl.PtrOffset(normOffset))

	// normalAttrib := uint32(gl.GetAttribLocation(program, gl.Str("aNormal\x00")))
	// gl.EnableVertexAttribArray(normalAttrib)
	// gl.VertexAttribPointer(normalAttrib, 3, gl.FLOAT, false, int32(stride), gl.PtrOffset(normOffset))

	// Uniforms
	// Projection
	projection := mgl32.Perspective(cam.Fovy, cam.Aspect, cam.Near, cam.Far)
	utils.SetUniformMat4(program, "uProjection", &projection[0])
	// Sampler/Texture Unit
	textureLocation := gl.GetUniformLocation(program, gl.Str("uTex\x00"))
	gl.Uniform1i(textureLocation, 0)

	lightColor := mgl32.Vec3{1, 1, 1}
	lightColourUniform := gl.GetUniformLocation(program, gl.Str("uLightColor\x00"))
	gl.Uniform3fv(lightColourUniform, 1, &lightColor[0])

	lightPos := mgl32.Vec3{3, 3, 3}
	lightPosUniform := gl.GetUniformLocation(program, gl.Str("uLightPos\x00"))
	gl.Uniform3fv(lightPosUniform, 1, &lightPos[0])

	// viewPosUniform := gl.GetUniformLocation(program, gl.Str("uViewPos\x00"))
	// gl.Uniform3fv(viewPosUniform, 1, &cam.Position[0])

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
		angle += elapsed * 0

		// Update view position, and model and view matrices
		utils.SetUniformVec3(program, "uViewPos", &cam.Position[0])

		model := mgl32.HomogRotate3D(float32(angle), mgl32.Vec3{0, 1, 0})
		utils.SetUniformMat4(program, "uModel", &model[0])

		view := mgl32.LookAtV(cam.Position, cam.Position.Add(cam.Forward), cam.Up)
		utils.SetUniformMat4(program, "uView", &view[0])

		// Clear, Draw, Swap, Poll
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.DrawElements(gl.TRIANGLES, int32(len(*indices)), gl.UNSIGNED_INT, gl.PtrOffset(0))
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
