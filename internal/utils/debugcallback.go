package utils

import (
	"fmt"
	"os"
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
)

// OpenGL debugger callback
func debugCb(
	source uint32,
	gltype uint32,
	id uint32,
	severity uint32,
	length int32,
	message string,
	userParam unsafe.Pointer) {

	switch source {
	case gl.DEBUG_SOURCE_API:
		fmt.Println("Source: API")
	case gl.DEBUG_SOURCE_WINDOW_SYSTEM:
		fmt.Println("Source: Window System")
	case gl.DEBUG_SOURCE_SHADER_COMPILER:
		fmt.Println("Source: Shader Compiler")
	case gl.DEBUG_SOURCE_THIRD_PARTY:
		fmt.Println("Source: Third Party")
	case gl.DEBUG_SOURCE_APPLICATION:
		fmt.Println("Source: Application")
	case gl.DEBUG_SOURCE_OTHER:
		fmt.Println("Source: Other")
	}

	switch gltype {
	case gl.DEBUG_TYPE_ERROR:
		fmt.Println("Type: Error")
	case gl.DEBUG_TYPE_DEPRECATED_BEHAVIOR:
		fmt.Println("Type: Deprecated Behaviour")
	case gl.DEBUG_TYPE_UNDEFINED_BEHAVIOR:
		fmt.Println("Type: Undefined Behaviour")
	case gl.DEBUG_TYPE_PORTABILITY:
		fmt.Println("Type: Portability")
	case gl.DEBUG_TYPE_PERFORMANCE:
		fmt.Println("Type: Performance")
	case gl.DEBUG_TYPE_MARKER:
		fmt.Println("Type: Marker")
	case gl.DEBUG_TYPE_PUSH_GROUP:
		fmt.Println("Type: Push Group")
	case gl.DEBUG_TYPE_POP_GROUP:
		fmt.Println("Type: Pop Group")
	case gl.DEBUG_TYPE_OTHER:
		fmt.Println("Type: Other")
	}

	switch severity {
	case gl.DEBUG_SEVERITY_HIGH:
		fmt.Println("Severity: high")
	case gl.DEBUG_SEVERITY_MEDIUM:
		fmt.Println("Severity: medium")
	case gl.DEBUG_SEVERITY_LOW:
		fmt.Println("Severity: low")
	case gl.DEBUG_SEVERITY_NOTIFICATION:
		fmt.Println("Severity: notification")
	}

	//msg := fmt.Sprintf("[GL_DEBUG] source %d gltype %d id %d severity %d length %d: %s\n", source, gltype, id, severity, length, message)
	msg := fmt.Sprintf("[GL_DEBUG] id %d length %d: %s\n", id, length, message)
	if severity == gl.DEBUG_SEVERITY_HIGH {
		panic(msg)
	}
	fmt.Fprintln(os.Stderr, msg)
}
