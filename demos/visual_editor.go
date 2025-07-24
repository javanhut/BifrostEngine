package main

import (
	"fmt"
	"log"
	"math"

	bmath "github.com/javanhut/BifrostEngine/m/v2/math"
	"github.com/javanhut/BifrostEngine/m/v2/renderer/core"
	"github.com/javanhut/BifrostEngine/m/v2/renderer/window"
	"github.com/javanhut/BifrostEngine/m/v2/ui"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type VisualEditor struct {
	renderer       *core.Renderer
	editor         *ui.Editor
	cameraDistance float32
	cameraAngleX   float32
	cameraAngleY   float32
	lastMouseX     float64
	lastMouseY     float64
	mousePressed   bool
	showMenu       bool
	statusMessage  string
}

func main() {
	// Initialize renderer
	renderer, err := core.New(1280, 720, "Bifrost Engine - Visual Editor")
	if err != nil {
		log.Fatal(err)
	}
	defer renderer.Cleanup()

	// Create editor
	visualEditor := &VisualEditor{
		renderer:       renderer,
		editor:         ui.NewEditor(),
		cameraDistance: 10.0,
		cameraAngleX:   0.3,
		cameraAngleY:   0.5,
		showMenu:       true, // Start with menu visible
		statusMessage:  "Welcome to Bifrost Engine!",
	}

	// Setup input callbacks
	setupMouseInput(renderer.GetWindow(), visualEditor)
	setupKeyboardInput(renderer.GetWindow(), visualEditor)
	
	// Set window title with instructions
	glfw.GetCurrentContext().SetTitle("Bifrost Engine - Press M for Menu, F1-F7 for Quick Actions")
	
	// Main loop
	for !renderer.ShouldClose() {
		// Update camera
		updateCamera(visualEditor)
		
		// Begin rendering
		renderer.BeginFrame()
		
		// Render grid if enabled
		if grid := visualEditor.editor.GetGrid(); grid.Visible {
			lines := grid.GetLines()
			renderer.DrawGrid(lines, grid.Color)
		}
		
		// Render 3D scene
		renderScene(renderer, visualEditor.editor)
		
		// Render visual overlay (using 3D objects as UI elements)
		renderVisualOverlay(renderer, visualEditor)
		
		// End frame
		renderer.EndFrame()
	}
}

func setupMouseInput(window *window.Window, editor *VisualEditor) {
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

func setupKeyboardInput(window *window.Window, editor *VisualEditor) {
	window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if action == glfw.Press {
			switch key {
			case glfw.KeyM:
				editor.showMenu = !editor.showMenu
				if editor.showMenu {
					editor.statusMessage = "Menu: ON - Use F1-F7 for quick actions"
				} else {
					editor.statusMessage = "Menu: OFF - Press M to show menu"
				}
			case glfw.Key1, glfw.KeyF1:
				editor.editor.AddObjectWithType("cube")
				editor.statusMessage = "Added Cube"
			case glfw.Key2, glfw.KeyF2:
				editor.editor.AddObjectWithType("sphere")
				editor.statusMessage = "Added Sphere"
			case glfw.Key3, glfw.KeyF3:
				editor.editor.AddObjectWithType("cylinder")
				editor.statusMessage = "Added Cylinder"
			case glfw.Key4, glfw.KeyF4:
				editor.editor.AddObjectWithType("plane")
				editor.statusMessage = "Added Plane"
			case glfw.Key5, glfw.KeyF5:
				editor.editor.AddObjectWithType("triangle")
				editor.statusMessage = "Added Triangle"
			case glfw.KeyF6:
				projCount := len(editor.editor.GetProjectManager().GetProjects())\n\t\t\t\tprojectName := fmt.Sprintf(\"Project_%d\", projCount+1)\n\t\t\t\teditor.editor.GetProjectManager().CreateProject(projectName)\n\t\t\t\teditor.statusMessage = fmt.Sprintf(\"Created %s\", projectName)
			case glfw.KeyF7:
				editor.editor.GetProjectManager().SaveCurrentProject()
				editor.statusMessage = "Project Saved"
			case glfw.KeyTab:
				objects := editor.editor.GetSceneObjects()
				if len(objects) > 0 {
					current := editor.editor.GetSelectedObject()
					next := (current + 1) % len(objects)
					editor.editor.SetSelectedObject(next)
					editor.statusMessage = fmt.Sprintf("Selected: %s", objects[next].Name)
				}
			case glfw.KeyDelete:
				objects := editor.editor.GetSceneObjects()
				if len(objects) > 0 {
					selectedIndex := editor.editor.GetSelectedObject()
					if selectedIndex < len(objects) {
						objName := objects[selectedIndex].Name
						editor.editor.DeleteObject(selectedIndex)
						editor.statusMessage = fmt.Sprintf("Deleted: %s", objName)
					}
				}
			case glfw.KeyG:
				grid := editor.editor.GetGrid()
				grid.Visible = !grid.Visible
				if grid.Visible {
					editor.statusMessage = "Grid: ON"
				} else {
					editor.statusMessage = "Grid: OFF"
				}
			}
		}
	})
}

func updateCamera(editor *VisualEditor) {
	// Orbit camera around origin
	x := editor.cameraDistance * float32(math.Cos(float64(editor.cameraAngleX))) * float32(math.Sin(float64(editor.cameraAngleY)))
	y := editor.cameraDistance * float32(math.Sin(float64(editor.cameraAngleX)))
	z := editor.cameraDistance * float32(math.Cos(float64(editor.cameraAngleX))) * float32(math.Cos(float64(editor.cameraAngleY)))
	
	camera := editor.renderer.GetCamera()
	camera.SetPosition(bmath.NewVector3(x, y, z))
	camera.SetTarget(bmath.NewVector3(0, 0, 0))
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

func renderVisualOverlay(renderer *core.Renderer, editor *VisualEditor) {
	if !editor.showMenu {
		return
	}
	
	// Position menu items in 3D space (upper corners of the view)
	camera := renderer.GetCamera()
	cameraPos := camera.GetPosition()
	
	// Calculate positions for menu items (floating in space)
	menuItems := []struct {
		pos   bmath.Vector3
		scale bmath.Vector3
		label string
	}{
		// Top row - Object creation
		{bmath.NewVector3(cameraPos.X-8, cameraPos.Y+4, cameraPos.Z-5), bmath.NewVector3(0.3, 0.3, 0.3), "F1:Cube"},
		{bmath.NewVector3(cameraPos.X-6, cameraPos.Y+4, cameraPos.Z-5), bmath.NewVector3(0.3, 0.3, 0.3), "F2:Sphere"},
		{bmath.NewVector3(cameraPos.X-4, cameraPos.Y+4, cameraPos.Z-5), bmath.NewVector3(0.3, 0.3, 0.3), "F3:Cylinder"},
		{bmath.NewVector3(cameraPos.X-2, cameraPos.Y+4, cameraPos.Z-5), bmath.NewVector3(0.3, 0.3, 0.3), "F4:Plane"},
		{bmath.NewVector3(cameraPos.X+0, cameraPos.Y+4, cameraPos.Z-5), bmath.NewVector3(0.3, 0.3, 0.3), "F5:Triangle"},
		
		// Bottom row - Project actions
		{bmath.NewVector3(cameraPos.X-6, cameraPos.Y+2, cameraPos.Z-5), bmath.NewVector3(0.3, 0.3, 0.3), "F6:NewProj"},
		{bmath.NewVector3(cameraPos.X-4, cameraPos.Y+2, cameraPos.Z-5), bmath.NewVector3(0.3, 0.3, 0.3), "F7:Save"},
		{bmath.NewVector3(cameraPos.X-2, cameraPos.Y+2, cameraPos.Z-5), bmath.NewVector3(0.3, 0.3, 0.3), "G:Grid"},
		{bmath.NewVector3(cameraPos.X+0, cameraPos.Y+2, cameraPos.Z-5), bmath.NewVector3(0.3, 0.3, 0.3), "Tab:Select"},
	}
	
	// Render menu items as small cubes
	for _, item := range menuItems {
		translation := bmath.NewTranslationMatrix(item.pos.X, item.pos.Y, item.pos.Z)
		scale := bmath.NewScaleMatrix(item.scale.X, item.scale.Y, item.scale.Z)
		model := translation.Multiply(scale)
		renderer.DrawCubeWithTransform(model)
	}
	
	// Render status indicator (a small colored cube that shows current status)
	statusPos := bmath.NewVector3(cameraPos.X+6, cameraPos.Y-3, cameraPos.Z-3)
	statusTranslation := bmath.NewTranslationMatrix(statusPos.X, statusPos.Y, statusPos.Z)
	statusScale := bmath.NewScaleMatrix(0.5, 0.5, 0.5)
	statusModel := statusTranslation.Multiply(statusScale)
	renderer.DrawCubeWithTransform(statusModel)
}