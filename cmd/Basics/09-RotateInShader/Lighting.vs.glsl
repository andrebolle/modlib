#version 460

in vec3 aPos;
in vec2 aUV;
in vec3 aNormal;

out vec3 fPos;
out vec2 fUV;
out vec3 fNormal;

uniform float uAngle;
uniform mat4 uModel;
uniform mat4 uView;
uniform mat4 uProjection;

mat4 translate3D(float Tx, float Ty, float Tz) {
	return mat4(1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, Tx, Ty, Tz, 1);
}

mat4 homogRotate3D(float angle, vec3 axis) {
	//x, y, z = axis[0], axis[1], axis[2]
    float x = axis.x;
    float y = axis.y;
    float z = axis.z;
	//s, c := float32(math.Sin(float64(angle))), float32(math.Cos(float64(angle)))
	float s = sin(angle);
	float c = cos(angle);
	float k = 1 - c;

	return mat4(x*x*k + c, x*y*k + z*s, x*z*k - y*s, 0, x*y*k - z*s, y*y*k + c, y*z*k + x*s, 0, x*z*k + y*s, y*z*k - x*s, z*z*k + c, 0, 0, 0, 0, 1);
}

void main()
{
    mat4 rotationMatrix = homogRotate3D(uAngle, vec3(0,0,1));
    fUV = aUV;
    fPos = vec3(uModel * vec4(aPos, 1.0));
    fNormal = mat3(transpose(inverse(uModel))) * aNormal;  
    gl_Position = uProjection * uView * vec4(fPos, 1.0);
}