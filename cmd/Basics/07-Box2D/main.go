// Copyright 2014 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Renders a textured spinning cube using GLFW 3 and OpenGL 4.1 core forward-compatible profile.
package main

import (
	_ "image/png"
	"math"
	"math/rand"

	"github.com/ByteArena/box2d"
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/purelazy/modlib/internal/utils"
)

func main() {

	// Box2D is tuned for meters, kilograms, and seconds.
	// Define the gravity vector.
	gravity := box2d.MakeB2Vec2(0.0, 0.0)

	// Construct a world object, which will hold and simulate the rigid bodies.
	world := box2d.MakeB2World(gravity)

	// Short names
	//kinetic := box2d.B2BodyType.B2_kinematicBody
	dynamic := box2d.B2BodyType.B2_dynamicBody

	// Timestep
	timeStep := 1.0 / 60.0
	velocityIterations := 8
	positionIterations := 3

	// A place to store bodies by name
	//characters := make(map[string]*box2d.B2Body)

	// 1. Create a bodydef. Initial Position, Type (Dynamic, Static,  Kinematic)
	// 2. Create the body from the def.
	// 3. Create a shape - Polygon, Chain, Circle
	// 4. Create a fixture - glues body and shape. denisty, friction, restitution

	// ----------------- Left
	leftBodyDef := box2d.MakeB2BodyDef()
	leftBodyDef.Position.Set(-20, 0)

	leftBody := world.CreateBody(&leftBodyDef)

	leftBox := box2d.MakeB2PolygonShape()

	leftBox.SetAsBox(19.0, 5000.0)

	leftBody.CreateFixture(&leftBox, 0.0)

	// ----------------- Right
	rightBodyDef := box2d.MakeB2BodyDef()
	rightBodyDef.Position.Set(20, 0)

	rightBody := world.CreateBody(&rightBodyDef)

	rightBox := box2d.MakeB2PolygonShape()

	rightBox.SetAsBox(10.0, 5000.0)

	rightBody.CreateFixture(&rightBox, 0.0)

	// ----------------- Ground
	groundBodyDef := box2d.MakeB2BodyDef()
	groundBodyDef.Position.Set(0.0, -20.0)

	groundBody := world.CreateBody(&groundBodyDef)

	groundBox := box2d.MakeB2PolygonShape()

	groundBox.SetAsBox(5000.0, 10.0)

	groundBody.CreateFixture(&groundBox, 0.0)

	// ----------------- Ceiling
	ceilingBodyDef := box2d.MakeB2BodyDef()
	ceilingBodyDef.Position.Set(0.0, 20.0)

	ceilingBody := world.CreateBody(&ceilingBodyDef)

	ceilingBox := box2d.MakeB2PolygonShape()

	ceilingBox.SetAsBox(5000.0, 10.0)

	ceilingBody.CreateFixture(&ceilingBox, 0.0)

	// ----------------- Box
	boxBodyDef := box2d.MakeB2BodyDef()
	boxBodyDef.Position.Set(0, 10)

	boxBodyDef.Type = dynamic
	boxBodyDef.AllowSleep = false
	boxBodyDef.LinearVelocity.Set(4, 2)
	// Body instance
	boxBody := world.CreateBody(&boxBodyDef)
	boxBody.SetTransform(box2d.B2Vec2{X: rand.Float64(), Y: rand.Float64()}, rand.Float64()*2*math.Pi)

	// Create a box shape
	boxShape := box2d.MakeB2PolygonShape()
	boxShape.SetAsBox(1.0, 1.0)

	// Fixture
	// A fixture binds a shape to a body and adds material properties such as density, friction, and restitution.
	// A fixture puts a shape into the collision system (broad-phase) so that it can collide with other shapes.
	fixtureDef := box2d.MakeB2FixtureDef()
	fixtureDef.Shape = &boxShape
	fixtureDef.Density = 1.0
	fixtureDef.Friction = 0.0
	fixtureDef.Restitution = 1
	boxBody.CreateFixtureFromDef(&fixtureDef)

	// physObjList := world.GetBodyList()
	// for i := 0; physObjList != nil; i++ {
	// 	fmt.Println(i)
	// 	physObjList = physObjList.GetNext()
	// }

	// Create the OpenGL context, window and camera
	window, cam := utils.GetWindowAndCamera(800, 600)
	defer window.Destroy()

	// Load textures
	cubemapTexture := utils.Cubemap(utils.Faces)
	modelTexture := utils.NewTexture("square.png")

	// Compile model and cubemap shaders
	lighting := utils.NewProgram(utils.ReadShader("Lighting.vs.glsl"), utils.ReadShader("Lighting.fs.glsl"))
	cubemapShader := utils.NewProgram(utils.ReadShader("cubemap.vs.glsl"), utils.ReadShader("cubemap.fs.glsl"))
	defer gl.DeleteProgram(lighting)
	defer gl.DeleteProgram(cubemapShader)

	// ------------------------- Compute and set static uniforms
	projection := mgl32.Perspective(cam.Fovy, cam.Aspect, cam.Near, cam.Far)

	cubeVAO, uModelLocation, uViewLocation, uViewPosLocation, indices := setupModel(lighting, &projection[0])

	skyboxVAO, uViewCubemapLocation := setupSkybox(cubemapShader, &projection[0])

	for !window.ShouldClose() {

		// View is used in multiple programs
		view := mgl32.LookAtV(cam.Position, cam.Position.Add(cam.Forward), cam.Up)

		{ // ----------------Physics
			// Reverse gravity
			if cam.Paused {
				boxBody.SetGravityScale(-1)
			} else {
				boxBody.SetGravityScale(1)
			}

			// Step through time
			world.Step(timeStep, velocityIterations, positionIterations)
		}

		{ // ----------------Draw the skybox (36 verts)
			gl.UseProgram(cubemapShader)
			// Drawing the skybox first will draw every pixel, so the screen does not
			// need to be cleared and not depth testing
			gl.Disable(gl.DEPTH_TEST)

			// Remove translation from the view matrix. i.e. the skybox stays in the same place.
			viewWithoutTranslation := view.Mat3().Mat4()
			gl.UniformMatrix4fv(uViewCubemapLocation, 1, false, &viewWithoutTranslation[0])

			gl.BindVertexArray(skyboxVAO)
			gl.ActiveTexture(gl.TEXTURE0)
			gl.BindTexture(gl.TEXTURE_CUBE_MAP, cubemapTexture)
			// Draw the VAO that is currently bound
			gl.DrawArrays(gl.TRIANGLES, 0, 36)
		}

		{ // ----------------Draw the model
			gl.UseProgram(lighting)
			gl.Enable(gl.CULL_FACE) // Only front-facing triangles will be drawn

			// Dynamic uniforms
			position := boxBody.GetPosition()
			angle := boxBody.GetAngle()
			rotate := mgl32.HomogRotate3D(float32(angle), mgl32.Vec3{0, 0, 1})

			translate := mgl32.Translate3D(float32(position.X), float32(position.Y), 0)
			model := translate.Mul4(rotate)

			gl.UniformMatrix4fv(uViewLocation, 1, false, &view[0])
			gl.UniformMatrix4fv(uModelLocation, 1, false, &model[0])
			gl.Uniform3fv(uViewPosLocation, 1, &cam.Position[0])

			gl.BindVertexArray(cubeVAO)
			gl.ActiveTexture(gl.TEXTURE0)
			gl.BindTexture(gl.TEXTURE_2D, modelTexture.ID)
			gl.DrawElements(gl.TRIANGLES, int32(len(*indices)), gl.UNSIGNED_INT, gl.PtrOffset(0))
		}

		// Swap and Poll
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
