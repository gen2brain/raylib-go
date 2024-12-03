//go:build linux && drm && !rgfw && !sdl && !sdl3 && !android
// +build linux,drm,!rgfw,!sdl,!sdl3,!android

package rl

/*
#include "raylib.h"
#include <stdlib.h>
*/
import "C"

import (
	"os"
	"unsafe"
)

// InitWindow - Initialize Window and OpenGL Graphics
func InitWindow(width int32, height int32, title string) {
	cwidth := (C.int)(width)
	cheight := (C.int)(height)

	ctitle := C.CString(title)
	defer C.free(unsafe.Pointer(ctitle))

	C.InitWindow(cwidth, cheight, ctitle)
}

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
