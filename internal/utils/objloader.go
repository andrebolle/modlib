package utils

import (
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

// OJBLoader OJBLoader
func OJBLoader(filename string) (*gwob.Obj, *[]uint32) {

	options := &gwob.ObjParserOptions{} // parser options

	obj, errObj := gwob.NewObjFromFile(filename, options) // parse/load OBJ

	//fmt.Println(o, errObj)
	if errObj != nil {
		panic(errObj)
	}

	// fmt.Println("Indices 			", len(obj.Indices))
	// fmt.Println("Coords   			", len(obj.Coord))
	// fmt.Println("Mtllib  			", obj.Mtllib)
	// fmt.Println("Groups  			", len(obj.Groups))
	// fmt.Println("BigIndexFound  	", obj.BigIndexFound)
	// fmt.Println("TextCoordFound 	", obj.TextCoordFound)
	// fmt.Println("NormCoordFound 	", obj.NormCoordFound)
	// fmt.Println("StrideSize         ", obj.StrideSize)
	// fmt.Println("StrideOffsetPosition", obj.StrideOffsetPosition)
	// fmt.Println("StrideOffsetTexture", obj.StrideOffsetTexture)
	// fmt.Println("StrideOffsetNormal", obj.StrideOffsetNormal)

	// fmt.Print(o.Indices)

	uIntIndices := make([]uint32, 0)
	// Convert index "ints" to "uints"
	for i := range obj.Indices {
		uIntIndices = append(uIntIndices, uint32(obj.Indices[i]))
	}

	//return &o.Coord, &o.Indices, o.StrideSize, o.StrideOffsetPosition, o.StrideOffsetTexture, o.StrideOffsetNormal
	//return &obj.Coord, &uIntIndices, obj.StrideSize, obj.StrideOffsetPosition, obj.StrideOffsetTexture, obj.StrideOffsetNormal
	return obj, &uIntIndices
}
