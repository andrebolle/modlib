package utils

import (
	"fmt"

	"github.com/udhos/gwob"
)

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
func OJBLoader(filename string) (*[]float32, *[]uint32, int, int, int, int) {

	options := &gwob.ObjParserOptions{} // parser options

	o, errObj := gwob.NewObjFromFile(filename, options) // parse/load OBJ

	//fmt.Println(o, errObj)
	if errObj != nil {
		panic(errObj)
	}

	fmt.Println("Indices 			", len(o.Indices))
	fmt.Println("Coords   			", len(o.Coord))
	fmt.Println("Mtllib  			", o.Mtllib)
	fmt.Println("Groups  			", len(o.Groups))
	fmt.Println("BigIndexFound  	", o.BigIndexFound)
	fmt.Println("TextCoordFound 	", o.TextCoordFound)
	fmt.Println("NormCoordFound 	", o.NormCoordFound)
	fmt.Println("StrideSize         ", o.StrideSize)
	fmt.Println("StrideOffsetPosition", o.StrideOffsetPosition)
	fmt.Println("StrideOffsetTexture", o.StrideOffsetTexture)
	fmt.Println("StrideOffsetNormal", o.StrideOffsetNormal)

	// fmt.Print(o.Indices)

	uIntIndices := make([]uint32, 0)
	// Convert index "ints" to "uints"
	for i := range o.Indices {
		uIntIndices = append(uIntIndices, uint32(o.Indices[i]))
	}

	//return &o.Coord, &o.Indices, o.StrideSize, o.StrideOffsetPosition, o.StrideOffsetTexture, o.StrideOffsetNormal
	return &o.Coord, &uIntIndices, o.StrideSize, o.StrideOffsetPosition, o.StrideOffsetTexture, o.StrideOffsetNormal

}
