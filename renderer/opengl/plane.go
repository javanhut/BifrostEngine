package opengl

func NewPlaneMesh() *Mesh {
	// Plane vertices (flat quad)
	vertices := []float32{
		// Position         Color
		-0.5, 0.0, -0.5,   0.8, 0.8, 0.8, // Bottom left - Light gray
		 0.5, 0.0, -0.5,   0.9, 0.9, 0.9, // Bottom right - Lighter gray
		 0.5, 0.0,  0.5,   1.0, 1.0, 1.0, // Top right - White
		-0.5, 0.0,  0.5,   0.7, 0.7, 0.7, // Top left - Gray
	}
	
	// Indices for the plane (2 triangles)
	indices := []uint32{
		0, 1, 2,  // First triangle
		2, 3, 0,  // Second triangle
	}
	
	return NewIndexedMesh(vertices, indices)
}