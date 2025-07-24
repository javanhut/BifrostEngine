package opengl

func NewTriangleMeshWithLighting() *Mesh {
	// Triangle vertices with position, color, UV coordinates, and normals
	// Format: [x, y, z, r, g, b, u, v, nx, ny, nz]
	// Calculate normal using cross product of two edges
	// Edge1: from bottom-left to bottom-right
	// Edge2: from bottom-left to top
	// Normal = Edge1 × Edge2 (pointing toward viewer)
	// For triangle in XY plane: normal = (0, 0, 1)
	
	vertices := []float32{
		// Position         Color           UV      Normal
		-0.5, -0.3, 0.0,  1.0, 0.0, 0.0,  0.0, 0.0,  0.0, 0.0, 1.0, // Bottom left - Red
		 0.5, -0.3, 0.0,  0.0, 1.0, 0.0,  1.0, 0.0,  0.0, 0.0, 1.0, // Bottom right - Green
		 0.0,  0.6, 0.0,  0.0, 0.0, 1.0,  0.5, 1.0,  0.0, 0.0, 1.0, // Top - Blue
	}
	
	return NewMeshWithLighting(vertices)
}