package opengl

import "math"

func NewCylinderMeshWithLighting() *Mesh {
	// Cylinder with lighting support: position, color, UV, and normals
	// Format: [x, y, z, r, g, b, u, v, nx, ny, nz]
	
	const (
		segments = 16
		radius   = 0.5
		height   = 1.0
	)
	
	vertices := make([]float32, 0)
	indices := make([]uint32, 0)
	vertexIndex := uint32(0)
	
	// Generate side vertices (cylindrical surface)
	for i := 0; i <= segments; i++ {
		angle := float64(i) * 2.0 * math.Pi / float64(segments)
		cos := float32(math.Cos(angle))
		sin := float32(math.Sin(angle))
		
		x := cos * radius
		z := sin * radius
		u := float32(i) / float32(segments)
		
		// Bottom vertex
		vertices = append(vertices,
			x, -height/2, z,           // Position
			0.8, 0.4, 0.2,            // Brown color
			u, 0.0,                   // UV
			cos, 0.0, sin)            // Normal (outward from cylinder axis)
		
		// Top vertex  
		vertices = append(vertices,
			x, height/2, z,            // Position
			0.9, 0.6, 0.3,            // Light brown color
			u, 1.0,                   // UV
			cos, 0.0, sin)            // Normal (outward from cylinder axis)
	}
	
	// Generate side indices
	for i := 0; i < segments; i++ {
		bottomLeft := vertexIndex + uint32(i*2)
		bottomRight := vertexIndex + uint32(i*2+2)
		topLeft := vertexIndex + uint32(i*2+1)
		topRight := vertexIndex + uint32(i*2+3)
		
		// Handle wrap-around for last segment
		if i == segments-1 {
			bottomRight = vertexIndex + 0
			topRight = vertexIndex + 1
		}
		
		// Two triangles per quad
		indices = append(indices, bottomLeft, topLeft, bottomRight)
		indices = append(indices, bottomRight, topLeft, topRight)
	}
	
	vertexIndex += uint32((segments + 1) * 2)
	
	// Bottom cap center
	vertices = append(vertices,
		0.0, -height/2, 0.0,       // Position
		0.6, 0.3, 0.1,            // Dark brown
		0.5, 0.5,                 // UV center
		0.0, -1.0, 0.0)           // Normal (downward)
	bottomCenter := vertexIndex
	vertexIndex++
	
	// Bottom cap rim vertices
	for i := 0; i <= segments; i++ {
		angle := float64(i) * 2.0 * math.Pi / float64(segments)
		cos := float32(math.Cos(angle))
		sin := float32(math.Sin(angle))
		
		x := cos * radius
		z := sin * radius
		u := (cos + 1.0) * 0.5
		v := (sin + 1.0) * 0.5
		
		vertices = append(vertices,
			x, -height/2, z,          // Position
			0.7, 0.35, 0.15,         // Medium brown
			u, v,                    // UV (radial)
			0.0, -1.0, 0.0)          // Normal (downward)
	}
	
	// Bottom cap indices
	for i := 0; i < segments; i++ {
		current := vertexIndex + uint32(i)
		next := vertexIndex + uint32(i+1)
		if i == segments-1 {
			next = vertexIndex
		}
		indices = append(indices, bottomCenter, next, current)
	}
	
	vertexIndex += uint32(segments + 1)
	
	// Top cap center
	vertices = append(vertices,
		0.0, height/2, 0.0,        // Position
		1.0, 0.7, 0.4,            // Light brown
		0.5, 0.5,                 // UV center
		0.0, 1.0, 0.0)            // Normal (upward)
	topCenter := vertexIndex
	vertexIndex++
	
	// Top cap rim vertices
	for i := 0; i <= segments; i++ {
		angle := float64(i) * 2.0 * math.Pi / float64(segments)
		cos := float32(math.Cos(angle))
		sin := float32(math.Sin(angle))
		
		x := cos * radius
		z := sin * radius
		u := (cos + 1.0) * 0.5
		v := (sin + 1.0) * 0.5
		
		vertices = append(vertices,
			x, height/2, z,           // Position
			0.9, 0.6, 0.3,           // Light brown
			u, v,                    // UV (radial)
			0.0, 1.0, 0.0)           // Normal (upward)
	}
	
	// Top cap indices
	for i := 0; i < segments; i++ {
		current := vertexIndex + uint32(i)
		next := vertexIndex + uint32(i+1)
		if i == segments-1 {
			next = vertexIndex
		}
		indices = append(indices, topCenter, current, next)
	}
	
	return NewIndexedMeshWithLighting(vertices, indices)
}