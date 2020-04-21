#version 460

in vec3 vert;
in vec2 vertTexCoord;
in vec3 aNormal;

uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;

out vec2 fragTexCoord;
out vec3 Normal;
out vec3 FragPos; 

void main() {
    fragTexCoord = vertTexCoord;
    Normal = aNormal;
    FragPos = vec3(model * vec4(vert, 1.0));
    gl_Position = projection * view * model * vec4(vert, 1);
}