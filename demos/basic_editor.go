package main

import (
	"fmt"
	"log"
	"math"

	bmath "github.com/javanhut/BifrostEngine/m/v2/math"
	"github.com/javanhut/BifrostEngine/m/v2/renderer/core"
	"github.com/javanhut/BifrostEngine/m/v2/renderer/window"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type Object struct {
	Name     string
	Position bmath.Vector3
	Rotation bmath.Vector3
	Scale    bmath.Vector3
	Type     string
	Visible  bool
}

type BasicEditor struct {
	objects        []Object
	selectedObject int
	camera         *core.Renderer
	cameraDistance float32
	cameraAngleX   float32
	cameraAngleY   float32
}

func main() {
	renderer, err := core.New(1200, 800, "Bifrost Engine - Basic Editor")
	if err != nil {
		log.Fatal(err)
	}
	defer renderer.Cleanup()

	editor := &BasicEditor{
		objects: []Object{
			{
				Name:     "Cube 1",
				Position: bmath.NewVector3(0, 0, 0),
				Rotation: bmath.NewVector3(0, 0, 0),
				Scale:    bmath.NewVector3(1, 1, 1),
				Type:     "cube",
				Visible:  true,
			},
			{
				Name:     "Triangle 1",
				Position: bmath.NewVector3(2, 0, 0),
				Rotation: bmath.NewVector3(0, 0, 0),
				Scale:    bmath.NewVector3(1, 1, 1),
				Type:     "triangle",
				Visible:  true,
			},
		},
		selectedObject: 0,
		camera:         renderer,
		cameraDistance: 5.0,
		cameraAngleX:   0.0,
		cameraAngleY:   0.0,
	}

	window := renderer.GetWindow()
	
	// Setup input callbacks
	setupInput(window, editor)
	
	fmt.Println("=== Bifrost Engine Editor ===")
	fmt.Println("Controls:")
	fmt.Println("  Mouse: Orbit camera around scene")
	fmt.Println("  Scroll: Zoom in/out")
	fmt.Println("  Tab: Switch selected object")
	fmt.Println("  WASD: Move selected object")
	fmt.Println("  QE: Rotate selected object")
	fmt.Println("  RF: Scale selected object")
	fmt.Println("  V: Toggle object visibility")
	fmt.Println("  1-6: Add objects (cube, triangle)")
	
	var lastMouseX, lastMouseY float64
	var mousePressed bool
	
	// Mouse callback for camera orbit
	window.SetMouseButtonCallback(func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
		if button == glfw.MouseButtonLeft {
			if action == glfw.Press {
				mousePressed = true
				lastMouseX, lastMouseY = w.GetCursorPos()
			} else if action == glfw.Release {
				mousePressed = false
			}
		}
	})
	
	window.SetCursorPosCallback(func(w *glfw.Window, xpos, ypos float64) {
		if mousePressed {
			deltaX := float32(xpos - lastMouseX)
			deltaY := float32(ypos - lastMouseY)
			
			editor.cameraAngleY += deltaX * 0.01
			editor.cameraAngleX += deltaY * 0.01
			
			// Clamp vertical angle
			if editor.cameraAngleX > 1.5 {
				editor.cameraAngleX = 1.5
			}
			if editor.cameraAngleX < -1.5 {
				editor.cameraAngleX = -1.5
			}
			
			updateCamera(editor)
			
			lastMouseX = xpos
			lastMouseY = ypos
		}
	})
	
	window.SetScrollCallback(func(w *glfw.Window, xoff, yoff float64) {
		editor.cameraDistance -= float32(yoff) * 0.5
		if editor.cameraDistance < 1.0 {
			editor.cameraDistance = 1.0
		}
		if editor.cameraDistance > 20.0 {
			editor.cameraDistance = 20.0
		}
		updateCamera(editor)
	})
	
	// Initial camera setup
	updateCamera(editor)
	
	for !renderer.ShouldClose() {
		renderer.BeginFrame()
		
		// Print current state
		printEditorState(editor)
		
		// Render all objects
		for i, obj := range editor.objects {
			if !obj.Visible {
				continue
			}
			
			// Create transform matrix
			translation := bmath.NewTranslationMatrix(obj.Position.X, obj.Position.Y, obj.Position.Z)
			rotationX := bmath.NewRotationX(bmath.Radians(obj.Rotation.X))
			rotationY := bmath.NewRotationY(bmath.Radians(obj.Rotation.Y))
			rotationZ := bmath.NewRotationZ(bmath.Radians(obj.Rotation.Z))
			scale := bmath.NewScaleMatrix(obj.Scale.X, obj.Scale.Y, obj.Scale.Z)
			
			rotation := rotationZ.Multiply(rotationY).Multiply(rotationX)
			model := translation.Multiply(rotation).Multiply(scale)
			
			// Highlight selected object (make it slightly bigger)
			if i == editor.selectedObject {
				highlightScale := bmath.NewScaleMatrix(1.1, 1.1, 1.1)
				model = model.Multiply(highlightScale)
			}
			
			// Render based on type
			switch obj.Type {
			case "cube":
				renderer.DrawCubeWithTransform(model)
			case "triangle":
				renderer.DrawTriangleWithTransform(model)
			}
		}
		
		renderer.EndFrame()
	}
}

