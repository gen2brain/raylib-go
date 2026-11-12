//go:build js && wasm

package rl

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"time"
	"unsafe"
)

// InitWindow - Initialize window and OpenGL context
func InitWindow(width int32, height int32, title string) {
	defer scratch.release(scratch.mark())
	cInitWindow(width, height, scratch.str(title))
}

// CloseWindow - Close window and unload OpenGL context
func CloseWindow() {
	cCloseWindow()
}

// WindowShouldClose - Check if application should close ([KeyEscape] pressed or windows close icon clicked)
func WindowShouldClose() bool {
	return cWindowShouldClose() != 0
}

// IsWindowReady - Check if window has been initialized successfully
func IsWindowReady() bool {
	return cIsWindowReady() != 0
}

// IsWindowFullscreen - Check if window is currently fullscreen
func IsWindowFullscreen() bool {
	return cIsWindowFullscreen() != 0
}

// IsWindowHidden - Check if window is currently hidden
func IsWindowHidden() bool {
	return cIsWindowHidden() != 0
}

// IsWindowMinimized - Check if window is currently minimized
func IsWindowMinimized() bool {
	return cIsWindowMinimized() != 0
}

// IsWindowMaximized - Check if window is currently maximized
func IsWindowMaximized() bool {
	return cIsWindowMaximized() != 0
}

// IsWindowFocused - Check if window is currently focused
func IsWindowFocused() bool {
	return cIsWindowFocused() != 0
}

// IsWindowResized - Check if window has been resized last frame
func IsWindowResized() bool {
	return cIsWindowResized() != 0
}

// IsWindowState - Check if one specific window flag is enabled
func IsWindowState(flag uint32) bool {
	return cIsWindowState(flag) != 0
}

// SetWindowState - Set window configuration state using flags
func SetWindowState(flags uint32) {
	cSetWindowState(flags)
}

// ClearWindowState - Clear window configuration state flags
func ClearWindowState(flags uint32) {
	cClearWindowState(flags)
}

// ToggleFullscreen - Toggle window state: fullscreen/windowed, resizes monitor to match window resolution
func ToggleFullscreen() {
	cToggleFullscreen()
}

// ToggleBorderlessWindowed - Toggle window state: borderless windowed, resizes window to match monitor resolution
func ToggleBorderlessWindowed() {
	cToggleBorderlessWindowed()
}

// MaximizeWindow - Set window state: maximized, if resizable
func MaximizeWindow() {
	cMaximizeWindow()
}

// MinimizeWindow - Set window state: minimized, if resizable
func MinimizeWindow() {
	cMinimizeWindow()
}

// RestoreWindow - Restore window from being minimized/maximized
func RestoreWindow() {
	cRestoreWindow()
}

// SetWindowIcon - Set icon for window (single image, RGBA 32bit, only PLATFORM_DESKTOP)
func SetWindowIcon(image Image) {
	defer scratch.release(scratch.mark())
	p0 := imageIn(&image)
	cSetWindowIcon(p0)
	imageDone(p0)
}

// SetWindowIcons - Set icon for window (multiple images, RGBA 32bit, only PLATFORM_DESKTOP)
func SetWindowIcons(images []Image, count int32) {
	defer scratch.release(scratch.mark())
	cs := make([]cImage, len(images))
	for i := range images {
		cs[i] = imageToC(&images[i])
	}
	cSetWindowIcons(scratchSlice(cs), count)
	for _, c := range cs {
		cFree(c.Data)
	}
}

// SetWindowTitle - Set title for window
func SetWindowTitle(title string) {
	defer scratch.release(scratch.mark())
	cSetWindowTitle(scratch.str(title))
}

// SetWindowPosition - Set window position on screen
func SetWindowPosition(x, y int) {
	cSetWindowPosition(int32(x), int32(y))
}

// SetWindowMonitor - Set monitor for the current window
func SetWindowMonitor(monitor int) {
	cSetWindowMonitor(int32(monitor))
}

// SetWindowMinSize - Set window minimum dimensions (for FLAG_WINDOW_RESIZABLE)
func SetWindowMinSize(width int, height int) {
	cSetWindowMinSize(int32(width), int32(height))
}

// SetWindowMaxSize - Set window maximum dimensions (for FLAG_WINDOW_RESIZABLE)
func SetWindowMaxSize(width int, height int) {
	cSetWindowMaxSize(int32(width), int32(height))
}

// SetWindowSize - Set window dimensions
func SetWindowSize(width int, height int) {
	cSetWindowSize(int32(width), int32(height))
}

// SetWindowOpacity - Set window opacity [0.0f..1.0f]
func SetWindowOpacity(opacity float32) {
	cSetWindowOpacity(opacity)
}

// SetWindowFocused - Set window focused
func SetWindowFocused() {
	cSetWindowFocused()
}

// GetWindowHandle - Get native window handle
func GetWindowHandle() unsafe.Pointer {
	return nil
}

// GetScreenWidth - Get current screen width
func GetScreenWidth() int {
	return int(cGetScreenWidth())
}

// GetScreenHeight - Get current screen height
func GetScreenHeight() int {
	return int(cGetScreenHeight())
}

// GetRenderWidth - Get current render width (it considers HiDPI)
func GetRenderWidth() int {
	return int(cGetRenderWidth())
}

// GetRenderHeight - Get current render height (it considers HiDPI)
func GetRenderHeight() int {
	return int(cGetRenderHeight())
}

// GetMonitorCount - Get number of connected monitors
func GetMonitorCount() int {
	return int(cGetMonitorCount())
}

// GetCurrentMonitor - Get current monitor where window is placed
func GetCurrentMonitor() int {
	return int(cGetCurrentMonitor())
}

// GetMonitorPosition - Get specified monitor position
func GetMonitorPosition(monitor int) Vector2 {
	ret := cGetMonitorPosition(int32(monitor))
	return *(*Vector2)(unsafe.Pointer(&ret))
}

// GetMonitorWidth - Get specified monitor width (current video mode used by monitor)
func GetMonitorWidth(monitor int) int {
	return int(cGetMonitorWidth(int32(monitor)))
}

// GetMonitorHeight - Get specified monitor height (current video mode used by monitor)
func GetMonitorHeight(monitor int) int {
	return int(cGetMonitorHeight(int32(monitor)))
}

// GetMonitorPhysicalWidth - Get specified monitor physical width in millimetres
func GetMonitorPhysicalWidth(monitor int) int {
	return int(cGetMonitorPhysicalWidth(int32(monitor)))
}

// GetMonitorPhysicalHeight - Get specified monitor physical height in millimetres
func GetMonitorPhysicalHeight(monitor int) int {
	return int(cGetMonitorPhysicalHeight(int32(monitor)))
}

// GetMonitorRefreshRate - Get specified monitor refresh rate
func GetMonitorRefreshRate(monitor int) int {
	return int(cGetMonitorRefreshRate(int32(monitor)))
}

// GetWindowPosition - Get window position XY on monitor
func GetWindowPosition() Vector2 {
	ret := cGetWindowPosition()
	return *(*Vector2)(unsafe.Pointer(&ret))
}

// GetWindowScaleDPI - Get window scale DPI factor
func GetWindowScaleDPI() Vector2 {
	ret := cGetWindowScaleDPI()
	return *(*Vector2)(unsafe.Pointer(&ret))
}

// GetMonitorName - Get the human-readable, UTF-8 encoded name of the specified monitor
func GetMonitorName(monitor int) string {
	return goString(cGetMonitorName(int32(monitor)))
}

// SetClipboardText - Set clipboard text content
func SetClipboardText(text string) {
	defer scratch.release(scratch.mark())
	cSetClipboardText(scratch.str(text))
}

// GetClipboardText - Get clipboard text content
func GetClipboardText() string {
	return goString(cGetClipboardText())
}

// GetClipboardImage - Get clipboard image content
//
// NOTE: This function doesn't work on all platforms/backends.
func GetClipboardImage() Image {
	cGetClipboardImage()
	return *imageRet()
}

// EnableEventWaiting - Enable waiting for events on EndDrawing(), no automatic event polling
func EnableEventWaiting() {
	cEnableEventWaiting()
}

// DisableEventWaiting - Disable waiting for events on EndDrawing(), automatic events polling
func DisableEventWaiting() {
	cDisableEventWaiting()
}

// ShowCursor - Shows cursor
func ShowCursor() {
	cShowCursor()
}

// HideCursor - Hides cursor
func HideCursor() {
	cHideCursor()
}

// IsCursorHidden - Check if cursor is not visible
func IsCursorHidden() bool {
	return cIsCursorHidden() != 0
}

// EnableCursor - Enables cursor (unlock cursor)
func EnableCursor() {
	cEnableCursor()
}

// DisableCursor - Disables cursor (lock cursor)
func DisableCursor() {
	cDisableCursor()
}

// IsCursorOnScreen - Check if cursor is on the screen
func IsCursorOnScreen() bool {
	return cIsCursorOnScreen() != 0
}

// ClearBackground - Set background color (framebuffer clear color)
func ClearBackground(col color.RGBA) {
	cClearBackground(colorToUint32(col))
}

// BeginDrawing - Setup canvas (framebuffer) to start drawing
func BeginDrawing() {
	cBeginDrawing()
}

// EndDrawing - End canvas drawing and swap buffers (double buffering)
func EndDrawing() {
	cEndDrawing()
	syncAutomationList()
	waitFrame()
}

// BeginMode2D - Begin 2D mode with custom camera (2D)
func BeginMode2D(camera Camera2D) {
	cBeginMode2D(camera.Offset.X, camera.Offset.Y, camera.Target.X, camera.Target.Y, camera.Rotation, camera.Zoom)
}

// EndMode2D - Ends 2D mode with custom camera
func EndMode2D() {
	cEndMode2D()
}

// BeginMode3D - Begin 3D mode with custom camera (3D)
func BeginMode3D(camera Camera3D) {
	cBeginMode3D(camera.Position.X, camera.Position.Y, camera.Position.Z, camera.Target.X, camera.Target.Y, camera.Target.Z, camera.Up.X, camera.Up.Y, camera.Up.Z, camera.Fovy, int32(camera.Projection))
}

// EndMode3D - Ends 3D mode and returns to default 2D orthographic mode
func EndMode3D() {
	cEndMode3D()
}

// BeginTextureMode - Begin drawing to render texture
func BeginTextureMode(target RenderTexture2D) {
	cBeginTextureMode(target.ID, target.Texture.ID, target.Texture.Width, target.Texture.Height, target.Texture.Mipmaps, int32(target.Texture.Format), target.Depth.ID, target.Depth.Width, target.Depth.Height, target.Depth.Mipmaps, int32(target.Depth.Format))
}

// EndTextureMode - Ends drawing to render texture
func EndTextureMode() {
	cEndTextureMode()
}

// BeginShaderMode - Begin custom shader drawing
func BeginShaderMode(shader Shader) {
	defer scratch.release(scratch.mark())
	cBeginShaderMode(shaderIn(shader))
}

// EndShaderMode - End custom shader drawing (use default shader)
func EndShaderMode() {
	cEndShaderMode()
}

// BeginBlendMode - Begin blending mode (alpha, additive, multiplied, subtract, custom)
func BeginBlendMode(mode BlendMode) {
	cBeginBlendMode(int32(mode))
}

// EndBlendMode - End blending mode (reset to default: alpha blending)
func EndBlendMode() {
	cEndBlendMode()
}

// BeginScissorMode - Begin scissor mode (define screen area for following drawing)
func BeginScissorMode(x, y, width, height int32) {
	cBeginScissorMode(x, y, width, height)
}

// EndScissorMode - End scissor mode
func EndScissorMode() {
	cEndScissorMode()
}

// BeginVrStereoMode - Begin stereo rendering (requires VR simulator)
func BeginVrStereoMode(config VrStereoConfig) {
	defer scratch.release(scratch.mark())
	cBeginVrStereoMode(scratchValue(&config))
}

// EndVrStereoMode - End stereo rendering (requires VR simulator)
func EndVrStereoMode() {
	cEndVrStereoMode()
}

// LoadVrStereoConfig - Load VR stereo config for VR simulator device parameters
func LoadVrStereoConfig(device VrDeviceInfo) VrStereoConfig {
	cLoadVrStereoConfig(device.HResolution, device.VResolution, device.HScreenSize, device.VScreenSize, device.EyeToScreenDistance, device.LensSeparationDistance, device.InterpupillaryDistance, device.LensDistortionValues[0], device.LensDistortionValues[1], device.LensDistortionValues[2], device.LensDistortionValues[3], device.ChromaAbCorrection[0], device.ChromaAbCorrection[1], device.ChromaAbCorrection[2], device.ChromaAbCorrection[3])
	return readRet[VrStereoConfig]()
}

// UnloadVrStereoConfig - Unload VR stereo config
func UnloadVrStereoConfig(config VrStereoConfig) {
	defer scratch.release(scratch.mark())
	cUnloadVrStereoConfig(scratchValue(&config))
}

// LoadShader - Load shader from files and bind default locations
func LoadShader(vsFileName string, fsFileName string) Shader {
	defer scratch.release(scratch.mark())
	cLoadShader(scratch.strNullable(vsFileName), scratch.strNullable(fsFileName))
	return shaderFromC(readRet[cShader]())
}

// LoadShaderFromMemory - Load shader from code strings and bind default locations
func LoadShaderFromMemory(vsCode string, fsCode string) Shader {
	defer scratch.release(scratch.mark())
	cLoadShaderFromMemory(scratch.strNullable(vsCode), scratch.strNullable(fsCode))
	return shaderFromC(readRet[cShader]())
}

// IsShaderValid - Check if a shader is valid (loaded on GPU)
func IsShaderValid(shader Shader) bool {
	defer scratch.release(scratch.mark())
	return cIsShaderValid(shaderIn(shader)) != 0
}

// GetShaderLocation - Get shader uniform location
func GetShaderLocation(shader Shader, uniformName string) int32 {
	defer scratch.release(scratch.mark())
	return cGetShaderLocation(shaderIn(shader), scratch.str(uniformName))
}

// GetShaderLocationAttrib - Get shader attribute location
func GetShaderLocationAttrib(shader Shader, attribName string) int32 {
	defer scratch.release(scratch.mark())
	return cGetShaderLocationAttrib(shaderIn(shader), scratch.str(attribName))
}

// SetShaderValue - Set shader uniform value (float)
func SetShaderValue(shader Shader, locIndex int32, value []float32, uniformType ShaderUniformDataType) {
	defer scratch.release(scratch.mark())
	cSetShaderValue(shaderIn(shader), locIndex, scratch.bytes(unsafe.Pointer(&value[0]), uint32(len(value))*4), int32(uniformType))
}

// SetShaderValueV - Set shader uniform value (float)
func SetShaderValueV(shader Shader, locIndex int32, value []float32, uniformType ShaderUniformDataType, count int32) {
	defer scratch.release(scratch.mark())
	cSetShaderValueV(shaderIn(shader), locIndex, scratch.bytes(unsafe.Pointer(&value[0]), uint32(len(value))*4), int32(uniformType), count)
}

// SetShaderValueMatrix - Set shader uniform value (matrix 4x4)
func SetShaderValueMatrix(shader Shader, locIndex int32, mat Matrix) {
	defer scratch.release(scratch.mark())
	cSetShaderValueMatrix(shaderIn(shader), locIndex, mat.M0, mat.M4, mat.M8, mat.M12, mat.M1, mat.M5, mat.M9, mat.M13, mat.M2, mat.M6, mat.M10, mat.M14, mat.M3, mat.M7, mat.M11, mat.M15)
}

// SetShaderValueTexture - Set shader uniform value for texture (sampler2d)
func SetShaderValueTexture(shader Shader, locIndex int32, texture Texture2D) {
	defer scratch.release(scratch.mark())
	cSetShaderValueTexture(shaderIn(shader), locIndex, texture.ID, texture.Width, texture.Height, texture.Mipmaps, int32(texture.Format))
}

// UnloadShader - Unload a custom shader from memory
func UnloadShader(shader Shader) {
	defer scratch.release(scratch.mark())
	cUnloadShader(shaderIn(shader))
	if shader.ID != crlGetShaderIdDefault() {
		delete(shaderLocs, shader.ID)
	}
}

// GetScreenToWorldRay - Get a ray trace from screen position (i.e mouse)
func GetScreenToWorldRay(position Vector2, camera Camera) Ray {
	cGetScreenToWorldRay(position.X, position.Y, camera.Position.X, camera.Position.Y, camera.Position.Z, camera.Target.X, camera.Target.Y, camera.Target.Z, camera.Up.X, camera.Up.Y, camera.Up.Z, camera.Fovy, int32(camera.Projection))
	return readRet[Ray]()
}

// GetScreenToWorldRayEx - Get a ray trace from screen position (i.e mouse) in a viewport
func GetScreenToWorldRayEx(position Vector2, camera Camera, width, height int32) Ray {
	cGetScreenToWorldRayEx(position.X, position.Y, camera.Position.X, camera.Position.Y, camera.Position.Z, camera.Target.X, camera.Target.Y, camera.Target.Z, camera.Up.X, camera.Up.Y, camera.Up.Z, camera.Fovy, int32(camera.Projection), width, height)
	return readRet[Ray]()
}

// GetWorldToScreen - Get the screen space position for a 3d world space position
func GetWorldToScreen(position Vector3, camera Camera) Vector2 {
	ret := cGetWorldToScreen(position.X, position.Y, position.Z, camera.Position.X, camera.Position.Y, camera.Position.Z, camera.Target.X, camera.Target.Y, camera.Target.Z, camera.Up.X, camera.Up.Y, camera.Up.Z, camera.Fovy, int32(camera.Projection))
	return *(*Vector2)(unsafe.Pointer(&ret))
}

// GetWorldToScreenEx - Get size position for a 3d world space position
func GetWorldToScreenEx(position Vector3, camera Camera, width int32, height int32) Vector2 {
	ret := cGetWorldToScreenEx(position.X, position.Y, position.Z, camera.Position.X, camera.Position.Y, camera.Position.Z, camera.Target.X, camera.Target.Y, camera.Target.Z, camera.Up.X, camera.Up.Y, camera.Up.Z, camera.Fovy, int32(camera.Projection), width, height)
	return *(*Vector2)(unsafe.Pointer(&ret))
}

// GetWorldToScreen2D - Get the screen space position for a 2d camera world space position
func GetWorldToScreen2D(position Vector2, camera Camera2D) Vector2 {
	ret := cGetWorldToScreen2D(position.X, position.Y, camera.Offset.X, camera.Offset.Y, camera.Target.X, camera.Target.Y, camera.Rotation, camera.Zoom)
	return *(*Vector2)(unsafe.Pointer(&ret))
}

