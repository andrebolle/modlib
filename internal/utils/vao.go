package utils

import (
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
)

// VAO Vertex Array Object
type VAO struct {
	id       uint32
	bound    bool
	drawMode int
}

// Attribute Attribute
func (*VAO) Attribute(program uint32, name string, size int32, xtype uint32, normalized bool, stride int32, offset unsafe.Pointer) {
	vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str(name+"\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, int32(stride), offset)
}

// NewVAOGenBind NewVAOGenBind
func NewVAOGenBind(vao *VAO) {
	gl.GenVertexArrays(1, &vao.id)
	gl.BindBuffer(gl.ARRAY_BUFFER, vao.id)
	vao.bound = true
}
