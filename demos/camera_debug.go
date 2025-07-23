package main

import (
	"fmt"
	"log"

	bmath "github.com/javanhut/BifrostEngine/m/v2/math"
	"github.com/javanhut/BifrostEngine/m/v2/renderer/core"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func main() {
	renderer, err := core.New(800, 600, "Camera Debug")
	if err != nil {
		log.Fatal(err)
	}
	defer renderer.Cleanup()

	fmt.Println("Controls:")
	fmt.Println("1 - No camera (identity matrices)")
	fmt.Println("2 - Orthographic projection only")
	fmt.Println("3 - View transformation only (with ortho)")  
	fmt.Println("4 - Full 3D camera (perspective + view)")
	fmt.Println("5 - Move camera back to z=5")
	
	mode := 1
	camera := renderer.GetCamera()
	
	// Set initial camera position at origin looking down -z
	camera.SetPosition(bmath.NewVector3(0, 0, 0))
	camera.SetTarget(bmath.NewVector3(0, 0, -1))
	
	window := renderer.GetWindow()
	
	// Key callback
	window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if action == glfw.Press {
			switch key {
			case glfw.Key1:
				mode = 1
				fmt.Println("Mode 1: Identity matrices")
			case glfw.Key2:
				mode = 2
				fmt.Println("Mode 2: Orthographic projection")
			case glfw.Key3:
				mode = 3
				fmt.Println("Mode 3: View transformation")
			case glfw.Key4:
				mode = 4
				fmt.Println("Mode 4: Full 3D camera")
			case glfw.Key5:
				camera.SetPosition(bmath.NewVector3(0, 0, 5))
				camera.SetTarget(bmath.NewVector3(0, 0, 0))
				mode = 4 // Switch to full camera mode to see the effect
				fmt.Println("Camera moved to z=5, switched to full camera mode")
			}
		}
	})
	
	for !renderer.ShouldClose() {
		renderer.BeginFrame()
		
		switch mode {
		case 1:
			renderer.DrawTriangleNoCamera()
		case 2:
			renderer.DrawTriangleProjectionOnly()
		case 3:
			renderer.DrawTriangleViewOnly()
		case 4:
			renderer.DrawTriangle()
		}
		
		renderer.EndFrame()
	}
}