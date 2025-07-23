package main

import (
	"log"
	"math"

	bmath "github.com/javanhut/BifrostEngine/m/v2/math"
	"github.com/javanhut/BifrostEngine/m/v2/renderer/core"
	"github.com/javanhut/BifrostEngine/m/v2/ui"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func main() {
	renderer, err := core.New(1200, 800, "Bifrost Engine - Editor")
	if err != nil {
		log.Fatal(err)
	}
	defer renderer.Cleanup()

	// Initialize ImGui
	window := renderer.GetWindow()
	glfwWindow := window.GetHandle() // We need to add this method
	
	editor := ui.NewEditor()
	
	// Main loop
	var deltaTime float32 = 0.016
	lastTime := glfw.GetTime()
	
	for !renderer.ShouldClose() {
		// Calculate delta time
		currentTime := glfw.GetTime()
		deltaTime = float32(currentTime - lastTime)
		lastTime = currentTime
		
		renderer.BeginFrame()
		
		// Update editor
		editor.Update(deltaTime)
		
		// Render 3D scene
		renderScene(renderer, editor)
		
		// Render UI (this would go here when ImGui is working)
		editor.Render()
		
		renderer.EndFrame()
	}
}

func renderScene(renderer *core.Renderer, editor *ui.Editor) {
	// Get scene objects from editor
	objects := editor.GetSceneObjects()
	
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
		
		// Combine transformations: T * R * S
		rotation := rotationZ.Multiply(rotationY).Multiply(rotationX)
		model := translation.Multiply(rotation).Multiply(scale)
		
		// Render based on object type
		switch obj.Type {
		case "cube":
			renderer.DrawCubeWithTransform(model)
		case "triangle":
			renderer.DrawTriangleWithTransform(model)
		}
		
		// Highlight selected object
		if i == editor.GetSelectedObject() {
			// Draw wireframe or outline for selected object
			// This would be implemented later
		}
	}
}