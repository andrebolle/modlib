package main

import (
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/purelazy/modlib/internal/utils"
)

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
