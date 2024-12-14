/*
Package raylib - Go bindings for raylib, a simple and easy-to-use library to enjoy videogames programming.

raylib is highly inspired by Borland BGI graphics lib and by XNA framework.
raylib could be useful for prototyping, tools development, graphic applications, embedded systems and education.

NOTE for ADVENTURERS: raylib is a programming library to learn videogames programming; no fancy interface, no visual helpers, no auto-debugging... just coding in the most pure spartan-programmers way.
*/
package rl

import (
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
	FrameCount uint32
	// Frequency (samples per second)
	SampleRate uint32
	// Bit depth (bits per sample): 8, 16, 32 (24 not supported)
	SampleSize uint32
	// Number of channels (1-mono, 2-stereo)
	Channels uint32
	// Buffer data pointer
	Data unsafe.Pointer
}

// NewWave - Returns new Wave
func NewWave(sampleCount, sampleRate, sampleSize, channels uint32, data []byte) Wave {
	d := unsafe.Pointer(&data[0])

	return Wave{sampleCount, sampleRate, sampleSize, channels, d}
}

// AudioCallback function.
type AudioCallback func(data []float32, frames int)

// Sound source type
type Sound struct {
	Stream     AudioStream
	FrameCount uint32
	_          [4]byte
}

// Music type (file streaming from memory)
// NOTE: Anything longer than ~10 seconds should be streamed
type Music struct {
	Stream     AudioStream
	FrameCount uint32
	Looping    bool
	CtxType    int32
	CtxData    unsafe.Pointer
}

// AudioStream type
// NOTE: Useful to create custom audio streams not bound to a specific file
type AudioStream struct {
	// Buffer
	Buffer *AudioBuffer
	// Processor
	Processor *AudioProcessor
	// Frequency (samples per second)
	SampleRate uint32
	// Bit depth (bits per sample): 8, 16, 32 (24 not supported)
	SampleSize uint32
	// Number of channels (1-mono, 2-stereo)
	Channels uint32
	_        [4]byte
}

type maDataConverter struct {
	FormatIn                uint32
	FormatOut               uint32
	ChannelsIn              uint32
	ChannelsOut             uint32
	SampleRateIn            uint32
	SampleRateOut           uint32
	DitherMode              uint32
	ExecutionPath           uint32
	ChannelConverter        maChannelConverter
	Resampler               maResampler
	HasPreFormatConversion  uint8
	HasPostFormatConversion uint8
	HasChannelConverter     uint8
	HasResampler            uint8
	IsPassthrough           uint8
	X_ownsHeap              uint8
	X_pHeap                 *byte
}

type maChannelConverter struct {
	Format         uint32
	ChannelsIn     uint32
	ChannelsOut    uint32
	MixingMode     uint32
	ConversionPath uint32
	PChannelMapIn  *uint8
	PChannelMapOut *uint8
	PShuffleTable  *uint8
	Weights        [8]byte
	X_pHeap        *byte
	X_ownsHeap     uint32
	Pad_cgo_0      [4]byte
}

type maResampler struct {
	PBackend         *byte
	PBackendVTable   *maResamplingBackendVtable
	PBackendUserData *byte
	Format           uint32
	Channels         uint32
	SampleRateIn     uint32
	SampleRateOut    uint32
	State            [136]byte
	X_pHeap          *byte
	X_ownsHeap       uint32
	Pad_cgo_0        [4]byte
}

type maResamplingBackendVtable struct {
	OnGetHeapSize                 *[0]byte
	OnInit                        *[0]byte
	OnUninit                      *[0]byte
	OnProcess                     *[0]byte
	OnSetRate                     *[0]byte
	OnGetInputLatency             *[0]byte
	OnGetOutputLatency            *[0]byte
	OnGetRequiredInputFrameCount  *[0]byte
	OnGetExpectedOutputFrameCount *[0]byte
	OnReset                       *[0]byte
}

