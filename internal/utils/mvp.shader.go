package utils

//MvpVertShader MvpVertShader
var MVPVertShader = `
#version 430 core

layout (location = 0) in vec4 position;

uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;

void main()
{
	gl_Position = projection * view * model * position;
}
` + "\x00"

//MVPFragShader MVPFragShader
var MVPFragShader = `
#version 430 core

out vec4 color;

vec4 red = vec4(0.2, 0.0, 0.0, 1.0);

void main() {
	color = red;
}
` + "\x00"
