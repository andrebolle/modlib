package utils

import (
	"fmt"
	"unsafe"

	"github.com/ByteArena/box2d"
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/udhos/gwob"
)

// Vao Vao
type Vao struct {
	Vao, Vbo, Ebo                      uint32
	Pos					*[]float32	
	UVs					*[]float32	
	Norms					*[]float32	
	PosAndAngle	*[]float32
	PosAndAngleOffset int
	Coord   *[]float32
	Indices                            *[]uint32
	AttrLocs                           map[string]uint32
	UniLocs                            map[string]int32
	Stride                             int32
	PosOffset, TexOffset, NormalOffset unsafe.Pointer
}

// GetVAOData GetVAOData
func GetVAOData(filename string) *Vao {

	options := &gwob.ObjParserOptions{} // parser options

	obj, errObj := gwob.NewObjFromFile(filename, options) // parse/load OBJ

	//fmt.Println(o, errObj)
	if errObj != nil {
		panic(errObj)
	}

	uIntIndices := make([]uint32, 0)
	// Convert index "ints" to "uints"
	for i := range obj.Indices {
		uIntIndices = append(uIntIndices, uint32(obj.Indices[i]))
	}
	
	vao := Vao{}

	vao.Coord = &obj.Coord // vertex data pos=(x,y,z) tex=(tx,ty) norm=(nx,ny,nz)
	vao.Indices = &uIntIndices
	vao.Stride = int32(obj.StrideSize)
	vao.PosOffset = gl.PtrOffset(obj.StrideOffsetPosition)
	vao.TexOffset = gl.PtrOffset(obj.StrideOffsetTexture)
	vao.NormalOffset = gl.PtrOffset(obj.StrideOffsetNormal)

	// Get Pos, UVs and Norms into their own slices
	vao.Pos, vao.UVs, vao.Norms = deInterlace(vao.Coord)

	return &vao
}


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

func findAttrAndUniLocs(vao *Vao, program uint32, projection *float32) {
		// Use program to get locations
		gl.UseProgram(program)

		// Get vertex attribute and uniform locations
		vertexAttributes := []string{"aPos", "aUV", "aNormal", "aInstancePosAngle"}
		vao.AttrLocs = attrLocs(program, vertexAttributes)
		uniforms := []string{"uModel", "uView", "uProjection",
			"uTex", "uViewPos", "uLightColor", "uLightPos"}
		vao.UniLocs = uniLocs(program, uniforms)
	
		// Compute and set static uniforms
		lightColor := mgl32.Vec3{1, 1, 1}
		lightPos := mgl32.Vec3{3, 3, -13}
		gl.UniformMatrix4fv(vao.UniLocs["uProjection"], 1, false, projection)
		gl.Uniform1i(vao.UniLocs["uTex"], 0)
		gl.Uniform3fv(vao.UniLocs["uLightPos"], 1, &lightPos[0])
		gl.Uniform3fv(vao.UniLocs["uLightColor"], 1, &lightColor[0])
}

// GetPositionAndAngle GetPositionAndAngle
func GetPositionAndAngle(world *box2d.B2World, name string) *[]float32 {
	posAndAngle := make([]float32,0)
	for b := world.GetBodyList(); b != nil; b = b.GetNext() {
		if b.GetUserData() == name {
			posAndAngle = append(posAndAngle, float32(b.GetPosition().X), float32(b.GetPosition().Y), float32(b.GetAngle()) )
		}
	}
	return &posAndAngle
}

