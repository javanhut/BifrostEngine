package engine

import (
	"fmt"
	"time"

	"github.com/javanhut/BifrostEngine/m/v2/input"
	"github.com/javanhut/BifrostEngine/m/v2/renderer/core"
	"github.com/javanhut/BifrostEngine/m/v2/scene"
)

// Engine is the core of Bifrost Engine
type Engine struct {
	renderer      *core.Renderer
	inputManager  *input.InputManager
	sceneManager  *scene.SceneManager
	renderSystem  *RenderSystem
	running       bool
	targetFPS     int
	currentFPS    float32
	deltaTime     float32
}

// Config holds engine configuration
type Config struct {
	WindowWidth  int
	WindowHeight int
	WindowTitle  string
	TargetFPS    int
}

// NewEngine creates a new engine instance
func NewEngine(config Config) (*Engine, error) {
	// Create renderer
	renderer, err := core.New(config.WindowWidth, config.WindowHeight, config.WindowTitle)
	if err != nil {
		return nil, fmt.Errorf("failed to create renderer: %w", err)
	}
	
	// Create input manager
	inputManager := input.NewInputManager(renderer.GetWindow().GetHandle())
	
	// Create scene manager
	sceneManager := scene.NewSceneManager()
	
	// Create engine
	engine := &Engine{
		renderer:     renderer,
		inputManager: inputManager,
		sceneManager: sceneManager,
		renderSystem: NewRenderSystem(renderer),
		running:      false,
		targetFPS:    config.TargetFPS,
	}
	
	// Create default scene
	defaultScene := sceneManager.CreateScene("default")
	sceneManager.SetActiveScene("default")
	
	// Add default systems
	defaultScene.AddSystem(engine.renderSystem)
	defaultScene.AddSystem(NewScriptSystem())
	
	return engine, nil
}

// Run starts the engine main loop
func (e *Engine) Run() {
	e.running = true
	
	targetFrameTime := time.Second / time.Duration(e.targetFPS)
	lastTime := time.Now()
	frameStart := time.Now()
	
	for e.running && !e.renderer.ShouldClose() {
		frameStart = time.Now()
		
		// Calculate delta time
		currentTime := time.Now()
		e.deltaTime = float32(currentTime.Sub(lastTime).Seconds())
		lastTime = currentTime
		
		// Update
		e.update(e.deltaTime)
		
		// Render
		e.render()
		
		// Frame rate limiting
		frameTime := time.Since(frameStart)
		if frameTime < targetFrameTime {
			time.Sleep(targetFrameTime - frameTime)
		}
		
		// Calculate FPS
		e.currentFPS = 1.0 / float32(time.Since(frameStart).Seconds())
	}
	
	e.cleanup()
}

// Stop stops the engine
func (e *Engine) Stop() {
	e.running = false
}

// update handles all updates
func (e *Engine) update(deltaTime float32) {
	// Update input
	e.inputManager.Update()
	
	// Update non-render systems only
	activeScene := e.sceneManager.GetActiveScene()
	if activeScene != nil {
		systems := activeScene.GetSystems()
		for _, system := range systems {
			// Skip render system in update phase
			if system.GetName() != "RenderSystem" {
				system.Update(activeScene, deltaTime)
			}
		}
		
		// Update entity components
		entities := activeScene.GetEntities()
		for _, entity := range entities {
			if entity.Active {
				for _, component := range entity.Components {
					component.Update(deltaTime)
				}
			}
		}
	}
}

// render handles rendering
func (e *Engine) render() {
	e.renderer.BeginFrame()
	
	// Now run render system
	activeScene := e.sceneManager.GetActiveScene()
	if activeScene != nil {
		systems := activeScene.GetSystems()
		for _, system := range systems {
			if system.GetName() == "RenderSystem" {
				system.Update(activeScene, 0) // deltaTime not needed for rendering
			}
		}
	}
	
	e.renderer.EndFrame()
}

// cleanup cleans up resources
func (e *Engine) cleanup() {
	e.renderer.Cleanup()
}

// GetRenderer returns the renderer
func (e *Engine) GetRenderer() *core.Renderer {
	return e.renderer
}

// GetInputManager returns the input manager
func (e *Engine) GetInputManager() *input.InputManager {
	return e.inputManager
}

// GetSceneManager returns the scene manager
func (e *Engine) GetSceneManager() *scene.SceneManager {
	return e.sceneManager
}

// GetFPS returns the current FPS
func (e *Engine) GetFPS() float32 {
	return e.currentFPS
}

// GetDeltaTime returns the frame delta time
func (e *Engine) GetDeltaTime() float32 {
	return e.deltaTime
}