package rl

import (
	"math"
)

// Vector2Zero - Vector with components value 0.0
func Vector2Zero() Vector2 {
	return NewVector2(0.0, 0.0)
}

// Vector2One - Vector with components value 1.0
func Vector2One() Vector2 {
	return NewVector2(1.0, 1.0)
}

// Vector2Add - Add two vectors (v1 + v2)
func Vector2Add(v1, v2 Vector2) Vector2 {
	return NewVector2(v1.X+v2.X, v1.Y+v2.Y)
}

// Vector2Subtract - Subtract two vectors (v1 - v2)
func Vector2Subtract(v1, v2 Vector2) Vector2 {
	return NewVector2(v1.X-v2.X, v1.Y-v2.Y)
}

// Vector2Length - Calculate vector length
func Vector2Length(v Vector2) float32 {
	return float32(math.Sqrt(float64((v.X * v.X) + (v.Y * v.Y))))
}

// Vector2DotProduct - Calculate two vectors dot product
func Vector2DotProduct(v1, v2 Vector2) float32 {
	return v1.X*v2.X + v1.Y*v2.Y
}

// Vector2Distance - Calculate distance between two vectors
func Vector2Distance(v1, v2 Vector2) float32 {
	return float32(math.Sqrt(float64((v1.X-v2.X)*(v1.X-v2.X) + (v1.Y-v2.Y)*(v1.Y-v2.Y))))
}

// Vector2Angle - Calculate angle between two vectors in X-axis
func Vector2Angle(v1, v2 Vector2) float32 {
	angle := float32(math.Atan2(float64(v2.Y-v1.Y), float64(v2.X-v1.X)) * (180.0 / float64(Pi)))

	if angle < 0 {
		angle += 360.0
	}

	return angle
}

// Vector2Scale - Scale vector (multiply by value)
func Vector2Scale(v Vector2, scale float32) Vector2 {
	return NewVector2(v.X*scale, v.Y*scale)
}

// Vector2Multiply - Multiply vector by vector
func Vector2Multiply(v1, v2 Vector2) Vector2 {
	return NewVector2(v1.X*v2.X, v1.Y*v2.Y)
}

// Vector2Negate - Negate vector
func Vector2Negate(v Vector2) Vector2 {
	return NewVector2(-v.X, -v.Y)
}

// Vector2Divide - Divide vector by vector
func Vector2DivideV(v1, v2 Vector2) Vector2 {
	return NewVector2(v1.X/v2.X, v1.Y/v2.Y)
}

// Vector2Normalize - Normalize provided vector
func Vector2Normalize(v Vector2) Vector2 {
	return Vector2Scale(v, 1/Vector2Length(v))
}

// Vector2Lerp - Calculate linear interpolation between two vectors
func Vector2Lerp(v1, v2 Vector2, amount float32) Vector2 {
	return NewVector2(v1.X+amount*(v2.X-v1.X), v1.Y+amount*(v2.Y-v1.Y))
}

// Vector2CrossProduct - Calculate two vectors cross product
func Vector2CrossProduct(v1, v2 Vector2) float32 {
	return v1.X*v2.Y - v1.Y*v2.X
}

// Vector2Cross - Calculate the cross product of a vector and a value
func Vector2Cross(value float32, vector Vector2) Vector2 {
	return NewVector2(-value*vector.Y, value*vector.X)
}

// Vector2LenSqr - Returns the len square root of a vector
func Vector2LenSqr(vector Vector2) float32 {
	return vector.X*vector.X + vector.Y*vector.Y
}

// Mat2Radians - Creates a matrix 2x2 from a given radians value
func Mat2Radians(radians float32) Mat2 {
	c := float32(math.Cos(float64(radians)))
	s := float32(math.Sin(float64(radians)))

	return NewMat2(c, -s, s, c)
}

// Mat2Set - Set values from radians to a created matrix 2x2
func Mat2Set(matrix *Mat2, radians float32) {
	cos := float32(math.Cos(float64(radians)))
	sin := float32(math.Sin(float64(radians)))

	matrix.M00 = cos
	matrix.M01 = -sin
	matrix.M10 = sin
	matrix.M11 = cos
}

// Mat2Transpose - Returns the transpose of a given matrix 2x2
func Mat2Transpose(matrix Mat2) Mat2 {
	return NewMat2(matrix.M00, matrix.M10, matrix.M01, matrix.M11)
}

// Mat2MultiplyVector2 - Multiplies a vector by a matrix 2x2
func Mat2MultiplyVector2(matrix Mat2, vector Vector2) Vector2 {
	return NewVector2(matrix.M00*vector.X+matrix.M01*vector.Y, matrix.M10*vector.X+matrix.M11*vector.Y)
}

