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

func (s *Shader) Delete() {
	gl.DeleteProgram(s.program)
}

func (s *Shader) GetProgramID() uint32 {
	return s.program
}