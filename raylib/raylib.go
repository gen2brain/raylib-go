/*
Package raylib - Go bindings for raylib, a simple and easy-to-use library to learn videogames programming.

raylib is highly inspired by Borland BGI graphics lib and by XNA framework.

raylib could be useful for prototyping, tools development, graphic applications, embedded systems and education.

NOTE for ADVENTURERS: raylib is a programming library to learn videogames programming; no fancy interface, no visual helpers, no auto-debugging... just coding in the most pure spartan-programmers way.

Example:

	package main

	import "github.com/icodealot/raylib-go-headless/raylib"

	func main() {
		rl.InitRaylib(800, 450)

		rl.SetTargetFPS(60)

		for !rl.RaylibShouldClose() {
			rl.BeginDrawing()

			rl.ClearBackground(rl.RayWhite)

			rl.DrawText("Congrats! You created your first window!", 190, 200, 20, rl.LightGray)

			rl.EndDrawing()
		}

		rl.CloseRaylib()
	}


*/
package rl

/*
#include "raylib.h"
#include <stdlib.h>
*/
import "C"
import (
	"image"
	"image/color"
	"io"
	"runtime"
	"unsafe"
)

func init() {
	// Make sure the main goroutine is bound to the main thread.
	runtime.LockOSThread()
}

// Wave type, defines audio wave data
type Wave struct {
	// Number of samples
	SampleCount uint32
	// Frequency (samples per second)
	SampleRate uint32
	// Bit depth (bits per sample): 8, 16, 32 (24 not supported)
	SampleSize uint32
	// Number of channels (1-mono, 2-stereo)
	Channels uint32
	// Buffer data pointer
	data unsafe.Pointer
}

// NewWave - Returns new Wave
func NewWave(sampleCount, sampleRate, sampleSize, channels uint32, data []byte) Wave {
	d := unsafe.Pointer(&data[0])
	return Wave{sampleCount, sampleRate, sampleSize, channels, d}
}

// newWaveFromPointer - Returns new Wave from pointer
func newWaveFromPointer(ptr unsafe.Pointer) Wave {
	return *(*Wave)(ptr)
}

// Sound source type
type Sound struct {
	Stream      AudioStream
	SampleCount uint32
	_           [4]byte
}

// newSoundFromPointer - Returns new Sound from pointer
func newSoundFromPointer(ptr unsafe.Pointer) Sound {
	return *(*Sound)(ptr)
}

// Music type (file streaming from memory)
// NOTE: Anything longer than ~10 seconds should be streamed
type Music struct {
	Stream      AudioStream
	SampleCount uint32
	Looping     bool
	CtxType     int32
	CtxData     unsafe.Pointer
}

// newMusicFromPointer - Returns new Music from pointer
func newMusicFromPointer(ptr unsafe.Pointer) Music {
	return *(*Music)(ptr)
}

// AudioStream type
// NOTE: Useful to create custom audio streams not bound to a specific file
type AudioStream struct {
	// Buffer
	Buffer *C.rAudioBuffer
	// Frequency (samples per second)
	SampleRate uint32
	// Bit depth (bits per sample): 8, 16, 32 (24 not supported)
	SampleSize uint32
	// Number of channels (1-mono, 2-stereo)
	Channels uint32
	_        [4]byte
}

// newAudioStreamFromPointer - Returns new AudioStream from pointer
func newAudioStreamFromPointer(ptr unsafe.Pointer) AudioStream {
	return *(*AudioStream)(ptr)
}

// CameraMode type
type CameraMode int32

// Camera system modes
const (
	CameraCustom CameraMode = iota
	CameraFree
	CameraOrbital
	CameraFirstPerson
	CameraThirdPerson
)

// CameraProjection type
type CameraProjection int32

// Camera projection modes
const (
	CameraPerspective CameraProjection = iota
	CameraOrthographic
)

// ShaderUniformDataType type
type ShaderUniformDataType int32

// ShaderUniformDataType enumeration
const (
	ShaderUniformFloat ShaderUniformDataType = iota
	ShaderUniformVec2
	ShaderUniformVec3
	ShaderUniformVec4
	ShaderUniformInt
	ShaderUniformIvec2
	ShaderUniformIvec3
	ShaderUniformIvec4
	ShaderUniformSampler2d
)

