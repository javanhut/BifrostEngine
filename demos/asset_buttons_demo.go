package main

import (
	"fmt"
	"log"

	bmath "github.com/javanhut/BifrostEngine/m/v2/math"
	"github.com/javanhut/BifrostEngine/m/v2/renderer/core"
	"github.com/javanhut/BifrostEngine/m/v2/ui"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func main() {
	renderer, err := core.New(1280, 720, "Bifrost Engine - Asset Buttons Test")
	if err != nil {
		log.Fatal(err)
	}
	defer renderer.Cleanup()

	fmt.Println("=== Asset Buttons Test ===")
	fmt.Println("Testing asset browser button functionality")
	fmt.Println()

	// Load test mesh
	objFilePath := "assets/test_cube.obj"
	fmt.Printf("Loading test mesh: %s\n", objFilePath)
	err = renderer.LoadMesh(objFilePath)
	if err != nil {
		log.Printf("Warning: Could not load test mesh: %v", err)
	} else {
		fmt.Println("✅ Test mesh loaded successfully!")
	}

	// Create editor with renderer
	editor := ui.NewEditor()
	editor.SetRenderer(renderer)

	// Test the editor's asset methods directly
	fmt.Println("\n=== Direct Editor Method Tests ===")
	
	// Test GetAssetStats
	meshCount, vertices, indices, err := editor.GetAssetStats()
	if err != nil {
		fmt.Printf("❌ GetAssetStats failed: %v\n", err)
	} else {
		fmt.Printf("✅ GetAssetStats: %d meshes, %d vertices, %d indices\n", meshCount, vertices, indices)
	}

	// Test GetLoadedMeshNames
	meshNames, err := editor.GetLoadedMeshNames()
	if err != nil {
		fmt.Printf("❌ GetLoadedMeshNames failed: %v\n", err)
	} else {
		fmt.Printf("✅ GetLoadedMeshNames: %v\n", meshNames)
	}

	// Test ClearAssetCache (but don't actually clear it yet)
	fmt.Println("\n=== Asset Cache Operations ===")
	fmt.Println("Before clear - checking assets...")
	meshNames, _ = editor.GetLoadedMeshNames()
	fmt.Printf("Loaded meshes before clear: %v\n", meshNames)

	// Now clear the cache
	err = editor.ClearAssetCache()
	if err != nil {
		fmt.Printf("❌ ClearAssetCache failed: %v\n", err)
	} else {
		fmt.Println("✅ ClearAssetCache succeeded!")
	}

	// Check if cache was cleared
	meshNames, _ = editor.GetLoadedMeshNames()
	fmt.Printf("Loaded meshes after clear: %v\n", meshNames)

	// Reload the asset to test again
	fmt.Println("\nReloading asset for GUI test...")
	err = renderer.LoadMesh(objFilePath)
	if err != nil {
		log.Printf("Could not reload test mesh: %v", err)
	}

	// Create GUI system
	guiSystem := ui.NewGUISystem(1280, 720, editor)
	defer guiSystem.Cleanup()

	fmt.Println("\n=== GUI Test ===")
	fmt.Println("Now testing with GUI - try the Assets menu!")
	fmt.Println("Controls:")
	fmt.Println("  Click 'Assets' menu")
	fmt.Println("  Try 'Asset Statistics' - should print to console")
	fmt.Println("  Try 'Browse Assets' - should open browser window")
	fmt.Println("  Try 'Clear Cache' - should clear and confirm")
	fmt.Println("  ESC - Exit")

	mouseX, mouseY := 0.0, 0.0
	leftClick := false

	// Setup input callbacks
	window := renderer.GetWindow()
	window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if action == glfw.Press && key == glfw.KeyEscape {
			w.SetShouldClose(true)
		}
	})

	window.SetMouseButtonCallback(func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
		if button == glfw.MouseButtonLeft {
			leftClick = (action == glfw.Press)
		}
	})

	window.SetCursorPosCallback(func(w *glfw.Window, xpos, ypos float64) {
		mouseX = xpos
		mouseY = ypos
	})

	// Main render loop
	for !renderer.ShouldClose() {
		renderer.BeginFrame()

		// Update GUI system with mouse state
		guiSystem.Update(mouseX, mouseY, leftClick)

		// Render a simple scene
		model := bmath.NewMatrix4Identity()
		renderer.DrawCubeWithLighting(model, false)

		// Render GUI on top
		guiSystem.Render()

		renderer.EndFrame()
		
		// Reset leftClick to prevent multiple triggers
		leftClick = false
	}

	fmt.Println("\nAsset buttons test completed!")
}