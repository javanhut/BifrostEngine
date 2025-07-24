package opengl

import "math"

func NewSphereMeshWithLighting() *Mesh {
	// Sphere with lighting support: position, color, UV, and normals
	// Format: [x, y, z, r, g, b, u, v, nx, ny, nz]
	
	const (
		latitudeBands  = 16
		longitudeBands = 16
		radius         = 0.5
	)
	
	vertices := make([]float32, 0)
	indices := make([]uint32, 0)
	
	// Generate vertices
	for lat := 0; lat <= latitudeBands; lat++ {
		theta := float64(lat) * math.Pi / float64(latitudeBands)
		sinTheta := math.Sin(theta)
		cosTheta := math.Cos(theta)
		
		for lng := 0; lng <= longitudeBands; lng++ {
			phi := float64(lng) * 2 * math.Pi / float64(longitudeBands)
			sinPhi := math.Sin(phi)
			cosPhi := math.Cos(phi)
			
			x := float32(cosPhi * sinTheta)
			y := float32(cosTheta)
			z := float32(sinPhi * sinTheta)
			
			// Position (scaled by radius)
			px := x * radius
			py := y * radius
			pz := z * radius
			
			// Color based on position (sphere coloring)
			r := (x + 1.0) * 0.5  // Normalize to 0-1
			g := (y + 1.0) * 0.5
			b := (z + 1.0) * 0.5
			
			// UV coordinates
			u := float32(lng) / float32(longitudeBands)
			v := float32(lat) / float32(latitudeBands)
			
			// Normal (same as normalized position for sphere)
			nx := x
			ny := y
			nz := z
			
			vertices = append(vertices, px, py, pz, r, g, b, u, v, nx, ny, nz)
		}
	}
	
	// Generate indices
	for lat := 0; lat < latitudeBands; lat++ {
		for lng := 0; lng < longitudeBands; lng++ {
			first := uint32(lat * (longitudeBands + 1) + lng)
			second := uint32(first + longitudeBands + 1)
			
			// First triangle
			indices = append(indices, first, second, first+1)
			// Second triangle
			indices = append(indices, second, second+1, first+1)
		}
	}
	
	return NewIndexedMeshWithLighting(vertices, indices)
}