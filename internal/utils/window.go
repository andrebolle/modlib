package utils

import (
	"fmt"
	"os"
	"runtime"
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
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

// GraphicsCardName GraphicsCardName
func GraphicsCardName() string {
	return gl.GoStr(gl.GetString(gl.RENDERER))
}

//CreateWindow You've guessed it!
func CreateWindow(title string, width, height int) *glfw.Window {
	if !(width != 0 && height != 0) {
		fmt.Println("Width and Height cannot be zero.")
		os.Exit(0)
	}

	if err := glfw.Init(); err != nil {
		panic(fmt.Errorf("Could not initialize GLFW: %v", err))
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 6)
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCompatProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	// Debug context
	glfw.WindowHint(glfw.OpenGLDebugContext, glfw.True)

	win, err := glfw.CreateWindow(width, height, title, nil, nil)

	if err != nil {
		panic(fmt.Errorf("Could not create OpenGL renderer: %v", err))
	}

	win.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}

	var flags int32
	gl.GetIntegerv(gl.CONTEXT_FLAGS, &flags)
	if flags&int32(gl.CONTEXT_FLAG_DEBUG_BIT) != 0 {
		fmt.Println("Debugging is available")

		gl.DebugMessageCallback(debugCb, unsafe.Pointer(nil))
		//gl.DebugMessageControl(gl.DONT_CARE, gl.DONT_CARE, gl.DONT_CARE, 0, nullptr, gl.TRUE)
		//gl.Enable(gl.DEBUG_OUTPUT)
		//gl.Disable(gl.DEBUG_OUTPUT)
		//gl.Enable(gl.DEBUG_OUTPUT_SYNCHRONOUS)
	}

	return win
}

// FullScreen FullScreen
func FullScreen() (int, int) {

	// Camera and Sceen choices (width, height, aspect ratio)
	monitor := glfw.GetPrimaryMonitor()
	vidMode := monitor.GetVideoMode()

	return vidMode.Width, vidMode.Height
}

// GetWindowAndCamera Boilerplate. LockOSThread, GLFW, Window & Callbacks
func GetWindowAndCamera(width int, height int) (*glfw.Window, *Camera) {
	// Lock this calling goroutine to its current operating system thread.
	runtime.LockOSThread()

	// Initializes the GLFW library
	if err := glfw.Init(); err != nil {
		panic(fmt.Errorf("I could not initialize glfw: %v", err))
	}

	if width == 0 || height == 0 {
		width, height = FullScreen()
	}

	// Set the camera aspect to the screen aspect
	cam := Cam()
	cam.Aspect = float32(width) / float32(height)

	// Print some info
	fmt.Println(width, "x ", height)
	fmt.Println("cam.Aspect", cam.Aspect)

	// Create a window
	win := CreateWindow(os.Args[0], width, height)

	// Print any useful info/help.
	fmt.Println(GraphicsCardName())
	fmt.Println("Use W,A,S,D to move forward, left, backward and right.")
	fmt.Println("Use E,C to move skyward and earthward.")

	// Set callbacks for user input
	SetWASDCallback(win, cam)
	SetPitchYawCallback(win, cam)

	return win, cam
}
