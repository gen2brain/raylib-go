//go:build !drm && !android
// +build !drm,!android

package rl

/*
#include "raylib.h"
#include <stdlib.h>
*/
import "C"

import (
	"unsafe"
)

// SetMain - Sets callback function
func SetMain(func()) {
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
	ret := C.IsFileDropped()
	v := bool(ret)

	return v
}

// LoadDroppedFiles - Load dropped filepaths
func LoadDroppedFiles() []string {
	ret := C.LoadDroppedFiles()
	defer C.UnloadDroppedFiles(ret)

	tmpslice := (*[1 << 24]*C.char)(unsafe.Pointer(ret.paths))[:ret.count:ret.count]
	gostrings := make([]string, ret.count)
	for i, s := range tmpslice {
		gostrings[i] = C.GoString(s)
	}

	return gostrings
}

// UnloadDroppedFiles - Unload dropped filepaths
func UnloadDroppedFiles() {
}
