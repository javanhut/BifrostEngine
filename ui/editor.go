package ui

import (
	"fmt"

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
	cameraPosition  bmath.Vector3
	fps            float32
}

func NewEditor() *Editor {
	return &Editor{
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
			{
				Name:     "Triangle 1", 
				Position: bmath.NewVector3(2, 0, 0),
				Rotation: bmath.NewVector3(0, 0, 0),
				Scale:    bmath.NewVector3(1, 1, 1),
				Color:    [3]float32{1.0, 0.5, 0.2},
				Type:     "triangle",
				Visible:  true,
			},
		},
		selectedObject: 0,
		showViewport:   true,
		showHierarchy:  true,
		showInspector:  true,
		showStats:      true,
		cameraPosition: bmath.NewVector3(0, 0, 3),
	}
}

func (e *Editor) Update(deltaTime float32) {
	e.fps = 1.0 / deltaTime
}

func (e *Editor) Render() {
	// Main menu bar
	if imgui.BeginMainMenuBar() {
		if imgui.BeginMenu("File") {
			if imgui.MenuItem("New Scene") {
				fmt.Println("New Scene")
			}
			if imgui.MenuItem("Open Scene") {
				fmt.Println("Open Scene")
			}
			if imgui.MenuItem("Save Scene") {
				fmt.Println("Save Scene")
			}
			imgui.Separator()
			if imgui.MenuItem("Exit") {
				fmt.Println("Exit requested")
			}
			imgui.EndMenu()
		}
		
		if imgui.BeginMenu("View") {
			imgui.MenuItemBoolPtr("Hierarchy", "", &e.showHierarchy)
			imgui.MenuItemBoolPtr("Inspector", "", &e.showInspector)
			imgui.MenuItemBoolPtr("Stats", "", &e.showStats)
			imgui.MenuItemBoolPtr("Demo Window", "", &e.showDemoWindow)
			imgui.EndMenu()
		}
		
		if imgui.BeginMenu("Objects") {
			if imgui.MenuItem("Add Cube") {
				e.AddObject("cube")
			}
			if imgui.MenuItem("Add Triangle") {
				e.AddObject("triangle")
			}
			imgui.EndMenu()
		}
		
		imgui.EndMainMenuBar()
	}
	
	// Render UI panels
	if e.showHierarchy {
		e.renderHierarchy()
	}
	
	if e.showInspector {
		e.renderInspector()
	}
	
	if e.showStats {
		e.renderStats()
	}
	
	if e.showDemoWindow {
		imgui.ShowDemoWindow(&e.showDemoWindow)
	}
}

func (e *Editor) renderHierarchy() {
	if imgui.Begin("Scene Hierarchy") {
		for i := range e.sceneObjects {
			obj := &e.sceneObjects[i]
			
			flags := imgui.TreeNodeFlagsLeaf | imgui.TreeNodeFlagsNoTreePushOnOpen
			if i == e.selectedObject {
				flags |= imgui.TreeNodeFlagsSelected
			}
			
			imgui.TreeNodeExFlags(obj.Name, flags)
			
			if imgui.IsItemClicked() {
				e.selectedObject = i
			}
			
			// Right-click context menu
			if imgui.BeginPopupContextItem() {
				if imgui.MenuItem("Delete") {
					e.DeleteObject(i)
				}
				if imgui.MenuItem("Duplicate") {
					e.DuplicateObject(i)
				}
				imgui.EndPopup()
			}
		}
	}
	imgui.End()
}

func (e *Editor) renderInspector() {
	if imgui.Begin("Inspector") {
		if e.selectedObject >= 0 && e.selectedObject < len(e.sceneObjects) {
			obj := &e.sceneObjects[e.selectedObject]
			
			imgui.Text("Object: " + obj.Name)
			imgui.Separator()
			
			// Name
			name := obj.Name
			if imgui.InputText("Name", &name) {
				obj.Name = name
			}
			
			// Type
			imgui.Text("Type: " + obj.Type)
			
			// Visibility
			imgui.Checkbox("Visible", &obj.Visible)
			
			imgui.Separator()
			
			// Transform
			if imgui.CollapsingHeader("Transform") {
				pos := [3]float32{obj.Position.X, obj.Position.Y, obj.Position.Z}
				if imgui.DragFloat3("Position", &pos, 0.1) {
					obj.Position = bmath.NewVector3(pos[0], pos[1], pos[2])
				}
				
				rot := [3]float32{obj.Rotation.X, obj.Rotation.Y, obj.Rotation.Z}
				if imgui.DragFloat3("Rotation", &rot, 1.0) {
					obj.Rotation = bmath.NewVector3(rot[0], rot[1], rot[2])
				}
				
				scale := [3]float32{obj.Scale.X, obj.Scale.Y, obj.Scale.Z}
				if imgui.DragFloat3("Scale", &scale, 0.1) {
					obj.Scale = bmath.NewVector3(scale[0], scale[1], scale[2])
				}
			}
			
			// Color
			if imgui.CollapsingHeader("Appearance") {
				imgui.ColorEdit3("Color", &obj.Color)
			}
		} else {
			imgui.Text("No object selected")
		}
	}
	imgui.End()
}

func (e *Editor) renderStats() {
	if imgui.Begin("Stats") {
		imgui.Text(fmt.Sprintf("FPS: %.1f", e.fps))
		imgui.Text(fmt.Sprintf("Objects: %d", len(e.sceneObjects)))
		imgui.Text(fmt.Sprintf("Selected: %d", e.selectedObject))
		
		imgui.Separator()
		
		pos := [3]float32{e.cameraPosition.X, e.cameraPosition.Y, e.cameraPosition.Z}
		imgui.Text("Camera Position:")
		imgui.DragFloat3("##CamPos", &pos, 0.1)
		e.cameraPosition = bmath.NewVector3(pos[0], pos[1], pos[2])
	}
	imgui.End()
}

func (e *Editor) AddObject(objectType string) {
	name := fmt.Sprintf("%s %d", objectType, len(e.sceneObjects)+1)
	
	newObj := SceneObject{
		Name:     name,
		Position: bmath.NewVector3(0, 0, 0),
		Rotation: bmath.NewVector3(0, 0, 0),
		Scale:    bmath.NewVector3(1, 1, 1),
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