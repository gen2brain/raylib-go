//go:build (!cgo || nocgo) && !js
// +build !cgo nocgo
// +build !js

package rl

import "github.com/jupiterrider/ffi"

var typeTexture2D = ffi.NewType(&ffi.TypeUint32, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32)
var typeImage = ffi.NewType(&ffi.TypePointer, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32)
var typeVector2 = ffi.NewType(&ffi.TypeFloat, &ffi.TypeFloat)
var typeVector3 = ffi.NewType(&ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeFloat)
var typeVector4 = ffi.NewType(&ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeFloat)
var typeColor = ffi.NewType(&ffi.TypeUint8, &ffi.TypeUint8, &ffi.TypeUint8, &ffi.TypeUint8)
var typeCamera2D = ffi.NewType(&typeVector2, &typeVector2, &ffi.TypeFloat, &ffi.TypeFloat)
var typeCamera3D = ffi.NewType(&typeVector3, &typeVector3, &typeVector3, &ffi.TypeFloat, &ffi.TypeSint32)
var typeRenderTexture2D = ffi.NewType(&ffi.TypeUint32, &typeTexture2D, &typeTexture2D)
var typeShader = ffi.NewType(&ffi.TypeUint32, &ffi.TypePointer)
var typeMatrix = ffi.NewType(
	&ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeFloat,
	&ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeFloat,
	&ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeFloat,
	&ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeFloat,
)
var typeVrStereoConfig = ffi.NewType(
	&typeMatrix, &typeMatrix,
	&typeMatrix, &typeMatrix,
	&ffi.TypeFloat, &ffi.TypeFloat,
	&ffi.TypeFloat, &ffi.TypeFloat,
	&ffi.TypeFloat, &ffi.TypeFloat,
	&ffi.TypeFloat, &ffi.TypeFloat,
	&ffi.TypeFloat, &ffi.TypeFloat,
	&ffi.TypeFloat, &ffi.TypeFloat,
)
var typeVrDeviceInfo = ffi.NewType(
	&ffi.TypeSint32, &ffi.TypeSint32,
	&ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeFloat,
	&ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeFloat,
	&ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeFloat,
)
var typeRay = ffi.NewType(&typeVector3, &typeVector3)
var typeFilePathList = ffi.NewType(&ffi.TypeUint32, &ffi.TypePointer)
var typeAutomationEvent = ffi.NewType(
	&ffi.TypeUint32, &ffi.TypeUint32,
	&ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32,
)
var typeAutomationEventList = ffi.NewType(&ffi.TypeUint32, &ffi.TypeUint32, &ffi.TypePointer)
var typeRectangle = ffi.NewType(&ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeFloat)
var typeFont = ffi.NewType(
	&ffi.TypeSint32,
	&ffi.TypeSint32,
	&ffi.TypeSint32,
	&typeTexture2D,
	&ffi.TypePointer,
	&ffi.TypePointer,
)
var typeNPatchInfo = ffi.NewType(
	&typeRectangle,
	&ffi.TypeSint32,
	&ffi.TypeSint32,
	&ffi.TypeSint32,
	&ffi.TypeSint32,
	&ffi.TypeSint32,
)
var typeGlyphInfo = ffi.NewType(
	&ffi.TypeSint32,
	&ffi.TypeSint32,
	&ffi.TypeSint32,
	&ffi.TypeSint32,
	&typeImage,
)
var typeModelSkeleton = ffi.NewType(&ffi.TypeSint32, &ffi.TypePointer, &ffi.TypePointer)
var typeModel = ffi.NewType(
	&typeMatrix,
	&ffi.TypeSint32,
	&ffi.TypeSint32,
	&ffi.TypePointer,
	&ffi.TypePointer,
	&ffi.TypePointer,
	&typeModelSkeleton,
	&ffi.TypePointer,
	&ffi.TypePointer,
)
var typeBoundingBox = ffi.NewType(&typeVector3, &typeVector3)
var typeMesh = ffi.NewType(
	&ffi.TypeSint32, &ffi.TypeSint32,
	&ffi.TypePointer, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypePointer,
	&ffi.TypeSint32, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypeUint32, &ffi.TypePointer,
)
var typeMaterial = ffi.NewType(
	&typeShader, &ffi.TypePointer,
	&ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeFloat,
)
var typeModelAnimation = ffi.NewType(
	&ffi.TypeSint8, &ffi.TypeSint8, &ffi.TypeSint8, &ffi.TypeSint8, &ffi.TypeSint8, &ffi.TypeSint8, &ffi.TypeSint8, &ffi.TypeSint8,
	&ffi.TypeSint8, &ffi.TypeSint8, &ffi.TypeSint8, &ffi.TypeSint8, &ffi.TypeSint8, &ffi.TypeSint8, &ffi.TypeSint8, &ffi.TypeSint8,
	&ffi.TypeSint8, &ffi.TypeSint8, &ffi.TypeSint8, &ffi.TypeSint8, &ffi.TypeSint8, &ffi.TypeSint8, &ffi.TypeSint8, &ffi.TypeSint8,
	&ffi.TypeSint8, &ffi.TypeSint8, &ffi.TypeSint8, &ffi.TypeSint8, &ffi.TypeSint8, &ffi.TypeSint8, &ffi.TypeSint8, &ffi.TypeSint8,
	&ffi.TypeSint32,
	&ffi.TypeSint32,
	&ffi.TypePointer,
)
var typeRayCollision = ffi.NewType(&ffi.TypeUint8, &ffi.TypeFloat, &typeVector3, &typeVector3)
var typeWave = ffi.NewType(
	&ffi.TypeUint32, &ffi.TypeUint32, &ffi.TypeUint32, &ffi.TypeUint32,
	&ffi.TypePointer,
)
var typeAudioStream = ffi.NewType(
	&ffi.TypePointer, &ffi.TypePointer,
	&ffi.TypeUint32, &ffi.TypeUint32, &ffi.TypeUint32,
)
var typeSound = ffi.NewType(&typeAudioStream, &ffi.TypeUint32)
var typeMusic = ffi.NewType(&typeAudioStream, &ffi.TypeUint32, &ffi.TypeUint8, &ffi.TypeSint32, &ffi.TypePointer)
var typeRenderBatch = ffi.NewType(&ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypeSint32, &ffi.TypeFloat)
