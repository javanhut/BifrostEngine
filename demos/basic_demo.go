package main

import (
	"log"

	"github.com/javanhut/BifrostEngine/m/v2/engine"
	bmath "github.com/javanhut/BifrostEngine/m/v2/math"
	"github.com/javanhut/BifrostEngine/m/v2/scene"
)

func main() {
	// Create engine
	eng, err := engine.NewEngine(engine.Config{
		WindowWidth:  800,
		WindowHeight: 600,
		WindowTitle:  "Basic Test",
		TargetFPS:    60,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Get active scene - DON'T add camera system
	activeScene := eng.GetSceneManager().GetActiveScene()
	
	// Create a simple cube without camera entity
	cube := activeScene.CreateEntity("Test Cube")
	cube.Transform.SetPosition(bmath.NewVector3(0, 0, 0))
	cube.AddComponent(scene.NewMeshComponent("cube"))
	
	// Run for a short time
	go func() {
		for i := 0; i < 180; i++ { // 3 seconds
			eng.GetSceneManager().Update(1.0/60.0)
		}
		eng.Stop()
	}()
	
	eng.Run()
}