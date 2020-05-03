#version 460

in vec3 aPos;
in vec2 aUV;
in vec3 aNormal;

out vec3 fPos;
out vec2 fUV;
out vec3 fNormal;


uniform mat4 uModel;
uniform mat4 uView;
uniform mat4 uProjection;

void main()
{
    fUV = aUV;
    fPos = vec3(uModel * vec4(aPos, 1.0));
    fNormal = mat3(transpose(inverse(uModel))) * aNormal;  
    gl_Position = uProjection * uView * vec4(fPos, 1.0);
}