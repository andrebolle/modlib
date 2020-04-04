package utils

import (
	"fmt"
	"os"
	"runtime"

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

	win, err := glfw.CreateWindow(width, height, title, nil, nil)

	if err != nil {
		panic(fmt.Errorf("Could not create OpenGL renderer: %v", err))
	}

	win.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}

	return win
}

// FullScreen FullScreen
func FullScreen() (*glfw.Window, *Camera) {
	// Camera and Sceen choices (width, height, aspect ratio)
	monitor := glfw.GetPrimaryMonitor()
	vidMode := monitor.GetVideoMode()

	// Set the camera aspect to the screen aspect
	cam := Cam()
	cam.Aspect = float32(vidMode.Width) / float32(vidMode.Height)

	// Print some info
	fmt.Println(vidMode.Width, "x ", vidMode.Height)
	fmt.Println("cam.Aspect", cam.Aspect)

	// Create a window
	win := CreateWindow(os.Args[0], vidMode.Width, vidMode.Height)
	return win, cam
}

// GetWindowAndCamera Boilerplate. LockOSThread, GLFW, Window & Callbacks
func GetWindowAndCamera() (*glfw.Window, *Camera) {
	// Lock this calling goroutine to its current operating system thread.
	runtime.LockOSThread()

	// Initializes the GLFW library
	if err := glfw.Init(); err != nil {
		panic(fmt.Errorf("I could not initialize glfw: %v", err))
	}

	// Create Window and OpenGL context
	win, cam := FullScreen()

	// Print any useful info/help.
	fmt.Println(GraphicsCardName())
	fmt.Println("Use W,A,S,D to move forward, left, backward and right.")
	fmt.Println("Use E,C to move skyward and earthward.")

	// Set callbacks for user input
	SetWASDCallback(win, cam)
	SetPitchYawCallback(win, cam)

	return win, cam
}
