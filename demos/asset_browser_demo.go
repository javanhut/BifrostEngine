package main

import (
	"fmt"
	"log"

	bmath "github.com/javanhut/BifrostEngine/m/v2/math"
	"github.com/javanhut/BifrostEngine/m/v2/renderer/core"
	"github.com/javanhut/BifrostEngine/m/v2/ui"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type AssetBrowserTest struct {
	renderer  *core.Renderer
	editor    *ui.Editor
	guiSystem *ui.GUISystem
	mouseX    float64
	mouseY    float64
	leftClick bool
}

func main() {
	renderer, err := core.New(1280, 720, "Bifrost Engine - Asset Browser Test")
	if err != nil {
		log.Fatal(err)
	}
	defer renderer.Cleanup()

	fmt.Println("=== Asset Browser Test ===")
	fmt.Println("Testing the asset browser integration")
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
	
	// Create GUI system
	guiSystem := ui.NewGUISystem(1280, 720, editor)
	defer guiSystem.Cleanup()

	test := &AssetBrowserTest{
		renderer:  renderer,
		editor:    editor,
		guiSystem: guiSystem,
	}

	// Setup input callbacks
	window := renderer.GetWindow()
	window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if action == glfw.Press {
			switch key {
			case glfw.KeyEscape:
				w.SetShouldClose(true)
			case glfw.KeyB:
				fmt.Println("Toggled asset browser via B key")
			}
		}
	})

	// Mouse callbacks
	window.SetMouseButtonCallback(func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
		if button == glfw.MouseButtonLeft {
			test.leftClick = (action == glfw.Press)
		}
	})

	window.SetCursorPosCallback(func(w *glfw.Window, xpos, ypos float64) {
		test.mouseX = xpos
		test.mouseY = ypos
	})

	fmt.Println("Controls:")
	fmt.Println("  Click 'Assets' menu to access asset browser")
	fmt.Println("  Try 'Browse Assets' to open the asset browser window")
	fmt.Println("  ESC - Exit")
	fmt.Println()

	// Main render loop
	for !renderer.ShouldClose() {
		renderer.BeginFrame()

		// Update GUI system with mouse state
		guiSystem.Update(test.mouseX, test.mouseY, test.leftClick)

		// Render a simple scene
		model := bmath.NewMatrix4Identity()
		renderer.DrawCubeWithLighting(model, false)

		// Render GUI on top
		guiSystem.Render()

		renderer.EndFrame()
	}

	fmt.Println("Asset browser test completed!")
}