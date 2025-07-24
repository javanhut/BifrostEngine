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
	renderer, err := core.New(1280, 720, "Bifrost Engine - Day/Night Cycle Demo")
	if err != nil {
		log.Fatal(err)
	}
	defer renderer.Cleanup()

	lightingSystem := renderer.GetLightingSystem()
	
	fmt.Println("=== Day/Night Cycle Demo ===")
	fmt.Println("Watch the sun move across the sky!")
	fmt.Println("Controls:")
	fmt.Println("  SPACE - Pause/Resume cycle")
	fmt.Println("  UP/DOWN - Speed up/slow down")
	fmt.Println("  R - Reset to noon")
	fmt.Println("  ESC - Exit")
	
	// Cycle parameters
	timeOfDay := float32(0.0) // 0 = midnight, 0.5 = noon, 1.0 = midnight again
	cycleSpeed := float32(0.2) // Speed multiplier
	paused := false
	
	// Setup input
	window := renderer.GetWindow()
	window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if action == glfw.Press {
			switch key {
			case glfw.KeyEscape:
				w.SetShouldClose(true)
			case glfw.KeySpace:
				paused = !paused
				fmt.Printf("Cycle %s\n", map[bool]string{true: "paused", false: "resumed"}[paused])
			case glfw.KeyUp:
				cycleSpeed += 0.1
				fmt.Printf("Speed: %.1fx\n", cycleSpeed)
			case glfw.KeyDown:
				cycleSpeed = float32(math.Max(0.1, float64(cycleSpeed-0.1)))
				fmt.Printf("Speed: %.1fx\n", cycleSpeed)
			case glfw.KeyR:
				timeOfDay = 0.5 // Reset to noon
				fmt.Println("Reset to noon")
			}
		}
	})

	// Create objects for lighting demonstration
	fmt.Println("\nRendering scene with:")
	fmt.Println("- Cube (left)")
	fmt.Println("- Sphere (center)")
	fmt.Println("- Cylinder (right)")
	fmt.Println("- Plane (ground)")

	// Main loop
	for !renderer.ShouldClose() {
		renderer.BeginFrame()
		
		// Update time of day
		if !paused {
			timeOfDay += 0.016 * cycleSpeed / 60.0 // 60 seconds = full cycle at 1x speed
			if timeOfDay > 1.0 {
				timeOfDay -= 1.0
			}
		}
		
		// Calculate sun position (moves in arc across sky)
		sunAngle := timeOfDay * 2.0 * math.Pi
		sunHeight := float32(math.Sin(float64(sunAngle))) // -1 to 1
		sunSide := float32(math.Cos(float64(sunAngle)))   // -1 to 1
		
		// Sun direction (normalized vector pointing toward sun)
		sunDirection := bmath.NewVector3(
			-sunSide,     // East (-1) to West (+1)
			-sunHeight,   // Down (-1) to Up (+1), negated because light direction points FROM sun
			-0.3,         // Slightly toward viewer
		)
		
		// Update sun properties based on time of day
		var sunColor [3]float32
		var sunIntensity float32
		var ambientIntensity float32
		
		if sunHeight > 0 { // Daytime
			// Golden hour near horizon, white at zenith
			horizonFactor := 1.0 - sunHeight // 1.0 at horizon, 0.0 at zenith
			sunColor = [3]float32{
				1.0,                              // Always full red
				0.7 + 0.3*(1.0-horizonFactor),   // More yellow at horizon
				0.4 + 0.6*(1.0-horizonFactor),   // More blue at zenith
			}
			sunIntensity = 0.3 + 0.7*sunHeight // Dimmer near horizon
			ambientIntensity = 0.1 + 0.4*sunHeight
		} else { // Nighttime
			// Moonlight (cool blue)
			moonHeight := -sunHeight // 0 to 1 during night
			sunColor = [3]float32{0.4, 0.6, 1.0} // Cool blue
			sunIntensity = 0.15 * moonHeight       // Very dim
			ambientIntensity = 0.05 * moonHeight
		}
		
		// Apply lighting changes
		lightingSystem.SetSunDirection(sunDirection)
		lightingSystem.SetSunColor(sunColor[0], sunColor[1], sunColor[2])
		lightingSystem.SetSunIntensity(sunIntensity)
		lightingSystem.SetAmbientLight([3]float32{0.2, 0.2, 0.3}, ambientIntensity)
		
		// Print current time info (every 60 frames)
		frameCount := int(timeOfDay * 60 * 60 / (0.016 * cycleSpeed))
		if frameCount%60 == 0 {
			timeStr := ""
			if sunHeight > 0.8 {
				timeStr = "High Noon"
			} else if sunHeight > 0.3 {
				timeStr = "Daytime" 
			} else if sunHeight > 0 {
				timeStr = "Golden Hour"
			} else if sunHeight > -0.3 {
				timeStr = "Twilight"
			} else {
				timeStr = "Night"
			}
			fmt.Printf("Time: %s | Sun: (%.2f, %.2f) | Intensity: %.2f\n", 
				timeStr, sunSide, sunHeight, sunIntensity)
		}
		
		// Render objects
		// Ground plane
		groundModel := bmath.NewMatrix4Identity()
		groundModel[13] = -1.0 // Y position (down)
		groundModel[0] = 3.0   // X scale (wider)
		groundModel[10] = 3.0  // Z scale (deeper)
		renderer.DrawPlaneWithLighting(groundModel, false)
		
		// Cube (left)
		cubeModel := bmath.NewMatrix4Identity()
		cubeModel[12] = -2.0 // X position
		renderer.DrawCubeWithLighting(cubeModel, false)
		
		// Sphere (center)
		sphereModel := bmath.NewMatrix4Identity()
		sphereModel[12] = 0.0 // X position
		renderer.DrawSphereWithLighting(sphereModel, false)
		
		// Cylinder (right)
		cylinderModel := bmath.NewMatrix4Identity()
		cylinderModel[12] = 2.0 // X position
		renderer.DrawCylinderWithLighting(cylinderModel, false)
		
		renderer.EndFrame()
	}
}