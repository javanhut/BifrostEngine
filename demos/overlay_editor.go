package main

import (
	"fmt"
	"log"
	"math"

	"github.com/go-gl/gl/v4.1-core/gl"
	bmath "github.com/javanhut/BifrostEngine/m/v2/math"
	"github.com/javanhut/BifrostEngine/m/v2/renderer/core"
	"github.com/javanhut/BifrostEngine/m/v2/renderer/window"
	"github.com/javanhut/BifrostEngine/m/v2/ui"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type OverlayEditor struct {
	renderer       *core.Renderer
	editor         *ui.Editor
	overlay        *ui.Overlay
	cameraDistance float32
	cameraAngleX   float32
	cameraAngleY   float32
	lastMouseX     float64
	lastMouseY     float64
	mousePressed   bool
}

func main() {
	// Initialize renderer
	renderer, err := core.New(1280, 720, "Bifrost Engine - Overlay Editor")
	if err != nil {
		log.Fatal(err)
	}
	defer renderer.Cleanup()

	// Create editor
	overlayEditor := &OverlayEditor{
		renderer:       renderer,
		editor:         ui.NewEditor(),
		overlay:        ui.NewOverlay(),
		cameraDistance: 10.0,
		cameraAngleX:   0.3,
		cameraAngleY:   0.5,
	}

	// Setup input callbacks
	setupMouseInput(renderer.GetWindow(), overlayEditor)
	setupKeyboardInput(renderer.GetWindow(), overlayEditor)
	
	fmt.Println("=== Bifrost Engine Overlay Editor ===")
	fmt.Println("Press M to toggle on-screen menu overlay")
	fmt.Println("Use function keys F1-F7 for quick actions")
	
	// Main loop
	for !renderer.ShouldClose() {
		// Update camera
		updateCamera(overlayEditor)
		
		// Begin rendering
		renderer.BeginFrame()
		
		// Render grid if enabled
		if grid := overlayEditor.editor.GetGrid(); grid.Visible {
			lines := grid.GetLines()
			renderer.DrawGrid(lines, grid.Color)
		}
		
		// Render 3D scene
		renderScene(renderer, overlayEditor.editor)
		
		// Render overlay
		renderOverlay(overlayEditor)
		
		// End frame
		renderer.EndFrame()
	}
}

func setupMouseInput(window *window.Window, editor *OverlayEditor) {
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

func setupKeyboardInput(window *window.Window, editor *OverlayEditor) {
	window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if action == glfw.Press {
			switch key {
			case glfw.KeyM:
				editor.overlay.ToggleMenu()
			case glfw.Key1, glfw.KeyF1:
				editor.editor.AddObjectWithType("cube")
			case glfw.Key2, glfw.KeyF2:
				editor.editor.AddObjectWithType("sphere")
			case glfw.Key3, glfw.KeyF3:
				editor.editor.AddObjectWithType("cylinder")
			case glfw.Key4, glfw.KeyF4:
				editor.editor.AddObjectWithType("plane")
			case glfw.Key5, glfw.KeyF5:
				editor.editor.AddObjectWithType("triangle")
			case glfw.KeyF6:
				editor.editor.GetProjectManager().CreateProject(fmt.Sprintf("Project_%d", len(editor.editor.GetProjectManager().GetProjects())+1))
			case glfw.KeyF7:
				editor.editor.GetProjectManager().SaveCurrentProject()
			case glfw.KeyTab:
				objects := editor.editor.GetSceneObjects()
				if len(objects) > 0 {
					current := editor.editor.GetSelectedObject()
					next := (current + 1) % len(objects)
					editor.editor.SetSelectedObject(next)
				}
			case glfw.KeyDelete:
				editor.editor.DeleteObject(editor.editor.GetSelectedObject())
			case glfw.KeyG:
				grid := editor.editor.GetGrid()
				grid.Visible = !grid.Visible
			}
		}
	})
}

