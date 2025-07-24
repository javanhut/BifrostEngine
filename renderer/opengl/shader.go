package opengl

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

type Shader struct {
	program uint32
}

const (
	DefaultVertexShader = `
#version 410 core
layout (location = 0) in vec3 aPos;
layout (location = 1) in vec3 aColor;

uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;

out vec3 vertexColor;

void main() {
    gl_Position = projection * view * model * vec4(aPos, 1.0);
    vertexColor = aColor;
}
` + "\x00"

	DefaultFragmentShader = `
#version 410 core
in vec3 vertexColor;
out vec4 FragColor;

void main() {
    FragColor = vec4(vertexColor, 1.0);
}
` + "\x00"

	SimpleVertexShader = `
#version 410 core
layout (location = 0) in vec3 aPos;
layout (location = 1) in vec3 aColor;

out vec3 vertexColor;

void main() {
    gl_Position = vec4(aPos, 1.0);
    vertexColor = aColor;
}
` + "\x00"

	MaterialVertexShader = `
#version 410 core
layout (location = 0) in vec3 aPos;
layout (location = 1) in vec3 aColor;
layout (location = 2) in vec2 aTexCoord;

uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;

out vec3 vertexColor;
out vec2 texCoord;

void main() {
    gl_Position = projection * view * model * vec4(aPos, 1.0);
    vertexColor = aColor;
    texCoord = aTexCoord;
}
` + "\x00"

	MaterialFragmentShader = `
#version 410 core
in vec3 vertexColor;
in vec2 texCoord;
out vec4 FragColor;

uniform bool useTexture;
uniform sampler2D diffuseTexture;
uniform vec3 materialColor;
uniform float alpha;

void main() {
    vec4 finalColor;
    
    if (useTexture) {
        // Sample texture and multiply with vertex color
        vec4 texColor = texture(diffuseTexture, texCoord);
        finalColor = texColor * vec4(vertexColor, 1.0);
    } else {
        // Use material color mixed with vertex color
        finalColor = vec4(materialColor * vertexColor, alpha);
    }
    
    FragColor = finalColor;
}
` + "\x00"

	// Lighting-enabled shaders with Phong lighting model
	LightingVertexShader = `
#version 410 core
layout (location = 0) in vec3 aPos;
layout (location = 1) in vec3 aColor;
layout (location = 2) in vec2 aTexCoord;
layout (location = 3) in vec3 aNormal;

uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;
uniform mat3 normalMatrix; // For transforming normals

out vec3 vertexColor;
out vec2 texCoord;
out vec3 fragPos;    // Fragment position in world space
out vec3 normal;     // Normal in world space

void main() {
    // Transform position to world space
    vec4 worldPos = model * vec4(aPos, 1.0);
    fragPos = worldPos.xyz;
    
    // Transform normal to world space (using normal matrix)
    normal = normalMatrix * aNormal;
    
    gl_Position = projection * view * worldPos;
    vertexColor = aColor;
    texCoord = aTexCoord;
}
` + "\x00"

	LightingFragmentShader = `
#version 410 core
in vec3 vertexColor;
in vec2 texCoord;
in vec3 fragPos;
in vec3 normal;
out vec4 FragColor;

// Material properties
uniform bool useTexture;
uniform sampler2D diffuseTexture;
uniform vec3 materialColor;
uniform float alpha;
uniform float shininess;

// Lighting uniforms
uniform vec3 viewPos; // Camera position
uniform vec3 ambientLight;

// Directional light
uniform bool dirLightEnabled;
uniform vec3 dirLightDirection;
uniform vec3 dirLightColor;
uniform float dirLightIntensity;

// Point lights (max 4 for now)
uniform int numPointLights;
uniform vec3 pointLightPositions[4];
uniform vec3 pointLightColors[4];
uniform float pointLightIntensities[4];
uniform vec3 pointLightAttenuations[4]; // constant, linear, quadratic

vec3 calculateDirectionalLight(vec3 norm, vec3 viewDir, vec3 baseColor) {
    if (!dirLightEnabled) return vec3(0.0);
    
    vec3 lightDir = normalize(-dirLightDirection);
    
    // Diffuse lighting
    float diff = max(dot(norm, lightDir), 0.0);
    vec3 diffuse = diff * dirLightColor * dirLightIntensity;
    
    // Specular lighting (Phong)
    vec3 reflectDir = reflect(-lightDir, norm);
    float spec = pow(max(dot(viewDir, reflectDir), 0.0), shininess);
    vec3 specular = spec * dirLightColor * dirLightIntensity * 0.5; // Reduce specular intensity
    
    return (diffuse + specular) * baseColor;
}

vec3 calculatePointLight(int index, vec3 norm, vec3 viewDir, vec3 baseColor) {
    if (index >= numPointLights) return vec3(0.0);
    
    vec3 lightPos = pointLightPositions[index];
    vec3 lightColor = pointLightColors[index];
    float lightIntensity = pointLightIntensities[index];
    vec3 attenuation = pointLightAttenuations[index];
    
    vec3 lightDir = normalize(lightPos - fragPos);
    float distance = length(lightPos - fragPos);
    
    // Attenuation
    float attenuationFactor = 1.0 / (attenuation.x + attenuation.y * distance + attenuation.z * distance * distance);
    
    // Diffuse lighting
    float diff = max(dot(norm, lightDir), 0.0);
    vec3 diffuse = diff * lightColor * lightIntensity * attenuationFactor;
    
    // Specular lighting
    vec3 reflectDir = reflect(-lightDir, norm);
    float spec = pow(max(dot(viewDir, reflectDir), 0.0), shininess);
    vec3 specular = spec * lightColor * lightIntensity * attenuationFactor * 0.5;
    
    return (diffuse + specular) * baseColor;
}

void main() {
    vec3 norm = normalize(normal);
    vec3 viewDir = normalize(viewPos - fragPos);
    
    // Get base color
    vec3 baseColor;
    if (useTexture) {
        vec4 texColor = texture(diffuseTexture, texCoord);
        baseColor = texColor.rgb * vertexColor;
    } else {
        baseColor = materialColor * vertexColor;
    }
    
    // Start with ambient lighting
    vec3 result = ambientLight * baseColor;
    
    // Add directional light
    result += calculateDirectionalLight(norm, viewDir, baseColor);
    
    // Add point lights
    for (int i = 0; i < numPointLights && i < 4; i++) {
        result += calculatePointLight(i, norm, viewDir, baseColor);
    }
    
    FragColor = vec4(result, alpha);
}
` + "\x00"
)

