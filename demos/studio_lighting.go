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
	renderer, err := core.New(1280, 720, "Bifrost Engine - Studio Lighting Demo")
	if err != nil {
		log.Fatal(err)
	}
	defer renderer.Cleanup()

	lightingSystem := renderer.GetLightingSystem()
	
	fmt.Println("=== Studio Lighting Demo ===")
	fmt.Println("Professional 3-point lighting setup with color temperature control")
	fmt.Println()
	
	// Disable sun for studio environment
	lightingSystem.SetSunIntensity(0.0)
	lightingSystem.SetAmbientLight([3]float32{0.05, 0.05, 0.08}, 0.05) // Very minimal ambient
	
	// Professional 3-point lighting setup
	
	// 1. KEY LIGHT (main light, warm white, front-left)
	keyLightIndex := lightingSystem.AddPointLight(
		bmath.NewVector3(-2.0, 3.0, 2.0),
		[3]float32{1.0, 0.9, 0.7}, // Warm white (3200K)
		2.5, // High intensity
	)
	
	// 2. FILL LIGHT (softer, cooler, front-right)
	fillLightIndex := lightingSystem.AddPointLight(
		bmath.NewVector3(2.0, 2.0, 1.5),
		[3]float32{0.8, 0.9, 1.0}, // Cool white (5600K)
		1.2, // Medium intensity
	)
	
	// 3. RIM LIGHT (back light for separation, slightly warm)
	rimLightIndex := lightingSystem.AddPointLight(
		bmath.NewVector3(0.0, 2.5, -2.5),
		[3]float32{1.0, 0.8, 0.6}, // Slightly warm
		1.8, // Medium-high intensity
	)
	
	// 4. BACKGROUND LIGHT (colored accent)
	backgroundLightIndex := lightingSystem.AddPointLight(
		bmath.NewVector3(-1.0, 1.0, -3.0),
		[3]float32{0.4, 0.7, 1.0}, // Blue accent
		0.8, // Lower intensity
	)
	
	fmt.Println("Studio Setup:")
	fmt.Println("  🔆 Key Light (warm, front-left) - Main illumination")
	fmt.Println("  💡 Fill Light (cool, front-right) - Softens shadows")
	fmt.Println("  ✨ Rim Light (warm, back) - Edge separation")
	fmt.Println("  🎨 Background Light (blue, back-left) - Colored accent")
	fmt.Println()
	
	lightingPresets := []struct {
		name string
		description string
		keyColor [3]float32
		fillColor [3]float32
		rimColor [3]float32
		backgroundColor [3]float32
	}{
		{
			"Classic Studio", "Warm key, cool fill",
			[3]float32{1.0, 0.9, 0.7},   // Warm key
			[3]float32{0.8, 0.9, 1.0},   // Cool fill
			[3]float32{1.0, 0.8, 0.6},   // Warm rim
			[3]float32{0.4, 0.7, 1.0},   // Blue background
		},
		{
			"Golden Hour", "Warm sunset lighting",
			[3]float32{1.0, 0.7, 0.4},   // Orange key
			[3]float32{1.0, 0.8, 0.6},   // Warm fill
			[3]float32{1.0, 0.6, 0.3},   // Deep orange rim
			[3]float32{0.8, 0.4, 0.2},   // Red background
		},
		{
			"Cool Corporate", "Professional blue tones",
			[3]float32{0.9, 0.95, 1.0},  // Cool white key
			[3]float32{0.8, 0.9, 1.0},   // Cool fill
			[3]float32{0.7, 0.8, 1.0},   // Cool rim
			[3]float32{0.3, 0.5, 0.8},   // Deep blue background
		},
		{
			"Neon Night", "Cyberpunk colors",
			[3]float32{1.0, 0.2, 0.8},   // Magenta key
			[3]float32{0.2, 1.0, 0.8},   // Cyan fill
			[3]float32{0.8, 0.2, 1.0},   // Purple rim
			[3]float32{0.2, 0.8, 0.2},   // Green background
		},
		{
			"Fire & Ice", "Contrasting warm/cool",
			[3]float32{1.0, 0.4, 0.1},   // Orange key
			[3]float32{0.1, 0.4, 1.0},   // Blue fill
			[3]float32{1.0, 0.6, 0.0},   // Yellow rim
			[3]float32{0.0, 0.8, 1.0},   // Cyan background
		},
	}
	
	currentPreset := 0
	animateObjects := true
	showIntensity := true
	
	fmt.Println("Controls:")
	fmt.Println("  SPACE - Cycle lighting presets")
	fmt.Println("  A - Toggle object animation")
	fmt.Println("  I - Toggle intensity display")
	fmt.Println("  1-4 - Toggle individual lights")
	fmt.Println("  + - Increase all intensities")
	fmt.Println("  - - Decrease all intensities")
	fmt.Println("  ESC - Exit")
	fmt.Printf("Current preset: %s\n", lightingPresets[currentPreset].name)
	
	globalIntensityMultiplier := float32(1.0)
	
	// Setup input
	window := renderer.GetWindow()
	window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if action == glfw.Press {
			switch key {
			case glfw.KeyEscape:
				w.SetShouldClose(true)
			case glfw.KeySpace:
				currentPreset = (currentPreset + 1) % len(lightingPresets)
				preset := lightingPresets[currentPreset]
				fmt.Printf("Preset: %s - %s\n", preset.name, preset.description)
			case glfw.KeyA:
				animateObjects = !animateObjects
				fmt.Printf("Animation: %s\n", map[bool]string{true: "enabled", false: "disabled"}[animateObjects])
			case glfw.KeyI:
				showIntensity = !showIntensity
				fmt.Printf("Intensity display: %s\n", map[bool]string{true: "enabled", false: "disabled"}[showIntensity])
			case glfw.Key1:
				lightingSystem.TogglePointLight(keyLightIndex)
				fmt.Printf("Key light: %t\n", lightingSystem.PointLights[keyLightIndex].Enabled)
			case glfw.Key2:
				lightingSystem.TogglePointLight(fillLightIndex)
				fmt.Printf("Fill light: %t\n", lightingSystem.PointLights[fillLightIndex].Enabled)
			case glfw.Key3:
				lightingSystem.TogglePointLight(rimLightIndex)
				fmt.Printf("Rim light: %t\n", lightingSystem.PointLights[rimLightIndex].Enabled)
			case glfw.Key4:
				lightingSystem.TogglePointLight(backgroundLightIndex)
				fmt.Printf("Background light: %t\n", lightingSystem.PointLights[backgroundLightIndex].Enabled)
			case glfw.KeyEqual, glfw.KeyKPAdd: // + key
				globalIntensityMultiplier += 0.2
				fmt.Printf("Global intensity: %.1fx\n", globalIntensityMultiplier)
			case glfw.KeyMinus, glfw.KeyKPSubtract: // - key
				globalIntensityMultiplier = float32(math.Max(0.1, float64(globalIntensityMultiplier-0.2)))
				fmt.Printf("Global intensity: %.1fx\n", globalIntensityMultiplier)
			}
		}
	})

	time := float32(0.0)
	lastIntensityPrint := float32(0.0)
	
	// Main loop
	for !renderer.ShouldClose() {
		renderer.BeginFrame()
		
		if animateObjects {
			time += 0.016 // ~60 FPS
		}
		
		// Apply current preset colors and intensities
		preset := lightingPresets[currentPreset]
		
		// Update light colors and intensities
		if keyLightIndex < len(lightingSystem.PointLights) {
			light := &lightingSystem.PointLights[keyLightIndex]
			light.Color = preset.keyColor
			light.Intensity = 2.5 * globalIntensityMultiplier
		}
		
		if fillLightIndex < len(lightingSystem.PointLights) {
			light := &lightingSystem.PointLights[fillLightIndex]
			light.Color = preset.fillColor
			light.Intensity = 1.2 * globalIntensityMultiplier
		}
		
		if rimLightIndex < len(lightingSystem.PointLights) {
			light := &lightingSystem.PointLights[rimLightIndex]
			light.Color = preset.rimColor
			light.Intensity = 1.8 * globalIntensityMultiplier
		}
		
		if backgroundLightIndex < len(lightingSystem.PointLights) {
			light := &lightingSystem.PointLights[backgroundLightIndex]
			light.Color = preset.backgroundColor
			light.Intensity = 0.8 * globalIntensityMultiplier
		}
		
		// Print light intensities periodically
		if showIntensity && time-lastIntensityPrint > 2.0 {
			fmt.Printf("Intensities: Key=%.1f Fill=%.1f Rim=%.1f BG=%.1f\n",
				lightingSystem.PointLights[keyLightIndex].Intensity,
				lightingSystem.PointLights[fillLightIndex].Intensity,
				lightingSystem.PointLights[rimLightIndex].Intensity,
				lightingSystem.PointLights[backgroundLightIndex].Intensity)
			lastIntensityPrint = time
		}
		
		// Render studio scene with rotating objects
		rotation := time * 0.3
		
		// Center sphere (main subject)
		sphereModel := bmath.NewMatrix4Identity()
		if animateObjects {
			// Gentle rotation
			sphereModel[0] = float32(math.Cos(float64(rotation)))
			sphereModel[2] = float32(math.Sin(float64(rotation)))
			sphereModel[8] = -float32(math.Sin(float64(rotation)))
			sphereModel[10] = float32(math.Cos(float64(rotation)))
		}
		sphereModel[13] = 0.2 // Slightly elevated
		renderer.DrawSphereWithLighting(sphereModel, false)
		
		// Left cube
		cubeModel := bmath.NewMatrix4Identity()
		cubeModel[12] = -2.8 // X position
		if animateObjects {
			// Counter-rotate
			cubeModel[0] = float32(math.Cos(float64(-rotation * 0.7)))
			cubeModel[2] = float32(math.Sin(float64(-rotation * 0.7)))
			cubeModel[8] = -float32(math.Sin(float64(-rotation * 0.7)))
			cubeModel[10] = float32(math.Cos(float64(-rotation * 0.7)))
		}
		renderer.DrawCubeWithLighting(cubeModel, false)
		
		// Right cylinder
		cylinderModel := bmath.NewMatrix4Identity()
		cylinderModel[12] = 2.8 // X position
		if animateObjects {
			// Slow Y rotation
			cylinderModel[0] = float32(math.Cos(float64(rotation * 0.5)))
			cylinderModel[2] = float32(math.Sin(float64(rotation * 0.5)))
			cylinderModel[8] = -float32(math.Sin(float64(rotation * 0.5)))
			cylinderModel[10] = float32(math.Cos(float64(rotation * 0.5)))
		}
		renderer.DrawCylinderWithLighting(cylinderModel, false)
		
		// Back pyramid
		pyramidModel := bmath.NewMatrix4Identity()
		pyramidModel[12] = 0.0  // X position
		pyramidModel[14] = -4.0 // Z position (back)
		pyramidModel[0] = 2.0   // X scale
		pyramidModel[5] = 2.0   // Y scale
		pyramidModel[10] = 2.0  // Z scale
		renderer.DrawPyramidWithLighting(pyramidModel, false)
		
		// Studio floor
		floorModel := bmath.NewMatrix4Identity()
		floorModel[13] = -1.2 // Y position (down)
		floorModel[0] = 6.0   // X scale (wide studio floor)
		floorModel[10] = 6.0  // Z scale (deep studio floor)
		renderer.DrawPlaneWithLighting(floorModel, false)
		
		renderer.EndFrame()
	}
}