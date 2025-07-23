package math

import (
	"math"
)

type Vector2 struct {
	X, Y float32
}

func NewVector2(x, y float32) Vector2 {
	return Vector2{X: x, Y: y}
}

func (v Vector2) Add(other Vector2) Vector2 {
	return Vector2{v.X + other.X, v.Y + other.Y}
}

func (v Vector2) Sub(other Vector2) Vector2 {
	return Vector2{v.X - other.X, v.Y - other.Y}
}

func (v Vector2) Mul(scalar float32) Vector2 {
	return Vector2{v.X * scalar, v.Y * scalar}
}

func (v Vector2) Div(scalar float32) Vector2 {
	return Vector2{v.X / scalar, v.Y / scalar}
}

func (v Vector2) Dot(other Vector2) float32 {
	return v.X*other.X + v.Y*other.Y
}

func (v Vector2) Length() float32 {
	return float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y)))
}

func (v Vector2) LengthSquared() float32 {
	return v.X*v.X + v.Y*v.Y
}

func (v Vector2) Normalize() Vector2 {
	length := v.Length()
	if length == 0 {
		return Vector2{0, 0}
	}
	return v.Div(length)
}

func (v Vector2) Distance(other Vector2) float32 {
	return v.Sub(other).Length()
}

func Lerp2(a, b Vector2, t float32) Vector2 {
	return Vector2{
		X: a.X + (b.X-a.X)*t,
		Y: a.Y + (b.Y-a.Y)*t,
	}
}