// GetScreenToWorld2D - Get the world space position for a 2d camera screen space position
func GetScreenToWorld2D(position Vector2, camera Camera2D) Vector2 {
	ret := cGetScreenToWorld2D(position.X, position.Y, camera.Offset.X, camera.Offset.Y, camera.Target.X, camera.Target.Y, camera.Rotation, camera.Zoom)
	return *(*Vector2)(unsafe.Pointer(&ret))
}

// GetCameraMatrix - Get camera transform matrix (view matrix)
func GetCameraMatrix(camera Camera) Matrix {
	cGetCameraMatrix(camera.Position.X, camera.Position.Y, camera.Position.Z, camera.Target.X, camera.Target.Y, camera.Target.Z, camera.Up.X, camera.Up.Y, camera.Up.Z, camera.Fovy, int32(camera.Projection))
	return readRet[Matrix]()
}

// GetCameraMatrix2D - Get camera 2d transform matrix
func GetCameraMatrix2D(camera Camera2D) Matrix {
	cGetCameraMatrix2D(camera.Offset.X, camera.Offset.Y, camera.Target.X, camera.Target.Y, camera.Rotation, camera.Zoom)
	return readRet[Matrix]()
}

// SetTargetFPS - Set target FPS (maximum)
func SetTargetFPS(fps int32) {
	targetFPS = fps
}

// GetFrameTime - Get time in seconds for last frame drawn (delta time)
func GetFrameTime() float32 {
	return cGetFrameTime()
}

// GetTime - Get elapsed time in seconds since [InitWindow]
func GetTime() float64 {
	return cGetTime()
}

// GetFPS - Get current FPS
func GetFPS() int32 {
	return cGetFPS()
}

// SwapScreenBuffer - Swap back buffer with front buffer (screen drawing)
func SwapScreenBuffer() {
	cSwapScreenBuffer()
}

// PollInputEvents - Register all input events
func PollInputEvents() {
	cPollInputEvents()
}

// WaitTime - Wait for some time (halt program execution)
func WaitTime(seconds float64) {
	time.Sleep(time.Duration(seconds * float64(time.Second)))
}

// GetRandomValue - Returns a random value between min and max (both included)
//
// Note: You can use go's math/rand package instead
func GetRandomValue(min, max int32) int32 {
	return cGetRandomValue(min, max)
}

// TakeScreenshot - Takes a screenshot of current screen (filename extension defines format)
func TakeScreenshot(fileName string) {
	defer scratch.release(scratch.mark())
	cTakeScreenshot(scratch.str(fileName))
}

// SetConfigFlags - Setup init configuration flags (view FLAGS)
func SetConfigFlags(flags uint32) {
	cSetConfigFlags(flags)
}

// OpenURL - Open URL with default system browser (if available)
func OpenURL(url string) {
	defer scratch.release(scratch.mark())
	cOpenURL(scratch.str(url))
}

// TraceLog - Trace log messages showing (DEBUG, INFO, WARNING, ERROR)
func TraceLog(logLevel TraceLogLevel, text string, args ...any) {
	defer scratch.release(scratch.mark())
	cTraceLogText(int32(logLevel), scratch.str(fmt.Sprintf(text, args...)))
}

// SetTraceLogLevel - Set the current threshold (minimum) log level
func SetTraceLogLevel(logLevel TraceLogLevel) {
	cSetTraceLogLevel(int32(logLevel))
}

// SetTraceLogCallback - set a call-back function for trace log
func SetTraceLogCallback(fn TraceLogCallbackFun) {
	traceLogCallback = fn
	cSetTraceLogCallbackEnabled(boolToInt32(fn != nil))
}

// IsFileDropped - Check if a file has been dropped into window
func IsFileDropped() bool {
	return cIsFileDropped() != 0
}

// LoadDroppedFiles - Load dropped filepaths
func LoadDroppedFiles() []string {
	defer scratch.release(scratch.mark())
	cLoadDroppedFiles()
	list := readRet[cFilePathList]()
	paths := goSlice[uint32](list.Paths, int(list.Count))
	result := make([]string, len(paths))
	for i, p := range paths {
		result[i] = goString(p)
	}
	cUnloadDroppedFiles(scratchValue(&list))
	return result
}

// UnloadDroppedFiles - Unload dropped filepaths
func UnloadDroppedFiles() {}

// LoadAutomationEventList - Load automation events list from file, NULL for empty list, capacity = MAX_AUTOMATION_EVENTS
func LoadAutomationEventList(fileName string) AutomationEventList {
	defer scratch.release(scratch.mark())
	cLoadAutomationEventList(scratch.str(fileName))
	return automationListFromC(readRet[cAutomationEventList]())
}

// UnloadAutomationEventList - Unload automation events list from file
func UnloadAutomationEventList(list *AutomationEventList) {
	if currentAutomation == list {
		currentAutomation = nil
	}
	e := automationLists[list.Events]
	if e == nil || list.Events == nil {
		return
	}
	cUnloadAutomationEventList(e.at)
	cFree(e.at)
	delete(automationLists, list.Events)
}

// ExportAutomationEventList - Export automation events list as text file
func ExportAutomationEventList(list AutomationEventList, fileName string) bool {
	defer scratch.release(scratch.mark())
	return cExportAutomationEventList(automationListIn(&list), scratch.str(fileName)) != 0
}

// SetAutomationEventList - Set automation event list to record to
func SetAutomationEventList(list *AutomationEventList) {
	currentAutomation = nil
	if list == nil {
		cSetAutomationEventList(0)
		return
	}
	e := automationLists[list.Events]
	if e == nil || list.Events == nil {
		c := cAutomationEventList{Capacity: list.Capacity, Count: list.Count, Events: cOwned(ptrSlice(list.Events, int(list.Capacity)))}
		e = &automationEntry{at: persist(&c), c: c}
		automationLists[list.Events] = e
	}
	automationListIn(list)
	currentAutomation = list
	cSetAutomationEventList(e.at)
}

// SetAutomationEventBaseFrame - Set automation event internal base frame to start recording
func SetAutomationEventBaseFrame(frame int) {
	cSetAutomationEventBaseFrame(int32(frame))
}

// StartAutomationEventRecording - Start recording automation events (AutomationEventList must be set)
func StartAutomationEventRecording() {
	cStartAutomationEventRecording()
}

// StopAutomationEventRecording - Stop recording automation events
func StopAutomationEventRecording() {
	cStopAutomationEventRecording()
}

// PlayAutomationEvent - Play a recorded automation event
func PlayAutomationEvent(event AutomationEvent) {
	cPlayAutomationEvent(event.Frame, event.Type, event.Params[0], event.Params[1], event.Params[2], event.Params[3])
}

// IsKeyPressed - Check if a key has been pressed once
func IsKeyPressed(key int32) bool {
	return cIsKeyPressed(key) != 0
}

// IsKeyPressedRepeat - Check if a key has been pressed again
func IsKeyPressedRepeat(key int32) bool {
	return cIsKeyPressedRepeat(key) != 0
}

// IsKeyDown - Check if a key is being pressed
func IsKeyDown(key int32) bool {
	return cIsKeyDown(key) != 0
}

// IsKeyReleased - Check if a key has been released once
func IsKeyReleased(key int32) bool {
	return cIsKeyReleased(key) != 0
}

// IsKeyUp - Check if a key is NOT being pressed
func IsKeyUp(key int32) bool {
	return cIsKeyUp(key) != 0
}

// GetKeyPressed - Get key pressed (keycode), call it multiple times for keys queued, returns 0 when the queue is empty
func GetKeyPressed() int32 {
	return cGetKeyPressed()
}

// GetCharPressed - Get char pressed (unicode), call it multiple times for chars queued, returns 0 when the queue is empty
func GetCharPressed() int32 {
	return cGetCharPressed()
}

// GetKeyName - Get name of a QWERTY key on the current keyboard layout
// (eg returns string 'q' for KEY_A on an AZERTY keyboard)
func GetKeyName(key int32) string {
	return goString(cGetKeyName(key))
}

// SetExitKey - Set a custom key to exit program (default is ESC)
func SetExitKey(key int32) {
	cSetExitKey(key)
}

// IsGamepadAvailable - Check if a gamepad is available
func IsGamepadAvailable(gamepad int32) bool {
	return cIsGamepadAvailable(gamepad) != 0
}

// GetGamepadName - Get gamepad internal name id
func GetGamepadName(gamepad int32) string {
	return goString(cGetGamepadName(gamepad))
}

// IsGamepadButtonPressed - Check if a gamepad button has been pressed once
func IsGamepadButtonPressed(gamepad, button int32) bool {
	return cIsGamepadButtonPressed(gamepad, button) != 0
}

// IsGamepadButtonDown - Check if a gamepad button is being pressed
func IsGamepadButtonDown(gamepad, button int32) bool {
	return cIsGamepadButtonDown(gamepad, button) != 0
}

// IsGamepadButtonReleased - Check if a gamepad button has been released once
func IsGamepadButtonReleased(gamepad, button int32) bool {
	return cIsGamepadButtonReleased(gamepad, button) != 0
}

// IsGamepadButtonUp - Check if a gamepad button is NOT being pressed
func IsGamepadButtonUp(gamepad, button int32) bool {
	return cIsGamepadButtonUp(gamepad, button) != 0
}

// GetGamepadButtonPressed - Get the last gamepad button pressed
func GetGamepadButtonPressed() int32 {
	return cGetGamepadButtonPressed()
}

// GetGamepadAxisCount - Get gamepad axis count for a gamepad
func GetGamepadAxisCount(gamepad int32) int32 {
	return cGetGamepadAxisCount(gamepad)
}

// GetGamepadAxisMovement - Get axis movement value for a gamepad axis
func GetGamepadAxisMovement(gamepad, axis int32) float32 {
	return cGetGamepadAxisMovement(gamepad, axis)
}

// SetGamepadMappings - Set internal gamepad mappings (SDL_GameControllerDB)
func SetGamepadMappings(mappings string) int32 {
	defer scratch.release(scratch.mark())
	return cSetGamepadMappings(scratch.str(mappings))
}

// SetGamepadVibration - Set gamepad vibration for both motors (duration in seconds)
func SetGamepadVibration(gamepad int32, leftMotor, rightMotor, duration float32) {
	cSetGamepadVibration(gamepad, leftMotor, rightMotor, duration)
}

// IsMouseButtonPressed - Check if a mouse button has been pressed once
func IsMouseButtonPressed(button MouseButton) bool {
	return cIsMouseButtonPressed(int32(button)) != 0
}

// IsMouseButtonDown - Check if a mouse button is being pressed
func IsMouseButtonDown(button MouseButton) bool {
	return cIsMouseButtonDown(int32(button)) != 0
}

// IsMouseButtonReleased - Check if a mouse button has been released once
func IsMouseButtonReleased(button MouseButton) bool {
	return cIsMouseButtonReleased(int32(button)) != 0
}

// IsMouseButtonUp - Check if a mouse button is NOT being pressed
func IsMouseButtonUp(button MouseButton) bool {
	return cIsMouseButtonUp(int32(button)) != 0
}

// GetMouseX - Get mouse position X
func GetMouseX() int32 {
	return cGetMouseX()
}

// GetMouseY - Get mouse position Y
func GetMouseY() int32 {
	return cGetMouseY()
}

// GetMousePosition - Get mouse position XY
func GetMousePosition() Vector2 {
	ret := cGetMousePosition()
	return *(*Vector2)(unsafe.Pointer(&ret))
}

// GetMouseDelta - Get mouse delta between frames
func GetMouseDelta() Vector2 {
	ret := cGetMouseDelta()
	return *(*Vector2)(unsafe.Pointer(&ret))
}

// SetMousePosition - Set mouse position XY
func SetMousePosition(x, y int) {
	cSetMousePosition(int32(x), int32(y))
}

// SetMouseOffset - Set mouse offset
func SetMouseOffset(offsetX, offsetY int) {
	cSetMouseOffset(int32(offsetX), int32(offsetY))
}

// SetMouseScale - Set mouse scaling
func SetMouseScale(scaleX, scaleY float32) {
	cSetMouseScale(scaleX, scaleY)
}

// GetMouseWheelMove - Get mouse wheel movement for X or Y, whichever is larger
func GetMouseWheelMove() float32 {
	return cGetMouseWheelMove()
}

// GetMouseWheelMoveV - Get mouse wheel movement for both X and Y
func GetMouseWheelMoveV() Vector2 {
	ret := cGetMouseWheelMoveV()
	return *(*Vector2)(unsafe.Pointer(&ret))
}

// SetMouseCursor - Set mouse cursor
func SetMouseCursor(cursor int32) {
	cSetMouseCursor(cursor)
}

// GetTouchX - Get touch position X for touch point 0 (relative to screen size)
func GetTouchX() int32 {
	return cGetTouchX()
}

// GetTouchY - Get touch position Y for touch point 0 (relative to screen size)
func GetTouchY() int32 {
	return cGetTouchY()
}

// GetTouchPosition - Get touch position XY for a touch point index (relative to screen size)
func GetTouchPosition(index int32) Vector2 {
	ret := cGetTouchPosition(index)
	return *(*Vector2)(unsafe.Pointer(&ret))
}

// GetTouchPointId - Get touch point identifier for given index
func GetTouchPointId(index int32) int32 {
	return cGetTouchPointId(index)
}

// GetTouchPointCount - Get number of touch points
func GetTouchPointCount() int32 {
	return cGetTouchPointCount()
}

// SetGesturesEnabled - Enable a set of gestures using flags
func SetGesturesEnabled(flags uint32) {
	cSetGesturesEnabled(flags)
}

// IsGestureDetected - Check if a gesture have been detected
func IsGestureDetected(gesture Gestures) bool {
	return cIsGestureDetected(uint32(gesture)) != 0
}

// GetGestureDetected - Get latest detected gesture
func GetGestureDetected() Gestures {
	return Gestures(cGetGestureDetected())
}

// GetGestureHoldDuration - Get gesture hold time in milliseconds
func GetGestureHoldDuration() float32 {
	return cGetGestureHoldDuration()
}

// GetGestureDragVector - Get gesture drag vector
func GetGestureDragVector() Vector2 {
	ret := cGetGestureDragVector()
	return *(*Vector2)(unsafe.Pointer(&ret))
}

// GetGestureDragAngle - Get gesture drag angle
func GetGestureDragAngle() float32 {
	return cGetGestureDragAngle()
}

// GetGesturePinchVector - Get gesture pinch delta
func GetGesturePinchVector() Vector2 {
	ret := cGetGesturePinchVector()
	return *(*Vector2)(unsafe.Pointer(&ret))
}

// GetGesturePinchAngle - Get gesture pinch angle
func GetGesturePinchAngle() float32 {
	return cGetGesturePinchAngle()
}

// SetShapesTexture - Set texture and rectangle to be used on shapes drawing
func SetShapesTexture(texture Texture2D, source Rectangle) {
	cSetShapesTexture(texture.ID, texture.Width, texture.Height, texture.Mipmaps, int32(texture.Format), source.X, source.Y, source.Width, source.Height)
}

// GetShapesTexture - Get texture that is used for shapes drawing
func GetShapesTexture() Texture2D {
	cGetShapesTexture()
	return readRet[Texture2D]()
}

// GetShapesTextureRectangle - Get texture source rectangle that is used for shapes drawing
func GetShapesTextureRectangle() Rectangle {
	cGetShapesTextureRectangle()
	return readRet[Rectangle]()
}

// DrawPixel - Draw a pixel
func DrawPixel(posX, posY int32, col color.RGBA) {
	cDrawPixel(posX, posY, colorToUint32(col))
}

// DrawPixelV - Draw a pixel (Vector version)
func DrawPixelV(position Vector2, col color.RGBA) {
	cDrawPixelV(position.X, position.Y, colorToUint32(col))
}

// DrawLine - Draw a line
func DrawLine(startPosX, startPosY, endPosX, endPosY int32, col color.RGBA) {
	cDrawLine(startPosX, startPosY, endPosX, endPosY, colorToUint32(col))
}

// DrawLineV - Draw a line (using gl lines)
func DrawLineV(startPos, endPos Vector2, col color.RGBA) {
	cDrawLineV(startPos.X, startPos.Y, endPos.X, endPos.Y, colorToUint32(col))
}

// DrawLineEx - Draw a line (using triangles/quads)
func DrawLineEx(startPos, endPos Vector2, thick float32, col color.RGBA) {
	cDrawLineEx(startPos.X, startPos.Y, endPos.X, endPos.Y, thick, colorToUint32(col))
}

// DrawLineStrip - Draw lines sequence (using gl lines)
func DrawLineStrip(points []Vector2, col color.RGBA) {
	defer scratch.release(scratch.mark())
	cDrawLineStrip(scratchSlice(points), int32(len(points)), colorToUint32(col))
}

// DrawLineBezier - Draw line segment cubic-bezier in-out interpolation
func DrawLineBezier(startPos, endPos Vector2, thick float32, col color.RGBA) {
	cDrawLineBezier(startPos.X, startPos.Y, endPos.X, endPos.Y, thick, colorToUint32(col))
}

// DrawLineDashed - Draw a dashed line
func DrawLineDashed(startPos, endPos Vector2, dashSize, spaceSize int32, col color.RGBA) {
	cDrawLineDashed(startPos.X, startPos.Y, endPos.X, endPos.Y, dashSize, spaceSize, colorToUint32(col))
}

// DrawCircle - Draw a color-filled circle
func DrawCircle(centerX, centerY int32, radius float32, col color.RGBA) {
	cDrawCircle(centerX, centerY, radius, colorToUint32(col))
}

// DrawCircleV - Draw a color-filled circle (Vector version)
func DrawCircleV(center Vector2, radius float32, col color.RGBA) {
	cDrawCircleV(center.X, center.Y, radius, colorToUint32(col))
}

// DrawCircleGradient - Draw a gradient-filled circle
func DrawCircleGradient(center Vector2, radius float32, inner, outer color.RGBA) {
	cDrawCircleGradient(center.X, center.Y, radius, colorToUint32(inner), colorToUint32(outer))
}

// DrawCircleSector - Draw a piece of a circle
func DrawCircleSector(center Vector2, radius, startAngle, endAngle float32, segments int32, col color.RGBA) {
	cDrawCircleSector(center.X, center.Y, radius, startAngle, endAngle, segments, colorToUint32(col))
}

