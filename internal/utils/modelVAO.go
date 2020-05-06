package utils

import (
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

func uniLocs(program uint32, names []string) map[string]int32 {
	uniLocs := map[string]int32{}
	for _, name := range names {
		uniLocs[name] = gl.GetUniformLocation(program, gl.Str(name+"\x00"))
	}
	return uniLocs
}

func attrLocs(program uint32, names []string) map[string]uint32 {
	attrLocs := map[string]uint32{}
	for _, name := range names {
		attrLocs[name] = uint32(gl.GetAttribLocation(program, gl.Str(name+"\x00")))
	}
	return attrLocs
}

// Vao Vao
type Vao struct {
	Vao, Vbo, Ebo                      uint32
	Indices                            *[]uint32
	AttrLocs                           map[string]uint32
	UniLocs                            map[string]int32
	Stride                             int32
	PosOffset, TexOffset, NormalOffset unsafe.Pointer
}

// SetupModel SetupModel
func SetupModel(file string, lighting uint32, projection *float32) *Vao {

	vao := Vao{}

	// Load the model geometry
	obj, indices := OJBLoader(file)

	vao.Indices = indices
	vao.Stride = int32(obj.StrideSize)
	vao.PosOffset = gl.PtrOffset(obj.StrideOffsetPosition)
	vao.TexOffset = gl.PtrOffset(obj.StrideOffsetTexture)
	vao.NormalOffset = gl.PtrOffset(obj.StrideOffsetNormal)

	// Use program to get locations
	gl.UseProgram(lighting)

	// Get vertex attribute and uniform locations
	vertexAttributes := []string{"aPos", "aUV", "aNormal"}
	vao.AttrLocs = attrLocs(lighting, vertexAttributes)
	uniforms := []string{"uAngle", "uModel", "uView", "uProjection",
		"uTex", "uViewPos", "uLightColor", "uLightPos", "uPosAngle"}
	vao.UniLocs = uniLocs(lighting, uniforms)

	// Compute and set static uniforms
	lightColor := mgl32.Vec3{1, 1, 1}
	lightPos := mgl32.Vec3{3, 3, -13}
	gl.UniformMatrix4fv(vao.UniLocs["uProjection"], 1, false, projection)
	gl.Uniform1i(vao.UniLocs["uTex"], 0)
	gl.Uniform3fv(vao.UniLocs["uLightPos"], 1, &lightPos[0])
	gl.Uniform3fv(vao.UniLocs["uLightColor"], 1, &lightColor[0])

	// Create & Bind VAO and its buffer
	gl.GenVertexArrays(1, &vao.Vao)
	gl.BindVertexArray(vao.Vao)
	gl.GenBuffers(1, &vao.Vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vao.Vbo)

	// For each atrribute {EnableVertexAttribArray, VertexAttribPointer}
	gl.EnableVertexAttribArray(vao.AttrLocs["aPos"])
	gl.VertexAttribPointer(vao.AttrLocs["aPos"], 3, gl.FLOAT, false, vao.Stride, vao.PosOffset)
	gl.EnableVertexAttribArray(vao.AttrLocs["aUV"])
	gl.VertexAttribPointer(vao.AttrLocs["aUV"], 2, gl.FLOAT, false, vao.Stride, vao.TexOffset)
	gl.EnableVertexAttribArray(vao.AttrLocs["aNormal"])
	gl.VertexAttribPointer(vao.AttrLocs["aNormal"], 3, gl.FLOAT, false, vao.Stride, vao.NormalOffset)

	gl.GenBuffers(1, &vao.Ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, vao.Ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(*vao.Indices)*4, unsafe.Pointer(&(*vao.Indices)[0]), gl.STATIC_DRAW)

	gl.BufferData(gl.ARRAY_BUFFER, len(obj.Coord)*4, gl.Ptr(&(obj.Coord)[0]), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	//return vao.cubeVAO, vao.indices, vao.uniLocs
	return &vao
}
