package opengl

import (
	"math"
)

func NewCylinderMesh() *Mesh {
	const segments = 16
	const height = 1.0
	const radius = 0.5
	
	vertices := make([]float32, 0)
	indices := make([]uint32, 0)
	
	// Bottom center vertex
	vertices = append(vertices, 0.0, -height/2, 0.0, 0.0, 0.0, 1.0) // Blue bottom
	bottomCenterIndex := uint32(0)
	
	// Top center vertex
	vertices = append(vertices, 0.0, height/2, 0.0, 1.0, 0.0, 0.0) // Red top
	topCenterIndex := uint32(1)
	
	vertexIndex := uint32(2)
	
	// Generate bottom and top circle vertices and side vertices
	for i := 0; i <= segments; i++ {
		angle := float64(i) * 2.0 * math.Pi / float64(segments)
		x := math.Cos(angle) * radius
		z := math.Sin(angle) * radius
		
		// Bottom circle vertex
		vertices = append(vertices, float32(x), -height/2, float32(z), 0.0, 0.5, 1.0) // Light blue
		bottomVertexIndex := vertexIndex
		vertexIndex++
		
		// Top circle vertex
		vertices = append(vertices, float32(x), height/2, float32(z), 1.0, 0.5, 0.0) // Orange
		topVertexIndex := vertexIndex
		vertexIndex++
		
		if i < segments {
			// Bottom face triangle
			nextBottomIndex := bottomVertexIndex + 2
			indices = append(indices, bottomCenterIndex, nextBottomIndex, bottomVertexIndex)
			
			// Top face triangle
			nextTopIndex := topVertexIndex + 2
			indices = append(indices, topCenterIndex, topVertexIndex, nextTopIndex)
			
			// Side face (2 triangles)
			indices = append(indices, bottomVertexIndex, topVertexIndex, nextBottomIndex)
			indices = append(indices, nextBottomIndex, topVertexIndex, nextTopIndex)
		}
	}
	
	return NewIndexedMesh(vertices, indices)
}