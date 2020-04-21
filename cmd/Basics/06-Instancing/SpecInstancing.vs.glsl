#version 460

uniform mat4 model, view, projection; uniform float uTime;

in vec3 aPos, aNormal; in vec2 aUV;
out vec3 fPos, fNormal; out vec2 fUV;

void main()
{
    vec3 instancePos = aPos + vec3(gl_InstanceID % 10 + fract(uTime), gl_InstanceID / 10 % 10, gl_InstanceID / 100) * 10.0;

    fUV = aUV;
    fPos = vec3(model * vec4(instancePos, 1.0));
    fNormal = mat3(transpose(inverse(model))) * aNormal;  
    gl_Position = projection * view * vec4(fPos, 1.0);
}