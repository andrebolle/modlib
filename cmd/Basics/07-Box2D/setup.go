package main

import (
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/purelazy/modlib/internal/utils"
)

func setupModel(lighting uint32, projection *float32) (uint32, int32, int32, int32, *[]uint32) {
	var uModelLocation, uViewLocation, uViewPosLocation int32

	// Load the model geometry
	floats, indices, stride, posOffset, texOffset, normOffset := utils.OJBLoader("cube.obj")

	// Use program to get locations
	gl.UseProgram(lighting)
	// ---------------------- Get locations
	aPosLocation := uint32(gl.GetAttribLocation(lighting, gl.Str("aPos\x00")))
	aUVLocation := uint32(gl.GetAttribLocation(lighting, gl.Str("aUV\x00")))
	aNormalLocation := uint32(gl.GetAttribLocation(lighting, gl.Str("aNormal\x00")))

	uModelLocation = gl.GetUniformLocation(lighting, gl.Str("uModel\x00"))
	uViewLocation = gl.GetUniformLocation(lighting, gl.Str("uView\x00"))
	uProjectionLocation := gl.GetUniformLocation(lighting, gl.Str("uProjection\x00"))
	uTexLocation := gl.GetUniformLocation(lighting, gl.Str("uTex\x00"))
	uViewPosLocation = gl.GetUniformLocation(lighting, gl.Str("uViewPos\x00"))
	uLightColourLocation := gl.GetUniformLocation(lighting, gl.Str("uLightColor\x00"))
	uLightPosLocation := gl.GetUniformLocation(lighting, gl.Str("uLightPos\x00"))

	// ------------------------- Compute and set static uniforms
	// := mgl32.Perspective(cam.Fovy, cam.Aspect, cam.Near, cam.Far)
	lightColor := mgl32.Vec3{1, 1, 1}
	lightPos := mgl32.Vec3{3, 3, 3}

	gl.UniformMatrix4fv(uProjectionLocation, 1, false, projection)
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

	return cubeVAO, uModelLocation, uViewLocation, uViewPosLocation, indices
}

func setupSkybox(cubemapShader uint32, projection *float32) (uint32, int32) {
	var uViewCubemapLocation int32
	//  --------------------------------------------------- Skybox
	// Use program to get locations
	gl.UseProgram(cubemapShader)
	// -------------------- Get locations
	uViewCubemapLocation = gl.GetUniformLocation(cubemapShader, gl.Str("uView\x00"))
	uProjectionCubemapLocation := gl.GetUniformLocation(cubemapShader, gl.Str("uProjection\x00"))
	uTexCubemapLocation := gl.GetUniformLocation(cubemapShader, gl.Str("uTex\x00"))

	// ------------------- Set static uniforms
	gl.UniformMatrix4fv(uProjectionCubemapLocation, 1, false, projection)
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
	return skyboxVAO, uViewCubemapLocation
}