// Vector3Zero - Vector with components value 0.0
func Vector3Zero() Vector3 {
	return NewVector3(0.0, 0.0, 0.0)
}

// Vector3One - Vector with components value 1.0
func Vector3One() Vector3 {
	return NewVector3(1.0, 1.0, 1.0)
}

// Vector3Add - Add two vectors
func Vector3Add(v1, v2 Vector3) Vector3 {
	return NewVector3(v1.X+v2.X, v1.Y+v2.Y, v1.Z+v2.Z)
}

// Vector3Multiply - Multiply vector by scalar
func Vector3Multiply(v Vector3, scalar float32) Vector3 {
	result := Vector3{}

	result.X = v.X * scalar
	result.Y = v.Y * scalar
	result.Z = v.Z * scalar

	return result
}

// Vector3MultiplyV - Multiply vector by vector
func Vector3MultiplyV(v1, v2 Vector3) Vector3 {
	result := Vector3{}

	result.X = v1.X * v2.X
	result.Y = v1.Y * v2.Y
	result.Z = v1.Z * v2.Z

	return result
}

// Vector3Subtract - Subtract two vectors
func Vector3Subtract(v1, v2 Vector3) Vector3 {
	return NewVector3(v1.X-v2.X, v1.Y-v2.Y, v1.Z-v2.Z)
}

// Vector3CrossProduct - Calculate two vectors cross product
func Vector3CrossProduct(v1, v2 Vector3) Vector3 {
	result := Vector3{}

	result.X = v1.Y*v2.Z - v1.Z*v2.Y
	result.Y = v1.Z*v2.X - v1.X*v2.Z
	result.Z = v1.X*v2.Y - v1.Y*v2.X

	return result
}

// Vector3Perpendicular - Calculate one vector perpendicular vector
func Vector3Perpendicular(v Vector3) Vector3 {
	result := Vector3{}

	min := math.Abs(float64(v.X))
	cardinalAxis := NewVector3(1.0, 0.0, 0.0)

	if math.Abs(float64(v.Y)) < min {
		min = math.Abs(float64(v.Y))
		cardinalAxis = NewVector3(0.0, 1.0, 0.0)
	}

	if math.Abs(float64(v.Z)) < min {
		cardinalAxis = NewVector3(0.0, 0.0, 1.0)
	}

	result = Vector3CrossProduct(v, cardinalAxis)

	return result
}

// Vector3Length - Calculate vector length
func Vector3Length(v Vector3) float32 {
	return float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y + v.Z*v.Z)))
}

// Vector3DotProduct - Calculate two vectors dot product
func Vector3DotProduct(v1, v2 Vector3) float32 {
	return v1.X*v2.X + v1.Y*v2.Y + v1.Z*v2.Z
}

// Vector3Distance - Calculate distance between two vectors
func Vector3Distance(v1, v2 Vector3) float32 {
	dx := v2.X - v1.X
	dy := v2.Y - v1.Y
	dz := v2.Z - v1.Z

	return float32(math.Sqrt(float64(dx*dx + dy*dy + dz*dz)))
}

// Vector3Scale - Scale provided vector
func Vector3Scale(v Vector3, scale float32) Vector3 {
	return NewVector3(v.X*scale, v.Y*scale, v.Z*scale)
}

// Vector3Negate - Negate provided vector (invert direction)
func Vector3Negate(v Vector3) Vector3 {
	return NewVector3(-v.X, -v.Y, -v.Z)
}

// Vector3Normalize - Normalize provided vector
func Vector3Normalize(v Vector3) Vector3 {
	result := v

	var length, ilength float32

	length = Vector3Length(v)

	if length == 0 {
		length = 1.0
	}

	ilength = 1.0 / length

	result.X *= ilength
	result.Y *= ilength
	result.Z *= ilength

	return result
}

// Vector3Transform - Transforms a Vector3 by a given Matrix
func Vector3Transform(v Vector3, mat Matrix) Vector3 {
	result := Vector3{}

	x := v.X
	y := v.Y
	z := v.Z

	result.X = mat.M0*x + mat.M4*y + mat.M8*z + mat.M12
	result.Y = mat.M1*x + mat.M5*y + mat.M9*z + mat.M13
	result.Z = mat.M2*x + mat.M6*y + mat.M10*z + mat.M14

	return result
}

// Vector3Lerp - Calculate linear interpolation between two vectors
func Vector3Lerp(v1, v2 Vector3, amount float32) Vector3 {
	result := Vector3{}

	result.X = v1.X + amount*(v2.X-v1.X)
	result.Y = v1.Y + amount*(v2.Y-v1.Y)
	result.Z = v1.Z + amount*(v2.Z-v1.Z)

	return result
}

