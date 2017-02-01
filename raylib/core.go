package raylib

/*
#include "raylib.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"
import "reflect"

// Some basic Defines
const (
	Pi      = 3.1415927
	Deg2rad = 0.017453292
	Rad2deg = 57.295776

	// Raylib Config Flags
	FlagFullscreenMode  = 1
	FlagResizableWindow = 2
	FlagShowLogo        = 4
	FlagShowMouseCursor = 8
	FlagCenteredMode    = 16
	FlagMsaa4xHint      = 32
	FlagVsyncHint       = 64

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

// Gestures
type Gestures int32

// Gestures type
// NOTE: It could be used as flags to enable only some gestures
const (
	GestureNone       Gestures = C.GESTURE_NONE
	GestureTap        Gestures = C.GESTURE_TAP
	GestureDoubletap  Gestures = C.GESTURE_DOUBLETAP
	GestureHold       Gestures = C.GESTURE_HOLD
	GestureDrag       Gestures = C.GESTURE_DRAG
	GestureSwipeRight Gestures = C.GESTURE_SWIPE_RIGHT
	GestureSwipeLeft  Gestures = C.GESTURE_SWIPE_LEFT
	GestureSwipeUp    Gestures = C.GESTURE_SWIPE_UP
	GestureSwipeDown  Gestures = C.GESTURE_SWIPE_DOWN
	GesturePinchIn    Gestures = C.GESTURE_PINCH_IN
	GesturePinchOut   Gestures = C.GESTURE_PINCH_OUT
)

// Camera mode
type CameraMode int32

// Camera system modes
const (
	CameraCustom      CameraMode = C.CAMERA_CUSTOM
	CameraFree        CameraMode = C.CAMERA_FREE
	CameraOrbital     CameraMode = C.CAMERA_ORBITAL
	CameraFirstPerson CameraMode = C.CAMERA_FIRST_PERSON
	CameraThirdPerson CameraMode = C.CAMERA_THIRD_PERSON
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

func (v *Vector2) cptr() *C.Vector2 {
	return (*C.Vector2)(unsafe.Pointer(v))
}

// Returns new Vector2
func NewVector2(x, y float32) Vector2 {
	return Vector2{x, y}
}

// Returns new Vector2 from pointer
func NewVector2FromPointer(ptr unsafe.Pointer) Vector2 {
	return *(*Vector2)(ptr)
}

// Vector3 type
type Vector3 struct {
	X float32
	Y float32
	Z float32
}

func (v *Vector3) cptr() *C.Vector3 {
	return (*C.Vector3)(unsafe.Pointer(v))
}

// Returns new Vector3
func NewVector3(X, Y, Z float32) Vector3 {
	return Vector3{X, Y, X}
}

// Returns new Vector3 from pointer
func NewVector3FromPointer(ptr unsafe.Pointer) Vector3 {
	return *(*Vector3)(ptr)
}

// Matrix type (OpenGL style 4x4 - right handed, column major)
type Matrix struct {
	M0, M4, M8, M12  float32
	M1, M5, M9, M13  float32
	M2, M6, M10, M14 float32
	M3, M7, M11, M15 float32
}

func (m *Matrix) cptr() *C.Matrix {
	return (*C.Matrix)(unsafe.Pointer(m))
}

// Returns new Matrix
func NewMatrix(m0, m4, m8, m12, m1, m5, m9, m13, m2, m6, m10, m14, m3, m7, m11, m15 float32) Matrix {
	return Matrix{m0, m4, m8, m12, m1, m5, m9, m13, m2, m6, m10, m14, m3, m7, m11, m15}
}

// Returns new Matrix from pointer
func NewMatrixFromPointer(ptr unsafe.Pointer) Matrix {
	return *(*Matrix)(ptr)
}

// Quaternion type
type Quaternion struct {
	X float32
	Y float32
	Z float32
	W float32
}

// Returns new Quaternion
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

func (color *Color) cptr() *C.Color {
	return (*C.Color)(unsafe.Pointer(color))
}

// Returns new Color
func NewColor(r, g, b, a uint8) Color {
	return Color{r, g, b, a}
}

// Returns new Color from pointer
func NewColorFromPointer(ptr unsafe.Pointer) Color {
	return *(*Color)(ptr)
}

// Rectangle type
type Rectangle struct {
	X      int32
	Y      int32
	Width  int32
	Height int32
}

func (r *Rectangle) cptr() *C.Rectangle {
	return (*C.Rectangle)(unsafe.Pointer(r))
}

// Returns new Rectangle
func NewRectangle(x, y, width, height int32) Rectangle {
	return Rectangle{x, y, width, height}
}

// Returns new Rectangle from pointer
func NewRectangleFromPointer(ptr unsafe.Pointer) Rectangle {
	return *(*Rectangle)(ptr)
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

func (c *Camera) cptr() *C.Camera {
	return (*C.Camera)(unsafe.Pointer(c))
}

// Returns new Camera
func NewCamera(pos, target, up Vector3, fovy float32) Camera {
	return Camera{pos, target, up, fovy}
}

// Returns new Camera from pointer
func NewCameraFromPointer(ptr unsafe.Pointer) Camera {
	return *(*Camera)(ptr)
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

func (c *Camera2D) cptr() *C.Camera2D {
	return (*C.Camera2D)(unsafe.Pointer(c))
}

// Returns new Camera2D
func NewCamera2D(offset, target Vector2, rotation, zoom float32) Camera2D {
	return Camera2D{offset, target, rotation, zoom}
}

// Returns new Camera2D from pointer
func NewCamera2DFromPointer(ptr unsafe.Pointer) Camera2D {
	return *(*Camera2D)(ptr)
}

// Bounding box
type BoundingBox struct {
	// Minimum vertex box-corner
	Min Vector3
	// Maximum vertex box-corner
	Max Vector3
}

func (b *BoundingBox) cptr() *C.BoundingBox {
	return (*C.BoundingBox)(unsafe.Pointer(b))
}

// Returns new BoundingBox
func NewBoundingBox(min, max Vector3) BoundingBox {
	return BoundingBox{min, max}
}

// Returns new BoundingBox from pointer
func NewBoundingBoxFromPointer(ptr unsafe.Pointer) BoundingBox {
	return *(*BoundingBox)(ptr)
}

// Close Window and Terminate Context
func CloseWindow() {
	C.CloseWindow()
}

// Detect if KEY_ESCAPE pressed or Close icon pressed
func WindowShouldClose() bool {
	ret := C.WindowShouldClose()
	v := bool(int(ret) == 1)
	return v
}

// Detect if window has been minimized (or lost focus)
func IsWindowMinimized() bool {
	ret := C.IsWindowMinimized()
	v := bool(int(ret) == 1)
	return v
}

// Fullscreen toggle (only PLATFORM_DESKTOP)
func ToggleFullscreen() {
	C.ToggleFullscreen()
}

// Get current screen width
func GetScreenWidth() int32 {
	ret := C.GetScreenWidth()
	v := (int32)(ret)
	return v
}

// Get current screen height
func GetScreenHeight() int32 {
	ret := C.GetScreenHeight()
	v := (int32)(ret)
	return v
}

// Sets Background Color
func ClearBackground(color Color) {
	ccolor := color.cptr()
	C.ClearBackground(*ccolor)
}

// Setup drawing canvas to start drawing
func BeginDrawing() {
	C.BeginDrawing()
}

// End canvas drawing and Swap Buffers (Double Buffering)
func EndDrawing() {
	C.EndDrawing()
}

// Initialize 2D mode with custom camera
func Begin2dMode(camera Camera2D) {
	ccamera := camera.cptr()
	C.Begin2dMode(*ccamera)
}

// Ends 2D mode custom camera usage
func End2dMode() {
	C.End2dMode()
}

// Initializes 3D mode for drawing (Camera setup)
func Begin3dMode(camera Camera) {
	ccamera := camera.cptr()
	C.Begin3dMode(*ccamera)
}

// Ends 3D mode and returns to default 2D orthographic mode
func End3dMode() {
	C.End3dMode()
}

// Initializes render texture for drawing
func BeginTextureMode(target RenderTexture2D) {
	ctarget := target.cptr()
	C.BeginTextureMode(*ctarget)
}

// Ends drawing to render texture
func EndTextureMode() {
	C.EndTextureMode()
}

// Returns a ray trace from mouse position
func GetMouseRay(mousePosition Vector2, camera Camera) Ray {
	cmousePosition := mousePosition.cptr()
	ccamera := camera.cptr()
	ret := C.GetMouseRay(*cmousePosition, *ccamera)
	v := NewRayFromPointer(unsafe.Pointer(&ret))
	return v
}

// Returns the screen space position from a 3d world space position
func GetWorldToScreen(position Vector3, camera Camera) Vector2 {
	cposition := position.cptr()
	ccamera := camera.cptr()
	ret := C.GetWorldToScreen(*cposition, *ccamera)
	v := NewVector2FromPointer(unsafe.Pointer(&ret))
	return v
}

// Returns camera transform matrix (view matrix)
func GetCameraMatrix(camera Camera) Matrix {
	ccamera := camera.cptr()
	ret := C.GetCameraMatrix(*ccamera)
	v := NewMatrixFromPointer(unsafe.Pointer(&ret))
	return v
}

// Set target FPS (maximum)
func SetTargetFPS(fps int32) {
	cfps := (C.int)(fps)
	C.SetTargetFPS(cfps)
}

// Returns current FPS
func GetFPS() float32 {
	ret := C.GetFPS()
	v := (float32)(ret)
	return v
}

// Returns time in seconds for one frame
func GetFrameTime() float32 {
	ret := C.GetFrameTime()
	v := (float32)(ret)
	return v
}

// Returns a Color struct from hexadecimal value
func GetColor(hexValue int32) Color {
	chexValue := (C.int)(hexValue)
	ret := C.GetColor(chexValue)
	v := NewColorFromPointer(unsafe.Pointer(&ret))
	return v
}

// Returns hexadecimal value for a Color
func GetHexValue(color Color) int32 {
	ccolor := color.cptr()
	ret := C.GetHexValue(*ccolor)
	v := (int32)(ret)
	return v
}

// Converts Color to float array and normalizes
func ColorToFloat(color Color) []float32 {
	ccolor := color.cptr()
	ret := C.ColorToFloat(*ccolor)

	var data []float32
	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&data)))
	sliceHeader.Cap = 4
	sliceHeader.Len = 4
	sliceHeader.Data = uintptr(unsafe.Pointer(ret))

	return data
}

// Converts Vector3 to float array
func VectorToFloat(vec Vector3) []float32 {
	cvec := vec.cptr()
	ret := C.VectorToFloat(*cvec)

	var data []float32
	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&data)))
	sliceHeader.Cap = 3
	sliceHeader.Len = 3
	sliceHeader.Data = uintptr(unsafe.Pointer(ret))

	return data
}

// Converts Matrix to float array
func MatrixToFloat(mat Matrix) []float32 {
	cmat := mat.cptr()
	ret := C.MatrixToFloat(*cmat)

	var data []float32
	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&data)))
	sliceHeader.Cap = 16
	sliceHeader.Len = 16
	sliceHeader.Data = uintptr(unsafe.Pointer(ret))

	return data
}

// Returns a random value between min and max (both included)
func GetRandomValue(min int32, max int32) int32 {
	cmin := (C.int)(min)
	cmax := (C.int)(max)
	ret := C.GetRandomValue(cmin, cmax)
	v := (int32)(ret)
	return v
}

// Color fade-in or fade-out, alpha goes from 0.0f to 1.0f
func Fade(color Color, alpha float32) Color {
	ccolor := color.cptr()
	calpha := (C.float)(alpha)
	ret := C.Fade(*ccolor, calpha)
	v := NewColorFromPointer(unsafe.Pointer(&ret))
	return v
}

// Setup some window configuration flags
func SetConfigFlags(flags byte) {
	cflags := (C.char)(flags)
	C.SetConfigFlags(cflags)
}

// Activates raylib logo at startup (can be done with flags)
func ShowLogo() {
	C.ShowLogo()
}

// Storage save integer value (to defined position)
func StorageSaveValue(position int32, value int32) {
	cposition := (C.int)(position)
	cvalue := (C.int)(value)
	C.StorageSaveValue(cposition, cvalue)
}

// Storage load integer value (from defined position)
func StorageLoadValue(position int32) int32 {
	cposition := (C.int)(position)
	ret := C.StorageLoadValue(cposition)
	v := (int32)(ret)
	return v
}

// Detect if a key has been pressed once
func IsKeyPressed(key int32) bool {
	ckey := (C.int)(key)
	ret := C.IsKeyPressed(ckey)
	v := bool(int(ret) == 1)
	return v
}

// Detect if a key is being pressed
func IsKeyDown(key int32) bool {
	ckey := (C.int)(key)
	ret := C.IsKeyDown(ckey)
	v := bool(int(ret) == 1)
	return v
}

// Detect if a key has been released once
func IsKeyReleased(key int32) bool {
	ckey := (C.int)(key)
	ret := C.IsKeyReleased(ckey)
	v := bool(int(ret) == 1)
	return v
}

// Detect if a key is NOT being pressed
func IsKeyUp(key int32) bool {
	ckey := (C.int)(key)
	ret := C.IsKeyUp(ckey)
	v := bool(int(ret) == 1)
	return v
}

// Get latest key pressed
func GetKeyPressed() int32 {
	ret := C.GetKeyPressed()
	v := (int32)(ret)
	return v
}

// Set a custom key to exit program (default is ESC)
func SetExitKey(key int32) {
	ckey := (C.int)(key)
	C.SetExitKey(ckey)
}

// Detect if a gamepad is available
func IsGamepadAvailable(gamepad int32) bool {
	cgamepad := (C.int)(gamepad)
	ret := C.IsGamepadAvailable(cgamepad)
	v := bool(int(ret) == 1)
	return v
}

// Check gamepad name (if available)
func IsGamepadName(gamepad int32, name string) bool {
	cgamepad := (C.int)(gamepad)
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	ret := C.IsGamepadName(cgamepad, cname)
	v := bool(int(ret) == 1)
	return v
}

// Return gamepad internal name id
func GetGamepadName(gamepad int32) string {
	cgamepad := (C.int)(gamepad)
	ret := C.GetGamepadName(cgamepad)
	v := C.GoString(ret)
	return v
}

// Detect if a gamepad button has been pressed once
func IsGamepadButtonPressed(gamepad int32, button int32) bool {
	cgamepad := (C.int)(gamepad)
	cbutton := (C.int)(button)
	ret := C.IsGamepadButtonPressed(cgamepad, cbutton)
	v := bool(int(ret) == 1)
	return v
}

// Detect if a gamepad button is being pressed
func IsGamepadButtonDown(gamepad int32, button int32) bool {
	cgamepad := (C.int)(gamepad)
	cbutton := (C.int)(button)
	ret := C.IsGamepadButtonDown(cgamepad, cbutton)
	v := bool(int(ret) == 1)
	return v
}

// Detect if a gamepad button has been released once
func IsGamepadButtonReleased(gamepad int32, button int32) bool {
	cgamepad := (C.int)(gamepad)
	cbutton := (C.int)(button)
	ret := C.IsGamepadButtonReleased(cgamepad, cbutton)
	v := bool(int(ret) == 1)
	return v
}

// Detect if a gamepad button is NOT being pressed
func IsGamepadButtonUp(gamepad int32, button int32) bool {
	cgamepad := (C.int)(gamepad)
	cbutton := (C.int)(button)
	ret := C.IsGamepadButtonUp(cgamepad, cbutton)
	v := bool(int(ret) == 1)
	return v
}

// Get the last gamepad button pressed
func GetGamepadButtonPressed() int32 {
	ret := C.GetGamepadButtonPressed()
	v := (int32)(ret)
	return v
}

// Return gamepad axis count for a gamepad
func GetGamepadAxisCount(gamepad int32) int32 {
	cgamepad := (C.int)(gamepad)
	ret := C.GetGamepadAxisCount(cgamepad)
	v := (int32)(ret)
	return v
}

// Return axis movement value for a gamepad axis
func GetGamepadAxisMovement(gamepad int32, axis int32) float32 {
	cgamepad := (C.int)(gamepad)
	caxis := (C.int)(axis)
	ret := C.GetGamepadAxisMovement(cgamepad, caxis)
	v := (float32)(ret)
	return v
}

// Detect if a mouse button has been pressed once
func IsMouseButtonPressed(button int32) bool {
	cbutton := (C.int)(button)
	ret := C.IsMouseButtonPressed(cbutton)
	v := bool(int(ret) == 1)
	return v
}

// Detect if a mouse button is being pressed
func IsMouseButtonDown(button int32) bool {
	cbutton := (C.int)(button)
	ret := C.IsMouseButtonDown(cbutton)
	v := bool(int(ret) == 1)
	return v
}

// Detect if a mouse button has been released once
func IsMouseButtonReleased(button int32) bool {
	cbutton := (C.int)(button)
	ret := C.IsMouseButtonReleased(cbutton)
	v := bool(int(ret) == 1)
	return v
}

// Detect if a mouse button is NOT being pressed
func IsMouseButtonUp(button int32) bool {
	cbutton := (C.int)(button)
	ret := C.IsMouseButtonUp(cbutton)
	v := bool(int(ret) == 1)
	return v
}

// Returns mouse position X
func GetMouseX() int32 {
	ret := C.GetMouseX()
	v := (int32)(ret)
	return v
}

// Returns mouse position Y
func GetMouseY() int32 {
	ret := C.GetMouseY()
	v := (int32)(ret)
	return v
}

// Returns mouse position XY
func GetMousePosition() Vector2 {
	ret := C.GetMousePosition()
	v := NewVector2FromPointer(unsafe.Pointer(&ret))
	return v
}

// Set mouse position XY
func SetMousePosition(position Vector2) {
	cposition := position.cptr()
	C.SetMousePosition(*cposition)
}

// Returns mouse wheel movement Y
func GetMouseWheelMove() int32 {
	ret := C.GetMouseWheelMove()
	v := (int32)(ret)
	return v
}

// Returns touch position X for touch point 0 (relative to screen size)
func GetTouchX() int32 {
	ret := C.GetTouchX()
	v := (int32)(ret)
	return v
}

// Returns touch position Y for touch point 0 (relative to screen size)
func GetTouchY() int32 {
	ret := C.GetTouchY()
	v := (int32)(ret)
	return v
}

// Returns touch position XY for a touch point index (relative to screen size)
func GetTouchPosition(index int32) Vector2 {
	cindex := (C.int)(index)
	ret := C.GetTouchPosition(cindex)
	v := NewVector2FromPointer(unsafe.Pointer(&ret))
	return v
}

// Enable a set of gestures using flags
func SetGesturesEnabled(gestureFlags uint32) {
	cgestureFlags := (C.uint)(gestureFlags)
	C.SetGesturesEnabled(cgestureFlags)
}

// Check if a gesture have been detected
func IsGestureDetected(gesture Gestures) bool {
	cgesture := (C.int)(gesture)
	ret := C.IsGestureDetected(cgesture)
	v := bool(int(ret) == 1)
	return v
}

// Get latest detected gesture
func GetGestureDetected() Gestures {
	ret := C.GetGestureDetected()
	v := (Gestures)(ret)
	return v
}

// Get touch points count
func GetTouchPointsCount() int32 {
	ret := C.GetTouchPointsCount()
	v := (int32)(ret)
	return v
}

// Get gesture hold time in milliseconds
func GetGestureHoldDuration() float32 {
	ret := C.GetGestureHoldDuration()
	v := (float32)(ret)
	return v
}

// Get gesture drag vector
func GetGestureDragVector() Vector2 {
	ret := C.GetGestureDragVector()
	v := NewVector2FromPointer(unsafe.Pointer(&ret))
	return v
}

// Get gesture drag angle
func GetGestureDragAngle() float32 {
	ret := C.GetGestureDragAngle()
	v := (float32)(ret)
	return v
}

// Get gesture pinch delta
func GetGesturePinchVector() Vector2 {
	ret := C.GetGesturePinchVector()
	v := NewVector2FromPointer(unsafe.Pointer(&ret))
	return v
}

// Get gesture pinch angle
func GetGesturePinchAngle() float32 {
	ret := C.GetGesturePinchAngle()
	v := (float32)(ret)
	return v
}

// Shows current FPS
func DrawFPS(posX int32, posY int32) {
	cposX := (C.int)(posX)
	cposY := (C.int)(posY)
	C.DrawFPS(cposX, cposY)
}