type AudioBuffer struct {
	Converter            maDataConverter
	Callback             *[0]byte
	Processor            *AudioProcessor
	Volume               float32
	Pitch                float32
	Pan                  float32
	Playing              bool
	Paused               bool
	Looping              bool
	Usage                int32
	IsSubBufferProcessed [2]bool
	SizeInFrames         uint32
	FrameCursorPos       uint32
	FramesProcessed      uint32
	Data                 *uint8
	Next                 *AudioBuffer
	Prev                 *AudioBuffer
}

type AudioProcessor struct {
	Process *[0]byte
	Next    *AudioProcessor
	Prev    *AudioProcessor
}

// AutomationEvent - Automation event
type AutomationEvent struct {
	Frame  uint32
	Type   uint32
	Params [4]int32
}

// AutomationEventList - Automation event list
type AutomationEventList struct {
	Capacity uint32
	Count    uint32
	// Events array (c array)
	//
	// Use AutomationEventList.GetEvents instead (go slice)
	Events *AutomationEvent
}

func (a *AutomationEventList) GetEvents() []AutomationEvent {
	return unsafe.Slice(a.Events, a.Count)
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
	// Set to support mouse passthrough, only supported when FLAG_WINDOW_UNDECORATED
	FlagWindowMousePassthrough = 0x00004000
	// Set to run program in borderless windowed mode
	FlagBorderlessWindowedMode = 0x00008000
	// Set to try enabling MSAA 4X
	FlagMsaa4xHint = 0x00000020
	// Set to try enabling interlaced video format (for V3D)
	FlagInterlacedHint = 0x00010000

	// KeyNull is used for no key pressed
	KeyNull = 0

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
	KeyMenu       = 5
	KeyVolumeUp   = 24
	KeyVolumeDown = 25

	MouseLeftButton   = MouseButtonLeft
	MouseRightButton  = MouseButtonRight
	MouseMiddleButton = MouseButtonMiddle
)

type MouseButton int32

// Mouse Buttons
const (
	MouseButtonLeft MouseButton = iota
	MouseButtonRight
	MouseButtonMiddle
	MouseButtonSide
	MouseButtonExtra
	MouseButtonForward
	MouseButtonBack
)

// Mouse cursor
type MouseCursor = int32

const (
	MouseCursorDefault      MouseCursor = iota // Default pointer shape
	MouseCursorArrow                           // Arrow shape
	MouseCursorIBeam                           // Text writing cursor shape
	MouseCursorCrosshair                       // Cross shape
	MouseCursorPointingHand                    // Pointing hand cursor
	MouseCursorResizeEW                        // Horizontal resize/move arrow shape
	MouseCursorResizeNS                        // Vertical resize/move arrow shape
	MouseCursorResizeNWSE                      // Top-left to bottom-right diagonal resize/move arrow shape
	MouseCursorResizeNESW                      // The top-right to bottom-left diagonal resize/move arrow shape
	MouseCursorResizeAll                       // The omni-directional resize/move cursor shape
	MouseCursorNotAllowed                      // The operation-not-allowed shape
)

