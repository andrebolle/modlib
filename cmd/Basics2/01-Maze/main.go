// Copyright 2014 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Renders a textured spinning cube using GLFW 3 and OpenGL 4.1 core forward-compatible profile.

// VS Code, left hand column: Green lines are new lines (since last commit), blue lines are changed from last commit,
// and red arrows mean deletion since last commit.
package main

import (
	_ "image/png"

	"github.com/ByteArena/box2d"
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/purelazy/modlib/internal/utils"
)

// App contains app stuff
type App struct {
	window                *glfw.Window
	cam                   *utils.Camera
	projection            mgl32.Mat4
	mazeWidth, mazeHeight int
	wallCount             int
	world                 *box2d.B2World
	nutVAO                *utils.Vao
}

func main() {

	// w The width (must be odd).
	// h The height (must be odd).
	app := App{mazeWidth: 7, mazeHeight: 7}

	// Set up Maze world
	m := mainMaze(app.mazeWidth, app.mazeHeight)
	app.world, app.wallCount = setupMaze(m)

	// Create the OpenGL context, window and camera
	app.window, app.cam = utils.GetWindowAndCamera(800, 600)
	defer app.window.Destroy()

	app.cam.Position = app.cam.StartPosition

	// Load Textures and Cubemap (aka Skybox)
	modelTexture := utils.NewTexture("square.png")
	gl.BindTexture(gl.TEXTURE_2D, modelTexture)
	cubemapTexture := utils.Cubemap(utils.Faces)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, cubemapTexture)

	// Compile model and cubemap shaders
	lighting := utils.NewProgram(utils.ReadShader("Lighting.vs.glsl"), utils.ReadShader("Lighting.fs.glsl"))
	cubemapShader := utils.NewProgram(utils.ReadShader("cubemap.vs.glsl"), utils.ReadShader("cubemap.fs.glsl"))
	defer gl.DeleteProgram(lighting)
	defer gl.DeleteProgram(cubemapShader)

	// ------------------------- Compute and set static uniforms
	app.projection = mgl32.Perspective(app.cam.Fovy, app.cam.Aspect, app.cam.Near, app.cam.Far)

	// Load Obj file
	app.nutVAO = utils.SetupModel("cubewithhole.obj", lighting, &app.projection[0], app.world)
	//sphereVAO := utils.SetupModel("sphere.obj", lighting, &projection[0], world)

	skyboxVAO, uViewCubemapLocation := setupSkybox(cubemapShader, &app.projection[0])

	// --------------------- Main Loop
	for !app.window.ShouldClose() {

		// View is used in multiple programs
		view := mgl32.LookAtV(app.cam.Position, app.cam.Position.Add(app.cam.Forward), app.cam.Up)

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
		// gl.DrawArrays(mode uint32, first int32, count int32)
		gl.DrawArrays(gl.TRIANGLES, 0, 36)

		// ----------------Draw the Box2D objects
		// Step through time

		//app.world.Step(1.0/60.0, 8, 3)

		gl.Enable(gl.DEPTH_TEST)
		gl.Clear(gl.DEPTH_BUFFER_BIT)
		gl.Enable(gl.CULL_FACE) // Only front-facing triangles will be drawn

		// Load program and set uniforms
		gl.UseProgram(lighting)
		gl.UniformMatrix4fv(app.nutVAO.UniLocs["uView"], 1, false, &view[0])
		gl.Uniform3fv(app.nutVAO.UniLocs["uViewPos"], 1, &app.cam.Position[0])

		// Bind VAO and VBO
		gl.BindVertexArray(app.nutVAO.Vao)
		gl.BindBuffer(gl.ARRAY_BUFFER, app.nutVAO.Vbo)

		// // Extract position and angle from Box2D world
		// posAndAngle := utils.GetPositionAndAngle(app.world, "box")

		// // Load posAndAngle into GPU
		// gl.BufferSubData(gl.ARRAY_BUFFER, app.nutVAO.PosAndAngleOffset, len(*posAndAngle)*4, gl.Ptr(*posAndAngle))

		// Draw boxCount instances
		drawAllIndices := int32(len(*app.nutVAO.Indices))
		gl.DrawElementsInstanced(gl.TRIANGLES, drawAllIndices, gl.UNSIGNED_INT, gl.PtrOffset(0), int32(app.wallCount))

		// Swap and Poll
		app.window.SwapBuffers()
		glfw.PollEvents()

	}
}
