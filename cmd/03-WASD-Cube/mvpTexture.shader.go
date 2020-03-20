package main

//VertexShader MVP and Texture Shader
var VertexShader = `
#version 330

uniform mat4 P;
uniform mat4 V;
uniform mat4 M;

in vec3 vert;
in vec2 vertTexCoord;

out vec2 fragTexCoord;

void main() {
    fragTexCoord = vertTexCoord;
    gl_Position = P * V * M * vec4(vert, 1);
}
` + "\x00"

//FragmentShader MVP and Texture Shader
var FragmentShader = `
#version 330

uniform sampler2D tex;

in vec2 fragTexCoord;

out vec4 outputColor;

void main() {
    outputColor = texture(tex, fragTexCoord);
}
` + "\x00"
