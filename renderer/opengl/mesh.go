package opengl

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

type Mesh struct {
	vao        uint32
	vbo        uint32
	ebo        uint32
	vertexCount int32
	indexCount  int32
	indexed     bool
}

func NewMesh(vertices []float32) *Mesh {
	mesh := &Mesh{
		vertexCount: int32(len(vertices) / 6), // 3 for position, 3 for color
		indexed:     false,
	}

	// Generate and bind VAO
	gl.GenVertexArrays(1, &mesh.vao)
	gl.BindVertexArray(mesh.vao)

	// Generate and bind VBO
	gl.GenBuffers(1, &mesh.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, mesh.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	// Position attribute
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	// Color attribute
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)

	// Unbind
	gl.BindVertexArray(0)

	return mesh
}

func NewIndexedMesh(vertices []float32, indices []uint32) *Mesh {
	mesh := &Mesh{
		vertexCount: int32(len(vertices) / 6),
		indexCount:  int32(len(indices)),
		indexed:     true,
	}

	// Generate and bind VAO
	gl.GenVertexArrays(1, &mesh.vao)
	gl.BindVertexArray(mesh.vao)

	// Generate and bind VBO
	gl.GenBuffers(1, &mesh.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, mesh.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	// Generate and bind EBO
	gl.GenBuffers(1, &mesh.ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, mesh.ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

	// Position attribute
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	// Color attribute
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)

	// Unbind
	gl.BindVertexArray(0)

	return mesh
}

func (m *Mesh) Draw() {
	gl.BindVertexArray(m.vao)
	if m.indexed {
		gl.DrawElements(gl.TRIANGLES, m.indexCount, gl.UNSIGNED_INT, gl.PtrOffset(0))
	} else {
		gl.DrawArrays(gl.TRIANGLES, 0, m.vertexCount)
	}
	gl.BindVertexArray(0)
}

func (m *Mesh) Delete() {
	gl.DeleteVertexArrays(1, &m.vao)
	gl.DeleteBuffers(1, &m.vbo)
	if m.indexed {
		gl.DeleteBuffers(1, &m.ebo)
	}
}

func NewMeshLines(vertices []float32) *Mesh {
	mesh := &Mesh{
		vertexCount: int32(len(vertices) / 6), // 3 for position, 3 for color
		indexed:     false,
	}

	// Generate and bind VAO
	gl.GenVertexArrays(1, &mesh.vao)
	gl.BindVertexArray(mesh.vao)

	// Generate and bind VBO
	gl.GenBuffers(1, &mesh.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, mesh.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	// Position attribute
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	// Color attribute
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)

	// Unbind
	gl.BindVertexArray(0)

	return mesh
}

func (m *Mesh) DrawLines() {
	gl.BindVertexArray(m.vao)
	gl.DrawArrays(gl.LINES, 0, m.vertexCount)
	gl.BindVertexArray(0)
}

// NewMeshWithUV creates a mesh with position, color, and UV coordinates
// Vertex format: [x, y, z, r, g, b, u, v] (8 floats per vertex)
func NewMeshWithUV(vertices []float32) *Mesh {
	mesh := &Mesh{
		vertexCount: int32(len(vertices) / 8), // 3 pos + 3 color + 2 UV
		indexed:     false,
	}

	// Generate and bind VAO
	gl.GenVertexArrays(1, &mesh.vao)
	gl.BindVertexArray(mesh.vao)

	// Generate and bind VBO
	gl.GenBuffers(1, &mesh.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, mesh.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	// Position attribute (location 0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 8*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	// Color attribute (location 1)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 8*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)

	// UV coordinate attribute (location 2)
	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, 8*4, gl.PtrOffset(6*4))
	gl.EnableVertexAttribArray(2)

	// Unbind
	gl.BindVertexArray(0)

	return mesh
}

// NewIndexedMeshWithUV creates an indexed mesh with position, color, and UV coordinates
// Vertex format: [x, y, z, r, g, b, u, v] (8 floats per vertex)
func NewIndexedMeshWithUV(vertices []float32, indices []uint32) *Mesh {
	mesh := &Mesh{
		vertexCount: int32(len(vertices) / 8), // 3 pos + 3 color + 2 UV
		indexCount:  int32(len(indices)),
		indexed:     true,
	}

	// Generate and bind VAO
	gl.GenVertexArrays(1, &mesh.vao)
	gl.BindVertexArray(mesh.vao)

	// Generate and bind VBO
	gl.GenBuffers(1, &mesh.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, mesh.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	// Generate and bind EBO
	gl.GenBuffers(1, &mesh.ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, mesh.ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

	// Position attribute (location 0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 8*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	// Color attribute (location 1)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 8*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)

	// UV coordinate attribute (location 2)
	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, 8*4, gl.PtrOffset(6*4))
	gl.EnableVertexAttribArray(2)

	// Unbind
	gl.BindVertexArray(0)

	return mesh
}
// NewMeshWithLighting creates a mesh with position, color, UV coordinates, and normals
// Vertex format: [x, y, z, r, g, b, u, v, nx, ny, nz] (11 floats per vertex)
func NewMeshWithLighting(vertices []float32) *Mesh {
	mesh := &Mesh{
		vertexCount: int32(len(vertices) / 11), // 3 pos + 3 color + 2 UV + 3 normal
		indexed:     false,
	}

	// Generate and bind VAO
	gl.GenVertexArrays(1, &mesh.vao)
	gl.BindVertexArray(mesh.vao)

	// Generate and bind VBO
	gl.GenBuffers(1, &mesh.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, mesh.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	// Position attribute (location 0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 11*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	// Color attribute (location 1)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 11*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)

	// UV coordinate attribute (location 2)
	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, 11*4, gl.PtrOffset(6*4))
	gl.EnableVertexAttribArray(2)

	// Normal attribute (location 3)
	gl.VertexAttribPointer(3, 3, gl.FLOAT, false, 11*4, gl.PtrOffset(8*4))
	gl.EnableVertexAttribArray(3)

	// Unbind
	gl.BindVertexArray(0)

	return mesh
}

// NewIndexedMeshWithLighting creates an indexed mesh with position, color, UV coordinates, and normals
// Vertex format: [x, y, z, r, g, b, u, v, nx, ny, nz] (11 floats per vertex)
func NewIndexedMeshWithLighting(vertices []float32, indices []uint32) *Mesh {
	mesh := &Mesh{
		vertexCount: int32(len(vertices) / 11), // 3 pos + 3 color + 2 UV + 3 normal
		indexCount:  int32(len(indices)),
		indexed:     true,
	}

	// Generate and bind VAO
	gl.GenVertexArrays(1, &mesh.vao)
	gl.BindVertexArray(mesh.vao)

	// Generate and bind VBO
	gl.GenBuffers(1, &mesh.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, mesh.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	// Generate and bind EBO
	gl.GenBuffers(1, &mesh.ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, mesh.ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

	// Position attribute (location 0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 11*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	// Color attribute (location 1)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 11*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)

	// UV coordinate attribute (location 2)
	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, 11*4, gl.PtrOffset(6*4))
	gl.EnableVertexAttribArray(2)

	// Normal attribute (location 3)
	gl.VertexAttribPointer(3, 3, gl.FLOAT, false, 11*4, gl.PtrOffset(8*4))
	gl.EnableVertexAttribArray(3)

	// Unbind
	gl.BindVertexArray(0)

	return mesh
}