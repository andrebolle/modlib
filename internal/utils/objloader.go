package utils

import (
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/udhos/gwob"
)

// Array of indices (+ve) defined with []int
/*
type Obj struct {
	Indices 				[]int
	Coord   				[]float32 // vertex data pos=(x,y,z) tex=(tx,ty) norm=(nx,ny,nz)
	Mtllib  				string
	Groups  				[]*Group

	BigIndexFound  			bool // index larger than 65535
	TextCoordFound 			bool // texture coord
	NormCoordFound 			bool // normal coord

	StrideSize           int // (px,py,pz),(tu,tv),(nx,ny,nz) = 8 x 4-byte floats = 32 bytes max
	StrideOffsetPosition int // 0
	StrideOffsetTexture  int // 3 x 4-byte floats
	StrideOffsetNormal   int // 5 x 4-byte floats
}
*/

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

// OJBLoader OJBLoader
func OJBLoader(filename string) *Vao {



	options := &gwob.ObjParserOptions{} // parser options

	// Obj holds parser result for .obj file.
// type Obj struct {
// 	Indices []int
// 	Coord   []float32 // vertex data pos=(x,y,z) tex=(tx,ty) norm=(nx,ny,nz)
// 	Mtllib  string
// 	Groups  []*Group

// 	BigIndexFound  bool // index larger than 65535
// 	TextCoordFound bool // texture coord
// 	NormCoordFound bool // normal coord

// 	StrideSize           int // (px,py,pz),(tu,tv),(nx,ny,nz) = 8 x 4-byte floats = 32 bytes max
// 	StrideOffsetPosition int // 0
// 	StrideOffsetTexture  int // 3 x 4-byte floats
// 	StrideOffsetNormal   int // 5 x 4-byte floats
// }

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

	return &vao
}
