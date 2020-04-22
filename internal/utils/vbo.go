package utils

import (
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
)

// VBO Vertex Buffer Object
type VBO struct {
	id    uint32
	bound bool
}

// NewBuffer NewBuffer
func NewBuffer(vbo *VBO, floats *[]float32) {
	gl.GenVertexArrays(1, &vbo.id)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo.id)
	gl.BufferData(gl.ARRAY_BUFFER, len(*floats)*4, unsafe.Pointer(&(*floats)[0]), gl.STATIC_DRAW)
	vbo.bound = true
}