// Vector3Reflect - Calculate reflected vector to normal
func Vector3Reflect(vector, normal Vector3) Vector3 {
	// I is the original vector
	// N is the normal of the incident plane
	// R = I - (2*N*( DotProduct[ I,N] ))

	result := Vector3{}

	dotProduct := Vector3DotProduct(vector, normal)

	result.X = vector.X - (2.0*normal.X)*dotProduct
	result.Y = vector.Y - (2.0*normal.Y)*dotProduct
	result.Z = vector.Z - (2.0*normal.Z)*dotProduct

	return result
}

// Vector3Min - Return min value for each pair of components
func Vector3Min(vec1, vec2 Vector3) Vector3 {
	result := Vector3{}

	result.X = float32(math.Min(float64(vec1.X), float64(vec2.X)))
	result.Y = float32(math.Min(float64(vec1.Y), float64(vec2.Y)))
	result.Z = float32(math.Min(float64(vec1.Z), float64(vec2.Z)))

	return result
}

// Vector3Max - Return max value for each pair of components
func Vector3Max(vec1, vec2 Vector3) Vector3 {
	result := Vector3{}

	result.X = float32(math.Max(float64(vec1.X), float64(vec2.X)))
	result.Y = float32(math.Max(float64(vec1.Y), float64(vec2.Y)))
	result.Z = float32(math.Max(float64(vec1.Z), float64(vec2.Z)))

	return result
}

// Vector3Barycenter - Barycenter coords for p in triangle abc
func Vector3Barycenter(p, a, b, c Vector3) Vector3 {
	v0 := Vector3Subtract(b, a)
	v1 := Vector3Subtract(c, a)
	v2 := Vector3Subtract(p, a)
	d00 := Vector3DotProduct(v0, v0)
	d01 := Vector3DotProduct(v0, v1)
	d11 := Vector3DotProduct(v1, v1)
	d20 := Vector3DotProduct(v2, v0)
	d21 := Vector3DotProduct(v2, v1)

	denom := d00*d11 - d01*d01

	result := Vector3{}

	result.Y = (d11*d20 - d01*d21) / denom
	result.Z = (d00*d21 - d01*d20) / denom
	result.X = 1.0 - (result.Z + result.Y)

	return result
}

// MatrixDeterminant - Compute matrix determinant
func MatrixDeterminant(mat Matrix) float32 {
	var result float32

	a00 := mat.M0
	a01 := mat.M1
	a02 := mat.M2
	a03 := mat.M3
	a10 := mat.M4
	a11 := mat.M5
	a12 := mat.M6
	a13 := mat.M7
	a20 := mat.M8
	a21 := mat.M9
	a22 := mat.M10
	a23 := mat.M11
	a30 := mat.M12
	a31 := mat.M13
	a32 := mat.M14
	a33 := mat.M15

	result = a30*a21*a12*a03 - a20*a31*a12*a03 - a30*a11*a22*a03 + a10*a31*a22*a03 +
		a20*a11*a32*a03 - a10*a21*a32*a03 - a30*a21*a02*a13 + a20*a31*a02*a13 +
		a30*a01*a22*a13 - a00*a31*a22*a13 - a20*a01*a32*a13 + a00*a21*a32*a13 +
		a30*a11*a02*a23 - a10*a31*a02*a23 - a30*a01*a12*a23 + a00*a31*a12*a23 +
		a10*a01*a32*a23 - a00*a11*a32*a23 - a20*a11*a02*a33 + a10*a21*a02*a33 +
		a20*a01*a12*a33 - a00*a21*a12*a33 - a10*a01*a22*a33 + a00*a11*a22*a33

	return result
}

// MatrixTrace - Returns the trace of the matrix (sum of the values along the diagonal)
func MatrixTrace(mat Matrix) float32 {
	return mat.M0 + mat.M5 + mat.M10 + mat.M15
}

// MatrixTranspose - Transposes provided matrix
func MatrixTranspose(mat Matrix) Matrix {
	var result Matrix

	result.M0 = mat.M0
	result.M1 = mat.M4
	result.M2 = mat.M8
	result.M3 = mat.M12
	result.M4 = mat.M1
	result.M5 = mat.M5
	result.M6 = mat.M9
	result.M7 = mat.M13
	result.M8 = mat.M2
	result.M9 = mat.M6
	result.M10 = mat.M10
	result.M11 = mat.M14
	result.M12 = mat.M3
	result.M13 = mat.M7
	result.M14 = mat.M11
	result.M15 = mat.M15

	return result
}