// Some basic Defines
const (
	Pi      = 3.1415927
	Deg2rad = 0.017453292
	Rad2deg = 57.295776

	// Raylib Config Flags

	// Set to try enabling V-Sync on GPU
	FlagVsyncHint = 0x00000040
	// Set to run program in fullscreen
	FlagFullscreenMode = 0x00000002
	// Set to allow resizable window
	FlagWindowResizable = 0x00000004
	// Set to disable window decoration (frame and buttons)
	FlagWindowUndecorated = 0x00000008
	// Set to hide window
	FlagWindowHidden = 0x00000080
	// Set to minimize window (iconify)
	FlagWindowMinimized = 0x00000200
	// Set to maximize window (expanded to monitor)
	FlagWindowMaximized = 0x00000400
	// Set to window non focused
	FlagWindowUnfocused = 0x00000800
	// Set to window always on top
	FlagWindowTopmost = 0x00001000
	// Set to allow windows running while minimized
	FlagWindowAlwaysRun = 0x00000100
	// Set to allow transparent window
	FlagWindowTransparent = 0x00000010
	// Set to support HighDPI
	FlagWindowHighdpi = 0x00002000
	// Set to try enabling MSAA 4X
	FlagMsaa4xHint = 0x00000020
	// Set to try enabling interlaced video format (for V3D)
	FlagInterlacedHint = 0x00010000

	// Keyboard Function Keys
	KeySpace        = 32
	KeyEscape       = 256
	KeyEnter        = 257
	KeyTab          = 258
	KeyBackspace    = 259
	KeyInsert       = 260
	KeyDelete       = 261
	KeyRight        = 262
	KeyLeft         = 263
	KeyDown         = 264
	KeyUp           = 265
	KeyPageUp       = 266
	KeyPageDown     = 267
	KeyHome         = 268
	KeyEnd          = 269
	KeyCapsLock     = 280
	KeyScrollLock   = 281
	KeyNumLock      = 282
	KeyPrintScreen  = 283
	KeyPause        = 284
	KeyF1           = 290
	KeyF2           = 291
	KeyF3           = 292
	KeyF4           = 293
	KeyF5           = 294
	KeyF6           = 295
	KeyF7           = 296
	KeyF8           = 297
	KeyF9           = 298
	KeyF10          = 299
	KeyF11          = 300
	KeyF12          = 301
	KeyLeftShift    = 340
	KeyLeftControl  = 341
	KeyLeftAlt      = 342
	KeyLeftSuper    = 343
	KeyRightShift   = 344
	KeyRightControl = 345
	KeyRightAlt     = 346
	KeyRightSuper   = 347
	KeyKbMenu       = 348
	KeyLeftBracket  = 91
	KeyBackSlash    = 92
	KeyRightBracket = 93
	KeyGrave        = 96

	// Keyboard Number Pad Keys
	KeyKp0        = 320
	KeyKp1        = 321
	KeyKp2        = 322
	KeyKp3        = 323
	KeyKp4        = 324
	KeyKp5        = 325
	KeyKp6        = 326
	KeyKp7        = 327
	KeyKp8        = 328
	KeyKp9        = 329
	KeyKpDecimal  = 330
	KeyKpDivide   = 331
	KeyKpMultiply = 332
	KeyKpSubtract = 333
	KeyKpAdd      = 334
	KeyKpEnter    = 335
	KeyKpEqual    = 336

	// Keyboard Alpha Numeric Keys
	KeyApostrophe = 39
	KeyComma      = 44
	KeyMinus      = 45
	KeyPeriod     = 46
	KeySlash      = 47
	KeyZero       = 48
	KeyOne        = 49
	KeyTwo        = 50
	KeyThree      = 51
	KeyFour       = 52
	KeyFive       = 53
	KeySix        = 54
	KeySeven      = 55
	KeyEight      = 56
	KeyNine       = 57
	KeySemicolon  = 59
	KeyEqual      = 61
	KeyA          = 65
	KeyB          = 66
	KeyC          = 67
	KeyD          = 68
	KeyE          = 69
	KeyF          = 70
	KeyG          = 71
	KeyH          = 72
	KeyI          = 73
	KeyJ          = 74
	KeyK          = 75
	KeyL          = 76
	KeyM          = 77
	KeyN          = 78
	KeyO          = 79
	KeyP          = 80
	KeyQ          = 81
	KeyR          = 82
	KeyS          = 83
	KeyT          = 84
	KeyU          = 85
	KeyV          = 86
	KeyW          = 87
	KeyX          = 88
	KeyY          = 89
	KeyZ          = 90

	// Android keys
	KeyBack       = 4
	KeyMenu       = 82
	KeyVolumeUp   = 24
	KeyVolumeDown = 25

	// Mouse Buttons
	MouseLeftButton   = 0
	MouseRightButton  = 1
	MouseMiddleButton = 2

	// Touch points registered
	MaxTouchPoints = 2

	// Gamepad Number
	GamepadPlayer1 = 0
	GamepadPlayer2 = 1
	GamepadPlayer3 = 2
	GamepadPlayer4 = 3

	// Gamepad Buttons/Axis

	// PS3 USB Controller Buttons
	GamepadPs3ButtonTriangle = 0
	GamepadPs3ButtonCircle   = 1
	GamepadPs3ButtonCross    = 2
	GamepadPs3ButtonSquare   = 3
	GamepadPs3ButtonL1       = 6
	GamepadPs3ButtonR1       = 7
	GamepadPs3ButtonL2       = 4
	GamepadPs3ButtonR2       = 5
	GamepadPs3ButtonStart    = 8
	GamepadPs3ButtonSelect   = 9
	GamepadPs3ButtonUp       = 24
	GamepadPs3ButtonRight    = 25
	GamepadPs3ButtonDown     = 26
	GamepadPs3ButtonLeft     = 27
	GamepadPs3ButtonPs       = 12

	// PS3 USB Controller Axis
	GamepadPs3AxisLeftX  = 0
	GamepadPs3AxisLeftY  = 1
	GamepadPs3AxisRightX = 2
	GamepadPs3AxisRightY = 5
	// [1..-1] (pressure-level)
	GamepadPs3AxisL2 = 3
	// [1..-1] (pressure-level)
	GamepadPs3AxisR2 = 4

	// Xbox360 USB Controller Buttons
	GamepadXboxButtonA      = 0
	GamepadXboxButtonB      = 1
	GamepadXboxButtonX      = 2
	GamepadXboxButtonY      = 3
	GamepadXboxButtonLb     = 4
	GamepadXboxButtonRb     = 5
	GamepadXboxButtonSelect = 6
	GamepadXboxButtonStart  = 7
	GamepadXboxButtonUp     = 10
	GamepadXboxButtonRight  = 11
	GamepadXboxButtonDown   = 12
	GamepadXboxButtonLeft   = 13
	GamepadXboxButtonHome   = 8

	// Android Gamepad Controller (SNES CLASSIC)
	GamepadAndroidDpadUp     = 19
	GamepadAndroidDpadDown   = 20
	GamepadAndroidDpadLeft   = 21
	GamepadAndroidDpadRight  = 22
	GamepadAndroidDpadCenter = 23

	GamepadAndroidButtonA  = 96
	GamepadAndroidButtonB  = 97
	GamepadAndroidButtonC  = 98
	GamepadAndroidButtonX  = 99
	GamepadAndroidButtonY  = 100
	GamepadAndroidButtonZ  = 101
	GamepadAndroidButtonL1 = 102
	GamepadAndroidButtonR1 = 103
	GamepadAndroidButtonL2 = 104
	GamepadAndroidButtonR2 = 105

	// Xbox360 USB Controller Axis
	// [-1..1] (left->right)
	GamepadXboxAxisLeftX = 0
	// [1..-1] (up->down)
	GamepadXboxAxisLeftY = 1
	// [-1..1] (left->right)
	GamepadXboxAxisRightX = 2
	// [1..-1] (up->down)
	GamepadXboxAxisRightY = 3
	// [-1..1] (pressure-level)
	GamepadXboxAxisLt = 4
	// [-1..1] (pressure-level)
	GamepadXboxAxisRt = 5
)

