package utils

import (
	"fmt"
	"image"
	"image/draw"
	"os"

	"github.com/go-gl/gl/v4.6-core/gl"
)

// LoadTexture LoadTexture
func LoadTexture(file string) (uint32, error) {
	imgFile, err := os.Open(file)
	if err != nil {
		return 0, fmt.Errorf("texture %q not found on disk: %v", file, err)
	}
	img, formatName, err := image.Decode(imgFile)
	if err != nil {
		return 0, err
	}

	fmt.Println("Image format is", formatName)
	fmt.Println("Image Bounds is", img.Bounds())

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return 0, fmt.Errorf("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	// A texture object is a data structure that contains the color data
	// for an image texture, and possibly for a set of mipmaps for the texture,
	// as well as the values of texture properties such as the minification
	// and magnification filters and the texture repeat mode.

	// Create a Texture Object
	var textureObject uint32
	gl.GenTextures(1, &textureObject)
	// Initially, texture unit number 0 is active.
	gl.ActiveTexture(gl.TEXTURE0)
	// Binding to a target (2D) allows us to set the attributes of textureObject
	gl.BindTexture(gl.TEXTURE_2D, textureObject)
	// Define filter, wrap and image
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	return textureObject, nil
}
