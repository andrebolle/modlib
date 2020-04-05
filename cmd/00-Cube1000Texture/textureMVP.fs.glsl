#version 330

uniform sampler2D tex;

in vec2 fragTexCoord;

out vec4 outputColor;

vec4 lightColor = vec4(1,1,1,1);
float ambientStrength = 0.3;
vec3 lightPos = vec3(3,3,3);

void main() {
    vec4 ambient = lightColor * ambientStrength;
    outputColor = texture(tex, fragTexCoord) * ambient;
}