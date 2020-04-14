#version 460

in vec3 aPos;
in vec2 aUV;
in vec3 aNormal;

out vec3 fPos;
out vec2 fUV;
out vec3 fNormal;


uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;

void main()
{
    fUV = aUV;
    fPos = vec3(model * vec4(aPos, 1.0));
    fNormal = mat3(transpose(inverse(model))) * aNormal;  
    gl_Position = projection * view * vec4(fPos, 1.0);
}