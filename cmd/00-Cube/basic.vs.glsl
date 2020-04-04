#version 460 core

layout (location = 0) in vec3 pos;

out vec3 vColor;

void main()
{
  vColor = vec3(1,pos.x,pos.x);
  gl_Position = vec4(pos,1.0);
}