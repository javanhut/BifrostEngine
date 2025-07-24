package main

import (
	"log"
	"math"

	"github.com/inkyblackness/imgui-go/v4"
	bmath "github.com/javanhut/BifrostEngine/m/v2/math"
	"github.com/javanhut/BifrostEngine/m/v2/renderer/core"
	"github.com/javanhut/BifrostEngine/m/v2/renderer/window"
	"github.com/javanhut/BifrostEngine/m/v2/ui"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type UIEditor struct {
	renderer       *core.Renderer
	editor         *ui.Editor
	imguiContext   *ui.ImGuiContext
	cameraDistance float32
	cameraAngleX   float32
	cameraAngleY   float32
	lastMouseX     float64
	lastMouseY     float64
	mousePressed   bool
}

func main() {
	// Initialize renderer
	renderer, err := core.New(1280, 720, "Bifrost Engine - UI Editor")
	if err != nil {
		log.Fatal(err)
	}
	defer renderer.Cleanup()

	// Initialize ImGui
	imguiContext := ui.NewImGuiContext(renderer.GetWindow().GetHandle())
	if imguiContext == nil {
		log.Fatal("Failed to initialize ImGui")
	}
	defer imguiContext.Destroy()

	// Create editor
	uiEditor := &UIEditor{
		renderer:       renderer,
		editor:         ui.NewEditor(),
		imguiContext:   imguiContext,
		cameraDistance: 10.0,
		cameraAngleX:   0.3,
		cameraAngleY:   0.5,
	}

	// Setup input callbacks
	setupMouseInput(renderer.GetWindow(), uiEditor)
	
	// Main loop
	lastTime := glfw.GetTime()
	
	for !renderer.ShouldClose() {
		// Calculate delta time
		currentTime := glfw.GetTime()
		deltaTime := float32(currentTime - lastTime)
		lastTime = currentTime
		
		// Start ImGui frame
		uiEditor.imguiContext.NewFrame()
		
		// Update editor
		uiEditor.editor.Update(deltaTime)
		
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
		
		// Render UI
		uiEditor.editor.Render()
		
		// Render ImGui
		uiEditor.imguiContext.Render()
		
		// End frame
		renderer.EndFrame()
	}
}

func setupMouseInput(window *window.Window, editor *UIEditor) {
	window.SetMouseButtonCallback(func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
		// Let ImGui handle input first
		io := imgui.CurrentIO()
		if io.WantCaptureMouse() {
			return
		}
		
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
		// Let ImGui handle input first
		io := imgui.CurrentIO()
		if io.WantCaptureMouse() {
			return
		}
		
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
		// Let ImGui handle input first
		io := imgui.CurrentIO()
		if io.WantCaptureMouse() {
			return
		}
		
		editor.cameraDistance -= float32(yoff) * 0.5
		if editor.cameraDistance < 1.0 {
			editor.cameraDistance = 1.0
		}
		if editor.cameraDistance > 50.0 {
			editor.cameraDistance = 50.0
		}
	})
}

func updateCamera(editor *UIEditor) {
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