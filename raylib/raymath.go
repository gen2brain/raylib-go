package rl

import (
	"math"
)

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

// Lerp - Calculate linear interpolation between two floats
func Lerp(start, end, amount float32) float32 {
	return start + amount*(end-start)
}

// Normalize - Normalize input value within input range
func Normalize(value, start, end float32) float32 {
	return (value - start) / (end - start)
}

// Remap - Remap input value within input range to output range
func Remap(value, inputStart, inputEnd, outputStart, outputEnd float32) float32 {
	return (value-inputStart)/(inputEnd-inputStart)*(outputEnd-outputStart) + outputStart
}

// Wrap - Wrap input value from min to max
func Wrap(value, min, max float32) float32 {
	return float32(float64(value) - float64(max-min)*math.Floor(float64((value-min)/(max-min))))
}

// FloatEquals - Check whether two given floats are almost equal
func FloatEquals(x, y float32) bool {
	return (math.Abs(float64(x-y)) <= 0.000001*math.Max(1.0, math.Max(math.Abs(float64(x)), math.Abs(float64(y)))))
}

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

// Vector2AddValue - Add vector and float value
func Vector2AddValue(v Vector2, add float32) Vector2 {
	return NewVector2(v.X+add, v.Y+add)
}

// Vector2Subtract - Subtract two vectors (v1 - v2)
func Vector2Subtract(v1, v2 Vector2) Vector2 {
	return NewVector2(v1.X-v2.X, v1.Y-v2.Y)
}

// Vector2SubtractValue - Subtract vector by float value
func Vector2SubtractValue(v Vector2, sub float32) Vector2 {
	return NewVector2(v.X-sub, v.Y-sub)
}

// Vector2Length - Calculate vector length
func Vector2Length(v Vector2) float32 {
	return float32(math.Sqrt(float64((v.X * v.X) + (v.Y * v.Y))))
}

// Vector2LengthSqr - Calculate vector square length
func Vector2LengthSqr(v Vector2) float32 {
	return v.X*v.X + v.Y*v.Y
}

// Vector2DotProduct - Calculate two vectors dot product
func Vector2DotProduct(v1, v2 Vector2) float32 {
	return v1.X*v2.X + v1.Y*v2.Y
}

// Vector2Distance - Calculate distance between two vectors
func Vector2Distance(v1, v2 Vector2) float32 {
	return float32(math.Sqrt(float64((v1.X-v2.X)*(v1.X-v2.X) + (v1.Y-v2.Y)*(v1.Y-v2.Y))))
}

// Vector2DistanceSqr - Calculate square distance between two vectors
func Vector2DistanceSqr(v1 Vector2, v2 Vector2) float32 {
	return (v1.X-v2.X)*(v1.X-v2.X) + (v1.Y-v2.Y)*(v1.Y-v2.Y)
}

// Vector2Angle - Calculate angle from two vectors in radians
func Vector2Angle(v1, v2 Vector2) float32 {
	result := math.Atan2(float64(v2.Y), float64(v2.X)) - math.Atan2(float64(v1.Y), float64(v1.X))

	return float32(result)
}

