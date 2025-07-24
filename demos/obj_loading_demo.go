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
	renderer, err := core.New(1280, 720, "Bifrost Engine - OBJ Loading Demo")
	if err != nil {
		log.Fatal(err)
	}
	defer renderer.Cleanup()

	fmt.Println("=== OBJ Loading Demo ===")
	fmt.Println("Testing asset loading system with .obj files")
	fmt.Println()

	// Load the test cube OBJ file
	objFilePath := "assets/test_cube.obj"
	fmt.Printf("Loading OBJ file: %s\n", objFilePath)
	
	err = renderer.LoadMesh(objFilePath)
	if err != nil {
		log.Fatalf("Failed to load OBJ file: %v", err)
	}
	
	fmt.Println("✅ OBJ file loaded successfully!")
	
	// Get asset statistics
	stats := renderer.GetAssetStats()
	fmt.Printf("Asset Stats: %d meshes, %d vertices, %d indices\n", 
		stats.LoadedMeshes, stats.TotalVertices, stats.TotalIndices)
	
	// List loaded meshes
	loadedMeshes := renderer.GetLoadedMeshes()
	fmt.Println("\nLoaded meshes:")
	for _, meshName := range loadedMeshes {
		fmt.Printf("  - %s\n", meshName)
	}
	fmt.Println()
	
	// Setup lighting for the demo
	lightingSystem := renderer.GetLightingSystem()
	lightingSystem.SetSunIntensity(0.8)
	lightingSystem.SetAmbientLight([3]float32{0.2, 0.2, 0.2}, 0.3)
	
	// Add a point light for better illumination
	lightingSystem.AddPointLight(
		bmath.NewVector3(2.0, 2.0, 2.0),
		[3]float32{1.0, 1.0, 1.0}, // White light
		1.5, // Intensity
	)

	fmt.Println("Controls:")
	fmt.Println("  SPACE - Toggle texture rendering")
	fmt.Println("  R - Reset camera")
	fmt.Println("  ESC - Exit")
	fmt.Println()

	useTextures := true
	time := float32(0.0)

	// Setup input
	window := renderer.GetWindow()
	window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if action == glfw.Press {
			switch key {
			case glfw.KeyEscape:
				w.SetShouldClose(true)
			case glfw.KeySpace:
				useTextures = !useTextures
				fmt.Printf("Textures: %s\n", map[bool]string{true: "enabled", false: "disabled"}[useTextures])
			case glfw.KeyR:
				// Reset camera position
				camera := renderer.GetCamera()
				camera.SetPosition(bmath.NewVector3(0, 0, 3))
				camera.SetTarget(bmath.NewVector3(0, 0, 0))
				fmt.Println("Camera reset to default position")
			}
		}
	})

	// Main render loop
	for !renderer.ShouldClose() {
		renderer.BeginFrame()
		
		time += 0.016

		// Render the loaded OBJ mesh with rotation animation
		objModel := bmath.NewMatrix4Identity()
		
		// Apply rotation around Y axis
		rotationY := time * 0.5
		objModel[0] = float32(math.Cos(float64(rotationY)))
		objModel[2] = float32(math.Sin(float64(rotationY)))
		objModel[8] = -float32(math.Sin(float64(rotationY)))
		objModel[10] = float32(math.Cos(float64(rotationY)))
		
		// Position the loaded mesh at origin
		objModel[12] = 0.0 // X
		objModel[13] = 0.0 // Y
		objModel[14] = 0.0 // Z
		
		// Draw the loaded OBJ mesh
		err = renderer.DrawLoadedMesh(objFilePath, objModel, useTextures)
		if err != nil {
			fmt.Printf("Error drawing loaded mesh: %v\n", err)
		}
		
		// Also render a comparison primitive cube for reference
		primitiveModel := bmath.NewMatrix4Identity()
		primitiveModel[12] = 3.0 // Offset to the right
		primitiveModel[13] = 0.0
		primitiveModel[14] = 0.0
		
		// Apply same rotation
		primitiveModel[0] = float32(math.Cos(float64(rotationY)))
		primitiveModel[2] = float32(math.Sin(float64(rotationY)))
		primitiveModel[8] = -float32(math.Sin(float64(rotationY)))
		primitiveModel[10] = float32(math.Cos(float64(rotationY)))
		
		renderer.DrawCubeWithLighting(primitiveModel, useTextures)
		
		// Render ground plane for reference
		groundModel := bmath.NewMatrix4Identity()
		groundModel[13] = -2.0 // Y position (below objects)
		groundModel[0] = 5.0   // X scale (wider)
		groundModel[10] = 5.0  // Z scale (deeper)
		renderer.DrawPlaneWithLighting(groundModel, useTextures)

		renderer.EndFrame()
	}

	fmt.Println("\nDemo completed!")
	fmt.Println("OBJ loading system is working correctly.")
}