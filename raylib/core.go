package raylib

/*
#include "raylib.h"
#include <stdlib.h>
*/
import "C"

import (
	"unsafe"
)

// cptr returns C pointer
func (v *Vector2) cptr() *C.Vector2 {
	return (*C.Vector2)(unsafe.Pointer(v))
}

// cptr returns C pointer
func (v *Vector3) cptr() *C.Vector3 {
	return (*C.Vector3)(unsafe.Pointer(v))
}

// cptr returns C pointer
func (m *Matrix) cptr() *C.Matrix {
	return (*C.Matrix)(unsafe.Pointer(m))
}

// cptr returns C pointer
func (color *Color) cptr() *C.Color {
	return (*C.Color)(unsafe.Pointer(color))
}

// cptr returns C pointer
func (r *Rectangle) cptr() *C.Rectangle {
	return (*C.Rectangle)(unsafe.Pointer(r))
}

// cptr returns C pointer
func (c *Camera) cptr() *C.Camera {
	return (*C.Camera)(unsafe.Pointer(c))
}

// cptr returns C pointer
func (c *Camera2D) cptr() *C.Camera2D {
	return (*C.Camera2D)(unsafe.Pointer(c))
}

// cptr returns C pointer
func (b *BoundingBox) cptr() *C.BoundingBox {
	return (*C.BoundingBox)(unsafe.Pointer(b))
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

// GetTime - Return time in seconds
func GetTime() float32 {
	ret := C.GetTime()
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

// ColorToInt - Returns hexadecimal value for a Color
func ColorToInt(color Color) int32 {
	ccolor := color.cptr()
	ret := C.ColorToInt(*ccolor)
	v := (int32)(ret)
	return v
}

// ColorToHSV - Returns HSV values for a Color
// NOTE: Hue is returned as degrees [0..360]
func ColorToHSV(color Color) Vector3 {
	ccolor := color.cptr()
	ret := C.ColorToHSV(*ccolor)
	v := newVector3FromPointer(unsafe.Pointer(&ret))
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
	cflags := (C.uchar)(flags)
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
