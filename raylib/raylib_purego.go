//go:build !cgo
// +build !cgo

package rl

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"reflect"
	"unsafe"

	"github.com/ebitengine/purego"
	"github.com/gen2brain/raylib-go/raylib/internal/convert"
	"github.com/jupiterrider/ffi"
)

var (
	// dll is the pointer to the shared library
	dll ffi.Lib = loadLibrary()

	// audioCallbacks is needed to have a reference between Go functions (keys) created by the user
	// and C function pointers (values) created by purego.NewCallback
	audioCallbacks map[uintptr]uintptr = make(map[uintptr]uintptr)

	// Window-related functions

	initWindow               = dll.MustPrep("InitWindow", &ffi.TypeVoid, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypePointer)
	closeWindow              = dll.MustPrep("CloseWindow", &ffi.TypeVoid)
	windowShouldClose        = dll.MustPrep("WindowShouldClose", &ffi.TypeUint8)
	isWindowReady            = dll.MustPrep("IsWindowReady", &ffi.TypeUint8)
	isWindowFullscreen       = dll.MustPrep("IsWindowFullscreen", &ffi.TypeUint8)
	isWindowHidden           = dll.MustPrep("IsWindowHidden", &ffi.TypeUint8)
	isWindowMinimized        = dll.MustPrep("IsWindowMinimized", &ffi.TypeUint8)
	isWindowMaximized        = dll.MustPrep("IsWindowMaximized", &ffi.TypeUint8)
	isWindowFocused          = dll.MustPrep("IsWindowFocused", &ffi.TypeUint8)
	isWindowResized          = dll.MustPrep("IsWindowResized", &ffi.TypeUint8)
	isWindowState            = dll.MustPrep("IsWindowState", &ffi.TypeUint8, &ffi.TypeUint32)
	setWindowState           = dll.MustPrep("SetWindowState", &ffi.TypeVoid, &ffi.TypeUint32)
	clearWindowState         = dll.MustPrep("ClearWindowState", &ffi.TypeVoid, &ffi.TypeUint32)
	toggleFullscreen         = dll.MustPrep("ToggleFullscreen", &ffi.TypeVoid)
	toggleBorderlessWindowed = dll.MustPrep("ToggleBorderlessWindowed", &ffi.TypeVoid)
	maximizeWindow           = dll.MustPrep("MaximizeWindow", &ffi.TypeVoid)
	minimizeWindow           = dll.MustPrep("MinimizeWindow", &ffi.TypeVoid)
	restoreWindow            = dll.MustPrep("RestoreWindow", &ffi.TypeVoid)
	setWindowIcon            = dll.MustPrep("SetWindowIcon", &ffi.TypeVoid, &typeImage)
	setWindowIcons           = dll.MustPrep("SetWindowIcons", &ffi.TypeVoid, &ffi.TypePointer, &ffi.TypeSint32)
	setWindowTitle           = dll.MustPrep("SetWindowTitle", &ffi.TypeVoid, &ffi.TypePointer)
	setWindowPosition        = dll.MustPrep("SetWindowPosition", &ffi.TypeVoid, &ffi.TypeSint32, &ffi.TypeSint32)
	setWindowMonitor         = dll.MustPrep("SetWindowMonitor", &ffi.TypeVoid, &ffi.TypeSint32)
	setWindowMinSize         = dll.MustPrep("SetWindowMinSize", &ffi.TypeVoid, &ffi.TypeSint32, &ffi.TypeSint32)
	setWindowMaxSize         = dll.MustPrep("SetWindowMaxSize", &ffi.TypeVoid, &ffi.TypeSint32, &ffi.TypeSint32)
	setWindowSize            = dll.MustPrep("SetWindowSize", &ffi.TypeVoid, &ffi.TypeSint32, &ffi.TypeSint32)
	setWindowOpacity         = dll.MustPrep("SetWindowOpacity", &ffi.TypeVoid, &ffi.TypeFloat)
	setWindowFocused         = dll.MustPrep("SetWindowFocused", &ffi.TypeVoid)
	getWindowHandle          = dll.MustPrep("GetWindowHandle", &ffi.TypePointer)
	getScreenWidth           = dll.MustPrep("GetScreenWidth", &ffi.TypeSint32)
	getScreenHeight          = dll.MustPrep("GetScreenHeight", &ffi.TypeSint32)
	getRenderWidth           = dll.MustPrep("GetRenderWidth", &ffi.TypeSint32)
	getRenderHeight          = dll.MustPrep("GetRenderHeight", &ffi.TypeSint32)
	getMonitorCount          = dll.MustPrep("GetMonitorCount", &ffi.TypeSint32)
	getCurrentMonitor        = dll.MustPrep("GetCurrentMonitor", &ffi.TypeSint32)
	getMonitorPosition       = dll.MustPrep("GetMonitorPosition", &typeVector2, &ffi.TypeSint32)
	getMonitorWidth          = dll.MustPrep("GetMonitorWidth", &ffi.TypeSint32, &ffi.TypeSint32)
	getMonitorHeight         = dll.MustPrep("GetMonitorHeight", &ffi.TypeSint32, &ffi.TypeSint32)
	getMonitorPhysicalWidth  = dll.MustPrep("GetMonitorPhysicalWidth", &ffi.TypeSint32, &ffi.TypeSint32)
	getMonitorPhysicalHeight = dll.MustPrep("GetMonitorPhysicalHeight", &ffi.TypeSint32, &ffi.TypeSint32)
	getMonitorRefreshRate    = dll.MustPrep("GetMonitorRefreshRate", &ffi.TypeSint32, &ffi.TypeSint32)
	getWindowPosition        = dll.MustPrep("GetWindowPosition", &typeVector2)
	getWindowScaleDPI        = dll.MustPrep("GetWindowScaleDPI", &typeVector2)
	getMonitorName           = dll.MustPrep("GetMonitorName", &ffi.TypePointer, &ffi.TypeSint32)
	setClipboardText         = dll.MustPrep("SetClipboardText", &ffi.TypeVoid, &ffi.TypePointer)
	getClipboardText         = dll.MustPrep("GetClipboardText", &ffi.TypePointer)
	getClipboardImage        = dll.MustPrep("GetClipboardImage", &typeImage)
	enableEventWaiting       = dll.MustPrep("EnableEventWaiting", &ffi.TypeVoid)
	disableEventWaiting      = dll.MustPrep("DisableEventWaiting", &ffi.TypeVoid)

	// Cursor-related functions

	showCursor       = dll.MustPrep("ShowCursor", &ffi.TypeVoid)
	hideCursor       = dll.MustPrep("HideCursor", &ffi.TypeVoid)
	isCursorHidden   = dll.MustPrep("IsCursorHidden", &ffi.TypeUint8)
	enableCursor     = dll.MustPrep("EnableCursor", &ffi.TypeVoid)
	disableCursor    = dll.MustPrep("DisableCursor", &ffi.TypeVoid)
	isCursorOnScreen = dll.MustPrep("IsCursorOnScreen", &ffi.TypeUint8)

	// Drawing-related functions

	clearBackground   = dll.MustPrep("ClearBackground", &ffi.TypeVoid, &typeColor)
	beginDrawing      = dll.MustPrep("BeginDrawing", &ffi.TypeVoid)
	endDrawing        = dll.MustPrep("EndDrawing", &ffi.TypeVoid)
	beginMode2D       = dll.MustPrep("BeginMode2D", &ffi.TypeVoid, &typeCamera2D)
	endMode2D         = dll.MustPrep("EndMode2D", &ffi.TypeVoid)
	beginMode3D       = dll.MustPrep("BeginMode3D", &ffi.TypeVoid, &typeCamera3D)
	endMode3D         = dll.MustPrep("EndMode3D", &ffi.TypeVoid)
	beginTextureMode  = dll.MustPrep("BeginTextureMode", &ffi.TypeVoid, &typeRenderTexture2D)
	endTextureMode    = dll.MustPrep("EndTextureMode", &ffi.TypeVoid)
	beginShaderMode   = dll.MustPrep("BeginShaderMode", &ffi.TypeVoid, &typeShader)
	endShaderMode     = dll.MustPrep("EndShaderMode", &ffi.TypeVoid)
	beginBlendMode    = dll.MustPrep("BeginBlendMode", &ffi.TypeVoid, &ffi.TypeSint32)
	endBlendMode      = dll.MustPrep("EndBlendMode", &ffi.TypeVoid)
	beginScissorMode  = dll.MustPrep("BeginScissorMode", &ffi.TypeVoid, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32)
	endScissorMode    = dll.MustPrep("EndScissorMode", &ffi.TypeVoid)
	beginVrStereoMode = dll.MustPrep("BeginVrStereoMode", &ffi.TypeVoid, &typeVrStereoConfig)
	endVrStereoMode   = dll.MustPrep("EndVrStereoMode", &ffi.TypeVoid)

	// VR stereo config functions for VR simulator

	loadVrStereoConfig   = dll.MustPrep("LoadVrStereoConfig", &typeVrStereoConfig, &typeVrDeviceInfo)
	unloadVrStereoConfig = dll.MustPrep("UnloadVrStereoConfig", &ffi.TypeVoid, &typeVrStereoConfig)

	// Shader management functions

	loadShader              = dll.MustPrep("LoadShader", &typeShader, &ffi.TypePointer, &ffi.TypePointer)
	loadShaderFromMemory    = dll.MustPrep("LoadShaderFromMemory", &typeShader, &ffi.TypePointer, &ffi.TypePointer)
	isShaderValid           = dll.MustPrep("IsShaderValid", &ffi.TypeUint8, &typeShader)
	getShaderLocation       = dll.MustPrep("GetShaderLocation", &ffi.TypeSint32, &typeShader, &ffi.TypePointer)
	getShaderLocationAttrib = dll.MustPrep("GetShaderLocationAttrib", &ffi.TypeSint32, &typeShader, &ffi.TypePointer)
	setShaderValue          = dll.MustPrep("SetShaderValue", &ffi.TypeVoid, &typeShader, &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypeSint32)
	setShaderValueV         = dll.MustPrep("SetShaderValueV", &ffi.TypeVoid, &typeShader, &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypeSint32, &ffi.TypeSint32)
	setShaderValueMatrix    = dll.MustPrep("SetShaderValueMatrix", &ffi.TypeVoid, &typeShader, &ffi.TypeSint32, &typeMatrix)
	setShaderValueTexture   = dll.MustPrep("SetShaderValueTexture", &ffi.TypeVoid, &typeShader, &ffi.TypeSint32, &typeTexture2D)
	unloadShader            = dll.MustPrep("UnloadShader", &ffi.TypeVoid, &typeShader)

	// Screen-space-related functions

	getScreenToWorldRay   = dll.MustPrep("GetScreenToWorldRay", &typeRay, &typeVector2, &typeCamera3D)
	getScreenToWorldRayEx = dll.MustPrep("GetScreenToWorldRayEx", &typeRay, &typeVector2, &typeCamera3D, &ffi.TypeSint32, &ffi.TypeSint32)
	getWorldToScreen      = dll.MustPrep("GetWorldToScreen", &typeVector2, &typeVector3, &typeCamera3D)
	getWorldToScreenEx    = dll.MustPrep("GetWorldToScreenEx", &typeVector2, &typeVector3, &typeCamera3D, &ffi.TypeSint32, &ffi.TypeSint32)
	getWorldToScreen2D    = dll.MustPrep("GetWorldToScreen2D", &typeVector2, &typeVector2, &typeCamera2D)
	getScreenToWorld2D    = dll.MustPrep("GetScreenToWorld2D", &typeVector2, &typeVector2, &typeCamera2D)
	getCameraMatrix       = dll.MustPrep("GetCameraMatrix", &typeMatrix, &typeCamera3D)
	getCameraMatrix2D     = dll.MustPrep("GetCameraMatrix2D", &typeMatrix, &typeCamera2D)

	// Timing-related functions

	setTargetFPS = dll.MustPrep("SetTargetFPS", &ffi.TypeVoid, &ffi.TypeSint32)
	getFrameTime = dll.MustPrep("GetFrameTime", &ffi.TypeFloat)
	getTime      = dll.MustPrep("GetTime", &ffi.TypeDouble)
	getFPS       = dll.MustPrep("GetFPS", &ffi.TypeSint32)

	// Custom frame control functions

	swapScreenBuffer = dll.MustPrep("SwapScreenBuffer", &ffi.TypeVoid)
	pollInputEvents  = dll.MustPrep("PollInputEvents", &ffi.TypeVoid)
	waitTime         = dll.MustPrep("WaitTime", &ffi.TypeVoid, &ffi.TypeDouble)

	// Random values generation functions

	getRandomValue = dll.MustPrep("GetRandomValue", &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32)

	// Misc. functions

	takeScreenshot = dll.MustPrep("TakeScreenshot", &ffi.TypeVoid, &ffi.TypePointer)
	setConfigFlags = dll.MustPrep("SetConfigFlags", &ffi.TypeVoid, &ffi.TypeUint32)
	openURL        = dll.MustPrep("OpenURL", &ffi.TypeVoid, &ffi.TypePointer)

	// Logging system

	setTraceLogLevel    = dll.MustPrep("SetTraceLogLevel", &ffi.TypeVoid, &ffi.TypeSint32)
	traceLog            = dll.MustPrepVar("TraceLog", 2, &ffi.TypeVoid, &ffi.TypeSint32, &ffi.TypePointer)
	setTraceLogCallback = dll.MustPrep("SetTraceLogCallback", &ffi.TypeVoid, &ffi.TypePointer)

	// Memory management, using internal allocators

	memFree = dll.MustPrep("MemFree", &ffi.TypeVoid, &ffi.TypePointer)

	// File system management functions

	isFileDropped      = dll.MustPrep("IsFileDropped", &ffi.TypeUint8)
	loadDroppedFiles   = dll.MustPrep("LoadDroppedFiles", &typeFilePathList)
	unloadDroppedFiles = dll.MustPrep("UnloadDroppedFiles", &ffi.TypeVoid, &typeFilePathList)

	// Automation events functionality

	loadAutomationEventList       = dll.MustPrep("LoadAutomationEventList", &typeAutomationEventList, &ffi.TypePointer)
	unloadAutomationEventList     = dll.MustPrep("UnloadAutomationEventList", &ffi.TypeVoid, &typeAutomationEventList)
	exportAutomationEventList     = dll.MustPrep("ExportAutomationEventList", &ffi.TypeUint8, &typeAutomationEventList, &ffi.TypePointer)
	setAutomationEventList        = dll.MustPrep("SetAutomationEventList", &ffi.TypeVoid, &ffi.TypePointer)
	setAutomationEventBaseFrame   = dll.MustPrep("SetAutomationEventBaseFrame", &ffi.TypeVoid, &ffi.TypeSint32)
	startAutomationEventRecording = dll.MustPrep("StartAutomationEventRecording", &ffi.TypeVoid)
	stopAutomationEventRecording  = dll.MustPrep("StopAutomationEventRecording", &ffi.TypeVoid)
	playAutomationEvent           = dll.MustPrep("PlayAutomationEvent", &ffi.TypeVoid, &typeAutomationEvent)

	// Input-related functions: keyboard

	isKeyPressed       = dll.MustPrep("IsKeyPressed", &ffi.TypeUint8, &ffi.TypeSint32)
	isKeyPressedRepeat = dll.MustPrep("IsKeyPressedRepeat", &ffi.TypeUint8, &ffi.TypeSint32)
	isKeyDown          = dll.MustPrep("IsKeyDown", &ffi.TypeUint8, &ffi.TypeSint32)
	isKeyReleased      = dll.MustPrep("IsKeyReleased", &ffi.TypeUint8, &ffi.TypeSint32)
	isKeyUp            = dll.MustPrep("IsKeyUp", &ffi.TypeUint8, &ffi.TypeSint32)
	getKeyPressed      = dll.MustPrep("GetKeyPressed", &ffi.TypeSint32)
	getCharPressed     = dll.MustPrep("GetCharPressed", &ffi.TypeSint32)
	getKeyName         = dll.MustPrep("GetKeyName", &ffi.TypePointer, &ffi.TypeSint32)
	setExitKey         = dll.MustPrep("SetExitKey", &ffi.TypeVoid, &ffi.TypeSint32)

	// Input-related functions: gamepads

	isGamepadAvailable      = dll.MustPrep("IsGamepadAvailable", &ffi.TypeUint8, &ffi.TypeSint32)
	getGamepadName          = dll.MustPrep("GetGamepadName", &ffi.TypePointer, &ffi.TypeSint32)
	isGamepadButtonPressed  = dll.MustPrep("IsGamepadButtonPressed", &ffi.TypeUint8, &ffi.TypeSint32, &ffi.TypeSint32)
	isGamepadButtonDown     = dll.MustPrep("IsGamepadButtonDown", &ffi.TypeUint8, &ffi.TypeSint32, &ffi.TypeSint32)
	isGamepadButtonReleased = dll.MustPrep("IsGamepadButtonReleased", &ffi.TypeUint8, &ffi.TypeSint32, &ffi.TypeSint32)
	isGamepadButtonUp       = dll.MustPrep("IsGamepadButtonUp", &ffi.TypeUint8, &ffi.TypeSint32, &ffi.TypeSint32)
	getGamepadButtonPressed = dll.MustPrep("GetGamepadButtonPressed", &ffi.TypeSint32)
	getGamepadAxisCount     = dll.MustPrep("GetGamepadAxisCount", &ffi.TypeSint32, &ffi.TypeSint32)
	getGamepadAxisMovement  = dll.MustPrep("GetGamepadAxisMovement", &ffi.TypeFloat, &ffi.TypeSint32, &ffi.TypeSint32)
	setGamepadMappings      = dll.MustPrep("SetGamepadMappings", &ffi.TypeSint32, &ffi.TypePointer)
	setGamepadVibration     = dll.MustPrep("SetGamepadVibration", &ffi.TypeVoid, &ffi.TypeSint32, &ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeFloat)

	// Input-related functions: mouse

	isMouseButtonPressed  = dll.MustPrep("IsMouseButtonPressed", &ffi.TypeUint8, &ffi.TypeSint32)
	isMouseButtonDown     = dll.MustPrep("IsMouseButtonDown", &ffi.TypeUint8, &ffi.TypeSint32)
	isMouseButtonReleased = dll.MustPrep("IsMouseButtonReleased", &ffi.TypeUint8, &ffi.TypeSint32)
	isMouseButtonUp       = dll.MustPrep("IsMouseButtonUp", &ffi.TypeUint8, &ffi.TypeSint32)
	getMouseX             = dll.MustPrep("GetMouseX", &ffi.TypeSint32)
	getMouseY             = dll.MustPrep("GetMouseY", &ffi.TypeSint32)
	getMousePosition      = dll.MustPrep("GetMousePosition", &typeVector2)
	getMouseDelta         = dll.MustPrep("GetMouseDelta", &typeVector2)
	setMousePosition      = dll.MustPrep("SetMousePosition", &ffi.TypeVoid, &ffi.TypeSint32, &ffi.TypeSint32)
	setMouseOffset        = dll.MustPrep("SetMouseOffset", &ffi.TypeVoid, &ffi.TypeSint32, &ffi.TypeSint32)
	setMouseScale         = dll.MustPrep("SetMouseScale", &ffi.TypeVoid, &ffi.TypeFloat, &ffi.TypeFloat)
	getMouseWheelMove     = dll.MustPrep("GetMouseWheelMove", &ffi.TypeFloat)
	getMouseWheelMoveV    = dll.MustPrep("GetMouseWheelMoveV", &typeVector2)
	setMouseCursor        = dll.MustPrep("SetMouseCursor", &ffi.TypeVoid, &ffi.TypeSint32)

	// Input-related functions: touch

	getTouchX          = dll.MustPrep("GetTouchX", &ffi.TypeSint32)
	getTouchY          = dll.MustPrep("GetTouchY", &ffi.TypeSint32)
	getTouchPosition   = dll.MustPrep("GetTouchPosition", &typeVector2, &ffi.TypeSint32)
	getTouchPointId    = dll.MustPrep("GetTouchPointId", &ffi.TypeSint32, &ffi.TypeSint32)
	getTouchPointCount = dll.MustPrep("GetTouchPointCount", &ffi.TypeSint32)

	// Gestures and Touch Handling Functions (Module: rgestures)

	setGesturesEnabled     = dll.MustPrep("SetGesturesEnabled", &ffi.TypeVoid, &ffi.TypeUint32)
	isGestureDetected      = dll.MustPrep("IsGestureDetected", &ffi.TypeUint8, &ffi.TypeUint32)
	getGestureDetected     = dll.MustPrep("GetGestureDetected", &ffi.TypeSint32)
	getGestureHoldDuration = dll.MustPrep("GetGestureHoldDuration", &ffi.TypeFloat)
	getGestureDragVector   = dll.MustPrep("GetGestureDragVector", &typeVector2)
	getGestureDragAngle    = dll.MustPrep("GetGestureDragAngle", &ffi.TypeFloat)
	getGesturePinchVector  = dll.MustPrep("GetGesturePinchVector", &typeVector2)
	getGesturePinchAngle   = dll.MustPrep("GetGesturePinchAngle", &ffi.TypeFloat)

	// Basic Shapes Drawing Functions (Module: shapes)

	setShapesTexture          = dll.MustPrep("SetShapesTexture", &ffi.TypeVoid, &typeTexture2D, &typeRectangle)
	getShapesTexture          = dll.MustPrep("GetShapesTexture", &typeTexture2D)
	getShapesTextureRectangle = dll.MustPrep("GetShapesTextureRectangle", &typeRectangle)

	// Basic shapes drawing functions

	drawPixel                   = dll.MustPrep("DrawPixel", &ffi.TypeVoid, &ffi.TypeSint32, &ffi.TypeSint32, &typeColor)
	drawPixelV                  = dll.MustPrep("DrawPixelV", &ffi.TypeVoid, &typeVector2, &typeColor)
	drawLine                    = dll.MustPrep("DrawLine", &ffi.TypeVoid, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32, &typeColor)
	drawLineV                   = dll.MustPrep("DrawLineV", &ffi.TypeVoid, &typeVector2, &typeVector2, &typeColor)
	drawLineEx                  = dll.MustPrep("DrawLineEx", &ffi.TypeVoid, &typeVector2, &typeVector2, &ffi.TypeFloat, &typeColor)
	drawLineStrip               = dll.MustPrep("DrawLineStrip", &ffi.TypeVoid, &ffi.TypePointer, &ffi.TypeSint32, &typeColor)
	drawLineBezier              = dll.MustPrep("DrawLineBezier", &ffi.TypeVoid, &typeVector2, &typeVector2, &ffi.TypeFloat, &typeColor)
	drawLineDashed              = dll.MustPrep("DrawLineDashed", &ffi.TypeVoid, &typeVector2, &typeVector2, &ffi.TypeSint32, &ffi.TypeSint32, &typeColor)
	drawCircle                  = dll.MustPrep("DrawCircle", &ffi.TypeVoid, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeFloat, &typeColor)
	drawCircleV                 = dll.MustPrep("DrawCircleV", &ffi.TypeVoid, &typeVector2, &ffi.TypeFloat, &typeColor)
	drawCircleGradient          = dll.MustPrep("DrawCircleGradient", &ffi.TypeVoid, &typeVector2, &ffi.TypeFloat, &typeColor, &typeColor)
	drawCircleSector            = dll.MustPrep("DrawCircleSector", &ffi.TypeVoid, &typeVector2, &ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeSint32, &typeColor)
	drawCircleSectorLines       = dll.MustPrep("DrawCircleSectorLines", &ffi.TypeVoid, &typeVector2, &ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeSint32, &typeColor)
	drawCircleLines             = dll.MustPrep("DrawCircleLines", &ffi.TypeVoid, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeFloat, &typeColor)
	drawCircleLinesV            = dll.MustPrep("DrawCircleLinesV", &ffi.TypeVoid, &typeVector2, &ffi.TypeFloat, &typeColor)
	drawEllipse                 = dll.MustPrep("DrawEllipse", &ffi.TypeVoid, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeFloat, &ffi.TypeFloat, &typeColor)
	drawEllipseV                = dll.MustPrep("DrawEllipseV", &ffi.TypeVoid, &typeVector2, &ffi.TypeFloat, &ffi.TypeFloat, &typeColor)
	drawEllipseLines            = dll.MustPrep("DrawEllipseLines", &ffi.TypeVoid, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeFloat, &ffi.TypeFloat, &typeColor)
	drawEllipseLinesV           = dll.MustPrep("DrawEllipseLinesV", &ffi.TypeVoid, &typeVector2, &ffi.TypeFloat, &ffi.TypeFloat, &typeColor)
	drawRing                    = dll.MustPrep("DrawRing", &ffi.TypeVoid, &typeVector2, &ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeSint32, &typeColor)
	drawRingLines               = dll.MustPrep("DrawRingLines", &ffi.TypeVoid, &typeVector2, &ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeSint32, &typeColor)
	drawRectangle               = dll.MustPrep("DrawRectangle", &ffi.TypeVoid, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32, &typeColor)
	drawRectangleV              = dll.MustPrep("DrawRectangleV", &ffi.TypeVoid, &typeVector2, &typeVector2, &typeColor)
	drawRectangleRec            = dll.MustPrep("DrawRectangleRec", &ffi.TypeVoid, &typeRectangle, &typeColor)
	drawRectanglePro            = dll.MustPrep("DrawRectanglePro", &ffi.TypeVoid, &typeRectangle, &typeVector2, &ffi.TypeFloat, &typeColor)
	drawRectangleGradientV      = dll.MustPrep("DrawRectangleGradientV", &ffi.TypeVoid, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32, &typeColor, &typeColor)
	drawRectangleGradientH      = dll.MustPrep("DrawRectangleGradientH", &ffi.TypeVoid, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32, &typeColor, &typeColor)
	drawRectangleGradientEx     = dll.MustPrep("DrawRectangleGradientEx", &ffi.TypeVoid, &typeRectangle, &typeColor, &typeColor, &typeColor, &typeColor)
	drawRectangleLines          = dll.MustPrep("DrawRectangleLines", &ffi.TypeVoid, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32, &typeColor)
	drawRectangleLinesEx        = dll.MustPrep("DrawRectangleLinesEx", &ffi.TypeVoid, &typeRectangle, &ffi.TypeFloat, &typeColor)
	drawRectangleRounded        = dll.MustPrep("DrawRectangleRounded", &ffi.TypeVoid, &typeRectangle, &ffi.TypeFloat, &ffi.TypeSint32, &typeColor)
	drawRectangleRoundedLines   = dll.MustPrep("DrawRectangleRoundedLines", &ffi.TypeVoid, &typeRectangle, &ffi.TypeFloat, &ffi.TypeSint32, &typeColor)
	drawRectangleRoundedLinesEx = dll.MustPrep("DrawRectangleRoundedLinesEx", &ffi.TypeVoid, &typeRectangle, &ffi.TypeFloat, &ffi.TypeSint32, &ffi.TypeFloat, &typeColor)
	drawTriangle                = dll.MustPrep("DrawTriangle", &ffi.TypeVoid, &typeVector2, &typeVector2, &typeVector2, &typeColor)
	drawTriangleLines           = dll.MustPrep("DrawTriangleLines", &ffi.TypeVoid, &typeVector2, &typeVector2, &typeVector2, &typeColor)
	drawTriangleFan             = dll.MustPrep("DrawTriangleFan", &ffi.TypeVoid, &ffi.TypePointer, &ffi.TypeSint32, &typeColor)
	drawTriangleStrip           = dll.MustPrep("DrawTriangleStrip", &ffi.TypeVoid, &ffi.TypePointer, &ffi.TypeSint32, &typeColor)
	drawPoly                    = dll.MustPrep("DrawPoly", &ffi.TypeVoid, &typeVector2, &ffi.TypeSint32, &ffi.TypeFloat, &ffi.TypeFloat, &typeColor)
	drawPolyLines               = dll.MustPrep("DrawPolyLines", &ffi.TypeVoid, &typeVector2, &ffi.TypeSint32, &ffi.TypeFloat, &ffi.TypeFloat, &typeColor)
	drawPolyLinesEx             = dll.MustPrep("DrawPolyLinesEx", &ffi.TypeVoid, &typeVector2, &ffi.TypeSint32, &ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeFloat, &typeColor)

	// Splines drawing functions

	drawSplineLinear                 = dll.MustPrep("DrawSplineLinear", &ffi.TypeVoid, &ffi.TypePointer, &ffi.TypeSint32, &ffi.TypeFloat, &typeColor)
	drawSplineBasis                  = dll.MustPrep("DrawSplineBasis", &ffi.TypeVoid, &ffi.TypePointer, &ffi.TypeSint32, &ffi.TypeFloat, &typeColor)
	drawSplineCatmullRom             = dll.MustPrep("DrawSplineCatmullRom", &ffi.TypeVoid, &ffi.TypePointer, &ffi.TypeSint32, &ffi.TypeFloat, &typeColor)
	drawSplineBezierQuadratic        = dll.MustPrep("DrawSplineBezierQuadratic", &ffi.TypeVoid, &ffi.TypePointer, &ffi.TypeSint32, &ffi.TypeFloat, &typeColor)
	drawSplineBezierCubic            = dll.MustPrep("DrawSplineBezierCubic", &ffi.TypeVoid, &ffi.TypePointer, &ffi.TypeSint32, &ffi.TypeFloat, &typeColor)
	drawSplineSegmentLinear          = dll.MustPrep("DrawSplineSegmentLinear", &ffi.TypeVoid, &typeVector2, &typeVector2, &ffi.TypeFloat, &typeColor)
	drawSplineSegmentBasis           = dll.MustPrep("DrawSplineSegmentBasis", &ffi.TypeVoid, &typeVector2, &typeVector2, &typeVector2, &typeVector2, &ffi.TypeFloat, &typeColor)
	drawSplineSegmentCatmullRom      = dll.MustPrep("DrawSplineSegmentCatmullRom", &ffi.TypeVoid, &typeVector2, &typeVector2, &typeVector2, &typeVector2, &ffi.TypeFloat, &typeColor)
	drawSplineSegmentBezierQuadratic = dll.MustPrep("DrawSplineSegmentBezierQuadratic", &ffi.TypeVoid, &typeVector2, &typeVector2, &typeVector2, &ffi.TypeFloat, &typeColor)
	drawSplineSegmentBezierCubic     = dll.MustPrep("DrawSplineSegmentBezierCubic", &ffi.TypeVoid, &typeVector2, &typeVector2, &typeVector2, &typeVector2, &ffi.TypeFloat, &typeColor)

	// Spline segment point evaluation functions, for a given t [0.0f .. 1.0f]

	getSplinePointLinear      = dll.MustPrep("GetSplinePointLinear", &typeVector2, &typeVector2, &typeVector2, &ffi.TypeFloat)
	getSplinePointBasis       = dll.MustPrep("GetSplinePointBasis", &typeVector2, &typeVector2, &typeVector2, &typeVector2, &typeVector2, &ffi.TypeFloat)
	getSplinePointCatmullRom  = dll.MustPrep("GetSplinePointCatmullRom", &typeVector2, &typeVector2, &typeVector2, &typeVector2, &typeVector2, &ffi.TypeFloat)
	getSplinePointBezierQuad  = dll.MustPrep("GetSplinePointBezierQuad", &typeVector2, &typeVector2, &typeVector2, &typeVector2, &ffi.TypeFloat)
	getSplinePointBezierCubic = dll.MustPrep("GetSplinePointBezierCubic", &typeVector2, &typeVector2, &typeVector2, &typeVector2, &typeVector2, &ffi.TypeFloat)

	// Basic shapes collision detection functions

	checkCollisionRecs          = dll.MustPrep("CheckCollisionRecs", &ffi.TypeUint8, &typeRectangle, &typeRectangle)
	checkCollisionCircles       = dll.MustPrep("CheckCollisionCircles", &ffi.TypeUint8, &typeVector2, &ffi.TypeFloat, &typeVector2, &ffi.TypeFloat)
	checkCollisionCircleRec     = dll.MustPrep("CheckCollisionCircleRec", &ffi.TypeUint8, &typeVector2, &ffi.TypeFloat, &typeRectangle)
	checkCollisionCircleLine    = dll.MustPrep("CheckCollisionCircleLine", &ffi.TypeUint8, &typeVector2, &ffi.TypeFloat, &typeVector2, &typeVector2)
	checkCollisionPointRec      = dll.MustPrep("CheckCollisionPointRec", &ffi.TypeUint8, &typeVector2, &typeRectangle)
	checkCollisionPointCircle   = dll.MustPrep("CheckCollisionPointCircle", &ffi.TypeUint8, &typeVector2, &typeVector2, &ffi.TypeFloat)
	checkCollisionPointTriangle = dll.MustPrep("CheckCollisionPointTriangle", &ffi.TypeUint8, &typeVector2, &typeVector2, &typeVector2, &typeVector2)
	checkCollisionPointLine     = dll.MustPrep("CheckCollisionPointLine", &ffi.TypeUint8, &typeVector2, &typeVector2, &typeVector2, &ffi.TypeSint32)
	checkCollisionPointPoly     = dll.MustPrep("CheckCollisionPointPoly", &ffi.TypeUint8, &typeVector2, &ffi.TypePointer, &ffi.TypeSint32)
	checkCollisionLines         = dll.MustPrep("CheckCollisionLines", &ffi.TypeUint8, &typeVector2, &typeVector2, &typeVector2, &typeVector2, &ffi.TypePointer)
	getCollisionRec             = dll.MustPrep("GetCollisionRec", &typeRectangle, &typeRectangle, &typeRectangle)

	// Image loading functions

	loadImage               = dll.MustPrep("LoadImage", &typeImage, &ffi.TypePointer)
	loadImageRaw            = dll.MustPrep("LoadImageRaw", &typeImage, &ffi.TypePointer, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32)
	loadImageAnim           = dll.MustPrep("LoadImageAnim", &typeImage, &ffi.TypePointer, &ffi.TypePointer)
	loadImageAnimFromMemory = dll.MustPrep("LoadImageAnimFromMemory", &typeImage, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypeSint32, &ffi.TypePointer)
	loadImageFromMemory     = dll.MustPrep("LoadImageFromMemory", &typeImage, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypeSint32)
	loadImageFromTexture    = dll.MustPrep("LoadImageFromTexture", &typeImage, &typeTexture2D)
	loadImageFromScreen     = dll.MustPrep("LoadImageFromScreen", &typeImage)
	isImageValid            = dll.MustPrep("IsImageValid", &ffi.TypeUint8, &typeImage)
	unloadImage             = dll.MustPrep("UnloadImage", &ffi.TypeVoid, &typeImage)
	exportImage             = dll.MustPrep("ExportImage", &ffi.TypeUint8, &typeImage, &ffi.TypePointer)
	exportImageToMemory     = dll.MustPrep("ExportImageToMemory", &ffi.TypePointer, &typeImage, &ffi.TypePointer, &ffi.TypePointer)

	// Image generation functions

	genImageColor          = dll.MustPrep("GenImageColor", &typeImage, &ffi.TypeSint32, &ffi.TypeSint32, &typeColor)
	genImageGradientLinear = dll.MustPrep("GenImageGradientLinear", &typeImage, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32, &typeColor, &typeColor)
	genImageGradientRadial = dll.MustPrep("GenImageGradientRadial", &typeImage, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeFloat, &typeColor, &typeColor)
	genImageGradientSquare = dll.MustPrep("GenImageGradientSquare", &typeImage, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeFloat, &typeColor, &typeColor)
	genImageChecked        = dll.MustPrep("GenImageChecked", &typeImage, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32, &typeColor, &typeColor)
	genImageWhiteNoise     = dll.MustPrep("GenImageWhiteNoise", &typeImage, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeFloat)
	genImagePerlinNoise    = dll.MustPrep("GenImagePerlinNoise", &typeImage, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeFloat)
	genImageCellular       = dll.MustPrep("GenImageCellular", &typeImage, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32)
	genImageText           = dll.MustPrep("GenImageText", &typeImage, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypePointer)

	// Image manipulation functions

	imageCopy              = dll.MustPrep("ImageCopy", &typeImage, &typeImage)
	imageFromImage         = dll.MustPrep("ImageFromImage", &typeImage, &typeImage, &typeRectangle)
	imageFromChannel       = dll.MustPrep("ImageFromChannel", &typeImage, &typeImage, &ffi.TypeSint32)
	imageText              = dll.MustPrep("ImageText", &typeImage, &ffi.TypePointer, &ffi.TypeSint32, &typeColor)
	imageTextEx            = dll.MustPrep("ImageTextEx", &typeImage, &typeFont, &ffi.TypePointer, &ffi.TypeFloat, &ffi.TypeFloat, &typeColor)
	imageFormat            = dll.MustPrep("ImageFormat", &ffi.TypeVoid, &ffi.TypePointer, &ffi.TypeSint32)
	imageToPOT             = dll.MustPrep("ImageToPOT", &ffi.TypeVoid, &ffi.TypePointer, &typeColor)
	imageCrop              = dll.MustPrep("ImageCrop", &ffi.TypeVoid, &ffi.TypePointer, &typeRectangle)
	imageAlphaCrop         = dll.MustPrep("ImageAlphaCrop", &ffi.TypeVoid, &ffi.TypePointer, &ffi.TypeFloat)
	imageAlphaClear        = dll.MustPrep("ImageAlphaClear", &ffi.TypeVoid, &ffi.TypePointer, &typeColor, &ffi.TypeFloat)
	imageAlphaMask         = dll.MustPrep("ImageAlphaMask", &ffi.TypeVoid, &ffi.TypePointer, &typeImage)
	imageAlphaPremultiply  = dll.MustPrep("ImageAlphaPremultiply", &ffi.TypeVoid, &ffi.TypePointer)
	imageBlurGaussian      = dll.MustPrep("ImageBlurGaussian", &ffi.TypeVoid, &ffi.TypePointer, &ffi.TypeSint32)
	imageKernelConvolution = dll.MustPrep("ImageKernelConvolution", &ffi.TypeVoid, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypeSint32)
	imageResize            = dll.MustPrep("ImageResize", &ffi.TypeVoid, &ffi.TypePointer, &ffi.TypeSint32, &ffi.TypeSint32)
	imageResizeNN          = dll.MustPrep("ImageResizeNN", &ffi.TypeVoid, &ffi.TypePointer, &ffi.TypeSint32, &ffi.TypeSint32)
	imageResizeCanvas      = dll.MustPrep("ImageResizeCanvas", &ffi.TypeVoid, &ffi.TypePointer, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32, &typeColor)
	imageMipmaps           = dll.MustPrep("ImageMipmaps", &ffi.TypeVoid, &ffi.TypePointer)
	imageDither            = dll.MustPrep("ImageDither", &ffi.TypeVoid, &ffi.TypePointer, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32)
	imageFlipVertical      = dll.MustPrep("ImageFlipVertical", &ffi.TypeVoid, &ffi.TypePointer)
	imageFlipHorizontal    = dll.MustPrep("ImageFlipHorizontal", &ffi.TypeVoid, &ffi.TypePointer)
	imageRotate            = dll.MustPrep("ImageRotate", &ffi.TypeVoid, &ffi.TypePointer, &ffi.TypeSint32)
	imageRotateCW          = dll.MustPrep("ImageRotateCW", &ffi.TypeVoid, &ffi.TypePointer)
	imageRotateCCW         = dll.MustPrep("ImageRotateCCW", &ffi.TypeVoid, &ffi.TypePointer)
	imageColorTint         = dll.MustPrep("ImageColorTint", &ffi.TypeVoid, &ffi.TypePointer, &typeColor)
	imageColorInvert       = dll.MustPrep("ImageColorInvert", &ffi.TypeVoid, &ffi.TypePointer)
	imageColorGrayscale    = dll.MustPrep("ImageColorGrayscale", &ffi.TypeVoid, &ffi.TypePointer)
	imageColorContrast     = dll.MustPrep("ImageColorContrast", &ffi.TypeVoid, &ffi.TypePointer, &ffi.TypeFloat)
	imageColorBrightness   = dll.MustPrep("ImageColorBrightness", &ffi.TypeVoid, &ffi.TypePointer, &ffi.TypeSint32)
	imageColorReplace      = dll.MustPrep("ImageColorReplace", &ffi.TypeVoid, &ffi.TypePointer, &typeColor, &typeColor)
	loadImageColors        = dll.MustPrep("LoadImageColors", &ffi.TypePointer, &typeImage)
	loadImagePalette       = dll.MustPrep("LoadImagePalette", &ffi.TypePointer, &typeImage, &ffi.TypeSint32, &ffi.TypePointer)
	unloadImageColors      = dll.MustPrep("UnloadImageColors", &ffi.TypeVoid, &ffi.TypePointer)
	unloadImagePalette     = dll.MustPrep("UnloadImagePalette", &ffi.TypeVoid, &ffi.TypePointer)
	getImageAlphaBorder    = dll.MustPrep("GetImageAlphaBorder", &typeRectangle, &typeImage, &ffi.TypeFloat)
	getImageColor          = dll.MustPrep("GetImageColor", &typeColor, &typeImage, &ffi.TypeSint32, &ffi.TypeSint32)

	// Image drawing functions

	imageClearBackground    = dll.MustPrep("ImageClearBackground", &ffi.TypeVoid, &ffi.TypePointer, &typeColor)
	imageDrawPixel          = dll.MustPrep("ImageDrawPixel", &ffi.TypeVoid, &ffi.TypePointer, &ffi.TypeSint32, &ffi.TypeSint32, &typeColor)
	imageDrawPixelV         = dll.MustPrep("ImageDrawPixelV", &ffi.TypeVoid, &ffi.TypePointer, &typeVector2, &typeColor)
	imageDrawLine           = dll.MustPrep("ImageDrawLine", &ffi.TypeVoid, &ffi.TypePointer, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32, &typeColor)
	imageDrawLineV          = dll.MustPrep("ImageDrawLineV", &ffi.TypeVoid, &ffi.TypePointer, &typeVector2, &typeVector2, &typeColor)
	imageDrawLineEx         = dll.MustPrep("ImageDrawLineEx", &ffi.TypeVoid, &ffi.TypePointer, &typeVector2, &typeVector2, &ffi.TypeSint32, &typeColor)
	imageDrawCircle         = dll.MustPrep("ImageDrawCircle", &ffi.TypeVoid, &ffi.TypePointer, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32, &typeColor)
	imageDrawCircleV        = dll.MustPrep("ImageDrawCircleV", &ffi.TypeVoid, &ffi.TypePointer, &typeVector2, &ffi.TypeSint32, &typeColor)
	imageDrawCircleLines    = dll.MustPrep("ImageDrawCircleLines", &ffi.TypeVoid, &ffi.TypePointer, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32, &typeColor)
	imageDrawCircleLinesV   = dll.MustPrep("ImageDrawCircleLinesV", &ffi.TypeVoid, &ffi.TypePointer, &typeVector2, &ffi.TypeSint32, &typeColor)
	imageDrawRectangle      = dll.MustPrep("ImageDrawRectangle", &ffi.TypeVoid, &ffi.TypePointer, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32, &typeColor)
	imageDrawRectangleV     = dll.MustPrep("ImageDrawRectangleV", &ffi.TypeVoid, &ffi.TypePointer, &typeVector2, &typeVector2, &typeColor)
	imageDrawRectangleRec   = dll.MustPrep("ImageDrawRectangleRec", &ffi.TypeVoid, &ffi.TypePointer, &typeRectangle, &typeColor)
	imageDrawRectangleLines = dll.MustPrep("ImageDrawRectangleLines", &ffi.TypeVoid, &ffi.TypePointer, &typeRectangle, &ffi.TypeSint32, &typeColor)
	imageDrawTriangle       = dll.MustPrep("ImageDrawTriangle", &ffi.TypeVoid, &ffi.TypePointer, &typeVector2, &typeVector2, &typeVector2, &typeColor)
	imageDrawTriangleEx     = dll.MustPrep("ImageDrawTriangleEx", &ffi.TypeVoid, &ffi.TypePointer, &typeVector2, &typeVector2, &typeVector2, &typeColor, &typeColor, &typeColor)
	imageDrawTriangleLines  = dll.MustPrep("ImageDrawTriangleLines", &ffi.TypeVoid, &ffi.TypePointer, &typeVector2, &typeVector2, &typeVector2, &typeColor)
	imageDrawTriangleFan    = dll.MustPrep("ImageDrawTriangleFan", &ffi.TypeVoid, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypeSint32, &typeColor)
	imageDrawTriangleStrip  = dll.MustPrep("ImageDrawTriangleStrip", &ffi.TypeVoid, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypeSint32, &typeColor)
	imageDraw               = dll.MustPrep("ImageDraw", &ffi.TypeVoid, &ffi.TypePointer, &typeImage, &typeRectangle, &typeRectangle, &typeColor)
	imageDrawText           = dll.MustPrep("ImageDrawText", &ffi.TypeVoid, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32, &typeColor)
	imageDrawTextEx         = dll.MustPrep("ImageDrawTextEx", &ffi.TypeVoid, &ffi.TypePointer, &typeFont, &ffi.TypePointer, &typeVector2, &ffi.TypeFloat, &ffi.TypeFloat, &typeColor)

	// Texture loading functions

	loadTexture          = dll.MustPrep("LoadTexture", &typeTexture2D, &ffi.TypePointer)
	loadTextureFromImage = dll.MustPrep("LoadTextureFromImage", &typeTexture2D, &typeImage)
	loadTextureCubemap   = dll.MustPrep("LoadTextureCubemap", &typeTexture2D, &typeImage, &ffi.TypeSint32)
	loadRenderTexture    = dll.MustPrep("LoadRenderTexture", &typeRenderTexture2D, &ffi.TypeSint32, &ffi.TypeSint32)
	isTextureValid       = dll.MustPrep("IsTextureValid", &ffi.TypeUint8, &typeTexture2D)
	unloadTexture        = dll.MustPrep("UnloadTexture", &ffi.TypeVoid, &typeTexture2D)
	isRenderTextureValid = dll.MustPrep("IsRenderTextureValid", &ffi.TypeUint8, &typeRenderTexture2D)
	unloadRenderTexture  = dll.MustPrep("UnloadRenderTexture", &ffi.TypeVoid, &typeRenderTexture2D)
	updateTexture        = dll.MustPrep("UpdateTexture", &ffi.TypeVoid, &typeTexture2D, &ffi.TypePointer)
	updateTextureRec     = dll.MustPrep("UpdateTextureRec", &ffi.TypeVoid, &typeTexture2D, &typeRectangle, &ffi.TypePointer)

	// Texture configuration functions

	genTextureMipmaps = dll.MustPrep("GenTextureMipmaps", &ffi.TypeVoid, &ffi.TypePointer)
	setTextureFilter  = dll.MustPrep("SetTextureFilter", &ffi.TypeVoid, &typeTexture2D, &ffi.TypeSint32)
	setTextureWrap    = dll.MustPrep("SetTextureWrap", &ffi.TypeVoid, &typeTexture2D, &ffi.TypeSint32)

	// Texture drawing functions

	drawTexture       = dll.MustPrep("DrawTexture", &ffi.TypeVoid, &typeTexture2D, &ffi.TypeSint32, &ffi.TypeSint32, &typeColor)
	drawTextureV      = dll.MustPrep("DrawTextureV", &ffi.TypeVoid, &typeTexture2D, &typeVector2, &typeColor)
	drawTextureEx     = dll.MustPrep("DrawTextureEx", &ffi.TypeVoid, &typeTexture2D, &typeVector2, &ffi.TypeFloat, &ffi.TypeFloat, &typeColor)
	drawTextureRec    = dll.MustPrep("DrawTextureRec", &ffi.TypeVoid, &typeTexture2D, &typeRectangle, &typeVector2, &typeColor)
	drawTexturePro    = dll.MustPrep("DrawTexturePro", &ffi.TypeVoid, &typeTexture2D, &typeRectangle, &typeRectangle, &typeVector2, &ffi.TypeFloat, &typeColor)
	drawTextureNPatch = dll.MustPrep("DrawTextureNPatch", &ffi.TypeVoid, &typeTexture2D, &typeNPatchInfo, &typeRectangle, &typeVector2, &ffi.TypeFloat, &typeColor)

	// Color/pixel related functions

	fade                = dll.MustPrep("Fade", &typeColor, &typeColor, &ffi.TypeFloat)
	colorToInt          = dll.MustPrep("ColorToInt", &ffi.TypeSint32, &typeColor)
	colorNormalize      = dll.MustPrep("ColorNormalize", &typeVector4, &typeColor)
	colorFromNormalized = dll.MustPrep("ColorFromNormalized", &typeColor, &typeVector4)
	colorToHSV          = dll.MustPrep("ColorToHSV", &typeVector3, &typeColor)
	colorFromHSV        = dll.MustPrep("ColorFromHSV", &typeColor, &ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeFloat)
	colorTint           = dll.MustPrep("ColorTint", &typeColor, &typeColor, &typeColor)
	colorBrightness     = dll.MustPrep("ColorBrightness", &typeColor, &typeColor, &ffi.TypeFloat)
	colorContrast       = dll.MustPrep("ColorContrast", &typeColor, &typeColor, &ffi.TypeFloat)
	colorAlpha          = dll.MustPrep("ColorAlpha", &typeColor, &typeColor, &ffi.TypeFloat)
	colorAlphaBlend     = dll.MustPrep("ColorAlphaBlend", &typeColor, &typeColor, &typeColor, &typeColor)
	colorLerp           = dll.MustPrep("ColorLerp", &typeColor, &typeColor, &typeColor, &ffi.TypeFloat)
	getColor            = dll.MustPrep("GetColor", &typeColor, &ffi.TypeUint32)
	getPixelColor       = dll.MustPrep("GetPixelColor", &typeColor, &ffi.TypePointer, &ffi.TypeSint32)
	setPixelColor       = dll.MustPrep("SetPixelColor", &ffi.TypeVoid, &ffi.TypePointer, &typeColor, &ffi.TypeSint32)
	getPixelDataSize    = dll.MustPrep("GetPixelDataSize", &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32)

	// Font loading/unloading functions

	getFontDefault     = dll.MustPrep("GetFontDefault", &typeFont)
	loadFont           = dll.MustPrep("LoadFont", &typeFont, &ffi.TypePointer)
	loadFontEx         = dll.MustPrep("LoadFontEx", &typeFont, &ffi.TypePointer, &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypeSint32)
	loadFontFromImage  = dll.MustPrep("LoadFontFromImage", &typeFont, &typeImage, &typeColor, &ffi.TypeSint32)
	loadFontFromMemory = dll.MustPrep("LoadFontFromMemory", &typeFont, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypeSint32)
	isFontValid        = dll.MustPrep("IsFontValid", &ffi.TypeUint8, &typeFont)
	loadFontData       = dll.MustPrep("LoadFontData", &ffi.TypePointer, &ffi.TypePointer, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypePointer)
	genImageFontAtlas  = dll.MustPrep("GenImageFontAtlas", &typeImage, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32)
	unloadFontData     = dll.MustPrep("UnloadFontData", &ffi.TypeVoid, &ffi.TypePointer, &ffi.TypeSint32)
	unloadFont         = dll.MustPrep("UnloadFont", &ffi.TypeVoid, &typeFont)

	// Text drawing functions

	drawFPS            = dll.MustPrep("DrawFPS", &ffi.TypeVoid, &ffi.TypeSint32, &ffi.TypeSint32)
	drawText           = dll.MustPrep("DrawText", &ffi.TypeVoid, &ffi.TypePointer, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32, &typeColor)
	drawTextEx         = dll.MustPrep("DrawTextEx", &ffi.TypeVoid, &typeFont, &ffi.TypePointer, &typeVector2, &ffi.TypeFloat, &ffi.TypeFloat, &typeColor)
	drawTextPro        = dll.MustPrep("DrawTextPro", &ffi.TypeVoid, &typeFont, &ffi.TypePointer, &typeVector2, &typeVector2, &ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeFloat, &typeColor)
	drawTextCodepoint  = dll.MustPrep("DrawTextCodepoint", &ffi.TypeVoid, &typeFont, &ffi.TypeSint32, &typeVector2, &ffi.TypeFloat, &typeColor)
	drawTextCodepoints = dll.MustPrep("DrawTextCodepoints", &ffi.TypeVoid, &typeFont, &ffi.TypePointer, &ffi.TypeSint32, &typeVector2, &ffi.TypeFloat, &ffi.TypeFloat, &typeColor)

	// Text font info functions

	setTextLineSpacing    = dll.MustPrep("SetTextLineSpacing", &ffi.TypeVoid, &ffi.TypeSint32)
	measureText           = dll.MustPrep("MeasureText", &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypeSint32)
	measureTextEx         = dll.MustPrep("MeasureTextEx", &typeVector2, &typeFont, &ffi.TypePointer, &ffi.TypeFloat, &ffi.TypeFloat)
	measureTextCodepoints = dll.MustPrep("MeasureTextCodepoints", &typeVector2, &typeFont, &ffi.TypePointer, &ffi.TypeSint32, &ffi.TypeFloat, &ffi.TypeFloat)
	getGlyphIndex         = dll.MustPrep("GetGlyphIndex", &ffi.TypeSint32, &typeFont, &ffi.TypeSint32)
	getGlyphInfo          = dll.MustPrep("GetGlyphInfo", &typeGlyphInfo, &typeFont, &ffi.TypeSint32)
	getGlyphAtlasRec      = dll.MustPrep("GetGlyphAtlasRec", &typeRectangle, &typeFont, &ffi.TypeSint32)

	// Basic geometric 3D shapes drawing functions

	drawLine3D          = dll.MustPrep("DrawLine3D", &ffi.TypeVoid, &typeVector3, &typeVector3, &typeColor)
	drawPoint3D         = dll.MustPrep("DrawPoint3D", &ffi.TypeVoid, &typeVector3, &typeColor)
	drawCircle3D        = dll.MustPrep("DrawCircle3D", &ffi.TypeVoid, &typeVector3, &ffi.TypeFloat, &typeVector3, &ffi.TypeFloat, &typeColor)
	drawTriangle3D      = dll.MustPrep("DrawTriangle3D", &ffi.TypeVoid, &typeVector3, &typeVector3, &typeVector3, &typeColor)
	drawTriangleStrip3D = dll.MustPrep("DrawTriangleStrip3D", &ffi.TypeVoid, &ffi.TypePointer, &ffi.TypeSint32, &typeColor)
	drawCube            = dll.MustPrep("DrawCube", &ffi.TypeVoid, &typeVector3, &ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeFloat, &typeColor)
	drawCubeV           = dll.MustPrep("DrawCubeV", &ffi.TypeVoid, &typeVector3, &typeVector3, &typeColor)
	drawCubeWires       = dll.MustPrep("DrawCubeWires", &ffi.TypeVoid, &typeVector3, &ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeFloat, &typeColor)
	drawCubeWiresV      = dll.MustPrep("DrawCubeWiresV", &ffi.TypeVoid, &typeVector3, &typeVector3, &typeColor)
	drawSphere          = dll.MustPrep("DrawSphere", &ffi.TypeVoid, &typeVector3, &ffi.TypeFloat, &typeColor)
	drawSphereEx        = dll.MustPrep("DrawSphereEx", &ffi.TypeVoid, &typeVector3, &ffi.TypeFloat, &ffi.TypeSint32, &ffi.TypeSint32, &typeColor)
	drawSphereWires     = dll.MustPrep("DrawSphereWires", &ffi.TypeVoid, &typeVector3, &ffi.TypeFloat, &ffi.TypeSint32, &ffi.TypeSint32, &typeColor)
	drawCylinder        = dll.MustPrep("DrawCylinder", &ffi.TypeVoid, &typeVector3, &ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeSint32, &typeColor)
	drawCylinderEx      = dll.MustPrep("DrawCylinderEx", &ffi.TypeVoid, &typeVector3, &typeVector3, &ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeSint32, &typeColor)
	drawCylinderWires   = dll.MustPrep("DrawCylinderWires", &ffi.TypeVoid, &typeVector3, &ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeSint32, &typeColor)
	drawCylinderWiresEx = dll.MustPrep("DrawCylinderWiresEx", &ffi.TypeVoid, &typeVector3, &typeVector3, &ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeSint32, &typeColor)
	drawCapsule         = dll.MustPrep("DrawCapsule", &ffi.TypeVoid, &typeVector3, &typeVector3, &ffi.TypeFloat, &ffi.TypeSint32, &ffi.TypeSint32, &typeColor)
	drawCapsuleWires    = dll.MustPrep("DrawCapsuleWires", &ffi.TypeVoid, &typeVector3, &typeVector3, &ffi.TypeFloat, &ffi.TypeSint32, &ffi.TypeSint32, &typeColor)
	drawPlane           = dll.MustPrep("DrawPlane", &ffi.TypeVoid, &typeVector3, &typeVector2, &typeColor)
	drawRay             = dll.MustPrep("DrawRay", &ffi.TypeVoid, &typeRay, &typeColor)
	drawGrid            = dll.MustPrep("DrawGrid", &ffi.TypeVoid, &ffi.TypeSint32, &ffi.TypeFloat)

	// Model management functions

	loadModel           = dll.MustPrep("LoadModel", &typeModel, &ffi.TypePointer)
	loadModelFromMesh   = dll.MustPrep("LoadModelFromMesh", &typeModel, &typeMesh)
	isModelValid        = dll.MustPrep("IsModelValid", &ffi.TypeUint8, &typeModel)
	unloadModel         = dll.MustPrep("UnloadModel", &ffi.TypeVoid, &typeModel)
	getModelBoundingBox = dll.MustPrep("GetModelBoundingBox", &typeBoundingBox, &typeModel)

	// Model drawing functions

	drawModel        = dll.MustPrep("DrawModel", &ffi.TypeVoid, &typeModel, &typeVector3, &ffi.TypeFloat, &typeColor)
	drawModelEx      = dll.MustPrep("DrawModelEx", &ffi.TypeVoid, &typeModel, &typeVector3, &typeVector3, &ffi.TypeFloat, &typeVector3, &typeColor)
	drawModelWires   = dll.MustPrep("DrawModelWires", &ffi.TypeVoid, &typeModel, &typeVector3, &ffi.TypeFloat, &typeColor)
	drawModelWiresEx = dll.MustPrep("DrawModelWiresEx", &ffi.TypeVoid, &typeModel, &typeVector3, &typeVector3, &ffi.TypeFloat, &typeVector3, &typeColor)
	drawBoundingBox  = dll.MustPrep("DrawBoundingBox", &ffi.TypeVoid, &typeBoundingBox, &typeColor)
	drawBillboard    = dll.MustPrep("DrawBillboard", &ffi.TypeVoid, &typeCamera3D, &typeTexture2D, &typeVector3, &ffi.TypeFloat, &typeColor)
	drawBillboardRec = dll.MustPrep("DrawBillboardRec", &ffi.TypeVoid, &typeCamera3D, &typeTexture2D, &typeRectangle, &typeVector3, &typeVector2, &typeColor)
	drawBillboardPro = dll.MustPrep("DrawBillboardPro", &ffi.TypeVoid, &typeCamera3D, &typeTexture2D, &typeRectangle, &typeVector3, &typeVector3, &typeVector2, &typeVector2, &ffi.TypeFloat, &typeColor)

	// Mesh management functions

	uploadMesh         = dll.MustPrep("UploadMesh", &ffi.TypeVoid, &ffi.TypePointer, &ffi.TypeUint8)
	updateMeshBuffer   = dll.MustPrep("UpdateMeshBuffer", &ffi.TypeVoid, &typeMesh, &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypeSint32, &ffi.TypeSint32)
	unloadMesh         = dll.MustPrep("UnloadMesh", &ffi.TypeVoid, &typeMesh)
	drawMesh           = dll.MustPrep("DrawMesh", &ffi.TypeVoid, &typeMesh, &typeMaterial, &typeMatrix)
	drawMeshInstanced  = dll.MustPrep("DrawMeshInstanced", &ffi.TypeVoid, &typeMesh, &typeMaterial, &ffi.TypePointer, &ffi.TypeSint32)
	getMeshBoundingBox = dll.MustPrep("GetMeshBoundingBox", &typeBoundingBox, &typeMesh)
	genMeshTangents    = dll.MustPrep("GenMeshTangents", &ffi.TypeVoid, &ffi.TypePointer)
	exportMesh         = dll.MustPrep("ExportMesh", &ffi.TypeUint8, &typeMesh, &ffi.TypePointer)

	// Mesh generation functions

	genMeshPoly       = dll.MustPrep("GenMeshPoly", &typeMesh, &ffi.TypeSint32, &ffi.TypeFloat)
	genMeshPlane      = dll.MustPrep("GenMeshPlane", &typeMesh, &ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeSint32, &ffi.TypeSint32)
	genMeshCube       = dll.MustPrep("GenMeshCube", &typeMesh, &ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeFloat)
	genMeshSphere     = dll.MustPrep("GenMeshSphere", &typeMesh, &ffi.TypeFloat, &ffi.TypeSint32, &ffi.TypeSint32)
	genMeshHemiSphere = dll.MustPrep("GenMeshHemiSphere", &typeMesh, &ffi.TypeFloat, &ffi.TypeSint32, &ffi.TypeSint32)
	genMeshCylinder   = dll.MustPrep("GenMeshCylinder", &typeMesh, &ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeSint32)
	genMeshCone       = dll.MustPrep("GenMeshCone", &typeMesh, &ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeSint32)
	genMeshTorus      = dll.MustPrep("GenMeshTorus", &typeMesh, &ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeSint32, &ffi.TypeSint32)
	genMeshKnot       = dll.MustPrep("GenMeshKnot", &typeMesh, &ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeSint32, &ffi.TypeSint32)
	genMeshHeightmap  = dll.MustPrep("GenMeshHeightmap", &typeMesh, &typeImage, &typeVector3)
	genMeshCubicmap   = dll.MustPrep("GenMeshCubicmap", &typeMesh, &typeImage, &typeVector3)

	// Material loading/unloading functions

	loadMaterials        = dll.MustPrep("LoadMaterials", &ffi.TypePointer, &ffi.TypePointer, &ffi.TypePointer)
	loadMaterialDefault  = dll.MustPrep("LoadMaterialDefault", &typeMaterial)
	isMaterialValid      = dll.MustPrep("IsMaterialValid", &ffi.TypeUint8, &typeMaterial)
	unloadMaterial       = dll.MustPrep("UnloadMaterial", &ffi.TypeVoid, &typeMaterial)
	setMaterialTexture   = dll.MustPrep("SetMaterialTexture", &ffi.TypeVoid, &ffi.TypePointer, &ffi.TypeSint32, &typeTexture2D)
	setModelMeshMaterial = dll.MustPrep("SetModelMeshMaterial", &ffi.TypeVoid, &ffi.TypePointer, &ffi.TypeSint32, &ffi.TypeSint32)

	// Model animations loading/unloading functions

	loadModelAnimations    = dll.MustPrep("LoadModelAnimations", &ffi.TypePointer, &ffi.TypePointer, &ffi.TypePointer)
	updateModelAnimation   = dll.MustPrep("UpdateModelAnimation", &ffi.TypeVoid, &typeModel, &typeModelAnimation, &ffi.TypeFloat)
	updateModelAnimationEx = dll.MustPrep("UpdateModelAnimationEx", &ffi.TypeVoid, &typeModel, &typeModelAnimation, &ffi.TypeFloat, &typeModelAnimation, &ffi.TypeFloat, &ffi.TypeFloat)
	unloadModelAnimations  = dll.MustPrep("UnloadModelAnimations", &ffi.TypeVoid, &ffi.TypePointer, &ffi.TypeSint32)
	isModelAnimationValid  = dll.MustPrep("IsModelAnimationValid", &ffi.TypeUint8, &typeModel, &typeModelAnimation)

	// Collision detection functions

	checkCollisionSpheres   = dll.MustPrep("CheckCollisionSpheres", &ffi.TypeUint8, &typeVector3, &ffi.TypeFloat, &typeVector3, &ffi.TypeFloat)
	checkCollisionBoxes     = dll.MustPrep("CheckCollisionBoxes", &ffi.TypeUint8, &typeBoundingBox, &typeBoundingBox)
	checkCollisionBoxSphere = dll.MustPrep("CheckCollisionBoxSphere", &ffi.TypeUint8, &typeBoundingBox, &typeVector3, &ffi.TypeFloat)
	getRayCollisionSphere   = dll.MustPrep("GetRayCollisionSphere", &typeRayCollision, &typeRay, &typeVector3, &ffi.TypeFloat)
	getRayCollisionBox      = dll.MustPrep("GetRayCollisionBox", &typeRayCollision, &typeRay, &typeBoundingBox)
	getRayCollisionMesh     = dll.MustPrep("GetRayCollisionMesh", &typeRayCollision, &typeRay, &typeMesh, &typeMatrix)
	getRayCollisionTriangle = dll.MustPrep("GetRayCollisionTriangle", &typeRayCollision, &typeRay, &typeVector3, &typeVector3, &typeVector3)
	getRayCollisionQuad     = dll.MustPrep("GetRayCollisionQuad", &typeRayCollision, &typeRay, &typeVector3, &typeVector3, &typeVector3, &typeVector3)

	// Audio device management functions

	initAudioDevice    = dll.MustPrep("InitAudioDevice", &ffi.TypeVoid)
	closeAudioDevice   = dll.MustPrep("CloseAudioDevice", &ffi.TypeVoid)
	isAudioDeviceReady = dll.MustPrep("IsAudioDeviceReady", &ffi.TypeUint8)
	setMasterVolume    = dll.MustPrep("SetMasterVolume", &ffi.TypeVoid, &ffi.TypeFloat)
	getMasterVolume    = dll.MustPrep("GetMasterVolume", &ffi.TypeFloat)

	// Wave/Sound loading/unloading functions

	loadWave           = dll.MustPrep("LoadWave", &typeWave, &ffi.TypePointer)
	loadWaveFromMemory = dll.MustPrep("LoadWaveFromMemory", &typeWave, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypeSint32)
	isWaveValid        = dll.MustPrep("IsWaveValid", &ffi.TypeUint8, &typeWave)
	loadSound          = dll.MustPrep("LoadSound", &typeSound, &ffi.TypePointer)
	loadSoundFromWave  = dll.MustPrep("LoadSoundFromWave", &typeSound, &typeWave)
	loadSoundAlias     = dll.MustPrep("LoadSoundAlias", &typeSound, &typeSound)
	isSoundValid       = dll.MustPrep("IsSoundValid", &ffi.TypeUint8, &typeSound)
	updateSound        = dll.MustPrep("UpdateSound", &ffi.TypeVoid, &typeSound, &ffi.TypePointer, &ffi.TypeSint32)
	unloadWave         = dll.MustPrep("UnloadWave", &ffi.TypeVoid, &typeWave)
	unloadSound        = dll.MustPrep("UnloadSound", &ffi.TypeVoid, &typeSound)
	unloadSoundAlias   = dll.MustPrep("UnloadSoundAlias", &ffi.TypeVoid, &typeSound)
	exportWave         = dll.MustPrep("ExportWave", &ffi.TypeUint8, &typeWave, &ffi.TypePointer)

	// Wave/Sound management functions

	playSound         = dll.MustPrep("PlaySound", &ffi.TypeVoid, &typeSound)
	stopSound         = dll.MustPrep("StopSound", &ffi.TypeVoid, &typeSound)
	pauseSound        = dll.MustPrep("PauseSound", &ffi.TypeVoid, &typeSound)
	resumeSound       = dll.MustPrep("ResumeSound", &ffi.TypeVoid, &typeSound)
	isSoundPlaying    = dll.MustPrep("IsSoundPlaying", &ffi.TypeUint8, &typeSound)
	setSoundVolume    = dll.MustPrep("SetSoundVolume", &ffi.TypeVoid, &typeSound, &ffi.TypeFloat)
	setSoundPitch     = dll.MustPrep("SetSoundPitch", &ffi.TypeVoid, &typeSound, &ffi.TypeFloat)
	setSoundPan       = dll.MustPrep("SetSoundPan", &ffi.TypeVoid, &typeSound, &ffi.TypeFloat)
	waveCopy          = dll.MustPrep("WaveCopy", &typeWave, &typeWave)
	waveCrop          = dll.MustPrep("WaveCrop", &ffi.TypeVoid, &ffi.TypePointer, &ffi.TypeSint32, &ffi.TypeSint32)
	waveFormat        = dll.MustPrep("WaveFormat", &ffi.TypeVoid, &ffi.TypePointer, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32)
	loadWaveSamples   = dll.MustPrep("LoadWaveSamples", &ffi.TypePointer, &typeWave)
	unloadWaveSamples = dll.MustPrep("UnloadWaveSamples", &ffi.TypeVoid, &ffi.TypePointer)

	// Music management functions

	loadMusicStream           = dll.MustPrep("LoadMusicStream", &typeMusic, &ffi.TypePointer)
	loadMusicStreamFromMemory = dll.MustPrep("LoadMusicStreamFromMemory", &typeMusic, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypeSint32)
	isMusicValid              = dll.MustPrep("IsMusicValid", &ffi.TypeUint8, &typeMusic)
	unloadMusicStream         = dll.MustPrep("UnloadMusicStream", &ffi.TypeVoid, &typeMusic)
	playMusicStream           = dll.MustPrep("PlayMusicStream", &ffi.TypeVoid, &typeMusic)
	isMusicStreamPlaying      = dll.MustPrep("IsMusicStreamPlaying", &ffi.TypeUint8, &typeMusic)
	updateMusicStream         = dll.MustPrep("UpdateMusicStream", &ffi.TypeVoid, &typeMusic)
	stopMusicStream           = dll.MustPrep("StopMusicStream", &ffi.TypeVoid, &typeMusic)
	pauseMusicStream          = dll.MustPrep("PauseMusicStream", &ffi.TypeVoid, &typeMusic)
	resumeMusicStream         = dll.MustPrep("ResumeMusicStream", &ffi.TypeVoid, &typeMusic)
	seekMusicStream           = dll.MustPrep("SeekMusicStream", &ffi.TypeVoid, &typeMusic, &ffi.TypeFloat)
	setMusicVolume            = dll.MustPrep("SetMusicVolume", &ffi.TypeVoid, &typeMusic, &ffi.TypeFloat)
	setMusicPitch             = dll.MustPrep("SetMusicPitch", &ffi.TypeVoid, &typeMusic, &ffi.TypeFloat)
	setMusicPan               = dll.MustPrep("SetMusicPan", &ffi.TypeVoid, &typeMusic, &ffi.TypeFloat)
	getMusicTimeLength        = dll.MustPrep("GetMusicTimeLength", &ffi.TypeFloat, &typeMusic)
	getMusicTimePlayed        = dll.MustPrep("GetMusicTimePlayed", &ffi.TypeFloat, &typeMusic)

	// AudioStream management functions

	loadAudioStream                 = dll.MustPrep("LoadAudioStream", &typeAudioStream, &ffi.TypeUint32, &ffi.TypeUint32, &ffi.TypeUint32)
	isAudioStreamValid              = dll.MustPrep("IsAudioStreamValid", &ffi.TypeUint8, &typeAudioStream)
	unloadAudioStream               = dll.MustPrep("UnloadAudioStream", &ffi.TypeVoid, &typeAudioStream)
	updateAudioStream               = dll.MustPrep("UpdateAudioStream", &ffi.TypeVoid, &typeAudioStream, &ffi.TypePointer, &ffi.TypeSint32)
	isAudioStreamProcessed          = dll.MustPrep("IsAudioStreamProcessed", &ffi.TypeUint8, &typeAudioStream)
	playAudioStream                 = dll.MustPrep("PlayAudioStream", &ffi.TypeVoid, &typeAudioStream)
	pauseAudioStream                = dll.MustPrep("PauseAudioStream", &ffi.TypeVoid, &typeAudioStream)
	resumeAudioStream               = dll.MustPrep("ResumeAudioStream", &ffi.TypeVoid, &typeAudioStream)
	isAudioStreamPlaying            = dll.MustPrep("IsAudioStreamPlaying", &ffi.TypeUint8, &typeAudioStream)
	stopAudioStream                 = dll.MustPrep("StopAudioStream", &ffi.TypeVoid, &typeAudioStream)
	setAudioStreamVolume            = dll.MustPrep("SetAudioStreamVolume", &ffi.TypeVoid, &typeAudioStream, &ffi.TypeFloat)
	setAudioStreamPitch             = dll.MustPrep("SetAudioStreamPitch", &ffi.TypeVoid, &typeAudioStream, &ffi.TypeFloat)
	setAudioStreamPan               = dll.MustPrep("SetAudioStreamPan", &ffi.TypeVoid, &typeAudioStream, &ffi.TypeFloat)
	setAudioStreamBufferSizeDefault = dll.MustPrep("SetAudioStreamBufferSizeDefault", &ffi.TypeVoid, &ffi.TypeSint32)
	setAudioStreamCallback          = dll.MustPrep("SetAudioStreamCallback", &ffi.TypeVoid, &typeAudioStream, &ffi.TypePointer)
	attachAudioStreamProcessor      = dll.MustPrep("AttachAudioStreamProcessor", &ffi.TypeVoid, &typeAudioStream, &ffi.TypePointer)
	detachAudioStreamProcessor      = dll.MustPrep("DetachAudioStreamProcessor", &ffi.TypeVoid, &typeAudioStream, &ffi.TypePointer)
	attachAudioMixedProcessor       = dll.MustPrep("AttachAudioMixedProcessor", &ffi.TypeVoid, &ffi.TypePointer)
	detachAudioMixedProcessor       = dll.MustPrep("DetachAudioMixedProcessor", &ffi.TypeVoid, &ffi.TypePointer)
)

