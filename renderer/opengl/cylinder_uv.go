package opengl

import (
	"math"
)

func NewCylinderMeshWithUV() *Mesh {
	const segments = 16
	const height = 1.0
	const radius = 0.5
	
	vertices := make([]float32, 0)
	indices := make([]uint32, 0)
	
	// Bottom center vertex (for bottom cap)
	vertices = append(vertices, 0.0, -height/2, 0.0, 0.0, 0.0, 1.0, 0.5, 0.5) // Center UV for radial pattern
	bottomCenterIndex := uint32(0)
	
	// Top center vertex (for top cap)
	vertices = append(vertices, 0.0, height/2, 0.0, 1.0, 0.0, 0.0, 0.5, 0.5) // Center UV for radial pattern
	topCenterIndex := uint32(1)
	
	vertexIndex := uint32(2)
	
	// Generate bottom circle, top circle, and side vertices
	for i := 0; i <= segments; i++ {
		angle := float64(i) * 2.0 * math.Pi / float64(segments)
		x := math.Cos(angle) * radius
		z := math.Sin(angle) * radius
		
		// UV coordinate for circular caps (radial mapping)
		capU := float32(0.5 + 0.5*math.Cos(angle))
		capV := float32(0.5 + 0.5*math.Sin(angle))
		
		// UV coordinate for side (cylindrical mapping)
		sideU := float32(i) / float32(segments)
		
		// Bottom circle vertex (for bottom cap)
		vertices = append(vertices, float32(x), -height/2, float32(z), 0.0, 0.5, 1.0, capU, capV)
		bottomVertexIndex := vertexIndex
		vertexIndex++
		
		// Top circle vertex (for top cap)
		vertices = append(vertices, float32(x), height/2, float32(z), 1.0, 0.5, 0.0, capU, capV)
		topVertexIndex := vertexIndex
		vertexIndex++
		
		// Side vertices (separate vertices for proper UV mapping)
		// Bottom side vertex
		vertices = append(vertices, float32(x), -height/2, float32(z), 0.5, 0.8, 0.2, sideU, 0.0)
		bottomSideIndex := vertexIndex
		vertexIndex++
		
		// Top side vertex  
		vertices = append(vertices, float32(x), height/2, float32(z), 0.8, 0.5, 0.2, sideU, 1.0)
		topSideIndex := vertexIndex
		vertexIndex++
		
		if i < segments {
			// Calculate next indices
			nextBottomIndex := bottomVertexIndex + 4 // Skip the side vertices
			nextTopIndex := topVertexIndex + 4
			nextBottomSideIndex := bottomSideIndex + 4
			nextTopSideIndex := topSideIndex + 4
			
			// Bottom face triangle
			indices = append(indices, bottomCenterIndex, nextBottomIndex, bottomVertexIndex)
			
			// Top face triangle
			indices = append(indices, topCenterIndex, topVertexIndex, nextTopIndex)
			
			// Side face (2 triangles using side vertices)
			indices = append(indices, bottomSideIndex, topSideIndex, nextBottomSideIndex)
			indices = append(indices, nextBottomSideIndex, topSideIndex, nextTopSideIndex)
		}
	}
	
	return NewIndexedMeshWithUV(vertices, indices)
}