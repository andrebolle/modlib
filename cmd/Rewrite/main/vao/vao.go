package vao

import (
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
)

// VAO VAO
type VAO struct {
	VAO         uint32
	Buffer      uint32
	Program     uint32
	VertexCount int32
}

// NewVAO NewVAO
func NewVAO(program uint32, floats *[]float32, size int32) *VAO {
	vao := new(VAO)
	vao.VertexCount = int32(len(*floats) / int(size))
	vao.Program = program
	gl.GenVertexArrays(1, &vao.VAO)
	gl.BindVertexArray(vao.VAO)
	gl.GenBuffers(1, &vao.Buffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, vao.Buffer)
	gl.BufferData(gl.ARRAY_BUFFER, len(*floats)*4, unsafe.Pointer(&(*floats)[0]), gl.STATIC_DRAW)
	gl.UseProgram(program)
	gl.VertexAttribPointer(0, size, gl.FLOAT, false, 0, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)
	return vao
}

// Draw Draw
func (vao *VAO) Draw() {
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.BindVertexArray(vao.VAO)
	gl.DrawArrays(gl.TRIANGLES, 0, vao.VertexCount)
	gl.Flush()
}
