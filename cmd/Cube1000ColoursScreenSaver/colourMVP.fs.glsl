#version 430 core

in vec3 vColour;
out vec4 colour;

void main() {
	colour = vec4(vColour, 1.0f);
}