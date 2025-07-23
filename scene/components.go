package scene

// MeshComponent represents a renderable mesh
type MeshComponent struct {
	MeshType string // "cube", "triangle", "sphere", etc.
	Visible  bool
	Color    [3]float32
}

// NewMeshComponent creates a new mesh component
func NewMeshComponent(meshType string) *MeshComponent {
	return &MeshComponent{
		MeshType: meshType,
		Visible:  true,
		Color:    [3]float32{1.0, 1.0, 1.0},
	}
}

// GetType returns the component type
func (m *MeshComponent) GetType() string {
	return "Mesh"
}

// Update updates the mesh component
func (m *MeshComponent) Update(deltaTime float32) {
	// Mesh components don't need per-frame updates
}

// SetColor sets the mesh color
func (m *MeshComponent) SetColor(r, g, b float32) {
	m.Color[0] = r
	m.Color[1] = g
	m.Color[2] = b
}

// CameraComponent represents a camera attached to an entity
type CameraComponent struct {
	FOV         float32
	NearPlane   float32
	FarPlane    float32
	AspectRatio float32
	Active      bool
}

// NewCameraComponent creates a new camera component
func NewCameraComponent(fov, nearPlane, farPlane, aspectRatio float32) *CameraComponent {
	return &CameraComponent{
		FOV:         fov,
		NearPlane:   nearPlane,
		FarPlane:    farPlane,
		AspectRatio: aspectRatio,
		Active:      false,
	}
}

// GetType returns the component type
func (c *CameraComponent) GetType() string {
	return "Camera"
}

// Update updates the camera component
func (c *CameraComponent) Update(deltaTime float32) {
	// Camera updates handled by render system
}

// LightComponent represents a light source
type LightComponent struct {
	Type      string // "directional", "point", "spot"
	Color     [3]float32
	Intensity float32
	Range     float32 // For point/spot lights
}

// NewLightComponent creates a new light component
func NewLightComponent(lightType string) *LightComponent {
	return &LightComponent{
		Type:      lightType,
		Color:     [3]float32{1.0, 1.0, 1.0},
		Intensity: 1.0,
		Range:     10.0,
	}
}

// GetType returns the component type
func (l *LightComponent) GetType() string {
	return "Light"
}

// Update updates the light component
func (l *LightComponent) Update(deltaTime float32) {
	// Light updates handled by render system
}

// ScriptComponent allows custom behavior
type ScriptComponent struct {
	Name     string
	OnStart  func(entity *Entity)
	OnUpdate func(entity *Entity, deltaTime float32)
	Started  bool
}

// NewScriptComponent creates a new script component
func NewScriptComponent(name string) *ScriptComponent {
	return &ScriptComponent{
		Name:    name,
		Started: false,
	}
}

// GetType returns the component type
func (s *ScriptComponent) GetType() string {
	return "Script"
}

// Update runs the script update function
func (s *ScriptComponent) Update(deltaTime float32) {
	// Script execution handled by script system
}