// DrawCircleSectorLines - Draw circle sector outline
func DrawCircleSectorLines(center Vector2, radius, startAngle, endAngle float32, segments int32, col color.RGBA) {
	cDrawCircleSectorLines(center.X, center.Y, radius, startAngle, endAngle, segments, colorToUint32(col))
}

// DrawCircleLines - Draw circle outline
func DrawCircleLines(centerX, centerY int32, radius float32, col color.RGBA) {
	cDrawCircleLines(centerX, centerY, radius, colorToUint32(col))
}

// DrawCircleLinesV - Draw circle outline (Vector version)
func DrawCircleLinesV(center Vector2, radius float32, col color.RGBA) {
	cDrawCircleLinesV(center.X, center.Y, radius, colorToUint32(col))
}

// DrawEllipse - Draw ellipse
func DrawEllipse(centerX, centerY int32, radiusH, radiusV float32, col color.RGBA) {
	cDrawEllipse(centerX, centerY, radiusH, radiusV, colorToUint32(col))
}

// DrawEllipseV - Draw ellipse (Vector version)
func DrawEllipseV(center Vector2, radiusH, radiusV float32, col color.RGBA) {
	cDrawEllipseV(center.X, center.Y, radiusH, radiusV, colorToUint32(col))
}

// DrawEllipseLines - Draw ellipse outline
func DrawEllipseLines(centerX, centerY int32, radiusH, radiusV float32, col color.RGBA) {
	cDrawEllipseLines(centerX, centerY, radiusH, radiusV, colorToUint32(col))
}

// DrawEllipseLinesV - Draw ellipse outline (Vector version)
func DrawEllipseLinesV(center Vector2, radiusH, radiusV float32, col color.RGBA) {
	cDrawEllipseLinesV(center.X, center.Y, radiusH, radiusV, colorToUint32(col))
}

// DrawRing - Draw ring
func DrawRing(center Vector2, innerRadius, outerRadius, startAngle, endAngle float32, segments int32, col color.RGBA) {
	cDrawRing(center.X, center.Y, innerRadius, outerRadius, startAngle, endAngle, segments, colorToUint32(col))
}

// DrawRingLinesEx - Draw ring outline with line thickness
func DrawRingLinesEx(center Vector2, innerRadius, outerRadius, startAngle, endAngle float32, segments int32, thick float32, col color.RGBA) {
	cDrawRingLinesEx(center.X, center.Y, innerRadius, outerRadius, startAngle, endAngle, segments, thick, colorToUint32(col))
}

// DrawCircleLinesEx - Draw circle outline with line thickness
func DrawCircleLinesEx(center Vector2, radius, thick float32, col color.RGBA) {
	cDrawCircleLinesEx(center.X, center.Y, radius, thick, colorToUint32(col))
}

// DrawCircleSectorLinesEx - Draw circle sector outline with thickness
func DrawCircleSectorLinesEx(center Vector2, radius, startAngle, endAngle float32, segments int32, thick float32, col color.RGBA) {
	cDrawCircleSectorLinesEx(center.X, center.Y, radius, startAngle, endAngle, segments, thick, colorToUint32(col))
}

// DrawEllipseLinesEx - Draw ellipse outline with line thickness
func DrawEllipseLinesEx(center Vector2, radiusH, radiusV, thick float32, col color.RGBA) {
	cDrawEllipseLinesEx(center.X, center.Y, radiusH, radiusV, thick, colorToUint32(col))
}

// DrawRingLines - Draw ring outline
func DrawRingLines(center Vector2, innerRadius, outerRadius, startAngle, endAngle float32, segments int32, col color.RGBA) {
	cDrawRingLines(center.X, center.Y, innerRadius, outerRadius, startAngle, endAngle, segments, colorToUint32(col))
}

// DrawRectangle - Draw a color-filled rectangle
func DrawRectangle(posX, posY, width, height int32, col color.RGBA) {
	cDrawRectangle(posX, posY, width, height, colorToUint32(col))
}

// DrawRectangleV - Draw a color-filled rectangle (Vector version)
func DrawRectangleV(position Vector2, size Vector2, col color.RGBA) {
	cDrawRectangleV(position.X, position.Y, size.X, size.Y, colorToUint32(col))
}

// DrawRectangleRec - Draw a color-filled rectangle
func DrawRectangleRec(rec Rectangle, col color.RGBA) {
	cDrawRectangleRec(rec.X, rec.Y, rec.Width, rec.Height, colorToUint32(col))
}

// DrawRectanglePro - Draw a color-filled rectangle with pro parameters
func DrawRectanglePro(rec Rectangle, origin Vector2, rotation float32, col color.RGBA) {
	cDrawRectanglePro(rec.X, rec.Y, rec.Width, rec.Height, origin.X, origin.Y, rotation, colorToUint32(col))
}

// DrawRectangleGradientV - Draw a vertical-gradient-filled rectangle
func DrawRectangleGradientV(posX, posY, width, height int32, top, bottom color.RGBA) {
	cDrawRectangleGradientV(posX, posY, width, height, colorToUint32(top), colorToUint32(bottom))
}

// DrawRectangleGradientH - Draw a horizontal-gradient-filled rectangle
func DrawRectangleGradientH(posX, posY, width, height int32, left, right color.RGBA) {
	cDrawRectangleGradientH(posX, posY, width, height, colorToUint32(left), colorToUint32(right))
}

// DrawRectangleGradientEx - Draw a gradient-filled rectangle with custom vertex colors
func DrawRectangleGradientEx(rec Rectangle, topLeft, bottomLeft, bottomRight, topRight color.RGBA) {
	cDrawRectangleGradientEx(rec.X, rec.Y, rec.Width, rec.Height, colorToUint32(topLeft), colorToUint32(bottomLeft), colorToUint32(bottomRight), colorToUint32(topRight))
}

// DrawRectangleLines - Draw rectangle outline
func DrawRectangleLines(posX, posY, width, height int32, col color.RGBA) {
	cDrawRectangleLines(posX, posY, width, height, colorToUint32(col))
}

// DrawRectangleLinesEx - Draw rectangle outline with extended parameters
func DrawRectangleLinesEx(rec Rectangle, lineThick float32, col color.RGBA) {
	cDrawRectangleLinesEx(rec.X, rec.Y, rec.Width, rec.Height, lineThick, colorToUint32(col))
}

// DrawRectangleRounded - Draw rectangle with rounded edges
func DrawRectangleRounded(rec Rectangle, roundness float32, segments int32, col color.RGBA) {
	cDrawRectangleRounded(rec.X, rec.Y, rec.Width, rec.Height, roundness, segments, colorToUint32(col))
}

// DrawRectangleRoundedLines - Draw rectangle lines with rounded edges
func DrawRectangleRoundedLines(rec Rectangle, roundness float32, segments int32, col color.RGBA) {
	cDrawRectangleRoundedLines(rec.X, rec.Y, rec.Width, rec.Height, roundness, segments, colorToUint32(col))
}

// DrawRectangleRoundedLinesEx - Draw rectangle with rounded edges outline
func DrawRectangleRoundedLinesEx(rec Rectangle, roundness float32, segments int32, lineThick float32, col color.RGBA) {
	cDrawRectangleRoundedLinesEx(rec.X, rec.Y, rec.Width, rec.Height, roundness, segments, lineThick, colorToUint32(col))
}

// DrawTriangle - Draw a color-filled triangle (vertex in counter-clockwise order!)
func DrawTriangle(v1, v2, v3 Vector2, col color.RGBA) {
	cDrawTriangle(v1.X, v1.Y, v2.X, v2.Y, v3.X, v3.Y, colorToUint32(col))
}

// DrawTriangleGradient - Draw triangle with interpolated colors, counter-clockwise vertex/color order
func DrawTriangleGradient(v1, v2, v3 Vector2, c1, c2, c3 color.RGBA) {
	cDrawTriangleGradient(v1.X, v1.Y, v2.X, v2.Y, v3.X, v3.Y, colorToUint32(c1), colorToUint32(c2), colorToUint32(c3))
}

// DrawTriangleLinesEx - Draw triangle outline with line thickness, counter-clockwise vertex order
func DrawTriangleLinesEx(v1, v2, v3 Vector2, thick float32, col color.RGBA) {
	cDrawTriangleLinesEx(v1.X, v1.Y, v2.X, v2.Y, v3.X, v3.Y, thick, colorToUint32(col))
}

// DrawTriangleLines - Draw triangle outline (vertex in counter-clockwise order!)
func DrawTriangleLines(v1, v2, v3 Vector2, col color.RGBA) {
	cDrawTriangleLines(v1.X, v1.Y, v2.X, v2.Y, v3.X, v3.Y, colorToUint32(col))
}

// DrawTriangleFan - Draw a triangle fan defined by points (first vertex is the center)
func DrawTriangleFan(points []Vector2, col color.RGBA) {
	defer scratch.release(scratch.mark())
	cDrawTriangleFan(scratchSlice(points), int32(len(points)), colorToUint32(col))
}

// DrawTriangleStrip - Draw a triangle strip defined by points
func DrawTriangleStrip(points []Vector2, col color.RGBA) {
	defer scratch.release(scratch.mark())
	cDrawTriangleStrip(scratchSlice(points), int32(len(points)), colorToUint32(col))
}

// DrawPoly - Draw a regular polygon (Vector version)
func DrawPoly(center Vector2, sides int32, radius, rotation float32, col color.RGBA) {
	cDrawPoly(center.X, center.Y, sides, radius, rotation, colorToUint32(col))
}

// DrawPolyLines - Draw a polygon outline of n sides
func DrawPolyLines(center Vector2, sides int32, radius, rotation float32, col color.RGBA) {
	cDrawPolyLines(center.X, center.Y, sides, radius, rotation, colorToUint32(col))
}

// DrawPolyLinesEx - Draw a polygon outline of n sides with extended parameters
func DrawPolyLinesEx(center Vector2, sides int32, radius float32, rotation float32, lineThick float32, col color.RGBA) {
	cDrawPolyLinesEx(center.X, center.Y, sides, radius, rotation, lineThick, colorToUint32(col))
}

// DrawSplineLinear - Draw spline: Linear, minimum 2 points
func DrawSplineLinear(points []Vector2, thick float32, col color.RGBA) {
	defer scratch.release(scratch.mark())
	cDrawSplineLinear(scratchSlice(points), int32(len(points)), thick, colorToUint32(col))
}

// DrawSplineBasis - Draw spline: B-Spline, minimum 4 points
func DrawSplineBasis(points []Vector2, thick float32, col color.RGBA) {
	defer scratch.release(scratch.mark())
	cDrawSplineBasis(scratchSlice(points), int32(len(points)), thick, colorToUint32(col))
}

// DrawSplineCatmullRom - Draw spline: Catmull-Rom, minimum 4 points
func DrawSplineCatmullRom(points []Vector2, thick float32, col color.RGBA) {
	defer scratch.release(scratch.mark())
	cDrawSplineCatmullRom(scratchSlice(points), int32(len(points)), thick, colorToUint32(col))
}

// DrawSplineBezierQuadratic - Draw spline: Quadratic Bezier, minimum 3 points (1 control point): [p1, c2, p3, c4...]
func DrawSplineBezierQuadratic(points []Vector2, thick float32, col color.RGBA) {
	defer scratch.release(scratch.mark())
	cDrawSplineBezierQuadratic(scratchSlice(points), int32(len(points)), thick, colorToUint32(col))
}

// DrawSplineBezierCubic - Draw spline: Cubic Bezier, minimum 4 points (2 control points): [p1, c2, c3, p4, c5, c6...]
func DrawSplineBezierCubic(points []Vector2, thick float32, col color.RGBA) {
	defer scratch.release(scratch.mark())
	cDrawSplineBezierCubic(scratchSlice(points), int32(len(points)), thick, colorToUint32(col))
}

// DrawSplineSegmentLinear - Draw spline segment: Linear, 2 points
func DrawSplineSegmentLinear(p1, p2 Vector2, thick float32, col color.RGBA) {
	cDrawSplineSegmentLinear(p1.X, p1.Y, p2.X, p2.Y, thick, colorToUint32(col))
}

// DrawSplineSegmentBasis - Draw spline segment: B-Spline, 4 points
func DrawSplineSegmentBasis(p1, p2, p3, p4 Vector2, thick float32, col color.RGBA) {
	cDrawSplineSegmentBasis(p1.X, p1.Y, p2.X, p2.Y, p3.X, p3.Y, p4.X, p4.Y, thick, colorToUint32(col))
}

// DrawSplineSegmentCatmullRom - Draw spline segment: Catmull-Rom, 4 points
func DrawSplineSegmentCatmullRom(p1, p2, p3, p4 Vector2, thick float32, col color.RGBA) {
	cDrawSplineSegmentCatmullRom(p1.X, p1.Y, p2.X, p2.Y, p3.X, p3.Y, p4.X, p4.Y, thick, colorToUint32(col))
}

// DrawSplineSegmentBezierQuadratic - Draw spline segment: Quadratic Bezier, 2 points, 1 control point
func DrawSplineSegmentBezierQuadratic(p1, c2, p3 Vector2, thick float32, col color.RGBA) {
	cDrawSplineSegmentBezierQuadratic(p1.X, p1.Y, c2.X, c2.Y, p3.X, p3.Y, thick, colorToUint32(col))
}

// DrawSplineSegmentBezierCubic - Draw spline segment: Cubic Bezier, 2 points, 2 control points
func DrawSplineSegmentBezierCubic(p1, c2, c3, p4 Vector2, thick float32, col color.RGBA) {
	cDrawSplineSegmentBezierCubic(p1.X, p1.Y, c2.X, c2.Y, c3.X, c3.Y, p4.X, p4.Y, thick, colorToUint32(col))
}

// GetSplinePointLinear - Get (evaluate) spline point: Linear
func GetSplinePointLinear(startPos, endPos Vector2, t float32) Vector2 {
	ret := cGetSplinePointLinear(startPos.X, startPos.Y, endPos.X, endPos.Y, t)
	return *(*Vector2)(unsafe.Pointer(&ret))
}

// GetSplinePointBasis - Get (evaluate) spline point: B-Spline
func GetSplinePointBasis(p1, p2, p3, p4 Vector2, t float32) Vector2 {
	ret := cGetSplinePointBasis(p1.X, p1.Y, p2.X, p2.Y, p3.X, p3.Y, p4.X, p4.Y, t)
	return *(*Vector2)(unsafe.Pointer(&ret))
}

// GetSplinePointCatmullRom - Get (evaluate) spline point: Catmull-Rom
func GetSplinePointCatmullRom(p1, p2, p3, p4 Vector2, t float32) Vector2 {
	ret := cGetSplinePointCatmullRom(p1.X, p1.Y, p2.X, p2.Y, p3.X, p3.Y, p4.X, p4.Y, t)
	return *(*Vector2)(unsafe.Pointer(&ret))
}

// GetSplinePointBezierQuadratic - Get (evaluate) spline point: Quadratic Bezier
func GetSplinePointBezierQuadratic(p1, c2, p3 Vector2, t float32) Vector2 {
	ret := cGetSplinePointBezierQuadratic(p1.X, p1.Y, c2.X, c2.Y, p3.X, p3.Y, t)
	return *(*Vector2)(unsafe.Pointer(&ret))
}

// GetSplinePointBezierCubic - Get (evaluate) spline point: Cubic Bezier
func GetSplinePointBezierCubic(p1, c2, c3, p4 Vector2, t float32) Vector2 {
	ret := cGetSplinePointBezierCubic(p1.X, p1.Y, c2.X, c2.Y, c3.X, c3.Y, p4.X, p4.Y, t)
	return *(*Vector2)(unsafe.Pointer(&ret))
}

// CheckCollisionRecs - Check collision between two rectangles
func CheckCollisionRecs(rec1, rec2 Rectangle) bool {
	return cCheckCollisionRecs(rec1.X, rec1.Y, rec1.Width, rec1.Height, rec2.X, rec2.Y, rec2.Width, rec2.Height) != 0
}

// CheckCollisionCircles - Check collision between two circles
func CheckCollisionCircles(center1 Vector2, radius1 float32, center2 Vector2, radius2 float32) bool {
	return cCheckCollisionCircles(center1.X, center1.Y, radius1, center2.X, center2.Y, radius2) != 0
}

// CheckCollisionCircleRec - Check collision between circle and rectangle
func CheckCollisionCircleRec(center Vector2, radius float32, rec Rectangle) bool {
	return cCheckCollisionCircleRec(center.X, center.Y, radius, rec.X, rec.Y, rec.Width, rec.Height) != 0
}

// CheckCollisionCircleLine - Check if circle collides with a line created betweeen two points [p1] and [p2]
func CheckCollisionCircleLine(center Vector2, radius float32, p1, p2 Vector2) bool {
	return cCheckCollisionCircleLine(center.X, center.Y, radius, p1.X, p1.Y, p2.X, p2.Y) != 0
}

// CheckCollisionPointRec - Check if point is inside rectangle
func CheckCollisionPointRec(point Vector2, rec Rectangle) bool {
	return cCheckCollisionPointRec(point.X, point.Y, rec.X, rec.Y, rec.Width, rec.Height) != 0
}

// CheckCollisionPointCircle - Check if point is inside circle
func CheckCollisionPointCircle(point Vector2, center Vector2, radius float32) bool {
	return cCheckCollisionPointCircle(point.X, point.Y, center.X, center.Y, radius) != 0
}

// CheckCollisionPointTriangle - Check if point is inside a triangle
func CheckCollisionPointTriangle(point, p1, p2, p3 Vector2) bool {
	return cCheckCollisionPointTriangle(point.X, point.Y, p1.X, p1.Y, p2.X, p2.Y, p3.X, p3.Y) != 0
}

// CheckCollisionPointLine - Check if point belongs to line created between two points [p1] and [p2] with defined margin in pixels [threshold]
func CheckCollisionPointLine(point, p1, p2 Vector2, threshold int32) bool {
	return cCheckCollisionPointLine(point.X, point.Y, p1.X, p1.Y, p2.X, p2.Y, threshold) != 0
}

// CheckCollisionPointPoly - Check if point is within a polygon described by array of vertices
//
// NOTE: Based on http://jeffreythompson.org/collision-detection/poly-point.php
func CheckCollisionPointPoly(point Vector2, points []Vector2) bool {
	defer scratch.release(scratch.mark())
	return cCheckCollisionPointPoly(point.X, point.Y, scratchSlice(points), int32(len(points))) != 0
}

