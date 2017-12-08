// +build js

//go:generate ./external/scripts/emcc-generate-js.sh

package raylib

import (
	"unsafe"

	"github.com/gopherjs/gopherjs/js"
)

// InitWindow - Initialize Window and OpenGL Graphics
func InitWindow(width int32, height int32, t interface{}) {
	js.Global.Get("Module").Call("_InitWindow", width, height, t.(string))
}

// SetCallbackFunc - Sets callback function
func SetCallbackFunc(func(unsafe.Pointer)) {
}

// SetMainLoop - Sets main loop function
func SetMainLoop(f func(), fps, simulateInfiniteLoop int) {
	js.Global.Get("Module").Call("_emscripten_set_main_loop", f, fps, simulateInfiniteLoop)
}

// ShowCursor - Shows cursor
func ShowCursor() {
}

// HideCursor - Hides cursor
func HideCursor() {
}

// IsCursorHidden - Returns true if cursor is not visible
func IsCursorHidden() bool {
	return false
}

// EnableCursor - Enables cursor
func EnableCursor() {
}

// DisableCursor - Disables cursor
func DisableCursor() {
}

// IsFileDropped - Check if a file have been dropped into window
func IsFileDropped() bool {
	return false
}

// GetDroppedFiles - Retrieve dropped files into window
func GetDroppedFiles(count *int32) (f []string) {
	return
}

// ClearDroppedFiles - Clear dropped files paths buffer
func ClearDroppedFiles() {
}

// InitAudioDevice - Initialize audio device and context
func InitAudioDevice() {
	js.Global.Get("Module").Call("_InitAudioDevice")
}

// CloseAudioDevice - Close the audio device and context
func CloseAudioDevice() {
	js.Global.Get("Module").Call("_CloseAudioDevice")
}

// IsAudioDeviceReady - Check if audio device has been initialized successfully
func IsAudioDeviceReady() bool {
	return js.Global.Get("Module").Call("_IsAudioDeviceReady").Bool()
}

// SetMasterVolume - Set master volume (listener)
func SetMasterVolume(volume float32) {
	js.Global.Get("Module").Call("_SetMasterVolume", volume)
}

// LoadWave - Load wave data from file into RAM
func LoadWave(fileName string) Wave {
	return newWaveFromPointer(unsafe.Pointer(js.Global.Get("Module").Call("_LoadWave", fileName).Unsafe()))
}

// LoadWaveEx - Load wave data from float array data (32bit)
func LoadWaveEx(data []byte, sampleCount int32, sampleRate int32, sampleSize int32, channels int32) Wave {
	return newWaveFromPointer(unsafe.Pointer(js.Global.Get("Module").Call("_LoadWaveEx", data, sampleCount, sampleRate, sampleSize, channels).Unsafe()))
}

// LoadSound - Load sound to memory
func LoadSound(fileName string) Sound {
	return newSoundFromPointer(unsafe.Pointer(js.Global.Get("Module").Call("_LoadSound", fileName).Unsafe()))
}

// LoadSoundFromWave - Load sound to memory from wave data
func LoadSoundFromWave(wave Wave) Sound {
	return newSoundFromPointer(unsafe.Pointer(js.Global.Get("Module").Call("_LoadSoundFromWave", wave).Unsafe()))
}

// UpdateSound - Update sound buffer with new data
func UpdateSound(sound Sound, data []byte, samplesCount int32) {
	js.Global.Get("Module").Call("_UpdateSound", sound, data, samplesCount)
}

// UnloadWave - Unload wave data
func UnloadWave(wave Wave) {
	js.Global.Get("Module").Call("_UnloadWave", wave)
}

// UnloadSound - Unload sound
func UnloadSound(sound Sound) {
	js.Global.Get("Module").Call("_UnloadSound", sound)
}

// PlaySound - Play a sound
func PlaySound(sound Sound) {
	js.Global.Get("Module").Call("_PlaySound", sound)
}

// PauseSound - Pause a sound
func PauseSound(sound Sound) {
	js.Global.Get("Module").Call("_PauseSound", sound)
}

// ResumeSound - Resume a paused sound
func ResumeSound(sound Sound) {
	js.Global.Get("Module").Call("_ResumeSound", sound)
}

// StopSound - Stop playing a sound
func StopSound(sound Sound) {
	js.Global.Get("Module").Call("_StopSound", sound)
}

// IsSoundPlaying - Check if a sound is currently playing
func IsSoundPlaying(sound Sound) bool {
	return js.Global.Get("Module").Call("_IsSoundPlaying", sound).Bool()
}

// SetSoundVolume - Set volume for a sound (1.0 is max level)
func SetSoundVolume(sound Sound, volume float32) {
	js.Global.Get("Module").Call("_SetSoundVolume", sound, volume)
}

// SetSoundPitch - Set pitch for a sound (1.0 is base level)
func SetSoundPitch(sound Sound, pitch float32) {
	js.Global.Get("Module").Call("_SetSoundPitch", sound, pitch)
}

// WaveFormat - Convert wave data to desired format
func WaveFormat(wave Wave, sampleRate int32, sampleSize int32, channels int32) {
	js.Global.Get("Module").Call("_WaveFormat", wave, sampleRate, sampleSize, channels)
}

// WaveCopy - Copy a wave to a new wave
func WaveCopy(wave Wave) Wave {
	return newWaveFromPointer(unsafe.Pointer(js.Global.Get("Module").Call("_WaveCopy", wave).Unsafe()))
}

// WaveCrop - Crop a wave to defined samples range
func WaveCrop(wave Wave, initSample int32, finalSample int32) {
	js.Global.Get("Module").Call("_WaveCrop", wave, initSample, finalSample)
}

// GetWaveData - Get samples data from wave as a floats array
func GetWaveData(wave Wave) []float32 {
	return js.Global.Get("Module").Call("_GetWaveData", wave).Interface().([]float32)
}

// LoadMusicStream - Load music stream from file
func LoadMusicStream(fileName string) Music {
	return newMusicFromPointer(unsafe.Pointer(js.Global.Get("Module").Call("_LoadMusicStream", fileName).Unsafe()))
}

// UnloadMusicStream - Unload music stream
func UnloadMusicStream(music Music) {
	js.Global.Get("Module").Call("_UnloadMusicStream", music)
}

// PlayMusicStream - Start music playing
func PlayMusicStream(music Music) {
	js.Global.Get("Module").Call("_PlayMusicStream", music)
}

// UpdateMusicStream - Updates buffers for music streaming
func UpdateMusicStream(music Music) {
	js.Global.Get("Module").Call("_UpdateMusicStream", music)
}

// StopMusicStream - Stop music playing
func StopMusicStream(music Music) {
	js.Global.Get("Module").Call("_StopMusicStream", music)
}

// PauseMusicStream - Pause music playing
func PauseMusicStream(music Music) {
	js.Global.Get("Module").Call("_PauseMusicStream", music)
}

// ResumeMusicStream - Resume playing paused music
func ResumeMusicStream(music Music) {
	js.Global.Get("Module").Call("_ResumeMusicStream", music)
}

// IsMusicPlaying - Check if music is playing
func IsMusicPlaying(music Music) bool {
	return js.Global.Get("Module").Call("_IsMusicPlaying", music).Bool()
}

// SetMusicVolume - Set volume for music (1.0 is max level)
func SetMusicVolume(music Music, volume float32) {
	js.Global.Get("Module").Call("_SetMusicVolume", music, volume)
}

// SetMusicPitch - Set pitch for a music (1.0 is base level)
func SetMusicPitch(music Music, pitch float32) {
	js.Global.Get("Module").Call("_SetMusicPitch", music, pitch)
}

// NOTE: If set to -1, means infinite loop
func SetMusicLoopCount(music Music, count int32) {
	js.Global.Get("Module").Call("_SetMusicLoopCount", music, count)
}

// GetMusicTimeLength - Get music time length (in seconds)
func GetMusicTimeLength(music Music) float32 {
	return float32(js.Global.Get("Module").Call("_GetMusicTimeLength", music).Float())
}

// GetMusicTimePlayed - Get current music time played (in seconds)
func GetMusicTimePlayed(music Music) float32 {
	return float32(js.Global.Get("Module").Call("_GetMusicTimePlayed", music).Float())
}

