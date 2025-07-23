package main

import (
	"fmt"
	"log"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/javanhut/BifrostEngine/m/v2/renderer/opengl"
	"github.com/javanhut/BifrostEngine/m/v2/renderer/window"
)

func main() {
	win, err := window.New(800, 600, "Simple Perspective Test")
	if err != nil {
		log.Fatal(err)
	}
	defer win.Destroy()

	if err := gl.Init(); err != nil {
		log.Fatal(err)
	}

	// Create shader with hardcoded perspective matrix
	vertexShader := `
#version 410 core
layout (location = 0) in vec3 aPos;
layout (location = 1) in vec3 aColor;

out vec3 vertexColor;

void main() {
    // Simple perspective projection matrix (hardcoded)
    float fov = 0.785398; // 45 degrees in radians
    float aspect = 800.0/600.0;
    float near = 1.0;
    float far = 100.0;
    
    float tanHalfFov = tan(fov / 2.0);
    
    mat4 projection = mat4(
        1.0/(aspect * tanHalfFov), 0, 0, 0,
        0, 1.0/tanHalfFov, 0, 0,
        0, 0, -(far + near)/(far - near), -(2.0 * far * near)/(far - near),
        0, 0, -1, 0
    );
    
    // Simple view matrix (camera at z=3 looking at origin)
    mat4 view = mat4(
        1, 0, 0, 0,
        0, 1, 0, 0,
        0, 0, 1, -3,
        0, 0, 0, 1
    );
    
    gl_Position = projection * view * vec4(aPos, 1.0);
    vertexColor = aColor;
}
` + "\x00"

	shader, err := opengl.NewShader(vertexShader, opengl.DefaultFragmentShader)
	if err != nil {
		log.Fatal("Shader error:", err)
	}
	defer shader.Delete()

	// Triangle at origin
	vertices := []float32{
		-0.5, -0.5, 0.0, 1.0, 0.0, 0.0,
		0.5, -0.5, 0.0, 0.0, 1.0, 0.0,
		0.0, 0.5, 0.0, 0.0, 0.0, 1.0,
	}
	triangle := opengl.NewMesh(vertices)
	defer triangle.Delete()

	fmt.Println("Rendering with hardcoded matrices...")

	for !win.ShouldClose() {
		gl.ClearColor(0.2, 0.2, 0.2, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		shader.Use()
		triangle.Draw()

		win.SwapBuffers()
		win.PollEvents()
	}
}