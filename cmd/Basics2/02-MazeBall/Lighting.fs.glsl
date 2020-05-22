#version 460

out vec4 FragColor;

in vec3 fPos;
in vec2 fUV;  
in vec3 fNormal;  

uniform sampler2D uTex;

uniform vec3 uLightPos; 
uniform vec3 uViewPos; 
uniform vec3 uLightColor;

void main()
{
    vec3 objectColor = texture(uTex, fUV).xyz;

    // Ambient light (Direction independent)
    float lightPower = 0.4;
    vec3 ambient = uLightColor * lightPower * 0;
  	
    // When light hits an object, an important fraction is reflected in all directions. This is the “diffuse component”.
    // This surface is illuminated differently according to the angle at which the light arrives.
    float distance = length( uLightPos - fPos );
    vec3 normal= normalize(fNormal);
    vec3 lightDir = normalize(uLightPos - fPos);
    float angleBetweenNormalAndLight = dot(normal, lightDir);
    // If the light is behind the triangle, normal and lightDir will on opposite sides of the surface, so dot(normal, lightDir) will be negative.
    float diffuseStrength = max(angleBetweenNormalAndLight, 0.0);
    // The 1000 is a hack for now.
    vec3 diffuse = diffuseStrength * uLightColor * lightPower * 1000 / (distance * distance);
    
    // specular
    // float specularStrength = 0.5;
    // vec3 viewDir = normalize(uViewPos - fPos);
    // vec3 reflectDir = reflect(-lightDir, normal);  
    // float spec = pow(max(dot(viewDir, reflectDir), 0.0), 32);
    // vec3 specular = specularStrength * spec * uLightColor;  
        
    // vec3 result = (ambient + diffuse + specular) * objectColor;
    vec3 result = (ambient + diffuse) * objectColor;

    FragColor = vec4(result, 1.0);
} 