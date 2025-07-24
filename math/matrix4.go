package math

import (
	"math"
)

type Matrix4 [16]float32

func NewMatrix4Identity() Matrix4 {
	return Matrix4{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
}

func (m Matrix4) Multiply(other Matrix4) Matrix4 {
	var result Matrix4
	
	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			sum := float32(0)
			for i := 0; i < 4; i++ {
				sum += m[row*4+i] * other[i*4+col]
			}
			result[row*4+col] = sum
		}
	}
	
	return result
}

func (m Matrix4) MultiplyVector3(v Vector3, w float32) Vector3 {
	return Vector3{
		X: m[0]*v.X + m[1]*v.Y + m[2]*v.Z + m[3]*w,
		Y: m[4]*v.X + m[5]*v.Y + m[6]*v.Z + m[7]*w,
		Z: m[8]*v.X + m[9]*v.Y + m[10]*v.Z + m[11]*w,
	}
}

func (m Matrix4) MultiplyVector4(v Vector4) Vector4 {
	return Vector4{
		X: m[0]*v.X + m[1]*v.Y + m[2]*v.Z + m[3]*v.W,
		Y: m[4]*v.X + m[5]*v.Y + m[6]*v.Z + m[7]*v.W,
		Z: m[8]*v.X + m[9]*v.Y + m[10]*v.Z + m[11]*v.W,
		W: m[12]*v.X + m[13]*v.Y + m[14]*v.Z + m[15]*v.W,
	}
}

func (m Matrix4) Inverse() Matrix4 {
	// Calculate determinant and adjugate matrix for 4x4 matrix inversion
	// Using Cramer's rule for matrix inversion
	
	// Calculate 2x2 determinants for adjugate matrix
	det2_01_01 := m[0]*m[5]  - m[1]*m[4]
	det2_01_02 := m[0]*m[6]  - m[2]*m[4]  
	det2_01_03 := m[0]*m[7]  - m[3]*m[4]
	det2_01_12 := m[1]*m[6]  - m[2]*m[5]
	det2_01_13 := m[1]*m[7]  - m[3]*m[5]
	det2_01_23 := m[2]*m[7]  - m[3]*m[6]
	
	// Calculate 3x3 determinants
	det3_201_012 := m[8]*det2_01_12 - m[9]*det2_01_02 + m[10]*det2_01_01
	det3_201_013 := m[8]*det2_01_13 - m[9]*det2_01_03 + m[11]*det2_01_01
	det3_201_023 := m[8]*det2_01_23 - m[10]*det2_01_03 + m[11]*det2_01_02
	det3_201_123 := m[9]*det2_01_23 - m[10]*det2_01_13 + m[11]*det2_01_12
	
	// Calculate the determinant of the matrix
	det := -det3_201_123*m[3] + det3_201_023*m[2] - det3_201_013*m[1] + det3_201_012*m[0]
	
	if det == 0 {
		// Matrix is singular, return identity
		return NewMatrix4Identity()
	}
	
	invDet := 1.0 / det
	
	// Calculate more 2x2 determinants needed for adjugate
	det2_03_01 := m[0]*m[13] - m[1]*m[12]
	det2_03_02 := m[0]*m[14] - m[2]*m[12]
	det2_03_03 := m[0]*m[15] - m[3]*m[12]
	det2_03_12 := m[1]*m[14] - m[2]*m[13]
	det2_03_13 := m[1]*m[15] - m[3]*m[13]
	det2_03_23 := m[2]*m[15] - m[3]*m[14]
	
	det2_13_01 := m[4]*m[13] - m[5]*m[12]
	det2_13_02 := m[4]*m[14] - m[6]*m[12]
	det2_13_03 := m[4]*m[15] - m[7]*m[12]
	det2_13_12 := m[5]*m[14] - m[6]*m[13]
	det2_13_13 := m[5]*m[15] - m[7]*m[13]
	det2_13_23 := m[6]*m[15] - m[7]*m[14]
	
	// Calculate adjugate matrix elements
	return Matrix4{
		 det3_201_123 * invDet,
		-det3_201_023 * invDet,
		 det3_201_013 * invDet,
		-det3_201_012 * invDet,
		
		-(m[8]*det2_03_23 - m[10]*det2_03_13 + m[11]*det2_03_12) * invDet,
		 (m[8]*det2_03_03 - m[9]*det2_03_13 + m[11]*det2_03_01) * invDet,
		-(m[8]*det2_03_02 - m[9]*det2_03_12 + m[10]*det2_03_01) * invDet,
		 (m[9]*det2_03_23 - m[10]*det2_03_13 + m[11]*det2_03_12) * invDet,
		
		 (m[8]*det2_13_23 - m[10]*det2_13_13 + m[11]*det2_13_12) * invDet,
		-(m[8]*det2_13_03 - m[9]*det2_13_13 + m[11]*det2_13_01) * invDet,
		 (m[8]*det2_13_02 - m[9]*det2_13_12 + m[10]*det2_13_01) * invDet,
		-(m[9]*det2_13_23 - m[10]*det2_13_13 + m[11]*det2_13_12) * invDet,
		
		-(m[4]*det2_01_23 - m[6]*det2_01_13 + m[7]*det2_01_12) * invDet,
		 (m[4]*det2_01_03 - m[5]*det2_01_13 + m[7]*det2_01_01) * invDet,
		-(m[4]*det2_01_02 - m[5]*det2_01_12 + m[6]*det2_01_01) * invDet,
		 (m[5]*det2_01_23 - m[6]*det2_01_13 + m[7]*det2_01_12) * invDet,
	}
}