// InitWindow - Initialize window and OpenGL context
func InitWindow(width int32, height int32, title string) {
	titlePtr := convert.ToBytePtr(title)
	initWindow.Call(nil, &width, &height, &titlePtr)
}

// CloseWindow - Close window and unload OpenGL context
func CloseWindow() {
	closeWindow.Call(nil)
}

// WindowShouldClose - Check if application should close ([KeyEscape] pressed or windows close icon clicked)
func WindowShouldClose() bool {
	var ret ffi.Arg
	windowShouldClose.Call(&ret)
	return ret.Bool()
}

// IsWindowReady - Check if window has been initialized successfully
func IsWindowReady() bool {
	var ret ffi.Arg
	isWindowReady.Call(&ret)
	return ret.Bool()
}

// IsWindowFullscreen - Check if window is currently fullscreen
func IsWindowFullscreen() bool {
	var ret ffi.Arg
	isWindowFullscreen.Call(&ret)
	return ret.Bool()
}

// IsWindowHidden - Check if window is currently hidden
func IsWindowHidden() bool {
	var ret ffi.Arg
	isWindowHidden.Call(&ret)
	return ret.Bool()
}

// IsWindowMinimized - Check if window is currently minimized
func IsWindowMinimized() bool {
	var ret ffi.Arg
	isWindowMinimized.Call(&ret)
	return ret.Bool()
}

