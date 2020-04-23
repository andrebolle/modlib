#version 460

out vec4 FragColor;


in vec3 fPos;
in vec2 fUV;  
in vec3 fNormal;  

uniform sampler2D uTex;
uniform vec3 uLightPos; 
uniform vec3 uViewPos; 
uniform vec3 uLightColor;
//uniform vec3 objectColor;

void main()
{
    vec3 objectColor = texture(uTex, fUV).xyz;

    // ambient
    float ambientStrength = 0.1;
    vec3 ambient = ambientStrength * uLightColor;
  	
    // diffuse 
    vec3 norm = normalize(fNormal);
    vec3 lightDir = normalize(uLightPos - fPos);
    float diff = max(dot(norm, lightDir), 0.0);
    vec3 diffuse = diff * uLightColor;
    
    // specular
    float specularStrength = 0.5;
    vec3 viewDir = normalize(uViewPos - fPos);
    vec3 reflectDir = reflect(-lightDir, norm);  
    float spec = pow(max(dot(viewDir, reflectDir), 0.0), 32);
    vec3 specular = specularStrength * spec * uLightColor;  
        
    vec3 result = (ambient + diffuse + specular) * objectColor;
    FragColor = vec4(result, 1.0);
} 