// InitAudioStream - Init audio stream (to stream raw audio pcm data)
func InitAudioStream(sampleRate uint32, sampleSize uint32, channels uint32) AudioStream {
	return newAudioStreamFromPointer(unsafe.Pointer(js.Global.Get("Module").Call("_InitAudioStream", sampleRate, sampleSize, channels).Unsafe()))
}

// UpdateAudioStream - Update audio stream buffers with data
func UpdateAudioStream(stream AudioStream, data []float32, samplesCount int32) {
	js.Global.Get("Module").Call("_UpdateAudioStream", stream, data, samplesCount)
}

// CloseAudioStream - Close audio stream and free memory
func CloseAudioStream(stream AudioStream) {
	js.Global.Get("Module").Call("_CloseAudioStream", stream)
}

// IsAudioBufferProcessed - Check if any audio stream buffers requires refill
func IsAudioBufferProcessed(stream AudioStream) bool {
	return js.Global.Get("Module").Call("_IsAudioBufferProcessed", stream).Bool()
}

// PlayAudioStream - Play audio stream
func PlayAudioStream(stream AudioStream) {
	js.Global.Get("Module").Call("_PlayAudioStream", stream)
}

// PauseAudioStream - Pause audio stream
func PauseAudioStream(stream AudioStream) {
	js.Global.Get("Module").Call("_PauseAudioStream", stream)
}

// ResumeAudioStream - Resume audio stream
func ResumeAudioStream(stream AudioStream) {
	js.Global.Get("Module").Call("_ResumeAudioStream", stream)
}

// StopAudioStream - Stop audio stream
func StopAudioStream(stream AudioStream) {
	js.Global.Get("Module").Call("_StopAudioStream", stream)
}

// SetCameraMode - Set camera mode (multiple camera modes available)
func SetCameraMode(camera Camera, mode CameraMode) {
	js.Global.Get("Module").Call("_SetCameraMode", camera, mode)
}

// UpdateCamera - Update camera position for selected mode
func UpdateCamera(camera *Camera) {
	js.Global.Get("Module").Call("_UpdateCamera", camera)
}

// SetCameraPanControl - Set camera pan key to combine with mouse movement (free camera)
func SetCameraPanControl(panKey int32) {
	js.Global.Get("Module").Call("_SetCameraPanControl", panKey)
}

// SetCameraAltControl - Set camera alt key to combine with mouse movement (free camera)
func SetCameraAltControl(altKey int32) {
	js.Global.Get("Module").Call("_SetCameraAltControl", altKey)
}

// SetCameraSmoothZoomControl - Set camera smooth zoom key to combine with mouse (free camera)
func SetCameraSmoothZoomControl(szKey int32) {
	js.Global.Get("Module").Call("_SetCameraSmoothZoomControl", szKey)
}

// SetCameraMoveControls - Set camera move controls (1st person and 3rd person cameras)
func SetCameraMoveControls(frontKey int32, backKey int32, rightKey int32, leftKey int32, upKey int32, downKey int32) {
	js.Global.Get("Module").Call("_SetCameraMoveControls", frontKey, backKey, rightKey, leftKey, upKey, downKey)
}

// CloseWindow - Close Window and Terminate Context
func CloseWindow() {
	js.Global.Get("Module").Call("_CloseWindow")
}

// WindowShouldClose - Detect if KEY_ESCAPE pressed or Close icon pressed
func WindowShouldClose() bool {
	return js.Global.Get("Module").Call("_WindowShouldClose").Bool()
}

// IsWindowMinimized - Detect if window has been minimized (or lost focus)
func IsWindowMinimized() bool {
	return js.Global.Get("Module").Call("_IsWindowMinimized").Bool()
}

// ToggleFullscreen - Fullscreen toggle (only PLATFORM_DESKTOP)
func ToggleFullscreen() {
	js.Global.Get("Module").Call("_ToggleFullscreen")
}

// SetWindowIcon - Set icon for window (only PLATFORM_DESKTOP)
func SetWindowIcon(image Image) {
	js.Global.Get("Module").Call("_SetWindowIcon", image)
}

// SetWindowTitle - Set title for window (only PLATFORM_DESKTOP)
func SetWindowTitle(title string) {
	js.Global.Get("Module").Call("_SetWindowTitle", title)
}

// SetWindowPosition - Set window position on screen (only PLATFORM_DESKTOP)
func SetWindowPosition(x, y int32) {
	js.Global.Get("Module").Call("_SetWindowPosition", x, y)
}

// SetWindowMonitor - Set monitor for the current window (fullscreen mode)
func SetWindowMonitor(monitor int32) {
	js.Global.Get("Module").Call("_SetWindowMonitor", monitor)
}

// GetScreenWidth - Get current screen width
func GetScreenWidth() int32 {
	return int32(js.Global.Get("Module").Call("_GetScreenWidth").Int())
}

// GetScreenHeight - Get current screen height
func GetScreenHeight() int32 {
	return int32(js.Global.Get("Module").Call("_GetScreenHeight").Int())
}

// ClearBackground - Sets Background Color
func ClearBackground(color Color) {
	js.Global.Get("Module").Call("_ClearBackground", color)
}

// BeginDrawing - Setup drawing canvas to start drawing
func BeginDrawing() {
	js.Global.Get("Module").Call("_BeginDrawing")
}

// EndDrawing - End canvas drawing and Swap Buffers (Double Buffering)
func EndDrawing() {
	js.Global.Get("Module").Call("_EndDrawing")
}

// Begin2dMode - Initialize 2D mode with custom camera
func Begin2dMode(camera Camera2D) {
	js.Global.Get("Module").Call("_Begin2dMode", camera)
}

// End2dMode - Ends 2D mode custom camera usage
func End2dMode() {
	js.Global.Get("Module").Call("_End2dMode")
}

// Begin3dMode - Initializes 3D mode for drawing (Camera setup)
func Begin3dMode(camera Camera) {
	js.Global.Get("Module").Call("_Begin3dMode", camera)
}

// End3dMode - Ends 3D mode and returns to default 2D orthographic mode
func End3dMode() {
	js.Global.Get("Module").Call("_End3dMode")
}

// BeginTextureMode - Initializes render texture for drawing
func BeginTextureMode(target RenderTexture2D) {
	js.Global.Get("Module").Call("_BeginTextureMode", target)
}

// EndTextureMode - Ends drawing to render texture
func EndTextureMode() {
	js.Global.Get("Module").Call("_EndTextureMode")
}

// GetMouseRay - Returns a ray trace from mouse position
func GetMouseRay(mousePosition Vector2, camera Camera) Ray {
	return newRayFromPointer(unsafe.Pointer(js.Global.Get("Module").Call("_GetMouseRay", mousePosition, camera).Unsafe()))
}

// GetWorldToScreen - Returns the screen space position from a 3d world space position
func GetWorldToScreen(position Vector3, camera Camera) Vector2 {
	return newVector2FromPointer(unsafe.Pointer(js.Global.Get("Module").Call("_GetWorldToScreen", position, camera).Unsafe()))
}

// GetCameraMatrix - Returns camera transform matrix (view matrix)
func GetCameraMatrix(camera Camera) Matrix {
	return newMatrixFromPointer(unsafe.Pointer(js.Global.Get("Module").Call("_GetCameraMatrix", camera).Unsafe()))
}

// SetTargetFPS - Set target FPS (maximum)
func SetTargetFPS(fps int32) {
	js.Global.Get("Module").Call("_SetTargetFPS")
}

// GetFPS - Returns current FPS
func GetFPS() float32 {
	return float32(js.Global.Get("Module").Call("_GetFPS").Float())
}

// GetFrameTime - Returns time in seconds for one frame
func GetFrameTime() float32 {
	return float32(js.Global.Get("Module").Call("_GetFrameTime").Float())
}

// GetColor - Returns a Color struct from hexadecimal value
func GetColor(hexValue int32) Color {
	return newColorFromPointer(unsafe.Pointer(js.Global.Get("Module").Call("_GetColor", hexValue).Unsafe()))
}

// GetHexValue - Returns hexadecimal value for a Color
func GetHexValue(color Color) int32 {
	return int32(js.Global.Get("Module").Call("_GetHexValue", color).Int())
}

