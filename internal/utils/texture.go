package utils

// https://github.com/SonarSystems/Modern-OpenGL-Tutorials/blob/master/%5BADVANCED%20OPENGL%5D/%5B17%5D%20Cubemap:Skybox/main.cpp

import (
	"fmt"
	"image"
	"image/draw"
	"unsafe"

	// Required in order to use jpeg files.
	_ "image/jpeg"
	"os"

	"github.com/go-gl/gl/v4.6-core/gl"
)

// Texture Texture
type Texture struct {
	ID     uint32
	Width  int32
	Height int32
	RGBA   *image.RGBA
}

// LoadRGBA LoadRGBA
func LoadRGBA(file string) *image.RGBA {
	// Open the file
	imgFile, err := os.Open(file)
	if err != nil {
		panic(err)
	}

	// Decode the image
	img, formatName, err := image.Decode(imgFile)
	if err != nil {
		panic(err)
	}

	fmt.Println("Image Filename:", file)
	fmt.Println("Image format is", formatName)
	fmt.Println("Image Bounds is", img.Bounds())

	// NewRGBA returns a new RGBA image with the given bounds.
	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		fmt.Errorf("unsupported stride")
		os.Exit(-1)
	}

	// Draw calls DrawMask with a nil mask.
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)
	return rgba
}

// NewTexture NewTexture
func NewTexture(file string) Texture {
	var texture Texture
	texture.RGBA = LoadRGBA(file)
	texture.Width = int32(texture.RGBA.Rect.Size().X)
	texture.Height = int32(texture.RGBA.Rect.Size().Y)

	gl.GenTextures(1, &texture.ID)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture.ID)

	// Set basic filter and wrap values
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB, texture.Width,
		texture.Height, 0, gl.RGBA, gl.UNSIGNED_BYTE, unsafe.Pointer(&texture.RGBA.Pix[0]))
	return texture
}
