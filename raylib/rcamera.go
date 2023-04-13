//go:build !android
// +build !android

package rl

/*
#include "raylib.h"
*/
import "C"

// UpdateCamera - Update camera position for selected mode
func UpdateCamera(camera *Camera, mode CameraMode) {
	ccamera := camera.cptr()
	C.UpdateCamera(ccamera, C.int(mode))
}

// UpdateCameraPro - Update camera movement/rotation
func UpdateCameraPro(camera *Camera, movement Vector3, rotation Vector3, zoom float32) {
	ccamera := camera.cptr()
	C.UpdateCameraPro(ccamera, *movement.cptr(), *rotation.cptr(), C.float(zoom))
}
