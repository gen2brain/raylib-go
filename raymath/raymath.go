// Package raymath - Some useful functions to work with Vector3, Matrix and Quaternions
package raymath

import (
	"math"

	"github.com/gen2brain/raylib-go/raylib"
)

// VectorAdd - Add two vectors
func VectorAdd(v1, v2 raylib.Vector3) raylib.Vector3 {
	result := raylib.Vector3{}

	result.X = v1.X + v2.X
	result.Y = v1.Y + v2.Y
	result.Z = v1.Z + v2.Z

	return result
}

// VectorSubtract - Subtract two vectors
func VectorSubtract(v1, v2 raylib.Vector3) raylib.Vector3 {
	result := raylib.Vector3{}

	result.X = v1.X - v2.X
	result.Y = v1.Y - v2.Y
	result.Z = v1.Z - v2.Z

	return result
}

// VectorCrossProduct - Calculate two vectors cross product
func VectorCrossProduct(v1, v2 raylib.Vector3) raylib.Vector3 {
	result := raylib.Vector3{}

	result.X = v1.Y*v2.Z - v1.Z*v2.Y
	result.Y = v1.Z*v2.X - v1.X*v2.Z
	result.Z = v1.X*v2.Y - v1.Y*v2.X

	return result
}

// VectorPerpendicular - Calculate one vector perpendicular vector
func VectorPerpendicular(v raylib.Vector3) raylib.Vector3 {
	result := raylib.Vector3{}

	min := math.Abs(float64(v.X))
	cardinalAxis := raylib.NewVector3(1.0, 0.0, 0.0)

	if math.Abs(float64(v.Y)) < min {
		min = math.Abs(float64(v.Y))
		cardinalAxis = raylib.NewVector3(0.0, 1.0, 0.0)
	}

	if math.Abs(float64(v.Z)) < min {
		cardinalAxis = raylib.NewVector3(0.0, 0.0, 1.0)
	}

	result = VectorCrossProduct(v, cardinalAxis)

	return result
}

// VectorDotProduct - Calculate two vectors dot product
func VectorDotProduct(v1, v2 raylib.Vector3) float32 {
	return v1.X*v2.X + v1.Y*v2.Y + v1.Z*v2.Z
}

// VectorLength - Calculate vector length
func VectorLength(v raylib.Vector3) float32 {
	return float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y + v.Z*v.Z)))
}

// VectorScale - Scale provided vector
func VectorScale(v *raylib.Vector3, scale float32) {
	v.X *= scale
	v.Y *= scale
	v.Z *= scale
}

// VectorNegate - Negate provided vector (invert direction)
func VectorNegate(v *raylib.Vector3) {
	v.X = -v.X
	v.Y = -v.Y
	v.Z = -v.Z
}

// VectorNormalize - Normalize provided vector
func VectorNormalize(v *raylib.Vector3) {
	var length, ilength float32

	length = VectorLength(*v)

	if length == 0 {
		length = 1.0
	}

	ilength = 1.0 / length

	v.X *= ilength
	v.Y *= ilength
	v.Z *= ilength
}

// VectorDistance - Calculate distance between two points
func VectorDistance(v1, v2 raylib.Vector3) float32 {
	var result float32

	dx := v2.X - v1.X
	dy := v2.Y - v1.Y
	dz := v2.Z - v1.Z

	result = float32(math.Sqrt(float64(dx*dx + dy*dy + dz*dz)))

	return result
}

// VectorLerp - Calculate linear interpolation between two vectors
func VectorLerp(v1, v2 raylib.Vector3, amount float32) raylib.Vector3 {
	result := raylib.Vector3{}

	result.X = v1.X + amount*(v2.X-v1.X)
	result.Y = v1.Y + amount*(v2.Y-v1.Y)
	result.Z = v1.Z + amount*(v2.Z-v1.Z)

	return result
}

// VectorReflect - Calculate reflected vector to normal
func VectorReflect(vector, normal raylib.Vector3) raylib.Vector3 {
	// I is the original vector
	// N is the normal of the incident plane
	// R = I - (2*N*( DotProduct[ I,N] ))

	result := raylib.Vector3{}

	dotProduct := VectorDotProduct(vector, normal)

	result.X = vector.X - (2.0*normal.X)*dotProduct
	result.Y = vector.Y - (2.0*normal.Y)*dotProduct
	result.Z = vector.Z - (2.0*normal.Z)*dotProduct

	return result
}

