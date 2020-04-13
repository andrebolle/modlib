#version 460

in vec3 vert;
in vec2 vertTexCoord;
in vec3 aNormal;

out vec3 FragPos;
out vec3 Normal;
out vec2 fragTexCoord;

uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;

void main()
{
    fragTexCoord = vertTexCoord;
    FragPos = vec3(model * vec4(vert, 1.0));
    Normal = mat3(transpose(inverse(model))) * aNormal;  
    gl_Position = projection * view * vec4(FragPos, 1.0);
}