package main

import (	
	"github.com/go-gl/gl/v4.6-core/gl"
"github.com/purelazy/modlib/internal/utils"
)

func shadowPass(width, height int32, fbo uint32, cam *utils.Camera) {
	// 1. first render to depth map
gl.Viewport(0, 0, width, height)
gl.BindFramebuffer(gl.FRAMEBUFFER, fbo)
gl.Clear(gl.DEPTH_BUFFER_BIT)
//ConfigureShaderAndMatrices()
// RenderScene()
// gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
// // 2. then render scene as normal with shadow mapping (using depth map)
// gl.Viewport(0, 0, SCR_WIDTH, SCR_HEIGHT)
// gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
// ConfigureShaderAndMatrices()
// gl.BindTexture(gl.TEXTURE_2D, depthMap)
// RenderScene()
}