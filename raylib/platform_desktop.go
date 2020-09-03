// +build !android,!arm

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
	ret := C.IsFileDropped()
	v := bool(ret)
	return v
}

// GetDroppedFiles - Retrieve dropped files into window
func GetDroppedFiles(count *int32) []string {
	ccount := (*C.int)(unsafe.Pointer(count))
	ret := C.GetDroppedFiles(ccount)

	tmpslice := (*[1 << 24]*C.char)(unsafe.Pointer(ret))[:*count:*count]
	gostrings := make([]string, *count)
	for i, s := range tmpslice {
		gostrings[i] = C.GoString(s)
	}

	return gostrings
}

// ClearDroppedFiles - Clear dropped files paths buffer
func ClearDroppedFiles() {
	C.ClearDroppedFiles()
}

// OpenAsset - Open asset
func OpenAsset(name string) (Asset, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	return f, nil
}

// IsWindowMaximized - Check if window has been maximized (only PLATFORM_DESKTOP)
func IsWindowMaximized() bool {
	ret := C.IsWindowMaximized()
	v := bool(ret)
	return v
}

// IsWindowFocused - Check if window has been focused
func IsWindowFocused() bool {
	ret := C.IsWindowFocused()
	v := bool(ret)
	return v
}

// UndecorateWindow - Undecorate the window (only PLATFORM_DESKTOP)
func UndecorateWindow() {
	C.UndecorateWindow()
}

// MaximizeWindow - Maximize the window, if resizable (only PLATFORM_DESKTOP)
func MaximizeWindow() {
	C.MaximizeWindow()
}

// RestoreWindow - Restore the window, if resizable (only PLATFORM_DESKTOP)
func RestoreWindow() {
	C.RestoreWindow()
}

// DecorateWindow - Decorate the window (only PLATFORM_DESKTOP)
func DecorateWindow() {
	C.DecorateWindow()
}

// GetMonitorRefreshRate - Get primary monitor refresh rate
func GetMonitorRefreshRate(monitor int) int {
	cmonitor := (C.int)(monitor)
	ret := C.GetMonitorRefreshRate(cmonitor)
	v := (int)(ret)
	return v
}

// GetWindowScaleDPI - Get window scale DPI factor
func GetWindowScaleDPI() Vector2 {
	ret := C.GetWindowScaleDPI()
	v := newVector2FromPointer(unsafe.Pointer(&ret))
	return v
}

// IsCursorOnScreen - Check if cursor is on the current screen.
func IsCursorOnScreen() bool {
	ret := C.IsCursorOnScreen()
	v := bool(ret)
	return v
}
