#version 330

// Sampler uniform points to a particular texture unit
// Sampling is the process of computing a color from an image texture and texture coordinates.
// Sampler variables must be declared as global uniform variables.
uniform sampler2D tex;

in vec2 fragTexCoord;

out vec4 outputColor;

vec4 lightColor = vec4(1,1,1,1);
float ambientStrength = 0.9;
vec3 lightPos = vec3(3,3,3);

void main() {
    vec4 ambient = lightColor * ambientStrength;
    outputColor = texture(tex, fragTexCoord) * ambient;
}