// ColorToFloat - Converts Color to float32 slice and normalizes
func ColorToFloat(color Color) []float32 {
	return js.Global.Get("Module").Call("_ColorToFloat", color).Interface().([]float32)
}

// Vector3ToFloat - Converts Vector3 to float32 slice
func Vector3ToFloat(vec Vector3) []float32 {
	return js.Global.Get("Module").Call("_Vector3ToFloat", vec).Interface().([]float32)
}

// MatrixToFloat - Converts Matrix to float32 slice
func MatrixToFloat(mat Matrix) []float32 {
	return js.Global.Get("Module").Call("_MatrixToFloat", mat).Interface().([]float32)
}

// GetRandomValue - Returns a random value between min and max (both included)
func GetRandomValue(min, max int32) int32 {
	return int32(js.Global.Get("Module").Call("_GetRandomValue", min, max).Int())
}

// Fade - Color fade-in or fade-out, alpha goes from 0.0f to 1.0f
func Fade(color Color, alpha float32) Color {
	return newColorFromPointer(unsafe.Pointer(js.Global.Get("Module").Call("_Fade", color, alpha).Unsafe()))
}

// ShowLogo - Activates raylib logo at startup (can be done with flags)
func ShowLogo() {
	js.Global.Get("Module").Call("_ShowLogo")
}

// SetConfigFlags - Setup some window configuration flags
func SetConfigFlags(flags byte) {
	js.Global.Get("Module").Call("_SetConfigFlags", flags)
}

// TakeScreenshot - Takes a screenshot of current screen (saved a .png)
func TakeScreenshot(name string) {
	js.Global.Get("Module").Call("_TakeScreenshot", name)
}

// StorageSaveValue - Storage save integer value (to defined position)
func StorageSaveValue(position, value int32) {
	js.Global.Get("Module").Call("_StorageSaveValue", position, value)
}

// StorageLoadValue - Storage load integer value (from defined position)
func StorageLoadValue(position int32) int32 {
	return int32(js.Global.Get("Module").Call("_StorageLoadValue", position).Int())
}

// IsKeyPressed - Detect if a key has been pressed once
func IsKeyPressed(key int32) bool {
	return js.Global.Get("Module").Call("_IsKeyPressed", key).Bool()
}

// IsKeyDown - Detect if a key is being pressed
func IsKeyDown(key int32) bool {
	return js.Global.Get("Module").Call("_IsKeyDown", key).Bool()
}

// IsKeyReleased - Detect if a key has been released once
func IsKeyReleased(key int32) bool {
	return js.Global.Get("Module").Call("_IsKeyReleased", key).Bool()
}

// IsKeyUp - Detect if a key is NOT being pressed
func IsKeyUp(key int32) bool {
	return js.Global.Get("Module").Call("_IsKeyUp", key).Bool()
}

// GetKeyPressed - Get latest key pressed
func GetKeyPressed() int32 {
	return int32(js.Global.Get("Module").Call("_GetKeyPressed").Int())
}

// SetExitKey - Set a custom key to exit program (default is ESC)
func SetExitKey(key int32) {
	js.Global.Get("Module").Call("_SetExitKey", key)
}

// IsGamepadAvailable - Detect if a gamepad is available
func IsGamepadAvailable(gamepad int32) bool {
	return js.Global.Get("Module").Call("_IsGamepadAvailable", gamepad).Bool()
}

// IsGamepadName - Check gamepad name (if available)
func IsGamepadName(gamepad int32, name string) bool {
	return js.Global.Get("Module").Call("_IsGamepadName", gamepad, name).Bool()
}

// GetGamepadName - Return gamepad internal name id
func GetGamepadName(gamepad int32) string {
	return js.Global.Get("Module").Call("_GetGamepadName", gamepad).String()
}

// IsGamepadButtonPressed - Detect if a gamepad button has been pressed once
func IsGamepadButtonPressed(gamepad, button int32) bool {
	return js.Global.Get("Module").Call("_IsGamepadButtonPressed", gamepad, button).Bool()
}

// IsGamepadButtonDown - Detect if a gamepad button is being pressed
func IsGamepadButtonDown(gamepad, button int32) bool {
	return js.Global.Get("Module").Call("_IsGamepadButtonDown", gamepad, button).Bool()
}

// IsGamepadButtonReleased - Detect if a gamepad button has been released once
func IsGamepadButtonReleased(gamepad, button int32) bool {
	return js.Global.Get("Module").Call("_IsGamepadButtonReleased", gamepad, button).Bool()
}

// IsGamepadButtonUp - Detect if a gamepad button is NOT being pressed
func IsGamepadButtonUp(gamepad, button int32) bool {
	return js.Global.Get("Module").Call("_IsGamepadButtonUp", gamepad, button).Bool()
}

// GetGamepadButtonPressed - Get the last gamepad button pressed
func GetGamepadButtonPressed() int32 {
	return int32(js.Global.Get("Module").Call("_GetGamepadButtonPressed").Int())
}

// GetGamepadAxisCount - Return gamepad axis count for a gamepad
func GetGamepadAxisCount(gamepad int32) int32 {
	return int32(js.Global.Get("Module").Call("_GetGamepadAxisCount", gamepad).Int())
}

// GetGamepadAxisMovement - Return axis movement value for a gamepad axis
func GetGamepadAxisMovement(gamepad, axis int32) float32 {
	return float32(js.Global.Get("Module").Call("_GetGamepadAxisMovement", gamepad, axis).Float())
}

// IsMouseButtonPressed - Detect if a mouse button has been pressed once
func IsMouseButtonPressed(button int32) bool {
	return js.Global.Get("Module").Call("_IsMouseButtonPressed", button).Bool()
}

// IsMouseButtonDown - Detect if a mouse button is being pressed
func IsMouseButtonDown(button int32) bool {
	return js.Global.Get("Module").Call("_IsMouseButtonDown", button).Bool()
}

// IsMouseButtonReleased - Detect if a mouse button has been released once
func IsMouseButtonReleased(button int32) bool {
	return js.Global.Get("Module").Call("_IsMouseButtonReleased", button).Bool()
}

// IsMouseButtonUp - Detect if a mouse button is NOT being pressed
func IsMouseButtonUp(button int32) bool {
	return js.Global.Get("Module").Call("_IsMouseButtonUp", button).Bool()
}

// GetMouseX - Returns mouse position X
func GetMouseX() int32 {
	return int32(js.Global.Get("Module").Call("_GetMouseX").Int())
}

// GetMouseY - Returns mouse position Y
func GetMouseY() int32 {
	return int32(js.Global.Get("Module").Call("_GetMouseY").Int())
}

// GetMousePosition - Returns mouse position XY
func GetMousePosition() Vector2 {
	return newVector2FromPointer(unsafe.Pointer(js.Global.Get("Module").Call("_GetMousePosition").Unsafe()))
}

// SetMousePosition - Set mouse position XY
func SetMousePosition(position Vector2) {
	js.Global.Get("Module").Call("_SetMousePosition", position)
}

// GetMouseWheelMove - Returns mouse wheel movement Y
func GetMouseWheelMove() int32 {
	return int32(js.Global.Get("Module").Call("_GetMouseWheelMove").Int())
}

// GetTouchX - Returns touch position X for touch point 0 (relative to screen size)
func GetTouchX() int32 {
	return int32(js.Global.Get("Module").Call("_GetTouchX").Int())
}

// GetTouchY - Returns touch position Y for touch point 0 (relative to screen size)
func GetTouchY() int32 {
	return int32(js.Global.Get("Module").Call("_GetTouchX").Int())
}

// GetTouchPosition - Returns touch position XY for a touch point index (relative to screen size)
func GetTouchPosition(index int32) Vector2 {
	return newVector2FromPointer(unsafe.Pointer(js.Global.Get("Module").Call("_GetTouchPosition", index).Unsafe()))
}

// SetGesturesEnabled - Enable a set of gestures using flags
func SetGesturesEnabled(gestureFlags uint32) {
	js.Global.Get("Module").Call("_SetGesturesEnabled", gestureFlags)
}

// IsGestureDetected - Check if a gesture have been detected
func IsGestureDetected(gesture Gestures) bool {
	return js.Global.Get("Module").Call("_IsGestureDetected", gesture).Bool()
}

