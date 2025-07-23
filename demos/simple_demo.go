package main

import (
	"log"
	"time"

	"github.com/javanhut/BifrostEngine/m/v2/renderer/core"
)

func main() {
	// Create renderer directly
	renderer, err := core.New(800, 600, "Simple Test")
	if err != nil {
		log.Fatal(err)
	}
	defer renderer.Cleanup()

	startTime := time.Now()
	
	for !renderer.ShouldClose() {
		elapsedTime := float32(time.Since(startTime).Seconds())
		
		renderer.BeginFrame()
		
		// Test basic cube rendering
		renderer.DrawCube()
		
		renderer.EndFrame()
		
		// Exit after 3 seconds
		if elapsedTime > 3.0 {
			break
		}
	}
}