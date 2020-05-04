package rl

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

// IsWindowReady - Check if window has been initialized successfully
func IsWindowReady() bool {
	ret := C.IsWindowReady()
	v := bool(ret)
	return v
}

// WindowShouldClose - Detect if KEY_ESCAPE pressed or Close icon pressed
func WindowShouldClose() bool {
	ret := C.WindowShouldClose()
	v := bool(ret)
	return v
}

// IsWindowMinimized - Detect if window has been minimized (or lost focus)
func IsWindowMinimized() bool {
	ret := C.IsWindowMinimized()
	v := bool(ret)
	return v
}

// IsWindowResized - Check if window has been resized
func IsWindowResized() bool {
	ret := C.IsWindowResized()
	v := bool(ret)
	return v
}

// IsWindowHidden - Check if window is currently hidden
func IsWindowHidden() bool {
	ret := C.IsWindowHidden()
	v := bool(ret)
	return v
}

// ToggleFullscreen - Fullscreen toggle (only PLATFORM_DESKTOP)
func ToggleFullscreen() {
	C.ToggleFullscreen()
}

// UnhideWindow - Show the window
func UnhideWindow() {
	C.UnhideWindow()
}

// HideWindow - Hide the window
func HideWindow() {
	C.HideWindow()
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
func SetWindowPosition(x, y int) {
	cx := (C.int)(x)
	cy := (C.int)(y)
	C.SetWindowPosition(cx, cy)
}

// SetWindowMonitor - Set monitor for the current window (fullscreen mode)
func SetWindowMonitor(monitor int) {
	cmonitor := (C.int)(monitor)
	C.SetWindowMonitor(cmonitor)
}

// SetWindowMinSize - Set window minimum dimensions (for FLAG_WINDOW_RESIZABLE)
func SetWindowMinSize(w, h int) {
	cw := (C.int)(w)
	ch := (C.int)(h)
	C.SetWindowMinSize(cw, ch)
}

// SetWindowSize - Set window dimensions
func SetWindowSize(w, h int) {
	cw := (C.int)(w)
	ch := (C.int)(h)
	C.SetWindowSize(cw, ch)
}

// GetScreenWidth - Get current screen width
func GetScreenWidth() int {
	ret := C.GetScreenWidth()
	v := (int)(ret)
	return v
}

// GetScreenHeight - Get current screen height
func GetScreenHeight() int {
	ret := C.GetScreenHeight()
	v := (int)(ret)
	return v
}

// GetMonitorCount - Get number of connected monitors
func GetMonitorCount() int {
	ret := C.GetMonitorCount()
	v := (int)(ret)
	return v
}

// GetMonitorWidth - Get primary monitor width
func GetMonitorWidth(monitor int) int {
	cmonitor := (C.int)(monitor)
	ret := C.GetMonitorWidth(cmonitor)
	v := (int)(ret)
	return v
}

// GetMonitorHeight - Get primary monitor height
func GetMonitorHeight(monitor int) int {
	cmonitor := (C.int)(monitor)
	ret := C.GetMonitorHeight(cmonitor)
	v := (int)(ret)
	return v
}

// GetMonitorPhysicalWidth - Get primary monitor physical width in millimetres
func GetMonitorPhysicalWidth(monitor int) int {
	cmonitor := (C.int)(monitor)
	ret := C.GetMonitorPhysicalWidth(cmonitor)
	v := (int)(ret)
	return v
}

// GetMonitorPhysicalHeight - Get primary monitor physical height in millimetres
func GetMonitorPhysicalHeight(monitor int) int {
	cmonitor := (C.int)(monitor)
	ret := C.GetMonitorPhysicalHeight(cmonitor)
	v := (int)(ret)
	return v
}

// GetMonitorName - Get the human-readable, UTF-8 encoded name of the primary monitor
func GetMonitorName(monitor int) string {
	cmonitor := (C.int)(monitor)
	ret := C.GetMonitorName(cmonitor)
	v := C.GoString(ret)
	return v
}

// GetClipboardText - Get clipboard text content
func GetClipboardText() string {
	ret := C.GetClipboardText()
	v := C.GoString(ret)
	return v
}

// SetClipboardText - Set clipboard text content
func SetClipboardText(data string) {
	cdata := C.CString(data)
	defer C.free(unsafe.Pointer(cdata))
	C.SetClipboardText(cdata)
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

// BeginMode2D - Initialize 2D mode with custom camera
func BeginMode2D(camera Camera2D) {
	ccamera := camera.cptr()
	C.BeginMode2D(*ccamera)
}

// EndMode2D - Ends 2D mode custom camera usage
func EndMode2D() {
	C.EndMode2D()
}

// BeginMode3D - Initializes 3D mode for drawing (Camera setup)
func BeginMode3D(camera Camera) {
	ccamera := camera.cptr()
	C.BeginMode3D(*ccamera)
}

// EndMode3D - Ends 3D mode and returns to default 2D orthographic mode
func EndMode3D() {
	C.EndMode3D()
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

// BeginScissorMode - Begins scissor mode (define screen area for following drawing)
func BeginScissorMode(x, y, width, height int32) {
	cx := (C.int)(x)
	cy := (C.int)(y)
	cwidth := (C.int)(width)
	cheight := (C.int)(height)
	C.BeginScissorMode(cx, cy, cwidth, cheight)
}

// EndScissorMode - Ends scissor mode
func EndScissorMode() {
	C.EndScissorMode()
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

// GetWorldToScreen2D - Returns the screen space position for a 2d camera world space position
func GetWorldToScreen2D(position Vector2, camera Camera2D) Vector2 {
	cposition := position.cptr()
	ccamera := camera.cptr()
	ret := C.GetWorldToScreen2D(*cposition, *ccamera)
	v := newVector2FromPointer(unsafe.Pointer(&ret))
	return v
}

// GetScreenToWorld2D - Returns the world space position for a 2d camera screen space position
func GetScreenToWorld2D(position Vector2, camera Camera2D) Vector2 {
	cposition := position.cptr()
	ccamera := camera.cptr()
	ret := C.GetScreenToWorld2D(*cposition, *ccamera)
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

// GetCameraMatrix2D - Returns camera 2d transform matrix
func GetCameraMatrix2D(camera Camera2D) Matrix {
	ccamera := camera.cptr()
	ret := C.GetCameraMatrix2D(*ccamera)
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

// ColorNormalize - Returns color normalized as float [0..1]
func ColorNormalize(color Color) Vector4 {
	result := Vector4{}
	result.X = float32(color.R) / 255
	result.Y = float32(color.G) / 255
	result.Z = float32(color.B) / 255
	result.W = float32(color.A) / 255

	return result
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

// SetConfigFlags - Setup some window configuration flags
func SetConfigFlags(flags byte) {
	cflags := (C.uint)(flags)
	C.SetConfigFlags(cflags)
}

// TakeScreenshot - Takes a screenshot of current screen (saved a .png)
func TakeScreenshot(name string) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	C.TakeScreenshot(cname)
}

// SaveStorageValue - Storage save integer value (to defined position)
func SaveStorageValue(position, value int32) {
	cposition := (C.uint)(position)
	cvalue := (C.int)(value)
	C.SaveStorageValue(cposition, cvalue)
}

// LoadStorageValue - Storage load integer value (from defined position)
func LoadStorageValue(position int32) int32 {
	cposition := (C.uint)(position)
	ret := C.LoadStorageValue(cposition)
	v := (int32)(ret)
	return v
}

// IsKeyPressed - Detect if a key has been pressed once
func IsKeyPressed(key int32) bool {
	ckey := (C.int)(key)
	ret := C.IsKeyPressed(ckey)
	v := bool(ret)
	return v
}

// IsKeyDown - Detect if a key is being pressed
func IsKeyDown(key int32) bool {
	ckey := (C.int)(key)
	ret := C.IsKeyDown(ckey)
	v := bool(ret)
	return v
}

// IsKeyReleased - Detect if a key has been released once
func IsKeyReleased(key int32) bool {
	ckey := (C.int)(key)
	ret := C.IsKeyReleased(ckey)
	v := bool(ret)
	return v
}

// IsKeyUp - Detect if a key is NOT being pressed
func IsKeyUp(key int32) bool {
	ckey := (C.int)(key)
	ret := C.IsKeyUp(ckey)
	v := bool(ret)
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
	v := bool(ret)
	return v
}

// IsGamepadName - Check gamepad name (if available)
func IsGamepadName(gamepad int32, name string) bool {
	cgamepad := (C.int)(gamepad)
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	ret := C.IsGamepadName(cgamepad, cname)
	v := bool(ret)
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
	v := bool(ret)
	return v
}

// IsGamepadButtonDown - Detect if a gamepad button is being pressed
func IsGamepadButtonDown(gamepad, button int32) bool {
	cgamepad := (C.int)(gamepad)
	cbutton := (C.int)(button)
	ret := C.IsGamepadButtonDown(cgamepad, cbutton)
	v := bool(ret)
	return v
}

// IsGamepadButtonReleased - Detect if a gamepad button has been released once
func IsGamepadButtonReleased(gamepad, button int32) bool {
	cgamepad := (C.int)(gamepad)
	cbutton := (C.int)(button)
	ret := C.IsGamepadButtonReleased(cgamepad, cbutton)
	v := bool(ret)
	return v
}

// IsGamepadButtonUp - Detect if a gamepad button is NOT being pressed
func IsGamepadButtonUp(gamepad, button int32) bool {
	cgamepad := (C.int)(gamepad)
	cbutton := (C.int)(button)
	ret := C.IsGamepadButtonUp(cgamepad, cbutton)
	v := bool(ret)
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
	v := bool(ret)
	return v
}

// IsMouseButtonDown - Detect if a mouse button is being pressed
func IsMouseButtonDown(button int32) bool {
	cbutton := (C.int)(button)
	ret := C.IsMouseButtonDown(cbutton)
	v := bool(ret)
	return v
}

// IsMouseButtonReleased - Detect if a mouse button has been released once
func IsMouseButtonReleased(button int32) bool {
	cbutton := (C.int)(button)
	ret := C.IsMouseButtonReleased(cbutton)
	v := bool(ret)
	return v
}

// IsMouseButtonUp - Detect if a mouse button is NOT being pressed
func IsMouseButtonUp(button int32) bool {
	cbutton := (C.int)(button)
	ret := C.IsMouseButtonUp(cbutton)
	v := bool(ret)
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
func SetMousePosition(x, y int) {
	cx := (C.int)(x)
	cy := (C.int)(y)
	C.SetMousePosition(cx, cy)
}

// SetMouseOffset - Set mouse offset
func SetMouseOffset(offsetX, offsetY int) {
	ox := (C.int)(offsetX)
	oy := (C.int)(offsetY)
	C.SetMouseOffset(ox, oy)
}

// SetMouseScale - Set mouse scaling
func SetMouseScale(scaleX, scaleY float32) {
	cscaleX := (C.float)(scaleX)
	cscaleY := (C.float)(scaleY)
	C.SetMouseScale(cscaleX, cscaleY)
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
