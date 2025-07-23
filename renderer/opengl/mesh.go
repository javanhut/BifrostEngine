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