package opengl

import (
	"math"
)

func NewSphereMesh() *Mesh {
	const segments = 16
	const rings = 16
	
	vertices := make([]float32, 0)
	indices := make([]uint32, 0)
	
	// Generate sphere vertices
	for ring := 0; ring <= rings; ring++ {
		phi := float64(ring) * math.Pi / float64(rings)
		
		for segment := 0; segment <= segments; segment++ {
			theta := float64(segment) * 2.0 * math.Pi / float64(segments)
			
			x := math.Sin(phi) * math.Cos(theta)
			y := math.Cos(phi)
			z := math.Sin(phi) * math.Sin(theta)
			
			// Position
			vertices = append(vertices, float32(x*0.5), float32(y*0.5), float32(z*0.5))
			
			// Color based on position for visual distinction
			r := float32(0.5 + 0.5*x)
			g := float32(0.5 + 0.5*y)
			b := float32(0.5 + 0.5*z)
			vertices = append(vertices, r, g, b)
		}
	}
	
	// Generate sphere indices
	for ring := 0; ring < rings; ring++ {
		for segment := 0; segment < segments; segment++ {
			current := uint32(ring*(segments+1) + segment)
			next := current + uint32(segments+1)
			
			// First triangle
			indices = append(indices, current, next, current+1)
			// Second triangle
			indices = append(indices, current+1, next, next+1)
		}
	}
	
	return NewIndexedMesh(vertices, indices)
}