// Vector2LineAngle - Calculate angle defined by a two vectors line
// NOTE: Parameters need to be normalized. Current implementation should be aligned with glm::angle
func Vector2LineAngle(start Vector2, end Vector2) float32 {
	return float32(-math.Atan2(float64(end.Y-start.Y), float64(end.X-start.X)))
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
func Vector2Divide(v1, v2 Vector2) Vector2 {
	return NewVector2(v1.X/v2.X, v1.Y/v2.Y)
}

// Vector2Normalize - Normalize provided vector
func Vector2Normalize(v Vector2) Vector2 {
	if l := Vector2Length(v); l > 0 {
		return Vector2Scale(v, 1/l)
	}
	return v
}

// Vector2Transform - Transforms a Vector2 by a given Matrix
func Vector2Transform(v Vector2, mat Matrix) Vector2 {
	var result = Vector2{}

	var x = v.X
	var y = v.Y
	var z float32

	result.X = mat.M0*x + mat.M4*y + mat.M8*z + mat.M12
	result.Y = mat.M1*x + mat.M5*y + mat.M9*z + mat.M13

	return result
}

// Vector2Lerp - Calculate linear interpolation between two vectors
func Vector2Lerp(v1, v2 Vector2, amount float32) Vector2 {
	return NewVector2(v1.X+amount*(v2.X-v1.X), v1.Y+amount*(v2.Y-v1.Y))
}

// Vector2Reflect - Calculate reflected vector to normal
func Vector2Reflect(v Vector2, normal Vector2) Vector2 {
	var result = Vector2{}

	dotProduct := v.X*normal.X + v.Y*normal.Y // Dot product

	result.X = v.X - 2.0*normal.X*dotProduct
	result.Y = v.Y - 2.0*normal.Y*dotProduct

	return result
}

// Vector2Rotate - Rotate vector by angle
func Vector2Rotate(v Vector2, angle float32) Vector2 {
	var result = Vector2{}

	cosres := float32(math.Cos(float64(angle)))
	sinres := float32(math.Sin(float64(angle)))

	result.X = v.X*cosres - v.Y*sinres
	result.Y = v.X*sinres + v.Y*cosres

	return result
}

// Vector2MoveTowards - Move Vector towards target
func Vector2MoveTowards(v Vector2, target Vector2, maxDistance float32) Vector2 {
	var result = Vector2{}

	dx := target.X - v.X
	dy := target.Y - v.Y
	value := dx*dx + dy*dy

	if value == 0 || maxDistance >= 0 && value <= maxDistance*maxDistance {
		return target
	}

	dist := float32(math.Sqrt(float64(value)))

	result.X = v.X + dx/dist*maxDistance
	result.Y = v.Y + dy/dist*maxDistance

	return result
}

// Vector2Invert - Invert the given vector
func Vector2Invert(v Vector2) Vector2 {
	return NewVector2(1.0/v.X, 1.0/v.Y)
}

// Vector2Clamp - Clamp the components of the vector between min and max values specified by the given vectors
func Vector2Clamp(v Vector2, min Vector2, max Vector2) Vector2 {
	var result = Vector2{}

	result.X = float32(math.Min(float64(max.X), math.Max(float64(min.X), float64(v.X))))
	result.Y = float32(math.Min(float64(max.Y), math.Max(float64(min.Y), float64(v.Y))))

	return result
}

// Vector2ClampValue - Clamp the magnitude of the vector between two min and max values
func Vector2ClampValue(v Vector2, min float32, max float32) Vector2 {
	var result = v

	length := v.X*v.X + v.Y*v.Y
	if length > 0.0 {
		length = float32(math.Sqrt(float64(length)))

		if length < min {
			scale := min / length
			result.X = v.X * scale
			result.Y = v.Y * scale
		} else if length > max {
			scale := max / length
			result.X = v.X * scale
			result.Y = v.Y * scale
		}
	}

	return result
}

// Vector2Equals - Check whether two given vectors are almost equal
func Vector2Equals(p Vector2, q Vector2) bool {
	return (math.Abs(float64(p.X-q.X)) <= 0.000001*math.Max(1.0, math.Max(math.Abs(float64(p.X)), math.Abs(float64(q.X)))) &&
		math.Abs(float64(p.Y-q.Y)) <= 0.000001*math.Max(1.0, math.Max(math.Abs(float64(p.Y)), math.Abs(float64(q.Y)))))
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

// Vector3AddValue - Add vector and float value
func Vector3AddValue(v Vector3, add float32) Vector3 {
	return NewVector3(v.X+add, v.Y+add, v.Z+add)
}

// Vector3Subtract - Subtract two vectors
func Vector3Subtract(v1, v2 Vector3) Vector3 {
	return NewVector3(v1.X-v2.X, v1.Y-v2.Y, v1.Z-v2.Z)
}

// Vector3SubtractValue - Subtract vector by float value
func Vector3SubtractValue(v Vector3, sub float32) Vector3 {
	return NewVector3(v.X-sub, v.Y-sub, v.Z-sub)
}

// Vector3Scale - Scale provided vector
func Vector3Scale(v Vector3, scale float32) Vector3 {
	return NewVector3(v.X*scale, v.Y*scale, v.Z*scale)
}

// Vector3Multiply - Multiply vector by vector
func Vector3Multiply(v1, v2 Vector3) Vector3 {
	result := Vector3{}

	result.X = v1.X * v2.X
	result.Y = v1.Y * v2.Y
	result.Z = v1.Z * v2.Z

	return result
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
	min := math.Abs(float64(v.X))
	cardinalAxis := NewVector3(1.0, 0.0, 0.0)

	if math.Abs(float64(v.Y)) < min {
		min = math.Abs(float64(v.Y))
		cardinalAxis = NewVector3(0.0, 1.0, 0.0)
	}

	if math.Abs(float64(v.Z)) < min {
		cardinalAxis = NewVector3(0.0, 0.0, 1.0)
	}

	result := Vector3CrossProduct(v, cardinalAxis)

	return result
}

// Vector3Length - Calculate vector length
func Vector3Length(v Vector3) float32 {
	return float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y + v.Z*v.Z)))
}

// Vector3LengthSqr - Calculate vector square length
func Vector3LengthSqr(v Vector3) float32 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
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

// Vector3DistanceSqr - Calculate square distance between two vectors
func Vector3DistanceSqr(v1 Vector3, v2 Vector3) float32 {
	var result float32

	dx := v2.X - v1.X
	dy := v2.Y - v1.Y
	dz := v2.Z - v1.Z
	result = dx*dx + dy*dy + dz*dz

	return result
}

