//go:build linux && drm && !android
// +build linux,drm,!android

package rl

/*
#include "raylib.h"
#include <stdlib.h>
*/
import "C"

import (
	"os"
)

// SetCallbackFunc - Sets callback function
func SetCallbackFunc(func()) {
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

// IsCursorHidden - Returns true if cursor is not visible
func IsCursorHidden() bool {
	ret := C.IsCursorHidden()
	v := bool(ret)
	return v
}

// IsCursorOnScreen - Check if cursor is on the current screen.
func IsCursorOnScreen() bool {
	ret := C.IsCursorOnScreen()
	v := bool(ret)
	return v
}

// EnableCursor - Enables cursor
func EnableCursor() {
	C.EnableCursor()
}

// DisableCursor - Disables cursor
func DisableCursor() {
	C.DisableCursor()
}

// IsFileDropped - Check if a file have been dropped into window
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

// OpenAsset - Open asset
func OpenAsset(name string) (Asset, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	return f, nil
}
