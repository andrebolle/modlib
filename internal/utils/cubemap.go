package utils

// https://learnopengl.com/Advanced-OpenGL/Cubemaps

import (
	"image"

	"github.com/go-gl/gl/v4.6-core/gl"
)

var skyboxVertices = []float32{
	// positions
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

// Cubemap Cubemap
type Cubemap struct {
	ID     uint32
	Vao    VAO
	Vbo    VBO
	Width  [6]int32
	Height [6]int32
	RGBA   [6]*image.RGBA
}

// NewCubemap NewCubemap
func NewCubemap(file []string, directory string) Cubemap {

	var cubemap Cubemap

	cubemap.Vao = NewArray()
	cubemap.Vbo = NewBuffer(&skyboxVertices)

	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))
	gl.GenTextures(1, &cubemap.ID)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, cubemap.ID)

	var width, height int32
	for i := 0; i < len(file); i++ {
		cubemap.RGBA[i] = LoadRGBA(directory + file[i])
		cubemap.Width[i] = int32(cubemap.RGBA[i].Rect.Size().X)
		cubemap.Height[i] = int32(cubemap.RGBA[i].Rect.Size().Y)
		gl.TexImage2D(uint32(gl.TEXTURE_CUBE_MAP_POSITIVE_X+i),
			0, gl.RGB, width, height, 0, gl.RGB, gl.UNSIGNED_BYTE, gl.Ptr(cubemap.RGBA[i].Pix))

	}
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)

	return cubemap
}