func updateCamera(editor *OverlayEditor) {
	// Orbit camera around origin
	x := editor.cameraDistance * float32(math.Cos(float64(editor.cameraAngleX))) * float32(math.Sin(float64(editor.cameraAngleY)))
	y := editor.cameraDistance * float32(math.Sin(float64(editor.cameraAngleX)))
	z := editor.cameraDistance * float32(math.Cos(float64(editor.cameraAngleX))) * float32(math.Cos(float64(editor.cameraAngleY)))
	
	camera := editor.renderer.GetCamera()
	camera.SetPosition(bmath.NewVector3(x, y, z))
	camera.SetTarget(bmath.NewVector3(0, 0, 0))
	
	// Update status
	objects := editor.editor.GetSceneObjects()
	statusText := fmt.Sprintf("Objects: %d", len(objects))
	if len(objects) > 0 && editor.editor.GetSelectedObject() < len(objects) {
		obj := objects[editor.editor.GetSelectedObject()]
		statusText = fmt.Sprintf("Selected: %s | Objects: %d | Grid: %v", 
			obj.Name, len(objects), editor.editor.GetGrid().Visible)
	}
	editor.overlay.SetStatus(statusText)
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

func renderOverlay(editor *OverlayEditor) {
	// Save OpenGL state
	gl.Disable(gl.DEPTH_TEST)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	
	// Get window size
	width, height := editor.renderer.GetWindow().GetSize()
	
	// Set up 2D projection
	gl.MatrixMode(gl.PROJECTION)
	gl.PushMatrix()
	gl.LoadIdentity()
	gl.Ortho(0, float64(width), float64(height), 0, -1, 1)
	
	gl.MatrixMode(gl.MODELVIEW)
	gl.PushMatrix()
	gl.LoadIdentity()
	
	// Always show status bar at bottom
	renderStatusBar(editor, width, height)
	
	// Show menu if enabled
	if editor.overlay.ShouldShowMenu() {
		renderMenu(editor, width, height)
	}
	
	// Restore matrices
	gl.PopMatrix()
	gl.MatrixMode(gl.PROJECTION)
	gl.PopMatrix()
	gl.MatrixMode(gl.MODELVIEW)
	
	// Restore OpenGL state
	gl.Enable(gl.DEPTH_TEST)
	gl.Disable(gl.BLEND)
}

func renderStatusBar(editor *OverlayEditor, width, height int) {
	// Draw semi-transparent background for status bar
	gl.Color4f(0.0, 0.0, 0.0, 0.7)
	gl.Begin(gl.QUADS)
	gl.Vertex2f(0, float32(height-30))
	gl.Vertex2f(float32(width), float32(height-30))
	gl.Vertex2f(float32(width), float32(height))
	gl.Vertex2f(0, float32(height))
	gl.End()
	
	// Note: In a real implementation, you'd render text here
	// For now, we just show the background bar
}

func renderMenu(editor *OverlayEditor, width, height int) {
	// Calculate menu dimensions
	menuWidth := 400
	menuHeight := 300
	x := (width - menuWidth) / 2
	y := (height - menuHeight) / 2
	
	// Draw semi-transparent background
	gl.Color4f(0.0, 0.0, 0.0, 0.8)
	gl.Begin(gl.QUADS)
	gl.Vertex2f(float32(x), float32(y))
	gl.Vertex2f(float32(x+menuWidth), float32(y))
	gl.Vertex2f(float32(x+menuWidth), float32(y+menuHeight))
	gl.Vertex2f(float32(x), float32(y+menuHeight))
	gl.End()
	
	// Draw border
	gl.Color4f(0.5, 0.5, 0.5, 1.0)
	gl.Begin(gl.LINE_LOOP)
	gl.Vertex2f(float32(x), float32(y))
	gl.Vertex2f(float32(x+menuWidth), float32(y))
	gl.Vertex2f(float32(x+menuWidth), float32(y+menuHeight))
	gl.Vertex2f(float32(x), float32(y+menuHeight))
	gl.End()
	
	// Draw title bar
	gl.Color4f(0.2, 0.4, 0.8, 1.0)
	gl.Begin(gl.QUADS)
	gl.Vertex2f(float32(x), float32(y))
	gl.Vertex2f(float32(x+menuWidth), float32(y))
	gl.Vertex2f(float32(x+menuWidth), float32(y+30))
	gl.Vertex2f(float32(x), float32(y+30))
	gl.End()
	
	// Draw menu items as colored rectangles (representing different options)
	menuItems := []struct{
		text string
		x, y, w, h int
		r, g, b float32
	}{
		{"F1: Cube", x+20, y+50, 80, 25, 0.8, 0.2, 0.2},
		{"F2: Sphere", x+120, y+50, 80, 25, 0.2, 0.8, 0.2},
		{"F3: Cylinder", x+220, y+50, 80, 25, 0.2, 0.2, 0.8},
		{"F4: Plane", x+320, y+50, 60, 25, 0.8, 0.8, 0.2},
		{"F5: Triangle", x+20, y+90, 80, 25, 0.8, 0.2, 0.8},
		{"F6: New Project", x+120, y+90, 100, 25, 0.2, 0.8, 0.8},
		{"F7: Save", x+240, y+90, 60, 25, 0.5, 0.5, 0.5},
		{"G: Grid", x+20, y+130, 60, 25, 0.6, 0.6, 0.2},
		{"Tab: Select", x+100, y+130, 80, 25, 0.2, 0.6, 0.6},
		{"Del: Remove", x+200, y+130, 80, 25, 0.6, 0.2, 0.2},
		{"M: Hide Menu", x+150, y+200, 100, 30, 0.4, 0.4, 0.4},
	}
	
	for _, item := range menuItems {
		// Draw button background
		gl.Color4f(item.r, item.g, item.b, 0.8)
		gl.Begin(gl.QUADS)
		gl.Vertex2f(float32(item.x), float32(item.y))
		gl.Vertex2f(float32(item.x+item.w), float32(item.y))
		gl.Vertex2f(float32(item.x+item.w), float32(item.y+item.h))
		gl.Vertex2f(float32(item.x), float32(item.y+item.h))
		gl.End()
		
		// Draw button border
		gl.Color4f(1.0, 1.0, 1.0, 0.5)
		gl.Begin(gl.LINE_LOOP)
		gl.Vertex2f(float32(item.x), float32(item.y))
		gl.Vertex2f(float32(item.x+item.w), float32(item.y))
		gl.Vertex2f(float32(item.x+item.w), float32(item.y+item.h))
		gl.Vertex2f(float32(item.x), float32(item.y+item.h))
		gl.End()
	}
}