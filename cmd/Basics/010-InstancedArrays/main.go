// Copyright 2014 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Renders a textured spinning cube using GLFW 3 and OpenGL 4.1 core forward-compatible profile.

// VS Code, left hand column: Green lines are new lines (since last commit), blue lines are changed from last commit,
// and red arrows mean deletion since last commit.
package main

import (
	"fmt"
	_ "image/png"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/purelazy/modlib/internal/utils"
)

func main() {

	// Create the OpenGL context, window and camera
	window, cam := utils.GetWindowAndCamera(840, 525)

	var uniformBlockSize int32
	gl.GetIntegerv(gl.MAX_UNIFORM_BLOCK_SIZE, &uniformBlockSize)
	fmt.Println("uniformBlockSize", uniformBlockSize)

	defer window.Destroy()

	// Set up Box2D world
	boxCount := 100
	world := setupPhysics(boxCount)

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

	nutVAO := utils.SetupModel("cube.obj", lighting, &projection[0], world)
	//sphereVAO := utils.SetupModel("sphere.obj", lighting, &projection[0], world)

	skyboxVAO, uViewCubemapLocation := setupSkybox(cubemapShader, &projection[0])

	// --------------------- Main Loop
	for !window.ShouldClose() {

		// View is used in multiple programs
		view := mgl32.LookAtV(cam.Position, cam.Position.Add(cam.Forward), cam.Up)
		
		// ----------------Draw the skybox (36 verts)
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
	


		// ----------------Draw the bodies

		// Step through time
		world.Step(1.0/60.0, 8, 3)

		gl.Enable(gl.DEPTH_TEST)
		gl.Clear(gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(lighting)
		gl.Enable(gl.CULL_FACE) // Only front-facing triangles will be drawn

		// Arm GPU with VAO and Render
		gl.UniformMatrix4fv(nutVAO.UniLocs["uView"], 1, false, &view[0])
		gl.Uniform3fv(nutVAO.UniLocs["uViewPos"], 1, &cam.Position[0])

		//start := time.Now()
		gl.BindVertexArray(nutVAO.Vao)
		gl.BindBuffer(gl.ARRAY_BUFFER, nutVAO.Vbo)
 
		// Get new position/angle data
		posAndAngle := utils.GetPositionAndAngle(world)
		fmt.Println("len(*posAndAngle) ", len(*posAndAngle), "floats")
		fmt.Println("PosAndAngleOffset: ", nutVAO.PosAndAngleOffset)

		// Write data to array
		gl.BufferSubData(gl.ARRAY_BUFFER, nutVAO.PosAndAngleOffset, len(*posAndAngle)*4, gl.Ptr(*posAndAngle))
		gl.DrawElementsInstanced(gl.TRIANGLES, int32(len(*nutVAO.Indices)), gl.UNSIGNED_INT, gl.PtrOffset(0), int32(boxCount))

		//elapsed := time.Since(start)
		//fmt.Println("Inner loop took ", elapsed)
		// Swap and Poll
		window.SwapBuffers()
		glfw.PollEvents()

	}
}
