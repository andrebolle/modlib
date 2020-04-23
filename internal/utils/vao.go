package utils

import (
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
)

// VAO Vertex Array Object
type VAO struct {
	ID       uint32
	Bound    bool
	DrawMode int
}

// Attribute Attribute
func Attribute(program uint32, name string, size int32, xtype uint32, normalized bool, stride int32, offset unsafe.Pointer) {
	vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str(name+"\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, size, gl.FLOAT, false, int32(stride), offset)
}

// NewArray NewArray
func NewArray() VAO {
	var vao VAO
	gl.GenVertexArrays(1, &vao.ID)
	gl.BindBuffer(gl.ARRAY_BUFFER, vao.ID)
	vao.Bound = true
	return vao
}
