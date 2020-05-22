package utils

import (
	"fmt"
	"os"
	"runtime"
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

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

// GetWindow GetWindow
func GetWindow(width int, height int) *glfw.Window {
	// Lock this calling goroutine to its current operating system thread.
	runtime.LockOSThread()

	// Initializes the GLFW library
	if err := glfw.Init(); err != nil {
		panic(fmt.Errorf("I could not initialize glfw: %v", err))
	}

	// I decided that non-positive width or height meant Fullscreen. Cool eh?
	if width == 0 || height == 0 {
		width, height = FullScreen()
	}
	fmt.Println("w,h", width, height)
	win := CreateWindow(os.Args[0], width, height)

	return win
}

// //Camera : position, rotation etc
// type Camera struct {
// 	Position      mgl32.Vec3
// 	Forward       mgl32.Vec3
// 	Right         mgl32.Vec3
// 	Projection    mgl32.Mat4
// 	View          mgl32.Mat4
// 	Up            mgl32.Vec3
// 	Aspect        float32 // width/height
// 	Fovy          float32
// 	Near          float32
// 	Far           float32
// 	Paused        bool
// 	StartLookAt   mgl32.Vec3
// 	StartPosition mgl32.Vec3
// }

// //Cam Cam
// func Cam() *Camera {

// 	return &Camera{
// 		Position:      mgl32.Vec3{0, 0, 4},
// 		Forward:       mgl32.Vec3{0, 0, -1},
// 		Right:         mgl32.Vec3{1, 0, 0},
// 		Projection:    mgl32.Ident4(),
// 		View:          mgl32.Ident4(),
// 		Up:            mgl32.Vec3{0, 1, 0},
// 		Aspect:        1, // width/height
// 		Fovy:          math.Pi / 4,
// 		Near:          0.01,
// 		Far:           1000,
// 		Paused:        false,
// 		StartLookAt:   mgl32.Vec3{0, 0, 0},
// 		StartPosition: mgl32.Vec3{15, 15, -60},
// 	}
// }
