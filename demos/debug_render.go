package main

import (
	"fmt"
	"log"

	"github.com/javanhut/BifrostEngine/m/v2/renderer/core"
)

func main() {
	renderer, err := core.New(800, 600, "Bifrost Engine - Debug Test")
	if err != nil {
		log.Fatal(err)
	}
	defer renderer.Cleanup()

	fmt.Println("Renderer initialized successfully")
	fmt.Println("Starting render loop...")
	
	frameCount := 0
	for !renderer.ShouldClose() {
		renderer.BeginFrame()
		
		// Draw with debug info
		renderer.DrawTriangle()
		
		renderer.EndFrame()
		
		frameCount++
		if frameCount == 1 {
			fmt.Println("First frame rendered")
		}
		if frameCount % 60 == 0 {
			fmt.Printf("Frame %d\n", frameCount)
		}
	}
}