// IsWindowMaximized - Check if window is currently maximized
func IsWindowMaximized() bool {
	var ret ffi.Arg
	isWindowMaximized.Call(&ret)
	return ret.Bool()
}

// IsWindowFocused - Check if window is currently focused
func IsWindowFocused() bool {
	var ret ffi.Arg
	isWindowFocused.Call(&ret)
	return ret.Bool()
}

// IsWindowResized - Check if window has been resized last frame
func IsWindowResized() bool {
	var ret ffi.Arg
	isWindowResized.Call(&ret)
	return ret.Bool()
}

// IsWindowState - Check if one specific window flag is enabled
func IsWindowState(flag uint32) bool {
	var ret ffi.Arg
	isWindowState.Call(&ret, &flag)
	return ret.Bool()
}

// SetWindowState - Set window configuration state using flags
func SetWindowState(flags uint32) {
	setWindowState.Call(nil, &flags)
}

// ClearWindowState - Clear window configuration state flags
func ClearWindowState(flags uint32) {
	clearWindowState.Call(nil, &flags)
}

// ToggleFullscreen - Toggle window state: fullscreen/windowed, resizes monitor to match window resolution
func ToggleFullscreen() {
	toggleFullscreen.Call(nil)
}

// ToggleBorderlessWindowed - Toggle window state: borderless windowed, resizes window to match monitor resolution
func ToggleBorderlessWindowed() {
	toggleBorderlessWindowed.Call(nil)
}

