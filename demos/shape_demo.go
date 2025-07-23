package main

import (
	"fmt"
	"log"
	"math"

	bmath "github.com/javanhut/BifrostEngine/m/v2/math"
	"github.com/javanhut/BifrostEngine/m/v2/renderer/core"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func main() {
	renderer, err := core.New(800, 600, "Bifrost Engine - Shape Demo")
	if err != nil {
		log.Fatal(err)
	}
	defer renderer.Cleanup()

	fmt.Println("=== 3D Shape Demo ===")
	fmt.Println("Press number keys to switch modes:")
	fmt.Println("1 - Static Triangle")
	fmt.Println("2 - Static Cube") 
	fmt.Println("3 - Rotating Triangle")
	fmt.Println("4 - Rotating Cube")
	fmt.Println("5 - Camera orbiting Triangle")
	fmt.Println("6 - Camera orbiting Cube")
	fmt.Println("ESC - Exit")

	mode := 1
	time := float32(0)
	camera := renderer.GetCamera()
	window := renderer.GetWindow()
	
	// Key callback
	window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if action == glfw.Press {
			switch key {
			case glfw.Key1:
				mode = 1
				fmt.Println("Mode 1: Static Triangle")
				// Reset camera for static view
				camera.SetPosition(bmath.NewVector3(0, 0, 3))
				camera.SetTarget(bmath.NewVector3(0, 0, 0))
			case glfw.Key2:
				mode = 2
				fmt.Println("Mode 2: Static Cube")
				camera.SetPosition(bmath.NewVector3(0, 0, 3))
				camera.SetTarget(bmath.NewVector3(0, 0, 0))
			case glfw.Key3:
				mode = 3
				fmt.Println("Mode 3: Rotating Triangle")
				camera.SetPosition(bmath.NewVector3(0, 0, 3))
				camera.SetTarget(bmath.NewVector3(0, 0, 0))
			case glfw.Key4:
				mode = 4
				fmt.Println("Mode 4: Rotating Cube")
				camera.SetPosition(bmath.NewVector3(0, 0, 3))
				camera.SetTarget(bmath.NewVector3(0, 0, 0))
			case glfw.Key5:
				mode = 5
				fmt.Println("Mode 5: Camera orbiting Triangle")
			case glfw.Key6:
				mode = 6
				fmt.Println("Mode 6: Camera orbiting Cube")
			case glfw.KeyEscape:
				w.SetShouldClose(true)
			}
		}
	})
	
	for !renderer.ShouldClose() {
		renderer.BeginFrame()
		
		time += 0.016 // ~60 FPS
		
		switch mode {
		case 1: // Static Triangle
			renderer.DrawTriangle()
			
		case 2: // Static Cube
			renderer.DrawCube()
			
		case 3: // Rotating Triangle
			renderer.DrawRotatingTriangle(time)
			
		case 4: // Rotating Cube
			renderer.DrawRotatingCube(time)
			
		case 5: // Camera orbiting Triangle
			radius := float32(4.0)
			x := radius * float32(math.Sin(float64(time)))
			z := radius * float32(math.Cos(float64(time)))
			y := float32(math.Sin(float64(time * 0.5)))
			
			camera.SetPosition(bmath.NewVector3(x, y, z))
			camera.SetTarget(bmath.NewVector3(0, 0, 0))
			renderer.DrawTriangle()
			
		case 6: // Camera orbiting Cube
			radius := float32(4.0)
			x := radius * float32(math.Sin(float64(time)))
			z := radius * float32(math.Cos(float64(time)))
			y := float32(math.Sin(float64(time * 0.5)))
			
			camera.SetPosition(bmath.NewVector3(x, y, z))
			camera.SetTarget(bmath.NewVector3(0, 0, 0))
			renderer.DrawCube()
		}
		
		renderer.EndFrame()
	}
}