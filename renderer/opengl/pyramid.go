package opengl

func NewPyramidMesh() *Mesh {
	// Pyramid vertices (square base with triangular sides meeting at apex)
	vertices := []float32{
		// Base vertices (square base on XZ plane)
		-0.5, 0.0, -0.5,  0.8, 0.6, 0.4, // Bottom left - Sandy brown
		 0.5, 0.0, -0.5,  0.9, 0.7, 0.5, // Bottom right - Light brown
		 0.5, 0.0,  0.5,  0.8, 0.6, 0.4, // Top right - Sandy brown
		-0.5, 0.0,  0.5,  0.9, 0.7, 0.5, // Top left - Light brown
		
		// Apex vertex (top of pyramid)
		 0.0, 0.8,  0.0,  1.0, 0.8, 0.6, // Apex - Golden
	}
	
	// Indices for pyramid faces
	indices := []uint32{
		// Base (2 triangles)
		0, 1, 2,  // First triangle
		2, 3, 0,  // Second triangle
		
		// Side faces (4 triangular faces)
		0, 4, 1,  // Front face
		1, 4, 2,  // Right face  
		2, 4, 3,  // Back face
		3, 4, 0,  // Left face
	}
	
	return NewIndexedMesh(vertices, indices)
}