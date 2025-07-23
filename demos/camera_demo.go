package main

import (
	"fmt"
	"log"
	"math"

	bmath "github.com/javanhut/BifrostEngine/m/v2/math"
	"github.com/javanhut/BifrostEngine/m/v2/renderer/core"
)

func main() {
	renderer, err := core.New(800, 600, "Bifrost Engine - 3D Camera Demo")
	if err != nil {
		log.Fatal(err)
	}
	defer renderer.Cleanup()

	fmt.Println("3D Camera Demo:")
	fmt.Println("- Camera orbits around the triangle")
	fmt.Println("- Triangle should appear to rotate as camera moves")
	fmt.Println("- Perspective should make it smaller when farther away")

	time := float32(0)
	camera := renderer.GetCamera()
	
	for !renderer.ShouldClose() {
		renderer.BeginFrame()
		
		// Rotate camera around the triangle in 3D
		time += 0.016 // ~60 FPS
		radius := float32(4.0)
		x := radius * float32(math.Sin(float64(time)))
		z := radius * float32(math.Cos(float64(time)))
		y := float32(math.Sin(float64(time * 0.5))) // Add vertical movement
		
		camera.SetPosition(bmath.NewVector3(x, y, z))
		camera.SetTarget(bmath.NewVector3(0, 0, 0))
		
		renderer.DrawTriangle()
		
		renderer.EndFrame()
	}
}