// CheckCollisionLines - Check the collision between two lines defined by two points each, returns collision point by reference
func CheckCollisionLines(startPos1, endPos1, startPos2, endPos2 Vector2, collisionPoint *Vector2) bool {
	defer scratch.release(scratch.mark())
	p4 := scratchPtr(collisionPoint)
	ret := cCheckCollisionLines(startPos1.X, startPos1.Y, endPos1.X, endPos1.Y, startPos2.X, startPos2.Y, endPos2.X, endPos2.Y, p4)
	readBack(collisionPoint, p4)
	return ret != 0
}

// GetCollisionRec - Get collision rectangle for two rectangles collision
func GetCollisionRec(rec1, rec2 Rectangle) Rectangle {
	cGetCollisionRec(rec1.X, rec1.Y, rec1.Width, rec1.Height, rec2.X, rec2.Y, rec2.Width, rec2.Height)
	return readRet[Rectangle]()
}

// LoadImage - Load image from file into CPU memory (RAM)
func LoadImage(fileName string) *Image {
	defer scratch.release(scratch.mark())
	cLoadImage(scratch.str(fileName))
	return imageRet()
}

// LoadImageRaw - Load image from RAW file data
func LoadImageRaw(fileName string, width, height int32, format PixelFormat, headerSize int32) *Image {
	defer scratch.release(scratch.mark())
	cLoadImageRaw(scratch.str(fileName), width, height, int32(format), headerSize)
	return imageRet()
}

// LoadImageAnim - Load image sequence from file (frames appended to image.data)
func LoadImageAnim(fileName string, frames *int32) *Image {
	defer scratch.release(scratch.mark())
	p := scratchPtr(frames)
	cLoadImageAnim(scratch.str(fileName), p)
	readBack(frames, p)
	return imageRet()
}

// LoadImageAnimFromMemory - Load image sequence from memory buffer
func LoadImageAnimFromMemory(fileType string, fileData []byte, dataSize int32, frames *int32) *Image {
	defer scratch.release(scratch.mark())
	p := scratchPtr(frames)
	cLoadImageAnimFromMemory(scratch.str(fileType), scratchSlice(fileData), dataSize, p)
	readBack(frames, p)
	return imageRet()
}

// LoadImageFromMemory - Load image from memory buffer, fileType refers to extension: i.e. '.png'
func LoadImageFromMemory(fileType string, fileData []byte, dataSize int32) *Image {
	defer scratch.release(scratch.mark())
	cLoadImageFromMemory(scratch.str(fileType), scratchSlice(fileData), dataSize)
	return imageRet()
}

// LoadImageFromTexture - Load image from GPU texture data
func LoadImageFromTexture(texture Texture2D) *Image {
	cLoadImageFromTexture(texture.ID, texture.Width, texture.Height, texture.Mipmaps, int32(texture.Format))
	return imageRet()
}

// LoadImageFromScreen - Load image from screen buffer and (screenshot)
func LoadImageFromScreen() *Image {
	cLoadImageFromScreen()
	return imageRet()
}

// IsImageValid - Check if an image is valid (data and parameters)
func IsImageValid(image *Image) bool {
	defer scratch.release(scratch.mark())
	p0 := imageIn(image)
	ret := cIsImageValid(p0)
	imageDone(p0)
	return ret != 0
}

// UnloadImage - Unload image from CPU memory (RAM)
func UnloadImage(image *Image) {
}

// ExportImage - Export image data to file, returns true on success
func ExportImage(image Image, fileName string) bool {
	defer scratch.release(scratch.mark())
	p0 := imageIn(&image)
	ret := cExportImage(p0, scratch.str(fileName))
	imageDone(p0)
	return ret != 0
}

// ExportImageToMemory - Export image to memory buffer
//
// The returned memory is a Go-managed slice. It doesn't need to be freed.
func ExportImageToMemory(image Image, fileType string) []byte {
	defer scratch.release(scratch.mark())
	img := imageIn(&image)
	var fileSize int32
	size := scratchValue(&fileSize)
	ret := cExportImageToMemory(img, scratch.str(fileType), size)
	imageDone(img)
	if ret == 0 {
		return nil
	}
	fileSize = readC[int32](size)
	result := goSlice[byte](ret, int(fileSize))
	cMemFree(ret)
	return result
}

// GenImageColor - Generate image: plain color
func GenImageColor(width, height int, col color.RGBA) *Image {
	cGenImageColor(int32(width), int32(height), colorToUint32(col))
	return imageRet()
}

// GenImageGradientLinear - Generate image: linear gradient, direction in degrees [0..360], 0=Vertical gradient
func GenImageGradientLinear(width, height, direction int, start, end color.RGBA) *Image {
	cGenImageGradientLinear(int32(width), int32(height), int32(direction), colorToUint32(start), colorToUint32(end))
	return imageRet()
}

// GenImageGradientRadial - Generate image: radial gradient
func GenImageGradientRadial(width, height int, density float32, inner, outer color.RGBA) *Image {
	cGenImageGradientRadial(int32(width), int32(height), density, colorToUint32(inner), colorToUint32(outer))
	return imageRet()
}

// GenImageGradientSquare - Generate image: square gradient
func GenImageGradientSquare(width, height int, density float32, inner, outer color.RGBA) *Image {
	cGenImageGradientSquare(int32(width), int32(height), density, colorToUint32(inner), colorToUint32(outer))
	return imageRet()
}

// GenImageChecked - Generate image: checked
func GenImageChecked(width, height, checksX, checksY int, col1, col2 color.RGBA) *Image {
	cGenImageChecked(int32(width), int32(height), int32(checksX), int32(checksY), colorToUint32(col1), colorToUint32(col2))
	return imageRet()
}

// GenImageWhiteNoise - Generate image: white noise
func GenImageWhiteNoise(width, height int, factor float32) *Image {
	cGenImageWhiteNoise(int32(width), int32(height), factor)
	return imageRet()
}

// GenImagePerlinNoise - Generate image: perlin noise
func GenImagePerlinNoise(width, height, offsetX, offsetY int, scale float32) *Image {
	cGenImagePerlinNoise(int32(width), int32(height), int32(offsetX), int32(offsetY), scale)
	return imageRet()
}

// GenImageCellular - Generate image: cellular algorithm, bigger tileSize means bigger cells
func GenImageCellular(width, height, tileSize int) *Image {
	cGenImageCellular(int32(width), int32(height), int32(tileSize))
	return imageRet()
}

// GenImageText - Generate image: grayscale image from text data
func GenImageText(width, height int, text string) *Image {
	defer scratch.release(scratch.mark())
	cGenImageText(int32(width), int32(height), scratch.str(text))
	return imageRet()
}

// ImageCopy - Create an image duplicate (useful for transformations)
func ImageCopy(image *Image) *Image {
	defer scratch.release(scratch.mark())
	p0 := imageIn(image)
	cImageCopy(p0)
	imageDone(p0)
	return imageRet()
}

// ImageFromImage - Create an image from another image piece
func ImageFromImage(image Image, rec Rectangle) Image {
	defer scratch.release(scratch.mark())
	p0 := imageIn(&image)
	cImageFromImage(p0, rec.X, rec.Y, rec.Width, rec.Height)
	imageDone(p0)
	return *imageRet()
}

// ImageFromChannel - Create an image from a selected channel of another image (GRAYSCALE)
func ImageFromChannel(image Image, selectedChannel int32) Image {
	defer scratch.release(scratch.mark())
	p0 := imageIn(&image)
	cImageFromChannel(p0, selectedChannel)
	imageDone(p0)
	return *imageRet()
}

// ImageText - Create an image from text (default font)
func ImageText(text string, fontSize int32, col color.RGBA) *Image {
	defer scratch.release(scratch.mark())
	cImageText(scratch.str(text), fontSize, colorToUint32(col))
	return imageRet()
}

// ImageTextEx - Create an image from text (custom sprite font)
func ImageTextEx(font Font, text string, fontSize, spacing float32, tint color.RGBA) *Image {
	defer scratch.release(scratch.mark())
	cImageTextEx(fontIn(&font), scratch.str(text), fontSize, spacing, colorToUint32(tint))
	return imageRet()
}

// ImageFormat - Convert image data to desired format
func ImageFormat(image *Image, newFormat PixelFormat) {
	defer scratch.release(scratch.mark())
	p0 := imageIn(image)
	cImageFormat(p0, int32(newFormat))
	imageBack(image, p0)
}

// ImageToPOT - Convert image to POT (power-of-two)
func ImageToPOT(image *Image, fill color.RGBA) {
	defer scratch.release(scratch.mark())
	p0 := imageIn(image)
	cImageToPOT(p0, colorToUint32(fill))
	imageBack(image, p0)
}

// ImageCrop - Crop an image to a defined rectangle
func ImageCrop(image *Image, crop Rectangle) {
	defer scratch.release(scratch.mark())
	p0 := imageIn(image)
	cImageCrop(p0, crop.X, crop.Y, crop.Width, crop.Height)
	imageBack(image, p0)
}

// ImageAlphaCrop - Crop image depending on alpha value
func ImageAlphaCrop(image *Image, threshold float32) {
	defer scratch.release(scratch.mark())
	p0 := imageIn(image)
	cImageAlphaCrop(p0, threshold)
	imageBack(image, p0)
}

// ImageAlphaClear - Clear alpha channel to desired color
func ImageAlphaClear(image *Image, col color.RGBA, threshold float32) {
	defer scratch.release(scratch.mark())
	p0 := imageIn(image)
	cImageAlphaClear(p0, colorToUint32(col), threshold)
	imageBack(image, p0)
}

// ImageAlphaMask - Apply alpha mask to image
func ImageAlphaMask(image, alphaMask *Image) {
	defer scratch.release(scratch.mark())
	p0 := imageIn(image)
	p1 := imageIn(alphaMask)
	cImageAlphaMask(p0, p1)
	imageBack(image, p0)
	imageDone(p1)
}

// ImageAlphaPremultiply - Premultiply alpha channel
func ImageAlphaPremultiply(image *Image) {
	defer scratch.release(scratch.mark())
	p0 := imageIn(image)
	cImageAlphaPremultiply(p0)
	imageBack(image, p0)
}

// ImageBlurGaussian - Apply Gaussian blur using a box blur approximation
func ImageBlurGaussian(image *Image, blurSize int32) {
	defer scratch.release(scratch.mark())
	p0 := imageIn(image)
	cImageBlurGaussian(p0, blurSize)
	imageBack(image, p0)
}

// ImageKernelConvolution - Apply custom square convolution kernel to image
func ImageKernelConvolution(image *Image, kernel []float32) {
	defer scratch.release(scratch.mark())
	p0 := imageIn(image)
	cImageKernelConvolution(p0, scratchSlice(kernel), int32(len(kernel)))
	imageBack(image, p0)
}

// ImageResize - Resize image (Bicubic scaling algorithm)
func ImageResize(image *Image, newWidth, newHeight int32) {
	defer scratch.release(scratch.mark())
	p0 := imageIn(image)
	cImageResize(p0, newWidth, newHeight)
	imageBack(image, p0)
}

// ImageResizeNN - Resize image (Nearest-Neighbor scaling algorithm)
func ImageResizeNN(image *Image, newWidth, newHeight int32) {
	defer scratch.release(scratch.mark())
	p0 := imageIn(image)
	cImageResizeNN(p0, newWidth, newHeight)
	imageBack(image, p0)
}

// ImageResizeCanvas - Resize canvas and fill with color
func ImageResizeCanvas(image *Image, newWidth, newHeight, offsetX, offsetY int32, fill color.RGBA) {
	defer scratch.release(scratch.mark())
	p0 := imageIn(image)
	cImageResizeCanvas(p0, newWidth, newHeight, offsetX, offsetY, colorToUint32(fill))
	imageBack(image, p0)
}

// ImageMipmaps - Compute all mipmap levels for a provided image
func ImageMipmaps(image *Image) {
	defer scratch.release(scratch.mark())
	p0 := imageIn(image)
	cImageMipmaps(p0)
	imageBack(image, p0)
}

// ImageDither - Dither image data to 16bpp or lower (Floyd-Steinberg dithering)
func ImageDither(image *Image, rBpp, gBpp, bBpp, aBpp int32) {
	defer scratch.release(scratch.mark())
	p0 := imageIn(image)
	cImageDither(p0, rBpp, gBpp, bBpp, aBpp)
	imageBack(image, p0)
}

// ImageFlipVertical - Flip image vertically
func ImageFlipVertical(image *Image) {
	defer scratch.release(scratch.mark())
	p0 := imageIn(image)
	cImageFlipVertical(p0)
	imageBack(image, p0)
}

// ImageFlipHorizontal - Flip image horizontally
func ImageFlipHorizontal(image *Image) {
	defer scratch.release(scratch.mark())
	p0 := imageIn(image)
	cImageFlipHorizontal(p0)
	imageBack(image, p0)
}

// ImageRotate - Rotate image by input angle in degrees (-359 to 359)
func ImageRotate(image *Image, degrees int32) {
	defer scratch.release(scratch.mark())
	p0 := imageIn(image)
	cImageRotate(p0, degrees)
	imageBack(image, p0)
}

// ImageRotateCW - Rotate image clockwise 90deg
func ImageRotateCW(image *Image) {
	defer scratch.release(scratch.mark())
	p0 := imageIn(image)
	cImageRotateCW(p0)
	imageBack(image, p0)
}

// ImageRotateCCW - Rotate image counter-clockwise 90deg
func ImageRotateCCW(image *Image) {
	defer scratch.release(scratch.mark())
	p0 := imageIn(image)
	cImageRotateCCW(p0)
	imageBack(image, p0)
}

// ImageColorTint - Modify image color: tint
func ImageColorTint(image *Image, col color.RGBA) {
	defer scratch.release(scratch.mark())
	p0 := imageIn(image)
	cImageColorTint(p0, colorToUint32(col))
	imageBack(image, p0)
}

// ImageColorInvert - Modify image color: invert
func ImageColorInvert(image *Image) {
	defer scratch.release(scratch.mark())
	p0 := imageIn(image)
	cImageColorInvert(p0)
	imageBack(image, p0)
}

// ImageColorGrayscale - Modify image color: grayscale
func ImageColorGrayscale(image *Image) {
	defer scratch.release(scratch.mark())
	p0 := imageIn(image)
	cImageColorGrayscale(p0)
	imageBack(image, p0)
}

// ImageColorContrast - Modify image color: contrast (-100 to 100)
func ImageColorContrast(image *Image, contrast int32) {
	defer scratch.release(scratch.mark())
	p0 := imageIn(image)
	cImageColorContrast(p0, contrast)
	imageBack(image, p0)
}

// ImageColorBrightness - Modify image color: brightness (-255 to 255)
func ImageColorBrightness(image *Image, brightness int32) {
	defer scratch.release(scratch.mark())
	p0 := imageIn(image)
	cImageColorBrightness(p0, brightness)
	imageBack(image, p0)
}

// ImageColorReplace - Modify image color: replace color
func ImageColorReplace(image *Image, col, replace color.RGBA) {
	defer scratch.release(scratch.mark())
	p0 := imageIn(image)
	cImageColorReplace(p0, colorToUint32(col), colorToUint32(replace))
	imageBack(image, p0)
}

// LoadImageColors - Get pixel data from image as a Color slice
func LoadImageColors(img *Image) []color.RGBA {
	defer scratch.release(scratch.mark())
	p := imageIn(img)
	ret := cLoadImageColors(p)
	imageDone(p)
	colors := goSlice[color.RGBA](ret, int(img.Width*img.Height))
	cFree(ret)
	return colors
}

// LoadImagePalette - Load colors palette from image as a Color array (RGBA - 32bit)
func LoadImagePalette(img Image, maxPaletteSize int32) []color.RGBA {
	defer scratch.release(scratch.mark())
	p := imageIn(&img)
	var colorCount int32
	count := scratchValue(&colorCount)
	ret := cLoadImagePalette(p, maxPaletteSize, count)
	imageDone(p)
	colors := goSlice[color.RGBA](ret, int(readC[int32](count)))
	cFree(ret)
	return colors
}

// UnloadImageColors - Unload color data loaded with LoadImageColors()
func UnloadImageColors(colors []color.RGBA) {
}

// UnloadImagePalette - Unload colors palette loaded with LoadImagePalette()
func UnloadImagePalette(colors []color.RGBA) {
}

// GetImageAlphaBorder - Get image alpha border rectangle
func GetImageAlphaBorder(image Image, threshold float32) Rectangle {
	defer scratch.release(scratch.mark())
	p0 := imageIn(&image)
	cGetImageAlphaBorder(p0, threshold)
	imageDone(p0)
	return readRet[Rectangle]()
}

// GetImageColor - Get image pixel color at (x, y) position
func GetImageColor(image Image, x, y int32) color.RGBA {
	defer scratch.release(scratch.mark())
	c := cImage{Width: image.Width, Height: image.Height, Mipmaps: image.Mipmaps, Format: int32(image.Format)}
	if pixel, n := imagePixel(&image, x, y); pixel != nil {
		c.Data, c.Width, c.Height, x, y = scratch.bytes(pixel, n), 1, 1, 0, 0
	}
	return uint32ToColor(cGetImageColor(scratchValue(&c), x, y))
}

// ImageClearBackground - Clear image background with given color
func ImageClearBackground(dst *Image, col color.RGBA) {
	defer scratch.release(scratch.mark())
	p0 := imageIn(dst)
	cImageClearBackground(p0, colorToUint32(col))
	imageBack(dst, p0)
}

// ImageDrawPixel - Draw pixel within an image
func ImageDrawPixel(dst *Image, posX, posY int32, col color.RGBA) {
	pixel, n := imagePixel(dst, posX, posY)
	if pixel == nil {
		return
	}
	defer scratch.release(scratch.mark())
	c := cImage{Data: scratch.bytes(pixel, n), Width: 1, Height: 1, Mipmaps: dst.Mipmaps, Format: int32(dst.Format)}
	cImageDrawPixel(scratchValue(&c), 0, 0, colorToUint32(col))
	copyOut(pixel, c.Data, n)
}

// ImageDrawPixelV - Draw pixel within an image (Vector version)
func ImageDrawPixelV(dst *Image, position Vector2, col color.RGBA) {
	ImageDrawPixel(dst, int32(position.X), int32(position.Y), col)
}

// ImageDrawLine - Draw line within an image
func ImageDrawLine(dst *Image, startPosX, startPosY, endPosX, endPosY int32, col color.RGBA) {
	defer scratch.release(scratch.mark())
	p0 := imageIn(dst)
	cImageDrawLine(p0, startPosX, startPosY, endPosX, endPosY, colorToUint32(col))
	imageBack(dst, p0)
}