// GetGestureDetected - Get latest detected gesture
func GetGestureDetected() Gestures {
	return Gestures(js.Global.Get("Module").Call("_GetGestureDetected").Int())
}

// GetTouchPointsCount - Get touch points count
func GetTouchPointsCount() int32 {
	return int32(js.Global.Get("Module").Call("_GetTouchPointsCount").Int())
}

// GetGestureHoldDuration - Get gesture hold time in milliseconds
func GetGestureHoldDuration() float32 {
	return float32(js.Global.Get("Module").Call("_GetGestureHoldDuration").Float())
}

// GetGestureDragVector - Get gesture drag vector
func GetGestureDragVector() Vector2 {
	return newVector2FromPointer(unsafe.Pointer(js.Global.Get("Module").Call("_GetGestureDragVector").Unsafe()))
}

// GetGestureDragAngle - Get gesture drag angle
func GetGestureDragAngle() float32 {
	return float32(js.Global.Get("Module").Call("_GetGestureDragAngle").Float())
}

// GetGesturePinchVector - Get gesture pinch delta
func GetGesturePinchVector() Vector2 {
	return newVector2FromPointer(unsafe.Pointer(js.Global.Get("Module").Call("_GetGesturePinchVector").Unsafe()))
}

// GetGesturePinchAngle - Get gesture pinch angle
func GetGesturePinchAngle() float32 {
	return float32(js.Global.Get("Module").Call("_GetGesturePinchAngle").Float())
}

// DrawLine3D - Draw a line in 3D world space
func DrawLine3D(startPos Vector3, endPos Vector3, color Color) {
	js.Global.Get("Module").Call("_DrawLine3D", startPos, endPos, color)
}

// DrawCircle3D - Draw a circle in 3D world space
func DrawCircle3D(center Vector3, radius float32, rotationAxis Vector3, rotationAngle float32, color Color) {
	js.Global.Get("Module").Call("_DrawCircle3D", center, radius, rotationAxis, rotationAngle, color)
}

// DrawCube - Draw cube
func DrawCube(position Vector3, width float32, height float32, length float32, color Color) {
	js.Global.Get("Module").Call("_DrawCube", position, width, height, length, color)
}

// DrawCubeV - Draw cube (Vector version)
func DrawCubeV(position Vector3, size Vector3, color Color) {
	js.Global.Get("Module").Call("_DrawCubeV", position, size, color)
}

// DrawCubeWires - Draw cube wires
func DrawCubeWires(position Vector3, width float32, height float32, length float32, color Color) {
	js.Global.Get("Module").Call("_DrawCubeWires", position, width, height, length, color)
}

// DrawCubeTexture - Draw cube textured
func DrawCubeTexture(texture Texture2D, position Vector3, width float32, height float32, length float32, color Color) {
	js.Global.Get("Module").Call("_DrawCubeTexture", texture, position, width, height, length, color)
}

// DrawSphere - Draw sphere
func DrawSphere(centerPos Vector3, radius float32, color Color) {
	js.Global.Get("Module").Call("_DrawSphere", centerPos, radius, color)
}

// DrawSphereEx - Draw sphere with extended parameters
func DrawSphereEx(centerPos Vector3, radius float32, rings int32, slices int32, color Color) {
	js.Global.Get("Module").Call("_DrawSphereEx", centerPos, radius, rings, slices, color)
}

// DrawSphereWires - Draw sphere wires
func DrawSphereWires(centerPos Vector3, radius float32, rings int32, slices int32, color Color) {
	js.Global.Get("Module").Call("_DrawSphereWires", centerPos, radius, rings, slices, color)
}

// DrawCylinder - Draw a cylinder/cone
func DrawCylinder(position Vector3, radiusTop float32, radiusBottom float32, height float32, slices int32, color Color) {
	js.Global.Get("Module").Call("_DrawCylinder", position, radiusTop, radiusBottom, height, slices, color)
}

// DrawCylinderWires - Draw a cylinder/cone wires
func DrawCylinderWires(position Vector3, radiusTop float32, radiusBottom float32, height float32, slices int32, color Color) {
	js.Global.Get("Module").Call("_DrawCylinderWires", position, radiusTop, radiusBottom, height, slices, color)
}

// DrawPlane - Draw a plane XZ
func DrawPlane(centerPos Vector3, size Vector2, color Color) {
	js.Global.Get("Module").Call("_DrawPlane", centerPos, size, color)
}

// DrawRay - Draw a ray line
func DrawRay(ray Ray, color Color) {
	js.Global.Get("Module").Call("_DrawRay", ray, color)
}

// DrawGrid - Draw a grid (centered at (0, 0, 0))
func DrawGrid(slices int32, spacing float32) {
	js.Global.Get("Module").Call("_DrawGrid", slices, spacing)
}

// DrawGizmo - Draw simple gizmo
func DrawGizmo(position Vector3) {
	js.Global.Get("Module").Call("_DrawGizmo", position)
}

// LoadMesh - Load mesh from file
func LoadMesh(fileName string) Mesh {
	return newMeshFromPointer(unsafe.Pointer(js.Global.Get("Module").Call("_LoadMesh", fileName).Unsafe()))
}

// LoadModel - Load model from file
func LoadModel(fileName string) Model {
	return newModelFromPointer(unsafe.Pointer(js.Global.Get("Module").Call("_LoadModel", fileName).Unsafe()))
}

// LoadModelFromMesh - Load model from mesh data
func LoadModelFromMesh(data Mesh) Model {
	return newModelFromPointer(unsafe.Pointer(js.Global.Get("Module").Call("_LoadModelFromMesh", data).Unsafe()))
}

// UnloadModel - Unload model from memory (RAM and/or VRAM)
func UnloadModel(model Model) {
	js.Global.Get("Module").Call("_UnloadModel", model)
}

// UnloadMesh - Unload mesh from memory (RAM and/or VRAM)
func UnloadMesh(mesh *Mesh) {
	js.Global.Get("Module").Call("_UnloadMesh", mesh)
}

// GenMeshPlane - Generate plane mesh (with subdivisions)
func GenMeshPlane(width, length float32, resX, resZ int) Mesh {
	return newMeshFromPointer(unsafe.Pointer(js.Global.Get("Module").Call("_GenMeshPlane", width, length, resX, resZ).Unsafe()))
}

// GenMeshCube - Generate cuboid mesh
func GenMeshCube(width, height, length float32) Mesh {
	return newMeshFromPointer(unsafe.Pointer(js.Global.Get("Module").Call("_GenMeshCube", width, height, length).Unsafe()))
}

// GenMeshSphere - Generate sphere mesh (standard sphere)
func GenMeshSphere(radius float32, rings, slices int) Mesh {
	return newMeshFromPointer(unsafe.Pointer(js.Global.Get("Module").Call("_GenMeshSphere", radius, rings, slices).Unsafe()))
}

// GenMeshHemiSphere - Generate half-sphere mesh (no bottom cap)
func GenMeshHemiSphere(radius float32, rings, slices int) Mesh {
	return newMeshFromPointer(unsafe.Pointer(js.Global.Get("Module").Call("_GenMeshHemiSphere", radius, rings, slices).Unsafe()))
}

// GenMeshCylinder - Generate cylinder mesh
func GenMeshCylinder(radius, height float32, slices int) Mesh {
	return newMeshFromPointer(unsafe.Pointer(js.Global.Get("Module").Call("_GenMeshCylinder", radius, height, slices).Unsafe()))
}

// GenMeshTorus - Generate torus mesh
func GenMeshTorus(radius, size float32, radSeg, sides int) Mesh {
	return newMeshFromPointer(unsafe.Pointer(js.Global.Get("Module").Call("_GenMeshTorus", radius, size, radSeg, sides).Unsafe()))
}

// GenMeshKnot - Generate trefoil knot mesh
func GenMeshKnot(radius, size float32, radSeg, sides int) Mesh {
	return newMeshFromPointer(unsafe.Pointer(js.Global.Get("Module").Call("_GenMeshKnot", radius, size, radSeg, sides).Unsafe()))
}

// GenMeshHeightmap - Generate heightmap mesh from image data
func GenMeshHeightmap(heightmap Image, size Vector3) Mesh {
	return newMeshFromPointer(unsafe.Pointer(js.Global.Get("Module").Call("_GenMeshHeightmap", heightmap, size).Unsafe()))
}

