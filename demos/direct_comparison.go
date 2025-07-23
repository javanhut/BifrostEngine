package main

import (
	"log"
	"time"

	"github.com/javanhut/BifrostEngine/m/v2/renderer/core"
	bmath "github.com/javanhut/BifrostEngine/m/v2/math"
)

func main() {
	// Create renderer directly
	renderer, err := core.New(800, 600, "Direct Comparison")
	if err != nil {
		log.Fatal(err)
	}
	defer renderer.Cleanup()

	startTime := time.Now()
	useTransform := false
	
	for !renderer.ShouldClose() {
		elapsedTime := float32(time.Since(startTime).Seconds())
		
		renderer.BeginFrame()
		
		// Switch between methods every 2 seconds
		if int(elapsedTime)%4 < 2 {
			if !useTransform {
				println("Using DrawCube()")
				useTransform = true
			}
			renderer.DrawCube()
		} else {
			if useTransform {
				println("Using DrawCubeWithTransform(identity)")
				useTransform = false
			}
			identity := bmath.NewMatrix4Identity()
			renderer.DrawCubeWithTransform(identity)
		}
		
		renderer.EndFrame()
		
		// Exit after 8 seconds
		if elapsedTime > 8.0 {
			break
		}
	}
}