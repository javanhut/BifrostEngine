package math

import (
	"math"
)

type Vector4 struct {
	X, Y, Z, W float32
}

func NewVector4(x, y, z, w float32) Vector4 {
	return Vector4{X: x, Y: y, Z: z, W: w}
}

func (v Vector4) Add(other Vector4) Vector4 {
	return Vector4{
		X: v.X + other.X,
		Y: v.Y + other.Y,
		Z: v.Z + other.Z,
		W: v.W + other.W,
	}
}

func (v Vector4) Sub(other Vector4) Vector4 {
	return Vector4{
		X: v.X - other.X,
		Y: v.Y - other.Y,
		Z: v.Z - other.Z,
		W: v.W - other.W,
	}
}

func (v Vector4) Scale(scalar float32) Vector4 {
	return Vector4{
		X: v.X * scalar,
		Y: v.Y * scalar,
		Z: v.Z * scalar,
		W: v.W * scalar,
	}
}

func (v Vector4) Dot(other Vector4) float32 {
	return v.X*other.X + v.Y*other.Y + v.Z*other.Z + v.W*other.W
}

func (v Vector4) Length() float32 {
	return float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y + v.Z*v.Z + v.W*v.W)))
}

func (v Vector4) Normalize() Vector4 {
	length := v.Length()
	if length == 0 {
		return Vector4{0, 0, 0, 0}
	}
	return Vector4{
		X: v.X / length,
		Y: v.Y / length,
		Z: v.Z / length,
		W: v.W / length,
	}
}