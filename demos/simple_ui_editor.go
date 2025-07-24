package main

import (
	"fmt"
	"log"
	"math"

	bmath "github.com/javanhut/BifrostEngine/m/v2/math"
	"github.com/javanhut/BifrostEngine/m/v2/renderer/core"
	renderer "github.com/javanhut/BifrostEngine/m/v2/renderer/window"
	"github.com/javanhut/BifrostEngine/m/v2/ui"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type SimpleUIEditor struct {
	renderer       *core.Renderer
	editor         *ui.Editor
	cameraDistance float32
	cameraAngleX   float32
	cameraAngleY   float32
	lastMouseX     float64
	lastMouseY     float64
	mousePressed   bool
	showMenu       bool
	menuSelection  int
}

func main() {
	// Initialize renderer
	renderer, err := core.New(1280, 720, "Bifrost Engine - Editor")
	if err != nil {
		log.Fatal(err)
	}
	defer renderer.Cleanup()

	// Create editor
	uiEditor := &SimpleUIEditor{
		renderer:       renderer,
		editor:         ui.NewEditor(),
		cameraDistance: 10.0,
		cameraAngleX:   0.3,
		cameraAngleY:   0.5,
	}

	// Setup input callbacks
	setupMouseInput(renderer.GetWindow(), uiEditor)
	setupKeyboardInput(renderer.GetWindow(), uiEditor)
	
	fmt.Println("=== Bifrost Engine Editor ===")
	fmt.Println("Controls:")
	fmt.Println("  Mouse: Orbit camera")
	fmt.Println("  Scroll: Zoom in/out")
	fmt.Println("  1-5: Add objects (cube, sphere, cylinder, plane, triangle)")
	fmt.Println("  Tab: Select next object")
	fmt.Println("  Delete: Remove selected object")
	fmt.Println("  G: Toggle grid")
	fmt.Println("  P: Open project menu (console)")
	fmt.Println()
	
	// Main loop
	for !renderer.ShouldClose() {
		// Update camera
		updateCamera(uiEditor)
		
		// Begin rendering
		renderer.BeginFrame()
		
		// Render grid if enabled
		if grid := uiEditor.editor.GetGrid(); grid.Visible {
			lines := grid.GetLines()
			renderer.DrawGrid(lines, grid.Color)
		}
		
		// Render 3D scene
		renderScene(renderer, uiEditor.editor)
		
		// Print status and menu
		printStatus(uiEditor.editor)
		printMenuOverlay(uiEditor)
		
		// End frame
		renderer.EndFrame()
	}
}

func setupMouseInput(window *renderer.Window, editor *SimpleUIEditor) {
	window.SetMouseButtonCallback(func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
		if button == glfw.MouseButtonLeft {
			if action == glfw.Press {
				editor.mousePressed = true
				editor.lastMouseX, editor.lastMouseY = w.GetCursorPos()
			} else if action == glfw.Release {
				editor.mousePressed = false
			}
		}
	})
	
	window.SetCursorPosCallback(func(w *glfw.Window, xpos, ypos float64) {
		if editor.mousePressed {
			deltaX := float32(xpos - editor.lastMouseX)
			deltaY := float32(ypos - editor.lastMouseY)
			
			editor.cameraAngleY += deltaX * 0.01
			editor.cameraAngleX += deltaY * 0.01
			
			// Clamp vertical angle
			if editor.cameraAngleX > 1.5 {
				editor.cameraAngleX = 1.5
			}
			if editor.cameraAngleX < -1.5 {
				editor.cameraAngleX = -1.5
			}
			
			editor.lastMouseX = xpos
			editor.lastMouseY = ypos
		}
	})
	
	window.SetScrollCallback(func(w *glfw.Window, xoff, yoff float64) {
		editor.cameraDistance -= float32(yoff) * 0.5
		if editor.cameraDistance < 1.0 {
			editor.cameraDistance = 1.0
		}
		if editor.cameraDistance > 50.0 {
			editor.cameraDistance = 50.0
		}
	})
}

func setupKeyboardInput(window *renderer.Window, editor *SimpleUIEditor) {
	window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if action == glfw.Press {
			switch key {
			case glfw.Key1:
				editor.editor.AddObjectWithType("cube")
				fmt.Println("Added cube")
			case glfw.Key2:
				editor.editor.AddObjectWithType("sphere")
				fmt.Println("Added sphere")
			case glfw.Key3:
				editor.editor.AddObjectWithType("cylinder")
				fmt.Println("Added cylinder")
			case glfw.Key4:
				editor.editor.AddObjectWithType("plane")
				fmt.Println("Added plane")
			case glfw.Key5:
				editor.editor.AddObjectWithType("triangle")
				fmt.Println("Added triangle")
			case glfw.KeyTab:
				objects := editor.editor.GetSceneObjects()
				if len(objects) > 0 {
					current := editor.editor.GetSelectedObject()
					next := (current + 1) % len(objects)
					editor.editor.SetSelectedObject(next)
					fmt.Printf("Selected: %s\n", objects[next].Name)
				}
			case glfw.KeyDelete:
				editor.editor.DeleteObject(editor.editor.GetSelectedObject())
				fmt.Println("Deleted object")
			case glfw.KeyG:
				grid := editor.editor.GetGrid()
				grid.Visible = !grid.Visible
				fmt.Printf("Grid: %v\n", grid.Visible)
			case glfw.KeyM:
				editor.showMenu = !editor.showMenu
				fmt.Printf("Menu: %v\n", editor.showMenu)
			case glfw.KeyP:
				fmt.Println("\n=== Project Menu ===")
				if proj := editor.editor.GetProjectManager().GetCurrentProject(); proj != nil {
					fmt.Printf("Current Project: %s\n", proj.Name)
				}
				fmt.Println("Press N for new project")
			case glfw.KeyN:
				if mods&glfw.ModControl != 0 {
					editor.editor.GetProjectManager().CreateProject(fmt.Sprintf("Project_%d", len(editor.editor.GetProjectManager().GetProjects())+1))
					fmt.Println("Created new project")
				}
			// Function keys for quick menu actions
			case glfw.KeyF1:
				editor.editor.AddObjectWithType("cube")
				fmt.Println("Added cube via menu")
			case glfw.KeyF2:
				editor.editor.AddObjectWithType("sphere")
				fmt.Println("Added sphere via menu")
			case glfw.KeyF3:
				editor.editor.AddObjectWithType("cylinder")
				fmt.Println("Added cylinder via menu")
			case glfw.KeyF4:
				editor.editor.AddObjectWithType("plane")
				fmt.Println("Added plane via menu")
			case glfw.KeyF5:
				editor.editor.AddObjectWithType("triangle")
				fmt.Println("Added triangle via menu")
			case glfw.KeyF6:
				editor.editor.GetProjectManager().CreateProject(fmt.Sprintf("Project_%d", len(editor.editor.GetProjectManager().GetProjects())+1))
				fmt.Println("Created new project via menu")
			case glfw.KeyF7:
				editor.editor.GetProjectManager().SaveCurrentProject()
				fmt.Println("Saved project via menu")
			}
		}
	})
}