// Some Basic Colors
// NOTE: Custom raylib color palette for amazing visuals on WHITE background
var (
	// Light Gray
	LightGray = NewColor(200, 200, 200, 255)
	// Gray
	Gray = NewColor(130, 130, 130, 255)
	// Dark Gray
	DarkGray = NewColor(80, 80, 80, 255)
	// Yellow
	Yellow = NewColor(253, 249, 0, 255)
	// Gold
	Gold = NewColor(255, 203, 0, 255)
	// Orange
	Orange = NewColor(255, 161, 0, 255)
	// Pink
	Pink = NewColor(255, 109, 194, 255)
	// Red
	Red = NewColor(230, 41, 55, 255)
	// Maroon
	Maroon = NewColor(190, 33, 55, 255)
	// Green
	Green = NewColor(0, 228, 48, 255)
	// Lime
	Lime = NewColor(0, 158, 47, 255)
	// Dark Green
	DarkGreen = NewColor(0, 117, 44, 255)
	// Sky Blue
	SkyBlue = NewColor(102, 191, 255, 255)
	// Blue
	Blue = NewColor(0, 121, 241, 255)
	// Dark Blue
	DarkBlue = NewColor(0, 82, 172, 255)
	// Purple
	Purple = NewColor(200, 122, 255, 255)
	// Violet
	Violet = NewColor(135, 60, 190, 255)
	// Dark Purple
	DarkPurple = NewColor(112, 31, 126, 255)
	// Beige
	Beige = NewColor(211, 176, 131, 255)
	// Brown
	Brown = NewColor(127, 106, 79, 255)
	// Dark Brown
	DarkBrown = NewColor(76, 63, 47, 255)
	// White
	White = NewColor(255, 255, 255, 255)
	// Black
	Black = NewColor(0, 0, 0, 255)
	// Blank (Transparent)
	Blank = NewColor(0, 0, 0, 0)
	// Magenta
	Magenta = NewColor(255, 0, 255, 255)
	// Ray White (RayLib Logo White)
	RayWhite = NewColor(245, 245, 245, 255)
)

