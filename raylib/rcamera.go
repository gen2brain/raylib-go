//go:build !android
// +build !android

package rl

/*
#include "raylib.h"
*/
import "C"

// SetCameraMode - Set camera mode (multiple camera modes available)
func SetCameraMode(camera Camera, mode CameraMode) {
	ccamera := camera.cptr()
	cmode := (C.int)(mode)
	C.SetCameraMode(*ccamera, cmode)
}

// UpdateCamera - Update camera position for selected mode
func UpdateCamera(camera *Camera) {
	ccamera := camera.cptr()
	C.UpdateCamera(ccamera)
}

// SetCameraPanControl - Set camera pan key to combine with mouse movement (free camera)
func SetCameraPanControl(panKey int32) {
	cpanKey := (C.int)(panKey)
	C.SetCameraPanControl(cpanKey)
}

// SetCameraAltControl - Set camera alt key to combine with mouse movement (free camera)
func SetCameraAltControl(altKey int32) {
	caltKey := (C.int)(altKey)
	C.SetCameraAltControl(caltKey)
}

// SetCameraSmoothZoomControl - Set camera smooth zoom key to combine with mouse (free camera)
func SetCameraSmoothZoomControl(szKey int32) {
	cszKey := (C.int)(szKey)
	C.SetCameraSmoothZoomControl(cszKey)
}

// SetCameraMoveControls - Set camera move controls (1st person and 3rd person cameras)
func SetCameraMoveControls(frontKey int32, backKey int32, rightKey int32, leftKey int32, upKey int32, downKey int32) {
	cfrontKey := (C.int)(frontKey)
	cbackKey := (C.int)(backKey)
	crightKey := (C.int)(rightKey)
	cleftKey := (C.int)(leftKey)
	cupKey := (C.int)(upKey)
	cdownKey := (C.int)(downKey)
	C.SetCameraMoveControls(cfrontKey, cbackKey, crightKey, cleftKey, cupKey, cdownKey)
}
