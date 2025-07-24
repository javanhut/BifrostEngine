package ui

import (
	"fmt"
	"reflect"

	"github.com/inkyblackness/imgui-go/v4"
	bmath "github.com/javanhut/BifrostEngine/m/v2/math"
)

type SceneObject struct {
	Name     string
	Position bmath.Vector3
	Rotation bmath.Vector3
	Scale    bmath.Vector3
	Color    [3]float32
	Type     string // "cube", "triangle", etc.
	Visible  bool
}

type Editor struct {
	sceneObjects    []SceneObject
	selectedObject  int
	showDemoWindow  bool
	showViewport    bool
	showHierarchy   bool
	showInspector   bool
	showStats       bool
	showToolbar     bool
	showProjectPanel bool
	cameraPosition  bmath.Vector3
	renderer        interface{} // Interface to avoid circular dependency
	fps            float32
	grid           *Grid
	projectManager *ProjectManager
	currentTool    string
}

func NewEditor() *Editor {
	e := &Editor{
		sceneObjects: []SceneObject{
			{
				Name:     "Cube 1",
				Position: bmath.NewVector3(0, 0, 0),
				Rotation: bmath.NewVector3(0, 0, 0),
				Scale:    bmath.NewVector3(1, 1, 1),
				Color:    [3]float32{1.0, 1.0, 1.0},
				Type:     "cube",
				Visible:  true,
			},
		},
		selectedObject: 0,
		showViewport:   true,
		showHierarchy:  true,
		showInspector:  true,
		showStats:      true,
		showToolbar:    true,
		showProjectPanel: false,
		cameraPosition: bmath.NewVector3(0, 0, 3),
		grid:          NewGrid(),
		projectManager: NewProjectManager(),
		currentTool:   "select",
	}
	
	// Create default project
	e.projectManager.CreateProject("Untitled Project")
	return e
}

func (e *Editor) Update(deltaTime float32) {
	e.fps = 1.0 / deltaTime
}

func (e *Editor) Render() {
	// Main menu bar
	if imgui.BeginMainMenuBar() {
		if imgui.BeginMenu("File") {
			if imgui.MenuItem("New Project") {
				e.showProjectPanel = true
			}
			if imgui.MenuItem("Save Project") {
				e.projectManager.SaveCurrentProject()
			}
			imgui.EndMenu()
		}
		
		if imgui.BeginMenu("Add Object") {
			if imgui.MenuItem("Cube") {
				e.AddObjectWithType("cube")
			}
			if imgui.MenuItem("Sphere") {
				e.AddObjectWithType("sphere")
			}
			if imgui.MenuItem("Cylinder") {
				e.AddObjectWithType("cylinder")
			}
			if imgui.MenuItem("Plane") {
				e.AddObjectWithType("plane")
			}
			if imgui.MenuItem("Triangle") {
				e.AddObjectWithType("triangle")
			}
			imgui.EndMenu()
		}
		
		if imgui.BeginMenu("View") {
			if imgui.MenuItem("Toggle Grid") {
				e.grid.Visible = !e.grid.Visible
			}
			imgui.EndMenu()
		}
		
		imgui.EndMainMenuBar()
	}
	
	// Project panel
	if e.showProjectPanel {
		e.renderProjectPanel()
	}
	
	// Status window
	e.renderStatusWindow()
}

func (e *Editor) renderProjectPanel() {
	if imgui.Begin("Project Manager") {
		imgui.Text("Project Management")
		imgui.Separator()
		
		if imgui.Button("New Project") {
			e.projectManager.CreateProject(fmt.Sprintf("Project_%d", len(e.projectManager.GetProjects())+1))
			fmt.Println("Created new project")
		}
		
		if imgui.Button("Close") {
			e.showProjectPanel = false
		}
		
		if project := e.projectManager.GetCurrentProject(); project != nil {
			imgui.Text("Current: " + project.Name)
		}
	}
	imgui.End()
}

func (e *Editor) renderStatusWindow() {
	if imgui.Begin("Status") {
		imgui.Text(fmt.Sprintf("Objects: %d", len(e.sceneObjects)))
		if len(e.sceneObjects) > 0 && e.selectedObject < len(e.sceneObjects) {
			obj := e.sceneObjects[e.selectedObject]
			imgui.Text(fmt.Sprintf("Selected: %s", obj.Name))
			imgui.Text(fmt.Sprintf("Position: (%.1f, %.1f, %.1f)", obj.Position.X, obj.Position.Y, obj.Position.Z))
		}
		imgui.Text(fmt.Sprintf("Grid: %v", e.grid.Visible))
	}
	imgui.End()
}

func (e *Editor) AddObject(objectType string) {
	e.AddObjectWithType(objectType)
}

