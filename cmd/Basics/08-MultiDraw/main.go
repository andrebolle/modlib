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

	cubeVAO, uModelLocation, uViewLocation, uViewPosLocation, indices := setupModel("cubewithhole.obj", lighting, &projection[0])

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
		gl.BindVertexArray(cubeVAO)
		for b := bodies; b != nil; b = b.GetNext() {
			if b.GetUserData() == "box" {

				// Version 1 of this box
				position := b.GetPosition()
				angle := b.GetAngle()
				rotate := mgl32.HomogRotate3D(float32(angle), mgl32.Vec3{0, 0, 1})
				translate := mgl32.Translate3D(float32(position.X), float32(position.Y), 0)
				model := translate.Mul4(rotate)

				// Set uniforms
				gl.UniformMatrix4fv(uViewLocation, 1, false, &view[0])
				gl.UniformMatrix4fv(uModelLocation, 1, false, &model[0])
				gl.Uniform3fv(uViewPosLocation, 1, &cam.Position[0])

				gl.DrawElements(gl.TRIANGLES, int32(len(*indices)), gl.UNSIGNED_INT, gl.PtrOffset(0))

				// Version 2 of this box
				translate = mgl32.Translate3D(float32(position.Y), float32(position.X), 20)
				model = translate.Mul4(rotate)

				gl.UniformMatrix4fv(uViewLocation, 1, false, &view[0])
				gl.UniformMatrix4fv(uModelLocation, 1, false, &model[0])
				gl.Uniform3fv(uViewPosLocation, 1, &cam.Position[0])

				gl.DrawElements(gl.TRIANGLES, int32(len(*indices)), gl.UNSIGNED_INT, gl.PtrOffset(0))

				// Version 3 of this box
				translate = mgl32.Translate3D(float32(-position.Y), float32(position.X), 40)
				model = translate.Mul4(rotate)

				gl.UniformMatrix4fv(uViewLocation, 1, false, &view[0])
				gl.UniformMatrix4fv(uModelLocation, 1, false, &model[0])
				gl.Uniform3fv(uViewPosLocation, 1, &cam.Position[0])

				gl.DrawElements(gl.TRIANGLES, int32(len(*indices)), gl.UNSIGNED_INT, gl.PtrOffset(0))

				// Version 4 of this box
				translate = mgl32.Translate3D(float32(-position.X), float32(position.Y), 60)
				model = translate.Mul4(rotate)

				// Set uniforms
				gl.UniformMatrix4fv(uViewLocation, 1, false, &view[0])
				gl.UniformMatrix4fv(uModelLocation, 1, false, &model[0])
				gl.Uniform3fv(uViewPosLocation, 1, &cam.Position[0])

				gl.DrawElements(gl.TRIANGLES, int32(len(*indices)), gl.UNSIGNED_INT, gl.PtrOffset(0))

			}
		}

		// Swap and Poll
		window.SwapBuffers()
		glfw.PollEvents()

	}
}
