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
	renderer, err := core.New(1280, 720, "Bifrost Engine - Lighting Demo")
	if err != nil {
		log.Fatal(err)
	}
	defer renderer.Cleanup()

	// Get the lighting system
	lightingSystem := renderer.GetLightingSystem()
	
	fmt.Println("=== Bifrost Engine Lighting Demo ===")
	fmt.Println("Current lighting setup:")
	fmt.Printf("Directional Lights: %d\n", len(lightingSystem.DirectionalLights))
	if len(lightingSystem.DirectionalLights) > 0 {
		light := lightingSystem.DirectionalLights[0]
		fmt.Printf("  Sun Direction: (%.2f, %.2f, %.2f)\n", light.Direction.X, light.Direction.Y, light.Direction.Z)
		fmt.Printf("  Sun Color: (%.2f, %.2f, %.2f)\n", light.Color[0], light.Color[1], light.Color[2])
		fmt.Printf("  Sun Intensity: %.2f\n", light.Intensity)
		fmt.Printf("  Sun Enabled: %t\n", light.Enabled)
	}
	
	// Add a point light for demonstration
	lightingSystem.AddPointLight(
		bmath.NewVector3(2.0, 2.0, 2.0), // Position above and to the right
		[3]float32{0.8, 0.4, 1.0},       // Purple color
		1.0,                             // Full intensity
	)
	
	fmt.Printf("Point Lights: %d\n", len(lightingSystem.PointLights))
	if len(lightingSystem.PointLights) > 0 {
		light := lightingSystem.PointLights[0]
		fmt.Printf("  Point Light Position: (%.2f, %.2f, %.2f)\n", light.Position.X, light.Position.Y, light.Position.Z)
		fmt.Printf("  Point Light Color: (%.2f, %.2f, %.2f)\n", light.Color[0], light.Color[1], light.Color[2])
		fmt.Printf("  Point Light Intensity: %.2f\n", light.Intensity)
	}
	
	fmt.Printf("Ambient Light: (%.2f, %.2f, %.2f) intensity %.2f\n", 
		lightingSystem.AmbientLight.Color[0], 
		lightingSystem.AmbientLight.Color[1], 
		lightingSystem.AmbientLight.Color[2], 
		lightingSystem.AmbientLight.Intensity)
	
	fmt.Println("\nControls:")
	fmt.Println("  SPACE - Toggle directional light on/off")
	fmt.Println("  P - Toggle point light on/off")
	fmt.Println("  ESC - Exit")
	
	// Setup input
	window := renderer.GetWindow()
	window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if action == glfw.Press {
			switch key {
			case glfw.KeyEscape:
				w.SetShouldClose(true)
			case glfw.KeySpace:
				lightingSystem.ToggleDirectionalLight(0)
				enabled := lightingSystem.DirectionalLights[0].Enabled
				fmt.Printf("Directional light: %t\n", enabled)
			case glfw.KeyP:
				if len(lightingSystem.PointLights) > 0 {
					lightingSystem.TogglePointLight(0)
					enabled := lightingSystem.PointLights[0].Enabled
					fmt.Printf("Point light: %t\n", enabled)
				}
			}
		}
	})

	// Main loop - render a lit cube and sphere
	time := float32(0.0)
	for !renderer.ShouldClose() {
		renderer.BeginFrame()
		
		time += 0.016 // ~60 FPS
		
		// Render a cube
		cubeModel := bmath.NewMatrix4Identity()
		cubeModel[12] = -1.5 // X position
		renderer.DrawCubeWithLighting(cubeModel, false)
		
		// Render a sphere
		sphereModel := bmath.NewMatrix4Identity()
		sphereModel[12] = 1.5 // X position
		renderer.DrawSphereWithLighting(sphereModel, false)
		
		// Animate point light position
		if len(lightingSystem.PointLights) > 0 {
			pointLight := &lightingSystem.PointLights[0]
			pointLight.Position.X = 2.0 * float32(math.Cos(float64(time*0.5)))
			pointLight.Position.Z = 2.0 * float32(math.Sin(float64(time*0.5)))
		}
		
		renderer.EndFrame()
	}
}