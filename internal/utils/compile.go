package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-gl/gl/v4.6-core/gl"
)

//Compile Compile a shader
func Compile(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)
	if shader == 0 {
		panic("gl.CreateShader returned zero")
	}

	csources, free := gl.Strs(source)
	// OpenGL copies the shader source code strings when glShaderSource is called,
	// so an application may free its copy of the source code strings immediately
	// after the function returns.
	// If length is NULL, each string is assumed to be null terminated.
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	// Check for errors
	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		// How many bytes to allocate
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))

		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))
		fmt.Println("CompileShader log")
		panic(fmt.Errorf("Failed to compile %v: %v", source, log))
	}

	return shader, nil
}

// NewProgram Create a Vertex --> Fragment shader
func NewProgram(vertexShaderSource, fragmentShaderSource string) uint32 {
	vertexShader, err := Compile(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		fmt.Println("Vertex shader did not compile")
		fmt.Println(err)
		panic(err)
	}
	defer gl.DeleteShader(vertexShader)

	fragmentShader, err := Compile(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		fmt.Println("Fragment shader did not compile")
		fmt.Println(err)
		panic(err)
	}
	defer gl.DeleteShader(fragmentShader)

	program := gl.CreateProgram()

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var linkedOK int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &linkedOK)
	if linkedOK == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		gl.DeleteProgram(program)

		fmt.Errorf("Failed to link program: %v", log)
		os.Exit(-1)
	}

	return program
}

// SetUniformMat4 SetUniformMat4
func SetUniformMat4(program uint32, name string, value *float32) (location int32) {
	location = gl.GetUniformLocation(program, gl.Str(name+"\x00"))
	gl.UniformMatrix4fv(location, 1, false, value)
	return
}

// SetUniformVec3 SetUniformVec3
func SetUniformVec3(program uint32, name string, value *float32) (location int32) {
	location = gl.GetUniformLocation(program, gl.Str(name+"\x00"))
	gl.Uniform3fv(location, 1, value)
	return
}