func (m Matrix4) Transpose() Matrix4 {
	return Matrix4{
		m[0], m[4], m[8], m[12],
		m[1], m[5], m[9], m[13],
		m[2], m[6], m[10], m[14],
		m[3], m[7], m[11], m[15],
	}
}

func NewTranslationMatrix(x, y, z float32) Matrix4 {
	return Matrix4{
		1, 0, 0, x,
		0, 1, 0, y,
		0, 0, 1, z,
		0, 0, 0, 1,
	}
}

func NewScaleMatrix(x, y, z float32) Matrix4 {
	return Matrix4{
		x, 0, 0, 0,
		0, y, 0, 0,
		0, 0, z, 0,
		0, 0, 0, 1,
	}
}

func NewRotationX(angle float32) Matrix4 {
	cos := float32(math.Cos(float64(angle)))
	sin := float32(math.Sin(float64(angle)))
	
	return Matrix4{
		1, 0, 0, 0,
		0, cos, -sin, 0,
		0, sin, cos, 0,
		0, 0, 0, 1,
	}
}

func NewRotationY(angle float32) Matrix4 {
	cos := float32(math.Cos(float64(angle)))
	sin := float32(math.Sin(float64(angle)))
	
	return Matrix4{
		cos, 0, sin, 0,
		0, 1, 0, 0,
		-sin, 0, cos, 0,
		0, 0, 0, 1,
	}
}

func NewRotationZ(angle float32) Matrix4 {
	cos := float32(math.Cos(float64(angle)))
	sin := float32(math.Sin(float64(angle)))
	
	return Matrix4{
		cos, -sin, 0, 0,
		sin, cos, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
}

func NewLookAt(eye, target, up Vector3) Matrix4 {
	zAxis := eye.Sub(target).Normalize()
	xAxis := up.Cross(zAxis).Normalize()
	yAxis := zAxis.Cross(xAxis)
	
	// OpenGL column-major look-at matrix
	return Matrix4{
		xAxis.X, yAxis.X, zAxis.X, 0,
		xAxis.Y, yAxis.Y, zAxis.Y, 0,
		xAxis.Z, yAxis.Z, zAxis.Z, 0,
		-xAxis.Dot(eye), -yAxis.Dot(eye), -zAxis.Dot(eye), 1,
	}
}

func NewPerspective(fov, aspect, near, far float32) Matrix4 {
	tanHalfFov := float32(math.Tan(float64(fov / 2)))
	
	// OpenGL column-major perspective matrix
	return Matrix4{
		1 / (aspect * tanHalfFov), 0, 0, 0,
		0, 1 / tanHalfFov, 0, 0,
		0, 0, -(far + near) / (far - near), -(2 * far * near) / (far - near),
		0, 0, -1, 0,
	}
}

func NewOrthographic(left, right, bottom, top, near, far float32) Matrix4 {
	return Matrix4{
		2 / (right - left), 0, 0, -(right + left) / (right - left),
		0, 2 / (top - bottom), 0, -(top + bottom) / (top - bottom),
		0, 0, -2 / (far - near), -(far + near) / (far - near),
		0, 0, 0, 1,
	}
}