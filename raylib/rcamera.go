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