// VectorTransform - Transforms a Vector3 by a given Matrix
func VectorTransform(v *raylib.Vector3, mat raylib.Matrix) {
	x := v.X
	y := v.Y
	z := v.Z

	v.X = mat.M0*x + mat.M4*y + mat.M8*z + mat.M12
	v.Y = mat.M1*x + mat.M5*y + mat.M9*z + mat.M13
	v.Z = mat.M2*x + mat.M6*y + mat.M10*z + mat.M14
}

// VectorZero - Return a Vector3 init to zero
func VectorZero() raylib.Vector3 {
	return raylib.NewVector3(0.0, 0.0, 0.0)
}

// VectorMin - Return min value for each pair of components
func VectorMin(vec1, vec2 raylib.Vector3) raylib.Vector3 {
	result := raylib.Vector3{}

	result.X = float32(math.Min(float64(vec1.X), float64(vec2.X)))
	result.Y = float32(math.Min(float64(vec1.Y), float64(vec2.Y)))
	result.Z = float32(math.Min(float64(vec1.Z), float64(vec2.Z)))

	return result
}

// VectorMax - Return max value for each pair of components
func VectorMax(vec1, vec2 raylib.Vector3) raylib.Vector3 {
	result := raylib.Vector3{}

	result.X = float32(math.Max(float64(vec1.X), float64(vec2.X)))
	result.Y = float32(math.Max(float64(vec1.Y), float64(vec2.Y)))
	result.Z = float32(math.Max(float64(vec1.Z), float64(vec2.Z)))

	return result
}

