package main

import (
	"fmt"
	"log"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/javanhut/BifrostEngine/m/v2/renderer/opengl"
	"github.com/javanhut/BifrostEngine/m/v2/renderer/window"
)

func main() {
	// Create window
	win, err := window.New(800, 600, "Simple Triangle Test")
	if err != nil {
		log.Fatal(err)
	}
	defer win.Destroy()

	// Initialize OpenGL
	if err := gl.Init(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("OpenGL Version: %s\n", gl.GoStr(gl.GetString(gl.VERSION)))

	// Create simple shader
	shader, err := opengl.NewShader(opengl.SimpleVertexShader, opengl.DefaultFragmentShader)
	if err != nil {
		log.Fatal("Shader error:", err)
	}
	defer shader.Delete()

	// Create triangle
	vertices := []float32{
		// Position       Color
		-0.5, -0.5, 0.0,  1.0, 0.0, 0.0,  // Bottom left - Red
		 0.5, -0.5, 0.0,  0.0, 1.0, 0.0,  // Bottom right - Green
		 0.0,  0.5, 0.0,  0.0, 0.0, 1.0,  // Top - Blue
	}
	triangle := opengl.NewMesh(vertices)
	defer triangle.Delete()

	// Check for OpenGL errors
	if err := gl.GetError(); err != 0 {
		fmt.Printf("OpenGL Error after setup: %d\n", err)
	}

	// Debug: Check if shader is valid
	fmt.Printf("Shader program ID: %d\n", shader.GetProgramID())

	fmt.Println("Setup complete, entering render loop...")
	frameCount := 0

	// Render loop
	for !win.ShouldClose() {
		// Clear
		gl.ClearColor(0.2, 0.2, 0.2, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Draw
		shader.Use()
		triangle.Draw()

		// Check for errors
		if err := gl.GetError(); err != 0 {
			fmt.Printf("OpenGL Error in draw: %d\n", err)
		}

		// Swap buffers
		win.SwapBuffers()
		win.PollEvents()
		
		frameCount++
		if frameCount == 1 {
			fmt.Println("First frame completed")
		}
	}
}