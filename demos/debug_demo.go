package main

import (
	"fmt"
	"log"

	"github.com/javanhut/BifrostEngine/m/v2/engine"
	bmath "github.com/javanhut/BifrostEngine/m/v2/math"
	"github.com/javanhut/BifrostEngine/m/v2/scene"
)

func main() {
	// Create engine with configuration
	eng, err := engine.NewEngine(engine.Config{
		WindowWidth:  800,
		WindowHeight: 600,
		WindowTitle:  "Debug Demo",
		TargetFPS:    60,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Get managers
	sceneManager := eng.GetSceneManager()
	
	// Get active scene
	activeScene := sceneManager.GetActiveScene()
	
	fmt.Println("Creating entities...")
	
	// Create a simple cube at origin where camera looks
	cube := activeScene.CreateEntity("Debug Cube")
	cube.Transform.SetPosition(bmath.NewVector3(0, 0, 0))  // At origin where camera looks
	cube.AddComponent(scene.NewMeshComponent("cube"))
	fmt.Printf("Created cube at position: %v\n", cube.Transform.Position)
	
	// Create camera entity
	cameraEntity := activeScene.CreateEntity("Main Camera")
	cameraEntity.Transform.SetPosition(bmath.NewVector3(0, 0, 3))  // Camera further back
	cameraComp := scene.NewCameraComponent(45, 0.1, 100, 4.0/3.0)
	cameraComp.Active = true
	cameraEntity.AddComponent(cameraComp)
	fmt.Printf("Created camera at position: %v\n", cameraEntity.Transform.Position)
	
	// Add camera system
	activeScene.AddSystem(engine.NewCameraSystem(eng.GetRenderer()))
	
	// Print systems
	entities := activeScene.GetEntities()
	fmt.Printf("Total entities: %d\n", len(entities))
	for _, entity := range entities {
		fmt.Printf("Entity: %s, Active: %t, Position: %v\n", 
			entity.Name, entity.Active, entity.Transform.Position)
		
		meshComp := entity.GetComponent("Mesh")
		if meshComp != nil {
			mesh := meshComp.(*scene.MeshComponent)
			fmt.Printf("  - Has mesh component: %s, Visible: %t\n", mesh.MeshType, mesh.Visible)
		}
		
		camComp := entity.GetComponent("Camera")
		if camComp != nil {
			cam := camComp.(*scene.CameraComponent)
			fmt.Printf("  - Has camera component: FOV: %f, Active: %t\n", cam.FOV, cam.Active)
		}
	}
	
	fmt.Println("Starting engine...")
	
	// Run the engine for a short time
	go func() {
		for i := 0; i < 120; i++ { // Run for 2 seconds at 60fps
			eng.GetSceneManager().Update(1.0/60.0)
		}
		eng.Stop()
	}()
	
	eng.Run()
	fmt.Println("Engine stopped")
}