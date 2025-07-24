package main

import (
	"fmt"
	"log"
	"math"

	bmath "github.com/javanhut/BifrostEngine/m/v2/math"
	"github.com/javanhut/BifrostEngine/m/v2/renderer/core"
	"github.com/javanhut/BifrostEngine/m/v2/ui"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func main() {
	renderer, err := core.New(1280, 720, "Bifrost Engine - Gizmo Demo")
	if err != nil {
		log.Fatal(err)
	}
	defer renderer.Cleanup()

	fmt.Println("=== Gizmo Demo ===")
	fmt.Println("This demo showcases the Transform Gizmo system")
	fmt.Println()
	fmt.Println("Controls:")
	fmt.Println("  Mouse - Orbit camera")
	fmt.Println("  Scroll - Zoom in/out")
	fmt.Println("  G - Toggle gizmo visibility")
	fmt.Println("  T - Toggle between translate/rotate/scale gizmos")
	fmt.Println("  1-5 - Switch gizmo scale size")
	fmt.Println("  WASD - Move object")
	fmt.Println("  ESC - Exit")
	fmt.Println()

	// Create editor and GUI
	editor := ui.NewEditor()
	editor.SetRenderer(renderer)
	guiSystem := ui.NewGUISystem(1280, 720, editor)
	defer guiSystem.Cleanup()

	// Add some test objects
	editor.AddObjectWithType("cube")
	editor.AddObjectWithType("sphere")
	editor.AddObjectWithType("cylinder")
	
	// Position objects
	objects := editor.GetSceneObjects()
	if len(objects) >= 3 {
		editor.UpdateObject(0, ui.SceneObject{
			Name:     "Red Cube",
			Position: bmath.NewVector3(-2, 0, 0),
			Rotation: bmath.NewVector3(0, 0, 0),
			Scale:    bmath.NewVector3(1, 1, 1),
			Color:    [3]float32{1.0, 0.0, 0.0},
			Type:     "cube",
			Visible:  true,
		})
		editor.UpdateObject(1, ui.SceneObject{
			Name:     "Green Sphere",
			Position: bmath.NewVector3(0, 0, 0),
			Rotation: bmath.NewVector3(0, 0, 0),
			Scale:    bmath.NewVector3(1, 1, 1),
			Color:    [3]float32{0.0, 1.0, 0.0},
			Type:     "sphere",
			Visible:  true,
		})
		editor.UpdateObject(2, ui.SceneObject{
			Name:     "Blue Cylinder",
			Position: bmath.NewVector3(2, 0, 0),
			Rotation: bmath.NewVector3(0, 0, 0),
			Scale:    bmath.NewVector3(1, 1, 1),
			Color:    [3]float32{0.0, 0.0, 1.0},
			Type:     "cylinder",
			Visible:  true,
		})
	}

	// Camera control variables
	cameraDistance := float32(8.0)
	cameraAngleX := float32(0.3)
	cameraAngleY := float32(0.5)
	lastMouseX, lastMouseY := 0.0, 0.0
	mousePressed := false
	
	// Gizmo control variables
	showGizmos := true
	gizmoType := core.GizmoTranslate
	gizmoScale := float32(1.0)

	// Set up input callbacks
	window := renderer.GetWindow()
	
	window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if action == glfw.Press {
			switch key {
			case glfw.KeyEscape:
				w.SetShouldClose(true)
			case glfw.KeyG:
				showGizmos = !showGizmos
				renderer.SetGizmoVisible(showGizmos)
				fmt.Printf("Gizmos: %v\n", showGizmos)
			case glfw.KeyT:
				switch gizmoType {
				case core.GizmoTranslate:
					gizmoType = core.GizmoRotate
					fmt.Println("Gizmo Type: Rotation")
				case core.GizmoRotate:
					gizmoType = core.GizmoScale
					fmt.Println("Gizmo Type: Scale")
				case core.GizmoScale:
					gizmoType = core.GizmoTranslate
					fmt.Println("Gizmo Type: Translation")
				}
				renderer.SetGizmoType(gizmoType)
			case glfw.Key1:
				gizmoScale = 0.5
				renderer.SetGizmoScale(gizmoScale)
				fmt.Printf("Gizmo Scale: %.1f\n", gizmoScale)
			case glfw.Key2:
				gizmoScale = 1.0
				renderer.SetGizmoScale(gizmoScale)
				fmt.Printf("Gizmo Scale: %.1f\n", gizmoScale)
			case glfw.Key3:
				gizmoScale = 1.5
				renderer.SetGizmoScale(gizmoScale)
				fmt.Printf("Gizmo Scale: %.1f\n", gizmoScale)
			case glfw.Key4:
				gizmoScale = 2.0
				renderer.SetGizmoScale(gizmoScale)
				fmt.Printf("Gizmo Scale: %.1f\n", gizmoScale)
			case glfw.Key5:
				gizmoScale = 3.0
				renderer.SetGizmoScale(gizmoScale)
				fmt.Printf("Gizmo Scale: %.1f\n", gizmoScale)
			case glfw.KeyW, glfw.KeyA, glfw.KeyS, glfw.KeyD:
				// Move selected object
				objects := editor.GetSceneObjects()
				selectedIndex := editor.GetSelectedObject()
				if selectedIndex >= 0 && selectedIndex < len(objects) {
					obj := objects[selectedIndex]
					step := float32(0.2)
					
					switch key {
					case glfw.KeyW:
						obj.Position.Z -= step
					case glfw.KeyS:
						obj.Position.Z += step
					case glfw.KeyA:
						obj.Position.X -= step
					case glfw.KeyD:
						obj.Position.X += step
					}
					
					editor.UpdateObject(selectedIndex, obj)
					fmt.Printf("Moved %s to (%.1f, %.1f, %.1f)\n", obj.Name, obj.Position.X, obj.Position.Y, obj.Position.Z)
				}
			case glfw.KeySpace:
				// Cycle through objects
				objects := editor.GetSceneObjects()
				selectedIndex := editor.GetSelectedObject()
				nextIndex := (selectedIndex + 1) % len(objects)
				editor.SetSelectedObject(nextIndex)
				fmt.Printf("Selected: %s\n", objects[nextIndex].Name)
			}
		}
	})

	window.SetMouseButtonCallback(func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
		if button == glfw.MouseButtonLeft {
			if action == glfw.Press {
				mouseX, mouseY := w.GetCursorPos()
				// Test gizmo interaction first
				selectedAxis := renderer.HandleGizmoMouseClick(mouseX, mouseY)
				if selectedAxis != core.GizmoAxisNone {
					fmt.Printf("Selected gizmo axis: %d\n", int(selectedAxis))
				} else {
					// No gizmo interaction, handle camera movement
					mousePressed = true
					lastMouseX, lastMouseY = mouseX, mouseY
				}
			} else {
				mousePressed = false
			}
		}
	})

	window.SetCursorPosCallback(func(w *glfw.Window, xpos, ypos float64) {
		// Always update gizmo hover state
		renderer.HandleGizmoMouseMove(xpos, ypos)
		
		// Handle camera movement if mouse is pressed and no gizmo is selected
		if mousePressed && renderer.GetGizmoSelectedAxis() == core.GizmoAxisNone {
			deltaX := float32(xpos - lastMouseX)
			deltaY := float32(ypos - lastMouseY)
			
			cameraAngleY += deltaX * 0.01
			cameraAngleX -= deltaY * 0.01
			
			// Clamp X rotation
			if cameraAngleX > 1.5 {
				cameraAngleX = 1.5
			}
			if cameraAngleX < -1.5 {
				cameraAngleX = -1.5
			}
			
			lastMouseX = xpos
			lastMouseY = ypos
		}
	})

	window.SetScrollCallback(func(w *glfw.Window, xoff, yoff float64) {
		cameraDistance -= float32(yoff) * 0.5
		if cameraDistance < 2.0 {
			cameraDistance = 2.0
		}
		if cameraDistance > 20.0 {
			cameraDistance = 20.0
		}
	})

	// Initialize gizmo settings
	renderer.SetGizmoVisible(showGizmos)
	renderer.SetGizmoType(gizmoType)
	renderer.SetGizmoScale(gizmoScale)

	fmt.Println("Demo started - try the controls!")

	// Main render loop
	for !renderer.ShouldClose() {
		// Update camera position
		x := cameraDistance * float32(math.Cos(float64(cameraAngleY))) * float32(math.Cos(float64(cameraAngleX)))
		y := cameraDistance * float32(math.Sin(float64(cameraAngleX)))
		z := cameraDistance * float32(math.Sin(float64(cameraAngleY))) * float32(math.Cos(float64(cameraAngleX)))
		
		camera := renderer.GetCamera()
		camera.SetPosition(bmath.NewVector3(x, y, z))
		camera.SetTarget(bmath.NewVector3(0, 0, 0))

		renderer.BeginFrame() // Gizmo demo uses default fill mode

		// Render grid
		if grid := editor.GetGrid(); grid.Visible {
			lines := grid.GetLines()
			renderer.DrawGrid(lines, grid.Color)
		}

		// Render objects
		objects := editor.GetSceneObjects()
		selectedIndex := editor.GetSelectedObject()
		
		for i, obj := range objects {
			if !obj.Visible {
				continue
			}
			
			// Create transform matrix
			model := bmath.NewMatrix4Identity()
			model[12] = obj.Position.X
			model[13] = obj.Position.Y
			model[14] = obj.Position.Z
			model[0] = obj.Scale.X
			model[5] = obj.Scale.Y
			model[10] = obj.Scale.Z
			
			// Highlight selected object
			if i == selectedIndex {
				highlightScale := bmath.NewScaleMatrix(1.1, 1.1, 1.1)
				model = model.Multiply(highlightScale)
			}
			
			// Render based on type
			switch obj.Type {
			case "cube":
				renderer.DrawCubeWithLighting(model, false)
			case "sphere":
				renderer.DrawSphereWithLighting(model, false)
			case "cylinder":
				renderer.DrawCylinderWithLighting(model, false)
			default:
				renderer.DrawCubeWithLighting(model, false)
			}
		}

		// Render gizmo for selected object
		if showGizmos && selectedIndex >= 0 && selectedIndex < len(objects) {
			selectedObj := objects[selectedIndex]
			renderer.RenderGizmo(selectedObj.Position)
		}

		// Render GUI overlay with info
		guiSystem.Update(0, 0, false) // Minimal GUI update
		
		renderer.EndFrame()
	}

	fmt.Println("\nGizmo demo completed!")
}