// ImageDrawLineV - Draw line within an image (Vector version)
func ImageDrawLineV(dst *Image, start, end Vector2, col color.RGBA) {
	defer scratch.release(scratch.mark())
	p0 := imageIn(dst)
	cImageDrawLineV(p0, start.X, start.Y, end.X, end.Y, colorToUint32(col))
	imageBack(dst, p0)
}

// ImageDrawLineEx - Draw a line defining thickness within an image
func ImageDrawLineEx(dst *Image, start, end Vector2, thick int32, col color.RGBA) {
	defer scratch.release(scratch.mark())
	p0 := imageIn(dst)
	cImageDrawLineEx(p0, start.X, start.Y, end.X, end.Y, thick, colorToUint32(col))
	imageBack(dst, p0)
}

// ImageDrawLineStrip - Draw a lines sequence within an image
func ImageDrawLineStrip(dst *Image, points []Vector2, col color.RGBA) {
	defer scratch.release(scratch.mark())
	p0 := imageIn(dst)
	cImageDrawLineStrip(p0, scratchSlice(points), int32(len(points)), colorToUint32(col))
	imageBack(dst, p0)
}

// ImageDrawCircle - Draw a filled circle within an image
func ImageDrawCircle(dst *Image, centerX, centerY, radius int32, col color.RGBA) {
	defer scratch.release(scratch.mark())
	p0 := imageIn(dst)
	cImageDrawCircle(p0, centerX, centerY, radius, colorToUint32(col))
	imageBack(dst, p0)
}

// ImageDrawCircleV - Draw a filled circle within an image (Vector version)
func ImageDrawCircleV(dst *Image, center Vector2, radius int32, col color.RGBA) {
	defer scratch.release(scratch.mark())
	p0 := imageIn(dst)
	cImageDrawCircleV(p0, center.X, center.Y, radius, colorToUint32(col))
	imageBack(dst, p0)
}

// ImageDrawCircleLines - Draw circle outline within an image
func ImageDrawCircleLines(dst *Image, centerX, centerY, radius int32, col color.RGBA) {
	defer scratch.release(scratch.mark())
	p0 := imageIn(dst)
	cImageDrawCircleLines(p0, centerX, centerY, radius, colorToUint32(col))
	imageBack(dst, p0)
}

// ImageDrawCircleLinesV - Draw circle outline within an image (Vector version)
func ImageDrawCircleLinesV(dst *Image, center Vector2, radius int32, col color.RGBA) {
	defer scratch.release(scratch.mark())
	p0 := imageIn(dst)
	cImageDrawCircleLinesV(p0, center.X, center.Y, radius, colorToUint32(col))
	imageBack(dst, p0)
}

// ImageDrawCircleGradient - Draw a gradient-filled circle within an image
func ImageDrawCircleGradient(dst *Image, center Vector2, radius float32, inner, outer color.RGBA) {
	defer scratch.release(scratch.mark())
	p0 := imageIn(dst)
	cImageDrawCircleGradient(p0, center.X, center.Y, radius, colorToUint32(inner), colorToUint32(outer))
	imageBack(dst, p0)
}

// ImageDrawRectangle - Draw rectangle within an image
func ImageDrawRectangle(dst *Image, posX, posY, width, height int32, col color.RGBA) {
	defer scratch.release(scratch.mark())
	p0 := imageIn(dst)
	cImageDrawRectangle(p0, posX, posY, width, height, colorToUint32(col))
	imageBack(dst, p0)
}

// ImageDrawRectangleV - Draw rectangle within an image (Vector version)
func ImageDrawRectangleV(dst *Image, position, size Vector2, col color.RGBA) {
	defer scratch.release(scratch.mark())
	p0 := imageIn(dst)
	cImageDrawRectangleV(p0, position.X, position.Y, size.X, size.Y, colorToUint32(col))
	imageBack(dst, p0)
}

// ImageDrawRectangleRec - Draw rectangle within an image
func ImageDrawRectangleRec(dst *Image, rec Rectangle, col color.RGBA) {
	defer scratch.release(scratch.mark())
	p0 := imageIn(dst)
	cImageDrawRectangleRec(p0, rec.X, rec.Y, rec.Width, rec.Height, colorToUint32(col))
	imageBack(dst, p0)
}

// ImageDrawRectangleLines - Draw rectangle lines within an image
func ImageDrawRectangleLines(dst *Image, posX, posY, width, height int32, col color.RGBA) {
	defer scratch.release(scratch.mark())
	p0 := imageIn(dst)
	cImageDrawRectangleLines(p0, posX, posY, width, height, colorToUint32(col))
	imageBack(dst, p0)
}

// ImageDrawRectangleLinesEx - Draw rectangle lines within an image with line thickness
func ImageDrawRectangleLinesEx(dst *Image, rec Rectangle, thick int, col color.RGBA) {
	defer scratch.release(scratch.mark())
	p0 := imageIn(dst)
	cImageDrawRectangleLinesEx(p0, rec.X, rec.Y, rec.Width, rec.Height, int32(thick), colorToUint32(col))
	imageBack(dst, p0)
}

// ImageDrawRectanglePro - Draw a color-filled rectangle with pro parameters within an image
func ImageDrawRectanglePro(dst *Image, rec Rectangle, origin Vector2, rotation float32, col color.RGBA) {
	defer scratch.release(scratch.mark())
	p0 := imageIn(dst)
	cImageDrawRectanglePro(p0, rec.X, rec.Y, rec.Width, rec.Height, origin.X, origin.Y, rotation, colorToUint32(col))
	imageBack(dst, p0)
}

// ImageDrawRectangleGradientEx - Draw rectangle with gradient colors within an image, counter-clockwise color order
func ImageDrawRectangleGradientEx(dst *Image, rec Rectangle, col1, col2, col3, col4 color.RGBA) {
	defer scratch.release(scratch.mark())
	p0 := imageIn(dst)
	cImageDrawRectangleGradientEx(p0, rec.X, rec.Y, rec.Width, rec.Height, colorToUint32(col1), colorToUint32(col2), colorToUint32(col3), colorToUint32(col4))
	imageBack(dst, p0)
}

// ImageDrawTriangle - Draw triangle within an image
func ImageDrawTriangle(dst *Image, v1, v2, v3 Vector2, col color.RGBA) {
	defer scratch.release(scratch.mark())
	p0 := imageIn(dst)
	cImageDrawTriangle(p0, v1.X, v1.Y, v2.X, v2.Y, v3.X, v3.Y, colorToUint32(col))
	imageBack(dst, p0)
}

// ImageDrawTriangleGradient - Draw triangle with interpolated colors within an image
func ImageDrawTriangleGradient(dst *Image, v1, v2, v3 Vector2, c1, c2, c3 color.RGBA) {
	defer scratch.release(scratch.mark())
	p0 := imageIn(dst)
	cImageDrawTriangleGradient(p0, v1.X, v1.Y, v2.X, v2.Y, v3.X, v3.Y, colorToUint32(c1), colorToUint32(c2), colorToUint32(c3))
	imageBack(dst, p0)
}

// ImageDrawTriangleLines - Draw triangle outline within an image
func ImageDrawTriangleLines(dst *Image, v1, v2, v3 Vector2, col color.RGBA) {
	defer scratch.release(scratch.mark())
	p0 := imageIn(dst)
	cImageDrawTriangleLines(p0, v1.X, v1.Y, v2.X, v2.Y, v3.X, v3.Y, colorToUint32(col))
	imageBack(dst, p0)
}

// ImageDrawTriangleFan - Draw a triangle fan defined by points within an image (first vertex is the center)
func ImageDrawTriangleFan(dst *Image, points []Vector2, col color.RGBA) {
	defer scratch.release(scratch.mark())
	p0 := imageIn(dst)
	cImageDrawTriangleFan(p0, scratchSlice(points), int32(len(points)), colorToUint32(col))
	imageBack(dst, p0)
}

// ImageDrawTriangleStrip - Draw a triangle strip defined by points within an image
func ImageDrawTriangleStrip(dst *Image, points []Vector2, col color.RGBA) {
	defer scratch.release(scratch.mark())
	p0 := imageIn(dst)
	cImageDrawTriangleStrip(p0, scratchSlice(points), int32(len(points)), colorToUint32(col))
	imageBack(dst, p0)
}

// ImageDrawImage - Draw an image within an image
func ImageDrawImage(dst *Image, src *Image, posX, posY int32, tint color.RGBA) {
	defer scratch.release(scratch.mark())
	p0 := imageIn(dst)
	p1 := imageIn(src)
	cImageDrawImage(p0, p1, posX, posY, colorToUint32(tint))
	imageBack(dst, p0)
	imageDone(p1)
}

// ImageDrawImageEx - Draw an image with scaling and rotation within an image
func ImageDrawImageEx(dst *Image, src *Image, position Vector2, rotation, scale float32, tint color.RGBA) {
	defer scratch.release(scratch.mark())
	p0 := imageIn(dst)
	p1 := imageIn(src)
	cImageDrawImageEx(p0, p1, position.X, position.Y, rotation, scale, colorToUint32(tint))
	imageBack(dst, p0)
	imageDone(p1)
}

// ImageDrawImageRec - Draw a part of an image defined by a rectangle within an image
func ImageDrawImageRec(dst *Image, src *Image, srcRec Rectangle, position Vector2, tint color.RGBA) {
	defer scratch.release(scratch.mark())
	p0 := imageIn(dst)
	p1 := imageIn(src)
	cImageDrawImageRec(p0, p1, srcRec.X, srcRec.Y, srcRec.Width, srcRec.Height, position.X, position.Y, colorToUint32(tint))
	imageBack(dst, p0)
	imageDone(p1)
}

// ImageDrawImagePro - Draw a part of an image defined by a rectangle into destination rectangle, with scaling and rotation, within an image
func ImageDrawImagePro(dst *Image, src *Image, srcRec, dstRec Rectangle, origin Vector2, rotation float32, tint color.RGBA) {
	defer scratch.release(scratch.mark())
	p0 := imageIn(dst)
	p1 := imageIn(src)
	cImageDrawImagePro(p0, p1, srcRec.X, srcRec.Y, srcRec.Width, srcRec.Height, dstRec.X, dstRec.Y, dstRec.Width, dstRec.Height, origin.X, origin.Y, rotation, colorToUint32(tint))
	imageBack(dst, p0)
	imageDone(p1)
}

// ImageDrawText - Draw text (using default font) within an image (destination)
func ImageDrawText(dst *Image, posX, posY int32, text string, fontSize int32, col color.RGBA) {
	defer scratch.release(scratch.mark())
	p0 := imageIn(dst)
	cImageDrawText(p0, scratch.str(text), posX, posY, fontSize, colorToUint32(col))
	imageBack(dst, p0)
}

// ImageDrawTextEx - Draw text (custom sprite font) within an image (destination)
func ImageDrawTextEx(dst *Image, position Vector2, font Font, text string, fontSize, spacing float32, tint color.RGBA) {
	defer scratch.release(scratch.mark())
	p0 := imageIn(dst)
	cImageDrawTextEx(p0, fontIn(&font), scratch.str(text), position.X, position.Y, fontSize, spacing, colorToUint32(tint))
	imageBack(dst, p0)
}

// ImageDrawTextPro - Draw text using Font and pro parameters (rotation)
func ImageDrawTextPro(dst *Image, font Font, text string, position, origin Vector2, rotation, fontSize, spacing float32, tint color.RGBA) {
	defer scratch.release(scratch.mark())
	p0 := imageIn(dst)
	cImageDrawTextPro(p0, fontIn(&font), scratch.str(text), position.X, position.Y, origin.X, origin.Y, rotation, fontSize, spacing, colorToUint32(tint))
	imageBack(dst, p0)
}

// LoadTexture - Load texture from file into GPU memory (VRAM)
func LoadTexture(fileName string) Texture2D {
	defer scratch.release(scratch.mark())
	cLoadTexture(scratch.str(fileName))
	return readRet[Texture2D]()
}

// LoadTextureFromImage - Load texture from image data
func LoadTextureFromImage(image *Image) Texture2D {
	defer scratch.release(scratch.mark())
	p0 := imageIn(image)
	cLoadTextureFromImage(p0)
	imageDone(p0)
	return readRet[Texture2D]()
}

// LoadTextureCubemap - Load cubemap from image, multiple image cubemap layouts supported
func LoadTextureCubemap(image *Image, layout int32) Texture2D {
	defer scratch.release(scratch.mark())
	p0 := imageIn(image)
	cLoadTextureCubemap(p0, layout)
	imageDone(p0)
	return readRet[Texture2D]()
}

// LoadRenderTexture - Load texture for rendering (framebuffer)
func LoadRenderTexture(width, height int32) RenderTexture2D {
	cLoadRenderTexture(width, height)
	return readRet[RenderTexture2D]()
}

// LoadRenderTextureEx - Load texture for rendering (framebuffer), with specific format
func LoadRenderTextureEx(width, height int32, format PixelFormat) RenderTexture2D {
	cLoadRenderTextureEx(width, height, int32(format))
	return readRet[RenderTexture2D]()
}

// IsTextureValid - Check if a texture is valid (loaded in GPU)
func IsTextureValid(texture Texture2D) bool {
	return cIsTextureValid(texture.ID, texture.Width, texture.Height, texture.Mipmaps, int32(texture.Format)) != 0
}

// UnloadTexture - Unload texture from GPU memory (VRAM)
func UnloadTexture(texture Texture2D) {
	cUnloadTexture(texture.ID, texture.Width, texture.Height, texture.Mipmaps, int32(texture.Format))
}

// IsRenderTextureValid - Check if a render texture is valid (loaded in GPU)
func IsRenderTextureValid(target RenderTexture2D) bool {
	return cIsRenderTextureValid(target.ID, target.Texture.ID, target.Texture.Width, target.Texture.Height, target.Texture.Mipmaps, int32(target.Texture.Format), target.Depth.ID, target.Depth.Width, target.Depth.Height, target.Depth.Mipmaps, int32(target.Depth.Format)) != 0
}

// UnloadRenderTexture - Unload render texture from GPU memory (VRAM)
func UnloadRenderTexture(target RenderTexture2D) {
	cUnloadRenderTexture(target.ID, target.Texture.ID, target.Texture.Width, target.Texture.Height, target.Texture.Mipmaps, int32(target.Texture.Format), target.Depth.ID, target.Depth.Width, target.Depth.Height, target.Depth.Mipmaps, int32(target.Depth.Format))
}

// UpdateTexture - Update GPU texture with new data
func UpdateTexture[T []color.RGBA | *image.RGBA | []byte](texture Texture2D, pixels T) {
	defer scratch.release(scratch.mark())
	cUpdateTexture(texture.ID, texture.Width, texture.Height, texture.Mipmaps, int32(texture.Format), texturePixels(pixels))
}

// UpdateTextureRec - Update GPU texture rectangle with new data
func UpdateTextureRec[T []color.RGBA | *image.RGBA | []byte](texture Texture2D, rec Rectangle, pixels T) {
	defer scratch.release(scratch.mark())
	cUpdateTextureRec(texture.ID, texture.Width, texture.Height, texture.Mipmaps, int32(texture.Format), rec.X, rec.Y, rec.Width, rec.Height, texturePixels(pixels))
}

// GenTextureMipmaps - Generate GPU mipmaps for a texture
func GenTextureMipmaps(texture *Texture2D) {
	defer scratch.release(scratch.mark())
	p0 := scratchPtr(texture)
	cGenTextureMipmaps(p0)
	readBack(texture, p0)
}

// SetTextureFilter - Set texture scaling filter mode
func SetTextureFilter(texture Texture2D, filter TextureFilterMode) {
	cSetTextureFilter(texture.ID, texture.Width, texture.Height, texture.Mipmaps, int32(texture.Format), int32(filter))
}

// SetTextureWrap - Set texture wrapping mode
func SetTextureWrap(texture Texture2D, wrap TextureWrapMode) {
	cSetTextureWrap(texture.ID, texture.Width, texture.Height, texture.Mipmaps, int32(texture.Format), int32(wrap))
}

// DrawTexture - Draw a Texture2D
func DrawTexture(texture Texture2D, posX int32, posY int32, tint color.RGBA) {
	cDrawTexture(texture.ID, texture.Width, texture.Height, texture.Mipmaps, int32(texture.Format), posX, posY, colorToUint32(tint))
}

// DrawTextureV - Draw a Texture2D with position defined as Vector2
func DrawTextureV(texture Texture2D, position Vector2, tint color.RGBA) {
	cDrawTextureV(texture.ID, texture.Width, texture.Height, texture.Mipmaps, int32(texture.Format), position.X, position.Y, colorToUint32(tint))
}

// DrawTextureEx - Draw a Texture2D with extended parameters
func DrawTextureEx(texture Texture2D, position Vector2, rotation, scale float32, tint color.RGBA) {
	cDrawTextureEx(texture.ID, texture.Width, texture.Height, texture.Mipmaps, int32(texture.Format), position.X, position.Y, rotation, scale, colorToUint32(tint))
}

// DrawTextureRec - Draw a part of a texture defined by a rectangle
func DrawTextureRec(texture Texture2D, source Rectangle, position Vector2, tint color.RGBA) {
	cDrawTextureRec(texture.ID, texture.Width, texture.Height, texture.Mipmaps, int32(texture.Format), source.X, source.Y, source.Width, source.Height, position.X, position.Y, colorToUint32(tint))
}

// DrawTexturePro - Draw a part of a texture defined by a rectangle with 'pro' parameters
func DrawTexturePro(texture Texture2D, source, dest Rectangle, origin Vector2, rotation float32, tint color.RGBA) {
	cDrawTexturePro(texture.ID, texture.Width, texture.Height, texture.Mipmaps, int32(texture.Format), source.X, source.Y, source.Width, source.Height, dest.X, dest.Y, dest.Width, dest.Height, origin.X, origin.Y, rotation, colorToUint32(tint))
}

// DrawTextureNPatch - Draws a texture (or part of it) that stretches or shrinks nicely
func DrawTextureNPatch(texture Texture2D, nPatchInfo NPatchInfo, dest Rectangle, origin Vector2, rotation float32, tint color.RGBA) {
	cDrawTextureNPatch(texture.ID, texture.Width, texture.Height, texture.Mipmaps, int32(texture.Format), nPatchInfo.Source.X, nPatchInfo.Source.Y, nPatchInfo.Source.Width, nPatchInfo.Source.Height, nPatchInfo.Left, nPatchInfo.Top, nPatchInfo.Right, nPatchInfo.Bottom, int32(nPatchInfo.Layout), dest.X, dest.Y, dest.Width, dest.Height, origin.X, origin.Y, rotation, colorToUint32(tint))
}

