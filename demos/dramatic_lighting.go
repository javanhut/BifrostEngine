package main

import (
	"fmt"
	"log"
	"math"

	bmath "github.com/javanhut/BifrostEngine/m/v2/math"
	"github.com/javanhut/BifrostEngine/m/v2/renderer/core"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func main() {
	renderer, err := core.New(1280, 720, "Bifrost Engine - Dramatic Lighting Demo")
	if err != nil {
		log.Fatal(err)
	}
	defer renderer.Cleanup()

	lightingSystem := renderer.GetLightingSystem()
	
	fmt.Println("=== Dramatic Lighting Demo ===")
	fmt.Println("Multiple colored lights creating dramatic effects!")
	fmt.Println()
	
	// Setup dramatic lighting scene
	// 1. Dim the sun for indoor/studio lighting feel
	lightingSystem.SetSunIntensity(0.2)
	lightingSystem.SetSunColor(0.9, 0.9, 1.0) // Cool white
	lightingSystem.SetAmbientLight([3]float32{0.1, 0.1, 0.15}, 0.1) // Very dim blue ambient
	
	// 2. Add multiple colored point lights
	// Red spotlight (left)
	redLightIndex := lightingSystem.AddPointLight(
		bmath.NewVector3(-3.0, 2.0, 1.0),
		[3]float32{1.0, 0.2, 0.2}, // Intense red
		2.0, // High intensity
	)
	
	// Blue spotlight (right)  
	blueLightIndex := lightingSystem.AddPointLight(
		bmath.NewVector3(3.0, 2.0, 1.0),
		[3]float32{0.2, 0.4, 1.0}, // Intense blue
		2.0, // High intensity
	)
	
	// Green accent light (back)
	greenLightIndex := lightingSystem.AddPointLight(
		bmath.NewVector3(0.0, 1.0, -2.0),
		[3]float32{0.2, 1.0, 0.4}, // Bright green
		1.5, // Medium intensity
	)
	
	// Purple rim light (above)
	purpleLightIndex := lightingSystem.AddPointLight(
		bmath.NewVector3(0.0, 4.0, 0.0),
		[3]float32{0.8, 0.2, 1.0}, // Purple
		1.0, // Medium intensity
	)
	
	fmt.Printf("Created %d point lights:\n", len(lightingSystem.PointLights))
	fmt.Println("  🔴 Red spotlight (left)")
	fmt.Println("  🔵 Blue spotlight (right)")
	fmt.Println("  🟢 Green accent (back)")
	fmt.Println("  🟣 Purple rim light (above)")
	fmt.Println()
	
	lightModes := []string{"All Lights", "Red Only", "Blue Only", "Green Only", "Purple Only", "No Point Lights"}
	currentMode := 0
	
	fmt.Println("Controls:")
	fmt.Println("  SPACE - Cycle through lighting modes")
	fmt.Println("  A - Toggle light animation")
	fmt.Println("  + - Increase intensity")
	fmt.Println("  - - Decrease intensity")
	fmt.Println("  ESC - Exit")
	fmt.Printf("Current mode: %s\n", lightModes[currentMode])
	
	animate := true
	globalIntensity := float32(1.0)
	
	// Setup input
	window := renderer.GetWindow()
	window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if action == glfw.Press {
			switch key {
			case glfw.KeyEscape:
				w.SetShouldClose(true)
			case glfw.KeySpace:
				currentMode = (currentMode + 1) % len(lightModes)
				fmt.Printf("Lighting mode: %s\n", lightModes[currentMode])
			case glfw.KeyA:
				animate = !animate
				fmt.Printf("Animation: %s\n", map[bool]string{true: "enabled", false: "disabled"}[animate])
			case glfw.KeyEqual, glfw.KeyKPAdd: // + key
				globalIntensity += 0.2
				fmt.Printf("Intensity: %.1f\n", globalIntensity)
			case glfw.KeyMinus, glfw.KeyKPSubtract: // - key
				globalIntensity = float32(math.Max(0.1, float64(globalIntensity-0.2)))
				fmt.Printf("Intensity: %.1f\n", globalIntensity)
			}
		}
	})

	time := float32(0.0)
	
	// Main loop
	for !renderer.ShouldClose() {
		renderer.BeginFrame()
		
		if animate {
			time += 0.016 // ~60 FPS
		}
		
		// Update light positions and enable/disable based on mode
		for i, light := range lightingSystem.PointLights {
			switch i {
			case redLightIndex:
				// Red light orbits horizontally
				light.Position.X = -3.0 + 0.5*float32(math.Sin(float64(time)))
				light.Position.Y = 2.0 + 0.3*float32(math.Cos(float64(time*1.5)))
				light.Intensity = 2.0 * globalIntensity
				light.Enabled = (currentMode == 0 || currentMode == 1)
				
			case blueLightIndex:
				// Blue light orbits opposite to red
				light.Position.X = 3.0 + 0.5*float32(math.Sin(float64(time+math.Pi)))
				light.Position.Y = 2.0 + 0.3*float32(math.Cos(float64(time*1.5+math.Pi)))
				light.Intensity = 2.0 * globalIntensity
				light.Enabled = (currentMode == 0 || currentMode == 2)
				
			case greenLightIndex:
				// Green light moves vertically
				light.Position.Y = 1.0 + 0.8*float32(math.Sin(float64(time*0.7)))
				light.Position.Z = -2.0 + 0.4*float32(math.Cos(float64(time*0.7)))
				light.Intensity = 1.5 * globalIntensity
				light.Enabled = (currentMode == 0 || currentMode == 3)
				
			case purpleLightIndex:
				// Purple light pulses
				intensityPulse := 0.5 + 0.5*float32(math.Sin(float64(time*2.0)))
				light.Intensity = (1.0 + intensityPulse) * globalIntensity
				light.Enabled = (currentMode == 0 || currentMode == 4)
			}
		}
		
		// Disable all point lights for "No Point Lights" mode
		if currentMode == 5 {
			for i := range lightingSystem.PointLights {
				lightingSystem.PointLights[i].Enabled = false
			}
		}
		
		// Render dramatic scene
		// Central sphere (main subject)
		sphereModel := bmath.NewMatrix4Identity()
		sphereModel[13] = 0.5 // Slightly elevated
		renderer.DrawSphereWithLighting(sphereModel, false)
		
		// Left cube
		cubeModel := bmath.NewMatrix4Identity()
		cubeModel[12] = -2.5 // X position
		cubeModel[13] = 0.0  // Y position
		renderer.DrawCubeWithLighting(cubeModel, false)
		
		// Right cylinder
		cylinderModel := bmath.NewMatrix4Identity()
		cylinderModel[12] = 2.5 // X position
		cylinderModel[13] = 0.0 // Y position
		renderer.DrawCylinderWithLighting(cylinderModel, false)
		
		// Background pyramid
		pyramidModel := bmath.NewMatrix4Identity()
		pyramidModel[12] = 0.0  // X position
		pyramidModel[13] = 0.0  // Y position
		pyramidModel[14] = -3.0 // Z position (back)
		pyramidModel[0] = 1.5   // X scale
		pyramidModel[5] = 1.5   // Y scale  
		pyramidModel[10] = 1.5  // Z scale
		renderer.DrawPyramidWithLighting(pyramidModel, false)
		
		// Ground plane
		groundModel := bmath.NewMatrix4Identity()
		groundModel[13] = -1.0 // Y position (down)
		groundModel[0] = 4.0   // X scale (wider)
		groundModel[10] = 4.0  // Z scale (deeper)
		renderer.DrawPlaneWithLighting(groundModel, false)
		
		renderer.EndFrame()
	}
}