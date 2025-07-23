package math

import (
	"math"
)

const (
	Pi      = float32(math.Pi)
	TwoPi   = float32(2 * math.Pi)
	HalfPi  = float32(math.Pi / 2)
	Deg2Rad = float32(math.Pi / 180)
	Rad2Deg = float32(180 / math.Pi)
)

func Radians(degrees float32) float32 {
	return degrees * Deg2Rad
}

func Degrees(radians float32) float32 {
	return radians * Rad2Deg
}

func Clamp(value, min, max float32) float32 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func Lerp(a, b, t float32) float32 {
	return a + (b-a)*t
}

func Min(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

func Max(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

func Abs(x float32) float32 {
	if x < 0 {
		return -x
	}
	return x
}

func Sign(x float32) float32 {
	if x < 0 {
		return -1
	}
	if x > 0 {
		return 1
	}
	return 0
}