// MatrixInvert - Invert provided matrix
func MatrixInvert(mat Matrix) Matrix {
	var result Matrix

	a00 := mat.M0
	a01 := mat.M1
	a02 := mat.M2
	a03 := mat.M3
	a10 := mat.M4
	a11 := mat.M5
	a12 := mat.M6
	a13 := mat.M7
	a20 := mat.M8
	a21 := mat.M9
	a22 := mat.M10
	a23 := mat.M11
	a30 := mat.M12
	a31 := mat.M13
	a32 := mat.M14
	a33 := mat.M15

	b00 := a00*a11 - a01*a10
	b01 := a00*a12 - a02*a10
	b02 := a00*a13 - a03*a10
	b03 := a01*a12 - a02*a11
	b04 := a01*a13 - a03*a11
	b05 := a02*a13 - a03*a12
	b06 := a20*a31 - a21*a30
	b07 := a20*a32 - a22*a30
	b08 := a20*a33 - a23*a30
	b09 := a21*a32 - a22*a31
	b10 := a21*a33 - a23*a31
	b11 := a22*a33 - a23*a32

	// Calculate the invert determinant (inlined to avoid double-caching)
	invDet := 1.0 / (b00*b11 - b01*b10 + b02*b09 + b03*b08 - b04*b07 + b05*b06)

	result.M0 = (a11*b11 - a12*b10 + a13*b09) * invDet
	result.M1 = (-a01*b11 + a02*b10 - a03*b09) * invDet
	result.M2 = (a31*b05 - a32*b04 + a33*b03) * invDet
	result.M3 = (-a21*b05 + a22*b04 - a23*b03) * invDet
	result.M4 = (-a10*b11 + a12*b08 - a13*b07) * invDet
	result.M5 = (a00*b11 - a02*b08 + a03*b07) * invDet
	result.M6 = (-a30*b05 + a32*b02 - a33*b01) * invDet
	result.M7 = (a20*b05 - a22*b02 + a23*b01) * invDet
	result.M8 = (a10*b10 - a11*b08 + a13*b06) * invDet
	result.M9 = (-a00*b10 + a01*b08 - a03*b06) * invDet
	result.M10 = (a30*b04 - a31*b02 + a33*b00) * invDet
	result.M11 = (-a20*b04 + a21*b02 - a23*b00) * invDet
	result.M12 = (-a10*b09 + a11*b07 - a12*b06) * invDet
	result.M13 = (a00*b09 - a01*b07 + a02*b06) * invDet
	result.M14 = (-a30*b03 + a31*b01 - a32*b00) * invDet
	result.M15 = (a20*b03 - a21*b01 + a22*b00) * invDet

	return result
}

// MatrixNormalize - Normalize provided matrix
func MatrixNormalize(mat Matrix) Matrix {
	var result Matrix

	det := MatrixDeterminant(mat)

	result.M0 /= det
	result.M1 /= det
	result.M2 /= det
	result.M3 /= det
	result.M4 /= det
	result.M5 /= det
	result.M6 /= det
	result.M7 /= det
	result.M8 /= det
	result.M9 /= det
	result.M10 /= det
	result.M11 /= det
	result.M12 /= det
	result.M13 /= det
	result.M14 /= det
	result.M15 /= det

	return result
}

// MatrixIdentity - Returns identity matrix
func MatrixIdentity() Matrix {
	return NewMatrix(
		1.0, 0.0, 0.0, 0.0,
		0.0, 1.0, 0.0, 0.0,
		0.0, 0.0, 1.0, 0.0,
		0.0, 0.0, 0.0, 1.0)
}

// MatrixAdd - Add two matrices
func MatrixAdd(left, right Matrix) Matrix {
	result := MatrixIdentity()

	result.M0 = left.M0 + right.M0
	result.M1 = left.M1 + right.M1
	result.M2 = left.M2 + right.M2
	result.M3 = left.M3 + right.M3
	result.M4 = left.M4 + right.M4
	result.M5 = left.M5 + right.M5
	result.M6 = left.M6 + right.M6
	result.M7 = left.M7 + right.M7
	result.M8 = left.M8 + right.M8
	result.M9 = left.M9 + right.M9
	result.M10 = left.M10 + right.M10
	result.M11 = left.M11 + right.M11
	result.M12 = left.M12 + right.M12
	result.M13 = left.M13 + right.M13
	result.M14 = left.M14 + right.M14
	result.M15 = left.M15 + right.M15

	return result
}

