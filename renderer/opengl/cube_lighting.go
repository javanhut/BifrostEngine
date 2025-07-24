package opengl

func NewCubeMeshWithLighting() *Mesh {
	// Cube vertices with position, color, UV coordinates, and normals
	// Format: [x, y, z, r, g, b, u, v, nx, ny, nz] - using indexed approach
	vertices := []float32{
		// Front face vertices (normal: 0, 0, 1)
		-0.5, -0.5,  0.5,  1.0, 0.0, 0.0,  0.0, 0.0,  0.0, 0.0, 1.0, // 0: Bottom left - Red
		 0.5, -0.5,  0.5,  1.0, 0.0, 0.0,  1.0, 0.0,  0.0, 0.0, 1.0, // 1: Bottom right - Red
		 0.5,  0.5,  0.5,  1.0, 0.0, 0.0,  1.0, 1.0,  0.0, 0.0, 1.0, // 2: Top right - Red
		-0.5,  0.5,  0.5,  1.0, 0.0, 0.0,  0.0, 1.0,  0.0, 0.0, 1.0, // 3: Top left - Red

		// Back face vertices (normal: 0, 0, -1)
		-0.5, -0.5, -0.5,  0.0, 1.0, 0.0,  1.0, 0.0,  0.0, 0.0, -1.0, // 4: Bottom left - Green
		-0.5,  0.5, -0.5,  0.0, 1.0, 0.0,  1.0, 1.0,  0.0, 0.0, -1.0, // 5: Top left - Green
		 0.5,  0.5, -0.5,  0.0, 1.0, 0.0,  0.0, 1.0,  0.0, 0.0, -1.0, // 6: Top right - Green
		 0.5, -0.5, -0.5,  0.0, 1.0, 0.0,  0.0, 0.0,  0.0, 0.0, -1.0, // 7: Bottom right - Green

		// Left face vertices (normal: -1, 0, 0)
		-0.5, -0.5, -0.5,  0.0, 0.0, 1.0,  0.0, 0.0,  -1.0, 0.0, 0.0, // 8: Bottom back - Blue
		-0.5, -0.5,  0.5,  0.0, 0.0, 1.0,  1.0, 0.0,  -1.0, 0.0, 0.0, // 9: Bottom front - Blue
		-0.5,  0.5,  0.5,  0.0, 0.0, 1.0,  1.0, 1.0,  -1.0, 0.0, 0.0, // 10: Top front - Blue
		-0.5,  0.5, -0.5,  0.0, 0.0, 1.0,  0.0, 1.0,  -1.0, 0.0, 0.0, // 11: Top back - Blue

		// Right face vertices (normal: 1, 0, 0)
		 0.5, -0.5, -0.5,  1.0, 1.0, 0.0,  1.0, 0.0,  1.0, 0.0, 0.0, // 12: Bottom back - Yellow
		 0.5,  0.5, -0.5,  1.0, 1.0, 0.0,  1.0, 1.0,  1.0, 0.0, 0.0, // 13: Top back - Yellow
		 0.5,  0.5,  0.5,  1.0, 1.0, 0.0,  0.0, 1.0,  1.0, 0.0, 0.0, // 14: Top front - Yellow
		 0.5, -0.5,  0.5,  1.0, 1.0, 0.0,  0.0, 0.0,  1.0, 0.0, 0.0, // 15: Bottom front - Yellow

		// Bottom face vertices (normal: 0, -1, 0)
		-0.5, -0.5, -0.5,  1.0, 0.0, 1.0,  0.0, 1.0,  0.0, -1.0, 0.0, // 16: Back left - Magenta
		 0.5, -0.5, -0.5,  1.0, 0.0, 1.0,  1.0, 1.0,  0.0, -1.0, 0.0, // 17: Back right - Magenta
		 0.5, -0.5,  0.5,  1.0, 0.0, 1.0,  1.0, 0.0,  0.0, -1.0, 0.0, // 18: Front right - Magenta
		-0.5, -0.5,  0.5,  1.0, 0.0, 1.0,  0.0, 0.0,  0.0, -1.0, 0.0, // 19: Front left - Magenta

		// Top face vertices (normal: 0, 1, 0)
		-0.5,  0.5, -0.5,  0.0, 1.0, 1.0,  0.0, 1.0,  0.0, 1.0, 0.0, // 20: Back left - Cyan
		-0.5,  0.5,  0.5,  0.0, 1.0, 1.0,  0.0, 0.0,  0.0, 1.0, 0.0, // 21: Front left - Cyan
		 0.5,  0.5,  0.5,  0.0, 1.0, 1.0,  1.0, 0.0,  0.0, 1.0, 0.0, // 22: Front right - Cyan
		 0.5,  0.5, -0.5,  0.0, 1.0, 1.0,  1.0, 1.0,  0.0, 1.0, 0.0, // 23: Back right - Cyan
	}

	// Indices for the cube faces (2 triangles per face) - proper winding order
	indices := []uint32{
		// Front face
		0, 1, 2,   2, 3, 0,
		// Back face  
		4, 5, 6,   6, 7, 4,
		// Left face
		8, 9, 10,  10, 11, 8,
		// Right face
		12, 13, 14, 14, 15, 12,
		// Bottom face
		16, 17, 18, 18, 19, 16,
		// Top face
		20, 21, 22, 22, 23, 20,
	}

	return NewIndexedMeshWithLighting(vertices, indices)
}