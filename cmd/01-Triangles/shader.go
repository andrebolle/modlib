package main

//BasicVS BasicVS
var BasicVS = `
#version 430 core

// If you change names here, change the name in GetAttribLocation calls.
layout (location = 0) in vec4 pos;

void main()
{
	gl_Position = pos;
}
` + "\x00"

//BasicFS BasicFS
var BasicFS = `
#version 430 core

out vec4 color;

void main() {
	color = vec4(1, 0, 1, 1.0);
}
` + "\x00"