// MatrixSubtract - Subtract two matrices (left - right)
func MatrixSubtract(left, right Matrix) Matrix {
	result := MatrixIdentity()

	result.M0 = left.M0 - right.M0
	result.M1 = left.M1 - right.M1
	result.M2 = left.M2 - right.M2
	result.M3 = left.M3 - right.M3
	result.M4 = left.M4 - right.M4
	result.M5 = left.M5 - right.M5
	result.M6 = left.M6 - right.M6
	result.M7 = left.M7 - right.M7
	result.M8 = left.M8 - right.M8
	result.M9 = left.M9 - right.M9
	result.M10 = left.M10 - right.M10
	result.M11 = left.M11 - right.M11
	result.M12 = left.M12 - right.M12
	result.M13 = left.M13 - right.M13
	result.M14 = left.M14 - right.M14
	result.M15 = left.M15 - right.M15

	return result
}

// MatrixTranslate - Returns translation matrix
func MatrixTranslate(x, y, z float32) Matrix {
	return NewMatrix(
		1.0, 0.0, 0.0, x,
		0.0, 1.0, 0.0, y,
		0.0, 0.0, 1.0, z,
		0, 0, 0, 1.0)
}

// MatrixRotate - Returns rotation matrix for an angle around an specified axis (angle in radians)
func MatrixRotate(axis Vector3, angle float32) Matrix {
	var result Matrix

	mat := MatrixIdentity()

	x := axis.X
	y := axis.Y
	z := axis.Z

	length := float32(math.Sqrt(float64(x*x + y*y + z*z)))

	if length != 1.0 && length != 0.0 {
		length = 1.0 / length
		x *= length
		y *= length
		z *= length
	}

	sinres := float32(math.Sin(float64(angle)))
	cosres := float32(math.Cos(float64(angle)))
	t := 1.0 - cosres

	// Cache some matrix values (speed optimization)
	a00 := mat.M0
	a01 := mat.M1
	a02 := mat.M2
	a03 := mat.M3
	a10 := mat.M4
	a11 := mat.M5
	a12 := mat.M6
	a13 := mat.M7
	a20 := mat.M8
	a21 := mat.M9
	a22 := mat.M10
	a23 := mat.M11

	// Construct the elements of the rotation matrix
	b00 := x*x*t + cosres
	b01 := y*x*t + z*sinres
	b02 := z*x*t - y*sinres
	b10 := x*y*t - z*sinres
	b11 := y*y*t + cosres
	b12 := z*y*t + x*sinres
	b20 := x*z*t + y*sinres
	b21 := y*z*t - x*sinres
	b22 := z*z*t + cosres

	// Perform rotation-specific matrix multiplication
	result.M0 = a00*b00 + a10*b01 + a20*b02
	result.M1 = a01*b00 + a11*b01 + a21*b02
	result.M2 = a02*b00 + a12*b01 + a22*b02
	result.M3 = a03*b00 + a13*b01 + a23*b02
	result.M4 = a00*b10 + a10*b11 + a20*b12
	result.M5 = a01*b10 + a11*b11 + a21*b12
	result.M6 = a02*b10 + a12*b11 + a22*b12
	result.M7 = a03*b10 + a13*b11 + a23*b12
	result.M8 = a00*b20 + a10*b21 + a20*b22
	result.M9 = a01*b20 + a11*b21 + a21*b22
	result.M10 = a02*b20 + a12*b21 + a22*b22
	result.M11 = a03*b20 + a13*b21 + a23*b22
	result.M12 = mat.M12
	result.M13 = mat.M13
	result.M14 = mat.M14
	result.M15 = mat.M15

	return result
}

// MatrixRotateX - Returns x-rotation matrix (angle in radians)
func MatrixRotateX(angle float32) Matrix {
	result := MatrixIdentity()

	cosres := float32(math.Cos(float64(angle)))
	sinres := float32(math.Sin(float64(angle)))

	result.M5 = cosres
	result.M6 = -sinres
	result.M9 = sinres
	result.M10 = cosres

	return result
}

// MatrixRotateY - Returns y-rotation matrix (angle in radians)
func MatrixRotateY(angle float32) Matrix {
	result := MatrixIdentity()

	cosres := float32(math.Cos(float64(angle)))
	sinres := float32(math.Sin(float64(angle)))

	result.M0 = cosres
	result.M2 = sinres
	result.M8 = -sinres
	result.M10 = cosres

	return result
}

// MatrixRotateZ - Returns z-rotation matrix (angle in radians)
func MatrixRotateZ(angle float32) Matrix {
	result := MatrixIdentity()

	cosres := float32(math.Cos(float64(angle)))
	sinres := float32(math.Sin(float64(angle)))

	result.M0 = cosres
	result.M1 = -sinres
	result.M4 = sinres
	result.M5 = cosres

	return result
}

