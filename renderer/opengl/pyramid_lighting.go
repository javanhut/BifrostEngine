package opengl

import "math"

func NewPyramidMeshWithLighting() *Mesh {
	// Pyramid vertices with position, color, UV coordinates, and normals
	// Format: [x, y, z, r, g, b, u, v, nx, ny, nz]
	// Pyramid with square base and triangular sides meeting at apex
	
	vertices := make([]float32, 0)
	indices := make([]uint32, 0)
	
	// Base vertices (square base on XZ plane, normal pointing down)
	baseVertices := []float32{
		// Position         Color           UV      Normal (downward)
		-0.5, 0.0, -0.5,  0.8, 0.6, 0.4,  0.0, 0.0,  0.0, -1.0, 0.0, // Bottom left
		 0.5, 0.0, -0.5,  0.9, 0.7, 0.5,  1.0, 0.0,  0.0, -1.0, 0.0, // Bottom right
		 0.5, 0.0,  0.5,  0.8, 0.6, 0.4,  1.0, 1.0,  0.0, -1.0, 0.0, // Top right
		-0.5, 0.0,  0.5,  0.9, 0.7, 0.5,  0.0, 1.0,  0.0, -1.0, 0.0, // Top left
	}
	vertices = append(vertices, baseVertices...)
	
	// Base indices (2 triangles)
	baseIndices := []uint32{
		0, 1, 2,  // First triangle
		2, 3, 0,  // Second triangle
	}
	indices = append(indices, baseIndices...)
	
	// Apex position
	apex := [3]float32{0.0, 0.8, 0.0}
	
	// Calculate normals for each triangular face
	// Face normals are calculated using cross product of edge vectors
	
	// Front face: base edge from (-0.5,0,0.5) to (0.5,0,0.5), apex at (0,0.8,0)
	// Edge1: (1, 0, 0), Edge2: (0.5, 0.8, -0.5)
	// Normal = Edge1 × Edge2 = (0, 0.5, 0.8) normalized
	frontNormalLen := float32(math.Sqrt(0.5*0.5 + 0.8*0.8))
	frontNormal := [3]float32{0.0, 0.5 / frontNormalLen, 0.8 / frontNormalLen}
	
	// Right face: base edge from (0.5,0,0.5) to (0.5,0,-0.5), apex at (0,0.8,0)  
	// Edge1: (0, 0, -1), Edge2: (-0.5, 0.8, 0.5)
	// Normal = Edge1 × Edge2 = (0.8, -0.5, 0) normalized
	rightNormalLen := float32(math.Sqrt(0.8*0.8 + 0.5*0.5))
	rightNormal := [3]float32{0.8 / rightNormalLen, -0.5 / rightNormalLen, 0.0}
	
	// Back face: base edge from (0.5,0,-0.5) to (-0.5,0,-0.5), apex at (0,0.8,0)
	// Edge1: (-1, 0, 0), Edge2: (-0.5, 0.8, 0.5) 
	// Normal = Edge1 × Edge2 = (0, -0.5, -0.8) normalized
	backNormal := [3]float32{0.0, -0.5 / frontNormalLen, -0.8 / frontNormalLen}
	
	// Left face: base edge from (-0.5,0,-0.5) to (-0.5,0,0.5), apex at (0,0.8,0)
	// Edge1: (0, 0, 1), Edge2: (0.5, 0.8, -0.5)
	// Normal = Edge1 × Edge2 = (-0.8, 0.5, 0) normalized  
	leftNormal := [3]float32{-0.8 / rightNormalLen, 0.5 / rightNormalLen, 0.0}
	
	// Front face vertices (base edge: left to right, then apex)
	frontFaceVertices := []float32{
		// Base left vertex for front face
		-0.5, 0.0,  0.5,  0.9, 0.7, 0.5,  0.0, 0.0,  frontNormal[0], frontNormal[1], frontNormal[2],
		// Base right vertex for front face  
		 0.5, 0.0,  0.5,  0.8, 0.6, 0.4,  1.0, 0.0,  frontNormal[0], frontNormal[1], frontNormal[2],
		// Apex for front face
		 apex[0], apex[1], apex[2],  1.0, 0.8, 0.6,  0.5, 1.0,  frontNormal[0], frontNormal[1], frontNormal[2],
	}
	vertices = append(vertices, frontFaceVertices...)
	indices = append(indices, 4, 5, 6) // Front face triangle
	
	// Right face vertices
	rightFaceVertices := []float32{
		// Base front vertex for right face
		 0.5, 0.0,  0.5,  0.8, 0.6, 0.4,  0.0, 0.0,  rightNormal[0], rightNormal[1], rightNormal[2],
		// Base back vertex for right face
		 0.5, 0.0, -0.5,  0.9, 0.7, 0.5,  1.0, 0.0,  rightNormal[0], rightNormal[1], rightNormal[2],
		// Apex for right face
		 apex[0], apex[1], apex[2],  1.0, 0.8, 0.6,  0.5, 1.0,  rightNormal[0], rightNormal[1], rightNormal[2],
	}
	vertices = append(vertices, rightFaceVertices...)
	indices = append(indices, 7, 8, 9) // Right face triangle
	
	// Back face vertices
	backFaceVertices := []float32{
		// Base right vertex for back face
		 0.5, 0.0, -0.5,  0.9, 0.7, 0.5,  0.0, 0.0,  backNormal[0], backNormal[1], backNormal[2],
		// Base left vertex for back face
		-0.5, 0.0, -0.5,  0.8, 0.6, 0.4,  1.0, 0.0,  backNormal[0], backNormal[1], backNormal[2],
		// Apex for back face
		 apex[0], apex[1], apex[2],  1.0, 0.8, 0.6,  0.5, 1.0,  backNormal[0], backNormal[1], backNormal[2],
	}
	vertices = append(vertices, backFaceVertices...)
	indices = append(indices, 10, 11, 12) // Back face triangle
	
	// Left face vertices
	leftFaceVertices := []float32{
		// Base back vertex for left face
		-0.5, 0.0, -0.5,  0.8, 0.6, 0.4,  0.0, 0.0,  leftNormal[0], leftNormal[1], leftNormal[2],
		// Base front vertex for left face
		-0.5, 0.0,  0.5,  0.9, 0.7, 0.5,  1.0, 0.0,  leftNormal[0], leftNormal[1], leftNormal[2],
		// Apex for left face
		 apex[0], apex[1], apex[2],  1.0, 0.8, 0.6,  0.5, 1.0,  leftNormal[0], leftNormal[1], leftNormal[2],
	}
	vertices = append(vertices, leftFaceVertices...)
	indices = append(indices, 13, 14, 15) // Left face triangle
	
	return NewIndexedMeshWithLighting(vertices, indices)
}