// Fade - Get color with alpha applied, alpha goes from 0.0f to 1.0f
func Fade(col color.RGBA, alpha float32) color.RGBA {
	return uint32ToColor(cFade(colorToUint32(col), alpha))
}

// ColorToInt - Get hexadecimal value for a Color (0xRRGGBBAA)
func ColorToInt(col color.RGBA) int32 {
	return cColorToInt(colorToUint32(col))
}

// ColorNormalize - Get Color normalized as float [0..1]
func ColorNormalize(col color.RGBA) Vector4 {
	cColorNormalize(colorToUint32(col))
	return readRet[Vector4]()
}

// ColorFromNormalized - Get Color from normalized values [0..1]
func ColorFromNormalized(normalized Vector4) color.RGBA {
	return uint32ToColor(cColorFromNormalized(normalized.X, normalized.Y, normalized.Z, normalized.W))
}

// ColorToHSV - Get HSV values for a Color, hue [0..360], saturation/value [0..1]
func ColorToHSV(col color.RGBA) Vector3 {
	cColorToHSV(colorToUint32(col))
	return readRet[Vector3]()
}

// ColorFromHSV - Get a Color from HSV values, hue [0..360], saturation/value [0..1]
func ColorFromHSV(hue, saturation, value float32) color.RGBA {
	return uint32ToColor(cColorFromHSV(hue, saturation, value))
}

// ColorTint - Get color multiplied with another color
func ColorTint(col color.RGBA, tint color.RGBA) color.RGBA {
	return uint32ToColor(cColorTint(colorToUint32(col), colorToUint32(tint)))
}

// ColorBrightness - Get color with brightness correction, brightness factor goes from -1.0f to 1.0f
func ColorBrightness(col color.RGBA, factor float32) color.RGBA {
	return uint32ToColor(cColorBrightness(colorToUint32(col), factor))
}

// ColorContrast - Get color with contrast correction, contrast values between -1.0f and 1.0f
func ColorContrast(col color.RGBA, contrast float32) color.RGBA {
	return uint32ToColor(cColorContrast(colorToUint32(col), contrast))
}

// ColorAlpha - Get color with alpha applied, alpha goes from 0.0f to 1.0f
func ColorAlpha(col color.RGBA, alpha float32) color.RGBA {
	return uint32ToColor(cColorAlpha(colorToUint32(col), alpha))
}

// ColorAlphaBlend - Get src alpha-blended into dst color with tint
func ColorAlphaBlend(dst, src, tint color.RGBA) color.RGBA {
	return uint32ToColor(cColorAlphaBlend(colorToUint32(dst), colorToUint32(src), colorToUint32(tint)))
}

// ColorLerp - Get color lerp interpolation between two colors, factor [0.0f..1.0f]
func ColorLerp(col1, col2 color.RGBA, factor float32) color.RGBA {
	return uint32ToColor(cColorLerp(colorToUint32(col1), colorToUint32(col2), factor))
}

// GetColor - Get Color structure from hexadecimal value
func GetColor(hexValue uint) color.RGBA {
	return uint32ToColor(cGetColor(uint32(hexValue)))
}

// GetPixelColor - Get Color from a source pixel pointer of certain format
func GetPixelColor(srcPtr unsafe.Pointer, format int32) color.RGBA {
	defer scratch.release(scratch.mark())
	return uint32ToColor(cGetPixelColor(scratch.bytes(srcPtr, uint32(cGetPixelDataSize(1, 1, format))), format))
}

// SetPixelColor - Set color formatted into destination pixel pointer
func SetPixelColor(dstPtr unsafe.Pointer, col color.RGBA, format int32) {
	defer scratch.release(scratch.mark())
	n := uint32(cGetPixelDataSize(1, 1, format))
	p := scratch.bytes(dstPtr, n)
	cSetPixelColor(p, colorToUint32(col), format)
	if p != 0 {
		copyOut(dstPtr, p, n)
	}
}

// GetPixelDataSize - Get pixel data size in bytes for certain format
func GetPixelDataSize(width, height, format int32) int32 {
	return cGetPixelDataSize(width, height, format)
}

// GetFontDefault - Get the default Font
func GetFontDefault() Font {
	cGetFontDefault()
	return fontFromC(readRet[cFont]())
}

// LoadFont - Load font from file into GPU memory (VRAM)
func LoadFont(fileName string) Font {
	defer scratch.release(scratch.mark())
	cLoadFont(scratch.str(fileName))
	return fontFromC(readRet[cFont]())
}

// LoadFontEx - Load Font from file with extended parameters
func LoadFontEx(fileName string, fontSize int32, fontChars []rune, runesNumber ...int32) Font {
	defer scratch.release(scratch.mark())
	var charsCount int32
	var chars uint32
	if fontChars != nil {
		chars = scratchSlice(fontChars)
		charsCount = int32(len(fontChars))
	} else if len(runesNumber) > 0 {
		charsCount = runesNumber[0]
	}
	cLoadFontEx(scratch.str(fileName), fontSize, chars, charsCount)
	return fontFromC(readRet[cFont]())
}

// LoadFontFromImage - Load font from Image (XNA style)
func LoadFontFromImage(image Image, key color.RGBA, firstChar rune) Font {
	defer scratch.release(scratch.mark())
	p0 := imageIn(&image)
	cLoadFontFromImage(p0, colorToUint32(key), int32(firstChar))
	imageDone(p0)
	return fontFromC(readRet[cFont]())
}

// LoadFontFromMemory - Load font from memory buffer, fileType refers to extension: i.e. '.ttf'
func LoadFontFromMemory(fileType string, fileData []byte, fontSize int32, codepoints []rune) Font {
	defer scratch.release(scratch.mark())
	cLoadFontFromMemory(scratch.str(fileType), scratchSlice(fileData), int32(len(fileData)), fontSize, scratchSlice(codepoints), int32(len(codepoints)))
	return fontFromC(readRet[cFont]())
}

// IsFontValid - Check if a font is valid (font data loaded, WARNING: GPU texture not checked)
func IsFontValid(font Font) bool {
	defer scratch.release(scratch.mark())
	return cIsFontValid(fontIn(&font)) != 0
}

// LoadFontData - Load font data for further use
func LoadFontData(fileData []byte, fontSize int32, codepoints []rune, codepointCount, typ int32) []GlyphInfo {
	defer scratch.release(scratch.mark())
	if codepointCount < 0 {
		codepointCount = 0
	}
	var glyphCount int32
	count := scratchValue(&glyphCount)
	ret := cLoadFontData(scratchSlice(fileData), int32(len(fileData)), fontSize, scratchSlice(codepoints), codepointCount, typ, count)
	glyphCount = readC[int32](count)
	if ret == 0 || glyphCount <= 0 {
		return nil
	}
	glyphs := glyphsFromC(ret, int(glyphCount))
	cUnloadFontData(ret, glyphCount)
	return glyphs
}

// GenImageFontAtlas - Generate image font atlas using chars info
func GenImageFontAtlas(glyphs []GlyphInfo, glyphRecs []*Rectangle, fontSize int32, padding int32, packMethod int32) Image {
	defer scratch.release(scratch.mark())
	recs := scratch.alloc(4)
	cGenImageFontAtlas(glyphsToC(glyphs), recs, int32(len(glyphs)), fontSize, padding, packMethod)
	p := readC[uint32](recs)
	if len(glyphRecs) > 0 {
		glyphRecs[0] = first(goSlice[Rectangle](p, len(glyphs)))
	}
	cFree(p)
	return *imageRet()
}

// UnloadFontData - Unload font chars info data (RAM)
func UnloadFontData(glyphs []GlyphInfo) {
}

// UnloadFont - Unload Font from GPU memory (VRAM)
func UnloadFont(font Font) {
	e := fonts[font.Texture.ID]
	if e == nil || font.Texture.ID == GetFontDefault().Texture.ID {
		return
	}
	cUnloadFont(fontIn(&font))
	cFree(e.at)
	delete(fonts, font.Texture.ID)
}

// DrawFPS - Draw current FPS
func DrawFPS(posX int32, posY int32) {
	cDrawFPS(posX, posY)
}

// DrawText - Draw text (using default font)
func DrawText(text string, posX int32, posY int32, fontSize int32, col color.RGBA) {
	defer scratch.release(scratch.mark())
	cDrawText(scratch.str(text), posX, posY, fontSize, colorToUint32(col))
}

// DrawTextEx - Draw text using font and additional parameters
func DrawTextEx(font Font, text string, position Vector2, fontSize float32, spacing float32, tint color.RGBA) {
	defer scratch.release(scratch.mark())
	cDrawTextEx(fontIn(&font), scratch.str(text), position.X, position.Y, fontSize, spacing, colorToUint32(tint))
}

// DrawTextPro - Draw text using Font and pro parameters (rotation)
func DrawTextPro(font Font, text string, position, origin Vector2, rotation, fontSize, spacing float32, tint color.RGBA) {
	defer scratch.release(scratch.mark())
	cDrawTextPro(fontIn(&font), scratch.str(text), position.X, position.Y, origin.X, origin.Y, rotation, fontSize, spacing, colorToUint32(tint))
}

// DrawTextCodepoint - Draw one character (codepoint)
func DrawTextCodepoint(font Font, codepoint rune, position Vector2, fontSize float32, tint color.RGBA) {
	defer scratch.release(scratch.mark())
	cDrawTextCodepoint(fontIn(&font), int32(codepoint), position.X, position.Y, fontSize, colorToUint32(tint))
}

// DrawTextCodepoints - Draw multiple character (codepoint)
func DrawTextCodepoints(font Font, codepoints []rune, position Vector2, fontSize float32, spacing float32, tint color.RGBA) {
	defer scratch.release(scratch.mark())
	cDrawTextCodepoints(fontIn(&font), scratchSlice(codepoints), int32(len(codepoints)), position.X, position.Y, fontSize, spacing, colorToUint32(tint))
}

// SetTextLineSpacing - Set vertical line spacing when drawing with line-breaks
func SetTextLineSpacing(spacing int) {
	cSetTextLineSpacing(int32(spacing))
}

// MeasureText - Measure string width for default font
func MeasureText(text string, fontSize int32) int32 {
	defer scratch.release(scratch.mark())
	return cMeasureText(scratch.str(text), fontSize)
}

// MeasureTextEx - Measure string size for Font
func MeasureTextEx(font Font, text string, fontSize float32, spacing float32) Vector2 {
	defer scratch.release(scratch.mark())
	ret := cMeasureTextEx(fontIn(&font), scratch.str(text), fontSize, spacing)
	return *(*Vector2)(unsafe.Pointer(&ret))
}

// MeasureTextCodepoints - Measure string size for an existing array of codepoints for Font
func MeasureTextCodepoints(font Font, codepoints []rune, fontSize float32, spacing float32) Vector2 {
	defer scratch.release(scratch.mark())
	ret := cMeasureTextCodepoints(fontIn(&font), scratchSlice(codepoints), int32(len(codepoints)), fontSize, spacing)
	return *(*Vector2)(unsafe.Pointer(&ret))
}

// GetGlyphIndex - Get glyph index position in font for a codepoint (unicode character), fallback to '?' if not found
func GetGlyphIndex(font Font, codepoint rune) int32 {
	defer scratch.release(scratch.mark())
	return cGetGlyphIndex(fontIn(&font), int32(codepoint))
}

// GetGlyphInfo - Get glyph font info data for a codepoint (unicode character), fallback to '?' if not found
func GetGlyphInfo(font Font, codepoint rune) GlyphInfo {
	glyphs := ptrSlice(font.Chars, int(font.CharsCount))
	if len(glyphs) == 0 {
		return GlyphInfo{}
	}
	return glyphs[GetGlyphIndex(font, codepoint)]
}

// GetGlyphAtlasRec - Get glyph rectangle in font atlas for a codepoint (unicode character), fallback to '?' if not found
func GetGlyphAtlasRec(font Font, codepoint rune) Rectangle {
	defer scratch.release(scratch.mark())
	cGetGlyphAtlasRec(fontIn(&font), int32(codepoint))
	return readRet[Rectangle]()
}

// DrawLine3D - Draw a line in 3D world space
func DrawLine3D(startPos Vector3, endPos Vector3, col color.RGBA) {
	cDrawLine3D(startPos.X, startPos.Y, startPos.Z, endPos.X, endPos.Y, endPos.Z, colorToUint32(col))
}

// DrawPoint3D - Draw a point in 3D space, actually a small line
func DrawPoint3D(position Vector3, col color.RGBA) {
	cDrawPoint3D(position.X, position.Y, position.Z, colorToUint32(col))
}

// DrawCircle3D - Draw a circle in 3D world space
func DrawCircle3D(center Vector3, radius float32, rotationAxis Vector3, rotationAngle float32, col color.RGBA) {
	cDrawCircle3D(center.X, center.Y, center.Z, radius, rotationAxis.X, rotationAxis.Y, rotationAxis.Z, rotationAngle, colorToUint32(col))
}

// DrawTriangle3D - Draw a color-filled triangle (vertex in counter-clockwise order!)
func DrawTriangle3D(v1 Vector3, v2 Vector3, v3 Vector3, col color.RGBA) {
	cDrawTriangle3D(v1.X, v1.Y, v1.Z, v2.X, v2.Y, v2.Z, v3.X, v3.Y, v3.Z, colorToUint32(col))
}

// DrawTriangleStrip3D - Draw a triangle strip defined by points
func DrawTriangleStrip3D(points []Vector3, col color.RGBA) {
	defer scratch.release(scratch.mark())
	cDrawTriangleStrip3D(scratchSlice(points), int32(len(points)), colorToUint32(col))
}

// DrawCube - Draw cube
func DrawCube(position Vector3, width float32, height float32, length float32, col color.RGBA) {
	cDrawCube(position.X, position.Y, position.Z, width, height, length, colorToUint32(col))
}

// DrawCubeV - Draw cube (Vector version)
func DrawCubeV(position Vector3, size Vector3, col color.RGBA) {
	cDrawCubeV(position.X, position.Y, position.Z, size.X, size.Y, size.Z, colorToUint32(col))
}

// DrawCubeWires - Draw cube wires
func DrawCubeWires(position Vector3, width float32, height float32, length float32, col color.RGBA) {
	cDrawCubeWires(position.X, position.Y, position.Z, width, height, length, colorToUint32(col))
}

// DrawCubeWiresV - Draw cube wires (Vector version)
func DrawCubeWiresV(position Vector3, size Vector3, col color.RGBA) {
	cDrawCubeWiresV(position.X, position.Y, position.Z, size.X, size.Y, size.Z, colorToUint32(col))
}

// DrawSphere - Draw sphere
func DrawSphere(centerPos Vector3, radius float32, col color.RGBA) {
	cDrawSphere(centerPos.X, centerPos.Y, centerPos.Z, radius, colorToUint32(col))
}

// DrawSphereEx - Draw sphere with extended parameters
func DrawSphereEx(centerPos Vector3, radius float32, rings int32, slices int32, col color.RGBA) {
	cDrawSphereEx(centerPos.X, centerPos.Y, centerPos.Z, radius, rings, slices, colorToUint32(col))
}

// DrawSphereWires - Draw sphere wires
func DrawSphereWires(centerPos Vector3, radius float32, rings int32, slices int32, col color.RGBA) {
	cDrawSphereWires(centerPos.X, centerPos.Y, centerPos.Z, radius, rings, slices, colorToUint32(col))
}

// DrawCylinder - Draw a cylinder/cone
func DrawCylinder(position Vector3, radiusTop float32, radiusBottom float32, height float32, slices int32, col color.RGBA) {
	cDrawCylinder(position.X, position.Y, position.Z, radiusTop, radiusBottom, height, slices, colorToUint32(col))
}

// DrawCylinderEx - Draw a cylinder with base at startPos and top at endPos
func DrawCylinderEx(startPos Vector3, endPos Vector3, startRadius float32, endRadius float32, sides int32, col color.RGBA) {
	cDrawCylinderEx(startPos.X, startPos.Y, startPos.Z, endPos.X, endPos.Y, endPos.Z, startRadius, endRadius, sides, colorToUint32(col))
}

// DrawCylinderWires - Draw a cylinder/cone wires
func DrawCylinderWires(position Vector3, radiusTop float32, radiusBottom float32, height float32, slices int32, col color.RGBA) {
	cDrawCylinderWires(position.X, position.Y, position.Z, radiusTop, radiusBottom, height, slices, colorToUint32(col))
}

// DrawCylinderWiresEx - Draw a cylinder wires with base at startPos and top at endPos
func DrawCylinderWiresEx(startPos Vector3, endPos Vector3, startRadius float32, endRadius float32, sides int32, col color.RGBA) {
	cDrawCylinderWiresEx(startPos.X, startPos.Y, startPos.Z, endPos.X, endPos.Y, endPos.Z, startRadius, endRadius, sides, colorToUint32(col))
}

// DrawCapsule - Draw a capsule with the center of its sphere caps at startPos and endPos
func DrawCapsule(startPos, endPos Vector3, radius float32, rings, slices int32, col color.RGBA) {
	cDrawCapsule(startPos.X, startPos.Y, startPos.Z, endPos.X, endPos.Y, endPos.Z, radius, rings, slices, colorToUint32(col))
}

// DrawCapsuleWires - Draw capsule wireframe with the center of its sphere caps at startPos and endPos
func DrawCapsuleWires(startPos, endPos Vector3, radius float32, rings, slices int32, col color.RGBA) {
	cDrawCapsuleWires(startPos.X, startPos.Y, startPos.Z, endPos.X, endPos.Y, endPos.Z, radius, rings, slices, colorToUint32(col))
}

// DrawPlane - Draw a plane XZ
func DrawPlane(centerPos Vector3, size Vector2, col color.RGBA) {
	cDrawPlane(centerPos.X, centerPos.Y, centerPos.Z, size.X, size.Y, colorToUint32(col))
}

// DrawRay - Draw a ray line
func DrawRay(ray Ray, col color.RGBA) {
	cDrawRay(ray.Position.X, ray.Position.Y, ray.Position.Z, ray.Direction.X, ray.Direction.Y, ray.Direction.Z, colorToUint32(col))
}

// DrawGrid - Draw a grid (centered at (0, 0, 0))
func DrawGrid(slices int32, spacing float32) {
	cDrawGrid(slices, spacing)
}

