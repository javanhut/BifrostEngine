package main

import (
	"fmt"
	"log"

	bmath "github.com/javanhut/BifrostEngine/m/v2/math"
	"github.com/javanhut/BifrostEngine/m/v2/renderer/core"
)

func main() {
	renderer, err := core.New(800, 600, "Camera Step Test")
	if err != nil {
		log.Fatal(err)
	}
	defer renderer.Cleanup()

	mode := 0
	frameCount := 0
	
	fmt.Println("Press window close button to cycle through modes:")
	fmt.Println("Mode 0: No camera (identity matrices)")
	fmt.Println("Mode 1: Camera at (0,0,2) looking at origin")
	fmt.Println("Mode 2: Camera at (0,0,5) looking at origin")
	
	for !renderer.ShouldClose() {
		renderer.BeginFrame()
		
		switch mode {
		case 0:
			// Draw without camera transformations
			renderer.DrawTriangleNoCamera()
		case 1:
			// Draw with camera closer
			camera := renderer.GetCamera()
			camera.SetPosition(bmath.NewVector3(0, 0, 2))
			camera.SetTarget(bmath.NewVector3(0, 0, 0))
			renderer.DrawTriangle()
		case 2:
			// Draw with camera farther
			camera := renderer.GetCamera()
			camera.SetPosition(bmath.NewVector3(0, 0, 5))
			camera.SetTarget(bmath.NewVector3(0, 0, 0))
			renderer.DrawTriangle()
		}
		
		renderer.EndFrame()
		
		frameCount++
		if frameCount % 180 == 0 { // Change mode every 3 seconds at 60fps
			mode = (mode + 1) % 3
			fmt.Printf("Switched to mode %d\n", mode)
		}
	}
}