func (e *Editor) AddObjectWithType(objectType string) {
	name := fmt.Sprintf("%s %d", objectType, len(e.sceneObjects)+1)
	
	// Get default size from template
	defaultScale := bmath.NewVector3(1, 1, 1)
	if template := GetObjectTemplate(ObjectType(objectType)); template != nil {
		defaultScale = template.DefaultSize
	}
	
	newObj := SceneObject{
		Name:     name,
		Position: bmath.NewVector3(0, 0, 0),
		Rotation: bmath.NewVector3(0, 0, 0),
		Scale:    defaultScale,
		Color:    [3]float32{1.0, 1.0, 1.0},
		Type:     objectType,
		Visible:  true,
	}
	
	e.sceneObjects = append(e.sceneObjects, newObj)
}

func (e *Editor) DeleteObject(index int) {
	if index >= 0 && index < len(e.sceneObjects) {
		e.sceneObjects = append(e.sceneObjects[:index], e.sceneObjects[index+1:]...)
		if e.selectedObject >= len(e.sceneObjects) {
			e.selectedObject = len(e.sceneObjects) - 1
		}
		if e.selectedObject < 0 {
			e.selectedObject = 0
		}
	}
}

func (e *Editor) DuplicateObject(index int) {
	if index >= 0 && index < len(e.sceneObjects) {
		original := e.sceneObjects[index]
		duplicate := original
		duplicate.Name += " Copy"
		duplicate.Position.X += 1.0 // Offset slightly
		e.sceneObjects = append(e.sceneObjects, duplicate)
	}
}

func (e *Editor) GetSceneObjects() []SceneObject {
	return e.sceneObjects
}

func (e *Editor) GetSelectedObject() int {
	return e.selectedObject
}

func (e *Editor) SetCameraPosition(pos bmath.Vector3) {
	e.cameraPosition = pos
}

func (e *Editor) GetGrid() *Grid {
	return e.grid
}

func (e *Editor) GetCurrentTool() string {
	return e.currentTool
}

func (e *Editor) SetSelectedObject(index int) {
	if index >= 0 && index < len(e.sceneObjects) {
		e.selectedObject = index
	}
}

func (e *Editor) GetProjectManager() *ProjectManager {
	return e.projectManager
}

func (e *Editor) UpdateObject(index int, obj SceneObject) {
	if index >= 0 && index < len(e.sceneObjects) {
		e.sceneObjects[index] = obj
	}
}

// SetRenderer sets the renderer for the editor (used for asset management)
func (e *Editor) SetRenderer(renderer interface{}) {
	e.renderer = renderer
}

// GetRenderer returns the renderer interface
func (e *Editor) GetRenderer() interface{} {
	return e.renderer
}

// GetAssetStats returns asset statistics safely
func (e *Editor) GetAssetStats() (int, int, int, error) {
	if e.renderer == nil {
		return 0, 0, 0, fmt.Errorf("no renderer available")
	}
	
	// Try multiple interface signatures since we can't know the exact return type
	var stats interface{}
	
	// Try the most general approach - use reflection to call the method
	rendererValue := reflect.ValueOf(e.renderer)
	getStatsMethod := rendererValue.MethodByName("GetAssetStats")
	if getStatsMethod.IsValid() {
		result := getStatsMethod.Call([]reflect.Value{})
		if len(result) > 0 {
			stats = result[0].Interface()
		}
	}
	
	if stats != nil {
		// Use reflection to access fields
		v := reflect.ValueOf(stats)
		if v.Kind() == reflect.Struct {
			loadedMeshes := v.FieldByName("LoadedMeshes")
			totalVertices := v.FieldByName("TotalVertices")
			totalIndices := v.FieldByName("TotalIndices")
			
			if loadedMeshes.IsValid() && totalVertices.IsValid() && totalIndices.IsValid() {
				return int(loadedMeshes.Int()), int(totalVertices.Int()), int(totalIndices.Int()), nil
			}
		}
	}
	
	return 0, 0, 0, fmt.Errorf("could not get asset stats")
}

// GetLoadedMeshNames returns list of loaded mesh names
func (e *Editor) GetLoadedMeshNames() ([]string, error) {
	if e.renderer == nil {
		return nil, fmt.Errorf("no renderer available")
	}
	
	if meshGetter, ok := e.renderer.(interface {
		GetLoadedMeshes() []string
	}); ok {
		return meshGetter.GetLoadedMeshes(), nil
	}
	
	return nil, fmt.Errorf("could not get loaded meshes")
}

// ClearAssetCache clears all loaded assets
func (e *Editor) ClearAssetCache() error {
	if e.renderer == nil {
		return fmt.Errorf("no renderer available")
	}
	
	// Use reflection to call GetAssetManager
	rendererValue := reflect.ValueOf(e.renderer)
	getAssetManagerMethod := rendererValue.MethodByName("GetAssetManager")
	if getAssetManagerMethod.IsValid() {
		result := getAssetManagerMethod.Call([]reflect.Value{})
		if len(result) > 0 {
			assetManager := result[0].Interface()
			
			// Use reflection to call UnloadAll on the asset manager
			amValue := reflect.ValueOf(assetManager)
			unloadAllMethod := amValue.MethodByName("UnloadAll")
			if unloadAllMethod.IsValid() {
				unloadAllMethod.Call([]reflect.Value{})
				return nil
			}
		}
	}
	
	return fmt.Errorf("could not clear asset cache")
}