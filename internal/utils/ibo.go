package utils

import (
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
)

// IBO Index Buffer Object
type IBO struct {
	id uint32
}

// NewIndices NewIndices
func NewIndices(ibo *IBO, indices *[]uint32) {
	gl.GenBuffers(1, &ibo.id)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ibo.id)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(*indices)*4, unsafe.Pointer(&(*indices)[0]), gl.STATIC_DRAW)
}
