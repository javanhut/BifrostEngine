package opengl

func NewPyramidMeshWithUV() *Mesh {
	// Pyramid vertices with UV coordinates (square base with triangular sides meeting at apex)
	// Format: [x, y, z, r, g, b, u, v]
	vertices := []float32{
		// Base vertices (square base on XZ plane)
		-0.5, 0.0, -0.5,  0.8, 0.6, 0.4,  0.0, 0.0, // Bottom left - Sandy brown
		 0.5, 0.0, -0.5,  0.9, 0.7, 0.5,  1.0, 0.0, // Bottom right - Light brown
		 0.5, 0.0,  0.5,  0.8, 0.6, 0.4,  1.0, 1.0, // Top right - Sandy brown
		-0.5, 0.0,  0.5,  0.9, 0.7, 0.5,  0.0, 1.0, // Top left - Light brown
		
		// Apex vertex (top of pyramid) - replicated for each face to have proper UV mapping
		// Front face apex
		 0.0, 0.8,  0.0,  1.0, 0.8, 0.6,  0.5, 1.0, // Apex - Golden
		// Right face apex
		 0.0, 0.8,  0.0,  1.0, 0.8, 0.6,  0.5, 1.0, // Apex - Golden
		// Back face apex
		 0.0, 0.8,  0.0,  1.0, 0.8, 0.6,  0.5, 1.0, // Apex - Golden
		// Left face apex
		 0.0, 0.8,  0.0,  1.0, 0.8, 0.6,  0.5, 1.0, // Apex - Golden
		
		// Additional base vertices for side faces (for proper UV mapping)
		// Front face base vertices
		-0.5, 0.0,  0.5,  0.9, 0.7, 0.5,  0.0, 0.0, // Bottom left
		 0.5, 0.0,  0.5,  0.8, 0.6, 0.4,  1.0, 0.0, // Bottom right
		
		// Right face base vertices  
		 0.5, 0.0,  0.5,  0.8, 0.6, 0.4,  0.0, 0.0, // Bottom left
		 0.5, 0.0, -0.5,  0.9, 0.7, 0.5,  1.0, 0.0, // Bottom right
		
		// Back face base vertices
		 0.5, 0.0, -0.5,  0.9, 0.7, 0.5,  0.0, 0.0, // Bottom left
		-0.5, 0.0, -0.5,  0.8, 0.6, 0.4,  1.0, 0.0, // Bottom right
		
		// Left face base vertices
		-0.5, 0.0, -0.5,  0.8, 0.6, 0.4,  0.0, 0.0, // Bottom left
		-0.5, 0.0,  0.5,  0.9, 0.7, 0.5,  1.0, 0.0, // Bottom right
	}
	
	// Indices for pyramid faces
	indices := []uint32{
		// Base (2 triangles)
		0, 1, 2,  // First triangle
		2, 3, 0,  // Second triangle
		
		// Side faces (4 triangular faces using replicated apex vertices)
		8, 4, 9,   // Front face
		10, 5, 11, // Right face  
		12, 6, 13, // Back face
		14, 7, 15, // Left face
	}
	
	return NewIndexedMeshWithUV(vertices, indices)
}