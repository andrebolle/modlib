package main

import (
	_ "image/png"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/purelazy/modlib/internal/utils"
)

func main() {

	// Model Data
	floats, indices, stride, posOffset, texOffset, normOffset := OJBLoader("suzanne.obj")

	// Window, Camera
	window, cam := utils.GetWindowAndCamera(800, 600)
	defer window.Destroy()

	// Program
	program := utils.NewProgram(utils.ReadShader("Lighting.vs.glsl"), utils.ReadShader("Lighting.fs.glsl"))
	defer gl.DeleteProgram(program)
	gl.UseProgram(program)

	// Program
	skybox := utils.NewProgram(utils.ReadShader("Skybox.vs.glsl"), utils.ReadShader("Skybox.fs.glsl"))
	defer gl.DeleteProgram(skybox)

	// Create Array, Buffer, Index Buffer, Texture
	vao := utils.NewArray()
	utils.NewBuffer(floats)
	utils.NewIndices(indices)
	tex := utils.NewTexture("square.png")
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, tex.Width, tex.Height, 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(tex.RGBA.Pix))

	// Attributes
	vao.Attribute(program, "aPos", 3, gl.FLOAT, false, int32(stride), gl.PtrOffset(posOffset))
	vao.Attribute(program, "aUV", 2, gl.FLOAT, false, int32(stride), gl.PtrOffset(texOffset))
	vao.Attribute(program, "aNormal", 3, gl.FLOAT, false, int32(stride), gl.PtrOffset(normOffset))

	// Uniforms
	// Projection
	projection := mgl32.Perspective(cam.Fovy, cam.Aspect, cam.Near, cam.Far)
	utils.SetUniformMat4(program, "uProjection", &projection[0])

	// Texture
	textureLocation := gl.GetUniformLocation(program, gl.Str("uTex\x00"))
	gl.Uniform1i(textureLocation, 0)

	// Skybox
	skyboxFiles := []string{"right.jpg", "left.jpg", "top.jpg", "bottom.jpg", "front.jpg", "back.jpg"}
	directory := "skybox/"
	utils.NewCubemap(skyboxFiles, directory)

	// Lighting
	lightColor := mgl32.Vec3{1, 1, 1}
	utils.SetUniformVec3(program, "uLightColor", &lightColor[0])
	lightPos := mgl32.Vec3{3, 3, 3}
	utils.SetUniformVec3(program, "uLightPos", &lightPos[0])

	// Pre Draw Setup
	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.CULL_FACE)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0, 0, 0, 1.0)

	// Rotation angle for animation
	angle := 0.0
	previousTime := glfw.GetTime()

	// Render Loop
	for !window.ShouldClose() {

		// Update the rotation angle
		time := glfw.GetTime()
		elapsed := time - previousTime
		previousTime = time
		angle += elapsed * 0.1

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
