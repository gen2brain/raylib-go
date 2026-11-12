//go:build cgo && !nocgo

package rl

import (
	"math"
	"testing"
	"time"
)

// tests in raymath_generated_test.go too

func BenchmarkVector3OrthoNormalize(b *testing.B) {
	var (
		perCCall  time.Duration
		perGoCall time.Duration
	)
	b.Run("c", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			v1 := NewVector3(1, 2, 3)
			v2 := NewVector3(1, 2, 3)
			cVector3OrthoNormalize(&v1, &v2)
		}
		perCCall = b.Elapsed() / time.Duration(b.N)
	})
	b.Run("go", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			v1 := NewVector3(1, 2, 3)
			v2 := NewVector3(1, 2, 3)
			Vector3OrthoNormalize(&v1, &v2)
		}
		perGoCall = b.Elapsed() / time.Duration(b.N)
	})
	if perCCall < perGoCall {
		b.Log("Go slower than C")
	}
}

func FuzzVector3OrthoNormalize(f *testing.F) {
	f.Add(
		float32(1), float32(2), float32(3),
		float32(4), float32(5), float32(6),
	)
	f.Fuzz(func(t *testing.T,
		v1X, v1Y, v1Z float32,
		v2X, v2Y, v2Z float32,
	) {
		want1 := NewVector3(v1X, v1Y, v1Z)
		want2 := NewVector3(v2X, v2Y, v2Z)
		inputScale := testScale(want1, want2)
		if testVector3Equals(Vector3CrossProduct(want1, want2), Vector3Zero(), 0) {
			t.SkipNow() // too unstable on macos
		}
		cVector3OrthoNormalize(&want1, &want2)

		got1 := NewVector3(v1X, v1Y, v1Z)
		got2 := NewVector3(v2X, v2Y, v2Z)
		Vector3OrthoNormalize(&got1, &got2)

		if !testVector3Equals(want1, got1, inputScale) {
			t.Errorf("v1: got %v; want %v", got1, want1)
		}
		if !testVector3Equals(want2, got2, inputScale) {
			t.Errorf("v2: got %v; want %v", got2, want2)
		}
	})
}

func BenchmarkQuaternionToAxisAngle(b *testing.B) {
	q := NewQuaternion(1, 2, 3, 4)
	var (
		perCCall  time.Duration
		perGoCall time.Duration
	)
	b.Run("c", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var axis Vector3
			var angle float32
			cQuaternionToAxisAngle(q, &axis, &angle)
		}
		perCCall = b.Elapsed() / time.Duration(b.N)
	})
	b.Run("go", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var outAxis Vector3
			var outAngle float32
			QuaternionToAxisAngle(q, &outAxis, &outAngle)
		}
		perGoCall = b.Elapsed() / time.Duration(b.N)
	})
	if perCCall < perGoCall {
		b.Log("Go slower than C")
	}
}

func FuzzQuaternionToAxisAngle(f *testing.F) {
	skip := NewVector3(1, 0, 0)
	f.Add(
		float32(1), float32(2), float32(3), float32(4),
	)
	f.Fuzz(func(t *testing.T,
		qX, qY, qZ, qW float32,
	) {
		q := NewQuaternion(qX, qY, qZ, qW)
		q = QuaternionNormalize(q)

		var wantAxis Vector3
		var wantAngle float32
		cQuaternionToAxisAngle(q, &wantAxis, &wantAngle)

		var gotAxis Vector3
		var gotAngle float32
		QuaternionToAxisAngle(q, &gotAxis, &gotAngle)

		axisScale := float32(1 / math.Max(math.Sqrt(math.Max(0, 1-float64(q.W)*float64(q.W))), epsilon))
		if !testVector3Equals(wantAxis, skip, 0) && // it's ok if our version has higher precision
			!testVector3Equals(wantAxis, gotAxis, axisScale) {
			t.Errorf("axis: got %v; want %v", gotAxis, wantAxis)
		}
		if !testFloat32Equals(wantAngle, gotAngle, axisScale) {
			t.Errorf("angle: got %v; want %v", gotAngle, wantAngle)
		}
	})
}

func BenchmarkMatrixDecompose(b *testing.B) {
	mat := NewMatrix(
		1, 2, 3, 4,
		1, 2, 3, 4,
		1, 2, 3, 4,
		1, 2, 3, 4,
	)
	var (
		perCCall  time.Duration
		perGoCall time.Duration
	)
	b.Run("c", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var translation Vector3
			var rotation Quaternion
			var scale Vector3
			cMatrixDecompose(mat, &translation, &rotation, &scale)
		}
		perCCall = b.Elapsed() / time.Duration(b.N)
	})
	b.Run("go", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var translation Vector3
			var rotation Quaternion
			var scale Vector3
			MatrixDecompose(mat, &translation, &rotation, &scale)
		}
		perGoCall = b.Elapsed() / time.Duration(b.N)
	})
	if perCCall < perGoCall {
		b.Log("Go slower than C")
	}
}

func FuzzMatrixDecompose(f *testing.F) {
	f.Add(
		float32(1), float32(2), float32(3), float32(4),
		float32(1), float32(2), float32(3), float32(4),
		float32(1), float32(2), float32(3), float32(4),
		float32(1), float32(2), float32(3), float32(4),
	)
	f.Add(
		float32(0), float32(-3), float32(0), float32(5),
		float32(2), float32(0), float32(0), float32(-6),
		float32(0), float32(0), float32(4), float32(7),
		float32(0), float32(0), float32(0), float32(1),
	)
	f.Fuzz(func(t *testing.T,
		m0, m4, m8, m12,
		m1, m5, m9, m13,
		m2, m6, m10, m14,
		m3, m7, m11, m15 float32,

	) {
		mat := NewMatrix(
			m0, m4, m8, m12,
			m1, m5, m9, m13,
			m2, m6, m10, m14,
			m3, m7, m11, m15,
		)

		var wantTranslation Vector3
		var wantRotation Quaternion
		var wantScale Vector3
		cMatrixDecompose(mat, &wantTranslation, &wantRotation, &wantScale)

		var gotTranslation Vector3
		var gotRotation Quaternion
		var gotScale Vector3
		MatrixDecompose(mat, &gotTranslation, &gotRotation, &gotScale)

		inputScale := testScale(mat)
		if !testVector3Equals(wantTranslation, gotTranslation, inputScale) {
			t.Errorf("translation: got %v; want %v", gotTranslation, wantTranslation)
		}
		sx, sy, sz := math.Abs(float64(gotScale.X)), math.Abs(float64(gotScale.Y)), math.Abs(float64(gotScale.Z))
		if !(math.Min(math.Min(sx, sy), sz) > 1e-3*math.Max(math.Max(sx, sy), sz)) {
			return // singular: rotation and scale are rounding noise
		}
		if !testQuaternionEquals(wantRotation, gotRotation, inputScale) {
			t.Errorf("rotation: got %v; want %v", gotRotation, wantRotation)
		}
		if !testVector3Equals(wantScale, gotScale, inputScale) {
			t.Errorf("scale: got %v; want %v", gotScale, wantScale)
		}
	})
}
