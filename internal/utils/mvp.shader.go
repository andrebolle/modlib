package utils

//MVPVertShader MVPVertShader
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

out vec4 colour;

void main() {
	colour = vec4(1, 1, 1, 1.0);
}
` + "\x00"
