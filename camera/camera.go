package camera

import (
	"math"
	
	bmath "github.com/javanhut/BifrostEngine/m/v2/math"
)

type Camera interface {
	GetViewMatrix() bmath.Matrix4
	GetProjectionMatrix() bmath.Matrix4
	GetViewProjectionMatrix() bmath.Matrix4
	Update(deltaTime float32)
}

type Camera3D struct {
	position   bmath.Vector3
	target     bmath.Vector3
	up         bmath.Vector3
	
	fov        float32
	aspect     float32
	near       float32
	far        float32
	
	viewMatrix       bmath.Matrix4
	projectionMatrix bmath.Matrix4
	dirty            bool
}

func NewCamera3D(position, target bmath.Vector3, fov, aspect, near, far float32) *Camera3D {
	cam := &Camera3D{
		position: position,
		target:   target,
		up:       bmath.Vector3Up,
		fov:      fov,
		aspect:   aspect,
		near:     near,
		far:      far,
		dirty:    true,
	}
	cam.updateMatrices()
	return cam
}

func (c *Camera3D) updateMatrices() {
	if c.dirty {
		c.viewMatrix = bmath.NewLookAt(c.position, c.target, c.up)
		c.projectionMatrix = bmath.NewPerspective(c.fov, c.aspect, c.near, c.far)
		c.dirty = false
	}
}

func (c *Camera3D) GetViewMatrix() bmath.Matrix4 {
	c.updateMatrices()
	return c.viewMatrix
}

func (c *Camera3D) GetProjectionMatrix() bmath.Matrix4 {
	c.updateMatrices()
	return c.projectionMatrix
}

func (c *Camera3D) GetViewProjectionMatrix() bmath.Matrix4 {
	c.updateMatrices()
	return c.projectionMatrix.Multiply(c.viewMatrix)
}

func (c *Camera3D) SetPosition(pos bmath.Vector3) {
	c.position = pos
	c.dirty = true
}

func (c *Camera3D) SetTarget(target bmath.Vector3) {
	c.target = target
	c.dirty = true
}

func (c *Camera3D) GetPosition() bmath.Vector3 {
	return c.position
}

func (c *Camera3D) GetTarget() bmath.Vector3 {
	return c.target
}

func (c *Camera3D) SetAspect(aspect float32) {
	c.aspect = aspect
	c.dirty = true
}

func (c *Camera3D) Move(delta bmath.Vector3) {
	c.position = c.position.Add(delta)
	c.target = c.target.Add(delta)
	c.dirty = true
}

func (c *Camera3D) Rotate(yaw, pitch float32) {
	direction := c.target.Sub(c.position)
	distance := direction.Length()
	
	// Current angles
	currentYaw := float32(math.Atan2(float64(direction.X), float64(direction.Z)))
	currentPitch := float32(math.Asin(float64(direction.Y / distance)))
	
	// New angles
	newYaw := currentYaw + yaw
	newPitch := bmath.Clamp(currentPitch+pitch, -bmath.HalfPi+0.01, bmath.HalfPi-0.01)
	
	// Calculate new direction
	direction = bmath.Vector3{
		X: distance * float32(math.Cos(float64(newPitch)) * math.Sin(float64(newYaw))),
		Y: distance * float32(math.Sin(float64(newPitch))),
		Z: distance * float32(math.Cos(float64(newPitch)) * math.Cos(float64(newYaw))),
	}
	
	c.target = c.position.Add(direction)
	c.dirty = true
}

func (c *Camera3D) Update(deltaTime float32) {
	// Override in subclasses for specific camera behaviors
}