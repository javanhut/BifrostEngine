package main

import (
	"fmt"
	"log"

	"github.com/javanhut/BifrostEngine/m/v2/renderer/core"
)

func main() {
	renderer, err := core.New(800, 600, "Bifrost Engine")
	if err != nil {
		log.Fatal(err)
	}
	defer renderer.Cleanup()

	fmt.Println("Drawing 3D cube with camera...")
	
	for !renderer.ShouldClose() {
		renderer.BeginFrame()

		// Draw the 3D cube
		renderer.DrawCube()

		renderer.EndFrame()
	}
}