// GenMeshCubicmap - Generate cubes-based map mesh from image data
func GenMeshCubicmap(cubicmap Image, size Vector3) Mesh {
	return newMeshFromPointer(unsafe.Pointer(js.Global.Get("Module").Call("_GenMeshCubicmap", cubicmap, size).Unsafe()))
}

// LoadMaterial - Load material data (.MTL)
func LoadMaterial(fileName string) Material {
	return newMaterialFromPointer(unsafe.Pointer(js.Global.Get("Module").Call("_LoadMaterial", fileName).Unsafe()))
}

// LoadMaterialDefault - Load default material (Supports: DIFFUSE, SPECULAR, NORMAL maps)
func LoadMaterialDefault() Material {
	return newMaterialFromPointer(unsafe.Pointer(js.Global.Get("Module").Call("_LoadMaterialDefault").Unsafe()))
}

// UnloadMaterial - Unload material textures from VRAM
func UnloadMaterial(material Material) {
	js.Global.Get("Module").Call("_UnloadMaterial", material)
}

// DrawModel - Draw a model (with texture if set)
func DrawModel(model Model, position Vector3, scale float32, tint Color) {
	js.Global.Get("Module").Call("_DrawModel", model, position, scale, tint)
}

// DrawModelEx - Draw a model with extended parameters
func DrawModelEx(model Model, position Vector3, rotationAxis Vector3, rotationAngle float32, scale Vector3, tint Color) {
	js.Global.Get("Module").Call("_DrawModelEx", model, position, rotationAxis, rotationAngle, scale, tint)
}

// DrawModelWires - Draw a model wires (with texture if set)
func DrawModelWires(model Model, position Vector3, scale float32, tint Color) {
	js.Global.Get("Module").Call("_DrawModelWires", model, position, scale, tint)
}

// DrawModelWiresEx - Draw a model wires (with texture if set) with extended parameters
func DrawModelWiresEx(model Model, position Vector3, rotationAxis Vector3, rotationAngle float32, scale Vector3, tint Color) {
	js.Global.Get("Module").Call("_DrawModelWiresEx", model, position, rotationAxis, rotationAngle, scale, tint)
}

// DrawBoundingBox - Draw bounding box (wires)
func DrawBoundingBox(box BoundingBox, color Color) {
	js.Global.Get("Module").Call("_DrawBoundingBox", box, color)
}

// DrawBillboard - Draw a billboard texture
func DrawBillboard(camera Camera, texture Texture2D, center Vector3, size float32, tint Color) {
	js.Global.Get("Module").Call("_DrawBillboard", camera, texture, center, size, tint)
}

// DrawBillboardRec - Draw a billboard texture defined by sourceRec
func DrawBillboardRec(camera Camera, texture Texture2D, sourceRec Rectangle, center Vector3, size float32, tint Color) {
	js.Global.Get("Module").Call("_DrawBillboardRec", camera, texture, sourceRec, center, size, tint)
}

// CalculateBoundingBox - Calculate mesh bounding box limits
func CalculateBoundingBox(mesh Mesh) BoundingBox {
	return newBoundingBoxFromPointer(unsafe.Pointer(js.Global.Get("Module").Call("_CalculateBoundingBox").Unsafe()))
}

// CheckCollisionSpheres - Detect collision between two spheres
func CheckCollisionSpheres(centerA Vector3, radiusA float32, centerB Vector3, radiusB float32) bool {
	return js.Global.Get("Module").Call("_CheckCollisionSpheres", centerA, radiusA, centerB, radiusB).Bool()
}

// CheckCollisionBoxes - Detect collision between two bounding boxes
func CheckCollisionBoxes(box1 BoundingBox, box2 BoundingBox) bool {
	return js.Global.Get("Module").Call("_CheckCollisionBoxes", box1, box2).Bool()
}

// CheckCollisionBoxSphere - Detect collision between box and sphere
func CheckCollisionBoxSphere(box BoundingBox, centerSphere Vector3, radiusSphere float32) bool {
	return js.Global.Get("Module").Call("_CheckCollisionBoxSphere", box, centerSphere, radiusSphere).Bool()
}

// CheckCollisionRaySphere - Detect collision between ray and sphere
func CheckCollisionRaySphere(ray Ray, spherePosition Vector3, sphereRadius float32) bool {
	return js.Global.Get("Module").Call("_CheckCollisionRaySphere", ray, spherePosition, sphereRadius).Bool()
}

// CheckCollisionRaySphereEx - Detect collision between ray and sphere with extended parameters and collision point detection
func CheckCollisionRaySphereEx(ray Ray, spherePosition Vector3, sphereRadius float32, collisionPoint Vector3) bool {
	return js.Global.Get("Module").Call("_CheckCollisionRaySphereEx", ray, spherePosition, sphereRadius, collisionPoint).Bool()
}

// CheckCollisionRayBox - Detect collision between ray and box
func CheckCollisionRayBox(ray Ray, box BoundingBox) bool {
	return js.Global.Get("Module").Call("_CheckCollisionRayBox", ray, box).Bool()
}

// LoadShader - Load a custom shader and bind default locations
func LoadShader(vsFileName string, fsFileName string) Shader {
	return newShaderFromPointer(unsafe.Pointer(js.Global.Get("Module").Call("_LoadShader", vsFileName, fsFileName).Unsafe()))
}

// UnloadShader - Unload a custom shader from memory
func UnloadShader(shader Shader) {
	js.Global.Get("Module").Call("_UnloadShader", shader)
}

// GetShaderDefault - Get default shader
func GetShaderDefault() Shader {
	return newShaderFromPointer(unsafe.Pointer(js.Global.Get("Module").Call("_GetShaderDefault").Unsafe()))
}

// GetTextureDefault - Get default texture
func GetTextureDefault() *Texture2D {
	v := newTexture2DFromPointer(unsafe.Pointer(js.Global.Get("Module").Call("_GetTextureDefault").Unsafe()))
	return &v
}

// GetShaderLocation - Get shader uniform location
func GetShaderLocation(shader Shader, uniformName string) int32 {
	return int32(js.Global.Get("Module").Call("_GetShaderLocation", shader, uniformName).Int())
}

// SetShaderValue - Set shader uniform value (float)
func SetShaderValue(shader Shader, uniformLoc int32, value []float32, size int32) {
	js.Global.Get("Module").Call("_SetShaderValue", shader, uniformLoc, value, size)
}

// SetShaderValuei - Set shader uniform value (int)
func SetShaderValuei(shader Shader, uniformLoc int32, value []int32, size int32) {
	js.Global.Get("Module").Call("_SetShaderValuei", shader, uniformLoc, value, size)
}

// SetShaderValueMatrix - Set shader uniform value (matrix 4x4)
func SetShaderValueMatrix(shader Shader, uniformLoc int32, mat Matrix) {
	js.Global.Get("Module").Call("_SetShaderValueMatrix", shader, uniformLoc, mat)
}

// SetMatrixProjection - Set a custom projection matrix (replaces internal projection matrix)
func SetMatrixProjection(proj Matrix) {
	js.Global.Get("Module").Call("_SetMatrixProjection", proj)
}

// SetMatrixModelview - Set a custom modelview matrix (replaces internal modelview matrix)
func SetMatrixModelview(view Matrix) {
	js.Global.Get("Module").Call("_SetMatrixModelview", view)
}

// GenTextureCubemap - Generate cubemap texture from HDR texture
func GenTextureCubemap(shader Shader, skyHDR Texture2D, size int) Texture2D {
	return newTexture2DFromPointer(unsafe.Pointer(js.Global.Get("Module").Call("_GenTextureCubemap", shader, skyHDR, size).Unsafe()))
}

// GenTextureIrradiance - Generate irradiance texture using cubemap data
func GenTextureIrradiance(shader Shader, cubemap Texture2D, size int) Texture2D {
	return newTexture2DFromPointer(unsafe.Pointer(js.Global.Get("Module").Call("_GenTextureIrradiance", shader, cubemap, size).Unsafe()))
}

// GenTexturePrefilter - Generate prefilter texture using cubemap data
func GenTexturePrefilter(shader Shader, cubemap Texture2D, size int) Texture2D {
	return newTexture2DFromPointer(unsafe.Pointer(js.Global.Get("Module").Call("_GenTexturePrefilter", shader, cubemap, size).Unsafe()))
}