// MatrixScale - Returns scaling matrix
func MatrixScale(x, y, z float32) Matrix {
	result := NewMatrix(
		x, 0.0, 0.0, 0.0,
		0.0, y, 0.0, 0.0,
		0.0, 0.0, z, 0.0,
		0.0, 0.0, 0.0, 1.0)

	return result
}

// MatrixMultiply - Returns two matrix multiplication
func MatrixMultiply(left, right Matrix) Matrix {
	var result Matrix

	result.M0 = left.M0*right.M0 + left.M1*right.M4 + left.M2*right.M8 + left.M3*right.M12
	result.M1 = left.M0*right.M1 + left.M1*right.M5 + left.M2*right.M9 + left.M3*right.M13
	result.M2 = left.M0*right.M2 + left.M1*right.M6 + left.M2*right.M10 + left.M3*right.M14
	result.M3 = left.M0*right.M3 + left.M1*right.M7 + left.M2*right.M11 + left.M3*right.M15
	result.M4 = left.M4*right.M0 + left.M5*right.M4 + left.M6*right.M8 + left.M7*right.M12
	result.M5 = left.M4*right.M1 + left.M5*right.M5 + left.M6*right.M9 + left.M7*right.M13
	result.M6 = left.M4*right.M2 + left.M5*right.M6 + left.M6*right.M10 + left.M7*right.M14
	result.M7 = left.M4*right.M3 + left.M5*right.M7 + left.M6*right.M11 + left.M7*right.M15
	result.M8 = left.M8*right.M0 + left.M9*right.M4 + left.M10*right.M8 + left.M11*right.M12
	result.M9 = left.M8*right.M1 + left.M9*right.M5 + left.M10*right.M9 + left.M11*right.M13
	result.M10 = left.M8*right.M2 + left.M9*right.M6 + left.M10*right.M10 + left.M11*right.M14
	result.M11 = left.M8*right.M3 + left.M9*right.M7 + left.M10*right.M11 + left.M11*right.M15
	result.M12 = left.M12*right.M0 + left.M13*right.M4 + left.M14*right.M8 + left.M15*right.M12
	result.M13 = left.M12*right.M1 + left.M13*right.M5 + left.M14*right.M9 + left.M15*right.M13
	result.M14 = left.M12*right.M2 + left.M13*right.M6 + left.M14*right.M10 + left.M15*right.M14
	result.M15 = left.M12*right.M3 + left.M13*right.M7 + left.M14*right.M11 + left.M15*right.M15

	return result
}

// MatrixFrustum - Returns perspective projection matrix
func MatrixFrustum(left, right, bottom, top, near, far float32) Matrix {
	var result Matrix

	rl := right - left
	tb := top - bottom
	fn := far - near

	result.M0 = (near * 2.0) / rl
	result.M1 = 0.0
	result.M2 = 0.0
	result.M3 = 0.0

	result.M4 = 0.0
	result.M5 = (near * 2.0) / tb
	result.M6 = 0.0
	result.M7 = 0.0

	result.M8 = right + left/rl
	result.M9 = top + bottom/tb
	result.M10 = -(far + near) / fn
	result.M11 = -1.0

	result.M12 = 0.0
	result.M13 = 0.0
	result.M14 = -(far * near * 2.0) / fn
	result.M15 = 0.0

	return result
}

// MatrixPerspective - Returns perspective projection matrix
func MatrixPerspective(fovy, aspect, near, far float32) Matrix {
	top := near * float32(math.Tan(float64(fovy*Pi)/360.0))
	right := top * aspect

	return MatrixFrustum(-right, right, -top, top, near, far)
}

// MatrixOrtho - Returns orthographic projection matrix
func MatrixOrtho(left, right, bottom, top, near, far float32) Matrix {
	var result Matrix

	rl := right - left
	tb := top - bottom
	fn := far - near

	result.M0 = 2.0 / rl
	result.M1 = 0.0
	result.M2 = 0.0
	result.M3 = 0.0
	result.M4 = 0.0
	result.M5 = 2.0 / tb
	result.M6 = 0.0
	result.M7 = 0.0
	result.M8 = 0.0
	result.M9 = 0.0
	result.M10 = -2.0 / fn
	result.M11 = 0.0
	result.M12 = -(left + right) / rl
	result.M13 = -(top + bottom) / tb
	result.M14 = -(far + near) / fn
	result.M15 = 1.0

	return result
}