// Vector2 type
type Vector2 struct {
	X float32
	Y float32
}

// NewVector2 - Returns new Vector2
func NewVector2(x, y float32) Vector2 {
	return Vector2{x, y}
}

// newVector2FromPointer - Returns new Vector2 from pointer
func newVector2FromPointer(ptr unsafe.Pointer) Vector2 {
	return *(*Vector2)(ptr)
}

// Vector3 type
type Vector3 struct {
	X float32
	Y float32
	Z float32
}

// NewVector3 - Returns new Vector3
func NewVector3(X, Y, Z float32) Vector3 {
	return Vector3{X, Y, Z}
}

// newVector3FromPointer - Returns new Vector3 from pointer
func newVector3FromPointer(ptr unsafe.Pointer) Vector3 {
	return *(*Vector3)(ptr)
}

// Vector4 type
type Vector4 struct {
	X float32
	Y float32
	Z float32
	W float32
}

// NewVector4 - Returns new Vector4
func NewVector4(X, Y, Z, W float32) Vector4 {
	return Vector4{X, Y, Z, W}
}

// newVector4FromPointer - Returns new Vector4 from pointer
func newVector4FromPointer(ptr unsafe.Pointer) Vector4 {
	return *(*Vector4)(ptr)
}

// Matrix type (OpenGL style 4x4 - right handed, column major)
type Matrix struct {
	M0, M4, M8, M12  float32
	M1, M5, M9, M13  float32
	M2, M6, M10, M14 float32
	M3, M7, M11, M15 float32
}

// NewMatrix - Returns new Matrix
func NewMatrix(m0, m4, m8, m12, m1, m5, m9, m13, m2, m6, m10, m14, m3, m7, m11, m15 float32) Matrix {
	return Matrix{m0, m4, m8, m12, m1, m5, m9, m13, m2, m6, m10, m14, m3, m7, m11, m15}
}

// newMatrixFromPointer - Returns new Matrix from pointer
func newMatrixFromPointer(ptr unsafe.Pointer) Matrix {
	return *(*Matrix)(ptr)
}

// Mat2 type (used for polygon shape rotation matrix)
type Mat2 struct {
	M00 float32
	M01 float32
	M10 float32
	M11 float32
}

// NewMat2 - Returns new Mat2
func NewMat2(m0, m1, m10, m11 float32) Mat2 {
	return Mat2{m0, m1, m10, m11}
}

// Quaternion type
type Quaternion struct {
	X float32
	Y float32
	Z float32
	W float32
}

// NewQuaternion - Returns new Quaternion
func NewQuaternion(x, y, z, w float32) Quaternion {
	return Quaternion{x, y, z, w}
}

// Color type, RGBA (32bit)
// TODO remove later, keep type for now to not break code
type Color = color.RGBA

// NewColor - Returns new Color
func NewColor(r, g, b, a uint8) color.RGBA {
	return color.RGBA{r, g, b, a}
}

// newColorFromPointer - Returns new Color from pointer
func newColorFromPointer(ptr unsafe.Pointer) color.RGBA {
	return *(*color.RGBA)(ptr)
}

// Rectangle type
type Rectangle struct {
	X      float32
	Y      float32
	Width  float32
	Height float32
}

// NewRectangle - Returns new Rectangle
func NewRectangle(x, y, width, height float32) Rectangle {
	return Rectangle{x, y, width, height}
}

// newRectangleFromPointer - Returns new Rectangle from pointer
func newRectangleFromPointer(ptr unsafe.Pointer) Rectangle {
	return *(*Rectangle)(ptr)
}