// GenTextureBRDF - Generate BRDF texture using cubemap data
func GenTextureBRDF(shader Shader, cubemap Texture2D, size int) Texture2D {
	return newTexture2DFromPointer(unsafe.Pointer(js.Global.Get("Module").Call("_GenTextureBRDF", shader, cubemap, size).Unsafe()))
}

// BeginShaderMode - Begin custom shader drawing
func BeginShaderMode(shader Shader) {
	js.Global.Get("Module").Call("_BeginShaderMode", shader)
}

// EndShaderMode - End custom shader drawing (use default shader)
func EndShaderMode() {
	js.Global.Get("Module").Call("_EndShaderMode")
}

// BeginBlendMode - Begin blending mode (alpha, additive, multiplied)
func BeginBlendMode(mode BlendMode) {
	js.Global.Get("Module").Call("_BeginBlendMode", mode)
}

// EndBlendMode - End blending mode (reset to default: alpha blending)
func EndBlendMode() {
	js.Global.Get("Module").Call("_EndBlendMode")
}

// GetVrDeviceInfo - Get VR device information for some standard devices
func GetVrDeviceInfo(vrDevice VrDevice) VrDeviceInfo {
	return newVrDeviceInfoFromPointer(unsafe.Pointer(js.Global.Get("Module").Call("_GetVrDeviceInfo", vrDevice).Unsafe()))
}

// InitVrSimulator - Init VR simulator for selected device
func InitVrSimulator(vrDeviceInfo VrDeviceInfo) {
	js.Global.Get("Module").Call("_InitVrSimulator", vrDeviceInfo)
}

// CloseVrSimulator - Close VR simulator for current device
func CloseVrSimulator() {
	js.Global.Get("Module").Call("_CloseVrSimulator")
}

// IsVrSimulatorReady - Detect if VR simulator is ready
func IsVrSimulatorReady() bool {
	return js.Global.Get("Module").Call("_IsVrSimulatorReady").Bool()
}

// UpdateVrTracking - Update VR tracking (position and orientation) and camera
func UpdateVrTracking(camera *Camera) {
	js.Global.Get("Module").Call("_UpdateVrTracking", camera)
}

// ToggleVrMode - Enable/Disable VR experience (device or simulator)
func ToggleVrMode() {
	js.Global.Get("Module").Call("_ToggleVrMode")
}

// BeginVrDrawing - Begin VR simulator stereo rendering
func BeginVrDrawing() {
	js.Global.Get("Module").Call("_BeginVrDrawing")
}

// EndVrDrawing - End VR simulator stereo rendering
func EndVrDrawing() {
	js.Global.Get("Module").Call("_EndVrDrawing")
}

// DrawPixel - Draw a pixel
func DrawPixel(posX, posY int32, color Color) {
	js.Global.Get("Module").Call("_DrawPixel", posX, posY, color)
}

// DrawPixelV - Draw a pixel (Vector version)
func DrawPixelV(position Vector2, color Color) {
	js.Global.Get("Module").Call("_DrawPixelV", position, color)
}

// DrawLine - Draw a line
func DrawLine(startPosX, startPosY, endPosX, endPosY int32, color Color) {
	js.Global.Get("Module").Call("_DrawLine", startPosX, startPosY, endPosX, endPosY, color)
}

// DrawLineV - Draw a line (Vector version)
func DrawLineV(startPos, endPos Vector2, color Color) {
	js.Global.Get("Module").Call("_DrawLineV", startPos, endPos, color)
}

// DrawLineEx - Draw a line defining thickness
func DrawLineEx(startPos, endPos Vector2, thick float32, color Color) {
	js.Global.Get("Module").Call("_DrawLineEx", startPos, endPos, thick, color)
}

// DrawLineBezier - Draw a line using cubic-bezier curves in-out
func DrawLineBezier(startPos, endPos Vector2, thick float32, color Color) {
	js.Global.Get("Module").Call("_DrawLineBezier", startPos, endPos, thick, color)
}

// DrawCircle - Draw a color-filled circle
func DrawCircle(centerX, centerY int32, radius float32, color Color) {
	js.Global.Get("Module").Call("_DrawCircle", centerX, centerY, radius, color)
}

// DrawCircleGradient - Draw a gradient-filled circle
func DrawCircleGradient(centerX, centerY int32, radius float32, color1, color2 Color) {
	js.Global.Get("Module").Call("_DrawCircleGradient", centerX, centerY, radius, color1, color2)
}

// DrawCircleV - Draw a color-filled circle (Vector version)
func DrawCircleV(center Vector2, radius float32, color Color) {
	js.Global.Get("Module").Call("_DrawCircleV", center, radius, color)
}

// DrawCircleLines - Draw circle outline
func DrawCircleLines(centerX, centerY int32, radius float32, color Color) {
	js.Global.Get("Module").Call("_DrawCircleLines", centerX, centerY, radius, color)
}

// DrawRectangle - Draw a color-filled rectangle
func DrawRectangle(posX, posY, width, height int32, color Color) {
	js.Global.Get("Module").Call("_DrawRectangle", posX, posY, width, height, color)
}

// DrawRectangleRec - Draw a color-filled rectangle
func DrawRectangleRec(rec Rectangle, color Color) {
	js.Global.Get("Module").Call("_DrawRectangleRec", rec, color)
}

// DrawRectanglePro - Draw a color-filled rectangle with pro parameters
func DrawRectanglePro(rec Rectangle, origin Vector2, rotation float32, color Color) {
	js.Global.Get("Module").Call("_DrawRectanglePro", rec, origin, rotation, color)
}

// DrawRectangleGradientV - Draw a vertical-gradient-filled rectangle
func DrawRectangleGradientV(posX, posY, width, height int32, color1, color2 Color) {
	js.Global.Get("Module").Call("_DrawRectangleGradientV", posX, posY, width, height, color1, color2)
}

// DrawRectangleGradientH - Draw a horizontal-gradient-filled rectangle
func DrawRectangleGradientH(posX, posY, width, height int32, color1, color2 Color) {
	js.Global.Get("Module").Call("_DrawRectangleGradientH", posX, posY, width, height, color1, color2)
}

// DrawRectangleGradientEx - Draw a gradient-filled rectangle with custom vertex colors
func DrawRectangleGradientEx(rec Rectangle, color1, color2, color3, color4 Color) {
	js.Global.Get("Module").Call("_DrawRectangleGradientEx", rec, color1, color2, color3, color4)
}

// DrawRectangleV - Draw a color-filled rectangle (Vector version)
func DrawRectangleV(position Vector2, size Vector2, color Color) {
	js.Global.Get("Module").Call("_DrawRectangleV", position, size, color)
}

// DrawRectangleLines - Draw rectangle outline
func DrawRectangleLines(posX, posY, width, height int32, color Color) {
	js.Global.Get("Module").Call("_DrawRectangleLines", posX, posY, width, height, color)
}

// DrawRectangleT - Draw rectangle using text character
func DrawRectangleT(posX, posY, width, height int32, color Color) {
	js.Global.Get("Module").Call("_DrawRectangleT", posX, posY, width, height, color)
}

// DrawTriangle - Draw a color-filled triangle
func DrawTriangle(v1, v2, v3 Vector2, color Color) {
	js.Global.Get("Module").Call("_DrawTriangle", v1, v2, v3, color)
}

// DrawTriangleLines - Draw triangle outline
func DrawTriangleLines(v1, v2, v3 Vector2, color Color) {
	js.Global.Get("Module").Call("_DrawTriangleLines", v1, v2, v3, color)
}

// DrawPoly - Draw a regular polygon (Vector version)
func DrawPoly(center Vector2, sides int32, radius, rotation float32, color Color) {
	js.Global.Get("Module").Call("_DrawPoly", center, sides, radius, rotation, color)
}

// DrawPolyEx - Draw a closed polygon defined by points
func DrawPolyEx(points []Vector2, numPoints int32, color Color) {
	js.Global.Get("Module").Call("_DrawPolyEx", points, numPoints, color)
}

// DrawPolyExLines - Draw polygon lines
func DrawPolyExLines(points []Vector2, numPoints int32, color Color) {
	js.Global.Get("Module").Call("_DrawPolyExLines", points, numPoints, color)
}