// MaximizeWindow - Set window state: maximized, if resizable
func MaximizeWindow() {
	maximizeWindow.Call(nil)
}

// MinimizeWindow - Set window state: minimized, if resizable
func MinimizeWindow() {
	minimizeWindow.Call(nil)
}

// RestoreWindow - Restore window from being minimized/maximized
func RestoreWindow() {
	restoreWindow.Call(nil)
}

// SetWindowIcon - Set icon for window (single image, RGBA 32bit, only PLATFORM_DESKTOP)
func SetWindowIcon(image Image) {
	setWindowIcon.Call(nil, &image)
}

// SetWindowIcons - Set icon for window (multiple images, RGBA 32bit, only PLATFORM_DESKTOP)
func SetWindowIcons(images []Image, count int32) {
	imagesPtr := &images[0]
	setWindowIcons.Call(nil, &imagesPtr, &count)
}

// SetWindowTitle - Set title for window
func SetWindowTitle(title string) {
	titlePtr := convert.ToBytePtr(title)
	setWindowTitle.Call(nil, &titlePtr)
}

// SetWindowPosition - Set window position on screen
func SetWindowPosition(x, y int) {
	posX, posY := int32(x), int32(y)
	setWindowPosition.Call(nil, &posX, &posY)
}

// SetWindowMonitor - Set monitor for the current window
func SetWindowMonitor(monitor int) {
	m := int32(monitor)
	setWindowMonitor.Call(nil, &m)
}

// SetWindowMinSize - Set window minimum dimensions (for FLAG_WINDOW_RESIZABLE)
func SetWindowMinSize(width int, height int) {
	w, h := int32(width), int32(height)
	setWindowMinSize.Call(nil, &w, &h)
}

// SetWindowMaxSize - Set window maximum dimensions (for FLAG_WINDOW_RESIZABLE)
func SetWindowMaxSize(width int, height int) {
	w, h := int32(width), int32(height)
	setWindowMaxSize.Call(nil, &w, &h)
}

// SetWindowSize - Set window dimensions
func SetWindowSize(width int, height int) {
	w, h := int32(width), int32(height)
	setWindowSize.Call(nil, &w, &h)
}

// SetWindowOpacity - Set window opacity [0.0f..1.0f]
func SetWindowOpacity(opacity float32) {
	setWindowOpacity.Call(nil, &opacity)
}

// SetWindowFocused - Set window focused
func SetWindowFocused() {
	setWindowFocused.Call(nil)
}

// GetWindowHandle - Get native window handle
func GetWindowHandle() unsafe.Pointer {
	var ret unsafe.Pointer
	getWindowHandle.Call(&ret)
	return ret
}

// GetScreenWidth - Get current screen width
func GetScreenWidth() int {
	var ret ffi.Arg
	getScreenWidth.Call(&ret)
	return int(ret)
}

// GetScreenHeight - Get current screen height
func GetScreenHeight() int {
	var ret ffi.Arg
	getScreenHeight.Call(&ret)
	return int(ret)
}

// GetRenderWidth - Get current render width (it considers HiDPI)
func GetRenderWidth() int {
	var ret ffi.Arg
	getRenderWidth.Call(&ret)
	return int(ret)
}

// GetRenderHeight - Get current render height (it considers HiDPI)
func GetRenderHeight() int {
	var ret ffi.Arg
	getRenderHeight.Call(&ret)
	return int(ret)
}

// GetMonitorCount - Get number of connected monitors
func GetMonitorCount() int {
	var ret ffi.Arg
	getMonitorCount.Call(&ret)
	return int(ret)
}

// GetCurrentMonitor - Get current monitor where window is placed
func GetCurrentMonitor() int {
	var ret ffi.Arg
	getCurrentMonitor.Call(&ret)
	return int(ret)
}

// GetMonitorPosition - Get specified monitor position
func GetMonitorPosition(monitor int) Vector2 {
	var ret Vector2
	m := int32(monitor)
	getMonitorPosition.Call(&ret, &m)
	return ret
}

// GetMonitorWidth - Get specified monitor width (current video mode used by monitor)
func GetMonitorWidth(monitor int) int {
	var ret ffi.Arg
	m := int32(monitor)
	getMonitorWidth.Call(&ret, &m)
	return int(ret)
}

// GetMonitorHeight - Get specified monitor height (current video mode used by monitor)
func GetMonitorHeight(monitor int) int {
	var ret ffi.Arg
	m := int32(monitor)
	getMonitorHeight.Call(&ret, &m)
	return int(ret)
}

// GetMonitorPhysicalWidth - Get specified monitor physical width in millimetres
func GetMonitorPhysicalWidth(monitor int) int {
	var ret ffi.Arg
	m := int32(monitor)
	getMonitorPhysicalWidth.Call(&ret, &m)
	return int(ret)
}

// GetMonitorPhysicalHeight - Get specified monitor physical height in millimetres
func GetMonitorPhysicalHeight(monitor int) int {
	var ret ffi.Arg
	m := int32(monitor)
	getMonitorPhysicalHeight.Call(&ret, &m)
	return int(ret)
}

// GetMonitorRefreshRate - Get specified monitor refresh rate
func GetMonitorRefreshRate(monitor int) int {
	var ret ffi.Arg
	m := int32(monitor)
	getMonitorRefreshRate.Call(&ret, &m)
	return int(ret)
}

// GetWindowPosition - Get window position XY on monitor
func GetWindowPosition() Vector2 {
	var ret Vector2
	getWindowPosition.Call(&ret)
	return ret
}

// GetWindowScaleDPI - Get window scale DPI factor
func GetWindowScaleDPI() Vector2 {
	var ret Vector2
	getWindowScaleDPI.Call(&ret)
	return ret
}

// GetMonitorName - Get the human-readable, UTF-8 encoded name of the specified monitor
func GetMonitorName(monitor int) string {
	var ret *byte
	m := int32(monitor)
	getMonitorName.Call(&ret, &m)
	return convert.ToString(ret)
}

// SetClipboardText - Set clipboard text content
func SetClipboardText(text string) {
	textPtr := convert.ToBytePtr(text)
	setClipboardText.Call(nil, &textPtr)
}

// GetClipboardText - Get clipboard text content
func GetClipboardText() string {
	var ret *byte
	getClipboardText.Call(&ret)
	return convert.ToString(ret)
}

// GetClipboardImage - Get clipboard image content
//
// NOTE: This function doesn't work on all platforms/backends.
func GetClipboardImage() Image {
	var ret Image
	getClipboardImage.Call(&ret)
	return ret
}

// EnableEventWaiting - Enable waiting for events on EndDrawing(), no automatic event polling
func EnableEventWaiting() {
	enableEventWaiting.Call(nil)
}

// DisableEventWaiting - Disable waiting for events on EndDrawing(), automatic events polling
func DisableEventWaiting() {
	disableEventWaiting.Call(nil)
}

// ShowCursor - Shows cursor
func ShowCursor() {
	showCursor.Call(nil)
}

// HideCursor - Hides cursor
func HideCursor() {
	hideCursor.Call(nil)
}

// IsCursorHidden - Check if cursor is not visible
func IsCursorHidden() bool {
	var ret ffi.Arg
	isCursorHidden.Call(&ret)
	return ret.Bool()
}

// EnableCursor - Enables cursor (unlock cursor)
func EnableCursor() {
	enableCursor.Call(nil)
}

// DisableCursor - Disables cursor (lock cursor)
func DisableCursor() {
	disableCursor.Call(nil)
}

// IsCursorOnScreen - Check if cursor is on the screen
func IsCursorOnScreen() bool {
	var ret ffi.Arg
	isCursorOnScreen.Call(&ret)
	return ret.Bool()
}

// ClearBackground - Set background color (framebuffer clear color)
func ClearBackground(col color.RGBA) {
	clearBackground.Call(nil, &col)
}

// BeginDrawing - Setup canvas (framebuffer) to start drawing
func BeginDrawing() {
	beginDrawing.Call(nil)
}

// EndDrawing - End canvas drawing and swap buffers (double buffering)
func EndDrawing() {
	endDrawing.Call(nil)
}

// BeginMode2D - Begin 2D mode with custom camera (2D)
func BeginMode2D(camera Camera2D) {
	beginMode2D.Call(nil, &camera)
}

// EndMode2D - Ends 2D mode with custom camera
func EndMode2D() {
	endMode2D.Call(nil)
}

// BeginMode3D - Begin 3D mode with custom camera (3D)
func BeginMode3D(camera Camera3D) {
	beginMode3D.Call(nil, &camera)
}

// EndMode3D - Ends 3D mode and returns to default 2D orthographic mode
func EndMode3D() {
	endMode3D.Call(nil)
}

// BeginTextureMode - Begin drawing to render texture
func BeginTextureMode(target RenderTexture2D) {
	beginTextureMode.Call(nil, &target)
}

// EndTextureMode - Ends drawing to render texture
func EndTextureMode() {
	endTextureMode.Call(nil)
}

// BeginShaderMode - Begin custom shader drawing
func BeginShaderMode(shader Shader) {
	beginShaderMode.Call(nil, &shader)
}

// EndShaderMode - End custom shader drawing (use default shader)
func EndShaderMode() {
	endShaderMode.Call(nil)
}

// BeginBlendMode - Begin blending mode (alpha, additive, multiplied, subtract, custom)
func BeginBlendMode(mode BlendMode) {
	beginBlendMode.Call(nil, &mode)
}

// EndBlendMode - End blending mode (reset to default: alpha blending)
func EndBlendMode() {
	endBlendMode.Call(nil)
}

// BeginScissorMode - Begin scissor mode (define screen area for following drawing)
func BeginScissorMode(x, y, width, height int32) {
	beginScissorMode.Call(nil, &x, &y, &width, &height)
}

// EndScissorMode - End scissor mode
func EndScissorMode() {
	endScissorMode.Call(nil)
}

// BeginVrStereoMode - Begin stereo rendering (requires VR simulator)
func BeginVrStereoMode(config VrStereoConfig) {
	beginVrStereoMode.Call(nil, &config)
}

// EndVrStereoMode - End stereo rendering (requires VR simulator)
func EndVrStereoMode() {
	endVrStereoMode.Call(nil)
}

// LoadVrStereoConfig - Load VR stereo config for VR simulator device parameters
func LoadVrStereoConfig(device VrDeviceInfo) VrStereoConfig {
	var ret VrStereoConfig
	loadVrStereoConfig.Call(&ret, &device)
	return ret
}

// UnloadVrStereoConfig - Unload VR stereo config
func UnloadVrStereoConfig(config VrStereoConfig) {
	unloadVrStereoConfig.Call(nil, &config)
}

// LoadShader - Load shader from files and bind default locations
func LoadShader(vsFileName string, fsFileName string) Shader {
	var ret Shader
	// "" becomes NULL to either load a vertex or fragment shader: https://github.com/gen2brain/raylib-go/issues/172
	vsFileNamePtr, fsFileNamePtr := convert.ToBytePtrNullable(vsFileName), convert.ToBytePtrNullable(fsFileName)
	loadShader.Call(&ret, &vsFileNamePtr, &fsFileNamePtr)
	return ret
}

// LoadShaderFromMemory - Load shader from code strings and bind default locations
func LoadShaderFromMemory(vsCode string, fsCode string) Shader {
	var ret Shader
	// "" becomes NULL to either load a vertex or fragment shader: https://github.com/gen2brain/raylib-go/issues/172
	vsCodePtr, fsCodePtr := convert.ToBytePtrNullable(vsCode), convert.ToBytePtrNullable(fsCode)
	loadShaderFromMemory.Call(&ret, &vsCodePtr, &fsCodePtr)
	return ret
}

// IsShaderValid - Check if a shader is valid (loaded on GPU)
func IsShaderValid(shader Shader) bool {
	var ret ffi.Arg
	isShaderValid.Call(&ret, &shader)
	return ret.Bool()
}

// GetShaderLocation - Get shader uniform location
func GetShaderLocation(shader Shader, uniformName string) int32 {
	var ret ffi.Arg
	uniformNamePtr := convert.ToBytePtr(uniformName)
	getShaderLocation.Call(&ret, &shader, &uniformNamePtr)
	return int32(ret)
}

// GetShaderLocationAttrib - Get shader attribute location
func GetShaderLocationAttrib(shader Shader, attribName string) int32 {
	var ret ffi.Arg
	attribNamePtr := convert.ToBytePtr(attribName)
	getShaderLocationAttrib.Call(&ret, &shader, &attribNamePtr)
	return int32(ret)
}

// SetShaderValue - Set shader uniform value
func SetShaderValue(shader Shader, locIndex int32, value []float32, uniformType ShaderUniformDataType) {
	valuePtr := &value[0]
	setShaderValue.Call(nil, &shader, &locIndex, &valuePtr, &uniformType)
}

// SetShaderValueV - Set shader uniform value vector
func SetShaderValueV(shader Shader, locIndex int32, value []float32, uniformType ShaderUniformDataType, count int32) {
	valuePtr := &value[0]
	setShaderValueV.Call(nil, &shader, &locIndex, &valuePtr, &uniformType, &count)
}

// SetShaderValueMatrix - Set shader uniform value (matrix 4x4)
func SetShaderValueMatrix(shader Shader, locIndex int32, mat Matrix) {
	setShaderValueMatrix.Call(nil, &shader, &locIndex, &mat)
}

// SetShaderValueTexture - Set shader uniform value for texture (sampler2d)
func SetShaderValueTexture(shader Shader, locIndex int32, texture Texture2D) {
	setShaderValueTexture.Call(nil, &shader, &locIndex, &texture)
}

// UnloadShader - Unload shader from GPU memory (VRAM)
func UnloadShader(shader Shader) {
	unloadShader.Call(nil, &shader)
}

// GetMouseRay - Get a ray trace from mouse position
//
// Deprecated: Use [GetScreenToWorldRay] instead.
func GetMouseRay(mousePosition Vector2, camera Camera) Ray {
	return GetScreenToWorldRay(mousePosition, camera)
}

// GetScreenToWorldRay - Get a ray trace from screen position (i.e mouse)
func GetScreenToWorldRay(position Vector2, camera Camera) Ray {
	var ret Ray
	getScreenToWorldRay.Call(&ret, &position, &camera)
	return ret
}

// GetScreenToWorldRayEx - Get a ray trace from screen position (i.e mouse) in a viewport
func GetScreenToWorldRayEx(position Vector2, camera Camera, width, height int32) Ray {
	var ret Ray
	getScreenToWorldRayEx.Call(&ret, &position, &camera, &width, &height)
	return ret
}

// GetWorldToScreen - Get the screen space position for a 3d world space position
func GetWorldToScreen(position Vector3, camera Camera) Vector2 {
	var ret Vector2
	getWorldToScreen.Call(&ret, &position, &camera)
	return ret
}

// GetWorldToScreenEx - Get size position for a 3d world space position
func GetWorldToScreenEx(position Vector3, camera Camera, width int32, height int32) Vector2 {
	var ret Vector2
	getWorldToScreenEx.Call(&ret, &position, &camera, &width, &height)
	return ret
}

// GetWorldToScreen2D - Get the screen space position for a 2d camera world space position
func GetWorldToScreen2D(position Vector2, camera Camera2D) Vector2 {
	var ret Vector2
	getWorldToScreen2D.Call(&ret, &position, &camera)
	return ret
}

// GetScreenToWorld2D - Get the world space position for a 2d camera screen space position
func GetScreenToWorld2D(position Vector2, camera Camera2D) Vector2 {
	var ret Vector2
	getScreenToWorld2D.Call(&ret, &position, &camera)
	return ret
}

// GetCameraMatrix - Get camera transform matrix (view matrix)
func GetCameraMatrix(camera Camera) Matrix {
	var ret Matrix
	getCameraMatrix.Call(&ret, &camera)
	return ret
}

// GetCameraMatrix2D - Get camera 2d transform matrix
func GetCameraMatrix2D(camera Camera2D) Matrix {
	var ret Matrix
	getCameraMatrix2D.Call(&ret, &camera)
	return ret
}

// SetTargetFPS - Set target FPS (maximum)
func SetTargetFPS(fps int32) {
	setTargetFPS.Call(nil, &fps)
}

// GetFrameTime - Get time in seconds for last frame drawn (delta time)
func GetFrameTime() float32 {
	var ret float32
	getFrameTime.Call(&ret)
	return ret
}

// GetTime - Get elapsed time in seconds since [InitWindow]
func GetTime() float64 {
	var ret float64
	getTime.Call(&ret)
	return ret
}

// GetFPS - Get current FPS
func GetFPS() int32 {
	var ret ffi.Arg
	getFPS.Call(&ret)
	return int32(ret)
}

// Custom frame control functions
// NOTE: SwapScreenBuffer and PollInputEvents are intended for advanced users that want full control over the frame processing
// By default EndDrawing() does this job: draws everything + SwapScreenBuffer() + manage frame timing + PollInputEvents()
// To avoid that behaviour and control frame processes manually you must recompile raylib with SUPPORT_CUSTOM_FRAME_CONTROL enabled in config.h
//
// See: https://github.com/gen2brain/raylib-go/issues/378

// SwapScreenBuffer - Swap back buffer with front buffer (screen drawing)
func SwapScreenBuffer() {
	swapScreenBuffer.Call(nil)
}

// PollInputEvents - Register all input events
func PollInputEvents() {
	pollInputEvents.Call(nil)
}

// WaitTime - Wait for some time (halt program execution)
func WaitTime(seconds float64) {
	waitTime.Call(nil, &seconds)
}

// GetRandomValue - Returns a random value between min and max (both included)
//
// Note: You can use go's math/rand package instead
func GetRandomValue(min, max int32) int32 {
	var ret ffi.Arg
	getRandomValue.Call(&ret, &min, &max)
	return int32(ret)
}

// TakeScreenshot - Takes a screenshot of current screen (filename extension defines format)
func TakeScreenshot(fileName string) {
	fileNamePtr := convert.ToBytePtr(fileName)
	takeScreenshot.Call(nil, &fileNamePtr)
}

// SetConfigFlags - Setup init configuration flags (view FLAGS)
func SetConfigFlags(flags uint32) {
	setConfigFlags.Call(nil, &flags)
}

// OpenURL - Open URL with default system browser (if available)
func OpenURL(url string) {
	urlPtr := convert.ToBytePtr(url)
	openURL.Call(nil, &urlPtr)
}

// TraceLog - Show trace log messages (LOG_DEBUG, LOG_INFO, LOG_WARNING, LOG_ERROR...)
func TraceLog(logLevel TraceLogLevel, text string, args ...any) {
	level := int32(logLevel)
	textPtr := convert.ToBytePtr(fmt.Sprintf(text, args...))
	traceLog.Call(nil, &level, &textPtr)
}

// SetTraceLogLevel - Set the current threshold (minimum) log level
func SetTraceLogLevel(logLevel TraceLogLevel) {
	level := int32(logLevel)
	setTraceLogLevel.Call(nil, &level)
}

// SetTraceLogCallback - Set custom trace log
func SetTraceLogCallback(fn TraceLogCallbackFun) {
	cb := traceLogCallbackWrapper(fn)
	setTraceLogCallback.Call(nil, &cb)
}

// IsFileDropped - Check if a file has been dropped into window
func IsFileDropped() bool {
	var ret ffi.Arg
	isFileDropped.Call(&ret)
	return ret.Bool()
}

// LoadDroppedFiles - Load dropped filepaths
func LoadDroppedFiles() []string {
	var filePathList = struct {
		count uint32
		paths **byte
	}{}
	loadDroppedFiles.Call(&filePathList)
	defer unloadDroppedFiles.Call(nil, &filePathList)

	paths := unsafe.Slice(filePathList.paths, filePathList.count)
	result := make([]string, filePathList.count)

	for i := range paths {
		result[i] = convert.ToString(paths[i])
	}

	return result
}

// UnloadDroppedFiles - Unload dropped filepaths
func UnloadDroppedFiles() {}

// LoadAutomationEventList - Load automation events list from file, NULL for empty list, capacity = MAX_AUTOMATION_EVENTS
func LoadAutomationEventList(fileName string) AutomationEventList {
	var ret AutomationEventList
	fileNamePtr := convert.ToBytePtr(fileName)
	loadAutomationEventList.Call(&ret, &fileNamePtr)
	return ret
}

// UnloadAutomationEventList - Unload automation events list from file
func UnloadAutomationEventList(list *AutomationEventList) {
	unloadAutomationEventList.Call(nil, list)
}

// ExportAutomationEventList - Export automation events list as text file
func ExportAutomationEventList(list AutomationEventList, fileName string) bool {
	var ret ffi.Arg
	fileNamePtr := convert.ToBytePtr(fileName)
	exportAutomationEventList.Call(&ret, &list, &fileNamePtr)
	return ret.Bool()
}

// SetAutomationEventList - Set automation event list to record to
func SetAutomationEventList(list *AutomationEventList) {
	setAutomationEventList.Call(nil, &list)
}

// SetAutomationEventBaseFrame - Set automation event internal base frame to start recording
func SetAutomationEventBaseFrame(frame int) {
	f := int32(frame)
	setAutomationEventBaseFrame.Call(nil, &f)
}

// StartAutomationEventRecording - Start recording automation events (AutomationEventList must be set)
func StartAutomationEventRecording() {
	startAutomationEventRecording.Call(nil)
}

// StopAutomationEventRecording - Stop recording automation events
func StopAutomationEventRecording() {
	stopAutomationEventRecording.Call(nil)
}

// PlayAutomationEvent - Play a recorded automation event
func PlayAutomationEvent(event AutomationEvent) {
	playAutomationEvent.Call(nil, &event)
}

// IsKeyPressed - Check if a key has been pressed once
func IsKeyPressed(key int32) bool {
	var ret ffi.Arg
	isKeyPressed.Call(&ret, &key)
	return ret.Bool()
}

// IsKeyPressedRepeat - Check if a key has been pressed again
func IsKeyPressedRepeat(key int32) bool {
	var ret ffi.Arg
	isKeyPressedRepeat.Call(&ret, &key)
	return ret.Bool()
}

// IsKeyDown - Check if a key is being pressed
func IsKeyDown(key int32) bool {
	var ret ffi.Arg
	isKeyDown.Call(&ret, &key)
	return ret.Bool()
}

// IsKeyReleased - Check if a key has been released once
func IsKeyReleased(key int32) bool {
	var ret ffi.Arg
	isKeyReleased.Call(&ret, &key)
	return ret.Bool()
}