// Gamepad Buttons
const (
	GamepadButtonUnknown        = iota // Unknown button, just for error checking
	GamepadButtonLeftFaceUp            // Gamepad left DPAD up button
	GamepadButtonLeftFaceRight         // Gamepad left DPAD right button
	GamepadButtonLeftFaceDown          // Gamepad left DPAD down button
	GamepadButtonLeftFaceLeft          // Gamepad left DPAD left button
	GamepadButtonRightFaceUp           // Gamepad right button up (i.e. PS3: Triangle, Xbox: Y)
	GamepadButtonRightFaceRight        // Gamepad right button right (i.e. PS3: Circle, Xbox: B)
	GamepadButtonRightFaceDown         // Gamepad right button down (i.e. PS3: Cross, Xbox: A)
	GamepadButtonRightFaceLeft         // Gamepad right button left (i.e. PS3: Square, Xbox: X)
	GamepadButtonLeftTrigger1          // Gamepad top/back trigger left (first), it could be a trailing button
	GamepadButtonLeftTrigger2          // Gamepad top/back trigger left (second), it could be a trailing button
	GamepadButtonRightTrigger1         // Gamepad top/back trigger right (first), it could be a trailing button
	GamepadButtonRightTrigger2         // Gamepad top/back trigger right (second), it could be a trailing button
	GamepadButtonMiddleLeft            // Gamepad center buttons, left one (i.e. PS3: Select)
	GamepadButtonMiddle                // Gamepad center buttons, middle one (i.e. PS3: PS, Xbox: XBOX)
	GamepadButtonMiddleRight           // Gamepad center buttons, right one (i.e. PS3: Start)
	GamepadButtonLeftThumb             // Gamepad joystick pressed button left
	GamepadButtonRightThumb            // Gamepad joystick pressed button right
)