// CheckCollisionRecs - Check collision between two rectangles
func CheckCollisionRecs(rec1, rec2 Rectangle) bool {
	return js.Global.Get("Module").Call("_CheckCollisionRecs", rec1, rec2).Bool()
}

// CheckCollisionCircles - Check collision between two circles
func CheckCollisionCircles(center1 Vector2, radius1 float32, center2 Vector2, radius2 float32) bool {
	return js.Global.Get("Module").Call("_CheckCollisionCircles", center1, radius1, center2, radius2).Bool()
}

// CheckCollisionCircleRec - Check collision between circle and rectangle
func CheckCollisionCircleRec(center Vector2, radius float32, rec Rectangle) bool {
	return js.Global.Get("Module").Call("_CheckCollisionCircleRec", center, radius, rec).Bool()
}

// GetCollisionRec - Get collision rectangle for two rectangles collision
func GetCollisionRec(rec1, rec2 Rectangle) Rectangle {
	return newRectangleFromPointer(unsafe.Pointer(js.Global.Get("Module").Call("_GetCollisionRec", rec1, rec2).Unsafe()))
}

// CheckCollisionPointRec - Check if point is inside rectangle
func CheckCollisionPointRec(point Vector2, rec Rectangle) bool {
	return js.Global.Get("Module").Call("_CheckCollisionPointRec", point, rec).Bool()
}

// CheckCollisionPointCircle - Check if point is inside circle
func CheckCollisionPointCircle(point Vector2, center Vector2, radius float32) bool {
	return js.Global.Get("Module").Call("_CheckCollisionPointCircle", point, center, radius).Bool()
}

// CheckCollisionPointTriangle - Check if point is inside a triangle
func CheckCollisionPointTriangle(point, p1, p2, p3 Vector2) bool {
	return js.Global.Get("Module").Call("_CheckCollisionPointTriangle", point, p1, p2, p3).Bool()
}

// GetDefaultFont - Get the default SpriteFont
func GetDefaultFont() SpriteFont {
	return newSpriteFontFromPointer(unsafe.Pointer(js.Global.Get("Module").Call("_GetDefaultFont").Unsafe()))
}

// LoadSpriteFont - Load a SpriteFont image into GPU memory (VRAM)
func LoadSpriteFont(fileName string) SpriteFont {
	return newSpriteFontFromPointer(unsafe.Pointer(js.Global.Get("Module").Call("_LoadSpriteFont", fileName).Unsafe()))
}

// LoadSpriteFontEx - Load SpriteFont from file with extended parameters
func LoadSpriteFontEx(fileName string, fontSize int32, charsCount int32, fontChars *int32) SpriteFont {
	return newSpriteFontFromPointer(unsafe.Pointer(js.Global.Get("Module").Call("_LoadSpriteFontEx", fileName, fontSize, charsCount, fontChars).Unsafe()))
}

// UnloadSpriteFont - Unload SpriteFont from GPU memory (VRAM)
func UnloadSpriteFont(spriteFont SpriteFont) {
	js.Global.Get("Module").Call("_UnloadSpriteFont", spriteFont)
}

// DrawText - Draw text (using default font)
func DrawText(text string, posX int32, posY int32, fontSize int32, color Color) {
	js.Global.Get("Module").Call("_DrawText", text, posX, posY, fontSize, color)
}

// DrawTextEx - Draw text using SpriteFont and additional parameters
func DrawTextEx(spriteFont SpriteFont, text string, position Vector2, fontSize float32, spacing int32, tint Color) {
	js.Global.Get("Module").Call("_DrawTextEx", spriteFont, text, position, fontSize, spacing, tint)
}

// MeasureText - Measure string width for default font
func MeasureText(text string, fontSize int32) int32 {
	return int32(js.Global.Get("Module").Call("_MeasureText", text, fontSize).Int())
}

// MeasureTextEx - Measure string size for SpriteFont
func MeasureTextEx(spriteFont SpriteFont, text string, fontSize float32, spacing int32) Vector2 {
	return newVector2FromPointer(unsafe.Pointer(js.Global.Get("Module").Call("_MeasureTextEx", spriteFont, text, fontSize, spacing).Unsafe()))
}

// DrawFPS - Shows current FPS
func DrawFPS(posX int32, posY int32) {
	js.Global.Get("Module").Call("_DrawFPS", posX, posY)
}

// LoadImage - Load an image into CPU memory (RAM)
func LoadImage(fileName string) *Image {
	return newImageFromPointer(unsafe.Pointer(js.Global.Get("Module").Call("_LoadImage", fileName).Unsafe()))
}

// LoadImageEx - Load image data from Color array data (RGBA - 32bit)
func LoadImageEx(pixels []Color, width, height int32) *Image {
	return newImageFromPointer(unsafe.Pointer(js.Global.Get("Module").Call("_LoadImageEx", pixels, width, height).Unsafe()))
}

// LoadImagePro - Load image from raw data with parameters
func LoadImagePro(data []byte, width, height int32, format TextureFormat) *Image {
	return newImageFromPointer(unsafe.Pointer(js.Global.Get("Module").Call("_LoadImagePro", data, width, height, format).Unsafe()))
}

// LoadImageRaw - Load image data from RAW file
func LoadImageRaw(fileName string, width, height int32, format TextureFormat, headerSize int32) *Image {
	return newImageFromPointer(unsafe.Pointer(js.Global.Get("Module").Call("_LoadImageRaw", fileName, width, height, format, headerSize).Unsafe()))
}

// LoadTexture - Load an image as texture into GPU memory
func LoadTexture(fileName string) Texture2D {
	return newTexture2DFromPointer(unsafe.Pointer(js.Global.Get("Module").Call("_LoadTexture", fileName).Unsafe()))
}

// LoadTextureFromImage - Load a texture from image data
func LoadTextureFromImage(image *Image) Texture2D {
	return newTexture2DFromPointer(unsafe.Pointer(js.Global.Get("Module").Call("_LoadTextureFromImage", image).Unsafe()))
}

// LoadRenderTexture - Load a texture to be used for rendering
func LoadRenderTexture(width, height int32) RenderTexture2D {
	return newRenderTexture2DFromPointer(unsafe.Pointer(js.Global.Get("Module").Call("_LoadRenderTexture", width, height).Unsafe()))
}

// UnloadImage - Unload image from CPU memory (RAM)
func UnloadImage(image *Image) {
	js.Global.Get("Module").Call("_UnloadImage", image)
}

// UnloadTexture - Unload texture from GPU memory
func UnloadTexture(texture Texture2D) {
	js.Global.Get("Module").Call("_UnloadTexture", texture)
}

// UnloadRenderTexture - Unload render texture from GPU memory
func UnloadRenderTexture(target RenderTexture2D) {
	js.Global.Get("Module").Call("_UnloadRenderTexture", target)
}

// GetImageData - Get pixel data from image
func GetImageData(image *Image) []byte {
	return js.Global.Get("Module").Call("_GetImageData", image).Interface().([]byte)
}

// GetTextureData - Get pixel data from GPU texture and return an Image
func GetTextureData(texture Texture2D) *Image {
	return newImageFromPointer(unsafe.Pointer(js.Global.Get("Module").Call("_GetTextureData", texture).Unsafe()))
}

// UpdateTexture - Update GPU texture with new data
func UpdateTexture(texture Texture2D, pixels []byte) {
	js.Global.Get("Module").Call("_UpdateTexture", texture, pixels)
}

// SaveImageAs - Save image to a PNG file
func SaveImageAs(name string, image Image) {
	js.Global.Get("Module").Call("_SaveImageAs", name, image)
}

// ImageToPOT - Convert image to POT (power-of-two)
func ImageToPOT(image *Image, fillColor Color) {
	js.Global.Get("Module").Call("_ImageToPot", image, fillColor)
}

// ImageFormat - Convert image data to desired format
func ImageFormat(image *Image, newFormat int32) {
	js.Global.Get("Module").Call("_ImageFormat", image, newFormat)
}

// ImageAlphaMask - Apply alpha mask to image
func ImageAlphaMask(image, alphaMask *Image) {
	js.Global.Get("Module").Call("_ImageAlphaMask", image, alphaMask)
}

