// Copyright 2014 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Renders a textured spinning cube using GLFW 3 and OpenGL 4.1 core forward-compatible profile.

// VS Code, left hand column: Green lines are new lines (since last commit), blue lines are changed from last commit,
// and red arrows mean deletion since last commit.
package main

import (
	_ "image/png"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/purelazy/modlib/internal/utils"
)

func main() {

	// Create the OpenGL context, window and camera
	//window, cam := utils.GetWindowAndCamera(1680-1, 1050-1)
	window, cam := utils.GetWindowAndCamera(840, 525)
	defer window.Destroy()

	// Set up Box2D world
	world := setupPhysics()

	// Load Textures and Cubemap (aka Skybox)
	modelTexture := utils.NewTexture("square.png")
	gl.BindTexture(gl.TEXTURE_2D, modelTexture.ID)
	cubemapTexture := utils.Cubemap(utils.Faces)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, cubemapTexture)

	// Compile model and cubemap shaders
	lighting := utils.NewProgram(utils.ReadShader("Lighting.vs.glsl"), utils.ReadShader("Lighting.fs.glsl"))
	cubemapShader := utils.NewProgram(utils.ReadShader("cubemap.vs.glsl"), utils.ReadShader("cubemap.fs.glsl"))
	defer gl.DeleteProgram(lighting)
	defer gl.DeleteProgram(cubemapShader)

	// ------------------------- Compute and set static uniforms
	projection := mgl32.Perspective(cam.Fovy, cam.Aspect, cam.Near, cam.Far)

	//cubeVAO, indices, uniLocs := utils.SetupModel("cubewithhole.obj", lighting, &projection[0])
	vao := utils.SetupModel("cubewithhole.obj", lighting, &projection[0])

	skyboxVAO, uViewCubemapLocation := setupSkybox(cubemapShader, &projection[0])

	for !window.ShouldClose() {

		// View is used in multiple programs
		view := mgl32.LookAtV(cam.Position, cam.Position.Add(cam.Forward), cam.Up)

		{ // ----------------Draw the skybox (36 verts)
			gl.UseProgram(cubemapShader)
			// Drawing the skybox first will draw every pixel, so the screen does not
			// need to be cleared and not depth testing
			gl.Disable(gl.DEPTH_TEST)

			// The skybox does not move, relative to the view. So all translation is set to zero
			viewWithoutTranslation := view.Mat3().Mat4()
			gl.UniformMatrix4fv(uViewCubemapLocation, 1, false, &viewWithoutTranslation[0])

			// Arm GPU with VAO and Render
			gl.BindVertexArray(skyboxVAO)
			gl.DrawArrays(gl.TRIANGLES, 0, 36)
		}

		// Step through time
		world.Step(1.0/60.0, 8, 3)

		bodies := world.GetBodyList()

		// ----------------Draw the bodies
		gl.Enable(gl.DEPTH_TEST)
		gl.Clear(gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(lighting)
		gl.Enable(gl.CULL_FACE) // Only front-facing triangles will be drawn

		// Arm GPU with VAO and Render
		gl.BindVertexArray(vao.CubeVAO)
		gl.UniformMatrix4fv(vao.UniLocs["uView"], 1, false, &view[0])
		gl.Uniform3fv(vao.UniLocs["uViewPos"], 1, &cam.Position[0])

		for b := bodies; b != nil; b = b.GetNext() {
			if b.GetUserData() == "box" {
				// Layer 1 of this box
				// Send Box2D info
				gl.Uniform4f(vao.UniLocs["uPosAngle"], float32(b.GetPosition().X), float32(b.GetPosition().Y), 0, float32(b.GetAngle()))
				gl.DrawElements(gl.TRIANGLES, int32(len(*vao.Indices)), gl.UNSIGNED_INT, gl.PtrOffset(0))

				// Layer 2 of this box
				gl.Uniform4f(vao.UniLocs["uPosAngle"], float32(b.GetPosition().Y), float32(b.GetPosition().X), 20, float32(b.GetAngle()))
				gl.DrawElements(gl.TRIANGLES, int32(len(*vao.Indices)), gl.UNSIGNED_INT, gl.PtrOffset(0))
			}
		}

		// Swap and Poll
		window.SwapBuffers()
		glfw.PollEvents()

	}
}
