package main

import (
	"fmt"
	"log"

	"github.com/javanhut/BifrostEngine/m/v2/renderer/core"
)

func main() {
	renderer, err := core.New(800, 600, "Bifrost Engine - 3D Cube Demo")
	if err != nil {
		log.Fatal(err)
	}
	defer renderer.Cleanup()

	fmt.Println("3D Cube Demo:")
	fmt.Println("- Cube rotates on multiple axes")
	fmt.Println("- Different colored faces show true 3D rotation")
	fmt.Println("- Red=Front, Green=Back, Blue=Left, Yellow=Right, Magenta=Top, Cyan=Bottom")

	time := float32(0)
	
	for !renderer.ShouldClose() {
		renderer.BeginFrame()
		
		// Draw rotating cube
		time += 0.016 // ~60 FPS
		renderer.DrawRotatingCube(time)
		
		renderer.EndFrame()
	}
}