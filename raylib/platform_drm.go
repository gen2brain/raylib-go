//go:build linux && drm && !android && !nocgo
// +build linux,drm,!android,!nocgo

package rl

/*
#include "raylib.h"
#include <stdlib.h>
*/
import "C"

// SetMain - Sets callback function
func SetMain(func()) {
	return
}

// ShowCursor - Shows cursor
func ShowCursor() {
	C.ShowCursor()
}

// HideCursor - Hides cursor
func HideCursor() {
	C.HideCursor()
}

// IsCursorHidden - Check if cursor is not visible
func IsCursorHidden() bool {
	ret := C.IsCursorHidden()
	v := bool(ret)
	return v
}

// IsCursorOnScreen - Check if cursor is on the screen
func IsCursorOnScreen() bool {
	ret := C.IsCursorOnScreen()
	v := bool(ret)
	return v
}

// EnableCursor - Enables cursor (unlock cursor)
func EnableCursor() {
	C.EnableCursor()
}

// DisableCursor - Disables cursor (lock cursor)
func DisableCursor() {
	C.DisableCursor()
}

// IsFileDropped - Check if a file has been dropped into window
func IsFileDropped() bool {
	return false
}

// LoadDroppedFiles - Load dropped filepaths
func LoadDroppedFiles() (files []string) {
	return
}

// UnloadDroppedFiles - Unload dropped filepaths
func UnloadDroppedFiles() {
	return
}