// IsKeyUp - Check if a key is NOT being pressed
func IsKeyUp(key int32) bool {
	var ret ffi.Arg
	isKeyUp.Call(&ret, &key)
	return ret.Bool()
}

// GetKeyPressed - Get key pressed (keycode), call it multiple times for keys queued, returns 0 when the queue is empty
func GetKeyPressed() int32 {
	var ret ffi.Arg
	getKeyPressed.Call(&ret)
	return int32(ret)
}

// GetCharPressed - Get char pressed (unicode), call it multiple times for chars queued, returns 0 when the queue is empty
func GetCharPressed() int32 {
	var ret ffi.Arg
	getCharPressed.Call(&ret)
	return int32(ret)
}

// GetKeyName - Get name of a QWERTY key on the current keyboard layout
// (eg returns string 'q' for KEY_A on an AZERTY keyboard)
func GetKeyName(key int32) string {
	var ret *byte
	getKeyName.Call(&ret, &key)
	return convert.ToString(ret)
}

// SetExitKey - Set a custom key to exit program (default is ESC)
func SetExitKey(key int32) {
	setExitKey.Call(nil, &key)
}

// IsGamepadAvailable - Check if a gamepad is available
func IsGamepadAvailable(gamepad int32) bool {
	var ret ffi.Arg
	isGamepadAvailable.Call(&ret, &gamepad)
	return ret.Bool()
}

// GetGamepadName - Get gamepad internal name id
func GetGamepadName(gamepad int32) string {
	var ret *byte
	getGamepadName.Call(&ret, &gamepad)
	return convert.ToString(ret)
}

// IsGamepadButtonPressed - Check if a gamepad button has been pressed once
func IsGamepadButtonPressed(gamepad, button int32) bool {
	var ret ffi.Arg
	isGamepadButtonPressed.Call(&ret, &gamepad, &button)
	return ret.Bool()
}

// IsGamepadButtonDown - Check if a gamepad button is being pressed
func IsGamepadButtonDown(gamepad, button int32) bool {
	var ret ffi.Arg
	isGamepadButtonDown.Call(&ret, &gamepad, &button)
	return ret.Bool()
}

// IsGamepadButtonReleased - Check if a gamepad button has been released once
func IsGamepadButtonReleased(gamepad, button int32) bool {
	var ret ffi.Arg
	isGamepadButtonReleased.Call(&ret, &gamepad, &button)
	return ret.Bool()
}

// IsGamepadButtonUp - Check if a gamepad button is NOT being pressed
func IsGamepadButtonUp(gamepad, button int32) bool {
	var ret ffi.Arg
	isGamepadButtonUp.Call(&ret, &gamepad, &button)
	return ret.Bool()
}

// GetGamepadButtonPressed - Get the last gamepad button pressed
func GetGamepadButtonPressed() int32 {
	var ret ffi.Arg
	getGamepadButtonPressed.Call(&ret)
	return int32(ret)
}

// GetGamepadAxisCount - Get gamepad axis count for a gamepad
func GetGamepadAxisCount(gamepad int32) int32 {
	var ret ffi.Arg
	getGamepadAxisCount.Call(&ret, &gamepad)
	return int32(ret)
}

// GetGamepadAxisMovement - Get axis movement value for a gamepad axis
func GetGamepadAxisMovement(gamepad, axis int32) float32 {
	var ret float32
	getGamepadAxisMovement.Call(&ret, &gamepad, &axis)
	return ret
}

// SetGamepadMappings - Set internal gamepad mappings (SDL_GameControllerDB)
func SetGamepadMappings(mappings string) int32 {
	var ret ffi.Arg
	mappingsPtr := convert.ToBytePtr(mappings)
	setGamepadMappings.Call(&ret, &mappingsPtr)
	return int32(ret)
}

// SetGamepadVibration - Set gamepad vibration for both motors (duration in seconds)
func SetGamepadVibration(gamepad int32, leftMotor, rightMotor, duration float32) {
	setGamepadVibration.Call(nil, &gamepad, &leftMotor, &rightMotor, &duration)
}

// IsMouseButtonPressed - Check if a mouse button has been pressed once
func IsMouseButtonPressed(button MouseButton) bool {
	var ret ffi.Arg
	isMouseButtonPressed.Call(&ret, &button)
	return ret.Bool()
}

// IsMouseButtonDown - Check if a mouse button is being pressed
func IsMouseButtonDown(button MouseButton) bool {
	var ret ffi.Arg
	isMouseButtonDown.Call(&ret, &button)
	return ret.Bool()
}

// IsMouseButtonReleased - Check if a mouse button has been released once
func IsMouseButtonReleased(button MouseButton) bool {
	var ret ffi.Arg
	isMouseButtonReleased.Call(&ret, &button)
	return ret.Bool()
}

// IsMouseButtonUp - Check if a mouse button is NOT being pressed
func IsMouseButtonUp(button MouseButton) bool {
	var ret ffi.Arg
	isMouseButtonUp.Call(&ret, &button)
	return ret.Bool()
}

// GetMouseX - Get mouse position X
func GetMouseX() int32 {
	var ret ffi.Arg
	getMouseX.Call(&ret)
	return int32(ret)
}

// GetMouseY - Get mouse position Y
func GetMouseY() int32 {
	var ret ffi.Arg
	getMouseY.Call(&ret)
	return int32(ret)
}

// GetMousePosition - Get mouse position XY
func GetMousePosition() Vector2 {
	var ret Vector2
	getMousePosition.Call(&ret)
	return ret
}

// GetMouseDelta - Get mouse delta between frames
func GetMouseDelta() Vector2 {
	var ret Vector2
	getMouseDelta.Call(&ret)
	return ret
}

// SetMousePosition - Set mouse position XY
func SetMousePosition(x, y int) {
	posX, poxY := int32(x), int32(y)
	setMousePosition.Call(nil, &posX, &poxY)
}

// SetMouseOffset - Set mouse offset
func SetMouseOffset(offsetX, offsetY int) {
	x, y := int32(offsetX), int32(offsetY)
	setMouseOffset.Call(nil, &x, &y)
}

// SetMouseScale - Set mouse scaling
func SetMouseScale(scaleX, scaleY float32) {
	setMouseScale.Call(nil, &scaleX, &scaleY)
}

// GetMouseWheelMove - Get mouse wheel movement for X or Y, whichever is larger
func GetMouseWheelMove() float32 {
	var ret float32
	getMouseWheelMove.Call(&ret)
	return ret
}

// GetMouseWheelMoveV - Get mouse wheel movement for both X and Y
func GetMouseWheelMoveV() Vector2 {
	var ret Vector2
	getMouseWheelMoveV.Call(&ret)
	return ret
}

// SetMouseCursor - Set mouse cursor
func SetMouseCursor(cursor int32) {
	setMouseCursor.Call(nil, &cursor)
}

// GetTouchX - Get touch position X for touch point 0 (relative to screen size)
func GetTouchX() int32 {
	var ret ffi.Arg
	getTouchX.Call(&ret)
	return int32(ret)
}

// GetTouchY - Get touch position Y for touch point 0 (relative to screen size)
func GetTouchY() int32 {
	var ret ffi.Arg
	getTouchY.Call(&ret)
	return int32(ret)
}

// GetTouchPosition - Get touch position XY for a touch point index (relative to screen size)
func GetTouchPosition(index int32) Vector2 {
	var ret Vector2
	getTouchPosition.Call(&ret, &index)
	return ret
}

// GetTouchPointId - Get touch point identifier for given index
func GetTouchPointId(index int32) int32 {
	var ret ffi.Arg
	getTouchPointId.Call(&ret, &index)
	return int32(ret)
}

// GetTouchPointCount - Get number of touch points
func GetTouchPointCount() int32 {
	var ret ffi.Arg
	getTouchPointCount.Call(&ret)
	return int32(ret)
}

// SetGesturesEnabled - Enable a set of gestures using flags
func SetGesturesEnabled(flags uint32) {
	setGesturesEnabled.Call(nil, &flags)
}

// IsGestureDetected - Check if a gesture have been detected
func IsGestureDetected(gesture Gestures) bool {
	var ret ffi.Arg
	isGestureDetected.Call(&ret, &gesture)
	return ret.Bool()
}

// GetGestureDetected - Get latest detected gesture
func GetGestureDetected() Gestures {
	var ret ffi.Arg
	getGestureDetected.Call(&ret)
	return Gestures(ret)
}

// GetGestureHoldDuration - Get gesture hold time in milliseconds
func GetGestureHoldDuration() float32 {
	var ret float32
	getGestureHoldDuration.Call(&ret)
	return ret
}

// GetGestureDragVector - Get gesture drag vector
func GetGestureDragVector() Vector2 {
	var ret Vector2
	getGestureDragVector.Call(&ret)
	return ret
}

// GetGestureDragAngle - Get gesture drag angle
func GetGestureDragAngle() float32 {
	var ret float32
	getGestureDragAngle.Call(&ret)
	return ret
}

// GetGesturePinchVector - Get gesture pinch delta
func GetGesturePinchVector() Vector2 {
	var ret Vector2
	getGesturePinchVector.Call(&ret)
	return ret
}

// GetGesturePinchAngle - Get gesture pinch angle
func GetGesturePinchAngle() float32 {
	var ret float32
	getGesturePinchAngle.Call(&ret)
	return ret
}

// SetShapesTexture - Set texture and rectangle to be used on shapes drawing
func SetShapesTexture(texture Texture2D, source Rectangle) {
	setShapesTexture.Call(nil, &texture, &source)
}

// GetShapesTexture - Get texture that is used for shapes drawing
func GetShapesTexture() Texture2D {
	var ret Texture2D
	getShapesTexture.Call(&ret)
	return ret
}

// GetShapesTextureRectangle - Get texture source rectangle that is used for shapes drawing
func GetShapesTextureRectangle() Rectangle {
	var ret Rectangle
	getShapesTextureRectangle.Call(&ret)
	return ret
}

// DrawPixel - Draw a pixel
func DrawPixel(posX, posY int32, col color.RGBA) {
	drawPixel.Call(nil, &posX, &posY, &col)
}

// DrawPixelV - Draw a pixel (Vector version)
func DrawPixelV(position Vector2, col color.RGBA) {
	drawPixelV.Call(nil, &position, &col)
}

// DrawLine - Draw a line
func DrawLine(startPosX, startPosY, endPosX, endPosY int32, col color.RGBA) {
	drawLine.Call(nil, &startPosX, &startPosY, &endPosX, &endPosY, &col)
}

// DrawLineV - Draw a line (using gl lines)
func DrawLineV(startPos, endPos Vector2, col color.RGBA) {
	drawLineV.Call(nil, &startPos, &endPos, &col)
}

// DrawLineEx - Draw a line (using triangles/quads)
func DrawLineEx(startPos, endPos Vector2, thick float32, col color.RGBA) {
	drawLineEx.Call(nil, &startPos, &endPos, &thick, &col)
}

// DrawLineStrip - Draw lines sequence (using gl lines)
func DrawLineStrip(points []Vector2, col color.RGBA) {
	pointCount := int32(len(points))
	pointsPtr := unsafe.SliceData(points)
	drawLineStrip.Call(nil, &pointsPtr, &pointCount, &col)
}

// DrawLineBezier - Draw line segment cubic-bezier in-out interpolation
func DrawLineBezier(startPos, endPos Vector2, thick float32, col color.RGBA) {
	drawLineBezier.Call(nil, &startPos, &endPos, &thick, &col)
}

// DrawLineDashed - Draw a dashed line
func DrawLineDashed(startPos, endPos Vector2, dashSize, spaceSize int32, col color.RGBA) {
	drawLineDashed.Call(nil, &startPos, &endPos, &dashSize, &spaceSize, &col)
}

// DrawCircle - Draw a color-filled circle
func DrawCircle(centerX, centerY int32, radius float32, col color.RGBA) {
	drawCircle.Call(nil, &centerX, &centerY, &radius, &col)
}

// DrawCircleV - Draw a color-filled circle (Vector version)
func DrawCircleV(center Vector2, radius float32, col color.RGBA) {
	drawCircleV.Call(nil, &center, &radius, &col)
}

// DrawCircleGradient - Draw a gradient-filled circle
func DrawCircleGradient(center Vector2, radius float32, inner, outer color.RGBA) {
	drawCircleGradient.Call(nil, &center, &radius, &inner, &outer)
}

// DrawCircleSector - Draw a piece of a circle
func DrawCircleSector(center Vector2, radius, startAngle, endAngle float32, segments int32, col color.RGBA) {
	drawCircleSector.Call(nil, &center, &radius, &startAngle, &endAngle, &segments, &col)
}

// DrawCircleSectorLines - Draw circle sector outline
func DrawCircleSectorLines(center Vector2, radius, startAngle, endAngle float32, segments int32, col color.RGBA) {
	drawCircleSectorLines.Call(nil, &center, &radius, &startAngle, &endAngle, &segments, &col)
}

// DrawCircleLines - Draw circle outline
func DrawCircleLines(centerX, centerY int32, radius float32, col color.RGBA) {
	drawCircleLines.Call(nil, &centerX, &centerY, &radius, &col)
}

// DrawCircleLinesV - Draw circle outline (Vector version)
func DrawCircleLinesV(center Vector2, radius float32, col color.RGBA) {
	drawCircleLinesV.Call(nil, &center, &radius, &col)
}

// DrawEllipse - Draw ellipse
func DrawEllipse(centerX, centerY int32, radiusH, radiusV float32, col color.RGBA) {
	drawEllipse.Call(nil, &centerX, &centerY, &radiusH, &radiusV, &col)
}

// DrawEllipseV - Draw ellipse (Vector version)
func DrawEllipseV(center Vector2, radiusH, radiusV float32, col color.RGBA) {
	drawEllipseV.Call(nil, &center, &radiusH, &radiusV, &col)
}

// DrawEllipseLines - Draw ellipse outline
func DrawEllipseLines(centerX, centerY int32, radiusH, radiusV float32, col color.RGBA) {
	drawEllipseLines.Call(nil, &centerX, &centerY, &radiusH, &radiusV, &col)
}

// DrawEllipseLinesV - Draw ellipse outline (Vector version)
func DrawEllipseLinesV(center Vector2, radiusH, radiusV float32, col color.RGBA) {
	drawEllipseLinesV.Call(nil, &center, &radiusH, &radiusV, &col)
}

// DrawRing - Draw ring
func DrawRing(center Vector2, innerRadius, outerRadius, startAngle, endAngle float32, segments int32, col color.RGBA) {
	drawRing.Call(nil, &center, &innerRadius, &outerRadius, &startAngle, &endAngle, &segments, &col)
}

// DrawRingLines - Draw ring outline
func DrawRingLines(center Vector2, innerRadius, outerRadius, startAngle, endAngle float32, segments int32, col color.RGBA) {
	drawRingLines.Call(nil, &center, &innerRadius, &outerRadius, &startAngle, &endAngle, &segments, &col)
}

// DrawRectangle - Draw a color-filled rectangle
func DrawRectangle(posX, posY, width, height int32, col color.RGBA) {
	drawRectangle.Call(nil, &posX, &posY, &width, &height, &col)
}

// DrawRectangleV - Draw a color-filled rectangle (Vector version)
func DrawRectangleV(position Vector2, size Vector2, col color.RGBA) {
	drawRectangleV.Call(nil, &position, &size, &col)
}

// DrawRectangleRec - Draw a color-filled rectangle
func DrawRectangleRec(rec Rectangle, col color.RGBA) {
	drawRectangleRec.Call(nil, &rec, &col)
}

// DrawRectanglePro - Draw a color-filled rectangle with pro parameters
func DrawRectanglePro(rec Rectangle, origin Vector2, rotation float32, col color.RGBA) {
	drawRectanglePro.Call(nil, &rec, &origin, &rotation, &col)
}

// DrawRectangleGradientV - Draw a vertical-gradient-filled rectangle
func DrawRectangleGradientV(posX, posY, width, height int32, top, bottom color.RGBA) {
	drawRectangleGradientV.Call(nil, &posX, &posY, &width, &height, &top, &bottom)
}

// DrawRectangleGradientH - Draw a horizontal-gradient-filled rectangle
func DrawRectangleGradientH(posX, posY, width, height int32, left, right color.RGBA) {
	drawRectangleGradientH.Call(nil, &posX, &posY, &width, &height, &left, &right)
}

// DrawRectangleGradientEx - Draw a gradient-filled rectangle with custom vertex colors
func DrawRectangleGradientEx(rec Rectangle, topLeft, bottomLeft, bottomRight, topRight color.RGBA) {
	drawRectangleGradientEx.Call(nil, &rec, &topLeft, &bottomLeft, &bottomRight, &topRight)
}

// DrawRectangleLines - Draw rectangle outline
func DrawRectangleLines(posX, posY, width, height int32, col color.RGBA) {
	drawRectangleLines.Call(nil, &posX, &posY, &width, &height, &col)
}

// DrawRectangleLinesEx - Draw rectangle outline with extended parameters
func DrawRectangleLinesEx(rec Rectangle, lineThick float32, col color.RGBA) {
	drawRectangleLinesEx.Call(nil, &rec, &lineThick, &col)
}

// DrawRectangleRounded - Draw rectangle with rounded edges
func DrawRectangleRounded(rec Rectangle, roundness float32, segments int32, col color.RGBA) {
	drawRectangleRounded.Call(nil, &rec, &roundness, &segments, &col)
}

// DrawRectangleRoundedLines - Draw rectangle lines with rounded edges
func DrawRectangleRoundedLines(rec Rectangle, roundness float32, segments int32, col color.RGBA) {
	drawRectangleRoundedLines.Call(nil, &rec, &roundness, &segments, &col)
}

// DrawRectangleRoundedLinesEx - Draw rectangle with rounded edges outline
func DrawRectangleRoundedLinesEx(rec Rectangle, roundness float32, segments int32, lineThick float32, col color.RGBA) {
	drawRectangleRoundedLinesEx.Call(nil, &rec, &roundness, &segments, &lineThick, &col)
}

// DrawTriangle - Draw a color-filled triangle (vertex in counter-clockwise order!)
func DrawTriangle(v1, v2, v3 Vector2, col color.RGBA) {
	drawTriangle.Call(nil, &v1, &v2, &v3, &col)
}

// DrawTriangleLines - Draw triangle outline (vertex in counter-clockwise order!)
func DrawTriangleLines(v1, v2, v3 Vector2, col color.RGBA) {
	drawTriangleLines.Call(nil, &v1, &v2, &v3, &col)
}

// DrawTriangleFan - Draw a triangle fan defined by points (first vertex is the center)
func DrawTriangleFan(points []Vector2, col color.RGBA) {
	pointCount := int32(len(points))
	pointsPtr := unsafe.SliceData(points)
	drawTriangleFan.Call(nil, &pointsPtr, &pointCount, &col)
}

// DrawTriangleStrip - Draw a triangle strip defined by points
func DrawTriangleStrip(points []Vector2, col color.RGBA) {
	pointCount := int32(len(points))
	pointsPtr := unsafe.SliceData(points)
	drawTriangleStrip.Call(nil, &pointsPtr, &pointCount, &col)
}

// DrawPoly - Draw a regular polygon (Vector version)
func DrawPoly(center Vector2, sides int32, radius, rotation float32, col color.RGBA) {
	drawPoly.Call(nil, &center, &sides, &radius, &rotation, &col)
}

// DrawPolyLines - Draw a polygon outline of n sides
func DrawPolyLines(center Vector2, sides int32, radius, rotation float32, col color.RGBA) {
	drawPolyLines.Call(nil, &center, &sides, &radius, &rotation, &col)
}

// DrawPolyLinesEx - Draw a polygon outline of n sides with extended parameters
func DrawPolyLinesEx(center Vector2, sides int32, radius float32, rotation float32, lineThick float32, col color.RGBA) {
	drawPolyLinesEx.Call(nil, &center, &sides, &radius, &rotation, &lineThick, &col)
}

// DrawSplineLinear - Draw spline: Linear, minimum 2 points
func DrawSplineLinear(points []Vector2, thick float32, col color.RGBA) {
	pointCount := int32(len(points))
	pointsPtr := unsafe.SliceData(points)
	drawSplineLinear.Call(nil, &pointsPtr, &pointCount, &thick, &col)
}

// DrawSplineBasis - Draw spline: B-Spline, minimum 4 points
func DrawSplineBasis(points []Vector2, thick float32, col color.RGBA) {
	pointCount := int32(len(points))
	pointsPtr := unsafe.SliceData(points)
	drawSplineBasis.Call(nil, &pointsPtr, &pointCount, &thick, &col)
}

// DrawSplineCatmullRom - Draw spline: Catmull-Rom, minimum 4 points
func DrawSplineCatmullRom(points []Vector2, thick float32, col color.RGBA) {
	pointCount := int32(len(points))
	pointsPtr := unsafe.SliceData(points)
	drawSplineCatmullRom.Call(nil, &pointsPtr, &pointCount, &thick, &col)
}

// DrawSplineBezierQuadratic - Draw spline: Quadratic Bezier, minimum 3 points (1 control point): [p1, c2, p3, c4...]
func DrawSplineBezierQuadratic(points []Vector2, thick float32, col color.RGBA) {
	pointCount := int32(len(points))
	pointsPtr := unsafe.SliceData(points)
	drawSplineBezierQuadratic.Call(nil, &pointsPtr, &pointCount, &thick, &col)
}

// DrawSplineBezierCubic - Draw spline: Cubic Bezier, minimum 4 points (2 control points): [p1, c2, c3, p4, c5, c6...]
func DrawSplineBezierCubic(points []Vector2, thick float32, col color.RGBA) {
	pointCount := int32(len(points))
	pointsPtr := unsafe.SliceData(points)
	drawSplineBezierCubic.Call(nil, &pointsPtr, &pointCount, &thick, &col)
}

// DrawSplineSegmentLinear - Draw spline segment: Linear, 2 points
func DrawSplineSegmentLinear(p1, p2 Vector2, thick float32, col color.RGBA) {
	drawSplineSegmentLinear.Call(nil, &p1, &p2, &thick, &col)
}

// DrawSplineSegmentBasis - Draw spline segment: B-Spline, 4 points
func DrawSplineSegmentBasis(p1, p2, p3, p4 Vector2, thick float32, col color.RGBA) {
	drawSplineSegmentBasis.Call(nil, &p1, &p2, &p3, &p4, &thick, &col)
}

// DrawSplineSegmentCatmullRom - Draw spline segment: Catmull-Rom, 4 points
func DrawSplineSegmentCatmullRom(p1, p2, p3, p4 Vector2, thick float32, col color.RGBA) {
	drawSplineSegmentCatmullRom.Call(nil, &p1, &p2, &p3, &p4, &thick, &col)
}

// DrawSplineSegmentBezierQuadratic - Draw spline segment: Quadratic Bezier, 2 points, 1 control point
func DrawSplineSegmentBezierQuadratic(p1, c2, p3 Vector2, thick float32, col color.RGBA) {
	drawSplineSegmentBezierQuadratic.Call(nil, &p1, &c2, &p3, &thick, &col)
}

// DrawSplineSegmentBezierCubic - Draw spline segment: Cubic Bezier, 2 points, 2 control points
func DrawSplineSegmentBezierCubic(p1, c2, c3, p4 Vector2, thick float32, col color.RGBA) {
	drawSplineSegmentBezierCubic.Call(nil, &p1, &c2, &c3, &p4, &thick, &col)
}

// GetSplinePointLinear - Get (evaluate) spline point: Linear
func GetSplinePointLinear(startPos, endPos Vector2, t float32) Vector2 {
	var ret Vector2
	getSplinePointLinear.Call(&ret, &startPos, &endPos, &t)
	return ret
}

// GetSplinePointBasis - Get (evaluate) spline point: B-Spline
func GetSplinePointBasis(p1, p2, p3, p4 Vector2, t float32) Vector2 {
	var ret Vector2
	getSplinePointBasis.Call(&ret, &p1, &p2, &p3, &p4, &t)
	return ret
}

// GetSplinePointCatmullRom - Get (evaluate) spline point: Catmull-Rom
func GetSplinePointCatmullRom(p1, p2, p3, p4 Vector2, t float32) Vector2 {
	var ret Vector2
	getSplinePointCatmullRom.Call(&ret, &p1, &p2, &p3, &p4, &t)
	return ret
}

// GetSplinePointBezierQuad - Get (evaluate) spline point: Quadratic Bezier
func GetSplinePointBezierQuad(p1, c2, p3 Vector2, t float32) Vector2 {
	var ret Vector2
	getSplinePointBezierQuad.Call(&ret, &p1, &c2, &p3, &t)
	return ret
}

