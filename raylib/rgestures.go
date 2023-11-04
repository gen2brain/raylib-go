package rl

/*
#include "raylib.h"
*/
import "C"
import "unsafe"

// SetGesturesEnabled - Enable a set of gestures using flags
func SetGesturesEnabled(gestureFlags uint32) {
	cgestureFlags := (C.uint)(gestureFlags)
	C.SetGesturesEnabled(cgestureFlags)
}

// IsGestureDetected - Check if a gesture have been detected
func IsGestureDetected(gesture Gestures) bool {
	cgesture := (C.uint)(gesture)
	ret := C.IsGestureDetected(cgesture)
	v := bool(ret)
	return v
}

// GetGestureDetected - Get latest detected gesture
func GetGestureDetected() Gestures {
	ret := C.GetGestureDetected()
	v := (Gestures)(ret)
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
