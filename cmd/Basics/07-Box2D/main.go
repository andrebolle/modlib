// Copyright 2014 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Renders a textured spinning cube using GLFW 3 and OpenGL 4.1 core forward-compatible profile.
package main

import (
	"fmt"
	_ "image/png"
	"math"
	"math/rand"
	"unsafe"

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

	// SHortnames
	//kinetic := box2d.B2BodyType.B2_kinematicBody
	dynamic := box2d.B2BodyType.B2_dynamicBody

	// Timestep
	timeStep := 1.0 / 60.0
	velocityIterations := 8
	positionIterations := 3
	//world.Step(timeStep, velocityIterations, positionIterations)

	// A place to store bodies by name
	//characters := make(map[string]*box2d.B2Body)

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

	for i := world.GetBodyList(); i != nil; i = i.GetNext() {
		fmt.Println(i)
	}

	// 1. Create a bodydef. Initial Position, Type (Dynamic, Static,  Kinematic)
	// 2. Create the body from the def.
	// 3. Create a shape - Polygon, Chain, Circle
	// 4. Create a fixture - glues body and shape. denisty, friction, restitution

	// Create the OpenGL context, window and camera
	window, cam := utils.GetWindowAndCamera(800, 600)
	defer window.Destroy()

	// Load the background and model textures
	cubemapTexture := utils.Cubemap(utils.Faces)
	modelTexture := utils.NewTexture("square.png")

	// Load the model geometry
	floats, indices, stride, posOffset, texOffset, normOffset := utils.OJBLoader("cube.obj")

	// Compile model and cubemap shaders
	lighting := utils.NewProgram(utils.ReadShader("Lighting.vs.glsl"), utils.ReadShader("Lighting.fs.glsl"))
	cubemapShader := utils.NewProgram(utils.ReadShader("cubemap.vs.glsl"), utils.ReadShader("cubemap.fs.glsl"))
	defer gl.DeleteProgram(lighting)
	defer gl.DeleteProgram(cubemapShader)

	// For each program

	// -------------------------------------- Cube ---------------------------------
	// Use program to get locations
	gl.UseProgram(lighting)
	// ---------------------- Get locations
	aPosLocation := uint32(gl.GetAttribLocation(lighting, gl.Str("aPos\x00")))
	aUVLocation := uint32(gl.GetAttribLocation(lighting, gl.Str("aUV\x00")))
	aNormalLocation := uint32(gl.GetAttribLocation(lighting, gl.Str("aNormal\x00")))

	uModelLocation := gl.GetUniformLocation(lighting, gl.Str("uModel\x00"))
	uViewLocation := gl.GetUniformLocation(lighting, gl.Str("uView\x00"))
	uProjectionLocation := gl.GetUniformLocation(lighting, gl.Str("uProjection\x00"))
	uTexLocation := gl.GetUniformLocation(lighting, gl.Str("uTex\x00"))
	uViewPosLocation := gl.GetUniformLocation(lighting, gl.Str("uViewPos\x00"))
	uLightColourLocation := gl.GetUniformLocation(lighting, gl.Str("uLightColor\x00"))
	uLightPosLocation := gl.GetUniformLocation(lighting, gl.Str("uLightPos\x00"))

	// ------------------------- Compute and set static uniforms
	projection := mgl32.Perspective(cam.Fovy, cam.Aspect, cam.Near, cam.Far)
	lightColor := mgl32.Vec3{1, 1, 1}
	lightPos := mgl32.Vec3{3, 3, 3}

	gl.UniformMatrix4fv(uProjectionLocation, 1, false, &projection[0])
	gl.Uniform1i(uTexLocation, 0)
	gl.Uniform3fv(uLightPosLocation, 1, &lightPos[0])
	gl.Uniform3fv(uLightColourLocation, 1, &lightColor[0])

	// -------------------------  VAO, EBO, VBO
	var cubeVAO, cubeVBO, cubeEBO uint32

	gl.GenVertexArrays(1, &cubeVAO)
	gl.BindVertexArray(cubeVAO)

	gl.GenBuffers(1, &cubeVBO)
	gl.BindBuffer(gl.ARRAY_BUFFER, cubeVBO)

	// For each atrribute {EnableVertexAttribArray, VertexAttribPointer}
	gl.EnableVertexAttribArray(aPosLocation)
	gl.VertexAttribPointer(aPosLocation, 3, gl.FLOAT, false, int32(stride), gl.PtrOffset(posOffset))
	gl.EnableVertexAttribArray(aUVLocation)
	gl.VertexAttribPointer(aUVLocation, 2, gl.FLOAT, false, int32(stride), gl.PtrOffset(texOffset))
	gl.EnableVertexAttribArray(aNormalLocation)
	gl.VertexAttribPointer(aNormalLocation, 3, gl.FLOAT, false, int32(stride), gl.PtrOffset(normOffset))

	gl.GenBuffers(1, &cubeEBO)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, cubeEBO)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(*indices)*4, unsafe.Pointer(&(*indices)[0]), gl.STATIC_DRAW)

	gl.BufferData(gl.ARRAY_BUFFER, len(*floats)*4, gl.Ptr(&(*floats)[0]), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	//  --------------------------------------------------- Skybox
	// Use program to get locations
	gl.UseProgram(cubemapShader)
	// -------------------- Get locations
	uViewCubemapLocation := gl.GetUniformLocation(cubemapShader, gl.Str("uView\x00"))
	uProjectionCubemapLocation := gl.GetUniformLocation(cubemapShader, gl.Str("uProjection\x00"))
	uTexCubemapLocation := gl.GetUniformLocation(cubemapShader, gl.Str("uTex\x00"))

	// ------------------- Set static uniforms
	gl.UniformMatrix4fv(uProjectionCubemapLocation, 1, false, &projection[0])
	gl.Uniform1i(uTexCubemapLocation, 0)

	// -------------------------  VAO, EBO, VBO
	var skyboxVAO, skyboxVBO uint32
	gl.GenVertexArrays(1, &skyboxVAO)
	gl.GenBuffers(1, &skyboxVBO)
	gl.BindVertexArray(skyboxVAO)
	gl.BindBuffer(gl.ARRAY_BUFFER, skyboxVBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(utils.SkyboxVertices)*4, gl.Ptr(&(utils.SkyboxVertices)[0]), gl.STATIC_DRAW)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))

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

		{ // ----------------Draw the skybox cube (36 verts)
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
