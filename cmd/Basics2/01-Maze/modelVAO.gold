package main

import (
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/purelazy/modlib/internal/utils"
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

func setupModel(file string, lighting uint32, projection *float32) (uint32, *[]uint32, map[string]int32) {

	// Load the model geometry
	floats, indices, stride, posOffset, texOffset, normOffset := utils.OJBLoader(file)

	// Use program to get locations
	gl.UseProgram(lighting)

	// Get vertex attribute locations
	vertexAttributeNames := []string{
		"aPos", "aUV", "aNormal"}
	attrLocs := attrLocs(lighting, vertexAttributeNames)

	// Get uniform locations
	uniformNames := []string{
		"uAngle", "uModel", "uView", "uProjection", "uTex", "uViewPos", "uLightColor", "uLightPos"}
	uniLocs := uniLocs(lighting, uniformNames)

	// Compute and set static uniforms
	lightColor := mgl32.Vec3{1, 1, 1}
	lightPos := mgl32.Vec3{3, 3, -13}

	gl.UniformMatrix4fv(uniLocs["uProjection"], 1, false, projection)
	gl.Uniform1i(uniLocs["uTex"], 0)
	gl.Uniform3fv(uniLocs["uLightPos"], 1, &lightPos[0])
	gl.Uniform3fv(uniLocs["uLightColor"], 1, &lightColor[0])

	// -------------------------  VAO, EBO, VBO
	var cubeVAO, cubeVBO, cubeEBO uint32

	gl.GenVertexArrays(1, &cubeVAO)
	gl.BindVertexArray(cubeVAO)

	gl.GenBuffers(1, &cubeVBO)
	gl.BindBuffer(gl.ARRAY_BUFFER, cubeVBO)

	// For each atrribute {EnableVertexAttribArray, VertexAttribPointer}
	gl.EnableVertexAttribArray(attrLocs["aPos"])
	gl.VertexAttribPointer(attrLocs["aPos"], 3, gl.FLOAT, false, int32(stride), gl.PtrOffset(posOffset))
	gl.EnableVertexAttribArray(attrLocs["aUV"])
	gl.VertexAttribPointer(attrLocs["aUV"], 2, gl.FLOAT, false, int32(stride), gl.PtrOffset(texOffset))
	gl.EnableVertexAttribArray(attrLocs["aNormal"])
	gl.VertexAttribPointer(attrLocs["aNormal"], 3, gl.FLOAT, false, int32(stride), gl.PtrOffset(normOffset))

	gl.GenBuffers(1, &cubeEBO)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, cubeEBO)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(*indices)*4, unsafe.Pointer(&(*indices)[0]), gl.STATIC_DRAW)

	gl.BufferData(gl.ARRAY_BUFFER, len(*floats)*4, gl.Ptr(&(*floats)[0]), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	return cubeVAO, indices, uniLocs
}
