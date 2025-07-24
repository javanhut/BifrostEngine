package opengl

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

// NewLinesMesh creates a mesh specifically for line rendering
func NewLinesMesh(vertices []float32, indices []uint32) *Mesh {
	var vao, vbo, ebo uint32
	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)
	gl.GenBuffers(1, &ebo)

	gl.BindVertexArray(vao)

	// Upload vertex data
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	// Upload index data
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

	// Position attribute (location = 0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	// Color attribute (location = 1)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)

	gl.BindVertexArray(0)

	return &Mesh{
		VAO:         vao,
		VBO:         vbo,
		EBO:         ebo,
		IndexCount:  int32(len(indices)),
		DrawMode:    gl.LINES,
		vertexCount: int32(len(vertices) / 6), // 6 floats per vertex (pos + color)
	}
}

// DrawAsLines renders the mesh specifically as lines (alternate method)
func (m *Mesh) DrawAsLines() {
	gl.BindVertexArray(m.VAO)
	gl.DrawElements(gl.LINES, m.IndexCount, gl.UNSIGNED_INT, nil)
	gl.BindVertexArray(0)
}