// Gamepad Axis
const (
	GamepadAxisLeftX        = iota // Gamepad left stick X axis
	GamepadAxisLeftY               // Gamepad left stick Y axis
	GamepadAxisRightX              // Gamepad right stick X axis
	GamepadAxisRightY              // Gamepad right stick Y axis
	GamepadAxisLeftTrigger         // Gamepad back trigger left, pressure level: [1..-1]
	GamepadAxisRightTrigger        // Gamepad back trigger right, pressure level: [1..-1]
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
func NewVector3(x, y, z float32) Vector3 {
	return Vector3{x, y, z}
}

// Vector4 type
type Vector4 struct {
	X float32
	Y float32
	Z float32
	W float32
}

// NewVector4 - Returns new Vector4
func NewVector4(x, y, z, w float32) Vector4 {
	return Vector4{x, y, z, w}
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

// Quaternion, 4 components (Vector4 alias)
type Quaternion = Vector4

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
	ShaderLocVertexPosition = iota
	ShaderLocVertexTexcoord01
	ShaderLocVertexTexcoord02
	ShaderLocVertexNormal
	ShaderLocVertexTangent
	ShaderLocVertexColor
	ShaderLocMatrixMvp
	ShaderLocMatrixView
	ShaderLocMatrixProjection
	ShaderLocMatrixModel
	ShaderLocMatrixNormal
	ShaderLocVectorView
	ShaderLocColorDiffuse
	ShaderLocColorSpecular
	ShaderLocColorAmbient
	ShaderLocMapAlbedo
	ShaderLocMapMetalness
	ShaderLocMapNormal
	ShaderLocMapRoughness
	ShaderLocMapOcclusion
	ShaderLocMapEmission
	ShaderLocMapHeight
	ShaderLocMapCubemap
	ShaderLocMapIrradiance
	ShaderLocMapPrefilter
	ShaderLocMapBrdf

	ShaderLocMapDiffuse  = ShaderLocMapAlbedo
	ShaderLocMapSpecular = ShaderLocMapMetalness
)

// ShaderUniformDataType type
type ShaderUniformDataType int32

// ShaderUniformDataType enumeration
const (
	// Shader uniform type: float
	ShaderUniformFloat ShaderUniformDataType = iota
	// Shader uniform type: vec2 (2 float)
	ShaderUniformVec2
	// Shader uniform type: vec3 (3 float)
	ShaderUniformVec3
	// Shader uniform type: vec4 (4 float)
	ShaderUniformVec4
	// Shader uniform type: int
	ShaderUniformInt
	// Shader uniform type: ivec2 (2 int)
	ShaderUniformIvec2
	// Shader uniform type: ivec2 (3 int)
	ShaderUniformIvec3
	// Shader uniform type: ivec2 (4 int)
	ShaderUniformIvec4
	// Shader uniform type: unsigned int
	ShaderUniformUint
	// Shader uniform type: uivec2 (2 unsigned int)
	ShaderUniformUivec2
	// Shader uniform type: uivec3 (3 unsigned int)
	ShaderUniformUivec3
	// Shader uniform type: uivec4 (4 unsigned int)
	ShaderUniformUivec4
	// Shader uniform type: sampler2d
	ShaderUniformSampler2d
)

// Material map index
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
	MapBrdf

	MapDiffuse  = MapAlbedo
	MapSpecular = MapMetalness
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
	// Bones animated transformation matrices
	BoneMatrices *Matrix
	// Number of bones
	BoneCount int32
	// OpenGL Vertex Array Object id
	VaoID uint32
	// OpenGL Vertex Buffer Objects id (7 types of vertex data)
	VboID *uint32
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

// GetMap - Get pointer to MaterialMap by map type
func (mt Material) GetMap(index int32) *MaterialMap {
	return (*MaterialMap)(unsafe.Pointer(uintptr(unsafe.Pointer(mt.Maps)) + uintptr(index)*unsafe.Sizeof(MaterialMap{})))
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

// Model is struct of model, meshes, materials and animation data
type Model struct {
	// Local transform matrix
	Transform Matrix
	// Number of meshes
	MeshCount int32
	// Number of materials
	MaterialCount int32
	// Meshes array (c array)
	//
	// Use Model.GetMeshes instead (go slice)
	Meshes *Mesh
	// Materials array (c array)
	//
	// Use Model.GetMaterials instead (go slice)
	Materials *Material
	// Mesh material number
	MeshMaterial *int32
	// Number of bones
	BoneCount int32
	// Bones information (skeleton) (c array)
	//
	// Use Model.GetBones instead (go slice)
	Bones *BoneInfo
	// Bones base transformation (pose) (c array)
	//
	// Use Model.GetBindPose instead (go slice)
	BindPose *Transform
}

// GetMeshes returns the meshes of a model as go slice
func (m Model) GetMeshes() []Mesh {
	return unsafe.Slice(m.Meshes, m.MeshCount)
}

// GetMaterials returns the materials of a model as go slice
func (m Model) GetMaterials() []Material {
	return unsafe.Slice(m.Materials, m.MaterialCount)
}

// GetBones returns the bones information (skeleton) of a model as go slice
func (m Model) GetBones() []BoneInfo {
	return unsafe.Slice(m.Bones, m.BoneCount)
}

// GetBindPose returns the bones base transformation of a model as go slice
func (m Model) GetBindPose() []Transform {
	return unsafe.Slice(m.BindPose, m.BoneCount)
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

// ModelAnimation type
type ModelAnimation struct {
	BoneCount  int32
	FrameCount int32
	Bones      *BoneInfo
	FramePoses **Transform
	Name       [32]uint8
}

// GetBones returns the bones information (skeleton) of a ModelAnimation as go slice
func (m ModelAnimation) GetBones() []BoneInfo {
	return unsafe.Slice(m.Bones, m.BoneCount)
}

// GetFramePose returns the Transform for a specific bone at a specific frame
func (m ModelAnimation) GetFramePose(frame, bone int) Transform {
	framePoses := unsafe.Slice(m.FramePoses, m.FrameCount)
	return unsafe.Slice(framePoses[frame], m.BoneCount)[bone]
}

// GetName returns the ModelAnimation's name as go string
func (m ModelAnimation) GetName() string {
	var end int
	for end = range m.Name {
		if m.Name[end] == 0 {
			break
		}
	}
	return string(m.Name[:end])
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

// BlendMode type
type BlendMode int32

// Color blending modes (pre-defined)
const (
	BlendAlpha            BlendMode = iota // Blend textures considering alpha (default)
	BlendAdditive                          // Blend textures adding colors
	BlendMultiplied                        // Blend textures multiplying colors
	BlendAddColors                         // Blend textures adding colors (alternative)
	BlendSubtractColors                    // Blend textures subtracting colors (alternative)
	BlendAlphaPremultiply                  // Blend premultiplied textures considering alpha
	BlendCustom                            // Blend textures using custom src/dst factors
	BlendCustomSeparate                    // Blend textures using custom rgb/alpha separate src/dst factors
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

// Font type, defines generation method
const (
	FontDefault = iota // Default font generation, anti-aliased
	FontBitmap         // Bitmap font generation, no anti-aliasing
	FontSdf            // SDF font generation, requires external shader
)

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

// Cubemap layouts
const (
	CubemapLayoutAutoDetect       = iota // Automatically detect layout type
	CubemapLayoutLineVertical            // Layout is defined by a vertical line with faces
	CubemapLayoutLineHorizontal          // Layout is defined by a horizontal line with faces
	CubemapLayoutCrossThreeByFour        // Layout is defined by a 3x4 cross with cubemap faces
	CubemapLayoutCrossFourByThree        // Layout is defined by a 4x3 cross with cubemap faces
)

// Image type, bpp always RGBA (32bit)
// NOTE: Data stored in CPU memory (RAM)
type Image struct {
	// Image raw Data
	Data unsafe.Pointer
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

// TraceLogCallbackFun - function that will recive the trace log messages
type TraceLogCallbackFun func(int, string)

// TraceLogLevel parameter of trace log message
type TraceLogLevel int

// Trace log level
// NOTE: Organized by priority level
const (
	// Display all logs
	LogAll TraceLogLevel = iota
	// Trace logging, intended for internal use only
	LogTrace
	// Debug logging, used for internal debugging, it should be disabled on release builds
	LogDebug
	// Info logging, used for program execution info
	LogInfo
	// Warning logging, used on recoverable failures
	LogWarning
	// Error logging, used on unrecoverable failures
	LogError
	// Fatal logging, used to abort program: exit(EXIT_FAILURE)
	LogFatal
	// Disable logging
	LogNone
)

// N-patch layout
type NPatchLayout int32

const (
	NPatchNinePatch            NPatchLayout = iota // Npatch layout: 3x3 tiles
	NPatchThreePatchVertical                       // Npatch layout: 1x3 tiles
	NPatchThreePatchHorizontal                     // Npatch layout: 3x1 tiles
)

// NPatchInfo type, n-patch layout info
type NPatchInfo struct {
	Source Rectangle    // Texture source rectangle
	Left   int32        // Left border offset
	Top    int32        // Top border offset
	Right  int32        // Right border offset
	Bottom int32        // Bottom border offset
	Layout NPatchLayout // Layout of the n-patch: 3x3, 1x3 or 3x1
}

// VrStereoConfig, VR stereo rendering configuration for simulator
type VrStereoConfig struct {
	Projection        [2]Matrix  // VR projection matrices (per eye)
	ViewOffset        [2]Matrix  // VR view offset matrices (per eye)
	LeftLensCenter    [2]float32 // VR left lens center
	RightLensCenter   [2]float32 // VR right lens center
	LeftScreenCenter  [2]float32 // VR left screen center
	RightScreenCenter [2]float32 // VR right screen center
	Scale             [2]float32 // VR distortion scale
	ScaleIn           [2]float32 // VR distortion scale in
}

// VrDeviceInfo, Head-Mounted-Display device parameters
type VrDeviceInfo struct {
	HResolution            int32      // Horizontal resolution in pixels
	VResolution            int32      // Vertical resolution in pixels
	HScreenSize            float32    // Horizontal size in meters
	VScreenSize            float32    // Vertical size in meters
	EyeToScreenDistance    float32    // Distance between eye and display in meters
	LensSeparationDistance float32    // Lens separation distance in meters
	InterpupillaryDistance float32    // IPD (distance between pupils) in meters
	LensDistortionValues   [4]float32 // Lens distortion constant parameters
	ChromaAbCorrection     [4]float32 // Chromatic aberration correction parameters
}
