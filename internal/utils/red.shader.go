package utils

//VertexShader VertexShader
var VertexShader = `
#version 430 core

layout (location = 0) in vec4 position;

void main()
{
	gl_Position = position;
}
` + "\x00"

//FragmentShader FragmentShader
var FragmentShader = `
#version 430 core

out vec4 color;

vec4 red = vec4(0, 1, 1, 1.0);

void main() {
	color = red;
}
` + "\x00"


