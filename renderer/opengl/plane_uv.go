package opengl

func NewPlaneMeshWithUV() *Mesh {
	// Plane vertices with UV coordinates (flat quad)
	// Format: [x, y, z, r, g, b, u, v]
	vertices := []float32{
		// Position         Color             UV
		-0.5, 0.0, -0.5,   0.8, 0.8, 0.8,   0.0, 0.0, // Bottom left
		 0.5, 0.0, -0.5,   0.9, 0.9, 0.9,   1.0, 0.0, // Bottom right
		 0.5, 0.0,  0.5,   1.0, 1.0, 1.0,   1.0, 1.0, // Top right
		-0.5, 0.0,  0.5,   0.7, 0.7, 0.7,   0.0, 1.0, // Top left
	}
	
	// Indices for the plane (2 triangles)
	indices := []uint32{
		0, 1, 2,  // First triangle
		2, 3, 0,  // Second triangle
	}
	
	return NewIndexedMeshWithUV(vertices, indices)
}