package engine

import (
	bmath "github.com/javanhut/BifrostEngine/m/v2/math"
	"github.com/javanhut/BifrostEngine/m/v2/renderer/core"
	. "github.com/javanhut/BifrostEngine/m/v2/scene"
)

// RenderSystem handles rendering of entities with mesh components
type RenderSystem struct {
	renderer *core.Renderer
}

// NewRenderSystem creates a new render system
func NewRenderSystem(renderer *core.Renderer) *RenderSystem {
	return &RenderSystem{
		renderer: renderer,
	}
}

// Update processes all entities with mesh components
func (rs *RenderSystem) Update(scene *Scene, deltaTime float32) {
	entities := scene.GetEntities()
	
	for _, entity := range entities {
		if !entity.Active {
			continue
		}
		
		// Check if entity has mesh component
		meshComp := entity.GetComponent("Mesh")
		if meshComp == nil {
			continue
		}
		
		// Type assert to MeshComponent
		mesh, ok := meshComp.(*MeshComponent)
		if !ok {
			continue
		}
		if !mesh.Visible {
			continue
		}
		
		// Get transform
		transform := entity.Transform
		worldMatrix := transform.GetWorldMatrix()
		
		// Render based on mesh type
		switch mesh.MeshType {
		case "cube":
			rs.renderer.DrawCubeWithTransform(worldMatrix)
		case "triangle":
			rs.renderer.DrawTriangleWithTransform(worldMatrix)
		}
	}
}

// GetName returns the system name
func (rs *RenderSystem) GetName() string {
	return "RenderSystem"
}

// ScriptSystem handles script components
type ScriptSystem struct{}

// NewScriptSystem creates a new script system
func NewScriptSystem() *ScriptSystem {
	return &ScriptSystem{}
}

// Update processes all entities with script components
func (ss *ScriptSystem) Update(scene *Scene, deltaTime float32) {
	entities := scene.GetEntities()
	
	for _, entity := range entities {
		if !entity.Active {
			continue
		}
		
		// Check if entity has script component
		scriptComp := entity.GetComponent("Script")
		if scriptComp == nil {
			continue
		}
		
		script, ok := scriptComp.(*ScriptComponent)
		if !ok {
			continue
		}
		
		// Run OnStart if not started
		if !script.Started && script.OnStart != nil {
			script.OnStart(entity)
			script.Started = true
		}
		
		// Run OnUpdate
		if script.OnUpdate != nil {
			script.OnUpdate(entity, deltaTime)
		}
	}
}

// GetName returns the system name
func (ss *ScriptSystem) GetName() string {
	return "ScriptSystem"
}

// CameraSystem handles camera components and updates the renderer's camera
type CameraSystem struct {
	renderer *core.Renderer
}

// NewCameraSystem creates a new camera system
func NewCameraSystem(renderer *core.Renderer) *CameraSystem {
	return &CameraSystem{
		renderer: renderer,
	}
}

// Update finds the active camera and updates the renderer
func (cs *CameraSystem) Update(scene *Scene, deltaTime float32) {
	entities := scene.GetEntities()
	
	for _, entity := range entities {
		if !entity.Active {
			continue
		}
		
		// Check if entity has camera component
		cameraComp := entity.GetComponent("Camera")
		if cameraComp == nil {
			continue
		}
		
		cam, ok := cameraComp.(*CameraComponent)
		if !ok {
			continue
		}
		if !cam.Active {
			continue
		}
		
		// Update renderer camera based on entity transform
		transform := entity.Transform
		worldPos := transform.GetWorldPosition()
		
		// Calculate look direction from rotation
		// This is simplified - in a real engine you'd use the forward vector
		camera := cs.renderer.GetCamera()
		camera.SetPosition(worldPos)
		
		// For now, always look at origin
		// In a real implementation, you'd calculate the target from rotation
		camera.SetTarget(bmath.NewVector3(0, 0, 0))
		
		break // Only use first active camera
	}
}

// GetName returns the system name
func (cs *CameraSystem) GetName() string {
	return "CameraSystem"
}