#version 330 core
out vec4 FragColor;

//in vec2 TexCoords;
in vec3 TexCoords;

//uniform sampler2D uTex;
uniform samplerCube uTex;

void main()
{    
    FragColor = texture(uTex, TexCoords);
    //FragColor = vec4(1);
}