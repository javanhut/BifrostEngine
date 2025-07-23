package scene

import (
	bmath "github.com/javanhut/BifrostEngine/m/v2/math"
)

// Transform component handles position, rotation, and scale
type Transform struct {
	Position     bmath.Vector3
	Rotation     bmath.Vector3 // Euler angles in degrees
	Scale        bmath.Vector3
	LocalMatrix  bmath.Matrix4
	WorldMatrix  bmath.Matrix4
	Parent       *Transform
	Children     []*Transform
	dirty        bool
}

// NewTransform creates a new transform component
func NewTransform() *Transform {
	return &Transform{
		Position:    bmath.NewVector3(0, 0, 0),
		Rotation:    bmath.NewVector3(0, 0, 0),
		Scale:       bmath.NewVector3(1, 1, 1),
		LocalMatrix: bmath.NewMatrix4Identity(),
		WorldMatrix: bmath.NewMatrix4Identity(),
		Children:    []*Transform{},
		dirty:       true,
	}
}

// GetType returns the component type
func (t *Transform) GetType() string {
	return "Transform"
}

// Update updates the transform matrices
func (t *Transform) Update(deltaTime float32) {
	if t.dirty {
		t.updateMatrices()
	}
}

// SetPosition sets the position and marks transform as dirty
func (t *Transform) SetPosition(pos bmath.Vector3) {
	t.Position = pos
	t.markDirty()
}

// SetRotation sets the rotation and marks transform as dirty
func (t *Transform) SetRotation(rot bmath.Vector3) {
	t.Rotation = rot
	t.markDirty()
}

// SetScale sets the scale and marks transform as dirty
func (t *Transform) SetScale(scale bmath.Vector3) {
	t.Scale = scale
	t.markDirty()
}

// Translate moves the transform by a delta
func (t *Transform) Translate(delta bmath.Vector3) {
	t.Position = t.Position.Add(delta)
	t.markDirty()
}

// Rotate rotates the transform by delta angles
func (t *Transform) Rotate(delta bmath.Vector3) {
	t.Rotation = t.Rotation.Add(delta)
	t.markDirty()
}

// SetParent sets the parent transform
func (t *Transform) SetParent(parent *Transform) {
	// Remove from old parent
	if t.Parent != nil {
		t.Parent.removeChild(t)
	}
	
	// Set new parent
	t.Parent = parent
	if parent != nil {
		parent.Children = append(parent.Children, t)
	}
	
	t.markDirty()
}

// GetWorldPosition returns the world position
func (t *Transform) GetWorldPosition() bmath.Vector3 {
	t.updateMatrices()
	return bmath.NewVector3(t.WorldMatrix[3], t.WorldMatrix[7], t.WorldMatrix[11])
}

// GetLocalMatrix returns the local transformation matrix
func (t *Transform) GetLocalMatrix() bmath.Matrix4 {
	t.updateMatrices()
	return t.LocalMatrix
}

// GetWorldMatrix returns the world transformation matrix
func (t *Transform) GetWorldMatrix() bmath.Matrix4 {
	t.updateMatrices()
	return t.WorldMatrix
}

// markDirty marks this transform and all children as needing update
func (t *Transform) markDirty() {
	t.dirty = true
	for _, child := range t.Children {
		child.markDirty()
	}
}

// updateMatrices recalculates the transformation matrices
func (t *Transform) updateMatrices() {
	if !t.dirty {
		return
	}
	
	// Build local matrix: T * R * S
	translation := bmath.NewTranslationMatrix(t.Position.X, t.Position.Y, t.Position.Z)
	rotationX := bmath.NewRotationX(bmath.Radians(t.Rotation.X))
	rotationY := bmath.NewRotationY(bmath.Radians(t.Rotation.Y))
	rotationZ := bmath.NewRotationZ(bmath.Radians(t.Rotation.Z))
	scale := bmath.NewScaleMatrix(t.Scale.X, t.Scale.Y, t.Scale.Z)
	
	rotation := rotationZ.Multiply(rotationY).Multiply(rotationX)
	t.LocalMatrix = translation.Multiply(rotation).Multiply(scale)
	
	// Calculate world matrix
	if t.Parent != nil {
		t.Parent.updateMatrices() // Ensure parent is updated
		t.WorldMatrix = t.Parent.WorldMatrix.Multiply(t.LocalMatrix)
	} else {
		t.WorldMatrix = t.LocalMatrix
	}
	
	t.dirty = false
}

// removeChild removes a child from the children list
func (t *Transform) removeChild(child *Transform) {
	for i, c := range t.Children {
		if c == child {
			t.Children = append(t.Children[:i], t.Children[i+1:]...)
			break
		}
	}
}