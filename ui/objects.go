package ui

import (
	bmath "github.com/javanhut/BifrostEngine/m/v2/math"
)

type ObjectType string

const (
	ObjectTypeCube     ObjectType = "cube"
	ObjectTypeSphere   ObjectType = "sphere"
	ObjectTypeCylinder ObjectType = "cylinder"
	ObjectTypePlane    ObjectType = "plane"
	ObjectTypeCone     ObjectType = "cone"
	ObjectTypeTorus    ObjectType = "torus"
	ObjectTypeTriangle ObjectType = "triangle"
	ObjectTypeLight    ObjectType = "light"
	ObjectTypeCamera   ObjectType = "camera"
)

type ObjectTemplate struct {
	Type        ObjectType
	DisplayName string
	Icon        string
	Category    string
	DefaultSize bmath.Vector3
}

var ObjectTemplates = []ObjectTemplate{
	{Type: ObjectTypeCube, DisplayName: "Cube", Icon: "cube", Category: "Primitives", DefaultSize: bmath.NewVector3(1, 1, 1)},
	{Type: ObjectTypeSphere, DisplayName: "Sphere", Icon: "sphere", Category: "Primitives", DefaultSize: bmath.NewVector3(1, 1, 1)},
	{Type: ObjectTypeCylinder, DisplayName: "Cylinder", Icon: "cylinder", Category: "Primitives", DefaultSize: bmath.NewVector3(1, 2, 1)},
	{Type: ObjectTypePlane, DisplayName: "Plane", Icon: "plane", Category: "Primitives", DefaultSize: bmath.NewVector3(10, 0.1, 10)},
	{Type: ObjectTypeCone, DisplayName: "Cone", Icon: "cone", Category: "Primitives", DefaultSize: bmath.NewVector3(1, 2, 1)},
	{Type: ObjectTypeTorus, DisplayName: "Torus", Icon: "torus", Category: "Primitives", DefaultSize: bmath.NewVector3(1, 0.3, 1)},
	{Type: ObjectTypeTriangle, DisplayName: "Triangle", Icon: "triangle", Category: "Primitives", DefaultSize: bmath.NewVector3(1, 1, 1)},
	{Type: ObjectTypeLight, DisplayName: "Light", Icon: "light", Category: "Lights", DefaultSize: bmath.NewVector3(0.2, 0.2, 0.2)},
	{Type: ObjectTypeCamera, DisplayName: "Camera", Icon: "camera", Category: "Cameras", DefaultSize: bmath.NewVector3(0.5, 0.5, 0.5)},
}

func GetObjectTemplate(objectType ObjectType) *ObjectTemplate {
	for _, template := range ObjectTemplates {
		if template.Type == objectType {
			return &template
		}
	}
	return nil
}

func GetObjectsByCategory(category string) []ObjectTemplate {
	var results []ObjectTemplate
	for _, template := range ObjectTemplates {
		if template.Category == category {
			results = append(results, template)
		}
	}
	return results
}

func GetCategories() []string {
	categoryMap := make(map[string]bool)
	var categories []string
	
	for _, template := range ObjectTemplates {
		if !categoryMap[template.Category] {
			categoryMap[template.Category] = true
			categories = append(categories, template.Category)
		}
	}
	
	return categories
}