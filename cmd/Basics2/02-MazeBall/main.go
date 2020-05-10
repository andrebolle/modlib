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
	width, height         int
	cam                   *utils.Camera
	projection            mgl32.Mat4
	mazeWidth, mazeHeight uint8
	wallCount             int
	world                 box2d.B2World
	nutVAO                *utils.Vao
}

func theBall(world *box2d.B2World) {

}

func main() {
	// Define the maze. Dimensions must be odd
	//app := App{width: 1680, height: 1050, mazeWidth: 31, mazeHeight: 31}
	app := App{width: 1680, height: 1050, mazeWidth: 31, mazeHeight: 31}

	// Camera
	app.cam = utils.Cam()
	app.cam.Aspect = float32(app.width) / float32(app.height)
	app.cam.Position = app.cam.StartPosition

	// Window and OpenGL
	app.window = utils.GetWindow(app.width, app.height)
	defer app.window.Destroy()

	// Callbacks
	utils.SetWASDCallback(app.window, app.cam)
	utils.SetPitchYawCallback(app.window, app.cam)

	// Create a Box2D world
	app.world = box2d.MakeB2World(box2d.MakeB2Vec2(0.0, 0.0))

	// Build a maze and return the number of walls to draw
	app.wallCount = buildMaze(&app.world, designMaze(app.mazeWidth, app.mazeHeight))

	// Add a ball
	addBox(&app.world, box2d.B2Vec2{X: 40, Y: 40}, box2d.B2Vec2{X: -6.8, Y: -6.8}, box2d.B2Vec2{X: .1, Y: .1}, dynamic)
	addBox(&app.world, box2d.B2Vec2{X: 30, Y: 24}, box2d.B2Vec2{X: -6.8, Y: -6.8}, box2d.B2Vec2{X: .1, Y: .1}, dynamic)
	addBox(&app.world, box2d.B2Vec2{X: 20, Y: 20}, box2d.B2Vec2{X: -6.8, Y: -6.8}, box2d.B2Vec2{X: .1, Y: .1}, dynamic)
	addBox(&app.world, box2d.B2Vec2{X: 24, Y: 10}, box2d.B2Vec2{X: -6.8, Y: -6.8}, box2d.B2Vec2{X: .1, Y: .1}, dynamic)
	//addStaticBox(app.world, box2d.B2Vec2{X: float64(22) * 2.1, Y: float64(22) * 2.1}, box2d.B2Vec2{X: 1, Y: 1})
	app.wallCount++

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
	app.nutVAO = utils.SetupModel("cube.obj", lighting, &app.projection[0], &app.world)
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

		app.world.Step(1.0/60.0, 8, 3)

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

		// Extract position and angle from Box2D world
		posAndAngle := utils.GetPositionAndAngle(&app.world, "box")

		// Load posAndAngle into GPU
		gl.BufferSubData(gl.ARRAY_BUFFER, app.nutVAO.PosAndAngleOffset, len(*posAndAngle)*4, gl.Ptr(*posAndAngle))

		// Draw boxCount instances
		indicesCount := int32(len(*app.nutVAO.Indices))
		gl.DrawElementsInstanced(gl.TRIANGLES, indicesCount, gl.UNSIGNED_INT, gl.PtrOffset(0), int32(app.wallCount+3))

		// Swap and Poll
		app.window.SwapBuffers()
		glfw.PollEvents()

	}
}
