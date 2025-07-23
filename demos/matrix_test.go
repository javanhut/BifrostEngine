package main

import (
	"fmt"
	"log"

	bmath "github.com/javanhut/BifrostEngine/m/v2/math"
	"github.com/javanhut/BifrostEngine/m/v2/renderer/core"
)

func main() {
	renderer, err := core.New(800, 600, "Matrix Test")
	if err != nil {
		log.Fatal(err)
	}
	defer renderer.Cleanup()

	fmt.Println("Testing matrices step by step...")
	
	// Print some debug info
	camera := renderer.GetCamera()
	
	view := camera.GetViewMatrix()
	projection := camera.GetProjectionMatrix()
	
	fmt.Printf("View matrix: %v\n", view)
	fmt.Printf("Projection matrix: %v\n", projection)
	
	for !renderer.ShouldClose() {
		renderer.BeginFrame()
		
		// Just draw with view matrix (what works)
		renderer.DrawTriangleViewOnly()
		
		renderer.EndFrame()
	}
}