package assets

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	bmath "github.com/javanhut/BifrostEngine/m/v2/math"
	"github.com/javanhut/BifrostEngine/m/v2/renderer/opengl"
)

// OBJMesh represents a loaded OBJ file mesh
type OBJMesh struct {
	Name     string
	Vertices []float32 // [x,y,z,r,g,b,u,v,nx,ny,nz] format
	Indices  []uint32
	Mesh     *opengl.Mesh
}

// OBJLoader handles loading .obj files
type OBJLoader struct {
	// Configuration
	DefaultColor    [3]float32 // Default color for vertices without materials
	FlipUV          bool       // Whether to flip V coordinate (some formats need this)
	GenerateNormals bool       // Generate normals if not present in file
}

// NewOBJLoader creates a new OBJ loader with default settings
func NewOBJLoader() *OBJLoader {
	return &OBJLoader{
		DefaultColor:    [3]float32{0.8, 0.8, 0.8}, // Light gray
		FlipUV:          true,                       // Common for OpenGL
		GenerateNormals: true,                       // Generate if missing
	}
}

// LoadOBJ loads an OBJ file and returns an OBJMesh
func (loader *OBJLoader) LoadOBJ(filepath string) (*OBJMesh, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open OBJ file: %v", err)
	}
	defer file.Close()

	// Temporary storage for parsing
	var positions []bmath.Vector3
	var texCoords []bmath.Vector2
	var normals []bmath.Vector3
	var faces []Face

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue // Skip empty lines and comments
		}

		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}

		switch parts[0] {
		case "v": // Vertex position
			if len(parts) >= 4 {
				pos, err := parseVector3(parts[1:4])
				if err != nil {
					return nil, fmt.Errorf("invalid vertex position: %v", err)
				}
				positions = append(positions, pos)
			}

		case "vt": // Texture coordinate
			if len(parts) >= 3 {
				tc, err := parseVector2(parts[1:3])
				if err != nil {
					return nil, fmt.Errorf("invalid texture coordinate: %v", err)
				}
				if loader.FlipUV {
					tc.Y = 1.0 - tc.Y // Flip V coordinate
				}
				texCoords = append(texCoords, tc)
			}

		case "vn": // Normal
			if len(parts) >= 4 {
				normal, err := parseVector3(parts[1:4])
				if err != nil {
					return nil, fmt.Errorf("invalid normal: %v", err)
				}
				normals = append(normals, normal.Normalize())
			}

		case "f": // Face
			if len(parts) >= 4 {
				face, err := parseFace(parts[1:])
				if err != nil {
					return nil, fmt.Errorf("invalid face: %v", err)
				}
				faces = append(faces, face)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	// Generate normals if not present and requested
	if len(normals) == 0 && loader.GenerateNormals {
		normals = loader.generateNormals(positions, faces)
	}

	// Build final vertex data and indices
	vertices, indices := loader.buildMeshData(positions, texCoords, normals, faces)

	// Create OpenGL mesh
	mesh := opengl.NewIndexedMeshWithLighting(vertices, indices)

	return &OBJMesh{
		Name:     extractFileName(filepath),
		Vertices: vertices,
		Indices:  indices,
		Mesh:     mesh,
	}, nil
}

// Face represents a face with vertex/texture/normal indices
type Face struct {
	Vertices []FaceVertex
}

// FaceVertex represents a single vertex reference in a face
type FaceVertex struct {
	PositionIndex int // 1-based index into positions array
	TexCoordIndex int // 1-based index into texCoords array (0 if not present)
	NormalIndex   int // 1-based index into normals array (0 if not present)
}

// parseVector3 parses a 3D vector from string parts
func parseVector3(parts []string) (bmath.Vector3, error) {
	if len(parts) < 3 {
		return bmath.Vector3{}, fmt.Errorf("need 3 components for vector3")
	}

	x, err := strconv.ParseFloat(parts[0], 32)
	if err != nil {
		return bmath.Vector3{}, err
	}

	y, err := strconv.ParseFloat(parts[1], 32)
	if err != nil {
		return bmath.Vector3{}, err
	}

	z, err := strconv.ParseFloat(parts[2], 32)
	if err != nil {
		return bmath.Vector3{}, err
	}

	return bmath.NewVector3(float32(x), float32(y), float32(z)), nil
}

// parseVector2 parses a 2D vector from string parts
func parseVector2(parts []string) (bmath.Vector2, error) {
	if len(parts) < 2 {
		return bmath.Vector2{}, fmt.Errorf("need 2 components for vector2")
	}

	x, err := strconv.ParseFloat(parts[0], 32)
	if err != nil {
		return bmath.Vector2{}, err
	}

	y, err := strconv.ParseFloat(parts[1], 32)
	if err != nil {
		return bmath.Vector2{}, err
	}

	return bmath.NewVector2(float32(x), float32(y)), nil
}

// parseFace parses a face definition (supports v, v/vt, v/vt/vn, v//vn formats)
func parseFace(parts []string) (Face, error) {
	var face Face

	for _, part := range parts {
		indices := strings.Split(part, "/")
		if len(indices) == 0 {
			return face, fmt.Errorf("empty face vertex")
		}

		var fv FaceVertex

		// Position index (required)
		if pos, err := strconv.Atoi(indices[0]); err == nil {
			fv.PositionIndex = pos
		} else {
			return face, fmt.Errorf("invalid position index: %s", indices[0])
		}

		// Texture coordinate index (optional)
		if len(indices) > 1 && indices[1] != "" {
			if tc, err := strconv.Atoi(indices[1]); err == nil {
				fv.TexCoordIndex = tc
			}
		}

		// Normal index (optional)
		if len(indices) > 2 && indices[2] != "" {
			if n, err := strconv.Atoi(indices[2]); err == nil {
				fv.NormalIndex = n
			}
		}

		face.Vertices = append(face.Vertices, fv)
	}

	return face, nil
}

// generateNormals generates face normals for faces without normals
func (loader *OBJLoader) generateNormals(positions []bmath.Vector3, faces []Face) []bmath.Vector3 {
	normals := make([]bmath.Vector3, len(positions))

	// Calculate face normals and accumulate at vertices
	for _, face := range faces {
		if len(face.Vertices) < 3 {
			continue
		}

		// Get first three vertices of face
		v1 := positions[face.Vertices[0].PositionIndex-1]
		v2 := positions[face.Vertices[1].PositionIndex-1]
		v3 := positions[face.Vertices[2].PositionIndex-1]

		// Calculate face normal using cross product
		edge1 := v2.Sub(v1)
		edge2 := v3.Sub(v1)
		faceNormal := edge1.Cross(edge2).Normalize()

		// Accumulate normal at each vertex of the face
		for _, fv := range face.Vertices {
			idx := fv.PositionIndex - 1
			if idx >= 0 && idx < len(normals) {
				normals[idx] = normals[idx].Add(faceNormal)
			}
		}
	}

	// Normalize accumulated normals
	for i := range normals {
		normals[i] = normals[i].Normalize()
	}

	return normals
}

// buildMeshData converts parsed OBJ data into vertex array and index array
func (loader *OBJLoader) buildMeshData(positions []bmath.Vector3, texCoords []bmath.Vector2, normals []bmath.Vector3, faces []Face) ([]float32, []uint32) {
	var vertices []float32
	var indices []uint32

	vertexIndex := uint32(0)

	for _, face := range faces {
		// Triangulate face (simple fan triangulation)
		for i := 1; i < len(face.Vertices)-1; i++ {
			// Create triangle from vertices 0, i, i+1
			for _, j := range []int{0, i, i + 1} {
				fv := face.Vertices[j]

				// Position
				var pos bmath.Vector3
				if fv.PositionIndex > 0 && fv.PositionIndex <= len(positions) {
					pos = positions[fv.PositionIndex-1]
				}

				// Color (use default)
				color := loader.DefaultColor

				// Texture coordinates
				var uv bmath.Vector2
				if fv.TexCoordIndex > 0 && fv.TexCoordIndex <= len(texCoords) {
					uv = texCoords[fv.TexCoordIndex-1]
				}

				// Normal
				var normal bmath.Vector3
				if fv.NormalIndex > 0 && fv.NormalIndex <= len(normals) {
					normal = normals[fv.NormalIndex-1]
				} else if len(normals) > 0 && fv.PositionIndex > 0 && fv.PositionIndex <= len(normals) {
					// Use position index for normal if no explicit normal index
					normal = normals[fv.PositionIndex-1]
				}

				// Append vertex data [x,y,z,r,g,b,u,v,nx,ny,nz]
				vertices = append(vertices,
					pos.X, pos.Y, pos.Z,           // Position
					color[0], color[1], color[2],  // Color
					uv.X, uv.Y,                    // UV
					normal.X, normal.Y, normal.Z,  // Normal
				)

				indices = append(indices, vertexIndex)
				vertexIndex++
			}
		}
	}

	return vertices, indices
}

// extractFileName extracts the filename without path and extension
func extractFileName(filepath string) string {
	// Get basename
	parts := strings.Split(filepath, "/")
	filename := parts[len(parts)-1]

	// Remove extension
	if dotIndex := strings.LastIndex(filename, "."); dotIndex != -1 {
		filename = filename[:dotIndex]
	}

	return filename
}

// Cleanup releases OpenGL resources
func (mesh *OBJMesh) Cleanup() {
	if mesh.Mesh != nil {
		mesh.Mesh.Delete()
	}
}