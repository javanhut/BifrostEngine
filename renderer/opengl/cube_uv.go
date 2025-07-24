package opengl

func NewCubeMeshWithUV() *Mesh {
	// Cube vertices with positions, colors, and UV coordinates
	// Format: [x, y, z, r, g, b, u, v] (8 floats per vertex)
	vertices := []float32{
		// Front face (red) - UV maps 0,0 to 1,1
		-0.5, -0.5,  0.5,  1.0, 0.0, 0.0,  0.0, 0.0, // Bottom left
		 0.5, -0.5,  0.5,  1.0, 0.0, 0.0,  1.0, 0.0, // Bottom right
		 0.5,  0.5,  0.5,  1.0, 0.0, 0.0,  1.0, 1.0, // Top right
		-0.5,  0.5,  0.5,  1.0, 0.0, 0.0,  0.0, 1.0, // Top left

		// Back face (green) - UV maps 0,0 to 1,1
		 0.5, -0.5, -0.5,  0.0, 1.0, 0.0,  0.0, 0.0, // Bottom left (flipped for back face)
		-0.5, -0.5, -0.5,  0.0, 1.0, 0.0,  1.0, 0.0, // Bottom right
		-0.5,  0.5, -0.5,  0.0, 1.0, 0.0,  1.0, 1.0, // Top right
		 0.5,  0.5, -0.5,  0.0, 1.0, 0.0,  0.0, 1.0, // Top left

		// Left face (blue) - UV maps 0,0 to 1,1
		-0.5, -0.5, -0.5,  0.0, 0.0, 1.0,  0.0, 0.0, // Bottom left
		-0.5, -0.5,  0.5,  0.0, 0.0, 1.0,  1.0, 0.0, // Bottom right
		-0.5,  0.5,  0.5,  0.0, 0.0, 1.0,  1.0, 1.0, // Top right
		-0.5,  0.5, -0.5,  0.0, 0.0, 1.0,  0.0, 1.0, // Top left

		// Right face (yellow) - UV maps 0,0 to 1,1
		 0.5, -0.5,  0.5,  1.0, 1.0, 0.0,  0.0, 0.0, // Bottom left
		 0.5, -0.5, -0.5,  1.0, 1.0, 0.0,  1.0, 0.0, // Bottom right
		 0.5,  0.5, -0.5,  1.0, 1.0, 0.0,  1.0, 1.0, // Top right
		 0.5,  0.5,  0.5,  1.0, 1.0, 0.0,  0.0, 1.0, // Top left

		// Top face (magenta) - UV maps 0,0 to 1,1
		-0.5,  0.5,  0.5,  1.0, 0.0, 1.0,  0.0, 0.0, // Bottom left
		 0.5,  0.5,  0.5,  1.0, 0.0, 1.0,  1.0, 0.0, // Bottom right
		 0.5,  0.5, -0.5,  1.0, 0.0, 1.0,  1.0, 1.0, // Top right
		-0.5,  0.5, -0.5,  1.0, 0.0, 1.0,  0.0, 1.0, // Top left

		// Bottom face (cyan) - UV maps 0,0 to 1,1
		-0.5, -0.5, -0.5,  0.0, 1.0, 1.0,  0.0, 0.0, // Bottom left
		 0.5, -0.5, -0.5,  0.0, 1.0, 1.0,  1.0, 0.0, // Bottom right
		 0.5, -0.5,  0.5,  0.0, 1.0, 1.0,  1.0, 1.0, // Top right
		-0.5, -0.5,  0.5,  0.0, 1.0, 1.0,  0.0, 1.0, // Top left
	}

	// Indices for the cube faces (2 triangles per face)
	indices := []uint32{
		// Front face
		0, 1, 2,   2, 3, 0,
		// Back face
		4, 5, 6,   6, 7, 4,
		// Left face
		8, 9, 10,  10, 11, 8,
		// Right face
		12, 13, 14, 14, 15, 12,
		// Top face
		16, 17, 18, 18, 19, 16,
		// Bottom face
		20, 21, 22, 22, 23, 20,
	}

	return NewIndexedMeshWithUV(vertices, indices)
}