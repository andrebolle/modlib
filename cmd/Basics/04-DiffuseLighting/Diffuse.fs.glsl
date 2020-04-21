#version 460

// Sampler uniform points to a particular texture unit
// Sampling is the process of computing a color from an image texture and texture coordinates.
// Sampler variables must be declared as global uniform variables.
uniform sampler2D tex;

in vec2 fragTexCoord;
in vec3 Normal;
in vec3 FragPos;  

out vec4 FragColor;

vec3 lightColor = vec3(1,1,1);
float ambientStrength = 0.1;
vec3 lightPos = vec3(3,3,3);

void main() {
    // The light floods the room with a certain strength
    vec3 ambient = ambientStrength * lightColor;

    
    // What do we need to calculate diffuse lighting?
    // 1. Normal vector: a vector that is perpendicular to the vertex' surface.
    vec3 norm = normalize(Normal);

    // 2. The light's direction from the frag's position.
    vec3 lightToFragDir = normalize(lightPos - FragPos); 
    float diffStrength = max(dot(norm, lightToFragDir), 0.0);
    vec3 diffuse = diffStrength * lightColor;
    vec3 result =  (ambient + diffuse) * texture(tex, fragTexCoord).xyz;
    FragColor = vec4(result, 1.0);
}