package utils

// https://learnopengl.com/Advanced-OpenGL/Cubemaps

import (
	"image"

	"github.com/go-gl/gl/v4.6-core/gl"
)

// Cubemap Cubemap
type Cubemap struct {
	ID     uint32
	Width  [6]int32
	Height [6]int32
	RGBA   [6]*image.RGBA
}

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

// NewCubemap NewCubemap
func NewCubemap(file []string, directory string) Cubemap {
	var texture Cubemap
	gl.GenTextures(1, &texture.ID)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, texture.ID)

	var width, height int32
	for i := 0; i < len(file); i++ {
		texture.RGBA[i] = LoadRGBA(directory + file[i])
		texture.Width[i] = int32(texture.RGBA[i].Rect.Size().X)
		texture.Height[i] = int32(texture.RGBA[i].Rect.Size().Y)
		gl.TexImage2D(uint32(gl.TEXTURE_CUBE_MAP_POSITIVE_X+i),
			0, gl.RGB, width, height, 0, gl.RGB, gl.UNSIGNED_BYTE, gl.Ptr(texture.RGBA[i].Pix))

	}
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)

	return texture
}
