package opengl

func NewTriangleMesh() *Mesh {
	// Triangle vertices (single triangle in 3D space)
	vertices := []float32{
		// Position         Color
		-0.5, -0.3, 0.0,  1.0, 0.0, 0.0, // Bottom left - Red
		 0.5, -0.3, 0.0,  0.0, 1.0, 0.0, // Bottom right - Green
		 0.0,  0.6, 0.0,  0.0, 0.0, 1.0, // Top - Blue
	}
	
	return NewMesh(vertices)
}