func NewShader(vertexSource, fragmentSource string) (*Shader, error) {
	vertexShader, err := CompileShader(vertexSource, gl.VERTEX_SHADER)
	if err != nil {
		return nil, err
	}

	fragmentShader, err := CompileShader(fragmentSource, gl.FRAGMENT_SHADER)
	if err != nil {
		gl.DeleteShader(vertexShader)
		return nil, err
	}

	program, err := CreateProgram(vertexShader, fragmentShader)
	if err != nil {
		return nil, err
	}

	return &Shader{program: program}, nil
}

func (s *Shader) Use() {
	gl.UseProgram(s.program)
}

func (s *Shader) SetMatrix4(name string, matrix *float32) {
	location := gl.GetUniformLocation(s.program, gl.Str(name + "\x00"))
	gl.UniformMatrix4fv(location, 1, false, matrix)
}

func (s *Shader) SetBool(name string, value bool) {
	location := gl.GetUniformLocation(s.program, gl.Str(name + "\x00"))
	if value {
		gl.Uniform1i(location, 1)
	} else {
		gl.Uniform1i(location, 0)
	}
}

func (s *Shader) SetInt(name string, value int32) {
	location := gl.GetUniformLocation(s.program, gl.Str(name + "\x00"))
	gl.Uniform1i(location, value)
}

func (s *Shader) SetFloat(name string, value float32) {
	location := gl.GetUniformLocation(s.program, gl.Str(name + "\x00"))
	gl.Uniform1f(location, value)
}

func (s *Shader) SetVec3(name string, x, y, z float32) {
	location := gl.GetUniformLocation(s.program, gl.Str(name + "\x00"))
	gl.Uniform3f(location, x, y, z)
}

func (s *Shader) SetMatrix3(name string, matrix *float32) {
	location := gl.GetUniformLocation(s.program, gl.Str(name + "\x00"))
	gl.UniformMatrix3fv(location, 1, false, matrix)
}

func (s *Shader) Delete() {
	gl.DeleteProgram(s.program)
}

func (s *Shader) GetProgramID() uint32 {
	return s.program
}