package main

import (
	"fmt"
	"log"

	"github.com/javanhut/BifrostEngine/m/v2/engine"
	"github.com/javanhut/BifrostEngine/m/v2/input"
	bmath "github.com/javanhut/BifrostEngine/m/v2/math"
	"github.com/javanhut/BifrostEngine/m/v2/scene"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func main() {
	// Create engine with configuration
	eng, err := engine.NewEngine(engine.Config{
		WindowWidth:  1200,
		WindowHeight: 800,
		WindowTitle:  "Bifrost Engine - Full Demo",
		TargetFPS:    60,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Get managers
	sceneManager := eng.GetSceneManager()
	inputManager := eng.GetInputManager()
	
	// Get active scene
	activeScene := sceneManager.GetActiveScene()
	
	// Create entities at different positions to test visibility
	createCube(activeScene, "Cube 1", bmath.NewVector3(0, 0, 0))
	createCube(activeScene, "Cube 2", bmath.NewVector3(2, 0, 0))
	createTriangle(activeScene, "Triangle 1", bmath.NewVector3(-2, 0, 0))
	
	// Create camera entity - positioned to look at cubes  
	cameraEntity := activeScene.CreateEntity("Main Camera")
	cameraEntity.Transform.SetPosition(bmath.NewVector3(0, 1, 4))  // Closer and slightly above
	cameraComp := scene.NewCameraComponent(45, 0.1, 100, 1.5)
	cameraComp.Active = true
	cameraEntity.AddComponent(cameraComp)
	
	// Add camera system
	activeScene.AddSystem(engine.NewCameraSystem(eng.GetRenderer()))
	
	// Create a rotating cube with script
	rotatingCube := createCube(activeScene, "Rotating Cube", bmath.NewVector3(0, 2, 0))
	script := scene.NewScriptComponent("Rotator")
	script.OnUpdate = func(entity *scene.Entity, deltaTime float32) {
		entity.Transform.Rotate(bmath.NewVector3(0, 60*deltaTime, 30*deltaTime))
	}
	rotatingCube.AddComponent(script)
	
	// Setup input callbacks
	var selectedEntity *scene.Entity
	selectedIndex := 0
	entities := activeScene.GetEntities()
	
	inputManager.SetCallbacks(&input.InputCallbacks{
		OnKeyPress: func(key glfw.Key) {
			switch key {
			case glfw.KeyTab:
				// Cycle through entities
				selectedIndex = (selectedIndex + 1) % len(entities)
				selectedEntity = entities[selectedIndex]
				fmt.Printf("Selected: %s\n", selectedEntity.Name)
				
			case glfw.Key1:
				// Add new cube at origin
				cube := createCube(activeScene, fmt.Sprintf("Cube %d", len(entities)), bmath.NewVector3(0, 0, 0))
				entities = append(entities, cube)
				
			case glfw.Key2:
				// Add new triangle at origin
				tri := createTriangle(activeScene, fmt.Sprintf("Triangle %d", len(entities)), bmath.NewVector3(0, 0, 0))
				entities = append(entities, tri)
			}
		},
		OnScroll: func(xOffset, yOffset float64) {
			// Zoom camera
			if cameraEntity != nil {
				pos := cameraEntity.Transform.Position
				pos.Z -= float32(yOffset) * 0.5
				cameraEntity.Transform.SetPosition(pos)
			}
		},
	})
	
	// Add real-time controls
	go func() {
		for {
			// Movement controls for selected entity
			if selectedEntity != nil && selectedEntity.Name != "Main Camera" {
				if inputManager.IsKeyHeld(glfw.KeyW) {
					selectedEntity.Transform.Translate(bmath.NewVector3(0, 0.05, 0))
				}
				if inputManager.IsKeyHeld(glfw.KeyS) {
					selectedEntity.Transform.Translate(bmath.NewVector3(0, -0.05, 0))
				}
				if inputManager.IsKeyHeld(glfw.KeyA) {
					selectedEntity.Transform.Translate(bmath.NewVector3(-0.05, 0, 0))
				}
				if inputManager.IsKeyHeld(glfw.KeyD) {
					selectedEntity.Transform.Translate(bmath.NewVector3(0.05, 0, 0))
				}
				if inputManager.IsKeyHeld(glfw.KeyQ) {
					selectedEntity.Transform.Rotate(bmath.NewVector3(0, -2, 0))
				}
				if inputManager.IsKeyHeld(glfw.KeyE) {
					selectedEntity.Transform.Rotate(bmath.NewVector3(0, 2, 0))
				}
			}
			
			// Camera controls
			if inputManager.IsKeyHeld(glfw.KeyLeft) {
				cameraEntity.Transform.Translate(bmath.NewVector3(-0.1, 0, 0))
			}
			if inputManager.IsKeyHeld(glfw.KeyRight) {
				cameraEntity.Transform.Translate(bmath.NewVector3(0.1, 0, 0))
			}
			if inputManager.IsKeyHeld(glfw.KeyUp) {
				cameraEntity.Transform.Translate(bmath.NewVector3(0, 0, -0.1))
			}
			if inputManager.IsKeyHeld(glfw.KeyDown) {
				cameraEntity.Transform.Translate(bmath.NewVector3(0, 0, 0.1))
			}
		}
	}()
	
	// Print controls
	fmt.Println("=== Bifrost Engine Demo ===")
	fmt.Println("Controls:")
	fmt.Println("  Tab - Cycle through objects")
	fmt.Println("  WASD - Move selected object")
	fmt.Println("  QE - Rotate selected object")
	fmt.Println("  Arrow Keys - Move camera")
	fmt.Println("  Scroll - Zoom camera")
	fmt.Println("  1 - Add cube")
	fmt.Println("  2 - Add triangle")
	
	// Run the engine
	eng.Run()
}

func createCube(activeScene *scene.Scene, name string, position bmath.Vector3) *scene.Entity {
	entity := activeScene.CreateEntity(name)
	entity.Transform.SetPosition(position)
	entity.AddComponent(scene.NewMeshComponent("cube"))
	return entity
}

func createTriangle(activeScene *scene.Scene, name string, position bmath.Vector3) *scene.Entity {
	entity := activeScene.CreateEntity(name)
	entity.Transform.SetPosition(position)
	entity.AddComponent(scene.NewMeshComponent("triangle"))
	return entity
}