// MatrixDeterminant - Compute matrix determinant
func MatrixDeterminant(mat raylib.Matrix) float32 {
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
func MatrixTrace(mat raylib.Matrix) float32 {
	return mat.M0 + mat.M5 + mat.M10 + mat.M15
}

// MatrixTranspose - Transposes provided matrix
func MatrixTranspose(mat *raylib.Matrix) {
	var temp raylib.Matrix

	temp.M0 = mat.M0
	temp.M1 = mat.M4
	temp.M2 = mat.M8
	temp.M3 = mat.M12
	temp.M4 = mat.M1
	temp.M5 = mat.M5
	temp.M6 = mat.M9
	temp.M7 = mat.M13
	temp.M8 = mat.M2
	temp.M9 = mat.M6
	temp.M10 = mat.M10
	temp.M11 = mat.M14
	temp.M12 = mat.M3
	temp.M13 = mat.M7
	temp.M14 = mat.M11
	temp.M15 = mat.M15

	mat = &temp
}

// MatrixInvert - Invert provided matrix
func MatrixInvert(mat *raylib.Matrix) {
	var temp raylib.Matrix

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

	temp.M0 = (a11*b11 - a12*b10 + a13*b09) * invDet
	temp.M1 = (-a01*b11 + a02*b10 - a03*b09) * invDet
	temp.M2 = (a31*b05 - a32*b04 + a33*b03) * invDet
	temp.M3 = (-a21*b05 + a22*b04 - a23*b03) * invDet
	temp.M4 = (-a10*b11 + a12*b08 - a13*b07) * invDet
	temp.M5 = (a00*b11 - a02*b08 + a03*b07) * invDet
	temp.M6 = (-a30*b05 + a32*b02 - a33*b01) * invDet
	temp.M7 = (a20*b05 - a22*b02 + a23*b01) * invDet
	temp.M8 = (a10*b10 - a11*b08 + a13*b06) * invDet
	temp.M9 = (-a00*b10 + a01*b08 - a03*b06) * invDet
	temp.M10 = (a30*b04 - a31*b02 + a33*b00) * invDet
	temp.M11 = (-a20*b04 + a21*b02 - a23*b00) * invDet
	temp.M12 = (-a10*b09 + a11*b07 - a12*b06) * invDet
	temp.M13 = (a00*b09 - a01*b07 + a02*b06) * invDet
	temp.M14 = (-a30*b03 + a31*b01 - a32*b00) * invDet
	temp.M15 = (a20*b03 - a21*b01 + a22*b00) * invDet

	mat = &temp
}

// MatrixNormalize - Normalize provided matrix
func MatrixNormalize(mat *raylib.Matrix) {
	det := MatrixDeterminant(*mat)

	mat.M0 /= det
	mat.M1 /= det
	mat.M2 /= det
	mat.M3 /= det
	mat.M4 /= det
	mat.M5 /= det
	mat.M6 /= det
	mat.M7 /= det
	mat.M8 /= det
	mat.M9 /= det
	mat.M10 /= det
	mat.M11 /= det
	mat.M12 /= det
	mat.M13 /= det
	mat.M14 /= det
	mat.M15 /= det
}

// MatrixIdentity - Returns identity matrix
func MatrixIdentity() raylib.Matrix {
	return raylib.NewMatrix(
		1.0, 0.0, 0.0, 0.0,
		0.0, 1.0, 0.0, 0.0,
		0.0, 0.0, 1.0, 0.0,
		0.0, 0.0, 0.0, 1.0)
}

// MatrixAdd - Add two matrices
func MatrixAdd(left, right raylib.Matrix) raylib.Matrix {
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
func MatrixSubtract(left, right raylib.Matrix) raylib.Matrix {
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
func MatrixTranslate(x, y, z float32) raylib.Matrix {
	return raylib.NewMatrix(
		1.0, 0.0, 0.0, 0.0,
		0.0, 1.0, 0.0, 0.0,
		0.0, 0.0, 1.0, 0.0,
		x, y, z, 1.0)
}

// MatrixRotate - Returns rotation matrix for an angle around an specified axis (angle in radians)
func MatrixRotate(axis raylib.Vector3, angle float32) raylib.Matrix {
	var result raylib.Matrix

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
func MatrixRotateX(angle float32) raylib.Matrix {
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
func MatrixRotateY(angle float32) raylib.Matrix {
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
func MatrixRotateZ(angle float32) raylib.Matrix {
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
func MatrixScale(x, y, z float32) raylib.Matrix {
	result := raylib.NewMatrix(
		x, 0.0, 0.0, 0.0,
		0.0, y, 0.0, 0.0,
		0.0, 0.0, z, 0.0,
		0.0, 0.0, 0.0, 1.0)

	return result
}

// MatrixMultiply - Returns two matrix multiplication
func MatrixMultiply(left, right raylib.Matrix) raylib.Matrix {
	var result raylib.Matrix

	result.M0 = right.M0*left.M0 + right.M1*left.M4 + right.M2*left.M8 + right.M3*left.M12
	result.M1 = right.M0*left.M1 + right.M1*left.M5 + right.M2*left.M9 + right.M3*left.M13
	result.M2 = right.M0*left.M2 + right.M1*left.M6 + right.M2*left.M10 + right.M3*left.M14
	result.M3 = right.M0*left.M3 + right.M1*left.M7 + right.M2*left.M11 + right.M3*left.M15
	result.M4 = right.M4*left.M0 + right.M5*left.M4 + right.M6*left.M8 + right.M7*left.M12
	result.M5 = right.M4*left.M1 + right.M5*left.M5 + right.M6*left.M9 + right.M7*left.M13
	result.M6 = right.M4*left.M2 + right.M5*left.M6 + right.M6*left.M10 + right.M7*left.M14
	result.M7 = right.M4*left.M3 + right.M5*left.M7 + right.M6*left.M11 + right.M7*left.M15
	result.M8 = right.M8*left.M0 + right.M9*left.M4 + right.M10*left.M8 + right.M11*left.M12
	result.M9 = right.M8*left.M1 + right.M9*left.M5 + right.M10*left.M9 + right.M11*left.M13
	result.M10 = right.M8*left.M2 + right.M9*left.M6 + right.M10*left.M10 + right.M11*left.M14
	result.M11 = right.M8*left.M3 + right.M9*left.M7 + right.M10*left.M11 + right.M11*left.M15
	result.M12 = right.M12*left.M0 + right.M13*left.M4 + right.M14*left.M8 + right.M15*left.M12
	result.M13 = right.M12*left.M1 + right.M13*left.M5 + right.M14*left.M9 + right.M15*left.M13
	result.M14 = right.M12*left.M2 + right.M13*left.M6 + right.M14*left.M10 + right.M15*left.M14
	result.M15 = right.M12*left.M3 + right.M13*left.M7 + right.M14*left.M11 + right.M15*left.M15

	return result
}

// MatrixFrustum - Returns perspective projection matrix
func MatrixFrustum(left, right, bottom, top, near, far float32) raylib.Matrix {
	var result raylib.Matrix

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
func MatrixPerspective(fovy, aspect, near, far float32) raylib.Matrix {
	top := near * float32(math.Tan(float64(fovy*raylib.Pi)/360.0))
	right := top * aspect

	return MatrixFrustum(-right, right, -top, top, near, far)
}

// MatrixOrtho - Returns orthographic projection matrix
func MatrixOrtho(left, right, bottom, top, near, far float32) raylib.Matrix {
	var result raylib.Matrix

	rl := (right - left)
	tb := (top - bottom)
	fn := (far - near)

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
func MatrixLookAt(eye, target, up raylib.Vector3) raylib.Matrix {
	var result raylib.Matrix

	z := VectorSubtract(eye, target)
	VectorNormalize(&z)
	x := VectorCrossProduct(up, z)
	VectorNormalize(&x)
	y := VectorCrossProduct(z, x)
	VectorNormalize(&y)

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
func QuaternionLength(quat raylib.Quaternion) float32 {
	return float32(math.Sqrt(float64(quat.X*quat.X + quat.Y*quat.Y + quat.Z*quat.Z + quat.W*quat.W)))
}

// QuaternionNormalize - Normalize provided quaternion
func QuaternionNormalize(q *raylib.Quaternion) {
	var length, ilength float32

	length = QuaternionLength(*q)

	if length == 0.0 {
		length = 1.0
	}

	ilength = 1.0 / length

	q.X *= ilength
	q.Y *= ilength
	q.Z *= ilength
	q.W *= ilength
}

// QuaternionInvert - Invert provided quaternion
func QuaternionInvert(quat *raylib.Quaternion) {
	length := QuaternionLength(*quat)
	lengthSq := length * length

	if lengthSq != 0.0 {
		i := 1.0 / lengthSq

		quat.X *= -i
		quat.Y *= -i
		quat.Z *= -i
		quat.W *= i
	}
}

// QuaternionMultiply - Calculate two quaternion multiplication
func QuaternionMultiply(q1, q2 raylib.Quaternion) raylib.Quaternion {
	var result raylib.Quaternion

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
func QuaternionSlerp(q1, q2 raylib.Quaternion, amount float32) raylib.Quaternion {
	var result raylib.Quaternion

	cosHalfTheta := q1.X*q2.X + q1.Y*q2.Y + q1.Z*q2.Z + q1.W*q2.W

	if math.Abs(float64(cosHalfTheta)) >= 1.0 {
		result = q1
	} else {
		halfTheta := float32(math.Acos(float64(cosHalfTheta)))
		sinHalfTheta := float32(math.Sqrt(float64(1.0 - cosHalfTheta*cosHalfTheta)))

		if math.Abs(float64(sinHalfTheta)) < 0.001 {
			result.X = (q1.X*0.5 + q2.X*0.5)
			result.Y = (q1.Y*0.5 + q2.Y*0.5)
			result.Z = (q1.Z*0.5 + q2.Z*0.5)
			result.W = (q1.W*0.5 + q2.W*0.5)
		} else {
			ratioA := float32(math.Sin(float64((1-amount)*halfTheta))) / sinHalfTheta
			ratioB := float32(math.Sin(float64(amount*halfTheta))) / sinHalfTheta

			result.X = (q1.X*ratioA + q2.X*ratioB)
			result.Y = (q1.Y*ratioA + q2.Y*ratioB)
			result.Z = (q1.Z*ratioA + q2.Z*ratioB)
			result.W = (q1.W*ratioA + q2.W*ratioB)
		}
	}

	return result
}

// QuaternionFromMatrix - Returns a quaternion for a given rotation matrix
func QuaternionFromMatrix(matrix raylib.Matrix) raylib.Quaternion {
	var result raylib.Quaternion

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
func QuaternionToMatrix(q raylib.Quaternion) raylib.Matrix {
	var result raylib.Matrix

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
func QuaternionFromAxisAngle(axis raylib.Vector3, angle float32) raylib.Quaternion {
	result := raylib.NewQuaternion(0.0, 0.0, 0.0, 1.0)

	if VectorLength(axis) != 0.0 {
		angle *= 0.5
	}

	VectorNormalize(&axis)

	sinres := float32(math.Sin(float64(angle)))
	cosres := float32(math.Cos(float64(angle)))

	result.X = axis.X * sinres
	result.Y = axis.Y * sinres
	result.Z = axis.Z * sinres
	result.W = cosres

	QuaternionNormalize(&result)

	return result
}

// QuaternionToAxisAngle - Returns the rotation angle and axis for a given quaternion
func QuaternionToAxisAngle(q raylib.Quaternion, outAxis *raylib.Vector3, outAngle *float32) {
	if math.Abs(float64(q.W)) > 1.0 {
		QuaternionNormalize(&q)
	}

	resAxis := raylib.NewVector3(0.0, 0.0, 0.0)

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
func QuaternionTransform(q *raylib.Quaternion, mat raylib.Matrix) {
	x := q.X
	y := q.Y
	z := q.Z
	w := q.W

	q.X = mat.M0*x + mat.M4*y + mat.M8*z + mat.M12*w
	q.Y = mat.M1*x + mat.M5*y + mat.M9*z + mat.M13*w
	q.Z = mat.M2*x + mat.M6*y + mat.M10*z + mat.M14*w
	q.W = mat.M3*x + mat.M7*y + mat.M11*z + mat.M15*w
}