// MatrixLookAt - Returns camera look-at matrix (view matrix)
func MatrixLookAt(eye, target, up Vector3) Matrix {
	var result Matrix

	z := Vector3Subtract(eye, target)
	z = Vector3Normalize(z)
	x := Vector3CrossProduct(up, z)
	x = Vector3Normalize(x)
	y := Vector3CrossProduct(z, x)
	y = Vector3Normalize(y)

	result.M0 = x.X
	result.M1 = x.Y
	result.M2 = x.Z
	result.M3 = -((x.X * eye.X) + (x.Y * eye.Y) + (x.Z * eye.Z))
	result.M4 = y.X
	result.M5 = y.Y
	result.M6 = y.Z
	result.M7 = -((y.X * eye.X) + (y.Y * eye.Y) + (y.Z * eye.Z))
	result.M8 = z.X
	result.M9 = z.Y
	result.M10 = z.Z
	result.M11 = -((z.X * eye.X) + (z.Y * eye.Y) + (z.Z * eye.Z))
	result.M12 = 0.0
	result.M13 = 0.0
	result.M14 = 0.0
	result.M15 = 1.0

	return result
}

// QuaternionLength - Compute the length of a quaternion
func QuaternionLength(quat Quaternion) float32 {
	return float32(math.Sqrt(float64(quat.X*quat.X + quat.Y*quat.Y + quat.Z*quat.Z + quat.W*quat.W)))
}

// QuaternionNormalize - Normalize provided quaternion
func QuaternionNormalize(q Quaternion) Quaternion {
	var result Quaternion

	var length, ilength float32

	length = QuaternionLength(q)

	if length == 0.0 {
		length = 1.0
	}

	ilength = 1.0 / length

	result.X *= ilength
	result.Y *= ilength
	result.Z *= ilength
	result.W *= ilength

	return result
}

// QuaternionInvert - Invert provided quaternion
func QuaternionInvert(quat Quaternion) Quaternion {
	result := quat

	length := QuaternionLength(quat)
	lengthSq := length * length

	if lengthSq != 0.0 {
		i := 1.0 / lengthSq

		result.X *= -i
		result.Y *= -i
		result.Z *= -i
		result.W *= i
	}

	return result
}

// QuaternionMultiply - Calculate two quaternion multiplication
func QuaternionMultiply(q1, q2 Quaternion) Quaternion {
	var result Quaternion

	qax := q1.X
	qay := q1.Y
	qaz := q1.Z
	qaw := q1.W
	qbx := q2.X
	qby := q2.Y
	qbz := q2.Z
	qbw := q2.W

	result.X = qax*qbw + qaw*qbx + qay*qbz - qaz*qby
	result.Y = qay*qbw + qaw*qby + qaz*qbx - qax*qbz
	result.Z = qaz*qbw + qaw*qbz + qax*qby - qay*qbx
	result.W = qaw*qbw - qax*qbx - qay*qby - qaz*qbz

	return result
}

// QuaternionSlerp - Calculates spherical linear interpolation between two quaternions
func QuaternionSlerp(q1, q2 Quaternion, amount float32) Quaternion {
	var result Quaternion

	cosHalfTheta := q1.X*q2.X + q1.Y*q2.Y + q1.Z*q2.Z + q1.W*q2.W

	if math.Abs(float64(cosHalfTheta)) >= 1.0 {
		result = q1
	} else {
		halfTheta := float32(math.Acos(float64(cosHalfTheta)))
		sinHalfTheta := float32(math.Sqrt(float64(1.0 - cosHalfTheta*cosHalfTheta)))

		if math.Abs(float64(sinHalfTheta)) < 0.001 {
			result.X = q1.X*0.5 + q2.X*0.5
			result.Y = q1.Y*0.5 + q2.Y*0.5
			result.Z = q1.Z*0.5 + q2.Z*0.5
			result.W = q1.W*0.5 + q2.W*0.5
		} else {
			ratioA := float32(math.Sin(float64((1-amount)*halfTheta))) / sinHalfTheta
			ratioB := float32(math.Sin(float64(amount*halfTheta))) / sinHalfTheta

			result.X = q1.X*ratioA + q2.X*ratioB
			result.Y = q1.Y*ratioA + q2.Y*ratioB
			result.Z = q1.Z*ratioA + q2.Z*ratioB
			result.W = q1.W*ratioA + q2.W*ratioB
		}
	}

	return result
}