// LoadModel - Load model from files (meshes and materials)
func LoadModel(fileName string) Model {
	defer scratch.release(scratch.mark())
	cLoadModel(scratch.str(fileName))
	return modelFromC(readRet[cModel]())
}

// LoadModelFromMesh - Load model from generated mesh (default material)
func LoadModelFromMesh(mesh Mesh) Model {
	defer scratch.release(scratch.mark())
	cLoadModelFromMesh(meshIn(&mesh))
	return modelFromC(readRet[cModel]())
}

// IsModelValid - Check if a model is valid (loaded in GPU, VAO/VBOs)
func IsModelValid(model Model) bool {
	defer scratch.release(scratch.mark())
	return cIsModelValid(modelIn(&model)) != 0
}

// UnloadModel - Unload model (including meshes) from memory (RAM and/or VRAM)
func UnloadModel(model Model) {
	defer scratch.release(scratch.mark())
	if models[model.Meshes] == nil {
		return
	}
	cUnloadModel(modelIn(&model))
	forgetModel(&model)
}

// GetModelBoundingBox - Compute model bounding box limits (considers all meshes)
func GetModelBoundingBox(model Model) BoundingBox {
	defer scratch.release(scratch.mark())
	cGetModelBoundingBox(modelIn(&model))
	return readRet[BoundingBox]()
}

// DrawModel - Draw a model (with texture if set)
func DrawModel(model Model, position Vector3, scale float32, tint color.RGBA) {
	defer scratch.release(scratch.mark())
	cDrawModel(modelIn(&model), position.X, position.Y, position.Z, scale, colorToUint32(tint))
}

// DrawModelEx - Draw a model with extended parameters
func DrawModelEx(model Model, position Vector3, rotationAxis Vector3, rotationAngle float32, scale Vector3, tint color.RGBA) {
	defer scratch.release(scratch.mark())
	cDrawModelEx(modelIn(&model), position.X, position.Y, position.Z, rotationAxis.X, rotationAxis.Y, rotationAxis.Z, rotationAngle, scale.X, scale.Y, scale.Z, colorToUint32(tint))
}

// DrawModelWires - Draw a model wires (with texture if set)
func DrawModelWires(model Model, position Vector3, scale float32, tint color.RGBA) {
	defer scratch.release(scratch.mark())
	cDrawModelWires(modelIn(&model), position.X, position.Y, position.Z, scale, colorToUint32(tint))
}

// DrawModelWiresEx - Draw a model wires (with texture if set) with extended parameters
func DrawModelWiresEx(model Model, position Vector3, rotationAxis Vector3, rotationAngle float32, scale Vector3, tint color.RGBA) {
	defer scratch.release(scratch.mark())
	cDrawModelWiresEx(modelIn(&model), position.X, position.Y, position.Z, rotationAxis.X, rotationAxis.Y, rotationAxis.Z, rotationAngle, scale.X, scale.Y, scale.Z, colorToUint32(tint))
}

// DrawBoundingBox - Draw bounding box (wires)
func DrawBoundingBox(box BoundingBox, col color.RGBA) {
	cDrawBoundingBox(box.Min.X, box.Min.Y, box.Min.Z, box.Max.X, box.Max.Y, box.Max.Z, colorToUint32(col))
}

// DrawBillboard - Draw a billboard texture
func DrawBillboard(camera Camera, texture Texture2D, position Vector3, scale float32, tint color.RGBA) {
	cDrawBillboard(camera.Position.X, camera.Position.Y, camera.Position.Z, camera.Target.X, camera.Target.Y, camera.Target.Z, camera.Up.X, camera.Up.Y, camera.Up.Z, camera.Fovy, int32(camera.Projection), texture.ID, texture.Width, texture.Height, texture.Mipmaps, int32(texture.Format), position.X, position.Y, position.Z, scale, colorToUint32(tint))
}

// DrawBillboardRec - Draw a billboard texture defined by source
func DrawBillboardRec(camera Camera, texture Texture2D, source Rectangle, position Vector3, size Vector2, tint color.RGBA) {
	cDrawBillboardRec(camera.Position.X, camera.Position.Y, camera.Position.Z, camera.Target.X, camera.Target.Y, camera.Target.Z, camera.Up.X, camera.Up.Y, camera.Up.Z, camera.Fovy, int32(camera.Projection), texture.ID, texture.Width, texture.Height, texture.Mipmaps, int32(texture.Format), source.X, source.Y, source.Width, source.Height, position.X, position.Y, position.Z, size.X, size.Y, colorToUint32(tint))
}

// DrawBillboardPro - Draw a billboard texture defined by source and rotation
func DrawBillboardPro(camera Camera, texture Texture2D, source Rectangle, position Vector3, up Vector3, size Vector2, origin Vector2, rotation float32, tint color.RGBA) {
	cDrawBillboardPro(camera.Position.X, camera.Position.Y, camera.Position.Z, camera.Target.X, camera.Target.Y, camera.Target.Z, camera.Up.X, camera.Up.Y, camera.Up.Z, camera.Fovy, int32(camera.Projection), texture.ID, texture.Width, texture.Height, texture.Mipmaps, int32(texture.Format), source.X, source.Y, source.Width, source.Height, position.X, position.Y, position.Z, up.X, up.Y, up.Z, size.X, size.Y, origin.X, origin.Y, rotation, colorToUint32(tint))
}

// UploadMesh - Upload vertex data into a VAO (if supported) and VBO
func UploadMesh(mesh *Mesh, dynamic bool) {
	if e, ok := meshes[mesh.VboID]; ok && mesh.VboID != nil {
		cUploadMesh(e.at, boolToInt32(dynamic))
		return
	}
	c := meshToC(mesh, true)
	at := persist(&c)
	cUploadMesh(at, boolToInt32(dynamic))
	c = readC[cMesh](at)
	mesh.VaoID = c.VaoId
	mesh.VboID = first(goSlice[uint32](c.VboId, maxMeshVertexBuffers))
	if mesh.VboID != nil {
		meshes[mesh.VboID] = meshEntry{at, true}
	}
}

// UpdateMeshBuffer - Update mesh vertex data in GPU for a specific buffer index
func UpdateMeshBuffer(mesh Mesh, index int, data []byte, offset int) {
	defer scratch.release(scratch.mark())
	cUpdateMeshBuffer(meshIn(&mesh), int32(index), scratchSlice(data), int32(len(data)), int32(offset))
}

// UnloadMesh - Unload mesh from memory (RAM and/or VRAM)
func UnloadMesh(mesh *Mesh) {
	e, ok := meshes[mesh.VboID]
	if !ok || mesh.VboID == nil {
		return
	}
	cUnloadMesh(e.at)
	if e.owned {
		cFree(e.at)
	}
	delete(meshes, mesh.VboID)
}

// DrawMesh - Draw a 3d mesh with material and transform
func DrawMesh(mesh Mesh, material Material, transform Matrix) {
	defer scratch.release(scratch.mark())
	cDrawMesh(meshIn(&mesh), materialIn(&material), transform.M0, transform.M4, transform.M8, transform.M12, transform.M1, transform.M5, transform.M9, transform.M13, transform.M2, transform.M6, transform.M10, transform.M14, transform.M3, transform.M7, transform.M11, transform.M15)
}

// DrawMeshInstanced - Draw multiple mesh instances with material and different transforms
func DrawMeshInstanced(mesh Mesh, material Material, transforms []Matrix, instances int) {
	defer scratch.release(scratch.mark())
	cDrawMeshInstanced(meshIn(&mesh), materialIn(&material), scratchSlice(transforms), int32(instances))
}

// GetMeshBoundingBox - Compute mesh bounding box limits
func GetMeshBoundingBox(mesh Mesh) BoundingBox {
	defer scratch.release(scratch.mark())
	cGetMeshBoundingBox(meshIn(&mesh))
	return readRet[BoundingBox]()
}

// GenMeshTangents - Compute mesh tangents
func GenMeshTangents(mesh *Mesh) {
	defer scratch.release(scratch.mark())
	_, registered := meshes[mesh.VboID]
	registered = registered && mesh.VboID != nil
	allocated := mesh.Tangents == nil
	at := meshIn(mesh)
	cGenMeshTangents(at)
	c := readC[cMesh](at)
	if mesh.Tangents == nil {
		mesh.Tangents = first(goSlice[float32](c.Tangents, int(c.VertexCount)*4))
	} else {
		copyOut(unsafe.Pointer(mesh.Tangents), c.Tangents, uint32(c.VertexCount)*16)
	}
	if registered {
		mesh.VaoID = c.VaoId
		copyOut(unsafe.Pointer(mesh.VboID), c.VboId, maxMeshVertexBuffers*4)
	} else if allocated {
		cFree(c.Tangents)
	}
}

// ExportMesh - Export mesh data to file, returns true on success
func ExportMesh(mesh Mesh, fileName string) bool {
	defer scratch.release(scratch.mark())
	return cExportMesh(meshIn(&mesh), scratch.str(fileName)) != 0
}

// GenMeshPoly - Generate polygonal mesh
func GenMeshPoly(sides int, radius float32) Mesh {
	cGenMeshPoly(int32(sides), radius)
	return meshRet()
}

// GenMeshPlane - Generate plane mesh (with subdivisions)
func GenMeshPlane(width, length float32, resX, resZ int) Mesh {
	cGenMeshPlane(width, length, int32(resX), int32(resZ))
	return meshRet()
}

// GenMeshCube - Generate cuboid mesh
func GenMeshCube(width, height, length float32) Mesh {
	cGenMeshCube(width, height, length)
	return meshRet()
}

// GenMeshSphere - Generate sphere mesh (standard sphere)
func GenMeshSphere(radius float32, rings, slices int) Mesh {
	cGenMeshSphere(radius, int32(rings), int32(slices))
	return meshRet()
}

// GenMeshHemiSphere - Generate half-sphere mesh (no bottom cap)
func GenMeshHemiSphere(radius float32, rings, slices int) Mesh {
	cGenMeshHemiSphere(radius, int32(rings), int32(slices))
	return meshRet()
}

// GenMeshCylinder - Generate cylinder mesh
func GenMeshCylinder(radius, height float32, slices int) Mesh {
	cGenMeshCylinder(radius, height, int32(slices))
	return meshRet()
}

// GenMeshCone - Generate cone/pyramid mesh
func GenMeshCone(radius, height float32, slices int) Mesh {
	cGenMeshCone(radius, height, int32(slices))
	return meshRet()
}

// GenMeshTorus - Generate torus mesh
func GenMeshTorus(radius, size float32, radSeg, sides int) Mesh {
	cGenMeshTorus(radius, size, int32(radSeg), int32(sides))
	return meshRet()
}

// GenMeshKnot - Generate trefoil knot mesh
func GenMeshKnot(radius, size float32, radSeg, sides int) Mesh {
	cGenMeshKnot(radius, size, int32(radSeg), int32(sides))
	return meshRet()
}

// GenMeshHeightmap - Generate heightmap mesh from image data
func GenMeshHeightmap(heightmap Image, size Vector3) Mesh {
	defer scratch.release(scratch.mark())
	p0 := imageIn(&heightmap)
	cGenMeshHeightmap(p0, size.X, size.Y, size.Z)
	imageDone(p0)
	return meshRet()
}

// GenMeshCubicmap - Generate cubes-based map mesh from image data
func GenMeshCubicmap(cubicmap Image, cubeSize Vector3) Mesh {
	defer scratch.release(scratch.mark())
	p0 := imageIn(&cubicmap)
	cGenMeshCubicmap(p0, cubeSize.X, cubeSize.Y, cubeSize.Z)
	imageDone(p0)
	return meshRet()
}

// LoadMaterials - Load materials from model file
func LoadMaterials(fileName string) []Material {
	defer scratch.release(scratch.mark())
	var materialCount int32
	count := scratchValue(&materialCount)
	ret := cLoadMaterials(scratch.str(fileName), count)
	cms := goSlice[cMaterial](ret, int(readC[int32](count)))
	mats := make([]Material, len(cms))
	for i := range cms {
		mats[i] = materialFromC(cms[i])
		cFree(cms[i].Maps)
	}
	cFree(ret)
	return mats
}

// LoadMaterialDefault - Load default material (Supports: DIFFUSE, SPECULAR, NORMAL maps)
func LoadMaterialDefault() Material {
	cLoadMaterialDefault()
	return materialRet()
}

// IsMaterialValid - Check if a material is valid (shader assigned, map textures loaded in GPU)
func IsMaterialValid(material Material) bool {
	defer scratch.release(scratch.mark())
	return cIsMaterialValid(materialIn(&material)) != 0
}

// UnloadMaterial - Unload material textures from VRAM
func UnloadMaterial(material Material) {
	defer scratch.release(scratch.mark())
	maps := uint32(0)
	if material.Maps != nil {
		maps = cMalloc(MaxMaterialMaps * uint32(unsafe.Sizeof(MaterialMap{})))
	}
	c := materialToC(&material, maps)
	cUnloadMaterial(scratchValue(&c))
	if material.Shader.ID != crlGetShaderIdDefault() {
		delete(shaderLocs, material.Shader.ID)
	}
}

// SetMaterialTexture - Set texture for a material map type (MAP_DIFFUSE, MAP_SPECULAR...)
func SetMaterialTexture(material *Material, mapType int32, texture Texture2D) {
	defer scratch.release(scratch.mark())
	c := materialToC(material, 0)
	cSetMaterialTexture(scratchValue(&c), mapType, texture.ID, texture.Width, texture.Height, texture.Mipmaps, int32(texture.Format))
	if material.Maps != nil {
		copyOut(unsafe.Pointer(material.Maps), c.Maps, MaxMaterialMaps*uint32(unsafe.Sizeof(MaterialMap{})))
	}
}

// SetModelMeshMaterial - Set material for a mesh
func SetModelMeshMaterial(model *Model, meshId int32, materialId int32) {
	defer scratch.release(scratch.mark())
	at := modelIn(model)
	cSetModelMeshMaterial(at, meshId, materialId)
	if model.MeshMaterial != nil {
		copyOut(unsafe.Pointer(model.MeshMaterial), readC[cModel](at).MeshMaterial, uint32(model.MeshCount)*4)
	}
}

// LoadModelAnimations - Load model animations from file
func LoadModelAnimations(fileName string) []ModelAnimation {
	defer scratch.release(scratch.mark())
	var animCount int32
	count := scratchValue(&animCount)
	ret := cLoadModelAnimations(scratch.str(fileName), count)
	if ret == 0 {
		return nil
	}
	return animsFromC(ret, readC[int32](count))
}

// UpdateModelAnimation - Update model animation pose (CPU)
func UpdateModelAnimation(model Model, anim ModelAnimation, frame float32) {
	defer scratch.release(scratch.mark())
	cUpdateModelAnimation(modelIn(&model), animIn(&anim), frame)
	modelPoseBack(&model)
}

// UpdateModelAnimationEx - Update model animation pose, blending two animations
func UpdateModelAnimationEx(model Model, animA ModelAnimation, frameA float32, animB ModelAnimation, frameB, blend float32) {
	defer scratch.release(scratch.mark())
	cUpdateModelAnimationEx(modelIn(&model), animIn(&animA), frameA, animIn(&animB), frameB, blend)
	modelPoseBack(&model)
}

// UnloadModelAnimations - Unload animation array data
func UnloadModelAnimations(animations []ModelAnimation) {
	if len(animations) == 0 {
		return
	}
	e, ok := anims[animations[0].KeyframePoses]
	if !ok || animations[0].KeyframePoses == nil {
		return
	}
	cUnloadModelAnimations(e.base, e.count)
	for k, v := range anims {
		if v.base == e.base {
			delete(anims, k)
		}
	}
}

// IsModelAnimationValid - Check model animation skeleton match
func IsModelAnimationValid(model Model, anim ModelAnimation) bool {
	defer scratch.release(scratch.mark())
	return cIsModelAnimationValid(modelIn(&model), animIn(&anim)) != 0
}

// CheckCollisionSpheres - Check collision between two spheres
func CheckCollisionSpheres(center1 Vector3, radius1 float32, center2 Vector3, radius2 float32) bool {
	return cCheckCollisionSpheres(center1.X, center1.Y, center1.Z, radius1, center2.X, center2.Y, center2.Z, radius2) != 0
}

// CheckCollisionBoxes - Check collision between two bounding boxes
func CheckCollisionBoxes(box1 BoundingBox, box2 BoundingBox) bool {
	return cCheckCollisionBoxes(box1.Min.X, box1.Min.Y, box1.Min.Z, box1.Max.X, box1.Max.Y, box1.Max.Z, box2.Min.X, box2.Min.Y, box2.Min.Z, box2.Max.X, box2.Max.Y, box2.Max.Z) != 0
}

// CheckCollisionBoxSphere - Check collision between box and sphere
func CheckCollisionBoxSphere(box BoundingBox, center Vector3, radius float32) bool {
	return cCheckCollisionBoxSphere(box.Min.X, box.Min.Y, box.Min.Z, box.Max.X, box.Max.Y, box.Max.Z, center.X, center.Y, center.Z, radius) != 0
}

// GetRayCollisionSphere - Get collision info between ray and sphere
func GetRayCollisionSphere(ray Ray, center Vector3, radius float32) RayCollision {
	cGetRayCollisionSphere(ray.Position.X, ray.Position.Y, ray.Position.Z, ray.Direction.X, ray.Direction.Y, ray.Direction.Z, center.X, center.Y, center.Z, radius)
	return readRet[RayCollision]()
}

// GetRayCollisionBox - Get collision info between ray and box
func GetRayCollisionBox(ray Ray, box BoundingBox) RayCollision {
	cGetRayCollisionBox(ray.Position.X, ray.Position.Y, ray.Position.Z, ray.Direction.X, ray.Direction.Y, ray.Direction.Z, box.Min.X, box.Min.Y, box.Min.Z, box.Max.X, box.Max.Y, box.Max.Z)
	return readRet[RayCollision]()
}

// GetRayCollisionMesh - Get collision info between ray and mesh
func GetRayCollisionMesh(ray Ray, mesh Mesh, transform Matrix) RayCollision {
	defer scratch.release(scratch.mark())
	cGetRayCollisionMesh(ray.Position.X, ray.Position.Y, ray.Position.Z, ray.Direction.X, ray.Direction.Y, ray.Direction.Z, meshIn(&mesh), transform.M0, transform.M4, transform.M8, transform.M12, transform.M1, transform.M5, transform.M9, transform.M13, transform.M2, transform.M6, transform.M10, transform.M14, transform.M3, transform.M7, transform.M11, transform.M15)
	return readRet[RayCollision]()
}