// ToInt32 converts rectangle to int32 variant
func (r *Rectangle) ToInt32() RectangleInt32 {
	rect := RectangleInt32{}
	rect.X = int32(r.X)
	rect.Y = int32(r.Y)
	rect.Width = int32(r.Width)
	rect.Height = int32(r.Height)

	return rect
}

// RectangleInt32 type
type RectangleInt32 struct {
	X      int32
	Y      int32
	Width  int32
	Height int32
}

// ToFloat32 converts rectangle to float32 variant
func (r *RectangleInt32) ToFloat32() Rectangle {
	rect := Rectangle{}
	rect.X = float32(r.X)
	rect.Y = float32(r.Y)
	rect.Width = float32(r.Width)
	rect.Height = float32(r.Height)

	return rect
}

// Camera3D type, defines a camera position/orientation in 3d space
type Camera3D struct {
	// Camera position
	Position Vector3
	// Camera target it looks-at
	Target Vector3
	// Camera up vector (rotation over its axis)
	Up Vector3
	// Camera field-of-view apperture in Y (degrees) in perspective, used as near plane width in orthographic
	Fovy float32
	// Camera type, controlling projection type, either CameraPerspective or CameraOrthographic.
	Projection CameraProjection
}

// Camera type fallback, defaults to Camera3D
type Camera = Camera3D

// NewCamera3D - Returns new Camera3D
func NewCamera3D(pos, target, up Vector3, fovy float32, ct CameraProjection) Camera3D {
	return Camera3D{pos, target, up, fovy, ct}
}

// newCamera3DFromPointer - Returns new Camera3D from pointer
func newCamera3DFromPointer(ptr unsafe.Pointer) Camera3D {
	return *(*Camera3D)(ptr)
}

// Camera2D type, defines a 2d camera
type Camera2D struct {
	// Camera offset (displacement from target)
	Offset Vector2
	// Camera target (rotation and zoom origin)
	Target Vector2
	// Camera rotation in degrees
	Rotation float32
	// Camera zoom (scaling), should be 1.0f by default
	Zoom float32
}

// NewCamera2D - Returns new Camera2D
func NewCamera2D(offset, target Vector2, rotation, zoom float32) Camera2D {
	return Camera2D{offset, target, rotation, zoom}
}

// newCamera2DFromPointer - Returns new Camera2D from pointer
func newCamera2DFromPointer(ptr unsafe.Pointer) Camera2D {
	return *(*Camera2D)(ptr)
}

// BoundingBox type
type BoundingBox struct {
	// Minimum vertex box-corner
	Min Vector3
	// Maximum vertex box-corner
	Max Vector3
}

// NewBoundingBox - Returns new BoundingBox
func NewBoundingBox(min, max Vector3) BoundingBox {
	return BoundingBox{min, max}
}

// newBoundingBoxFromPointer - Returns new BoundingBox from pointer
func newBoundingBoxFromPointer(ptr unsafe.Pointer) BoundingBox {
	return *(*BoundingBox)(ptr)
}

// Asset file
type Asset interface {
	io.ReadSeeker
	io.Closer
}

// Gestures type
type Gestures int32

// Gestures types
// NOTE: It could be used as flags to enable only some gestures
const (
	GestureNone       Gestures = 0
	GestureTap        Gestures = 1
	GestureDoubletap  Gestures = 2
	GestureHold       Gestures = 4
	GestureDrag       Gestures = 8
	GestureSwipeRight Gestures = 16
	GestureSwipeLeft  Gestures = 32
	GestureSwipeUp    Gestures = 64
	GestureSwipeDown  Gestures = 128
	GesturePinchIn    Gestures = 256
	GesturePinchOut   Gestures = 512
)

// Shader location point type
const (
	LocVertexPosition = iota
	LocVertexTexcoord01
	LocVertexTexcoord02
	LocVertexNormal
	LocVertexTangent
	LocVertexColor
	LocMatrixMvp
	LocMatrixView
	LocMatrixProjection
	LocMatrixModel
	LocMatrixNormal
	LocVectorView
	LocColorDiffuse
	LocColorSpecular
	LocColorAmbient
	LocMapAlbedo
	LocMapMetalness
	LocMapNormal
	LocMapRoughness
	LocMapOcclusion
	LocMapEmission
	LocMapHeight
	LocMapCubemap
	LocMapIrradiance
	LocMapPrefilter
	LocMapBrdf
)

