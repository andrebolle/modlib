package utils

// https://github.com/SonarSystems/Modern-OpenGL-Tutorials/blob/master/%5BADVANCED%20OPENGL%5D/%5B17%5D%20Cubemap:Skybox/main.cpp

import (
	"fmt"
	"image"
	"image/draw"
	"os"

	"github.com/go-gl/gl/v4.6-core/gl"
)

// Texture Texture
type Texture struct {
	ID      uint32
	Width   int32
	Height  int32
	Texture []uint8
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
func NewTexture(tex *Texture, file string) {

	rgba := LoadRGBA(file)
	tex.Texture = rgba.Pix
	tex.Width = int32(rgba.Rect.Size().X)
	tex.Height = int32(rgba.Rect.Size().Y)
	tex.Texture = rgba.Pix

	gl.GenTextures(1, &tex.ID)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, tex.ID)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

}