// GetRayCollisionTriangle - Get collision info between ray and triangle
func GetRayCollisionTriangle(ray Ray, p1, p2, p3 Vector3) RayCollision {
	cGetRayCollisionTriangle(ray.Position.X, ray.Position.Y, ray.Position.Z, ray.Direction.X, ray.Direction.Y, ray.Direction.Z, p1.X, p1.Y, p1.Z, p2.X, p2.Y, p2.Z, p3.X, p3.Y, p3.Z)
	return readRet[RayCollision]()
}

// GetRayCollisionQuad - Get collision info between ray and quad
func GetRayCollisionQuad(ray Ray, p1, p2, p3, p4 Vector3) RayCollision {
	cGetRayCollisionQuad(ray.Position.X, ray.Position.Y, ray.Position.Z, ray.Direction.X, ray.Direction.Y, ray.Direction.Z, p1.X, p1.Y, p1.Z, p2.X, p2.Y, p2.Z, p3.X, p3.Y, p3.Z, p4.X, p4.Y, p4.Z)
	return readRet[RayCollision]()
}

// InitAudioDevice - Initialize audio device and context
func InitAudioDevice() {
	cInitAudioDevice()
}

// CloseAudioDevice - Close the audio device and context
func CloseAudioDevice() {
	cCloseAudioDevice()
}

// IsAudioDeviceReady - Check if audio device has been initialized successfully
func IsAudioDeviceReady() bool {
	return cIsAudioDeviceReady() != 0
}

// SetMasterVolume - Set master volume (listener)
func SetMasterVolume(volume float32) {
	cSetMasterVolume(volume)
}

// GetMasterVolume - Get master volume (listener)
func GetMasterVolume() float32 {
	return cGetMasterVolume()
}

// LoadWave - Load wave data from file
func LoadWave(fileName string) Wave {
	defer scratch.release(scratch.mark())
	cLoadWave(scratch.str(fileName))
	return waveFromC(readRet[cWave]())
}

// LoadWaveFromMemory - Load wave from memory buffer, fileType refers to extension: i.e. '.wav'
func LoadWaveFromMemory(fileType string, fileData []byte, dataSize int32) Wave {
	defer scratch.release(scratch.mark())
	cLoadWaveFromMemory(scratch.str(fileType), scratchSlice(fileData), dataSize)
	return waveFromC(readRet[cWave]())
}

// IsWaveValid - Checks if wave data is valid (data loaded and parameters)
func IsWaveValid(wave Wave) bool {
	defer scratch.release(scratch.mark())
	p0 := waveIn(&wave)
	ret := cIsWaveValid(p0)
	waveDone(p0)
	return ret != 0
}

// LoadSound - Load sound from file
func LoadSound(fileName string) Sound {
	defer scratch.release(scratch.mark())
	cLoadSound(scratch.str(fileName))
	return soundFromC(readRet[cSound]())
}

// LoadSoundFromWave - Load sound from wave data
func LoadSoundFromWave(wave Wave) Sound {
	defer scratch.release(scratch.mark())
	p0 := waveIn(&wave)
	cLoadSoundFromWave(p0)
	waveDone(p0)
	return soundFromC(readRet[cSound]())
}

// LoadSoundAlias - Create a new sound that shares the same sample data as the source sound, does not own the sound data
func LoadSoundAlias(source Sound) Sound {
	defer scratch.release(scratch.mark())
	cLoadSoundAlias(soundIn(&source))
	return soundFromC(readRet[cSound]())
}

// IsSoundValid - Checks if a sound is valid (data loaded and buffers initialized)
func IsSoundValid(sound Sound) bool {
	defer scratch.release(scratch.mark())
	return cIsSoundValid(soundIn(&sound)) != 0
}

// UpdateSound - Update sound buffer with new data (default data format: 32 bit float, stereo)
func UpdateSound(sound Sound, data []byte, sampleCount int32) {
	defer scratch.release(scratch.mark())
	cUpdateSound(soundIn(&sound), scratchSlice(data), sampleCount)
}

// UnloadWave - Unload wave data
func UnloadWave(wave Wave) {
}

// UnloadSound - Unload sound
func UnloadSound(sound Sound) {
	defer scratch.release(scratch.mark())
	cUnloadSound(soundIn(&sound))
	forgetAudioStream(&sound.Stream)
}

// UnloadSoundAlias - Unload a sound alias (does not deallocate sample data)
func UnloadSoundAlias(alias Sound) {
	defer scratch.release(scratch.mark())
	cUnloadSoundAlias(soundIn(&alias))
	forgetAudioStream(&alias.Stream)
}

// ExportWave - Export wave data to file, returns true on success
func ExportWave(wave Wave, fileName string) bool {
	defer scratch.release(scratch.mark())
	p0 := waveIn(&wave)
	ret := cExportWave(p0, scratch.str(fileName))
	waveDone(p0)
	return ret != 0
}

// PlaySound - Play a sound
func PlaySound(sound Sound) {
	defer scratch.release(scratch.mark())
	cPlaySound(soundIn(&sound))
}

// StopSound - Stop playing a sound
func StopSound(sound Sound) {
	defer scratch.release(scratch.mark())
	cStopSound(soundIn(&sound))
}

// PauseSound - Pause a sound
func PauseSound(sound Sound) {
	defer scratch.release(scratch.mark())
	cPauseSound(soundIn(&sound))
}

// ResumeSound - Resume a paused sound
func ResumeSound(sound Sound) {
	defer scratch.release(scratch.mark())
	cResumeSound(soundIn(&sound))
}

// IsSoundPlaying - Check if a sound is currently playing
func IsSoundPlaying(sound Sound) bool {
	defer scratch.release(scratch.mark())
	return cIsSoundPlaying(soundIn(&sound)) != 0
}

// SetSoundVolume - Set volume for a sound (1.0 is max level)
func SetSoundVolume(sound Sound, volume float32) {
	defer scratch.release(scratch.mark())
	cSetSoundVolume(soundIn(&sound), volume)
}

// SetSoundPitch - Set pitch for a sound (1.0 is base level)
func SetSoundPitch(sound Sound, pitch float32) {
	defer scratch.release(scratch.mark())
	cSetSoundPitch(soundIn(&sound), pitch)
}

// SetSoundPan - Set pan for a sound (-1.0 left, 0.0 center, 1.0 right)
func SetSoundPan(sound Sound, pan float32) {
	defer scratch.release(scratch.mark())
	cSetSoundPan(soundIn(&sound), pan)
}

// WaveCopy - Copy a wave to a new wave
func WaveCopy(wave Wave) Wave {
	defer scratch.release(scratch.mark())
	p0 := waveIn(&wave)
	cWaveCopy(p0)
	waveDone(p0)
	return waveFromC(readRet[cWave]())
}

// WaveCrop - Crop a wave to defined frames range
func WaveCrop(wave *Wave, initFrame int32, finalFrame int32) {
	defer scratch.release(scratch.mark())
	p0 := waveIn(wave)
	cWaveCrop(p0, initFrame, finalFrame)
	waveBack(wave, p0)
}

// WaveFormat - Convert wave data to desired format
func WaveFormat(wave *Wave, sampleRate int32, sampleSize int32, channels int32) {
	defer scratch.release(scratch.mark())
	p0 := waveIn(wave)
	cWaveFormat(p0, sampleRate, sampleRate, channels)
	waveBack(wave, p0)
}

// LoadWaveSamples - Get samples data from wave as a floats array
func LoadWaveSamples(wave Wave) []float32 {
	defer scratch.release(scratch.mark())
	p := waveIn(&wave)
	ret := cLoadWaveSamples(p)
	waveDone(p)
	samples := goSlice[float32](ret, int(wave.FrameCount*wave.Channels))
	cFree(ret)
	return samples
}

// UnloadWaveSamples - Unload samples data loaded with LoadWaveSamples()
func UnloadWaveSamples(samples []float32) {
}

// LoadMusicStream - Load music stream from file
func LoadMusicStream(fileName string) Music {
	defer scratch.release(scratch.mark())
	cLoadMusicStream(scratch.str(fileName))
	return musicFromC(readRet[cMusic]())
}

// LoadMusicStreamFromMemory - Load music stream from data
func LoadMusicStreamFromMemory(fileType string, data []byte, dataSize int32) Music {
	defer scratch.release(scratch.mark())
	p := cOwned(data)
	cLoadMusicStreamFromMemory(scratch.str(fileType), p, dataSize)
	music := musicFromC(readRet[cMusic]())
	if music.CtxData != nil {
		musicData[music.CtxData] = p
	} else {
		cFree(p)
	}
	return music
}

// IsMusicValid - Checks if a music stream is valid (context and buffers initialized)
func IsMusicValid(music Music) bool {
	defer scratch.release(scratch.mark())
	return cIsMusicValid(musicIn(&music)) != 0
}

// UnloadMusicStream - Unload music stream
func UnloadMusicStream(music Music) {
	defer scratch.release(scratch.mark())
	cUnloadMusicStream(musicIn(&music))
	forgetAudioStream(&music.Stream)
	dropHandle((*byte)(music.CtxData))
	if p, ok := musicData[music.CtxData]; ok {
		cFree(p)
		delete(musicData, music.CtxData)
	}
}

// PlayMusicStream - Start music playing
func PlayMusicStream(music Music) {
	defer scratch.release(scratch.mark())
	cPlayMusicStream(musicIn(&music))
}

// IsMusicStreamPlaying - Check if music is playing
func IsMusicStreamPlaying(music Music) bool {
	defer scratch.release(scratch.mark())
	return cIsMusicStreamPlaying(musicIn(&music)) != 0
}

// UpdateMusicStream - Updates buffers for music streaming
func UpdateMusicStream(music Music) {
	defer scratch.release(scratch.mark())
	cUpdateMusicStream(musicIn(&music))
}

// StopMusicStream - Stop music playing
func StopMusicStream(music Music) {
	defer scratch.release(scratch.mark())
	cStopMusicStream(musicIn(&music))
}

// PauseMusicStream - Pause music playing
func PauseMusicStream(music Music) {
	defer scratch.release(scratch.mark())
	cPauseMusicStream(musicIn(&music))
}

// ResumeMusicStream - Resume playing paused music
func ResumeMusicStream(music Music) {
	defer scratch.release(scratch.mark())
	cResumeMusicStream(musicIn(&music))
}

// SeekMusicStream - Seek music to a position (in seconds)
func SeekMusicStream(music Music, position float32) {
	defer scratch.release(scratch.mark())
	cSeekMusicStream(musicIn(&music), position)
}

// SetMusicVolume - Set volume for music (1.0 is max level)
func SetMusicVolume(music Music, volume float32) {
	defer scratch.release(scratch.mark())
	cSetMusicVolume(musicIn(&music), volume)
}

// SetMusicPitch - Set pitch for a music (1.0 is base level)
func SetMusicPitch(music Music, pitch float32) {
	defer scratch.release(scratch.mark())
	cSetMusicPitch(musicIn(&music), pitch)
}

// SetMusicPan - Set pan for a music (-1.0 left, 0.0 center, 1.0 right)
func SetMusicPan(music Music, pan float32) {
	defer scratch.release(scratch.mark())
	cSetMusicPan(musicIn(&music), pan)
}

// GetMusicTimeLength - Get music time length (in seconds)
func GetMusicTimeLength(music Music) float32 {
	defer scratch.release(scratch.mark())
	return cGetMusicTimeLength(musicIn(&music))
}

// GetMusicTimePlayed - Get current music time played (in seconds)
func GetMusicTimePlayed(music Music) float32 {
	defer scratch.release(scratch.mark())
	return cGetMusicTimePlayed(musicIn(&music))
}

// LoadAudioStream - Load audio stream (to stream raw audio pcm data)
func LoadAudioStream(sampleRate uint32, sampleSize uint32, channels uint32) AudioStream {
	cLoadAudioStream(sampleRate, sampleSize, channels)
	return audioStreamFromC(readRet[cAudioStream]())
}

// IsAudioStreamValid - Checks if an audio stream is valid (buffers initialized)
func IsAudioStreamValid(stream AudioStream) bool {
	defer scratch.release(scratch.mark())
	return cIsAudioStreamValid(audioStreamIn(&stream)) != 0
}

// UnloadAudioStream - Unload audio stream and free memory
func UnloadAudioStream(stream AudioStream) {
	defer scratch.release(scratch.mark())
	cUnloadAudioStream(audioStreamIn(&stream))
	forgetAudioStream(&stream)
}

// UpdateAudioStream - Update audio stream buffers with data
func UpdateAudioStream[T float32 | int16](stream AudioStream, data []T) {
	defer scratch.release(scratch.mark())
	cUpdateAudioStream(audioStreamIn(&stream), scratchSlice(data), int32(len(data)))
}

// IsAudioStreamProcessed - Check if any audio stream buffers requires refill
func IsAudioStreamProcessed(stream AudioStream) bool {
	defer scratch.release(scratch.mark())
	return cIsAudioStreamProcessed(audioStreamIn(&stream)) != 0
}

// PlayAudioStream - Play audio stream
func PlayAudioStream(stream AudioStream) {
	defer scratch.release(scratch.mark())
	cPlayAudioStream(audioStreamIn(&stream))
}

// PauseAudioStream - Pause audio stream
func PauseAudioStream(stream AudioStream) {
	defer scratch.release(scratch.mark())
	cPauseAudioStream(audioStreamIn(&stream))
}

// ResumeAudioStream - Resume audio stream
func ResumeAudioStream(stream AudioStream) {
	defer scratch.release(scratch.mark())
	cResumeAudioStream(audioStreamIn(&stream))
}

// IsAudioStreamPlaying - Check if audio stream is playing
func IsAudioStreamPlaying(stream AudioStream) bool {
	defer scratch.release(scratch.mark())
	return cIsAudioStreamPlaying(audioStreamIn(&stream)) != 0
}

// StopAudioStream - Stop audio stream
func StopAudioStream(stream AudioStream) {
	defer scratch.release(scratch.mark())
	cStopAudioStream(audioStreamIn(&stream))
}

// SetAudioStreamVolume - Set volume for audio stream (1.0 is max level)
func SetAudioStreamVolume(stream AudioStream, volume float32) {
	defer scratch.release(scratch.mark())
	cSetAudioStreamVolume(audioStreamIn(&stream), volume)
}

// SetAudioStreamPitch - Set pitch for audio stream (1.0 is base level)
func SetAudioStreamPitch(stream AudioStream, pitch float32) {
	defer scratch.release(scratch.mark())
	cSetAudioStreamPitch(audioStreamIn(&stream), pitch)
}

// SetAudioStreamPan - Set pan for audio stream (-1.0 to 1.0 range, 0.0 is centered)
func SetAudioStreamPan(stream AudioStream, pan float32) {
	defer scratch.release(scratch.mark())
	cSetAudioStreamPan(audioStreamIn(&stream), pan)
}

// SetAudioStreamBufferSizeDefault - Default size for new audio streams
func SetAudioStreamBufferSizeDefault(size int32) {
	cSetAudioStreamBufferSizeDefault(size)
}

// SetAudioStreamCallback - Audio thread callback to request new data
func SetAudioStreamCallback(stream AudioStream, callback AudioCallback) {
	slot := streamCallbackSlot(&stream, callback)
	if slot == -2 {
		return
	}
	defer scratch.release(scratch.mark())
	cSetAudioStreamCallbackSlot(audioStreamIn(&stream), slot)
}

// AttachAudioStreamProcessor - Attach audio stream processor to stream, receives frames x 2 samples as 'float' (stereo)
func AttachAudioStreamProcessor(stream AudioStream, processor AudioCallback) {
	slot := acquireAudioStreamProcessor(stream, processor)
	if slot < 0 {
		return
	}
	defer scratch.release(scratch.mark())
	cAttachAudioStreamProcessorSlot(audioStreamIn(&stream), int32(slot))
}

// DetachAudioStreamProcessor - Detach audio stream processor from stream
func DetachAudioStreamProcessor(stream AudioStream, processor AudioCallback) {
	slot := findAudioStreamProcessor(stream, processor)
	if slot < 0 {
		return
	}
	defer scratch.release(scratch.mark())
	cDetachAudioStreamProcessorSlot(audioStreamIn(&stream), int32(slot))
	releaseAudioStreamProcessor(slot)
}

// AttachAudioMixedProcessor - Attach audio stream processor to the entire audio pipeline, receives frames x 2 samples as 'float' (stereo)
func AttachAudioMixedProcessor(processor AudioCallback) {
	if addAudioMixedProcessor(processor) {
		cAttachAudioMixedProcessorSlot()
	}
}

// DetachAudioMixedProcessor - Detach audio stream processor from the entire audio pipeline
func DetachAudioMixedProcessor(processor AudioCallback) {
	if removeAudioMixedProcessor(processor) {
		cDetachAudioMixedProcessorSlot()
	}
}

// SetMain - Sets callback function
func SetMain(func()) {
}

// NewImageFromImage - Returns new Image from Go image.Image
func NewImageFromImage(img image.Image) *Image {
	size := img.Bounds().Size()

	ret := GenImageColor(size.X, size.Y, White)

	for y := 0; y < size.Y; y++ {
		for x := 0; x < size.X; x++ {
			col := img.At(x, y)
			r, g, b, a := col.RGBA()
			rcolor := NewColor(uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a>>8))
			ImageDrawPixel(ret, int32(x), int32(y), rcolor)
		}
	}

	return ret
}

// ToImage converts a Image to Go image.Image
func (i *Image) ToImage() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, int(i.Width), int(i.Height)))

	// Get pixel data from image (RGBA 32bit)
	ret := LoadImageColors(i)
	pixels := (*[1 << 24]uint8)(unsafe.Pointer(unsafe.SliceData(ret)))[0 : i.Width*i.Height*4]

	img.Pix = pixels

	return img
}

// HomeDir - Returns user home directory
// NOTE: On Android this returns internal data path and must be called after InitWindow
func HomeDir() string {
	if homeDir, err := os.UserHomeDir(); err == nil {
		return homeDir
	}
	return ""
}

func texturePixels[T []color.RGBA | *image.RGBA | []byte](pixels T) uint32 {
	switch p := any(pixels).(type) {
	case []color.RGBA:
		return scratchSlice(p)
	case *image.RGBA:
		return scratchSlice(p.Pix)
	case []byte:
		return scratchSlice(p)
	}
	return 0
}
