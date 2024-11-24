package rl

/*
#include "raylib.h"
#include <stdlib.h>
*/
import "C"

import (
	"image/color"
	"unsafe"
)

// newVector2FromPointer - Returns new Vector2 from pointer
func newVector2FromPointer(ptr unsafe.Pointer) Vector2 {
	return *(*Vector2)(ptr)
}

// cptr returns C pointer
func (v *Vector2) cptr() *C.Vector2 {
	return (*C.Vector2)(unsafe.Pointer(v))
}

// newVector3FromPointer - Returns new Vector3 from pointer
func newVector3FromPointer(ptr unsafe.Pointer) Vector3 {
	return *(*Vector3)(ptr)
}

// cptr returns C pointer
func (v *Vector3) cptr() *C.Vector3 {
	return (*C.Vector3)(unsafe.Pointer(v))
}

// newVector4FromPointer - Returns new Vector4 from pointer
func newVector4FromPointer(ptr unsafe.Pointer) Vector4 {
	return *(*Vector4)(ptr)
}

// cptr returns C pointer
func (v *Vector4) cptr() *C.Vector4 {
	return (*C.Vector4)(unsafe.Pointer(v))
}

// newMatrixFromPointer - Returns new Matrix from pointer
func newMatrixFromPointer(ptr unsafe.Pointer) Matrix {
	return *(*Matrix)(ptr)
}

// cptr returns C pointer
func (m *Matrix) cptr() *C.Matrix {
	return (*C.Matrix)(unsafe.Pointer(m))
}

// newColorFromPointer - Returns new Color from pointer
func newColorFromPointer(ptr unsafe.Pointer) color.RGBA {
	return *(*color.RGBA)(ptr)
}

// colorCptr returns color C pointer
func colorCptr(col color.RGBA) *C.Color {
	return (*C.Color)(unsafe.Pointer(&col))
}

// newRectangleFromPointer - Returns new Rectangle from pointer
func newRectangleFromPointer(ptr unsafe.Pointer) Rectangle {
	return *(*Rectangle)(ptr)
}

// cptr returns C pointer
func (r *Rectangle) cptr() *C.Rectangle {
	return (*C.Rectangle)(unsafe.Pointer(r))
}

// newCamera3DFromPointer - Returns new Camera3D from pointer
func newCamera3DFromPointer(ptr unsafe.Pointer) Camera3D {
	return *(*Camera3D)(ptr)
}

// cptr returns C pointer
func (c *Camera) cptr() *C.Camera {
	return (*C.Camera)(unsafe.Pointer(c))
}

// newCamera2DFromPointer - Returns new Camera2D from pointer
func newCamera2DFromPointer(ptr unsafe.Pointer) Camera2D {
	return *(*Camera2D)(ptr)
}

// cptr returns C pointer
func (c *Camera2D) cptr() *C.Camera2D {
	return (*C.Camera2D)(unsafe.Pointer(c))
}

// newBoundingBoxFromPointer - Returns new BoundingBox from pointer
func newBoundingBoxFromPointer(ptr unsafe.Pointer) BoundingBox {
	return *(*BoundingBox)(ptr)
}

// cptr returns C pointer
func (b *BoundingBox) cptr() *C.BoundingBox {
	return (*C.BoundingBox)(unsafe.Pointer(b))
}

// newShaderFromPointer - Returns new Shader from pointer
func newShaderFromPointer(ptr unsafe.Pointer) Shader {
	return *(*Shader)(ptr)
}

// cptr returns C pointer
func (s *Shader) cptr() *C.Shader {
	return (*C.Shader)(unsafe.Pointer(s))
}

// newAutomationEventFromPointer - Returns new AutomationEvent from pointer
func newAutomationEventFromPointer(ptr unsafe.Pointer) AutomationEvent {
	return *(*AutomationEvent)(ptr)
}