// Material map type
const (
	MapAlbedo = iota
	MapMetalness
	MapNormal
	MapRoughness
	MapOcclusion
	MapEmission
	MapHeight
	MapBrdg
	MapCubemap
	MapIrradiance
	MapPrefilter
)

// Material map type
const (
	MapDiffuse     = MapAlbedo
	MapSpecular    = MapMetalness
	LocMapDiffuse  = LocMapAlbedo
	LocMapSpecular = LocMapMetalness
)

// Shader and material limits
const (
	// Maximum number of predefined locations stored in shader struct
	MaxShaderLocations = 32
	// Maximum number of texture maps stored in shader struct
	MaxMaterialMaps = 12
)

// Mesh - Vertex data definning a mesh
type Mesh struct {
	// Number of vertices stored in arrays
	VertexCount int32
	// Number of triangles stored (indexed or not)
	TriangleCount int32
	// Vertex position (XYZ - 3 components per vertex) (shader-location = 0)
	Vertices *float32
	// Vertex texture coordinates (UV - 2 components per vertex) (shader-location = 1)
	Texcoords *float32
	// Vertex second texture coordinates (useful for lightmaps) (shader-location = 5)
	Texcoords2 *float32
	// Vertex normals (XYZ - 3 components per vertex) (shader-location = 2)
	Normals *float32
	// Vertex tangents (XYZ - 3 components per vertex) (shader-location = 4)
	Tangents *float32
	// Vertex colors (RGBA - 4 components per vertex) (shader-location = 3)
	Colors *uint8
	// Vertex indices (in case vertex data comes indexed)
	Indices *uint16
	// AnimVertices
	AnimVertices *float32
	// AnimNormals
	AnimNormals *float32
	// BoneIds
	BoneIds *int32
	// BoneWeights
	BoneWeights *float32
	// OpenGL Vertex Array Object id
	VaoID uint32
	// OpenGL Vertex Buffer Objects id (7 types of vertex data)
	VboID *uint32
}

// newMeshFromPointer - Returns new Mesh from pointer
func newMeshFromPointer(ptr unsafe.Pointer) Mesh {
	return *(*Mesh)(ptr)
}

// Material type
type Material struct {
	// Shader
	Shader Shader
	// Maps
	Maps *MaterialMap
	// Generic parameters (if required)
	Params [4]float32
}

// newMaterialFromPointer - Returns new Material from pointer
func newMaterialFromPointer(ptr unsafe.Pointer) Material {
	return *(*Material)(ptr)
}

// GetMap - Get pointer to MaterialMap by map type
func (mt Material) GetMap(index int32) *MaterialMap {
	return (*MaterialMap)(unsafe.Pointer(uintptr(unsafe.Pointer(mt.Maps)) + uintptr(index)*uintptr(unsafe.Sizeof(MaterialMap{}))))
}

// MaterialMap type
type MaterialMap struct {
	// Texture
	Texture Texture2D
	// Color
	Color color.RGBA
	// Value
	Value float32
}

// Model, meshes, materials and animation data
type Model struct {
	// Local transform matrix
	Transform     Matrix
	MeshCount     int32
	MaterialCount int32
	Meshes        *Mesh
	Materials     *Material
	MeshMaterial  *int32
	BoneCount     int32
	Bones         *BoneInfo
	BindPose      *Transform
}

// newModelFromPointer - Returns new Model from pointer
func newModelFromPointer(ptr unsafe.Pointer) Model {
	return *(*Model)(ptr)
}

// BoneInfo type
type BoneInfo struct {
	Name   [32]int8
	Parent int32
}

// Transform type
type Transform struct {
	Translation Vector3
	Rotation    Vector4
	Scale       Vector3
}

// Ray type (useful for raycast)
type Ray struct {
	// Ray position (origin)
	Position Vector3
	// Ray direction
	Direction Vector3
}

// NewRay - Returns new Ray
func NewRay(position, direction Vector3) Ray {
	return Ray{position, direction}
}

// newRayFromPointer - Returns new Ray from pointer
func newRayFromPointer(ptr unsafe.Pointer) Ray {
	return *(*Ray)(ptr)
}

// ModelAnimation type
type ModelAnimation struct {
	BoneCount  int32
	FrameCount int32
	Bones      *BoneInfo
	FramePoses **Transform
}

// newModelAnimationFromPointer - Returns new ModelAnimation from pointer
func newModelAnimationFromPointer(ptr unsafe.Pointer) ModelAnimation {
	return *(*ModelAnimation)(ptr)
}

// RayCollision type - ray hit information
type RayCollision struct {
	Hit      bool
	Distance float32
	Point    Vector3
	Normal   Vector3
}

