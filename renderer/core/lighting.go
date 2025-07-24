package core

import (
	bmath "github.com/javanhut/BifrostEngine/m/v2/math"
)

// DirectionalLight represents a directional light source (like the sun)
type DirectionalLight struct {
	Direction bmath.Vector3 // Direction the light is pointing
	Color     [3]float32    // RGB color of the light
	Intensity float32       // Light intensity (0.0 to 1.0+)
	Enabled   bool         // Whether this light is active
}

// PointLight represents a point light source
type PointLight struct {
	Position    bmath.Vector3 // World position of the light
	Color       [3]float32    // RGB color of the light
	Intensity   float32       // Light intensity
	Attenuation [3]float32    // Constant, linear, quadratic attenuation
	Enabled     bool         // Whether this light is active
}

// AmbientLight represents global ambient lighting
type AmbientLight struct {
	Color     [3]float32 // RGB color of ambient light
	Intensity float32    // Ambient intensity (0.0 to 1.0)
}

// LightingSystem manages all lights in the scene
type LightingSystem struct {
	DirectionalLights []DirectionalLight
	PointLights      []PointLight
	AmbientLight     AmbientLight
	MaxDirectionalLights int
	MaxPointLights       int
}

// NewLightingSystem creates a new lighting system with default settings
func NewLightingSystem() *LightingSystem {
	ls := &LightingSystem{
		DirectionalLights: make([]DirectionalLight, 0),
		PointLights:      make([]PointLight, 0),
		MaxDirectionalLights: 4,
		MaxPointLights:       8,
	}
	
	// Set default ambient lighting
	ls.AmbientLight = AmbientLight{
		Color:     [3]float32{0.2, 0.2, 0.2}, // Soft gray
		Intensity: 0.3,
	}
	
	// Add default sun light
	ls.AddDirectionalLight(
		bmath.NewVector3(-0.5, -1.0, -0.3).Normalize(), // Coming from upper right
		[3]float32{1.0, 0.95, 0.8}, // Warm white/yellow
		0.8,
	)
	
	return ls
}

// AddDirectionalLight adds a new directional light
func (ls *LightingSystem) AddDirectionalLight(direction bmath.Vector3, color [3]float32, intensity float32) int {
	if len(ls.DirectionalLights) >= ls.MaxDirectionalLights {
		return -1 // Max lights reached
	}
	
	light := DirectionalLight{
		Direction: direction.Normalize(),
		Color:     color,
		Intensity: intensity,
		Enabled:   true,
	}
	
	ls.DirectionalLights = append(ls.DirectionalLights, light)
	return len(ls.DirectionalLights) - 1
}

// AddPointLight adds a new point light
func (ls *LightingSystem) AddPointLight(position bmath.Vector3, color [3]float32, intensity float32) int {
	if len(ls.PointLights) >= ls.MaxPointLights {
		return -1 // Max lights reached
	}
	
	light := PointLight{
		Position:    position,
		Color:       color,
		Intensity:   intensity,
		Attenuation: [3]float32{1.0, 0.09, 0.032}, // Constant, linear, quadratic
		Enabled:     true,
	}
	
	ls.PointLights = append(ls.PointLights, light)
	return len(ls.PointLights) - 1
}

// SetAmbientLight updates the ambient lighting
func (ls *LightingSystem) SetAmbientLight(color [3]float32, intensity float32) {
	ls.AmbientLight.Color = color
	ls.AmbientLight.Intensity = intensity  
}

// GetDirectionalLight returns a directional light by index
func (ls *LightingSystem) GetDirectionalLight(index int) *DirectionalLight {
	if index >= 0 && index < len(ls.DirectionalLights) {
		return &ls.DirectionalLights[index]
	}
	return nil
}

// GetPointLight returns a point light by index
func (ls *LightingSystem) GetPointLight(index int) *PointLight {
	if index >= 0 && index < len(ls.PointLights) {
		return &ls.PointLights[index]
	}
	return nil
}

// ToggleDirectionalLight enables/disables a directional light
func (ls *LightingSystem) ToggleDirectionalLight(index int) {
	if light := ls.GetDirectionalLight(index); light != nil {
		light.Enabled = !light.Enabled
	}
}

// TogglePointLight enables/disables a point light
func (ls *LightingSystem) TogglePointLight(index int) {
	if light := ls.GetPointLight(index); light != nil {
		light.Enabled = !light.Enabled
	}
}

// SetSunDirection updates the main directional light direction
func (ls *LightingSystem) SetSunDirection(direction bmath.Vector3) {
	if len(ls.DirectionalLights) > 0 {
		ls.DirectionalLights[0].Direction = direction.Normalize()
	}
}

// SetSunColor updates the main directional light color
func (ls *LightingSystem) SetSunColor(r, g, b float32) {
	if len(ls.DirectionalLights) > 0 {
		ls.DirectionalLights[0].Color = [3]float32{r, g, b}
	}
}

// SetSunIntensity updates the main directional light intensity
func (ls *LightingSystem) SetSunIntensity(intensity float32) {
	if len(ls.DirectionalLights) > 0 {
		ls.DirectionalLights[0].Intensity = intensity
	}
}

// GetSunLight returns the main directional light (sun)
func (ls *LightingSystem) GetSunLight() *DirectionalLight {
	if len(ls.DirectionalLights) > 0 {
		return &ls.DirectionalLights[0]
	}
	return nil
}