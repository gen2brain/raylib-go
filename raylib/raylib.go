/*
Package raylib - Go bindings for raylib, a simple and easy-to-use library to learn videogames programming.

raylib is highly inspired by Borland BGI graphics lib and by XNA framework.

raylib could be useful for prototyping, tools development, graphic applications, embedded systems and education.

NOTE for ADVENTURERS: raylib is a programming library to learn videogames programming; no fancy interface, no visual helpers, no auto-debugging... just coding in the most pure spartan-programmers way.

Example:

	package main

	import "github.com/gen2brain/raylib-go/raylib"

	func main() {
		raylib.InitWindow(800, 450, "raylib [core] example - basic window")

		raylib.SetTargetFPS(60)

		for !raylib.WindowShouldClose() {
			raylib.BeginDrawing()

			raylib.ClearBackground(raylib.RayWhite)

			raylib.DrawText("Congrats! You created your first window!", 190, 200, 20, raylib.LightGray)

			raylib.EndDrawing()
		}

		raylib.CloseWindow()
	}


*/
package raylib

import (
	"image"
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

// Sound source type
type Sound struct {
	// Audio source id
	Source uint32
	// Audio buffer id
	Buffer uint32
	// Audio format specifier
	Format int32
}

// NewSound - Returns new Sound
func NewSound(source, buffer uint32, format int32) Sound {
	return Sound{source, buffer, format}
}

// Music type (file streaming from memory)
// NOTE: Anything longer than ~10 seconds should be streamed
type Music struct {
	CtxType      uint32
	_            [4]byte
	ctxOgg       unsafe.Pointer
	ctxFlac      unsafe.Pointer
	ctxXm        unsafe.Pointer
	ctxMod       unsafe.Pointer
	Stream       AudioStream
	LoopCount    int32
	TotalSamples uint32
	SamplesLeft  uint32
}

// AudioStream type
// NOTE: Useful to create custom audio streams not bound to a specific file
type AudioStream struct {
	// Frequency (samples per second)
	SampleRate uint32
	// Bit depth (bits per sample): 8, 16, 32 (24 not supported)
	SampleSize uint32
	// Number of channels (1-mono, 2-stereo)
	Channels uint32
	// Audio format specifier
	Format int32
	// Audio source id
	Source uint32
	// Audio buffers (double buffering)
	Buffers [2]uint32
}

// NewAudioStream - Returns new AudioStream
func NewAudioStream(sampleRate, sampleSize, channels uint32, format int32, source uint32, buffers [2]uint32) AudioStream {
	return AudioStream{sampleRate, sampleSize, channels, format, source, buffers}
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

// Some basic Defines
const (
	Pi      = 3.1415927
	Deg2rad = 0.017453292
	Rad2deg = 57.295776

	// Raylib Config Flags

	// Set to show raylib logo at startup
	FlagShowLogo = 1
	// Set to run program in fullscreen
	FlagFullscreenMode = 2
	// Set to allow resizable window
	FlagWindowResizable = 4
	// Set to show window decoration (frame and buttons)
	FlagWindowDecorated = 8
	// Set to allow transparent window
	FlagWindowTransparent = 16
	// Set to try enabling MSAA 4X
	FlagMsaa4xHint = 32
	// Set to try enabling V-Sync on GPU
	FlagVsyncHint = 64

	// Keyboard Function Keys
	KeySpace        = 32
	KeyEscape       = 256
	KeyEnter        = 257
	KeyBackspace    = 259
	KeyRight        = 262
	KeyLeft         = 263
	KeyDown         = 264
	KeyUp           = 265
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
	KeyRightShift   = 344
	KeyRightControl = 345
	KeyRightAlt     = 346

	// Keyboard Alpha Numeric Keys
	KeyZero  = 48
	KeyOne   = 49
	KeyTwo   = 50
	KeyThree = 51
	KeyFour  = 52
	KeyFive  = 53
	KeySix   = 54
	KeySeven = 55
	KeyEight = 56
	KeyNine  = 57
	KeyA     = 65
	KeyB     = 66
	KeyC     = 67
	KeyD     = 68
	KeyE     = 69
	KeyF     = 70
	KeyG     = 71
	KeyH     = 72
	KeyI     = 73
	KeyJ     = 74
	KeyK     = 75
	KeyL     = 76
	KeyM     = 77
	KeyN     = 78
	KeyO     = 79
	KeyP     = 80
	KeyQ     = 81
	KeyR     = 82
	KeyS     = 83
	KeyT     = 84
	KeyU     = 85
	KeyV     = 86
	KeyW     = 87
	KeyX     = 88
	KeyY     = 89
	KeyZ     = 90

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
type Color struct {
	R uint8
	G uint8
	B uint8
	A uint8
}

// NewColor - Returns new Color
func NewColor(r, g, b, a uint8) Color {
	return Color{r, g, b, a}
}

// Rectangle type
type Rectangle struct {
	X      int32
	Y      int32
	Width  int32
	Height int32
}

// NewRectangle - Returns new Rectangle
func NewRectangle(x, y, width, height int32) Rectangle {
	return Rectangle{x, y, width, height}
}

// Camera type, defines a camera position/orientation in 3d space
type Camera struct {
	// Camera position
	Position Vector3
	// Camera target it looks-at
	Target Vector3
	// Camera up vector (rotation over its axis)
	Up Vector3
	// Camera field-of-view apperture in Y (degrees)
	Fovy float32
}

// NewCamera - Returns new Camera
func NewCamera(pos, target, up Vector3, fovy float32) Camera {
	return Camera{pos, target, up, fovy}
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
	LocMatrixModel
	LocMatrixView
	LocMatrixProjection
	LocVectorView
	LocColorDiffuse
	LocColorSpecular
	LocColorAmbient
	LocMapAlbedo
	LocMapMetalness
	LocMapNormal
	LocMapRoughness
	LocMapOccusion
	LocMapEmission
	LocMapHeight
	LocMapCubemap
	LocMapIrradiance
	LocMapPrefilter
	LocMapBrdf
)

// Material map type
const (
	// MapDiffuse
	MapAlbedo = iota
	MapMetalness
	MapNormal
	MapRoughness
	MapOcclusion
	MapEmission
	MapHeight
	// NOTE: Uses GL_TEXTURE_CUBE_MAP
	MapCubemap
	// NOTE: Uses GL_TEXTURE_CUBE_MAP
	MapIrradiance
	// NOTE: Uses GL_TEXTURE_CUBE_MAP
	MapPrefilter
	MapBrdf
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
	Vertices *[]float32
	// Vertex texture coordinates (UV - 2 components per vertex) (shader-location = 1)
	Texcoords *[]float32
	// Vertex second texture coordinates (useful for lightmaps) (shader-location = 5)
	Texcoords2 *[]float32
	// Vertex normals (XYZ - 3 components per vertex) (shader-location = 2)
	Normals *[]float32
	// Vertex tangents (XYZ - 3 components per vertex) (shader-location = 4)
	Tangents *[]float32
	// Vertex colors (RGBA - 4 components per vertex) (shader-location = 3)
	Colors *[]uint8
	// Vertex indices (in case vertex data comes indexed)
	Indices *[]uint16
	// OpenGL Vertex Array Object id
	VaoID uint32
	// OpenGL Vertex Buffer Objects id (7 types of vertex data)
	VboID [7]uint32
}

// NewMesh - Returns new Mesh
func NewMesh(vertexCount, triangleCount int32, vertices, texcoords, texcoords2, normals, tangents *[]float32, colors *[]uint8, indices *[]uint16, vaoID uint32, vboID [7]uint32) Mesh {
	return Mesh{vertexCount, triangleCount, vertices, texcoords, texcoords2, normals, tangents, colors, indices, vaoID, vboID}
}

// Material type
type Material struct {
	// Shader
	Shader Shader
	// Maps
	Maps [MaxMaterialMaps]MaterialMap
	// Padding
	_ [4]byte
	// Generic parameters (if required)
	Params *[]float32
}

// NewMaterial - Returns new Material
func NewMaterial(shader Shader, maps [MaxMaterialMaps]MaterialMap, params *[]float32) Material {
	return Material{shader, maps, [4]byte{}, params}
}

// MaterialMap type
type MaterialMap struct {
	// Texture
	Texture Texture2D
	// Color
	Color Color
	// Value
	Value float32
}

// Model type
type Model struct {
	// Vertex data buffers (RAM and VRAM)
	Mesh Mesh
	// Local transform matrix
	Transform Matrix
	// Shader and textures data
	Material Material
	// Padding
	_ [4]byte
}

// NewModel - Returns new Model
func NewModel(mesh Mesh, transform Matrix, material Material) Model {
	return Model{mesh, transform, material, [4]byte{}}
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

// VrDevice type
type VrDevice int32

// Head Mounted Display devices
const (
	HmdDefaultDevice VrDevice = iota
	HmdOculusRiftDk2
	HmdOculusRiftCv1
	HmdOculusGo
	HmdValveHtcVive
	HmdSonyPsvr
)

// VrDeviceInfo - Head-Mounted-Display device parameters
type VrDeviceInfo struct {
	// HMD horizontal resolution in pixels
	hResolution int
	// HMD vertical resolution in pixels
	vResolution int
	// HMD horizontal size in meters
	hScreenSize float32
	// HMD vertical size in meters
	vScreenSize float32
	// HMD screen center in meters
	vScreenCenter float32
	// HMD distance between eye and display in meters
	eyeToScreenDistance float32
	// HMD lens separation distance in meters
	lensSeparationDistance float32
	// HMD IPD (distance between pupils) in meters
	interpupillaryDistance float32
	// HMD lens distortion constant parameters
	lensDistortionValues [4]float32
	// HMD chromatic aberration correction parameters
	chromaAbCorrection [4]float32
}

// NewVrDeviceInfo - Returns new VrDeviceInfo
func NewVrDeviceInfo(hResolution, vResolution int, hScreenSize, vScreenSize, vScreenCenter, eyeToScreenDistance,
	lensSeparationDistance, interpupillaryDistance float32, lensDistortionValues, chromaAbCorrection [4]float32) VrDeviceInfo {

	return VrDeviceInfo{hResolution, vResolution, hScreenSize, vScreenSize, vScreenCenter, eyeToScreenDistance,
		lensSeparationDistance, interpupillaryDistance, lensDistortionValues, chromaAbCorrection}
}

// BlendMode type
type BlendMode int32

// Color blending modes (pre-defined)
const (
	BlendAlpha BlendMode = iota
	BlendAdditive
	BlendMultiplied
)

// Shader type (generic shader)
type Shader struct {
	// Shader program id
	ID uint32
	// Shader locations array
	Locs [MaxShaderLocations]int32
}

// NewShader - Returns new Shader
func NewShader(id uint32, locs [MaxShaderLocations]int32) Shader {
	return Shader{id, locs}
}

// CharInfo - SpriteFont character info
type CharInfo struct {
	// Character value (Unicode)
	Value int32
	// Character rectangle in sprite font
	Rec Rectangle
	// Character offset X when drawing
	OffsetX int32
	// Character offset Y when drawing
	OffsetY int32
	// Character advance position X
	AdvanceX int32
}

// NewCharInfo - Returns new SpriteFont
func NewCharInfo(value int32, rec Rectangle, offsetX, offsetY, advanceX int32) CharInfo {
	return CharInfo{value, rec, offsetX, offsetY, advanceX}
}

// SpriteFont type, includes texture and charSet array data
type SpriteFont struct {
	// Font texture
	Texture Texture2D
	// Base size (default chars height)
	BaseSize int32
	// Number of characters
	CharsCount int32
	// Characters info data
	Chars *CharInfo
}

// NewSpriteFont - Returns new SpriteFont
func NewSpriteFont(texture Texture2D, baseSize, charsCount int32, chars *CharInfo) SpriteFont {
	return SpriteFont{texture, baseSize, charsCount, chars}
}

// TextureFormat - Texture format
type TextureFormat int32

// Texture formats
// NOTE: Support depends on OpenGL version and platform
const (
	// 8 bit per pixel (no alpha)
	UncompressedGrayscale TextureFormat = iota + 1
	// 16 bpp (2 channels)
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
	WrapMirror
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
	// Data format (TextureFormat)
	Format TextureFormat
}

// ToImage converts a Image to Go image.Image
func (i *Image) ToImage() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, int(i.Width), int(i.Height)))

	// Get pixel data from image (RGBA 32bit)
	pixels := GetImageData(i)

	img.Pix = pixels

	return img
}

// NewImage - Returns new Image
func NewImage(data []byte, width, height, mipmaps int32, format TextureFormat) *Image {
	d := unsafe.Pointer(&data[0])
	return &Image{d, width, height, mipmaps, format}
}

// NewImageFromImage - Returns new Image from Go image.Image
func NewImageFromImage(img image.Image) *Image {
	size := img.Bounds().Size()
	pixels := make([]Color, size.X*size.Y)

	for y := 0; y < size.Y; y++ {
		for x := 0; x < size.X; x++ {
			color := img.At(x, y)
			r, g, b, a := color.RGBA()
			pixels[x+y*size.Y] = NewColor(uint8(r), uint8(g), uint8(b), uint8(a))
		}
	}

	return LoadImageEx(pixels, int32(size.X), int32(size.Y))
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
	// Data format (TextureFormat)
	Format TextureFormat
}

// NewTexture2D - Returns new Texture2D
func NewTexture2D(id uint32, width, height, mipmaps int32, format TextureFormat) Texture2D {
	return Texture2D{id, width, height, mipmaps, format}
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

// Log message types
const (
	LogInfo = iota
	LogWarning
	LogError
	LogDebug
)