// GetSplinePointBezierCubic - Get (evaluate) spline point: Cubic Bezier
func GetSplinePointBezierCubic(p1, c2, c3, p4 Vector2, t float32) Vector2 {
	var ret Vector2
	getSplinePointBezierCubic.Call(&ret, &p1, &c2, &c3, &p4, &t)
	return ret
}

// CheckCollisionRecs - Check collision between two rectangles
func CheckCollisionRecs(rec1, rec2 Rectangle) bool {
	var ret ffi.Arg
	checkCollisionRecs.Call(&ret, &rec1, &rec2)
	return ret.Bool()
}

// CheckCollisionCircles - Check collision between two circles
func CheckCollisionCircles(center1 Vector2, radius1 float32, center2 Vector2, radius2 float32) bool {
	var ret ffi.Arg
	checkCollisionCircles.Call(&ret, &center1, &radius1, &center2, &radius2)
	return ret.Bool()
}

// CheckCollisionCircleRec - Check collision between circle and rectangle
func CheckCollisionCircleRec(center Vector2, radius float32, rec Rectangle) bool {
	var ret ffi.Arg
	checkCollisionCircleRec.Call(&ret, &center, &radius, &rec)
	return ret.Bool()
}

// CheckCollisionCircleLine - Check if circle collides with a line created betweeen two points [p1] and [p2]
func CheckCollisionCircleLine(center Vector2, radius float32, p1, p2 Vector2) bool {
	var ret ffi.Arg
	checkCollisionCircleLine.Call(&ret, &center, &radius, &p1, &p2)
	return ret.Bool()
}

// CheckCollisionPointRec - Check if point is inside rectangle
func CheckCollisionPointRec(point Vector2, rec Rectangle) bool {
	var ret ffi.Arg
	checkCollisionPointRec.Call(&ret, &point, &rec)
	return ret.Bool()
}

// CheckCollisionPointCircle - Check if point is inside circle
func CheckCollisionPointCircle(point Vector2, center Vector2, radius float32) bool {
	var ret ffi.Arg
	checkCollisionPointCircle.Call(&ret, &point, &center, &radius)
	return ret.Bool()
}

// CheckCollisionPointTriangle - Check if point is inside a triangle
func CheckCollisionPointTriangle(point, p1, p2, p3 Vector2) bool {
	var ret ffi.Arg
	checkCollisionPointTriangle.Call(&ret, &point, &p1, &p2, &p3)
	return ret.Bool()
}

// CheckCollisionPointLine - Check if point belongs to line created between two points [p1] and [p2] with defined margin in pixels [threshold]
func CheckCollisionPointLine(point, p1, p2 Vector2, threshold int32) bool {
	var ret ffi.Arg
	checkCollisionPointLine.Call(&ret, &point, &p1, &p2, &threshold)
	return ret.Bool()
}

// CheckCollisionPointPoly - Check if point is within a polygon described by array of vertices
//
// NOTE: Based on http://jeffreythompson.org/collision-detection/poly-point.php
func CheckCollisionPointPoly(point Vector2, points []Vector2) bool {
	var ret ffi.Arg
	pointCount := int32(len(points))
	pointsPtr := unsafe.SliceData(points)
	checkCollisionPointPoly.Call(&ret, &point, &pointsPtr, &pointCount)
	return ret.Bool()
}

// CheckCollisionLines - Check the collision between two lines defined by two points each, returns collision point by reference
func CheckCollisionLines(startPos1, endPos1, startPos2, endPos2 Vector2, collisionPoint *Vector2) bool {
	var ret ffi.Arg
	checkCollisionLines.Call(&ret, &startPos1, &endPos1, &startPos2, &endPos2, &collisionPoint)
	return ret.Bool()
}

// GetCollisionRec - Get collision rectangle for two rectangles collision
func GetCollisionRec(rec1, rec2 Rectangle) Rectangle {
	var ret Rectangle
	getCollisionRec.Call(&ret, &rec1, &rec2)
	return ret
}

// LoadImage - Load image from file into CPU memory (RAM)
func LoadImage(fileName string) *Image {
	var ret Image
	fileNamePtr := convert.ToBytePtr(fileName)
	loadImage.Call(&ret, &fileNamePtr)
	return &ret
}

// LoadImageRaw - Load image from RAW file data
func LoadImageRaw(fileName string, width, height int32, format PixelFormat, headerSize int32) *Image {
	var ret Image
	fileNamePtr := convert.ToBytePtr(fileName)
	loadImageRaw.Call(&ret, &fileNamePtr, &width, &height, &format, &headerSize)
	return &ret
}

// LoadImageAnim - Load image sequence from file (frames appended to image.data)
func LoadImageAnim(fileName string, frames *int32) *Image {
	var ret Image
	fileNamePtr := convert.ToBytePtr(fileName)
	loadImageAnim.Call(&ret, &fileNamePtr, &frames)
	return &ret
}

// LoadImageAnimFromMemory - Load image sequence from memory buffer
func LoadImageAnimFromMemory(fileType string, fileData []byte, dataSize int32, frames *int32) *Image {
	var ret Image
	fileTypePtr := convert.ToBytePtr(fileType)
	fileDataPtr := unsafe.SliceData(fileData)
	loadImageAnimFromMemory.Call(&ret, &fileTypePtr, &fileDataPtr, &dataSize, &frames)
	return &ret
}

// LoadImageFromMemory - Load image from memory buffer, fileType refers to extension: i.e. '.png'
func LoadImageFromMemory(fileType string, fileData []byte, dataSize int32) *Image {
	var ret Image
	fileTypePtr := convert.ToBytePtr(fileType)
	fileDataPtr := unsafe.SliceData(fileData)
	loadImageFromMemory.Call(&ret, &fileTypePtr, &fileDataPtr, &dataSize)
	return &ret
}

// LoadImageFromTexture - Load image from GPU texture data
func LoadImageFromTexture(texture Texture2D) *Image {
	var ret Image
	loadImageFromTexture.Call(&ret, &texture)
	return &ret
}

// LoadImageFromScreen - Load image from screen buffer and (screenshot)
func LoadImageFromScreen() *Image {
	var ret Image
	loadImageFromScreen.Call(&ret)
	return &ret
}

// IsImageValid - Check if an image is valid (data and parameters)
func IsImageValid(image *Image) bool {
	var ret ffi.Arg
	isImageValid.Call(&ret, image)
	return ret.Bool()
}

// UnloadImage - Unload image from CPU memory (RAM)
func UnloadImage(image *Image) {
	unloadImage.Call(nil, image)
}

// ExportImage - Export image data to file, returns true on success
func ExportImage(image Image, fileName string) bool {
	var ret ffi.Arg
	fileNamePtr := convert.ToBytePtr(fileName)
	exportImage.Call(&ret, &image, &fileNamePtr)
	return ret.Bool()
}

// ExportImageToMemory - Export image to memory buffer
//
// The returned memory is a Go-managed slice. It doesn't need to be freed.
func ExportImageToMemory(image Image, fileType string) []byte {
	var ret *byte
	fileTypePtr := convert.ToBytePtr(fileType)
	var fileSize int32
	fileSizePtr := &fileSize
	exportImageToMemory.Call(&ret, &image, &fileTypePtr, &fileSizePtr)
	defer memFree.Call(nil, &ret) // free the memory

	if ret == nil {
		return nil
	}
	result := make([]byte, fileSize)
	copy(result, unsafe.Slice(ret, fileSize))
	return result
}

// GenImageColor - Generate image: plain color
func GenImageColor(width, height int, col color.RGBA) *Image {
	var ret Image
	w, h := int32(width), int32(height)
	genImageColor.Call(&ret, &w, &h, &col)
	return &ret
}

// GenImageGradientLinear - Generate image: linear gradient, direction in degrees [0..360], 0=Vertical gradient
func GenImageGradientLinear(width, height, direction int, start, end color.RGBA) *Image {
	var ret Image
	w, h, dir := int32(width), int32(height), int32(direction)
	genImageGradientLinear.Call(&ret, &w, &h, &dir, &start, &end)
	return &ret
}

// GenImageGradientRadial - Generate image: radial gradient
func GenImageGradientRadial(width, height int, density float32, inner, outer color.RGBA) *Image {
	var ret Image
	w, h := int32(width), int32(height)
	genImageGradientRadial.Call(&ret, &w, &h, &density, &inner, &outer)
	return &ret
}

// GenImageGradientSquare - Generate image: square gradient
func GenImageGradientSquare(width, height int, density float32, inner, outer color.RGBA) *Image {
	var ret Image
	w, h := int32(width), int32(height)
	genImageGradientSquare.Call(&ret, &w, &h, &density, &inner, &outer)
	return &ret
}

// GenImageChecked - Generate image: checked
func GenImageChecked(width, height, checksX, checksY int, col1, col2 color.RGBA) *Image {
	var ret Image
	w, h, cX, cY := int32(width), int32(height), int32(checksX), int32(checksY)
	genImageChecked.Call(&ret, &w, &h, &cX, &cY, &col1, &col2)
	return &ret
}

// GenImageWhiteNoise - Generate image: white noise
func GenImageWhiteNoise(width, height int, factor float32) *Image {
	var ret Image
	w, h := int32(width), int32(height)
	genImageWhiteNoise.Call(&ret, &w, &h, &factor)
	return &ret
}

// GenImagePerlinNoise - Generate image: perlin noise
func GenImagePerlinNoise(width, height, offsetX, offsetY int, scale float32) *Image {
	var ret Image
	w, h, oX, oY := int32(width), int32(height), int32(offsetX), int32(offsetY)
	genImagePerlinNoise.Call(&ret, &w, &h, &oX, &oY, &scale)
	return &ret
}

// GenImageCellular - Generate image: cellular algorithm, bigger tileSize means bigger cells
func GenImageCellular(width, height, tileSize int) *Image {
	var ret Image
	w, h, tS := int32(width), int32(height), int32(tileSize)
	genImageCellular.Call(&ret, &w, &h, &tS)
	return &ret
}

// GenImageText - Generate image: grayscale image from text data
func GenImageText(width, height int, text string) *Image {
	var ret Image
	w, h, textPtr := int32(width), int32(height), convert.ToBytePtr(text)
	genImageText.Call(&ret, &w, &h, &textPtr)
	return &ret
}

// ImageCopy - Create an image duplicate (useful for transformations)
func ImageCopy(image *Image) *Image {
	var ret Image
	imageCopy.Call(&ret, image)
	return &ret
}

// ImageFromImage - Create an image from another image piece
func ImageFromImage(image Image, rec Rectangle) Image {
	var ret Image
	imageFromImage.Call(&ret, &image, &rec)
	return ret
}

// ImageFromChannel - Create an image from a selected channel of another image (GRAYSCALE)
func ImageFromChannel(image Image, selectedChannel int32) Image {
	var ret Image
	imageFromChannel.Call(&ret, &image, &selectedChannel)
	return ret
}

// ImageText - Create an image from text (default font)
func ImageText(text string, fontSize int32, col color.RGBA) *Image {
	var ret Image
	textPtr := convert.ToBytePtr(text)
	imageText.Call(&ret, &textPtr, &fontSize, &col)
	return &ret
}

// ImageTextEx - Create an image from text (custom sprite font)
func ImageTextEx(font Font, text string, fontSize, spacing float32, tint color.RGBA) *Image {
	var ret Image
	textPtr := convert.ToBytePtr(text)
	imageTextEx.Call(&ret, &font, &textPtr, &fontSize, &spacing, &tint)
	return &ret
}

// ImageFormat - Convert image data to desired format
func ImageFormat(image *Image, newFormat PixelFormat) {
	imageFormat.Call(nil, &image, &newFormat)
}

// ImageToPOT - Convert image to POT (power-of-two)
func ImageToPOT(image *Image, fill color.RGBA) {
	imageToPOT.Call(nil, &image, &fill)
}

// ImageCrop - Crop an image to a defined rectangle
func ImageCrop(image *Image, crop Rectangle) {
	imageCrop.Call(nil, &image, &crop)
}

// ImageAlphaCrop - Crop image depending on alpha value
func ImageAlphaCrop(image *Image, threshold float32) {
	imageAlphaCrop.Call(nil, &image, &threshold)
}

// ImageAlphaClear - Clear alpha channel to desired color
func ImageAlphaClear(image *Image, col color.RGBA, threshold float32) {
	imageAlphaClear.Call(nil, &image, &col, &threshold)
}

// ImageAlphaMask - Apply alpha mask to image
func ImageAlphaMask(image, alphaMask *Image) {
	imageAlphaMask.Call(nil, &image, alphaMask)
}

// ImageAlphaPremultiply - Premultiply alpha channel
func ImageAlphaPremultiply(image *Image) {
	imageAlphaPremultiply.Call(nil, &image)
}

// ImageBlurGaussian - Apply Gaussian blur using a box blur approximation
func ImageBlurGaussian(image *Image, blurSize int32) {
	imageBlurGaussian.Call(nil, &image, &blurSize)
}

// ImageKernelConvolution - Apply custom square convolution kernel to image
func ImageKernelConvolution(image *Image, kernel []float32) {
	kernelPtr := unsafe.SliceData(kernel)
	kernelSize := int32(len(kernel))
	imageKernelConvolution.Call(nil, &image, &kernelPtr, &kernelSize)
}

// ImageResize - Resize image (Bicubic scaling algorithm)
func ImageResize(image *Image, newWidth, newHeight int32) {
	imageResize.Call(nil, &image, &newWidth, &newHeight)
}

// ImageResizeNN - Resize image (Nearest-Neighbor scaling algorithm)
func ImageResizeNN(image *Image, newWidth, newHeight int32) {
	imageResizeNN.Call(nil, &image, &newWidth, &newHeight)
}

// ImageResizeCanvas - Resize canvas and fill with color
func ImageResizeCanvas(image *Image, newWidth, newHeight, offsetX, offsetY int32, fill color.RGBA) {
	imageResizeCanvas.Call(nil, &image, &newWidth, &newHeight, &offsetX, &offsetY, &fill)
}

// ImageMipmaps - Compute all mipmap levels for a provided image
func ImageMipmaps(image *Image) {
	imageMipmaps.Call(nil, &image)
}

// ImageDither - Dither image data to 16bpp or lower (Floyd-Steinberg dithering)
func ImageDither(image *Image, rBpp, gBpp, bBpp, aBpp int32) {
	imageDither.Call(nil, &image, &rBpp, &gBpp, &bBpp, &aBpp)
}

// ImageFlipVertical - Flip image vertically
func ImageFlipVertical(image *Image) {
	imageFlipVertical.Call(nil, &image)
}

// ImageFlipHorizontal - Flip image horizontally
func ImageFlipHorizontal(image *Image) {
	imageFlipHorizontal.Call(nil, &image)
}

// ImageRotate - Rotate image by input angle in degrees (-359 to 359)
func ImageRotate(image *Image, degrees int32) {
	imageRotate.Call(nil, &image, &degrees)
}

// ImageRotateCW - Rotate image clockwise 90deg
func ImageRotateCW(image *Image) {
	imageRotateCW.Call(nil, &image)
}

// ImageRotateCCW - Rotate image counter-clockwise 90deg
func ImageRotateCCW(image *Image) {
	imageRotateCCW.Call(nil, &image)
}

// ImageColorTint - Modify image color: tint
func ImageColorTint(image *Image, col color.RGBA) {
	imageColorTint.Call(nil, &image, &col)
}

// ImageColorInvert - Modify image color: invert
func ImageColorInvert(image *Image) {
	imageColorInvert.Call(nil, &image)
}

// ImageColorGrayscale - Modify image color: grayscale
func ImageColorGrayscale(image *Image) {
	imageColorGrayscale.Call(nil, &image)
}

// ImageColorContrast - Modify image color: contrast (-100 to 100)
func ImageColorContrast(image *Image, contrast float32) {
	imageColorContrast.Call(nil, &image, &contrast)
}

// ImageColorBrightness - Modify image color: brightness (-255 to 255)
func ImageColorBrightness(image *Image, brightness int32) {
	imageColorBrightness.Call(nil, &image, &brightness)
}

// ImageColorReplace - Modify image color: replace color
func ImageColorReplace(image *Image, col, replace color.RGBA) {
	imageColorReplace.Call(nil, &image, &col, &replace)
}

// LoadImageColors - Load color data from image as a Color array (RGBA - 32bit)
//
// NOTE: Memory allocated should be freed using UnloadImageColors()
func LoadImageColors(img *Image) []color.RGBA {
	var ret *color.RGBA
	loadImageColors.Call(&ret, img)
	return unsafe.Slice(ret, img.Width*img.Height)
}

// LoadImagePalette - Load colors palette from image as a Color array (RGBA - 32bit)
//
// NOTE: Memory allocated should be freed using UnloadImagePalette()
func LoadImagePalette(img Image, maxPaletteSize int32) []color.RGBA {
	var colorCount int32
	colorCountPtr := &colorCount
	var ret *color.RGBA
	loadImagePalette.Call(&ret, &img, &maxPaletteSize, &colorCountPtr)
	return unsafe.Slice(ret, colorCount)
}

// UnloadImageColors - Unload color data loaded with LoadImageColors()
func UnloadImageColors(colors []color.RGBA) {
	colorsPtr := unsafe.SliceData(colors)
	unloadImageColors.Call(nil, &colorsPtr)
}

// UnloadImagePalette - Unload colors palette loaded with LoadImagePalette()
func UnloadImagePalette(colors []color.RGBA) {
	colorsPtr := unsafe.SliceData(colors)
	unloadImagePalette.Call(nil, &colorsPtr)
}

// GetImageAlphaBorder - Get image alpha border rectangle
func GetImageAlphaBorder(image Image, threshold float32) Rectangle {
	var ret Rectangle
	getImageAlphaBorder.Call(&ret, &image, &threshold)
	return ret
}

// GetImageColor - Get image pixel color at (x, y) position
func GetImageColor(image Image, x, y int32) color.RGBA {
	var ret color.RGBA
	getImageColor.Call(&ret, &image, &x, &y)
	return ret
}

// ImageClearBackground - Clear image background with given color
func ImageClearBackground(dst *Image, col color.RGBA) {
	imageClearBackground.Call(nil, &dst, &col)
}

// ImageDrawPixel - Draw pixel within an image
func ImageDrawPixel(dst *Image, posX, posY int32, col color.RGBA) {
	imageDrawPixel.Call(nil, &dst, &posX, &posY, &col)
}

// ImageDrawPixelV - Draw pixel within an image (Vector version)
func ImageDrawPixelV(dst *Image, position Vector2, col color.RGBA) {
	imageDrawPixelV.Call(nil, &dst, &position, &col)
}

// ImageDrawLine - Draw line within an image
func ImageDrawLine(dst *Image, startPosX, startPosY, endPosX, endPosY int32, col color.RGBA) {
	imageDrawLine.Call(nil, &dst, &startPosX, &startPosY, &endPosX, &endPosY, &col)
}

// ImageDrawLineV - Draw line within an image (Vector version)
func ImageDrawLineV(dst *Image, start, end Vector2, col color.RGBA) {
	imageDrawLineV.Call(nil, &dst, &start, &end, &col)
}

// ImageDrawLineEx - Draw a line defining thickness within an image
func ImageDrawLineEx(dst *Image, start, end Vector2, thick int32, col color.RGBA) {
	imageDrawLineEx.Call(nil, &dst, &start, &end, &thick, &col)
}

// ImageDrawCircle - Draw a filled circle within an image
func ImageDrawCircle(dst *Image, centerX, centerY, radius int32, col color.RGBA) {
	imageDrawCircle.Call(nil, &dst, &centerX, &centerY, &radius, &col)
}

// ImageDrawCircleV - Draw a filled circle within an image (Vector version)
func ImageDrawCircleV(dst *Image, center Vector2, radius int32, col color.RGBA) {
	imageDrawCircleV.Call(nil, &dst, &center, &radius, &col)
}

// ImageDrawCircleLines - Draw circle outline within an image
func ImageDrawCircleLines(dst *Image, centerX, centerY, radius int32, col color.RGBA) {
	imageDrawCircleLines.Call(nil, &dst, &centerX, &centerY, &radius, &col)
}

// ImageDrawCircleLinesV - Draw circle outline within an image (Vector version)
func ImageDrawCircleLinesV(dst *Image, center Vector2, radius int32, col color.RGBA) {
	imageDrawCircleLinesV.Call(nil, &dst, &center, &radius, &col)
}

// ImageDrawRectangle - Draw rectangle within an image
func ImageDrawRectangle(dst *Image, posX, posY, width, height int32, col color.RGBA) {
	imageDrawRectangle.Call(nil, &dst, &posX, &posY, &width, &height, &col)
}

// ImageDrawRectangleV - Draw rectangle within an image (Vector version)
func ImageDrawRectangleV(dst *Image, position, size Vector2, col color.RGBA) {
	imageDrawRectangleV.Call(nil, &dst, &position, &size, &col)
}

// ImageDrawRectangleRec - Draw rectangle within an image
func ImageDrawRectangleRec(dst *Image, rec Rectangle, col color.RGBA) {
	imageDrawRectangleRec.Call(nil, &dst, &rec, &col)
}

// ImageDrawRectangleLines - Draw rectangle lines within an image
func ImageDrawRectangleLines(dst *Image, rec Rectangle, thick int, col color.RGBA) {
	t := int32(thick)
	imageDrawRectangleLines.Call(nil, &dst, &rec, &t, &col)
}

// ImageDrawTriangle - Draw triangle within an image
func ImageDrawTriangle(dst *Image, v1, v2, v3 Vector2, col color.RGBA) {
	imageDrawTriangle.Call(nil, &dst, &v1, &v2, &v3, &col)
}

// ImageDrawTriangleEx - Draw triangle with interpolated colors within an image
func ImageDrawTriangleEx(dst *Image, v1, v2, v3 Vector2, c1, c2, c3 color.RGBA) {
	imageDrawTriangleEx.Call(nil, &dst, &v1, &v2, &v3, &c1, &c2, &c3)
}

// ImageDrawTriangleLines - Draw triangle outline within an image
func ImageDrawTriangleLines(dst *Image, v1, v2, v3 Vector2, col color.RGBA) {
	imageDrawTriangleLines.Call(nil, &dst, &v1, &v2, &v3, &col)
}

// ImageDrawTriangleFan - Draw a triangle fan defined by points within an image (first vertex is the center)
func ImageDrawTriangleFan(dst *Image, points []Vector2, col color.RGBA) {
	pointCount := int32(len(points))
	pointsPtr := unsafe.SliceData(points)
	imageDrawTriangleFan.Call(nil, &dst, &pointsPtr, &pointCount, &col)
}

// ImageDrawTriangleStrip - Draw a triangle strip defined by points within an image
func ImageDrawTriangleStrip(dst *Image, points []Vector2, col color.RGBA) {
	pointCount := int32(len(points))
	pointsPtr := unsafe.SliceData(points)
	imageDrawTriangleStrip.Call(nil, &dst, &pointsPtr, &pointCount, &col)
}

// ImageDraw - Draw a source image within a destination image (tint applied to source)
func ImageDraw(dst, src *Image, srcRec, dstRec Rectangle, tint color.RGBA) {
	imageDraw.Call(nil, &dst, src, &srcRec, &dstRec, &tint)
}

// ImageDrawText - Draw text (using default font) within an image (destination)
func ImageDrawText(dst *Image, posX, posY int32, text string, fontSize int32, col color.RGBA) {
	textPtr := convert.ToBytePtr(text)
	imageDrawText.Call(nil, &dst, &textPtr, &posX, &posY, &fontSize, &col)
}

// ImageDrawTextEx - Draw text (custom sprite font) within an image (destination)
func ImageDrawTextEx(dst *Image, position Vector2, font Font, text string, fontSize, spacing float32, tint color.RGBA) {
	textPtr := convert.ToBytePtr(text)
	imageDrawTextEx.Call(nil, &dst, &font, &textPtr, &position, &fontSize, &spacing, &tint)
}

// LoadTexture - Load texture from file into GPU memory (VRAM)
func LoadTexture(fileName string) Texture2D {
	var ret Texture2D
	fileNamePtr := convert.ToBytePtr(fileName)
	loadTexture.Call(&ret, &fileNamePtr)
	return ret
}

// LoadTextureFromImage - Load texture from image data
func LoadTextureFromImage(image *Image) Texture2D {
	var ret Texture2D
	loadTextureFromImage.Call(&ret, image)
	return ret
}

// LoadTextureCubemap - Load cubemap from image, multiple image cubemap layouts supported
func LoadTextureCubemap(image *Image, layout int32) Texture2D {
	var ret Texture2D
	loadTextureCubemap.Call(&ret, image, &layout)
	return ret
}

// LoadRenderTexture - Load texture for rendering (framebuffer)
func LoadRenderTexture(width, height int32) RenderTexture2D {
	var ret RenderTexture2D
	loadRenderTexture.Call(&ret, &width, &height)
	return ret
}