// cptr returns C pointer
func (a *AutomationEvent) cptr() *C.AutomationEvent {
	return (*C.AutomationEvent)(unsafe.Pointer(a))
}

// newAutomationEventListFromPointer - Returns new AutomationEventList from pointer
func newAutomationEventListFromPointer(ptr unsafe.Pointer) AutomationEventList {
	return *(*AutomationEventList)(ptr)
}

// cptr returns C pointer
func (a *AutomationEventList) cptr() *C.AutomationEventList {
	return (*C.AutomationEventList)(unsafe.Pointer(a))
}

// WindowShouldClose - Check if KeyEscape pressed or Close icon pressed
func WindowShouldClose() bool {
	ret := C.WindowShouldClose()
	v := bool(ret)
	return v
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

// IsWindowFullscreen - Check if window is currently fullscreen
func IsWindowFullscreen() bool {
	ret := C.IsWindowFullscreen()
	v := bool(ret)
	return v
}

// IsWindowHidden - Check if window is currently hidden
func IsWindowHidden() bool {
	ret := C.IsWindowHidden()
	v := bool(ret)
	return v
}

// IsWindowMinimized - Check if window is currently minimized
func IsWindowMinimized() bool {
	ret := C.IsWindowMinimized()
	v := bool(ret)
	return v
}

// IsWindowMaximized - Check if window is currently maximized
func IsWindowMaximized() bool {
	ret := C.IsWindowMaximized()
	v := bool(ret)
	return v
}

// IsWindowFocused - Check if window is currently focused
func IsWindowFocused() bool {
	ret := C.IsWindowFocused()
	v := bool(ret)
	return v
}

// IsWindowResized - Check if window has been resized
func IsWindowResized() bool {
	ret := C.IsWindowResized()
	v := bool(ret)
	return v
}

// IsWindowState - Check if one specific window flag is enabled
func IsWindowState(flag uint32) bool {
	cflag := (C.uint)(flag)
	ret := C.IsWindowState(cflag)
	v := bool(ret)
	return v
}

// SetWindowState - Set window configuration state using flags
func SetWindowState(flags uint32) {
	cflags := (C.uint)(flags)
	C.SetWindowState(cflags)
}

// ClearWindowState - Clear window configuration state flags
func ClearWindowState(flags uint32) {
	cflags := (C.uint)(flags)
	C.ClearWindowState(cflags)
}

// ToggleFullscreen - Fullscreen toggle (only PLATFORM_DESKTOP)
func ToggleFullscreen() {
	C.ToggleFullscreen()
}

// ToggleBorderlessWindowed - Borderless fullscreen toggle (only PLATFORM_DESKTOP)
func ToggleBorderlessWindowed() {
	C.ToggleBorderlessWindowed()
}

// MaximizeWindow - Set window state: maximized, if resizable
func MaximizeWindow() {
	C.MaximizeWindow()
}

// MinimizeWindow - Set window state: minimized, if resizable
func MinimizeWindow() {
	C.MinimizeWindow()
}

// RestoreWindow - Set window state: not minimized/maximized
func RestoreWindow() {
	C.RestoreWindow()
}

// SetWindowIcon - Set icon for window (single image, RGBA 32bit, only PLATFORM_DESKTOP)
func SetWindowIcon(image Image) {
	cimage := image.cptr()
	C.SetWindowIcon(*cimage)
}

// SetWindowIcons - Set icon for window (multiple images, RGBA 32bit, only PLATFORM_DESKTOP)
func SetWindowIcons(images []Image, count int32) {
	cimages := (&images[0]).cptr()
	cimagesCount := C.int(count)
	C.SetWindowIcons(cimages, cimagesCount)
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

// SetWindowMaxSize - Set window maximum dimensions (for FLAG_WINDOW_RESIZABLE)
func SetWindowMaxSize(w, h int) {
	cw := (C.int)(w)
	ch := (C.int)(h)
	C.SetWindowMaxSize(cw, ch)
}

// SetWindowSize - Set window dimensions
func SetWindowSize(w, h int) {
	cw := (C.int)(w)
	ch := (C.int)(h)
	C.SetWindowSize(cw, ch)
}

// SetWindowOpacity - Set window opacity [0.0f..1.0f] (only PLATFORM_DESKTOP)
func SetWindowOpacity(opacity float32) {
	copacity := (C.float)(opacity)
	C.SetWindowOpacity(copacity)
}

// GetWindowHandle - Get native window handle
func GetWindowHandle() unsafe.Pointer {
	v := C.GetWindowHandle()
	return v
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

// GetRenderWidth - Get current render width (it considers HiDPI)
func GetRenderWidth() int {
	ret := C.GetRenderWidth()
	v := (int)(ret)
	return v
}

// GetRenderHeight - Get current render height (it considers HiDPI)
func GetRenderHeight() int {
	ret := C.GetRenderHeight()
	v := (int)(ret)
	return v
}

// GetMonitorCount - Get number of connected monitors
func GetMonitorCount() int {
	ret := C.GetMonitorCount()
	v := (int)(ret)
	return v
}

// GetCurrentMonitor - Get current monitor where window is placed
func GetCurrentMonitor() int {
	ret := C.GetCurrentMonitor()
	v := (int)(ret)
	return v
}

// GetMonitorPosition - Get specified monitor position
func GetMonitorPosition(monitor int) Vector2 {
	cmonitor := (C.int)(monitor)
	ret := C.GetMonitorPosition(cmonitor)
	v := newVector2FromPointer(unsafe.Pointer(&ret))
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

// GetMonitorRefreshRate - Get specified monitor refresh rate
func GetMonitorRefreshRate(monitor int) int {
	cmonitor := (C.int)(monitor)
	ret := C.GetMonitorRefreshRate(cmonitor)
	v := (int)(ret)
	return v
}

// GetWindowPosition - Get window position XY on monitor
func GetWindowPosition() Vector2 {
	ret := C.GetWindowPosition()
	v := newVector2FromPointer(unsafe.Pointer(&ret))
	return v
}

// GetWindowScaleDPI - Get window scale DPI factor
func GetWindowScaleDPI() Vector2 {
	ret := C.GetWindowScaleDPI()
	v := newVector2FromPointer(unsafe.Pointer(&ret))
	return v
}

// GetMonitorName - Get the human-readable, UTF-8 encoded name of the primary monitor
func GetMonitorName(monitor int) string {
	cmonitor := (C.int)(monitor)
	ret := C.GetMonitorName(cmonitor)
	v := C.GoString(ret)
	return v
}

// SetClipboardText - Set clipboard text content
func SetClipboardText(data string) {
	cdata := C.CString(data)
	defer C.free(unsafe.Pointer(cdata))
	C.SetClipboardText(cdata)
}

// GetClipboardText - Get clipboard text content
func GetClipboardText() string {
	ret := C.GetClipboardText()
	v := C.GoString(ret)
	return v
}

// GetClipboardImage - Get clipboard image content
//
// Only works with SDL3 backend or Windows with GLFW/RGFW
func GetClipboardImage() Image {
	ret := C.GetClipboardImage()
	v := newImageFromPointer(unsafe.Pointer(&ret))
	return *v
}

// EnableEventWaiting - Enable waiting for events on EndDrawing(), no automatic event polling
func EnableEventWaiting() {
	C.EnableEventWaiting()
}

// DisableEventWaiting - Disable waiting for events on EndDrawing(), automatic events polling
func DisableEventWaiting() {
	C.DisableEventWaiting()
}

// ClearBackground - Sets Background Color
func ClearBackground(col color.RGBA) {
	ccolor := colorCptr(col)
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

// LoadShader - Load a custom shader and bind default locations
func LoadShader(vsFileName string, fsFileName string) Shader {
	cvsFileName := C.CString(vsFileName)
	defer C.free(unsafe.Pointer(cvsFileName))

	cfsFileName := C.CString(fsFileName)
	defer C.free(unsafe.Pointer(cfsFileName))

	if vsFileName == "" {
		cvsFileName = nil
	}

	if fsFileName == "" {
		cfsFileName = nil
	}

	ret := C.LoadShader(cvsFileName, cfsFileName)
	v := newShaderFromPointer(unsafe.Pointer(&ret))

	return v
}

// LoadShaderFromMemory - Load shader from code strings and bind default locations
func LoadShaderFromMemory(vsCode string, fsCode string) Shader {
	cvsCode := C.CString(vsCode)
	defer C.free(unsafe.Pointer(cvsCode))

	cfsCode := C.CString(fsCode)
	defer C.free(unsafe.Pointer(cfsCode))

	if vsCode == "" {
		cvsCode = nil
	}

	if fsCode == "" {
		cfsCode = nil
	}

	ret := C.LoadShaderFromMemory(cvsCode, cfsCode)
	v := newShaderFromPointer(unsafe.Pointer(&ret))

	return v
}

// IsShaderValid - Check if a shader is valid (loaded on GPU)
func IsShaderValid(shader Shader) bool {
	cshader := shader.cptr()
	ret := C.IsShaderValid(*cshader)
	v := bool(ret)
	return v
}

// GetShaderLocation - Get shader uniform location
func GetShaderLocation(shader Shader, uniformName string) int32 {
	cshader := shader.cptr()
	cuniformName := C.CString(uniformName)
	defer C.free(unsafe.Pointer(cuniformName))

	ret := C.GetShaderLocation(*cshader, cuniformName)
	v := (int32)(ret)
	return v
}

// GetShaderLocationAttrib - Get shader attribute location
func GetShaderLocationAttrib(shader Shader, attribName string) int32 {
	cshader := shader.cptr()
	cuniformName := C.CString(attribName)
	defer C.free(unsafe.Pointer(cuniformName))

	ret := C.GetShaderLocationAttrib(*cshader, cuniformName)
	v := (int32)(ret)
	return v
}

// SetShaderValue - Set shader uniform value (float)
func SetShaderValue(shader Shader, locIndex int32, value []float32, uniformType ShaderUniformDataType) {
	cshader := shader.cptr()
	clocIndex := (C.int)(locIndex)
	cvalue := (*C.float)(unsafe.Pointer(&value[0]))
	cuniformType := (C.int)(uniformType)
	C.SetShaderValue(*cshader, clocIndex, unsafe.Pointer(cvalue), cuniformType)
}

// SetShaderValueV - Set shader uniform value (float)
func SetShaderValueV(shader Shader, locIndex int32, value []float32, uniformType ShaderUniformDataType, count int32) {
	cshader := shader.cptr()
	clocIndex := (C.int)(locIndex)
	cvalue := (*C.float)(unsafe.Pointer(&value[0]))
	cuniformType := (C.int)(uniformType)
	ccount := (C.int)(count)
	C.SetShaderValueV(*cshader, clocIndex, unsafe.Pointer(cvalue), cuniformType, ccount)
}

// SetShaderValueMatrix - Set shader uniform value (matrix 4x4)
func SetShaderValueMatrix(shader Shader, locIndex int32, mat Matrix) {
	cshader := shader.cptr()
	clocIndex := (C.int)(locIndex)
	cmat := mat.cptr()
	C.SetShaderValueMatrix(*cshader, clocIndex, *cmat)
}

// SetShaderValueTexture - Set shader uniform value for texture (sampler2d)
func SetShaderValueTexture(shader Shader, locIndex int32, texture Texture2D) {
	cshader := shader.cptr()
	clocIndex := (C.int)(locIndex)
	ctexture := texture.cptr()
	C.SetShaderValueTexture(*cshader, clocIndex, *ctexture)
}

// UnloadShader - Unload a custom shader from memory
func UnloadShader(shader Shader) {
	cshader := shader.cptr()
	C.UnloadShader(*cshader)
}

// GetMouseRay - Get a ray trace from mouse position
//
// Deprecated: Use [GetScreenToWorldRay] instead.
func GetMouseRay(mousePosition Vector2, camera Camera) Ray {
	return GetScreenToWorldRay(mousePosition, camera)
}

// GetScreenToWorldRay - Get a ray trace from screen position (i.e mouse)
func GetScreenToWorldRay(position Vector2, camera Camera) Ray {
	cposition := position.cptr()
	ccamera := camera.cptr()
	ret := C.GetScreenToWorldRay(*cposition, *ccamera)
	v := newRayFromPointer(unsafe.Pointer(&ret))
	return v
}

// GetScreenToWorldRayEx - Get a ray trace from screen position (i.e mouse) in a viewport
func GetScreenToWorldRayEx(position Vector2, camera Camera, width, height int32) Ray {
	cposition := position.cptr()
	ccamera := camera.cptr()
	cwidth := (C.int)(width)
	cheight := (C.int)(height)
	ret := C.GetScreenToWorldRayEx(*cposition, *ccamera, cwidth, cheight)
	v := newRayFromPointer(unsafe.Pointer(&ret))
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

// GetWorldToScreen - Returns the screen space position from a 3d world space position
func GetWorldToScreen(position Vector3, camera Camera) Vector2 {
	cposition := position.cptr()
	ccamera := camera.cptr()
	ret := C.GetWorldToScreen(*cposition, *ccamera)
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

// GetWorldToScreenEx - Get size position for a 3d world space position
func GetWorldToScreenEx(position Vector3, camera Camera, width int32, height int32) Vector2 {
	cposition := position.cptr()
	ccamera := camera.cptr()
	cwidth := (C.int)(width)
	cheight := (C.int)(height)
	ret := C.GetWorldToScreenEx(*cposition, *ccamera, cwidth, cheight)
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

// SetTargetFPS - Set target FPS (maximum)
func SetTargetFPS(fps int32) {
	cfps := (C.int)(fps)
	C.SetTargetFPS(cfps)
}

// GetFPS - Returns current FPS
func GetFPS() int32 {
	ret := C.GetFPS()
	v := (int32)(ret)
	return v
}

// GetFrameTime - Returns time in seconds for one frame
func GetFrameTime() float32 {
	ret := C.GetFrameTime()
	v := (float32)(ret)
	return v
}

// GetTime - Return time in seconds
func GetTime() float64 {
	ret := C.GetTime()
	v := (float64)(ret)
	return v
}

// Custom frame control functions
// NOTE: SwapScreenBuffer and PollInputEvents are intended for advanced users that want full control over the frame processing
// By default EndDrawing() does this job: draws everything + SwapScreenBuffer() + manage frame timing + PollInputEvents()
// To avoid that behaviour and control frame processes manually you can either enable in config.h: SUPPORT_CUSTOM_FRAME_CONTROL
// or add CGO_CFLAGS="-DSUPPORT_CUSTOM_FRAME_CONTROL=1" to your build

// SwapScreenBuffer - Swap back buffer to front buffer
func SwapScreenBuffer() {
	C.SwapScreenBuffer()
}

// Register all input events
func PollInputEvents() {
	C.PollInputEvents()
}

// WaitTime - Wait for some time (halt program execution)
func WaitTime(seconds float64) {
	cseconds := (C.double)(seconds)
	C.WaitTime(cseconds)
}

// Fade - Returns color with alpha applied, alpha goes from 0.0f to 1.0f
func Fade(col color.RGBA, alpha float32) color.RGBA {
	ccolor := colorCptr(col)
	calpha := (C.float)(alpha)
	ret := C.Fade(*ccolor, calpha)
	v := newColorFromPointer(unsafe.Pointer(&ret))
	return v
}

// ColorToInt - Get hexadecimal value for a Color (0xRRGGBBAA)
func ColorToInt(col color.RGBA) int32 {
	ccolor := colorCptr(col)
	ret := C.ColorToInt(*ccolor)
	v := (int32)(ret)
	return v
}

// ColorNormalize - Returns color normalized as float [0..1]
func ColorNormalize(col color.RGBA) Vector4 {
	result := Vector4{}
	r, g, b, a := col.R, col.G, col.B, col.A
	result.X = float32(r) / 255
	result.Y = float32(g) / 255
	result.Z = float32(b) / 255
	result.W = float32(a) / 255

	return result
}

// ColorFromNormalized - Returns Color from normalized values [0..1]
func ColorFromNormalized(normalized Vector4) color.RGBA {
	cnormalized := normalized.cptr()
	ret := C.ColorFromNormalized(*cnormalized)
	v := newColorFromPointer(unsafe.Pointer(&ret))
	return v
}

// ColorToHSV - Returns HSV values for a Color, hue [0..360], saturation/value [0..1]
func ColorToHSV(col color.RGBA) Vector3 {
	ccolor := colorCptr(col)
	ret := C.ColorToHSV(*ccolor)
	v := newVector3FromPointer(unsafe.Pointer(&ret))
	return v
}

// ColorFromHSV - Returns a Color from HSV values, hue [0..360], saturation/value [0..1]
func ColorFromHSV(hue, saturation, value float32) color.RGBA {
	chue := (C.float)(hue)
	csaturation := (C.float)(saturation)
	cvalue := (C.float)(value)
	ret := C.ColorFromHSV(chue, csaturation, cvalue)
	v := newColorFromPointer(unsafe.Pointer(&ret))
	return v
}

// ColorTint - Get color multiplied with another color
func ColorTint(col color.RGBA, tint color.RGBA) color.RGBA {
	ccolor := colorCptr(col)
	ctint := colorCptr(tint)
	ret := C.ColorTint(*ccolor, *ctint)
	v := newColorFromPointer(unsafe.Pointer(&ret))
	return v
}

// ColorBrightness - Get color with brightness correction, brightness factor goes from -1.0f to 1.0f
func ColorBrightness(col color.RGBA, factor float32) color.RGBA {
	ccolor := colorCptr(col)
	cfactor := C.float(factor)
	ret := C.ColorBrightness(*ccolor, cfactor)
	v := newColorFromPointer(unsafe.Pointer(&ret))
	return v
}

// ColorContrast - Get color with contrast correction, contrast values between -1.0f and 1.0f
func ColorContrast(col color.RGBA, contrast float32) color.RGBA {
	ccolor := colorCptr(col)
	ccontrast := C.float(contrast)
	ret := C.ColorContrast(*ccolor, ccontrast)
	v := newColorFromPointer(unsafe.Pointer(&ret))
	return v
}

// ColorAlpha - Returns color with alpha applied, alpha goes from 0.0f to 1.0f
func ColorAlpha(col color.RGBA, alpha float32) color.RGBA {
	return Fade(col, alpha)
}

// ColorAlphaBlend - Returns src alpha-blended into dst color with tint
func ColorAlphaBlend(src, dst, tint color.RGBA) color.RGBA {
	csrc := colorCptr(src)
	cdst := colorCptr(dst)
	ctint := colorCptr(tint)
	ret := C.ColorAlphaBlend(*csrc, *cdst, *ctint)
	v := newColorFromPointer(unsafe.Pointer(&ret))
	return v
}

// ColorLerp - Get color lerp interpolation between two colors, factor [0.0f..1.0f]
func ColorLerp(col1, col2 color.RGBA, factor float32) color.RGBA {
	ccol1 := colorCptr(col1)
	ccol2 := colorCptr(col2)
	ret := C.ColorLerp(*ccol1, *ccol2, C.float(factor))
	v := newColorFromPointer(unsafe.Pointer(&ret))
	return v
}

// GetColor - Returns a Color struct from hexadecimal value
func GetColor(hexValue uint) color.RGBA {
	chexValue := (C.uint)(hexValue)
	ret := C.GetColor(chexValue)
	v := newColorFromPointer(unsafe.Pointer(&ret))
	return v
}

// GetPixelDataSize - Get pixel data size in bytes for certain format
func GetPixelDataSize(width, height, format int32) int32 {
	cwidth := (C.int)(width)
	cheight := (C.int)(height)
	cformat := (C.int)(format)
	ret := C.GetPixelDataSize(cwidth, cheight, cformat)
	v := (int32)(ret)
	return v
}

// Vector3ToFloat - Converts Vector3 to float32 slice
func Vector3ToFloat(vec Vector3) []float32 {
	data := make([]float32, 0)
	data[0] = vec.X
	data[1] = vec.Y
	data[2] = vec.Z

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

// OpenURL - Open URL with default system browser (if available)
func OpenURL(url string) {
	curl := C.CString(url)
	defer C.free(unsafe.Pointer(curl))
	C.OpenURL(curl)
}

// SetConfigFlags - Setup some window configuration flags
func SetConfigFlags(flags uint32) {
	cflags := (C.uint)(flags)
	C.SetConfigFlags(cflags)
}

// TakeScreenshot - Takes a screenshot of current screen (saved a .png)
func TakeScreenshot(name string) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	C.TakeScreenshot(cname)
}

// LoadAutomationEventList - Load automation events list from file, NULL for empty list, capacity = MAX_AUTOMATION_EVENTS
func LoadAutomationEventList(fileName string) AutomationEventList {
	cfileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cfileName))

	ret := C.LoadAutomationEventList(cfileName)
	v := newAutomationEventListFromPointer(unsafe.Pointer(&ret))

	return v
}

// UnloadAutomationEventList - Unload automation events list from file
func UnloadAutomationEventList(list *AutomationEventList) {
	C.UnloadAutomationEventList(*list.cptr())
}

// ExportAutomationEventList - Export automation events list as text file
func ExportAutomationEventList(list AutomationEventList, fileName string) bool {
	cfileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cfileName))

	ret := C.ExportAutomationEventList(*list.cptr(), cfileName)
	v := bool(ret)

	return v
}

// SetAutomationEventList - Set automation event list to record to
func SetAutomationEventList(list *AutomationEventList) {
	C.SetAutomationEventList(list.cptr())
}

// SetAutomationEventBaseFrame - Set automation event internal base frame to start recording
func SetAutomationEventBaseFrame(frame int) {
	cframe := (C.int)(frame)
	C.SetAutomationEventBaseFrame(cframe)
}

// StartAutomationEventRecording - Start recording automation events (AutomationEventList must be set)
func StartAutomationEventRecording() {
	C.StartAutomationEventRecording()
}

// StopAutomationEventRecording - Stop recording automation events
func StopAutomationEventRecording() {
	C.StopAutomationEventRecording()
}

// PlayAutomationEvent - Play a recorded automation event
func PlayAutomationEvent(event AutomationEvent) {
	C.PlayAutomationEvent(*event.cptr())
}

// IsKeyPressed - Detect if a key has been pressed once
func IsKeyPressed(key int32) bool {
	ckey := (C.int)(key)
	ret := C.IsKeyPressed(ckey)
	v := bool(ret)
	return v
}

// IsKeyPressedRepeat - Detect if a key has been pressed again (Only PLATFORM_DESKTOP)
func IsKeyPressedRepeat(key int32) bool {
	ckey := (C.int)(key)
	ret := C.IsKeyPressedRepeat(ckey)
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

// GetCharPressed - Get the last char pressed
func GetCharPressed() int32 {
	ret := C.GetCharPressed()
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

// SetGamepadMappings - Set internal gamepad mappings (SDL_GameControllerDB)
func SetGamepadMappings(mappings string) int32 {
	cmappings := C.CString(mappings)
	defer C.free(unsafe.Pointer(cmappings))
	ret := C.SetGamepadMappings(cmappings)
	v := (int32)(ret)
	return v
}

// SetGamepadVibration - Set gamepad vibration for both motors (duration in seconds)
func SetGamepadVibration(gamepad int32, leftMotor, rightMotor, duration float32) {
	C.SetGamepadVibration(C.int(gamepad), C.float(leftMotor), C.float(rightMotor), C.float(duration))
}

// IsMouseButtonPressed - Detect if a mouse button has been pressed once
func IsMouseButtonPressed(button MouseButton) bool {
	cbutton := (C.int)(button)
	ret := C.IsMouseButtonPressed(cbutton)
	v := bool(ret)
	return v
}

// IsMouseButtonDown - Detect if a mouse button is being pressed
func IsMouseButtonDown(button MouseButton) bool {
	cbutton := (C.int)(button)
	ret := C.IsMouseButtonDown(cbutton)
	v := bool(ret)
	return v
}

// IsMouseButtonReleased - Detect if a mouse button has been released once
func IsMouseButtonReleased(button MouseButton) bool {
	cbutton := (C.int)(button)
	ret := C.IsMouseButtonReleased(cbutton)
	v := bool(ret)
	return v
}

// IsMouseButtonUp - Detect if a mouse button is NOT being pressed
func IsMouseButtonUp(button MouseButton) bool {
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

// GetMouseDelta - Get mouse delta between frames
func GetMouseDelta() Vector2 {
	ret := C.GetMouseDelta()
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

// GetMouseWheelMove - Get mouse wheel movement for X or Y, whichever is larger
func GetMouseWheelMove() float32 {
	ret := C.GetMouseWheelMove()
	v := (float32)(ret)
	return v
}

// GetMouseWheelMoveV - Get mouse wheel movement for both X and Y
func GetMouseWheelMoveV() Vector2 {
	ret := C.GetMouseWheelMoveV()
	v := newVector2FromPointer(unsafe.Pointer(&ret))
	return v
}

// SetMouseCursor - Set mouse cursor
func SetMouseCursor(cursor int32) {
	ccursor := (C.int)(cursor)
	C.SetMouseCursor(ccursor)
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

// GetTouchPointId - Get touch point identifier for given index
func GetTouchPointId(index int32) int32 {
	cindex := (C.int)(index)
	ret := C.GetTouchPointId(cindex)
	v := (int32)(ret)
	return v
}

// GetTouchPointCount - Get number of touch points
func GetTouchPointCount() int32 {
	ret := C.GetTouchPointCount()
	v := (int32)(ret)
	return v
}

// BeginVrStereoMode - Begin stereo rendering (requires VR simulator)
func BeginVrStereoMode(config VrStereoConfig) {
	C.BeginVrStereoMode(*(*C.VrStereoConfig)(unsafe.Pointer(&config)))
}

// EndVrStereoMode - End stereo rendering (requires VR simulator)
func EndVrStereoMode() {
	C.EndVrStereoMode()
}

// LoadVrStereoConfig - Load VR stereo config for VR simulator device parameters
func LoadVrStereoConfig(device VrDeviceInfo) VrStereoConfig {
	ret := C.LoadVrStereoConfig(*(*C.VrDeviceInfo)(unsafe.Pointer(&device)))
	return *(*VrStereoConfig)(unsafe.Pointer(&ret))
}

// UnloadVrStereoConfig - Unload VR stereo config
func UnloadVrStereoConfig(config VrStereoConfig) {
	C.UnloadVrStereoConfig(*(*C.VrStereoConfig)(unsafe.Pointer(&config)))
}
