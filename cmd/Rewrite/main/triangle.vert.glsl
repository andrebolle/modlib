#version 430 core

// layout(qualifier1​, qualifier2​ = value, ...) variable definition
//layout(location = 0) in vec4 aPosition;
layout(location = 0) in vec2 aPosition;

void main()
{
    gl_Position = vec4(aPosition, 0, 1);
}