// IsTextureValid - Check if a texture is valid (loaded in GPU)
func IsTextureValid(texture Texture2D) bool {
	var ret ffi.Arg
	isTextureValid.Call(&ret, &texture)
	return ret.Bool()
}

// UnloadTexture - Unload texture from GPU memory (VRAM)
func UnloadTexture(texture Texture2D) {
	unloadTexture.Call(nil, &texture)
}

// IsRenderTextureValid - Check if a render texture is valid (loaded in GPU)
func IsRenderTextureValid(target RenderTexture2D) bool {
	var ret ffi.Arg
	isRenderTextureValid.Call(&ret, &target)
	return ret.Bool()
}

// UnloadRenderTexture - Unload render texture from GPU memory (VRAM)
func UnloadRenderTexture(target RenderTexture2D) {
	unloadRenderTexture.Call(nil, &target)
}

// UpdateTexture - Update GPU texture with new data (pixels should be able to fill texture)
func UpdateTexture[T []color.RGBA | *image.RGBA | []byte](texture Texture2D, pixels T) {
	var cpixels unsafe.Pointer
	switch p := any(pixels).(type) {
	case []color.RGBA:
		cpixels = unsafe.Pointer(&p[0])
	case *image.RGBA:
		cpixels = unsafe.Pointer(&p.Pix[0])
	case []byte:
		cpixels = unsafe.Pointer(&p[0])
	}
	updateTexture.Call(nil, &texture, &cpixels)
}

// UpdateTextureRec - Update GPU texture rectangle with new data (pixels and rec should fit in texture)
func UpdateTextureRec[T []color.RGBA | *image.RGBA | []byte](texture Texture2D, rec Rectangle, pixels T) {
	var cpixels unsafe.Pointer
	switch p := any(pixels).(type) {
	case []color.RGBA:
		cpixels = unsafe.Pointer(&p[0])
	case *image.RGBA:
		cpixels = unsafe.Pointer(&p.Pix[0])
	case []byte:
		cpixels = unsafe.Pointer(&p[0])
	}
	updateTextureRec.Call(nil, &texture, &rec, &cpixels)
}

// GenTextureMipmaps - Generate GPU mipmaps for a texture
func GenTextureMipmaps(texture *Texture2D) {
	genTextureMipmaps.Call(nil, &texture)
}

// SetTextureFilter - Set texture scaling filter mode
func SetTextureFilter(texture Texture2D, filter TextureFilterMode) {
	setTextureFilter.Call(nil, &texture, &filter)
}

// SetTextureWrap - Set texture wrapping mode
func SetTextureWrap(texture Texture2D, wrap TextureWrapMode) {
	setTextureWrap.Call(nil, &texture, &wrap)
}

// DrawTexture - Draw a Texture2D
func DrawTexture(texture Texture2D, posX int32, posY int32, tint color.RGBA) {
	drawTexture.Call(nil, &texture, &posX, &posY, &tint)
}

// DrawTextureV - Draw a Texture2D with position defined as Vector2
func DrawTextureV(texture Texture2D, position Vector2, tint color.RGBA) {
	drawTextureV.Call(nil, &texture, &position, &tint)
}

// DrawTextureEx - Draw a Texture2D with extended parameters
func DrawTextureEx(texture Texture2D, position Vector2, rotation, scale float32, tint color.RGBA) {
	drawTextureEx.Call(nil, &texture, &position, &rotation, &scale, &tint)
}

// DrawTextureRec - Draw a part of a texture defined by a rectangle
func DrawTextureRec(texture Texture2D, source Rectangle, position Vector2, tint color.RGBA) {
	drawTextureRec.Call(nil, &texture, &source, &position, &tint)
}

// DrawTexturePro - Draw a part of a texture defined by a rectangle with 'pro' parameters
func DrawTexturePro(texture Texture2D, source, dest Rectangle, origin Vector2, rotation float32, tint color.RGBA) {
	drawTexturePro.Call(nil, &texture, &source, &dest, &origin, &rotation, &tint)
}

// DrawTextureNPatch - Draws a texture (or part of it) that stretches or shrinks nicely
func DrawTextureNPatch(texture Texture2D, nPatchInfo NPatchInfo, dest Rectangle, origin Vector2, rotation float32, tint color.RGBA) {
	drawTextureNPatch.Call(nil, &texture, &nPatchInfo, &dest, &origin, &rotation, &tint)
}

// Fade - Get color with alpha applied, alpha goes from 0.0f to 1.0f
func Fade(col color.RGBA, alpha float32) color.RGBA {
	var ret color.RGBA
	fade.Call(&ret, &col, &alpha)
	return ret
}

// ColorToInt - Get hexadecimal value for a Color (0xRRGGBBAA)
func ColorToInt(col color.RGBA) int32 {
	var ret ffi.Arg
	colorToInt.Call(&ret, &col)
	return int32(ret)
}

// ColorNormalize - Get Color normalized as float [0..1]
func ColorNormalize(col color.RGBA) Vector4 {
	var ret Vector4
	colorNormalize.Call(&ret, &col)
	return ret
}

// ColorFromNormalized - Get Color from normalized values [0..1]
func ColorFromNormalized(normalized Vector4) color.RGBA {
	var ret color.RGBA
	colorFromNormalized.Call(&ret, &normalized)
	return ret
}

// ColorToHSV - Get HSV values for a Color, hue [0..360], saturation/value [0..1]
func ColorToHSV(col color.RGBA) Vector3 {
	var ret Vector3
	colorToHSV.Call(&ret, &col)
	return ret
}

// ColorFromHSV - Get a Color from HSV values, hue [0..360], saturation/value [0..1]
func ColorFromHSV(hue, saturation, value float32) color.RGBA {
	var ret color.RGBA
	colorFromHSV.Call(&ret, &hue, &saturation, &value)
	return ret
}

// ColorTint - Get color multiplied with another color
func ColorTint(col color.RGBA, tint color.RGBA) color.RGBA {
	var ret color.RGBA
	colorTint.Call(&ret, &col, &tint)
	return ret
}

// ColorBrightness - Get color with brightness correction, brightness factor goes from -1.0f to 1.0f
func ColorBrightness(col color.RGBA, factor float32) color.RGBA {
	var ret color.RGBA
	colorBrightness.Call(&ret, &col, &factor)
	return ret
}

// ColorContrast - Get color with contrast correction, contrast values between -1.0f and 1.0f
func ColorContrast(col color.RGBA, contrast float32) color.RGBA {
	var ret color.RGBA
	colorContrast.Call(&ret, &col, &contrast)
	return ret
}

// ColorAlpha - Get color with alpha applied, alpha goes from 0.0f to 1.0f
func ColorAlpha(col color.RGBA, alpha float32) color.RGBA {
	var ret color.RGBA
	colorAlpha.Call(&ret, &col, &alpha)
	return ret
}

// ColorAlphaBlend - Get src alpha-blended into dst color with tint
func ColorAlphaBlend(dst, src, tint color.RGBA) color.RGBA {
	var ret color.RGBA
	colorAlphaBlend.Call(&ret, &dst, &src, &tint)
	return ret
}

// ColorLerp - Get color lerp interpolation between two colors, factor [0.0f..1.0f]
func ColorLerp(col1, col2 color.RGBA, factor float32) color.RGBA {
	var ret color.RGBA
	colorLerp.Call(&ret, &col1, &col2, &factor)
	return ret
}

// GetColor - Get Color structure from hexadecimal value
func GetColor(hexValue uint) color.RGBA {
	var ret color.RGBA
	getColor.Call(&ret, &hexValue)
	return ret
}

// GetPixelColor - Get Color from a source pixel pointer of certain format
func GetPixelColor(srcPtr unsafe.Pointer, format int32) color.RGBA {
	var ret color.RGBA
	getPixelColor.Call(&ret, &srcPtr, &format)
	return ret
}

// SetPixelColor - Set color formatted into destination pixel pointer
func SetPixelColor(dstPtr unsafe.Pointer, col color.RGBA, format int32) {
	setPixelColor.Call(nil, &dstPtr, &col, &format)
}

// GetPixelDataSize - Get pixel data size in bytes for certain format
func GetPixelDataSize(width, height, format int32) int32 {
	var ret ffi.Arg
	getPixelDataSize.Call(&ret, &width, &height, &format)
	return int32(ret)
}

// GetFontDefault - Get the default Font
func GetFontDefault() Font {
	var ret Font
	getFontDefault.Call(&ret)
	return ret
}

// LoadFont - Load font from file into GPU memory (VRAM)
func LoadFont(fileName string) Font {
	var ret Font
	fileNamePtr := convert.ToBytePtr(fileName)
	loadFont.Call(&ret, &fileNamePtr)
	return ret
}

// LoadFontEx - Load Font from file with extended parameters.
// Use nil for fontChars and 0 for runesNumber to load the default character set.
// Font size is provided in pixels height.
func LoadFontEx(fileName string, fontSize int32, fontChars []rune, runesNumber ...int32) Font {
	var ret Font
	fileNamePtr := convert.ToBytePtr(fileName)
	var charsCount int32
	var fontCharsPtr *rune
	if fontChars != nil {
		fontCharsPtr = unsafe.SliceData(fontChars)
		charsCount = int32(len(fontChars))
	} else {
		if len(runesNumber) > 0 {
			charsCount = int32(runesNumber[0])
		}
	}

	loadFontEx.Call(&ret, &fileNamePtr, &fontSize, &fontCharsPtr, &charsCount)
	return ret
}

// LoadFontFromImage - Load font from Image (XNA style)
func LoadFontFromImage(image Image, key color.RGBA, firstChar rune) Font {
	var ret Font
	loadFontFromImage.Call(&ret, &image, &key, &firstChar)
	return ret
}

// LoadFontFromMemory - Load font from memory buffer, fileType refers to extension: i.e. '.ttf'
func LoadFontFromMemory(fileType string, fileData []byte, fontSize int32, codepoints []rune) Font {
	var ret Font
	dataSize := int32(len(fileData))
	fileDataPtr := unsafe.SliceData(fileData)
	codepointCount := int32(len(codepoints))
	codepointsPtr := unsafe.SliceData(codepoints)
	fileTypePtr := convert.ToBytePtr(fileType)
	loadFontFromMemory.Call(&ret, &fileTypePtr, &fileDataPtr, &dataSize, &fontSize, &codepointsPtr, &codepointCount)
	return ret
}

// IsFontValid - Check if a font is valid (font data loaded, WARNING: GPU texture not checked)
func IsFontValid(font Font) bool {
	var ret ffi.Arg
	isFontValid.Call(&ret, &font)
	return ret.Bool()
}

// LoadFontData - Load font data for further use
func LoadFontData(fileData []byte, fontSize int32, codepoints []rune, codepointCount, typ int32) []GlyphInfo {
	dataSize := int32(len(fileData))
	fileDataPtr := unsafe.SliceData(fileData)
	if codepointCount < 0 {
		codepointCount = 0
	}
	codepointsPtr := unsafe.SliceData(codepoints)
	var glyphCount int32
	glyphCountPtr := &glyphCount
	var ret *GlyphInfo
	loadFontData.Call(&ret, &fileDataPtr, &dataSize, &fontSize, &codepointsPtr, &codepointCount, &typ, &glyphCountPtr)
	if ret == nil || glyphCount <= 0 {
		return nil
	}
	return unsafe.Slice(ret, glyphCount)
}

// GenImageFontAtlas - Generate image font atlas using chars info
func GenImageFontAtlas(glyphs []GlyphInfo, glyphRecs []*Rectangle, fontSize int32, padding int32, packMethod int32) Image {
	var ret Image
	glyphCount := int32(len(glyphs))
	glyphsPtr := unsafe.SliceData(glyphs)
	glyphRecsPtr := unsafe.SliceData(glyphRecs)
	genImageFontAtlas.Call(&ret, &glyphsPtr, &glyphRecsPtr, &glyphCount, &fontSize, &padding, &packMethod)
	return ret
}

// UnloadFontData - Unload font chars info data (RAM)
func UnloadFontData(glyphs []GlyphInfo) {
	glyphCount := int32(len(glyphs))
	glyphsPtr := unsafe.SliceData(glyphs)
	unloadFontData.Call(nil, &glyphsPtr, &glyphCount)
}

// UnloadFont - Unload font from GPU memory (VRAM)
func UnloadFont(font Font) {
	unloadFont.Call(nil, &font)
}

// DrawFPS - Draw current FPS
func DrawFPS(posX int32, posY int32) {
	drawFPS.Call(nil, &posX, &posY)
}

// DrawText - Draw text (using default font)
func DrawText(text string, posX int32, posY int32, fontSize int32, col color.RGBA) {
	textPtr := convert.ToBytePtr(text)
	drawText.Call(nil, &textPtr, &posX, &posY, &fontSize, &col)
}

// DrawTextEx - Draw text using font and additional parameters
func DrawTextEx(font Font, text string, position Vector2, fontSize float32, spacing float32, tint color.RGBA) {
	textPtr := convert.ToBytePtr(text)
	drawTextEx.Call(nil, &font, &textPtr, &position, &fontSize, &spacing, &tint)
}

// DrawTextPro - Draw text using Font and pro parameters (rotation)
func DrawTextPro(font Font, text string, position, origin Vector2, rotation, fontSize, spacing float32, tint color.RGBA) {
	textPtr := convert.ToBytePtr(text)
	drawTextPro.Call(nil, &font, &textPtr, &position, &origin, &rotation, &fontSize, &spacing, &tint)
}

// DrawTextCodepoint - Draw one character (codepoint)
func DrawTextCodepoint(font Font, codepoint rune, position Vector2, fontSize float32, tint color.RGBA) {
	drawTextCodepoint.Call(nil, &font, &codepoint, &position, &fontSize, &tint)
}

// DrawTextCodepoints - Draw multiple character (codepoint)
func DrawTextCodepoints(font Font, codepoints []rune, position Vector2, fontSize float32, spacing float32, tint color.RGBA) {
	codepointCount := int32(len(codepoints))
	codepointsPtr := unsafe.SliceData(codepoints)
	drawTextCodepoints.Call(nil, &font, &codepointsPtr, &codepointCount, &position, &fontSize, &spacing, &tint)
}

// SetTextLineSpacing - Set vertical line spacing when drawing with line-breaks
func SetTextLineSpacing(spacing int) {
	s := int32(spacing)
	setTextLineSpacing.Call(nil, &s)
}

// MeasureText - Measure string width for default font
func MeasureText(text string, fontSize int32) int32 {
	textPtr := convert.ToBytePtr(text)
	var ret ffi.Arg
	measureText.Call(&ret, &textPtr, &fontSize)
	return int32(ret)
}

// MeasureTextEx - Measure string size for Font
func MeasureTextEx(font Font, text string, fontSize float32, spacing float32) Vector2 {
	textPtr := convert.ToBytePtr(text)
	var ret Vector2
	measureTextEx.Call(&ret, &font, &textPtr, &fontSize, &spacing)
	return ret
}

// MeasureTextCodepoints - Measure string size for an existing array of codepoints for Font
func MeasureTextCodepoints(font Font, codepoints []rune, fontSize float32, spacing float32) Vector2 {
	length := int32(len(codepoints))
	codepointsPtr := unsafe.SliceData(codepoints)
	var ret Vector2
	measureTextCodepoints.Call(&ret, &font, &codepointsPtr, &length, &fontSize, &spacing)
	return ret
}

// GetGlyphIndex - Get glyph index position in font for a codepoint (unicode character), fallback to '?' if not found
func GetGlyphIndex(font Font, codepoint rune) int32 {
	var ret ffi.Arg
	getGlyphIndex.Call(&ret, &font, &codepoint)
	return int32(ret)
}

// GetGlyphInfo - Get glyph font info data for a codepoint (unicode character), fallback to '?' if not found
func GetGlyphInfo(font Font, codepoint rune) GlyphInfo {
	var ret GlyphInfo
	getGlyphInfo.Call(&ret, &font, &codepoint)
	return ret
}

// GetGlyphAtlasRec - Get glyph rectangle in font atlas for a codepoint (unicode character), fallback to '?' if not found
func GetGlyphAtlasRec(font Font, codepoint rune) Rectangle {
	var ret Rectangle
	getGlyphAtlasRec.Call(&ret, &font, &codepoint)
	return ret
}

// DrawLine3D - Draw a line in 3D world space
func DrawLine3D(startPos Vector3, endPos Vector3, col color.RGBA) {
	drawLine3D.Call(nil, &startPos, &endPos, &col)
}

// DrawPoint3D - Draw a point in 3D space, actually a small line
func DrawPoint3D(position Vector3, col color.RGBA) {
	drawPoint3D.Call(nil, &position, &col)
}

// DrawCircle3D - Draw a circle in 3D world space
func DrawCircle3D(center Vector3, radius float32, rotationAxis Vector3, rotationAngle float32, col color.RGBA) {
	drawCircle3D.Call(nil, &center, &radius, &rotationAxis, &rotationAngle, &col)
}

// DrawTriangle3D - Draw a color-filled triangle (vertex in counter-clockwise order!)
func DrawTriangle3D(v1 Vector3, v2 Vector3, v3 Vector3, col color.RGBA) {
	drawTriangle3D.Call(nil, &v1, &v2, &v3, &col)
}

// DrawTriangleStrip3D - Draw a triangle strip defined by points
func DrawTriangleStrip3D(points []Vector3, col color.RGBA) {
	pointCount := int32(len(points))
	pointsPtr := unsafe.SliceData(points)
	drawTriangleStrip3D.Call(nil, &pointsPtr, &pointCount, &col)
}

// DrawCube - Draw cube
func DrawCube(position Vector3, width float32, height float32, length float32, col color.RGBA) {
	drawCube.Call(nil, &position, &width, &height, &length, &col)
}

// DrawCubeV - Draw cube (Vector version)
func DrawCubeV(position Vector3, size Vector3, col color.RGBA) {
	drawCubeV.Call(nil, &position, &size, &col)
}

// DrawCubeWires - Draw cube wires
func DrawCubeWires(position Vector3, width float32, height float32, length float32, col color.RGBA) {
	drawCubeWires.Call(nil, &position, &width, &height, &length, &col)
}

// DrawCubeWiresV - Draw cube wires (Vector version)
func DrawCubeWiresV(position Vector3, size Vector3, col color.RGBA) {
	drawCubeWiresV.Call(nil, &position, &size, &col)
}

// DrawSphere - Draw sphere
func DrawSphere(centerPos Vector3, radius float32, col color.RGBA) {
	drawSphere.Call(nil, &centerPos, &radius, &col)
}

// DrawSphereEx - Draw sphere with extended parameters
func DrawSphereEx(centerPos Vector3, radius float32, rings int32, slices int32, col color.RGBA) {
	drawSphereEx.Call(nil, &centerPos, &radius, &rings, &slices, &col)
}

// DrawSphereWires - Draw sphere wires
func DrawSphereWires(centerPos Vector3, radius float32, rings int32, slices int32, col color.RGBA) {
	drawSphereWires.Call(nil, &centerPos, &radius, &rings, &slices, &col)
}

// DrawCylinder - Draw a cylinder/cone
func DrawCylinder(position Vector3, radiusTop float32, radiusBottom float32, height float32, slices int32, col color.RGBA) {
	drawCylinder.Call(nil, &position, &radiusTop, &radiusBottom, &height, &slices, &col)
}

// DrawCylinderEx - Draw a cylinder with base at startPos and top at endPos
func DrawCylinderEx(startPos Vector3, endPos Vector3, startRadius float32, endRadius float32, sides int32, col color.RGBA) {
	drawCylinderEx.Call(nil, &startPos, &endPos, &startRadius, &endRadius, &sides, &col)
}

// DrawCylinderWires - Draw a cylinder/cone wires
func DrawCylinderWires(position Vector3, radiusTop float32, radiusBottom float32, height float32, slices int32, col color.RGBA) {
	drawCylinderWires.Call(nil, &position, &radiusTop, &radiusBottom, &height, &slices, &col)
}

// DrawCylinderWiresEx - Draw a cylinder wires with base at startPos and top at endPos
func DrawCylinderWiresEx(startPos Vector3, endPos Vector3, startRadius float32, endRadius float32, sides int32, col color.RGBA) {
	drawCylinderWiresEx.Call(nil, &startPos, &endPos, &startRadius, &endRadius, &sides, &col)
}

// DrawCapsule - Draw a capsule with the center of its sphere caps at startPos and endPos
func DrawCapsule(startPos, endPos Vector3, radius float32, slices, rings int32, col color.RGBA) {
	drawCapsule.Call(nil, &startPos, &endPos, &radius, &slices, &rings, &col)
}

// DrawCapsuleWires - Draw capsule wireframe with the center of its sphere caps at startPos and endPos
func DrawCapsuleWires(startPos, endPos Vector3, radius float32, slices, rings int32, col color.RGBA) {
	drawCapsuleWires.Call(nil, &startPos, &endPos, &radius, &slices, &rings, &col)
}

// DrawPlane - Draw a plane XZ
func DrawPlane(centerPos Vector3, size Vector2, col color.RGBA) {
	drawPlane.Call(nil, &centerPos, &size, &col)
}

// DrawRay - Draw a ray line
func DrawRay(ray Ray, col color.RGBA) {
	drawRay.Call(nil, &ray, &col)
}

// DrawGrid - Draw a grid (centered at (0, 0, 0))
func DrawGrid(slices int32, spacing float32) {
	drawGrid.Call(nil, &slices, &spacing)
}

// LoadModel - Load model from files (meshes and materials)
func LoadModel(fileName string) Model {
	var ret Model
	fileNamePtr := convert.ToBytePtr(fileName)
	loadModel.Call(&ret, &fileNamePtr)
	return ret
}

// LoadModelFromMesh - Load model from generated mesh (default material)
func LoadModelFromMesh(mesh Mesh) Model {
	var ret Model
	loadModelFromMesh.Call(&ret, &mesh)
	return ret
}

// IsModelValid - Check if a model is valid (loaded in GPU, VAO/VBOs)
func IsModelValid(model Model) bool {
	var ret ffi.Arg
	isModelValid.Call(&ret, &model)
	return ret.Bool()
}

// UnloadModel - Unload model (including meshes) from memory (RAM and/or VRAM)
func UnloadModel(model Model) {
	unloadModel.Call(nil, &model)
}

// GetModelBoundingBox - Compute model bounding box limits (considers all meshes)
func GetModelBoundingBox(model Model) BoundingBox {
	var ret BoundingBox
	getModelBoundingBox.Call(&ret, &model)
	return ret
}

// DrawModel - Draw a model (with texture if set)
func DrawModel(model Model, position Vector3, scale float32, tint color.RGBA) {
	drawModel.Call(nil, &model, &position, &scale, &tint)
}

// DrawModelEx - Draw a model with extended parameters
func DrawModelEx(model Model, position Vector3, rotationAxis Vector3, rotationAngle float32, scale Vector3, tint color.RGBA) {
	drawModelEx.Call(nil, &model, &position, &rotationAxis, &rotationAngle, &scale, &tint)
}

// DrawModelWires - Draw a model wires (with texture if set)
func DrawModelWires(model Model, position Vector3, scale float32, tint color.RGBA) {
	drawModelWires.Call(nil, &model, &position, &scale, &tint)
}

// DrawModelWiresEx - Draw a model wires (with texture if set) with extended parameters
func DrawModelWiresEx(model Model, position Vector3, rotationAxis Vector3, rotationAngle float32, scale Vector3, tint color.RGBA) {
	drawModelWiresEx.Call(nil, &model, &position, &rotationAxis, &rotationAngle, &scale, &tint)
}

// DrawBoundingBox - Draw bounding box (wires)
func DrawBoundingBox(box BoundingBox, col color.RGBA) {
	drawBoundingBox.Call(nil, &box, &col)
}

// DrawBillboard - Draw a billboard texture
func DrawBillboard(camera Camera, texture Texture2D, position Vector3, scale float32, tint color.RGBA) {
	drawBillboard.Call(nil, &camera, &texture, &position, &scale, &tint)
}

// DrawBillboardRec - Draw a billboard texture defined by source
func DrawBillboardRec(camera Camera, texture Texture2D, source Rectangle, position Vector3, size Vector2, tint color.RGBA) {
	drawBillboardRec.Call(nil, &camera, &texture, &source, &position, &size, &tint)
}

// DrawBillboardPro - Draw a billboard texture defined by source and rotation
func DrawBillboardPro(camera Camera, texture Texture2D, source Rectangle, position Vector3, up Vector3, size Vector2, origin Vector2, rotation float32, tint color.RGBA) {
	drawBillboardPro.Call(nil, &camera, &texture, &source, &position, &up, &size, &origin, &rotation, &tint)
}

// UploadMesh - Upload mesh vertex data in GPU and provide VAO/VBO ids
func UploadMesh(mesh *Mesh, dynamic bool) {
	uploadMesh.Call(nil, &mesh, &dynamic)
}

