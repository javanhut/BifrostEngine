package opengl

func NewTriangleMeshWithUV() *Mesh {
	// Triangle vertices with UV coordinates (single triangle in 3D space)
	// Format: [x, y, z, r, g, b, u, v]
	vertices := []float32{
		// Position         Color           UV
		-0.5, -0.3, 0.0,  1.0, 0.0, 0.0,  0.0, 0.0, // Bottom left - Red
		 0.5, -0.3, 0.0,  0.0, 1.0, 0.0,  1.0, 0.0, // Bottom right - Green
		 0.0,  0.6, 0.0,  0.0, 0.0, 1.0,  0.5, 1.0, // Top - Blue
	}
	
	return NewMeshWithUV(vertices)
}