func setupInput(window *window.Window, editor *BasicEditor) {
	window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if action == glfw.Press || action == glfw.Repeat {
			switch key {
			case glfw.KeyTab:
				editor.selectedObject = (editor.selectedObject + 1) % len(editor.objects)
				
			// Movement
			case glfw.KeyW:
				editor.objects[editor.selectedObject].Position.Y += 0.1
			case glfw.KeyS:
				editor.objects[editor.selectedObject].Position.Y -= 0.1
			case glfw.KeyA:
				editor.objects[editor.selectedObject].Position.X -= 0.1
			case glfw.KeyD:
				editor.objects[editor.selectedObject].Position.X += 0.1
				
			// Rotation
			case glfw.KeyQ:
				editor.objects[editor.selectedObject].Rotation.Y += 15
			case glfw.KeyE:
				editor.objects[editor.selectedObject].Rotation.Y -= 15
				
			// Scale
			case glfw.KeyR:
				obj := &editor.objects[editor.selectedObject]
				obj.Scale = obj.Scale.Mul(1.1)
			case glfw.KeyF:
				obj := &editor.objects[editor.selectedObject]
				obj.Scale = obj.Scale.Mul(0.9)
				
			// Visibility
			case glfw.KeyV:
				editor.objects[editor.selectedObject].Visible = !editor.objects[editor.selectedObject].Visible
				
			// Add objects
			case glfw.Key1:
				addObject(editor, "cube")
			case glfw.Key2:
				addObject(editor, "triangle")
			}
		}
	})
}

func addObject(editor *BasicEditor, objectType string) {
	name := fmt.Sprintf("%s %d", objectType, len(editor.objects)+1)
	newObj := Object{
		Name:     name,
		Position: bmath.NewVector3(0, 0, 0),
		Rotation: bmath.NewVector3(0, 0, 0),
		Scale:    bmath.NewVector3(1, 1, 1),
		Type:     objectType,
		Visible:  true,
	}
	editor.objects = append(editor.objects, newObj)
}

func updateCamera(editor *BasicEditor) {
	// Orbit camera around origin
	x := editor.cameraDistance * float32(math.Cos(float64(editor.cameraAngleX))) * float32(math.Sin(float64(editor.cameraAngleY)))
	y := editor.cameraDistance * float32(math.Sin(float64(editor.cameraAngleX)))
	z := editor.cameraDistance * float32(math.Cos(float64(editor.cameraAngleX))) * float32(math.Cos(float64(editor.cameraAngleY)))
	
	camera := editor.camera.GetCamera()
	camera.SetPosition(bmath.NewVector3(x, y, z))
	camera.SetTarget(bmath.NewVector3(0, 0, 0))
}

var frameCount = 0

func printEditorState(editor *BasicEditor) {
	// Print state every 60 frames (roughly once per second)
	frameCount++
	if frameCount%60 == 0 {
		obj := editor.objects[editor.selectedObject]
		fmt.Printf("\rSelected: %s | Pos: (%.1f,%.1f,%.1f) | Rot: (%.0f,%.0f,%.0f) | Scale: (%.1f,%.1f,%.1f) | Visible: %t",
			obj.Name, obj.Position.X, obj.Position.Y, obj.Position.Z,
			obj.Rotation.X, obj.Rotation.Y, obj.Rotation.Z,
			obj.Scale.X, obj.Scale.Y, obj.Scale.Z, obj.Visible)
	}
}