// UpdateMeshBuffer - Update mesh vertex data in GPU for a specific buffer index
func UpdateMeshBuffer(mesh Mesh, index int, data []byte, offset int) {
	dataSize := int32(len(data))
	dataPtr := unsafe.SliceData(data)
	o := int32(offset)
	i := int32(index)
	updateMeshBuffer.Call(nil, &mesh, &i, &dataPtr, &dataSize, &o)
}

// UnloadMesh - Unload mesh data from CPU and GPU
func UnloadMesh(mesh *Mesh) {
	unloadMesh.Call(nil, &mesh)
}

// DrawMesh - Draw a 3d mesh with material and transform
func DrawMesh(mesh Mesh, material Material, transform Matrix) {
	drawMesh.Call(nil, &mesh, &material, &transform)
}

// DrawMeshInstanced - Draw multiple mesh instances with material and different transforms
func DrawMeshInstanced(mesh Mesh, material Material, transforms []Matrix, instances int) {
	i := int32(instances)
	transformsPtr := unsafe.SliceData(transforms)
	drawMeshInstanced.Call(nil, &mesh, &material, &transformsPtr, &i)
}

// GetMeshBoundingBox - Compute mesh bounding box limits
func GetMeshBoundingBox(mesh Mesh) BoundingBox {
	var ret BoundingBox
	getMeshBoundingBox.Call(&ret, &mesh)
	return ret
}

// GenMeshTangents - Compute mesh tangents
func GenMeshTangents(mesh *Mesh) {
	genMeshTangents.Call(nil, &mesh)
}

// ExportMesh - Export mesh data to file, returns true on success
func ExportMesh(mesh Mesh, fileName string) bool {
	fileNamePtr := convert.ToBytePtr(fileName)
	var ret ffi.Arg
	exportMesh.Call(&ret, &mesh, &fileNamePtr)
	return ret.Bool()
}

// GenMeshPoly - Generate polygonal mesh
func GenMeshPoly(sides int, radius float32) Mesh {
	var ret Mesh
	s := int32(sides)
	genMeshPoly.Call(&ret, &s, &radius)
	return ret
}

// GenMeshPlane - Generate plane mesh (with subdivisions)
func GenMeshPlane(width, length float32, resX, resZ int) Mesh {
	var ret Mesh
	x, z := int32(resX), int32(resZ)
	genMeshPlane.Call(&ret, &width, &length, &x, &z)
	return ret
}

// GenMeshCube - Generate cuboid mesh
func GenMeshCube(width, height, length float32) Mesh {
	var ret Mesh
	genMeshCube.Call(&ret, &width, &height, &length)
	return ret
}

// GenMeshSphere - Generate sphere mesh (standard sphere)
func GenMeshSphere(radius float32, rings, slices int) Mesh {
	var ret Mesh
	r, s := int32(rings), int32(slices)
	genMeshSphere.Call(&ret, &radius, &r, &s)
	return ret
}

// GenMeshHemiSphere - Generate half-sphere mesh (no bottom cap)
func GenMeshHemiSphere(radius float32, rings, slices int) Mesh {
	var ret Mesh
	r, s := int32(rings), int32(slices)
	genMeshHemiSphere.Call(&ret, &radius, &r, &s)
	return ret
}

// GenMeshCylinder - Generate cylinder mesh
func GenMeshCylinder(radius, height float32, slices int) Mesh {
	var ret Mesh
	s := int32(slices)
	genMeshCylinder.Call(&ret, &radius, &height, &s)
	return ret
}

// GenMeshCone - Generate cone/pyramid mesh
func GenMeshCone(radius, height float32, slices int) Mesh {
	var ret Mesh
	s := int32(slices)
	genMeshCone.Call(&ret, &radius, &height, &s)
	return ret
}

// GenMeshTorus - Generate torus mesh
func GenMeshTorus(radius, size float32, radSeg, sides int) Mesh {
	var ret Mesh
	r, s := int32(radSeg), int32(sides)
	genMeshTorus.Call(&ret, &radius, &size, &r, &s)
	return ret
}

// GenMeshKnot - Generate trefoil knot mesh
func GenMeshKnot(radius, size float32, radSeg, sides int) Mesh {
	var ret Mesh
	r, s := int32(radSeg), int32(sides)
	genMeshKnot.Call(&ret, &radius, &size, &r, &s)
	return ret
}

// GenMeshHeightmap - Generate heightmap mesh from image data
func GenMeshHeightmap(heightmap Image, size Vector3) Mesh {
	var ret Mesh
	genMeshHeightmap.Call(&ret, &heightmap, &size)
	return ret
}

// GenMeshCubicmap - Generate cubes-based map mesh from image data
func GenMeshCubicmap(cubicmap Image, cubeSize Vector3) Mesh {
	var ret Mesh
	genMeshCubicmap.Call(&ret, &cubicmap, &cubeSize)
	return ret
}

// LoadMaterials - Load materials from model file
func LoadMaterials(fileName string) []Material {
	var materialCount int32
	materialCountPtr := &materialCount
	var ret *Material
	fileNamePtr := convert.ToBytePtr(fileName)
	loadMaterials.Call(&ret, &fileNamePtr, &materialCountPtr)
	return unsafe.Slice(ret, materialCount)
}

// LoadMaterialDefault - Load default material (Supports: DIFFUSE, SPECULAR, NORMAL maps)
func LoadMaterialDefault() Material {
	var ret Material
	loadMaterialDefault.Call(&ret)
	return ret
}

// IsMaterialValid - Check if a material is valid (shader assigned, map textures loaded in GPU)
func IsMaterialValid(material Material) bool {
	var ret ffi.Arg
	isMaterialValid.Call(&ret, &material)
	return ret.Bool()
}

// UnloadMaterial - Unload material from GPU memory (VRAM)
func UnloadMaterial(material Material) {
	unloadMaterial.Call(nil, &material)
}

// SetMaterialTexture - Set texture for a material map type (MATERIAL_MAP_DIFFUSE, MATERIAL_MAP_SPECULAR...)
func SetMaterialTexture(material *Material, mapType int32, texture Texture2D) {
	setMaterialTexture.Call(nil, &material, &mapType, &texture)
}

// SetModelMeshMaterial - Set material for a mesh
func SetModelMeshMaterial(model *Model, meshId int32, materialId int32) {
	setModelMeshMaterial.Call(nil, &model, &meshId, &materialId)
}

// LoadModelAnimations - Load model animations from file
func LoadModelAnimations(fileName string) []ModelAnimation {
	fileNamePtr := convert.ToBytePtr(fileName)
	var animCount int32
	animCountPtr := &animCount
	var ret *ModelAnimation
	loadModelAnimations.Call(&ret, &fileNamePtr, &animCountPtr)
	return unsafe.Slice(ret, animCount)
}

// UpdateModelAnimation - Update model animation pose (vertex buffers and bone matrices)
func UpdateModelAnimation(model Model, anim ModelAnimation, frame float32) {
	updateModelAnimation.Call(nil, &model, &anim, &frame)
}

// UpdateModelAnimationEx - Update model animation pose, blending two animations
func UpdateModelAnimationEx(model Model, animA ModelAnimation, frameA float32, animB ModelAnimation, frameB, blend float32) {
	updateModelAnimationEx.Call(nil, &model, &animA, &frameA, &animB, &frameB, &blend)
}

// UnloadModelAnimations - Unload animation array data
func UnloadModelAnimations(animations []ModelAnimation) {
	animCount := int32(len(animations))
	animationsPtr := unsafe.SliceData(animations)
	unloadModelAnimations.Call(nil, &animationsPtr, &animCount)
}

// IsModelAnimationValid - Check model animation skeleton match
func IsModelAnimationValid(model Model, anim ModelAnimation) bool {
	var ret ffi.Arg
	isModelAnimationValid.Call(&ret, &model, &anim)
	return ret.Bool()
}

// CheckCollisionSpheres - Check collision between two spheres
func CheckCollisionSpheres(center1 Vector3, radius1 float32, center2 Vector3, radius2 float32) bool {
	var ret ffi.Arg
	checkCollisionSpheres.Call(&ret, &center1, &radius1, &center2, &radius2)
	return ret.Bool()
}

// CheckCollisionBoxes - Check collision between two bounding boxes
func CheckCollisionBoxes(box1 BoundingBox, box2 BoundingBox) bool {
	var ret ffi.Arg
	checkCollisionBoxes.Call(&ret, &box1, &box2)
	return ret.Bool()
}

// CheckCollisionBoxSphere - Check collision between box and sphere
func CheckCollisionBoxSphere(box BoundingBox, center Vector3, radius float32) bool {
	var ret ffi.Arg
	checkCollisionBoxSphere.Call(&ret, &box, &center, &radius)
	return ret.Bool()
}

// GetRayCollisionSphere - Get collision info between ray and sphere
func GetRayCollisionSphere(ray Ray, center Vector3, radius float32) RayCollision {
	var ret RayCollision
	getRayCollisionSphere.Call(&ret, &ray, &center, &radius)
	return ret
}

// GetRayCollisionBox - Get collision info between ray and box
func GetRayCollisionBox(ray Ray, box BoundingBox) RayCollision {
	var ret RayCollision
	getRayCollisionBox.Call(&ret, &ray, &box)
	return ret
}

// GetRayCollisionMesh - Get collision info between ray and mesh
func GetRayCollisionMesh(ray Ray, mesh Mesh, transform Matrix) RayCollision {
	var ret RayCollision
	getRayCollisionMesh.Call(&ret, &ray, &mesh, &transform)
	return ret
}

// GetRayCollisionTriangle - Get collision info between ray and triangle
func GetRayCollisionTriangle(ray Ray, p1, p2, p3 Vector3) RayCollision {
	var ret RayCollision
	getRayCollisionTriangle.Call(&ret, &ray, &p1, &p2, &p3)
	return ret
}

// GetRayCollisionQuad - Get collision info between ray and quad
func GetRayCollisionQuad(ray Ray, p1, p2, p3, p4 Vector3) RayCollision {
	var ret RayCollision
	getRayCollisionQuad.Call(&ret, &ray, &p1, &p2, &p3, &p4)
	return ret
}

// InitAudioDevice - Initialize audio device and context
func InitAudioDevice() {
	initAudioDevice.Call(nil)
}

// CloseAudioDevice - Close the audio device and context
func CloseAudioDevice() {
	closeAudioDevice.Call(nil)
}

// IsAudioDeviceReady - Check if audio device has been initialized successfully
func IsAudioDeviceReady() bool {
	var ret ffi.Arg
	isAudioDeviceReady.Call(&ret)
	return ret.Bool()
}

// SetMasterVolume - Set master volume (listener)
func SetMasterVolume(volume float32) {
	setMasterVolume.Call(nil, &volume)
}

// GetMasterVolume - Get master volume (listener)
func GetMasterVolume() float32 {
	var ret float32
	getMasterVolume.Call(&ret)
	return ret
}

// LoadWave - Load wave data from file
func LoadWave(fileName string) Wave {
	fileNamePtr := convert.ToBytePtr(fileName)
	var ret Wave
	loadWave.Call(&ret, &fileNamePtr)
	return ret
}

// LoadWaveFromMemory - Load wave from memory buffer, fileType refers to extension: i.e. '.wav'
func LoadWaveFromMemory(fileType string, fileData []byte, dataSize int32) Wave {
	var ret Wave
	fileTypePtr := convert.ToBytePtr(fileType)
	fileDataPtr := unsafe.SliceData(fileData)
	loadWaveFromMemory.Call(&ret, &fileTypePtr, &fileDataPtr, &dataSize)
	return ret
}

// IsWaveValid - Checks if wave data is valid (data loaded and parameters)
func IsWaveValid(wave Wave) bool {
	var ret ffi.Arg
	isWaveValid.Call(&ret, &wave)
	return ret.Bool()
}

// LoadSound - Load sound from file
func LoadSound(fileName string) Sound {
	fileNamePtr := convert.ToBytePtr(fileName)
	var ret Sound
	loadSound.Call(&ret, &fileNamePtr)
	return ret
}

// LoadSoundFromWave - Load sound from wave data
func LoadSoundFromWave(wave Wave) Sound {
	var ret Sound
	loadSoundFromWave.Call(&ret, &wave)
	return ret
}

// LoadSoundAlias - Create a new sound that shares the same sample data as the source sound, does not own the sound data
func LoadSoundAlias(source Sound) Sound {
	var ret Sound
	loadSoundAlias.Call(&ret, &source)
	return ret
}

// IsSoundValid - Checks if a sound is valid (data loaded and buffers initialized)
func IsSoundValid(sound Sound) bool {
	var ret ffi.Arg
	isSoundValid.Call(&ret, &sound)
	return ret.Bool()
}

// UpdateSound - Update sound buffer with new data (default data format: 32 bit float, stereo)
func UpdateSound(sound Sound, data []byte, sampleCount int32) {
	dataPtr := unsafe.SliceData(data)
	updateSound.Call(nil, &sound, &dataPtr, &sampleCount)
}

// UnloadWave - Unload wave data
func UnloadWave(wave Wave) {
	unloadWave.Call(nil, &wave)
}

// UnloadSound - Unload sound
func UnloadSound(sound Sound) {
	unloadSound.Call(nil, &sound)
}

// UnloadSoundAlias - Unload a sound alias (does not deallocate sample data)
func UnloadSoundAlias(alias Sound) {
	unloadSoundAlias.Call(nil, &alias)
}

// ExportWave - Export wave data to file, returns true on success
func ExportWave(wave Wave, fileName string) bool {
	fileNamePtr := convert.ToBytePtr(fileName)
	var ret ffi.Arg
	exportWave.Call(&ret, &wave, &fileNamePtr)
	return ret.Bool()
}

// PlaySound - Play a sound
func PlaySound(sound Sound) {
	playSound.Call(nil, &sound)
}

// StopSound - Stop playing a sound
func StopSound(sound Sound) {
	stopSound.Call(nil, &sound)
}

// PauseSound - Pause a sound
func PauseSound(sound Sound) {
	pauseSound.Call(nil, &sound)
}

// ResumeSound - Resume a paused sound
func ResumeSound(sound Sound) {
	resumeSound.Call(nil, &sound)
}

// IsSoundPlaying - Check if a sound is currently playing
func IsSoundPlaying(sound Sound) bool {
	var ret ffi.Arg
	isSoundPlaying.Call(&ret, &sound)
	return ret.Bool()
}

// SetSoundVolume - Set volume for a sound (1.0 is max level)
func SetSoundVolume(sound Sound, volume float32) {
	setSoundVolume.Call(nil, &sound, &volume)
}

// SetSoundPitch - Set pitch for a sound (1.0 is base level)
func SetSoundPitch(sound Sound, pitch float32) {
	setSoundPitch.Call(nil, &sound, &pitch)
}

// SetSoundPan - Set pan for a sound (-1.0 left, 0.0 center, 1.0 right)
func SetSoundPan(sound Sound, pan float32) {
	setSoundPan.Call(nil, &sound, &pan)
}

// WaveCopy - Copy a wave to a new wave
func WaveCopy(wave Wave) Wave {
	var ret Wave
	waveCopy.Call(&ret, &wave)
	return ret
}

// WaveCrop - Crop a wave to defined frames range
func WaveCrop(wave *Wave, initFrame int32, finalFrame int32) {
	waveCrop.Call(nil, &wave, &initFrame, &finalFrame)
}

// WaveFormat - Convert wave data to desired format
func WaveFormat(wave *Wave, sampleRate int32, sampleSize int32, channels int32) {
	waveFormat.Call(nil, &wave, &sampleRate, &sampleRate, &channels)
}

// LoadWaveSamples - Load samples data from wave as a 32bit float data array
func LoadWaveSamples(wave Wave) []float32 {
	var ret *float32
	loadWaveSamples.Call(&ret, &wave)
	return unsafe.Slice(ret, wave.FrameCount*wave.Channels)
}

// UnloadWaveSamples - Unload samples data loaded with LoadWaveSamples()
func UnloadWaveSamples(samples []float32) {
	samplesPtr := unsafe.SliceData(samples)
	unloadWaveSamples.Call(nil, &samplesPtr)
}

// LoadMusicStream - Load music stream from file
func LoadMusicStream(fileName string) Music {
	fileNamePtr := convert.ToBytePtr(fileName)
	var ret Music
	loadMusicStream.Call(&ret, &fileNamePtr)
	return ret
}

// LoadMusicStreamFromMemory - Load music stream from data
func LoadMusicStreamFromMemory(fileType string, data []byte, dataSize int32) Music {
	var ret Music
	fileTypePtr := convert.ToBytePtr(fileType)
	dataPtr := unsafe.SliceData(data)
	loadMusicStreamFromMemory.Call(&ret, &fileTypePtr, &dataPtr, &dataSize)
	return ret
}

// IsMusicValid - Checks if a music stream is valid (context and buffers initialized)
func IsMusicValid(music Music) bool {
	var ret ffi.Arg
	isMusicValid.Call(&ret, &music)
	return ret.Bool()
}

// UnloadMusicStream - Unload music stream
func UnloadMusicStream(music Music) {
	unloadMusicStream.Call(nil, &music)
}

// PlayMusicStream - Start music playing
func PlayMusicStream(music Music) {
	playMusicStream.Call(nil, &music)
}

// IsMusicStreamPlaying - Check if music is playing
func IsMusicStreamPlaying(music Music) bool {
	var ret ffi.Arg
	isMusicStreamPlaying.Call(&ret, &music)
	return ret.Bool()
}

// UpdateMusicStream - Updates buffers for music streaming
func UpdateMusicStream(music Music) {
	updateMusicStream.Call(nil, &music)
}

// StopMusicStream - Stop music playing
func StopMusicStream(music Music) {
	stopMusicStream.Call(nil, &music)
}

// PauseMusicStream - Pause music playing
func PauseMusicStream(music Music) {
	pauseMusicStream.Call(nil, &music)
}

// ResumeMusicStream - Resume playing paused music
func ResumeMusicStream(music Music) {
	resumeMusicStream.Call(nil, &music)
}

// SeekMusicStream - Seek music to a position (in seconds)
func SeekMusicStream(music Music, position float32) {
	seekMusicStream.Call(nil, &music, &position)
}

// SetMusicVolume - Set volume for music (1.0 is max level)
func SetMusicVolume(music Music, volume float32) {
	setMusicVolume.Call(nil, &music, &volume)
}

// SetMusicPitch - Set pitch for a music (1.0 is base level)
func SetMusicPitch(music Music, pitch float32) {
	setMusicPitch.Call(nil, &music, &pitch)
}

// SetMusicPan - Set pan for a music (-1.0 left, 0.0 center, 1.0 right)
func SetMusicPan(music Music, pan float32) {
	setMusicPan.Call(nil, &music, &pan)
}

// GetMusicTimeLength - Get music time length (in seconds)
func GetMusicTimeLength(music Music) float32 {
	var ret float32
	getMusicTimeLength.Call(&ret, &music)
	return ret
}

// GetMusicTimePlayed - Get current music time played (in seconds)
func GetMusicTimePlayed(music Music) float32 {
	var ret float32
	getMusicTimePlayed.Call(&ret, &music)
	return ret
}

// LoadAudioStream - Load audio stream (to stream raw audio pcm data)
func LoadAudioStream(sampleRate uint32, sampleSize uint32, channels uint32) AudioStream {
	var ret AudioStream
	loadAudioStream.Call(&ret, &sampleRate, &sampleSize, &channels)
	return ret
}

// IsAudioStreamValid - Checks if an audio stream is valid (buffers initialized)
func IsAudioStreamValid(stream AudioStream) bool {
	var ret ffi.Arg
	isAudioStreamValid.Call(&ret, &stream)
	return ret.Bool()
}

// UnloadAudioStream - Unload audio stream and free memory
func UnloadAudioStream(stream AudioStream) {
	unloadAudioStream.Call(nil, &stream)
}

// UpdateAudioStream - Update audio stream buffers with data
func UpdateAudioStream[T float32 | int16](stream AudioStream, data []T) {
	var dataPtr unsafe.Pointer
	var frameCount int32
	switch d := any(data).(type) {
	case []float32:
		dataPtr = unsafe.Pointer(&d[0])
		frameCount = int32(len(d))
	case []int16:
		dataPtr = unsafe.Pointer(&d[0])
		frameCount = int32(len(d))
	}
	updateAudioStream.Call(nil, &stream, &dataPtr, &frameCount)
}

// IsAudioStreamProcessed - Check if any audio stream buffers requires refill
func IsAudioStreamProcessed(stream AudioStream) bool {
	var ret ffi.Arg
	isAudioStreamProcessed.Call(&ret, &stream)
	return ret.Bool()
}

// PlayAudioStream - Play audio stream
func PlayAudioStream(stream AudioStream) {
	playAudioStream.Call(nil, &stream)
}

// PauseAudioStream - Pause audio stream
func PauseAudioStream(stream AudioStream) {
	pauseAudioStream.Call(nil, &stream)
}

// ResumeAudioStream - Resume audio stream
func ResumeAudioStream(stream AudioStream) {
	resumeAudioStream.Call(nil, &stream)
}

// IsAudioStreamPlaying - Check if audio stream is playing
func IsAudioStreamPlaying(stream AudioStream) bool {
	var ret ffi.Arg
	isAudioStreamPlaying.Call(&ret, &stream)
	return ret.Bool()
}

// StopAudioStream - Stop audio stream
func StopAudioStream(stream AudioStream) {
	stopAudioStream.Call(nil, &stream)
}

// SetAudioStreamVolume - Set volume for audio stream (1.0 is max level)
func SetAudioStreamVolume(stream AudioStream, volume float32) {
	setAudioStreamVolume.Call(nil, &stream, &volume)
}

// SetAudioStreamPitch - Set pitch for audio stream (1.0 is base level)
func SetAudioStreamPitch(stream AudioStream, pitch float32) {
	setAudioStreamPitch.Call(nil, &stream, &pitch)
}

// SetAudioStreamPan - Set pan for audio stream (-1.0 to 1.0 range, 0.0 is centered)
func SetAudioStreamPan(stream AudioStream, pan float32) {
	setAudioStreamPan.Call(nil, &stream, &pan)
}

// SetAudioStreamBufferSizeDefault - Default size for new audio streams
func SetAudioStreamBufferSizeDefault(size int32) {
	setAudioStreamBufferSizeDefault.Call(nil, &size)
}

// SetAudioStreamCallback - Audio thread callback to request new data
func SetAudioStreamCallback(stream AudioStream, callback AudioCallback) {
	channels := int(stream.Channels)
	if channels < 1 {
		channels = 1
	}
	fn := purego.NewCallback(func(bufferData unsafe.Pointer, frames int32) uintptr {
		// Buffer holds frames*channels samples, interleaved per stream channels.
		callback(unsafe.Slice((*float32)(bufferData), int(frames)*channels), int(frames))
		return 0
	})
	setAudioStreamCallback.Call(nil, &stream, &fn)
}

// AttachAudioStreamProcessor - Attach audio stream processor to stream, receives frames x 2 samples as 'float' (stereo)
func AttachAudioStreamProcessor(stream AudioStream, processor AudioCallback) {
	fn := purego.NewCallback(func(bufferData unsafe.Pointer, frames int32) uintptr {
		processor(unsafe.Slice((*float32)(bufferData), frames), int(frames))
		return 0
	})
	ptr := uintptr(reflect.ValueOf(processor).UnsafePointer())
	audioCallbacks[ptr] = fn
	attachAudioStreamProcessor.Call(nil, &stream, &fn)
}

// DetachAudioStreamProcessor - Detach audio stream processor from stream
func DetachAudioStreamProcessor(stream AudioStream, processor AudioCallback) {
	ptr := uintptr(reflect.ValueOf(processor).UnsafePointer())
	fn := audioCallbacks[ptr]
	detachAudioStreamProcessor.Call(nil, &stream, &fn)
}

// AttachAudioMixedProcessor - Attach audio stream processor to the entire audio pipeline, receives frames x 2 samples as 'float' (stereo)
func AttachAudioMixedProcessor(processor AudioCallback) {
	fn := purego.NewCallback(func(bufferData unsafe.Pointer, frames int32) uintptr {
		processor(unsafe.Slice((*float32)(bufferData), frames), int(frames))
		return 0
	})
	ptr := uintptr(reflect.ValueOf(processor).UnsafePointer())
	audioCallbacks[ptr] = fn
	attachAudioMixedProcessor.Call(nil, &fn)
}

// DetachAudioMixedProcessor - Detach audio stream processor from the entire audio pipeline
func DetachAudioMixedProcessor(processor AudioCallback) {
	ptr := uintptr(reflect.ValueOf(processor).UnsafePointer())
	fn := audioCallbacks[ptr]
	detachAudioMixedProcessor.Call(nil, &fn)
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