// NewRayCollision - Returns new RayCollision
func NewRayCollision(hit bool, distance float32, point, normal Vector3) RayCollision {
	return RayCollision{hit, distance, point, normal}
}

// newRayCollisionFromPointer - Returns new RayCollision from pointer
func newRayCollisionFromPointer(ptr unsafe.Pointer) RayCollision {
	return *(*RayCollision)(ptr)
}

// BlendMode type
type BlendMode int32

// Color blending modes (pre-defined)
const (
	BlendAlpha          BlendMode = iota // Blend textures considering alpha (default)
	BlendAdditive                        // Blend textures adding colors
	BlendMultiplied                      // Blend textures multiplying colors
	BlendAddColors                       // Blend textures adding colors (alternative)
	BlendSubtractColors                  // Blend textures subtracting colors (alternative)
	BlendCustom                          // Blend textures using custom src/dst factors (use SetBlendModeCustom())
)

// Shader type (generic shader)
type Shader struct {
	// Shader program id
	ID uint32
	// Shader locations array
	Locs *int32
}

// NewShader - Returns new Shader
func NewShader(id uint32, locs *int32) Shader {
	return Shader{id, locs}
}

// newShaderFromPointer - Returns new Shader from pointer
func newShaderFromPointer(ptr unsafe.Pointer) Shader {
	return *(*Shader)(ptr)
}

// GetLocation - Get shader value's location
func (sh Shader) GetLocation(index int32) int32 {
	return *(*int32)(unsafe.Pointer(uintptr(unsafe.Pointer(sh.Locs)) + uintptr(index*4)))
}

// UpdateLocation - Update shader value's location
func (sh Shader) UpdateLocation(index int32, loc int32) {
	*(*int32)(unsafe.Pointer(uintptr(unsafe.Pointer(sh.Locs)) + uintptr(index*4))) = loc
}

// GlyphInfo - Font character info
type GlyphInfo struct {
	// Character value (Unicode)
	Value int32
	// Character offset X when drawing
	OffsetX int32
	// Character offset Y when drawing
	OffsetY int32
	// Character advance position X
	AdvanceX int32
	// Character image data
	Image Image
}

// NewGlyphInfo - Returns new CharInfo
func NewGlyphInfo(value int32, offsetX, offsetY, advanceX int32, image Image) GlyphInfo {
	return GlyphInfo{value, offsetX, offsetY, advanceX, image}
}

// newGlyphInfoFromPointer - Returns new GlyphInfo from pointer
func newGlyphInfoFromPointer(ptr unsafe.Pointer) GlyphInfo {
	return *(*GlyphInfo)(ptr)
}

// Font type, includes texture and charSet array data
type Font struct {
	// Base size (default chars height)
	BaseSize int32
	// Number of characters
	CharsCount int32
	// Padding around the chars
	CharsPadding int32
	// Characters texture atlas
	Texture Texture2D
	// Characters rectangles in texture
	Recs *Rectangle
	// Characters info data
	Chars *GlyphInfo
}

// newFontFromPointer - Returns new Font from pointer
func newFontFromPointer(ptr unsafe.Pointer) Font {
	return *(*Font)(ptr)
}

// PixelFormat - Texture format
type PixelFormat int32

// Texture formats
// NOTE: Support depends on OpenGL version and platform
const (
	// 8 bit per pixel (no alpha)
	UncompressedGrayscale PixelFormat = iota + 1
	// 8*2 bpp (2 channels)
	UncompressedGrayAlpha
	// 16 bpp
	UncompressedR5g6b5
	// 24 bpp
	UncompressedR8g8b8
	// 16 bpp (1 bit alpha)
	UncompressedR5g5b5a1
	// 16 bpp (4 bit alpha)
	UncompressedR4g4b4a4
	// 32 bpp
	UncompressedR8g8b8a8
	// 32 bpp (1 channel - float)
	UncompressedR32
	// 32*3 bpp (3 channels - float)
	UncompressedR32g32b32
	// 32*4 bpp (4 channels - float)
	UncompressedR32g32b32a32
	// 4 bpp (no alpha)
	CompressedDxt1Rgb
	// 4 bpp (1 bit alpha)
	CompressedDxt1Rgba
	// 8 bpp
	CompressedDxt3Rgba
	// 8 bpp
	CompressedDxt5Rgba
	// 4 bpp
	CompressedEtc1Rgb
	// 4 bpp
	CompressedEtc2Rgb
	// 8 bpp
	CompressedEtc2EacRgba
	// 4 bpp
	CompressedPvrtRgb
	// 4 bpp
	CompressedPvrtRgba
	// 8 bpp
	CompressedAstc4x4Rgba
	// 2 bpp
	CompressedAstc8x8Rgba
)

