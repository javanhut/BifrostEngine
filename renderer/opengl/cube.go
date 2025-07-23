package opengl

func NewCubeMesh() *Mesh {
	// Cube vertices with positions and colors
	// Each face has different colors to see rotation clearly
	vertices := []float32{
		// Front face (red)
		-0.5, -0.5,  0.5,  1.0, 0.0, 0.0,
		 0.5, -0.5,  0.5,  1.0, 0.0, 0.0,
		 0.5,  0.5,  0.5,  1.0, 0.0, 0.0,
		-0.5,  0.5,  0.5,  1.0, 0.0, 0.0,

		// Back face (green)
		-0.5, -0.5, -0.5,  0.0, 1.0, 0.0,
		 0.5, -0.5, -0.5,  0.0, 1.0, 0.0,
		 0.5,  0.5, -0.5,  0.0, 1.0, 0.0,
		-0.5,  0.5, -0.5,  0.0, 1.0, 0.0,

		// Left face (blue)
		-0.5, -0.5, -0.5,  0.0, 0.0, 1.0,
		-0.5, -0.5,  0.5,  0.0, 0.0, 1.0,
		-0.5,  0.5,  0.5,  0.0, 0.0, 1.0,
		-0.5,  0.5, -0.5,  0.0, 0.0, 1.0,

		// Right face (yellow)
		 0.5, -0.5, -0.5,  1.0, 1.0, 0.0,
		 0.5, -0.5,  0.5,  1.0, 1.0, 0.0,
		 0.5,  0.5,  0.5,  1.0, 1.0, 0.0,
		 0.5,  0.5, -0.5,  1.0, 1.0, 0.0,

		// Top face (magenta)
		-0.5,  0.5, -0.5,  1.0, 0.0, 1.0,
		 0.5,  0.5, -0.5,  1.0, 0.0, 1.0,
		 0.5,  0.5,  0.5,  1.0, 0.0, 1.0,
		-0.5,  0.5,  0.5,  1.0, 0.0, 1.0,

		// Bottom face (cyan)
		-0.5, -0.5, -0.5,  0.0, 1.0, 1.0,
		 0.5, -0.5, -0.5,  0.0, 1.0, 1.0,
		 0.5, -0.5,  0.5,  0.0, 1.0, 1.0,
		-0.5, -0.5,  0.5,  0.0, 1.0, 1.0,
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

	return NewIndexedMesh(vertices, indices)
}