// ImageDither - Dither image data to 16bpp or lower (Floyd-Steinberg dithering)
func ImageDither(image *Image, rBpp, gBpp, bBpp, aBpp int32) {
	js.Global.Get("Module").Call("_ImageDither", image, rBpp, gBpp, bBpp, aBpp)
}

// ImageCopy - Create an image duplicate (useful for transformations)
func ImageCopy(image *Image) *Image {
	return newImageFromPointer(unsafe.Pointer(js.Global.Get("Module").Call("_ImageCopy", image).Unsafe()))
}

// ImageCrop - Crop an image to a defined rectangle
func ImageCrop(image *Image, crop Rectangle) {
	js.Global.Get("Module").Call("_ImageCrop", image, crop)
}

// ImageResize - Resize an image (bilinear filtering)
func ImageResize(image *Image, newWidth, newHeight int32) {
	js.Global.Get("Module").Call("_ImageResize", image, newWidth, newHeight)
}

// ImageResizeNN - Resize an image (Nearest-Neighbor scaling algorithm)
func ImageResizeNN(image *Image, newWidth, newHeight int32) {
	js.Global.Get("Module").Call("_ImageResizeNN", image, newWidth, newHeight)
}

// ImageText - Create an image from text (default font)
func ImageText(text string, fontSize int32, color Color) *Image {
	return newImageFromPointer(unsafe.Pointer(js.Global.Get("Module").Call("_ImageText", text, fontSize, color).Unsafe()))
}

// ImageTextEx - Create an image from text (custom sprite font)
func ImageTextEx(font SpriteFont, text string, fontSize float32, spacing int32, tint Color) *Image {
	return newImageFromPointer(unsafe.Pointer(js.Global.Get("Module").Call("_ImageTextEx", font, text, fontSize, spacing, tint).Unsafe()))
}

// ImageDraw - Draw a source image within a destination image
func ImageDraw(dst, src *Image, srcRec, dstRec Rectangle) {
	js.Global.Get("Module").Call("_ImageDraw", dst, src, srcRec, dstRec)
}

// ImageDrawText - Draw text (default font) within an image (destination)
func ImageDrawText(dst *Image, position Vector2, text string, fontSize int32, color Color) {
	js.Global.Get("Module").Call("_ImageDrawText", dst, position, text, fontSize, color)
}

// ImageDrawTextEx - Draw text (custom sprite font) within an image (destination)
func ImageDrawTextEx(dst *Image, position Vector2, font SpriteFont, text string, fontSize float32, spacing int32, color Color) {
	js.Global.Get("Module").Call("_ImageDrawTextEx", dst, position, font, text, fontSize, spacing, color)
}

// ImageFlipVertical - Flip image vertically
func ImageFlipVertical(image *Image) {
	js.Global.Get("Module").Call("_ImageFlipVertical", image)
}

// ImageFlipHorizontal - Flip image horizontally
func ImageFlipHorizontal(image *Image) {
	js.Global.Get("Module").Call("_ImageFlipHorizontal", image)
}

// ImageColorTint - Modify image color: tint
func ImageColorTint(image *Image, color Color) {
	js.Global.Get("Module").Call("_ImageColorTint", image, color)
}

// ImageColorInvert - Modify image color: invert
func ImageColorInvert(image *Image) {
	js.Global.Get("Module").Call("_ImageColorInvert", image)
}

// ImageColorGrayscale - Modify image color: grayscale
func ImageColorGrayscale(image *Image) {
	js.Global.Get("Module").Call("_ImageColorGrayscale", image)
}

// ImageColorContrast - Modify image color: contrast (-100 to 100)
func ImageColorContrast(image *Image, contrast float32) {
	js.Global.Get("Module").Call("_ImageColorContrast", image, contrast)
}

// ImageColorBrightness - Modify image color: brightness (-255 to 255)
func ImageColorBrightness(image *Image, brightness int32) {
	js.Global.Get("Module").Call("_ImageColorBrightness", image, brightness)
}

// GenImageColor - Generate image: plain color
func GenImageColor(width, height int, color Color) *Image {
	return newImageFromPointer(unsafe.Pointer(js.Global.Get("Module").Call("_GenImageColor", width, height, color).Unsafe()))
}

// GenImageGradientV - Generate image: vertical gradient
func GenImageGradientV(width, height int, top, bottom Color) *Image {
	return newImageFromPointer(unsafe.Pointer(js.Global.Get("Module").Call("_GenImageGradientV", width, height, top, bottom).Unsafe()))
}

// GenImageGradientH - Generate image: horizontal gradient
func GenImageGradientH(width, height int, left, right Color) *Image {
	return newImageFromPointer(unsafe.Pointer(js.Global.Get("Module").Call("_GenImageGradientH", width, height, left, right).Unsafe()))
}

// GenImageGradientRadial - Generate image: radial gradient
func GenImageGradientRadial(width, height int, density float32, inner, outer Color) *Image {
	return newImageFromPointer(unsafe.Pointer(js.Global.Get("Module").Call("_GenImageGradientRadial", width, height, density, inner, outer).Unsafe()))
}

// GenImageChecked - Generate image: checked
func GenImageChecked(width, height, checksX, checksY int, col1, col2 Color) *Image {
	return newImageFromPointer(unsafe.Pointer(js.Global.Get("Module").Call("_GenImageChecked", width, height, checksX, checksY, col1, col2).Unsafe()))
}

// GenImageWhiteNoise - Generate image: white noise
func GenImageWhiteNoise(width, height int, factor float32) *Image {
	return newImageFromPointer(unsafe.Pointer(js.Global.Get("Module").Call("_GenImageWhiteNoise", width, height, factor).Unsafe()))
}

// GenImagePerlinNoise - Generate image: perlin noise
func GenImagePerlinNoise(width, height int, scale float32) *Image {
	return newImageFromPointer(unsafe.Pointer(js.Global.Get("Module").Call("_GenImagePerlinNoise", width, height, scale).Unsafe()))
}

// GenImageCellular - Generate image: cellular algorithm. Bigger tileSize means bigger cells
func GenImageCellular(width, height, tileSize int) *Image {
	return newImageFromPointer(unsafe.Pointer(js.Global.Get("Module").Call("_GenImageCellular", width, height, tileSize).Unsafe()))
}

// GenTextureMipmaps - Generate GPU mipmaps for a texture
func GenTextureMipmaps(texture *Texture2D) {
	js.Global.Get("Module").Call("_GenTextureMipmaps", texture)
}

// SetTextureFilter - Set texture scaling filter mode
func SetTextureFilter(texture Texture2D, filterMode TextureFilterMode) {
	js.Global.Get("Module").Call("_SetTextureFilter", texture, filterMode)
}

// SetTextureWrap - Set texture wrapping mode
func SetTextureWrap(texture Texture2D, wrapMode TextureWrapMode) {
	js.Global.Get("Module").Call("_SetTextureWrap", texture, wrapMode)
}

// DrawTexture - Draw a Texture2D
func DrawTexture(texture Texture2D, posX int32, posY int32, tint Color) {
	js.Global.Get("Module").Call("_DrawTexture", texture, posX, posY, tint)
}

// DrawTextureV - Draw a Texture2D with position defined as Vector2
func DrawTextureV(texture Texture2D, position Vector2, tint Color) {
	js.Global.Get("Module").Call("_DrawTextureV", texture, position, tint)
}

// DrawTextureEx - Draw a Texture2D with extended parameters
func DrawTextureEx(texture Texture2D, position Vector2, rotation, scale float32, tint Color) {
	js.Global.Get("Module").Call("_DrawTextureEx", texture, position, rotation, scale, tint)
}

// DrawTextureRec - Draw a part of a texture defined by a rectangle
func DrawTextureRec(texture Texture2D, sourceRec Rectangle, position Vector2, tint Color) {
	js.Global.Get("Module").Call("_DrawTextureRec", texture, sourceRec, position, tint)
}

// DrawTexturePro - Draw a part of a texture defined by a rectangle with 'pro' parameters
func DrawTexturePro(texture Texture2D, sourceRec, destRec Rectangle, origin Vector2, rotation float32, tint Color) {
	js.Global.Get("Module").Call("_DrawTexturePro", texture, sourceRec, destRec, origin, rotation, tint)
}
