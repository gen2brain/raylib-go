// +build !android

package raylib

/*
#include "raylib.h"
*/
import "C"
import "unsafe"

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

// NewCameraFromPointer - Returns new Camera from pointer
func NewCameraFromPointer(ptr unsafe.Pointer) Camera {
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

// NewCamera2DFromPointer - Returns new Camera2D from pointer
func NewCamera2DFromPointer(ptr unsafe.Pointer) Camera2D {
	return *(*Camera2D)(ptr)
}

// CameraMode type
type CameraMode int32

// Camera system modes
const (
	CameraCustom      CameraMode = C.CAMERA_CUSTOM
	CameraFree        CameraMode = C.CAMERA_FREE
	CameraOrbital     CameraMode = C.CAMERA_ORBITAL
	CameraFirstPerson CameraMode = C.CAMERA_FIRST_PERSON
	CameraThirdPerson CameraMode = C.CAMERA_THIRD_PERSON
)

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