// QuaternionFromMatrix - Returns a quaternion for a given rotation matrix
func QuaternionFromMatrix(matrix Matrix) Quaternion {
	var result Quaternion

	trace := MatrixTrace(matrix)

	if trace > 0.0 {
		s := float32(math.Sqrt(float64(trace+1)) * 2.0)
		invS := 1.0 / s

		result.W = s * 0.25
		result.X = (matrix.M6 - matrix.M9) * invS
		result.Y = (matrix.M8 - matrix.M2) * invS
		result.Z = (matrix.M1 - matrix.M4) * invS
	} else {
		m00 := matrix.M0
		m11 := matrix.M5
		m22 := matrix.M10

		if m00 > m11 && m00 > m22 {
			s := float32(math.Sqrt(float64(1.0+m00-m11-m22)) * 2.0)
			invS := 1.0 / s

			result.W = (matrix.M6 - matrix.M9) * invS
			result.X = s * 0.25
			result.Y = (matrix.M4 + matrix.M1) * invS
			result.Z = (matrix.M8 + matrix.M2) * invS
		} else if m11 > m22 {
			s := float32(math.Sqrt(float64(1.0+m11-m00-m22)) * 2.0)
			invS := 1.0 / s

			result.W = (matrix.M8 - matrix.M2) * invS
			result.X = (matrix.M4 + matrix.M1) * invS
			result.Y = s * 0.25
			result.Z = (matrix.M9 + matrix.M6) * invS
		} else {
			s := float32(math.Sqrt(float64(1.0+m22-m00-m11)) * 2.0)
			invS := 1.0 / s

			result.W = (matrix.M1 - matrix.M4) * invS
			result.X = (matrix.M8 + matrix.M2) * invS
			result.Y = (matrix.M9 + matrix.M6) * invS
			result.Z = s * 0.25
		}
	}

	return result
}

// QuaternionToMatrix - Returns a matrix for a given quaternion
func QuaternionToMatrix(q Quaternion) Matrix {
	var result Matrix

	x := q.X
	y := q.Y
	z := q.Z
	w := q.W

	x2 := x + x
	y2 := y + y
	z2 := z + z

	xx := x * x2
	xy := x * y2
	xz := x * z2

	yy := y * y2
	yz := y * z2
	zz := z * z2

	wx := w * x2
	wy := w * y2
	wz := w * z2

	result.M0 = 1.0 - (yy + zz)
	result.M1 = xy - wz
	result.M2 = xz + wy
	result.M3 = 0.0
	result.M4 = xy + wz
	result.M5 = 1.0 - (xx + zz)
	result.M6 = yz - wx
	result.M7 = 0.0
	result.M8 = xz - wy
	result.M9 = yz + wx
	result.M10 = 1.0 - (xx + yy)
	result.M11 = 0.0
	result.M12 = 0.0
	result.M13 = 0.0
	result.M14 = 0.0
	result.M15 = 1.0

	return result
}

// QuaternionFromAxisAngle - Returns rotation quaternion for an angle and axis
func QuaternionFromAxisAngle(axis Vector3, angle float32) Quaternion {
	result := NewQuaternion(0.0, 0.0, 0.0, 1.0)

	if Vector3Length(axis) != 0.0 {
		angle *= 0.5
	}

	axis = Vector3Normalize(axis)

	sinres := float32(math.Sin(float64(angle)))
	cosres := float32(math.Cos(float64(angle)))

	result.X = axis.X * sinres
	result.Y = axis.Y * sinres
	result.Z = axis.Z * sinres
	result.W = cosres

	result = QuaternionNormalize(result)

	return result
}

// QuaternionToAxisAngle - Returns the rotation angle and axis for a given quaternion
func QuaternionToAxisAngle(q Quaternion, outAxis *Vector3, outAngle *float32) {
	if math.Abs(float64(q.W)) > 1.0 {
		q = QuaternionNormalize(q)
	}

	resAxis := NewVector3(0.0, 0.0, 0.0)

	resAngle := 2.0 * float32(math.Acos(float64(q.W)))
	den := float32(math.Sqrt(float64(1.0 - q.W*q.W)))

	if den > 0.0001 {
		resAxis.X = q.X / den
		resAxis.Y = q.Y / den
		resAxis.Z = q.Z / den
	} else {
		// This occurs when the angle is zero.
		// Not a problem: just set an arbitrary normalized axis.
		resAxis.X = 1.0
	}

	*outAxis = resAxis
	*outAngle = resAngle
}

// QuaternionTransform - Transform a quaternion given a transformation matrix
func QuaternionTransform(q Quaternion, mat Matrix) Quaternion {
	var result Quaternion

	x := q.X
	y := q.Y
	z := q.Z
	w := q.W

	result.X = mat.M0*x + mat.M4*y + mat.M8*z + mat.M12*w
	result.Y = mat.M1*x + mat.M5*y + mat.M9*z + mat.M13*w
	result.Z = mat.M2*x + mat.M6*y + mat.M10*z + mat.M14*w
	result.W = mat.M3*x + mat.M7*y + mat.M11*z + mat.M15*w

	return result
}

// Clamp - Clamp float value
func Clamp(value, min, max float32) float32 {
	var res float32
	if value < min {
		res = min
	} else {
		res = value
	}

	if res > max {
		return max
	}

	return res
}
