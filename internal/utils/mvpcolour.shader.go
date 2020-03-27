package utils

//MVPColourVertShader MVPColourVertShader
var MVPColourVertShader = `
#version 430 core

layout (location = 0) in vec4 position;
layout (location = 1) in vec3 colour;

uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;

out vec3 vColour;

void main()
{
	vColour = colour;
	gl_Position = projection * view * model * position;
}
` + "\x00"

//MVPColourFragShader MVPColourFragShader
var MVPColourFragShader = `
#version 430 core

in vec3 vColour;
out vec4 colour;

void main() {
	colour = vec4(vColour, 1.0f);
}
` + "\x00"
