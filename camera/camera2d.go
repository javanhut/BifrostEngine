package camera

import (
	bmath "github.com/javanhut/BifrostEngine/m/v2/math"
)

type Camera2D struct {
	position   bmath.Vector2
	zoom       float32
	rotation   float32
	
	viewportWidth  float32
	viewportHeight float32
	
	viewMatrix       bmath.Matrix4
	projectionMatrix bmath.Matrix4
	dirty            bool
}

func NewCamera2D(viewportWidth, viewportHeight float32) *Camera2D {
	cam := &Camera2D{
		position:       bmath.NewVector2(0, 0),
		zoom:          1.0,
		rotation:      0,
		viewportWidth: viewportWidth,
		viewportHeight: viewportHeight,
		dirty:         true,
	}
	cam.updateMatrices()
	return cam
}

func (c *Camera2D) updateMatrices() {
	if c.dirty {
		// Build view matrix (inverse of camera transform)
		translation := bmath.NewTranslationMatrix(-c.position.X, -c.position.Y, 0)
		rotation := bmath.NewRotationZ(-c.rotation)
		scale := bmath.NewScaleMatrix(c.zoom, c.zoom, 1)
		
		c.viewMatrix = scale.Multiply(rotation).Multiply(translation)
		
		// Orthographic projection for 2D
		halfWidth := c.viewportWidth / 2
		halfHeight := c.viewportHeight / 2
		c.projectionMatrix = bmath.NewOrthographic(-halfWidth, halfWidth, -halfHeight, halfHeight, -1, 1)
		
		c.dirty = false
	}
}

func (c *Camera2D) GetViewMatrix() bmath.Matrix4 {
	c.updateMatrices()
	return c.viewMatrix
}

func (c *Camera2D) GetProjectionMatrix() bmath.Matrix4 {
	c.updateMatrices()
	return c.projectionMatrix
}

func (c *Camera2D) GetViewProjectionMatrix() bmath.Matrix4 {
	c.updateMatrices()
	return c.projectionMatrix.Multiply(c.viewMatrix)
}

func (c *Camera2D) SetPosition(pos bmath.Vector2) {
	c.position = pos
	c.dirty = true
}

func (c *Camera2D) Move(delta bmath.Vector2) {
	c.position = c.position.Add(delta)
	c.dirty = true
}

func (c *Camera2D) SetZoom(zoom float32) {
	c.zoom = bmath.Clamp(zoom, 0.1, 10.0)
	c.dirty = true
}

func (c *Camera2D) Zoom(delta float32) {
	c.SetZoom(c.zoom + delta)
}

func (c *Camera2D) SetRotation(rotation float32) {
	c.rotation = rotation
	c.dirty = true
}

func (c *Camera2D) Rotate(delta float32) {
	c.rotation += delta
	c.dirty = true
}

func (c *Camera2D) SetViewportSize(width, height float32) {
	c.viewportWidth = width
	c.viewportHeight = height
	c.dirty = true
}

func (c *Camera2D) ScreenToWorld(screenPos bmath.Vector2) bmath.Vector2 {
	// Convert screen coordinates to normalized device coordinates
	ndcX := (screenPos.X / c.viewportWidth) * 2.0 - 1.0
	ndcY := 1.0 - (screenPos.Y / c.viewportHeight) * 2.0
	
	// Note: In a real implementation, you'd want to add matrix inversion to the math library
	// For now, we'll do a simplified version assuming no rotation
	
	worldX := (ndcX * c.viewportWidth / 2) / c.zoom + c.position.X
	worldY := (ndcY * c.viewportHeight / 2) / c.zoom + c.position.Y
	
	return bmath.NewVector2(worldX, worldY)
}

func (c *Camera2D) Update(deltaTime float32) {
	// Override for specific camera behaviors like smooth follow
}