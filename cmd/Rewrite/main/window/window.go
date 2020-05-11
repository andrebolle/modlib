package window

import "github.com/go-gl/glfw/v3.3/glfw"

// NewWindow initializes glfw and returns a Window to use.
func NewWindow(width, height int) *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4) // OR 2
	glfw.WindowHint(glfw.ContextVersionMinor, 6)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "Hello OpenGL!", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	return window
}

// //CreateWindow You've guessed it!
// func CreateWindow(title string, width, height int) *glfw.Window {
// 	if !(width != 0 && height != 0) {
// 		fmt.Println("Width and Height cannot be zero.")
// 		os.Exit(0)
// 	}

// 	if err := glfw.Init(); err != nil {
// 		panic(fmt.Errorf("Could not initialize GLFW: %v", err))
// 	}

// 	glfw.WindowHint(glfw.ContextVersionMajor, 4)
// 	glfw.WindowHint(glfw.ContextVersionMinor, 6)
// 	glfw.WindowHint(glfw.Resizable, glfw.True)
// 	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCompatProfile)
// 	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

// 	// Debug context
// 	glfw.WindowHint(glfw.OpenGLDebugContext, glfw.True)

// 	win, err := glfw.CreateWindow(width, height, title, nil, nil)

// 	if err != nil {
// 		panic(fmt.Errorf("Could not create OpenGL renderer: %v", err))
// 	}

// 	win.MakeContextCurrent()

// 	if err := gl.Init(); err != nil {
// 		panic(err)
// 	}

// 	var flags int32
// 	gl.GetIntegerv(gl.CONTEXT_FLAGS, &flags)
// 	if flags&int32(gl.CONTEXT_FLAG_DEBUG_BIT) != 0 {
// 		fmt.Println("Debugging is available")

// 		gl.DebugMessageCallback(debugCb, unsafe.Pointer(nil))
// 		//gl.DebugMessageControl(gl.DONT_CARE, gl.DONT_CARE, gl.DONT_CARE, 0, nullptr, gl.TRUE)
// 		//gl.Enable(gl.DEBUG_OUTPUT)
// 		//gl.Disable(gl.DEBUG_OUTPUT)
// 		//gl.Enable(gl.DEBUG_OUTPUT_SYNCHRONOUS)
// 	}

// 	return win
// }
