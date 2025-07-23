package scene

import (
	"fmt"
	"sync"
)

// SceneManager manages all scenes in the engine
type SceneManager struct {
	scenes       map[string]*Scene
	activeScene  *Scene
	mu           sync.RWMutex
}

// Scene represents a game scene with entities
type Scene struct {
	name      string
	entities  map[uint64]*Entity
	systems   []System
	active    bool
	nextID    uint64
	mu        sync.RWMutex
}

// Entity represents a game object
type Entity struct {
	ID         uint64
	Name       string
	Active     bool
	Components map[string]Component
	Transform  *Transform
}

// Component interface for all components
type Component interface {
	GetType() string
	Update(deltaTime float32)
}

// System interface for systems that process entities
type System interface {
	Update(scene *Scene, deltaTime float32)
	GetName() string
}

// NewSceneManager creates a new scene manager
func NewSceneManager() *SceneManager {
	return &SceneManager{
		scenes: make(map[string]*Scene),
	}
}

// CreateScene creates a new scene
func (sm *SceneManager) CreateScene(name string) *Scene {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	
	scene := &Scene{
		name:     name,
		entities: make(map[uint64]*Entity),
		systems:  []System{},
		active:   false,
		nextID:   1,
	}
	
	sm.scenes[name] = scene
	return scene
}

// SetActiveScene sets the active scene
func (sm *SceneManager) SetActiveScene(name string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	
	scene, exists := sm.scenes[name]
	if !exists {
		return fmt.Errorf("scene '%s' not found", name)
	}
	
	if sm.activeScene != nil {
		sm.activeScene.active = false
	}
	
	sm.activeScene = scene
	scene.active = true
	return nil
}

// GetActiveScene returns the current active scene
func (sm *SceneManager) GetActiveScene() *Scene {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.activeScene
}

// Update updates the active scene
func (sm *SceneManager) Update(deltaTime float32) {
	if sm.activeScene != nil {
		sm.activeScene.Update(deltaTime)
	}
}

// CreateEntity creates a new entity in the scene
func (s *Scene) CreateEntity(name string) *Entity {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	entity := &Entity{
		ID:         s.nextID,
		Name:       name,
		Active:     true,
		Components: make(map[string]Component),
		Transform:  NewTransform(),
	}
	
	s.entities[s.nextID] = entity
	s.nextID++
	
	return entity
}

// RemoveEntity removes an entity from the scene
func (s *Scene) RemoveEntity(id uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.entities, id)
}

// GetEntity retrieves an entity by ID
func (s *Scene) GetEntity(id uint64) *Entity {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.entities[id]
}

// GetEntities returns all entities in the scene
func (s *Scene) GetEntities() []*Entity {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	entities := make([]*Entity, 0, len(s.entities))
	for _, entity := range s.entities {
		entities = append(entities, entity)
	}
	return entities
}

// AddSystem adds a system to the scene
func (s *Scene) AddSystem(system System) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.systems = append(s.systems, system)
}

// GetSystems returns all systems in the scene
func (s *Scene) GetSystems() []System {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.systems
}

// Update updates all systems and entities in the scene
func (s *Scene) Update(deltaTime float32) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	// Update all systems
	for _, system := range s.systems {
		system.Update(s, deltaTime)
	}
	
	// Update all entity components
	for _, entity := range s.entities {
		if entity.Active {
			for _, component := range entity.Components {
				component.Update(deltaTime)
			}
		}
	}
}

// AddComponent adds a component to the entity
func (e *Entity) AddComponent(component Component) {
	e.Components[component.GetType()] = component
}

// GetComponent retrieves a component by type
func (e *Entity) GetComponent(componentType string) Component {
	return e.Components[componentType]
}

// HasComponent checks if entity has a component
func (e *Entity) HasComponent(componentType string) bool {
	_, exists := e.Components[componentType]
	return exists
}