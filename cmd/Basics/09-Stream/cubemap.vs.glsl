#version 330 core
layout (location = 0) in vec3 aPos;
//layout (location = 1) in vec2 aTexCoords;

//out vec2 TexCoords;
out vec3 TexCoords;

//uniform mat4 uModel;
uniform mat4 uView;
uniform mat4 uProjection;

void main()
{
    //TexCoords = aTexCoords;    
    TexCoords = aPos;    
    //gl_Position = uProjection * uView * uModel * vec4(aPos, 1.0);
    gl_Position = uProjection * uView * vec4(aPos, 1.0);
}