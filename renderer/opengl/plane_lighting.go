package opengl

func NewPlaneMeshWithLighting() *Mesh {
	// Plane vertices with position, color, UV coordinates, and normals
	// Format: [x, y, z, r, g, b, u, v, nx, ny, nz]
	// Plane in XZ plane facing upward (normal: 0, 1, 0)
	vertices := []float32{
		// Triangle 1
		-0.5, 0.0, -0.5,  0.5, 0.5, 0.5,  0.0, 0.0,  0.0, 1.0, 0.0, // Bottom left - Gray
		 0.5, 0.0, -0.5,  0.7, 0.7, 0.7,  1.0, 0.0,  0.0, 1.0, 0.0, // Bottom right - Light gray
		 0.5, 0.0,  0.5,  0.9, 0.9, 0.9,  1.0, 1.0,  0.0, 1.0, 0.0, // Top right - Almost white
		
		// Triangle 2  
		 0.5, 0.0,  0.5,  0.9, 0.9, 0.9,  1.0, 1.0,  0.0, 1.0, 0.0, // Top right - Almost white
		-0.5, 0.0,  0.5,  0.3, 0.3, 0.3,  0.0, 1.0,  0.0, 1.0, 0.0, // Top left - Dark gray
		-0.5, 0.0, -0.5,  0.5, 0.5, 0.5,  0.0, 0.0,  0.0, 1.0, 0.0, // Bottom left - Gray
	}
	
	return NewMeshWithLighting(vertices)
}