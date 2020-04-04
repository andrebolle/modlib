package main

// Run this ...
// go get github.com/g3n/engine@master

import (
	"fmt"
	"math"
	"time"

	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/util/helper"
	"github.com/g3n/engine/window"
	"github.com/purelazy/modlib/cmd/lsysg3n/lsys"
)

func main() {

	var hilbert = lsys.L3D{
		Seed:  "X",
		Angle: math.Pi / 2,
		Rules: map[rune]string{
			'X': "^<XF^<XFX-F^>>XFX&F+>>XFX-F>X->",
			'F': "F",
			'-': "-",
			'+': "+",
			'^': "^",
			'&': "&",
			'<': "<",
			'>': ">",
		}}

	// Generate the L-system string
	lstring := lsys.GenLString3D(hilbert, 3)
	//fmt.Println(lstring)

	// Create first frame
	floatArray, coordCount := lsys.Lsystem3D(lstring, hilbert.Angle)
	points := coordCount / 3
	fmt.Println(floatArray[0], points)

	// Convert coordinates to slice of math32.Vector3s
	vector3s := make([]math32.Vector3, 0, 50000)
	for p := 0; p < coordCount; p += 3 {
		vector3s = append(vector3s, math32.Vector3{X: floatArray[p], Y: floatArray[p+1], Z: floatArray[p+2]})
	}

	// --------------------------------

	// App returns the Application singleton, creating it the first time.
	a := app.App()

	// NewNode returns a pointer to a new Node.
	scene := core.NewNode()

	// Set the scene to be managed by the gui manager
	gui.Manager().Set(scene)

	// Create perspective camera
	cam := camera.New(1)
	cam.SetPosition(0, 0, 3)
	scene.Add(cam)

	// Set up orbit control for the camera
	camera.NewOrbitControl(cam)

	// Set up callback to update viewport and camera aspect ratio when the window is resized
	onResize := func(evname string, ev interface{}) {
		// Get framebuffer size and update viewport accordingly
		width, height := a.GetSize()
		a.Gls().Viewport(0, 0, int32(width), int32(height))
		// Update the camera's aspect ratio
		cam.SetAspect(float32(width) / float32(height))
	}
	a.Subscribe(window.OnWindowSize, onResize)
	onResize("", nil)

	// Create a blue torus and add it to the scene
	geom := geometry.NewTube(vector3s, 0.1, 20, true)
	//geom := megeometry.MyNewTube(vector3s, 0.1, 20, true)
	//geom := geometry.NewTube([]math32.Vector3{{X: 0, Y: 0, Z: 0}, {X: 0.5, Y: 0, Z: 0}, {X: 0.5, Y: 0.5, Z: 0}}, 0.2, 12, false)
	// geom := NewTube(path []math32.Vector3, radius float32, radialSegments int, close bool)
	//geom := geometry.NewTorus(1, .4, 12, 32, math32.Pi*2)
	mat := material.NewStandard(math32.NewColor("Grey"))
	mesh := graphic.NewMesh(geom, mat)
	scene.Add(mesh)

	// Create and add a button to the scene
	btn := gui.NewButton("Make Red")
	btn.SetPosition(100, 40)
	btn.SetSize(40, 40)
	btn.Subscribe(gui.OnClick, func(name string, ev interface{}) {
		mat.SetColor(math32.NewColor("DarkRed"))
	})
	scene.Add(btn)

	// Create and add lights to the scene
	scene.Add(light.NewAmbient(&math32.Color{R: 1.0, G: 1.0, B: 1.0}, 0.2))
	pointLight := light.NewPoint(&math32.Color{R: 1, G: 1, B: 1}, 5.0)
	pointLight.SetPosition(3, 6, 3)
	scene.Add(pointLight)

	// Create and add an axis helper to the scene
	scene.Add(helper.NewAxes(0.5))

	// Set background color to gray
	a.Gls().ClearColor(0.5, 0.5, 0.5, 1.0)

	// Run the application
	a.Run(func(renderer *renderer.Renderer, deltaTime time.Duration) {
		//a.Gls().FrontFace(gls.CW)
		a.Gls().Disable(gls.CULL_FACE)
		//a.Gls().CullFace(gls.FRONT)
		a.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
		//a.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
		renderer.Render(scene, cam)
	})
}