func updateCamera(editor *SimpleUIEditor) {
	// Orbit camera around origin
	x := editor.cameraDistance * float32(math.Cos(float64(editor.cameraAngleX))) * float32(math.Sin(float64(editor.cameraAngleY)))
	y := editor.cameraDistance * float32(math.Sin(float64(editor.cameraAngleX)))
	z := editor.cameraDistance * float32(math.Cos(float64(editor.cameraAngleX))) * float32(math.Cos(float64(editor.cameraAngleY)))
	
	camera := editor.renderer.GetCamera()
	camera.SetPosition(bmath.NewVector3(x, y, z))
	camera.SetTarget(bmath.NewVector3(0, 0, 0))
	
	// Update editor's camera position for UI display
	editor.editor.SetCameraPosition(bmath.NewVector3(x, y, z))
}

func renderScene(renderer *core.Renderer, editor *ui.Editor) {
	objects := editor.GetSceneObjects()
	selectedIndex := editor.GetSelectedObject()
	
	for i, obj := range objects {
		if !obj.Visible {
			continue
		}
		
		// Create transform matrix
		translation := bmath.NewTranslationMatrix(obj.Position.X, obj.Position.Y, obj.Position.Z)
		rotationX := bmath.NewRotationX(bmath.Radians(obj.Rotation.X))
		rotationY := bmath.NewRotationY(bmath.Radians(obj.Rotation.Y))
		rotationZ := bmath.NewRotationZ(bmath.Radians(obj.Rotation.Z))
		scale := bmath.NewScaleMatrix(obj.Scale.X, obj.Scale.Y, obj.Scale.Z)
		
		// Combine transformations
		rotation := rotationZ.Multiply(rotationY).Multiply(rotationX)
		model := translation.Multiply(rotation).Multiply(scale)
		
		// Highlight selected object
		if i == selectedIndex {
			highlightScale := bmath.NewScaleMatrix(1.05, 1.05, 1.05)
			model = model.Multiply(highlightScale)
		}
		
		// Render based on object type
		switch obj.Type {
		case "cube":
			renderer.DrawCubeWithTransform(model)
		case "triangle":
			renderer.DrawTriangleWithTransform(model)
		case "sphere":
			// For now, render as cube until sphere mesh is implemented
			renderer.DrawCubeWithTransform(model)
		case "cylinder":
			// For now, render as cube until cylinder mesh is implemented
			renderer.DrawCubeWithTransform(model)
		case "plane":
			// For now, render as flat cube
			planeScale := bmath.NewScaleMatrix(obj.Scale.X, 0.1, obj.Scale.Z)
			rotation := rotationZ.Multiply(rotationY).Multiply(rotationX)
			planeModel := translation.Multiply(rotation).Multiply(planeScale)
			renderer.DrawCubeWithTransform(planeModel)
		default:
			renderer.DrawCubeWithTransform(model)
		}
	}
}

var frameCount = 0
var menuCounter = 0

func printStatus(editor *ui.Editor) {
	frameCount++
	if frameCount%60 == 0 {
		objects := editor.GetSceneObjects()
		if len(objects) > 0 && editor.GetSelectedObject() < len(objects) {
			obj := objects[editor.GetSelectedObject()]
			fmt.Printf("\rSelected: %s | Pos: (%.1f,%.1f,%.1f) | Objects: %d",
				obj.Name, obj.Position.X, obj.Position.Y, obj.Position.Z, len(objects))
		}
	}
}

func printMenuOverlay(editor *SimpleUIEditor) {
	if !editor.showMenu {
		return
	}
	
	// Print menu overlay every few seconds to show options
	menuCounter++
	if menuCounter%180 == 0 { // Every 3 seconds
		fmt.Println("\n=== MENU OVERLAY ===")
		fmt.Println("F1: Add Cube    F2: Add Sphere   F3: Add Cylinder")
		fmt.Println("F4: Add Plane   F5: Add Triangle")
		fmt.Println("F6: New Project F7: Save Project")
		fmt.Println("Press M to hide menu")
		fmt.Println("=====================")
	}
}