// TextureFilterMode - Texture filter mode
type TextureFilterMode int32

// Texture parameters: filter mode
// NOTE 1: Filtering considers mipmaps if available in the texture
// NOTE 2: Filter is accordingly set for minification and magnification
const (
	// No filter, just pixel aproximation
	FilterPoint TextureFilterMode = iota
	// Linear filtering
	FilterBilinear
	// Trilinear filtering (linear with mipmaps)
	FilterTrilinear
	// Anisotropic filtering 4x
	FilterAnisotropic4x
	// Anisotropic filtering 8x
	FilterAnisotropic8x
	// Anisotropic filtering 16x
	FilterAnisotropic16x
)

// TextureWrapMode - Texture wrap mode
type TextureWrapMode int32

// Texture parameters: wrap mode
const (
	WrapRepeat TextureWrapMode = iota
	WrapClamp
	WrapMirrorRepeat
	WrapMirrorClamp
)

// Image type, bpp always RGBA (32bit)
// NOTE: Data stored in CPU memory (RAM)
type Image struct {
	// Image raw data
	data unsafe.Pointer
	// Image base width
	Width int32
	// Image base height
	Height int32
	// Mipmap levels, 1 by default
	Mipmaps int32
	// Data format (PixelFormat)
	Format PixelFormat
}

// NewImage - Returns new Image
func NewImage(data []byte, width, height, mipmaps int32, format PixelFormat) *Image {
	d := unsafe.Pointer(&data[0])
	return &Image{d, width, height, mipmaps, format}
}

// newImageFromPointer - Returns new Image from pointer
func newImageFromPointer(ptr unsafe.Pointer) *Image {
	return (*Image)(ptr)
}

// NewImageFromImage - Returns new Image from Go image.Image
func NewImageFromImage(img image.Image) *Image {
	size := img.Bounds().Size()

	cx := (C.int)(size.X)
	cy := (C.int)(size.Y)
	ccolor := colorCptr(White)
	ret := C.GenImageColor(cx, cy, *ccolor)

	for y := 0; y < size.Y; y++ {
		for x := 0; x < size.X; x++ {
			color := img.At(x, y)
			r, g, b, a := color.RGBA()
			rcolor := NewColor(uint8(r), uint8(g), uint8(b), uint8(a))
			ccolor = colorCptr(rcolor)

			cx = (C.int)(x)
			cy = (C.int)(y)
			C.ImageDrawPixel(&ret, cx, cy, *ccolor)
		}
	}
	v := newImageFromPointer(unsafe.Pointer(&ret))
	return v
}

// Texture2D type, bpp always RGBA (32bit)
// NOTE: Data stored in GPU memory
type Texture2D struct {
	// OpenGL texture id
	ID uint32
	// Texture base width
	Width int32
	// Texture base height
	Height int32
	// Mipmap levels, 1 by default
	Mipmaps int32
	// Data format (PixelFormat)
	Format PixelFormat
}

// NewTexture2D - Returns new Texture2D
func NewTexture2D(id uint32, width, height, mipmaps int32, format PixelFormat) Texture2D {
	return Texture2D{id, width, height, mipmaps, format}
}

// newTexture2DFromPointer - Returns new Texture2D from pointer
func newTexture2DFromPointer(ptr unsafe.Pointer) Texture2D {
	return *(*Texture2D)(ptr)
}

// RenderTexture2D type, for texture rendering
type RenderTexture2D struct {
	// Render texture (fbo) id
	ID uint32
	// Color buffer attachment texture
	Texture Texture2D
	// Depth buffer attachment texture
	Depth Texture2D
}

// NewRenderTexture2D - Returns new RenderTexture2D
func NewRenderTexture2D(id uint32, texture, depth Texture2D) RenderTexture2D {
	return RenderTexture2D{id, texture, depth}
}

// newRenderTexture2DFromPointer - Returns new RenderTexture2D from pointer
func newRenderTexture2DFromPointer(ptr unsafe.Pointer) RenderTexture2D {
	return *(*RenderTexture2D)(ptr)
}

// Log message types
const (
	LogAll = iota
	LogTrace
	LogDebug
	LogInfo
	LogWarning
	LogError
	LogFatal
	LogNone
)

var logTypeFlags = LogInfo | LogWarning | LogError
