package utils

import "github.com/go-gl/gl/v4.6-core/gl"

var skyboxVertices = []float32{
	// Positions
	-1.0, 1.0, -1.0,
	-1.0, -1.0, -1.0,
	1.0, -1.0, -1.0,
	1.0, -1.0, -1.0,
	1.0, 1.0, -1.0,
	-1.0, 1.0, -1.0,

	-1.0, -1.0, 1.0,
	-1.0, -1.0, -1.0,
	-1.0, 1.0, -1.0,
	-1.0, 1.0, -1.0,
	-1.0, 1.0, 1.0,
	-1.0, -1.0, 1.0,

	1.0, -1.0, -1.0,
	1.0, -1.0, 1.0,
	1.0, 1.0, 1.0,
	1.0, 1.0, 1.0,
	1.0, 1.0, -1.0,
	1.0, -1.0, -1.0,

	-1.0, -1.0, 1.0,
	-1.0, 1.0, 1.0,
	1.0, 1.0, 1.0,
	1.0, 1.0, 1.0,
	1.0, -1.0, 1.0,
	-1.0, -1.0, 1.0,

	-1.0, 1.0, -1.0,
	1.0, 1.0, -1.0,
	1.0, 1.0, 1.0,
	1.0, 1.0, 1.0,
	-1.0, 1.0, 1.0,
	-1.0, 1.0, -1.0,

	-1.0, -1.0, -1.0,
	-1.0, -1.0, 1.0,
	1.0, -1.0, -1.0,
	1.0, -1.0, -1.0,
	-1.0, -1.0, 1.0,
	1.0, -1.0, 1.0,
}

func loadCubemap(faces []string) uint32 {
	var textureID uint32
	gl.GenTextures(1, &textureID)

	var imageWidth, imageHeight int
	var image *ubyte

	gl.BindTexture(gl.TEXTURE_CUBE_MAP, textureID)

	for i = 0; i < len(faces); i++ {
		image := SOIL_load_image(faces[i], &imageWidth, &imageHeight, 0, SOIL_LOAD_RGB)
		glTexImage2D(gl.TEXTURE_CUBE_MAP_POSITIVE_X+i, 0, gl.RGB, imageWidth, imageHeight, 0, gl.RGB, gl.UNSIGNED_BYTE, image)
		SOIL_free_image_data(image)
	}
	glTexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	glTexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	glTexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	glTexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	glTexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)
	glBindTexture(gl.TEXTURE_CUBE_MAP, 0)

	return textureID
}

var skyboxVAO, skyboxVBO uint32

// InitSkyBox InitSkyBox
func InitSkyBox() {
	// Setup skybox VAO

	gl.GenVertexArrays(1, &skyboxVAO)
	gl.GenBuffers(1, &skyboxVBO)
	gl.BindVertexArray(skyboxVAO)
	gl.BindBuffer(gl.ARRAY_BUFFER, skyboxVBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(skyboxVertices)*4, gl.Ptr(skyboxVertices), gl.STATIC_DRAW)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))
	gl.BindVertexArray(0)
	//doink!
}
