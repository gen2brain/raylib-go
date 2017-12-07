package raylib

/*
#include "raylib.h"
#include <stdlib.h>
*/
import "C"

import (
	"io"
	"unsafe"
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

func (v *Vector2) cptr() *C.Vector2 {
	return (*C.Vector2)(unsafe.Pointer(v))
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

func (v *Vector3) cptr() *C.Vector3 {
	return (*C.Vector3)(unsafe.Pointer(v))
}

// NewVector3 - Returns new Vector3
func NewVector3(X, Y, Z float32) Vector3 {
	return Vector3{X, Y, Z}
}

// newVector3FromPointer - Returns new Vector3 from pointer
func newVector3FromPointer(ptr unsafe.Pointer) Vector3 {
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
type Color struct {
	R uint8
	G uint8
	B uint8
	A uint8
}

func (color *Color) cptr() *C.Color {
	return (*C.Color)(unsafe.Pointer(color))
}

// NewColor - Returns new Color
func NewColor(r, g, b, a uint8) Color {
	return Color{r, g, b, a}
}

// newColorFromPointer - Returns new Color from pointer
func newColorFromPointer(ptr unsafe.Pointer) Color {
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

// NewRectangle - Returns new Rectangle
func NewRectangle(x, y, width, height int32) Rectangle {
	return Rectangle{x, y, width, height}
}

// newRectangleFromPointer - Returns new Rectangle from pointer
func newRectangleFromPointer(ptr unsafe.Pointer) Rectangle {
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

// NewCamera - Returns new Camera
func NewCamera(pos, target, up Vector3, fovy float32) Camera {
	return Camera{pos, target, up, fovy}
}

// newCameraFromPointer - Returns new Camera from pointer
func newCameraFromPointer(ptr unsafe.Pointer) Camera {
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

func (b *BoundingBox) cptr() *C.BoundingBox {
	return (*C.BoundingBox)(unsafe.Pointer(b))
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

// CloseWindow - Close Window and Terminate Context
func CloseWindow() {
	C.CloseWindow()
}

// WindowShouldClose - Detect if KEY_ESCAPE pressed or Close icon pressed
func WindowShouldClose() bool {
	ret := C.WindowShouldClose()
	v := bool(int(ret) == 1)
	return v
}

// IsWindowMinimized - Detect if window has been minimized (or lost focus)
func IsWindowMinimized() bool {
	ret := C.IsWindowMinimized()
	v := bool(int(ret) == 1)
	return v
}

// ToggleFullscreen - Fullscreen toggle (only PLATFORM_DESKTOP)
func ToggleFullscreen() {
	C.ToggleFullscreen()
}

// SetWindowIcon - Set icon for window (only PLATFORM_DESKTOP)
func SetWindowIcon(image Image) {
	cimage := image.cptr()
	C.SetWindowIcon(*cimage)
}

// SetWindowTitle - Set title for window (only PLATFORM_DESKTOP)
func SetWindowTitle(title string) {
	ctitle := C.CString(title)
	defer C.free(unsafe.Pointer(ctitle))
	C.SetWindowTitle(ctitle)
}

// SetWindowPosition - Set window position on screen (only PLATFORM_DESKTOP)
func SetWindowPosition(x, y int32) {
	cx := (C.int)(x)
	cy := (C.int)(y)
	C.SetWindowPosition(cx, cy)
}

// SetWindowMonitor - Set monitor for the current window (fullscreen mode)
func SetWindowMonitor(monitor int32) {
	cmonitor := (C.int)(monitor)
	C.SetWindowMonitor(cmonitor)
}

// GetScreenWidth - Get current screen width
func GetScreenWidth() int32 {
	ret := C.GetScreenWidth()
	v := (int32)(ret)
	return v
}

// GetScreenHeight - Get current screen height
func GetScreenHeight() int32 {
	ret := C.GetScreenHeight()
	v := (int32)(ret)
	return v
}

// ClearBackground - Sets Background Color
func ClearBackground(color Color) {
	ccolor := color.cptr()
	C.ClearBackground(*ccolor)
}

// BeginDrawing - Setup drawing canvas to start drawing
func BeginDrawing() {
	C.BeginDrawing()
}

// EndDrawing - End canvas drawing and Swap Buffers (Double Buffering)
func EndDrawing() {
	C.EndDrawing()
}

// Begin2dMode - Initialize 2D mode with custom camera
func Begin2dMode(camera Camera2D) {
	ccamera := camera.cptr()
	C.Begin2dMode(*ccamera)
}

// End2dMode - Ends 2D mode custom camera usage
func End2dMode() {
	C.End2dMode()
}

// Begin3dMode - Initializes 3D mode for drawing (Camera setup)
func Begin3dMode(camera Camera) {
	ccamera := camera.cptr()
	C.Begin3dMode(*ccamera)
}

// End3dMode - Ends 3D mode and returns to default 2D orthographic mode
func End3dMode() {
	C.End3dMode()
}

// BeginTextureMode - Initializes render texture for drawing
func BeginTextureMode(target RenderTexture2D) {
	ctarget := target.cptr()
	C.BeginTextureMode(*ctarget)
}

// EndTextureMode - Ends drawing to render texture
func EndTextureMode() {
	C.EndTextureMode()
}

// GetMouseRay - Returns a ray trace from mouse position
func GetMouseRay(mousePosition Vector2, camera Camera) Ray {
	cmousePosition := mousePosition.cptr()
	ccamera := camera.cptr()
	ret := C.GetMouseRay(*cmousePosition, *ccamera)
	v := newRayFromPointer(unsafe.Pointer(&ret))
	return v
}

// GetWorldToScreen - Returns the screen space position from a 3d world space position
func GetWorldToScreen(position Vector3, camera Camera) Vector2 {
	cposition := position.cptr()
	ccamera := camera.cptr()
	ret := C.GetWorldToScreen(*cposition, *ccamera)
	v := newVector2FromPointer(unsafe.Pointer(&ret))
	return v
}

// GetCameraMatrix - Returns camera transform matrix (view matrix)
func GetCameraMatrix(camera Camera) Matrix {
	ccamera := camera.cptr()
	ret := C.GetCameraMatrix(*ccamera)
	v := newMatrixFromPointer(unsafe.Pointer(&ret))
	return v
}

// SetTargetFPS - Set target FPS (maximum)
func SetTargetFPS(fps int32) {
	cfps := (C.int)(fps)
	C.SetTargetFPS(cfps)
}

// GetFPS - Returns current FPS
func GetFPS() float32 {
	ret := C.GetFPS()
	v := (float32)(ret)
	return v
}

// GetFrameTime - Returns time in seconds for one frame
func GetFrameTime() float32 {
	ret := C.GetFrameTime()
	v := (float32)(ret)
	return v
}

// GetColor - Returns a Color struct from hexadecimal value
func GetColor(hexValue int32) Color {
	chexValue := (C.int)(hexValue)
	ret := C.GetColor(chexValue)
	v := newColorFromPointer(unsafe.Pointer(&ret))
	return v
}

// GetHexValue - Returns hexadecimal value for a Color
func GetHexValue(color Color) int32 {
	ccolor := color.cptr()
	ret := C.GetHexValue(*ccolor)
	v := (int32)(ret)
	return v
}

// ColorToFloat - Converts Color to float32 slice and normalizes
func ColorToFloat(color Color) []float32 {
	data := make([]float32, 0)
	data[0] = float32(color.R) / 255
	data[0] = float32(color.G) / 255
	data[0] = float32(color.B) / 255
	data[0] = float32(color.A) / 255

	return data
}

// Vector3ToFloat - Converts Vector3 to float32 slice
func Vector3ToFloat(vec Vector3) []float32 {
	data := make([]float32, 0)
	data[0] = vec.X
	data[1] = vec.Y
	data[2] = vec.Z

	return data
}

// MatrixToFloat - Converts Matrix to float32 slice
func MatrixToFloat(mat Matrix) []float32 {
	data := make([]float32, 0)

	data[0] = mat.M0
	data[1] = mat.M4
	data[2] = mat.M8
	data[3] = mat.M12
	data[4] = mat.M1
	data[5] = mat.M5
	data[6] = mat.M9
	data[7] = mat.M13
	data[8] = mat.M2
	data[9] = mat.M6
	data[10] = mat.M10
	data[11] = mat.M14
	data[12] = mat.M3
	data[13] = mat.M7
	data[14] = mat.M11
	data[15] = mat.M15

	return data
}

// GetRandomValue - Returns a random value between min and max (both included)
func GetRandomValue(min, max int32) int32 {
	cmin := (C.int)(min)
	cmax := (C.int)(max)
	ret := C.GetRandomValue(cmin, cmax)
	v := (int32)(ret)
	return v
}

// Fade - Color fade-in or fade-out, alpha goes from 0.0f to 1.0f
func Fade(color Color, alpha float32) Color {
	ccolor := color.cptr()
	calpha := (C.float)(alpha)
	ret := C.Fade(*ccolor, calpha)
	v := newColorFromPointer(unsafe.Pointer(&ret))
	return v
}

// ShowLogo - Activates raylib logo at startup (can be done with flags)
func ShowLogo() {
	C.ShowLogo()
}

// SetConfigFlags - Setup some window configuration flags
func SetConfigFlags(flags byte) {
	cflags := (C.char)(flags)
	C.SetConfigFlags(cflags)
}

// TakeScreenshot - Takes a screenshot of current screen (saved a .png)
func TakeScreenshot(name string) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	C.TakeScreenshot(cname)
}

// StorageSaveValue - Storage save integer value (to defined position)
func StorageSaveValue(position, value int32) {
	cposition := (C.int)(position)
	cvalue := (C.int)(value)
	C.StorageSaveValue(cposition, cvalue)
}

// StorageLoadValue - Storage load integer value (from defined position)
func StorageLoadValue(position int32) int32 {
	cposition := (C.int)(position)
	ret := C.StorageLoadValue(cposition)
	v := (int32)(ret)
	return v
}

// IsKeyPressed - Detect if a key has been pressed once
func IsKeyPressed(key int32) bool {
	ckey := (C.int)(key)
	ret := C.IsKeyPressed(ckey)
	v := bool(int(ret) == 1)
	return v
}

// IsKeyDown - Detect if a key is being pressed
func IsKeyDown(key int32) bool {
	ckey := (C.int)(key)
	ret := C.IsKeyDown(ckey)
	v := bool(int(ret) == 1)
	return v
}

// IsKeyReleased - Detect if a key has been released once
func IsKeyReleased(key int32) bool {
	ckey := (C.int)(key)
	ret := C.IsKeyReleased(ckey)
	v := bool(int(ret) == 1)
	return v
}

// IsKeyUp - Detect if a key is NOT being pressed
func IsKeyUp(key int32) bool {
	ckey := (C.int)(key)
	ret := C.IsKeyUp(ckey)
	v := bool(int(ret) == 1)
	return v
}

// GetKeyPressed - Get latest key pressed
func GetKeyPressed() int32 {
	ret := C.GetKeyPressed()
	v := (int32)(ret)
	return v
}

// SetExitKey - Set a custom key to exit program (default is ESC)
func SetExitKey(key int32) {
	ckey := (C.int)(key)
	C.SetExitKey(ckey)
}

// IsGamepadAvailable - Detect if a gamepad is available
func IsGamepadAvailable(gamepad int32) bool {
	cgamepad := (C.int)(gamepad)
	ret := C.IsGamepadAvailable(cgamepad)
	v := bool(int(ret) == 1)
	return v
}

// IsGamepadName - Check gamepad name (if available)
func IsGamepadName(gamepad int32, name string) bool {
	cgamepad := (C.int)(gamepad)
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	ret := C.IsGamepadName(cgamepad, cname)
	v := bool(int(ret) == 1)
	return v
}

// GetGamepadName - Return gamepad internal name id
func GetGamepadName(gamepad int32) string {
	cgamepad := (C.int)(gamepad)
	ret := C.GetGamepadName(cgamepad)
	v := C.GoString(ret)
	return v
}

// IsGamepadButtonPressed - Detect if a gamepad button has been pressed once
func IsGamepadButtonPressed(gamepad, button int32) bool {
	cgamepad := (C.int)(gamepad)
	cbutton := (C.int)(button)
	ret := C.IsGamepadButtonPressed(cgamepad, cbutton)
	v := bool(int(ret) == 1)
	return v
}

// IsGamepadButtonDown - Detect if a gamepad button is being pressed
func IsGamepadButtonDown(gamepad, button int32) bool {
	cgamepad := (C.int)(gamepad)
	cbutton := (C.int)(button)
	ret := C.IsGamepadButtonDown(cgamepad, cbutton)
	v := bool(int(ret) == 1)
	return v
}

// IsGamepadButtonReleased - Detect if a gamepad button has been released once
func IsGamepadButtonReleased(gamepad, button int32) bool {
	cgamepad := (C.int)(gamepad)
	cbutton := (C.int)(button)
	ret := C.IsGamepadButtonReleased(cgamepad, cbutton)
	v := bool(int(ret) == 1)
	return v
}

// IsGamepadButtonUp - Detect if a gamepad button is NOT being pressed
func IsGamepadButtonUp(gamepad, button int32) bool {
	cgamepad := (C.int)(gamepad)
	cbutton := (C.int)(button)
	ret := C.IsGamepadButtonUp(cgamepad, cbutton)
	v := bool(int(ret) == 1)
	return v
}

// GetGamepadButtonPressed - Get the last gamepad button pressed
func GetGamepadButtonPressed() int32 {
	ret := C.GetGamepadButtonPressed()
	v := (int32)(ret)
	return v
}

// GetGamepadAxisCount - Return gamepad axis count for a gamepad
func GetGamepadAxisCount(gamepad int32) int32 {
	cgamepad := (C.int)(gamepad)
	ret := C.GetGamepadAxisCount(cgamepad)
	v := (int32)(ret)
	return v
}

// GetGamepadAxisMovement - Return axis movement value for a gamepad axis
func GetGamepadAxisMovement(gamepad, axis int32) float32 {
	cgamepad := (C.int)(gamepad)
	caxis := (C.int)(axis)
	ret := C.GetGamepadAxisMovement(cgamepad, caxis)
	v := (float32)(ret)
	return v
}

// IsMouseButtonPressed - Detect if a mouse button has been pressed once
func IsMouseButtonPressed(button int32) bool {
	cbutton := (C.int)(button)
	ret := C.IsMouseButtonPressed(cbutton)
	v := bool(int(ret) == 1)
	return v
}

// IsMouseButtonDown - Detect if a mouse button is being pressed
func IsMouseButtonDown(button int32) bool {
	cbutton := (C.int)(button)
	ret := C.IsMouseButtonDown(cbutton)
	v := bool(int(ret) == 1)
	return v
}

// IsMouseButtonReleased - Detect if a mouse button has been released once
func IsMouseButtonReleased(button int32) bool {
	cbutton := (C.int)(button)
	ret := C.IsMouseButtonReleased(cbutton)
	v := bool(int(ret) == 1)
	return v
}

// IsMouseButtonUp - Detect if a mouse button is NOT being pressed
func IsMouseButtonUp(button int32) bool {
	cbutton := (C.int)(button)
	ret := C.IsMouseButtonUp(cbutton)
	v := bool(int(ret) == 1)
	return v
}

// GetMouseX - Returns mouse position X
func GetMouseX() int32 {
	ret := C.GetMouseX()
	v := (int32)(ret)
	return v
}

// GetMouseY - Returns mouse position Y
func GetMouseY() int32 {
	ret := C.GetMouseY()
	v := (int32)(ret)
	return v
}

// GetMousePosition - Returns mouse position XY
func GetMousePosition() Vector2 {
	ret := C.GetMousePosition()
	v := newVector2FromPointer(unsafe.Pointer(&ret))
	return v
}

// SetMousePosition - Set mouse position XY
func SetMousePosition(position Vector2) {
	cposition := position.cptr()
	C.SetMousePosition(*cposition)
}

// GetMouseWheelMove - Returns mouse wheel movement Y
func GetMouseWheelMove() int32 {
	ret := C.GetMouseWheelMove()
	v := (int32)(ret)
	return v
}

// GetTouchX - Returns touch position X for touch point 0 (relative to screen size)
func GetTouchX() int32 {
	ret := C.GetTouchX()
	v := (int32)(ret)
	return v
}

// GetTouchY - Returns touch position Y for touch point 0 (relative to screen size)
func GetTouchY() int32 {
	ret := C.GetTouchY()
	v := (int32)(ret)
	return v
}

// GetTouchPosition - Returns touch position XY for a touch point index (relative to screen size)
func GetTouchPosition(index int32) Vector2 {
	cindex := (C.int)(index)
	ret := C.GetTouchPosition(cindex)
	v := newVector2FromPointer(unsafe.Pointer(&ret))
	return v
}
