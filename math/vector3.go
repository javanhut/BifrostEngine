package math

import (
	"math"
)

type Vector3 struct {
	X, Y, Z float32
}

func NewVector3(x, y, z float32) Vector3 {
	return Vector3{X: x, Y: y, Z: z}
}

func (v Vector3) Add(other Vector3) Vector3 {
	return Vector3{v.X + other.X, v.Y + other.Y, v.Z + other.Z}
}

func (v Vector3) Sub(other Vector3) Vector3 {
	return Vector3{v.X - other.X, v.Y - other.Y, v.Z - other.Z}
}

func (v Vector3) Mul(scalar float32) Vector3 {
	return Vector3{v.X * scalar, v.Y * scalar, v.Z * scalar}
}

func (v Vector3) Div(scalar float32) Vector3 {
	return Vector3{v.X / scalar, v.Y / scalar, v.Z / scalar}
}

func (v Vector3) Dot(other Vector3) float32 {
	return v.X*other.X + v.Y*other.Y + v.Z*other.Z
}

func (v Vector3) Cross(other Vector3) Vector3 {
	return Vector3{
		X: v.Y*other.Z - v.Z*other.Y,
		Y: v.Z*other.X - v.X*other.Z,
		Z: v.X*other.Y - v.Y*other.X,
	}
}

func (v Vector3) Length() float32 {
	return float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y + v.Z*v.Z)))
}

func (v Vector3) LengthSquared() float32 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

func (v Vector3) Normalize() Vector3 {
	length := v.Length()
	if length == 0 {
		return Vector3{0, 0, 0}
	}
	return v.Div(length)
}

func (v Vector3) Distance(other Vector3) float32 {
	return v.Sub(other).Length()
}

func Lerp3(a, b Vector3, t float32) Vector3 {
	return Vector3{
		X: a.X + (b.X-a.X)*t,
		Y: a.Y + (b.Y-a.Y)*t,
		Z: a.Z + (b.Z-a.Z)*t,
	}
}

var (
	Vector3Zero    = Vector3{0, 0, 0}
	Vector3One     = Vector3{1, 1, 1}
	Vector3Up      = Vector3{0, 1, 0}
	Vector3Down    = Vector3{0, -1, 0}
	Vector3Left    = Vector3{-1, 0, 0}
	Vector3Right   = Vector3{1, 0, 0}
	Vector3Forward = Vector3{0, 0, -1}
	Vector3Back    = Vector3{0, 0, 1}
)