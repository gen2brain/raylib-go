package raylib

/*
#include "raylib.h"
*/
import "C"
import "unsafe"

// Gestures type
type Gestures int32

// Gestures types
// NOTE: It could be used as flags to enable only some gestures
const (
	GestureNone       Gestures = 0
	GestureTap        Gestures = 1
	GestureDoubletap  Gestures = 2
	GestureHold       Gestures = 4
	GestureDrag       Gestures = 8
	GestureSwipeRight Gestures = 16
	GestureSwipeLeft  Gestures = 32
	GestureSwipeUp    Gestures = 64
	GestureSwipeDown  Gestures = 128
	GesturePinchIn    Gestures = 256
	GesturePinchOut   Gestures = 512
)

// SetGesturesEnabled - Enable a set of gestures using flags
func SetGesturesEnabled(gestureFlags uint32) {
	cgestureFlags := (C.uint)(gestureFlags)
	C.SetGesturesEnabled(cgestureFlags)
}

// IsGestureDetected - Check if a gesture have been detected
func IsGestureDetected(gesture Gestures) bool {
	cgesture := (C.int)(gesture)
	ret := C.IsGestureDetected(cgesture)
	v := bool(int(ret) == 1)
	return v
}

// GetGestureDetected - Get latest detected gesture
func GetGestureDetected() Gestures {
	ret := C.GetGestureDetected()
	v := (Gestures)(ret)
	return v
}

// GetTouchPointsCount - Get touch points count
func GetTouchPointsCount() int32 {
	ret := C.GetTouchPointsCount()
	v := (int32)(ret)
	return v
}

// GetGestureHoldDuration - Get gesture hold time in milliseconds
func GetGestureHoldDuration() float32 {
	ret := C.GetGestureHoldDuration()
	v := (float32)(ret)
	return v
}

// GetGestureDragVector - Get gesture drag vector
func GetGestureDragVector() Vector2 {
	ret := C.GetGestureDragVector()
	v := newVector2FromPointer(unsafe.Pointer(&ret))
	return v
}

// GetGestureDragAngle - Get gesture drag angle
func GetGestureDragAngle() float32 {
	ret := C.GetGestureDragAngle()
	v := (float32)(ret)
	return v
}

// GetGesturePinchVector - Get gesture pinch delta
func GetGesturePinchVector() Vector2 {
	ret := C.GetGesturePinchVector()
	v := newVector2FromPointer(unsafe.Pointer(&ret))
	return v
}

// GetGesturePinchAngle - Get gesture pinch angle
func GetGesturePinchAngle() float32 {
	ret := C.GetGesturePinchAngle()
	v := (float32)(ret)
	return v
}
