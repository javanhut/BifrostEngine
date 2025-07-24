package main

import (
	"fmt"
	"log"
	"math"
	"strings"

	bmath "github.com/javanhut/BifrostEngine/m/v2/math"
	"github.com/javanhut/BifrostEngine/m/v2/renderer/core"
	"github.com/javanhut/BifrostEngine/m/v2/renderer/window"
	"github.com/javanhut/BifrostEngine/m/v2/ui"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type GUIOverlayEditor struct {
	renderer       *core.Renderer
	editor         *ui.Editor
	guiSystem      *ui.GUISystem
	cameraDistance float32
	cameraAngleX   float32
	cameraAngleY   float32
	lastMouseX     float64
	lastMouseY     float64
	mousePressed   bool
	lastLeftClick  bool
	transformMode  string // "select", "move", "transform"
	objectDragging bool
	dragStartPos   bmath.Vector3
	axisConstraint string // "", "x", "y", "z" - constrains movement to specific axis
}

func main() {
	// Initialize renderer
	renderer, err := core.New(1280, 720, "Bifrost Engine - GUI Overlay Editor")
	if err != nil {
		log.Fatal(err)
	}
	defer renderer.Cleanup()

	// Create editor
	editor := ui.NewEditor()
	guiEditor := &GUIOverlayEditor{
		renderer:       renderer,
		editor:         editor,
		guiSystem:      ui.NewGUISystem(1280, 720, editor),
		cameraDistance: 10.0,
		cameraAngleX:   0.3,
		cameraAngleY:   0.5,
		transformMode:  "select",
		objectDragging: false,
	}
	defer guiEditor.guiSystem.Cleanup()

	// Setup input callbacks
	setupInput(renderer.GetWindow(), guiEditor)
	
	fmt.Println("=== Bifrost Engine GUI Overlay Editor ===")
	fmt.Println("Real GUI with menus and dropdowns!")
	fmt.Println("Click on 'File' and 'Objects' in the menu bar")
	fmt.Println("Controls:")
	fmt.Println("  Q - Select Mode (click objects to select)")
	fmt.Println("  M - Move Mode (drag with mouse OR use arrow keys)")
	fmt.Println("  T - Transform Mode (drag to scale/rotate object)")
	fmt.Println("  R - Reset Camera to default position")
	fmt.Println("Arrow Keys:")
	fmt.Println("  ↑↓ - Move along Y-axis (or constrained axis)")
	fmt.Println("  ←→ - Move along X-axis (or constrained axis)")
	fmt.Println("Axis Constraints:")
	fmt.Println("  X - Constrain to X-axis (arrows become ↑↓=±X, ←→=±X)")
	fmt.Println("  Y - Constrain to Y-axis (arrows become ↑↓=±Y, ←→=±Y)")
	fmt.Println("  Z - Constrain to Z-axis (arrows become ↑↓=±Z, ←→=±Z)")
	fmt.Println("Object properties table in bottom-right corner")
	
	// Main loop
	for !renderer.ShouldClose() {
		// Update camera
		updateCamera(guiEditor)
		
		// Get current mouse state
		mouseX, mouseY := renderer.GetWindow().GetHandle().GetCursorPos()
		leftClick := renderer.GetWindow().GetHandle().GetMouseButton(glfw.MouseButtonLeft) == glfw.Press
		
		// Detect new left clicks (not held down)
		newLeftClick := leftClick && !guiEditor.lastLeftClick
		guiEditor.lastLeftClick = leftClick
		
		// Update GUI system
		guiEditor.guiSystem.Update(mouseX, mouseY, newLeftClick)
		guiEditor.guiSystem.SetCurrentMode(guiEditor.transformMode)
		
		// Begin rendering
		renderer.BeginFrame()
		
		// Render grid if enabled
		if grid := guiEditor.editor.GetGrid(); grid.Visible {
			lines := grid.GetLines()
			renderer.DrawGrid(lines, grid.Color)
		}
		
		// Render 3D scene
		renderScene(renderer, guiEditor.editor)
		
		// Render GUI overlay (includes stats table)
		guiEditor.guiSystem.Render()
		
		// End frame
		renderer.EndFrame()
	}
}

func setupInput(window *window.Window, editor *GUIOverlayEditor) {
	window.SetMouseButtonCallback(func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
		if button == glfw.MouseButtonLeft {
			if action == glfw.Press {
				mouseX, mouseY := w.GetCursorPos()
				
				// Check if click is in GUI area (top 30 pixels for menu bar)
				if mouseY <= 30 {
					// Click is in GUI area, let GUI handle it
					return
				}
				
				// Click is in 3D viewport area
				if editor.transformMode == "select" {
					// Try to select an object
					selectedObj := selectObjectAtMousePos(editor, mouseX, mouseY)
					if selectedObj >= 0 {
						editor.editor.SetSelectedObject(selectedObj)
						fmt.Printf("Selected object: %s\n", editor.editor.GetSceneObjects()[selectedObj].Name)
					}
				} else {
					// Start object transformation
					objects := editor.editor.GetSceneObjects()
					selectedIndex := editor.editor.GetSelectedObject()
					if selectedIndex < len(objects) {
						editor.objectDragging = true
						editor.dragStartPos = objects[selectedIndex].Position
					}
				}
				
				editor.mousePressed = true
				editor.lastMouseX, editor.lastMouseY = mouseX, mouseY
			} else if action == glfw.Release {
				editor.mousePressed = false
				editor.objectDragging = false
			}
		}
	})
	
	window.SetCursorPosCallback(func(w *glfw.Window, xpos, ypos float64) {
		if editor.mousePressed && ypos > 30 {
			deltaX := float32(xpos - editor.lastMouseX)
			deltaY := float32(ypos - editor.lastMouseY)
			
			if editor.objectDragging && editor.transformMode != "select" {
				// Transform the selected object
				objects := editor.editor.GetSceneObjects()
				selectedIndex := editor.editor.GetSelectedObject()
				if selectedIndex < len(objects) {
					obj := &objects[selectedIndex]
					
					switch editor.transformMode {
					case "move":
						// Move object based on mouse movement in camera-relative space
						// Get camera vectors for proper 3D movement
						camera := editor.renderer.GetCamera()
						cameraPos := camera.GetPosition()
						cameraTarget := camera.GetTarget()
						
						// Calculate camera right and up vectors
						forward := cameraTarget.Sub(cameraPos).Normalize()
						worldUp := bmath.NewVector3(0, 1, 0)
						right := forward.Cross(worldUp).Normalize()
						up := right.Cross(forward).Normalize()
						
						// Convert mouse movement to world space
						moveSensitivity := float32(0.02)
						moveX := deltaX * moveSensitivity
						moveY := deltaY * moveSensitivity
						
						// Apply axis constraints if any
						var movement bmath.Vector3
						if editor.axisConstraint == "x" {
							// Move only along world X axis
							movement = bmath.NewVector3(moveX*2, 0, 0)
						} else if editor.axisConstraint == "y" {
							// Move only along world Y axis  
							movement = bmath.NewVector3(0, -moveY*2, 0)
						} else if editor.axisConstraint == "z" {
							// Move only along world Z axis
							movement = bmath.NewVector3(0, 0, -moveX*2)
						} else {
							// Free movement in camera-relative directions
							rightMovement := right.Mul(moveX)
							upMovement := up.Mul(-moveY) // Flip Y for intuitive movement
							movement = rightMovement.Add(upMovement)
						}
						
						// Apply movement to object position
						obj.Position = obj.Position.Add(movement)
						
					case "transform":
						// Transform mode: scale with vertical movement, rotate with horizontal
						transformSensitivity := float32(0.01)
						
						// Scale with vertical mouse movement
						scaleChange := -deltaY * transformSensitivity
						newScale := obj.Scale.X + scaleChange
						if newScale > 0.1 { // Prevent negative scaling
							obj.Scale.X = newScale
							obj.Scale.Y = newScale
							obj.Scale.Z = newScale
						}
						
						// Rotate with horizontal mouse movement
						rotateSensitivity := float32(0.5)
						obj.Rotation.Y += deltaX * rotateSensitivity
					}
					
					// Update the object in the editor
					editor.editor.UpdateObject(selectedIndex, *obj)
				}
			} else if !editor.objectDragging {
				// Rotate camera when not transforming objects
				editor.cameraAngleY += deltaX * 0.01
				editor.cameraAngleX += deltaY * 0.01
				
				// Clamp vertical angle
				if editor.cameraAngleX > 1.5 {
					editor.cameraAngleX = 1.5
				}
				if editor.cameraAngleX < -1.5 {
					editor.cameraAngleX = -1.5
				}
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
	
	// Keyboard shortcuts still work
	window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if action == glfw.Press {
			switch key {
			case glfw.Key1, glfw.KeyF1:
				editor.editor.AddObjectWithType("cube")
				fmt.Println("Added Cube (F1)")
			case glfw.Key2, glfw.KeyF2:
				editor.editor.AddObjectWithType("sphere")
				fmt.Println("Added Sphere (F2)")
			case glfw.Key3, glfw.KeyF3:
				editor.editor.AddObjectWithType("cylinder")
				fmt.Println("Added Cylinder (F3)")
			case glfw.Key4, glfw.KeyF4:
				editor.editor.AddObjectWithType("plane")
				fmt.Println("Added Plane (F4)")
			case glfw.Key5, glfw.KeyF5:
				editor.editor.AddObjectWithType("triangle")
				fmt.Println("Added Triangle (F5)")
			case glfw.Key6, glfw.KeyF6:
				editor.editor.AddObjectWithType("pyramid")
				fmt.Println("Added Pyramid (F6)")
			case glfw.KeyQ:
				editor.transformMode = "select"
				fmt.Println("Mode: Select (click to select objects)")
			case glfw.KeyM:
				editor.transformMode = "move"
				fmt.Println("Mode: Move (drag or use arrow keys to move)")
			case glfw.KeyT:
				editor.transformMode = "transform"
				fmt.Println("Mode: Transform (drag to scale/rotate)")
			case glfw.KeyX:
				if editor.axisConstraint == "x" {
					editor.axisConstraint = ""
					fmt.Println("Axis Constraint: None")
				} else {
					editor.axisConstraint = "x"
					fmt.Println("Axis Constraint: X-axis only")
				}
			case glfw.KeyY:
				if editor.axisConstraint == "y" {
					editor.axisConstraint = ""
					fmt.Println("Axis Constraint: None")
				} else {
					editor.axisConstraint = "y"
					fmt.Println("Axis Constraint: Y-axis only")
				}
			case glfw.KeyZ:
				if editor.axisConstraint == "z" {
					editor.axisConstraint = ""
					fmt.Println("Axis Constraint: None")
				} else {
					editor.axisConstraint = "z"
					fmt.Println("Axis Constraint: Z-axis only")
				}
			case glfw.KeyG:
				grid := editor.editor.GetGrid()
				grid.Visible = !grid.Visible
				fmt.Printf("Grid: %v\n", grid.Visible)
			case glfw.KeyTab:
				objects := editor.editor.GetSceneObjects()
				if len(objects) > 0 {
					current := editor.editor.GetSelectedObject()
					next := (current + 1) % len(objects)
					editor.editor.SetSelectedObject(next)
					fmt.Printf("Selected: %s\n", objects[next].Name)
				}
			case glfw.KeyDelete:
				objects := editor.editor.GetSceneObjects()
				if len(objects) > 0 {
					selectedIndex := editor.editor.GetSelectedObject()
					if selectedIndex < len(objects) {
						objName := objects[selectedIndex].Name
						editor.editor.DeleteObject(selectedIndex)
						fmt.Printf("Deleted: %s\n", objName)
					}
				}
			case glfw.KeyUp:
				moveObjectWithArrowKey(editor, "up")
			case glfw.KeyDown:
				moveObjectWithArrowKey(editor, "down")
			case glfw.KeyLeft:
				moveObjectWithArrowKey(editor, "left")
			case glfw.KeyRight:
				moveObjectWithArrowKey(editor, "right")
			case glfw.KeyR:
				// Reset viewport camera to default position
				editor.cameraDistance = 10.0
				editor.cameraAngleX = 0.3
				editor.cameraAngleY = 0.5
				fmt.Println("Camera reset to default position")
			}
		}
	})
}

func updateCamera(editor *GUIOverlayEditor) {
	// Orbit camera around origin
	x := editor.cameraDistance * float32(math.Cos(float64(editor.cameraAngleX))) * float32(math.Sin(float64(editor.cameraAngleY)))
	y := editor.cameraDistance * float32(math.Sin(float64(editor.cameraAngleX)))
	z := editor.cameraDistance * float32(math.Cos(float64(editor.cameraAngleX))) * float32(math.Cos(float64(editor.cameraAngleY)))
	
	camera := editor.renderer.GetCamera()
	camera.SetPosition(bmath.NewVector3(x, y, z))
	camera.SetTarget(bmath.NewVector3(0, 0, 0))
}

func selectObjectAtMousePos(editor *GUIOverlayEditor, mouseX, mouseY float64) int {
	// Simple bounding box selection for now
	// In a real implementation, this would use raycasting
	objects := editor.editor.GetSceneObjects()
	
	// Convert mouse coordinates to normalized device coordinates
	windowWidth, windowHeight := 1280.0, 720.0
	ndcX := (2.0 * mouseX / windowWidth) - 1.0
	ndcY := 1.0 - (2.0 * mouseY / windowHeight)
	
	// For simplicity, select objects based on screen position proximity
	// This is a basic implementation - a proper system would use 3D raycasting
	for i, obj := range objects {
		if !obj.Visible {
			continue
		}
		
		// Project object position to screen space (simplified)
		// Distance from camera affects selection area
		objDistance := obj.Position.Distance(bmath.NewVector3(0, 0, 0))
		selectionRadius := 0.2 / (objDistance * 0.1 + 1.0) // Smaller radius for distant objects
		
		objScreenX := obj.Position.X * 0.1 // Simplified projection
		objScreenY := obj.Position.Y * 0.1
		
		// Check if mouse is within selection radius
		distanceToObj := math.Sqrt(math.Pow(float64(objScreenX)-ndcX, 2) + math.Pow(float64(objScreenY)-ndcY, 2))
		if distanceToObj < float64(selectionRadius) {
			return i
		}
	}
	
	return -1 // No object selected
}

func renderScene(renderer *core.Renderer, editor *ui.Editor) {
	objects := editor.GetSceneObjects()
	selectedIndex := editor.GetSelectedObject()
	
	for i, obj := range objects {
		if !obj.Visible {
			continue
		}
		
		// Create transform matrix (simplified for debugging)
		// Start with identity and apply only translation
		model := bmath.NewMatrix4Identity()
		
		// Apply translation directly
		model[12] = obj.Position.X  // X translation
		model[13] = obj.Position.Y  // Y translation  
		model[14] = obj.Position.Z  // Z translation
		
		// Apply uniform scale
		model[0] = obj.Scale.X   // X scale
		model[5] = obj.Scale.Y   // Y scale
		model[10] = obj.Scale.Z  // Z scale
		
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
			renderer.DrawTriangleMeshWithTransform(model)
		case "sphere":
			renderer.DrawSphereWithTransform(model)
		case "cylinder":
			renderer.DrawCylinderWithTransform(model)
		case "plane":
			renderer.DrawPlaneWithTransform(model)
		case "pyramid":
			renderer.DrawPyramidWithTransform(model)
		default:
			renderer.DrawCubeWithTransform(model)
		}
	}
}


func moveObjectWithArrowKey(editor *GUIOverlayEditor, direction string) {
	objects := editor.editor.GetSceneObjects()
	selectedIndex := editor.editor.GetSelectedObject()
	
	if selectedIndex >= 0 && selectedIndex < len(objects) {
		obj := &objects[selectedIndex]
		
		// Movement step size
		step := float32(0.1)
		
		// Apply movement based on direction and axis constraint
		var movement bmath.Vector3
		
		switch direction {
		case "up":
			if editor.axisConstraint == "x" {
				movement = bmath.NewVector3(step, 0, 0)
			} else if editor.axisConstraint == "z" {
				movement = bmath.NewVector3(0, 0, -step)
			} else {
				movement = bmath.NewVector3(0, step, 0) // Default Y up
			}
		case "down":
			if editor.axisConstraint == "x" {
				movement = bmath.NewVector3(-step, 0, 0)
			} else if editor.axisConstraint == "z" {
				movement = bmath.NewVector3(0, 0, step)
			} else {
				movement = bmath.NewVector3(0, -step, 0) // Default Y down
			}
		case "left":
			if editor.axisConstraint == "y" {
				movement = bmath.NewVector3(0, -step, 0)
			} else if editor.axisConstraint == "z" {
				movement = bmath.NewVector3(0, 0, -step)
			} else {
				movement = bmath.NewVector3(-step, 0, 0) // Default X left
			}
		case "right":
			if editor.axisConstraint == "y" {
				movement = bmath.NewVector3(0, step, 0)
			} else if editor.axisConstraint == "z" {
				movement = bmath.NewVector3(0, 0, step)
			} else {
				movement = bmath.NewVector3(step, 0, 0) // Default X right
			}
		}
		
		// Apply movement directly to the object in the slice
		obj.Position = obj.Position.Add(movement)
		
		// Print feedback
		constraintText := ""
		if editor.axisConstraint != "" {
			constraintText = fmt.Sprintf(" [%s-axis]", strings.ToUpper(editor.axisConstraint))
		}
		fmt.Printf("Moved %s %s%s - Position: (%.2f, %.2f, %.2f)\n", 
			obj.Name, direction, constraintText, obj.Position.X, obj.Position.Y, obj.Position.Z)
	}
}