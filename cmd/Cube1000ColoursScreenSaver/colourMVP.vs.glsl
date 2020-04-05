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