func deInterlace(posUVsNorms *[]float32) (*[]float32, *[]float32, *[]float32) {
	positions, uvs, norms := make([]float32, 0), make([]float32, 0), make([]float32, 0)

	fmt.Println("posUVsNorms", len(*posUVsNorms))

	for i := 0; i < len(*posUVsNorms); {
		positions = append(positions, (*posUVsNorms)[i], (*posUVsNorms)[i+1], (*posUVsNorms)[i+2] )
		i += 3
		uvs = append(uvs, (*posUVsNorms)[i], (*posUVsNorms)[i+1] )
		i += 2
		norms = append(norms, (*posUVsNorms)[i], (*posUVsNorms)[i+1], (*posUVsNorms)[i+2] )
		i += 3
	}

	return &positions, &uvs, &norms
}
// SetupModel SetupModel
func SetupModel(file string, program uint32, projection *float32, world *box2d.B2World) *Vao {

	// Short name
	null := unsafe.Pointer(nil)

	// Load the model geometry and return a VAO
	vao := GetVAOData(file)

	// Put all Box2D position and angle values in their own slice
	vao.PosAndAngle = GetPositionAndAngle(world, "box")

	// This function must be called before vao.AttrLocs is used anywhere
	findAttrAndUniLocs(vao, program, projection) 

	// Create & Bind VAO and its buffer
	gl.GenVertexArrays(1, &vao.Vao)
	gl.BindVertexArray(vao.Vao)
	gl.GenBuffers(1, &vao.Vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vao.Vbo)

	fmt.Println("Array Buffer Size", (len(*vao.Pos) + len(*vao.UVs) + len(*vao.Norms) + len(*vao.PosAndAngle))*4, "bytes" )

	// Allocate memory for array buffer
	gl.BufferData(gl.ARRAY_BUFFER, (len(*vao.Pos) + len(*vao.UVs) + len(*vao.Norms) + len(*vao.PosAndAngle))*4, null, gl.DYNAMIC_DRAW)

	// Copy Pos, UVs, Norms, Position and Angle data to ARRAY_BUFFER
	offset := 0
	gl.BufferSubData(gl.ARRAY_BUFFER, offset, len(*vao.Pos)*4, gl.Ptr(*vao.Pos))
	offset += len(*vao.Pos)*4
	gl.BufferSubData(gl.ARRAY_BUFFER, offset, len(*vao.UVs)*4, gl.Ptr(*vao.UVs))
	offset += len(*vao.UVs)*4
	gl.BufferSubData(gl.ARRAY_BUFFER, offset, len(*vao.Norms)*4, gl.Ptr(*vao.Norms))
	offset += len(*vao.Norms)*4
	gl.BufferSubData(gl.ARRAY_BUFFER, offset, len(*vao.PosAndAngle)*4, gl.Ptr(*vao.PosAndAngle))
	vao.PosAndAngleOffset = offset

	// Define Vertex Shader Attributes
	gl.VertexAttribPointer(vao.AttrLocs["aPos"], 3, gl.FLOAT, false, 0, null)
	gl.VertexAttribPointer(vao.AttrLocs["aUV"], 2, gl.FLOAT, false, 0, gl.PtrOffset(len(*vao.Pos)*4))
	gl.VertexAttribPointer(vao.AttrLocs["aNormal"], 3, gl.FLOAT, false, 0, gl.PtrOffset((len(*vao.Pos) + len(*vao.UVs))*4))
	gl.VertexAttribPointer(vao.AttrLocs["aInstancePosAngle"], 3, gl.FLOAT, false, 0, 
		gl.PtrOffset((len(*vao.Pos) + len(*vao.UVs) + len(*vao.Norms))*4))

	// Enable all vertex attributes
	gl.EnableVertexAttribArray(vao.AttrLocs["aPos"])
	gl.EnableVertexAttribArray(vao.AttrLocs["aUV"])
	gl.EnableVertexAttribArray(vao.AttrLocs["aNormal"])
	gl.EnableVertexAttribArray(vao.AttrLocs["aInstancePosAngle"])

	// Modify the rate at which generic vertex attributes advance during instanced rendering
	// A "PosAngle" for every 1 instance
	gl.VertexAttribDivisor(vao.AttrLocs["aInstancePosAngle"], 1)

	// Copy indices
	gl.GenBuffers(1, &vao.Ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, vao.Ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(*vao.Indices)*4, gl.Ptr(*vao.Indices), gl.STATIC_DRAW)

	return vao
}