// Vector3Angle - Calculate angle between two vectors
func Vector3Angle(v1 Vector3, v2 Vector3) float32 {
	var result float32

	cross := Vector3{X: v1.Y*v2.Z - v1.Z*v2.Y, Y: v1.Z*v2.X - v1.X*v2.Z, Z: v1.X*v2.Y - v1.Y*v2.X}
	length := float32(math.Sqrt(float64(cross.X*cross.X + cross.Y*cross.Y + cross.Z*cross.Z)))
	dot := v1.X*v2.X + v1.Y*v2.Y + v1.Z*v2.Z
	result = float32(math.Atan2(float64(length), float64(dot)))

	return result
}

// Vector3Negate - Negate provided vector (invert direction)
func Vector3Negate(v Vector3) Vector3 {
	return NewVector3(-v.X, -v.Y, -v.Z)
}

// Vector3Divide - Divide vector by vector
func Vector3Divide(v1 Vector3, v2 Vector3) Vector3 {
	return NewVector3(v1.X/v2.X, v1.Y/v2.Y, v1.Z/v2.Z)
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

// Vector3Project - Calculate the projection of the vector v1 on to v2
func Vector3Project(v1, v2 Vector3) Vector3 {
	result := Vector3{}

	v1dv2 := (v1.X*v2.X + v1.Y*v2.Y + v1.Z*v2.Z)
	v2dv2 := (v2.X*v2.X + v2.Y*v2.Y + v2.Z*v2.Z)

	mag := v1dv2 / v2dv2

	result.X = v2.X * mag
	result.Y = v2.Y * mag
	result.Z = v2.Z * mag

	return result
}

// Vector3Reject - Calculate the rejection of the vector v1 on to v2
func Vector3Reject(v1, v2 Vector3) Vector3 {
	result := Vector3{}

	v1dv2 := (v1.X*v2.X + v1.Y*v2.Y + v1.Z*v2.Z)
	v2dv2 := (v2.X*v2.X + v2.Y*v2.Y + v2.Z*v2.Z)

	mag := v1dv2 / v2dv2

	result.X = v1.X - (v2.X * mag)
	result.Y = v1.Y - (v2.Y * mag)
	result.Z = v1.Z - (v2.Z * mag)

	return result
}

// Vector3OrthoNormalize - Orthonormalize provided vectors
// Makes vectors normalized and orthogonal to each other
// Gram-Schmidt function implementation
func Vector3OrthoNormalize(v1, v2 *Vector3) {
	Vector3Normalize(*v1)

	vn1 := Vector3CrossProduct(*v1, *v2)
	Vector3Normalize(vn1)

	vn2 := Vector3CrossProduct(vn1, *v1)
	*v2 = vn2
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

// Vector3RotateByQuaternion - Transform a vector by quaternion rotation
func Vector3RotateByQuaternion(v Vector3, q Quaternion) Vector3 {
	var result Vector3

	result.X = v.X*(q.X*q.X+q.W*q.W-q.Y*q.Y-q.Z*q.Z) + v.Y*(2*q.X*q.Y-2*q.W*q.Z) + v.Z*(2*q.X*q.Z+2*q.W*q.Y)
	result.Y = v.X*(2*q.W*q.Z+2*q.X*q.Y) + v.Y*(q.W*q.W-q.X*q.X+q.Y*q.Y-q.Z*q.Z) + v.Z*(-2*q.W*q.X+2*q.Y*q.Z)
	result.Z = v.X*(-2*q.W*q.Y+2*q.X*q.Z) + v.Y*(2*q.W*q.X+2*q.Y*q.Z) + v.Z*(q.W*q.W-q.X*q.X-q.Y*q.Y+q.Z*q.Z)

	return result
}

// Vector3RotateByAxisAngle - Rotates a vector around an axis
func Vector3RotateByAxisAngle(v Vector3, axis Vector3, angle float32) Vector3 {
	// Using Euler-Rodrigues Formula
	// Ref.: https://en.wikipedia.org/w/index.php?title=Euler%E2%80%93Rodrigues_formula

	result := v

	// Vector3Normalize(axis);
	length := float32(math.Sqrt(float64(axis.X*axis.X + axis.Y*axis.Y + axis.Z*axis.Z)))
	if length == 0.0 {
		length = 1.0
	}
	ilength := 1.0 / length
	axis.X *= ilength
	axis.Y *= ilength
	axis.Z *= ilength

	angle /= 2.0
	a := float32(math.Sin(float64(angle)))
	b := axis.X * a
	c := axis.Y * a
	d := axis.Z * a
	a = float32(math.Cos(float64(angle)))
	w := NewVector3(b, c, d)

	// Vector3CrossProduct(w, v)
	wv := NewVector3(w.Y*v.Z-w.Z*v.Y, w.Z*v.X-w.X*v.Z, w.X*v.Y-w.Y*v.X)

	// Vector3CrossProduct(w, wv)
	wwv := NewVector3(w.Y*wv.Z-w.Z*wv.Y, w.Z*wv.X-w.X*wv.Z, w.X*wv.Y-w.Y*wv.X)

	// Vector3Scale(wv, 2*a)
	a *= 2
	wv.X *= a
	wv.Y *= a
	wv.Z *= a

	// Vector3Scale(wwv, 2)
	wwv.X *= 2
	wwv.Y *= 2
	wwv.Z *= 2

	result.X += wv.X
	result.Y += wv.Y
	result.Z += wv.Z

	result.X += wwv.X
	result.Y += wwv.Y
	result.Z += wwv.Z

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

// Vector3Unproject - Projects a Vector3 from screen space into object space
// NOTE: We are avoiding calling other raymath functions despite available
func Vector3Unproject(source Vector3, projection Matrix, view Matrix) Vector3 {
	var result = Vector3{}

	// Calculate unprojected matrix (multiply view matrix by projection matrix) and invert it
	var matViewProj = Matrix{ // MatrixMultiply(view, projection);
		M0:  view.M0*projection.M0 + view.M1*projection.M4 + view.M2*projection.M8 + view.M3*projection.M12,
		M4:  view.M0*projection.M1 + view.M1*projection.M5 + view.M2*projection.M9 + view.M3*projection.M13,
		M8:  view.M0*projection.M2 + view.M1*projection.M6 + view.M2*projection.M10 + view.M3*projection.M14,
		M12: view.M0*projection.M3 + view.M1*projection.M7 + view.M2*projection.M11 + view.M3*projection.M15,
		M1:  view.M4*projection.M0 + view.M5*projection.M4 + view.M6*projection.M8 + view.M7*projection.M12,
		M5:  view.M4*projection.M1 + view.M5*projection.M5 + view.M6*projection.M9 + view.M7*projection.M13,
		M9:  view.M4*projection.M2 + view.M5*projection.M6 + view.M6*projection.M10 + view.M7*projection.M14,
		M13: view.M4*projection.M3 + view.M5*projection.M7 + view.M6*projection.M11 + view.M7*projection.M15,
		M2:  view.M8*projection.M0 + view.M9*projection.M4 + view.M10*projection.M8 + view.M11*projection.M12,
		M6:  view.M8*projection.M1 + view.M9*projection.M5 + view.M10*projection.M9 + view.M11*projection.M13,
		M10: view.M8*projection.M2 + view.M9*projection.M6 + view.M10*projection.M10 + view.M11*projection.M14,
		M14: view.M8*projection.M3 + view.M9*projection.M7 + view.M10*projection.M11 + view.M11*projection.M15,
		M3:  view.M12*projection.M0 + view.M13*projection.M4 + view.M14*projection.M8 + view.M15*projection.M12,
		M7:  view.M12*projection.M1 + view.M13*projection.M5 + view.M14*projection.M9 + view.M15*projection.M13,
		M11: view.M12*projection.M2 + view.M13*projection.M6 + view.M14*projection.M10 + view.M15*projection.M14,
		M15: view.M12*projection.M3 + view.M13*projection.M7 + view.M14*projection.M11 + view.M15*projection.M15}

	// Calculate inverted matrix -> MatrixInvert(matViewProj);
	// Cache the matrix values (speed optimization)
	var a00 = matViewProj.M0
	var a01 = matViewProj.M1
	var a02 = matViewProj.M2
	var a03 = matViewProj.M3
	var a10 = matViewProj.M4
	var a11 = matViewProj.M5
	var a12 = matViewProj.M6
	var a13 = matViewProj.M7
	var a20 = matViewProj.M8
	var a21 = matViewProj.M9
	var a22 = matViewProj.M10
	var a23 = matViewProj.M11
	var a30 = matViewProj.M12
	var a31 = matViewProj.M13
	var a32 = matViewProj.M14
	var a33 = matViewProj.M15

	var b00 = a00*a11 - a01*a10
	var b01 = a00*a12 - a02*a10
	var b02 = a00*a13 - a03*a10
	var b03 = a01*a12 - a02*a11
	var b04 = a01*a13 - a03*a11
	var b05 = a02*a13 - a03*a12
	var b06 = a20*a31 - a21*a30
	var b07 = a20*a32 - a22*a30
	var b08 = a20*a33 - a23*a30
	var b09 = a21*a32 - a22*a31
	var b10 = a21*a33 - a23*a31
	var b11 = a22*a33 - a23*a32

	// Calculate the invert determinant (inlined to avoid double-caching)
	var invDet = 1.0 / (b00*b11 - b01*b10 + b02*b09 + b03*b08 - b04*b07 + b05*b06)

	var matViewProjInv = Matrix{
		M0:  (a11*b11 - a12*b10 + a13*b09) * invDet,
		M4:  (-a01*b11 + a02*b10 - a03*b09) * invDet,
		M8:  (a31*b05 - a32*b04 + a33*b03) * invDet,
		M12: (-a21*b05 + a22*b04 - a23*b03) * invDet,
		M1:  (-a10*b11 + a12*b08 - a13*b07) * invDet,
		M5:  (a00*b11 - a02*b08 + a03*b07) * invDet,
		M9:  (-a30*b05 + a32*b02 - a33*b01) * invDet,
		M13: (a20*b05 - a22*b02 + a23*b01) * invDet,
		M2:  (a10*b10 - a11*b08 + a13*b06) * invDet,
		M6:  (-a00*b10 + a01*b08 - a03*b06) * invDet,
		M10: (a30*b04 - a31*b02 + a33*b00) * invDet,
		M14: (-a20*b04 + a21*b02 - a23*b00) * invDet,
		M3:  (-a10*b09 + a11*b07 - a12*b06) * invDet,
		M7:  (a00*b09 - a01*b07 + a02*b06) * invDet,
		M11: (-a30*b03 + a31*b01 - a32*b00) * invDet,
		M15: (a20*b03 - a21*b01 + a22*b00) * invDet}

	// Create quaternion from source point
	var quat = Quaternion{X: source.X, Y: source.Y, Z: source.Z, W: 1.0}

	// Multiply quat point by unprojecte matrix
	var qtransformed = Quaternion{ // QuaternionTransform(quat, matViewProjInv)
		X: matViewProjInv.M0*quat.X + matViewProjInv.M4*quat.Y + matViewProjInv.M8*quat.Z + matViewProjInv.M12*quat.W,
		Y: matViewProjInv.M1*quat.X + matViewProjInv.M5*quat.Y + matViewProjInv.M9*quat.Z + matViewProjInv.M13*quat.W,
		Z: matViewProjInv.M2*quat.X + matViewProjInv.M6*quat.Y + matViewProjInv.M10*quat.Z + matViewProjInv.M14*quat.W,
		W: matViewProjInv.M3*quat.X + matViewProjInv.M7*quat.Y + matViewProjInv.M11*quat.Z + matViewProjInv.M15*quat.W}

	// Normalized world points in vectors
	result.X = qtransformed.X / qtransformed.W
	result.Y = qtransformed.Y / qtransformed.W
	result.Z = qtransformed.Z / qtransformed.W

	return result
}

// Vector3ToFloatV - Get Vector3 as float array
func Vector3ToFloatV(v Vector3) [3]float32 {
	var result [3]float32

	result[0] = v.X
	result[1] = v.Y
	result[2] = v.Z

	return result
}

// Vector3Invert - Invert the given vector
func Vector3Invert(v Vector3) Vector3 {
	return NewVector3(1.0/v.X, 1.0/v.Y, 1.0/v.Z)
}

// Vector3Clamp - Clamp the components of the vector between min and max values specified by the given vectors
func Vector3Clamp(v Vector3, min Vector3, max Vector3) Vector3 {
	var result = Vector3{}

	result.X = float32(math.Min(float64(max.X), math.Max(float64(min.X), float64(v.X))))
	result.Y = float32(math.Min(float64(max.Y), math.Max(float64(min.Y), float64(v.Y))))
	result.Z = float32(math.Min(float64(max.Z), math.Max(float64(min.Z), float64(v.Z))))

	return result
}

// Vector3ClampValue - Clamp the magnitude of the vector between two values
func Vector3ClampValue(v Vector3, min float32, max float32) Vector3 {
	var result = v

	length := v.X*v.X + v.Y*v.Y + v.Z*v.Z
	if length > 0.0 {
		length = float32(math.Sqrt(float64(length)))

		if length < min {
			scale := min / length
			result.X = v.X * scale
			result.Y = v.Y * scale
			result.Z = v.Z * scale
		} else if length > max {
			scale := max / length
			result.X = v.X * scale
			result.Y = v.Y * scale
			result.Z = v.Z * scale
		}
	}

	return result
}

// Vector3Equals - Check whether two given vectors are almost equal
func Vector3Equals(p Vector3, q Vector3) bool {
	return (math.Abs(float64(p.X-q.X)) <= 0.000001*math.Max(1.0, math.Max(math.Abs(float64(p.X)), math.Abs(float64(q.X)))) &&
		math.Abs(float64(p.Y-q.Y)) <= 0.000001*math.Max(1.0, math.Max(math.Abs(float64(p.Y)), math.Abs(float64(q.Y)))) &&
		math.Abs(float64(p.Z-q.Z)) <= 0.000001*math.Max(1.0, math.Max(math.Abs(float64(p.Z)), math.Abs(float64(q.Z)))))
}

// Vector3Refract - Compute the direction of a refracted ray
//
// v: normalized direction of the incoming ray
// n: normalized normal vector of the interface of two optical media
// r: ratio of the refractive index of the medium from where the ray comes to the refractive index of the medium on the other side of the surface
func Vector3Refract(v Vector3, n Vector3, r float32) Vector3 {
	var result = Vector3{}

	dot := v.X*n.X + v.Y*n.Y + v.Z*n.Z
	d := 1.0 - r*r*(1.0-dot*dot)

	if d >= 0.0 {
		d = float32(math.Sqrt(float64(d)))
		v.X = r*v.X - (r*dot+d)*n.X
		v.Y = r*v.Y - (r*dot+d)*n.Y
		v.Z = r*v.Z - (r*dot+d)*n.Z

		result = v
	}

	return result
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

// MatrixIdentity - Returns identity matrix
func MatrixIdentity() Matrix {
	return NewMatrix(
		1.0, 0.0, 0.0, 0.0,
		0.0, 1.0, 0.0, 0.0,
		0.0, 0.0, 1.0, 0.0,
		0.0, 0.0, 0.0, 1.0)
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

// MatrixRotateXYZ - Get xyz-rotation matrix (angles in radians)
func MatrixRotateXYZ(ang Vector3) Matrix {
	result := MatrixIdentity()

	cosz := float32(math.Cos(float64(-ang.Z)))
	sinz := float32(math.Sin(float64(-ang.Z)))
	cosy := float32(math.Cos(float64(-ang.Y)))
	siny := float32(math.Sin(float64(-ang.Y)))
	cosx := float32(math.Cos(float64(-ang.X)))
	sinx := float32(math.Sin(float64(-ang.X)))

	result.M0 = cosz * cosy
	result.M4 = (cosz * siny * sinx) - (sinz * cosx)
	result.M8 = (cosz * siny * cosx) + (sinz * sinx)

	result.M1 = sinz * cosy
	result.M5 = (sinz * siny * sinx) + (cosz * cosx)
	result.M9 = (sinz * siny * cosx) - (cosz * sinx)

	result.M2 = -siny
	result.M6 = cosy * sinx
	result.M10 = cosy * cosx

	return result
}

// MatrixRotateZYX - Get zyx-rotation matrix
// NOTE: Angle must be provided in radians
func MatrixRotateZYX(angle Vector3) Matrix {
	var result = Matrix{}

	var cz = float32(math.Cos(float64(angle.Z)))
	var sz = float32(math.Sin(float64(angle.Z)))
	var cy = float32(math.Cos(float64(angle.Y)))
	var sy = float32(math.Sin(float64(angle.Y)))
	var cx = float32(math.Cos(float64(angle.X)))
	var sx = float32(math.Sin(float64(angle.X)))

	result.M0 = cz * cy
	result.M4 = cz*sy*sx - cx*sz
	result.M8 = sz*sx + cz*cx*sy
	result.M12 = float32(0)

	result.M1 = cy * sz
	result.M5 = cz*cx + sz*sy*sx
	result.M9 = cx*sz*sy - cz*sx
	result.M13 = float32(0)

	result.M2 = -sy
	result.M6 = cy * sx
	result.M10 = cy * cx
	result.M14 = float32(0)

	result.M3 = float32(0)
	result.M7 = float32(0)
	result.M11 = float32(0)
	result.M15 = float32(1)

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

// MatrixToFloatV - Get float array of matrix data
func MatrixToFloatV(mat Matrix) [16]float32 {
	var result [16]float32

	result[0] = mat.M0
	result[1] = mat.M1
	result[2] = mat.M2
	result[3] = mat.M3
	result[4] = mat.M4
	result[5] = mat.M5
	result[6] = mat.M6
	result[7] = mat.M7
	result[8] = mat.M8
	result[9] = mat.M9
	result[10] = mat.M10
	result[11] = mat.M11
	result[12] = mat.M12
	result[13] = mat.M13
	result[14] = mat.M14
	result[15] = mat.M15

	return result
}

// MatrixToFloat - Converts Matrix to float32 slice
func MatrixToFloat(mat Matrix) []float32 {
	data := MatrixToFloatV(mat)
	return data[:]
}

// QuaternionAdd - Add two quaternions
func QuaternionAdd(q1 Quaternion, q2 Quaternion) Quaternion {
	var result = Quaternion{X: q1.X + q2.X, Y: q1.Y + q2.Y, Z: q1.Z + q2.Z, W: q1.W + q2.W}

	return result
}

// QuaternionAddValue - Add quaternion and float value
func QuaternionAddValue(q Quaternion, add float32) Quaternion {
	var result = Quaternion{X: q.X + add, Y: q.Y + add, Z: q.Z + add, W: q.W + add}

	return result
}

// QuaternionSubtract - Subtract two quaternions
func QuaternionSubtract(q1 Quaternion, q2 Quaternion) Quaternion {
	var result = Quaternion{X: q1.X - q2.X, Y: q1.Y - q2.Y, Z: q1.Z - q2.Z, W: q1.W - q2.W}

	return result
}

// QuaternionSubtractValue - Subtract quaternion and float value
func QuaternionSubtractValue(q Quaternion, sub float32) Quaternion {
	var result = Quaternion{X: q.X - sub, Y: q.Y - sub, Z: q.Z - sub, W: q.W - sub}

	return result
}

// QuaternionIdentity - Get identity quaternion
func QuaternionIdentity() Quaternion {
	var result = Quaternion{W: 1.0}

	return result
}

// QuaternionLength - Compute the length of a quaternion
func QuaternionLength(quat Quaternion) float32 {
	return float32(math.Sqrt(float64(quat.X*quat.X + quat.Y*quat.Y + quat.Z*quat.Z + quat.W*quat.W)))
}

// QuaternionNormalize - Normalize provided quaternion
func QuaternionNormalize(q Quaternion) Quaternion {
	result := q

	length := QuaternionLength(q)

	if length != 0.0 {
		result.X /= length
		result.Y /= length
		result.Z /= length
		result.W /= length
	}

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

// QuaternionScale - Scale quaternion by float value
func QuaternionScale(q Quaternion, mul float32) Quaternion {
	var result = Quaternion{}

	result.X = q.X * mul
	result.Y = q.Y * mul
	result.Z = q.Z * mul
	result.W = q.W * mul

	return result
}

// QuaternionDivide - Divide two quaternions
func QuaternionDivide(q1 Quaternion, q2 Quaternion) Quaternion {
	var result = Quaternion{X: q1.X / q2.X, Y: q1.Y / q2.Y, Z: q1.Z / q2.Z, W: q1.W / q2.W}

	return result
}

// QuaternionLerp - Calculate linear interpolation between two quaternions
func QuaternionLerp(q1 Quaternion, q2 Quaternion, amount float32) Quaternion {
	var result = Quaternion{}

	result.X = q1.X + amount*(q2.X-q1.X)
	result.Y = q1.Y + amount*(q2.Y-q1.Y)
	result.Z = q1.Z + amount*(q2.Z-q1.Z)
	result.W = q1.W + amount*(q2.W-q1.W)

	return result
}

// QuaternionNlerp - Calculate slerp-optimized interpolation between two quaternions
func QuaternionNlerp(q1 Quaternion, q2 Quaternion, amount float32) Quaternion {
	var result = Quaternion{}

	// QuaternionLerp(q1, q2, amount)
	result.X = q1.X + amount*(q2.X-q1.X)
	result.Y = q1.Y + amount*(q2.Y-q1.Y)
	result.Z = q1.Z + amount*(q2.Z-q1.Z)
	result.W = q1.W + amount*(q2.W-q1.W)

	// QuaternionNormalize(q);
	q := result
	length := float32(math.Sqrt(float64(q.X*q.X + q.Y*q.Y + q.Z*q.Z + q.W*q.W)))
	if length == 0.0 {
		length = 1.0
	}
	ilength := 1.0 / length

	result.X = q.X * ilength
	result.Y = q.Y * ilength
	result.Z = q.Z * ilength
	result.W = q.W * ilength

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

// QuaternionFromVector3ToVector3 - Calculate quaternion based on the rotation from one vector to another
func QuaternionFromVector3ToVector3(from Vector3, to Vector3) Quaternion {
	var result = Quaternion{}

	cos2Theta := from.X*to.X + from.Y*to.Y + from.Z*to.Z                                                       // Vector3DotProduct(from, to)
	cross := Vector3{X: from.Y*to.Z - from.Z*to.Y, Y: from.Z*to.X - from.X*to.Z, Z: from.X*to.Y - from.Y*to.X} // Vector3CrossProduct(from, to)

	result.X = cross.X
	result.Y = cross.Y
	result.Z = cross.Z
	result.W = 1.0 + cos2Theta

	// QuaternionNormalize(q);
	// NOTE: Normalize to essentially nlerp the original and identity to 0.5
	q := result
	length := float32(math.Sqrt(float64(q.X*q.X + q.Y*q.Y + q.Z*q.Z + q.W*q.W)))
	if length == 0.0 {
		length = 1.0
	}
	ilength := 1.0 / length

	result.X = q.X * ilength
	result.Y = q.Y * ilength
	result.Z = q.Z * ilength
	result.W = q.W * ilength

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

// QuaternionFromEuler - Get the quaternion equivalent to Euler angles
// NOTE: Rotation order is ZYX
func QuaternionFromEuler(pitch, yaw, roll float32) Quaternion {
	var result Quaternion

	x0 := float32(math.Cos(float64(pitch * 0.5)))
	x1 := float32(math.Sin(float64(pitch * 0.5)))
	y0 := float32(math.Cos(float64(yaw * 0.5)))
	y1 := float32(math.Sin(float64(yaw * 0.5)))
	z0 := float32(math.Cos(float64(roll * 0.5)))
	z1 := float32(math.Sin(float64(roll * 0.5)))

	result.X = x1*y0*z0 - x0*y1*z1
	result.Y = x0*y1*z0 + x1*y0*z1
	result.Z = x0*y0*z1 - x1*y1*z0
	result.W = x0*y0*z0 + x1*y1*z1

	return result
}

// QuaternionToEuler - Get the Euler angles equivalent to quaternion (roll, pitch, yaw)
// NOTE: Angles are returned in a Vector3 struct in radians
func QuaternionToEuler(q Quaternion) Vector3 {
	var result Vector3

	// Roll (x-axis rotation)
	x0 := 2.0 * (q.W*q.X + q.Y*q.Z)
	x1 := 1.0 - 2.0*(q.X*q.X+q.Y*q.Y)
	result.X = float32(math.Atan2(float64(x0), float64(x1)))

	// Pitch (y-axis rotation)
	y0 := 2.0 * (q.W*q.Y - q.Z*q.X)
	y0 = Clamp(y0, -1.0, 1.0)
	result.Y = float32(math.Asin(float64(y0)))

	// Yaw (z-axis rotation)
	z0 := 2.0 * (q.W*q.Z + q.X*q.Y)
	z1 := 1.0 - 2.0*(q.Y*q.Y+q.Z*q.Z)
	result.Z = float32(math.Atan2(float64(z0), float64(z1)))

	return result
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

// QuaternionEquals - Check whether two given quaternions are almost equal
func QuaternionEquals(p, q Quaternion) bool {
	return (math.Abs(float64(p.X-q.X)) <= 0.000001*math.Max(1.0, math.Max(math.Abs(float64(p.X)), math.Abs(float64(q.X)))) &&
		math.Abs(float64(p.Y-q.Y)) <= 0.000001*math.Max(1.0, math.Max(math.Abs(float64(p.Y)), math.Abs(float64(q.Y)))) &&
		math.Abs(float64(p.Z-q.Z)) <= 0.000001*math.Max(1.0, math.Max(math.Abs(float64(p.Z)), math.Abs(float64(q.Z)))) &&
		math.Abs(float64(p.W-q.W)) <= 0.000001*math.Max(1.0, math.Max(math.Abs(float64(p.W)), math.Abs(float64(q.W)))) ||
		math.Abs(float64(p.X+q.X)) <= 0.000001*math.Max(1.0, math.Max(math.Abs(float64(p.X)), math.Abs(float64(q.X)))) &&
			math.Abs(float64(p.Y+q.Y)) <= 0.000001*math.Max(1.0, math.Max(math.Abs(float64(p.Y)), math.Abs(float64(q.Y)))) &&
			math.Abs(float64(p.Z+q.Z)) <= 0.000001*math.Max(1.0, math.Max(math.Abs(float64(p.Z)), math.Abs(float64(q.Z)))) &&
			math.Abs(float64(p.W+q.W)) <= 0.000001*math.Max(1.0, math.Max(math.Abs(float64(p.W)), math.Abs(float64(q.W)))))
}

// MatrixDecompose - Decompose a transformation matrix into its rotational, translational and scaling components
func MatrixDecompose(mat Matrix, translation *Vector3, rotation *Quaternion, scale *Vector3) {
	// Extract translation.
	translation.X = mat.M12
	translation.Y = mat.M13
	translation.Z = mat.M14

	// Extract upper-left for determinant computation
	a := mat.M0
	b := mat.M4
	c := mat.M8
	d := mat.M1
	e := mat.M5
	f := mat.M9
	g := mat.M2
	h := mat.M6
	i := mat.M10
	A := e*i - f*h
	B := f*g - d*i
	C := d*h - e*g

	// Extract scale
	det := a*A + b*B + c*C
	abc := NewVector3(a, b, c)
	def := NewVector3(d, e, f)
	ghi := NewVector3(g, h, i)

	scalex := Vector3Length(abc)
	scaley := Vector3Length(def)
	scalez := Vector3Length(ghi)
	s := NewVector3(scalex, scaley, scalez)

	if det < 0 {
		s = Vector3Negate(s)
	}

	*scale = s

	// Remove scale from the matrix if it is not close to zero
	clone := mat
	if !FloatEquals(det, 0) {
		clone.M0 /= s.X
		clone.M5 /= s.Y
		clone.M10 /= s.Z

		// Extract rotation
		*rotation = QuaternionFromMatrix(clone)
	} else {
		// Set to identity if close to zero
		*rotation = QuaternionIdentity()
	}
}
