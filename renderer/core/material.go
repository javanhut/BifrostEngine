package core

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"strings"
	
	"github.com/go-gl/gl/v4.1-core/gl"
)

// Material represents rendering material properties
type Material struct {
	Name       string
	DiffuseColor [3]float32 // RGB color
	SpecularColor [3]float32 // Specular highlight color
	Shininess  float32      // Specular shininess factor
	DiffuseTexture uint32    // OpenGL texture ID (0 = no texture)
	UseTexture bool         // Whether to use texture or color
}

// TextureManager handles texture loading and caching
type TextureManager struct {
	textures map[string]uint32 // filename -> texture ID
}

// NewTextureManager creates a new texture manager
func NewTextureManager() *TextureManager {
	return &TextureManager{
		textures: make(map[string]uint32),
	}
}

// LoadTexture loads a texture from file and returns OpenGL texture ID
func (tm *TextureManager) LoadTexture(filepath string) (uint32, error) {
	// Check if already loaded
	if texID, exists := tm.textures[filepath]; exists {
		return texID, nil
	}
	
	// Open file
	file, err := os.Open(filepath)
	if err != nil {
		return 0, fmt.Errorf("failed to open texture file %s: %w", filepath, err)
	}
	defer file.Close()
	
	// Decode image based on file extension
	var img image.Image
	ext := strings.ToLower(filepath[strings.LastIndex(filepath, "."):])
	
	switch ext {
	case ".png":
		img, err = png.Decode(file)
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(file)
	default:
		return 0, fmt.Errorf("unsupported texture format: %s", ext)
	}
	
	if err != nil {
		return 0, fmt.Errorf("failed to decode texture %s: %w", filepath, err)
	}
	
	// Convert to RGBA
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	rgba := image.NewRGBA(bounds)
	
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			rgba.Set(x, y, img.At(x, y))
		}
	}
	
	// Generate OpenGL texture
	var textureID uint32
	gl.GenTextures(1, &textureID)
	gl.BindTexture(gl.TEXTURE_2D, textureID)
	
	// Upload texture data
	gl.TexImage2D(
		gl.TEXTURE_2D, 0, gl.RGBA,
		int32(width), int32(height), 0,
		gl.RGBA, gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix),
	)
	
	// Set texture parameters
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	
	// Generate mipmaps
	gl.GenerateMipmap(gl.TEXTURE_2D)
	
	// Unbind texture
	gl.BindTexture(gl.TEXTURE_2D, 0)
	
	// Cache the texture
	tm.textures[filepath] = textureID
	
	fmt.Printf("Loaded texture: %s (ID: %d, Size: %dx%d)\n", filepath, textureID, width, height)
	return textureID, nil
}

// DeleteTexture removes a texture from OpenGL and cache
func (tm *TextureManager) DeleteTexture(filepath string) {
	if texID, exists := tm.textures[filepath]; exists {
		gl.DeleteTextures(1, &texID)
		delete(tm.textures, filepath)
	}
}

// Cleanup removes all textures
func (tm *TextureManager) Cleanup() {
	for filepath, texID := range tm.textures {
		gl.DeleteTextures(1, &texID)
		delete(tm.textures, filepath)
	}
}

// NewMaterial creates a new material with default values
func NewMaterial(name string) *Material {
	return &Material{
		Name:          name,
		DiffuseColor:  [3]float32{1.0, 1.0, 1.0}, // White
		SpecularColor: [3]float32{1.0, 1.0, 1.0}, // White
		Shininess:     32.0,
		DiffuseTexture: 0,
		UseTexture:    false,
	}
}

// SetDiffuseColor sets the material's diffuse color
func (m *Material) SetDiffuseColor(r, g, b float32) {
	m.DiffuseColor = [3]float32{r, g, b}
}

// SetSpecularColor sets the material's specular color
func (m *Material) SetSpecularColor(r, g, b float32) {
	m.SpecularColor = [3]float32{r, g, b}
}

// SetTexture assigns a texture to this material
func (m *Material) SetTexture(textureID uint32) {
	m.DiffuseTexture = textureID
	m.UseTexture = textureID != 0
}

// Bind activates this material for rendering
func (m *Material) Bind() {
	if m.UseTexture && m.DiffuseTexture != 0 {
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, m.DiffuseTexture)
	} else {
		gl.BindTexture(gl.TEXTURE_2D, 0)
	}
}

// LoadCheckerboardTexture creates a procedural checkerboard texture
func (tm *TextureManager) LoadCheckerboardTexture() (uint32, error) {
	// Check if already loaded
	if texID, exists := tm.textures["checkerboard"]; exists {
		return texID, nil
	}
	
	// Create checkerboard pattern (8x8 squares, 64x64 pixels total)
	const size = 64
	const squareSize = 8
	pixels := make([]byte, size*size*4) // RGBA
	
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			// Determine if this pixel is black or white
			squareX := x / squareSize
			squareY := y / squareSize
			isWhite := (squareX+squareY)%2 == 0
			
			pixelIndex := (y*size + x) * 4
			if isWhite {
				pixels[pixelIndex] = 255   // R
				pixels[pixelIndex+1] = 255 // G  
				pixels[pixelIndex+2] = 255 // B
			} else {
				pixels[pixelIndex] = 0     // R
				pixels[pixelIndex+1] = 0   // G
				pixels[pixelIndex+2] = 0   // B
			}
			pixels[pixelIndex+3] = 255 // A (full opacity)
		}
	}
	
	// Generate OpenGL texture
	var textureID uint32
	gl.GenTextures(1, &textureID)
	gl.BindTexture(gl.TEXTURE_2D, textureID)
	
	// Upload texture data
	gl.TexImage2D(
		gl.TEXTURE_2D, 0, gl.RGBA,
		size, size, 0,
		gl.RGBA, gl.UNSIGNED_BYTE,
		gl.Ptr(pixels),
	)
	
	// Set texture parameters
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	
	// Unbind texture
	gl.BindTexture(gl.TEXTURE_2D, 0)
	
	// Cache the texture
	tm.textures["checkerboard"] = textureID
	
	fmt.Printf("Generated checkerboard texture (ID: %d, Size: %dx%d)\n", textureID, size, size)
	return textureID, nil
}

// BindTexture binds a texture to the specified texture unit
func (tm *TextureManager) BindTexture(textureID uint32, unit int32) {
	gl.ActiveTexture(gl.TEXTURE0 + uint32(unit))
	gl.BindTexture(gl.TEXTURE_2D, textureID)
}