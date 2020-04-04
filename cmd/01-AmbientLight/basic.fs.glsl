#version 460 core

in vec3 vColor;

out vec3 color;

vec3 lightColor = vec3(1,1,1);
float ambientStrength = 0.3f;

void main() {
   vec3 ambient = lightColor * ambientStrength;
   color = ambient * vColor;
}